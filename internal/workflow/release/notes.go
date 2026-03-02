package release

import (
	"fmt"
	"strings"
)

var requiredNoteSections = []string{
	"## Summary",
	"## Highlights",
	"## Breaking Changes",
	"## Migration Notes",
	"## Verification",
	"## GitHub Changes",
	"## Changelog",
	"## Full Changelog",
}

func buildReleaseNotes(version, generated string) string {
	return buildReleaseNotesWithChangelog(version, generated, ChangelogSections{}, "", "")
}

func buildReleaseNotesWithChangelog(version, generated string, sections ChangelogSections, changelogURL, fullChangelog string) string {
	if strings.TrimSpace(fullChangelog) == "" {
		fullChangelog = fmt.Sprintf("https://github.com/cafaye/cleo/commits/%s", version)
	}
	lines := []string{
		"## Summary",
		strings.TrimSpace(sections.Summary),
		"",
		"## Highlights",
		strings.TrimSpace(sections.Highlights),
		"",
		"## Breaking Changes",
		strings.TrimSpace(sections.BreakingChanges),
		"",
		"## Migration Notes",
		strings.TrimSpace(sections.MigrationNotes),
		"",
		"## Verification",
		strings.TrimSpace(sections.Verification),
		"",
		"## GitHub Changes",
		strings.TrimSpace(generated),
		"",
		"## Changelog",
		changelogURL,
		"",
		"## Full Changelog",
		fullChangelog,
		"",
	}
	return strings.Join(lines, "\n")
}

func validateReleaseNotes(body string) error {
	for _, section := range requiredNoteSections {
		if !strings.Contains(body, section) {
			return fmt.Errorf("release notes missing required section: %s", section)
		}
	}
	return nil
}
