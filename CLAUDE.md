# Kanban Backend (Go)

## Purpose

This is a learning project. The goal is to understand Go and what a backend is actually doing, not to ship features quickly. The Kanban app is the vehicle; each milestone introduces a new backend concept.

Optimise for my understanding, not for speed or completeness.

## How to help

- Explain the concept and the trade-offs before writing code.
- Prefer hints, guidance and small examples. Let me write the code unless I explicitly ask you to.
- When you do write code, keep it small and explain why it is idiomatic Go.
- Don't silently fix bugs. Tell me what is wrong and why, then let me fix it.
- When reviewing my code, point out non-idiomatic Go and explain the better pattern.
- Connect new work back to the layers it touches: TCP → HTTP → router → middleware → handler → store → Postgres → response.
- Stay within the current milestone. If something belongs to a later milestone, say so rather than building it now.

## Technical constraints

- Standard library first: `net/http`, `database/sql`, `encoding/json`, `context`, `log/slog`.
- No ORM and no web framework for now. The verbosity of `database/sql` is intentional; it shows what the database interaction really looks like. Tools like sqlc or pgx come later, once I understand what they abstract.
- Postgres is the database, accessed through the `lib/pq` driver.
- Pass `r.Context()` through to database calls (`QueryContext`, `ExecContext`, `BeginTx`).
- Don't introduce external infrastructure (Redis, RabbitMQ, etc.) until the roadmap reaches it. Learn the in-process version first (e.g. Go channels as a job queue).

## Architecture

- Request flow: HTTP request → router → middleware → handler → store → Postgres.
- **Handlers** (`backend/internal/handlers`): decode JSON, validate, call the store, choose a status code, encode the response.
- **Store** (`backend/internal/store`): all SQL lives here. Handlers never write SQL. Handlers depend on the `Store` interface in `store.go`, not on `PostgresStore`, so tests can use a fake.
- **Models** (`backend/internal/models`): shared structs with `json` and `validate` tags (go-playground/validator).
- Each layer has one responsibility; no giant handlers.

## API conventions

- Status codes: 200, 201, 204 (e.g. DELETE), 400, 401, 403, 404, 409, 500.
- Path params identify a resource (`/boards/123`); query params control how it is returned (`/boards?archived=false&page=2`).
- Responses set `Content-Type: application/json`.
- List endpoints return `[]`, never `null`, when empty.
- All errors use one format:

  ```json
  {
    "error": {
      "code": "BOARD_NOT_FOUND",
      "message": "Board not found"
    }
  }
  ```

- Database constraints (NOT NULL, UNIQUE, FOREIGN KEY) enforce what must never happen; application code checks what should happen. Map constraint violations to sensible API errors (e.g. UNIQUE → 409).

### Known gaps against these conventions

- Errors are still plain text via `http.Error`, not the JSON format above.
- `GET /cards` and `GET /columns?board_id=` return `null` when empty.
- Logging uses `log`/`fmt`, not `log/slog`.

## Roadmap

Work through these in order. Don't add several at once.

1. **CRUD** (done)
2. **HTTP fundamentals:** status codes, JSON encode/decode, validation, path/query params, consistent errors, Content-Type, request lifecycle, `context.Context`
3. **Authentication / authorization:** users, password hashing (bcrypt or Argon2id), sessions + cookies (sessions before JWTs), auth middleware, ownership checks (403 vs 404), database constraints
4. **Transactions:** a dedicated MoveCard operation that updates the card's column and writes a movement-history row in one transaction (`BeginTx` / `Commit` / `Rollback`)
5. **Concurrency:** race conditions (two users moving the same card), `SELECT ... FOR UPDATE`, isolation levels, optimistic vs pessimistic concurrency, goroutines, mutexes, channels
6. **WebSockets:** real-time board updates, hub/broadcast pattern, connection lifecycle, heartbeats, disconnect cleanup, backpressure
7. **Background workers:** in-process job queue (`chan Job` + worker goroutines) for audit events and notifications
8. **Caching / Redis**
9. **Observability:** structured logging with `log/slog`, request IDs, latency and query timing, error tracking, metrics
10. **Testing**
11. **Docker + deployment**
12. **Production architecture**

## Current milestone

<!-- Update this as you progress -->

Working on: step 2, HTTP fundamentals.

- Done: user CRUD endpoints (`/users`) with unique emails (409 on duplicates). No passwords or auth yet.
- Next: authentication → authorization → middleware → context.
- Note: `MoveCard` already runs in a transaction and writes to `card_movements` (step 4 groundwork).

## Commands

Run all backend commands from `backend/`.

| Task | Command |
| --- | --- |
| Start Postgres | `docker compose up -d` |
| Run the API (port 8080) | `go run .` |
| Test | `go test ./...` |
| Format | `gofmt -w .` |
| Vet | `go vet ./...` |
| Frontend (from `frontend/`) | `pnpm install && pnpm dev` |

Database connection settings come from `backend/.env` (`DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_SSLMODE`).

There are no migrations yet. The API creates `card_movements` on startup; `boards`, `columns`, `cards` and `users` must be created by hand using the SQL in `README.md`. When adding a table, add its SQL to the README too.

## Project layout

```
backend/
  main.go                 wiring: env, store, handlers, routes
  internal/models/        structs + validation tags
  internal/store/         Store interface, PostgresStore, one file per resource
  internal/handlers/      one handler struct per resource, plus handlers_test.go
  docker-compose.yml      local Postgres
frontend/                 React + TypeScript (Vite), proxies API routes to :8080
```

New resources follow the board pattern: model → `Store` interface methods → `store/<resource>.go` → `handlers/<resource>.go` → routes in `main.go` → fake-store methods and tests in `handlers_test.go`.
