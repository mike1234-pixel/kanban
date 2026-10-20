# Learning Log

Quick reference for the roadmap concepts I've learned.

## 1. CRUD

**Create, Read, Update, Delete**: the four operations on a stored resource. Each maps to an HTTP method; the method says *what to do*, the path says *to what*.

| Operation | Request | Success |
| --- | --- | --- |
| Create | `POST /boards` | 201 |
| Read | `GET /boards`, `GET /boards/{id}` | 200 |
| Update | `PUT /boards/{id}` | 200 |
| Delete | `DELETE /boards/{id}` | 204 |

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
- Struct tags map names: `json:"created_at"`. Only exported (Capitalised) fields are included.
- A `nil` slice encodes as `null`; an empty slice as `[]`.
- Unknown fields are ignored; wrong types are a decode error (400).

### Validation

Decoding proves the JSON is readable; validation proves it makes sense.

- Rules are struct tags: `validate:"required,email,max=254"`.
- Order: decode → validate → store. Invalid data never reaches the database.
- `required` can't tell a missing `int` from `0`; use `*int` when zero is valid.
- Database constraints are the last line of defence, and catch what validation can't see (e.g. a duplicate email → 409).

### Path/query params

- **Path params identify *which* resource:** `GET /boards/{id}`, read with `r.PathValue("id")`.
- **Query params filter or shape the result:** `/columns?board_id=…`, read with `r.URL.Query().Get("board_id")`. Returns `""` if missing.
- The router matches path params and methods (wrong method → 405); it ignores query params, so the handler checks them.
- Both arrive as strings; convert numbers with `strconv.Atoi`.

### Consistent errors

Every error has one shape, so the client handles them all the same way:

```json
{"error": {"code": "BOARD_NOT_FOUND", "message": "Board not found"}}
```

- **Codes** are for code to branch on; **messages** are for people.
- Codes are constants in `handlers/respond.go`; `writeError` builds every error.
- 500s log the real error and send a generic message, never internal details.
- The store returns sentinel errors (`store.ErrEmailTaken`); the handler checks them with `errors.Is` and picks the status.

### Content-Type

A header that tells the client how to read the body: `application/json`, `text/plain`, `text/html`.

- Order matters: set headers → `WriteHeader` → write body. Headers set later are silently ignored.
- If you don't set it, Go guesses, and JSON comes out as `text/plain`.
- `http.Error` always sends `text/plain`, so it can't send JSON errors.

### `context.Context`

Carries a cancellation signal (plus deadlines and request values) down a call chain.

- `r.Context()` is cancelled when the client disconnects or the handler returns.
- Pass it to the database (`QueryContext`, `ExecContext`, `BeginTx`) so abandoned queries stop and free their connection.
- Conventions: first parameter, named `ctx`; never stored in a struct; use `context.Background()` when there's no request.

### Not yet covered

- Request lifecycle
