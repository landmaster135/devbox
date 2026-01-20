# Cron Workflow

Command line utility that keeps a predefined set of background jobs running with [github.com/go-co-op/gocron/v2](https://github.com/go-co-op/gocron). The binary launches a single scheduler, feeds it with the workflows declared in `workflow/core.go`, and blocks until the process receives `SIGINT` or `SIGTERM`.

## Workflows included

| Description | Frequency | Timezone | Purpose |
|-------------|-----------|----------|---------|
| Heartbeat monitor (every minute) | `*/1 * * * *` | Asia/Tokyo (default) | Emits a short log entry so external monitors can confirm the service is alive. |
| Hourly state snapshot | `0 * * * *` | UTC | Simulates a longer-running aggregation job, demonstrating how to respect cancellation through `context.Context`. |

Add or update workflows inside `workflow/core.go` and recompile the CLI to apply the changes.

## Build

```bash
# From the repository root
go build -o bin/cron-workflow ./cmd/cli/cron-workflow
```

## Run

```bash
# Runs the scheduler until Ctrl+C
go run ./cmd/cli/cron-workflow
```

No CLI flags are available other than `-h` / `--help`.

## Sample output

```
$ go run ./cmd/cli/cron-workflow
2026/01/20 10:00:00.123456 registered workflow "Heartbeat monitor (every minute)" (cron=CRON_TZ=Asia/Tokyo */1 * * * *)
2026/01/20 10:00:00.123789 registered workflow "Hourly state snapshot" (cron=CRON_TZ=UTC 0 * * * *)
2026/01/20 10:00:00.123900 scheduler started (2 workflow(s)). waiting for termination signal...
2026/01/20 10:01:00.125678 workflow "Heartbeat monitor (every minute)" completed
2026/01/20 10:02:00.130201 workflow "Heartbeat monitor (every minute)" completed
^C
2026/01/20 10:02:10.000001 signal received: interrupt. shutting down...
2026/01/20 10:02:10.000400 scheduler stopped cleanly
```

The scheduler keeps running indefinitely; terminate it with Ctrl+C or send a termination signal from your supervisor.
