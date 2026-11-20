# Kanban Backend (Go)

## Purpose

This is a learning project. The goal is to understand Go and what a backend is actually doing, not to ship features quickly. The Kanban app is the vehicle; each milestone introduces a new backend concept.

Optimise for my understanding of concepts, not for hand-writing code. Pace matters too: the goal is to cover the whole roadmap in about 6 months (by April 2027), roughly two weeks per step.

## How to help

Workflow for each roadmap step (theory first, then code):

1. **Write up** the step's concepts in `LEARNING.md` and add matching questions to `QUIZ.md`, before any code. Use worked examples from this codebase.
2. **Revise:** I read the entries and ask questions. Explain one concept at a time, briefly, with its trade-offs.
3. **Quiz** me on the step's questions to reinforce the concepts, until I feel comfortable.
4. **Implement** the concepts in code yourself, in small steps, so each diff is easy to review. I don't hand-write the code; I learn by reading the diffs. After each diff, check my understanding with one question about what changed.
5. **Update** the entries if the code ends up differing from what they say. Go syntax and SQL the diffs introduce go in `LEARNING_GO.md` / `QUIZ_GO.md` and `LEARNING_SQL.md` / `QUIZ_SQL.md`.

Also:

- Keep explanations bitesized. Not walls of text.

- Core concepts first, not minutiae. Leave out edge cases, caveats and extra clarifications unless they're essential; they muddy the key idea.

- When I check my understanding or answer a quiz question, start with a clear "Correct" or "Wrong". If wrong, give only the correction.

- Explain one concept at a time and check that I have understood the concept before moving on.

- Explain briefly why the code you write is idiomatic Go.

- Name design patterns explicitly. When code or an entry uses one (repository, dependency injection, middleware/decorator, worker pool…), say its name and add a `🧩 **Pattern: name.**` line (the 🧩 marks it as a design pattern) to the entry that explains it. No separate patterns file.

- Don't silently fix bugs.

- Connect new work back to the layers it touches: TCP → HTTP → router → middleware → handler → store → Postgres → response.

- Reinforce basic concepts. I wish to learn by repetition.

- Stay within the current milestone. If something belongs to a later milestone, say so rather than building it now.

## Technical constraints

- Standard library first: `net/http`, `database/sql`, `encoding/json`, `context`, `log/slog`.

- No ORM and no web framework for now. The verbosity of `database/sql` is intentional; it shows what the database interaction really looks like. Step 6 converted only `store/sessions.go` to sqlc, for comparison; every other store file stays hand-written `database/sql`.

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

Work through the roadmap in order, one step at a time.

Ticks in `ROADMAP.md` show how well I know each step. I add them by hand; don't change them.

- ✅ done
- ✅✅ secure knowledge (e.g. passed its quiz 100%)
- ✅✅✅ strong knowledge

@ROADMAP.md

## Learning log

`LEARNING.md` is a quick reference for the roadmap concepts I've learned. Keep it up to date:

- At the start of each step, add an entry for each of its concepts (see the workflow in "How to help").
- Group entries under a heading per roadmap step (e.g. `## 2. HTTP Fundamentals`), using the concept names from the roadmap.
- Go syntax goes in `LEARNING_GO.md`, and SQL in `LEARNING_SQL.md` (see below).
- It's a quick reference, not a tutorial: each entry is at most 300 words (aim much shorter), mostly a one-line definition plus a few key bullet points.
- If the code changes in a way that makes an entry wrong, update the entry.
- Keep the Contents list at the top in step with the headings: add a link for every new `##` or `###`. The same goes for `LEARNING_GO.md` and `LEARNING_SQL.md`.

## Go and SQL references

`LEARNING_GO.md` (Go syntax) and `LEARNING_SQL.md` (SQL and Postgres) are quick references, not on the roadmap. They cover what this codebase uses, with examples taken from it.

- When a diff uses Go syntax or SQL that isn't in its reference file yet (e.g. goroutines, `switch`, `LEFT JOIN`, `GROUP BY`), point it out and add an entry.
- Same format as `LEARNING.md`: a one-line definition, a short example, a few bullets.

## Quiz

`QUIZ.md` holds recall questions on the concepts in `LEARNING.md`; `QUIZ_GO.md` and `QUIZ_SQL.md` do the same for `LEARNING_GO.md` and `LEARNING_SQL.md`.

