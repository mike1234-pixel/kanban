package store

import (
	"context"
	"database/sql"
	"fmt"
	"kanban/internal/models"
)

// SaveCard inserts the card only if its column belongs to ownerID. It
// returns false (nothing inserted) if the column is missing or someone else's.
func (s *PostgresStore) SaveCard(ctx context.Context, card models.Card, ownerID string) (bool, error) {
	query := `
		INSERT INTO cards (id, column_id, title, description, "order")
		SELECT $1, $2, $3, $4, $5
		WHERE EXISTS (
			SELECT 1 FROM columns col JOIN boards b ON b.id = col.board_id
			WHERE col.id = $2 AND b.owner_id = $6)`
	res, err := s.db.ExecContext(ctx, query, card.ID, card.ColumnID, card.Title, card.Description, card.Order, ownerID)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}

// #region Get Cards
func (s *PostgresStore) GetCards(ctx context.Context, ownerID string) ([]models.Card, error) {

	query := `
		SELECT ca.id, ca.column_id, ca.title, ca.description, ca."order"
		FROM cards ca
		JOIN columns col ON col.id = ca.column_id
		JOIN boards b ON b.id = col.board_id
		WHERE b.owner_id = $1`

	rows, err := s.db.QueryContext(ctx, query, ownerID)

	if err != nil {
		return nil, err
	}

	defer rows.Close() // release the network connection back to the pool. Defer executes just before the method returns.

	cards := []models.Card{} // for less verbose mapping, consider sqlx - helpers for mapping SQL results into Go structs, sqlc - generates go code from SQL, ORM - generates SQL from go code

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
func (s *PostgresStore) GetCardByID(ctx context.Context, id, ownerID string) (models.Card, error) {

	query := `
		SELECT ca.id, ca.column_id, ca.title, ca.description, ca."order"
		FROM cards ca
		JOIN columns col ON col.id = ca.column_id
		JOIN boards b ON b.id = col.board_id
		WHERE ca.id = $1 AND b.owner_id = $2`

	var c models.Card

	err := s.db.QueryRowContext(ctx, query, id, ownerID).Scan(&c.ID, &c.ColumnID, &c.Title, &c.Description, &c.Order)

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
// UpdateCard only updates the card if both its current column and its new
// column belong to ownerID; otherwise it returns false.
func (s *PostgresStore) UpdateCard(ctx context.Context, id string, updated models.Card, ownerID string) (bool, error) {
	query := `
		UPDATE cards ca
		SET column_id = $1, title = $2, description = $3, "order" = $4
		WHERE ca.id = $5
		  AND EXISTS (
			SELECT 1 FROM columns col JOIN boards b ON b.id = col.board_id
			WHERE col.id = ca.column_id AND b.owner_id = $6)
		  AND EXISTS (
			SELECT 1 FROM columns col JOIN boards b ON b.id = col.board_id
			WHERE col.id = $1 AND b.owner_id = $6)`
	result, err := s.db.ExecContext(ctx, query, updated.ColumnID, updated.Title, updated.Description, updated.Order, id, ownerID)

	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

// MoveCard changes a card's column and records the move atomically. It returns
// false if the card is missing or someone else's, and ErrColumnNotFound if the
// destination column is.
func (s *PostgresStore) MoveCard(ctx context.Context, id, columnID, ownerID string) (bool, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	var fromColumnID string
	err = tx.QueryRowContext(ctx, `
		SELECT ca.column_id FROM cards ca
		WHERE ca.id = $1 AND EXISTS (
			SELECT 1 FROM columns col JOIN boards b ON b.id = col.board_id
			WHERE col.id = ca.column_id AND b.owner_id = $2)
		FOR UPDATE`, id, ownerID).Scan(&fromColumnID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	var canMoveThere bool
	err = tx.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM columns col JOIN boards b ON b.id = col.board_id
			WHERE col.id = $1 AND b.owner_id = $2)`, columnID, ownerID).Scan(&canMoveThere)
	if err != nil {
		return false, err
	}
	if !canMoveThere {
		return false, ErrColumnNotFound
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
func (s *PostgresStore) DeleteCard(ctx context.Context, id, ownerID string) (bool, error) {
	query := `
		DELETE FROM cards ca
		WHERE ca.id = $1 AND EXISTS (
			SELECT 1 FROM columns col JOIN boards b ON b.id = col.board_id
			WHERE col.id = ca.column_id AND b.owner_id = $2)`

	result, err := s.db.ExecContext(ctx, query, id, ownerID)

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
