# Quiz

Recall questions on the concepts in `LEARNING.md`. Try answering before opening each answer.

## 1. CRUD

**1.** What does CRUD stand for, and which HTTP method goes with each?

<details><summary>Answer</summary>

Create (`POST`), Read (`GET`), Update (`PUT`), Delete (`DELETE`).
</details>

**2.** What does *idempotent* mean, and which of `GET`, `POST`, `PUT`, `DELETE` is not idempotent?

<details><summary>Answer</summary>

Repeating the request has the same effect as doing it once. `POST` isn't: two creates can make two resources.
</details>

**3.** An `UPDATE` or `DELETE` matches no rows. Postgres doesn't return an error. How does the store know nothing was found?

<details><summary>Answer</summary>

`RowsAffected()` returns 0.
</details>

## 2. HTTP Fundamentals

### Status codes

**4.** What do the three status code classes 2xx, 4xx and 5xx mean?

<details><summary>Answer</summary>

2xx: success. 4xx: the client's fault. 5xx: the server's fault.
</details>

**5.** What do 201, 204, 404 and 409 mean?

<details><summary>Answer</summary>

201 Created, 204 No Content, 404 Not Found, 409 Conflict.
</details>

**6.** How many times can a handler set the status code?

<details><summary>Answer</summary>

Once. Later calls to `WriteHeader` are ignored.
</details>

### JSON encode/decode

**7.** What does decoding do, and what does encoding do?

<details><summary>Answer</summary>

Decoding turns JSON into a Go struct (for requests). Encoding turns a Go struct into JSON (for responses).
</details>

**8.** What does a `nil` slice encode as, and what does an empty slice encode as?

<details><summary>Answer</summary>

`nil` encodes as `null`; empty encodes as `[]`.
</details>

**9.** Which struct fields can `encoding/json` see?

<details><summary>Answer</summary>

Only exported fields, the ones starting with a capital letter.
</details>

### Validation

**10.** What's the difference between decoding and validating?

<details><summary>Answer</summary>

Decoding proves the JSON is readable; validation proves the values make sense (lengths, formats, required fields).
</details>

**11.** In what order does a create handler decode, validate and call the store?

<details><summary>Answer</summary>

Decode → validate → store.
</details>

**12.** Validation can't catch a duplicate email. What does?

<details><summary>Answer</summary>

The database's `UNIQUE` constraint. The API returns 409 Conflict.
</details>

### Path/query params

**13.** What does URL stand for, and what are the parts of `http://localhost:8080/boards/123?page=2`?

<details><summary>Answer</summary>

Uniform Resource Locator: it locates a resource. `http` is the scheme, `localhost:8080` the host and port, `/boards/123` the path, `?page=2` the query.
</details>

**14.** What's a path param for, and what's a query param for?

<details><summary>Answer</summary>

A path param identifies *which* resource (`/boards/{id}`). A query param filters or shapes the result (`/cards?status=completed`).
</details>

**15.** What does a missing query param return in Go?

<details><summary>Answer</summary>

An empty string, `""`.
</details>

**16.** What type do path and query params arrive as?

<details><summary>Answer</summary>

Strings. Numbers need converting, e.g. with `strconv.Atoi`.
</details>

### Consistent errors

**17.** In an error response, what's `code` for, and what's `message` for?

<details><summary>Answer</summary>

`code` is for the client's code to branch on; `message` is for people to read.
</details>

**18.** Where does the HTTP status code go in a response, and where does the application error code go?

<details><summary>Answer</summary>

The HTTP status code is in the status line, the first line of the response. The application error code is in the JSON body.
</details>

**19.** Why do 500 errors send a generic message instead of the real error?

<details><summary>Answer</summary>

The real error can leak database details to attackers, and the client can't act on it anyway. It goes to the logs instead.
</details>

**20.** Which layer picks the HTTP status code: the handler or the store?

<details><summary>Answer</summary>

The handler. The store only talks to the database and reports what happened.
</details>

**21.** Why does the store turn Postgres error `23505` into `ErrBoardExists`, instead of the handler checking for `23505`?

