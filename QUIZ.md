# Quiz

Recall questions on the concepts in `LEARNING.md`. Try answering before opening each answer.

## 1. CRUD

**1.** What does CRUD stand for, and which HTTP method goes with each?

- [x] Create `POST`, Read `GET`, Update `PUT`, Delete `DELETE`
- [ ] Create `PUT`, Read `GET`, Update `POST`, Delete `DELETE`
- [ ] Create `POST`, Read `GET`, Update `PATCH`, Delete `POST`
- [ ] Copy `GET`, Read `GET`, Update `PUT`, Drop `DELETE`

<details><summary>Answer</summary>

Create (`POST`), Read (`GET`), Update (`PUT`), Delete (`DELETE`).
</details>

**2.** What does *idempotent* mean, and which of `GET`, `POST`, `PUT`, `DELETE` is not idempotent?

- [ ] Repeating it has the same effect as once; `DELETE` isn't
- [ ] It never changes data; `PUT` isn't
- [ ] It always returns the same body; `GET` isn't
- [x] Repeating it has the same effect as once; `POST` isn't

<details><summary>Answer</summary>

Repeating the request has the same effect as doing it once. `POST` isn't: two creates can make two resources.
</details>

**3.** An `UPDATE` or `DELETE` matches no rows. Postgres doesn't return an error. How does the store know nothing was found?

- [ ] `Exec` returns `sql.ErrNoRows`
- [ ] `Exec` returns a `nil` result
- [x] `RowsAffected()` returns 0
- [ ] Postgres returns error code `23505`

<details><summary>Answer</summary>

`RowsAffected()` returns 0.
</details>

## 2. HTTP Fundamentals

### Status codes

**4.** What do the three status code classes 2xx, 4xx and 5xx mean?

- [ ] 2xx success, 4xx server's fault, 5xx client's fault
- [x] 2xx success, 4xx client's fault, 5xx server's fault
- [ ] 2xx created, 4xx redirect, 5xx timeout
- [ ] 2xx success, 4xx warning, 5xx client's fault

<details><summary>Answer</summary>

2xx: success. 4xx: the client's fault. 5xx: the server's fault.
</details>

**5.** What do 201, 204, 404 and 409 mean?

- [x] 201 Created, 204 No Content, 404 Not Found, 409 Conflict
- [ ] 201 OK, 204 Deleted, 404 Not Found, 409 Forbidden
- [ ] 201 Created, 204 Not Modified, 404 Gone, 409 Conflict
- [ ] 201 Accepted, 204 No Content, 404 Bad Request, 409 Duplicate

<details><summary>Answer</summary>

201 Created, 204 No Content, 404 Not Found, 409 Conflict.
</details>

**6.** How many times can a handler set the status code?

- [ ] As many times as it likes; the last one wins
- [ ] Twice: once for headers, once for the body
- [ ] Never; Go always sends 200
- [x] Once; later `WriteHeader` calls are ignored

<details><summary>Answer</summary>

Once. Later calls to `WriteHeader` are ignored.
</details>

### JSON encode/decode

**7.** What does decoding do, and what does encoding do?

- [ ] Decoding: Go struct → JSON. Encoding: JSON → Go struct
- [ ] Decoding decrypts the body; encoding encrypts it
- [x] Decoding: JSON → Go struct. Encoding: Go struct → JSON
- [ ] Decoding checks the JSON; encoding saves it to the database

<details><summary>Answer</summary>

Decoding turns JSON into a Go struct (for requests). Encoding turns a Go struct into JSON (for responses).
</details>

**8.** What does a `nil` slice encode as, and what does an empty slice encode as?

- [ ] Both encode as `[]`
- [x] `nil` → `null`; empty → `[]`
- [ ] Both encode as `null`
- [ ] `nil` → `[]`; empty → `null`

<details><summary>Answer</summary>

`nil` encodes as `null`; empty encodes as `[]`.
</details>

**9.** Which struct fields can `encoding/json` see?

- [x] Only exported fields (capital first letter)
- [ ] Only unexported fields (lowercase first letter)
- [ ] Every field, exported or not
- [ ] Only fields with a `json` struct tag

<details><summary>Answer</summary>

Only exported fields, the ones starting with a capital letter.
</details>

### Validation

**10.** What's the difference between decoding and validating?

- [x] Decoding: JSON is readable. Validating: values make sense
- [ ] Decoding: values make sense. Validating: JSON is readable
- [ ] They're the same step with different names
- [ ] Decoding checks types; validating checks the database

<details><summary>Answer</summary>

Decoding proves the JSON is readable; validation proves the values make sense (lengths, formats, required fields).
</details>

**11.** In what order does a create handler decode, validate and call the store?

- [ ] Validate → decode → store
- [ ] Store → decode → validate
- [ ] Decode → store → validate
- [x] Decode → validate → store

<details><summary>Answer</summary>

Decode → validate → store.
</details>

**12.** Validation can't catch a duplicate email. What does?

