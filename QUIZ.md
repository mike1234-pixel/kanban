# Quiz

Recall questions on the concepts in `LEARNING.md`. Try answering before opening each answer.

## 1. CRUD

**1.** What does CRUD stand for, and which HTTP method goes with each?

- [x] Create `POST`, Read `GET`, Update `PUT`, Delete `DELETE`
- [ ] Create `PUT`, Read `GET`, Update `POST`, Delete `DELETE`
- [ ] Create `POST`, Read `GET`, Update `PATCH`, Delete `POST`
- [ ] Copy `GET`, Read `GET`, Update `PUT`, Drop `DELETE`

<details><summary>Answer</summary>

Example for boards: `POST /boards` creates one, `GET /boards/123` reads it, `PUT /boards/123` updates it, `DELETE /boards/123` deletes it.

**The method says what to do; the path says which resource.** That's why the same URL can take four methods.
</details>

**2.** What does *idempotent* mean, and which of `GET`, `POST`, `PUT`, `DELETE` is not idempotent?

- [ ] Repeating it has the same effect as once; `DELETE` isn't
- [ ] It never changes data; `PUT` isn't
- [ ] It always returns the same body; `GET` isn't
- [x] Repeating it has the same effect as once; `POST` isn't

<details><summary>Answer</summary>

Example: sending `DELETE /boards/123` twice leaves the board deleted either way. Sending `POST /boards` twice **creates two boards**.

That's why it's safe for a browser to retry a `GET`, `PUT` or `DELETE`, but not a `POST`.
</details>

**3.** An `UPDATE` or `DELETE` matches no rows. Postgres doesn't return an error. How does the store know nothing was found?

- [ ] `Exec` returns `sql.ErrNoRows`
- [ ] `Exec` returns a `nil` result
- [x] `RowsAffected()` returns 0
- [ ] Postgres returns error code `23505`

<details><summary>Answer</summary>

Example: `DELETE FROM boards WHERE id = 'abc'` with no board `abc` runs fine and deletes 0 rows. That's not an error to Postgres.

So the store checks **`RowsAffected() == 0`** and returns `false`, and the handler turns that into `404`.
</details>

## 2. HTTP Fundamentals

### Status codes

**4.** What do the three status code classes 2xx, 4xx and 5xx mean?

- [ ] 2xx success, 4xx server's fault, 5xx client's fault
- [x] 2xx success, 4xx client's fault, 5xx server's fault
- [ ] 2xx created, 4xx redirect, 5xx timeout
- [ ] 2xx success, 4xx warning, 5xx client's fault

<details><summary>Answer</summary>

Example: `201` means the board was created. `400` means the JSON you sent was broken, so **you fix your request**. `500` means the database failed, so **the server needs fixing**; retrying the same request may work later.

The first digit tells the client who's at fault before it reads anything else.
</details>

**5.** What do 201, 204, 404 and 409 mean?

- [x] 201 Created, 204 No Content, 404 Not Found, 409 Conflict
- [ ] 201 OK, 204 Deleted, 404 Not Found, 409 Forbidden
- [ ] 201 Created, 204 Not Modified, 404 Gone, 409 Conflict
- [ ] 201 Accepted, 204 No Content, 404 Bad Request, 409 Duplicate

<details><summary>Answer</summary>

Examples from this API: `POST /boards` returns **`201`** with the new board. `DELETE /boards/123` returns **`204`** with an empty body. `GET /boards/999` returns **`404`**. Creating a user with an email that's taken returns **`409`**.
</details>

**6.** How many times can a handler set the status code?

- [ ] As many times as it likes; the last one wins
- [ ] Twice: once for headers, once for the body
- [ ] Never; Go always sends 200
- [x] Once; later `WriteHeader` calls are ignored

<details><summary>Answer</summary>

The status line is the first thing sent to the client, so once it's gone it can't be changed.

Example: a handler calls `WriteHeader(201)`, then later `WriteHeader(500)`. The client still gets **`201`**, and Go logs `superfluous response.WriteHeader call`.
</details>

### JSON encode/decode

**7.** What does decoding do, and what does encoding do?

- [ ] Decoding: Go struct → JSON. Encoding: JSON → Go struct
- [ ] Decoding decrypts the body; encoding encrypts it
- [x] Decoding: JSON → Go struct. Encoding: Go struct → JSON
- [ ] Decoding checks the JSON; encoding saves it to the database

<details><summary>Answer</summary>

Example: the client sends `{"title":"Sprint"}`. **Decoding** fills `models.Board{Title: "Sprint"}`. The handler saves it, then **encoding** turns the saved board back into `{"id":"…","title":"Sprint"}` for the response.

In TypeScript terms: decoding is `JSON.parse`, encoding is `JSON.stringify`.
</details>

**8.** What does a `nil` slice encode as, and what does an empty slice encode as?

- [ ] Both encode as `[]`
- [x] `nil` → `null`; empty → `[]`
- [ ] Both encode as `null`
- [ ] `nil` → `[]`; empty → `null`

<details><summary>Answer</summary>

Example: `var cards []Card` is `nil` and encodes as **`null`**. `cards := []Card{}` encodes as **`[]`**.

A frontend doing `data.map(…)` crashes on `null`, which is why list endpoints always start with an empty slice.
</details>

**9.** Which struct fields can `encoding/json` see?

- [x] Only exported fields (capital first letter)
- [ ] Only unexported fields (lowercase first letter)
- [ ] Every field, exported or not
- [ ] Only fields with a `json` struct tag

<details><summary>Answer</summary>

`encoding/json` is a different package, and **other packages can only see capitalised names**.

Example: in `type Board struct { ID string; title string }`, only `ID` appears in the JSON. `title` is silently left out.
</details>

### Validation

**10.** What's the difference between decoding and validating?

- [x] Decoding: JSON is readable. Validating: values make sense
- [ ] Decoding: values make sense. Validating: JSON is readable
- [ ] They're the same step with different names
- [ ] Decoding checks types; validating checks the database

<details><summary>Answer</summary>

Example: `{"title": ""}` **decodes fine**: it's valid JSON. But an empty title makes no sense, so **validation** rejects it with `400 VALIDATION_FAILED`.

`{"title": ` (cut off) fails earlier, at decoding, with `400 INVALID_JSON`.
</details>

**11.** In what order does a create handler decode, validate and call the store?

- [ ] Validate → decode → store
- [ ] Store → decode → validate
- [ ] Decode → store → validate
- [x] Decode → validate → store

<details><summary>Answer</summary>

Each step needs the one before it: you can't validate a struct until it's decoded, and **the store should only ever receive valid data**.

Example: `POST /boards` with `{"title":""}` stops at validation, so no SQL runs at all.
</details>

**12.** Validation can't catch a duplicate email. What does?

- [ ] A `validate:"unique"` tag, returned as 400
- [ ] The handler's `if` check, returned as 404
- [x] The database's `UNIQUE` constraint, returned as 409
- [ ] The `NOT NULL` constraint, returned as 500

