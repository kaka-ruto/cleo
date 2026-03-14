package skill

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cafaye/cleo/internal/skills"
)

type Command struct {
	out      io.Writer
	resolver skills.Resolver
}

func New() (*Command, error) {
	resolver, err := skills.NewResolver()
	if err != nil {
		return nil, err
	}
	return &Command{out: os.Stdout, resolver: resolver}, nil
}

func newForTest(out io.Writer, resolver skills.Resolver) *Command {
	return &Command{out: out, resolver: resolver}
}

func (c *Command) Execute(name string, args []string) error {
	switch name {
	case "list":
		return c.list()
	case "use":
		if len(args) == 0 {
			return errors.New("usage: cleo skill use <name>")
		}
		return c.use(args[0])
	case "customize":
		if len(args) == 0 {
			return errors.New("usage: cleo skill customize <name>")
		}
		return c.customize(args[0])
	case "check":
		skillName := ""
		if len(args) > 0 {
			skillName = args[0]
		}
		return c.check(skillName)
	case "install":
		return c.install(args)
	case "sync":
		return c.sync(args)
	default:
		return fmt.Errorf("unknown skill command: %s", name)
	}
}

func (c *Command) list() error {
	list, err := c.resolver.List()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		fmt.Fprintln(c.out, "No skills found.")
		return nil
	}
	for _, s := range list {
		fmt.Fprintf(c.out, "%s\t%s\t%s\n", s.Name, s.Origin, s.Path)
	}
	return nil
}

func (c *Command) use(name string) error {
	src, body, err := c.resolver.Resolve(name)
	if err != nil {
		return err
	}
	if err := skills.ValidateForUse(body, src.Name); err != nil {
		return fmt.Errorf("%s: %w", src.Path, err)
	}
	fmt.Fprintf(c.out, "# source: %s (%s)\n\n", src.Path, src.Origin)
	fmt.Fprintln(c.out, strings.TrimSpace(string(body)))
	return nil
}

func (c *Command) customize(name string) error {
	path, err := c.resolver.Customize(name)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.out, "Customized skill written to %s\n", path)
	return nil
}

func (c *Command) check(name string) error {
	rows, err := c.resolver.Check(name)
	if err != nil {
		return err
	}
	if name != "" {
		fmt.Fprintf(c.out, "Skill %s is valid.\n", name)
		return nil
	}
	fmt.Fprintf(c.out, "Checked %d skill(s): all valid.\n", len(rows))
	return nil
}

func (c *Command) install(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: cleo skill install <name> [--global|--project]")
	}
	name := args[0]
	root, err := c.installRoot(args[1:])
	if err != nil {
		return err
	}
	src, body, err := c.resolver.Resolve(name)
	if err != nil {
		return err
	}
	target, err := writeSkill(root, src.Name, body)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.out, "Installed skill %s to %s\n", src.Name, target)
	return nil
}

func (c *Command) sync(args []string) error {
	root, err := c.installRoot(args)
	if err != nil {
		return err
	}
	builtins := skills.BuiltinList()
	if len(builtins) == 0 {
		fmt.Fprintln(c.out, "No builtin skills to sync.")
		return nil
	}
	count := 0
	for _, s := range builtins {
		body, err := skills.ReadBuiltin(s.Name)
		if err != nil {
			return err
		}
		if _, err := writeSkill(root, s.Name, body); err != nil {
			return err
		}
		count++
	}
	fmt.Fprintf(c.out, "Synced %d skill(s) to %s\n", count, root)
	return nil
}

func (c *Command) installRoot(flags []string) (string, error) {
	project := false
	global := false
	for _, f := range flags {
		switch f {
		case "--project":
			project = true
		case "--global":
			global = true
		default:
			return "", fmt.Errorf("unknown flag: %s", f)
		}
	}
	if project && global {
		return "", errors.New("choose only one target: --global or --project")
	}
	if project {
		return filepath.Join(c.resolver.Cwd, ".agents", "skills"), nil
	}
	return filepath.Join(c.resolver.Home, ".agents", "skills"), nil
}

func writeSkill(root string, name string, body []byte) (string, error) {
	target := filepath.Join(root, name, "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return "", fmt.Errorf("create skill directory: %w", err)
	}
	if err := os.WriteFile(target, body, 0o644); err != nil {
		return "", fmt.Errorf("write skill: %w", err)
	}
	return target, nil
}
