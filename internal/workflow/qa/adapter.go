package qa

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/cafaye/cleo/internal/config"
	"github.com/cafaye/cleo/internal/qaassets"
	"github.com/cafaye/cleo/internal/taskstore"
)

type Adapter struct {
	store    *taskstore.Store
	repoKey  string
	cfg      *config.Config
	loader   qaassets.Loader
	now      func() time.Time
	runCmdFn func(command string, vars map[string]string) error
}

func NewAdapter(store *taskstore.Store, repoKey string, cfg *config.Config) *Adapter {
	return &Adapter{
		store:   store,
		repoKey: strings.TrimSpace(repoKey),
		cfg:     cfg,
		loader: qaassets.Loader{
			ProfilesDir:     cfg.QA.ProfilesDir,
			RunbooksDir:     cfg.QA.RunbooksDir,
			EnvironmentsDir: cfg.QA.EnvironmentsDir,
		},
		now:      time.Now,
		runCmdFn: runCommand,
	}
}

func (a *Adapter) Start(source string, ref string, goals string) (int64, error) {
	s, err := a.store.StartSession(context.Background(), source, ref, goals, a.now())
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

func (a *Adapter) Plan(sessionID int64, env string, profiles []string) (string, error) {
	if _, err := a.store.Session(context.Background(), sessionID); err != nil {
		return "", err
	}
	resolved, envCfg, runbooks, err := a.resolveAssets(env, profiles)
	if err != nil {
		return "", err
	}
	var lines []string
	lines = append(lines, fmt.Sprintf("QA plan for session %d", sessionID))
	lines = append(lines, fmt.Sprintf("environment=%s", envCfg.Name))
	lines = append(lines, fmt.Sprintf("profiles=%s", strings.Join(resolved, ",")))
	lines = append(lines, "runbooks:")
	for _, runbook := range runbooks {
		lines = append(lines, fmt.Sprintf("- %s (%d checks)", runbook.Name, len(runbook.Checks)))
	}
	return strings.Join(lines, "\n"), nil
}

func (a *Adapter) Run(sessionID int64, env string, profiles []string) (string, error) {
	if _, err := a.store.Session(context.Background(), sessionID); err != nil {
		return "", err
	}
	resolved, envCfg, runbooks, err := a.resolveAssets(env, profiles)
	if err != nil {
		return "", err
	}
	checked := 0
	passed := 0
	failed := 0
	created := 0
	var taskRefs []string
	for _, runbook := range runbooks {
		for _, check := range runbook.Checks {
			checked++
			if strings.TrimSpace(check.Command) == "" {
				passed++
				continue
			}
			err := a.runCmdFn(check.Command, envCfg.Vars)
			if err == nil {
				passed++
				continue
			}
			failed++
			title := firstNonEmpty(check.FailureTitle, fmt.Sprintf("QA check failed: %s", firstNonEmpty(check.Title, check.ID)))
			details := firstNonEmpty(check.FailureDetails, fmt.Sprintf("runbook=%s check=%s command=%q error=%s", runbook.Name, firstNonEmpty(check.ID, check.Title), check.Command, err.Error()))
			severity := firstNonEmpty(strings.TrimSpace(check.Severity), "medium")
			taskID, isCreated, logErr := a.LogIssue(sessionID, title, details, severity)
			if logErr != nil {
				return "", logErr
			}
			if isCreated {
				created++
			}
			taskRefs = append(taskRefs, fmt.Sprintf("#%d", taskID))
		}
	}
	if failed > 0 {
		if err := a.Finish(sessionID, "fail"); err != nil {
			return "", err
		}
	} else {
		if err := a.Finish(sessionID, "pass"); err != nil {
			return "", err
		}
	}
	var lines []string
	lines = append(lines, fmt.Sprintf("QA run complete for session %d", sessionID))
	lines = append(lines, fmt.Sprintf("environment=%s", envCfg.Name))
	lines = append(lines, fmt.Sprintf("profiles=%s", strings.Join(resolved, ",")))
	lines = append(lines, fmt.Sprintf("checks=%d passed=%d failed=%d", checked, passed, failed))
	lines = append(lines, fmt.Sprintf("new_tasks=%d", created))
	if len(taskRefs) > 0 {
		lines = append(lines, fmt.Sprintf("task_refs=%s", strings.Join(taskRefs, ",")))
	}
	return strings.Join(lines, "\n"), nil
}

func (a *Adapter) resolveAssets(env string, profiles []string) ([]string, qaassets.Environment, []qaassets.Runbook, error) {
	envName := strings.TrimSpace(env)
	if envName == "" {
		envName = strings.TrimSpace(a.cfg.QA.DefaultEnv)
	}
	envCfg, err := a.loader.LoadEnvironment(envName)
	if err != nil {
		return nil, qaassets.Environment{}, nil, err
	}
	resolvedProfiles := profiles
	if len(resolvedProfiles) == 0 {
		resolvedProfiles = append([]string{}, a.cfg.QA.DefaultProfiles...)
	}
	if len(resolvedProfiles) == 0 {
		return nil, qaassets.Environment{}, nil, fmt.Errorf("no QA profiles provided and qa.default_profiles is empty")
	}
	runbookNames := map[string]struct{}{}
	for _, profileName := range resolvedProfiles {
		profile, profileErr := a.loader.LoadProfile(profileName)
		if profileErr != nil {
			return nil, qaassets.Environment{}, nil, profileErr
		}
		for _, rb := range profile.Runbooks {
			runbookNames[strings.TrimSpace(rb)] = struct{}{}
		}
	}
	sorted := make([]string, 0, len(runbookNames))
	for name := range runbookNames {
		sorted = append(sorted, name)
	}
	sort.Strings(sorted)
	runbooks := make([]qaassets.Runbook, 0, len(sorted))
	for _, name := range sorted {
		runbook, runbookErr := a.loader.LoadRunbook(name)
		if runbookErr != nil {
			return nil, qaassets.Environment{}, nil, runbookErr
		}
		runbooks = append(runbooks, runbook)
	}
	return resolvedProfiles, envCfg, runbooks, nil
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

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func runCommand(command string, vars map[string]string) error {
	cmd := exec.Command("bash", "-lc", command)
	env := append([]string{}, os.Environ()...)
	for key, value := range vars {
		env = append(env, key+"="+value)
	}
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}
