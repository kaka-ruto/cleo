package qaassets

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoaderLoadsProfileEnvironmentAndRunbook(t *testing.T) {
	root := t.TempDir()
	profiles := filepath.Join(root, "profiles")
	runbooks := filepath.Join(root, "runbooks")
	envs := filepath.Join(root, "environments")
	for _, dir := range []string{profiles, runbooks, envs} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(profiles, "core.yml"), []byte("name: core\nrunbooks:\n  - local\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(runbooks, "local.yml"), []byte("name: local\nchecks:\n  - id: t1\n    title: test\n    command: go test ./...\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(envs, "local.yml"), []byte("name: local\nvars:\n  GOFLAGS: \"\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	loader := Loader{ProfilesDir: profiles, RunbooksDir: runbooks, EnvironmentsDir: envs}
	profile, err := loader.LoadProfile("core")
	if err != nil {
		t.Fatalf("LoadProfile error: %v", err)
	}
	if len(profile.Runbooks) != 1 || profile.Runbooks[0] != "local" {
		t.Fatalf("unexpected profile runbooks: %#v", profile.Runbooks)
	}
	if _, err := loader.LoadRunbook("local"); err != nil {
		t.Fatalf("LoadRunbook error: %v", err)
	}
	if _, err := loader.LoadEnvironment("local"); err != nil {
		t.Fatalf("LoadEnvironment error: %v", err)
	}
}
