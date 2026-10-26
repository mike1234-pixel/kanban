# SQL Quiz

Recall questions on `LEARNING_SQL.md`. Try answering before opening each answer.

## Tables

### `CREATE TABLE`

**1.** What does `CREATE TABLE IF NOT EXISTS` do differently from `CREATE TABLE`?

<details><summary>Answer</summary>

It does nothing if the table already exists, instead of failing.
</details>

### Column types

**2.** What's the difference between `VARCHAR(100)` and `TEXT`?

<details><summary>Answer</summary>

`VARCHAR(100)` allows up to 100 characters; `TEXT` has no length limit.
</details>

**3.** Which type stores a moment in time, with its time zone?

<details><summary>Answer</summary>

`TIMESTAMPTZ`.
</details>

### Constraints

**4.** What two rules does `PRIMARY KEY` combine, and what does it come with?

<details><summary>Answer</summary>

Unique and not null. It comes with an index.
</details>

**5.** What does a `CHECK` constraint do? Give the example from this project.

<details><summary>Answer</summary>

Rejects any row that makes the expression false, e.g. `CHECK (position >= 0)`.
</details>

**6.** What does `DEFAULT NOW()` do?

<details><summary>Answer</summary>

Fills in the current time when an insert doesn't give a value.
</details>

**7.** What do Postgres error codes starting with `23` mean, and which one is a unique violation?

<details><summary>Answer</summary>

`23…` codes are all constraint violations. `23505` is a unique violation (`23503` is a foreign key violation).
</details>

**8.** How should Go code refer to a Postgres error code like `23505`?

<details><summary>Answer</summary>

Through a named constant, e.g. `uniqueViolation = "23505"`, never as a bare string.
</details>

### Quoting names

**9.** Why is the cards column written `"order"`, with double quotes?

<details><summary>Answer</summary>

`ORDER` is an SQL keyword; quotes let it be used as a column name.
</details>

## Queries

### `SELECT`

**10.** Without `ORDER BY`, what order do rows come back in?

<details><summary>Answer</summary>

No guaranteed order.
</details>

**11.** What does `WHERE` do?

<details><summary>Answer</summary>

Filters which rows the query uses.
</details>

### `INSERT`, `UPDATE`, `DELETE`

**12.** What does `RETURNING` do on an `INSERT`?

<details><summary>Answer</summary>

Sends back columns from the inserted row, such as values the database filled in (`created_at`).
</details>

**13.** What happens if you run `DELETE FROM users` without a `WHERE`?

<details><summary>Answer</summary>

Every row in the table is deleted.
</details>

### Placeholders

**14.** What are `$1`, `$2` for?

<details><summary>Answer</summary>

Placeholders: values sent separately from the SQL text, filled in by the database.
</details>

**15.** What attack do placeholders prevent, and how?

<details><summary>Answer</summary>

SQL injection. The value is only ever treated as data, never as part of the SQL.
</details>

### `JOIN`

**16.** What does a `JOIN` do?

<details><summary>Answer</summary>

Combines rows from two tables in one query, matching them on a condition (e.g. `u.id = s.user_id`).
</details>

**17.** In `FROM sessions s`, what is `s`?

<details><summary>Answer</summary>

An alias: a short name for the table, so `s.user_id` means `sessions.user_id`.
</details>

### Writes that check another table

**18.** What does `INSERT … SELECT … WHERE EXISTS (…)` do when the `EXISTS` is false?

<details><summary>Answer</summary>

Inserts nothing: the `SELECT` produces no rows. `RowsAffected()` is 0.
</details>

**19.** What does `EXISTS (subquery)` return, and why is it often written `SELECT 1`?

<details><summary>Answer</summary>

True if the subquery returns at least one row. The selected value doesn't matter, only whether a row exists, so `1` is a placeholder.
</details>

**20.** How do you join another table into an `UPDATE` or a `DELETE` in Postgres?

<details><summary>Answer</summary>

`UPDATE … FROM other_table WHERE …` and `DELETE FROM … USING other_table WHERE …`.
</details>

## Performance

### Indexes

**21.** What is an index for?

<details><summary>Answer</summary>

Finding rows by a column's value without reading every row. It's a sorted copy of that column's values, each with its row's location.
</details>

**22.** Does a foreign key column get an index automatically?

<details><summary>Answer</summary>

No. `PRIMARY KEY` and `UNIQUE` do; a foreign key like `cards.column_id` doesn't, so you add one with `CREATE INDEX`.
</details>

**23.** What's the trade-off of adding an index?

<details><summary>Answer</summary>

Faster reads, but slightly slower writes and more storage.
</details>

### `EXPLAIN`

**24.** What does `EXPLAIN` show?

<details><summary>Answer</summary>

How Postgres plans to run a query, e.g. whether it uses an index or scans every row.
</details>

**25.** In a query plan, which is slow on a big table: an index scan or a sequential scan?

<details><summary>Answer</summary>

A sequential scan (`Seq Scan`): it reads every row.
</details>

**26.** In a query plan, what's the difference between `Index Cond` and `Filter`?

<details><summary>Answer</summary>

`Index Cond` is answered by an index. `Filter` is checked on each row, one by one.
</details>

**27.** What does `EXPLAIN ANALYZE` do that `EXPLAIN` doesn't, and why be careful with it?

<details><summary>Answer</summary>

It actually runs the query and shows real times. On a `DELETE` or `UPDATE` it really changes the data.
</details>

## Go's `database/sql`

**28.** When do you use `QueryRowContext`, `QueryContext` and `ExecContext`?

<details><summary>Answer</summary>

`QueryRowContext`: one row. `QueryContext`: many rows. `ExecContext`: writes with no rows back.
</details>

**29.** What does `QueryRowContext(...).Scan` return when no row matches?

<details><summary>Answer</summary>

`sql.ErrNoRows`.
</details>

**30.** Why does `Scan` take pointers, like `Scan(&u.ID, &u.Name)`?

<details><summary>Answer</summary>

So it can fill in the original variables, not copies.
</details>
