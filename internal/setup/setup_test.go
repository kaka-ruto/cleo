package setup

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteConfigStatesNoConfigFile(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	wizard := &Wizard{Stdout: w}
	if err := wizard.writeConfig(); err != nil {
		t.Fatal(err)
	}
	_ = w.Close()

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	text := string(out)
	if !strings.Contains(text, "No cleo.yml file is used.") {
		t.Fatalf("expected no-config message, got %q", text)
	}
	if !strings.Contains(text, "infers repo context from git") {
		t.Fatalf("expected git inference message, got %q", text)
	}
}

func TestPathContains(t *testing.T) {
	original := os.Getenv("PATH")
	t.Cleanup(func() { _ = os.Setenv("PATH", original) })
	if err := os.Setenv("PATH", "/usr/bin:/tmp/cleo/bin"); err != nil {
		t.Fatal(err)
	}
	if !pathContains("/tmp/cleo/bin") {
		t.Fatal("expected PATH match")
	}
	if pathContains("/tmp/missing/bin") {
		t.Fatal("did not expect PATH match")
	}
}

func TestCopyExecutableCopiesFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dst := filepath.Join(dir, "dst")
	if err := os.WriteFile(src, []byte("cleo-binary"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := copyExecutable(src, dst); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "cleo-binary" {
		t.Fatalf("unexpected copied data: %s", string(body))
	}
}