- [ ] A `validate:"unique"` tag, returned as 400
- [ ] The handler's `if` check, returned as 404
- [x] The database's `UNIQUE` constraint, returned as 409
- [ ] The `NOT NULL` constraint, returned as 500

<details><summary>Answer</summary>

The database's `UNIQUE` constraint. The API returns 409 Conflict.
</details>

### Path/query params

**13.** What does URL stand for, and what are the parts of `http://localhost:8080/boards/123?page=2`?

- [ ] Universal Request Link: protocol, server, route, params
- [x] Uniform Resource Locator: scheme, host and port, path, query
- [ ] Uniform Resource Locator: protocol, domain, folder, file
- [ ] Unique Resource Label: scheme, path, body, query

<details><summary>Answer</summary>

Uniform Resource Locator: it locates a resource. `http` is the scheme, `localhost:8080` the host and port, `/boards/123` the path, `?page=2` the query.
</details>

**14.** What's a path param for, and what's a query param for?

- [x] Path param: which resource. Query param: filter or shape it
- [ ] Path param: filter or shape it. Query param: which resource
- [ ] Path param: for `GET`. Query param: for `POST`
- [ ] Both identify the resource; query params are optional

<details><summary>Answer</summary>

A path param identifies *which* resource (`/boards/{id}`). A query param filters or shapes the result (`/cards?status=completed`).
</details>

**15.** What does a missing query param return in Go?

- [ ] `nil`
- [ ] An error
- [ ] The default value from the route
- [x] An empty string, `""`

<details><summary>Answer</summary>

An empty string, `""`.
</details>

**16.** What type do path and query params arrive as?

- [ ] Their real types, e.g. `int`
- [ ] `any`
- [x] Strings
- [ ] Bytes

<details><summary>Answer</summary>

Strings. Numbers need converting, e.g. with `strconv.Atoi`.
</details>

### Consistent errors

**17.** In an error response, what's `code` for, and what's `message` for?

- [ ] `code`: the HTTP status. `message`: the Go error text
- [x] `code`: for the client's code. `message`: for people
- [ ] `code`: for people. `message`: for the client's code
- [ ] `code`: for logging. `message`: for the database

<details><summary>Answer</summary>

`code` is for the client's code to branch on; `message` is for people to read.
</details>

**18.** Where does the HTTP status code go in a response, and where does the application error code go?

- [x] Status code: the status line. Error code: the JSON body
- [ ] Both go in the JSON body
- [ ] Status code: the JSON body. Error code: a header
- [ ] Both go in the response headers

<details><summary>Answer</summary>

The HTTP status code is in the status line, the first line of the response. The application error code is in the JSON body.
</details>

**19.** Why do 500 errors send a generic message instead of the real error?

- [ ] The real error is too long to send
- [ ] JSON can't encode Go errors
- [ ] The browser hides 500 messages anyway
- [x] It can leak database details, and the client can't act on it

<details><summary>Answer</summary>

The real error can leak database details to attackers, and the client can't act on it anyway. It goes to the logs instead.
</details>

**20.** Which layer picks the HTTP status code: the handler or the store?

- [ ] The store
- [ ] The router
- [x] The handler
- [ ] Postgres

<details><summary>Answer</summary>

The handler. The store only talks to the database and reports what happened.
</details>

**21.** Why does the store turn Postgres error `23505` into `ErrBoardExists`, instead of the handler checking for `23505`?

- [ ] Handlers can't import `lib/pq`
- [x] Only the store should know Postgres exists
- [ ] `23505` changes between Postgres versions
- [ ] It makes the query run faster

<details><summary>Answer</summary>

Only the store should know Postgres exists. Handlers depend on the `Store` interface, so a test fake or another database can return the same error.
</details>

### Content-Type

**22.** What does the `Content-Type` header tell the client?

- [x] The format of the response body, e.g. `application/json`
- [ ] How long the response body is
- [ ] Which client the response is for
- [ ] Whether the response was cached

<details><summary>Answer</summary>

What format the response body is in, e.g. `application/json`, `text/plain` or `text/html`.
</details>

**23.** In what order must a handler set headers, write the status and write the body?

- [ ] Status → headers → body
- [ ] Body → headers → status
- [ ] Any order; Go sorts them
- [x] Headers → status → body

<details><summary>Answer</summary>

Headers → status (`WriteHeader`) → body.
</details>

### Request lifecycle

**24.** What do TCP and TLS stand for?

- [ ] Transfer Connection Protocol; Transport Link Security
- [ ] Transmission Control Protocol; Total Line Security
- [x] Transmission Control Protocol; Transport Layer Security
- [ ] Transport Control Packet; Trusted Layer Service

<details><summary>Answer</summary>

TCP: Transmission Control Protocol. TLS: Transport Layer Security.
</details>

**25.** What are the three steps of the TCP handshake?

- [ ] ACK → SYN → SYN-ACK
- [x] SYN → SYN-ACK → ACK
- [ ] HELLO → OK → READY
- [ ] SYN → ACK → FIN

<details><summary>Answer</summary>

SYN (synchronize) → SYN-ACK (synchronize-acknowledge) → ACK (acknowledge).
</details>

