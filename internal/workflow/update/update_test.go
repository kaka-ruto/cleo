package update

import (
	"os"
	"path/filepath"
	"testing"
)

func TestChecksumForAsset(t *testing.T) {
	content := "abc123  cleo_v1.0.0_linux_amd64.tar.gz\nfff111  checksums.txt\n"
	v, err := checksumForAsset(content, "cleo_v1.0.0_linux_amd64.tar.gz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "abc123" {
		t.Fatalf("unexpected checksum: %s", v)
	}
}

func TestFileSHA256(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "x.txt")
	if err := os.WriteFile(path, []byte("abc"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	sum, err := fileSHA256(path)
	if err != nil {
		t.Fatalf("hash file: %v", err)
	}
	if sum != "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad" {
		t.Fatalf("unexpected hash: %s", sum)
	}
}

func TestFindAssetURL(t *testing.T) {
	assets := []githubAsset{
		{Name: "checksums.txt", BrowserDownloadURL: "https://example/checksums.txt"},
	}
	u, err := findAssetURL(assets, "checksums.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u == "" {
		t.Fatal("expected url")
	}
}
