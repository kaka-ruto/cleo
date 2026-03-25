# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Summary

- Document upcoming changes here before the next release.

### Highlights

- Add highlights for unreleased work.

### Breaking Changes

- None.

### Migration Notes

- None.

### Verification

- Add verification commands/results for unreleased work.

## [v0.2.11]

### Summary

- Prepared Cleo for public launch with a significantly improved README and completed repository owner migration to `kaka-ruto`.

### Highlights

- Rewrote README for launch quality:
  - Complete capability map across `setup`, `update`, `pr`, `qa`, `task`, `skill`, `release`, `cost`, `version/help`
  - Real artifact-oriented output examples (not command echo examples)
  - Practical quick-start and workflow context
  - Updated install/uninstall links and references to `kaka-ruto/cleo`
- Completed owner migration across the codebase:
  - Module path updated to `github.com/kaka-ruto/cleo`
  - Internal imports and GitHub references updated from `cafaye` to `kaka-ruto`
  - Git remotes/workflow references aligned to new owner path
- Removed legacy `docs/YOUR_AGENTS.md` and cleaned related README references to keep documentation skill-first.

### Breaking Changes

- None.

### Migration Notes

- If you have old remotes, update them to:
  - `git@github.com:kaka-ruto/cleo.git`
- If any internal tooling references old module paths, update imports to:
  - `github.com/kaka-ruto/cleo/...`

### Verification

- `go test ./...`
- `go run ./cmd/cleo release latest`
- `go run ./cmd/cleo help`

## [v0.2.10]

### Summary

- Cleo can now discover, install, and uninstall skills directly from remote registries, including OpenAI curated skills and community skill packs.

### Highlights

- New built-in registries available out of the box:
  - `openai` (`openai/skills`, path `skills/.curated`)
  - `superpowers` (`obra/superpowers`, path `skills`)
  - `superpowers-ruby` (`lucianghinda/superpowers-ruby`, path `skills`)
- New discovery commands:
  - `cleo skill registry list`
  - `cleo skill registry skills <registry> [--search <term>]`
- New remote installation command:
  - `cleo skill install <name> --registry <registry> [--global|--project] [--force]`
- New removal command:
  - `cleo skill uninstall <name> [--global|--project]`
- New custom registry management:
  - `cleo skill registry add <name> --repo <owner/repo> --path <path> [--ref <ref>] [--description <text>]`
  - `cleo skill registry remove <name>`
- Custom registries persist in user scope at `~/.agents/skills/registries.yml`.
- Existing customization behavior remains:
  - `cleo skill customize <name>` still creates/uses project-local overrides in `.agents/skills/<name>/SKILL.md`.

### Breaking Changes

- None.

### Migration Notes

- No migration required.
- Recommended upgrade path:
  - Use built-ins immediately: `cleo skill registry list`
  - Add team registries for internal skills:
    - `cleo skill registry add <name> --repo <owner/repo> --path <path> [--ref <ref>]`
  - Install to project scope when sharing with a team:
    - `cleo skill install <skill> --registry <name> --project`

### Verification

- `go test ./...`
- `go run ./cmd/cleo skill registry list`
- `go run ./cmd/cleo skill registry skills openai --search frontend`
- `go run ./cmd/cleo skill install frontend-skill --registry openai --project`
- `go run ./cmd/cleo skill uninstall frontend-skill --project`

## [v0.2.9]

### Summary

- Added a first-class builtin `cleo` skill and switched onboarding to a skill-first model with automatic installation.

### Highlights

- Added bundled `cleo` skill guidance at `internal/skills/builtin/cleo/SKILL.md`.
- Updated setup/update post-migrations to ensure `.agents/skills/cleo/SKILL.md` exists.
- Preserved user customization by leaving existing `.agents/skills/cleo/SKILL.md` files unchanged.
- Added tests for migration behavior and skill command coverage for the new builtin skill.
- Updated README and `docs/YOUR_AGENTS.md` to remove copy-paste-heavy agent setup guidance.

### Breaking Changes

- None.

### Migration Notes

- Run `cleo setup` (or `cleo update`) in each repository to auto-install the builtin `cleo` skill.
- Existing custom project skill overrides at `.agents/skills/cleo/SKILL.md` are preserved.

### Verification

- `go test ./...`

## [v0.2.8]

### Summary

- Removed `cleo.yml` entirely; Cleo now infers repository context from git and runs with built-in policy defaults.

### Highlights

