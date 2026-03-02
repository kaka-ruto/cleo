package pr

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cafaye/cleo/internal/ghcli"
)

func (s *Service) Run(pr string, dry bool) error {
	out, err := s.gh.Run("pr", "view", pr, "--repo", s.repo(), "--json", "body")
	if err != nil {
		return err
	}
	var payload struct {
		Body string `json:"body"`
	}
	if err := ghcli.DecodeJSON(out, &payload); err != nil {
		return err
	}
	cmds, err := parseCommands(payload.Body, s.cfg.PR.PostMerge.Markers.Start, s.cfg.PR.PostMerge.Markers.End, s.cfg.PR.PostMerge.AllowNone)
	if err != nil {
		return err
	}
	for _, cmd := range cmds {
		if denied(cmd, s.cfg.PR.PostMerge.CommandDenylist) {
			return fmt.Errorf("denied post-merge command: %s", cmd)
		}
		if s.cfg.PR.PostMerge.RequireCommandAllowlist && !allowed(cmd, s.cfg.PR.PostMerge.CommandAllowlistPrefixes) {
			return fmt.Errorf("command not in allowlist: %s", cmd)
		}
	}
	if len(cmds) == 0 {
		fmt.Printf("No post-merge commands for PR #%s.\n", pr)
		return nil
	}
	for _, cmd := range cmds {
		fmt.Printf("+ %s\n", cmd)
		if dry {
			continue
		}
		r := exec.Command("bash", "-lc", cmd)
		r.Stdout = os.Stdout
		r.Stderr = os.Stderr
		if err := r.Run(); err != nil {
			return fmt.Errorf("post-merge command failed: %w", err)
		}
	}
	return nil
}

func parseCommands(body, startMarker, endMarker string, allowNone bool) ([]string, error) {
	start := strings.Index(body, startMarker)
	end := strings.Index(body, endMarker)
	if start < 0 || end < 0 || end <= start {
		return nil, fmt.Errorf("missing required post-merge command markers")
	}
	block := body[start+len(startMarker) : end]
	commands := []string{}
	for _, line := range strings.Split(block, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "- `") || !strings.HasSuffix(line, "`") {
			continue
		}
		cmd := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "- `"), "`"))
		if cmd == "" || strings.EqualFold(cmd, "none") {
			continue
		}
		commands = append(commands, cmd)
	}
	if len(commands) == 0 && !allowNone {
		return nil, fmt.Errorf("no post-merge commands found and allow_none=false")
	}
	return commands, nil
}

func denied(cmd string, denylist []string) bool {
	for _, v := range denylist {
		if v != "" && strings.Contains(cmd, v) {
			return true
		}
	}
	return false
}

func allowed(cmd string, prefixes []string) bool {
	for _, v := range prefixes {
		if strings.HasPrefix(strings.TrimSpace(cmd), v) {
			return true
		}
	}
	return false
}
