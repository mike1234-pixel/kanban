package store

import (
	"context"
	"time"
)

// #region Create Session
// CreateSession stores a hash of the session token, never the token itself,
// so a leaked sessions table can't be used to log in.
func (s *PostgresStore) CreateSession(ctx context.Context, tokenHash, userID string, expiresAt time.Time) error {
	query := `INSERT INTO sessions (token_hash, user_id, expires_at) VALUES ($1, $2, $3)`
	_, err := s.db.ExecContext(ctx, query, tokenHash, userID, expiresAt)
	return err
}

//#endregion
