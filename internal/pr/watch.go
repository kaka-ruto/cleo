package pr

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kaka-ruto/cleo/internal/ghcli"
)

func (s *Service) Watch(ref string) error {
	sha, err := s.resolveSHA(ref)
	if err != nil {
		return err
	}
	deadline := time.Now().Add(time.Duration(s.cfg.PR.DeployWatch.TimeoutSeconds) * time.Second)
	poll := time.Duration(s.cfg.PR.DeployWatch.PollIntervalSeconds) * time.Second
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for deploy workflow for sha=%s", sha)
		}
		runID, err := s.findWorkflowRun(sha)
		if err != nil {
			return err
		}
		if runID > 0 {
			_, err := s.gh.Run("run", "watch", strconv.Itoa(runID), "--repo", s.repo(), "--exit-status")
			return err
		}
		time.Sleep(poll)
	}
}

func (s *Service) resolveSHA(ref string) (string, error) {
	if _, err := strconv.Atoi(ref); err != nil {
		return ref, nil
	}
	out, err := s.gh.Run("pr", "view", ref, "--repo", s.repo(), "--json", "mergeCommit")
	if err != nil {
		return "", err
	}
	var payload struct {
		MergeCommit struct {
			OID string `json:"oid"`
		} `json:"mergeCommit"`
	}
	if err := ghcli.DecodeJSON(out, &payload); err != nil {
		return "", err
	}
	if payload.MergeCommit.OID == "" {
		return "", fmt.Errorf("PR #%s has no merge commit yet", ref)
	}
	return payload.MergeCommit.OID, nil
}

func (s *Service) findWorkflowRun(sha string) (int, error) {
	out, err := s.gh.Run("run", "list", "--repo", s.repo(), "--workflow", s.cfg.PR.DeployWatch.Workflow, "--branch", s.cfg.PR.DeployWatch.Branch, "--limit", "30", "--json", "databaseId,headSha")
	if err != nil {
		return 0, err
	}
	var runs []struct {
		DatabaseID int    `json:"databaseId"`
		HeadSHA    string `json:"headSha"`
	}
	if err := ghcli.DecodeJSON(out, &runs); err != nil {
		return 0, err
	}
	for _, r := range runs {
		if r.HeadSHA == sha {
			return r.DatabaseID, nil
		}
	}
	return 0, nil
}
