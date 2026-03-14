package skill

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cafaye/cleo/internal/skills"
)

func TestExecuteListIncludesCEO(t *testing.T) {
	var out bytes.Buffer
	r := skills.Resolver{Cwd: t.TempDir(), Home: t.TempDir()}
	cmd := newForTest(&out, r)
	if err := cmd.Execute("list", nil); err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(out.String(), "ceo") {
		t.Fatalf("expected ceo in list, got: %s", out.String())
	}
}

func TestExecuteUsePrintsSkill(t *testing.T) {
	var out bytes.Buffer
	r := skills.Resolver{Cwd: t.TempDir(), Home: t.TempDir()}
	cmd := newForTest(&out, r)
	if err := cmd.Execute("use", []string{"ceo"}); err != nil {
		t.Fatalf("use: %v", err)
	}
	if !strings.Contains(out.String(), "# Mega Plan Review Mode") {
		t.Fatalf("expected imported ceo skill body, got: %s", out.String())
	}
}

func TestExecuteCustomizeWritesProjectOverride(t *testing.T) {
	dir := t.TempDir()
	var out bytes.Buffer
	r := skills.Resolver{Cwd: dir, Home: t.TempDir()}
	cmd := newForTest(&out, r)
	if err := cmd.Execute("customize", []string{"ceo"}); err != nil {
		t.Fatalf("customize: %v", err)
	}
	path := filepath.Join(dir, ".agents", "skills", "ceo", "SKILL.md")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %s: %v", path, err)
	}
}

func TestExecuteUnknownCommand(t *testing.T) {
	var out bytes.Buffer
	r := skills.Resolver{Cwd: t.TempDir(), Home: t.TempDir()}
	cmd := newForTest(&out, r)
	err := cmd.Execute("wat", nil)
	if err == nil || !strings.Contains(err.Error(), "unknown skill command") {
		t.Fatalf("expected unknown command error, got %v", err)
	}
}

func TestExecuteInstallGlobalWritesSkill(t *testing.T) {
	var out bytes.Buffer
	home := t.TempDir()
	r := skills.Resolver{Cwd: t.TempDir(), Home: home}
	cmd := newForTest(&out, r)
	if err := cmd.Execute("install", []string{"ceo", "--global"}); err != nil {
		t.Fatalf("install: %v", err)
	}
	path := filepath.Join(home, ".agents", "skills", "ceo", "SKILL.md")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %s: %v", path, err)
	}
}

func TestExecuteSyncProjectWritesBuiltins(t *testing.T) {
	var out bytes.Buffer
	cwd := t.TempDir()
	r := skills.Resolver{Cwd: cwd, Home: t.TempDir()}
	cmd := newForTest(&out, r)
	if err := cmd.Execute("sync", []string{"--project"}); err != nil {
		t.Fatalf("sync: %v", err)
	}
	path := filepath.Join(cwd, ".agents", "skills", "ceo", "SKILL.md")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %s: %v", path, err)
	}
}
