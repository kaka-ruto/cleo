package pr

import (
	"strings"
	"testing"
)

func TestRenderIncludesMarkersAndSections(t *testing.T) {
	body := Render("Summary", "Why", "- A", "- B", "Low", "Revert", "alice", "version: 1\nname: t\ncriteria: []", []string{"bin/kamal logs"})
	checks := []string{
		"## Summary",
		"## Why",
		"## What Changed",
		"## How To Test",
		"## Acceptance Criteria",
		"<!-- cleo-ac:start -->",
		"<!-- cleo-ac:end -->",
		"## Post-Merge Production Commands",
		"<!-- post-merge-commands:start -->",
		"<!-- post-merge-commands:end -->",
	}
	for _, c := range checks {
		if !strings.Contains(body, c) {
			t.Fatalf("expected body to contain %q", c)
		}
	}
}

func TestRenderDefaults(t *testing.T) {
	body := Render("", "", "", "", "", "", "", "", nil)
	if !strings.Contains(body, "TBD") {
		t.Fatal("expected defaults to include TBD")
	}
	if !strings.Contains(body, "- `None`") {
		t.Fatal("expected default None command")
	}
	if !strings.Contains(body, "actors: [core]") {
		t.Fatal("expected default AC scaffold")
	}
	if !strings.Contains(body, "given: Replace with setup state and actor context") {
		t.Fatal("expected default BDD AC scaffold")
	}
}
