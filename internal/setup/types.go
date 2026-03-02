package setup

import "os"

type Options struct {
	AutoYes        bool
	NonInteractive bool
	SkipAuth       bool
}

type Wizard struct {
	Stdout  *os.File
	Stderr  *os.File
	Stdin   *os.File
	Options Options
}

func NewWizard(options Options) *Wizard {
	return &Wizard{Stdout: os.Stdout, Stderr: os.Stderr, Stdin: os.Stdin, Options: options}
}
