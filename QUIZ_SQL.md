# SQL Quiz

Recall questions on `LEARNING_SQL.md`. Try answering before opening each answer.

## Tables

### `CREATE TABLE`

**1.** What does `CREATE TABLE IF NOT EXISTS` do differently from `CREATE TABLE`?

- [x] Does nothing if the table exists, instead of failing
- [ ] Replaces the table if it exists
- [ ] Adds any missing columns to the table
- [ ] Creates the table only in tests

<details><summary>Answer</summary>

Example: the API runs `CREATE TABLE IF NOT EXISTS card_movements …` on every startup. The first `go run .` creates it; **every later run does nothing**.

Plain `CREATE TABLE` would fail on the second run with `relation "card_movements" already exists`.
</details>

**1a.** You add a foreign key to a `CREATE TABLE IF NOT EXISTS`. Why doesn't your existing database get it?

- [ ] Foreign keys only apply to new rows
- [x] `IF NOT EXISTS` never changes a table that's already there
- [ ] Postgres needs restarting first
- [ ] Foreign keys can only be added with `CREATE INDEX`

<details><summary>Answer</summary>

Example: your database already has `card_movements` from last month. You add `REFERENCES cards(id)` to the `CREATE TABLE`, restart, and **nothing changes**, because the table exists so the whole statement is skipped.

You run **`ALTER TABLE card_movements ADD FOREIGN KEY …`** once by hand. Step 6's migrations fix this properly.
</details>

### Column types

**2.** What's the difference between `VARCHAR(100)` and `TEXT`?

- [ ] `VARCHAR(100)`: exactly 100 characters. `TEXT`: up to 255
- [ ] `VARCHAR(100)`: 100 words. `TEXT`: 100 lines
- [ ] They're the same; `TEXT` is older
- [x] `VARCHAR(100)`: up to 100 characters. `TEXT`: no limit

<details><summary>Answer</summary>

Example: `title VARCHAR(100)` **rejects a 101-character title**, so the database enforces the limit. `description TEXT` accepts any length.

Use `VARCHAR(n)` when there's a real maximum; `TEXT` when there isn't. They're equally fast in Postgres.
</details>

**3.** Which type stores a moment in time, with its time zone?

- [ ] `TIMESTAMP`
- [ ] `DATE`
- [x] `TIMESTAMPTZ`
- [ ] `TIME`

<details><summary>Answer</summary>

Example: Ada in London creates a card at 09:00. Bo in New York sees **04:00 his time**: the same moment. `TIMESTAMPTZ` stores the exact moment, not a clock reading.

Plain `TIMESTAMP` would store "09:00" with no zone, and Bo couldn't tell when it really happened.
</details>

### Constraints

**4.** What two rules does `PRIMARY KEY` combine, and what does it come with?

- [ ] Unique and indexed; comes with a default
- [x] Unique and not null; comes with an index
- [ ] Not null and a foreign key; comes with an index
- [ ] Unique and auto-numbered; comes with a trigger

<details><summary>Answer</summary>

Example: `users.id UUID PRIMARY KEY` means no two users share an `id`, and **no user can have no `id`**.

The index (`users_pkey`) makes `WHERE id = $1` a quick lookup, which matters because nearly every query finds rows by ID.
</details>

**5.** What does a `CHECK` constraint do? Give the example from this project.

- [x] Rejects rows that make it false, e.g. `CHECK (position >= 0)`
- [ ] Checks a value exists in another table, e.g. `CHECK (column_id)`
- [ ] Warns about bad rows but keeps them
- [ ] Fixes bad values, e.g. turns `-1` into `0`

<details><summary>Answer</summary>

Example: `columns.position INTEGER CHECK (position >= 0)`. Inserting `2` works; inserting `-1` **fails with `23514`**, even if the Go code has a bug.

It's a rule you write yourself, so the database can enforce anything a true/false expression can say.
</details>

**6.** What does `DEFAULT NOW()` do?

- [ ] Updates the time on every change
- [ ] Makes the column required
- [ ] Sets the time zone
- [x] Fills in the current time when an insert gives no value

<details><summary>Answer</summary>

Example: `INSERT INTO boards (id, title, owner_id) VALUES (…)` doesn't mention `created_at`, so **Postgres fills in the current time**.

The Go code never sets `created_at`; it gets the value back with `RETURNING created_at`.
</details>

**7.** What do Postgres error codes starting with `23` mean, and which one is a unique violation?

- [ ] Connection errors; `23505` is a timeout
- [ ] Constraint violations; `23503` is unique
- [x] Constraint violations; `23505` is unique
- [ ] Syntax errors; `23505` is a typo

<details><summary>Answer</summary>

Postgres groups its codes by class: **`23` is the "a constraint was broken" class**.

