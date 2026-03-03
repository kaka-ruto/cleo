package config

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
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

func Load(path string) (*Config, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	cfg := &Config{}
	dec := yaml.NewDecoder(bytes.NewReader(body))
	dec.KnownFields(true)
	if err := dec.Decode(cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	cfg.applyDefaults()
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
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
	if len(c.PR.PostMerge.CommandAllowlistPrefixes) == 0 {
		c.PR.PostMerge.CommandAllowlistPrefixes = []string{"bin/kamal"}
	}
	if c.Release.TagPrefix == "" {
		c.Release.TagPrefix = "v"
	}
	if c.Release.ChangelogFile == "" {
		c.Release.ChangelogFile = "CHANGELOG.md"
	}
	if c.Release.BinaryName == "" {
		c.Release.BinaryName = "cleo"
	}
	if c.Release.BuildTarget == "" {
		c.Release.BuildTarget = "./cmd/cleo"
	}
	if c.QA.ActorsDir == "" {
		c.QA.ActorsDir = ".cleo/qa/actors"
	}
	if c.QA.EvidenceDir == "" {
		c.QA.EvidenceDir = ".cleo/evidence"
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
		return fmt.Errorf("github.owner and github.repo are required")
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
