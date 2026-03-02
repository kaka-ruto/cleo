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

## Setup Wizard

Run a guided setup with dependency checks, optional installs, GitHub auth, and `cleo.yml` generation:

```bash
cleo setup
```

## Build

```bash
go build ./cmd/cleo
```

## Commands

```bash
cleo pr status <pr>
cleo pr gate <pr>
cleo pr checks <pr>
cleo pr watch <pr|sha>
cleo pr run <pr> [--dry]
cleo pr create [--title ...] [--summary ...] [--why ...] [--what ...] [--test ...] [--risk ...] [--rollback ...] [--owner ...] [--cmd ...] [--draft]
cleo pr merge <pr> [--no-watch] [--no-run] [--no-rebase] [--delete-branch]
cleo pr rebase <pr>
cleo pr retarget <pr> --base <branch>
cleo pr batch [--from <pr>] [--no-watch] [--no-run] [--no-rebase]
```

`cleo pr` commands resolve the target repository from `cleo.yml` (`github.owner` + `github.repo`); there is no user-facing `--repo` flag.

## Tests

```bash
go test ./...
```
