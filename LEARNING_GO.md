# Learning Go

Quick reference for the Go syntax used in this project. Examples come from the codebase. Backend concepts live in `LEARNING.md`.

## Basics

### Packages and imports

Every file starts with `package <name>`. A folder is one package; `main` is the program's entry point.

```go
import (
	"net/http"                    // standard library
	"kanban/internal/models"      // this project
	"github.com/lib/pq"           // third party (listed in go.mod)
)
```

- Use an import as `http.StatusOK`: the package name, a dot, then the name.
- `_ "github.com/lib/pq"` imports a package only for its side effects (e.g. registering a driver).
- `internal/` packages can only be imported from inside this project.

### Exported vs unexported

A capital first letter means visible outside the package; lowercase means private to it.

```go
type User struct { Name string }   // exported: other packages (and encoding/json) can see it
const codeUserNotFound = "..."     // unexported: only package handlers can use it
```

### Variables: `:=` vs `var`

```go
id := r.PathValue("id")     // declare and assign; type is inferred. Inside functions only.
var user models.User        // declare with the zero value, assign later
```

- `:=` needs at least one new variable on the left.
- `var` is used when you want the zero value first, e.g. before decoding into it.

### Zero values

Every variable always has a value. If you don't give it one, Go uses the type's zero value:

| Type | Zero value |
| --- | --- |
| `int` | `0` |
| `string` | `""` |
| `bool` | `false` |
| pointer, slice, map, interface, error | `nil` |
| struct | every field set to its own zero value |

### Constants

```go
const (
	codeInvalidJSON      = "INVALID_JSON"
	codeValidationFailed = "VALIDATION_FAILED"
)
```

- Fixed at compile time; can't be changed. A `const ( ... )` block groups related ones.

### Multiple return values

Functions can return several values. The last is usually an `error`.

```go
createdAt, err := h.Store.SaveUser(r.Context(), user)
```

- `_` discards a value you don't need: `_, err = io.ReadFull(...)`.

### `if` with a short statement

```go
if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
	// err only exists inside this if
}
```

- Run something, then test it, in one line. Keeps `err` scoped to the `if`.

### `for`

Go's only loop keyword.

```go
for rows.Next() { ... }           // like "while"
for i := 0; i < 3; i++ { ... }    // classic counter
for i, card := range cards { ... } // each index and element
```

### Type conversion

`T(value)` converts between types: `[]byte(payload.Password)`, `string(hash)`.

## Types

### Structs

A named group of fields: Go's main way to model data.

```go
user := models.User{ID: id, Name: "Ada"}   // struct literal; missing fields get zero values
user.Email = "ada@example.com"             // field access with a dot
```

### Struct tags

Text after a field that other packages read, e.g. `encoding/json` and the validator.

```go
Email        string `json:"email" validate:"required,email"`  // JSON key "email", plus rules
PasswordHash string `json:"-"`                                // never in JSON
```

### Anonymous structs

A struct type with no name, used once:

```go
var payload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
```

- Handy for request bodies that don't match a model exactly.

### Embedding

Writing a type with no field name pulls in its fields (and methods).

```go
var payload struct {
	models.User               // payload.Name, payload.Email… all work
	Password string
}
```

- Embedding is not inheritance: `payload` *contains* a `User` (`payload.User`); it isn't one.

### Pointers

A pointer holds the address of a value. `&x` gets the address; `*T` is the type "pointer to T".

```go
json.NewDecoder(r.Body).Decode(&user)   // pass a pointer so Decode can fill user in
Position *int `json:"position"`          // nil = missing; 0 = really zero
return &PostgresStore{db: db}, nil       // return a pointer to a new struct
```

- Without a pointer, a function gets a **copy** and can't change your variable.

### Slices

A growable list: `[]T`.

```go
cards := []models.Card{}       // empty, encodes as []
var cards []models.Card        // nil, encodes as null
cards = append(cards, c)       // append returns the new slice; always reassign
len(cards)                     // length
```

### Maps

Key → value lookup: `map[K]V`.

```go
colMap := make(map[string]int)   // make creates an empty, usable map
colMap[c.ID] = len(columns)      // set
i, ok := colMap[id]              // get; ok is false if the key is missing
```

- Writing to a `nil` map panics, so create it with `make` or a literal (`map[string]int{}`).

### `any`

`any` means a value of any type (short for `interface{}`): `func writeJSON(w http.ResponseWriter, status int, v any)`.

## Functions and methods

### Methods and receivers

A method is a function attached to a type. The receiver goes before the name.

```go
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) { ... }
```

- **Pointer receiver** (`*UserHandler`): the method works on the original, not a copy. The default choice for structs.
- **Value receiver** (`testStore`): gets a copy. Fine for small types that aren't changed.

### Functions as values

Functions can be passed around and stored like any other value.

```go
mux.HandleFunc("POST /users", userHandler.CreateUser)   // pass a method without calling it (no ())
saveUser: func(models.User) error { return store.ErrEmailTaken }   // anonymous function in a test fake
```

- A **closure** is an anonymous function that uses variables from around it, e.g. setting `called = true` in a test.

### `defer`

Runs a call when the surrounding function returns, however it returns.

```go
defer rows.Close()   // the rows are always closed, even after an early return
```

- Put it right after acquiring the resource, so cleanup can't be forgotten.

## Interfaces

A set of method signatures. Any type with those methods satisfies the interface automatically; there's no `implements` keyword.

```go
type Store interface {
	SaveUser(ctx context.Context, user models.User) (time.Time, error)
	// ...
}
```

- `*PostgresStore` and the test's `testStore` both satisfy `Store`, so handlers work with either.
- Embedding an interface in a struct (`testStore` embeds `store.Store`) satisfies it with no methods written; calling one you didn't write panics.

## Errors

### Errors are values

Functions return an `error`; `nil` means success. Check it straight away:

```go
if err != nil {
	return nil, err
}
```

- No exceptions: errors go up the call chain by being returned.

### Sentinel errors

A named error value that callers can check for.

```go
var ErrEmailTaken = errors.New("email already in use")
if errors.Is(err, store.ErrEmailTaken) { ... }
```

### Wrapping

`fmt.Errorf` with `%w` adds context but keeps the original error inside:

```go
return nil, fmt.Errorf("failed to ping postgres: %w", err)
```

- `errors.Is` sees through wrapping to a specific error **value**.
- `errors.As` finds an error of a specific **type** and gives you it: `var pqErr *pq.Error; errors.As(err, &pqErr)`.

## Testing

```go
func TestCreateUser(t *testing.T) {
	t.Run("rejects an invalid email", func(t *testing.T) { ... })
}
```

- Tests live in `_test.go` files; test functions start with `Test` and take `*testing.T`.
- `t.Run` makes a named subtest.
- `t.Fatalf` fails and stops the test; `t.Errorf` fails and keeps going.
- `httptest.NewRequest` and `httptest.NewRecorder` call a handler without a real server.
- Run with `go test ./...` (`./...` means this folder and every folder below it).
