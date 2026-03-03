package qa

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/cafaye/cleo/internal/config"
	"github.com/cafaye/cleo/internal/qaaction"
	"github.com/cafaye/cleo/internal/qacatalog"
	"github.com/cafaye/cleo/internal/qacontract"
	"github.com/cafaye/cleo/internal/qatools"
	"github.com/cafaye/cleo/internal/taskstore"
	"gopkg.in/yaml.v3"
)

type Adapter struct {
	store    *taskstore.Store
	repoKey  string
	cfg      *config.Config
	registry qaaction.Registry
	loader   qacatalog.Loader
	now      func() time.Time
}

const (
	qaSummaryStart = "<!-- cleo-qa-results:start -->"
	qaSummaryEnd   = "<!-- cleo-qa-results:end -->"
	qaPolicyStart  = "<!-- cleo-qa-policy:start -->"
	qaPolicyEnd    = "<!-- cleo-qa-policy:end -->"
)

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

func (a *Adapter) Init() error {
	return qacatalog.EnsureQAKit(".")
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
	if err := os.MkdirAll(a.sessionEvidenceDir(s.ID), 0o755); err != nil {
		return 0, fmt.Errorf("create QA evidence dir: %w", err)
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

func (a *Adapter) Report(sessionID int64, publish string, ref string) (string, error) {
	session, err := a.store.Session(context.Background(), sessionID)
	if err != nil {
		return "", err
	}
	doc, err := qacontract.LoadString(session.ACText)
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
	lines = append(lines, fmt.Sprintf("evidence_dir=%s", a.sessionEvidenceDir(session.ID)))
	lines = append(lines, "")
	if len(tasks) == 0 {
		lines = append(lines, "No tasks logged.")
	} else {
		lines = append(lines, "Tasks:")
		for _, task := range tasks {
			lines = append(lines, fmt.Sprintf("- #%d [%s] %s (status=%s occurrences=%d)", task.ID, task.Severity, task.Title, task.Status, task.Occurrences))
		}
	}

	text := strings.Join(lines, "\n")
	if strings.TrimSpace(publish) == "" {
		return text, nil
	}
	if strings.TrimSpace(publish) != "pr" {
		return "", fmt.Errorf("--publish must be pr")
	}
	prRef := strings.TrimSpace(ref)
	if prRef == "" {
		if strings.TrimSpace(session.Source) != "pr" || strings.TrimSpace(session.Ref) == "" {
			return "", fmt.Errorf("publish=pr requires --ref when session source is not pr")
		}
		prRef = strings.TrimSpace(session.Ref)
	}
	md := a.renderPRReportMarkdown(session, doc, tasks)
	if err := a.publishReportToPR(prRef, md); err != nil {
		return "", err
	}
	return text + "\n\nPublished QA report to PR " + prRef, nil
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

func (a *Adapter) Run(sessionID int64, mode string) (string, error) {
	session, err := a.store.Session(context.Background(), sessionID)
	if err != nil {
		return "", err
	}
	doc, _, err := a.loadSessionContract(sessionID)
	if err != nil {
		return "", err
	}
	resolvedMode := strings.TrimSpace(mode)
	if resolvedMode == "" {
		resolvedMode = "auto"
	}
	if resolvedMode == "pr" {
		if strings.TrimSpace(session.Source) != "pr" || strings.TrimSpace(session.Ref) == "" {
			return "", fmt.Errorf("--mode pr requires a PR-backed session")
		}
		resolvedMode, err = a.qaModeFromPRPolicy(strings.TrimSpace(session.Ref))
		if err != nil {
			return "", err
		}
	}
	if resolvedMode != "auto" && resolvedMode != "manual" {
		return "", fmt.Errorf("--mode must be auto|manual|pr")
	}
	if resolvedMode == "manual" && !a.cfg.QAManualEnabled() {
		return "", fmt.Errorf("manual QA mode is disabled in cleo.yml (qa.manual.enabled=false)")
	}
	var lines []string
	lines = append(lines, fmt.Sprintf("QA AC guidance for session %d", sessionID))
	lines = append(lines, fmt.Sprintf("name=%s", emptyDefault(doc.Name, "unnamed")))
	lines = append(lines, fmt.Sprintf("mode=%s", resolvedMode))
	lines = append(lines, fmt.Sprintf("criteria=%d", len(doc.Criteria)))
	lines = append(lines, fmt.Sprintf("evidence_dir=%s", a.sessionEvidenceDir(sessionID)))
	lines = append(lines, "")
	for _, criterion := range doc.Criteria {
		lines = append(lines, fmt.Sprintf("criterion %s: %s", criterion.ID, criterion.Title))
		lines = append(lines, fmt.Sprintf("  actors: %s", strings.Join(criterion.Actors, ",")))
		lines = append(lines, fmt.Sprintf("  surface: %s", criterion.Surface))
		lines = append(lines, fmt.Sprintf("  environment: %s", emptyDefault(criterion.Environment, "local")))
		lines = append(lines, fmt.Sprintf("  given: %s", criterion.Given))
		lines = append(lines, fmt.Sprintf("  when: %s", criterion.When))
		lines = append(lines, "  then:")
		for idx, expected := range criterion.Then {
			lines = append(lines, fmt.Sprintf("    %d. %s", idx+1, expected))
		}
		if len(criterion.Evidence) > 0 {
			lines = append(lines, fmt.Sprintf("  evidence_required=%s", strings.Join(criterion.Evidence, ",")))
		}
		lines = append(lines, "")
	}
	if resolvedMode == "auto" {
		lines = append(lines, "Mode auto:")
		lines = append(lines, "1. Ensure each criterion has automated test coverage for given/when/then behavior.")
		lines = append(lines, "2. Run test suites that cover these behaviors (at minimum `make test`; use targeted suites as needed).")
		lines = append(lines, "3. For missing or weak coverage, log QA findings with exact criterion IDs and coverage gaps.")
	} else {
		lines = append(lines, "Mode manual:")
		lines = append(lines, "1. Execute exploratory/manual checks for each criterion across declared actors and surfaces.")
		lines = append(lines, "2. Capture required evidence artifacts (screenshots, videos, API responses, logs).")
		lines = append(lines, "3. Log defects and behavior mismatches with criterion IDs and artifact references.")
	}
	lines = append(lines, "Then log findings with `cleo qa log`, and close session with `cleo qa finish`.")
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

func (a *Adapter) Scaffold(title string) (string, error) {
	name := strings.TrimSpace(title)
	if name == "" {
		name = "Acceptance Criteria"
	}
	block := fmt.Sprintf(`version: 1
name: %s
criteria:
  - id: c1
    title: Replace with criterion title
    severity: medium
    actors: [core]
    surface: web
    environment: local
    given: Replace with setup state and actor context
    when: Replace with user/system action under test
    then:
      - Replace with observable expected outcome
    evidence_required:
      - replace_with_evidence_artifact`, name)
	return block, nil
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
	body, err := a.prBody(ref)
	if err != nil {
		return "", err
	}
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

func (a *Adapter) qaModeFromPRPolicy(ref string) (string, error) {
	body, err := a.prBody(ref)
	if err != nil {
		return "", err
	}
	block := extractMarkerBlock(body, qaPolicyStart, qaPolicyEnd)
	if strings.TrimSpace(block) == "" {
		return "auto", nil
	}
	var policy struct {
		Mode string `yaml:"mode"`
	}
	if err := yaml.Unmarshal([]byte(block), &policy); err != nil {
		return "", fmt.Errorf("parse PR QA policy: %w", err)
	}
	mode := strings.TrimSpace(policy.Mode)
	if mode == "" {
		mode = "auto"
	}
	if mode != "auto" && mode != "manual" {
		return "", fmt.Errorf("PR QA policy mode must be auto|manual")
	}
	return mode, nil
}

func (a *Adapter) prBody(ref string) (string, error) {
	repo := strings.TrimSpace(a.cfg.GitHub.Owner) + "/" + strings.TrimSpace(a.cfg.GitHub.Repo)
	cmd := exec.Command("gh", "pr", "view", strings.TrimSpace(ref), "--repo", repo, "--json", "body", "--jq", ".body")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("read PR body for AC: %s", strings.TrimSpace(string(out)))
	}
	return string(out), nil
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

func (a *Adapter) renderPRReportMarkdown(session taskstore.Session, doc qacontract.Document, tasks []taskstore.Task) string {
	var lines []string
	lines = append(lines, fmt.Sprintf("### QA Session %d", session.ID))
	lines = append(lines, fmt.Sprintf("- Verdict: `%s`", emptyDefault(session.Verdict, "pending")))
	lines = append(lines, fmt.Sprintf("- Source: `%s` `%s`", session.Source, session.Ref))
	lines = append(lines, fmt.Sprintf("- Goals: %s", session.Goals))
	lines = append(lines, fmt.Sprintf("- Evidence Dir: `%s`", a.sessionEvidenceDir(session.ID)))
	lines = append(lines, "")
	lines = append(lines, "#### BDD Results")
	checked := len(tasks) == 0 && strings.TrimSpace(session.Verdict) == "pass"
	box := " "
	if checked {
		box = "x"
	}
	for _, c := range doc.Criteria {
		lines = append(lines, fmt.Sprintf("- [%s] %s: %s", box, c.ID, c.Title))
		lines = append(lines, fmt.Sprintf("  - Given: %s", c.Given))
		lines = append(lines, fmt.Sprintf("  - When: %s", c.When))
		lines = append(lines, "  - Then:")
		for _, outcome := range c.Then {
			lines = append(lines, fmt.Sprintf("    - [%s] %s", box, outcome))
		}
	}
	lines = append(lines, "")
	if len(tasks) == 0 {
		lines = append(lines, "#### Findings")
		lines = append(lines, "- None")
	} else {
		lines = append(lines, "#### Findings")
		for _, task := range tasks {
			lines = append(lines, fmt.Sprintf("- #%d [%s] %s", task.ID, task.Severity, task.Title))
		}
	}
	return strings.Join(lines, "\n")
}

func (a *Adapter) sessionEvidenceDir(sessionID int64) string {
	return filepath.Join(a.cfg.QAEvidenceDir(), "qa", fmt.Sprintf("session-%d", sessionID))
}

func (a *Adapter) publishReportToPR(prRef, report string) error {
	repo := strings.TrimSpace(a.cfg.GitHub.Owner) + "/" + strings.TrimSpace(a.cfg.GitHub.Repo)
	if _, err := exec.Command("gh", "pr", "comment", prRef, "--repo", repo, "--body", report).CombinedOutput(); err != nil {
		return fmt.Errorf("post QA PR comment: %w", err)
	}
	view := exec.Command("gh", "pr", "view", prRef, "--repo", repo, "--json", "body", "--jq", ".body")
	out, err := view.CombinedOutput()
	if err != nil {
		return fmt.Errorf("read PR body for QA summary update: %s", strings.TrimSpace(string(out)))
	}
	body := string(out)
	summaryBlock := qaSummaryStart + "\n" + report + "\n" + qaSummaryEnd
	next := upsertMarkerBlock(body, qaSummaryStart, qaSummaryEnd, summaryBlock)
	tmp, err := os.CreateTemp("", "cleo-qa-pr-body-*.md")
	if err != nil {
		return fmt.Errorf("create temp PR body file: %w", err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.WriteString(next); err != nil {
		return fmt.Errorf("write temp PR body file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp PR body file: %w", err)
	}
	edit := exec.Command("gh", "pr", "edit", prRef, "--repo", repo, "--body-file", tmp.Name())
	if out, err := edit.CombinedOutput(); err != nil {
		return fmt.Errorf("update PR body QA summary: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func upsertMarkerBlock(body, start, end, block string) string {
	si := strings.Index(body, start)
	ei := strings.Index(body, end)
	if si >= 0 && ei > si {
		ei += len(end)
		return strings.TrimSpace(body[:si]) + "\n\n" + block + "\n\n" + strings.TrimSpace(body[ei:])
	}
	trimmed := strings.TrimSpace(body)
	if trimmed == "" {
		return block + "\n"
	}
	return trimmed + "\n\n## QA Results\n" + block + "\n"
}

func extractMarkerBlock(body, start, end string) string {
	si := strings.Index(body, start)
	ei := strings.Index(body, end)
	if si < 0 || ei <= si {
		return ""
	}
	return strings.TrimSpace(body[si+len(start) : ei])
}
