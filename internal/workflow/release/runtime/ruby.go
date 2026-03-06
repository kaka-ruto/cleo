package runtime

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func DetectRuby(root string) bool {
	_, ok := RubyGemspecPath(root)
	return ok
}

func RubyGemspecPath(root string) (string, bool) {
	matches, err := filepath.Glob(filepath.Join(root, "*.gemspec"))
	if err != nil || len(matches) == 0 {
		return "", false
	}
	sort.Strings(matches)
	return matches[0], true
}

func BuildRubyReleaseArtifacts(root, version string) ([]string, error) {
	gemspec, ok := RubyGemspecPath(root)
	if !ok {
		return nil, fmt.Errorf("no gemspec found in %s", root)
	}
	distDir := filepath.Join(root, "dist", "release", version)
	if err := os.RemoveAll(distDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(distDir, 0o755); err != nil {
		return nil, err
	}

	gemName, err := expectedRubyGemName(gemspec, version)
	if err != nil {
		return nil, err
	}
	gemPath := filepath.Join(distDir, gemName)
	if _, err := runLocalDirEnv(root, []string{"BUNDLE_GEMFILE="}, "gem", "build", filepath.Base(gemspec), "--output", gemPath); err != nil {
		return nil, err
	}

	sum, err := sha256File(gemPath)
	if err != nil {
		return nil, err
	}
	checksumPath := filepath.Join(distDir, "checksums.txt")
	line := sum + "  " + filepath.Base(gemPath) + "\n"
	if err := os.WriteFile(checksumPath, []byte(line), 0o644); err != nil {
		return nil, err
	}
	return []string{gemPath, checksumPath}, nil
}

func VerifyRubyAssets(assetNames []string) error {
	hasGem := false
	hasChecksums := false
	for _, name := range assetNames {
		if strings.HasSuffix(strings.TrimSpace(name), ".gem") {
			hasGem = true
		}
		if strings.TrimSpace(name) == "checksums.txt" {
			hasChecksums = true
		}
	}
	if !hasGem {
		return fmt.Errorf("missing release asset: *.gem")
	}
	if !hasChecksums {
		return fmt.Errorf("missing release asset: checksums.txt")
	}
	return nil
}

func runLocalDirEnv(dir string, env []string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), env...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s %s: %s", name, strings.Join(args, " "), strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func expectedRubyGemName(gemspecPath, version string) (string, error) {
	base := strings.TrimSuffix(filepath.Base(gemspecPath), ".gemspec")
	v := strings.TrimPrefix(version, "v")
	if strings.TrimSpace(v) == "" {
		return "", fmt.Errorf("invalid release version: %s", version)
	}
	return fmt.Sprintf("%s-%s.gem", base, v), nil
}

func sha256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
