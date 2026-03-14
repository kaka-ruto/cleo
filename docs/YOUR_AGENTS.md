# AGENTS.md

## Core Rules

- Prefer `cleo` workflow commands over raw `gh` or manual steps.
- Use `--non-interactive` in automation/agent runs to avoid stalls.
- If `cleo` does not support a required action, document the fallback command clearly.
- Treat command output as source of truth for next steps.

## PR Workflows

- Start PR work with `cleo pr help`, then use `cleo pr ...` commands.
- Prefer built-in PR guardrails before merge:
  - `cleo pr doctor`
  - `cleo pr gate <pr>`
  - `cleo pr checks <pr>`
- Use `cleo pr merge <pr>` instead of manual merge flows when possible.

## Release Workflows

- Start release work with `cleo release help`, then use `cleo release ...`.
- Standard release sequence:
  1. `cleo release plan --version vX.Y.Z`
  2. `cleo release cut --version vX.Y.Z`
  3. `cleo release publish --version vX.Y.Z`
  4. `cleo release verify --version vX.Y.Z`

## QA Workflows

- Start QA work with `cleo qa help`, then use `cleo qa ...`.
- Keep acceptance criteria in BDD form (`given/when/then`).
- Typical session flow:
  1. `cleo qa start --source <...> --ref <...> --goals <...>`
  2. `cleo qa plan --session <id>`
  3. `cleo qa run --session <id> --mode <auto|manual|pr>`
  4. `cleo qa finish --session <id> --verdict <pass|fail|blocked>`
  5. `cleo qa report --session <id> [--publish pr --ref <pr>]`
- QA evidence root comes from `qa.evidence_dir` (default `.cleo/evidence`) and session outputs go to `.cleo/evidence/qa/session-<id>`.

## Skills Workflows

- Start skill work with `cleo skill help`, then use `cleo skill ...`.
- Skills are instruction overlays for the current response, not separate tasks.
- When the user says "use <x> skill", "apply <x> skill", or equivalent phrasing:
  1. Run `cleo skill use <x>`.
  2. Apply the resolved `SKILL.md` immediately to the current response.
  3. Follow its structure and constraints.
- If a skill cannot be resolved, report that briefly and continue with best-effort default behavior.
- Do not invent skill content; `cleo skill use <x>` output is the source of truth.
- If users want project-specific behavior, suggest `cleo skill customize <x>` and then use the customized skill.

## Setup and Update

- Use `cleo setup` to bootstrap repository workflow assets.
- Use `cleo update` to update Cleo from releases.
- After `cleo update`, Cleo applies safe additive setup migrations for existing repos (missing config defaults + QA kit assets) without overwriting existing config values.

## Agent Self-Discovery

- Use built-in help before acting:
  - `cleo help`
  - `cleo pr help`
  - `cleo release help`
  - `cleo qa help`
  - `cleo skill help`

## Intent Mapping

- "create/open PR" -> `cleo pr create`
- "validate PR" -> `cleo pr doctor`, `cleo pr gate`, `cleo pr checks`
- "merge PR" -> `cleo pr merge <pr>`
- "prepare release" -> `cleo release plan --version vX.Y.Z`
- "cut release tag" -> `cleo release cut --version vX.Y.Z`
- "publish release" -> `cleo release publish --version vX.Y.Z`
- "verify release" -> `cleo release verify --version vX.Y.Z`
- "start QA from PR AC block" -> `cleo qa start --source pr --ref <pr> --goals <text>`
- "bootstrap reusable QA kit files" -> `cleo qa init`
- "plan QA from BDD AC" -> `cleo qa plan --session <id>`
- "run QA guidance (default automated-coverage mode)" -> `cleo qa run --session <id> --mode auto`
- "run QA manual checks (if enabled)" -> `cleo qa run --session <id> --mode manual`
- "run QA guidance using PR policy" -> `cleo qa run --session <id> --mode pr`
- "publish QA report to PR (comment + latest summary block)" -> `cleo qa report --session <id> --publish pr --ref <pr>`
- "use ceo skill" / "use <x> skill" -> `cleo skill use <x>` then apply that SKILL.md as the instruction overlay for the current response
