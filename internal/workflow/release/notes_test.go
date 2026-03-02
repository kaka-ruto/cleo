package release

import (
	"strings"
	"testing"
)

func TestBuildReleaseNotesContainsSections(t *testing.T) {
	body := buildReleaseNotes("v1.2.3", "## What's Changed\n- test")
	for _, section := range requiredNoteSections {
		if !strings.Contains(body, section) {
			t.Fatalf("missing section %q", section)
		}
	}
}

func TestValidateReleaseNotes(t *testing.T) {
	body := buildReleaseNotes("v1.2.3", "x")
	if err := validateReleaseNotes(body); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if err := validateReleaseNotes("## Summary\nmissing"); err == nil {
		t.Fatal("expected validation error")
	}
}
