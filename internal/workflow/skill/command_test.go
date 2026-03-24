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

func TestExecuteUsePrintsCleoSkill(t *testing.T) {
	var out bytes.Buffer
	r := skills.Resolver{Cwd: t.TempDir(), Home: t.TempDir()}
	cmd := newForTest(&out, r)
	if err := cmd.Execute("use", []string{"cleo"}); err != nil {
		t.Fatalf("use: %v", err)
	}
	if !strings.Contains(out.String(), "# Cleo Workflow Skill") {
		t.Fatalf("expected imported cleo skill body, got: %s", out.String())
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
	cleoPath := filepath.Join(cwd, ".agents", "skills", "cleo", "SKILL.md")
	if _, err := os.Stat(cleoPath); err != nil {
		t.Fatalf("expected %s: %v", cleoPath, err)
	}
}

func TestExecuteUninstallGlobalRemovesSkill(t *testing.T) {
	var out bytes.Buffer
	home := t.TempDir()
	r := skills.Resolver{Cwd: t.TempDir(), Home: home}
	cmd := newForTest(&out, r)
	if err := cmd.Execute("install", []string{"ceo", "--global"}); err != nil {
		t.Fatalf("install: %v", err)
	}
	if err := cmd.Execute("uninstall", []string{"ceo", "--global"}); err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	path := filepath.Join(home, ".agents", "skills", "ceo")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected uninstall of %s", path)
	}
}

func TestExecuteRegistryList(t *testing.T) {
	var out bytes.Buffer
	r := skills.Resolver{Cwd: t.TempDir(), Home: t.TempDir()}
	cmd := newForTest(&out, r)
	if err := cmd.Execute("registry", []string{"list"}); err != nil {
		t.Fatalf("registry list: %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "openai") || !strings.Contains(got, "superpowers") || !strings.Contains(got, "superpowers-ruby") {
		t.Fatalf("missing expected registries: %s", got)
	}
}

func TestExecuteRegistryAddAndRemove(t *testing.T) {
	var out bytes.Buffer
	home := t.TempDir()
	r := skills.Resolver{Cwd: t.TempDir(), Home: home}
	cmd := newForTest(&out, r)
	if err := cmd.Execute("registry", []string{"add", "team", "--repo", "acme/skills", "--path", "skills"}); err != nil {
		t.Fatalf("registry add: %v", err)
	}
	out.Reset()
	if err := cmd.Execute("registry", []string{"list"}); err != nil {
		t.Fatalf("registry list: %v", err)
	}
	if !strings.Contains(out.String(), "team") {
		t.Fatalf("expected custom registry in list, got: %s", out.String())
	}
	out.Reset()
	if err := cmd.Execute("registry", []string{"remove", "team"}); err != nil {
		t.Fatalf("registry remove: %v", err)
	}
	out.Reset()
	if err := cmd.Execute("registry", []string{"list"}); err != nil {
		t.Fatalf("registry list after remove: %v", err)
	}
	if strings.Contains(out.String(), "team") {
		t.Fatalf("expected custom registry removed, got: %s", out.String())
	}
}
