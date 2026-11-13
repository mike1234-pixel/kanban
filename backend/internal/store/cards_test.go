package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"

	"kanban/internal/models"

	"github.com/joho/godotenv"
)

// testStore connects to the Postgres from backend/.env. It skips the test if
// there's no database, so go test ./... still works without Docker.
func testStore(t *testing.T) *PostgresStore {
	t.Helper()
	godotenv.Load("../../.env")
	if os.Getenv("DB_HOST") == "" {
		t.Skip("no database configured (backend/.env)")
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_SSLMODE"))
	s, err := NewPostgresStore(connStr)
	if err != nil {
		t.Skipf("database not reachable: %v", err)
	}
	t.Cleanup(func() { s.db.Close() })
	return s
}

// board is the test data seedBoard creates: one card in "todo".
type board struct {
	ownerID, todo, doing, done, cardID string
}

// seedBoard creates a user who owns a board with three columns and one card.
// Deleting the user at the end removes everything else through ON DELETE CASCADE.
func seedBoard(t *testing.T, s *PostgresStore) board {
	t.Helper()
	ctx := context.Background()
	var b board

	insert := func(dest *string, query string, args ...any) {
		t.Helper()
		if err := s.db.QueryRowContext(ctx, query, args...).Scan(dest); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}

	insert(&b.ownerID, `INSERT INTO users (id, name, email, password_hash)
		VALUES (gen_random_uuid(), 'Race Test', gen_random_uuid() || '@example.com', 'x') RETURNING id`)
	t.Cleanup(func() { s.db.Exec(`DELETE FROM users WHERE id = $1`, b.ownerID) })

	var boardID string
	insert(&boardID, `INSERT INTO boards (id, title, owner_id) VALUES (gen_random_uuid(), 'Race board', $1) RETURNING id`, b.ownerID)
	insert(&b.todo, `INSERT INTO columns (id, board_id, title, position) VALUES (gen_random_uuid(), $1, 'todo', 0) RETURNING id`, boardID)
	insert(&b.doing, `INSERT INTO columns (id, board_id, title, position) VALUES (gen_random_uuid(), $1, 'doing', 1) RETURNING id`, boardID)
	insert(&b.done, `INSERT INTO columns (id, board_id, title, position) VALUES (gen_random_uuid(), $1, 'done', 2) RETURNING id`, boardID)
	insert(&b.cardID, `INSERT INTO cards (id, column_id, title, "order") VALUES (gen_random_uuid(), $1, 'c1', 0) RETURNING id`, b.todo)

	return b
}

// TestMoveCardConcurrent moves the same card from two goroutines at once:
// one to "doing", one to "done". Whichever runs second must see the first
// move, so the history chains up: todo → X, then X → Y. Without FOR UPDATE,
// both can read "todo" and the second row records the wrong from-column.
// It repeats because a race only shows when the two moves overlap.
func TestMoveCardConcurrent(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()

	for round := range 20 {
		b := seedBoard(t, s)

		var wg sync.WaitGroup
		var errDoing, errDone error
		var foundDoing, foundDone bool
		wg.Add(2)
		go func() {
			defer wg.Done()
			foundDoing, errDoing = s.MoveCard(ctx, b.cardID, b.doing, b.ownerID)
		}()
		go func() {
			defer wg.Done()
			foundDone, errDone = s.MoveCard(ctx, b.cardID, b.done, b.ownerID)
		}()
		wg.Wait() // block until both goroutines have called Done
		if errDoing != nil || errDone != nil {
			t.Fatalf("round %d: move failed: %v, %v", round, errDoing, errDone)
		}
		if !foundDoing || !foundDone {
			t.Fatalf("round %d: a move reported the card as not found (found: doing=%v, done=%v)", round, foundDoing, foundDone)
		}

		rows, err := s.db.QueryContext(ctx,
			`SELECT from_column_id, to_column_id FROM card_movements WHERE card_id = $1 ORDER BY id`, b.cardID)
		if err != nil {
			t.Fatal(err)
		}
		var moves [][2]string
		for rows.Next() {
			var from, to string
			if err := rows.Scan(&from, &to); err != nil {
				t.Fatal(err)
			}
			moves = append(moves, [2]string{from, to})
		}
		rows.Close()

		if len(moves) != 2 {
			t.Fatalf("round %d: want 2 movement rows, got %d", round, len(moves))
		}
		if moves[0][0] != b.todo {
			t.Errorf("round %d: first move should start in todo", round)
		}
		if moves[1][0] != moves[0][1] {
			t.Errorf("round %d: lost update: second move starts in %s, but the first move ended in %s",
				round, columnName(b, moves[1][0]), columnName(b, moves[0][1]))
		}
	}
}

// columnName turns a column ID back into "todo", "doing" or "done" for messages.
func columnName(b board, id string) string {
	switch id {
	case b.todo:
		return "todo"
	case b.doing:
		return "doing"
	case b.done:
		return "done"
	}
	return id
}

// TestMoveCardOtherOwner checks someone else's card is reported as not found.
func TestMoveCardOtherOwner(t *testing.T) {
	s := testStore(t)
	ada := seedBoard(t, s)
	bo := seedBoard(t, s)

	found, err := s.MoveCard(context.Background(), ada.cardID, ada.doing, bo.ownerID)
	if err != nil {
		t.Fatal(err)
	}
	if found {
		t.Fatal("Bo moved Ada's card; want not found")
	}
}

// TestUpdateCardVersion saves two edits that both read version 1. The first
// wins and bumps the card to version 2; the second must get ErrVersionConflict.
func TestUpdateCardVersion(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	b := seedBoard(t, s)

	ada := models.Card{ColumnID: b.todo, Title: "Ada's title", Version: 1}
	if ok, err := s.UpdateCard(ctx, b.cardID, ada, b.ownerID); !ok || err != nil {
		t.Fatalf("first save: ok=%v err=%v", ok, err)
	}

	bo := models.Card{ColumnID: b.todo, Title: "Bo's title", Version: 1} // stale: read before Ada saved
	if _, err := s.UpdateCard(ctx, b.cardID, bo, b.ownerID); !errors.Is(err, ErrVersionConflict) {
		t.Fatalf("second save: want ErrVersionConflict, got %v", err)
	}

	card, err := s.GetCardByID(ctx, b.cardID, b.ownerID)
	if err != nil {
		t.Fatal(err)
	}
	if card.Title != "Ada's title" || card.Version != 2 {
		t.Fatalf("want Ada's title at version 2, got %q at version %d", card.Title, card.Version)
	}
}

// TestMoveCardBumpsVersion: an edit read before a move must conflict, or
// saving it would put the card back in its old column.
func TestMoveCardBumpsVersion(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	b := seedBoard(t, s)

	if _, err := s.MoveCard(ctx, b.cardID, b.doing, b.ownerID); err != nil {
		t.Fatal(err)
	}
	staleEdit := models.Card{ColumnID: b.todo, Title: "edited", Version: 1} // read while the card was in todo
	if _, err := s.UpdateCard(ctx, b.cardID, staleEdit, b.ownerID); !errors.Is(err, ErrVersionConflict) {
		t.Fatalf("want ErrVersionConflict, got %v", err)
	}
}
