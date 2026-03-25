package pr

import (
	"strings"
	"testing"

	"github.com/kaka-ruto/cleo/internal/config"
)

func TestGatePassesForGreenPR(t *testing.T) {
	cfg := testConfig()
	f := newFakeRunner()
	f.when([]string{"pr", "view", "12", "--repo", "kaka-ruto/cleo", "--json", "number,title,url,state,isDraft,mergeable,reviewDecision,baseRefName,headRefName,statusCheckRollup"}, `{"number":12,"title":"T","url":"u","state":"OPEN","isDraft":false,"mergeable":"MERGEABLE","reviewDecision":"APPROVED","baseRefName":"master","headRefName":"feat","statusCheckRollup":[{"name":"ci","workflowName":"CI","status":"COMPLETED","conclusion":"SUCCESS"}]}`)
	svc := NewServiceWithRunner(cfg, f)
	if err := svc.Gate("12"); err != nil {
		t.Fatalf("expected gate pass, got %v", err)
	}
}

func TestGateFailsForDraft(t *testing.T) {
	cfg := testConfig()
	f := newFakeRunner()
	f.when([]string{"pr", "view", "7", "--repo", "kaka-ruto/cleo", "--json", "number,title,url,state,isDraft,mergeable,reviewDecision,baseRefName,headRefName,statusCheckRollup"}, `{"number":7,"title":"T","url":"u","state":"OPEN","isDraft":true,"mergeable":"MERGEABLE","reviewDecision":"APPROVED","baseRefName":"master","headRefName":"feat","statusCheckRollup":[]}`)
	svc := NewServiceWithRunner(cfg, f)
	err := svc.Gate("7")
	if err == nil || !strings.Contains(err.Error(), "draft") {
		t.Fatalf("expected draft error, got %v", err)
	}
}

func TestGateFailsForPendingChecks(t *testing.T) {
	cfg := testConfig()
	f := newFakeRunner()
	f.when([]string{"pr", "view", "9", "--repo", "kaka-ruto/cleo", "--json", "number,title,url,state,isDraft,mergeable,reviewDecision,baseRefName,headRefName,statusCheckRollup"}, `{"number":9,"title":"T","url":"u","state":"OPEN","isDraft":false,"mergeable":"MERGEABLE","reviewDecision":"APPROVED","baseRefName":"master","headRefName":"feat","statusCheckRollup":[{"name":"ci","workflowName":"CI","status":"IN_PROGRESS","conclusion":"","url":"https://example/check"}]}`)
	svc := NewServiceWithRunner(cfg, f)
	err := svc.Gate("9")
	if err == nil || !strings.Contains(err.Error(), "pending checks") {
		t.Fatalf("expected pending checks error, got %v", err)
	}
	if !strings.Contains(err.Error(), "cleo pr watch 9") {
		t.Fatalf("expected watch hint in error, got %v", err)
	}
}

func TestGateFailsWhenNoChecksReported(t *testing.T) {
	cfg := testConfig()
	f := newFakeRunner()
	f.when([]string{"pr", "view", "10", "--repo", "kaka-ruto/cleo", "--json", "number,title,url,state,isDraft,mergeable,reviewDecision,baseRefName,headRefName,statusCheckRollup"}, `{"number":10,"title":"T","url":"u","state":"OPEN","isDraft":false,"mergeable":"MERGEABLE","reviewDecision":"APPROVED","baseRefName":"master","headRefName":"feat","statusCheckRollup":[]}`)
	svc := NewServiceWithRunner(cfg, f)
	err := svc.Gate("10")
	if err == nil || !strings.Contains(err.Error(), "no status checks reported") {
		t.Fatalf("expected no checks error, got %v", err)
	}
}

func testConfig() *config.Config {
	cfg := &config.Config{}
	cfg.Version = 1
	cfg.GitHub.Owner = "kaka-ruto"
	cfg.GitHub.Repo = "cleo"
	cfg.GitHub.BaseBranch = "master"
	cfg.GitHub.MergeMethod = "merge"
	cfg.PR.RequireNonDraft = true
	cfg.PR.RequireMergeable = true
	cfg.PR.BlockRequestedChanges = true
	cfg.PR.Checks.Mode = "required"
	cfg.PR.PostMerge.Markers.Start = "<!-- post-merge-commands:start -->"
	cfg.PR.PostMerge.Markers.End = "<!-- post-merge-commands:end -->"
	cfg.PR.PostMerge.AllowNone = true
	cfg.PR.Stack.AutoDetectNextPR = true
	cfg.PR.Stack.ForceWithLease = true
	cfg.PR.DeployWatch.Workflow = "Deploy to Production"
	cfg.PR.DeployWatch.Branch = "master"
	cfg.PR.DeployWatch.TimeoutSeconds = 5
	cfg.PR.DeployWatch.PollIntervalSeconds = 1
	return cfg
}
