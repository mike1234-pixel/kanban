# Learning SQL

Quick reference for the SQL and Postgres used in this project. Examples come from the codebase and `README.md`. Backend concepts live in `LEARNING.md`; Go syntax in `LEARNING_GO.md`.

## Tables

### `CREATE TABLE`

Defines a table: its columns, their types and their constraints.

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email VARCHAR(254) NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

- If a table with that name already exists:

  | Statement | Result |
  | --- | --- |
  | `CREATE TABLE users …` | Fails: `ERROR: relation "users" already exists` |
  | `CREATE TABLE IF NOT EXISTS users …` | Does nothing, no error |

  If it doesn't exist, both create it. A name can never be used by two tables.
- The API uses `IF NOT EXISTS` for `card_movements` on startup; without it, the second `go run .` would fail.
- `ALTER TABLE users ADD COLUMN …` changes an existing table.

### Column types

| Type | Holds | Used for |
| --- | --- | --- |
| `UUID` | A 128-bit unique ID | Every `id` |
| `VARCHAR(n)` | Text up to `n` characters | `title`, `email` |
| `TEXT` | Text of any length | `description`, `password_hash` |
| `INTEGER` | Whole numbers | `position`, `"order"` |
| `TIMESTAMPTZ` | A moment in time, with time zone | `created_at`, `expires_at` |

### Constraints

Rules the database enforces on every write, whatever the Go code does.

| Constraint | Means | Violation |
| --- | --- | --- |
| `PRIMARY KEY` | Unique and not null; the row's identity. Comes with an index. | `23505` |
| `NOT NULL` | Must have a value | `23502` |
| `UNIQUE` | No two rows share the value | `23505` |
| `REFERENCES t(col)` | Foreign key: must match a row in `t` | `23503` |
| `CHECK (position >= 0)` | Must make the expression true | `23514` |
| `DEFAULT NOW()` | Value used when the insert doesn't give one | (not a rule) |

- `ON DELETE CASCADE` on a foreign key: deleting the parent row also deletes the rows that reference it (board → columns → cards; user → sessions).
- **Error codes starting with `23` are all constraint violations.** Worth knowing by heart: `23505` (unique), the most common. Look the rest up in the Postgres docs ("Appendix A: PostgreSQL Error Codes").
- In Go, never write a code bare: give it a named constant, e.g. `uniqueViolation = "23505"` in `store/postgres.go`.

### Quoting names

`"order"` has double quotes because `ORDER` is an SQL keyword (`ORDER BY`). Quotes let a keyword be used as a column name. Avoid it in new tables: pick a name like `position`.

## Queries

### `SELECT`

```sql
SELECT id, title, created_at FROM boards ORDER BY created_at DESC
```

- `WHERE` filters rows: `WHERE board_id = $1`.
- `ORDER BY … ASC | DESC` sorts. Without it, the row order isn't guaranteed.
- `AND` / `OR` combine conditions: `WHERE token_hash = $1 AND expires_at > NOW()`.
- `NOW()` is the current time, according to the database.

### `INSERT`, `UPDATE`, `DELETE`

```sql
INSERT INTO boards (id, title) VALUES ($1, $2) RETURNING created_at
UPDATE users SET name = $1, email = $2 WHERE id = $3
DELETE FROM users WHERE id = $1
```

- `RETURNING` sends back columns from the written row, e.g. values the database filled in, like `created_at`.
- An `UPDATE` or `DELETE` **without `WHERE`** changes every row in the table.

### Placeholders: `$1`, `$2`

Values are sent separately from the SQL text, never pasted into it:

```go
s.db.QueryRowContext(ctx, `SELECT … FROM users WHERE email = $1`, email)
```

- This prevents **SQL injection**: if `email` were pasted in, an input like `' OR '1'='1` would become part of the query. As a placeholder value, it's only ever data.
- Postgres numbers them `$1, $2…`; some other databases use `?`.

### `JOIN`

`JOIN` builds one result from two tables. For each row of the first table, Postgres finds the rows of the second table that satisfy the `ON` condition, and outputs one row containing the columns of both.

**Example.** Two tables:

```
sessions                     users
token_hash | user_id         id | name
-----------+--------         ---+-----
a3e6…      | 1               1  | Ada
9f2c…      | 2               2  | Bo
77b1…      | 1
```

The query:

```sql
SELECT s.token_hash, u.name
FROM sessions s
JOIN users u ON u.id = s.user_id
```

What Postgres does:

1. `FROM sessions s`: take each row of `sessions`.
2. `JOIN users u ON u.id = s.user_id`: for that row, find the `users` row whose `id` equals the session's `user_id`.
3. Output one row with the selected columns from both.

Result:

```
token_hash | name
-----------+-----
a3e6…      | Ada     ← session a3e6… has user_id 1 → users row 1
9f2c…      | Bo      ← session 9f2c… has user_id 2 → users row 2
77b1…      | Ada     ← session 77b1… has user_id 1 → users row 1
```

Adding `WHERE s.token_hash = $1` keeps only the row for one token. That is `GetSessionUser`: one query returns the user for a session. Without `JOIN` it would take two queries: read `user_id` from `sessions`, then read that user from `users`.

