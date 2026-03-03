package setup

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

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
		return w.fallbackInstall(bin, err)
	}
	if !hasCommand(bin) {
		return w.fallbackInstall(bin, fmt.Errorf("installation completed but %q is still not available", bin))
	}
	return nil
}

func (w *Wizard) ensurePlaywrightRuntime() error {
	if !hasCommand("go") {
		return fmt.Errorf("go is required to install Playwright runtime")
	}
	fmt.Fprintln(w.Stdout, "Ensuring Playwright Chromium runtime for QA browser actions...")
	if err := runStreaming(w.Stdin, w.Stdout, w.Stderr, "go", "run", "github.com/playwright-community/playwright-go/cmd/playwright@v0.5200.1", "install", "chromium"); err != nil {
		return fmt.Errorf("install playwright runtime: %w", err)
	}
	return nil
}

func (w *Wizard) fallbackInstall(bin string, installErr error) error {
	if bin != "gum" || !hasCommand("go") {
		return installErr
	}
	if err := w.installGumWithGo(); err != nil {
		return fmt.Errorf("%v; fallback install failed: %w", installErr, err)
	}
	return nil
}

func (w *Wizard) installGumWithGo() error {
	gobin := os.ExpandEnv("$HOME/.local/bin")
	cmd := exec.Command("go", "install", "github.com/charmbracelet/gum@latest")
	cmd.Env = append(os.Environ(), "GOBIN="+gobin)
	cmd.Stdin = w.Stdin
	cmd.Stdout = w.Stdout
	cmd.Stderr = w.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	path := os.Getenv("PATH")
	if !strings.Contains(path, gobin) {
		_ = os.Setenv("PATH", gobin+":"+path)
	}
	if !hasCommand("gum") {
		return fmt.Errorf("gum not found after go install")
	}
	fmt.Fprintf(w.Stdout, "Installed gum with go install into %s\n", gobin)
	return nil
}

func (w *Wizard) ensureGitHubAuth() error {
	if err := runStreaming(w.Stdin, w.Stdout, w.Stderr, "gh", "auth", "status"); err == nil {
		return nil
	}
	if w.Options.NonInteractive {
		fmt.Fprintln(w.Stdout, "GitHub auth is missing. Continuing non-interactive setup.")
		fmt.Fprintln(w.Stdout, "Next step: run `gh auth login` before using PR commands.")
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

func installCommand(bin string) (string, []string, error) {
	switch runtime.GOOS {
	case "darwin":
		if hasCommand("brew") {
			return "brew", []string{"install", bin}, nil
		}
		return "", nil, fmt.Errorf("homebrew is required to install %s automatically", bin)
	case "linux":
		pkg := linuxPackageName(bin)
		if hasCommand("apt-get") {
			return "sudo", []string{"apt-get", "install", "-y", pkg}, nil
		}
		if hasCommand("dnf") {
			return "sudo", []string{"dnf", "install", "-y", pkg}, nil
		}
		if hasCommand("yum") {
			return "sudo", []string{"yum", "install", "-y", pkg}, nil
		}
		return "", nil, fmt.Errorf("no supported package manager found for auto-install on linux")
	default:
		return "", nil, fmt.Errorf("auto-install not supported on %s", runtime.GOOS)
	}
}

func discoverRepoSlug() (string, error) {
	out, err := exec.Command("gh", "repo", "view", "--json", "nameWithOwner", "--jq", ".nameWithOwner").CombinedOutput()
	if err == nil {
		return strings.TrimSpace(string(out)), nil
	}
	repo, gitErr := discoverRepoSlugFromGitRemote()
	if gitErr == nil {
		return repo, nil
	}
	return "", fmt.Errorf("detect repo slug: %s", strings.TrimSpace(string(out)))
}

func linuxPackageName(bin string) string {
	switch bin {
	case "go":
		return "golang-go"
	case "gh":
		return "gh"
	case "node":
		return "nodejs"
	default:
		return bin
	}
}

func discoverRepoSlugFromGitRemote() (string, error) {
	out, err := exec.Command("git", "config", "--get", "remote.origin.url").CombinedOutput()
	if err != nil {
		return "", err
	}
	url := strings.TrimSpace(string(out))
	url = strings.TrimSuffix(url, ".git")
	url = strings.TrimPrefix(url, "git@github.com:")
	url = strings.TrimPrefix(url, "https://github.com/")
	parts := strings.Split(url, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid remote origin url: %s", url)
	}
	return parts[0] + "/" + parts[1], nil
}
