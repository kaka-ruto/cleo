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
	if err := ensureDefaultCoreActor(root); err != nil {
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

func ensureDefaultCoreActor(root string) error {
	path := filepath.Join(root, ".cleo", "qa", "actors", "core.yml")
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create actors dir: %w", err)
	}
	if err := os.WriteFile(path, []byte(defaultCoreActor), 0o644); err != nil {
		return fmt.Errorf("write default core actor: %w", err)
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
    types:
      - opened
      - synchronize
      - reopened
      - ready_for_review
  workflow_dispatch:

permissions:
  contents: read
  pull-requests: read

concurrency:
  group: ${{ github.workflow }}-${{ github.event.pull_request.number || github.ref }}
  cancel-in-progress: true

jobs:
  qa:
    if: |
      github.event_name == 'pull_request' ||
      github.event_name == 'workflow_dispatch'
    runs-on: ubuntu-latest
    env:
      GH_TOKEN: ${{ github.token }}
      PR_NUMBER: ${{ github.event.pull_request.number }}
      SHA: ${{ github.event.pull_request.head.sha || github.sha }}
    steps:
      - uses: actions/checkout@v4
      - run: |
          sudo apt-get update
          sudo apt-get install -y gh jq curl tar gzip nodejs
      - shell: bash
        run: |
          set -euo pipefail
          api="https://api.github.com/repos/kaka-ruto/cleo/releases/latest"
          version="$(curl -fsSL "$api" | jq -r '.tag_name')"
          asset="cleo_${version}_linux_amd64.tar.gz"
          url="$(curl -fsSL "$api" | jq -r --arg name "$asset" '.assets[] | select(.name==$name) | .browser_download_url')"
          curl -fsSL "$url" -o /tmp/cleo.tgz
          tar -xzf /tmp/cleo.tgz -C /tmp
          sudo install -m 0755 /tmp/cleo /usr/local/bin/cleo
      - if: github.event_name == 'pull_request'
        shell: bash
        run: |
          set -euo pipefail
          for _ in $(seq 1 60); do
            data="$(gh pr view "$PR_NUMBER" --json statusCheckRollup --jq '.statusCheckRollup')"
            pending="$(jq '[.[] | select((.workflowName // "" | ascii_downcase) != "qa") | select((.status // "") != "COMPLETED")] | length' <<< "$data")"
            failed="$(jq '[.[] | select((.workflowName // "" | ascii_downcase) != "qa") | select((.status // "") == "COMPLETED") | select((.conclusion // "SUCCESS") != "SUCCESS" and (.conclusion // "") != "NEUTRAL")] | length' <<< "$data")"
            if [[ "$failed" -gt 0 ]]; then exit 1; fi
            if [[ "$pending" -eq 0 ]]; then break; fi
            sleep 10
          done
      - if: github.event_name == 'pull_request'
        shell: bash
        run: |
          set -euo pipefail
          body="$(gh pr view "$PR_NUMBER" --json body --jq .body)"
          if ! grep -q "<!-- cleo-ac:start -->" <<< "$body" || ! grep -q "<!-- cleo-ac:end -->" <<< "$body"; then
            echo "No AC markers; skipping."
            exit 0
          fi
          policy="$(awk '/<!-- cleo-qa-policy:start -->/{flag=1;next}/<!-- cleo-qa-policy:end -->/{flag=0}flag' <<< "$body")"
          mode="$(awk -F: '/^mode:/{print $2; exit}' <<< "$policy" | xargs || true)"
          if [[ -z "$mode" ]]; then mode="auto"; fi
          short_sha="$(printf "%s" "$SHA" | cut -c1-7)"
          goals="CI-gated QA run for PR #${PR_NUMBER} (${short_sha})"
          sid="$(cleo qa start --source pr --ref "$PR_NUMBER" --goals "$goals" | awk '{print $4}')"
          cleo qa plan --session "$sid"
          cleo qa doctor --session "$sid"
          set +e
          cleo qa run --session "$sid" --mode "$mode"
          run_status=$?
          set -e
          verdict="pass"
          if [[ "$run_status" -ne 0 ]]; then verdict="fail"; fi
          cleo qa finish --session "$sid" --verdict "$verdict"
          cleo qa report --session "$sid"
          exit "$run_status"
`

const defaultCoreActor = `name: core
description: Default actor profile for core QA flows.
surfaces:
  - web
  - api
auth_profile: none
`
