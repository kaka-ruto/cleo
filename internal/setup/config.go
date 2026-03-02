package setup

import (
	"fmt"
	"os"
	"strings"
)

func (w *Wizard) writeConfig() error {
	if _, err := os.Stat("cleo.yml"); err == nil {
		overwrite, err := w.confirm("cleo.yml already exists. Overwrite?")
		if err != nil {
			return err
		}
		if !overwrite {
			fmt.Fprintln(w.Stdout, "Keeping existing cleo.yml")
			return nil
		}
	}
	repo, err := discoverRepoSlug()
	if err != nil {
		return err
	}
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid GitHub repo slug: %s", repo)
	}
	if err := os.WriteFile("cleo.yml", []byte(defaultConfig(parts[0], parts[1])), 0o644); err != nil {
		return err
	}
	fmt.Fprintln(w.Stdout, "Wrote cleo.yml")
	return nil
}

func defaultConfig(owner, repo string) string {
	return fmt.Sprintf(`version: 1

github:
  owner: %s
  repo: %s
  host: github.com
  base_branch: main
  merge_method: merge
  delete_branch_on_merge: false

pr:
  required_approvals: 1
  require_non_draft: true
  require_mergeable: true
  block_if_requested_changes: true
  allow_self_approval: false

  checks:
    mode: required
    required: []
    ignore: []
    treat_neutral_as_pass: false
    timeout_seconds: 1800
    poll_interval_seconds: 10

  deploy_watch:
    enabled: true
    workflow: Deploy to Production
    branch: main
    timeout_seconds: 2700
    poll_interval_seconds: 10

  post_merge:
    enabled: true
    require_command_allowlist: false
    command_allowlist_prefixes:
      - "bin/kamal"
    command_denylist:
      - "rm -rf /"
    markers:
      start: "<!-- post-merge-commands:start -->"
      end: "<!-- post-merge-commands:end -->"
    allow_none: true

  stack:
    rebase_next_after_merge: true
    auto_detect_next_pr: true
    force_with_lease: true

safety:
  require_explicit_apply: true
`, owner, repo)
}
