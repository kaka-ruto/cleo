package state

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cafaye/cleo/internal/config"
)

func DBPath(cfg *config.Config) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	host := clean(cfg.GitHub.Host)
	owner := clean(cfg.GitHub.Owner)
	repo := clean(cfg.GitHub.Repo)
	if host == "" || owner == "" || repo == "" {
		return "", fmt.Errorf("github.host, github.owner, and github.repo are required for state path")
	}
	return filepath.Join(home, ".cleo", "state", "repos", host, owner, repo, "cleo.db"), nil
}

func RepoKey(cfg *config.Config) string {
	return strings.TrimSpace(cfg.GitHub.Host) + "/" + strings.TrimSpace(cfg.GitHub.Owner) + "/" + strings.TrimSpace(cfg.GitHub.Repo)
}

func clean(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "\\", "_")
	return s
}
