package setup

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (w *Wizard) title(text string) {
	if hasCommand("gum") {
		_ = runStreaming(w.Stdin, w.Stdout, w.Stderr, "gum", "style", "--border", "rounded", "--padding", "1 2", text)
		return
	}
	fmt.Fprintln(w.Stdout, text)
}

func (w *Wizard) confirm(question string) (bool, error) {
	if w.Options.AutoYes {
		fmt.Fprintf(w.Stdout, "%s [auto: yes]\n", question)
		return true, nil
	}
	if w.Options.NonInteractive || !isTerminal(w.Stdin) {
		fmt.Fprintf(w.Stdout, "%s [auto: no]\n", question)
		return false, nil
	}
	if hasCommand("gum") && isTerminal(w.Stdout) {
		err := runStreaming(w.Stdin, w.Stdout, w.Stderr, "gum", "confirm", question)
		if err == nil {
			return true, nil
		}
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return false, nil
		}
		return false, err
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

func pathContains(dir string) bool {
	for _, entry := range filepath.SplitList(os.Getenv("PATH")) {
		if filepath.Clean(entry) == filepath.Clean(dir) {
			return true
		}
	}
	return false
}

func copyExecutable(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Chmod(0o755)
}
