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

`go run .` looks for **`package main` with a `func main()`** and starts there.

Example: `backend/main.go` is `package main`. Every other folder (`store`, `handlers`) has its own name and can only be imported, not run.
</details>

**2.** What does `_` in front of an import mean?

- [ ] Import it but keep its names private
- [ ] Import it only in tests
- [ ] Import it under the name `_`, used as `_.Open`
- [x] Import it only for its side effects, e.g. a driver

<details><summary>Answer</summary>

Example: `import _ "github.com/lib/pq"`. You never write `pq.Something`, but importing it **registers the `postgres` driver**, so `sql.Open("postgres", …)` works.

Without the `_`, Go refuses to compile, because unused imports are an error.
</details>

### Exported vs unexported

**3.** What makes a name visible outside its package?

- [ ] Adding the `export` keyword
- [ ] Starting it with an underscore
- [x] Starting it with a capital letter
- [ ] Listing it in `go.mod`

<details><summary>Answer</summary>

There's no `public`/`private` keyword: **the first letter decides**.

Example: in package `store`, `ErrEmailTaken` can be used as `store.ErrEmailTaken` from `handlers`. A lowercase `uniqueViolation` can only be used inside `store`.
</details>

**4.** Why can't `encoding/json` see a lowercase struct field?

- [ ] Lowercase fields can't hold strings
- [x] It's unexported, and `encoding/json` is another package
- [ ] It needs a `json` struct tag first
- [ ] JSON keys must be capitalised

<details><summary>Answer</summary>

`encoding/json` reads your struct from outside your package, so it follows the same rule as any other package: **lowercase names are hidden**.

Example: `type Board struct { title string }` encodes as `{}`. No error, the field is just missing.
</details>

### Variables

**5.** What's the difference between `x := 5` and `var x int`?

- [x] `:=` declares and assigns; `var x int` declares with zero value `0`
- [ ] `:=` assigns only; `var x int` declares and assigns `5`
- [ ] `:=` makes a constant; `var x int` makes a variable
- [ ] They're identical; `var` is just older

<details><summary>Answer</summary>

Example: `count := 5` makes an `int` set to `5`, with **the type worked out from the value**. `var count int` makes an `int` set to `0`.

Use `:=` when you have a value now; `var` when you'll fill it in later, like `var c models.Card` before `Scan(&c.ID, …)`.
</details>

**6.** Where can't you use `:=`?

- [ ] Inside `for` loops
- [ ] Inside `if` statements
- [ ] Inside methods
- [x] At package level, outside functions

<details><summary>Answer</summary>

`:=` only works inside functions. **At package level, every declaration starts with a keyword** (`var`, `const`, `func`, `type`).

Example: `var dummyHash, _ = bcrypt.GenerateFromPassword(…)` in `auth.go` is at package level, so it uses `var`.
</details>

### Zero values

**7.** What is a zero value?

- [ ] The value `0`, for every type
- [ ] A variable that holds nothing, like `undefined`
- [x] The value Go gives a variable you don't set
- [ ] The first value a function returns

<details><summary>Answer</summary>

Example: `var n int` is `0`, not undefined. `var c models.Card` is a card with every field empty.

There's no `undefined` like in TypeScript: **every variable always has a usable value**.
</details>

**8.** What are the zero values of `int`, `string` and `bool`?

- [ ] `nil`, `nil` and `nil`
- [x] `0`, `""` and `false`
- [ ] `0`, `nil` and `false`
- [ ] `0`, `" "` and `nil`

<details><summary>Answer</summary>

Example: `var c models.Card` gives `c.Order == 0`, `c.Title == ""`.

That's why `GetCardByID` can return `models.Card{}` for "not found": the handler checks **`card.ID == ""`**.
</details>

**9.** What is the zero value of a pointer, slice, map or error?

- [x] `nil`
- [ ] `0`
- [ ] An empty value, e.g. `[]`
- [ ] `undefined`

<details><summary>Answer</summary>

These types **point to something**, and `nil` means they don't point to anything yet.

Example: `var err error` is `nil`, meaning no error. `var m map[string]int` is `nil`, so writing to it panics.
</details>

