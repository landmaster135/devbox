# Entrypoints
This project provides various entrypoints.

## Installation
Install by any entrypoints

### CLI & MCP Tools
Build CLI & MCP tools:
```bash
./scripts/build.sh
```

Confirm compilable distributions.
```bash
go tool dist list
```

### Taskfile
Install various pipelines:
```bash
mv pkg/taskfile /path/to/taskfile
mv pkg/bin /path/to/taskfile/pkg/bin
```

### MCP Tools
Build MCP tools:
```bash
./scripts/build_mcp_tools.sh
```

### Git Hooks
Set Git hooks with built binary files:
```bash
cd /path/to/project
$HOME/devbox/pkg/bash/setup-git-pre-commit-hooks.sh
```

### RESTful API
Run RESTful API:
```bash
go run ./cmd/http/main.go
```

### gRPC API
Run gRPC API:
```bash
go run ./cmd/grpc/main.go
```

### Containerize
Build container of RESTful API:
```bash
task build:dev

task build:staging

task build:prod
```
