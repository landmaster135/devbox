Create Git commit message in English. Remember displaying that message to user.

<detailed_sequence_of_steps>

# Git Commit Message Creation

## Workflow
1. If `git_dir` directory is not provided, ask for the directory.
2. Set “staged_only” to false if there are no specified instructions. 
3. Retrieve Git diff in the specified directory with just using the agent: `get_git_diff`. 
4. Create Git commit message in English on a single line with tags.
5. Display that commit message to user.

## Requirements
- Git commit message must have one of the following tags: "feat:", "refactor:", "fix:", "test:", or "doc:" at the beginning of the message without parentheses.
- Do NOT use tags: "chore:".
- Sample for Git commit message is here:
  - refactor: split PC stats API endpoints into separate files for better maintainability
  - feat: add mackerels update endpoint with PATCH /mackerels/update
  - fix: improve SQL injection protection and refactor query building in web clips service
  - feat: add mackerels masters append API endpoint for managing category, flavor, and delicious meaning datakill 

</detailed_sequence_of_steps>