- `s` and `u` are **aliases** for the table names. `s.user_id` means `sessions.user_id`. They are needed when both tables have a column with the same name (`created_at`).
- A row with no match is left out of the result. This is called an inner join, and it's what `JOIN` means on its own.
- The `ON` condition is usually the foreign key: `sessions.user_id REFERENCES users(id)`.

### Writes that check another table

These let one statement both check ownership and write, so there's no separate "check, then write" step.

**`INSERT … SELECT … WHERE EXISTS`** inserts a row only if a condition holds:

```sql
INSERT INTO columns (id, board_id, title, position)
SELECT $1, $2, $3, $4
WHERE EXISTS (SELECT 1 FROM boards WHERE id = $2 AND owner_id = $5)
```

The `SELECT` produces one row (the values) if the `EXISTS` is true, and no rows otherwise. No rows means nothing is inserted, and `RowsAffected()` is 0.

**`EXISTS (subquery)`** is true if the subquery returns at least one row. `SELECT 1` is a convention: the value doesn't matter, only whether a row exists. `SELECT EXISTS (…)` returns that true/false as a value you can `Scan` into a Go `bool`.

**`UPDATE … FROM`** and **`DELETE … USING`** join another table into the write:

```sql
UPDATE columns c SET title = $1
FROM boards b
WHERE c.id = $2 AND b.id = c.board_id AND b.owner_id = $3

DELETE FROM columns c
USING boards b
WHERE c.id = $1 AND b.id = c.board_id AND b.owner_id = $2
```

Only rows that match the join and the `WHERE` are updated or deleted.

## Performance

### Indexes

An index is a separate, sorted structure that Postgres keeps alongside a table. It stores the values of one column (or several), each with the location of its row. Because it's sorted, Postgres can find a value in a few steps instead of checking every row.

**Without an index**, `WHERE column_id = $1` makes Postgres read every row of `cards` and test each one:

```
row 1: column_id = b1ee…?  yes → keep
row 2: column_id = b1ee…?  no
row 3: column_id = b1ee…?  no
… and so on, for every row in the table
```

**With an index on `id`**, `WHERE id = $1` looks the value up in the sorted index, gets the row's location, and reads only that row.

- With 3 rows the difference is nothing. With 1,000,000 rows, a full read checks 1,000,000 rows; an index lookup checks about 20 entries.
- `PRIMARY KEY` and `UNIQUE` create an index automatically: `cards_pkey`, `users_email_key`.
- **A foreign key does not.** `cards.column_id` has no index, so `WHERE column_id = $1` reads every row. You add one with `CREATE INDEX cards_column_id_idx ON cards (column_id);`.
- Add indexes for columns you often filter, join or sort on.
- Cost: every `INSERT`, `UPDATE` and `DELETE` must also update each index on the table, and indexes take disk space. So don't index every column.

### `EXPLAIN`

`EXPLAIN` before a query shows the plan Postgres would use to run it, without running it. Real output from this project's database:

```sql
EXPLAIN SELECT * FROM cards WHERE id = 'c2ee…';
```

```
Index Scan using cards_pkey on cards  (cost=0.14..8.16 rows=1 width=878)
  Index Cond: (id = 'c2ee…'::uuid)
```

`id` has an index (`cards_pkey`), so Postgres uses it. `Index Cond` is the condition the index answers.

```sql
EXPLAIN SELECT * FROM cards WHERE column_id = 'b1ee…';
```

```
Seq Scan on cards  (cost=0.00..11.12 rows=1 width=878)
  Filter: (column_id = 'b1ee…'::uuid)
```

`column_id` has no index, so Postgres reads every row (`Seq Scan`, a sequential scan) and tests each one (`Filter`).

How to read a plan:

| Part | Meaning |
| --- | --- |
| `Index Scan using <index>` | Finds rows through that index |
| `Seq Scan` | Reads every row of the table |
| `Index Cond` | Condition answered by the index |
| `Filter` | Condition checked row by row |
| `cost=0.14..8.16` | Postgres's estimate of the work (startup..total), in its own units. Compare plans with it; it isn't milliseconds. |
| `rows=1` | Estimated number of rows returned |

- `EXPLAIN ANALYZE` actually runs the query and adds real times. Careful: on a `DELETE` or `UPDATE` it really changes the data.
- On a tiny table, Postgres may choose a `Seq Scan` even when an index exists, because reading a few rows directly is cheaper. Plans only become meaningful with realistic amounts of data.

## Go's `database/sql`

| Method | Use for | Returns |
| --- | --- | --- |
| `QueryRowContext` | Exactly one row expected | A row to `Scan`; `sql.ErrNoRows` if none |
| `QueryContext` | Many rows | `rows`: loop with `rows.Next()`, then check `rows.Err()`, and `defer rows.Close()` |
| `ExecContext` | Writes with no rows back | A result: `RowsAffected()` |

- `Scan(&a, &b)` copies the columns, in order, into Go variables. It needs pointers to fill them in.

Coming later: transactions (`BeginTx`, `COMMIT`, `ROLLBACK`, step 4) and `SELECT … FOR UPDATE` (step 5).
