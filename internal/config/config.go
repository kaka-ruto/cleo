package config

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Config struct {
	Version int `yaml:"version"`
	GitHub  struct {
		Owner               string `yaml:"owner"`
		Repo                string `yaml:"repo"`
		Host                string `yaml:"host"`
		BaseBranch          string `yaml:"base_branch"`
		MergeMethod         string `yaml:"merge_method"`
		DeleteBranchOnMerge bool   `yaml:"delete_branch_on_merge"`
	} `yaml:"github"`
	PR struct {
		RequiredApprovals     int  `yaml:"required_approvals"`
		RequireNonDraft       bool `yaml:"require_non_draft"`
		RequireMergeable      bool `yaml:"require_mergeable"`
		BlockRequestedChanges bool `yaml:"block_if_requested_changes"`
		AllowSelfApproval     bool `yaml:"allow_self_approval"`
		Checks                struct {
			Mode                string   `yaml:"mode"`
			Required            []string `yaml:"required"`
			Ignore              []string `yaml:"ignore"`
			TreatNeutralAsPass  bool     `yaml:"treat_neutral_as_pass"`
			TimeoutSeconds      int      `yaml:"timeout_seconds"`
			PollIntervalSeconds int      `yaml:"poll_interval_seconds"`
		} `yaml:"checks"`
		DeployWatch struct {
			Enabled             bool   `yaml:"enabled"`
			Workflow            string `yaml:"workflow"`
			Branch              string `yaml:"branch"`
			TimeoutSeconds      int    `yaml:"timeout_seconds"`
			PollIntervalSeconds int    `yaml:"poll_interval_seconds"`
		} `yaml:"deploy_watch"`
		PostMerge struct {
			Enabled                  bool     `yaml:"enabled"`
			RequireCommandAllowlist  bool     `yaml:"require_command_allowlist"`
			CommandAllowlistPrefixes []string `yaml:"command_allowlist_prefixes"`
			CommandDenylist          []string `yaml:"command_denylist"`
			Markers                  struct {
				Start string `yaml:"start"`
				End   string `yaml:"end"`
			} `yaml:"markers"`
			AllowNone bool `yaml:"allow_none"`
		} `yaml:"post_merge"`
		Stack struct {
			RebaseNextAfterMerge bool `yaml:"rebase_next_after_merge"`
			AutoDetectNextPR     bool `yaml:"auto_detect_next_pr"`
			ForceWithLease       bool `yaml:"force_with_lease"`
		} `yaml:"stack"`
	} `yaml:"pr"`
	Safety struct {
		RequireExplicitApply bool `yaml:"require_explicit_apply"`
	} `yaml:"safety"`
	Release struct {
		TagPrefix     string `yaml:"tag_prefix"`
		ChangelogFile string `yaml:"changelog_file"`
		BinaryName    string `yaml:"binary_name"`
		BuildTarget   string `yaml:"build_target"`
		GenerateNotes bool   `yaml:"generate_notes"`
		DefaultDraft  bool   `yaml:"default_draft"`
	} `yaml:"release"`
	QA struct {
		ActorsDir     string   `yaml:"actors_dir"`
		EvidenceDir   string   `yaml:"evidence_dir"`
		DefaultActors []string `yaml:"default_actors"`
		Manual        struct {
			Enabled *bool `yaml:"enabled"`
		} `yaml:"manual"`
	} `yaml:"qa"`
}

func LoadProject() (*Config, error) {
	cfg := &Config{}
	cfg.inferFromGit()
	cfg.applyDefaults()
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) inferFromGit() {
	host, owner, repo, err := discoverRepoFromGitRemote()
	if err == nil {
		if c.GitHub.Host == "" {
			c.GitHub.Host = host
		}
		if c.GitHub.Owner == "" {
			c.GitHub.Owner = owner
		}
		if c.GitHub.Repo == "" {
			c.GitHub.Repo = repo
		}
	}
	if c.GitHub.BaseBranch == "" {
		if branch, branchErr := discoverBaseBranchFromGit(); branchErr == nil {
			c.GitHub.BaseBranch = branch
		}
	}
}

