package store

import (
	"context"
	"database/sql"
	"kanban/internal/models"
)

// #region Save Board
func (s *PostgresStore) SaveBoard(ctx context.Context, board models.Board) error {
	query := `INSERT INTO boards (id, title) VALUES ($1, $2)`
	_, err := s.db.ExecContext(ctx, query, board.ID, board.Title)
	return err
}

//#endregion

// #region Get Boards
func (s *PostgresStore) GetBoards(ctx context.Context) ([]models.Board, error) {
	query := `SELECT id, title, created_at FROM boards ORDER BY created_at DESC`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []models.Board
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
func (s *PostgresStore) GetBoardByID(ctx context.Context, id string) (models.Board, error) {
	// 1. Fetch Board
	var b models.Board
	query := `SELECT id, title, created_at FROM boards WHERE id = $1`
	err := s.db.QueryRowContext(ctx, query, id).Scan(&b.ID, &b.Title, &b.CreatedAt)
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
	query := `UPDATE boards SET title = $1 WHERE id = $2`
	res, err := s.db.ExecContext(ctx, query, board.Title, board.ID)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	return rows > 0, err
}

func (s *PostgresStore) DeleteBoard(ctx context.Context, id string) (bool, error) {
	query := `DELETE FROM boards WHERE id = $1`
	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}

//#endregion
