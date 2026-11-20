package store

import (
	"context"
	"database/sql"
	"kanban/internal/models"
	"kanban/internal/store/queries"
	"time"
)

// The SQL for these lives in queries/sessions.sql; sqlc generates the Go that
// runs it and scans the rows. Each method here only converts between the
// generated types and the models.

// #region Create Session
// CreateSession stores a hash of the session token, never the token itself,
// so a leaked sessions table can't be used to log in.
func (s *PostgresStore) CreateSession(ctx context.Context, tokenHash, userID string, expiresAt time.Time) error {
	return s.q.CreateSession(ctx, queries.CreateSessionParams{
		TokenHash: tokenHash,
		UserID:    userID,
		ExpiresAt: expiresAt,
	})
}

//#endregion

// #region Get Session User
// GetSessionUser returns the user a session belongs to, or an empty User if
// the session doesn't exist or has expired.
func (s *PostgresStore) GetSessionUser(ctx context.Context, tokenHash string) (models.User, error) {
	row, err := s.q.GetSessionUser(ctx, tokenHash)
	if err == sql.ErrNoRows {
		return models.User{}, nil
	}
	if err != nil {
		return models.User{}, err
	}
	return models.User{ID: row.ID, Name: row.Name, Email: row.Email, CreatedAt: row.CreatedAt}, nil
}

//#endregion

// #region Delete Session
// DeleteSession removes a session. Deleting one that doesn't exist isn't an
// error: either way, the token no longer works.
func (s *PostgresStore) DeleteSession(ctx context.Context, tokenHash string) error {
	return s.q.DeleteSession(ctx, tokenHash)
}

//#endregion
