# AGENTS.md Template

Copy this into your repository `AGENTS.md` and adjust names/paths.

## Workflow Rules

- Prefer `cleo` workflow commands over raw `gh` or manual steps.
- For PR tasks, start with `cleo pr help` and use `cleo pr ...`.
- For release tasks, start with `cleo release help` and use `cleo release ...`.
- In automation/agent runs, use `--non-interactive` to avoid prompts.
- If `cleo` does not support a required action, document the fallback command clearly.

## Intent Mapping

- "create/open PR" -> `cleo pr create`
- "validate PR" -> `cleo pr doctor`, `cleo pr gate`, `cleo pr checks`
- "merge PR" -> `cleo pr merge <pr>`
- "prepare release" -> `cleo release plan --version vX.Y.Z`
- "cut release tag" -> `cleo release cut --version vX.Y.Z`
- "publish release" -> `cleo release publish --version vX.Y.Z`
- "verify release" -> `cleo release verify --version vX.Y.Z`

## Agent Self-Discovery

- Use built-in help before acting:
  - `cleo help`
  - `cleo pr help`
  - `cleo release help`
- Prefer workflow command outputs as source of truth for next steps.
