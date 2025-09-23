# Repository Guidelines

## Project Structure & Module Organization
Top-level layout centers around `cmd/`, `internal/`, `pkg/`, `scripts/`, and `util/`. Each CLI tool lives in `cmd/cli/<tool>` with matching domain and usecase packages under `internal/<tool>/...`, following the Clean Architecture layering documented in `.clinerules`. Deployment artifacts land in `pkg/bin/<tool>/<platform>/`, while shared assets and sample fixtures stay in `assets/` and `sample_data/`. Review `docs/service_implementation_status.md` before extending a tool to avoid duplicated work.

## Build, Test, and Development Commands
Use `./scripts/create_project_files.sh <package>` to scaffold new tool packages. `./scripts/build.sh` orchestrates every tool-specific `build_*.sh` script (excluding MCP) and fills `pkg/bin/`. Run `./scripts/build_mcp_tools.sh` for MCP-only binaries. Service entry points run via `go run ./cmd/http/main.go` (REST) and `go run ./cmd/grpc/main.go` (gRPC). When iterating on a single CLI, call its dedicated script, e.g. `./scripts/build_image_converter.sh`.

## Coding Style & Naming Conventions
Target Go 1.23 and format with `go fmt ./...` (or `goimports`) before committing. Keep command packages thin; business logic belongs in `internal/*/usecases`, with `domain` models and `interfaces` adapters separated. Directory names and scripts should stay hyphenated (`cmd/cli/git-diff-recorder`, `scripts/build_git_diff_recorder.sh`). Exported Go identifiers use PascalCase, private helpers lowerCamelCase. Shell scripts remain bash-compatible and include brief comments only for non-obvious flows.

## Testing Guidelines
Unit tests rely on standard tooling: `go test -v ./... -coverpkg=./... -covermode=count -coverprofile=coverage.out`. Inspect coverage with `go tool cover -html=coverage.out -o coverage.html` and keep the generated badge in sync. Scope tool-specific checks with commands such as `go test ./internal/git_diff_recorder/...`. Aim for ≥90% coverage (see `.clinerules`). Multimedia and OCR paths depend on FFmpeg and Tesseract, mirroring the CI matrix.

## Commit & Pull Request Guidelines
Follow the Conventional Commit pattern seen in history (`feat:`, `fix:`, `refactor:`, `test:`), keeping subjects imperative and ≤72 characters. Include test evidence or updated artifacts (`coverage.out`, binaries) alongside code changes. PRs should summarize the change, list validation commands, link issues, and attach screenshots or sample outputs for I/O tools. Ensure coverage files and README badges reflect the latest run; the `pull-request-stats` workflow will publish diff metrics, so keep submissions focused.

## Environment & Dependencies
Align with Go 1.23.x (`env.yml`, CI) and prefer module-managed dependencies. Install FFmpeg and Tesseract locally when working on conversion or OCR tools to match CI prerequisites. Provide configuration via environment variables, keeping secrets out of source control and using `sample_data/` for redacted fixtures.
