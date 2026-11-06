# Go Quiz

Recall questions on the syntax in `LEARNING_GO.md`. Try answering before opening each answer.

## Basics

### Packages and imports

**1.** What package name does a Go program's entry point use?

- [x] `main`
- [ ] `app`
- [ ] `start`
- [ ] `index`

<details><summary>Answer</summary>

`main`.
</details>

**2.** What does `_` in front of an import mean?

- [ ] Import it but keep its names private
- [ ] Import it only in tests
- [ ] Import it under the name `_`, used as `_.Open`
- [x] Import it only for its side effects, e.g. a driver

<details><summary>Answer</summary>

Import it only for its side effects, e.g. registering a database driver. You don't use its names directly.
</details>

### Exported vs unexported

**3.** What makes a name visible outside its package?

- [ ] Adding the `export` keyword
- [ ] Starting it with an underscore
- [x] Starting it with a capital letter
- [ ] Listing it in `go.mod`

<details><summary>Answer</summary>

Starting it with a capital letter.
</details>

**4.** Why can't `encoding/json` see a lowercase struct field?

- [ ] Lowercase fields can't hold strings
- [x] It's unexported, and `encoding/json` is another package
- [ ] It needs a `json` struct tag first
- [ ] JSON keys must be capitalised

<details><summary>Answer</summary>

It's unexported, and `encoding/json` is a different package.
</details>

### Variables

**5.** What's the difference between `x := 5` and `var x int`?

- [x] `:=` declares and assigns; `var x int` declares with zero value `0`
- [ ] `:=` assigns only; `var x int` declares and assigns `5`
- [ ] `:=` makes a constant; `var x int` makes a variable
- [ ] They're identical; `var` is just older

<details><summary>Answer</summary>

`:=` declares and assigns, inferring the type. `var x int` declares `x` with the zero value, `0`.
</details>

**6.** Where can't you use `:=`?

- [ ] Inside `for` loops
- [ ] Inside `if` statements
- [ ] Inside methods
- [x] At package level, outside functions

<details><summary>Answer</summary>

Outside functions, at package level. Use `var` there.
</details>

### Zero values

**7.** What is a zero value?

- [ ] The value `0`, for every type
- [ ] A variable that holds nothing, like `undefined`
- [x] The value Go gives a variable you don't set
- [ ] The first value a function returns

<details><summary>Answer</summary>

The value Go gives a variable when you don't set one. Every variable always has a value.
</details>

**8.** What are the zero values of `int`, `string` and `bool`?

- [ ] `nil`, `nil` and `nil`
- [x] `0`, `""` and `false`
- [ ] `0`, `nil` and `false`
- [ ] `0`, `" "` and `nil`

<details><summary>Answer</summary>

`0`, `""` and `false`.
</details>

**9.** What is the zero value of a pointer, slice, map or error?

- [x] `nil`
- [ ] `0`
- [ ] An empty value, e.g. `[]`
- [ ] `undefined`

<details><summary>Answer</summary>

`nil`.
</details>

### Constants

**10.** What keyword declares a value that can never change?

- [x] `const`
- [ ] `let`
- [ ] `final`
- [ ] `var`

<details><summary>Answer</summary>

`const`.
</details>

### Multiple return values

**11.** By convention, which return value is the error?

- [ ] The first one
- [ ] The second one
- [ ] Any; it's named `err`
- [x] The last one

<details><summary>Answer</summary>

The last one.
</details>

**12.** How do you ignore a return value you don't need?

- [ ] Assign it to `nil`
- [ ] Leave it out of the left side
- [x] Assign it to `_`
- [ ] Assign it to `ignore`

<details><summary>Answer</summary>

Assign it to `_`.
</details>

### `if` and `for`

**13.** In `if err := f(); err != nil { }`, where can you use `err`?

- [ ] Anywhere in the function
- [x] Only inside the `if` (and any `else`)
- [ ] Only on the next line
- [ ] Anywhere in the file

<details><summary>Answer</summary>

Only inside the `if` (and any `else`).
</details>

**14.** How many loop keywords does Go have?

