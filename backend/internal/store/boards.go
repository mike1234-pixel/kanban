package store

import (
	"context"
	"database/sql"
	"kanban/internal/models"
	"time"
)

// #region Save Board
func (s *PostgresStore) SaveBoard(ctx context.Context, board models.Board) (time.Time, error) {
	query := `INSERT INTO boards (id, title, owner_id) VALUES ($1, $2, $3) RETURNING created_at`
	var createdAt time.Time
	err := s.db.QueryRowContext(ctx, query, board.ID, board.Title, board.OwnerID).Scan(&createdAt)
	if isUniqueViolation(err) {
		return time.Time{}, ErrBoardExists
	}
	return createdAt, err
}

//#endregion

// #region Get Boards
func (s *PostgresStore) GetBoards(ctx context.Context, ownerID string) ([]models.Board, error) {
	query := `SELECT id, title, created_at FROM boards WHERE owner_id = $1 ORDER BY created_at DESC`
	rows, err := s.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	boards := []models.Board{}
	for rows.Next() {
		var b models.Board
		b.Columns = []models.Column{}
		if err := rows.Scan(&b.ID, &b.Title, &b.CreatedAt); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}
	return boards, rows.Err()
}

//#endregion

// #region Get Board By Id
// GetBoardByID returns an empty Board both when the board doesn't exist and
// when it belongs to someone else, so callers can't tell the two apart.
func (s *PostgresStore) GetBoardByID(ctx context.Context, id, ownerID string) (models.Board, error) {
	// 1. Fetch Board
	var b models.Board
	query := `SELECT id, title, created_at FROM boards WHERE id = $1 AND owner_id = $2`
	err := s.db.QueryRowContext(ctx, query, id, ownerID).Scan(&b.ID, &b.Title, &b.CreatedAt)
	if err == sql.ErrNoRows {
		return models.Board{}, nil
	}
	if err != nil {
		return models.Board{}, err
	}

	// 2. Load nested Columns with Cards
	cols, err := s.GetColumnsWithCards(ctx, b.ID)
	if err != nil {
		return models.Board{}, err
	}
	b.Columns = cols

	return b, nil
}

//#endregion

// #region Update Board
func (s *PostgresStore) UpdateBoard(ctx context.Context, board models.Board) (bool, error) {
	query := `UPDATE boards SET title = $1 WHERE id = $2 AND owner_id = $3`
	res, err := s.db.ExecContext(ctx, query, board.Title, board.ID, board.OwnerID)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	return rows > 0, err
}

func (s *PostgresStore) DeleteBoard(ctx context.Context, id, ownerID string) (bool, error) {
	query := `DELETE FROM boards WHERE id = $1 AND owner_id = $2`
	res, err := s.db.ExecContext(ctx, query, id, ownerID)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}

//#endregion