<details><summary>Answer</summary>

Only the store should know Postgres exists. Handlers depend on the `Store` interface, so a test fake or another database can return the same error.
</details>

### Content-Type

**22.** What does the `Content-Type` header tell the client?

<details><summary>Answer</summary>

What format the response body is in, e.g. `application/json`, `text/plain` or `text/html`.
</details>

**23.** In what order must a handler set headers, write the status and write the body?

<details><summary>Answer</summary>

Headers → status (`WriteHeader`) → body.
</details>

### Request lifecycle

**24.** What do TCP and TLS stand for?

<details><summary>Answer</summary>

TCP: Transmission Control Protocol. TLS: Transport Layer Security.
</details>

**25.** What are the three steps of the TCP handshake?

<details><summary>Answer</summary>

SYN (synchronize) → SYN-ACK (synchronize-acknowledge) → ACK (acknowledge).
</details>

**26.** What does TLS add on top of TCP?

<details><summary>Answer</summary>

Encryption, and server identity via a certificate. Together they make HTTPS.
</details>

**27.** Why does Go's server give each connection its own goroutine?

<details><summary>Answer</summary>

So requests are handled concurrently, and one slow client doesn't block the others.
</details>

**28.** Which part of the server decides which handler a request goes to?

<details><summary>Answer</summary>

The mux (multiplexer), also called the router.
</details>

### `context.Context`

**29.** What two things cancel a request's context?

<details><summary>Answer</summary>

The client disconnecting, or the handler returning.
</details>

**30.** Why pass `r.Context()` to database calls?

<details><summary>Answer</summary>

So the database abandons the query if the client disconnects, freeing the connection.
</details>

**31.** By convention, where does a context parameter go, and what is it named?

<details><summary>Answer</summary>

First in the parameter list, named `ctx`.
</details>

## 3. Authentication / Authorization

### Password hashing

**32.** Why must a password hash be slow?

<details><summary>Answer</summary>

Attackers guess passwords and hash each guess. A slow hash (about 100ms) makes guessing millions of passwords impractical, while one login barely notices.
</details>

**33.** What is a salt, and what attack does it stop?

<details><summary>Answer</summary>

A random value mixed into each password before hashing, so identical passwords get different hashes. It stops precomputed lookup tables (rainbow tables).
</details>

**34.** Why doesn't the salt need to be secret?

<details><summary>Answer</summary>

It's random, so it can't be guessed in advance to build a lookup table. Once stolen, it only lets an attacker guess one user's password at a time.
</details>

**35.** Can anything stop an attacker with a stolen database from guessing passwords? What makes guessing impractical?

<details><summary>Answer</summary>

No: they can always hash a guess with a user's salt and compare. The salt makes them guess each user separately, bcrypt's slowness makes each guess take about 100ms, and a strong password means far too many guesses are needed.
</details>

**36.** How does login check a password, if a hash can't be reversed?

<details><summary>Answer</summary>

It hashes what the user typed with the same salt and cost (read from the stored hash) and compares the result with the stored hash.
</details>

**37.** Which struct tag keeps a field out of JSON responses?

<details><summary>Answer</summary>

`json:"-"`.
</details>

### Sessions and cookies

**38.** HTTP is stateless. How does the server know who sent a request after login?

<details><summary>Answer</summary>

Login creates a session and sets a cookie holding a random session token. The browser sends the cookie on every request, and the server looks the token up in the `sessions` table.
</details>

**39.** Why store a hash of the session token rather than the token itself?

<details><summary>Answer</summary>

A token is usable as it is: anyone with it is logged in. If the table leaks, hashes can't be used.
</details>

**40.** Why is SHA-256 fine for session tokens, when passwords need bcrypt?

<details><summary>Answer</summary>

Tokens are long and random, so they can't be guessed; slowness adds nothing and would slow every request. Passwords are guessable, so they need a slow hash.
</details>

**41.** What do the cookie flags `HttpOnly`, `Secure` and `SameSite` do?

<details><summary>Answer</summary>

`HttpOnly`: JavaScript can't read it. `Secure`: only sent over HTTPS. `SameSite`: not sent on requests started by other sites.
</details>

