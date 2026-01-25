package workflow

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	usecases "github.com/landmaster135/devbox/internal/cron_workflow/usecases"
	workflowCreator "github.com/landmaster135/devbox/internal/cron_workflow/usecases/workflow_creator"

	textGenerator "github.com/landmaster135/devbox/internal/datetime_calculator/usecases/text_generator"
	discordWebhook "github.com/landmaster135/devbox/internal/discord_webhook/usecases"
	machineInfo "github.com/landmaster135/devbox/internal/machine_info/usecases"
	postgres "github.com/landmaster135/devbox/internal/postgresql/usecases"
	weatherNotificator "github.com/landmaster135/devbox/internal/weather_notificator/usecases"
)

type WorkflowHandler struct {
	Creator *workflowCreator.WorkflowCreator
}

func (wh *WorkflowHandler) GetCreator() *workflowCreator.WorkflowCreator {
	return wh.Creator
}

func NewWorkflowHandler(creator *workflowCreator.WorkflowCreator) *WorkflowHandler {
	return &WorkflowHandler{
		Creator: creator,
	}
}

type WorkflowHandlerRepository interface {
	GetCreator() *workflowCreator.WorkflowCreator
	KeepHeartbeat(ctx context.Context) error
	RetrievePCInfo(ctx context.Context) error
	NotifyWeather(ctx context.Context) error
	NotifyDailyHeading(ctx context.Context) error
	DumpPostgreSQLNotification(ctx context.Context) error
}

// List returns all configured workflows.
func List() ([]usecases.Workflow, error) {
	const tz = "Asia/Tokyo"

	wc, err := workflowCreator.NewWorkflowCreatorDefault(tz)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize WorkflowCreator: %w", err)
	}
	wh := NewWorkflowHandler(wc)

	heartbeatWorkflow := usecases.NewWorkflow(
		"Heartbeat monitor",
		"*/1 * * * *",
		wh.GetCreator().Timezone,
		wh.KeepHeartbeat,
	)
	weatherWorkflow := usecases.NewWorkflow(
		"Daily Tokyo weather notification",
		"0 1 * * 0-6",
		wh.GetCreator().Timezone,
		wh.NotifyWeather,
	)
	dailyHeadingWorkflow := usecases.NewWorkflow(
		"Daily heading Discord notification",
		"1 0 * * 0-6",
		wh.GetCreator().Timezone,
		wh.NotifyDailyHeading,
	)
	postgresDumpWorkflow := usecases.NewWorkflow(
		"Daily PostgreSQL dump with notification",
		"0 2 * * 0-6",
		wh.GetCreator().Timezone,
		wh.DumpPostgreSQLNotification,
	)
	pcInfoWorkflow := usecases.NewWorkflow(
		"Ubuntu PC info snapshot",
		"*/10 * * * 0-6",
		wh.GetCreator().Timezone,
		wh.RetrievePCInfo,
	)

	return []usecases.Workflow{
		*heartbeatWorkflow,
		*weatherWorkflow,
		*dailyHeadingWorkflow,
		*postgresDumpWorkflow,
		*pcInfoWorkflow,
	}, nil
}

func (wh *WorkflowHandler) KeepHeartbeat(ctx context.Context) error {
	heartOwner, err := getEnvVars(wh.GetCreator().EnvRepo, EnvKeyHeartOwner)
	if err != nil {
		return fmt.Errorf("resolve heart owner from %s: %w", "HEART_OWNER", err)
	}

	now, err := wh.GetCreator().TimeRepo.Now(wh.GetCreator().Timezone)
	if err != nil {
		return fmt.Errorf("resolve current time: %w", err)
	}
	timestamp := now.Format("20060102150405")
	statusFile := filepath.Join(wh.GetCreator().VolumeDir, fmt.Sprintf("heartbeat-%s.status", timestamp))

	message := fmt.Sprintf("[heartbeat] alive: %s (owner=%s)", time.Now().Format(time.RFC3339), heartOwner)
	log.Printf("%s", message)
	log.Printf("[heartbeat] writing status file: %s", statusFile)
	if err := wh.GetCreator().FileRepo.Write(statusFile, true, message+"\n"); err != nil {
		return fmt.Errorf("write heartbeat status file: %w", err)
	}
	return nil
}

