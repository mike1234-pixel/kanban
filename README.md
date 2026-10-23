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

The API expects the `boards`, `columns`, `cards`, and `users` tables to exist. The current repository does not include database migrations. For a fresh local database, create the tables with:

```sql
CREATE TABLE boards (
  id UUID PRIMARY KEY,
  title VARCHAR(100) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE columns (
  id UUID PRIMARY KEY,
  board_id UUID NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
  title VARCHAR(100) NOT NULL,
  position INTEGER NOT NULL CHECK (position >= 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE cards (
  id UUID PRIMARY KEY,
  column_id UUID NOT NULL REFERENCES columns(id) ON DELETE CASCADE,
  title VARCHAR(150) NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  "order" INTEGER NOT NULL CHECK ("order" >= 0)
);

CREATE TABLE users (
  id UUID PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(254) NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sessions (
  token_hash TEXT PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

The API creates the `card_movements` table automatically when it connects.

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