**26.** What does TLS add on top of TCP?

- [x] Encryption, and server identity via a certificate
- [ ] Faster connections and compression
- [ ] Guaranteed delivery and ordering
- [ ] User login and sessions

<details><summary>Answer</summary>

Encryption, and server identity via a certificate. Together they make HTTPS.
</details>

**27.** Why does Go's server give each connection its own goroutine?

- [ ] Because each connection needs its own database
- [ ] So connections can share memory safely
- [ ] Because Go has no threads
- [x] So one slow client doesn't block the others

<details><summary>Answer</summary>

So requests are handled concurrently, and one slow client doesn't block the others.
</details>

**28.** Which part of the server decides which handler a request goes to?

- [ ] The middleware
- [ ] The store
- [x] The mux (router)
- [ ] The `http.Server` TLS layer

<details><summary>Answer</summary>

The mux (multiplexer), also called the router.
</details>

### `context.Context`

**29.** What two things cancel a request's context?

- [ ] The database closing, or the server starting
- [x] The client disconnecting, or the handler returning
- [ ] A 404, or a 500
- [ ] Calling `defer`, or a panic

<details><summary>Answer</summary>

The client disconnecting, or the handler returning.
</details>

**30.** Why pass `r.Context()` to database calls?

- [x] So the query is abandoned if the client disconnects
- [ ] So the query runs inside a transaction
- [ ] So the database knows which user is logged in
- [ ] So the query result is cached

<details><summary>Answer</summary>

So the database abandons the query if the client disconnects, freeing the connection.
</details>

**31.** By convention, where does a context parameter go, and what is it named?

- [ ] Last, named `ctx`
- [ ] First, named `context`
- [ ] Last, named `c`
- [x] First, named `ctx`

<details><summary>Answer</summary>

First in the parameter list, named `ctx`.
</details>

## 3. Authentication / Authorization

### Password hashing

**32.** Why must a password hash be slow?

- [ ] Slow hashes are harder to reverse
- [ ] It stops two users having the same hash
- [x] Each guess costs ~100ms, so mass guessing is impractical
- [ ] It hides the password from the database

<details><summary>Answer</summary>

Attackers guess passwords and hash each guess. A slow hash (about 100ms) makes guessing millions of passwords impractical, while one login barely notices.
</details>

**33.** What is a salt, and what attack does it stop?

- [ ] A secret key on the server; stops stolen databases
- [x] A random value per password; stops precomputed lookup tables
- [ ] A random value per password; stops slow guessing
- [ ] Extra characters users must add; stops weak passwords

<details><summary>Answer</summary>

A random value mixed into each password before hashing, so identical passwords get different hashes. It stops precomputed lookup tables (rainbow tables).
</details>

**34.** Why doesn't the salt need to be secret?

- [x] It's random, so lookup tables can't be built in advance
- [ ] It's encrypted with the server's key
- [ ] It's deleted after the hash is made
- [ ] Attackers can't read the `users` table

<details><summary>Answer</summary>

It's random, so it can't be guessed in advance to build a lookup table. Once stolen, it only lets an attacker guess one user's password at a time.
</details>

**35.** Can anything stop an attacker with a stolen database from guessing passwords? What makes guessing impractical?

- [ ] Yes; the salt stops guessing completely
- [ ] Yes; bcrypt hashes can't be compared
- [ ] No; nothing makes guessing any slower
- [x] No; salt, bcrypt's slowness and strong passwords make it impractical

<details><summary>Answer</summary>

No: they can always hash a guess with a user's salt and compare. The salt makes them guess each user separately, bcrypt's slowness makes each guess take about 100ms, and a strong password means far too many guesses are needed.
</details>

**36.** How does login check a password, if a hash can't be reversed?

- [ ] Reverses the stored hash and compares the passwords
- [ ] Asks bcrypt to decrypt the stored password
- [x] Hashes the typed password with the same salt and cost, then compares
- [ ] Compares the typed password with the salt

<details><summary>Answer</summary>

It hashes what the user typed with the same salt and cost (read from the stored hash) and compares the result with the stored hash.
</details>

**37.** Which struct tag keeps a field out of JSON responses?

- [ ] `json:"omit"`
- [x] `json:"-"`
- [ ] `json:"private"`
- [ ] `json:",omitempty"`

<details><summary>Answer</summary>

`json:"-"`.
</details>

### Sessions and cookies

**38.** HTTP is stateless. How does the server know who sent a request after login?

- [x] A cookie with a session token, looked up in `sessions`
- [ ] The browser's IP address
- [ ] The user ID in the URL
- [ ] The password sent with every request

<details><summary>Answer</summary>

Login creates a session and sets a cookie holding a random session token. The browser sends the cookie on every request, and the server looks the token up in the `sessions` table.
</details>

**39.** Why store a hash of the session token rather than the token itself?

- [ ] Hashes take less space in the table
- [ ] Tokens can't be stored in Postgres
- [ ] Hashing makes the lookup faster
- [x] A leaked token is usable as it is; a hash isn't

<details><summary>Answer</summary>

