# Release Runbook

Use this sequence for GitHub releases:

1. Plan

```bash
cleo release list --limit 10
cleo release latest
cleo release plan --version vX.Y.Z
```

2. Cut tag

```bash
cleo release cut --version vX.Y.Z
```

3. Publish release

```bash
cleo release publish --version vX.Y.Z --final
```

Explicit Go flow (equivalent in Go repos):

```bash
cleo release go publish --version vX.Y.Z --final
```

4. Verify

```bash
cleo release verify --version vX.Y.Z
```

## Notes

- Run from a clean working tree.
- Version must use configured prefix (`release.tag_prefix`, default `v`).
- `publish` supports `--draft` and `--no-notes`.
- In Go repositories (`go.mod`), `publish` auto-builds multi-arch tarballs and `checksums.txt`.
- In Ruby gem repositories (`*.gemspec`), `publish` auto-builds `.gem` + `checksums.txt`.
- `cleo update` pulls from latest GitHub release artifacts and verifies checksums.
- Release notes are generated in a fixed format with required sections:
  - Summary
  - Highlights
  - Breaking Changes
  - Migration Notes
  - Verification
  - GitHub Changes
  - Full Changelog
- Preferred headings inside each `## [vX.Y.Z]` changelog entry:
  - `### Summary`
  - `### Highlights`
  - `### Breaking Changes`
  - `### Migration Notes`
  - `### Verification`
- If sections are missing, cleo warns and uses guidance defaults.
- Agents can provide explicit section data at publish time:
  - `--summary`
  - `--highlights`
  - `--breaking`
  - `--migration`
  - `--verification`
- Cross-project packaging config:
  - `release.binary_name`
  - `release.build_target`
  - `release.changelog_file`
