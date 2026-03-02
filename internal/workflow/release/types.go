package release

import "github.com/cafaye/cleo/internal/config"

type Input struct {
	Name string
	Args []string
}

type Plan struct {
	Name        string
	Version     string
	Description string
	ReadOnly    bool
}

type Result struct {
	Name    string
	Version string
}

type Verification struct {
	Checked bool
	Reason  string
}

type Actions interface {
	CheckGitClean() error
	EnsureReleaseMissing(version string) error
	Cut(version string) error
	Publish(version string, draft bool, generateNotes bool) error
	Verify(version string) error
	List(limit int) error
	Latest() error
}

type Options struct {
	DefaultDraft  bool
	GenerateNotes bool
	TagPrefix     string
}

func NewOptions(cfg *config.Config) Options {
	return Options{
		DefaultDraft:  cfg.Release.DefaultDraft,
		GenerateNotes: cfg.Release.GenerateNotes,
		TagPrefix:     cfg.Release.TagPrefix,
	}
}
