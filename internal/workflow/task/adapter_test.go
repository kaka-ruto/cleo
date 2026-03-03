package task

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cafaye/cleo/internal/config"
	"github.com/cafaye/cleo/internal/taskstore"
)

func TestWorkFromBaseBranchCreatesTaskBranch(t *testing.T) {
	store := openTestStore(t)
	defer func() { _ = store.Close() }()
	id := seedTask(t, store)

	cfg := &config.Config{}
	cfg.GitHub.BaseBranch = "master"
	adapter := NewAdapter(store, cfg)
	current := "master"
	adapter.runLocalFn = fakeGitRunner(&current)

	text, err := adapter.Work(id, WorkOptions{})
	if err != nil {
		t.Fatalf("Work error: %v", err)
	}
	if !strings.Contains(text, "Work lane: new-branch") {
		t.Fatalf("expected new-branch lane, got %q", text)
	}
	task, err := store.Task(context.Background(), id)
	if err != nil {
		t.Fatalf("task lookup error: %v", err)
	}
	if task.WorkBranch == "" || task.WorkBranch == "master" {
		t.Fatalf("expected task work branch to be set, got %q", task.WorkBranch)
	}
}

func TestWorkFromFeatureBranchKeepsCurrentBranch(t *testing.T) {
	store := openTestStore(t)
	defer func() { _ = store.Close() }()
	id := seedTask(t, store)

	cfg := &config.Config{}
	cfg.GitHub.BaseBranch = "master"
	adapter := NewAdapter(store, cfg)
	current := "feature/checkout-fix"
	adapter.runLocalFn = fakeGitRunner(&current)

	text, err := adapter.Work(id, WorkOptions{})
	if err != nil {
		t.Fatalf("Work error: %v", err)
	}
	if !strings.Contains(text, "Work lane: in-place") {
		t.Fatalf("expected in-place lane, got %q", text)
	}
	task, err := store.Task(context.Background(), id)
	if err != nil {
		t.Fatalf("task lookup error: %v", err)
	}
	if task.WorkBranch != "feature/checkout-fix" {
		t.Fatalf("expected branch to remain feature branch, got %q", task.WorkBranch)
	}
}

func TestWorkWithForceNewBranchFromFeatureBranchCreatesTaskBranch(t *testing.T) {
	store := openTestStore(t)
	defer func() { _ = store.Close() }()
	id := seedTask(t, store)

	cfg := &config.Config{}
	cfg.GitHub.BaseBranch = "master"
	adapter := NewAdapter(store, cfg)
	current := "feature/checkout-fix"
	adapter.runLocalFn = fakeGitRunner(&current)

	text, err := adapter.Work(id, WorkOptions{ForceNewBranch: true})
	if err != nil {
		t.Fatalf("Work error: %v", err)
	}
	if !strings.Contains(text, "Work lane: new-branch") {
		t.Fatalf("expected new-branch lane, got %q", text)
	}
	task, err := store.Task(context.Background(), id)
	if err != nil {
		t.Fatalf("task lookup error: %v", err)
	}
	if task.WorkBranch == "feature/checkout-fix" {
		t.Fatalf("expected branch switch for force new branch, got %q", task.WorkBranch)
	}
}

func TestWorkWithForceInPlaceOnBaseBranchFails(t *testing.T) {
	store := openTestStore(t)
	defer func() { _ = store.Close() }()
	id := seedTask(t, store)

	cfg := &config.Config{}
	cfg.GitHub.BaseBranch = "master"
	adapter := NewAdapter(store, cfg)
	current := "master"
	adapter.runLocalFn = fakeGitRunner(&current)

	_, err := adapter.Work(id, WorkOptions{ForceInPlace: true})
	if err == nil {
		t.Fatal("expected force in-place to fail on base branch")
	}
	if !strings.Contains(err.Error(), "--in-place cannot be used on base branch") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func openTestStore(t *testing.T) *taskstore.Store {
	t.Helper()
	path := filepath.Join(t.TempDir(), "cleo.db")
	store, err := taskstore.Open(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	return store
}

func seedTask(t *testing.T, store *taskstore.Store) int64 {
	t.Helper()
	now := time.Date(2026, 3, 3, 12, 0, 0, 0, time.UTC)
	session, err := store.StartSession(context.Background(), "pr", "4", "qa", "", now)
	if err != nil {
		t.Fatalf("start session: %v", err)
	}
	task, _, err := store.UpsertOpenTask(context.Background(), taskstore.Task{
		SessionID: session.ID,
		RepoKey:   "github.com/cafaye/cleo",
		Title:     "Checkout fails",
		Details:   "500 on submit",
		Severity:  "high",
		DedupeKey: "seed",
	}, now)
	if err != nil {
		t.Fatalf("seed task: %v", err)
	}
	return task.ID
}

func fakeGitRunner(current *string) func(name string, args ...string) (string, error) {
	return func(name string, args ...string) (string, error) {
		if name != "git" {
			return "", nil
		}
		if len(args) == 2 && args[0] == "branch" && args[1] == "--show-current" {
			return *current + "\n", nil
		}
		if len(args) == 3 && args[0] == "checkout" && args[1] == "-b" {
			*current = args[2]
			return "", nil
		}
		if len(args) == 2 && args[0] == "checkout" {
			*current = args[1]
			return "", nil
		}
		return "", nil
	}
}
