---
name: changelog-generator
description: Create project-grade changelogs that mirror /CHANGELOG.md (version heading, PR line, Features/Improvements/Refactor/Bug Fixes/Documentation blocks, optional tables) by organizing commits from a user-specified start datetime. Always ask the user for the exact datetime (ISO 8601 or YYYY-MM-DD HH:MM) before starting.
---

# Changelog Creation Workflow

This skill walks you through producing release notes that follow `/CHANGELOG.md`. Commits from a target period are grouped into user-facing sections and rewritten in clear language.

## 0. Pre-flight Checklist

1. Briefly confirm the presence of `.agents/`, `.agents/instructions.md`, and `.agents/test_results.md` in the repo root.
2. Ask the user before touching git:
   - "From which exact datetime (ISO 8601 or `YYYY-MM-DD HH:MM`) should the changelog start?"
   - If relevant, ask for end datetime, target version, or PR number constraints.
3. Normalize the answer so it can drop directly into commands like `git log --since="<datetime>"`.

## 1. Gather Commits

- Typical command: `git log --since="<ISO 8601>" --pretty=format:"%h%x09%an%x09%ad%x09%s"`.
- Add `--until`, `--author`, or `-- <path>` when the brief calls for them.
- When comparing releases, combine `git log <from>..<to>` with the `--since/--until` filters as needed.

## 2. Analyze and Classify

Map commits onto the sections that `/CHANGELOG.md` already uses.

| Section | Common keywords | Example impact |
| --- | --- | --- |
| Features | feat, add, introduce | New CLIs or net-new capabilities |
| Improvements | perf, improve, polish | UX, speed, or reliability boosts |
| Refactor | refactor, reorganize | Structure or naming overhauls |
| Bug Fixes | fix, hotfix, revert | Bugs, regressions, exception handling |
| Documentation | docs, README | README/guide/Skill updates |
| Other | chore, test | Usually omit; include only if user-facing |

- If a commit spans multiple categories, pick the one that best conveys user value.
- Drop internal noise (CI plumbing, version bumps without substance) unless explicitly requested.

## 3. Output Format

Follow the exact structure from `/CHANGELOG.md`.

```
## v<version> — <YYYY-MM-DD>
PR: #<number> (use `PR: pending` if unknown)

### Features
- **component (scope)**: Reader-friendly blurb.

### Improvements
- ...
```

- When a section needs a long list of CLIs or services, reuse the table layout (`| service | cli | description |`).
- Style: Start each bullet with a bolded component (or CLI path) and keep explanations succinct.
- If a section has no applicable commits, omit it unless the user demands all sections—then write `- No notable updates.`

## 4. Writing Guidelines

- Lead with the customer benefit, then describe the mechanism or implementation.
- Quantify impact when possible ("2× faster", "adds ISO 8601 partial reads").
- Mention concrete paths like `cmd/cli/...` only when they clarify scope.
- Keep the narrative in Japanese when drafting the final changelog, matching the rest of `/CHANGELOG.md`.

## 5. Finalization

1. sanity-check that the number of described bullets roughly matches the commit count you analyzed.
2. Provide whatever the user asked for:
   - The final Markdown entry
   - Snippets from `git show <hash>` if they want supporting detail
   - Summaries of PR numbers, tags, or blast radius
3. Insert the new entry at the top of `/CHANGELOG.md` unless the user specifies another location.
4. If you ran tests (`go test`, scripts, etc.), record results and update `.agents/test_results.md` accordingly.

## 6. Reference Template

```markdown
## v0.07 — 2026-02-04
PR: #<pending>

### Features
- **new-cli (CLI)**: ...

### Improvements
- ...

### Bug Fixes
- ...

### Documentation
- ...
```

## 7. Conversation Flow Example

1. "Please provide the start datetime for this changelog in ISO 8601 (e.g., 2025-12-01T00:00:00Z)."
2. Follow up: "Do you also have an end datetime, version tag, or related PR I should target?"
3. Begin the workflow only after the user supplies the required datetime.