### Constants

**10.** What keyword declares a value that can never change?

- [x] `const`
- [ ] `let`
- [ ] `final`
- [ ] `var`

<details><summary>Answer</summary>

Example: `const uniqueViolation = "23505"`. Trying `uniqueViolation = "x"` **won't compile**.

It gives a magic value a name, so code reads `pqErr.Code == uniqueViolation` instead of `== "23505"`.
</details>

### Multiple return values

**11.** By convention, which return value is the error?

- [ ] The first one
- [ ] The second one
- [ ] Any; it's named `err`
- [x] The last one

<details><summary>Answer</summary>

Example: `func (s *PostgresStore) GetBoard(…) (models.Board, error)`.

Because it's always last, **anyone reading `x, err := f()` knows where to look** without checking the function.
</details>

**12.** How do you ignore a return value you don't need?

- [ ] Assign it to `nil`
- [ ] Leave it out of the left side
- [x] Assign it to `_`
- [ ] Assign it to `ignore`

<details><summary>Answer</summary>

Go won't compile if you declare a variable and never use it. **`_` throws the value away.**

Example: `res, _ := …` would ignore an error, which is rarely a good idea. `for _, card := range cards` ignores the index, which is common.
</details>

### `if` and `for`

**13.** In `if err := f(); err != nil { }`, where can you use `err`?

- [ ] Anywhere in the function
- [x] Only inside the `if` (and any `else`)
- [ ] Only on the next line
- [ ] Anywhere in the file

<details><summary>Answer</summary>

Example: `if err := tx.Commit(); err != nil { return err }`. After the closing `}`, `err` **no longer exists**.

That keeps a short-lived variable from leaking into the rest of the function and clashing with other `err`s.
</details>

**14.** How many loop keywords does Go have?

- [x] One: `for`
- [ ] Two: `for` and `while`
- [ ] Three: `for`, `while` and `do`
- [ ] None: you use `range`

<details><summary>Answer</summary>

`for` does every loop job:

- `for i := 0; i < 3; i++ { }`: a counting loop
- `for rows.Next() { }`: a `while` loop
- `for { }`: forever

**One keyword, three forms.**
</details>

**15.** How do you loop over every element of a slice with its index?

- [ ] `for item in items { }`
- [ ] `items.forEach((item, i) => { })`
- [ ] `for (i = 0; item of items) { }`
- [x] `for i, item := range items { }`

<details><summary>Answer</summary>

Example: for `cards := []string{"a", "b"}`, the loop runs with `i = 0, item = "a"` then `i = 1, item = "b"`.

`for _, item := range cards` when you don't need **the index**, which is the usual case.
</details>

### Type conversion

**16.** How do you convert a string to a byte slice?

- [ ] `s.bytes()`
- [ ] `byte[](s)`
- [x] `[]byte(s)`
- [ ] `bytes(s)`

<details><summary>Answer</summary>

Go never converts types for you; **you write the type as a function**.

Example: bcrypt works on bytes, so `bcrypt.GenerateFromPassword([]byte(payload.Password), …)`. `string(b)` turns bytes back into text.
</details>

## Types

### Structs

**17.** What happens to fields you leave out of a struct literal?

- [ ] They're `nil`
- [x] They get their zero values
- [ ] It's a compile error
- [ ] They're `undefined`

<details><summary>Answer</summary>

Example: `models.Card{Title: "Fix bug"}` gives `Description: ""` and `Order: 0`.

You **only write the fields you care about**; the rest are never undefined.
</details>

### Struct tags

**18.** What does the struct tag `json:"created_at"` do?

- [x] Sets the field's JSON key to `created_at`
- [ ] Sets the database column name
- [ ] Fills the field with the current time
- [ ] Makes the field required

<details><summary>Answer</summary>

Go names are `CreatedAt`, but JSON convention is `created_at`. **The tag maps one to the other**, both ways.

Example: `` CreatedAt time.Time `json:"created_at"` `` encodes as `"created_at": "2026-11-11T…"`, and decodes from it too.
</details>

