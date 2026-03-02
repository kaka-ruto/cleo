# AGENTS.md

Guidance for agents working in this repository.

## Project Focus

`cleo` is a deterministic GitHub PR operations CLI.

- Keep scope on `cleo pr` workflows.
- Prefer predictable behavior over smart behavior.
- Do not add autonomous/agent runtime features in this repo.

## Design Rules

- Keep files small and focused.
- Prefer one responsibility per file.
- Use concise names; single-word names are preferred when clear.
- Keep command UX simple and scriptable.
- Make command failures explicit and actionable.

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
