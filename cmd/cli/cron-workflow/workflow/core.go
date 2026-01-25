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
	workflowCreator "github.com/landmaster135/devbox/internal/cron_workflow/usecases/workflow_creator"

	textGenerator "github.com/landmaster135/devbox/internal/datetime_calculator/usecases/text_generator"
	discordWebhook "github.com/landmaster135/devbox/internal/discord_webhook/usecases"
	machineInfo "github.com/landmaster135/devbox/internal/machine_info/usecases"
	postgres "github.com/landmaster135/devbox/internal/postgresql/usecases"
	weatherNotificator "github.com/landmaster135/devbox/internal/weather_notificator/usecases"
)

// List returns all configured workflows.
func List() ([]usecases.Workflow, error) {
	const tz = "Asia/Tokyo"
	envRepo := infraEnv.NewRepository()
	filesystemRepo := infraFilesystem.NewRepository()
	timeRepo := infraTime.NewRepository()
	wc, err := workflowCreator.NewWorkflowCreator(tz, envRepo, filesystemRepo, timeRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize WorkflowCreator: %w", err)
	}

	heartbeatWorkflow, err := createHeartbeatWorkflow(wc)
	if err != nil {
		return nil, fmt.Errorf("failed to create Heartbeat Workflow: %w", err)
	}
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
	pcInfoWorkflow, err := createPCInfoWorkflow(wc)
	if err != nil {
		return nil, fmt.Errorf("failed to create PC Info Workflow: %w", err)
	}

	return []usecases.Workflow{
		*heartbeatWorkflow,
		*weatherWorkflow,
		*dailyHeadingWorkflow,
		*postgresDumpWorkflow,
		*pcInfoWorkflow,
	}, nil
}

func createHeartbeatWorkflow(c *workflowCreator.WorkflowCreator) (*usecases.Workflow, error) {
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

func createWeatherNotificationWorkflow(c *workflowCreator.WorkflowCreator) (*usecases.Workflow, error) {
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

func createDailyHeadingNotificationWorkflow(c *workflowCreator.WorkflowCreator) (*usecases.Workflow, error) {
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

func createPostgreSQLDumpNotificationWorkflow(c *workflowCreator.WorkflowCreator) (*usecases.Workflow, error) {
	const (
		cronExp        = "0 2 * * 0-6"
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

	headerGen := func(header string, summaries []string) string {
		return fmt.Sprintf("%s\n%s", header, strings.Join(summaries, "\n"))
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

				_, minResult, err := postgres.HandleToDumpAllTables(ctx, target.dbURL, target.outputDir, format, nil, &concurrency, "markdown", target.name)
				if err != nil {
					return fmt.Errorf("dump %s database: %w", target.name, err)
				}

				log.Printf("[postgres-dump] completed %s dump into %s", target.name, target.outputDir)
				dumpSummaries = append(dumpSummaries, minResult)
			}

			content := notification
			if len(dumpSummaries) > 0 {
				content = headerGen(notification, dumpSummaries)
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

func createPCInfoWorkflow(c *workflowCreator.WorkflowCreator) (*usecases.Workflow, error) {
	const (
		workflowName     = "Ubuntu PC info snapshot"
		cronExp          = "*/10 * * * 0-6"
		networkInterface = "eth0"
	)

	outDirEnv, err := getEnvVars(c.EnvRepo, EnvKeyPCInfoOutputDirectory)
	if err != nil {
		return nil, fmt.Errorf("resolve PC info output directory: %w", err)
	}
	trimmedOutDir := strings.TrimSpace(outDirEnv)
	if trimmedOutDir == "" {
		return nil, fmt.Errorf("PC info output directory is empty (env=%s)", EnvKeyPCInfoOutputDirectory)
	}
	service := machineInfo.NewMachineInfoService()

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

			outputDir := filepath.Join(c.VolumeDir, trimmedOutDir)
			if err := c.FileRepo.EnsureDir(outputDir); err != nil {
				return fmt.Errorf("prepare PC info output directory: %w", err)
			}

			result, _, outputPath, err := service.CollectAndSaveUbuntuInfo(networkInterface, outputDir)
			if err != nil {
				return fmt.Errorf("collect Ubuntu PC info: %w", err)
			}
			if result != nil {
				for _, warning := range result.Warnings {
					log.Printf("[pc-info] warning: %s", warning)
				}
				if result.Info != nil {
					log.Printf(
						"[pc-info] CPU=%s temp=%.2fC mem_used=%.2fMB mem_total=%.2fMB path=%s",
						strings.TrimSpace(result.Info.CPUName),
						result.Info.CPUTemperature,
						result.Info.MemoryUsageMB,
						result.Info.MemoryTotalMB,
						outputPath,
					)
					return nil
				}
			}

			log.Printf("[pc-info] exported machine info to %s", outputPath)
			return nil
		},
	), nil
}
