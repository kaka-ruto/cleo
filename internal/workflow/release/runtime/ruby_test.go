package runtime

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetectRuby(t *testing.T) {
	dir := t.TempDir()
	if DetectRuby(dir) {
		t.Fatal("expected false when gemspec is missing")
	}
	if err := os.WriteFile(filepath.Join(dir, "demo.gemspec"), []byte("Gem::Specification.new do |s| end\n"), 0o644); err != nil {
		t.Fatalf("write gemspec: %v", err)
	}
	if !DetectRuby(dir) {
		t.Fatal("expected true when gemspec exists")
	}
}

func TestVerifyRubyAssets(t *testing.T) {
	if err := VerifyRubyAssets([]string{"checksums.txt", "demo-1.0.0.gem"}); err != nil {
		t.Fatalf("expected valid ruby assets: %v", err)
	}
	if err := VerifyRubyAssets([]string{"checksums.txt"}); err == nil {
		t.Fatal("expected missing gem error")
	}
	if err := VerifyRubyAssets([]string{"demo-1.0.0.gem"}); err == nil {
		t.Fatal("expected missing checksums error")
	}
}

func TestBuildRubyReleaseArtifacts(t *testing.T) {
	dir := t.TempDir()
	gemspec := `Gem::Specification.new do |spec|
  spec.name = "demo_gem"
  spec.version = "1.2.3"
  spec.summary = "demo"
  spec.authors = ["test"]
  spec.files = []
end
`
	if err := os.WriteFile(filepath.Join(dir, "demo_gem.gemspec"), []byte(gemspec), 0o644); err != nil {
		t.Fatalf("write gemspec: %v", err)
	}

	assets, err := BuildRubyReleaseArtifacts(dir, "v1.2.3")
	if err != nil {
		t.Fatalf("build ruby release artifacts: %v", err)
	}
	if len(assets) != 2 {
		t.Fatalf("expected 2 assets, got %d", len(assets))
	}

	gemFile := filepath.Base(assets[0])
	if !strings.HasSuffix(gemFile, ".gem") {
		t.Fatalf("expected first asset to be .gem, got %s", gemFile)
	}
	if _, err := os.Stat(assets[0]); err != nil {
		t.Fatalf("expected gem artifact to exist: %v", err)
	}
	if _, err := os.Stat(assets[1]); err != nil {
		t.Fatalf("expected checksums artifact to exist: %v", err)
	}
	checksums, err := os.ReadFile(assets[1])
	if err != nil {
		t.Fatalf("read checksums: %v", err)
	}
	if !strings.Contains(string(checksums), gemFile) {
		t.Fatalf("checksums missing gem file entry: %s", string(checksums))
	}
}
