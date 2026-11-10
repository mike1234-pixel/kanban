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

It does nothing if the table already exists, instead of failing.
</details>

**1a.** You add a foreign key to a `CREATE TABLE IF NOT EXISTS`. Why doesn't your existing database get it?

- [ ] Foreign keys only apply to new rows
- [x] `IF NOT EXISTS` never changes a table that's already there
- [ ] Postgres needs restarting first
- [ ] Foreign keys can only be added with `CREATE INDEX`

<details><summary>Answer</summary>

`IF NOT EXISTS` skips a table that already exists, so it's never changed. An existing table needs **`ALTER TABLE … ADD FOREIGN KEY`**.
</details>

### Column types

**2.** What's the difference between `VARCHAR(100)` and `TEXT`?

- [ ] `VARCHAR(100)`: exactly 100 characters. `TEXT`: up to 255
- [ ] `VARCHAR(100)`: 100 words. `TEXT`: 100 lines
- [ ] They're the same; `TEXT` is older
- [x] `VARCHAR(100)`: up to 100 characters. `TEXT`: no limit

<details><summary>Answer</summary>

`VARCHAR(100)` allows up to 100 characters; `TEXT` has no length limit.
</details>

**3.** Which type stores a moment in time, with its time zone?

- [ ] `TIMESTAMP`
- [ ] `DATE`
- [x] `TIMESTAMPTZ`
- [ ] `TIME`

<details><summary>Answer</summary>

`TIMESTAMPTZ`.
</details>

### Constraints

**4.** What two rules does `PRIMARY KEY` combine, and what does it come with?

- [ ] Unique and indexed; comes with a default
- [x] Unique and not null; comes with an index
- [ ] Not null and a foreign key; comes with an index
- [ ] Unique and auto-numbered; comes with a trigger

<details><summary>Answer</summary>

Unique and not null. It comes with an index.
</details>

**5.** What does a `CHECK` constraint do? Give the example from this project.

- [x] Rejects rows that make it false, e.g. `CHECK (position >= 0)`
- [ ] Checks a value exists in another table, e.g. `CHECK (column_id)`
- [ ] Warns about bad rows but keeps them
- [ ] Fixes bad values, e.g. turns `-1` into `0`

<details><summary>Answer</summary>

Rejects any row that makes the expression false, e.g. `CHECK (position >= 0)`.
</details>

**6.** What does `DEFAULT NOW()` do?

- [ ] Updates the time on every change
- [ ] Makes the column required
- [ ] Sets the time zone
- [x] Fills in the current time when an insert gives no value

<details><summary>Answer</summary>

Fills in the current time when an insert doesn't give a value.
</details>

**7.** What do Postgres error codes starting with `23` mean, and which one is a unique violation?

- [ ] Connection errors; `23505` is a timeout
- [ ] Constraint violations; `23503` is unique
- [x] Constraint violations; `23505` is unique
- [ ] Syntax errors; `23505` is a typo

<details><summary>Answer</summary>

`23…` codes are all constraint violations. `23505` is a unique violation (`23503` is a foreign key violation).
</details>

**8.** How should Go code refer to a Postgres error code like `23505`?

- [ ] As the bare string `"23505"`
- [x] Through a named constant, e.g. `uniqueViolation`
- [ ] As the number `23505`
- [ ] By matching the error message text

<details><summary>Answer</summary>

Through a named constant, e.g. `uniqueViolation = "23505"`, never as a bare string.
</details>

### Quoting names

**9.** Why is the cards column written `"order"`, with double quotes?

- [x] `ORDER` is an SQL keyword
- [ ] Column names must be quoted
- [ ] It makes the column case-insensitive
- [ ] It marks the column as sorted

<details><summary>Answer</summary>

`ORDER` is an SQL keyword; quotes let it be used as a column name.
</details>

## Queries

### `SELECT`

**10.** Without `ORDER BY`, what order do rows come back in?

- [x] No guaranteed order
- [ ] The order rows were inserted
- [ ] Sorted by primary key
- [ ] Alphabetical

<details><summary>Answer</summary>

No guaranteed order.
</details>

**11.** What does `WHERE` do?

- [ ] Picks which columns to return
- [ ] Sorts the rows
- [ ] Names the table
- [x] Filters which rows the query uses

<details><summary>Answer</summary>

Filters which rows the query uses.
</details>

### `INSERT`, `UPDATE`, `DELETE`

**12.** What does `RETURNING` do on an `INSERT`?

- [ ] Undoes the insert
- [ ] Returns the number of rows
- [x] Sends back columns from the inserted row
- [ ] Returns the SQL that ran

<details><summary>Answer</summary>

Sends back columns from the inserted row, such as values the database filled in (`created_at`).
</details>

**13.** What happens if you run `DELETE FROM users` without a `WHERE`?

- [ ] Nothing happens
- [x] Every row is deleted
- [ ] Postgres asks for confirmation
- [ ] Only the first row is deleted

<details><summary>Answer</summary>

Every row in the table is deleted.
</details>

### Placeholders

**14.** What are `$1`, `$2` for?

- [x] Placeholders for values sent separately
- [ ] Column numbers
- [ ] Variables stored in Postgres
- [ ] Row numbers in the result

<details><summary>Answer</summary>

Placeholders: values sent separately from the SQL text, filled in by the database.
</details>

**15.** What attack do placeholders prevent, and how?

- [ ] Brute force; queries are slowed down
- [ ] XSS; the value is escaped as HTML
- [ ] CSRF; the query needs a token
- [x] SQL injection; the value is only ever data

<details><summary>Answer</summary>

SQL injection. The value is only ever treated as data, never as part of the SQL.
</details>

### `JOIN`

**16.** What does a `JOIN` do?

