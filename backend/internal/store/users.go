package store

import (
	"context"
	"database/sql"
	"kanban/internal/models"
	"time"
)

// #region Save User
func (s *PostgresStore) SaveUser(ctx context.Context, user models.User) (time.Time, error) {
	query := `INSERT INTO users (id, name, email) VALUES ($1, $2, $3) RETURNING created_at`
	var createdAt time.Time
	err := s.db.QueryRowContext(ctx, query, user.ID, user.Name, user.Email).Scan(&createdAt)
	if isUniqueViolation(err) {
		return time.Time{}, ErrEmailTaken
	}
	return createdAt, err
}

//#endregion

// #region Get Users
func (s *PostgresStore) GetUsers(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, name, email, created_at FROM users ORDER BY created_at DESC`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

//#endregion

// #region Get User By Id
func (s *PostgresStore) GetUserByID(ctx context.Context, id string) (models.User, error) {
	var u models.User
	query := `SELECT id, name, email, created_at FROM users WHERE id = $1`
	err := s.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return models.User{}, nil
	}
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

//#endregion

// #region Update User
func (s *PostgresStore) UpdateUser(ctx context.Context, user models.User) (bool, error) {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
	res, err := s.db.ExecContext(ctx, query, user.Name, user.Email, user.ID)
	if isUniqueViolation(err) {
		return false, ErrEmailTaken
	}
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	return rows > 0, err
}

//#endregion

// #region Delete User
func (s *PostgresStore) DeleteUser(ctx context.Context, id string) (bool, error) {
	query := `DELETE FROM users WHERE id = $1`
	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}

//#endregion
