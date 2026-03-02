package release

import (
	"strings"
	"testing"
)

func TestBuildReleaseNotesContainsSections(t *testing.T) {
	sections := ChangelogSections{
		Summary:         "- summary",
		Highlights:      "- highlights",
		BreakingChanges: "- none",
		MigrationNotes:  "- none",
		Verification:    "- test",
	}
	body := buildReleaseNotesWithChangelog("v1.2.3", "## What's Changed\n- test", sections, "https://example/changelog", "https://example/commits")
	for _, section := range requiredNoteSections {
		if !strings.Contains(body, section) {
			t.Fatalf("missing section %q", section)
		}
	}
}

func TestValidateReleaseNotes(t *testing.T) {
	sections := ChangelogSections{
		Summary:         "- summary",
		Highlights:      "- highlights",
		BreakingChanges: "- none",
		MigrationNotes:  "- none",
		Verification:    "- test",
	}
	body := buildReleaseNotesWithChangelog("v1.2.3", "x", sections, "https://example/changelog", "https://example/commits")
	if err := validateReleaseNotes(body); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if err := validateReleaseNotes("## Summary\nmissing"); err == nil {
		t.Fatal("expected validation error")
	}
}