<details><summary>Answer</summary>

Validation only sees the request; it can't know who's already in the `users` table.

Example: two sign-ups for `ada@example.com` arrive at once. Both pass validation, but **Postgres's `UNIQUE` rejects the second** (`23505`). The store returns `ErrEmailTaken`, and the handler returns `409 USER_EMAIL_TAKEN`.
</details>

### Path/query params

**13.** What does URL stand for, and what are the parts of `http://localhost:8080/boards/123?page=2`?

- [ ] Universal Request Link: protocol, server, route, params
- [x] Uniform Resource Locator: scheme, host and port, path, query
- [ ] Uniform Resource Locator: protocol, domain, folder, file
- [ ] Unique Resource Label: scheme, path, body, query

<details><summary>Answer</summary>

Each part answers a question: **how** to talk (`http`), **which machine** (`localhost:8080`), **which resource** (`/boards/123`), and **how to return it** (`?page=2`).

The router only matches on the method and the path; the query is read by the handler.
</details>

**14.** What's a path param for, and what's a query param for?

- [x] Path param: which resource. Query param: filter or shape it
- [ ] Path param: filter or shape it. Query param: which resource
- [ ] Path param: for `GET`. Query param: for `POST`
- [ ] Both identify the resource; query params are optional

<details><summary>Answer</summary>

Example: `/boards/123` names board 123, and **without the `123` there's no resource**. `/boards?archived=false&page=2` is still the boards list, just **filtered and paged**.

Rule of thumb: if leaving it out changes *what* you get, it's a path param; if it changes *how much or which order*, it's a query param.
</details>

**15.** What does a missing query param return in Go?

- [ ] `nil`
- [ ] An error
- [ ] The default value from the route
- [x] An empty string, `""`

<details><summary>Answer</summary>

Example: for `/boards` (no `?page=`), `r.URL.Query().Get("page")` returns **`""`, not an error**.

So the handler must decide what empty means, e.g. "use page 1".
</details>

**16.** What type do path and query params arrive as?

- [ ] Their real types, e.g. `int`
- [ ] `any`
- [x] Strings
- [ ] Bytes

<details><summary>Answer</summary>

A URL is just text, so Go can't know `123` is meant to be a number.

Example: `?page=2` gives you **the string `"2"`**. `strconv.Atoi("2")` turns it into `2`, and returns an error for `?page=abc`, which the handler turns into `400`.
</details>

### Consistent errors

**17.** In an error response, what's `code` for, and what's `message` for?

- [ ] `code`: the HTTP status. `message`: the Go error text
- [x] `code`: for the client's code. `message`: for people
- [ ] `code`: for people. `message`: for the client's code
- [ ] `code`: for logging. `message`: for the database

<details><summary>Answer</summary>

Example: `{"code":"BOARD_NOT_FOUND","message":"Board not found"}`. The frontend checks **`code === "BOARD_NOT_FOUND"`**, so you can reword the message later without breaking it.

The `message` can be shown to the user as it is.
</details>

**18.** Where does the HTTP status code go in a response, and where does the application error code go?

- [x] Status code: the status line. Error code: the JSON body
- [ ] Both go in the JSON body
- [ ] Status code: the JSON body. Error code: a header
- [ ] Both go in the response headers

<details><summary>Answer</summary>

Example response:

`HTTP/1.1 404 Not Found` ← **status code**, read by browsers, proxies and `fetch`'s `res.ok`
`{"error":{"code":"BOARD_NOT_FOUND", …}}` ← **error code**, read by your frontend

The status says the general kind of failure; the error code says exactly which one.
</details>

**19.** Why do 500 errors send a generic message instead of the real error?

- [ ] The real error is too long to send
- [ ] JSON can't encode Go errors
- [ ] The browser hides 500 messages anyway
- [x] It can leak database details, and the client can't act on it

<details><summary>Answer</summary>

Example: a real error like `pq: relation "boardz" does not exist` tells an attacker your table names. The user can't fix it anyway.

So the client gets **`500 INTERNAL_ERROR` with a generic message**, and the real error goes to `slog` for you to read.
</details>

**20.** Which layer picks the HTTP status code: the handler or the store?

- [ ] The store
- [ ] The router
- [x] The handler
- [ ] Postgres

<details><summary>Answer</summary>

Status codes are HTTP, and **only the handler knows about HTTP**.

Example: the store returns `false` for "no rows matched". The same `false` might be `404` from `GET /boards/123` but something else from another endpoint. The handler decides.
</details>

**21.** Why does the store turn Postgres error `23505` into `ErrBoardExists`, instead of the handler checking for `23505`?

- [ ] Handlers can't import `lib/pq`
- [x] Only the store should know Postgres exists
- [ ] `23505` changes between Postgres versions
- [ ] It makes the query run faster

<details><summary>Answer</summary>

`23505` is a Postgres detail. If handlers checked it, **every handler would depend on Postgres**.

Example: the fake store in `handlers_test.go` can return `ErrBoardExists` without any database, and the handler still returns `409`.
</details>

### Content-Type

**22.** What does the `Content-Type` header tell the client?

- [x] The format of the response body, e.g. `application/json`
- [ ] How long the response body is
- [ ] Which client the response is for
- [ ] Whether the response was cached

<details><summary>Answer</summary>

The body is just bytes; this header says how to read them.

Example: with **`Content-Type: application/json`**, `fetch`'s `res.json()` and tools like curl know to parse it as JSON. Without it, they may treat it as plain text.
</details>

**23.** In what order must a handler set headers, write the status and write the body?

- [ ] Status → headers → body
- [ ] Body → headers → status
- [ ] Any order; Go sorts them
- [x] Headers → status → body

<details><summary>Answer</summary>

They're sent in that order on the wire. **Once the status or body starts, the headers have already gone.**

Example: setting `Content-Type` after `w.Write(body)` does nothing, because the headers were already sent.
</details>

### Request lifecycle

**24.** What do TCP and TLS stand for?

- [ ] Transfer Connection Protocol; Transport Link Security
- [ ] Transmission Control Protocol; Total Line Security
- [x] Transmission Control Protocol; Transport Layer Security
- [ ] Transport Control Packet; Trusted Layer Service

<details><summary>Answer</summary>

**TCP** gets bytes from one machine to another reliably and in order. **TLS** encrypts those bytes.

HTTP is sent over TCP; HTTPS is HTTP sent over TLS over TCP.
</details>

**25.** What are the three steps of the TCP handshake?

- [ ] ACK → SYN → SYN-ACK
- [x] SYN → SYN-ACK → ACK
- [ ] HELLO → OK → READY
- [ ] SYN → ACK → FIN

<details><summary>Answer</summary>

Example: the browser sends **SYN** ("I want to connect"), the server replies **SYN-ACK** ("OK, I'm ready"), the browser sends **ACK** ("got it").