Example: two sign-ups with `ada@example.com` give **`23505`** (unique). Deleting a user who has sessions without `ON DELETE CASCADE` gives `23503` (foreign key).
</details>

**8.** How should Go code refer to a Postgres error code like `23505`?

- [ ] As the bare string `"23505"`
- [x] Through a named constant, e.g. `uniqueViolation`
- [ ] As the number `23505`
- [ ] By matching the error message text

<details><summary>Answer</summary>

`if pqErr.Code == "23505"` leaves the reader guessing what `23505` is. **`if pqErr.Code == uniqueViolation` explains itself.**

It also means a typo like `"2305"` can only happen in one place.
</details>

### Quoting names

**9.** Why is the cards column written `"order"`, with double quotes?

- [x] `ORDER` is an SQL keyword
- [ ] Column names must be quoted
- [ ] It makes the column case-insensitive
- [ ] It marks the column as sorted

<details><summary>Answer</summary>

Example: `SELECT id, order FROM cards` fails, because Postgres reads `order` as the start of `ORDER BY`. **`"order"` tells it this is a name**, not the keyword.

That's why new tables use names like `position` instead.
</details>

## Queries

### `SELECT`

**10.** Without `ORDER BY`, what order do rows come back in?

- [x] No guaranteed order
- [ ] The order rows were inserted
- [ ] Sorted by primary key
- [ ] Alphabetical

<details><summary>Answer</summary>

Example: `SELECT * FROM boards` might return Ada's boards oldest first today, and **a different order after an update**, because Postgres returns rows in whatever order is quickest for it.

If the order matters, say so: `ORDER BY created_at DESC`.
</details>

**11.** What does `WHERE` do?

- [ ] Picks which columns to return
- [ ] Sorts the rows
- [ ] Names the table
- [x] Filters which rows the query uses

<details><summary>Answer</summary>

Example: `SELECT * FROM cards WHERE column_id = $1` **only returns cards in that column**. On an `UPDATE` or `DELETE`, `WHERE` decides which rows change.

Without it, every row is used.
</details>

### `INSERT`, `UPDATE`, `DELETE`

**12.** What does `RETURNING` do on an `INSERT`?

- [ ] Undoes the insert
- [ ] Returns the number of rows
- [x] Sends back columns from the inserted row
- [ ] Returns the SQL that ran

<details><summary>Answer</summary>

Example: `INSERT INTO boards (id, title, owner_id) VALUES ($1, $2, $3) RETURNING created_at`. The database sets `created_at`, and **`RETURNING` sends it back in the same query**.

Without it, you'd need a second `SELECT` to find out what the database filled in.
</details>

**13.** What happens if you run `DELETE FROM users` without a `WHERE`?

- [ ] Nothing happens
- [x] Every row is deleted
- [ ] Postgres asks for confirmation
- [ ] Only the first row is deleted

<details><summary>Answer</summary>

Example: you mean `DELETE FROM users WHERE id = $1` but forget the `WHERE`. **Every user is deleted**, and through `ON DELETE CASCADE`, all their boards and sessions too.

Postgres doesn't warn you. It's why the store always includes the `WHERE`.
</details>

### Placeholders

**14.** What are `$1`, `$2` for?

- [x] Placeholders for values sent separately
- [ ] Column numbers
- [ ] Variables stored in Postgres
- [ ] Row numbers in the result

<details><summary>Answer</summary>

Example: `QueryRowContext(ctx, "SELECT … WHERE email = $1", email)`. **The SQL text and the value travel separately**, and Postgres puts the value in place of `$1`.

`$2` is the second argument, `$3` the third, and so on.
</details>

**15.** What attack do placeholders prevent, and how?

- [ ] Brute force; queries are slowed down
- [ ] XSS; the value is escaped as HTML
- [ ] CSRF; the query needs a token
- [x] SQL injection; the value is only ever data

<details><summary>Answer</summary>

Example: if code pasted the email into the SQL, the input `' OR '1'='1` would make it `WHERE email = '' OR '1'='1'`, which **matches every user**.

As a placeholder value, that input is just a strange email address that matches nobody.
</details>

### `JOIN`

**16.** What does a `JOIN` do?

- [ ] Adds the rows of one table to the end of another
- [ ] Merges two tables into a new one permanently
- [x] Combines rows from two tables, matched on a condition
- [ ] Copies columns from one table to another

<details><summary>Answer</summary>

Example: `sessions` has `user_id`, and you want the user's name. `JOIN users u ON u.id = s.user_id` gives **one row with columns from both tables**: `token_hash` and `name`.

That's `GetSessionUser`: one query instead of two.
</details>

**17.** In `FROM sessions s`, what is `s`?

