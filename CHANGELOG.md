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