Only then can the HTTP request be sent. It's why a new connection costs a round trip before any data moves.
</details>

**26.** What does TLS add on top of TCP?

- [x] Encryption, and server identity via a certificate
- [ ] Faster connections and compression
- [ ] Guaranteed delivery and ordering
- [ ] User login and sessions

<details><summary>Answer</summary>

Without TLS, anyone on the network (e.g. café Wi-Fi) can **read your session cookie and use it**.

TLS encrypts the traffic, and the **certificate proves** you're talking to the real server, not an impostor.
</details>

**27.** Why does Go's server give each connection its own goroutine?

- [ ] Because each connection needs its own database
- [ ] So connections can share memory safely
- [ ] Because Go has no threads
- [x] So one slow client doesn't block the others

<details><summary>Answer</summary>

Example: Ada is on a slow phone connection, uploading slowly. Bo's `GET /boards` **doesn't have to wait for Ada**, because his connection has its own goroutine.

Goroutines are cheap, so a Go server can run thousands at once.
</details>

**28.** Which part of the server decides which handler a request goes to?

- [ ] The middleware
- [ ] The store
- [x] The mux (router)
- [ ] The `http.Server` TLS layer

<details><summary>Answer</summary>

Example: `POST /cards/c1/move` arrives. **The mux matches the method and path** against `mux.Handle("POST /cards/{id}/move", …)` and calls `cardHandler.MoveCard`.

No match means it returns `404` before any of your code runs.
</details>

### `context.Context`

**29.** What two things cancel a request's context?

- [ ] The database closing, or the server starting
- [x] The client disconnecting, or the handler returning
- [ ] A 404, or a 500
- [ ] Calling `defer`, or a panic

<details><summary>Answer</summary>

Example: Ada closes the tab mid-request. The context is cancelled, so **any query still running is abandoned**.

When the handler returns, the context is cancelled too, so nothing keeps working on a finished request.
</details>

**30.** Why pass `r.Context()` to database calls?

- [x] So the query is abandoned if the client disconnects
- [ ] So the query runs inside a transaction
- [ ] So the database knows which user is logged in
- [ ] So the query result is cached

<details><summary>Answer</summary>

Example: a slow report query takes 10 seconds. The user closes the tab after 1. With `r.Context()`, **Postgres stops after 1 second** and the connection goes back to the pool.

Without it, the query runs all 10 seconds for a response nobody will read.
</details>

**31.** By convention, where does a context parameter go, and what is it named?

- [ ] Last, named `ctx`
- [ ] First, named `context`
- [ ] Last, named `c`
- [x] First, named `ctx`

<details><summary>Answer</summary>

Example: `func (s *PostgresStore) GetBoard(ctx context.Context, id string)`.

It's a convention all Go code follows, so **anyone reading a function sees straight away that it can be cancelled**.
</details>

## 3. Authentication / Authorization

### Password hashing

**32.** Why must a password hash be slow?

- [ ] Slow hashes are harder to reverse
- [ ] It stops two users having the same hash
- [x] Each guess costs ~100ms, so mass guessing is impractical
- [ ] It hides the password from the database

<details><summary>Answer</summary>

Example: an attacker has a stolen hash and tries guesses. At 100ms each, **10 million guesses take about 12 days** per user. A fast hash like SHA-256 would do them in under a second.

One login waits 100ms once, which nobody notices.
</details>

**33.** What is a salt, and what attack does it stop?

- [ ] A secret key on the server; stops stolen databases
- [x] A random value per password; stops precomputed lookup tables
- [ ] A random value per password; stops slow guessing
- [ ] Extra characters users must add; stops weak passwords

<details><summary>Answer</summary>

Example: Ada and Bo both use `password123`. Without a salt, their hashes are identical, and **one lookup table cracks both**. With a salt each, their hashes differ.

A lookup table only works if it was built for that exact salt, which is impossible when every user's salt is random.
</details>

**34.** Why doesn't the salt need to be secret?

- [x] It's random, so lookup tables can't be built in advance
- [ ] It's encrypted with the server's key
- [ ] It's deleted after the hash is made
- [ ] Attackers can't read the `users` table

<details><summary>Answer</summary>

Its job is to make every user's hash different, **not to be unknown**.

Example: the bcrypt hash `$2a$10$N9qo8uLOickgx2ZMRZoMye…` stores the salt (`N9qo8uLOickgx2ZMRZoMye`) right inside it. Login needs it to check the password.
</details>

**35.** Can anything stop an attacker with a stolen database from guessing passwords? What makes guessing impractical?

- [ ] Yes; the salt stops guessing completely
- [ ] Yes; bcrypt hashes can't be compared
- [ ] No; nothing makes guessing any slower
- [x] No; salt, bcrypt's slowness and strong passwords make it impractical

<details><summary>Answer</summary>

Example: an attacker with the database tries `password123` against Ada's hash. If that's her password, they're in. **Nothing stops a guess from working**; you can only make guessing slow.

That's why a long, unusual password is the final defence.
</details>

**36.** How does login check a password, if a hash can't be reversed?

- [ ] Reverses the stored hash and compares the passwords
- [ ] Asks bcrypt to decrypt the stored password
- [x] Hashes the typed password with the same salt and cost, then compares
- [ ] Compares the typed password with the salt

<details><summary>Answer</summary>

Example: Ada types `correct-horse`. `bcrypt.CompareHashAndPassword` **reads the salt and cost from her stored hash**, hashes `correct-horse` the same way, and checks the two results match.

The password itself is never stored or decrypted.
</details>

**37.** Which struct tag keeps a field out of JSON responses?

- [ ] `json:"omit"`
- [x] `json:"-"`
- [ ] `json:"private"`
- [ ] `json:",omitempty"`

<details><summary>Answer</summary>

Example: `` PasswordHash string `json:"-"` `` means **encoding skips the field**, so `GET /users/123` never sends the hash, even by accident.
</details>

### Sessions and cookies

**38.** HTTP is stateless. How does the server know who sent a request after login?

- [x] A cookie with a session token, looked up in `sessions`
- [ ] The browser's IP address
- [ ] The user ID in the URL
- [ ] The password sent with every request

<details><summary>Answer</summary>

Example: Ada logs in, and the response sets `session=Xy7…`. On `GET /boards`, the browser **sends the cookie automatically**. The server hashes `Xy7…`, finds the row in `sessions`, and knows it's Ada.

Each request carries its own proof of who's sending it.
</details>

**39.** Why store a hash of the session token rather than the token itself?

- [ ] Hashes take less space in the table
- [ ] Tokens can't be stored in Postgres
- [ ] Hashing makes the lookup faster
- [x] A leaked token is usable as it is; a hash isn't

<details><summary>Answer</summary>

Example: someone copies the `sessions` table. A raw token `Xy7…` could be put straight into a cookie and **they'd be logged in as Ada**. Its SHA-256 can't be: the server would hash it again and find no match.
</details>

