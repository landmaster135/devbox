# Cron Workflow HTTP Handler

Serves GUI for CRON workflow.

## Prerequisites
- Go 1.25.5 or later
- templ v0.3.977 or later

## Build
Run HTTP server
```bash
go run ./cmd/http/main.go
```

Generate HTML page from .templ boilerplate files.
e.g.
```bash
go run github.com/a-h/templ/cmd/templ@v0.3.977 generate internal/templ_components/core/style/style.templ
```
