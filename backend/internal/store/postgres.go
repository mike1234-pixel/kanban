package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq" // also registers the PostgreSQL driver with database/sql
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
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS card_movements (
			id BIGSERIAL PRIMARY KEY,
			card_id UUID NOT NULL,
			from_column_id UUID NOT NULL,
			to_column_id UUID NOT NULL,
			moved_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create card movement history table: %w", err)
	}

	return &PostgresStore{db: db}, nil
}
