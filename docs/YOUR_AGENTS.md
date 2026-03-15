# AGENTS.md (Optional Bootstrap)

`cleo setup` and `cleo update` auto-install the builtin `cleo` skill at `.agents/skills/cleo/SKILL.md`.

If your runtime still needs an `AGENTS.md` file to set defaults, keep it minimal:

```md
## Cleo Default
- Apply the `cleo` skill for repository workflows before using raw `gh` or manual git commands.
- Use `cleo` command output as the source of truth for next actions.
```

That is enough to replace the older copy/paste-heavy template.
