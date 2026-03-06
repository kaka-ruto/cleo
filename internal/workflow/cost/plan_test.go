package cost

import "testing"

func TestBuildPlanEstimateDefaults(t *testing.T) {
	p, err := BuildPlan(Input{Name: "estimate"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "estimate" {
		t.Fatalf("unexpected plan: %#v", p)
	}
}

func TestBuildPlanRejectsBadRatesSource(t *testing.T) {
	_, err := BuildPlan(Input{Name: "estimate", Args: []string{"--rates-source", "live"}})
	if err != nil {
		t.Fatalf("expected live to be valid, got: %v", err)
	}
	_, err = BuildPlan(Input{Name: "estimate", Args: []string{"--rates-source", "nope"}})
	if err == nil {
		t.Fatal("expected invalid source error")
	}
}

func TestBuildPlanManualRequiresRate(t *testing.T) {
	_, err := BuildPlan(Input{Name: "estimate", Args: []string{"--rates-source", "manual"}})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
