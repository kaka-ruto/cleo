# AGENTS.md

## Workflow Rules

- Prefer `cleo` workflow commands over raw `gh` or manual steps.
- For PR tasks, start with `cleo pr help` and use `cleo pr ...`.
- For release tasks, start with `cleo release help` and use `cleo release ...`.
- For QA tasks, start with `cleo qa help` and keep ACs in BDD form (`given/when/then`).
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
- "start QA from PR AC block" -> `cleo qa start --source pr --ref <pr> --goals <text>`
- "bootstrap reusable QA kit files" -> `cleo qa init`
- "plan QA from BDD AC" -> `cleo qa plan --session <id>`
- "run QA guidance (default automated-coverage mode)" -> `cleo qa run --session <id> --mode auto`
- "run QA manual checks (if enabled)" -> `cleo qa run --session <id> --mode manual`
- "run QA guidance using PR policy" -> `cleo qa run --session <id> --mode pr`
- "publish QA report to PR (comment + latest summary block)" -> `cleo qa report --session <id> --publish pr --ref <pr>`
- QA evidence root comes from `qa.evidence_dir` (default `.cleo/evidence`) and session outputs go to `.cleo/evidence/qa/session-<id>`.

## Agent Self-Discovery

- Use built-in help before acting:
  - `cleo help`
  - `cleo pr help`
  - `cleo release help`
- Prefer workflow command outputs as source of truth for next steps.
