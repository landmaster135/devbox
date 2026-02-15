# README

## Symbolic links
`CLAUDE.md` is a symbolic link to `AGENTS.md`
```bash
ln -s AGENTS.md CLAUDE.md
```

`.claude/` is a symbolic link to `.codex/`
```bash
ln -s ../.codex/commands .claude/commands
ln -s ../.codex/skills   .claude/skills
```
