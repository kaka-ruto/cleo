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
	"## Full Changelog",
}

func buildReleaseNotes(version, generated string) string {
	changelog := fmt.Sprintf("https://github.com/cafaye/cleo/commits/%s", version)
	lines := []string{
		"## Summary",
		"- Initial release for this version.",
		"",
		"## Highlights",
		"- See GitHub Changes for merged PR details.",
		"",
		"## Breaking Changes",
		"- None.",
		"",
		"## Migration Notes",
		"- None.",
		"",
		"## Verification",
		"- Release artifacts uploaded and checksums generated.",
		"",
		"## GitHub Changes",
		strings.TrimSpace(generated),
		"",
		"## Full Changelog",
		changelog,
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
