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
for i := range 10 { ... }          // i = 0, 1, … 9
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

- **Two ways to reach an embedded field:**
  - `payload.Name`: the short form. Go finds `Name` inside the embedded `User`. This is the usual way, and the point of embedding.
  - `payload.User.Name`: the full form. The embedded field has no name you wrote, so Go names it after its type (`User`).
- Use the full form when both structs have a field with the same name: the short form then means the outer struct's field.
- Embedding is not inheritance: `payload` *contains* a `User` (`payload.User`); it isn't one.

### Pointers

A pointer is a variable that holds **where** another value is stored (its address), not the value itself.

Three pieces of syntax:

| You write | Meaning |
| --- | --- |
| `&x` | the address of `x` |
| `*int` | a type: "pointer to an int" |
| `*p` | the value at the address in `p` |

```go
x := 5
p := &x      // p holds x's address. p's type is *int
*p           // 5: the value at that address
*p = 10      // change the value at that address
x            // 10: x changed, because p points to x
```

**Why use one:** a function gets a copy of what you pass it. Pass a pointer and it can change your original.

```go
func double(n int)      { n = n * 2 }     // changes its own copy only
func doublePtr(n *int)  { *n = *n * 2 }   // changes the original

x := 5
double(x)      // x is still 5
doublePtr(&x)  // x is now 10
```

A pointer that points to nothing is `nil`. Reading `*p` when `p` is `nil` panics.

In this codebase:

```go
json.NewDecoder(r.Body).Decode(&user)   // pass a pointer so Decode can fill user in
Position *int `json:"position"`          // nil = missing; 0 = really zero
return &PostgresStore{db: db}, nil       // return a pointer to a new struct
```

### Arrays

A fixed-length list. The length is part of the type: `[3]int` and `[4]int` are different types.

```go
var scores [3]int        // [0 0 0]: zero values, not nil
scores[0] = 10           // [10 0 0]
copy := scores           // copies all 3 values
copy[0] = 99             // scores is still [10 0 0]
```

Assigning an array, or passing it to a function, copies every element. Arrays can't grow, so Go code rarely uses them directly; slices are built on top of them.

### Slices

A growable list: `[]T`. A slice doesn't hold the elements itself. It points to an underlying array and records two numbers: its **length** (elements in use) and its **capacity** (room in the array before it must grow).

```go
cards := []models.Card{}       // empty, encodes as []
var cards []models.Card        // nil, encodes as null
cards = append(cards, c)       // append returns the new slice; always reassign
len(cards)                     // length
```

```go
s := make([]int, 0, 10)   // length 0, capacity 10: room for 10 before growing
s = append(s, 1, 2, 3)    // [1 2 3], len 3, cap 10
```

When `append` runs out of capacity, it copies everything into a new, bigger array. That's why its result must be assigned back.

**Slicing** takes part of a slice with `s[from:to]`. `from` is included and `to` isn't:

```go
s := []int{10, 20, 30, 40}
mid := s[1:3]   // [20 30]
mid[0] = 99     // s is now [10 99 30 40]
```

`mid` doesn't copy anything. It points into the same array as `s`, so changing `mid[0]` changes `s[1]`. Use `copy` or `slices.Clone` when you need an independent copy.

`for i, v := range s` loops over a slice in order.

### Maps

Key → value lookup: `map[K]V`.

```go
ages := make(map[string]int)   // keys are strings, values are ints; make creates an empty, usable map

ages["ada"] = 36               // store 36 under the key "ada"
ages["bo"] = 29                // store 29 under the key "bo"

ages["ada"]                    // 36
age, ok := ages["cy"]          // 0, false: ok is false if the key is missing
delete(ages, "bo")             // remove a key; nothing happens if it isn't there
len(ages)                      // 1: the number of keys
```

- Writing to a `nil` map panics, so create it with `make` or a literal (`map[string]int{}`).
- `for key, value := range m` loops over a map, but **in a random order** that changes from run to run. If the order matters, collect the keys into a slice and sort it.
- Like a slice, a map refers to its data: assigning it or passing it to a function shares the same map, not a copy.

### Choosing a data structure