**19.** Which struct tag keeps a field out of JSON completely?

- [ ] `json:"omitempty"`
- [ ] `json:"hidden"`
- [ ] `json:""`
- [x] `json:"-"`

<details><summary>Answer</summary>

Example: `` PasswordHash string `json:"-"` `` on `models.User`. Encoding **skips it**, so a response can never include the hash.

Decoding skips it too, so a client can't set it by sending `"PasswordHash"`.
</details>

### Anonymous structs and embedding

**20.** What is an anonymous struct, and when is it handy?

- [ ] A struct with no fields; handy as a placeholder
- [ ] A struct from another package
- [x] A struct with no name; handy when you need it once
- [ ] A struct that can hold any type

<details><summary>Answer</summary>

Example: `POST /users` needs a password, but `models.User` has no `Password` field. So the handler decodes into:

`` var payload struct { models.User; Password string `json:"password"` } ``

**A type you only need in one place doesn't need a name.**
</details>

**21.** What does embedding a type in a struct do?

- [ ] Copies the type's code into the file
- [x] Pulls its fields and methods into the outer struct
- [ ] Makes the outer struct inherit from it
- [ ] Hides the type's fields from JSON

<details><summary>Answer</summary>

Example: `payload` embeds `models.User`, so `payload.Name` and `payload.Email` work **as if they were declared on `payload`**.

It saves you copying every `User` field into the anonymous struct.
</details>

**22.** Is embedding the same as inheritance?

- [x] No; the outer struct contains it, it isn't one
- [ ] Yes; it's Go's word for inheritance
- [ ] Yes, but only for methods
- [ ] No; embedding only copies fields

<details><summary>Answer</summary>

Example: `payload` embeds `models.User`, but you **can't pass `payload` to a function that wants a `models.User`**. You pass `payload.User`.

It's composition: `payload` *has* a user inside it; it isn't one.
</details>

**22a.** `payload` embeds `models.User`. What are the two ways to read the user's `Name`?

- [ ] `payload.Name` or `payload->User.Name`
- [x] `payload.Name` or `payload.User.Name`
- [ ] `payload.User.Name` only
- [ ] `payload.Name` only

<details><summary>Answer</summary>

The embedded field's name **is its type's name**, `User`. So the full form is `payload.User.Name`.

`payload.Name` is the shortcut Go gives you, and it's what you'll usually write.
</details>

### Pointers

**23.** What do `&x` and `*T` mean?

- [ ] `&x`: value of `x`. `*T`: a copy of `T`
- [ ] `&x`: pointer type. `*T`: address of `T`
- [ ] `&x`: both `x`s. `*T`: every `T`
- [x] `&x`: address of `x`. `*T`: the type pointer to `T`

<details><summary>Answer</summary>

Example: `x := 5`, `p := &x`. `p` holds **the address where `x` lives**, and its type is `*int`.

The `*` in a type means "pointer to"; the `&` in front of a variable means "address of".
</details>

**24.** Why does `Decode(&user)` take a pointer?

- [ ] Because pointers are faster to send
- [ ] Because `user` might be `nil`
- [x] So it can fill in the original `user`
- [ ] So it can validate `user`

<details><summary>Answer</summary>

Go passes everything by copy. `Decode(user)` would fill **a copy that's thrown away**, leaving your `user` empty.

`Decode(&user)` gives it the address, so it writes into your variable. `Scan(&c.ID)` works the same way.
</details>

**25.** Why use `*int` instead of `int` for an optional field?

- [ ] `*int` holds bigger numbers
- [x] `nil` can mean missing; `0` might be a real value
- [ ] `int` can't be sent in JSON
- [ ] `*int` uses less memory

<details><summary>Answer</summary>

Example: `{"order": 0}` means "move to the top". `{}` means "don't change it". With `int`, **both decode to `0`** and look the same.

With `*int`, `{}` gives `nil` and `{"order": 0}` gives a pointer to `0`.
</details>

**25a.** `p` is a pointer. What does `*p` give you?

- [ ] The address of `p`
- [ ] A copy of the pointer
- [x] The value at the address `p` holds
- [ ] The type of `p`

