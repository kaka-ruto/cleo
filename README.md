# cleo

Cleo is a deterministic CLI for software delivery with humans and coding agents.

It turns PRs, QA, tasks, skills, cost estimation, and releases into explicit workflows so teams ship faster with less chaos.

## Why Cleo

Teams adopting coding agents often struggle with:

- inconsistent output quality
- skipped QA and release discipline
- unclear ownership of follow-up work
- ad-hoc shell sequences that are hard to audit

Cleo provides one workflow surface with predictable commands, guardrails, and verifiable outcomes.

## Install

Install once per machine:

```bash
curl -fsSL https://raw.githubusercontent.com/kaka-ruto/cleo/master/install.sh | bash
```

Non-interactive install:

```bash
NON_INTERACTIVE=1 curl -fsSL https://raw.githubusercontent.com/kaka-ruto/cleo/master/install.sh | bash
```

Update:

```bash
cleo update
```

## Quick Start

In each repository (once):

```bash
cd /path/to/repo
cleo setup
```

Then verify core workflows:

```bash
cleo pr doctor
cleo qa init
cleo task list
cleo release latest
cleo skill list
```

## Complete Capability Map

### `setup`

Purpose: bootstrap per-repository workflow assets and prerequisites.

Commands:

- `cleo setup`
- `cleo setup --non-interactive`

Example:

```bash
cleo setup --non-interactive
```

### `update`

Purpose: update your installed `cleo` binary from GitHub releases.

Commands:

- `cleo update`
- `cleo update --non-interactive`

Example:

```bash
cleo update
```

### `pr`

Purpose: run deterministic PR operations with safety checks.

Capabilities:

- PR summary and status visibility
- merge-readiness gating
- GitHub check inspection and watch mode
- local environment doctor checks
- PR creation with structured metadata
- safe merge/rebase/retarget/batch flows
- post-merge command execution support

Commands:

- `cleo pr status <pr>`
- `cleo pr gate <pr>`
- `cleo pr checks <pr>`
- `cleo pr watch <pr|sha>`
- `cleo pr doctor`
- `cleo pr run <pr> [--dry]`
- `cleo pr create [flags]`
- `cleo pr merge <pr> [flags]`
- `cleo pr rebase <pr>`
- `cleo pr retarget <pr> --base <branch>`
- `cleo pr batch [--from <pr>] [flags]`

Examples:

```bash
cleo pr status 123
cleo pr gate 123
cleo pr merge 123 --delete-branch
```

### `qa`

Purpose: manage full QA sessions from scaffold to final report.

Capabilities:

- initialize reusable QA kit files
- scaffold acceptance criteria templates
- start structured QA sessions by source (`branch`, `pr`, `request`)
- run planning and diagnostics for QA sessions
- execute QA (`auto`, `manual`, `pr` modes)
- capture findings with severity
- finish with verdict (`pass`, `fail`, `blocked`)
- publish or print QA reports

Commands:

- `cleo qa init`
- `cleo qa scaffold [--title <text>]`
- `cleo qa start --source <branch|pr|request> --ref <name|id|text> --goals <text> [--ac <yaml>]`
- `cleo qa doctor --session <id>`
- `cleo qa plan --session <id>`
- `cleo qa run --session <id> [--mode <auto|manual|pr>]`
- `cleo qa log --session <id> --title <text> --details <text> [--severity <low|medium|high|critical>]`
- `cleo qa finish --session <id> --verdict <pass|fail|blocked>`
- `cleo qa report --session <id> [--publish <pr>] [--ref <pr>]`

Examples:

```bash
cleo qa start --source pr --ref 123 --goals "Validate checkout happy/failure paths"
cleo qa run --session 20260325-abc123 --mode auto
cleo qa finish --session 20260325-abc123 --verdict pass
```

### `task`

Purpose: track and execute follow-up work identified by QA/review workflows.

Capabilities:

- list by status
- inspect single task details
- claim ownership
- execute with branch strategy
- close completed tasks

Commands:

- `cleo task list [--status <open|in_progress|closed>]`
- `cleo task show --id <task-id>`
- `cleo task claim --id <task-id>`
- `cleo task work --id <task-id> [--new-branch|--in-place]`
- `cleo task close --id <task-id>`

Examples:

```bash
cleo task list --status open
cleo task claim --id TSK-42
cleo task work --id TSK-42 --new-branch
```