- Deleted repo `cleo.yml` and removed runtime file-based config loading.
- `LoadProject` now resolves config from:
  - git remote (`remote.origin.url`) for host/owner/repo
  - git origin HEAD/default branch for base branch
  - built-in defaults for policy values
- Updated setup/update flows to stop creating or migrating `cleo.yml`.
- Updated docs and agent guidance to reflect no config-file workflow.
- Added/updated tests for:
  - git-only config inference
  - failure mode when git remote is missing
  - setup/update behavior with no `cleo.yml`

### Breaking Changes

- `cleo.yml` is no longer supported or read by Cleo.

### Migration Notes

- Remove existing `cleo.yml`; it is no longer used.
- Ensure repository has a valid origin remote:
  - `git remote add origin git@github.com:<owner>/<repo>.git`
- Ensure origin default branch is discoverable when needed:
  - `git symbolic-ref refs/remotes/origin/HEAD refs/remotes/origin/<branch>`

### Verification

- `go test ./...`
- `go run ./cmd/cleo task list`
- `go run ./cmd/cleo release latest`
- `go run ./cmd/cleo qa init`
- `go run ./cmd/cleo setup --non-interactive`

## [v0.2.7]

### Summary

- Added `cleo skill install` and `cleo skill sync` for easier distribution of skills in `.agents/skills`.

### Highlights

- Added `cleo skill install <name> [--global|--project]` to install one resolved skill into `.agents/skills`.
- Added `cleo skill sync [--global|--project]` to materialize all bundled skills into `.agents/skills`.
- Updated skill help output to include new install/sync options.
- Updated AGENTS guidance and README examples for global/project skill distribution workflows.

### Breaking Changes

- None.

### Migration Notes

- No migration required.
- Teams can populate project skills with:
  - `cleo skill sync --project`
- Users can install specific skills globally with:
  - `cleo skill install ceo --global`

### Verification

- `go test ./...`
- `go run ./cmd/cleo help skill`
- `go run ./cmd/cleo skill sync --project`
- `go run ./cmd/cleo skill install ceo --global`

## [v0.2.6]

### Summary

- Standardized Cleo skill overrides to `.agents/skills` for both project and user scopes.

### Highlights

- Removed `.cleo/skills` lookup and customization paths from the skill resolver.
- Skill resolution now uses:
  - `<project>/.agents/skills/<name>/SKILL.md`
  - `~/.agents/skills/<name>/SKILL.md`
  - built-in bundled skills
- Updated `cleo skill customize` target path to `.agents/skills/<name>/SKILL.md`.
- Updated CLI help, AGENTS guidance, and template docs to use `.agents/skills`.
- Updated resolver/workflow tests to validate `.agents/skills` behavior.

### Breaking Changes

- Existing project overrides under `.cleo/skills/...` are no longer loaded.

### Migration Notes

- Move any existing overrides from `.cleo/skills/<name>/SKILL.md` to `.agents/skills/<name>/SKILL.md`.
- Re-run:
  - `cleo skill list`
  - `cleo skill use <name>`
  to validate expected sources.

### Verification

- `go test ./...`
- `go run ./cmd/cleo help skill`
- `go run ./cmd/cleo skill list`
- `go run ./cmd/cleo skill customize ceo`

## [v0.2.5]

### Summary

- Added first-class `cleo skill` workflows with built-in CEO skill support and project override customization.

### Highlights

- Added top-level `skill` command family:
  - `cleo skill list`
  - `cleo skill use <name>`
  - `cleo skill customize <name>`
  - `cleo skill check [name]`
- Added built-in `ceo` skill bundled with Cleo under embedded skill assets.
- Added deterministic skill resolution order across project/user locations, then built-ins.
- Added skill frontmatter validation checks (`name` and `description` required).
- Updated CLI help, README, and agent templates to support natural-language requests like `use <x> skill`.
- Added automated tests for resolver behavior and skill workflow command execution.

### Breaking Changes

- None.

### Migration Notes

- No migration required.
- Teams can start using skills immediately via:
  - `cleo skill list`
  - `cleo skill use ceo`
- To share custom skill behavior in a repository:
  - `cleo skill customize ceo`
  - commit `.agents/skills/ceo/SKILL.md`

### Verification

- `go test ./...`
- `go run ./cmd/cleo help skill`
- `go run ./cmd/cleo skill list`
- `go run ./cmd/cleo skill use ceo`

## [v0.2.4]

### Summary

- Improved `cleo cost estimate` output formatting with aligned Markdown tables and added multi-format output controls.