- [x] One: `for`
- [ ] Two: `for` and `while`
- [ ] Three: `for`, `while` and `do`
- [ ] None: you use `range`

<details><summary>Answer</summary>

One: `for`. It also does the job of `while`.
</details>

**15.** How do you loop over every element of a slice with its index?

- [ ] `for item in items { }`
- [ ] `items.forEach((item, i) => { })`
- [ ] `for (i = 0; item of items) { }`
- [x] `for i, item := range items { }`

<details><summary>Answer</summary>

`for i, item := range items { }`.
</details>

### Type conversion

**16.** How do you convert a string to a byte slice?

- [ ] `s.bytes()`
- [ ] `byte[](s)`
- [x] `[]byte(s)`
- [ ] `bytes(s)`

<details><summary>Answer</summary>

`[]byte(s)`. The reverse is `string(b)`.
</details>

## Types

### Structs

**17.** What happens to fields you leave out of a struct literal?

- [ ] They're `nil`
- [x] They get their zero values
- [ ] It's a compile error
- [ ] They're `undefined`

<details><summary>Answer</summary>

They get their zero values.
</details>

### Struct tags

**18.** What does the struct tag `json:"created_at"` do?

- [x] Sets the field's JSON key to `created_at`
- [ ] Sets the database column name
- [ ] Fills the field with the current time
- [ ] Makes the field required

<details><summary>Answer</summary>

Sets the field's JSON key to `created_at`.
</details>

**19.** Which struct tag keeps a field out of JSON completely?

- [ ] `json:"omitempty"`
- [ ] `json:"hidden"`
- [ ] `json:""`
- [x] `json:"-"`

<details><summary>Answer</summary>

`json:"-"`.
</details>

### Anonymous structs and embedding

**20.** What is an anonymous struct, and when is it handy?

- [ ] A struct with no fields; handy as a placeholder
- [ ] A struct from another package
- [x] A struct with no name; handy when you need it once
- [ ] A struct that can hold any type

<details><summary>Answer</summary>

A struct type with no name, declared where it's used. Handy for a request body that doesn't match a model exactly.
</details>

**21.** What does embedding a type in a struct do?

- [ ] Copies the type's code into the file
- [x] Pulls its fields and methods into the outer struct
- [ ] Makes the outer struct inherit from it
- [ ] Hides the type's fields from JSON

<details><summary>Answer</summary>

Pulls its fields and methods into the outer struct, so `payload.Name` works when `payload` embeds `models.User`.
</details>

**22.** Is embedding the same as inheritance?

- [x] No; the outer struct contains it, it isn't one
- [ ] Yes; it's Go's word for inheritance
- [ ] Yes, but only for methods
- [ ] No; embedding only copies fields

<details><summary>Answer</summary>

No. The outer struct *contains* the embedded value; it isn't one.
</details>

**22a.** `payload` embeds `models.User`. What are the two ways to read the user's `Name`?

- [ ] `payload.Name` or `payload->User.Name`
- [x] `payload.Name` or `payload.User.Name`
- [ ] `payload.User.Name` only
- [ ] `payload.Name` only

<details><summary>Answer</summary>

`payload.Name` (short form, the usual one) or `payload.User.Name` (full form: the embedded field is named after its type).
</details>

### Pointers

**23.** What do `&x` and `*T` mean?

- [ ] `&x`: value of `x`. `*T`: a copy of `T`
- [ ] `&x`: pointer type. `*T`: address of `T`
- [ ] `&x`: both `x`s. `*T`: every `T`
- [x] `&x`: address of `x`. `*T`: the type pointer to `T`

<details><summary>Answer</summary>

`&x` is the address of `x`. `*T` is the type "pointer to T".
</details>

**24.** Why does `Decode(&user)` take a pointer?

- [ ] Because pointers are faster to send
- [ ] Because `user` might be `nil`
- [x] So it can fill in the original `user`
- [ ] So it can validate `user`

<details><summary>Answer</summary>

So it can fill in the original `user`. Without a pointer it would get a copy.
</details>