### `skill`

Purpose: resolve, manage, install, validate, and customize skills.

Capabilities:

- list resolved skills and source precedence
- print resolved `SKILL.md` for active use
- customize skill locally in `.agents/skills/<name>/SKILL.md`
- validate one or all skills
- install/uninstall skills at project/global scope
- sync builtin skills
- discover remote registries
- browse registry skills
- add/remove custom registries

Commands:

- `cleo skill list`
- `cleo skill use <name>`
- `cleo skill customize <name>`
- `cleo skill check [name]`
- `cleo skill install <name> [--global|--project] [--registry <name>] [--force]`
- `cleo skill uninstall <name> [--global|--project]`
- `cleo skill sync [--global|--project]`
- `cleo skill registry [list]`
- `cleo skill registry skills <registry> [--search <term>]`
- `cleo skill registry add <name> --repo <owner/repo> --path <path> [--ref <ref>] [--description <text>]`
- `cleo skill registry remove <name>`

Built-in registries:

- `openai` (`openai/skills`, `skills/.curated`)
- `superpowers` (`obra/superpowers`, `skills`)
- `superpowers-ruby` (`lucianghinda/superpowers-ruby`, `skills`)

Examples:

```bash
cleo skill registry list
cleo skill registry skills openai --search frontend
cleo skill install frontend-skill --registry openai --global
cleo skill customize cleo
cleo skill uninstall frontend-skill --global
```

### `release`

Purpose: deterministic release lifecycle.

Capabilities:

- list and inspect existing releases
- validate preconditions before tagging
- cut and push version tags
- publish GitHub releases with notes
- verify released artifact visibility
- run explicit Go release flow

Commands:

- `cleo release list [--limit N]`
- `cleo release latest`
- `cleo release plan --version <vX.Y.Z>`
- `cleo release cut --version <vX.Y.Z>`
- `cleo release publish --version <vX.Y.Z> [flags]`
- `cleo release verify --version <vX.Y.Z>`
- `cleo release go <command>`

Examples:

```bash
cleo release plan --version v0.3.0
cleo release cut --version v0.3.0
cleo release publish --version v0.3.0 --final
cleo release verify --version v0.3.0
```

Release order:

1. `plan`
2. `cut`
3. `publish`
4. `verify`

Packaging behavior:

- Go repos (`go.mod`) auto-attach multi-arch tarballs + checksums
- Ruby gem repos (`*.gemspec`) auto-attach `.gem` + checksums

### `cost`

Purpose: estimate engineering cost using codebase metrics and rate sources.

Capabilities:

- estimate from local code metrics
- choose rate source (`cached`, `live`, `manual`)
- choose output format (`markdown`, `plain`, `json`)

Commands:

- `cleo cost estimate`

Examples:

```bash
cleo cost estimate --path . --rates-source cached
cleo cost estimate --rates-source live --country Kenya
cleo cost estimate --rates-source manual --hourly-rate 160 --format json
```

### `version` and `help`

Purpose: discover current version and command usage.

Commands:

- `cleo version`
- `cleo help`
- `cleo help pr`
- `cleo help qa`
- `cleo help task`
- `cleo help skill`
- `cleo help release`
- `cleo help cost`

Example:

```bash
cleo version
cleo help skill
```

## Configuration Model

Cleo does not require `cleo.yml`.

Repo context is inferred from git:

- `remote.origin.url` for host/owner/repo
- origin default branch for base branch
- built-in defaults for operational behavior

## Agent Integration Model

`cleo setup` and `cleo update` auto-install the builtin Cleo skill at:

- `.agents/skills/cleo/SKILL.md`

This keeps agent workflow guidance co-located and versionable with the repo.

## Developer Setup

Build:

```bash
go build ./cmd/cleo
```

Test:

```bash
go test ./...
```

Quality commands:

```bash
make fmt
make lint
make shellcheck
make test
make smoke
make clean
make quality
make ci-status
make install-git-hooks
```

## Uninstall

```bash
curl -fsSL https://raw.githubusercontent.com/kaka-ruto/cleo/master/uninstall.sh | bash
```

Remove Go toolchain and logs too:

```bash
REMOVE_GO=1 REMOVE_LOGS=1 NON_INTERACTIVE=1 curl -fsSL https://raw.githubusercontent.com/kaka-ruto/cleo/master/uninstall.sh | bash
```
