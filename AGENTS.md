# AGENTS.md

Guidance for agents working in this repository.

## Project Focus

`cleo` is a deterministic engineering workflow CLI for GitHub-centered teams.

- Keep scope on deterministic workflow families such as `cleo pr`, `cleo release`, `cleo qa`, and `cleo task`.
- Prefer `cleo` workflows over raw `gh`/manual steps when a workflow command exists.
- Prefer predictable behavior over smart behavior.
- Agent-assist capabilities are allowed when they are bounded by explicit workflow contracts, deterministic plan/run/verify phases, and clear operator control.

## Workflow Preference

- PR work: start with `cleo pr help`, then use `cleo pr ...`.
- Release work: start with `cleo release help`, then use `cleo release ...`.
- QA work: start with `cleo qa help`, then use `cleo qa ...`.
- Do not jump to raw `gh` unless `cleo` has no equivalent.
- Use `--non-interactive` in agent/automation contexts to avoid stalls.
- After every significant improvement, publish a release using `cleo release plan|cut|publish|verify`.

## Local Commands

- Prefer Make targets over direct script or raw tool invocation.
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
- Use `go test ./...` only for targeted fast checks while iterating; before handoff, run `make quality`.

## Intent Mapping

- "create/open PR" -> `cleo pr create`
- "check PR health" -> `cleo pr doctor`, `cleo pr gate`, `cleo pr checks`
- "merge PR" -> `cleo pr merge <pr>`
- "new release" -> `cleo release plan|cut|publish|verify`
- Go release explicitly -> `cleo release go plan|cut|publish|verify`
- "start QA from PR AC block" -> `cleo qa start --source pr --ref <pr> --goals <text>`
- "bootstrap QA kit in a repo" -> `cleo qa init`
- "plan QA from AC" -> `cleo qa plan --session <id>`
- "run QA guidance (default automated-coverage mode)" -> `cleo qa run --session <id> --mode auto`
- "run QA manual checks (if enabled)" -> `cleo qa run --session <id> --mode manual`
- "run QA using PR policy mode" -> `cleo qa run --session <id> --mode pr`
- "log QA findings as tasks" -> `cleo qa log --session <id> --title <text> --details <text>`
- "publish QA report to PR (comment history + latest body summary)" -> `cleo qa report --session <id> --publish pr --ref <pr>`

## QA Contract

- Acceptance Criteria in PRs and `cleo qa start --ac` must be BDD-style YAML.
- Each criterion should declare:
  - `id`, `title`, `severity`, `actors`
  - `surface` (`web|api|mobile|cli`)
  - `environment` (default to `local` unless stated)
  - `given`, `when`, `then` (behavior contract)
  - `evidence_required` (artifacts expected from QA run)
- QA session evidence is written under `qa.evidence_dir/qa/session-<id>` (default `.cleo/evidence/qa/session-<id>`).
- AC defines behavior expectations, not executable command scripts.
- Prefer concrete, observable `then` outcomes and explicit evidence items (for example screenshot, video, API response).
- QA run modes:
  - `auto` (default): verify behaviors are sufficiently covered by automated tests.
  - `manual`: execute manual/exploratory checks and collect artifacts.
- Manual mode is configurable in `cleo.yml` via `qa.manual.enabled`.
- PR QA publishing:
  - CI-triggered QA keeps results in workflow logs and uploaded artifacts.
  - PR comment/body publishing is optional and should not be the default automation path.
- PR QA policy:
  - Policy is read from PR body between `<!-- cleo-qa-policy:start -->` and `<!-- cleo-qa-policy:end -->`.
  - `mode` defaults to `auto` when absent.

## Design Rules

- Keep files small and focused.
- Prefer one responsibility per file.
- Use concise names; single-word names are preferred when clear.
- Keep command UX simple and scriptable.
- Make command failures explicit and actionable.

## Size Limits

- Go files: max ~250 LOC.
- Shell scripts: max ~250 LOC.
- Functions/methods (Go and shell): max ~50 LOC.
- If a file or function approaches these limits, split by responsibility.

## Rule Of Thumb

- If a file needs many comments to explain flow, split it.
- If a function has multiple responsibilities or deep branching, split it.
- If a file crosses soft limits, document a concrete reason.

## Configuration

- `cleo.yml` at repo root is required for command execution.
- Use strict YAML parsing (`KnownFields`).
- Add safe defaults for common teams.
- Avoid framework-specific assumptions (no Rails-only rules).

## Command Principles

- Keep onboarding polished with `cleo setup`:
  - sequential checks
  - explicit confirmations before installs
  - clear progress output
- Read-only commands: `status`, `gate`, `checks`, `watch`.
- Mutating commands: `create`, `merge`, `run`, `rebase`, `retarget`, `batch`.
- Mutations must validate preconditions first where possible.
- Support dry-run mode for risky operations.

## Testing

- Every new behavior needs automated Go tests.
- Test parser/validation logic heavily.
- Test GitHub CLI integration via fakes; avoid network in unit tests.
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
