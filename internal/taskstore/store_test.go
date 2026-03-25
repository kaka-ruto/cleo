package taskstore

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func TestUpsertOpenTaskDedupesByRepoAndKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cleo.db")
	store, err := Open(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer func() { _ = store.Close() }()

	now := time.Date(2026, 3, 3, 12, 0, 0, 0, time.UTC)
	session, err := store.StartSession(context.Background(), "pr", "123", "checkout regression", "", now)
	if err != nil {
		t.Fatalf("start session: %v", err)
	}

	first, created, err := store.UpsertOpenTask(context.Background(), Task{
		SessionID: session.ID,
		RepoKey:   "github.com/kaka-ruto/cleo",
		Title:     "Checkout button fails",
		Details:   "clicking checkout returns 500",
		Severity:  "high",
		DedupeKey: "abc",
	}, now)
	if err != nil {
		t.Fatalf("first upsert: %v", err)
	}
	if !created {
		t.Fatalf("expected first insert to create task")
	}

	second, created, err := store.UpsertOpenTask(context.Background(), Task{
		SessionID: session.ID,
		RepoKey:   "github.com/kaka-ruto/cleo",
		Title:     "Checkout button fails",
		Details:   "clicking checkout returns 500",
		Severity:  "high",
		DedupeKey: "abc",
	}, now.Add(1*time.Minute))
	if err != nil {
		t.Fatalf("second upsert: %v", err)
	}
	if created {
		t.Fatalf("expected second upsert to dedupe")
	}
	if first.ID != second.ID {
		t.Fatalf("expected same task id, got first=%d second=%d", first.ID, second.ID)
	}
	if second.Occurrences != 2 {
		t.Fatalf("expected occurrences=2, got %d", second.Occurrences)
	}
}
