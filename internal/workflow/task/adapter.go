package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kaka-ruto/cleo/internal/config"
	"github.com/kaka-ruto/cleo/internal/taskstore"
)

type Adapter struct {
	store      *taskstore.Store
	cfg        *config.Config
	now        func() time.Time
	runLocalFn func(name string, args ...string) (string, error)
}

func NewAdapter(store *taskstore.Store, cfg *config.Config) *Adapter {
	return &Adapter{store: store, cfg: cfg, now: time.Now, runLocalFn: runLocal}
}

func (a *Adapter) List(status string) (string, error) {
	tasks, err := a.store.ListTasks(context.Background(), status)
	if err != nil {
		return "", err
	}
	if len(tasks) == 0 {
		return "No tasks found.", nil
	}
	var b strings.Builder
	for _, task := range tasks {
		branch := task.WorkBranch
		if strings.TrimSpace(branch) == "" {
			branch = "-"
		}
		session, sessionErr := a.store.Session(context.Background(), task.SessionID)
		source := "unknown"
		if sessionErr == nil {
			source = fmt.Sprintf("%s:%s", session.Source, session.Ref)
		}
		fmt.Fprintf(&b, "#%d [%s] %s status=%s occurrences=%d source=%s branch=%s\n", task.ID, task.Severity, task.Title, task.Status, task.Occurrences, source, branch)
	}
	return strings.TrimRight(b.String(), "\n"), nil
}

func (a *Adapter) Show(id int64) (string, error) {
	task, err := a.store.Task(context.Background(), id)
	if err != nil {
		return "", err
	}
	branch := task.WorkBranch
	if strings.TrimSpace(branch) == "" {
		branch = "-"
	}
	session, sessionErr := a.store.Session(context.Background(), task.SessionID)
	source := "unknown"
	if sessionErr == nil {
		source = fmt.Sprintf("%s:%s", session.Source, session.Ref)
	}
	return fmt.Sprintf("Task #%d\nrepo=%s\nsource=%s\nstatus=%s severity=%s\noccurrences=%d\nwork_branch=%s\n\ntitle: %s\n\ndetails:\n%s", task.ID, task.RepoKey, source, task.Status, task.Severity, task.Occurrences, branch, task.Title, task.Details), nil
}

func (a *Adapter) Claim(id int64) error {
	return a.store.UpdateTaskStatus(context.Background(), id, "in_progress", a.now())
}

func (a *Adapter) Close(id int64) error {
	return a.store.UpdateTaskStatus(context.Background(), id, "closed", a.now())
}

func (a *Adapter) Work(id int64, opts WorkOptions) (string, error) {
	task, err := a.store.Task(context.Background(), id)
	if err != nil {
		return "", err
	}
	if err := a.store.UpdateTaskStatus(context.Background(), id, "in_progress", a.now()); err != nil {
		return "", err
	}
	current, err := a.runLocalFn("git", "branch", "--show-current")
	if err != nil {
		return "", err
	}
	currentBranch := strings.TrimSpace(current)
	if currentBranch == "" {
		return "", fmt.Errorf("cannot determine current branch")
	}

	onBaseBranch := currentBranch == strings.TrimSpace(a.cfg.GitHub.BaseBranch)
	if opts.ForceNewBranch && opts.ForceInPlace {
		return "", fmt.Errorf("--new-branch and --in-place cannot be used together")
	}
	if opts.ForceInPlace && onBaseBranch {
		return "", fmt.Errorf("--in-place cannot be used on base branch %q", strings.TrimSpace(a.cfg.GitHub.BaseBranch))
	}

	workBranch := currentBranch
	lane := "in-place"
	useNewBranch := onBaseBranch
	if opts.ForceNewBranch {
		useNewBranch = true
	}
	if opts.ForceInPlace {
		useNewBranch = false
	}
	if useNewBranch {
		workBranch = branchForTask(task.ID, task.Title)
		if _, err := a.runLocalFn("git", "checkout", "-b", workBranch); err != nil {
			if _, checkoutErr := a.runLocalFn("git", "checkout", workBranch); checkoutErr != nil {
				return "", err
			}
		}
		lane = "new-branch"
	}

	if err := a.store.SetTaskWorkBranch(context.Background(), id, workBranch, a.now()); err != nil {
		return "", err
	}
	eventPayload := fmt.Sprintf("lane=%s current=%s work=%s", lane, currentBranch, workBranch)
	if err := a.store.AddTaskEvent(context.Background(), id, "work_started", eventPayload, a.now()); err != nil {
		return "", err
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("Task %d is now in progress.", id))
	lines = append(lines, fmt.Sprintf("Work lane: %s", lane))
	lines = append(lines, fmt.Sprintf("Working branch: %s", workBranch))
	lines = append(lines, "")
	lines = append(lines, "Required auditability:")
	lines = append(lines, fmt.Sprintf("- Include commit trailer: Task: #%d", id))
	lines = append(lines, "- Keep commit scope limited to this task.")
	lines = append(lines, "")
	lines = append(lines, "Suggested next steps:")
	lines = append(lines, "1. Implement and commit with the required trailer.")
	lines = append(lines, "2. Run make quality.")
	if lane == "new-branch" {
		lines = append(lines, "3. Push branch and open PR with cleo pr create.")
	} else {
		lines = append(lines, "3. Update the current PR or open one with cleo pr create if needed.")
	}
	return strings.Join(lines, "\n"), nil
}
