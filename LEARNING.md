# Learning Log

Quick reference for the roadmap concepts I've learned.

## 1. CRUD

**Create, Read, Update, Delete**: the four operations on a stored resource. Each maps to an HTTP method; the method says _what to do_, the path says _to what_.

| Operation | Request                           | Success |
| --------- | --------------------------------- | ------- |
| Create    | `POST /boards`                    | 201     |
| Read      | `GET /boards`, `GET /boards/{id}` | 200     |
| Update    | `PUT /boards/{id}`                | 200     |
| Delete    | `DELETE /boards/{id}`             | 204     |

- **Idempotent** means repeating the request has the same effect as doing it once. `GET`, `PUT` and `DELETE` are idempotent; `POST` isn't.
- Not found is detected with `sql.ErrNoRows` (reads) or `RowsAffected() == 0` (updates and deletes).

## 2. HTTP Fundamentals

### Status codes

The first line of every response; the client checks it before reading the body.

- **2xx success:** 200 OK, 201 Created, 204 No Content (no body).
- **4xx client's fault:** 400 bad input, 404 not found, 409 conflicts with existing state.
- **5xx server's fault:** 500, the request might work later.
- Ask: whose fault is it, and could retrying help?
- `w.WriteHeader` works once only. If you never call it, the first write sends 200.

### JSON encode/decode

Converting between JSON text and Go structs.

- **Decode:** `json.NewDecoder(r.Body).Decode(&v)`, pass a pointer so it can fill `v`.
- **Encode:** `json.NewEncoder(w).Encode(v)`.
- **Struct tags** set a field's name in JSON. Go field names are `CamelCase`, while JSON keys here are `snake_case`, so the tag links the two:

  ```go
  CreatedAt time.Time `json:"created_at"`   // Go: CreatedAt  ↔  JSON: "created_at"
  ```

- **Only exported fields are encoded or decoded.** A field is exported if its name starts with a capital letter, which makes it visible outside its package. `encoding/json` is a different package from yours, so it can't see lowercase fields:

  ```go
  Title string   // exported: included in JSON
  notes string   // unexported: ignored
  ```
- A `nil` slice encodes as `null`; an empty slice as `[]`.
- Unknown fields are ignored; wrong types are a decode error (400).

### Validation

Decoding proves the JSON is readable; validation proves it makes sense.

- Rules are struct tags: `validate:"required,email,max=254"`.
- Order: decode → validate → store. Invalid data never reaches the database.
- `required` can't tell a missing `int` from `0`; use `*int` when zero is valid.
- Database constraints are the last line of defence, and catch what validation can't see (e.g. a duplicate email → 409).

### Path/query params

- **Path params identify _which_ resource:** `GET /boards/{id}`, read with `r.PathValue("id")`.
- **Query params filter or shape the result:** `/columns?board_id=…`, read with `r.URL.Query().Get("board_id")`. Returns `""` if missing.
- The router matches path params and methods (wrong method → 405); it ignores query params, so the handler checks them.
- Both arrive as strings; convert numbers with `strconv.Atoi`.

### Consistent errors

Every error has one shape, so the client handles them all the same way:

```json
{ "error": { "code": "BOARD_NOT_FOUND", "message": "Board not found" } }
```

An error response carries two different codes, in two different places:

```
HTTP/1.1 404 Not Found            ← status line: HTTP status code
Content-Type: application/json    ← headers

{"error": {"code": "BOARD_NOT_FOUND", "message": "Board not found"}}   ← body: application error code
```

- The **HTTP status code** (404) is in the status line, before the headers. It gives the general category, which any HTTP client understands.
- The **application error code** (`BOARD_NOT_FOUND`) is in the JSON body. It says exactly what went wrong, which only clients of this API understand.
- **Codes** are for code to branch on; **messages** are for people.
- Codes are constants in `handlers/respond.go`; `writeError` builds every error.
- 500s log the real error and send a generic message, never internal details.
- The store returns sentinel errors (`store.ErrEmailTaken`); the handler checks them with `errors.Is` and picks the status.
- Why not check `23505` in the handler? Only the store knows Postgres exists. Handlers depend on the `Store` interface, so a test fake (or a different database) can return `ErrBoardExists` with no Postgres involved.

### Content-Type

A header that tells the client how to read the body: `application/json`, `text/plain`, `text/html`.

