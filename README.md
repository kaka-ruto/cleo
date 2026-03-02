# cleo

Deterministic CLI for GitHub PR operations.

## One-Command Install

```bash
curl -fsSL https://raw.githubusercontent.com/cafaye/cleo/master/install.sh | bash
```

Non-interactive mode:

```bash
NON_INTERACTIVE=1 curl -fsSL https://raw.githubusercontent.com/cafaye/cleo/master/install.sh | bash
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

Run a guided setup with dependency checks, optional installs, GitHub auth, and `cleo.yml` generation:

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

## Tests

```bash
go test ./...
```
