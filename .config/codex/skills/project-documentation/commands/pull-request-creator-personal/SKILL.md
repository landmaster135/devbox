---
name: pull-request-creator-personal
description: Creates GitHub pull requests with properly formatted titles that pass the check-pr-title CI validation. Use when creating PRs, submitting changes for review, or when the user says /pr or asks to create a pull request.
allowed-tools: Bash(git:*), Read, Grep, Glob
---

# Create Pull Request

Creates GitHub PRs with titles.

## Pre-flight Checklist

1. **Push confirmation**: At the very start, explicitly ask the user whether the current branch still needs to be pushed (or force-pushed) before opening the PR.
2. **Repository target**: Confirm the GitHub repository URL where the pull request must be opened. If not provided, ask before proceeding.
3. **CHANGELOG baseline**: Ask whether there is an existing CHANGELOG entry/file/version to use as the primary source for PR summary. If yes, request the concrete path/version heading to reference.
4. **Commit history start point (fallback)**: Only when no CHANGELOG baseline is provided, request the date/time of the first Git commit to include in the PR as an absolute string (ISO 8601 or `YYYY-MM-DD HH:MM TZ`).
5. **Commit window confirmation (fallback)**: Only when no CHANGELOG baseline is provided, confirm that the PR summary will be based on commits after the provided start date/time from the current directory.
6. **Template baseline**: Confirm that the PR description will be composed from `.github/PULL_REQUEST_TEMPLATE/default.md` in the project root.
7. **Clarifications**: As the final information-gathering checkpoint, if users provide the following ambiguous scopes, follow up until the scope is concrete. Do not proceed without:
  - push decision
  - repository URLs
  - changelog decision (`use changelog` or `use git commits`)
  - changelog scope (file path/version heading) if changelog is used
  - commit-window agreement if git commits are used
  - date/time values if git commits are used

## PR Title Format

```
<type>: <summary>
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

### Summary Rules

- Use imperative present tense: "Add" not "Added"
- Capitalize first letter
- No period at the end
- No ticket IDs (e.g., N8N-1234)
- Add `(no-changelog)` suffix to exclude from changelog

### Summary Section
- Describe what the PR does
- Explain how to test the changes
- Include screenshots/videos for UI changes
- If a CHANGELOG baseline is provided, prioritize it as the summary source and use git history only for supplemental verification
- Reference the chosen CHANGELOG headings or git history window in prose if it adds clarity

### Examples

Feature in editor
```
feat: Add workflow performance metrics display
```

Bug fix
```
fix: Resolve memory leak in execution engine
```

Breaking change (add exclamation mark before colon)
```
feat!: Remove deprecated v1 endpoints
```

No changelog entry
```
refactor: Simplify error handling (no-changelog)
```

No scope (affects multiple areas)
```
chore: Update dependencies to latest versions
```

### Validation

The PR title must match this pattern:
```
^(feat|fix|perf|test|docs|refactor|build|ci|chore): [A-Z].+[^.]$
```

## PR Body

### Commit Classification Rules

- Classify each tool/service area by its dominant commit type in the selected commit window.
- Do not duplicate a tool across sections when it is already represented as `feat` in the same PR summary.
- If a tool has `feat:` commits, list it under feature-related sections only, and do not also list it as `refactor` or `test`.
- Use `refactor` and `test` sections for tools that are not primarily introduced as `feat` in that PR window.

### Related Links Section
- Link to GitHub issues using keywords to auto-close:
  - `closes #123` / `fixes #123` / `resolves #123`

## PR Creation
1. Use the `create_pull_request_with_current_branch` tool to create a new pull request.
2. Use the `list_pull_requests` tool to verify that the pull request was created successfully.