- Questions test recall of concepts, not predicting what code does by reading it.
- Keep each question short, with a short answer in a collapsible `<details>` block.
- Every question also has four options for the quiz app, as a task list between the question and the `<details>`: `- [x]` marks the one correct option, `- [ ]` the three wrong ones. Make the wrong ones believable (common mistakes, including ones I've made), keep all four a similar length, and vary where the correct one sits.
- The `<details>` answer must add insight, not restate the correct option: explain why it's right, with a short concrete example (real values from this codebase), in 2–4 lines.
- Put the key idea in **bold** in answers when one phrase matters most: the quiz app shows bold with a yellow highlighter.
- When an entry is added to any `LEARNING*.md` file, add 2-3 matching questions to its quiz file under the same heading.
- When quizzing me: one question at a time, keep a running score, and briefly explain any miss before moving on.

## Current milestone

<!-- Update this as you progress -->

Working on: step 7 (collaboration). Next: write its `LEARNING.md` entries and `QUIZ.md` questions. Steps 1–6 are complete.

- Step 7 changes the schema (`workspaces`, `workspace_members`, moving `boards.owner_id` to `workspace_id`), so every change is a new numbered migration, never an edit to `00001` or `00002`.

- Done (step 6): concepts in `LEARNING.md`, quiz questions 99–118 (plus 110a, 110b, 114a).
  - Done: goose migrations. `00001_initial_schema.sql` (all tables, `IF NOT EXISTS`) runs on startup from `NewPostgresStore` via `//go:embed`; the README's hand-run SQL is gone.
  - Done: `00002_add_cards_column_id_index.sql`. `EXPLAIN ANALYZE` with 100,000 cards: `Seq Scan` 3.6 ms → `Index Scan` 0.025 ms.
  - Done: sqlc. `sqlc.yaml` reads the schema from the migrations; `store/queries/sessions.sql` → generated `sessions.sql.go`; `sessions.go` wraps it. sqlc is pinned as a `tool` in `go.mod`. `TestSessionLifecycle` covers it.
  - Done: #9 fixed. `GetColumnsWithCards` fetches only the board's cards (join to `columns`, `WHERE col.board_id = $1`); `TestGetColumnsWithCardsOneBoard`.

- Done: user CRUD endpoints (`/users`) with unique emails (409 on duplicates).
- Done: password hashing with bcrypt. `POST /users` requires a password (8–72 characters) and stores only its hash in `users.password_hash`; tested live.
- Done: sessions + cookies. `POST /login` sets an `HttpOnly`, `Secure`, `SameSite=Lax` `session` cookie; `sessions` stores only the token's SHA-256 (`hashToken` in `handlers/auth.go`). Nothing reads the cookie yet.
- Done: auth middleware. `RequireAuth` (`handlers/middleware.go`) protects `GET /me`; handlers read the user with `userFromContext`. Board routes aren't protected yet (the frontend has no login).
- Done: logout. `POST /logout` deletes the session and clears the cookie; always 204.
- Done: ownership checks. `boards.owner_id`; every board, column and card query filters by the logged-in user (columns and cards via joins to `boards`); someone else's resource returns 404. All board, column and card routes are behind `RequireAuth`.
- Step 3 is complete. Test login: `ada@example.com` / `correct-horse` (owns the existing boards).
- The frontend has a minimal login page (`AuthGate` checks `GET /me`; `LoginPage`; Log out in the sidebar). Use Chrome or Firefox locally: the `Secure` cookie isn't sent over `http://localhost` in Safari.
- Done (step 4): `MoveCard` runs in one transaction (`BeginTx` / `defer tx.Rollback()` / `Commit`), updating the card and writing to `card_movements`. `card_movements` has foreign keys with `ON DELETE CASCADE`; moving a card into its current column writes nothing.
- Done (step 5): concepts revised and quizzed (revisit later).
  - Pessimistic: `MoveCard` locks the card with a join-free `SELECT … FOR UPDATE`, then checks ownership (#8 fixed). `TestMoveCardConcurrent` (`store/cards_test.go`) moves one card from two goroutines against the real Postgres.
  - Optimistic: `PUT /cards/{id}` sends `version`; `ErrVersionConflict` → `409 CARD_VERSION_CONFLICT`; `MoveCard` bumps the version too; the frontend sends it. `cards.version` is in the README (with the `ALTER`). Tested live.
  - Go concurrency is only in tests so far: `TestUpdateCardConcurrent` (`handlers/concurrency_test.go`) uses goroutines, a `WaitGroup`, a `sync.Mutex` in an in-memory store, and a buffered results channel. No production code uses goroutines, mutexes or channels yet (the job queue is step 9).
  - Store tests that need Postgres skip themselves when `backend/.env` is missing.

## Commands

Run all backend commands from `backend/`.

| Task                        | Command                    |
| --------------------------- | -------------------------- |
| Start Postgres              | `docker compose up -d`     |
| Run the API (port 8080)     | `go run .`                 |
| Test                        | `go test ./...`            |
| Format                      | `gofmt -w .`               |
| Vet                         | `go vet ./...`             |
| Regenerate sqlc code        | `go tool sqlc generate`    |
| Frontend (from `frontend/`) | `pnpm install && pnpm dev` |
| Quiz app (from `quiz/`)     | `pnpm install && pnpm dev` (port 5174) |

Database connection settings come from `backend/.env` (`DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_SSLMODE`).

The schema lives in goose migrations (`backend/internal/store/migrations/`), which `NewPostgresStore` runs on startup. Change the schema by adding a new numbered `.sql` file; never edit one that has already run.

## Project layout

```
backend/
  main.go                 wiring: env, store, handlers, routes
  internal/models/        structs + validation tags
  internal/store/         Store interface, PostgresStore, one file per resource
  internal/handlers/      one handler struct per resource, plus handlers_test.go
  docker-compose.yml      local Postgres
frontend/                 React + TypeScript (Vite), proxies API routes to :8080
quiz/                     multiple-choice quiz app that reads the QUIZ*.md files; keep their format
                          (`##` section, `###` topic, `**N.**` question, `- [x]`/`- [ ]` options, <details> answer)
```

New resources follow the board pattern: model → `Store` interface methods → `store/<resource>.go` → `handlers/<resource>.go` → routes in `main.go` → fake-store methods and tests in `handlers_test.go`.