### Sessions vs JWTs

**42.** What's the main difference between a session token and a JWT?

<details><summary>Answer</summary>

A session token is a random ID the server looks up; a JWT contains the user's details, signed, so the server trusts it without a lookup.
</details>

**43.** How do you log a user out everywhere with sessions, and why is it hard with JWTs?

<details><summary>Answer</summary>

Delete all their session rows. A JWT has nothing stored to delete, so it stays valid until it expires (unless you add a blocklist).
</details>

**44.** Is a JWT's payload secret?

<details><summary>Answer</summary>

No. JWTs are usually signed, not encrypted: anyone can read the payload, but can't change it without breaking the signature.
</details>

**45.** How does the server check a JWT's signature?

<details><summary>Answer</summary>

It recalculates the signature from the token's header and payload using its secret key, and compares it with the token's signature. A mismatch means the token was changed (or not made by this server).
</details>

**46.** Why can't an attacker change a JWT's payload and make a new valid signature?

<details><summary>Answer</summary>

Making the signature needs the server's secret key, which never leaves the server.
</details>

**47.** Does a JWT's payload change from one request to the next? Who makes a new one?

<details><summary>Answer</summary>

No. The client sends the same token, created at login, on every request; only the URL and body change. If the payload needs to change, the server issues a new token with a new signature.
</details>

**48.** What are access and refresh tokens, and what problem do they solve?

<details><summary>Answer</summary>

An access token is a short-lived JWT used on every request; a refresh token is long-lived, stored in the database, and swapped for new access tokens. They fix JWT revocation: delete the refresh token and the user is logged out once the access token expires.
</details>

**49.** Why don't sessions need refresh tokens?

<details><summary>Answer</summary>

Every request already checks the `sessions` table, so deleting the row logs the user out immediately.
</details>

### User enumeration

**50.** What is user enumeration?

<details><summary>Answer</summary>

Using an app's responses (messages or timing) to find out which emails have accounts.
</details>

**51.** What two things does login keep identical for an unknown email and a wrong password?

<details><summary>Answer</summary>

The error message (and status), and the response time.
</details>

### Database constraints

**52.** What happens if you delete a user who still has sessions, when `sessions.user_id` is `REFERENCES users(id)` without `ON DELETE CASCADE`?

<details><summary>Answer</summary>

Postgres refuses the delete with a foreign key violation (`23503`).
</details>

**53.** What does `ON DELETE CASCADE` do?

<details><summary>Answer</summary>

Deleting the referenced row (the user) also deletes the rows that reference it (their sessions).
</details>

### Auth middleware

**54.** What is middleware?

<details><summary>Answer</summary>

A function that wraps a handler and runs on every request to it, before (and after) the handler. In Go: `func(next http.Handler) http.Handler`.
</details>

**55.** How does middleware reject a request?

<details><summary>Answer</summary>

It writes the error response and returns without calling `next.ServeHTTP`, so the handler never runs.
</details>

**56.** How does the auth middleware pass the logged-in user to the handler?

<details><summary>Answer</summary>

It adds the user to the request's context with `context.WithValue`, and passes on `r.WithContext(ctx)`. The handler reads it back from `r.Context()`.
</details>

**57.** Why put the session check in middleware rather than in each handler?

<details><summary>Answer</summary>

It's written once, it can't be forgotten on a new route, and handlers keep their one job.
</details>

### Logout

**58.** What are the two halves of logout, and which one actually logs the user out?

<details><summary>Answer</summary>

Deleting the session row (server) and clearing the cookie (browser). Deleting the row is what logs them out: the token then matches nothing.
</details>

**59.** Why is logout `POST`, not `GET`?

<details><summary>Answer</summary>

A `GET` should never change anything. Browsers, link previews and other sites can send `GET` requests on their own.
</details>

**60.** Why does logout return `204` even when there's no cookie or session?

<details><summary>Answer</summary>

The result is the same either way: logged out. That makes it idempotent.
</details>

### Ownership checks (403 vs 404)

