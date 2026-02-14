---
name: go-naming-convention
description: Review and refactor Go naming conventions while preserving behavior. Use when the user asks to fix identifier casing, exported/unexported naming, acronyms, interface names, package/file naming, or naming-only refactors in Go code.
---

# Go Naming Convention Refactoring Skill

Refactor names only. Keep processing logic intact.

## Workflow

1. Confirm target files or package scope. Ask for targets if missing.
2. Inspect identifiers, package names, directory names, and file names.
3. Build a rename plan first, then apply renames consistently across all references.
4. Keep behavior unchanged. Do not mix anti-pattern fixes or logic rewrites unless explicitly requested.
5. Run formatting and tests for affected scope after renaming.

## Naming Rules

Apply `PascalCase` to:
- Public structs
- Public interfaces
- Public methods
- Public functions
- Public variables
- Public constants
- Public fields
- Any identifier that must be accessible outside the package

Apply `camelCase` to:
- Private methods
- Private functions
- Private variables
- Private constants
- Private fields
- Local variables

Use short names for:
- Loop counters (`i`, `j`, `k`)
- Temporary variables with short scope
- Method receivers (typically 1-2 characters)
- Error variables (`err`)

Use lowercase for:
- Package names (single lowercase word; use kebab-case only if one word is not practical)
- Directory names (single lowercase word; use kebab-case only if one word is not practical)
- File names (lowercase with underscores)

## Special Cases

- Keep acronyms uppercase when public (`HTTPClient`, `URLParser`).
- Keep acronyms lowercase when private (`httpClient`, `urlParser`).
- Name single-method interfaces with `method + er` style.
- Ensure Go test files end with `_test.go`.

## Refactoring Guardrails

- Rename definitions and all usages together.
- Update cross-file and cross-package references.
- Avoid introducing API breaks unless the user asked for it.
- If API changes are required, report them clearly before finalizing.

## Validation Checklist

- Run `gofmt` on changed files.
- Run `go test` for affected packages.
- Confirm no compile errors remain after renaming.
- Provide a concise rename summary (`old -> new`) in the final report.
