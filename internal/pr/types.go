package pr

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kaka-ruto/cleo/internal/config"
	"github.com/kaka-ruto/cleo/internal/ghcli"
)

type Runner interface {
	Run(args ...string) (string, error)
}

type Service struct {
	gh  Runner
	cfg *config.Config
}

func NewService(cfg *config.Config) *Service {
	return &Service{gh: ghcli.New(), cfg: cfg}
}

func NewServiceWithRunner(cfg *config.Config, runner Runner) *Service {
	return &Service{gh: runner, cfg: cfg}
}

type PRView struct {
	Number            int     `json:"number"`
	Title             string  `json:"title"`
	URL               string  `json:"url"`
	State             string  `json:"state"`
	IsDraft           bool    `json:"isDraft"`
	Mergeable         string  `json:"mergeable"`
	ReviewDecision    string  `json:"reviewDecision"`
	BaseRefName       string  `json:"baseRefName"`
	HeadRefName       string  `json:"headRefName"`
	StatusCheckRollup []Check `json:"statusCheckRollup"`
}

type Check struct {
	Name         string `json:"name"`
	WorkflowName string `json:"workflowName"`
	Status       string `json:"status"`
	Conclusion   string `json:"conclusion"`
	URL          string `json:"url"`
}

func (s *Service) repo() string {
	return s.cfg.GitHub.Owner + "/" + s.cfg.GitHub.Repo
}

func (s *Service) Get(pr string) (*PRView, error) {
	out, err := s.gh.Run("pr", "view", pr, "--repo", s.repo(), "--json", "number,title,url,state,isDraft,mergeable,reviewDecision,baseRefName,headRefName,statusCheckRollup")
	if err != nil {
		return nil, err
	}
	var view PRView
	if err := ghcli.DecodeJSON(out, &view); err != nil {
		return nil, err
	}
	return &view, nil
}

func valueOr(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}

func parseNumber(raw string) (int, error) {
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid PR number: %s", raw)
	}
	return n, nil
}
