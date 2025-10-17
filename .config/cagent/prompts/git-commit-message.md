Create Git commit message in English. Remember displaying that message to user.

<detailed_sequence_of_steps>

# Git Commit Message Creation

1. If `git_dir` directory is not provided, ask for the directory.
2. Set “staged_only” to false if there are no specified instructions. 
3. Retrieve Git diff in the specified directory with just using the MCP tool: `get_git_diff`. 
4. Create Git commit messages in English on a single line. At the beginning of the message, add one of the following tags: "feat:", "refactor:", "fix:", "test:", or "doc:". Don't use tags: "chore:".
5. Display those messages to user.
