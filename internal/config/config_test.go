package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadProjectAppliesDefaultsWithGitInference(t *testing.T) {
	dir := t.TempDir()
	initGitRepo(t, dir)
	configureGitRemote(t, dir, "git@github.com:kaka-ruto/cleo.git")
	setOriginHead(t, dir, "main")
	t.Chdir(dir)

	cfg, err := LoadProject()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.GitHub.Host != "github.com" {
		t.Fatalf("expected inferred host github.com, got %s", cfg.GitHub.Host)
	}
	if cfg.GitHub.Owner != "kaka-ruto" || cfg.GitHub.Repo != "cleo" {
		t.Fatalf("expected inferred repo kaka-ruto/cleo, got %s/%s", cfg.GitHub.Owner, cfg.GitHub.Repo)
	}
	if cfg.GitHub.BaseBranch != "main" {
		t.Fatalf("expected inferred base branch main, got %s", cfg.GitHub.BaseBranch)
	}
	if cfg.PR.RequiredApprovals != 1 {
		t.Fatalf("expected required approvals default 1, got %d", cfg.PR.RequiredApprovals)
	}
	if !cfg.PR.RequireNonDraft || !cfg.PR.RequireMergeable || !cfg.PR.BlockRequestedChanges {
		t.Fatal("expected strict PR gate defaults enabled")
	}
	if !cfg.PR.DeployWatch.Enabled {
		t.Fatal("expected deploy watch enabled by default")
	}
	if !cfg.PR.PostMerge.Enabled || !cfg.PR.PostMerge.AllowNone {
		t.Fatal("expected post-merge defaults enabled")
	}
	if len(cfg.PR.PostMerge.CommandAllowlistPrefixes) != 1 || cfg.PR.PostMerge.CommandAllowlistPrefixes[0] != "bin/kamal" {
		t.Fatalf("expected default command allowlist [bin/kamal], got %#v", cfg.PR.PostMerge.CommandAllowlistPrefixes)
	}
	if len(cfg.PR.PostMerge.CommandDenylist) != 1 || cfg.PR.PostMerge.CommandDenylist[0] != "rm -rf /" {
		t.Fatalf("expected default command denylist [rm -rf /], got %#v", cfg.PR.PostMerge.CommandDenylist)
	}
	if !cfg.PR.Stack.RebaseNextAfterMerge || !cfg.PR.Stack.AutoDetectNextPR || !cfg.PR.Stack.ForceWithLease {
		t.Fatal("expected PR stack defaults enabled")
	}
	if !cfg.Safety.RequireExplicitApply {
		t.Fatal("expected safety.require_explicit_apply enabled by default")
	}
	if !cfg.Release.GenerateNotes {
		t.Fatal("expected release.generate_notes enabled by default")
	}
	if !cfg.QAManualEnabled() {
		t.Fatal("expected qa manual mode enabled by default")
	}
	if cfg.QAEvidenceDir() != ".cleo/evidence" {
		t.Fatalf("expected default QA evidence dir .cleo/evidence, got %s", cfg.QAEvidenceDir())
	}
	if len(cfg.QA.DefaultActors) != 1 || cfg.QA.DefaultActors[0] != "core" {
		t.Fatalf("expected default QA actors [core], got %#v", cfg.QA.DefaultActors)
	}
}

func TestLoadProjectFailsWhenRepoNotInferable(t *testing.T) {
	dir := t.TempDir()
	initGitRepo(t, dir)
	t.Chdir(dir)

	_, err := LoadProject()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "configure git remote.origin.url") {
		t.Fatalf("expected remote origin guidance, got %v", err)
	}
}

func TestLoadProjectIgnoresLegacyCleoYml(t *testing.T) {
	dir := t.TempDir()
	initGitRepo(t, dir)
	configureGitRemote(t, dir, "https://github.com/kaka-ruto/cleo.git")
	setOriginHead(t, dir, "main")
	t.Chdir(dir)

	badConfig := "not: [valid:\n"
	if err := os.WriteFile(filepath.Join(dir, "cleo.yml"), []byte(badConfig), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadProject()
	if err != nil {
		t.Fatalf("expected cleo.yml to be ignored, got error: %v", err)
	}
	if cfg.GitHub.Owner != "kaka-ruto" || cfg.GitHub.Repo != "cleo" {
		t.Fatalf("expected inferred repo kaka-ruto/cleo, got %s/%s", cfg.GitHub.Owner, cfg.GitHub.Repo)
	}
}

func TestLoadProjectDefaultsBuildTargetForRootMainGo(t *testing.T) {
	dir := t.TempDir()
	initGitRepo(t, dir)
	configureGitRemote(t, dir, "git@github.com:cafaye/cafaye-cli.git")
	setOriginHead(t, dir, "main")
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\nfunc main(){}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)

	cfg, err := LoadProject()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Release.BuildTarget != "." {
		t.Fatalf("expected inferred build target '.', got %s", cfg.Release.BuildTarget)
	}
}

