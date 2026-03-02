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
		"pr",
		"cleo pr status 123",
	)
}