<details><summary>Answer</summary>

Example: `x := 5`, `p := &x`. `*p` is **`5`**. Writing `*p = 10` changes `x` to `10`.

In a type, `*` means pointer; in front of a pointer variable, it means "go to the value".
</details>

**25b.** What happens if you read `*p` when `p` is `nil`?

- [ ] It returns `0`
- [ ] It returns `nil`
- [x] It panics
- [ ] It's a compile error

<details><summary>Answer</summary>

`nil` points to nothing, so **there's no value to read**.

Example: `var p *int`, then `fmt.Println(*p)` stops the program with `invalid memory address or nil pointer dereference`. Check `if p != nil` first.
</details>

### Arrays

**26.** What's the difference between an array and a slice?

- [x] An array has a fixed length; a slice can grow
- [ ] An array can grow; a slice has a fixed length
- [ ] An array holds one type; a slice holds any type
- [ ] They're the same; slice is the newer name

<details><summary>Answer</summary>

Example: `[3]int` holds **exactly 3**, and `[4]int` is a different type. `[]int` holds any number and grows with `append`.

You'll almost always use slices. Arrays mostly exist as what slices are built on.
</details>

**27.** What happens when you assign one array to another variable and change the copy?

- [ ] The original changes too; they share data
- [ ] It's a compile error
- [ ] Both become `nil`
- [x] The original doesn't change; arrays are copied

<details><summary>Answer</summary>

Example: `a := [3]int{1, 2, 3}`, `b := a`, `b[0] = 99`. `a` is still `[1 2 3]`: **`b` got its own copy** of every element.

Slices are different: `b := s` shares the same elements.
</details>

### Slices

**28.** Why must you write `cards = append(cards, c)`, not just `append(cards, c)`?

- [ ] `append` only works with `=`
- [ ] Otherwise `c` is added twice
- [x] `append` returns the new slice; it doesn't update the old variable
- [ ] `append` deletes the old slice

<details><summary>Answer</summary>

When the slice is full, `append` **makes a new, bigger array** and returns a slice pointing to it. Your old `cards` variable still points to the old one.

Writing `append(cards, c)` on its own won't compile, because Go rejects an unused result.
</details>

**29.** How do you get a slice's length?

- [ ] `s.length`
- [x] `len(s)`
- [ ] `s.len()`
- [ ] `size(s)`

<details><summary>Answer</summary>

Example: `cards := []string{"a", "b"}`, `len(cards)` is **`2`**.

It's a built-in function, not a method: `len(cards)`, not `cards.length` like in TypeScript.
</details>

**30.** What are a slice's length and capacity?

- [x] Length: elements in use. Capacity: room before `append` grows it
- [ ] Length: room in the array. Capacity: elements in use
- [ ] Length: bytes used. Capacity: maximum bytes
- [ ] They're always the same number

<details><summary>Answer</summary>

Example: `s := make([]int, 2, 5)` has 2 elements in use and **room for 3 more**. The next 3 `append`s reuse the array; the 4th makes a bigger one and copies everything.

Length is what you see; capacity is spare room.
</details>

**31.** Does slicing (`s[1:3]`) copy the elements? What follows from that?

- [ ] Yes; the new slice is independent
- [ ] Yes, but only the first element
- [ ] No; slicing deletes the other elements
- [x] No; it shares the array, so changes show in the original

<details><summary>Answer</summary>

Example: `a := []int{1, 2, 3}`, `b := a[0:2]`, `b[0] = 99`. Now **`a` is `[99 2 3]` too**.

A slice is a window onto an array, so two slices can share elements. `slices.Clone(a)` gives you a separate copy.
</details>

### Maps

**32.** How do you create an empty, usable map?

- [ ] `var m map[string]int`
- [ ] `new map[string]int`
- [x] `make(map[string]int)`
- [ ] `map[string]int(nil)`

<details><summary>Answer</summary>

`var m map[string]int` is `nil`, and **writing to it panics**. `make` creates the map so it can be written to.

Example: `ages := make(map[string]int)`, then `ages["Ada"] = 36` works.
</details>