**25.** Why use `*int` instead of `int` for an optional field?

- [ ] `*int` holds bigger numbers
- [x] `nil` can mean missing; `0` might be a real value
- [ ] `int` can't be sent in JSON
- [ ] `*int` uses less memory

<details><summary>Answer</summary>

`nil` means "missing", which an `int` can't express: its zero value `0` might be a real value.
</details>

**25a.** `p` is a pointer. What does `*p` give you?

- [ ] The address of `p`
- [ ] A copy of the pointer
- [x] The value at the address `p` holds
- [ ] The type of `p`

<details><summary>Answer</summary>

The value at the address `p` holds.
</details>

**25b.** What happens if you read `*p` when `p` is `nil`?

- [ ] It returns `0`
- [ ] It returns `nil`
- [x] It panics
- [ ] It's a compile error

<details><summary>Answer</summary>

It panics.
</details>

### Arrays

**26.** What's the difference between an array and a slice?

- [x] An array has a fixed length; a slice can grow
- [ ] An array can grow; a slice has a fixed length
- [ ] An array holds one type; a slice holds any type
- [ ] They're the same; slice is the newer name

<details><summary>Answer</summary>

An array has a fixed length that's part of its type (`[3]int`). A slice can grow, and points to an underlying array.
</details>

**27.** What happens when you assign one array to another variable and change the copy?

- [ ] The original changes too; they share data
- [ ] It's a compile error
- [ ] Both become `nil`
- [x] The original doesn't change; arrays are copied

<details><summary>Answer</summary>

The original doesn't change. Assigning an array copies every element.
</details>

### Slices

**28.** Why must you write `cards = append(cards, c)`, not just `append(cards, c)`?

- [ ] `append` only works with `=`
- [ ] Otherwise `c` is added twice
- [x] `append` returns the new slice; it doesn't update the old variable
- [ ] `append` deletes the old slice

<details><summary>Answer</summary>

`append` returns the new slice; the old variable isn't updated.
</details>

**29.** How do you get a slice's length?

- [ ] `s.length`
- [x] `len(s)`
- [ ] `s.len()`
- [ ] `size(s)`

<details><summary>Answer</summary>

`len(s)`.
</details>

**30.** What are a slice's length and capacity?

- [x] Length: elements in use. Capacity: room before `append` grows it
- [ ] Length: room in the array. Capacity: elements in use
- [ ] Length: bytes used. Capacity: maximum bytes
- [ ] They're always the same number

<details><summary>Answer</summary>

Length: how many elements are in use. Capacity: how much room the underlying array has before `append` must make a bigger one.
</details>

**31.** Does slicing (`s[1:3]`) copy the elements? What follows from that?

- [ ] Yes; the new slice is independent
- [ ] Yes, but only the first element
- [ ] No; slicing deletes the other elements
- [x] No; it shares the array, so changes show in the original

<details><summary>Answer</summary>

No. The new slice points into the same underlying array, so changing an element through it changes the original too. Use `copy` or `slices.Clone` for an independent copy.
</details>

### Maps

**32.** How do you create an empty, usable map?

- [ ] `var m map[string]int`
- [ ] `new map[string]int`
- [x] `make(map[string]int)`
- [ ] `map[string]int(nil)`

<details><summary>Answer</summary>

`make(map[string]int)` or `map[string]int{}`.
</details>

**33.** How do you tell a missing key apart from a zero value?

- [ ] `m[key] == nil`
- [x] `v, ok := m[key]`; `ok` is false if missing
- [ ] `m.has(key)`
- [ ] `v := m[key]`; `v` is `undefined` if missing

<details><summary>Answer</summary>

`v, ok := m[key]`: `ok` is `false` if the key is missing.
</details>

**34.** What happens if you write to a `nil` map?

- [x] It panics
- [ ] The map is created automatically
- [ ] It's ignored
- [ ] It's a compile error

<details><summary>Answer</summary>

It panics.
</details>

**35.** How do you remove a key from a map?

- [ ] `m[key] = nil`
- [ ] `m.remove(key)`
- [ ] `remove(m, key)`
- [x] `delete(m, key)`

