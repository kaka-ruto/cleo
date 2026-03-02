package setup

import (
	"fmt"
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
