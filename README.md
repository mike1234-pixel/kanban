# Boardroom

A Kanban board app with a React and TypeScript frontend, a Go HTTP API, and PostgreSQL storage. The frontend supports creating boards and columns, adding and editing cards, and moving cards between columns.

## Requirements

- Go 1.26 or newer
- Node.js and pnpm
- Docker Compose, or a local PostgreSQL server

## Run locally

### 1. Configure the database

Create `backend/.env` with the connection settings used by both Docker Compose and the API:

```dotenv
DB_USER=postgres
DB_PASSWORD=your-local-password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=kanban
DB_SSLMODE=disable
```

Start PostgreSQL from the backend directory:

```sh
cd backend
docker compose up -d
```

You don't create any tables by hand. When the API starts, it runs the migrations in `backend/internal/store/migrations/` with [goose](https://github.com/pressly/goose), so a fresh database gets every table and an existing one gets only the migrations it hasn't had yet. goose records which have run in its `goose_db_version` table.

To change the schema, add a new numbered file to that folder (e.g. `00002_add_cards_column_id_index.sql`). Never edit a migration that has already run.

### 2. Start the API

In a terminal, from the repository root:

```sh
cd backend
go run .
```

The API listens on `http://localhost:8080`.

### 3. Start the frontend

In another terminal, from the repository root:

```sh
cd frontend
pnpm install
pnpm dev
```

Open the local URL printed by Vite (usually `http://localhost:5173`). Vite proxies `/boards`, `/columns`, and `/cards` requests to the API on port 8080.

## Quiz app

`quiz/` is a separate multiple-choice app for the learning quizzes. It reads `QUIZ.md`, `QUIZ_GO.md` and `QUIZ_SQL.md` straight from the repository root, so new questions appear without any extra step. It doesn't need the API or Postgres.

```sh
cd quiz
pnpm install
pnpm dev
```

Open `http://localhost:5174`. Results are saved in the browser's local storage.

## Production build

Build the frontend from `frontend/`:

```sh
pnpm build
```

This runs the TypeScript build and writes the static site to `frontend/dist/`. The Go API is started with `go run .` from `backend/`.

## Run backend tests

From the repository root, run all Go tests in the backend module:

```sh
cd backend
go test ./...
```

The handler tests use an in-memory fake store, so they do not require PostgreSQL or a `.env` file.

## API routes

| Method | Route | Purpose |
| --- | --- | --- |
| `GET`, `POST` | `/boards` | List or create boards |
| `GET`, `PUT`, `DELETE` | `/boards/{id}` | Read, update, or delete a board |
| `GET` | `/columns?board_id={id}` | List columns and cards for a board |
| `POST` | `/columns` | Create a column |
| `GET`, `PUT`, `DELETE` | `/columns/{id}` | Read, update, or delete a column |
| `GET`, `POST` | `/cards` | List or create cards |
| `GET`, `PUT`, `DELETE` | `/cards/{id}` | Read, update, or delete a card |
| `POST` | `/cards/{id}/move` | Move a card to another column |
| `GET`, `POST` | `/users` | List or create users |
| `GET`, `PUT`, `DELETE` | `/users/{id}` | Read, update, or delete a user |