**33.** How do you tell a missing key apart from a zero value?

- [ ] `m[key] == nil`
- [x] `v, ok := m[key]`; `ok` is false if missing
- [ ] `m.has(key)`
- [ ] `v := m[key]`; `v` is `undefined` if missing

<details><summary>Answer</summary>

Example: `ages := map[string]int{"Bo": 0}`. Both `ages["Bo"]` and `ages["Cy"]` give `0`.

`v, ok := ages["Cy"]` gives **`ok = false`**, so you know Cy isn't there, rather than being 0.
</details>

**34.** What happens if you write to a `nil` map?

- [x] It panics
- [ ] The map is created automatically
- [ ] It's ignored
- [ ] It's a compile error

<details><summary>Answer</summary>

A `nil` map has **no storage to write into**.

Example: `var ages map[string]int`, then `ages["Ada"] = 36` stops the program with `assignment to entry in nil map`. Reading a `nil` map is fine: it gives the zero value.
</details>

**35.** How do you remove a key from a map?

- [ ] `m[key] = nil`
- [ ] `m.remove(key)`
- [ ] `remove(m, key)`
- [x] `delete(m, key)`

<details><summary>Answer</summary>

Example: `delete(ages, "Ada")`. If `"Ada"` isn't there, **it does nothing**, so you don't need to check first.
</details>

**36.** In what order does `for k, v := range m` visit a map's keys?

- [ ] The order the keys were added
- [ ] Alphabetical order
- [x] A random order that can change between runs
- [ ] Reverse order

<details><summary>Answer</summary>

Go **deliberately randomises it**, so code can't rely on an order that might change.

Example: ranging over `{"Ada": 36, "Bo": 41}` might give Bo first one run and Ada first the next. Collect the keys, `sort` them, then loop over those.
</details>

### Choosing a data structure

**37.** When would you use a struct, a slice and a map?

- [ ] Struct: ordered list. Slice: lookup by key. Map: named fields
- [x] Struct: named fields. Slice: ordered list. Map: lookup by key
- [ ] Struct: lookup by key. Slice: named fields. Map: ordered list
- [ ] Struct: any data. Slice: numbers only. Map: strings only

<details><summary>Answer</summary>

Examples from this codebase:

- **Struct:** `models.Card` with `ID`, `Title`, `Order`.
- **Slice:** `[]models.Card` for a column's cards, in order.
- **Map:** `colMap[c.ID]` to find a column by its ID in one step, without looping.
</details>

**38.** Which are copied when assigned, and which are shared: struct, array, slice, map?

- [x] Copied: struct, array. Shared: slice, map
- [ ] Copied: slice, map. Shared: struct, array
- [ ] All four are copied
- [ ] All four are shared

<details><summary>Answer</summary>

Example: `b := a` with a struct gives **two separate cards**; changing `b.Title` leaves `a` alone. With a map, `b := a` gives **two names for one map**; `b["x"] = 1` shows in `a`.

That's why passing a map to a function lets it change your map.
</details>

### `any`

**39.** What does `any` mean?

- [ ] A value that can be `nil`
- [ ] Any number type
- [ ] A value whose type is decided at compile time
- [x] A value of any type; short for `interface{}`

<details><summary>Answer</summary>

Example: `var v any = "hello"`, then `v = 42`. Both work, because **`any` accepts every type**.

The cost: you have to check the type (a type assertion) before you can use the value. It's like `unknown` in TypeScript.
</details>

### Type assertions

**40.** What does `v, ok := x.(models.User)` do?

- [ ] Converts `x` into a `models.User`
- [ ] Checks if `v` is a `models.User`
- [x] Checks if `x` holds a `models.User`; gives it back and `ok`
- [ ] Creates a new `models.User` from `x`

<details><summary>Answer</summary>

Example: `userFromContext` gets `any` back from `ctx.Value(userKey)`. `user, ok := v.(models.User)` **checks it really is a user** before using it.

If it is, `ok` is `true` and `user` is Ada. If not, `ok` is `false` and `user` is an empty `models.User{}`.
</details>

**41.** Why use the `, ok` form of a type assertion?

