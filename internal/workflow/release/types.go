package release

import "github.com/kaka-ruto/cleo/internal/config"

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
	ValidateChangelog(version string) error
	Cut(version string) error
	Publish(version string, draft bool, generateNotes bool, notes NoteOverrides) error
	Verify(version string) error
	List(limit int) error
	Latest() error
}

type NoteOverrides struct {
	Summary         string
	Highlights      string
	BreakingChanges string
	MigrationNotes  string
	Verification    string
}

type Options struct {
	DefaultDraft  bool
	GenerateNotes bool
	TagPrefix     string
	ChangelogFile string
	BinaryName    string
	BuildTarget   string
}

func NewOptions(cfg *config.Config) Options {
	return Options{
		DefaultDraft:  cfg.Release.DefaultDraft,
		GenerateNotes: cfg.Release.GenerateNotes,
		TagPrefix:     cfg.Release.TagPrefix,
		ChangelogFile: cfg.Release.ChangelogFile,
		BinaryName:    cfg.Release.BinaryName,
		BuildTarget:   cfg.Release.BuildTarget,
	}
}
