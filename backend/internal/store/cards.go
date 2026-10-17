package store

import (
	"context"
	"database/sql"
	"fmt"
	"kanban/internal/models"
)

func (s *PostgresStore) SaveCard(ctx context.Context, card models.Card) error {

	query := `
		INSERT INTO cards (id, column_id, title, description, "order")
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := s.db.ExecContext(ctx, query, card.ID, card.ColumnID, card.Title, card.Description, card.Order)

	return err // Go functions that perform side effects (like database writes, file updates, or HTTP requests) follow a standard convention: if the function returns (error) and that value is nil, the operation succeeded completely.
}

// #region Get Cards
func (s *PostgresStore) GetCards(ctx context.Context) ([]models.Card, error) {

	query := `SELECT id, column_id, title, description, "order" FROM cards`

	rows, err := s.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close() // release the network connection back to the pool. Defer executes just before the method returns.

	var cards []models.Card // for less verbose mapping, consider sqlx - helpers for mapping SQL results into Go structs, sqlc - generates go code from SQL, ORM - generates SQL from go code

	for rows.Next() {
		var c models.Card

		if err := rows.Scan(&c.ID, &c.ColumnID, &c.Title, &c.Description, &c.Order); err != nil {
			return nil, err
		}

		cards = append(cards, c)
	}

	if err := rows.Err(); err != nil { // After iteration has stopped, check for streaming/network errors that occurred during iteration
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return cards, nil // when to return more than just nil - when the database generates data during the insert that your application needs to know (e.g., auto-incrementing IDs or default timestamps)
}

//#endregion

// #region Get Card By Id
func (s *PostgresStore) GetCardByID(ctx context.Context, id string) (models.Card, error) {

	query := `SELECT id, column_id, title, description, "order" FROM cards WHERE id = $1`

	var c models.Card

	err := s.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.ColumnID, &c.Title, &c.Description, &c.Order)

	if err == sql.ErrNoRows { // query ran successfully, but there wasn't a card matching what I asked for.
		return models.Card{}, nil // return empty card to signal not found
	}

	if err != nil {
		return models.Card{}, err
	}

	return c, nil
}

//#endregion

// #region Update Card
func (s *PostgresStore) UpdateCard(ctx context.Context, id string, updated models.Card) (bool, error) {
	query := `
		UPDATE cards 
		SET column_id = $1, title = $2, description = $3, "order" = $4
		WHERE id = $5
	`
	result, err := s.db.ExecContext(ctx, query, updated.ColumnID, updated.Title, updated.Description, updated.Order, id)

	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

// MoveCard changes a card's column and records the move atomically.
func (s *PostgresStore) MoveCard(ctx context.Context, id string, columnID string) (bool, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	var fromColumnID string
	err = tx.QueryRowContext(ctx, `SELECT column_id FROM cards WHERE id = $1 FOR UPDATE`, id).Scan(&fromColumnID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE cards SET column_id = $1 WHERE id = $2`, columnID, id); err != nil {
		return false, err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO card_movements (card_id, from_column_id, to_column_id) VALUES ($1, $2, $3)`, id, fromColumnID, columnID); err != nil {
		return false, err
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

//#endregion

// #region Delete Card
func (s *PostgresStore) DeleteCard(ctx context.Context, id string) (bool, error) {
	query := `DELETE FROM cards WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

//#endregion
