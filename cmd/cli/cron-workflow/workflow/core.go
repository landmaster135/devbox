package workflow

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	infraEnv "github.com/landmaster135/devbox/internal/cron_workflow/infrastructure/env"
	infraFilesystem "github.com/landmaster135/devbox/internal/cron_workflow/infrastructure/filesystem"
	infraTime "github.com/landmaster135/devbox/internal/cron_workflow/infrastructure/time"
	usecases "github.com/landmaster135/devbox/internal/cron_workflow/usecases"
	textGenerator "github.com/landmaster135/devbox/internal/datetime_calculator/usecases/text_generator"
	discordWebhook "github.com/landmaster135/devbox/internal/discord_webhook/usecases"
	postgres "github.com/landmaster135/devbox/internal/postgresql/usecases"
	weatherNotificator "github.com/landmaster135/devbox/internal/weather_notificator/usecases"
)

// List returns all configured workflows.
func List() ([]usecases.Workflow, error) {
	const tz = "Asia/Tokyo"
	envRepo := infraEnv.NewRepository()
	filesystemRepo := infraFilesystem.NewRepository()
	timeRepo := infraTime.NewRepository()
	wc, err := usecases.NewWorkflowCreator(tz, envRepo, filesystemRepo, timeRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize WorkflowCreator: %w", err)
	}

	heartbeatWorkflow, err := createHeartbeatWorkflow(wc)
	if err != nil {
		return nil, fmt.Errorf("failed to create Heartbeat Workflow: %w", err)
	}
	hourlyHealthSnapshotWorkflow := createHourlyHealthSnapshotWorkflow(wc)
	weatherWorkflow, err := createWeatherNotificationWorkflow(wc)
	if err != nil {
		return nil, fmt.Errorf("failed to create Weather Notification Workflow: %w", err)
	}
	dailyHeadingWorkflow, err := createDailyHeadingNotificationWorkflow(wc)
	if err != nil {
		return nil, fmt.Errorf("failed to create Daily Heading Workflow: %w", err)
	}
	postgresDumpWorkflow, err := createPostgreSQLDumpNotificationWorkflow(wc)
	if err != nil {
		return nil, fmt.Errorf("failed to create PostgreSQL Dump Workflow: %w", err)
	}

	return []usecases.Workflow{
		*heartbeatWorkflow,
		*hourlyHealthSnapshotWorkflow,
		*weatherWorkflow,
		*dailyHeadingWorkflow,
		*postgresDumpWorkflow,
	}, nil
}

func createHeartbeatWorkflow(c *usecases.WorkflowCreator) (*usecases.Workflow, error) {
	heartOwner, err := getEnvVars(c.EnvRepo, EnvKeyHeartOwner)
	if err != nil {
		return nil, fmt.Errorf("resolve heart owner from %s: %w", "HEART_OWNER", err)
	}

	return usecases.NewWorkflow(
		"Heartbeat monitor (every minute)",
		"*/1 * * * *",
		c.Timezone,
		func(ctx context.Context) error {
			now, err := c.TimeRepo.Now(c.Timezone)
			if err != nil {
				return fmt.Errorf("resolve current time: %w", err)
			}
			timestamp := now.Format("20060102150405")
			statusFile := filepath.Join(c.VolumeDir, fmt.Sprintf("heartbeat-%s.status", timestamp))

			message := fmt.Sprintf("[heartbeat] alive: %s (owner=%s)", time.Now().Format(time.RFC3339), heartOwner)
			log.Printf("%s", message)
			log.Printf("[heartbeat] writing status file: %s", statusFile)
			if err := c.FileRepo.Write(statusFile, true, message+"\n"); err != nil {
				return fmt.Errorf("write heartbeat status file: %w", err)
			}
			return nil
		},
	), nil
}

