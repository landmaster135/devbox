---
name: bash-script-senior
description: Implement and refactor production-grade Bash scripts with safe defaults, subroutine-first structure, deterministic output naming, and practical verification workflows.
---

# Bash Shell Script Senior

Use this skill when creating or refactoring bash scripts for build and operations tasks.

## Core policy

- Prioritize safe behavior and predictable outputs.
- Prefer small subroutines over one large script body.
- Keep long-running operations observable.
- Validate by actually running the script.

## Implementation workflow

1. Inspect nearby scripts in the same directory and follow local conventions.
2. Clarify inputs, outputs, and output file naming before editing.
3. Implement with function-based structure:
  - `show_help`
  - `check_requirements`
  - domain functions such as `zip_compiled_tools`
  - `main "$@"`
4. Verify with real execution and artifact checks.

## Required Bash baseline

- Shebang: `#!/usr/bin/env bash`
- Strict mode: `set -euo pipefail`
- Resolve script root from `BASH_SOURCE[0]`
- Quote variable expansions
- Use `local` in functions
- Use `[[ ... ]]` for conditionals

## Artifact and overwrite rules

- Prefer timestamped output names for generated artifacts, for example:
  - `compiled_tools_$(date '+%Y%m%d%H%M%S').zip`
- If fixed output filename is mandatory, explain overwrite behavior and choose a non-destructive update strategy.

## Progress visibility rules

- Prefer to show progress by default for long-running jobs.
- For `zip`, prefer progress-friendly options like `-db -dc` when appropriate.

## CLI behavior requirements

- Support `--help`.
- Check required commands with `command -v`.
- Validate required input paths before processing.
- Write clear errors to stderr and exit non-zero on failure.

## Verification checklist

- `bash -n path/to/script.sh`
- Run the script with representative inputs.
- Run with `--help`.
- Confirm output artifacts and paths.
- Report exactly what changed.

## Anti-patterns

- Monolithic scripts without subroutines.
- Unquoted variables.
- Silent destructive behavior not requested by the user.
- Silent long-running operations that appear stalled.
- Hard-coded absolute paths when script-relative paths are available.

## Minimal template

```bash
#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"

function show_help() {
  cat <<EOF
Usage: ./scripts/ops/example.sh [--help]
EOF
}

function check_requirements() {
  command -v zip >/dev/null 2>&1 || {
    echo "zip command not found" >&2
    exit 1
  }
}

function run_main_logic() {
  :
}

function main() {
  if [[ "${1:-}" == "--help" ]]; then
    show_help
    exit 0
  fi

  check_requirements
  run_main_logic
}

main "$@"
```