**40.** Why is SHA-256 fine for session tokens, when passwords need bcrypt?

- [ ] SHA-256 is slower than bcrypt
- [ ] Tokens expire, so they don't need hashing
- [x] Tokens are long and random, so they can't be guessed
- [ ] bcrypt can't hash random bytes

<details><summary>Answer</summary>

A session token is 32 random bytes: there are far too many possible tokens to guess one. **Slowness only helps against guessable input**, like passwords.

And the token is hashed on every request, so 100ms of bcrypt would slow every page.
</details>

**41.** What do the cookie flags `HttpOnly`, `Secure` and `SameSite` do?

- [ ] `HttpOnly`: HTTP only. `Secure`: encrypted. `SameSite`: same server only
- [x] `HttpOnly`: no JS. `Secure`: HTTPS only. `SameSite`: not on other sites' requests
- [ ] `HttpOnly`: no HTTPS. `Secure`: signed. `SameSite`: one tab only
- [ ] `HttpOnly`: no JS. `Secure`: encrypted value. `SameSite`: no cross-origin JS

<details><summary>Answer</summary>

Each one blocks a different attack:

- **`HttpOnly`:** an injected script can't do `document.cookie` and steal the session.
- **`Secure`:** it's never sent over plain HTTP, where anyone on the network could read it.
- **`SameSite=Lax`:** a form on `evil.com` can't make a `POST` with Ada's cookie.
</details>

### Sessions vs JWTs

**42.** What's the main difference between a session token and a JWT?

- [x] Session: a random ID looked up. JWT: signed details, no lookup
- [ ] Session: signed details. JWT: a random ID looked up
- [ ] Session: stored in the browser. JWT: stored in the database
- [ ] They're the same; JWT is the newer name

<details><summary>Answer</summary>

Example: a session token `Xy7…` means nothing until the server **looks it up** in `sessions`. A JWT says `{"user_id":"123"}` plus a signature, so the server **trusts it without a database query**.

Session: one lookup per request, easy to revoke. JWT: no lookup, hard to revoke.
</details>

**43.** How do you log a user out everywhere with sessions, and why is it hard with JWTs?

- [ ] Change their password; JWTs can't be changed
- [ ] Clear their cookies; JWTs aren't cookies
- [ ] Delete the user; JWTs are copied to other servers
- [x] Delete their session rows; a JWT has nothing stored to delete

<details><summary>Answer</summary>

Example: Ada's laptop is stolen. With sessions, `DELETE FROM sessions WHERE user_id = 'ada'` **logs her out everywhere at once**.

A JWT is checked by its signature alone, so the stolen one keeps working until it expires.
</details>

**44.** Is a JWT's payload secret?

- [ ] Yes; it's encrypted with the server's key
- [ ] Yes; only the browser can read it
- [x] No; it's signed, not encrypted, so anyone can read it
- [ ] No; it's sent as plain text without a signature

<details><summary>Answer</summary>

The payload is just base64, which anyone can decode. Paste a JWT into jwt.io and you'll see `{"user_id":"123"}`.

The signature only stops it being **changed**, not **read**. So never put secrets in it.
</details>

**45.** How does the server check a JWT's signature?

- [ ] Decrypts the signature to get the payload
- [x] Recalculates it from header + payload with its secret, and compares
- [ ] Looks the token up in the `sessions` table
- [ ] Checks the signature against the user's password hash

<details><summary>Answer</summary>

Example: the token is `header.payload.signature`. The server works out `HMAC(header.payload, secret)` and **compares it with the signature part**.

If someone changed `user_id` from 123 to 456, the recalculated signature won't match, and the token is rejected.
</details>

**46.** Why can't an attacker change a JWT's payload and make a new valid signature?

- [x] A valid signature needs the server's secret key
- [ ] The payload is encrypted
- [ ] The browser blocks changes to cookies
- [ ] The signature includes the user's password

<details><summary>Answer</summary>

Example: an attacker changes `"user_id":"123"` to `"456"`. To make a matching signature, they'd need to work out `HMAC(…, secret)`, and **without the secret they can't**.

So the signature proves the server created the token and nobody changed it.
</details>

**47.** Does a JWT's payload change from one request to the next? Who makes a new one?

- [ ] Yes; the client re-signs it on every request
- [ ] Yes; the payload holds the current URL
- [ ] No; and it can never be replaced
- [x] No; the same token each time. The server issues a new one

<details><summary>Answer</summary>

Example: Ada logs in and gets a token for `user_id 123`. **Every request sends that exact token.** If her role changes, the old token still says the old role until the server **issues a new one**.

Like a cookie, it's the same proof each time.
</details>

**48.** What are access and refresh tokens, and what problem do they solve?

- [ ] Two copies of the same JWT; fixes lost tokens
- [ ] A JWT + a password; fixes weak passwords
- [x] Short-lived JWT + long-lived stored token; fixes revocation
- [ ] A session + a cookie; fixes CSRF

<details><summary>Answer</summary>

Example: the access token lasts 15 minutes, the refresh token 30 days. Every 15 minutes, the client swaps the refresh token for a new access token.

To log Ada out, delete her refresh token: **within 15 minutes, her access token expires and can't be renewed**.
</details>

**49.** Why don't sessions need refresh tokens?

- [ ] Sessions never expire
- [x] Every request checks `sessions`, so deleting the row logs out now
- [ ] Sessions are stored in the browser
- [ ] Refresh tokens only work with cookies

<details><summary>Answer</summary>

With JWTs, the refresh token exists only so there's something stored to delete. Sessions **are already stored**.

Example: `DELETE FROM sessions WHERE token_hash = …` and Ada's next request gets `401`, straight away.
</details>

### User enumeration

**50.** What is user enumeration?

- [x] Using responses to find out which emails have accounts
- [ ] Listing every user with `GET /users`
- [ ] Numbering users with sequential IDs
- [ ] Guessing a user's password many times

<details><summary>Answer</summary>

Example: login says "No account with that email" for `bo@example.com` but "Wrong password" for `ada@example.com`. An attacker now knows **Ada has an account** and can target her.

Even the time a response takes can give it away.
</details>

**51.** What two things does login keep identical for an unknown email and a wrong password?

- [ ] The cookie, and the session row
- [ ] The email, and the password hash
- [ ] The status code, and the request body
- [x] The error message (and status), and the response time

<details><summary>Answer</summary>

Example: for an unknown email, login **still runs bcrypt** on a `dummyHash`, so it takes ~100ms like a real check. Both cases return `401 INVALID_CREDENTIALS`.

Without the dummy hash, unknown emails would answer much faster, and give themselves away.
</details>

### Database constraints

**52.** What happens if you delete a user who still has sessions, when `sessions.user_id` is `REFERENCES users(id)` without `ON DELETE CASCADE`?

