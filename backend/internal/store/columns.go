package store

import (
	"context"
	"database/sql"
	"kanban/internal/models"
)

// #region Save Column
// SaveColumn inserts the column only if its board belongs to ownerID. It
// returns false (nothing inserted) if the board is missing or someone else's.
func (s *PostgresStore) SaveColumn(ctx context.Context, col models.Column, ownerID string) (bool, error) {
	query := `
		INSERT INTO columns (id, board_id, title, position)
		SELECT $1, $2, $3, $4
		WHERE EXISTS (SELECT 1 FROM boards WHERE id = $2 AND owner_id = $5)`
	res, err := s.db.ExecContext(ctx, query, col.ID, col.BoardID, col.Title, col.Position, ownerID)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}

//#endregion

// #region Get Columns With Cards
func (s *PostgresStore) GetColumnsWithCards(ctx context.Context, boardID, ownerID string) ([]models.Column, error) {
	// GetColumnsWithCards performs a two-step load to construct the full board hierarchy

	// 1. Fetch the board's columns, only if the board belongs to ownerID
	colQuery := `
		SELECT c.id, c.board_id, c.title, c.position
		FROM columns c
		JOIN boards b ON b.id = c.board_id
		WHERE c.board_id = $1 AND b.owner_id = $2
		ORDER BY c.position ASC`
	colRows, err := s.db.QueryContext(ctx, colQuery, boardID, ownerID)
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
func (s *PostgresStore) GetColumnByID(ctx context.Context, id, ownerID string) (models.Column, error) {
	query := `
		SELECT c.id, c.board_id, c.title, c.position, c.created_at
		FROM columns c
		JOIN boards b ON b.id = c.board_id
		WHERE c.id = $1 AND b.owner_id = $2`
	var c models.Column
	err := s.db.QueryRowContext(ctx, query, id, ownerID).Scan(&c.ID, &c.BoardID, &c.Title, &c.Position, &c.CreatedAt)
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
func (s *PostgresStore) UpdateColumn(ctx context.Context, col models.Column, ownerID string) (bool, error) {
	query := `
		UPDATE columns c
		SET title = $1, position = $2
		FROM boards b
		WHERE c.id = $3 AND b.id = c.board_id AND b.owner_id = $4`
	res, err := s.db.ExecContext(ctx, query, col.Title, col.Position, col.ID, ownerID)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	return rows > 0, err
}

//#endregion

// #region Delete Column
func (s *PostgresStore) DeleteColumn(ctx context.Context, id, ownerID string) (bool, error) {
	query := `
		DELETE FROM columns c
		USING boards b
		WHERE c.id = $1 AND b.id = c.board_id AND b.owner_id = $2`
	res, err := s.db.ExecContext(ctx, query, id, ownerID)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}

//#endregion
