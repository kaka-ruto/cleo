# AGENTS.md

Guidance for agents working in this repository.

## Project Focus

`cleo` is a deterministic engineering workflow CLI for GitHub-centered teams.

- Keep scope on deterministic workflow families such as `cleo pr`, `cleo release`, `cleo qa`, `cleo task`, and `cleo skill`.
- Prefer `cleo` workflows over raw `gh`/manual steps when a workflow command exists.
- Prefer predictable behavior over smart behavior.
- Agent-assist capabilities are allowed when they are bounded by explicit workflow contracts, deterministic plan/run/verify phases, and clear operator control.

## Core Rules

- Use `--non-interactive` in agent/automation contexts to avoid stalls.
- If `cleo` does not support a required action, document the fallback command clearly.
- Prefer command output as source of truth for next steps.

## PR Workflows

- Start PR work with `cleo pr help`, then use `cleo pr ...`.
- Preferred validation path before merge:
  - `cleo pr doctor`
  - `cleo pr gate <pr>`
  - `cleo pr checks <pr>`
- Prefer `cleo pr merge <pr>` over manual merge flows where possible.

## Release Workflows

- Start release work with `cleo release help`, then use `cleo release ...`.
- Standard flow:
  1. `cleo release plan --version vX.Y.Z`
  2. `cleo release cut --version vX.Y.Z`
  3. `cleo release publish --version vX.Y.Z`
  4. `cleo release verify --version vX.Y.Z`
- After significant improvements, publish a release following the standard flow.
- Every release must include a proper changelog entry before running `cleo release publish`.

## QA Workflows

- Start QA work with `cleo qa help`, then use `cleo qa ...`.
- Acceptance Criteria in PRs and `cleo qa start --ac` must be BDD-style YAML.
- Each criterion should declare:
  - `id`, `title`, `severity`, `actors`
  - `surface` (`web|api|mobile|cli`)
  - `environment` (default `local` unless stated)
  - `given`, `when`, `then`
  - `evidence_required`
- QA run modes:
  - `auto` (default): automated coverage verification
  - `manual`: exploratory/manual checks
  - `pr`: resolved from PR policy block
- QA evidence root comes from `qa.evidence_dir` (default `.cleo/evidence`) and session outputs go to `.cleo/evidence/qa/session-<id>`.

## Skills Workflows

- Start skill work with `cleo skill help`, then use `cleo skill ...`.
- Skills are instruction overlays for the current response, not separate tasks.
- When the user says "use <x> skill" or equivalent phrasing:
  1. Run `cleo skill use <x>`.
  2. Apply the resolved `SKILL.md` immediately to the current response.
  3. Follow its structure and constraints.
- If a skill cannot be resolved, report that briefly and continue best-effort.
- Do not invent skill content; `cleo skill use <x>` output is the source of truth.
- For project-specific behavior, use `cleo skill customize <x>` and commit `.cleo/skills/<x>/SKILL.md` when team sharing is desired.

## Setup and Update

- `cleo update` should only apply safe additive migrations (no destructive config overwrites).
- `cleo setup` should preserve existing `cleo.yml` and apply safe additive defaults/assets.

## Local Commands

- Prefer Make targets over raw tool invocation.
- Canonical commands:
  - `make fmt`
  - `make lint`
  - `make shellcheck`
  - `make test`
  - `make smoke`
  - `make quality`
  - `make ci-status`
  - `make clean`
  - `make install-git-hooks`
- Default pre-PR validation: run `make quality`.

## Intent Mapping

- "create/open PR" -> `cleo pr create`
- "check PR health" -> `cleo pr doctor`, `cleo pr gate`, `cleo pr checks`
- "merge PR" -> `cleo pr merge <pr>`
- "new release" -> `cleo release plan|cut|publish|verify`
- "start QA from PR AC block" -> `cleo qa start --source pr --ref <pr> --goals <text>`
- "bootstrap QA kit in a repo" -> `cleo qa init`
- "plan QA from AC" -> `cleo qa plan --session <id>`
- "run QA guidance (default automated-coverage mode)" -> `cleo qa run --session <id> --mode auto`
- "run QA manual checks (if enabled)" -> `cleo qa run --session <id> --mode manual`
- "run QA using PR policy mode" -> `cleo qa run --session <id> --mode pr`
- "log QA findings as tasks" -> `cleo qa log --session <id> --title <text> --details <text>`
- "publish QA report to PR" -> `cleo qa report --session <id> --publish pr --ref <pr>`
- "use ceo skill" / "use <x> skill" -> `cleo skill use <x>` then apply that `SKILL.md` as instruction overlay for the current response

## Agent Self-Discovery

- Use built-in help before acting:
  - `cleo help`
  - `cleo pr help`
  - `cleo release help`
  - `cleo qa help`
  - `cleo task help`
  - `cleo cost help`
  - `cleo skill help`

## Design Rules

- Keep files small and focused.
- Prefer one responsibility per file.
- Use concise names when clear.
- Keep command UX simple and scriptable.
- Make command failures explicit and actionable.

## Configuration

- `cleo.yml` at repo root is required for command execution.
- Use strict YAML parsing (`KnownFields`).
- Add safe defaults for common teams.
- Avoid framework-specific assumptions.

## Testing

- Every new behavior needs automated Go tests.
- Test parser/validation logic heavily.
- Test integrations via fakes when feasible.
- Keep tests deterministic and fast.

## Safety

- Post-merge command execution must use marker block parsing.
- Enforce allowlist/denylist policy from config.
- Do not execute arbitrary commands without explicit config allowance.

## Quality Bar

Before finishing changes:

1. `go test ./...` passes.
2. No giant files or mixed concerns.
3. Error messages are clear and operational.
