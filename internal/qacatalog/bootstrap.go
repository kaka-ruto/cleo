package qacatalog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	acStart        = "<!-- cleo-ac:start -->"
	acEnd          = "<!-- cleo-ac:end -->"
	qaResultsStart = "<!-- cleo-qa-results:start -->"
	qaResultsEnd   = "<!-- cleo-qa-results:end -->"
	qaPolicyStart  = "<!-- cleo-qa-policy:start -->"
	qaPolicyEnd    = "<!-- cleo-qa-policy:end -->"
)

func EnsureQAKit(root string) error {
	if err := ensureQAWorkflow(root); err != nil {
		return err
	}
	if err := ensurePRTemplate(root); err != nil {
		return err
	}
	return nil
}

func ensureQAWorkflow(root string) error {
	path := filepath.Join(root, ".github", "workflows", "qa.yml")
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create workflow dir: %w", err)
	}
	if err := os.WriteFile(path, []byte(defaultWorkflow), 0o644); err != nil {
		return fmt.Errorf("write qa workflow: %w", err)
	}
	return nil
}

func ensurePRTemplate(root string) error {
	path := filepath.Join(root, ".github", "pull_request_template.md")
	body := ""
	if raw, err := os.ReadFile(path); err == nil {
		body = string(raw)
	}
	if strings.TrimSpace(body) == "" {
		body = defaultPRTemplate
	} else {
		if !strings.Contains(body, acStart) || !strings.Contains(body, acEnd) {
			body += "\n\n## Acceptance Criteria\n" + defaultACBlock
		}
		if !strings.Contains(body, qaResultsStart) || !strings.Contains(body, qaResultsEnd) {
			body += "\n\n## QA Results\n" + defaultQAResultsBlock
		}
		if !strings.Contains(body, qaPolicyStart) || !strings.Contains(body, qaPolicyEnd) {
			body += "\n\n## QA Policy\n" + defaultQAPolicyBlock
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create template dir: %w", err)
	}
	if err := os.WriteFile(path, []byte(strings.TrimSpace(body)+"\n"), 0o644); err != nil {
		return fmt.Errorf("write PR template: %w", err)
	}
	return nil
}

const defaultACBlock = `<!-- cleo-ac:start -->
version: 1
name: Acceptance Criteria
criteria:
  - id: c1
    title: Replace with criterion title
    severity: medium
    actors: [core]
    surface: web
    environment: local
    given: Replace with setup state and actor context
    when: Replace with user/system action under test
    then:
      - Replace with observable expected outcome
    evidence_required:
      - replace_with_evidence_artifact
<!-- cleo-ac:end -->`

const defaultQAResultsBlock = `<!-- cleo-qa-results:start -->
_No QA report published yet._
<!-- cleo-qa-results:end -->`

const defaultQAPolicyBlock = `<!-- cleo-qa-policy:start -->
mode: auto
workflow: qa
<!-- cleo-qa-policy:end -->`

const defaultPRTemplate = `## Summary

## Why

## What Changed

## How To Test

## Risk

## Rollback

## Ownership
- Primary:
- Backup:

## Acceptance Criteria
` + defaultACBlock + `

## QA Results
` + defaultQAResultsBlock + `

## QA Policy
` + defaultQAPolicyBlock + `
`

const defaultWorkflow = `name: qa

on:
  pull_request:
    types: [opened, synchronize, reopened, ready_for_review]
  workflow_run:
    workflows: [ci]
    types: [completed]
  workflow_dispatch:

permissions:
  contents: read
  pull-requests: read

jobs:
  qa:
    if: |
      github.event_name == 'pull_request' ||
      github.event_name == 'workflow_dispatch' ||
      (github.event_name == 'workflow_run' && github.event.workflow_run.conclusion == 'success')
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true
      - run: |
          sudo apt-get update
          sudo apt-get install -y gh shellcheck nodejs
      - run: echo "Use repository's committed qa.yml for full QA execution logic."
`
