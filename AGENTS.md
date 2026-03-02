# AGENTS.md

Guidance for agents working in this repository.

## Project Focus

`cleo` is a deterministic GitHub PR operations CLI.

- Keep scope on `cleo pr` workflows.
- Prefer `cleo` workflows over raw `gh`/manual steps when a workflow command exists.
- Prefer predictable behavior over smart behavior.
- Do not add autonomous/agent runtime features in this repo.

## Workflow Preference

- PR work: start with `cleo pr help`, then use `cleo pr ...`.
- Release work: start with `cleo release help`, then use `cleo release ...`.
- Do not jump to raw `gh` unless `cleo` has no equivalent.
- Use `--non-interactive` in agent/automation contexts to avoid stalls.
- After every significant improvement, publish a release using `cleo release plan|cut|publish|verify`.

## Intent Mapping

- "create/open PR" -> `cleo pr create`
- "check PR health" -> `cleo pr doctor`, `cleo pr gate`, `cleo pr checks`
- "merge PR" -> `cleo pr merge <pr>`
- "new release" -> `cleo release plan|cut|publish|verify`
- Go release explicitly -> `cleo release go plan|cut|publish|verify`

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
