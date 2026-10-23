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

**13.** What's a path param for, and what's a query param for?

<details><summary>Answer</summary>

A path param identifies *which* resource (`/boards/{id}`). A query param filters or shapes the result (`/cards?status=completed`).
</details>

**14.** What does a missing query param return in Go?

<details><summary>Answer</summary>

An empty string, `""`.
</details>

**15.** What type do path and query params arrive as?

<details><summary>Answer</summary>

Strings. Numbers need converting, e.g. with `strconv.Atoi`.
</details>

### Consistent errors

**16.** In an error response, what's `code` for, and what's `message` for?

<details><summary>Answer</summary>

`code` is for the client's code to branch on; `message` is for people to read.
</details>

**17.** Where does the HTTP status code go in a response, and where does the application error code go?

<details><summary>Answer</summary>

The HTTP status code is in the status line, the first line of the response. The application error code is in the JSON body.
</details>

**18.** Why do 500 errors send a generic message instead of the real error?

<details><summary>Answer</summary>

The real error can leak database details to attackers, and the client can't act on it anyway. It goes to the logs instead.
</details>

**19.** Which layer picks the HTTP status code: the handler or the store?

<details><summary>Answer</summary>

The handler. The store only talks to the database and reports what happened.
</details>

**20.** Why does the store turn Postgres error `23505` into `ErrBoardExists`, instead of the handler checking for `23505`?

<details><summary>Answer</summary>

Only the store should know Postgres exists. Handlers depend on the `Store` interface, so a test fake or another database can return the same error.
</details>

### Content-Type

**21.** What does the `Content-Type` header tell the client?

<details><summary>Answer</summary>

What format the response body is in, e.g. `application/json`, `text/plain` or `text/html`.
</details>

**22.** In what order must a handler set headers, write the status and write the body?

<details><summary>Answer</summary>

Headers → status (`WriteHeader`) → body.
</details>

### Request lifecycle

**23.** What do TCP and TLS stand for?

<details><summary>Answer</summary>

TCP: Transmission Control Protocol. TLS: Transport Layer Security.
</details>

**24.** What are the three steps of the TCP handshake?

<details><summary>Answer</summary>

SYN (synchronize) → SYN-ACK (synchronize-acknowledge) → ACK (acknowledge).
</details>

**25.** What does TLS add on top of TCP?

<details><summary>Answer</summary>

Encryption, and server identity via a certificate. Together they make HTTPS.
</details>

**26.** Why does Go's server give each connection its own goroutine?

<details><summary>Answer</summary>

So requests are handled concurrently, and one slow client doesn't block the others.
</details>

**27.** Which part of the server decides which handler a request goes to?

<details><summary>Answer</summary>

The mux (multiplexer), also called the router.
</details>

### `context.Context`

**28.** What two things cancel a request's context?

<details><summary>Answer</summary>

The client disconnecting, or the handler returning.
</details>

**29.** Why pass `r.Context()` to database calls?

<details><summary>Answer</summary>

So the database abandons the query if the client disconnects, freeing the connection.
</details>

**30.** By convention, where does a context parameter go, and what is it named?

<details><summary>Answer</summary>

First in the parameter list, named `ctx`.
</details>

## 3. Authentication / Authorization

### Password hashing

**31.** Why must a password hash be slow?

<details><summary>Answer</summary>

Attackers guess passwords and hash each guess. A slow hash (about 100ms) makes guessing millions of passwords impractical, while one login barely notices.
</details>

**32.** What is a salt, and what attack does it stop?

<details><summary>Answer</summary>

A random value mixed into each password before hashing, so identical passwords get different hashes. It stops precomputed lookup tables (rainbow tables).
</details>

**33.** How does login check a password, if a hash can't be reversed?

<details><summary>Answer</summary>

It hashes what the user typed with the same salt and cost (read from the stored hash) and compares the result with the stored hash.
</details>

**34.** Which struct tag keeps a field out of JSON responses?

<details><summary>Answer</summary>

`json:"-"`.
</details>

### Sessions and cookies

**35.** HTTP is stateless. How does the server know who sent a request after login?

<details><summary>Answer</summary>

Login creates a session and sets a cookie holding a random session token. The browser sends the cookie on every request, and the server looks the token up in the `sessions` table.
</details>

**36.** Why store a hash of the session token rather than the token itself?

<details><summary>Answer</summary>

A token is usable as it is: anyone with it is logged in. If the table leaks, hashes can't be used.
</details>

**37.** Why is SHA-256 fine for session tokens, when passwords need bcrypt?

<details><summary>Answer</summary>

Tokens are long and random, so they can't be guessed; slowness adds nothing and would slow every request. Passwords are guessable, so they need a slow hash.
</details>

**38.** What do the cookie flags `HttpOnly`, `Secure` and `SameSite` do?

<details><summary>Answer</summary>

`HttpOnly`: JavaScript can't read it. `Secure`: only sent over HTTPS. `SameSite`: not sent on requests started by other sites.
</details>

### Sessions vs JWTs

**39.** What's the main difference between a session token and a JWT?

<details><summary>Answer</summary>

A session token is a random ID the server looks up; a JWT contains the user's details, signed, so the server trusts it without a lookup.
</details>

**40.** How do you log a user out everywhere with sessions, and why is it hard with JWTs?

<details><summary>Answer</summary>

Delete all their session rows. A JWT has nothing stored to delete, so it stays valid until it expires (unless you add a blocklist).
</details>

**41.** Is a JWT's payload secret?

<details><summary>Answer</summary>

No. JWTs are usually signed, not encrypted: anyone can read the payload, but can't change it without breaking the signature.
</details>

### User enumeration

**42.** What is user enumeration?

<details><summary>Answer</summary>

Using an app's responses (messages or timing) to find out which emails have accounts.
</details>

**43.** What two things does login keep identical for an unknown email and a wrong password?

<details><summary>Answer</summary>

The error message (and status), and the response time.
</details>

### Database constraints

**44.** What happens if you delete a user who still has sessions, when `sessions.user_id` is `REFERENCES users(id)` without `ON DELETE CASCADE`?

<details><summary>Answer</summary>

Postgres refuses the delete with a foreign key violation (`23503`).
</details>

**45.** What does `ON DELETE CASCADE` do?

<details><summary>Answer</summary>

Deleting the referenced row (the user) also deletes the rows that reference it (their sessions).
</details>
