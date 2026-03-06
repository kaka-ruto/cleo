package help

import (
	"bytes"
	"strings"
	"testing"
)

func testContainsAll(t *testing.T, text string, parts ...string) {
	t.Helper()
	for _, part := range parts {
		if !strings.Contains(text, part) {
			t.Fatalf("expected output to contain %q, got:\n%s", part, text)
		}
	}
}

func TestPrintRoot(t *testing.T) {
	var out bytes.Buffer
	PrintRoot(&out)
	testContainsAll(t, out.String(),
		"usage: cleo <command>",
		"setup",
		"update",
		"pr",
		"release",
		"cost",
		"cleo pr status 123",
	)
}

func TestPrintCommandUpdate(t *testing.T) {
	var out bytes.Buffer
	ok := PrintCommand(&out, "update")
	if !ok {
		t.Fatal("expected update command help")
	}
	if !strings.Contains(out.String(), "usage: cleo update") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestPrintCommandCost(t *testing.T) {
	var out bytes.Buffer
	ok := PrintCommand(&out, "cost")
	if !ok {
		t.Fatal("expected cost command help")
	}
	if !strings.Contains(out.String(), "usage: cleo cost") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}
