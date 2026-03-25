package state

import (
	"path/filepath"
	"testing"

	"github.com/kaka-ruto/cleo/internal/config"
)

func TestDBPathUsesRepoScopedHierarchy(t *testing.T) {
	t.Setenv("HOME", "/tmp/testhome")
	cfg := &config.Config{}
	cfg.GitHub.Host = "github.com"
	cfg.GitHub.Owner = "kaka-ruto"
	cfg.GitHub.Repo = "cleo"
	path, err := DBPath(cfg)
	if err != nil {
		t.Fatalf("DBPath error: %v", err)
	}
	expected := filepath.Join("/tmp/testhome", ".cleo", "state", "repos", "github.com", "kaka-ruto", "cleo", "cleo.db")
	if path != expected {
		t.Fatalf("expected %q, got %q", expected, path)
	}
}