- [ ] Adds the rows of one table to the end of another
- [ ] Merges two tables into a new one permanently
- [x] Combines rows from two tables, matched on a condition
- [ ] Copies columns from one table to another

<details><summary>Answer</summary>

Combines rows from two tables in one query, matching them on a condition (e.g. `u.id = s.user_id`).
</details>

**17.** In `FROM sessions s`, what is `s`?

- [ ] A second `sessions` table
- [x] An alias for `sessions`
- [ ] The plural of session
- [ ] A schema name

<details><summary>Answer</summary>

An alias: a short name for the table, so `s.user_id` means `sessions.user_id`.
</details>

### Writes that check another table

**18.** What does `INSERT … SELECT … WHERE EXISTS (…)` do when the `EXISTS` is false?

- [x] Inserts nothing; `RowsAffected()` is 0
- [ ] Inserts an empty row
- [ ] Returns an error
- [ ] Inserts the row anyway

<details><summary>Answer</summary>

Inserts nothing: the `SELECT` produces no rows. `RowsAffected()` is 0.
</details>

**19.** What does `EXISTS (subquery)` return, and why is it often written `SELECT 1`?

- [ ] The number of rows; `1` limits it to one
- [ ] The first row; `1` means the first column
- [ ] True if exactly one row; `1` is that count
- [x] True if any row is returned; the `1` doesn't matter

<details><summary>Answer</summary>

True if the subquery returns at least one row. The selected value doesn't matter, only whether a row exists, so `1` is a placeholder.
</details>

**20.** How do you join another table into an `UPDATE` or a `DELETE` in Postgres?

- [ ] `UPDATE … JOIN` and `DELETE … JOIN`
- [ ] You can't; use two queries
- [x] `UPDATE … FROM` and `DELETE … USING`
- [ ] `UPDATE … WITH` and `DELETE … WITH`

<details><summary>Answer</summary>

`UPDATE … FROM other_table WHERE …` and `DELETE FROM … USING other_table WHERE …`.
</details>

## Performance

### Indexes

**21.** What is an index for?

- [ ] Numbering the rows in a table
- [x] Finding rows by a value without reading every row
- [ ] Stopping duplicate values
- [ ] Linking two tables together

<details><summary>Answer</summary>

Finding rows by a column's value without reading every row. It's a sorted copy of that column's values, each with its row's location.
</details>

**22.** Does a foreign key column get an index automatically?

- [x] No; you add one with `CREATE INDEX`
- [ ] Yes, always
- [ ] Only if it's `NOT NULL`
- [ ] Only with `ON DELETE CASCADE`

<details><summary>Answer</summary>

No. `PRIMARY KEY` and `UNIQUE` do; a foreign key like `cards.column_id` doesn't, so you add one with `CREATE INDEX`.
</details>

**23.** What's the trade-off of adding an index?

- [ ] Faster writes; slower reads
- [ ] Faster everything; no downside
- [ ] Less storage; slower reads
- [x] Faster reads; slower writes and more storage

<details><summary>Answer</summary>

Faster reads, but slightly slower writes and more storage.
</details>

### `EXPLAIN`

**24.** What does `EXPLAIN` show?

- [ ] What a query's result means
- [ ] Why a query failed
- [x] How Postgres plans to run a query
- [ ] Which user ran a query

<details><summary>Answer</summary>

How Postgres plans to run a query, e.g. whether it uses an index or scans every row.
</details>

**25.** In a query plan, which is slow on a big table: an index scan or a sequential scan?

- [ ] An index scan
- [x] A sequential scan (`Seq Scan`)
- [ ] They're the same speed
- [ ] Neither; plans don't affect speed

<details><summary>Answer</summary>

A sequential scan (`Seq Scan`): it reads every row.
</details>

**26.** In a query plan, what's the difference between `Index Cond` and `Filter`?

- [x] `Index Cond`: answered by an index. `Filter`: checked row by row
- [ ] `Index Cond`: checked row by row. `Filter`: answered by an index
- [ ] `Index Cond`: a `JOIN`. `Filter`: a `WHERE`
- [ ] They're the same

<details><summary>Answer</summary>

`Index Cond` is answered by an index. `Filter` is checked on each row, one by one.
</details>

**27.** What does `EXPLAIN ANALYZE` do that `EXPLAIN` doesn't, and why be careful with it?

- [ ] It shows the plan in more detail without running it
- [ ] It runs the query in a test database
- [ ] It rolls the query back automatically
- [x] It really runs the query, so a `DELETE` really deletes

<details><summary>Answer</summary>

It actually runs the query and shows real times. On a `DELETE` or `UPDATE` it really changes the data.
</details>

## Go's `database/sql`

**28.** When do you use `QueryRowContext`, `QueryContext` and `ExecContext`?

- [ ] `QueryRow`: writes. `Query`: one row. `Exec`: many rows
- [ ] `QueryRow`: many rows. `Query`: writes. `Exec`: one row
- [x] `QueryRow`: one row. `Query`: many rows. `Exec`: no rows back
- [ ] They're the same; the names are just style

<details><summary>Answer</summary>

`QueryRowContext`: one row. `QueryContext`: many rows. `ExecContext`: writes with no rows back.
</details>

**29.** What does `QueryRowContext(...).Scan` return when no row matches?

- [ ] `nil`
- [x] `sql.ErrNoRows`
- [ ] An empty struct and `nil`
- [ ] `0`

<details><summary>Answer</summary>

`sql.ErrNoRows`.
</details>

**30.** Why does `Scan` take pointers, like `Scan(&u.ID, &u.Name)`?

- [x] So it can fill in the original variables
- [ ] Because pointers are faster
- [ ] So it can check for `NULL`
- [ ] Because columns are pointers in Postgres

<details><summary>Answer</summary>

So it can fill in the original variables, not copies.
</details>