- [ ] A second `sessions` table
- [x] An alias for `sessions`
- [ ] The plural of session
- [ ] A schema name

<details><summary>Answer</summary>

Example: `FROM sessions s JOIN users u ON u.id = s.user_id`. Both tables have `created_at`, so **`s.created_at` and `u.created_at` say which one you mean**.

Aliases also keep long queries short to read.
</details>

### Writes that check another table

**18.** What does `INSERT … SELECT … WHERE EXISTS (…)` do when the `EXISTS` is false?

- [x] Inserts nothing; `RowsAffected()` is 0
- [ ] Inserts an empty row
- [ ] Returns an error
- [ ] Inserts the row anyway

<details><summary>Answer</summary>

Example: Bo tries to add a column to Ada's board. `EXISTS (… WHERE owner_id = 'bo')` is false, so the `SELECT` gives no rows and **there's nothing to insert**.

The store sees `RowsAffected() == 0` and the handler returns `404`. One statement checks and writes.
</details>

**19.** What does `EXISTS (subquery)` return, and why is it often written `SELECT 1`?

- [ ] The number of rows; `1` limits it to one
- [ ] The first row; `1` means the first column
- [ ] True if exactly one row; `1` is that count
- [x] True if any row is returned; the `1` doesn't matter

<details><summary>Answer</summary>

Example: `EXISTS (SELECT 1 FROM boards WHERE id = $2 AND owner_id = $5)` is **true if Ada owns that board**, false if not.

`1` is just something to select. `SELECT title` would work the same; only whether a row exists matters.
</details>

**20.** How do you join another table into an `UPDATE` or a `DELETE` in Postgres?

- [ ] `UPDATE … JOIN` and `DELETE … JOIN`
- [ ] You can't; use two queries
- [x] `UPDATE … FROM` and `DELETE … USING`
- [ ] `UPDATE … WITH` and `DELETE … WITH`

<details><summary>Answer</summary>

Example: `DELETE FROM columns c USING boards b WHERE c.id = $1 AND b.id = c.board_id AND b.owner_id = $2`.

**`USING boards b` brings in the board**, so the `WHERE` can check its owner. Bo's ID matches nothing, so nothing is deleted.
</details>

## Performance

### Indexes

**21.** What is an index for?

- [ ] Numbering the rows in a table
- [x] Finding rows by a value without reading every row
- [ ] Stopping duplicate values
- [ ] Linking two tables together

<details><summary>Answer</summary>

Example: `WHERE id = $1` on a million cards. Without an index, Postgres checks **all 1,000,000 rows**. With the `cards_pkey` index, it checks about **20 entries**, because the index is sorted.

The bigger the table, the bigger the difference.
</details>

**22.** Does a foreign key column get an index automatically?

- [x] No; you add one with `CREATE INDEX`
- [ ] Yes, always
- [ ] Only if it's `NOT NULL`
- [ ] Only with `ON DELETE CASCADE`

<details><summary>Answer</summary>

Example: `cards.column_id REFERENCES columns(id)` has **no index**, so `WHERE column_id = $1` reads every card.

Postgres only creates indexes for `PRIMARY KEY` and `UNIQUE`. You add `CREATE INDEX cards_column_id_idx ON cards (column_id)` yourself (step 6).
</details>

**23.** What's the trade-off of adding an index?

- [ ] Faster writes; slower reads
- [ ] Faster everything; no downside
- [ ] Less storage; slower reads
- [x] Faster reads; slower writes and more storage

<details><summary>Answer</summary>

Example: an index on `cards.column_id` makes listing a column's cards fast. But **every `INSERT` of a card now also updates the index**, and it takes disk space.

So index columns you search by, not every column.
</details>

### `EXPLAIN`

**24.** What does `EXPLAIN` show?

- [ ] What a query's result means
- [ ] Why a query failed
- [x] How Postgres plans to run a query
- [ ] Which user ran a query

<details><summary>Answer</summary>

Example: `EXPLAIN SELECT * FROM cards WHERE column_id = '…'` shows **`Seq Scan on cards`**, meaning Postgres would read every row.

It's how you find out why a query is slow, without running it.
</details>

**25.** In a query plan, which is slow on a big table: an index scan or a sequential scan?

- [ ] An index scan
- [x] A sequential scan (`Seq Scan`)
- [ ] They're the same speed
- [ ] Neither; plans don't affect speed

<details><summary>Answer</summary>

Example: with a million cards, `Seq Scan` **reads all million rows**. An `Index Scan` goes straight to the few it needs.

On a 3-row table, a `Seq Scan` is fine, and Postgres may even prefer it.
</details>

**26.** In a query plan, what's the difference between `Index Cond` and `Filter`?

