# Kanban Backend (Go)

## Purpose

This is a learning project. The goal is to understand Go and what a backend is actually doing, not to ship features quickly. The Kanban app is the vehicle; each milestone introduces a new backend concept.

Optimise for my understanding of concepts, not for hand-writing code. Pace matters too: the goal is to cover the whole roadmap in about 6 months (by April 2027), roughly two weeks per step.

## How to help

Workflow for each concept:

1. **Explain** the concept and its trade-offs briefly. I may want to probe further before moving on.
2. **Implement** it yourself in small steps, so each diff is easy to review. I don't hand-write the code; I learn by reading the diffs.
3. **Check** my understanding with one question about what changed.
4. **Record** it in `LEARNING.md` and `QUIZ.md` (backend concepts) or `LEARNING_GO.md` and `QUIZ_GO.md` (Go syntax).

Also:

- Keep explanations bitesized. Not walls of text.

- Explain one concept at a time and check that I have understood the concept before moving on.

- Explain briefly why the code you write is idiomatic Go.

- Don't silently fix bugs.

- Connect new work back to the layers it touches: TCP → HTTP → router → middleware → handler → store → Postgres → response.

- Reinforce basic concepts. I wish to learn by repetition.

- Stay within the current milestone. If something belongs to a later milestone, say so rather than building it now.

## Technical constraints

- Standard library first: `net/http`, `database/sql`, `encoding/json`, `context`, `log/slog`.

- No ORM and no web framework for now. The verbosity of `database/sql` is intentional; it shows what the database interaction really looks like. Tools like sqlc or pgx come later (roadmap step 6), once I understand what they abstract.

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

- Write responses with the helpers in `handlers/respond.go`, never `http.Error`: `writeJSON` for success, `writeError` for client errors, `writeInternalError` for 500s.
- Error codes are constants in `respond.go`. Resource-specific failures get one code per resource (`BOARD_NOT_FOUND`, `USER_EMAIL_TAKEN`); failures that mean the same thing on every endpoint share a code (`INVALID_JSON`, `VALIDATION_FAILED`, `INTERNAL_ERROR`).
- 500s log the real error with `slog` (method, path, error) and send the client a generic message. Never put `err.Error()` from the store in a response.
- Error messages come from the handler, not from an error's text, so changing an internal error can't change the API.
- Database constraints (NOT NULL, UNIQUE, FOREIGN KEY) enforce what must never happen; application code checks what should happen. Map constraint violations to sensible API errors (e.g. UNIQUE → 409).

### Known issues

Bugs and gaps are tracked in GitHub Issues: https://github.com/mike1234-pixel/kanban/issues

- When you find a problem you aren't fixing now, list it and offer to open an issue (via `gh`) with what's wrong, where it is, and the fix. Create it only after I confirm.
- Label issues `bug` or `enhancement`.
- When a commit fixes an issue, include `fixes #N` in the commit message.

## Roadmap

Work through the roadmap in order, one step at a time:

@ROADMAP.md

## Learning log

`LEARNING.md` is a quick reference for the roadmap concepts I've learned. Keep it up to date:

- When I've learned a concept (you've checked my understanding), add an entry for it.
- Group entries under a heading per roadmap step (e.g. `## 2. HTTP Fundamentals`), using the concept names from the roadmap.
- Go syntax goes in `LEARNING_GO.md` instead (see below).
- It's a quick reference, not a tutorial: each entry is at most 300 words (aim much shorter), mostly a one-line definition plus a few key bullet points.
- List concepts from the current step that I haven't learned yet under "Not yet covered", and remove them from that list as they get entries.
- If the code changes in a way that makes an entry wrong, update the entry.

## Go reference

`LEARNING_GO.md` is a quick reference for Go syntax, not on the roadmap. It covers the Go used in this codebase, with examples taken from it.

- When a diff uses Go syntax that isn't in `LEARNING_GO.md` yet (e.g. goroutines, channels, `switch`), point it out and add an entry.
- Same format as `LEARNING.md`: a one-line definition, a short example, a few bullets.

## Quiz

`QUIZ.md` holds recall questions on the concepts in `LEARNING.md`; `QUIZ_GO.md` does the same for `LEARNING_GO.md`.

- Questions test recall of concepts, not predicting what code does by reading it.
- Keep each question short, with a short answer in a collapsible `<details>` block.
- When an entry is added to `LEARNING.md` or `LEARNING_GO.md`, add 2-3 matching questions to its quiz file under the same heading.
- When quizzing me: one question at a time, keep a running score, and briefly explain any miss before moving on.

## Current milestone

<!-- Update this as you progress -->

Working on: step 3, authentication / authorization.

- Done: user CRUD endpoints (`/users`) with unique emails (409 on duplicates).
- In progress: password hashing with bcrypt. `users.password_hash` exists (`TEXT NOT NULL`); `POST /users` doesn't fill it yet, so it returns 500 until that's done.
- Next: sessions + cookies → auth middleware → ownership checks.
- Note: `MoveCard` already runs in a transaction and writes to `card_movements` (step 4 groundwork).

## Commands

Run all backend commands from `backend/`.

| Task                        | Command                    |
| --------------------------- | -------------------------- |
| Start Postgres              | `docker compose up -d`     |
| Run the API (port 8080)     | `go run .`                 |
| Test                        | `go test ./...`            |
| Format                      | `gofmt -w .`               |
| Vet                         | `go vet ./...`             |
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
