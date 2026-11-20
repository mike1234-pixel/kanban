package store

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq" // also registers the PostgreSQL driver with database/sql
	"github.com/pressly/goose/v3"
)

// uniqueViolation is the Postgres error code for a unique constraint violation.
const uniqueViolation = "23505"

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == uniqueViolation
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25) // cap active connections to prevent crashing Postgres

	db.SetMaxIdleConns(10) // keep warm connections ready so incoming requests don't wait for TCP handshakes

	db.SetConnMaxLifetime(5 * time.Minute) // retire old connections to clear stale state and handle DB restarts gracefully

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}
	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

// migrations holds the .sql files in migrations/, built into the binary, so
// the API doesn't need the folder next to it when it runs.
//
//go:embed migrations/*.sql
var migrations embed.FS

// migrate runs every migration this database hasn't had yet, in order. goose
// records the ones it has run in its goose_db_version table.
func migrate(db *sql.DB) error {
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, "migrations")
}
