package setup

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
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
		filepath.Join(dir, ".agents", "skills", "cleo", "SKILL.md"),
	} {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("expected %s to exist: %v", p, err)
		}
	}

	body, err := os.ReadFile(filepath.Join(dir, ".agents", "skills", "cleo", "SKILL.md"))
	if err != nil {
		t.Fatalf("read cleo skill: %v", err)
	}
	if !strings.Contains(string(body), "name: cleo") {
		t.Fatalf("expected builtin cleo skill content, got: %s", string(body))
	}

	if _, err := os.Stat(filepath.Join(dir, "cleo.yml")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected no cleo.yml file to be created, got err=%v", err)
	}
}

func TestApplyPostUpdateMigrationsDoesNotOverwriteExistingCleoSkill(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	target := filepath.Join(dir, ".agents", "skills", "cleo", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatalf("mkdir skill dir: %v", err)
	}
	const custom = "custom-skill-body"
	if err := os.WriteFile(target, []byte(custom), 0o644); err != nil {
		t.Fatalf("seed custom skill: %v", err)
	}

	if err := ApplyPostUpdateMigrations(nil); err != nil {
		t.Fatalf("ApplyPostUpdateMigrations() error = %v", err)
	}

	body, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read skill: %v", err)
	}
	if string(body) != custom {
		t.Fatalf("expected existing skill to remain unchanged, got: %s", string(body))
	}
}
