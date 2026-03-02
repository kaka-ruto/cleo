package runtime

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectGo(t *testing.T) {
	dir := t.TempDir()
	if DetectGo(dir) {
		t.Fatal("expected false when go.mod is missing")
	}
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	if !DetectGo(dir) {
		t.Fatal("expected true when go.mod exists")
	}
}

func TestExpectedGoAssetNames(t *testing.T) {
	names := ExpectedGoAssetNames("v1.2.3", "cleo")
	if len(names) != len(DefaultGoTargets)+1 {
		t.Fatalf("unexpected asset count: %d", len(names))
	}
	want := "cleo_v1.2.3_linux_amd64.tar.gz"
	found := false
	for _, n := range names {
		if n == want {
			found = true
		}
	}
	if !found {
		t.Fatalf("missing expected asset %q", want)
	}
}
