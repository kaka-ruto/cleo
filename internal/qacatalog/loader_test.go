package qacatalog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoaderLoadsActorEnvironmentAndRunbook(t *testing.T) {
	root := t.TempDir()
	actors := filepath.Join(root, "actors")
	runbooks := filepath.Join(root, "runbooks")
	envs := filepath.Join(root, "environments")
	for _, dir := range []string{actors, runbooks, envs} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(actors, "core.yml"), []byte("name: core\nrunbooks:\n  - local\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(runbooks, "local.yml"), []byte("name: local\nchecks:\n  - id: t1\n    title: test\n    goal: validate behavior\n    how_to_test: run deterministic checks\n    expected_result: no regressions\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(envs, "local.yml"), []byte("name: local\nvars:\n  GOFLAGS: \"\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	loader := Loader{ActorsDir: actors, RunbooksDir: runbooks, EnvironmentsDir: envs}
	actor, err := loader.LoadActor("core")
	if err != nil {
		t.Fatalf("LoadActor error: %v", err)
	}
	if len(actor.Runbooks) != 1 || actor.Runbooks[0] != "local" {
		t.Fatalf("unexpected actor runbooks: %#v", actor.Runbooks)
	}
	if _, err := loader.LoadRunbook("local"); err != nil {
		t.Fatalf("LoadRunbook error: %v", err)
	}
	if _, err := loader.LoadEnvironment("local"); err != nil {
		t.Fatalf("LoadEnvironment error: %v", err)
	}
}