- [x] `Index Cond`: answered by an index. `Filter`: checked row by row
- [ ] `Index Cond`: checked row by row. `Filter`: answered by an index
- [ ] `Index Cond`: a `JOIN`. `Filter`: a `WHERE`
- [ ] They're the same

<details><summary>Answer</summary>

Example:

- `Index Cond: (id = '…')`: Postgres **used the index** to jump to the row.
- `Filter: (column_id = '…')`: Postgres **read each row and tested it**, which means no index helped.

Seeing `Filter` on a big table is a hint you need an index.
</details>

**27.** What does `EXPLAIN ANALYZE` do that `EXPLAIN` doesn't, and why be careful with it?

- [ ] It shows the plan in more detail without running it
- [ ] It runs the query in a test database
- [ ] It rolls the query back automatically
- [x] It really runs the query, so a `DELETE` really deletes

<details><summary>Answer</summary>

Example: `EXPLAIN ANALYZE DELETE FROM cards WHERE …` **deletes the cards for real**, then shows how long it took.

Wrap it in `BEGIN; … ROLLBACK;` to measure a write without keeping its changes.
</details>

## Go's `database/sql`

**28.** When do you use `QueryRowContext`, `QueryContext` and `ExecContext`?

- [ ] `QueryRow`: writes. `Query`: one row. `Exec`: many rows
- [ ] `QueryRow`: many rows. `Query`: writes. `Exec`: one row
- [x] `QueryRow`: one row. `Query`: many rows. `Exec`: no rows back
- [ ] They're the same; the names are just style

<details><summary>Answer</summary>

Examples from the store:

- `GetCardByID` uses **`QueryRowContext`**: one card.
- `GetCards` uses **`QueryContext`**: many cards, read in a `rows.Next()` loop.
- `DeleteCard` uses **`ExecContext`**: nothing comes back except `RowsAffected()`.
</details>

**29.** What does `QueryRowContext(...).Scan` return when no row matches?

- [ ] `nil`
- [x] `sql.ErrNoRows`
- [ ] An empty struct and `nil`
- [ ] `0`

<details><summary>Answer</summary>

Example: `GetCardByID` asks for card `abc`, which doesn't exist. The query works, but **there's nothing to scan**, so `Scan` returns `sql.ErrNoRows`.

The store checks `err == sql.ErrNoRows` and treats it as "not found", not as a failure.
</details>

**30.** Why does `Scan` take pointers, like `Scan(&u.ID, &u.Name)`?

- [x] So it can fill in the original variables
- [ ] Because pointers are faster
- [ ] So it can check for `NULL`
- [ ] Because columns are pointers in Postgres

<details><summary>Answer</summary>

Go passes values by copy. `Scan(u.ID)` would get **a copy of the empty `ID`** and couldn't change yours.

`Scan(&u.ID, &u.Name)` gives it the addresses, so it can write the row's values into your variables. Same idea as `Decode(&user)`.
</details>

## Transactions and locking

### `BEGIN`, `COMMIT`, `ROLLBACK`

**31.** Without `BEGIN`, when does an `UPDATE` commit?

- [ ] When the connection closes
- [x] Straight away: each statement is its own transaction
- [ ] At the next `COMMIT`
- [ ] When the next statement runs

<details><summary>Answer</summary>

Example: `UPDATE cards …` then `INSERT INTO card_movements …`, with no `BEGIN`. The `UPDATE` is **saved the moment it finishes**, so if the `INSERT` fails, the card has moved with no history row.

`BEGIN` makes both wait for one `COMMIT`.
</details>

### `SELECT … FOR UPDATE`

**32.** Why should a `FOR UPDATE` query avoid joining other tables?

- [ ] Joins can't be used with `FOR UPDATE`
- [ ] It would lock every joined table
- [ ] Joins are too slow inside a transaction
- [x] After waiting, only the locked row is re-read, so joined rows can be stale

<details><summary>Answer</summary>

Example: `MoveCard` waited while another move changed `c1` from `todo` to `doing`. Postgres **re-read the card (`doing`) but kept the old joined column (`todo`)**, so the `EXISTS` check failed and the card looked missing (#8).

Lock the one row on its own, then check the rest in a second query.
</details>

### `version = version + 1`

**33.** In `UPDATE cards SET version = version + 1 WHERE id = $1 AND version = $2`, what does `version = $2` do?

- [x] Only matches if the row is still at the version the client read
- [ ] Sets the version to `$2`
- [ ] Adds `$2` to the version
- [ ] Locks the row until `$2` changes

<details><summary>Answer</summary>

Example: Ada read version 4, but Bo saved first, so the row is at version 5. **`version = 4` matches nothing**, `RowsAffected()` is 0, and Ada gets `409`.

If nobody saved in between, it matches, and `SET version = version + 1` makes it 5.
</details>
