# cleo

Deterministic CLI for GitHub PR operations.

## How Setup Works

`cleo` has two setup layers:

1. Global setup (one-time per machine): install `cleo` binary and dependencies.
2. Repository setup (one-time per repo): create/update that repo's `cleo.yml`.

After global install, you can use the same `cleo` command in any repository. You only run `cleo setup` again when entering a new repo for the first time.

## One-Command Install

```bash
curl -fsSL https://raw.githubusercontent.com/cafaye/cleo/master/install.sh | bash
```

Non-interactive mode:

```bash
NON_INTERACTIVE=1 curl -fsSL https://raw.githubusercontent.com/cafaye/cleo/master/install.sh | bash
```

## Quick Start Across Repositories

Install once:

```bash
curl -fsSL https://raw.githubusercontent.com/cafaye/cleo/master/install.sh | bash
```

In each repository (once):

```bash
cd /path/to/repo
cleo setup
```

Use normally after that:

```bash
cleo pr status <pr>
cleo pr doctor
```

## Help

```bash
cleo help
cleo pr help
cleo pr help merge
```

## Agent Setup

Add workflow preferences to your project `AGENTS.md` so agents default to `cleo` commands.

Use this copy-ready template:

- [YOUR_AGENTS.md](/Users/kaka/Code/Cafaye/cleo/docs/YOUR_AGENTS.md)

## Update cleo

Update from the latest GitHub release:

```bash
cleo update
```

Non-interactive update:

```bash
cleo update --non-interactive
```

## One-Command Uninstall

```bash
curl -fsSL https://raw.githubusercontent.com/cafaye/cleo/master/uninstall.sh | bash
```

Options:

```bash
NON_INTERACTIVE=1 SCAN_ROOTS="$HOME/Code,$HOME/work" curl -fsSL https://raw.githubusercontent.com/cafaye/cleo/master/uninstall.sh | bash
```

- `NON_INTERACTIVE=1`: auto-confirms prompts.
- `SCAN_ROOTS`: comma-separated directories scanned for `cleo.yml` cleanup.

## Setup Wizard

Run a guided per-repo setup with dependency checks, optional installs, GitHub auth, and `cleo.yml` generation:

```bash
cleo setup
cleo setup --non-interactive
```

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

For reuse in other projects, configure release packaging in `cleo.yml`:

- `release.binary_name` (artifact/binary name)
- `release.build_target` (Go build target path, for example `./cmd/mycli`)
- `release.changelog_file` (notes source path)

## Tests

```bash
go test ./...
```

## Uninstall

```bash
curl -fsSL https://raw.githubusercontent.com/cafaye/cleo/master/uninstall.sh | bash
```

Remove Go toolchain and logs too:

```bash
REMOVE_GO=1 REMOVE_LOGS=1 NON_INTERACTIVE=1 curl -fsSL https://raw.githubusercontent.com/cafaye/cleo/master/uninstall.sh | bash
```