### Highlights

- Added explicit output format selection for cost reports:
  - `--format markdown|plain|json` (default: `markdown`)
- Improved Markdown output readability:
  - rendered key sections as tables
  - aligned table columns for consistent terminal display regardless of content length
  - preserved right-aligned numeric columns for scanability
- Updated help and README examples to document output format usage.
- Added/updated cost workflow tests for format behavior.

### Breaking Changes

- None.

### Migration Notes

- No migration required.
- Existing usage remains valid:
  - `cleo cost estimate` still defaults to Markdown output.

### Verification

- `go test ./internal/workflow/cost`
- `go run ./cmd/cleo cost estimate --path . --format markdown`
- `go run ./cmd/cleo cost estimate --path . --format plain`
- `go run ./cmd/cleo cost estimate --path . --format json`

## [v0.2.3]

### Summary

- Added a new project-agnostic `cleo cost estimate` workflow to estimate engineering and team-loaded delivery cost from repository metrics.

### Highlights

- Added top-level `cost` command family and `estimate` subcommand:
  - `cleo cost help`
  - `cleo cost estimate`
- Implemented language-agnostic codebase scanning with file-type classification:
  - code, tests, docs, config
  - mixed-language detection and language mix reporting
- Added rate source modes:
  - `cached` (default deterministic US benchmark baseline)
  - `manual` (`--hourly-rate` override)
  - `live` (country-adjusted rates via live economic data)
- Added country-aware live market estimation:
  - `--country <name|ISO2>`
  - dynamic country resolution without hardcoded country list
  - live PPP data lookup and baseline scaling for low/avg/high rate bands
- Added help text, README docs, and workflow tests for the new command family.

### Breaking Changes

- None.

### Migration Notes

- No migration required.
- Teams can adopt immediately:
  - `cleo cost estimate`
  - `cleo cost estimate --rates-source live --country <country>`

### Verification

- `go test ./internal/workflow/cost ./cmd/cleo/help ./cmd/cleo`
- `go run ./cmd/cleo cost estimate --path .`
- `go run ./cmd/cleo cost estimate --path . --rates-source live --country Kenya`

## [v0.2.2]

### Summary

- Hardened QA bootstrap and setup ergonomics, and made default QA workflow packaging project-agnostic.

### Highlights

- Made packaged `qa.yml` reusable across non-Go repositories by removing hard Go toolchain assumptions and using released `cleo` binary installation.
- Updated default QA workflow behavior to gate QA execution on non-QA check success and run via PR policy/AC markers.
- Changed `cleo setup` to always preserve existing `cleo.yml` without overwrite prompts.
- Extended QA kit bootstrap to provision default actor profile:
  - `.cleo/qa/actors/core.yml` (created if missing, never overwritten).

### Breaking Changes

- None.

### Migration Notes

- Existing repositories can run `cleo qa init` (or `cleo setup`) to provision missing QA assets, including `core.yml`.
- Existing custom actor/config files remain unchanged.

### Verification

- `go test ./internal/qacatalog -v`
- `go test ./internal/setup -v`
- `make quality`

## [v0.2.1]

### Summary

- Improved upgrade safety by adding additive post-update/setup migrations and automatic QA kit bootstrap for existing repos.

### Highlights

- Added safe post-update migrations for existing repositories:
  - Fill only missing `qa` config defaults in `cleo.yml` (no overwrite of existing values).
  - Ensure QA kit assets exist in-repo.
- Wired migrations into both flows:
  - `cleo update` now runs post-update migrations automatically.
  - `cleo setup` now runs the same migrations after config handling.
- Added reusable QA kit bootstrap under `qacatalog` and exposed it through:
  - `cleo qa init`
  - shared setup/update migration path.
- Added setup migration tests covering:
  - missing-key backfill,
  - no-op when keys are already present.

### Breaking Changes

- None.

### Migration Notes

- Existing users should run `cleo update` and then `cleo setup` (or `cleo qa init`) in active repos to ensure QA kit/config defaults are present.
- Existing `cleo.yml` values are preserved; only missing keys are added.

### Verification

- `go test ./internal/setup ./internal/workflow/update -v`
- `cleo setup --non-interactive`
- `make quality`

## [v0.2.0]

### Summary

- Introduced a reusable, CI-gated QA system for PRs with BDD acceptance criteria and policy-driven execution.

### Highlights