- [ ] The sessions are deleted too
- [ ] The user is deleted and sessions keep a dangling `user_id`
- [x] Postgres refuses it with a foreign key violation (`23503`)
- [ ] Postgres sets `sessions.user_id` to `NULL`

<details><summary>Answer</summary>

Example: `DELETE FROM users WHERE id = 'ada'` while she has two session rows. Those rows would point at a user who doesn't exist, so **Postgres refuses** with `23503`.

The foreign key protects you from orphan rows.
</details>

**53.** What does `ON DELETE CASCADE` do?

- [ ] Deleting a child row deletes its parent
- [x] Deleting the parent row deletes the rows that reference it
- [ ] It stops the parent row being deleted
- [ ] It deletes the whole table

<details><summary>Answer</summary>

Example: deleting Ada's user row **also deletes her two sessions** in the same statement. Deleting a board deletes its columns, and through those, its cards.

You don't need code to clean up the children.
</details>

### Auth middleware

**54.** What is middleware?

- [x] A function that wraps a handler and runs on every request to it
- [ ] A layer between the store and Postgres
- [ ] A goroutine that runs in the background
- [ ] A struct that holds the database pool

<details><summary>Answer</summary>

Example: `RequireAuth` wraps `cardHandler.MoveCard`. Each request to it **runs `RequireAuth` first**, which checks the session, then calls `next.ServeHTTP` to run the handler.

It's like a wrapper component in React: it runs around what it wraps.
</details>

**55.** How does middleware reject a request?

- [ ] Returns an `error` value to the router
- [ ] Calls `panic`, which the server turns into a 401
- [ ] Calls `next.ServeHTTP` with an empty request
- [x] Writes an error and returns without calling `next.ServeHTTP`

<details><summary>Answer</summary>

Example: no cookie, so `RequireAuth` writes `401 UNAUTHENTICATED` and **returns without calling `next.ServeHTTP`**. The `MoveCard` handler never runs.

Forgetting that `return` would write the error and then run the handler anyway.
</details>

**56.** How does the auth middleware pass the logged-in user to the handler?

- [ ] Sets a global `currentUser` variable
- [ ] Adds it to the request body as JSON
- [x] Adds it to the request's context with `context.WithValue`
- [ ] Writes it to a response header

<details><summary>Answer</summary>

Example: `RequireAuth` looks up the session, finds Ada, and does `ctx := context.WithValue(r.Context(), userKey, ada)`. It passes on `r.WithContext(ctx)`.

The handler calls **`userFromContext(r.Context())`** and gets Ada without touching the database.
</details>

**57.** Why put the session check in middleware rather than in each handler?

- [ ] Middleware runs faster than handlers
- [x] Written once, can't be forgotten, handlers keep one job
- [ ] Handlers can't read cookies
- [ ] Go requires auth to be in middleware

<details><summary>Answer</summary>

Example: you add `DELETE /columns/{id}` and forget to check the session. **Anyone can delete columns.** With `protected(…)` on the route, the check is applied by the router.

Handlers stay about one thing: the request itself.
</details>

**57a.** Which design pattern is middleware a form of?

- [ ] Repository
- [ ] Singleton
- [x] Decorator
- [ ] Observer

<details><summary>Answer</summary>

**Decorator: wrap something to add behaviour without changing it.** `RequireAuth(next)` returns a new handler that checks the session, then calls `next`. The `Me` handler itself has no auth code at all.
</details>

### Logout

**58.** What are the two halves of logout, and which one actually logs the user out?

- [x] Delete the session row and clear the cookie; the row is what counts
- [ ] Delete the session row and clear the cookie; the cookie is what counts
- [ ] Clear the cookie and expire the JWT; the JWT is what counts
- [ ] Delete the user and the session; the user is what counts

<details><summary>Answer</summary>

Example: Ada logs out but only the cookie is cleared. Someone who copied `session=Xy7…` earlier **can still use it**, because the row is still in `sessions`.

Deleting the row makes the token useless everywhere, whoever has it.
</details>

**59.** Why is logout `POST`, not `GET`?

- [ ] `GET` requests can't carry cookies
- [ ] `POST` is faster than `GET`
- [ ] Browsers cache `POST` responses
- [x] A `GET` should never change anything; others can send them on their own

<details><summary>Answer</summary>

Example: if logout were `GET /logout`, a page on `evil.com` with `<img src="https://yourapp/logout">` would **log Ada out just by being viewed**. A link preview could do it too.

`POST` requests aren't sent by images or previews.
</details>

**60.** Why does logout return `204` even when there's no cookie or session?

- [ ] 204 is the only code logout can return
- [ ] Returning 401 would log the user back in
- [x] The result is the same either way, so it's idempotent
- [ ] The handler can't tell whether there was a cookie

<details><summary>Answer</summary>

Example: Ada clicks Log out in two tabs. The second request has no valid session, but **she's logged out either way**, so `204` is the honest answer.

An error here would only confuse the frontend.
</details>

### Ownership checks (403 vs 404)

**61.** What's the difference between authentication and authorization?

- [ ] Authentication: are you allowed? Authorization: who are you?
- [x] Authentication: who are you? Authorization: are you allowed?
- [ ] They mean the same thing
- [ ] Authentication: for APIs. Authorization: for browsers

<details><summary>Answer</summary>

Example: Bo's cookie proves he's Bo: that's **authentication**, failing gives `401`. Bo asks for Ada's board, which he's not allowed to see: that's **authorization**.

Authentication always comes first: you can't decide what someone's allowed until you know who they are.
</details>

**62.** Why does Bo get `404`, not `403`, for Ada's board?

- [x] `403` would confirm the board exists
- [ ] `403` only applies to admins
- [ ] `404` is faster to send
- [ ] Bo isn't logged in

<details><summary>Answer</summary>

Example: with `403`, Bo could try `/boards/1`, `/boards/2`, … and **learn which IDs exist** from the different answers.

With `404` for both "not yours" and "doesn't exist", he learns nothing.
</details>

**63.** When someone creates a board, where must the server get its `owner_id` (the user ID) from?

- [ ] The request body
- [ ] A query param
- [ ] The `boards` table
- [x] The session (the logged-in user)

<details><summary>Answer</summary>

Example: if the body could set `owner_id`, Bo could send `{"title":"x","owner_id":"ada"}` and **create a board in Ada's account**.

The session is the only thing the client can't fake, so `userFromContext` gives the owner.
</details>

**64.** Columns and cards have no `owner_id`. How are they checked?

- [ ] They're checked in Go after loading
- [ ] They aren't checked; anyone can edit them
- [x] Their queries join up to `boards` and check `owner_id`
- [ ] Each has a hidden `owner_id` column

<details><summary>Answer</summary>

Ownership is stored once, on the board. A card belongs to Ada if **its column's board is Ada's**.

Example: `GetCardByID` joins `cards → columns → boards` and adds `WHERE b.owner_id = $2`. Bo's ID matches no rows, so he gets `404`.
</details>

## 4. Transactions

### What a transaction is (ACID)

