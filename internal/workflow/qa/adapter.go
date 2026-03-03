package qa

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/cafaye/cleo/internal/config"
	"github.com/cafaye/cleo/internal/qaaction"
	"github.com/cafaye/cleo/internal/qacatalog"
	"github.com/cafaye/cleo/internal/qacontract"
	"github.com/cafaye/cleo/internal/qatools"
	"github.com/cafaye/cleo/internal/taskstore"
)

type Adapter struct {
	store    *taskstore.Store
	repoKey  string
	cfg      *config.Config
	registry qaaction.Registry
	loader   qacatalog.Loader
	now      func() time.Time
}

func NewAdapter(store *taskstore.Store, repoKey string, cfg *config.Config) *Adapter {
	return &Adapter{
		store:    store,
		repoKey:  strings.TrimSpace(repoKey),
		cfg:      cfg,
		registry: qaaction.NewRegistry(),
		loader:   qacatalog.Loader{ActorsDir: cfg.QA.ActorsDir},
		now:      time.Now,
	}
}

func (a *Adapter) Start(source string, ref string, goals string, ac string) (int64, error) {
	acText, err := a.resolveAC(source, ref, ac)
	if err != nil {
		return 0, err
	}
	s, err := a.store.StartSession(context.Background(), source, ref, goals, acText, a.now())
	if err != nil {
		return 0, err
	}
	return s.ID, nil
}

func (a *Adapter) LogIssue(sessionID int64, title string, details string, severity string) (int64, bool, error) {
	if severity == "" {
		severity = "medium"
	}
	if _, err := a.store.Session(context.Background(), sessionID); err != nil {
		return 0, false, err
	}
	key := dedupeKey(a.repoKey, title, details)
	task, created, err := a.store.UpsertOpenTask(context.Background(), taskstore.Task{
		SessionID: sessionID,
		RepoKey:   a.repoKey,
		Title:     title,
		Details:   details,
		Severity:  severity,
		DedupeKey: key,
	}, a.now())
	if err != nil {
		return 0, false, err
	}
	return task.ID, created, nil
}

func (a *Adapter) Finish(sessionID int64, verdict string) error {
	if err := validateVerdict(verdict); err != nil {
		return err
	}
	return a.store.FinishSession(context.Background(), sessionID, verdict, a.now())
}

func (a *Adapter) Report(sessionID int64) (string, error) {
	session, err := a.store.Session(context.Background(), sessionID)
	if err != nil {
		return "", err
	}
	tasks, err := a.store.TasksBySession(context.Background(), sessionID)
	if err != nil {
		return "", err
	}
	var lines []string
	lines = append(lines, fmt.Sprintf("QA session %d", session.ID))
	lines = append(lines, fmt.Sprintf("source=%s ref=%s verdict=%s", session.Source, session.Ref, emptyDefault(session.Verdict, "pending")))
	lines = append(lines, fmt.Sprintf("goals=%s", session.Goals))
	lines = append(lines, "")
	if len(tasks) == 0 {
		lines = append(lines, "No tasks logged.")
		return strings.Join(lines, "\n"), nil
	}
	lines = append(lines, "Tasks:")
	for _, task := range tasks {
		lines = append(lines, fmt.Sprintf("- #%d [%s] %s (status=%s occurrences=%d)", task.ID, task.Severity, task.Title, task.Status, task.Occurrences))
	}
	return strings.Join(lines, "\n"), nil
}

func (a *Adapter) Plan(sessionID int64) (string, error) {
	doc, actors, err := a.loadSessionContract(sessionID)
	if err != nil {
		return "", err
	}
	var lines []string
	lines = append(lines, fmt.Sprintf("QA AC plan for session %d", sessionID))
	lines = append(lines, fmt.Sprintf("name=%s", emptyDefault(doc.Name, "unnamed")))
	lines = append(lines, fmt.Sprintf("criteria=%d", len(doc.Criteria)))
	lines = append(lines, fmt.Sprintf("actors=%s", strings.Join(actors, ",")))
	lines = append(lines, fmt.Sprintf("required_tools=%s", strings.Join(a.registry.ToolSummary(doc), ",")))
	for _, criterion := range doc.Criteria {
		lines = append(lines, fmt.Sprintf("- [%s] %s (%s) actors=%s", criterion.ID, criterion.Title, emptyDefault(criterion.Severity, "medium"), strings.Join(criterion.Actors, ",")))
	}
	return strings.Join(lines, "\n"), nil
}