func TestLoadProjectDefaultsBuildTargetForCleoLayout(t *testing.T) {
	dir := t.TempDir()
	initGitRepo(t, dir)
	configureGitRemote(t, dir, "git@github.com:kaka-ruto/cleo.git")
	setOriginHead(t, dir, "master")
	if err := os.MkdirAll(filepath.Join(dir, "cmd", "cleo"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)

	cfg, err := LoadProject()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Release.BuildTarget != "./cmd/cleo" {
		t.Fatalf("expected inferred build target './cmd/cleo', got %s", cfg.Release.BuildTarget)
	}
}

func TestLoadProjectDefaultsBinaryNameFromRepo(t *testing.T) {
	dir := t.TempDir()
	initGitRepo(t, dir)
	configureGitRemote(t, dir, "git@github.com:cafaye/cafaye-cli.git")
	setOriginHead(t, dir, "master")
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\nfunc main(){}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)

	cfg, err := LoadProject()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Release.BinaryName != "cafaye" {
		t.Fatalf("expected inferred binary name 'cafaye', got %s", cfg.Release.BinaryName)
	}
}

func TestInferDefaultBinaryName(t *testing.T) {
	tests := []struct {
		repo string
		want string
	}{
		{repo: "cleo", want: "cleo"},
		{repo: "cafaye-cli", want: "cafaye"},
		{repo: "mytool", want: "mytool"},
		{repo: "", want: "cleo"},
	}
	for _, tc := range tests {
		if got := inferDefaultBinaryName(tc.repo); got != tc.want {
			t.Fatalf("inferDefaultBinaryName(%q)=%q want %q", tc.repo, got, tc.want)
		}
	}
}

func TestParseRepoFromRemoteURL(t *testing.T) {
	tests := []struct {
		name  string
		raw   string
		host  string
		owner string
		repo  string
	}{
		{name: "ssh", raw: "git@github.com:kaka-ruto/cleo.git", host: "github.com", owner: "kaka-ruto", repo: "cleo"},
		{name: "https", raw: "https://github.com/kaka-ruto/cleo.git", host: "github.com", owner: "kaka-ruto", repo: "cleo"},
		{name: "enterprise", raw: "git@ghe.local:team/repo.git", host: "ghe.local", owner: "team", repo: "repo"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			host, owner, repo, err := parseRepoFromRemoteURL(tc.raw)
			if err != nil {
				t.Fatalf("parseRepoFromRemoteURL error: %v", err)
			}
			if host != tc.host || owner != tc.owner || repo != tc.repo {
				t.Fatalf("expected %s/%s/%s, got %s/%s/%s", tc.host, tc.owner, tc.repo, host, owner, repo)
			}
		})
	}
}

func TestParseRepoFromRemoteURLRejectsInvalid(t *testing.T) {
	for _, raw := range []string{"", "github.com/kaka-ruto/cleo", "https://github.com/kaka-ruto"} {
		if _, _, _, err := parseRepoFromRemoteURL(raw); err == nil {
			t.Fatalf("expected parse failure for %q", raw)
		}
	}
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	runCmd(t, dir, "git", "init")
	runCmd(t, dir, "git", "config", "user.name", "Cleo Test")
	runCmd(t, dir, "git", "config", "user.email", "cleo@example.com")
}

func configureGitRemote(t *testing.T, dir, remoteURL string) {
	t.Helper()
	runCmd(t, dir, "git", "remote", "add", "origin", remoteURL)
}

func setOriginHead(t *testing.T, dir, branch string) {
	t.Helper()
	runCmd(t, dir, "git", "symbolic-ref", "refs/remotes/origin/HEAD", fmt.Sprintf("refs/remotes/origin/%s", branch))
}

func runCmd(t *testing.T, dir, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s %v failed: %v (%s)", name, args, err, string(out))
	}
}