<details><summary>Answer</summary>

`delete(m, key)`. Nothing happens if the key isn't there.
</details>

**36.** In what order does `for k, v := range m` visit a map's keys?

- [ ] The order the keys were added
- [ ] Alphabetical order
- [x] A random order that can change between runs
- [ ] Reverse order

<details><summary>Answer</summary>

A random order that can change between runs. Sort the keys first if order matters.
</details>

### Choosing a data structure

**37.** When would you use a struct, a slice and a map?

- [ ] Struct: ordered list. Slice: lookup by key. Map: named fields
- [x] Struct: named fields. Slice: ordered list. Map: lookup by key
- [ ] Struct: lookup by key. Slice: named fields. Map: ordered list
- [ ] Struct: any data. Slice: numbers only. Map: strings only

<details><summary>Answer</summary>

Struct: a fixed set of named fields. Slice: an ordered list that can grow. Map: looking values up by a key.
</details>

**38.** Which are copied when assigned, and which are shared: struct, array, slice, map?

- [x] Copied: struct, array. Shared: slice, map
- [ ] Copied: slice, map. Shared: struct, array
- [ ] All four are copied
- [ ] All four are shared

<details><summary>Answer</summary>

Structs and arrays are copied. Slices and maps are shared: both variables see the same data.
</details>

### `any`

**39.** What does `any` mean?

- [ ] A value that can be `nil`
- [ ] Any number type
- [ ] A value whose type is decided at compile time
- [x] A value of any type; short for `interface{}`

<details><summary>Answer</summary>

A value of any type. It's short for `interface{}`.
</details>

### Type assertions

**40.** What does `v, ok := x.(models.User)` do?

- [ ] Converts `x` into a `models.User`
- [ ] Checks if `v` is a `models.User`
- [x] Checks if `x` holds a `models.User`; gives it back and `ok`
- [ ] Creates a new `models.User` from `x`

<details><summary>Answer</summary>

Checks whether `x` (an `any` or interface value) holds a `models.User`. If so, `v` is that user and `ok` is `true`; if not, `v` is the zero value and `ok` is `false`.
</details>

**41.** Why use the `, ok` form of a type assertion?

- [ ] Without it, `v` is always `nil`
- [x] Without it, a wrong type panics
- [ ] Without it, the code doesn't compile
- [ ] It makes the assertion faster

<details><summary>Answer</summary>

Without it, a value of the wrong type causes a panic.
</details>

## Functions and methods

### Parameters that share a type

**42.** What does `func f(tokenHash, userID string)` mean?

- [x] Both are strings
- [ ] `tokenHash` has no type
- [ ] `tokenHash` is `any`
- [ ] `userID` is optional

<details><summary>Answer</summary>

Both `tokenHash` and `userID` are strings. Neighbouring parameters of the same type can share one type name.
</details>

### Methods and receivers

**43.** What is a receiver?

- [ ] A function's return value
- [ ] A channel that gets messages
- [ ] The parameter a callback receives
- [x] The value a method is attached to

<details><summary>Answer</summary>

The value a method is attached to, written before the method name: `func (h *UserHandler) CreateUser(...)`.
</details>

**44.** What's the difference between a pointer receiver and a value receiver?

- [ ] Pointer receiver: a copy. Value receiver: the original
- [ ] Pointer receivers can't change fields
- [x] Pointer receiver: the original. Value receiver: a copy
- [ ] Value receivers are only for interfaces

<details><summary>Answer</summary>

A pointer receiver works on the original; a value receiver gets a copy.
</details>

### Functions as values

**45.** In `mux.HandleFunc("POST /users", userHandler.CreateUser)`, why is there no `()` after `CreateUser`?

- [ ] Methods never use `()` in Go
- [x] It's passed as a value to call later; `()` would call it now
- [ ] The router adds the `()` itself
- [ ] `()` is only for functions with parameters

<details><summary>Answer</summary>

The method is passed as a value for the router to call later; `()` would call it now.
</details>

**46.** What is a closure?