func (a *Adapter) Run(sessionID int64) (string, error) {
	doc, _, err := a.loadSessionContract(sessionID)
	if err != nil {
		return "", err
	}
	var lines []string
	lines = append(lines, fmt.Sprintf("QA AC guidance for session %d", sessionID))
	lines = append(lines, fmt.Sprintf("name=%s", emptyDefault(doc.Name, "unnamed")))
	lines = append(lines, fmt.Sprintf("criteria=%d", len(doc.Criteria)))
	lines = append(lines, "")
	for _, criterion := range doc.Criteria {
		lines = append(lines, fmt.Sprintf("criterion %s: %s", criterion.ID, criterion.Title))
		lines = append(lines, fmt.Sprintf("  actors: %s", strings.Join(criterion.Actors, ",")))
		lines = append(lines, fmt.Sprintf("  goal: %s", criterion.Acceptance.Goal))
		lines = append(lines, fmt.Sprintf("  expected_result: %s", criterion.Acceptance.ExpectedResult))
		lines = append(lines, fmt.Sprintf("  surface: %s", criterion.Execution.Surface))
		if len(criterion.Execution.Preconditions) > 0 {
			lines = append(lines, "  preconditions:")
			keys := make([]string, 0, len(criterion.Execution.Preconditions))
			for key := range criterion.Execution.Preconditions {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			for _, key := range keys {
				lines = append(lines, fmt.Sprintf("    - %s=%s", key, criterion.Execution.Preconditions[key]))
			}
		}
		lines = append(lines, "  steps:")
		for idx, step := range criterion.Execution.Steps {
			lines = append(lines, fmt.Sprintf("    %d. action=%s params=%v", idx+1, step.Action, step.Params))
		}
		if len(criterion.Evidence) > 0 {
			lines = append(lines, fmt.Sprintf("  evidence_required=%s", strings.Join(criterion.Evidence, ",")))
		}
		lines = append(lines, "")
	}
	lines = append(lines, "Run the criteria using configured tools, then log failures with `cleo qa log` and finish with `cleo qa finish`.")
	return strings.Join(lines, "\n"), nil
}

func (a *Adapter) Doctor(sessionID int64) (string, error) {
	doc, _, err := a.loadSessionContract(sessionID)
	if err != nil {
		return "", err
	}
	checks := qatools.Doctor(a.registry.ToolSummary(doc))
	var lines []string
	lines = append(lines, "QA tool doctor")
	for _, check := range checks {
		lines = append(lines, fmt.Sprintf("- %s: %s (%s)", check.Name, check.Status, check.Detail))
	}
	return strings.Join(lines, "\n"), nil
}

func (a *Adapter) loadSessionContract(sessionID int64) (qacontract.Document, []string, error) {
	session, err := a.store.Session(context.Background(), sessionID)
	if err != nil {
		return qacontract.Document{}, nil, err
	}
	doc, err := qacontract.LoadString(session.ACText)
	if err != nil {
		return qacontract.Document{}, nil, err
	}
	if err := a.registry.Validate(doc); err != nil {
		return qacontract.Document{}, nil, err
	}
	actors, err := a.validateActors(doc)
	if err != nil {
		return qacontract.Document{}, nil, err
	}
	return doc, actors, nil
}

func (a *Adapter) validateActors(doc qacontract.Document) ([]string, error) {
	set := map[string]struct{}{}
	for _, criterion := range doc.Criteria {
		for _, actor := range criterion.Actors {
			name := strings.TrimSpace(actor)
			if name == "" {
				continue
			}
			if _, err := a.loader.LoadActor(name); err != nil {
				return nil, err
			}
			set[name] = struct{}{}
		}
	}
	out := make([]string, 0, len(set))
	for name := range set {
		out = append(out, name)
	}
	sort.Strings(out)
	return out, nil
}

func (a *Adapter) resolveAC(source string, ref string, inline string) (string, error) {
	if strings.TrimSpace(inline) != "" {
		if _, err := qacontract.LoadString(inline); err != nil {
			return "", err
		}
		return inline, nil
	}
	if strings.TrimSpace(source) == "pr" {
		return a.acFromPR(ref)
	}
	return "", fmt.Errorf("acceptance criteria are required: provide --ac or use source=pr with AC block in PR body")
}

func (a *Adapter) acFromPR(ref string) (string, error) {
	repo := strings.TrimSpace(a.cfg.GitHub.Owner) + "/" + strings.TrimSpace(a.cfg.GitHub.Repo)
	cmd := exec.Command("gh", "pr", "view", strings.TrimSpace(ref), "--repo", repo, "--json", "body", "--jq", ".body")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("read PR body for AC: %s", strings.TrimSpace(string(out)))
	}
	body := string(out)
	start := "<!-- cleo-ac:start -->"
	end := "<!-- cleo-ac:end -->"
	si := strings.Index(body, start)
	ei := strings.Index(body, end)
	if si < 0 || ei <= si {
		return "", fmt.Errorf("PR does not contain AC block markers %q and %q", start, end)
	}
	ac := strings.TrimSpace(body[si+len(start) : ei])
	if ac == "" {
		return "", fmt.Errorf("PR AC block is empty")
	}
	if _, err := qacontract.LoadString(ac); err != nil {
		return "", err
	}
	return ac, nil
}

func dedupeKey(repoKey string, title string, details string) string {
	normalized := strings.ToLower(strings.TrimSpace(repoKey) + "|" + strings.TrimSpace(title) + "|" + strings.TrimSpace(details))
	sum := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(sum[:])
}

func validateVerdict(v string) error {
	switch strings.TrimSpace(v) {
	case "pass", "fail", "blocked":
		return nil
	default:
		return fmt.Errorf("--verdict must be pass|fail|blocked")
	}
}

func emptyDefault(v string, d string) string {
	if strings.TrimSpace(v) == "" {
		return d
	}
	return v
}