- [ ] Without it, `v` is always `nil`
- [x] Without it, a wrong type panics
- [ ] Without it, the code doesn't compile
- [ ] It makes the assertion faster

<details><summary>Answer</summary>

Example: `n := v.(int)` when `v` holds `"hello"` **stops the program with a panic**.

`n, ok := v.(int)` gives `ok = false` instead, and the code can return an error like `401`.
</details>

## Functions and methods

### Parameters that share a type

**42.** What does `func f(tokenHash, userID string)` mean?

- [x] Both are strings
- [ ] `tokenHash` has no type
- [ ] `tokenHash` is `any`
- [ ] `userID` is optional

<details><summary>Answer</summary>

Example: `func (s *PostgresStore) CreateSession(ctx context.Context, tokenHash, userID string, …)`.

`tokenHash, userID string` is **shorthand for `tokenHash string, userID string`**: the type applies to every name before it.
</details>

### Methods and receivers

**43.** What is a receiver?

- [ ] A function's return value
- [ ] A channel that gets messages
- [ ] The parameter a callback receives
- [x] The value a method is attached to

<details><summary>Answer</summary>

Example: in `func (h *UserHandler) CreateUser(w, r)`, `h` is the receiver. Inside, **`h.store` reaches the handler's store**.

It's like `this` in a TypeScript class method, but named and written explicitly.
</details>

**44.** What's the difference between a pointer receiver and a value receiver?

- [ ] Pointer receiver: a copy. Value receiver: the original
- [ ] Pointer receivers can't change fields
- [x] Pointer receiver: the original. Value receiver: a copy
- [ ] Value receivers are only for interfaces

<details><summary>Answer</summary>

Example: `func (c Counter) Add()` changes **a copy**, so `c.n++` is lost when it returns. `func (c *Counter) Add()` changes **the original**.

Store and handler methods use pointer receivers, so all requests share one store.
</details>

### Functions as values

**45.** In `mux.HandleFunc("POST /users", userHandler.CreateUser)`, why is there no `()` after `CreateUser`?

- [ ] Methods never use `()` in Go
- [x] It's passed as a value to call later; `()` would call it now
- [ ] The router adds the `()` itself
- [ ] `()` is only for functions with parameters

<details><summary>Answer</summary>

Example: `mux.HandleFunc("POST /users", userHandler.CreateUser)` gives the router **the function itself**. The router calls it with `(w, r)` each time a request arrives.

`userHandler.CreateUser()` would try to run it right now, at startup, with no request.
</details>

**46.** What is a closure?

- [x] A function that uses variables from the code around it
- [ ] A function that closes a connection
- [ ] A function that runs when the program ends
- [ ] A function with no return value

<details><summary>Answer</summary>

Example: `protected := func(h http.HandlerFunc) http.Handler { return authHandler.RequireAuth(h) }`. The function **uses `authHandler` from `main`**, even though it's not a parameter.

Same as a TypeScript arrow function using a variable from outside it.
</details>

### `defer`

**47.** When does a deferred call run?

- [ ] At the end of the current line
- [ ] When the program exits
- [ ] After a 1-second delay
- [x] When the surrounding function returns

<details><summary>Answer</summary>

Example: `defer rows.Close()` near the top of `GetCards`. Whether the function returns normally or on an error, **`rows.Close()` runs just before it exits**.

If there are several defers, they run in reverse order.
</details>

**48.** Where should `defer rows.Close()` go?

- [ ] At the end of the function
- [ ] Before the query runs
- [x] Right after successfully getting `rows`
- [ ] Inside the `for rows.Next()` loop

<details><summary>Answer</summary>

Example:

`rows, err := s.db.QueryContext(…)`
`if err != nil { return nil, err }`
`defer rows.Close()`

**After the error check**, because if the query failed, `rows` is `nil` and there's nothing to close.
</details>

## Interfaces

**49.** How does a type satisfy an interface in Go?

- [ ] With the `implements` keyword
- [x] By having all its methods; it's automatic
- [ ] By embedding the interface
- [ ] By having the same fields

<details><summary>Answer</summary>

