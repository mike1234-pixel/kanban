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

### Arrays

**26.** What's the difference between an array and a slice?

<details><summary>Answer</summary>

An array has a fixed length that's part of its type (`[3]int`). A slice can grow, and points to an underlying array.
</details>

**27.** What happens when you assign one array to another variable and change the copy?

<details><summary>Answer</summary>

The original doesn't change. Assigning an array copies every element.
</details>

### Slices

**28.** Why must you write `cards = append(cards, c)`, not just `append(cards, c)`?

<details><summary>Answer</summary>

`append` returns the new slice; the old variable isn't updated.
</details>

**29.** How do you get a slice's length?

<details><summary>Answer</summary>

`len(s)`.
</details>

**30.** What are a slice's length and capacity?

<details><summary>Answer</summary>

Length: how many elements are in use. Capacity: how much room the underlying array has before `append` must make a bigger one.
</details>

**31.** Does slicing (`s[1:3]`) copy the elements? What follows from that?

<details><summary>Answer</summary>

No. The new slice points into the same underlying array, so changing an element through it changes the original too. Use `copy` or `slices.Clone` for an independent copy.
</details>

### Maps

**32.** How do you create an empty, usable map?

<details><summary>Answer</summary>

`make(map[string]int)` or `map[string]int{}`.
</details>

**33.** How do you tell a missing key apart from a zero value?

<details><summary>Answer</summary>

`v, ok := m[key]`: `ok` is `false` if the key is missing.
</details>

**34.** What happens if you write to a `nil` map?

<details><summary>Answer</summary>

It panics.
</details>

**35.** How do you remove a key from a map?

<details><summary>Answer</summary>

`delete(m, key)`. Nothing happens if the key isn't there.
</details>

**36.** In what order does `for k, v := range m` visit a map's keys?

<details><summary>Answer</summary>

A random order that can change between runs. Sort the keys first if order matters.
</details>

### Choosing a data structure

**37.** When would you use a struct, a slice and a map?

<details><summary>Answer</summary>

Struct: a fixed set of named fields. Slice: an ordered list that can grow. Map: looking values up by a key.
</details>

**38.** Which are copied when assigned, and which are shared: struct, array, slice, map?

<details><summary>Answer</summary>

Structs and arrays are copied. Slices and maps are shared: both variables see the same data.
</details>

### `any`

**39.** What does `any` mean?

<details><summary>Answer</summary>

A value of any type. It's short for `interface{}`.
</details>

### Type assertions

**40.** What does `v, ok := x.(models.User)` do?

<details><summary>Answer</summary>

Checks whether `x` (an `any` or interface value) holds a `models.User`. If so, `v` is that user and `ok` is `true`; if not, `v` is the zero value and `ok` is `false`.
</details>

**41.** Why use the `, ok` form of a type assertion?

<details><summary>Answer</summary>

Without it, a value of the wrong type causes a panic.
</details>

## Functions and methods

### Parameters that share a type

**42.** What does `func f(tokenHash, userID string)` mean?

<details><summary>Answer</summary>

Both `tokenHash` and `userID` are strings. Neighbouring parameters of the same type can share one type name.
</details>

### Methods and receivers

**43.** What is a receiver?

<details><summary>Answer</summary>

The value a method is attached to, written before the method name: `func (h *UserHandler) CreateUser(...)`.
</details>

**44.** What's the difference between a pointer receiver and a value receiver?

<details><summary>Answer</summary>

A pointer receiver works on the original; a value receiver gets a copy.
</details>

### Functions as values

**45.** In `mux.HandleFunc("POST /users", userHandler.CreateUser)`, why is there no `()` after `CreateUser`?

<details><summary>Answer</summary>

The method is passed as a value for the router to call later; `()` would call it now.
</details>

**46.** What is a closure?

<details><summary>Answer</summary>

An anonymous function that uses variables from the code around it.
</details>

### `defer`

**47.** When does a deferred call run?

<details><summary>Answer</summary>

When the surrounding function returns, however it returns.
</details>

**48.** Where should `defer rows.Close()` go?

<details><summary>Answer</summary>

Right after successfully getting `rows`, so cleanup can't be forgotten.
</details>

## Interfaces

**49.** How does a type satisfy an interface in Go?

<details><summary>Answer</summary>

By having all its methods. It's automatic; there's no `implements` keyword.
</details>

**50.** Why do handlers depend on the `Store` interface rather than `*PostgresStore`?

<details><summary>Answer</summary>

So tests can pass in a fake store instead of a real database.
</details>

## Errors

**51.** What does a `nil` error mean?

<details><summary>Answer</summary>

Success.
</details>

**52.** Go has no exceptions. How do errors travel up the call chain?

<details><summary>Answer</summary>

They're returned, and each caller checks `if err != nil`.
</details>

**53.** What does `%w` do in `fmt.Errorf`?

<details><summary>Answer</summary>

Wraps the original error, adding context while keeping it findable with `errors.Is` and `errors.As`.
</details>

**54.** What's the difference between `errors.Is` and `errors.As`?

<details><summary>Answer</summary>

`errors.Is` checks for a specific error **value** (e.g. `ErrEmailTaken`). `errors.As` finds an error of a specific **type** (e.g. `*pq.Error`) and gives it to you.
</details>

## Testing

**55.** How must test files and test functions be named?

<details><summary>Answer</summary>

Files end in `_test.go`; functions start with `Test` and take `t *testing.T`.
</details>

**56.** What's the difference between `t.Fatalf` and `t.Errorf`?

<details><summary>Answer</summary>

`t.Fatalf` fails and stops the test; `t.Errorf` fails and keeps going.
</details>

**57.** What does `./...` mean in `go test ./...`?

<details><summary>Answer</summary>

This folder and every folder below it.
</details>
