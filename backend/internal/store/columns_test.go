package store

import (
	"context"
	"testing"
)

// TestGetColumnsWithCardsOneBoard checks a board only gets its own cards,
// with another user's board and card in the database too (#9).
func TestGetColumnsWithCardsOneBoard(t *testing.T) {
	s := testStore(t)
	ada := seedBoard(t, s)
	seedBoard(t, s) // someone else's board, with its own card

	var boardID string
	if err := s.db.QueryRow(`SELECT board_id FROM columns WHERE id = $1`, ada.todo).Scan(&boardID); err != nil {
		t.Fatal(err)
	}

	columns, err := s.GetColumnsWithCards(context.Background(), boardID, ada.ownerID)
	if err != nil {
		t.Fatal(err)
	}
	cards := 0
	for _, col := range columns {
		for _, card := range col.Cards {
			if card.ID != ada.cardID {
				t.Fatalf("got card %s, which isn't on this board", card.ID)
			}
			cards++
		}
	}
	if len(columns) != 3 || cards != 1 {
		t.Fatalf("got %d columns and %d cards; want 3 and 1", len(columns), cards)
	}
}
