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

- **URL = Uniform Resource Locator:** the address of a resource, saying where it is and how to reach it. In `http://localhost:8080/boards/123?page=2`: `http` is the scheme, `localhost:8080` the host and port, `/boards/123` the path, `?page=2` the query.
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
  - Why it can be public: bcrypt generates the salt randomly, rather than a person choosing it, so an attacker can't guess it in advance and build a lookup table for it. Once they've stolen the database they can read every salt, but that only lets them guess passwords one user at a time, with bcrypt's slow hash on every guess.
- **Nothing stops guessing; it can only be made impractical.** An attacker with the stolen database can always take a guess, hash it with a user's salt and compare. The defences decide how long that takes:

  | Defence           | What it does                                                                                                                                                |
  | ----------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- |
  | Salt              | Removes the shortcuts: pre-built lookup tables don't work, and each user must be guessed separately (500 users with `password123` get 500 different hashes) |
  | bcrypt's slowness | Slows each guess down: about 10 guesses a second per CPU core, instead of billions with SHA-256                                                             |
  | A strong password | Makes the number of guesses needed enormous                                                                                                                 |

  Worked example at 10 guesses a second: `password123` is in the first few thousand guesses of any list, so it falls in minutes. A random 12-character password has about 10²³ possibilities, so it would take billions of years.

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
- **Session token:** 32 random bytes from `crypto/rand` (never `math/rand`, which is predictable). Possessing it _is_ being logged in.
- **Store only its hash** (`token_hash`). A leaked table must not contain usable tokens. SHA-256 is fine here: a random token can't be guessed, so bcrypt's slowness would only slow every request.
- **Cookie flags:** `HttpOnly` (JS can't read it: XSS), `Secure` (HTTPS only), `SameSite` (not sent on other sites' requests: CSRF), `Expires`/`Max-Age`.
- **Logout** = delete the session row; it stops working immediately. Log out everywhere: `DELETE FROM sessions WHERE user_id = $1`.
- `http.SetCookie` writes a header, so it must come before the body.

### Sessions vs JWTs

A **JWT** (JSON Web Token) holds the user's details itself (`{"user_id":…,"exp":…}`), signed by the server. The server checks the signature and trusts the contents, with no lookup.

|              | Sessions                       | JWT                                                                          |
| ------------ | ------------------------------ | ---------------------------------------------------------------------------- |
| Client holds | A random ID, meaningless alone | The data itself, signed                                                      |
| Per request  | One database lookup            | Verify the signature only                                                    |
| Revoke       | Delete the row: instant        | Valid until it expires (unless you add a blocklist, which is a lookup again) |
| Stale data   | No, looked up fresh            | Yes, e.g. an old role stays in old tokens                                    |

- JWTs are usually **signed, not encrypted**: anyone can read the payload; the signature only stops changes.

  **What a signature is.** A JWT is three base64 strings joined by dots: `header.payload.signature`. The payload is only encoded, not hidden: decoding `eyJ1c2VyX2lkIjoiMSJ9` gives `{"user_id":"1"}`, which anyone can do. The signature is a hash (commonly HMAC-SHA256) of the header and payload combined with a **secret key** that only the server knows:

  ```
  signature = HMAC-SHA256(header + "." + payload, secretKey)
  ```

  **The token never changes between requests.** The server creates the token once, at login, and the client sends back exactly the same token on every request. What changes between requests is the URL and the request body; the token just rides alongside, like our session cookie:

  ```
  Login     → server creates token: {"user_id":"1"}, signature X
  Request 1 → client sends the same token: {"user_id":"1"}, X   ✅
  Request 2 → client sends the same token: {"user_id":"1"}, X   ✅
  Attacker  → edits the payload:   {"user_id":"2"}, X   ❌ signature doesn't match
  ```

  So any payload that differs from what the server signed must have been changed by someone else. When the payload legitimately needs to change (e.g. the token expires), the **server** issues a whole new token with a new signature. The client never edits a token; it only stores it and sends it back.

  To check a token, the server recalculates the signature from the header and payload it received, using its secret key, and compares it with the signature on the token. If they match, the payload hasn't been changed since the server signed it.

  **Worked example.** The server issues a token with payload `{"user_id":"1"}` and signature `X`. An attacker changes the payload to `{"user_id":"2"}` but keeps `X`. The server recalculates the signature over `{"user_id":"2"}` and gets `Y`, a different value, so it rejects the token. To make a matching signature for the new payload, the attacker would need the secret key, which never leaves the server.

  **Common mix-ups:**
  1. The signature is made from the **header and payload**, not "payload and signature". The signature doesn't exist yet when it's being calculated.
  2. The **secret is never attached** to the JWT. Only the result of the hash, the signature, goes on the token. The secret never leaves the server. If it were attached, anyone could use it to sign their own tokens.
  3. On the way back, the server **hashes the header and payload with its secret**, then compares the result with the signature on the token.

  ```
  Creating:   signature = hash(header + payload, secret)
              token     = header.payload.signature        ← secret is NOT in here

  Checking:   hash(header + payload received, secret) == signature received?
              yes → unchanged, made by this server
              no  → rejected
  ```

  **One secret for every user.** The secret is a single value stored in an environment variable on the server (e.g. `JWT_SECRET` in `.env`, read at startup like `DB_PASSWORD`), so the same secret signs every user's token. Each token's signature still differs, because each payload differs: `hash({"user_id":"1"}, secret)` gives `X`, `hash({"user_id":"2"}, secret)` gives `Y`. Never put the secret in code, git or the database: anyone who has it can forge a token for any user.

  So a signature proves the payload is **unchanged and came from the server**, but it doesn't hide it. Never put secrets in a JWT payload.