func discoverRepoFromGitRemote() (host, owner, repo string, err error) {
	out, err := exec.Command("git", "config", "--get", "remote.origin.url").CombinedOutput()
	if err != nil {
		return "", "", "", err
	}
	return parseRepoFromRemoteURL(strings.TrimSpace(string(out)))
}

func parseRepoFromRemoteURL(raw string) (host, owner, repo string, err error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", "", "", fmt.Errorf("empty git remote origin url")
	}
	pathPart := ""
	if strings.Contains(raw, "://") {
		u, parseErr := url.Parse(raw)
		if parseErr != nil {
			return "", "", "", parseErr
		}
		host = strings.TrimSpace(u.Hostname())
		pathPart = strings.TrimPrefix(u.Path, "/")
	} else {
		withoutUser := raw
		if at := strings.Index(withoutUser, "@"); at >= 0 {
			withoutUser = withoutUser[at+1:]
		}
		split := strings.SplitN(withoutUser, ":", 2)
		if len(split) != 2 {
			return "", "", "", fmt.Errorf("invalid git remote origin url: %s", raw)
		}
		host = strings.TrimSpace(split[0])
		pathPart = strings.TrimPrefix(split[1], "/")
	}
	pathPart = strings.TrimSuffix(pathPart, ".git")
	segments := strings.Split(pathPart, "/")
	if len(segments) < 2 {
		return "", "", "", fmt.Errorf("invalid git remote origin url: %s", raw)
	}
	owner = strings.TrimSpace(segments[len(segments)-2])
	repo = strings.TrimSpace(segments[len(segments)-1])
	if host == "" || owner == "" || repo == "" {
		return "", "", "", fmt.Errorf("invalid git remote origin url: %s", raw)
	}
	return host, owner, repo, nil
}

func discoverBaseBranchFromGit() (string, error) {
	out, err := exec.Command("git", "symbolic-ref", "--short", "refs/remotes/origin/HEAD").CombinedOutput()
	if err == nil {
		branch := parseBaseBranch(strings.TrimSpace(string(out)))
		if branch != "" {
			return branch, nil
		}
	}
	remoteShow, showErr := exec.Command("git", "remote", "show", "origin").CombinedOutput()
	if showErr != nil {
		return "", showErr
	}
	for _, line := range strings.Split(string(remoteShow), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "HEAD branch:") {
			branch := strings.TrimSpace(strings.TrimPrefix(line, "HEAD branch:"))
			if branch != "" && branch != "(unknown)" {
				return branch, nil
			}
		}
	}
	return "", fmt.Errorf("base branch not detected from git")
}

func parseBaseBranch(ref string) string {
	ref = strings.TrimSpace(ref)
	ref = strings.TrimPrefix(ref, "refs/remotes/origin/")
	ref = strings.TrimPrefix(ref, "origin/")
	return strings.TrimSpace(ref)
}

