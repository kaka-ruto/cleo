package pr

import (
	"strings"
	"testing"
)

func TestDoctorPassesWhenWorkflowExists(t *testing.T) {
	cfg := testConfig()
	cfg.PR.DeployWatch.Enabled = true
	cfg.PR.DeployWatch.Workflow = "Deploy to Production"

	f := newFakeRunner()
	f.when([]string{"auth", "status"}, "ok")
	f.when([]string{"repo", "view", "kaka-ruto/cleo", "--json", "nameWithOwner"}, `{"nameWithOwner":"kaka-ruto/cleo"}`)
	f.when([]string{"workflow", "list", "--repo", "kaka-ruto/cleo", "--json", "name,path,state"}, `[{"name":"Deploy to Production","path":".github/workflows/deploy.yml","state":"active"}]`)

	svc := NewServiceWithRunner(cfg, f)
	if err := svc.Doctor(); err != nil {
		t.Fatalf("expected doctor success, got %v", err)
	}
}

func TestDoctorFailsWhenWorkflowMissing(t *testing.T) {
	cfg := testConfig()
	cfg.PR.DeployWatch.Enabled = true
	cfg.PR.DeployWatch.Workflow = "Deploy to Production"

	f := newFakeRunner()
	f.when([]string{"auth", "status"}, "ok")
	f.when([]string{"repo", "view", "kaka-ruto/cleo", "--json", "nameWithOwner"}, `{"nameWithOwner":"kaka-ruto/cleo"}`)
	f.when([]string{"workflow", "list", "--repo", "kaka-ruto/cleo", "--json", "name,path,state"}, `[{"name":"ci","path":".github/workflows/ci.yml","state":"active"}]`)

	svc := NewServiceWithRunner(cfg, f)
	err := svc.Doctor()
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected missing workflow error, got %v", err)
	}
}