- **Access and refresh tokens (a JWT pattern):** they exist to fix JWT's revocation problem.
  - **Access token:** a short-lived JWT (e.g. 15 minutes), sent on every request and checked with no database lookup.
  - **Refresh token:** a long-lived random token (e.g. 30 days) **stored in the database**. When the access token expires, the client sends the refresh token to get a new one.
  - Logout deletes the refresh token. The current access token keeps working until it expires (at most 15 minutes), then no new one can be issued.
  - Sessions don't need this: every request already checks the `sessions` table, so deleting the row logs the user out immediately. A refresh token is essentially a session; the pattern is "JWTs for speed on most requests, plus a session-like token for revocation".
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

### Ownership checks (403 vs 404)

Authentication answers "who are you?"; authorization answers "are you allowed to do this?". Being logged in isn't enough: Bo must not read or change Ada's boards, even if he knows their IDs.

Each board has an `owner_id`. Every query includes the logged-in user, taken from the session through `currentUser`, never from the request body. `GET /boards/b1` runs `WHERE id = $1 AND owner_id = $2`. If Bo asks for Ada's board, that matches no rows, exactly as if the board didn't exist, so he gets `404 BOARD_NOT_FOUND`. Updates and deletes do the same with `RowsAffected() == 0`.

Columns and cards have no owner of their own, so their queries join up to the board: `cards → columns → boards`, checking `boards.owner_id`. Creating one uses `INSERT … SELECT … WHERE EXISTS (…)`, which inserts nothing unless the parent board or column is yours. Updating a card checks both its current column and its new one, so a card can't be moved onto someone else's board.

The response for something that exists but isn't yours is `404`, not `403`. A `403 Forbidden` would confirm the resource exists, which leaks information in the same way as user enumeration. `403` is for when the user may know the thing exists but not do this action, such as a read-only shared board. This app has no sharing, so it always uses `404`.

Putting the check inside the SQL, rather than loading the row and comparing owners in Go, means one query and nothing to forget.

Step 3 complete.

## 4. Transactions

### What a transaction is (ACID)

A transaction groups several SQL statements so they succeed or fail as one unit.

Example: moving card `c1` from `todo` to `done` is two writes:

1. `UPDATE cards SET column_id = 'done' WHERE id = 'c1'`
2. `INSERT INTO card_movements (card_id, from_column_id, to_column_id) VALUES ('c1', 'todo', 'done')`

Without a transaction, if the server crashes between 1 and 2, the card is in `done` but the history has no record of the move. With a transaction, either both rows change or neither does.

ACID is the four guarantees a transaction gives:

- **Atomicity:** all or nothing. A failure part-way undoes every earlier statement.
- **Consistency:** constraints (`NOT NULL`, `UNIQUE`, `FOREIGN KEY`) hold at the end, or the transaction fails. _Essentially, if any constraint violation occurs, then the whole transaction fails_.
- **Isolation:** other connections don't see half-finished changes. How strictly is set by the isolation level (step 5).
- **Durability:** once `COMMIT` returns, the change survives a crash or power cut.

A single statement is already its own transaction. You only need an explicit one when two or more statements must stay in step.

### `BeginTx` / `Commit` / `Rollback`

The three calls that start and end a transaction in `database/sql`:

- `tx, err := s.db.BeginTx(ctx, nil)` takes one connection from the pool, sends `BEGIN`, and returns a `*sql.Tx` tied to that connection. `nil` means default options (Postgres's default isolation level, `READ COMMITTED`).
- `tx.Commit()` sends `COMMIT`: every change becomes permanent and visible to other connections at the same moment.
- `tx.Rollback()` sends `ROLLBACK`: every change since `BEGIN` is thrown away, as if it never ran.

After either one, the connection goes back to the pool and the `tx` can't be used again (calls return `sql.ErrTxDone`).

`ctx` is `r.Context()`. If the client disconnects or times out before `Commit`, `database/sql` rolls the transaction back automatically.

### `defer tx.Rollback()`

Call `defer tx.Rollback()` straight after `BeginTx` succeeds, so every early `return` rolls back.

```go
tx, err := s.db.BeginTx(ctx, nil)
if err != nil {
    return false, err
}
defer tx.Rollback()
```

`MoveCard` has six `return`s between `BeginTx` and `Commit`. Without the `defer`, each would need its own `tx.Rollback()`, and forgetting one leaves the transaction open: the connection stays checked out, and any row locks stay held.

If `Commit` succeeds, the deferred `Rollback` still runs, but it does nothing and returns `sql.ErrTxDone`, which is ignored. That's why it's safe to always defer it.

### Every query on `tx`, not `db`

Inside a transaction, every statement must go through `tx` (`tx.QueryRowContext`, `tx.ExecContext`), never `s.db`.

`s.db` is a pool of connections. A transaction lives on one connection. `s.db.ExecContext(...)` borrows a different connection, so that statement:

- runs outside the transaction and commits on its own straight away
- isn't undone by `Rollback`
- can't see the transaction's uncommitted changes
- can get stuck waiting on a row the transaction has locked

The method names are the same on both (`ExecContext`, `QueryRowContext`), so the mistake compiles and looks right. Check the receiver.

### `MoveCard` (`POST /cards/{id}/move`)

A dedicated endpoint for moving a card, separate from `PUT /cards/{id}`, because a move does more than change one field: it also records history.

Request: `POST /cards/c1/move` with `{"column_id": "done"}`. Success is `204`.

Inside one transaction, `store.MoveCard`:

1. Reads the card's current column (`todo`), checking via `boards.owner_id` that it's yours.
2. Checks the destination column exists and is on one of your boards.
3. `UPDATE cards` to the new column.
4. `INSERT INTO card_movements` with `from_column_id = todo`, `to_column_id = done`.
5. `Commit`.

`card_movements` is an append-only history: rows are added, never updated. Step 1 needs the old column before it's overwritten, which is why the read happens inside the same transaction.

Step 1 also uses `SELECT ... FOR UPDATE`, which locks the card row until commit. That's for two people moving the same card at once (step 5).

### Mapping transaction failures to API errors

The store returns three kinds of result, and the handler turns each into a response:

| Store returns       | Meaning                                 | Response                      |
| ------------------- | --------------------------------------- | ----------------------------- |
| `true, nil`         | committed                               | `204`                         |
| `false, nil`        | card missing or not yours               | `404 CARD_NOT_FOUND`          |
| `ErrColumnNotFound` | destination column missing or not yours | `404 COLUMN_NOT_FOUND`        |
| any other error     | `BeginTx`, a query, or `Commit` failed  | `500 INTERNAL_ERROR` (logged) |

In every failure case the deferred `Rollback` has already undone the work, so the database is unchanged. The client never gets a half-done move.

`ErrColumnNotFound` is a sentinel error: a package-level `var` the handler checks with `errors.Is`. That way the store reports _what_ went wrong without knowing about HTTP status codes.

Go syntax lives in `LEARNING_GO.md`.
