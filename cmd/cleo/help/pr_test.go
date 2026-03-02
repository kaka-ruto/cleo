package help

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintPR(t *testing.T) {
	var out bytes.Buffer
	PrintPR(&out)
	testContainsAll(t, out.String(),
		"usage: cleo pr <command>",
		"status <pr>",
		"merge <pr>",
		"help [command]",
	)
}

func TestPrintPRCommandKnown(t *testing.T) {
	var out bytes.Buffer
	ok := PrintPRCommand(&out, "merge")
	if !ok {
		t.Fatal("expected known command to return true")
	}
	if got := strings.TrimSpace(out.String()); got != "usage: cleo pr merge <pr> [--no-watch] [--no-run] [--no-rebase] [--delete-branch]" {
		t.Fatalf("unexpected merge usage: %q", got)
	}
}

func TestPrintPRCommandUnknown(t *testing.T) {
	var out bytes.Buffer
	ok := PrintPRCommand(&out, "nope")
	if ok {
		t.Fatal("expected unknown command to return false")
	}
	if out.Len() != 0 {
		t.Fatalf("expected no output for unknown command, got: %q", out.String())
	}
}
