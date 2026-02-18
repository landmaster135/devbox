# devbox
![Go](https://img.shields.io/badge/Go-1.25-%2300ADD8?logo=go)
![Coverage](https://img.shields.io/badge/Coverage-57.7%25-yellow)
![License](https://img.shields.io/badge/license-MIT-blue)

![thumbnail](assets/thumbnail.webp)

Provides utilities for development.

# Usage
- Go 1.25.5 or later

# Development

## Generate Go packages
```bash
./scripts/create_project_files.sh <PACKAGE_NAME>
```

## Generate shell scripts to build
```bash
cd devbox
./scripts/initialize/script_generator_to_build.sh <TOOL_NAME>
```
# Installation
Install by any entrypoints

## CLI & MCP Tools
Build CLI & MCP tools:
```bash
./scripts/build.sh
```

Confirm compilable distributions.
```bash
go tool dist list
```

## Taskfile
Install various pipelines:
```bash
mv pkg/taskfile /path/to/taskfile
mv pkg/bin /path/to/taskfile/pkg/bin
```

## MCP Tools
Build MCP tools:
```bash
./scripts/build_mcp_tools.sh
```

## Git Hooks
Set Git hooks with built binary files:
```bash
cd /path/to/project
$HOME/devbox/pkg/bash/setup-git-pre-commit-hooks.sh
```

## RESTful API
Run RESTful API:
```bash
go run ./cmd/http/main.go
```

## gRPC API
Run gRPC API:
```bash
go run ./cmd/grpc/main.go
```

## Containerize
Build container of RESTful API:
```bash
task build:dev

task build:staging

task build:prod
```


# Service Implementing Status
[Here](./docs/service_implementation_status.md)

# License
MIT License
