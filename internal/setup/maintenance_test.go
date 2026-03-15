package setup

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestApplyPostUpdateMigrationsEnsuresQAKitWithoutConfig(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	if err := ApplyPostUpdateMigrations(nil); err != nil {
		t.Fatalf("ApplyPostUpdateMigrations() error = %v", err)
	}

	for _, p := range []string{
		filepath.Join(dir, ".github", "workflows", "qa.yml"),
		filepath.Join(dir, ".github", "pull_request_template.md"),
		filepath.Join(dir, ".cleo", "qa", "actors", "core.yml"),
	} {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("expected %s to exist: %v", p, err)
		}
	}

	if _, err := os.Stat(filepath.Join(dir, "cleo.yml")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected no cleo.yml file to be created, got err=%v", err)
	}
}
