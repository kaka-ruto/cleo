# cleo

Cleo is the master CLI for humans and agents.

## Positioning

### What It Is

- A deterministic CLI that turns common engineering work into explicit, repeatable workflows.
- A bridge between human intent and agent execution.
- A quality-and-delivery harness that standardizes how PRs, QA, tasks, and releases are done.

### What It Does Today

- `cleo pr`: structured PR creation, checks, gating, and merge safety.
- `cleo qa`: BDD-style acceptance criteria, policy-driven QA runs, and CI-integrated QA execution.
- `cleo task`: captures and tracks follow-up work from QA/results.
- `cleo release`: plan, cut, publish, verify with release discipline.
- `cleo skill`: resolve, use, validate, and customize agent skills.
- `cleo setup` / `cleo update`: safe, additive bootstrap and maintenance.

### Why It Exists

- Teams are adopting coding agents quickly, but execution quality is inconsistent.
- Cleo gives agents rails: clear contracts, predictable outputs, and auditable steps.
- It reduces handholding while improving reliability.

### What It Intends To Be

- The default control plane for agentic software delivery.
- A reusable standard any repo can adopt, regardless of stack.
- A system where humans define intent and policy, and agents execute with rigor.

## How Setup Works

`cleo` has two setup layers:

1. Global setup (one-time per machine): install `cleo` binary and dependencies.
2. Repository setup (one-time per repo): initialize/maintain that repo's workflow config and QA kit assets.

After global install, you can use the same `cleo` command in any repository. You only run `cleo setup` again when entering a new repo for the first time.

## One-Command Install

```bash
curl -fsSL https://raw.githubusercontent.com/kaka-ruto/cleo/master/install.sh | bash
```

Non-interactive mode:

```bash
NON_INTERACTIVE=1 curl -fsSL https://raw.githubusercontent.com/kaka-ruto/cleo/master/install.sh | bash
```

## Quick Start Across Repositories

Install once:

```bash
curl -fsSL https://raw.githubusercontent.com/kaka-ruto/cleo/master/install.sh | bash
```

In each repository (once):

```bash
cd /path/to/repo
cleo setup
```

Use normally after that:

```bash
cleo pr status <pr>
cleo qa init
cleo qa scaffold
cleo pr doctor
cleo task list
cleo release latest
```

## Help

```bash
cleo help
cleo pr help
cleo pr help merge
cleo skill help
```

## Agent Setup

`cleo setup` and `cleo update` now auto-install the builtin `cleo` skill at:

- `.agents/skills/cleo/SKILL.md`

This removes the need for large copy/paste `AGENTS.md` templates.

## Update cleo

Update from the latest GitHub release:

```bash
cleo update
```

Non-interactive update:

```bash
cleo update --non-interactive
```

`cleo update` checks latest GitHub release first; if already current, it exits quickly with an up-to-date message.

## One-Command Uninstall

```bash
curl -fsSL https://raw.githubusercontent.com/kaka-ruto/cleo/master/uninstall.sh | bash
```

Options:

```bash
NON_INTERACTIVE=1 SCAN_ROOTS="$HOME/Code,$HOME/work" curl -fsSL https://raw.githubusercontent.com/kaka-ruto/cleo/master/uninstall.sh | bash
```

- `NON_INTERACTIVE=1`: auto-confirms prompts.
- `SCAN_ROOTS`: comma-separated directories scanned for legacy `cleo.yml` cleanup.

## Setup Wizard

Run a guided per-repo setup with dependency checks, optional installs, GitHub auth, and additive repository bootstrap:

```bash
cleo setup
cleo setup --non-interactive
```

`cleo setup` does not create or read `cleo.yml`; it infers repo context from git and applies safe additive QA kit assets.

## Build

```bash
go build ./cmd/cleo
```

## Developer Commands

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

Logs default to `~/.cleo/logs/<repo>`. Override per run with:

```bash
LOG_DIR=logs make quality
```

## PR Commands

