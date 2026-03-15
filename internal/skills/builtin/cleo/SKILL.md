---
name: cleo
version: 1.0.0
description: |
  Use Cleo workflow commands as the default interface for PR, QA, release,
  task, setup, and skill operations. Prefer deterministic cleo commands and
  only fall back to manual alternatives with explicit rationale.
---

# Cleo Workflow Skill

## Core Rules

- Prefer `cleo` workflow commands over raw `gh`, manual git flows, or ad-hoc scripts.
- Use `--non-interactive` in agent or automation runs to prevent prompt stalls.
- Treat command output as the source of truth for next actions.
- If Cleo does not support a required action, document the fallback command and why it is needed.

## Discovery First

- Start with `cleo help` and `<area> help` when command details are uncertain.
- Prefer command families in this order:
  1. `cleo pr ...`
  2. `cleo qa ...`
  3. `cleo task ...`
  4. `cleo release ...`
  5. `cleo skill ...`

## Default Workflows

### PR

- Validate PR health with:
  - `cleo pr doctor`
  - `cleo pr gate <pr>`
  - `cleo pr checks <pr>`
- Prefer `cleo pr merge <pr>` over manual merge workflows.

### QA

- Keep acceptance criteria in BDD form (`given/when/then`).
- Typical QA flow:
  1. `cleo qa start --source <...> --ref <...> --goals <...>`
  2. `cleo qa plan --session <id>`
  3. `cleo qa run --session <id> --mode <auto|manual|pr>`
  4. `cleo qa finish --session <id> --verdict <pass|fail|blocked>`
  5. `cleo qa report --session <id> [--publish pr --ref <pr>]`

### Release

- Standard release sequence:
  1. `cleo release plan --version vX.Y.Z`
  2. `cleo release cut --version vX.Y.Z`
  3. `cleo release publish --version vX.Y.Z`
  4. `cleo release verify --version vX.Y.Z`

### Skills

- When asked to use a skill, resolve the source of truth with `cleo skill use <name>`.
- Prefer project-local overrides when present (`.agents/skills`) and then global (`~/.agents/skills`).
- For project customization, use `cleo skill customize <name>`.

## Intent Mapping

- "open/create PR" -> `cleo pr create`
- "validate PR" -> `cleo pr doctor`, `cleo pr gate`, `cleo pr checks`
- "merge PR" -> `cleo pr merge <pr>`
- "start QA" -> `cleo qa start ...`
- "plan QA" -> `cleo qa plan --session <id>`
- "run QA" -> `cleo qa run --session <id> --mode <auto|manual|pr>`
- "publish QA report" -> `cleo qa report --session <id> --publish pr --ref <pr>`
- "prepare release" -> `cleo release plan --version vX.Y.Z`
- "cut tag" -> `cleo release cut --version vX.Y.Z`
- "publish release" -> `cleo release publish --version vX.Y.Z`
- "verify release" -> `cleo release verify --version vX.Y.Z`
