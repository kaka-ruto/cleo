package pr

import (
	"fmt"
	"os"
	"strings"
)

func (s *Service) Create(title, summary, why, what, test, risk, rollback, owner, ac string, cmds []string, draft bool) error {
	branch, err := runLocal("git", "branch", "--show-current")
	if err != nil {
		return err
	}
	head := strings.TrimSpace(branch)
	if head == "" {
		return fmt.Errorf("cannot determine current branch")
	}
	if title == "" {
		title = summary
	}
	if title == "" {
		return fmt.Errorf("title or summary is required")
	}
	body := Render(summary, why, what, test, risk, rollback, owner, ac, cmds)
	tmp, err := os.CreateTemp("", "cleo-pr-body-*.md")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.WriteString(body); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	args := []string{"pr", "create", "--repo", s.repo(), "--base", s.cfg.GitHub.BaseBranch, "--head", head, "--title", title, "--body-file", tmp.Name()}
	if draft {
		args = append(args, "--draft")
	}
	_, err = s.gh.Run(args...)
	return err
}

func Render(summary, why, what, test, risk, rollback, owner, ac string, cmds []string) string {
	fields := withDefaults(summary, why, what, test, risk, rollback, owner)
	return fmt.Sprintf(prBodyTemplate, fields.summary, fields.why, fields.what, fields.test, fields.risk, fields.rollback, fields.owner, renderACBlock(ac), renderCommandLines(cmds))
}

type prFields struct {
	summary  string
	why      string
	what     string
	test     string
	risk     string
	rollback string
	owner    string
}

func withDefaults(summary, why, what, test, risk, rollback, owner string) prFields {
	return prFields{
		summary:  first(summary, "TBD"),
		why:      first(why, "TBD"),
		what:     first(what, "- TBD"),
		test:     first(test, "- TBD"),
		risk:     first(risk, "Low"),
		rollback: first(rollback, "Revert this PR"),
		owner:    first(owner, "TBD"),
	}
}

func renderCommandLines(cmds []string) string {
	if len(cmds) == 0 {
		return "- `None`"
	}
	lines := make([]string, 0, len(cmds))
	for _, cmd := range cmds {
		lines = append(lines, "- `"+cmd+"`")
	}
	return strings.Join(lines, "\n")
}

func first(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

const prBodyTemplate = `## Summary
%s

## Why
%s

## What Changed
%s

## How To Test
%s

## Risk
%s

## Rollback
%s

## Ownership
- Primary: %s
- Backup: TBD

## Acceptance Criteria
%s

## Observability
- Expected signals: TBD
- Dashboard/alerts: TBD

## Post-Merge Production Commands
<!-- post-merge-commands:start -->
%s
<!-- post-merge-commands:end -->
`
