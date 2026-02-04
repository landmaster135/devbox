---
name: create-pr
description: Creates GitHub pull requests with properly formatted titles that pass the check-pr-title CI validation. Use when creating PRs, submitting changes for review, or when the user says /pr or asks to create a pull request.
allowed-tools: Bash(git:*), Bash(gh:*), Read, Grep, Glob
---

# Create Pull Request

Creates GitHub PRs with titles that pass n8n's `check-pr-title` CI validation.

## Information Gathering

1. **Repository target**: Confirm the GitHub repository URL where the pull request must be opened. If the user does not provide it up front, explicitly ask for the URL before proceeding.
2. **Push confirmation**: At the very start, explicitly ask the user whether the current branch still needs to be pushed (or force-pushed) before opening the PR.
3. **Source selection**: Determine whether the PR description should be based on the CHANGELOG or on Git commit history.
  - If the user already specified the source, acknowledge it and proceed.
  - If the user did not specify, explicitly ask: "Should I base the PR summary on the CHANGELOG or on Git commit history?"
4. **CHANGELOG context**: When the user chooses CHANGELOG, ask which headings/sections (e.g., `## Added`, `## Fixed`) should be used to craft the summary before drafting any text.
5. **Git history window**: When the user chooses Git commit history, request the starting point for the log as an absolute date/time string (ISO 8601 or `YYYY-MM-DD HH:MM TZ`). Explain that this value feeds into commands such as `git log --since "<date/time>"` to scope the summary accurately.
6. **Clarifications**: If users provide ambiguous headings, date ranges, or repository URLs, follow up until the scope is concrete. Do not proceed without the required section names, start date/time, and repo URL.

## PR Title Format

```
<type>(<scope>): <summary>
```

### Types (required)

| Type       | Description                                      | Changelog |
|------------|--------------------------------------------------|-----------|
| `feat`     | New feature                                      | Yes       |
| `fix`      | Bug fix                                          | Yes       |
| `perf`     | Performance improvement                          | Yes       |
| `test`     | Adding/correcting tests                          | No        |
| `docs`     | Documentation only                               | No        |
| `refactor` | Code change (no bug fix or feature)              | No        |
| `build`    | Build system or dependencies                     | No        |
| `ci`       | CI configuration                                 | No        |
| `chore`    | Routine tasks, maintenance                       | No        |

### Scopes (optional but recommended)

- `API` - Public API changes
- `benchmark` - Benchmark CLI changes
- `core` - Core/backend/private API
- `editor` - Editor UI changes
- `* Node` - Specific node (e.g., `Slack Node`, `GitHub Node`)

### Summary Rules

- Use imperative present tense: "Add" not "Added"
- Capitalize first letter
- No period at the end
- No ticket IDs (e.g., N8N-1234)
- Add `(no-changelog)` suffix to exclude from changelog

## Steps

1. **Collect required context** (see Information Gathering section) before inspecting the repository.
2. **Check current state**:
  ```bash
  git status
  git diff --stat
  git log origin/master..HEAD --oneline
  ```
3. **Analyze changes** to determine:
  - Type: What kind of change is this?
  - Scope: Which package/area is affected?
  - Summary: What does the change do?
  - Source-based summary: Use the specified CHANGELOG headings or `git log --since "<date/time>"` output as instructed by the user.
4. **Push branch if needed** (Skip the push if not needed, as the changes might already be reflected in the remote repository.):
  ```bash
  git push -u origin HEAD
  ```
5. **Create PR** by invoking the GitHub MCP tools directly (without relying on the gh CLI or external config references):
  - Build the PR body from `.github/pull_request_template.md`, filling in each section the same way you would when using the CLI.
  - Invoke the `create_pull_request_with_current_branch` tool with the composed title/body (set `draft: false` unless the user directs otherwise).
  - Immediately call the `list_pull_requests` tool to confirm the new PR appears in the repository and report the resulting URL/reference back to the user.

## PR Body Guidelines

Based on `.github/pull_request_template.md`:

### Summary Section
- Describe what the PR does
- Explain how to test the changes
- Include screenshots/videos for UI changes
- Reference the chosen CHANGELOG headings or git history window in prose if it adds clarity

### Related Links Section
- Link to Linear ticket: `https://linear.app/n8n/issue/[TICKET-ID]`
- Link to GitHub issues using keywords to auto-close:
  - `closes #123` / `fixes #123` / `resolves #123`
- Link to Community forum posts if applicable

### Checklist
All items should be addressed before merging:
- PR title follows conventions
- Docs updated or follow-up ticket created
- Tests included (bugs need regression tests, features need coverage)
- `release/backport` label added if urgent fix needs backporting

## Examples

### Feature in editor
```
feat(editor): Add workflow performance metrics display
```

### Bug fix in core
```
fix(core): Resolve memory leak in execution engine
```

### Node-specific change
```
fix(Slack Node): Handle rate limiting in message send
```

### Breaking change (add exclamation mark before colon)
```
feat(API)!: Remove deprecated v1 endpoints
```

### No changelog entry
```
refactor(core): Simplify error handling (no-changelog)
```

### No scope (affects multiple areas)
```
chore: Update dependencies to latest versions
```

## Validation

The PR title must match this pattern:
```
^(feat|fix|perf|test|docs|refactor|build|ci|chore|revert)(\([a-zA-Z0-9 ]+( Node)?\))?!?: [A-Z].+[^.]$
```

Key validation rules:
- Type must be one of the allowed types
- Scope is optional but must be in parentheses if present
- Exclamation mark for breaking changes goes before the colon
- Summary must start with capital letter
- Summary must not end with a period