**65.** What does a transaction guarantee about the statements inside it?

- [ ] They run faster together
- [x] They succeed or fail as one unit
- [ ] They run one after another in a set order
- [ ] Other users can't read the tables until it ends

<details><summary>Answer</summary>

Example: `MoveCard` updates the card **and** inserts a history row. If the insert fails, the update is undone too.

You never end up with a card moved but no record of the move.
</details>

**66.** What do the four letters of ACID stand for?

- [x] Atomicity, Consistency, Isolation, Durability
- [ ] Atomicity, Concurrency, Integrity, Durability
- [ ] Accuracy, Consistency, Isolation, Delivery
- [ ] Atomicity, Commit, Isolation, Data

<details><summary>Answer</summary>

Each letter with `MoveCard` as the example:

- **Atomicity:** if the history insert fails, the card's move is undone too.
- **Consistency:** a foreign key violation fails the whole transaction.
- **Isolation:** Bo doesn't see the card moved until you commit.
- **Durability:** once `Commit` returns, a server crash can't lose the move.
</details>

**67.** When do you need an explicit transaction?

- [ ] For every `SELECT`
- [ ] Only when the table has a foreign key
- [ ] Whenever you use `r.Context()`
- [x] When two or more statements must stay in step

<details><summary>Answer</summary>

Example: `DELETE FROM boards WHERE id = $1` is one statement, so it's **already all or nothing**. `MoveCard`'s `UPDATE` and `INSERT` are two, so they need `BeginTx`.

Wrapping single statements in transactions adds work for no gain.
</details>

### `BeginTx` / `Commit` / `Rollback`

**68.** What does `BeginTx` do with the connection pool?

- [ ] Locks the whole pool until `Commit`
- [ ] Opens a new connection outside the pool
- [x] Takes one connection, sends `BEGIN`, and ties the `*sql.Tx` to it
- [ ] Sends `BEGIN` on every connection in the pool

<details><summary>Answer</summary>

Every statement in a transaction must run on **the same connection**, because that's where Postgres keeps the transaction.

Example: `MoveCard`'s `SELECT`, `UPDATE` and `INSERT` all go down the one connection that `BeginTx` took. Until `Commit`, no other request can use it.
</details>

**69.** What's the difference between `Commit` and `Rollback`?

- [ ] `Commit` throws changes away; `Rollback` keeps them
- [x] `Commit` keeps every change; `Rollback` throws them away
- [ ] `Commit` saves one statement; `Rollback` undoes the last one
- [ ] `Commit` starts it; `Rollback` ends it

<details><summary>Answer</summary>

Example: after the `UPDATE` and `INSERT` in `MoveCard`, **`Commit`** makes both visible to Bo at once. If either failed, **`Rollback`** puts the card back where it was, as if nothing happened.
</details>

**70.** What happens to an open transaction if the request's context is cancelled?

- [x] `database/sql` rolls it back automatically
- [ ] It's committed automatically
- [ ] It stays open until the server restarts
- [ ] Postgres keeps the changes made so far

<details><summary>Answer</summary>

Example: Ada closes the tab halfway through `MoveCard`. The context is cancelled, so `database/sql` **rolls back the transaction** and returns the connection to the pool.

Nothing is left half-done.
</details>

### `defer tx.Rollback()`

**71.** Why call `defer tx.Rollback()` straight after `BeginTx`?

- [ ] So the transaction commits faster
- [ ] Because `Commit` doesn't work without it
- [ ] So the rollback runs before the first query
- [x] So every early `return` rolls back

<details><summary>Answer</summary>

`MoveCard` has many `return`s. Without `defer`, you'd need `tx.Rollback()` before every error return, and **missing one leaves the transaction open**, holding its locks and connection.

`defer` runs on every return, so you can't forget.
</details>

**72.** Why is it safe for the deferred `Rollback` to run after a successful `Commit`?

- [ ] It undoes the commit
- [ ] It panics, but `defer` catches it
- [x] It does nothing; it returns `sql.ErrTxDone`, which is ignored
- [ ] It commits a second time

<details><summary>Answer</summary>

Example: `MoveCard` commits, then returns, so the deferred `Rollback` runs. **The transaction is already over**, so it just returns `sql.ErrTxDone` and changes nothing.

That's why one `defer` covers both the failure and the success paths.
</details>

### Every query on `tx`, not `db`

**73.** Inside a transaction, what goes wrong if one statement uses `s.db` instead of `tx`?

- [ ] It fails with a compile error
- [x] It runs outside the transaction, so `Rollback` doesn't undo it
- [ ] It joins the transaction automatically
- [ ] It runs twice

<details><summary>Answer</summary>

Example: in `MoveCard`, if the history `INSERT` used `s.db`, it would **commit straight away on another connection**. If the transaction then rolled back, the card wouldn't move, but the history row would still say it did.
</details>

**74.** Why doesn't the compiler catch this mistake?

- [x] `*sql.DB` and `*sql.Tx` have the same method names
- [ ] Go doesn't type-check SQL strings
- [ ] `s.db` and `tx` are the same type
- [ ] The compiler doesn't check `defer`

<details><summary>Answer</summary>

Both types have `ExecContext`, `QueryRowContext` and `QueryContext` with the same arguments, so **`s.db.ExecContext(…)` compiles just as well as `tx.ExecContext(…)`**.

The mistake only shows up when a rollback doesn't undo everything.
</details>

### `MoveCard` (`POST /cards/{id}/move`)

**75.** Why does moving a card have its own endpoint instead of using `PUT /cards/{id}`?

- [ ] `PUT` can't change `column_id`
- [ ] `PUT` requests can't have a body
- [ ] Moves need a different status code
- [x] A move does two writes that must be in one transaction

<details><summary>Answer</summary>

`PUT /cards/{id}` just updates the card's fields. A move also has to **write a `card_movements` row in the same transaction**.

A separate `POST /cards/{id}/move` makes that a named operation the store can do atomically.
</details>

**76.** Why is the card's current column read inside the transaction?

- [ ] Postgres requires a `SELECT` before every `UPDATE`
- [ ] To check the card's title hasn't changed
- [x] It's needed for `from_column_id` before the `UPDATE` overwrites it
- [ ] Reading outside a transaction isn't allowed

<details><summary>Answer</summary>

Example: the card is in `todo`, and you're moving it to `doing`. The history row needs `from_column_id = todo`, but **after the `UPDATE` the card says `doing`**.

Reading it in the same transaction (with `FOR UPDATE`) means nobody can move it in between.
</details>

**77.** What does "append-only" mean for `card_movements`?

- [ ] Rows are added to the end of the file on disk
- [x] The app only adds rows; it never changes them
- [ ] Rows can be updated but not deleted
- [ ] Rows are deleted after each move

<details><summary>Answer</summary>

Example: card `c1` goes `todo → doing → done`. That's **two rows, never edited**. To find where it's been, you read the rows in order.

