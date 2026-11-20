package store

import (
	"context"
	"testing"
	"time"
)

// TestSessionLifecycle creates a session, reads its user back, then deletes
// it. It checks the sqlc-generated queries behind sessions.go.
func TestSessionLifecycle(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	b := seedBoard(t, s) // only for its user; the cleanup deletes the session too

	if err := s.CreateSession(ctx, "test-hash", b.ownerID, time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	user, err := s.GetSessionUser(ctx, "test-hash")
	if err != nil {
		t.Fatal(err)
	}
	if user.ID != b.ownerID || user.Name != "Race Test" {
		t.Fatalf("got user %q (%q), want %q", user.ID, user.Name, b.ownerID)
	}

	if err := s.DeleteSession(ctx, "test-hash"); err != nil {
		t.Fatal(err)
	}
	if user, err := s.GetSessionUser(ctx, "test-hash"); err != nil || user.ID != "" {
		t.Fatalf("after delete: got user %q, err %v; want none", user.ID, err)
	}
}
