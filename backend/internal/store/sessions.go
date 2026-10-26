package store

import (
	"context"
	"database/sql"
	"kanban/internal/models"
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

// #region Get Session User
// GetSessionUser returns the user a session belongs to, or an empty User if
// the session doesn't exist or has expired.
func (s *PostgresStore) GetSessionUser(ctx context.Context, tokenHash string) (models.User, error) {
	var u models.User
	query := `
		SELECT u.id, u.name, u.email, u.created_at
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.token_hash = $1 AND s.expires_at > NOW()`
	err := s.db.QueryRowContext(ctx, query, tokenHash).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return models.User{}, nil
	}
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

//#endregion

// #region Delete Session
// DeleteSession removes a session. Deleting one that doesn't exist isn't an
// error: either way, the token no longer works.
func (s *PostgresStore) DeleteSession(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM sessions WHERE token_hash = $1`
	_, err := s.db.ExecContext(ctx, query, tokenHash)
	return err
}

//#endregion