func (c *Config) applyDefaults() {
	if c.Version == 0 {
		c.Version = 1
	}
	if c.GitHub.Host == "" {
		c.GitHub.Host = "github.com"
	}
	if c.GitHub.BaseBranch == "" {
		c.GitHub.BaseBranch = "master"
	}
	if c.GitHub.MergeMethod == "" {
		c.GitHub.MergeMethod = "merge"
	}
	if c.PR.RequiredApprovals == 0 {
		c.PR.RequiredApprovals = 1
	}
	c.PR.RequireNonDraft = true
	c.PR.RequireMergeable = true
	c.PR.BlockRequestedChanges = true
	if c.PR.Checks.Mode == "" {
		c.PR.Checks.Mode = "required"
	}
	if c.PR.Checks.TimeoutSeconds == 0 {
		c.PR.Checks.TimeoutSeconds = 1800
	}
	if c.PR.Checks.PollIntervalSeconds == 0 {
		c.PR.Checks.PollIntervalSeconds = 10
	}
	if c.PR.DeployWatch.Workflow == "" {
		c.PR.DeployWatch.Workflow = "Deploy to Production"
	}
	c.PR.DeployWatch.Enabled = true
	if c.PR.DeployWatch.Branch == "" {
		c.PR.DeployWatch.Branch = c.GitHub.BaseBranch
	}
	if c.PR.DeployWatch.TimeoutSeconds == 0 {
		c.PR.DeployWatch.TimeoutSeconds = 2700
	}
	if c.PR.DeployWatch.PollIntervalSeconds == 0 {
		c.PR.DeployWatch.PollIntervalSeconds = 10
	}
	if c.PR.PostMerge.Markers.Start == "" {
		c.PR.PostMerge.Markers.Start = "<!-- post-merge-commands:start -->"
	}
	if c.PR.PostMerge.Markers.End == "" {
		c.PR.PostMerge.Markers.End = "<!-- post-merge-commands:end -->"
	}
	c.PR.PostMerge.Enabled = true
	c.PR.PostMerge.AllowNone = true
	if len(c.PR.PostMerge.CommandAllowlistPrefixes) == 0 {
		c.PR.PostMerge.CommandAllowlistPrefixes = []string{"bin/kamal"}
	}
	if len(c.PR.PostMerge.CommandDenylist) == 0 {
		c.PR.PostMerge.CommandDenylist = []string{"rm -rf /"}
	}
	c.PR.Stack.RebaseNextAfterMerge = true
	c.PR.Stack.AutoDetectNextPR = true
	c.PR.Stack.ForceWithLease = true
	c.Safety.RequireExplicitApply = true
	if c.Release.TagPrefix == "" {
		c.Release.TagPrefix = "v"
	}
	if c.Release.ChangelogFile == "" {
		c.Release.ChangelogFile = "CHANGELOG.md"
	}
	if c.Release.BinaryName == "" {
		c.Release.BinaryName = inferDefaultBinaryName(c.GitHub.Repo)
	}
	if c.Release.BuildTarget == "" {
		c.Release.BuildTarget = inferDefaultBuildTarget()
	}
	c.Release.GenerateNotes = true
	if c.QA.ActorsDir == "" {
		c.QA.ActorsDir = ".cleo/qa/actors"
	}
	if c.QA.EvidenceDir == "" {
		c.QA.EvidenceDir = ".cleo/evidence"
	}
	if len(c.QA.DefaultActors) == 0 {
		c.QA.DefaultActors = []string{"core"}
	}
}

func (c *Config) QAManualEnabled() bool {
	if c.QA.Manual.Enabled == nil {
		return true
	}
	return *c.QA.Manual.Enabled
}

func (c *Config) QAEvidenceDir() string {
	return c.QA.EvidenceDir
}

func (c *Config) validate() error {
	if c.Version != 1 {
		return fmt.Errorf("version must be 1")
	}
	if c.GitHub.Owner == "" || c.GitHub.Repo == "" {
		return fmt.Errorf("github.owner and github.repo are required (configure git remote.origin.url)")
	}
	switch c.GitHub.MergeMethod {
	case "merge", "squash", "rebase":
	default:
		return fmt.Errorf("github.merge_method must be merge|squash|rebase")
	}
	if c.PR.Checks.Mode != "required" && c.PR.Checks.Mode != "all" {
		return fmt.Errorf("pr.checks.mode must be required|all")
	}
	if c.Release.TagPrefix == "" {
		return fmt.Errorf("release.tag_prefix is required")
	}
	if c.Release.BinaryName == "" {
		return fmt.Errorf("release.binary_name is required")
	}
	if c.Release.BuildTarget == "" {
		return fmt.Errorf("release.build_target is required")
	}
	return nil
}

func inferDefaultBuildTarget() string {
	if pathExists("cmd", "cleo") {
		return "./cmd/cleo"
	}
	if pathExists("main.go") {
		return "."
	}
	if pathExists("cmd", "main.go") {
		return "./cmd"
	}
	return "./cmd/cleo"
}

func inferDefaultBinaryName(repo string) string {
	name := strings.TrimSpace(repo)
	if name == "" {
		return "cleo"
	}
	if strings.HasSuffix(name, "-cli") && len(name) > len("-cli") {
		name = strings.TrimSuffix(name, "-cli")
	}
	return name
}

func pathExists(parts ...string) bool {
	_, err := os.Stat(filepath.Join(parts...))
	return err == nil
}
