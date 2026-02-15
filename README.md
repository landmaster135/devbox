# devbox
![Go](https://img.shields.io/badge/Go-1.25-%2300ADD8?logo=go)
![Coverage](https://img.shields.io/badge/Coverage-57.4%25-yellow)
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
./pkg/bin/cli/linux_amd64/script-generator-to-build <TOOL_NAME>
```

# Build

## Common Tools
```bash
./scripts/build.sh
```

Confirm compilable distributions.
```bash
go tool dist list
```

## MCP Tools
```bash
./scripts/build_mcp_tools.sh
```

## Git Hooks
Set Git hooks with built binary files
```bash
cd /path/to/dir
$HOME/devbox/pkg/bash/setup-git-pre-commit-hooks.sh
```

## RESTful API
```bash
go run ./cmd/http/main.go
```

## gRPC API
```bash
go run ./cmd/grpc/main.go
```

# Service Implementing Status
[Here](./docs/service_implementation_status.md)

# License
MIT License
