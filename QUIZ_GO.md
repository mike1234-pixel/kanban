# Go Quiz

Recall questions on the syntax in `LEARNING_GO.md`. Try answering before opening each answer.

## Basics

### Packages and imports

**1.** What package name does a Go program's entry point use?

<details><summary>Answer</summary>

`main`.
</details>

**2.** What does `_` in front of an import mean?

<details><summary>Answer</summary>

Import it only for its side effects, e.g. registering a database driver. You don't use its names directly.
</details>

### Exported vs unexported

**3.** What makes a name visible outside its package?

<details><summary>Answer</summary>

Starting it with a capital letter.
</details>

**4.** Why can't `encoding/json` see a lowercase struct field?

<details><summary>Answer</summary>

It's unexported, and `encoding/json` is a different package.
</details>

### Variables

**5.** What's the difference between `x := 5` and `var x int`?

<details><summary>Answer</summary>

`:=` declares and assigns, inferring the type. `var x int` declares `x` with the zero value, `0`.
</details>

**6.** Where can't you use `:=`?

<details><summary>Answer</summary>

Outside functions, at package level. Use `var` there.
</details>

### Zero values

**7.** What is a zero value?

<details><summary>Answer</summary>

The value Go gives a variable when you don't set one. Every variable always has a value.
</details>

**8.** What are the zero values of `int`, `string` and `bool`?

<details><summary>Answer</summary>

`0`, `""` and `false`.
</details>

**9.** What is the zero value of a pointer, slice, map or error?

<details><summary>Answer</summary>

`nil`.
</details>

### Constants

**10.** What keyword declares a value that can never change?

<details><summary>Answer</summary>

`const`.
</details>

### Multiple return values

**11.** By convention, which return value is the error?

<details><summary>Answer</summary>

The last one.
</details>

**12.** How do you ignore a return value you don't need?

<details><summary>Answer</summary>

Assign it to `_`.
</details>

### `if` and `for`

**13.** In `if err := f(); err != nil { }`, where can you use `err`?

<details><summary>Answer</summary>

Only inside the `if` (and any `else`).
</details>

**14.** How many loop keywords does Go have?

<details><summary>Answer</summary>

One: `for`. It also does the job of `while`.
</details>

**15.** How do you loop over every element of a slice with its index?

<details><summary>Answer</summary>

`for i, item := range items { }`.
</details>

### Type conversion

**16.** How do you convert a string to a byte slice?

<details><summary>Answer</summary>

`[]byte(s)`. The reverse is `string(b)`.
</details>

## Types

### Structs

**17.** What happens to fields you leave out of a struct literal?

<details><summary>Answer</summary>

They get their zero values.
</details>

### Struct tags

**18.** What does the struct tag `json:"created_at"` do?

<details><summary>Answer</summary>

Sets the field's JSON key to `created_at`.
</details>

**19.** Which struct tag keeps a field out of JSON completely?

<details><summary>Answer</summary>

`json:"-"`.
</details>

### Anonymous structs and embedding

**20.** What is an anonymous struct, and when is it handy?

<details><summary>Answer</summary>

A struct type with no name, declared where it's used. Handy for a request body that doesn't match a model exactly.
</details>

**21.** What does embedding a type in a struct do?

<details><summary>Answer</summary>

Pulls its fields and methods into the outer struct, so `payload.Name` works when `payload` embeds `models.User`.
</details>

**22.** Is embedding the same as inheritance?

<details><summary>Answer</summary>

No. The outer struct *contains* the embedded value; it isn't one.
</details>

### Pointers

**23.** What do `&x` and `*T` mean?

<details><summary>Answer</summary>

`&x` is the address of `x`. `*T` is the type "pointer to T".
</details>

**24.** Why does `Decode(&user)` take a pointer?

<details><summary>Answer</summary>

So it can fill in the original `user`. Without a pointer it would get a copy.
</details>

**25.** Why use `*int` instead of `int` for an optional field?

