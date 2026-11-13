package handlers

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"kanban/internal/models"
	"kanban/internal/store"
)

// memCardStore keeps cards in a map and does the same version check as
// PostgresStore.UpdateCard, but in Go memory instead of in the database.
// Requests run in their own goroutines, so mu guards cards: only one goroutine
// at a time can read or write it.
type memCardStore struct {
	store.Store
	mu    sync.Mutex
	cards map[string]models.Card
}

func (s *memCardStore) UpdateCard(_ context.Context, id string, updated models.Card, _ string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock() // every return unlocks, like defer tx.Rollback()

	// The read, the version check and the write happen under one lock, so no
	// other save can slip in between them.
	current, ok := s.cards[id]
	if !ok {
		return false, nil
	}
	if current.Version != updated.Version {
		return false, store.ErrVersionConflict
	}
	updated.ID = id
	updated.Version++
	s.cards[id] = updated
	return true, nil
}

func (s *memCardStore) GetCardByID(_ context.Context, id, _ string) (models.Card, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cards[id], nil
}

// TestUpdateCardConcurrent sends 10 saves of the same card at once, all with
// version 1, like 10 tabs that opened the editor together. Exactly one should
// win; the other 9 must get 409. Run it with go test -race.
func TestUpdateCardConcurrent(t *testing.T) {
	const cardID = "c2ee0000-0000-4000-8000-000000000001"
	mem := &memCardStore{cards: map[string]models.Card{
		cardID: {ID: cardID, ColumnID: "a5fc3aed-6138-4e57-a8ec-3059d343eb83", Title: "Fix bug", Version: 1},
	}}
	handler := CardHandler{Store: mem, Validate: testValidator()}

	const saves = 10
	codes := make([]int, saves) // each goroutine writes only its own index
	var wg sync.WaitGroup
	for i := range saves {
		wg.Add(1)
		go func() {
			defer wg.Done()
			body := fmt.Sprintf(`{"id":%q,"column_id":"a5fc3aed-6138-4e57-a8ec-3059d343eb83","title":"Save %d","order":0,"version":1}`, cardID, i)
			request := httptest.NewRequest("PUT", "/cards/"+cardID, strings.NewReader(body))
			request.SetPathValue("id", cardID)
			response := httptest.NewRecorder()
			handler.UpdateCard(response, asUser(request))
			codes[i] = response.Code
		}()
	}
	wg.Wait()

	wins, conflicts := 0, 0
	for _, code := range codes {
		switch code {
		case 200:
			wins++
		case 409:
			conflicts++
		}
	}
	if wins != 1 || conflicts != saves-1 {
		t.Fatalf("got %d wins and %d conflicts; want 1 and %d", wins, conflicts, saves-1)
	}
}
