package release

import "testing"

func TestBuildPlanRequiresVersion(t *testing.T) {
	_, err := BuildPlan(Input{Name: "plan", Args: nil}, Options{TagPrefix: "v"})
	if err == nil {
		t.Fatal("expected version error")
	}
}

func TestBuildPlanVersionPrefix(t *testing.T) {
	_, err := BuildPlan(Input{Name: "cut", Args: []string{"1.2.3"}}, Options{TagPrefix: "v"})
	if err == nil {
		t.Fatal("expected prefix validation error")
	}
}

func TestBuildPlanOK(t *testing.T) {
	p, err := BuildPlan(Input{Name: "verify", Args: []string{"--version", "v1.2.3"}}, Options{TagPrefix: "v"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Version != "v1.2.3" {
		t.Fatalf("unexpected version: %s", p.Version)
	}
}

func TestBuildPlanList(t *testing.T) {
	p, err := BuildPlan(Input{Name: "list"}, Options{TagPrefix: "v"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p.ReadOnly {
		t.Fatal("list should be read-only")
	}
}
