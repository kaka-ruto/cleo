package setup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureConfigDefaultsAddsMissingQAKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cleo.yml")
	body := "version: 1\ngithub:\n  owner: cafaye\n  repo: cleo\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	changed, err := ensureConfigDefaults(path)
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("expected config migration to change file")
	}
	out, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(out)
	for _, needle := range []string{"actors_dir: .cleo/qa/actors", "evidence_dir: .cleo/evidence", "enabled: true"} {
		if !strings.Contains(text, needle) {
			t.Fatalf("expected migrated config to contain %q", needle)
		}
	}
}

func TestEnsureConfigDefaultsNoChangeWhenPresent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cleo.yml")
	body := `version: 1
github:
  owner: cafaye
  repo: cleo
qa:
  actors_dir: .cleo/qa/actors
  evidence_dir: .cleo/evidence
  manual:
    enabled: true
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	changed, err := ensureConfigDefaults(path)
	if err != nil {
		t.Fatal(err)
	}
	if changed {
		t.Fatal("expected no config migration changes")
	}
}
