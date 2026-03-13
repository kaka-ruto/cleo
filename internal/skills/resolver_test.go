package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveFallsBackToBuiltin(t *testing.T) {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatal(err)
	}
	r := Resolver{Cwd: dir, Home: home}
	src, body, err := r.Resolve("ceo")
	if err != nil {
		t.Fatalf("resolve ceo: %v", err)
	}
	if src.Origin != "builtin" {
		t.Fatalf("expected builtin origin, got %s", src.Origin)
	}
	if !strings.Contains(string(body), "name: plan-ceo-review") {
		t.Fatalf("expected imported ceo skill body")
	}
}

func TestResolvePrefersProjectOverride(t *testing.T) {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	if err := os.MkdirAll(filepath.Join(dir, ".cleo", "skills", "ceo"), 0o755); err != nil {
		t.Fatal(err)
	}
	overridePath := filepath.Join(dir, ".cleo", "skills", "ceo", "SKILL.md")
	body := "---\nname: ceo\ndescription: project ceo\n---\n\nproject"
	if err := os.WriteFile(overridePath, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	r := Resolver{Cwd: dir, Home: home}
	src, got, err := r.Resolve("ceo")
	if err != nil {
		t.Fatal(err)
	}
	if src.Origin != "project" {
		t.Fatalf("expected project origin, got %s", src.Origin)
	}
	if string(got) != body {
		t.Fatalf("expected override body")
	}
}

func TestCustomizeWritesProjectOverride(t *testing.T) {
	dir := t.TempDir()
	r := Resolver{Cwd: dir, Home: filepath.Join(dir, "home")}
	path, err := r.Customize("ceo")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(path, ".cleo/skills/ceo/SKILL.md") {
		t.Fatalf("unexpected path: %s", path)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "name: plan-ceo-review") {
		t.Fatal("expected imported ceo content")
	}
}

func TestResolveFounderNotFound(t *testing.T) {
	dir := t.TempDir()
	r := Resolver{Cwd: dir, Home: filepath.Join(dir, "home")}
	_, _, err := r.Resolve("founder")
	if err == nil || !strings.Contains(err.Error(), "skill not found") {
		t.Fatalf("expected not found for founder, got %v", err)
	}
}

func TestCheckRejectsMissingFrontmatter(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".cleo", "skills", "bad"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".cleo", "skills", "bad", "SKILL.md"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	r := Resolver{Cwd: dir, Home: filepath.Join(dir, "home")}
	_, err := r.Check("bad")
	if err == nil || !strings.Contains(err.Error(), "missing frontmatter start") {
		t.Fatalf("expected frontmatter error, got %v", err)
	}
}