- [x] A function that uses variables from the code around it
- [ ] A function that closes a connection
- [ ] A function that runs when the program ends
- [ ] A function with no return value

<details><summary>Answer</summary>

An anonymous function that uses variables from the code around it.
</details>

### `defer`

**47.** When does a deferred call run?

- [ ] At the end of the current line
- [ ] When the program exits
- [ ] After a 1-second delay
- [x] When the surrounding function returns

<details><summary>Answer</summary>

When the surrounding function returns, however it returns.
</details>

**48.** Where should `defer rows.Close()` go?

- [ ] At the end of the function
- [ ] Before the query runs
- [x] Right after successfully getting `rows`
- [ ] Inside the `for rows.Next()` loop

<details><summary>Answer</summary>

Right after successfully getting `rows`, so cleanup can't be forgotten.
</details>

## Interfaces

**49.** How does a type satisfy an interface in Go?

- [ ] With the `implements` keyword
- [x] By having all its methods; it's automatic
- [ ] By embedding the interface
- [ ] By having the same fields

<details><summary>Answer</summary>

By having all its methods. It's automatic; there's no `implements` keyword.
</details>

**50.** Why do handlers depend on the `Store` interface rather than `*PostgresStore`?

- [x] So tests can pass in a fake store
- [ ] Because interfaces are faster
- [ ] Because `*PostgresStore` is unexported
- [ ] So handlers can write SQL

<details><summary>Answer</summary>

So tests can pass in a fake store instead of a real database.
</details>

## Errors

**51.** What does a `nil` error mean?

- [ ] An unknown error
- [ ] The error was ignored
- [ ] The function panicked
- [x] Success

<details><summary>Answer</summary>

Success.
</details>

**52.** Go has no exceptions. How do errors travel up the call chain?

- [ ] They're thrown and caught with `try/catch`
- [ ] They're logged and the program stops
- [x] They're returned, and each caller checks `if err != nil`
- [ ] They're sent on a channel

<details><summary>Answer</summary>

They're returned, and each caller checks `if err != nil`.
</details>

**53.** What does `%w` do in `fmt.Errorf`?

- [ ] Writes the error to the log
- [x] Wraps the error so `errors.Is`/`As` can still find it
- [ ] Formats the error as a warning
- [ ] Replaces the error with a new one

<details><summary>Answer</summary>

Wraps the original error, adding context while keeping it findable with `errors.Is` and `errors.As`.
</details>

**54.** What's the difference between `errors.Is` and `errors.As`?

- [x] `Is`: a specific value. `As`: a specific type, put in your variable
- [ ] `Is`: a specific type. `As`: a specific value
- [ ] `Is`: compares messages. `As`: converts to a string
- [ ] They're the same; `As` is newer

<details><summary>Answer</summary>

`errors.Is` checks for a specific error **value** (e.g. `ErrEmailTaken`). `errors.As` checks if an error is a specific **type** (e.g. `*pq.Error`) and, if so, puts it in your variable.
</details>

## Testing

**55.** How must test files and test functions be named?

- [ ] Files start `test_`; functions end `Test`
- [ ] Files are in a `tests/` folder; any function name
- [ ] Files end `.test.go`; functions start `it`
- [x] Files end `_test.go`; functions start `Test`

<details><summary>Answer</summary>

Files end in `_test.go`; functions start with `Test` and take `t *testing.T`.
</details>

**56.** What's the difference between `t.Fatalf` and `t.Errorf`?

- [ ] `Fatalf` keeps going; `Errorf` stops the test
- [ ] `Fatalf` stops every test; `Errorf` only logs
- [x] `Fatalf` stops the test; `Errorf` keeps going
- [ ] They're the same

<details><summary>Answer</summary>

`t.Fatalf` fails and stops the test; `t.Errorf` fails and keeps going.
</details>

**57.** What does `./...` mean in `go test ./...`?

- [ ] Only the current folder
- [x] This folder and every folder below it
- [ ] The parent folder
- [ ] Every Go file on the computer

<details><summary>Answer</summary>

This folder and every folder below it.
</details>
