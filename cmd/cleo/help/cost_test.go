package help

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintCost(t *testing.T) {
	var out bytes.Buffer
	PrintCost(&out)
	text := out.String()
	testContainsAll(t, text,
		"usage: cleo cost <command>",
		"estimate",
		"cleo cost estimate",
	)
}

func TestPrintCostCommandKnown(t *testing.T) {
	var out bytes.Buffer
	ok := PrintCostCommand(&out, "estimate")
	if !ok {
		t.Fatal("expected estimate help")
	}
	if !strings.Contains(out.String(), "usage: cleo cost estimate") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestPrintCostCommandUnknown(t *testing.T) {
	var out bytes.Buffer
	ok := PrintCostCommand(&out, "nope")
	if ok {
		t.Fatal("expected unknown command")
	}
}
