# Quiz

Questions on the concepts in `LEARNING.md`. Try answering before opening each answer.

## 1. CRUD

**1.** Which HTTP method and success status would you use to delete a card?

<details><summary>Answer</summary>

`DELETE /cards/{id}`, returning 204 No Content.
</details>

**2.** What does *idempotent* mean, and which of `GET`, `POST`, `PUT`, `DELETE` is not idempotent?

<details><summary>Answer</summary>

Repeating the request has the same effect as doing it once. `POST` isn't: two creates can make two resources.
</details>

**3.** How does the store tell the handler that a `DELETE` didn't find the row?

<details><summary>Answer</summary>

`RowsAffected()` returns 0, so the store returns `false` and the handler sends 404.
</details>

## 2. HTTP Fundamentals

### Status codes

**4.** A client sends a valid request, but the database is down. 4xx or 5xx? Why?

<details><summary>Answer</summary>

5xx (500). It's the server's fault, and the same request might work later.
</details>

**5.** What status does a handler send if it writes a body but never calls `w.WriteHeader`?

<details><summary>Answer</summary>

200 OK. The first write sends it automatically.
</details>

**6.** Why does `DeleteBoard` call `w.WriteHeader(http.StatusNoContent)` instead of `writeJSON`?

<details><summary>Answer</summary>

A 204 response must have no body.
</details>

### JSON encode/decode

**7.** Why does `Decode` need `&board` and not `board`?

<details><summary>Answer</summary>

It fills in the struct, so it needs a pointer to it. A copy would be filled and then thrown away.
</details>

**8.** `GET /users` returned `null` on an empty table. What caused it, and what fixed it?

<details><summary>Answer</summary>

`var users []models.User` is a `nil` slice, which encodes as `null`. `users := []models.User{}` is an empty slice, which encodes as `[]`.
</details>

**9.** A struct has a field `title string` (lowercase). What happens when you decode `{"title":"Sprint"}` into it?

<details><summary>Answer</summary>

Nothing. Unexported fields are invisible to `encoding/json`, so it stays empty.
</details>

### Validation

**10.** What's the difference between decoding and validating?

<details><summary>Answer</summary>

Decoding proves the JSON is readable; validation proves the values make sense (lengths, formats, required fields).
</details>

**11.** Why is `Position` a `*int` in `UpdateColumn`'s payload?

<details><summary>Answer</summary>

`required` can't tell a missing `int` from `0`. With a pointer, missing is `nil` (rejected) and `0` is a pointer to zero (allowed).
</details>

**12.** Why can't a `validate` tag enforce unique emails? What does instead?

<details><summary>Answer</summary>

Validation only sees the request, not what's in the table. The `UNIQUE` constraint catches it, and the store maps Postgres error `23505` to `ErrEmailTaken` → 409.
</details>

### Path/query params

**13.** Path or query param: the ID of a board? A filter for archived boards?

<details><summary>Answer</summary>

The ID is a path param (`/boards/{id}`); it identifies the resource. The filter is a query param (`/boards?archived=false`); it shapes the result.
</details>

**14.** What does `r.URL.Query().Get("board_id")` return if the parameter is missing?

<details><summary>Answer</summary>

An empty string, `""`.
</details>

**15.** Several handlers check `if id == ""` after `r.PathValue("id")`. Can it ever be true through the router? What about in tests?

<details><summary>Answer</summary>

Not through the router: `{id}` only matches a non-empty path segment, so `/boards/` never reaches `GetBoard`. In tests that call a handler directly without `SetPathValue`, it is `""`.
</details>

### Consistent errors

**16.** In an error response, what's `code` for, and what's `message` for?

<details><summary>Answer</summary>

`code` is for the client's code to branch on; `message` is for people to read.
</details>

**17.** Why does `writeInternalError` send "Something went wrong" instead of the real error?

<details><summary>Answer</summary>

The real error can leak database details to attackers, and the client can't act on it anyway. It goes to the logs for you instead.
</details>

**18.** The store can't choose HTTP status codes. How does it tell the handler "this email is taken"?

<details><summary>Answer</summary>

It returns the sentinel error `store.ErrEmailTaken`; the handler checks `errors.Is(err, store.ErrEmailTaken)` and chooses 409.
</details>

### Content-Type

**19.** What goes wrong if you set `Content-Type` after calling `w.WriteHeader`?

<details><summary>Answer</summary>

Nothing visible, and that's the problem: headers are already sent, so the new value is silently ignored.
</details>

**20.** Why did the handlers stop using `http.Error`?

<details><summary>Answer</summary>

It always sends `text/plain`, so it can't send JSON errors.
</details>

### Request lifecycle

**21.** Put these in order: mux, TCP handshake, handler, parse, goroutine, TLS handshake, accept.

<details><summary>Answer</summary>

TCP handshake → TLS handshake (HTTPS only) → accept → goroutine → parse → mux → handler.
</details>

**22.** You run `curl localhost:8080/boards`. Which handshakes happen?

<details><summary>Answer</summary>

Only TCP. It's `http://`, so there's no TLS.
</details>

**23.** `r.Body` isn't read when your handler starts. Which line in your code reads it?

<details><summary>Answer</summary>

`json.NewDecoder(r.Body).Decode(&board)` (or the equivalent in each create/update handler).
</details>

**24.** Why must anything your handlers share be safe for concurrent use?

<details><summary>Answer</summary>

Each connection gets its own goroutine, so many handlers run at the same time.
</details>

### `context.Context`

**25.** When is `r.Context()` cancelled?

<details><summary>Answer</summary>

When the client disconnects, or when the handler returns.
</details>

**26.** What happens to a slow query if the client disconnects and you used `QueryContext(ctx, ...)`? And with plain `Query(...)`?

<details><summary>Answer</summary>

With the context, the query is cancelled and its connection goes back to the pool. Without it, Postgres keeps working on a result nobody will read.
</details>

**27.** Why shouldn't you store a context in a struct?

<details><summary>Answer</summary>

A context belongs to one call chain (one request). A struct can outlive that request.
</details>