<details><summary>Answer</summary>

`nil` means "missing", which an `int` can't express: its zero value `0` might be a real value.
</details>

### Slices

**26.** Why must you write `cards = append(cards, c)`, not just `append(cards, c)`?

<details><summary>Answer</summary>

`append` returns the new slice; the old variable isn't updated.
</details>

**27.** How do you get a slice's length?

<details><summary>Answer</summary>

`len(s)`.
</details>

### Maps

**28.** How do you create an empty, usable map?

<details><summary>Answer</summary>

`make(map[string]int)` or `map[string]int{}`.
</details>

**29.** How do you tell a missing key apart from a zero value?

<details><summary>Answer</summary>

`v, ok := m[key]`: `ok` is `false` if the key is missing.
</details>

**30.** What happens if you write to a `nil` map?

<details><summary>Answer</summary>

It panics.
</details>

### `any`

**31.** What does `any` mean?

<details><summary>Answer</summary>

A value of any type. It's short for `interface{}`.
</details>

## Functions and methods

### Parameters that share a type

**32.** What does `func f(tokenHash, userID string)` mean?

<details><summary>Answer</summary>

Both `tokenHash` and `userID` are strings. Neighbouring parameters of the same type can share one type name.
</details>

### Methods and receivers

**33.** What is a receiver?

<details><summary>Answer</summary>

The value a method is attached to, written before the method name: `func (h *UserHandler) CreateUser(...)`.
</details>

**34.** What's the difference between a pointer receiver and a value receiver?

<details><summary>Answer</summary>

A pointer receiver works on the original; a value receiver gets a copy.
</details>

### Functions as values

**35.** In `mux.HandleFunc("POST /users", userHandler.CreateUser)`, why is there no `()` after `CreateUser`?

<details><summary>Answer</summary>

The method is passed as a value for the router to call later; `()` would call it now.
</details>

**36.** What is a closure?

<details><summary>Answer</summary>

An anonymous function that uses variables from the code around it.
</details>

### `defer`

**37.** When does a deferred call run?

<details><summary>Answer</summary>

When the surrounding function returns, however it returns.
</details>

**38.** Where should `defer rows.Close()` go?

<details><summary>Answer</summary>

Right after successfully getting `rows`, so cleanup can't be forgotten.
</details>

## Interfaces

**39.** How does a type satisfy an interface in Go?

<details><summary>Answer</summary>

By having all its methods. It's automatic; there's no `implements` keyword.
</details>

**40.** Why do handlers depend on the `Store` interface rather than `*PostgresStore`?

<details><summary>Answer</summary>

So tests can pass in a fake store instead of a real database.
</details>

## Errors

**41.** What does a `nil` error mean?

<details><summary>Answer</summary>

Success.
</details>

**42.** Go has no exceptions. How do errors travel up the call chain?

<details><summary>Answer</summary>

They're returned, and each caller checks `if err != nil`.
</details>

**43.** What does `%w` do in `fmt.Errorf`?

<details><summary>Answer</summary>

Wraps the original error, adding context while keeping it findable with `errors.Is` and `errors.As`.
</details>

**44.** What's the difference between `errors.Is` and `errors.As`?

<details><summary>Answer</summary>

`errors.Is` checks for a specific error **value** (e.g. `ErrEmailTaken`). `errors.As` finds an error of a specific **type** (e.g. `*pq.Error`) and gives it to you.
</details>

## Testing

**45.** How must test files and test functions be named?

<details><summary>Answer</summary>

Files end in `_test.go`; functions start with `Test` and take `t *testing.T`.
</details>

**46.** What's the difference between `t.Fatalf` and `t.Errorf`?

<details><summary>Answer</summary>

`t.Fatalf` fails and stops the test; `t.Errorf` fails and keeps going.
</details>

**47.** What does `./...` mean in `go test ./...`?

<details><summary>Answer</summary>

This folder and every folder below it.
</details>