- Order matters: set headers → `WriteHeader` → write body. Headers set later are silently ignored.
- If you don't set it, Go guesses, and JSON comes out as `text/plain`.
- `http.Error` always sends `text/plain`, so it can't send JSON errors.

### Request lifecycle

What happens between a client connecting and the response arriving.

```
TCP handshake → (TLS handshake) → accept → goroutine → parse → mux → handler → response
```

- **TCP (Transmission Control Protocol) handshake** opens a reliable connection in three steps:
  1. **SYN** (synchronize), client: "I want to talk to you."
  2. **SYN-ACK** (synchronize-acknowledge), server: "I hear you, and I want to talk to you too."
  3. **ACK** (acknowledge), client: "Got it, connection established."
- **TLS (Transport Layer Security) handshake** runs on top of TCP, for HTTPS only: encryption plus server identity. This project uses plain HTTP (`ListenAndServe`), so local traffic is unencrypted. `localhost` traffic never leaves the machine, though.
- **One goroutine per connection.** The server accepts connections in a loop; a slow client doesn't block others.
- **Parse:** raw bytes become an `*http.Request`. `r.Body` is a stream, read later by the handler.
- **Route:** the mux matches method and path, then calls your handler with `(w, r)`.
- **Respond:** the handler writes to `w`; when it returns, the response is finished and `r.Context()` is cancelled.
- **Keep-alive:** the connection usually stays open for the client's next request.
- Handlers run concurrently, so shared state must be safe for concurrent use. `*sql.DB` is safe for concurrent use.

### `context.Context`

Carries a cancellation signal (plus deadlines and request values) down a call chain.

- `r.Context()` is cancelled when the client disconnects or the handler returns.
- Pass it to the database (`QueryContext`, `ExecContext`, `BeginTx`) so abandoned queries stop and free their connection.
- Conventions: first parameter, named `ctx`; never stored in a struct; use `context.Background()` when there's no request.

## 3. Authentication / Authorization

### Password hashing

Never store passwords, only a hash of them. A hash is one-way: the same input always gives the same output, but the output can't be turned back into the password.

```
sign-up:  store bcrypt(password)
login:    bcrypt(what they typed) == stored hash?  →  correct password
```

- **Slow on purpose.** Attackers guess passwords and hash each guess. Fast hashes (SHA-256) allow billions of guesses a second; bcrypt and Argon2id take about 100ms each. The cost factor (`bcrypt.DefaultCost` = 10) sets how slow.
- **Salt:** a random value mixed in before hashing, so identical passwords get different hashes. It defeats precomputed lookup tables (rainbow tables). It isn't secret.
- **bcrypt stores everything in one string:** `$2a$10$<salt><hash>` holds the version, cost, salt and hash, so `CompareHashAndPassword` can repeat the hash at login.
- **bcrypt vs Argon2id:** Argon2id is newer and memory-hard (resists GPUs); bcrypt is simpler and still safe. This project uses bcrypt.
- bcrypt only uses the first 72 **bytes**, so validate the password length (here 8–72) and handle `bcrypt.ErrPasswordTooLong`.
- In Go: the tag `json:"-"` on `PasswordHash` keeps the hash out of every JSON response. The column is `password_hash TEXT NOT NULL`.

### Sessions and cookies

HTTP is stateless: each request arrives with no memory of earlier ones. A session lets the server recognise a logged-in user without resending the password.

```
POST /login  → check password → create session → Set-Cookie: session=<token>
GET /boards  → browser sends Cookie: session=<token> → look up session → user
```

- **Session:** a server-side record mapping a token to a user and an expiry (`sessions` table).
- **Session token:** 32 random bytes from `crypto/rand` (never `math/rand`, which is predictable). Possessing it *is* being logged in.
- **Store only its hash** (`token_hash`). A leaked table must not contain usable tokens. SHA-256 is fine here: a random token can't be guessed, so bcrypt's slowness would only slow every request.
- **Cookie flags:** `HttpOnly` (JS can't read it: XSS), `Secure` (HTTPS only), `SameSite` (not sent on other sites' requests: CSRF), `Expires`/`Max-Age`.
- **Logout** = delete the session row; it stops working immediately. Log out everywhere: `DELETE FROM sessions WHERE user_id = $1`.
- `http.SetCookie` writes a header, so it must come before the body.

