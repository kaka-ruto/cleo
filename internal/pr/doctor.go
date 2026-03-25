package pr

import (
	"fmt"

	"github.com/kaka-ruto/cleo/internal/ghcli"
)

func (s *Service) Doctor() error {
	if _, err := s.gh.Run("auth", "status"); err != nil {
		return fmt.Errorf("github auth check failed: %w", err)
	}
	if _, err := s.gh.Run("repo", "view", s.repo(), "--json", "nameWithOwner"); err != nil {
		return fmt.Errorf("repo access check failed for %s: %w", s.repo(), err)
	}
	if s.cfg.PR.DeployWatch.Enabled {
		if err := s.checkDeployWorkflow(); err != nil {
			return err
		}
	}
	fmt.Println("PR doctor passed.")
	return nil
}

func (s *Service) checkDeployWorkflow() error {
	out, err := s.gh.Run("workflow", "list", "--repo", s.repo(), "--json", "name,path,state")
	if err != nil {
		return fmt.Errorf("deploy workflow check failed: %w", err)
	}
	var workflows []struct {
		Name  string `json:"name"`
		Path  string `json:"path"`
		State string `json:"state"`
	}
	if err := ghcli.DecodeJSON(out, &workflows); err != nil {
		return err
	}
	for _, wf := range workflows {
		if wf.Name == s.cfg.PR.DeployWatch.Workflow {
			if wf.State == "disabled_manually" || wf.State == "disabled_inactivity" {
				return fmt.Errorf("deploy workflow %q exists but is disabled", s.cfg.PR.DeployWatch.Workflow)
			}
			return nil
		}
	}
	return fmt.Errorf("deploy workflow %q not found in %s", s.cfg.PR.DeployWatch.Workflow, s.repo())
}
