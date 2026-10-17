package store

import (
	"context"
	"database/sql"
	"kanban/internal/models"
)

// #region Save Column
func (s *PostgresStore) SaveColumn(ctx context.Context, col models.Column) error {

	query := `INSERT INTO columns (id, board_id, title, position) VALUES ($1, $2, $3, $4)`

	_, err := s.db.ExecContext(ctx, query, col.ID, col.BoardID, col.Title, col.Position)

	return err
}

//#endregion

// #region Get Columns With Cards
func (s *PostgresStore) GetColumnsWithCards(ctx context.Context, boardID string) ([]models.Column, error) {
	// GetColumnsWithCards performs a two-step load to construct the full board hierarchy

	// 1. Fetch columns filtered by board_id
	colQuery := `SELECT id, board_id, title, position FROM columns WHERE board_id = $1 ORDER BY position ASC`
	colRows, err := s.db.QueryContext(ctx, colQuery, boardID)
	if err != nil {
		return nil, err
	}
	defer colRows.Close()

	var columns []models.Column
	colMap := make(map[string]int)

	for colRows.Next() {
		var c models.Column
		c.Cards = []models.Card{}
		if err := colRows.Scan(&c.ID, &c.BoardID, &c.Title, &c.Position); err != nil {
			return nil, err
		}
		colMap[c.ID] = len(columns)
		columns = append(columns, c)
	}
	if err := colRows.Err(); err != nil {
		return nil, err
	}

	if len(columns) == 0 {
		return []models.Column{}, nil
	}

	// 2. Fetch cards
	cardQuery := `SELECT id, column_id, title, description, "order" FROM cards ORDER BY "order" ASC`
	cardRows, err := s.db.QueryContext(ctx, cardQuery)
	if err != nil {
		return nil, err
	}
	defer cardRows.Close()

	for cardRows.Next() {
		var card models.Card
		if err := cardRows.Scan(&card.ID, &card.ColumnID, &card.Title, &card.Description, &card.Order); err != nil {
			return nil, err
		}
		if idx, exists := colMap[card.ColumnID]; exists {
			columns[idx].Cards = append(columns[idx].Cards, card)
		}
	}

	return columns, cardRows.Err()
}

//#endregion

// #region Get Column By Id
func (s *PostgresStore) GetColumnByID(ctx context.Context, id string) (models.Column, error) {
	query := `SELECT id, board_id, title, position, created_at FROM columns WHERE id = $1`
	var c models.Column
	err := s.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.BoardID, &c.Title, &c.Position, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return models.Column{}, nil
	}
	if err != nil {
		return models.Column{}, err
	}
	return c, nil
}

//#endregion

// #region Update Column
func (s *PostgresStore) UpdateColumn(ctx context.Context, col models.Column) (bool, error) {
	query := `
		UPDATE columns 
		SET title = $1, position = $2 
		WHERE id = $3
	`
	res, err := s.db.ExecContext(ctx, query, col.Title, col.Position, col.ID)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	return rows > 0, err
}

//#endregion

// #region Delete Column
func (s *PostgresStore) DeleteColumn(ctx context.Context, id string) (bool, error) {
	query := `DELETE FROM columns WHERE id = $1`
	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}

//#endregion