| Use | When | Example in this project |
| --- | --- | --- |
| Struct | A fixed set of named fields, possibly of different types | `models.User` |
| Slice | An ordered list that can grow | `[]models.Card` in a column |
| Map | Looking values up by a key | `colMap`: column ID → position in the slice |
| Array | A fixed number of items, known when writing the code | Rare; mostly inside slices |

Assigning a struct or an array copies its values. Assigning a slice or a map copies only the reference, so both variables see the same data.

### `any`

`any` means a value of any type (short for `interface{}`): `func writeJSON(w http.ResponseWriter, status int, v any)`.

### Type assertions

A type assertion asks: "is the value inside this `any` really a `string` (or some other type)? If so, give it to me as one."

**Why you need it:** once a value is stored as `any`, Go no longer knows its type, so you can't use it as that type:

```go
var v any = "hello"

v + " world"          // doesn't compile: Go doesn't know v is a string
```

**The assertion checks the type and gives the value back with its type:**

```go
s, ok := v.(string)   // s = "hello", ok = true: it was a string
n, ok := v.(int)      // n = 0,       ok = false: it wasn't an int, so n gets int's zero value

s + " world"          // works: s is a string
```

If the check fails, you get the type's zero value and `ok = false`.

In this codebase: `ctx.Value(userKey)` returns `any`, so `userFromContext` uses `ctx.Value(userKey).(models.User)` to get the user back out.

- Always use the `, ok` form. Without it (`ctx.Value(userKey).(models.User)` alone), a wrong type panics.
- Same comma-ok shape as a map lookup: `v, ok := m[key]`.

## Functions and methods

### Parameters that share a type

Neighbouring parameters with the same type can share one type name:

```go
func CreateSession(ctx context.Context, tokenHash, userID string, expiresAt time.Time) error
//                                      └── both strings ──┘
```

- `tokenHash, userID string` means the same as `tokenHash string, userID string`.
- Also works in `var` declarations and struct fields: `var storedHash, storedUserID string`.

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
- **Error vs exception:** an exception (TypeScript's `throw`) jumps out of the function by itself, up to whichever `try/catch` is waiting. A Go error is just a return value: nothing jumps, and each caller has to check it and pass it on.

  ```ts
  const user = getUser(id)           // TypeScript: if this throws, the next line never runs
  ```

  ```go
  user, err := getUser(id)           // Go: the error comes back like any other value
  if err != nil {
  	return err                     // you pass it up yourself
  }
  ```

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
- `errors.As` checks if an error is a specific **type** and, if so, puts it in your variable: `var pqErr *pq.Error; errors.As(err, &pqErr)`.

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
- In a helper, `t.Helper()` makes a failure point at the test's line, not the helper's.
- `t.Cleanup(func() { … })` runs when the test ends: the test version of `defer`.
- `t.Skip("reason")` skips a test, e.g. the store tests when there's no database.
- `go test -race ./...` adds the race detector: it fails the test if goroutines touch the same memory without protection.

## Concurrency

Used in the concurrency tests (`store/cards_test.go`, `handlers/concurrency_test.go`). Concepts are in `LEARNING.md`, step 5.

### `go func() { … }()`

Starts an anonymous function in a new goroutine. The `()` at the end calls it.

```go
go func() {
	defer wg.Done()
	codes[i] = response.Code
}()
```

- The line returns straight away; the function runs at the same time as the code after it.
- It's a closure, so it can use variables from around it (`i`, `codes`). Since Go 1.22, each loop round gets its own `i`.

### `sync.WaitGroup`

Waits for a group of goroutines to finish.

```go
var wg sync.WaitGroup
wg.Add(1)             // one more goroutine to wait for
go func() {
	defer wg.Done()   // this one has finished
	// ...
}()
wg.Wait()             // pause until every Add has had a Done
```

- Call `Add` before `go`, not inside the goroutine, or `Wait` might run before it's counted.
- Its zero value is ready to use.

### `sync.Mutex`

Lets one goroutine at a time run the code between `Lock` and `Unlock`.

```go
type memCardStore struct {
	mu    sync.Mutex
	cards map[string]models.Card
}

func (s *memCardStore) UpdateCard(…) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// read, check and write s.cards
}
```

- Put the mutex next to the data it guards, in the same struct.
- `defer Unlock()` straight after `Lock()`, so every return unlocks.
- Its zero value is unlocked and ready. Use a pointer receiver: copying a struct copies its mutex, and the copy doesn't guard anything.