A token is usable as it is: anyone with it is logged in. If the table leaks, hashes can't be used.
</details>

**40.** Why is SHA-256 fine for session tokens, when passwords need bcrypt?

- [ ] SHA-256 is slower than bcrypt
- [ ] Tokens expire, so they don't need hashing
- [x] Tokens are long and random, so they can't be guessed
- [ ] bcrypt can't hash random bytes

<details><summary>Answer</summary>

Tokens are long and random, so they can't be guessed; slowness adds nothing and would slow every request. Passwords are guessable, so they need a slow hash.
</details>

**41.** What do the cookie flags `HttpOnly`, `Secure` and `SameSite` do?

- [ ] `HttpOnly`: HTTP only. `Secure`: encrypted. `SameSite`: same server only
- [x] `HttpOnly`: no JS. `Secure`: HTTPS only. `SameSite`: not on other sites' requests
- [ ] `HttpOnly`: no HTTPS. `Secure`: signed. `SameSite`: one tab only
- [ ] `HttpOnly`: no JS. `Secure`: encrypted value. `SameSite`: no cross-origin JS

<details><summary>Answer</summary>

`HttpOnly`: JavaScript can't read it. `Secure`: only sent over HTTPS. `SameSite`: not sent on requests started by other sites.
</details>

### Sessions vs JWTs

**42.** What's the main difference between a session token and a JWT?

- [x] Session: a random ID looked up. JWT: signed details, no lookup
- [ ] Session: signed details. JWT: a random ID looked up
- [ ] Session: stored in the browser. JWT: stored in the database
- [ ] They're the same; JWT is the newer name

<details><summary>Answer</summary>

A session token is a random ID the server looks up; a JWT contains the user's details, signed, so the server trusts it without a lookup.
</details>

**43.** How do you log a user out everywhere with sessions, and why is it hard with JWTs?

- [ ] Change their password; JWTs can't be changed
- [ ] Clear their cookies; JWTs aren't cookies
- [ ] Delete the user; JWTs are copied to other servers
- [x] Delete their session rows; a JWT has nothing stored to delete

<details><summary>Answer</summary>

Delete all their session rows. A JWT has nothing stored to delete, so it stays valid until it expires (unless you add a blocklist).
</details>

**44.** Is a JWT's payload secret?

- [ ] Yes; it's encrypted with the server's key
- [ ] Yes; only the browser can read it
- [x] No; it's signed, not encrypted, so anyone can read it
- [ ] No; it's sent as plain text without a signature

<details><summary>Answer</summary>

No. JWTs are usually signed, not encrypted: anyone can read the payload, but can't change it without breaking the signature.
</details>

**45.** How does the server check a JWT's signature?

- [ ] Decrypts the signature to get the payload
- [x] Recalculates it from header + payload with its secret, and compares
- [ ] Looks the token up in the `sessions` table
- [ ] Checks the signature against the user's password hash

<details><summary>Answer</summary>

It recalculates the signature from the token's header and payload using its secret key, and compares it with the token's signature. A mismatch means the token was changed (or not made by this server).
</details>

**46.** Why can't an attacker change a JWT's payload and make a new valid signature?

- [x] A valid signature needs the server's secret key
- [ ] The payload is encrypted
- [ ] The browser blocks changes to cookies
- [ ] The signature includes the user's password

<details><summary>Answer</summary>

Making the signature needs the server's secret key, which never leaves the server.
</details>

**47.** Does a JWT's payload change from one request to the next? Who makes a new one?

- [ ] Yes; the client re-signs it on every request
- [ ] Yes; the payload holds the current URL
- [ ] No; and it can never be replaced
- [x] No; the same token each time. The server issues a new one

<details><summary>Answer</summary>

No. The client sends the same token, created at login, on every request; only the URL and body change. If the payload needs to change, the server issues a new token with a new signature.
</details>

**48.** What are access and refresh tokens, and what problem do they solve?

- [ ] Two copies of the same JWT; fixes lost tokens
- [ ] A JWT + a password; fixes weak passwords
- [x] Short-lived JWT + long-lived stored token; fixes revocation
- [ ] A session + a cookie; fixes CSRF

<details><summary>Answer</summary>

An access token is a short-lived JWT used on every request; a refresh token is long-lived, stored in the database, and swapped for new access tokens. They fix JWT revocation: delete the refresh token and the user is logged out once the access token expires.
</details>

**49.** Why don't sessions need refresh tokens?

- [ ] Sessions never expire
- [x] Every request checks `sessions`, so deleting the row logs out now
- [ ] Sessions are stored in the browser
- [ ] Refresh tokens only work with cookies

<details><summary>Answer</summary>

Every request already checks the `sessions` table, so deleting the row logs the user out immediately.
</details>

### User enumeration

**50.** What is user enumeration?

- [x] Using responses to find out which emails have accounts
- [ ] Listing every user with `GET /users`
- [ ] Numbering users with sequential IDs
- [ ] Guessing a user's password many times

<details><summary>Answer</summary>

Using an app's responses (messages or timing) to find out which emails have accounts.
</details>