Example: `PostgresStore` has every method `Store` lists, so **it's a `Store` automatically**. The fake store in `handlers_test.go` is one too.

No `implements Store` line is needed; Go checks the methods.
</details>

**50.** Why do handlers depend on the `Store` interface rather than `*PostgresStore`?

- [x] So tests can pass in a fake store
- [ ] Because interfaces are faster
- [ ] Because `*PostgresStore` is unexported
- [ ] So handlers can write SQL

<details><summary>Answer</summary>

Example: `NewBoardHandler(store)` accepts anything with the `Store` methods. Tests pass a **fake that keeps boards in a map**, so they run without Postgres.

The handler can't tell the difference, which is the point.
</details>

## Errors

**51.** What does a `nil` error mean?

- [ ] An unknown error
- [ ] The error was ignored
- [ ] The function panicked
- [x] Success

<details><summary>Answer</summary>

Example: `board, err := s.GetBoard(…)`. If `err == nil`, **`board` is safe to use**. If not, don't trust `board`.

It's why almost every call is followed by `if err != nil`.
</details>

**52.** Go has no exceptions. How do errors travel up the call chain?

- [ ] They're thrown and caught with `try/catch`
- [ ] They're logged and the program stops
- [x] They're returned, and each caller checks `if err != nil`
- [ ] They're sent on a channel

<details><summary>Answer</summary>

Example: `MoveCard` returns an error from the `INSERT`. The handler **checks `if err != nil`** and returns `500`.

Nothing is thrown, so you can see every place an error can come from by reading the code.
</details>

**53.** What does `%w` do in `fmt.Errorf`?

- [ ] Writes the error to the log
- [x] Wraps the error so `errors.Is`/`As` can still find it
- [ ] Formats the error as a warning
- [ ] Replaces the error with a new one

<details><summary>Answer</summary>

Example: `fmt.Errorf("error during row iteration: %w", err)` gives a message with context, but **`errors.Is(err, …)` can still find the original** inside.

With `%v` instead, the original error would become plain text and couldn't be found.
</details>

**54.** What's the difference between `errors.Is` and `errors.As`?

- [x] `Is`: a specific value. `As`: a specific type, put in your variable
- [ ] `Is`: a specific type. `As`: a specific value
- [ ] `Is`: compares messages. `As`: converts to a string
- [ ] They're the same; `As` is newer

<details><summary>Answer</summary>

Example:

- `errors.Is(err, store.ErrEmailTaken)`: is this **that exact error**?
- `var pqErr *pq.Error; errors.As(err, &pqErr)`: is this **a Postgres error**? If so, `pqErr.Code` is now readable.
</details>

## Testing

**55.** How must test files and test functions be named?

- [ ] Files start `test_`; functions end `Test`
- [ ] Files are in a `tests/` folder; any function name
- [ ] Files end `.test.go`; functions start `it`
- [x] Files end `_test.go`; functions start `Test`

<details><summary>Answer</summary>

`go test` **only looks for these names**.

Example: `handlers_test.go` contains `func TestCreateBoard(t *testing.T)`. Name the function `CreateBoardTest` and it never runs, with no warning.
</details>

**56.** What's the difference between `t.Fatalf` and `t.Errorf`?

- [ ] `Fatalf` keeps going; `Errorf` stops the test
- [ ] `Fatalf` stops every test; `Errorf` only logs
- [x] `Fatalf` stops the test; `Errorf` keeps going
- [ ] They're the same

<details><summary>Answer</summary>

Example: if `CreateBoard` returns an error, the next checks would read a broken board, so use **`t.Fatalf` to stop**.

When checking several separate fields, `t.Errorf` reports every wrong one in a single run.
</details>

**57.** What does `./...` mean in `go test ./...`?

- [ ] Only the current folder
- [x] This folder and every folder below it
- [ ] The parent folder
- [ ] Every Go file on the computer

<details><summary>Answer</summary>

Example: from `backend/`, `go test ./...` runs the tests in `internal/handlers`, `internal/store` and **every other folder below**.

`go test .` would only run tests in `backend/` itself.
</details>