A log you can edit can't be trusted as history.
</details>

**77a.** What does `MoveCard` do when the card is already in the destination column?

- [ ] Writes a history row from `todo` to `todo`
- [ ] Returns `409 Conflict`
- [x] Returns success and writes nothing
- [ ] Returns `404 COLUMN_NOT_FOUND`

<details><summary>Answer</summary>

Example: `c1` is in `doing`, and the client asks to move it to `doing`. Writing a `doing → doing` row would be **history of a move that didn't happen**.

So `MoveCard` returns `true` early, and the deferred `Rollback` ends the empty transaction.
</details>

### Mapping transaction failures to API errors

**78.** After a failed move, what state is the database in?

- [x] Unchanged; the `Rollback` undid any partial work
- [ ] Half-changed; the card moved but no history row
- [ ] Locked until someone commits
- [ ] The card is deleted

<details><summary>Answer</summary>

Example: the card's `UPDATE` worked, then the history `INSERT` failed. **The deferred `Rollback` undoes the `UPDATE`**, so the card is still in its old column, and there's no history row.

The client can safely try again.
</details>

**79.** Why does the store return `ErrColumnNotFound` instead of a status code?

- [ ] Status codes can't be returned from Go functions
- [ ] `errors.Is` only works with status codes
- [ ] The handler doesn't know the column ID
- [x] The store shouldn't know about HTTP

<details><summary>Answer</summary>

`404` is HTTP, and the store only talks to the database.

Example: `MoveCard` returns **`ErrColumnNotFound`**. The handler checks `errors.Is(err, store.ErrColumnNotFound)` and writes `404 COLUMN_NOT_FOUND`. A test fake can return the same error with no database.
</details>

## 5. Concurrency

### Race conditions

**80.** What is a race condition?

- [ ] Two requests that ask for the same URL
- [x] A bug where the result depends on which of two simultaneous things happens first
- [ ] A query that runs slower than the request's timeout
- [ ] Two goroutines that finish at exactly the same time

<details><summary>Answer</summary>

Example: Ada has two tabs open. Both move card `c1` at the same moment. Whichever request writes last decides where `c1` ends up, and the history may be wrong.

The code is correct when run once. It **only breaks when two runs overlap**, which is why race conditions are hard to spot in testing.
</details>

**81.** Why can two requests run `MoveCard` at the same moment?

- [ ] Postgres runs every query twice
- [ ] The frontend sends every request twice
- [ ] `MoveCard` starts its own goroutine
- [x] Go's server handles each request in its own goroutine

<details><summary>Answer</summary>

When Ada's two tabs each send `POST /cards/c1/move`, `net/http` **starts a goroutine for each request**. Both run `MoveCard` at the same time, so both can read `c1` before either writes.

You didn't write `go` anywhere; the server does it for you.
</details>

**82.** What is a lost update?

- [x] A write silently built on data that was already out of date
- [ ] An `UPDATE` that matches no rows
- [ ] A write that's rolled back by `Rollback`
- [ ] An update the client never receives a response for

<details><summary>Answer</summary>

Example: both tabs read `c1` as `todo`. Tab 1 moves it to `doing`. Tab 2 still thinks it's in `todo`, so it records **`todo → done` instead of `doing → done`**.

Nothing errors. The bad data is saved silently, so you only notice later.
</details>

### `SELECT ... FOR UPDATE`

**83.** What does `SELECT ... FOR UPDATE` do?

- [ ] Selects the rows and updates them in one statement
- [ ] Locks the whole table until the transaction ends
- [x] Locks the selected rows until the transaction ends, so others wait
- [ ] Marks the rows to be updated by the next query

<details><summary>Answer</summary>

Example: tab 1 runs `SELECT … FOR UPDATE` on `c1`. Tab 2 runs the same query and **waits** until tab 1 commits. Then tab 2 reads `c1` as `doing`, the up-to-date value, and records `doing → done`.

The second request is made to queue behind the first, which removes the race.
</details>

**84.** Why does `FOR UPDATE` only work inside a transaction?

- [ ] Postgres refuses `FOR UPDATE` outside one
- [ ] It needs `BEGIN` to know which table to lock
- [ ] It only works with `SERIALIZABLE`
- [x] The lock is held until `Commit` or `Rollback`, which a transaction provides

<details><summary>Answer</summary>

Without `BEGIN`, each statement is its own transaction, so the lock is **released as soon as the `SELECT` finishes**.

Then the `UPDATE` that follows is unprotected: another request can change the card in between, and you're back to the race.
</details>

**85.** With `FOR UPDATE` on card `c1`, does a move of card `c2` have to wait?

- [x] No; only the selected row (`c1`) is locked
- [ ] Yes; the whole `cards` table is locked
- [ ] Yes, if `c2` is in the same column
- [ ] No; `c2` is moved first

<details><summary>Answer</summary>

No. `FOR UPDATE` **locks only the rows the `SELECT` returned**, so moving `c1` doesn't block a move of `c2`.

That's why it's cheap: only requests for the same card wait for each other.
</details>

**85a.** Why does `MoveCard` lock the card with a query that has no join?

- [ ] Joins can't be used inside a transaction
- [ ] `FOR UPDATE` locks every table in a join
- [x] After waiting, Postgres re-reads only the locked row, so a join can go stale
- [ ] A join makes the lock last after `Commit`

<details><summary>Answer</summary>