**51.** What two things does login keep identical for an unknown email and a wrong password?

- [ ] The cookie, and the session row
- [ ] The email, and the password hash
- [ ] The status code, and the request body
- [x] The error message (and status), and the response time

<details><summary>Answer</summary>

The error message (and status), and the response time.
</details>

### Database constraints

**52.** What happens if you delete a user who still has sessions, when `sessions.user_id` is `REFERENCES users(id)` without `ON DELETE CASCADE`?

- [ ] The sessions are deleted too
- [ ] The user is deleted and sessions keep a dangling `user_id`
- [x] Postgres refuses it with a foreign key violation (`23503`)
- [ ] Postgres sets `sessions.user_id` to `NULL`

<details><summary>Answer</summary>

Postgres refuses the delete with a foreign key violation (`23503`).
</details>

**53.** What does `ON DELETE CASCADE` do?

- [ ] Deleting a child row deletes its parent
- [x] Deleting the parent row deletes the rows that reference it
- [ ] It stops the parent row being deleted
- [ ] It deletes the whole table

<details><summary>Answer</summary>

Deleting the referenced row (the user) also deletes the rows that reference it (their sessions).
</details>

### Auth middleware

**54.** What is middleware?

- [x] A function that wraps a handler and runs on every request to it
- [ ] A layer between the store and Postgres
- [ ] A goroutine that runs in the background
- [ ] A struct that holds the database pool

<details><summary>Answer</summary>

A function that wraps a handler and runs on every request to it, before (and after) the handler. In Go: `func(next http.Handler) http.Handler`.
</details>

**55.** How does middleware reject a request?

- [ ] Returns an `error` value to the router
- [ ] Calls `panic`, which the server turns into a 401
- [ ] Calls `next.ServeHTTP` with an empty request
- [x] Writes an error and returns without calling `next.ServeHTTP`

<details><summary>Answer</summary>

It writes the error response and returns without calling `next.ServeHTTP`, so the handler never runs.
</details>

**56.** How does the auth middleware pass the logged-in user to the handler?

- [ ] Sets a global `currentUser` variable
- [ ] Adds it to the request body as JSON
- [x] Adds it to the request's context with `context.WithValue`
- [ ] Writes it to a response header

<details><summary>Answer</summary>

It adds the user to the request's context with `context.WithValue`, and passes on `r.WithContext(ctx)`. The handler reads it back from `r.Context()`.
</details>

**57.** Why put the session check in middleware rather than in each handler?

- [ ] Middleware runs faster than handlers
- [x] Written once, can't be forgotten, handlers keep one job
- [ ] Handlers can't read cookies
- [ ] Go requires auth to be in middleware

<details><summary>Answer</summary>

It's written once, it can't be forgotten on a new route, and handlers keep their one job.
</details>

### Logout

**58.** What are the two halves of logout, and which one actually logs the user out?

- [x] Delete the session row and clear the cookie; the row is what counts
- [ ] Delete the session row and clear the cookie; the cookie is what counts
- [ ] Clear the cookie and expire the JWT; the JWT is what counts
- [ ] Delete the user and the session; the user is what counts

<details><summary>Answer</summary>

Deleting the session row (server) and clearing the cookie (browser). Deleting the row is what logs them out: the token then matches nothing.
</details>

**59.** Why is logout `POST`, not `GET`?

- [ ] `GET` requests can't carry cookies
- [ ] `POST` is faster than `GET`
- [ ] Browsers cache `POST` responses
- [x] A `GET` should never change anything; others can send them on their own

<details><summary>Answer</summary>

A `GET` should never change anything. Browsers, link previews and other sites can send `GET` requests on their own.
</details>

**60.** Why does logout return `204` even when there's no cookie or session?

- [ ] 204 is the only code logout can return
- [ ] Returning 401 would log the user back in
- [x] The result is the same either way, so it's idempotent
- [ ] The handler can't tell whether there was a cookie

<details><summary>Answer</summary>

The result is the same either way: logged out. That makes it idempotent.
</details>

### Ownership checks (403 vs 404)

**61.** What's the difference between authentication and authorization?

- [ ] Authentication: are you allowed? Authorization: who are you?
- [x] Authentication: who are you? Authorization: are you allowed?
- [ ] They mean the same thing
- [ ] Authentication: for APIs. Authorization: for browsers

<details><summary>Answer</summary>

Authentication: who are you? Authorization: are you allowed to do this?
</details>

**62.** Why does Bo get `404`, not `403`, for Ada's board?

- [x] `403` would confirm the board exists
- [ ] `403` only applies to admins
- [ ] `404` is faster to send
- [ ] Bo isn't logged in

<details><summary>Answer</summary>

`403` would confirm the board exists, which leaks information. `404` makes "not yours" look the same as "doesn't exist".
</details>

**63.** When someone creates a board, where must the server get its `owner_id` (the user ID) from?

- [ ] The request body
- [ ] A query param
- [ ] The `boards` table
- [x] The session (the logged-in user)

<details><summary>Answer</summary>

The session (the logged-in user), never the request body.
</details>