### Sessions vs JWTs

A **JWT** (JSON Web Token) holds the user's details itself (`{"user_id":…,"exp":…}`), signed by the server. The server checks the signature and trusts the contents, with no lookup.

| | Sessions | JWT |
| --- | --- | --- |
| Client holds | A random ID, meaningless alone | The data itself, signed |
| Per request | One database lookup | Verify the signature only |
| Revoke | Delete the row: instant | Valid until it expires (unless you add a blocklist, which is a lookup again) |
| Stale data | No, looked up fresh | Yes, e.g. an old role stays in old tokens |

- JWTs are usually **signed, not encrypted**: anyone can read the payload; the signature only stops changes.
- JWTs suit many services verifying users without a shared database. For one browser app with one backend, sessions are simpler and logout really works.

### User enumeration

Using an app's responses to find out which emails have accounts.

- Login returns the same `401 INVALID_CREDENTIALS` "Invalid email or password" for an unknown email and a wrong password.
- It also takes the same time: for an unknown email it still runs bcrypt against a dummy hash. Otherwise the instant reply would give it away.
- 401 Unauthorized means "we don't know who you are".

### Database constraints

The database enforces what must never happen, even if the Go code has a bug.

- `NOT NULL`, `UNIQUE` (`23505` → 409), and **foreign keys**: `user_id UUID REFERENCES users(id)` means a session must point at a real user.
- Deleting a row that others still reference fails (`23503`), unless the foreign key says what to do: `ON DELETE CASCADE` deletes the referencing rows too.
- No foreign key at all would allow orphaned rows.

### Auth middleware

Middleware is code that runs before a handler, on every request to that handler. In Go it's a function that takes a handler and returns a new handler that wraps it: `func(next http.Handler) http.Handler`. The new handler can do some work first, then call `next.ServeHTTP(w, r)` to run the original handler, or skip that call to stop the request there.

`RequireAuth` (`handlers/middleware.go`) uses this to protect routes. `Me` is not part of it: it's an ordinary handler that returns the logged-in user's details, which a frontend can call to check whether someone is logged in. When a request for `GET /me` arrives, the router passes it to `RequireAuth` first. `RequireAuth` reads the `session` cookie, hashes its token, and calls `GetSessionUser`. If there's no cookie, or the query finds no unexpired session, it writes `401 UNAUTHENTICATED` and returns. Because it never calls `next`, the `Me` handler doesn't run.

If the session is valid, the handler needs to know who the user is. `RequireAuth` passes that on through the request's context: `context.WithValue` creates a new context holding the user, `r.WithContext` creates a copy of the request carrying that context, and `next` is called with the copy. `Me` then reads the user back with `userFromContext(r.Context())`.

The route is protected by wrapping it when it's registered:

```go
mux.Handle("GET /me", auth.RequireAuth(http.HandlerFunc(auth.Me)))
```

The check is written once, a new route can't skip it by accident, and handlers keep their own job. If `userFromContext` finds no user, the route was registered without `RequireAuth`. That's a bug in the code, so `Me` returns a 500.

Context values should only hold data about the current request, such as the user. The key's type is unexported, so other packages can't read or overwrite it.

### Logout

Logging out has a server half and a browser half. `POST /logout` deletes the session row (`DELETE FROM sessions WHERE token_hash = $1`), then sends `Set-Cookie: session=; Max-Age=0` so the browser deletes the cookie. The server half is the one that matters: once the row is gone, the token matches nothing and `RequireAuth` returns 401, even if a copy of the cookie is sent again. Clearing the cookie only stops the browser sending a token that no longer works.

Logout always returns `204`, with or without a cookie, because the result is the same either way. That makes it idempotent. It isn't wrapped in `RequireAuth`, so a user whose session has already expired still gets their cookie cleared.

It uses `POST`, not `GET`, because a `GET` should never change anything: browsers, link previews and other sites can trigger `GET` requests on their own. `DELETE /session` would also be valid REST, but `POST /logout` pairs with `POST /login`, and the session ID can't go in a URL anyway.

In Go, `MaxAge: -1` is what sends `Max-Age=0`, because `MaxAge: 0` means "don't send one". The cookie's `Path` must match the original, or the browser won't delete it.

Not yet covered: ownership checks (403 vs 404).

Go syntax lives in `LEARNING_GO.md`.
