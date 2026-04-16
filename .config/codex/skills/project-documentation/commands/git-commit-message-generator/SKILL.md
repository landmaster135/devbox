---
name: git-commit-message-generator
description: Create Git commit message with tag in English. Don't use parentheses. git_dir='.' staged_only=true. Remember displaying that message to user.
---

# Git Commit Message Creation

## Workflow
1. If `git_dir` directory is not provided, ask for the directory.
2. Set “staged_only” to false if there are no specified instructions.
3. Retrieve Git diff in the specified directory with just using the agent: `get_git_diff`.
4. Create Git commit message in English on a single line with tags.
5. Display that commit message to user.

## Requirements
- Add one of the following prefixes to the beginning of Git commit messages: "feat:", "refactor:", "fix:", "test:", or "docs:"
- Do NOT use tags: "chore:".
- Sample for Git commit message is here:
  - refactor: split PC stats API endpoints into separate files for better maintainability
  - feat: add mackerels update endpoint with PATCH /mackerels/update
  - fix: improve SQL injection protection and refactor query building in web clips service
  - feat: add mackerels masters append API endpoint for managing category, flavor, and delicious meaning datakill