**64.** Columns and cards have no `owner_id`. How are they checked?

- [ ] They're checked in Go after loading
- [ ] They aren't checked; anyone can edit them
- [x] Their queries join up to `boards` and check `owner_id`
- [ ] Each has a hidden `owner_id` column

<details><summary>Answer</summary>

Their queries join up to the board (`cards → columns → boards`) and check `boards.owner_id`.
</details>

## 4. Transactions

### What a transaction is (ACID)

**65.** What does a transaction guarantee about the statements inside it?

- [ ] They run faster together
- [x] They succeed or fail as one unit
- [ ] They run one after another in a set order
- [ ] Other users can't read the tables until it ends

<details><summary>Answer</summary>

They succeed or fail as one unit: either every change is kept, or none is.
</details>

**66.** What do the four letters of ACID stand for?

- [x] Atomicity, Consistency, Isolation, Durability
- [ ] Atomicity, Concurrency, Integrity, Durability
- [ ] Accuracy, Consistency, Isolation, Delivery
- [ ] Atomicity, Commit, Isolation, Data

<details><summary>Answer</summary>

Atomicity (all or nothing), Consistency (constraints still hold), Isolation (others don't see half-finished changes), Durability (committed changes survive a crash).
</details>

**67.** When do you need an explicit transaction?

- [ ] For every `SELECT`
- [ ] Only when the table has a foreign key
- [ ] Whenever you use `r.Context()`
- [x] When two or more statements must stay in step

<details><summary>Answer</summary>

When two or more statements must stay in step. A single statement is already its own transaction.
</details>

### `BeginTx` / `Commit` / `Rollback`

**68.** What does `BeginTx` do with the connection pool?

- [ ] Locks the whole pool until `Commit`
- [ ] Opens a new connection outside the pool
- [x] Takes one connection, sends `BEGIN`, and ties the `*sql.Tx` to it
- [ ] Sends `BEGIN` on every connection in the pool

<details><summary>Answer</summary>

It takes one connection from the pool, sends `BEGIN` on it, and ties the returned `*sql.Tx` to that connection until `Commit` or `Rollback`.
</details>

**69.** What's the difference between `Commit` and `Rollback`?

- [ ] `Commit` throws changes away; `Rollback` keeps them
- [x] `Commit` keeps every change; `Rollback` throws them away
- [ ] `Commit` saves one statement; `Rollback` undoes the last one
- [ ] `Commit` starts it; `Rollback` ends it

<details><summary>Answer</summary>

`Commit` makes every change permanent and visible to others. `Rollback` throws every change since `BEGIN` away.
</details>

**70.** What happens to an open transaction if the request's context is cancelled?

- [x] `database/sql` rolls it back automatically
- [ ] It's committed automatically
- [ ] It stays open until the server restarts
- [ ] Postgres keeps the changes made so far

<details><summary>Answer</summary>

`database/sql` rolls it back automatically.
</details>

### `defer tx.Rollback()`

**71.** Why call `defer tx.Rollback()` straight after `BeginTx`?

- [ ] So the transaction commits faster
- [ ] Because `Commit` doesn't work without it
- [ ] So the rollback runs before the first query
- [x] So every early `return` rolls back

<details><summary>Answer</summary>

So every early `return` rolls back, without writing `Rollback` before each one. A forgotten rollback leaves the connection checked out and locks held.
</details>

**72.** Why is it safe for the deferred `Rollback` to run after a successful `Commit`?

- [ ] It undoes the commit
- [ ] It panics, but `defer` catches it
- [x] It does nothing; it returns `sql.ErrTxDone`, which is ignored
- [ ] It commits a second time

<details><summary>Answer</summary>

It does nothing and returns `sql.ErrTxDone`, which is ignored.
</details>

### Every query on `tx`, not `db`

**73.** Inside a transaction, what goes wrong if one statement uses `s.db` instead of `tx`?

- [ ] It fails with a compile error
- [x] It runs outside the transaction, so `Rollback` doesn't undo it
- [ ] It joins the transaction automatically
- [ ] It runs twice

<details><summary>Answer</summary>

It runs on a different connection, outside the transaction: it commits on its own, isn't undone by `Rollback`, and can't see the transaction's uncommitted changes.
</details>

**74.** Why doesn't the compiler catch this mistake?

- [x] `*sql.DB` and `*sql.Tx` have the same method names
- [ ] Go doesn't type-check SQL strings
- [ ] `s.db` and `tx` are the same type
- [ ] The compiler doesn't check `defer`

<details><summary>Answer</summary>

`*sql.DB` and `*sql.Tx` have the same method names (`ExecContext`, `QueryRowContext`).
</details>

### `MoveCard` (`POST /cards/{id}/move`)

**75.** Why does moving a card have its own endpoint instead of using `PUT /cards/{id}`?

- [ ] `PUT` can't change `column_id`
- [ ] `PUT` requests can't have a body
- [ ] Moves need a different status code
- [x] A move does two writes that must be in one transaction

<details><summary>Answer</summary>

A move does two writes: it changes the card's column and records a row in `card_movements`. Both must happen in one transaction.
</details>

**76.** Why is the card's current column read inside the transaction?

- [ ] Postgres requires a `SELECT` before every `UPDATE`
- [ ] To check the card's title hasn't changed
- [x] It's needed for `from_column_id` before the `UPDATE` overwrites it
- [ ] Reading outside a transaction isn't allowed

<details><summary>Answer</summary>

The history row needs `from_column_id`, and it must be read before the `UPDATE` overwrites it, as part of the same unit of work.
</details>

**77.** What does "append-only" mean for `card_movements`?

- [ ] Rows are added to the end of the file on disk
- [x] The app only adds rows; it never changes them
- [ ] Rows can be updated but not deleted
- [ ] Rows are deleted after each move

<details><summary>Answer</summary>

The app only ever adds rows; it never updates them. It's a history log. (Rows are deleted only when their card is, through `ON DELETE CASCADE`.)
</details>

**77a.** What does `MoveCard` do when the card is already in the destination column?

- [ ] Writes a history row from `todo` to `todo`
- [ ] Returns `409 Conflict`
- [x] Returns success and writes nothing
- [ ] Returns `404 COLUMN_NOT_FOUND`

<details><summary>Answer</summary>

Returns success (`204`) and writes nothing: the card is already where the client asked. The deferred `Rollback` ends the transaction.
</details>

### Mapping transaction failures to API errors

**78.** After a failed move, what state is the database in?

- [x] Unchanged; the `Rollback` undid any partial work
- [ ] Half-changed; the card moved but no history row
- [ ] Locked until someone commits
- [ ] The card is deleted

<details><summary>Answer</summary>

Unchanged. The deferred `Rollback` undid any partial work.
</details>

**79.** Why does the store return `ErrColumnNotFound` instead of a status code?

- [ ] Status codes can't be returned from Go functions
- [ ] `errors.Is` only works with status codes
- [ ] The handler doesn't know the column ID
- [x] The store shouldn't know about HTTP

<details><summary>Answer</summary>

The store shouldn't know about HTTP. It reports what went wrong with a sentinel error; the handler checks it with `errors.Is` and picks `404`.
</details>

## 5. Concurrency

### Race conditions

**80.** What is a race condition?

- [ ] Two requests that ask for the same URL
- [x] A bug where the result depends on which of two simultaneous things happens first
- [ ] A query that runs slower than the request's timeout
- [ ] Two goroutines that finish at exactly the same time

<details><summary>Answer</summary>

A bug where the result depends on **which of two things running at the same time happens first**.
</details>

**81.** Why can two requests run `MoveCard` at the same moment?

- [ ] Postgres runs every query twice
- [ ] The frontend sends every request twice
- [ ] `MoveCard` starts its own goroutine
- [x] Go's server handles each request in its own goroutine

<details><summary>Answer</summary>

Go's server handles **each request in its own goroutine**, so requests run concurrently.
</details>

**82.** What is a lost update?

- [x] A write silently built on data that was already out of date
- [ ] An `UPDATE` that matches no rows
- [ ] A write that's rolled back by `Rollback`
- [ ] An update the client never receives a response for

<details><summary>Answer</summary>

A write **built on stale data**, so another change is silently overwritten or recorded wrongly. E.g. both tabs read "To Do", so the second move records To Do → Done instead of In Progress → Done.
</details>

### `SELECT ... FOR UPDATE`

**83.** What does `SELECT ... FOR UPDATE` do?

- [ ] Selects the rows and updates them in one statement
- [ ] Locks the whole table until the transaction ends
- [x] Locks the selected rows until the transaction ends, so others wait
- [ ] Marks the rows to be updated by the next query

<details><summary>Answer</summary>

**Locks the rows it returns until the transaction ends.** Another transaction trying to lock the same rows waits.
</details>

**84.** Why does `FOR UPDATE` only work inside a transaction?

- [ ] Postgres refuses `FOR UPDATE` outside one
- [ ] It needs `BEGIN` to know which table to lock
- [ ] It only works with `SERIALIZABLE`
- [x] The lock is held until `Commit` or `Rollback`, which a transaction provides

<details><summary>Answer</summary>

The lock lasts **until the transaction ends**. Without one, the statement is its own transaction and the lock is released straight away.
</details>

**85.** With `FOR UPDATE` on card `c1`, does a move of card `c2` have to wait?

- [x] No; only the selected row (`c1`) is locked
- [ ] Yes; the whole `cards` table is locked
- [ ] Yes, if `c2` is in the same column
- [ ] No; `c2` is moved first

<details><summary>Answer</summary>

No. **Only the rows selected are locked**, so other cards aren't affected.
</details>

### Isolation levels

**86.** What does an isolation level set?

- [ ] How many transactions can run at once
- [x] How much one transaction sees of others' changes while it runs
- [ ] Which users can read a table
- [ ] How long a lock is held

<details><summary>Answer</summary>

**How much one transaction can see of other transactions' changes** while it runs.
</details>

**87.** What is Postgres's default isolation level, and what does each statement see?

- [ ] `SERIALIZABLE`: as if transactions ran one at a time
- [ ] `REPEATABLE READ`: data as it was when the transaction started
- [x] `READ COMMITTED`: everything committed before that statement started
- [ ] `READ UNCOMMITTED`: every change, committed or not

<details><summary>Answer</summary>

**`READ COMMITTED`**: each statement sees everything committed before that statement started.
</details>

**88.** What's the cost of a stricter level like `SERIALIZABLE`?

- [ ] Queries can't use indexes
- [ ] Only one connection can be open
- [x] Postgres may cancel a transaction (`40001`), and the app must retry it
- [ ] Committed changes might be lost

<details><summary>Answer</summary>

Postgres may **cancel a transaction with a serialization error (`40001`)**, so the app must retry it.
</details>

### Optimistic vs pessimistic concurrency

**89.** What's the difference between pessimistic and optimistic concurrency?

- [ ] Pessimistic retries; optimistic waits
- [x] Pessimistic locks first so others wait; optimistic checks at write time and fails on a change
- [ ] Pessimistic is for reads; optimistic is for writes
- [ ] Pessimistic uses goroutines; optimistic uses channels

<details><summary>Answer</summary>

**Pessimistic: lock first, others wait** (`FOR UPDATE`). **Optimistic: no lock; check at write time** that nothing changed, and fail if it did.
</details>

**90.** With a `version` column, how does the store know someone else changed the card first?

- [ ] Postgres returns error `23505`
- [ ] The `version` column becomes `NULL`
- [ ] The `SELECT` returns two rows
- [x] `UPDATE … WHERE version = $3` matches no rows, so `RowsAffected()` is 0

<details><summary>Answer</summary>

The `UPDATE` includes `WHERE version = <the version read>`. If it changed, **no rows match and `RowsAffected()` is 0**. The API returns `409 Conflict`.
</details>

**91.** Why do edit forms use optimistic, not pessimistic, concurrency?

- [x] A lock can't be held while the user spends minutes editing
- [ ] Forms can't send `FOR UPDATE`
- [ ] Optimistic is always faster
- [ ] Pessimistic locks don't work over HTTP

<details><summary>Answer</summary>

**A lock can't be held for the minutes a user spends editing.** Optimistic needs no lock; it just checks the version on save.
</details>

### Goroutines

**92.** What does `go sendEmail(user)` do?

- [ ] Sends the email, then continues
- [ ] Schedules the email for when `main` returns
- [x] Runs `sendEmail` at the same time; the line returns straight away
- [ ] Runs `sendEmail` in a separate program

<details><summary>Answer</summary>

Starts `sendEmail` **in a new goroutine**, running at the same time. The line returns straight away without waiting.
</details>

**93.** What happens to running goroutines when `main` returns?

- [ ] The program waits for them to finish
- [x] They all stop
- [ ] They keep running in the background
- [ ] They finish their current line, then stop

<details><summary>Answer</summary>

**They all stop.** The program doesn't wait for goroutines; use `sync.WaitGroup` to wait.
</details>

### Mutexes

**94.** What does a mutex do?

- [ ] Locks a database row
- [ ] Sends a value between goroutines
- [ ] Starts a goroutine
- [x] Lets only one goroutine at a time run a piece of code

<details><summary>Answer</summary>

Lets **only one goroutine at a time** run the code between `Lock()` and `Unlock()`.
</details>

**95.** Why write `defer mu.Unlock()` straight after `mu.Lock()`?

- [x] So every return unlocks, like `defer tx.Rollback()`
- [ ] So the mutex unlocks before the code runs
- [ ] Because `Unlock` must be the last line of the file
- [ ] So other goroutines can run at the same time

<details><summary>Answer</summary>

So **every return unlocks it**: the same pattern as `defer tx.Rollback()`. A forgotten unlock makes every other goroutine wait forever.
</details>

**96.** What does `go test -race` do?

- [ ] Runs the tests in parallel to make them faster
- [ ] Times how long each test takes
- [x] Reports memory shared between goroutines without protection
- [ ] Runs each test twice

<details><summary>Answer</summary>

Runs the **race detector**, which reports memory used by several goroutines without protection.
</details>

### Channels

**97.** What is a channel?

- [ ] A database connection
- [x] A typed pipe one goroutine sends values into and another receives from
- [ ] A list that only goroutines can read
- [ ] A lock shared between goroutines

<details><summary>Answer</summary>

A **typed pipe between goroutines**: one sends with `ch <- v`, another receives with `v := <-ch`.
</details>

**98.** What's the difference between `make(chan T)` and `make(chan T, 10)`?

- [ ] The first holds one value; the second holds ten
- [ ] The first is for sending; the second is for receiving
- [ ] The first can be closed; the second can't
- [x] Unbuffered: the sender waits for a receiver. Buffered: holds 10 before the sender waits

<details><summary>Answer</summary>

**Unbuffered** (`make(chan T)`): a send waits until someone receives. **Buffered** (`make(chan T, 10)`): holds up to 10 values before a send waits.
</details>
