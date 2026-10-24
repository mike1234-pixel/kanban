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

- `CREATE TABLE IF NOT EXISTS` does nothing if the table is already there (used for `card_movements` on startup).
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

Combines rows from two tables in one query, matching on a condition:

```sql
SELECT u.id, u.name, u.email
FROM sessions s
JOIN users u ON u.id = s.user_id
WHERE s.token_hash = $1
```

- `s` and `u` are **aliases**: short names, so `s.user_id` means `sessions.user_id`.
- `JOIN` (an inner join) keeps only rows that match on both sides.
- It's usually the same link a foreign key protects.

## Performance

### Indexes

A lookup structure (like a book's index) that finds rows without scanning the whole table.

- `PRIMARY KEY` and `UNIQUE` create one automatically.
- Columns you filter on often (`WHERE …`) are candidates for one: `CREATE INDEX … ON columns (board_id)`.
- Trade-off: faster reads, slightly slower writes, more storage.

### `EXPLAIN`

Shows how Postgres plans to run a query, without running it.

```
Index Scan using sessions_pkey on sessions   ← uses an index: good
Seq Scan on cards                            ← reads every row: slow on big tables
```

## Go's `database/sql`

| Method | Use for | Returns |
| --- | --- | --- |
| `QueryRowContext` | Exactly one row expected | A row to `Scan`; `sql.ErrNoRows` if none |
| `QueryContext` | Many rows | `rows`: loop with `rows.Next()`, then check `rows.Err()`, and `defer rows.Close()` |
| `ExecContext` | Writes with no rows back | A result: `RowsAffected()` |

- `Scan(&a, &b)` copies the columns, in order, into Go variables. It needs pointers to fill them in.

Coming later: transactions (`BeginTx`, `COMMIT`, `ROLLBACK`, step 4) and `SELECT … FOR UPDATE` (step 5).
