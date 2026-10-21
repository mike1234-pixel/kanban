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

## Go basics

### Zero values

In Go, every variable always has a value. If you don't give it one, Go fills in that type's zero value:

| Type | Zero value |
| --- | --- |
| `int` | `0` |
| `string` | `""` (empty string) |
| `bool` | `false` |
| pointer, slice, map | `nil` |
| struct | every field set to its own zero value |