- Added end-to-end QA/task workflows for agent-driven execution and follow-up work tracking.
- Standardized Acceptance Criteria to BDD contract fields (`given`, `when`, `then`) with actor/surface/environment metadata.
- Added QA run modes:
  - `auto` (default): validate automated test coverage against BDD criteria.
  - `manual`: optional exploratory/manual execution path.
  - `pr`: resolve mode from PR QA policy block.
- Added PR QA policy markers and parsing (`cleo-qa-policy:start/end`) to support per-PR mode/workflow settings.
- Added QA reporting integration:
  - `cleo qa report` supports PR publishing when explicitly requested.
  - PR checks now surface QA workflow presence relative to AC/policy.
- Added reusable QA bootstrap:
  - `cleo qa init` installs reusable QA kit files.
  - `cleo setup` now bootstraps QA kit assets on first setup.
- Added CI-gated GitHub QA workflow (`.github/workflows/qa.yml`) that runs QA only after CI success.
- Added deterministic local evidence conventions (`qa.evidence_dir`, default `.cleo/evidence`) and per-session evidence directories.
- Added Playwright Go browser automation setup checks and runtime installation support for QA tooling.

### Breaking Changes

- None.

### Migration Notes

- Repositories should include PR AC markers and BDD-formatted criteria for automated QA execution.
- Repositories can run `cleo qa init` to scaffold QA workflow/template assets where missing.
- Default QA automation flow now favors workflow logs/check results over automatic PR comment/body publishing.

### Verification

- `make quality`
- `cleo qa init`
- `cleo qa scaffold`
- `cleo qa start --source pr --ref <pr> --goals "<goal>"`
- `cleo qa run --session <id> --mode auto`
- `cleo qa finish --session <id> --verdict pass`
- `cleo qa report --session <id>`
- `cleo pr checks <pr>`

## [v0.1.4]

### Summary

- Improved PR check reliability and update UX clarity for agent-driven workflows.

### Highlights

- Improved `cleo pr gate` to block on pending or missing checks with explicit `cleo pr watch <pr>` guidance.
- Improved `cleo pr checks` diagnostics with pending/failed summaries and traceability hints.
- Improved `cleo update` logs with current/latest version visibility and clear progress messaging.
- Added repository agent rule to publish a release after each significant improvement.

### Breaking Changes

- None.

### Migration Notes

- None.

### Verification

- `go test ./...`
- `cleo update`

## [v0.1.3]

### Summary

- Made release packaging reusable across projects without hardcoded cleo-only assumptions.

### Highlights

- Added release config keys for cross-project packaging:
  - `release.binary_name`
  - `release.build_target`
  - `release.changelog_file`
- Generalized Go release artifact naming and build target selection.
- Generalized release scripts to accept configurable binary/build targets.

### Breaking Changes

- None.

### Migration Notes

- Existing `cleo.yml` keeps working with defaults.
- Other projects can set `release.binary_name` and `release.build_target`.

### Verification

- `go test ./...`
- `scripts/release/build-assets.sh v0.1.3 /tmp/cleo-release-test`
- `scripts/release/verify-assets.sh v0.1.3 /tmp/cleo-release-test`

## [v0.1.2]

### Summary

- Improved release authoring ergonomics so agents can always produce structured notes.

### Highlights

- Added release inspection commands: `cleo release list` and `cleo release latest`.
- Added warning-first fallback behavior when changelog note sections are missing.
- Added publish-time release note override flags (`--summary`, `--highlights`, `--breaking`, `--migration`, `--verification`).

### Breaking Changes

- None.

### Migration Notes

- None.

### Verification

- `go test ./...`
- `cleo release help publish`

## [v0.1.1]

### Summary

- Established end-to-end GitHub release automation and release-based updater behavior.

### Highlights

- Added artifact build + checksum verification + publish automation.
- Added explicit Go release path (`cleo release go ...`) and runtime split.
- Improved agent workflow guidance and copy-ready `docs/YOUR_AGENTS.md`.

### Breaking Changes

- None.

### Migration Notes

- None.

### Verification

- Release workflow run succeeded for `v0.1.1`.

## [v0.1.0]

### Summary

- Initial public release of `cleo` workflow-driven PR automation.

### Highlights

- Added PR workflow commands (`status`, `gate`, `checks`, `doctor`, `run`, `create`, `merge`, `rebase`, `retarget`, `batch`).
- Added setup wizard and one-command install/uninstall scripts.
- Introduced modular `plan -> run -> verify` workflow architecture.

### Breaking Changes

- None.

### Migration Notes

- None.

### Verification

- Release artifacts and checksums were published.