**61.** What's the difference between authentication and authorization?

<details><summary>Answer</summary>

Authentication: who are you? Authorization: are you allowed to do this?
</details>

**62.** Why does Bo get `404`, not `403`, for Ada's board?

<details><summary>Answer</summary>

`403` would confirm the board exists, which leaks information. `404` makes "not yours" look the same as "doesn't exist".
</details>

**63.** When someone creates a board, where must the server get its `owner_id` (the user ID) from?

<details><summary>Answer</summary>

The session (the logged-in user), never the request body.
</details>

**64.** Columns and cards have no `owner_id`. How are they checked?

<details><summary>Answer</summary>

Their queries join up to the board (`cards → columns → boards`) and check `boards.owner_id`.
</details>

## 4. Transactions

### What a transaction is (ACID)

**65.** What does a transaction guarantee about the statements inside it?

<details><summary>Answer</summary>

They succeed or fail as one unit: either every change is kept, or none is.
</details>

**66.** What do the four letters of ACID stand for?

<details><summary>Answer</summary>

Atomicity (all or nothing), Consistency (constraints still hold), Isolation (others don't see half-finished changes), Durability (committed changes survive a crash).
</details>

**67.** When do you need an explicit transaction?

<details><summary>Answer</summary>

When two or more statements must stay in step. A single statement is already its own transaction.
</details>

### `BeginTx` / `Commit` / `Rollback`

**68.** What does `BeginTx` do with the connection pool?

<details><summary>Answer</summary>

It takes one connection from the pool, sends `BEGIN` on it, and ties the returned `*sql.Tx` to that connection until `Commit` or `Rollback`.
</details>

**69.** What's the difference between `Commit` and `Rollback`?

<details><summary>Answer</summary>

`Commit` makes every change permanent and visible to others. `Rollback` throws every change since `BEGIN` away.
</details>

**70.** What happens to an open transaction if the request's context is cancelled?

<details><summary>Answer</summary>

`database/sql` rolls it back automatically.
</details>

### `defer tx.Rollback()`

**71.** Why call `defer tx.Rollback()` straight after `BeginTx`?

<details><summary>Answer</summary>

So every early `return` rolls back, without writing `Rollback` before each one. A forgotten rollback leaves the connection checked out and locks held.
</details>

**72.** Why is it safe for the deferred `Rollback` to run after a successful `Commit`?

<details><summary>Answer</summary>

It does nothing and returns `sql.ErrTxDone`, which is ignored.
</details>

### Every query on `tx`, not `db`

**73.** Inside a transaction, what goes wrong if one statement uses `s.db` instead of `tx`?

<details><summary>Answer</summary>

It runs on a different connection, outside the transaction: it commits on its own, isn't undone by `Rollback`, and can't see the transaction's uncommitted changes.
</details>

**74.** Why doesn't the compiler catch this mistake?

<details><summary>Answer</summary>

`*sql.DB` and `*sql.Tx` have the same method names (`ExecContext`, `QueryRowContext`).
</details>

### `MoveCard` (`POST /cards/{id}/move`)

**75.** Why does moving a card have its own endpoint instead of using `PUT /cards/{id}`?

<details><summary>Answer</summary>

A move does two writes: it changes the card's column and records a row in `card_movements`. Both must happen in one transaction.
</details>

**76.** Why is the card's current column read inside the transaction?

<details><summary>Answer</summary>

The history row needs `from_column_id`, and it must be read before the `UPDATE` overwrites it, as part of the same unit of work.
</details>

**77.** What does "append-only" mean for `card_movements`?

<details><summary>Answer</summary>

Rows are only ever added, never updated or deleted. It's a history log.
</details>

### Mapping transaction failures to API errors

**78.** After a failed move, what state is the database in?

<details><summary>Answer</summary>

Unchanged. The deferred `Rollback` undid any partial work.
</details>

**79.** Why does the store return `ErrColumnNotFound` instead of a status code?

<details><summary>Answer</summary>

The store shouldn't know about HTTP. It reports what went wrong with a sentinel error; the handler checks it with `errors.Is` and picks `404`.
</details>
