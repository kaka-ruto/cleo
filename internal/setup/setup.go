package setup

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Wizard struct {
	Stdout *os.File
	Stderr *os.File
	Stdin  *os.File
}

func NewWizard() *Wizard {
	return &Wizard{Stdout: os.Stdout, Stderr: os.Stderr, Stdin: os.Stdin}
}

func (w *Wizard) Run() error {
	w.title("Cleo Setup")
	if err := w.checkOrInstall("git"); err != nil {
		return err
	}
	if err := w.checkOrInstall("gh"); err != nil {
		return err
	}
	if err := w.checkOrInstall("gum"); err != nil {
		return err
	}
	if err := w.ensureGitHubAuth(); err != nil {
		return err
	}
	if err := w.writeConfig(); err != nil {
		return err
	}
	fmt.Fprintln(w.Stdout, "Setup complete. Next: cleo pr status <pr>")
	return nil
}

func (w *Wizard) title(text string) {
	if hasCommand("gum") {
		_ = runStreaming(w.Stdin, w.Stdout, w.Stderr, "gum", "style", "--border", "rounded", "--padding", "1 2", text)
		return
	}
	fmt.Fprintln(w.Stdout, text)
}

func (w *Wizard) checkOrInstall(bin string) error {
	if hasCommand(bin) {
		fmt.Fprintf(w.Stdout, "[ok] %s\n", bin)
		return nil
	}
	fmt.Fprintf(w.Stdout, "[missing] %s\n", bin)
	ok, err := w.confirm(fmt.Sprintf("Install %s now?", bin))
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("required dependency %q not installed", bin)
	}
	cmd, args, err := installCommand(bin)
	if err != nil {
		return err
	}
	fmt.Fprintf(w.Stdout, "Installing %s with: %s %s\n", bin, cmd, strings.Join(args, " "))
	if err := runStreaming(w.Stdin, w.Stdout, w.Stderr, cmd, args...); err != nil {
		return err
	}
	if !hasCommand(bin) {
		return fmt.Errorf("installation completed but %q is still not available", bin)
	}
	return nil
}

func (w *Wizard) ensureGitHubAuth() error {
	if err := runStreaming(w.Stdin, w.Stdout, w.Stderr, "gh", "auth", "status"); err == nil {
		return nil
	}
	ok, err := w.confirm("GitHub auth is missing. Run `gh auth login` now?")
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("github auth required")
	}
	return runStreaming(w.Stdin, w.Stdout, w.Stderr, "gh", "auth", "login")
}

func (w *Wizard) writeConfig() error {
	if _, err := os.Stat("cleo.yml"); err == nil {
		overwrite, err := w.confirm("cleo.yml already exists. Overwrite?")
		if err != nil {
			return err
		}
		if !overwrite {
			fmt.Fprintln(w.Stdout, "Keeping existing cleo.yml")
			return nil
		}
	}
	repo, err := discoverRepoSlug()
	if err != nil {
		return err
	}
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid GitHub repo slug: %s", repo)
	}
	content := defaultConfig(parts[0], parts[1])
	if err := os.WriteFile("cleo.yml", []byte(content), 0o644); err != nil {
		return err
	}
	fmt.Fprintln(w.Stdout, "Wrote cleo.yml")
	return nil
}

func (w *Wizard) confirm(question string) (bool, error) {
	if hasCommand("gum") && isTerminal(w.Stdin) && isTerminal(w.Stdout) {
		err := runStreaming(w.Stdin, w.Stdout, w.Stderr, "gum", "confirm", question)
		if err == nil {
			return true, nil
		}
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return false, nil
		}
		return false, err
	}
	if !isTerminal(w.Stdin) {
		fmt.Fprintf(w.Stdout, "%s [auto: no]\n", question)
		return false, nil
	}
	fmt.Fprintf(w.Stdout, "%s [y/N]: ", question)
	r := bufio.NewReader(w.Stdin)
	line, err := r.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			fmt.Fprintln(w.Stdout, "no")
			return false, nil
		}
		return false, err
	}
	ans := strings.ToLower(strings.TrimSpace(line))
	return ans == "y" || ans == "yes", nil
}

func isTerminal(file *os.File) bool {
	if file == nil {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}

func hasCommand(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func runStreaming(stdin *os.File, stdout *os.File, stderr *os.File, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func installCommand(bin string) (string, []string, error) {
	switch runtime.GOOS {
	case "darwin":
		if hasCommand("brew") {
			return "brew", []string{"install", bin}, nil
		}
		return "", nil, fmt.Errorf("homebrew is required to install %s automatically", bin)
	case "linux":
		if hasCommand("apt-get") {
			return "sudo", []string{"apt-get", "install", "-y", bin}, nil
		}
		if hasCommand("dnf") {
			return "sudo", []string{"dnf", "install", "-y", bin}, nil
		}
		if hasCommand("yum") {
			return "sudo", []string{"yum", "install", "-y", bin}, nil
		}
		return "", nil, fmt.Errorf("no supported package manager found for auto-install on linux")
	default:
		return "", nil, fmt.Errorf("auto-install not supported on %s", runtime.GOOS)
	}
}

func discoverRepoSlug() (string, error) {
	out, err := exec.Command("gh", "repo", "view", "--json", "nameWithOwner", "--jq", ".nameWithOwner").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("detect repo slug: %s", strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

func defaultConfig(owner, repo string) string {
	return fmt.Sprintf(`version: 1

github:
  owner: %s
  repo: %s
  host: github.com
  base_branch: main
  merge_method: merge
  delete_branch_on_merge: false

pr:
  required_approvals: 1
  require_non_draft: true
  require_mergeable: true
  block_if_requested_changes: true
  allow_self_approval: false

  checks:
    mode: required
    required: []
    ignore: []
    treat_neutral_as_pass: false
    timeout_seconds: 1800
    poll_interval_seconds: 10

  deploy_watch:
    enabled: true
    workflow: Deploy to Production
    branch: main
    timeout_seconds: 2700
    poll_interval_seconds: 10

  post_merge:
    enabled: true
    require_command_allowlist: false
    command_allowlist_prefixes:
      - "bin/kamal"
    command_denylist:
      - "rm -rf /"
    markers:
      start: "<!-- post-merge-commands:start -->"
      end: "<!-- post-merge-commands:end -->"
    allow_none: true

  stack:
    rebase_next_after_merge: true
    auto_detect_next_pr: true
    force_with_lease: true

safety:
  require_explicit_apply: true
`, owner, repo)
}