Example: tab 2 waits while tab 1 moves `c1` from `todo` to `doing`. When tab 2 wakes, Postgres **re-reads the card (`doing`) but reuses the column it joined before (`todo`)**. `col.id = ca.column_id` no longer matches, so the query returns no rows and the move says "not found" (#8).

Locking with `SELECT column_id FROM cards WHERE id = $1 FOR UPDATE`, then checking ownership separately, leaves nothing to go stale.
</details>

### Isolation levels

**86.** What does an isolation level set?

- [ ] How many transactions can run at once
- [x] How much one transaction sees of others' changes while it runs
- [ ] Which users can read a table
- [ ] How long a lock is held

<details><summary>Answer</summary>

Example: your transaction reads `c1` as `todo`. Meanwhile someone else commits a move to `doing`. When you read `c1` again, do you see `todo` or `doing`?

**The isolation level decides.** Stricter levels hide others' commits until your transaction ends.
</details>

**87.** What is Postgres's default isolation level, and what does it see?

- [ ] `SERIALIZABLE`: a frozen snapshot, and clashes are cancelled
- [ ] `REPEATABLE READ`: a frozen snapshot from when the transaction started
- [x] `READ COMMITTED`: other people's changes as soon as they commit
- [ ] `READ UNCOMMITTED`: every change, committed or not

<details><summary>Answer</summary>

Example: your transaction reads `c1` as `todo`. Someone commits a move to `doing`. Your next read in the same transaction **sees `doing`**.

Each statement sees the latest committed data. You never see changes that haven't been committed yet.
</details>

**88.** What's the cost of a stricter level like `SERIALIZABLE`?

- [ ] Queries can't use indexes
- [ ] Only one connection can be open
- [x] Postgres may cancel a transaction (`40001`), and the app must retry it
- [ ] Committed changes might be lost

<details><summary>Answer</summary>

Under `SERIALIZABLE`, Postgres won't let two transactions produce a result that couldn't happen if they ran one after the other.

Example: both tabs read `c1` as `todo`. Tab 1 moves it to `doing` and commits. When tab 2 tries to move it, Postgres **cancels tab 2 with `40001`**. Your Go code must catch that error and run tab 2's transaction again from the start.

The cost is that extra retry code, plus the wasted work.
</details>

**88a.** `READ COMMITTED` gives the most up-to-date data. So what do stricter levels give you?

- [ ] Even newer data
- [ ] Faster queries
- [x] All the reads in one transaction agree with each other
- [ ] No need for transactions

<details><summary>Answer</summary>

Example: under `READ COMMITTED`, you count the cards in `todo` and get 3. Someone commits a move. You list the cards in `todo` and get 2. The two reads contradict each other.

Under `REPEATABLE READ`, both reads see 3: **every read in the transaction sees the same snapshot**.
</details>

### Optimistic vs pessimistic concurrency

**89.** What's the difference between pessimistic and optimistic concurrency?

- [ ] Pessimistic retries; optimistic waits
- [x] Pessimistic locks first so others wait; optimistic checks at write time and fails on a change
- [ ] Pessimistic is for reads; optimistic is for writes
- [ ] Pessimistic uses goroutines; optimistic uses channels

<details><summary>Answer</summary>

Pessimistic: `MoveCard` uses `FOR UPDATE`, so a second move **waits** for the first.

Optimistic: no lock. The `UPDATE` checks the `version` column hasn't changed since you read it, and **fails with `409`** if it has. The user then reloads and tries again.

Use pessimistic when clashes are common and quick; optimistic when they're rare or slow.
</details>

**90.** With a `version` column, how does the store know someone else changed the card first?

- [ ] Postgres returns error `23505`
- [ ] The `version` column becomes `NULL`
- [ ] The `SELECT` returns two rows
- [x] `UPDATE … WHERE version = $3` matches no rows, so `RowsAffected()` is 0

<details><summary>Answer</summary>

Example: Ada loads card `c1` at `version = 4`. Bo saves it first, so it's now `version = 5`. Ada's save runs:

`UPDATE cards SET title = $1, version = version + 1 WHERE id = $2 AND version = 4`

No row has `version = 4` any more, so **`RowsAffected()` is 0**. The store reports a conflict and the handler returns `409`.
</details>

**91.** Why do edit forms use optimistic, not pessimistic, concurrency?

- [x] A lock can't be held while the user spends minutes editing
- [ ] Forms can't send `FOR UPDATE`
- [ ] Optimistic is always faster
- [ ] Pessimistic locks don't work over HTTP

<details><summary>Answer</summary>

A transaction, and its lock, ends when the request ends. Ada opens the form in one request and saves it minutes later in another, so **no lock can last between the two**.

Even if it could, it would block Bo from touching the card while Ada typed. A version check costs nothing until the save.
</details>

### Goroutines

**92.** What does `go sendEmail(user)` do?

- [ ] Sends the email, then continues
- [ ] Schedules the email for when `main` returns
- [x] Runs `sendEmail` at the same time; the line returns straight away
- [ ] Runs `sendEmail` in a separate program

<details><summary>Answer</summary>

Example: in `CreateUser`, sending the welcome email takes 2 seconds. `go sendEmail(user)` runs it **in the background**, so the handler returns `201` straight away instead of making the user wait.

The catch: the handler can't use the email's result or error.
</details>

**93.** What happens to running goroutines when `main` returns?

- [ ] The program waits for them to finish
- [x] They all stop
- [ ] They keep running in the background
- [ ] They finish their current line, then stop

<details><summary>Answer</summary>

Example: `main` runs `go sendEmail(user)` and then returns. The email **is never sent**: the goroutine is stopped mid-way.

`sync.WaitGroup` makes `main` wait: `wg.Add(1)` before starting, `wg.Done()` when the goroutine finishes, `wg.Wait()` to block until all are done.
</details>

### Mutexes

**94.** What does a mutex do?

- [ ] Locks a database row
- [ ] Sends a value between goroutines
- [ ] Starts a goroutine
- [x] Lets only one goroutine at a time run a piece of code

<details><summary>Answer</summary>

Example: two goroutines each run `count++` on `count = 5`. Both read 5 and both write 6, so **one increment is lost**.

With `mu.Lock()` before and `mu.Unlock()` after, the second goroutine waits; the result is 7. It's the Go version of `FOR UPDATE`, but for memory instead of a row.
</details>

**95.** Why write `defer mu.Unlock()` straight after `mu.Lock()`?

- [x] So every return unlocks, like `defer tx.Rollback()`
- [ ] So the mutex unlocks before the code runs
- [ ] Because `Unlock` must be the last line of the file
- [ ] So other goroutines can run at the same time

<details><summary>Answer</summary>

Example: between `Lock()` and `Unlock()` there's an `if err != nil { return err }`. Without `defer`, that return **skips the unlock**, and every other goroutine waits forever.

`defer` runs on every return, so you can't forget.
</details>

**96.** What does `go test -race` do?

- [ ] Runs the tests in parallel to make them faster
- [ ] Times how long each test takes
- [x] Reports memory shared between goroutines without protection
- [ ] Runs each test twice

<details><summary>Answer</summary>

You run `go test -race ./...`. If two goroutines touch the same variable without a mutex or channel, the test fails with **`WARNING: DATA RACE`** and the two lines of code involved.

It's how you find races that would otherwise only show up now and then in production.
</details>

### Channels

**97.** What is a channel?

- [ ] A database connection
- [x] A typed pipe one goroutine sends values into and another receives from
- [ ] A list that only goroutines can read
- [ ] A lock shared between goroutines

<details><summary>Answer</summary>

Example: a goroutine computes a total and sends it with `results <- 5`. `main` waits at `answer := <-results` until it arrives.

The value **moves from one goroutine to the other**, so they never share a variable, and you don't need a mutex.
</details>

**98.** What's the difference between `make(chan T)` and `make(chan T, 10)`?

- [ ] The first holds one value; the second holds ten
- [ ] The first is for sending; the second is for receiving
- [ ] The first can be closed; the second can't
- [x] Unbuffered: the sender waits for a receiver. Buffered: holds 10 before the sender waits

<details><summary>Answer</summary>

Example: a handler sends jobs to a worker. With `make(chan string)`, each send **waits until the worker takes it**. With `make(chan string, 10)`, the handler can drop off 10 jobs and carry on; it only waits when all 10 slots are full.

Buffered channels are how step 9's job queue lets requests return quickly.
</details>
