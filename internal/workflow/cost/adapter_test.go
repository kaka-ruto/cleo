package cost

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEstimateProducesReport(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main_test.go"), []byte("package main\n\nfunc TestX() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	report, err := Estimate([]string{"--path", dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(report, "# Cleo Cost Estimate") {
		t.Fatalf("unexpected report:\n%s", report)
	}
	if !strings.Contains(report, "Rates Source: cached") {
		t.Fatalf("expected default rates source in report:\n%s", report)
	}
}

func TestEstimateManualRate(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.py"), []byte("print('x')\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	report, err := Estimate([]string{"--path", dir, "--rates-source", "manual", "--hourly-rate", "150"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(report, "Rates Source: manual") {
		t.Fatalf("unexpected report:\n%s", report)
	}
}