```bash
cleo pr status <pr>
cleo pr gate <pr>
cleo pr checks <pr>
cleo pr doctor
cleo pr watch <pr|sha>
cleo pr run <pr> [--dry]
cleo pr create [--title ...] [--summary ...] [--why ...] [--what ...] [--test ...] [--risk ...] [--rollback ...] [--owner ...] [--cmd ...] [--draft]
cleo pr merge <pr> [--no-watch] [--no-run] [--no-rebase] [--delete-branch]
cleo pr rebase <pr>
cleo pr retarget <pr> --base <branch>
cleo pr batch [--from <pr>] [--no-watch] [--no-run] [--no-rebase]
```

## QA Commands

```bash
cleo qa init
cleo qa scaffold [--title <text>]
cleo qa start --source <branch|pr|request> --ref <name|id|text> --goals <text> [--ac <yaml>]
cleo qa plan --session <id>
cleo qa doctor --session <id>
cleo qa run --session <id> [--mode <auto|manual|pr>]
cleo qa log --session <id> --title <text> --details <text> [--severity <low|medium|high|critical>]
cleo qa finish --session <id> --verdict <pass|fail|blocked>
cleo qa report --session <id> [--publish <pr>] [--ref <pr>]
```

## Task Commands

```bash
cleo task list
cleo task show --task <id>
cleo task claim --task <id>
cleo task close --task <id>
cleo task work --task <id>
```

## Skill Commands

```bash
cleo skill list
cleo skill use cleo
cleo skill use ceo
cleo skill customize cleo
cleo skill customize ceo
cleo skill install cleo --global
cleo skill install ceo --global
cleo skill sync --project
cleo skill check
cleo skill check cleo
cleo skill check ceo
```

## Cost Commands

```bash
cleo cost help
cleo cost estimate
cleo cost estimate --path . --rates-source cached
cleo cost estimate --rates-source live --country Germany
cleo cost estimate --format markdown
cleo cost estimate --format plain
cleo cost estimate --format json
cleo cost estimate --rates-source manual --hourly-rate 160
```

## Release Commands

```bash
cleo release help
cleo release go help
cleo release list --limit 10
cleo release latest
cleo release plan --version v0.1.0
cleo release cut --version v0.1.0
cleo release publish --version v0.1.0 [--draft|--final] [--no-notes] [--summary "..."] [--highlights "..."] [--breaking "..."] [--migration "..."] [--verification "..."]
cleo release verify --version v0.1.0
cleo release go publish --version v0.1.0 [--draft|--final] [--no-notes]
```

Release workflow follows the same deterministic pattern:

1. `plan` validates preconditions.
2. `cut` creates and pushes the tag.
3. `publish` creates the GitHub release.
4. `verify` confirms release visibility.

`publish` generates release notes in an enforced sectioned format using the matching `CHANGELOG.md` version entry and includes GitHub-generated change notes inside that template. If changelog sections are missing, cleo warns and uses guidance defaults; agents can provide explicit section text via publish flags.

For Go repositories (`go.mod` present), `publish` automatically:

- builds `linux/darwin` binaries for `amd64/arm64`
- packages tarballs
- generates `checksums.txt`
- uploads artifacts to the GitHub release

For Ruby gem repositories (`*.gemspec` present), `publish` automatically:

- builds the gem into `dist/release/<version>/`
- writes `checksums.txt`
- uploads the `.gem` and checksum file to the GitHub release

Release packaging currently uses built-in defaults:

- `binary_name`: `cleo`
- `build_target`: `./cmd/cleo`
- `changelog_file`: `CHANGELOG.md`

## Tests

```bash
go test ./...
```

## Uninstall

```bash
curl -fsSL https://raw.githubusercontent.com/kaka-ruto/cleo/master/uninstall.sh | bash
```

Remove Go toolchain and logs too:

```bash
REMOVE_GO=1 REMOVE_LOGS=1 NON_INTERACTIVE=1 curl -fsSL https://raw.githubusercontent.com/kaka-ruto/cleo/master/uninstall.sh | bash
```
