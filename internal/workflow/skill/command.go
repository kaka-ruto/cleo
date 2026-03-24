package skill

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cafaye/cleo/internal/skills"
	"github.com/cafaye/cleo/internal/skills/registry"
)

type Command struct {
	out      io.Writer
	resolver skills.Resolver
	registry registry.Client
}

func New() (*Command, error) {
	resolver, err := skills.NewResolver()
	if err != nil {
		return nil, err
	}
	return &Command{out: os.Stdout, resolver: resolver, registry: registry.NewClient()}, nil
}

func newForTest(out io.Writer, resolver skills.Resolver) *Command {
	return &Command{out: out, resolver: resolver, registry: registry.NewClient()}
}

func newForTestWithRegistry(out io.Writer, resolver skills.Resolver, rc registry.Client) *Command {
	return &Command{out: out, resolver: resolver, registry: rc}
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
	case "uninstall":
		return c.uninstall(args)
	case "registry":
		return c.registryCmd(args)
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
		return errors.New("usage: cleo skill install <name> [--global|--project] [--registry <name>] [--force]")
	}
	name := args[0]
	root, err := c.installRoot(args[1:])
	if err != nil {
		return err
	}
	opts, err := parseInstallOptions(args[1:])
	if err != nil {
		return err
	}
	if opts.Registry != "" {
		def, ok, err := registry.ResolveDefinition(c.resolver.Home, opts.Registry)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("unknown registry: %s", opts.Registry)
		}
		target, err := c.registry.InstallSkill(def, name, root, opts.Force)
		if err != nil {
			return err
		}
		fmt.Fprintf(c.out, "Installed skill %s from registry %s to %s\n", strings.ToLower(name), def.Name, target)
		return nil
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

func (c *Command) uninstall(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: cleo skill uninstall <name> [--global|--project]")
	}
	root, err := c.installRoot(args[1:])
	if err != nil {
		return err
	}
	name := strings.TrimSpace(strings.ToLower(args[0]))
	if name == "" {
		return errors.New("skill name is required")
	}
	target := filepath.Join(root, name)
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return fmt.Errorf("skill not installed: %s", name)
	}
	if err := os.RemoveAll(target); err != nil {
		return fmt.Errorf("remove skill: %w", err)
	}
	fmt.Fprintf(c.out, "Uninstalled skill %s from %s\n", name, target)
	return nil
}

func (c *Command) registryCmd(args []string) error {
	if len(args) == 0 || args[0] == "list" {
		rows, err := registry.AllDefinitions(c.resolver.Home)
		if err != nil {
			return err
		}
		for _, r := range rows {
			fmt.Fprintf(c.out, "%s\t%s\t%s@%s:%s\n", r.Name, r.Description, r.Repo, r.Ref, r.Path)
		}
		return nil
	}
	if args[0] == "add" {
		return c.registryAdd(args[1:])
	}
	if args[0] == "remove" {
		return c.registryRemove(args[1:])
	}
	if args[0] != "skills" {
		return fmt.Errorf("unknown registry command: %s", args[0])
	}
	if len(args) < 2 {
		return errors.New("usage: cleo skill registry skills <registry> [--search <term>]")
	}
	def, ok, err := registry.ResolveDefinition(c.resolver.Home, args[1])
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("unknown registry: %s", args[1])
	}
	search := ""
	if len(args) > 2 {
		if len(args) != 4 || args[2] != "--search" {
			return errors.New("usage: cleo skill registry skills <registry> [--search <term>]")
		}
		search = strings.TrimSpace(strings.ToLower(args[3]))
	}
	skills, err := c.registry.ListSkills(def)
	if err != nil {
		return err
	}
	names := make([]string, 0, len(skills))
	for _, s := range skills {
		if search != "" && !strings.Contains(s.Name, search) {
			continue
		}
		names = append(names, s.Name)
	}
	sort.Strings(names)
	for _, n := range names {
		fmt.Fprintln(c.out, n)
	}
	return nil
}

func (c *Command) registryAdd(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: cleo skill registry add <name> --repo <owner/repo> --path <path> [--ref <ref>] [--description <text>]")
	}
	d := registry.Definition{Name: args[0]}
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--repo":
			if i+1 >= len(args) {
				return errors.New("--repo requires a value")
			}
			d.Repo = args[i+1]
			i++
		case "--path":
			if i+1 >= len(args) {
				return errors.New("--path requires a value")
			}
			d.Path = args[i+1]
			i++
		case "--ref":
			if i+1 >= len(args) {
				return errors.New("--ref requires a value")
			}
			d.Ref = args[i+1]
			i++
		case "--description":
			if i+1 >= len(args) {
				return errors.New("--description requires a value")
			}
			d.Description = args[i+1]
			i++
		default:
			return fmt.Errorf("unknown flag: %s", args[i])
		}
	}
	if err := registry.UpsertCustom(c.resolver.Home, d); err != nil {
		return err
	}
	fmt.Fprintf(c.out, "Saved registry %s (%s@%s:%s)\n", strings.ToLower(strings.TrimSpace(d.Name)), strings.TrimSpace(d.Repo), strings.TrimSpace(d.Ref), strings.TrimSpace(d.Path))
	return nil
}

func (c *Command) registryRemove(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: cleo skill registry remove <name>")
	}
	removed, err := registry.RemoveCustom(c.resolver.Home, args[0])
	if err != nil {
		return err
	}
	if !removed {
		return fmt.Errorf("registry not found: %s", args[0])
	}
	fmt.Fprintf(c.out, "Removed registry %s\n", strings.ToLower(strings.TrimSpace(args[0])))
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
	for i := 0; i < len(flags); i++ {
		f := flags[i]
		switch f {
		case "--project":
			project = true
		case "--global":
			global = true
		case "--force":
		case "--registry":
			if i+1 >= len(flags) {
				return "", errors.New("--registry requires a value")
			}
			i++
		default:
			if strings.HasPrefix(f, "--registry=") {
				continue
			}
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

type installOptions struct {
	Registry string
	Force    bool
}

func parseInstallOptions(flags []string) (installOptions, error) {
	out := installOptions{}
	for i := 0; i < len(flags); i++ {
		f := flags[i]
		switch f {
		case "--global", "--project":
		case "--force":
			out.Force = true
		case "--registry":
			if i+1 >= len(flags) {
				return out, errors.New("--registry requires a value")
			}
			out.Registry = strings.ToLower(strings.TrimSpace(flags[i+1]))
			i++
		default:
			if strings.HasPrefix(f, "--registry=") {
				out.Registry = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(f, "--registry=")))
				continue
			}
			return out, fmt.Errorf("unknown flag: %s", f)
		}
	}
	return out, nil
}