func createHourlyHealthSnapshotWorkflow(c *usecases.WorkflowCreator) *usecases.Workflow {
	return usecases.NewWorkflow(
		"Hourly state snapshot",
		"15 * * * *",
		c.Timezone,
		func(ctx context.Context) error {
			select {
			case <-time.After(2 * time.Second):
				log.Printf("[snapshot] captured at UTC=%s", time.Now().UTC().Format(time.RFC3339))
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	)
}

func createWeatherNotificationWorkflow(c *usecases.WorkflowCreator) (*usecases.Workflow, error) {
	const (
		city    = "Tokyo"
		maxDays = 3
		cronExp = "0 1 * * 0-6"
	)

	webhookURL, err := getEnvVars(c.EnvRepo, EnvKeyDiscordWebhookURLForWeather)
	if err != nil {
		return nil, fmt.Errorf("resolve Discord webhook URL: %w", err)
	}
	apiKey, err := getEnvVars(c.EnvRepo, EnvKeyOpenWeatherAPIKey)
	if err != nil {
		return nil, fmt.Errorf("resolve OpenWeather API key: %w", err)
	}
	service := weatherNotificator.NewWeatherNotificatorService()

	return usecases.NewWorkflow(
		"Daily Tokyo weather notification",
		cronExp,
		c.Timezone,
		func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if err := service.HandleWeatherNotification(apiKey, city, maxDays, webhookURL); err != nil {
				return fmt.Errorf("send weather notification: %w", err)
			}

			log.Printf("[weather] dispatched %s forecast to Discord", city)
			return nil
		},
	), nil
}

func createDailyHeadingNotificationWorkflow(c *usecases.WorkflowCreator) (*usecases.Workflow, error) {
	const (
		cronExp   = "1 0 * * 0-6"
		dayOffset = 0
	)

	webhookURL, err := getEnvVars(c.EnvRepo, EnvKeyDiscordWebhookURLForDailyTemplate)
	if err != nil {
		return nil, fmt.Errorf("resolve daily heading Discord webhook URL: %w", err)
	}

	service := discordWebhook.NewDefaultDiscordWebhookService()

	return usecases.NewWorkflow(
		"Daily heading Discord notification",
		cronExp,
		c.Timezone,
		func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			content := textGenerator.GenerateDailyHeading(dayOffset)
			if err := service.SendNotification(ctx, webhookURL, "テンプレートあゆ", content, "none", "", "", ""); err != nil {
				return fmt.Errorf("send daily heading notification: %w", err)
			}

			log.Printf("[daily-heading] dispatched heading content to Discord")
			return nil
		},
	), nil
}

func createPostgreSQLDumpNotificationWorkflow(c *usecases.WorkflowCreator) (*usecases.Workflow, error) {
	const (
		cronExp        = "0 2 * * *"
		workflowName   = "Daily PostgreSQL dump notification"
		format         = "sql"
		notification   = "PostgreSQLのダンプが完了しました"
		embedType      = "postgres"
		embedText      = "最新バックアップ"
		workerParallel = 3
	)

	webhookURL, err := getEnvVars(c.EnvRepo, EnvKeyDiscordWebhookURLForDailyTemplate)
	if err != nil {
		return nil, fmt.Errorf("resolve Discord webhook URL for PostgreSQL dump: %w", err)
	}
	stagingDBURL, err := getEnvVars(c.EnvRepo, EnvKeyDBURL01Staging)
	if err != nil {
		return nil, fmt.Errorf("resolve staging DB URL: %w", err)
	}
	stagingDirEnv, err := getEnvVars(c.EnvRepo, EnvKeyDBDirectory01Staging)
	if err != nil {
		return nil, fmt.Errorf("resolve staging dump directory: %w", err)
	}
	productDBURL, err := getEnvVars(c.EnvRepo, EnvKeyDBURL01Product)
	if err != nil {
		return nil, fmt.Errorf("resolve production DB URL: %w", err)
	}
	productDirEnv, err := getEnvVars(c.EnvRepo, EnvKeyDBDirectory01Product)
	if err != nil {
		return nil, fmt.Errorf("resolve production dump directory: %w", err)
	}

	stagingOutputDir := filepath.Join(c.VolumeDir, stagingDirEnv)
	productOutputDir := filepath.Join(c.VolumeDir, productDirEnv)

	service := discordWebhook.NewDefaultDiscordWebhookService()
	concurrency := workerParallel

	targets := []struct {
		name      string
		dbURL     string
		outputDir string
	}{
		{name: "staging", dbURL: stagingDBURL, outputDir: stagingOutputDir},
		{name: "production", dbURL: productDBURL, outputDir: productOutputDir},
	}

	return usecases.NewWorkflow(
		workflowName,
		cronExp,
		c.Timezone,
		func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			var dumpSummaries []string
			for _, target := range targets {
				if err := c.FileRepo.EnsureDir(target.outputDir); err != nil {
					return fmt.Errorf("prepare dump directory for %s: %w", target.name, err)
				}

				result, err := postgres.HandleToDumpAllTables(ctx, target.dbURL, target.outputDir, format, nil, &concurrency)
				if err != nil {
					return fmt.Errorf("dump %s database: %w", target.name, err)
				}

				log.Printf("[postgres-dump] completed %s dump into %s", target.name, target.outputDir)
				dumpSummaries = append(dumpSummaries, fmt.Sprintf("[%s]\n%s", target.name, result))
			}

			content := notification
			if len(dumpSummaries) > 0 {
				content = fmt.Sprintf("%s\n%s", notification, strings.Join(dumpSummaries, "\n\n"))
			}

			if err := service.SendNotification(
				ctx,
				webhookURL,
				"",
				content,
				embedType,
				embedText,
				"",
				"",
			); err != nil {
				return fmt.Errorf("send PostgreSQL dump notification: %w", err)
			}

			log.Printf("[postgres-dump] dispatched Discord notification for PostgreSQL backups")
			return nil
		},
	), nil
}