func (wh *WorkflowHandler) RetrievePCInfo(ctx context.Context) error {
	const (
		networkInterface = "eth0"
	)

	outDirEnv, err := getEnvVars(wh.GetCreator().EnvRepo, EnvKeyPCInfoOutputDirectory)
	if err != nil {
		return fmt.Errorf("resolve PC info output directory: %w", err)
	}
	trimmedOutDir := strings.TrimSpace(outDirEnv)
	if trimmedOutDir == "" {
		return fmt.Errorf("PC info output directory is empty (env=%s)", EnvKeyPCInfoOutputDirectory)
	}
	service := machineInfo.NewMachineInfoService()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	outputDir := filepath.Join(wh.GetCreator().VolumeDir, trimmedOutDir)
	if err := wh.GetCreator().FileRepo.EnsureDir(outputDir); err != nil {
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
}

func (wh *WorkflowHandler) NotifyWeather(ctx context.Context) error {
	const (
		city    = "Tokyo"
		maxDays = 3
	)

	webhookURL, err := getEnvVars(wh.GetCreator().EnvRepo, EnvKeyDiscordWebhookURLForWeather)
	if err != nil {
		return fmt.Errorf("resolve Discord webhook URL: %w", err)
	}
	apiKey, err := getEnvVars(wh.GetCreator().EnvRepo, EnvKeyOpenWeatherAPIKey)
	if err != nil {
		return fmt.Errorf("resolve OpenWeather API key: %w", err)
	}
	service := weatherNotificator.NewWeatherNotificatorService()

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
}

func (wh *WorkflowHandler) NotifyDailyHeading(ctx context.Context) error {
	const dayOffset = 0

	webhookURL, err := getEnvVars(wh.GetCreator().EnvRepo, EnvKeyDiscordWebhookURLForDailyTemplate)
	if err != nil {
		return fmt.Errorf("resolve daily heading Discord webhook URL: %w", err)
	}

	service := discordWebhook.NewDefaultDiscordWebhookService()

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
}

func (wh *WorkflowHandler) DumpPostgreSQLNotification(ctx context.Context) error {
	const (
		format         = "sql"
		notification   = "PostgreSQLのダンプが完了しました"
		embedType      = "postgres"
		embedText      = "最新バックアップ"
		workerParallel = 3
	)

	creator := wh.GetCreator()
	webhookURL, err := getEnvVars(creator.EnvRepo, EnvKeyDiscordWebhookURLForDailyTemplate)
	if err != nil {
		return fmt.Errorf("resolve Discord webhook URL for PostgreSQL dump: %w", err)
	}
	stagingDBURL, err := getEnvVars(creator.EnvRepo, EnvKeyDBURL01Staging)
	if err != nil {
		return fmt.Errorf("resolve staging DB URL: %w", err)
	}
	stagingDirEnv, err := getEnvVars(creator.EnvRepo, EnvKeyDBDirectory01Staging)
	if err != nil {
		return fmt.Errorf("resolve staging dump directory: %w", err)
	}
	productDBURL, err := getEnvVars(creator.EnvRepo, EnvKeyDBURL01Product)
	if err != nil {
		return fmt.Errorf("resolve production DB URL: %w", err)
	}
	productDirEnv, err := getEnvVars(creator.EnvRepo, EnvKeyDBDirectory01Product)
	if err != nil {
		return fmt.Errorf("resolve production dump directory: %w", err)
	}

	stagingOutputDir := filepath.Join(creator.VolumeDir, stagingDirEnv)
	productOutputDir := filepath.Join(creator.VolumeDir, productDirEnv)

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

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	var dumpSummaries []string
	for _, target := range targets {
		if err := creator.FileRepo.EnsureDir(target.outputDir); err != nil {
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
}
