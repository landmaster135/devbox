package workflow

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	usecases "github.com/landmaster135/devbox/internal/cron_workflow/usecases"
	workflowCreator "github.com/landmaster135/devbox/internal/cron_workflow/usecases/workflow_creator"
	logging "github.com/landmaster135/devbox/internal/logging"

	textGenerator "github.com/landmaster135/devbox/internal/datetime_calculator/usecases/text_generator"
	discordWebhook "github.com/landmaster135/devbox/internal/discord_webhook/usecases"
	machineInfo "github.com/landmaster135/devbox/internal/machine_info/usecases"
	postgres "github.com/landmaster135/devbox/internal/postgresql/usecases"
	weatherNotificator "github.com/landmaster135/devbox/internal/weather_notificator/usecases"
)

const postgresAttachmentsTable = "public.attachments"

type WorkflowHandler struct {
	Creator *workflowCreator.WorkflowCreator
	logger  *logging.StructuredLogger
}

func (wh *WorkflowHandler) GetCreator() *workflowCreator.WorkflowCreator {
	return wh.Creator
}

func NewWorkflowHandler(creator *workflowCreator.WorkflowCreator, logger *logging.StructuredLogger) *WorkflowHandler {
	return &WorkflowHandler{
		Creator: creator,
		logger:  logging.Ensure(logger),
	}
}

type WorkflowHandlerRepository interface {
	GetCreator() *workflowCreator.WorkflowCreator
	KeepHeartbeat(ctx context.Context) error
	RetrievePCInfo(ctx context.Context) error
	NotifyWeather(ctx context.Context) error
	NotifyDailyHeading(ctx context.Context) error
	DumpPostgreSQLNotification(ctx context.Context) error
	DumpPostgreSQLNotificationForMemos(ctx context.Context) error
}

type postgresDumpTarget struct {
	name               string
	dbURL              string
	outputDir          string
	excludeTableData   []string
	extraSQLDumpTables []string
}

func newPostgresDumpTarget(name, dbURL, outputDir string, splitAttachments bool) postgresDumpTarget {
	target := postgresDumpTarget{
		name:      name,
		dbURL:     dbURL,
		outputDir: outputDir,
	}
	if splitAttachments {
		target.excludeTableData = []string{postgresAttachmentsTable}
		target.extraSQLDumpTables = []string{postgresAttachmentsTable}
	}
	return target
}

// List returns all configured workflows.
func List(logger *logging.StructuredLogger) ([]usecases.Workflow, error) {
	const tz = "Asia/Tokyo"

	wc, err := workflowCreator.NewWorkflowCreatorDefault(tz)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize WorkflowCreator: %w", err)
	}
	wh := NewWorkflowHandler(wc, logging.Ensure(logger))

	heartbeatWorkflow := usecases.NewWorkflow(
		"Heartbeat monitor",
		"0 0 * * 0-6",
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
		"50 0 * * 0-6",
		wh.GetCreator().Timezone,
		wh.DumpPostgreSQLNotification,
	)
	postgresDumpWorkflowForMemos := usecases.NewWorkflow(
		"Daily PostgreSQL dump for memos with notification",
		"5 2 * * 0-6",
		wh.GetCreator().Timezone,
		wh.DumpPostgreSQLNotificationForMemos,
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
		*postgresDumpWorkflowForMemos,
		*pcInfoWorkflow,
	}, nil
}

func (wh *WorkflowHandler) KeepHeartbeat(ctx context.Context) error {
	creator := wh.GetCreator()
	heartOwner, err := getEnvVars(creator.EnvRepo, EnvKeyHeartOwner)
	if err != nil {
		return fmt.Errorf("resolve heart owner from %s: %w", "HEART_OWNER", err)
	}

	now, err := creator.TimeRepo.Now(creator.Timezone)
	if err != nil {
		return fmt.Errorf("resolve current time: %w", err)
	}
	timestamp := now.Format("20060102150405")
	statusFile := filepath.Join(creator.VolumeDir, fmt.Sprintf("heartbeat-%s.status", timestamp))

	heartbeatLogger := wh.logger.WithTags("heartbeat")
	message := fmt.Sprintf("alive: %s (owner=%s)", time.Now().Format(time.RFC3339), heartOwner)
	heartbeatLogger.Infof("%s", message)
	heartbeatLogger.Infof("writing status file: %s", statusFile)
	if err := creator.FileRepo.Write(statusFile, true, message+"\n"); err != nil {
		return fmt.Errorf("write heartbeat status file: %w", err)
	}
	return nil
}

func (wh *WorkflowHandler) RetrievePCInfo(ctx context.Context) error {
	const (
		networkInterface = "eth0"
	)

	creator := wh.GetCreator()
	outDirEnv, err := getEnvVars(creator.EnvRepo, EnvKeyPCInfoOutputDirectory)
	if err != nil {
		return fmt.Errorf("resolve PC info output directory: %w", err)
	}
	trimmedOutDir := strings.TrimSpace(outDirEnv)
	if trimmedOutDir == "" {
		return fmt.Errorf("PC info output directory is empty (env=%s)", EnvKeyPCInfoOutputDirectory)
	}
	memoryNames, err := getOptionalEnvVars(creator.EnvRepo, EnvKeyPCInfoMemoryNamesOfNAS01)
	if err != nil {
		return fmt.Errorf("resolve PC memory names: %w", err)
	}
	memoryManufacturers, err := getOptionalEnvVars(creator.EnvRepo, EnvKeyPCInfoMemoryManufacturersOfNAS01)
	if err != nil {
		return fmt.Errorf("resolve PC memory manufacturers: %w", err)
	}
	hostnameOverride, err := getOptionalEnvVars(creator.EnvRepo, EnvKeyPCInfoHostnameOfNAS01)
	if err != nil {
		return fmt.Errorf("resolve PC hostname override: %w", err)
	}
	service := machineInfo.NewMachineInfoService()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	outputDir := filepath.Join(creator.VolumeDir, trimmedOutDir)
	if err := creator.FileRepo.EnsureDir(outputDir); err != nil {
		return fmt.Errorf("prepare PC info output directory: %w", err)
	}

	result, _, outputPath, err := service.CollectAndSaveUbuntuInfo(
		networkInterface,
		memoryManufacturers,
		memoryNames,
		outputDir,
		strings.TrimSpace(hostnameOverride),
		wh.GetCreator().Timezone,
	)
	if err != nil {
		return fmt.Errorf("collect Ubuntu PC info: %w", err)
	}
	pcLogger := wh.logger.WithTags("pc-info")
	if result != nil {
		for _, warning := range result.Warnings {
			pcLogger.Warnf("warning: %s", warning)
		}
		if result.Info != nil {
			pcLogger.Infof(
				"CPU=%s temp=%.2fC mem_used=%dMB mem_total=%dMB path=%s",
				strings.TrimSpace(result.Info.CPUName),
				result.Info.CPUTemperature,
				result.Info.MemoryUsageMB,
				result.Info.MemoryTotalMB,
				outputPath,
			)
			return nil
		}
	}

	pcLogger.Infof("exported machine info to %s", outputPath)
	return nil
}

func (wh *WorkflowHandler) NotifyWeather(ctx context.Context) error {
	const (
		city    = "Tokyo"
		maxDays = 3
	)

	creator := wh.GetCreator()
	webhookURL, err := getEnvVars(creator.EnvRepo, EnvKeyDiscordWebhookURLForWeather)
	if err != nil {
		return fmt.Errorf("resolve Discord webhook URL: %w", err)
	}
	apiKey, err := getEnvVars(creator.EnvRepo, EnvKeyOpenWeatherAPIKey)
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

	wh.logger.WithTags("weather").Infof("dispatched %s forecast to Discord", city)
	return nil
}

func (wh *WorkflowHandler) NotifyDailyHeading(ctx context.Context) error {
	const dayOffset = 0

	creator := wh.GetCreator()
	webhookURL, err := getEnvVars(creator.EnvRepo, EnvKeyDiscordWebhookURLForDailyTemplate)
	if err != nil {
		return fmt.Errorf("resolve daily heading Discord webhook URL: %w", err)
	}

	service := discordWebhook.NewDefaultDiscordWebhookService()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	content := textGenerator.GenerateDailyHeading(dayOffset, creator.Timezone)
	if err := service.SendNotification(ctx, webhookURL, "テンプレートあゆ", content, "none", "", "", ""); err != nil {
		return fmt.Errorf("send daily heading notification: %w", err)
	}

	wh.logger.WithTags("daily-heading").Infof("dispatched heading content to Discord")
	return nil
}

func (wh *WorkflowHandler) DumpPostgreSQLNotification(ctx context.Context) error {
	creator := wh.GetCreator()
	notification := "PostgreSQLのダンプが完了しました"
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
	targets := []postgresDumpTarget{
		newPostgresDumpTarget("staging", stagingDBURL, stagingOutputDir, true),
		newPostgresDumpTarget("production", productDBURL, productOutputDir, true),
	}

	return wh.dumpPostgreSQLAndNotify(ctx, notification, targets)
}

func (wh *WorkflowHandler) DumpPostgreSQLNotificationForMemos(ctx context.Context) error {
	creator := wh.GetCreator()
	notification := "Memos PostgreSQLのダンプが完了しました"
	memosStagingDBURL, err := getEnvVars(creator.EnvRepo, EnvKeyDBURL01MemosStaging)
	if err != nil {
		return fmt.Errorf("resolve memos staging DB URL: %w", err)
	}
	memosStagingDirEnv, err := getEnvVars(creator.EnvRepo, EnvKeyDBDirectory01MemosStaging)
	if err != nil {
		return fmt.Errorf("resolve memos staging dump directory: %w", err)
	}
	memosProdDBURL, err := getEnvVars(creator.EnvRepo, EnvKeyDBURL01MemosProd)
	if err != nil {
		return fmt.Errorf("resolve memos production DB URL: %w", err)
	}
	memosProdDirEnv, err := getEnvVars(creator.EnvRepo, EnvKeyDBDirectory01MemosProd)
	if err != nil {
		return fmt.Errorf("resolve memos production dump directory: %w", err)
	}

	memosStagingOutputDir := filepath.Join(creator.VolumeDir, memosStagingDirEnv)
	memosProdOutputDir := filepath.Join(creator.VolumeDir, memosProdDirEnv)
	targets := []postgresDumpTarget{
		newPostgresDumpTarget("memos-staging", memosStagingDBURL, memosStagingOutputDir, false),
		newPostgresDumpTarget("memos-prod", memosProdDBURL, memosProdOutputDir, false),
	}

	return wh.dumpPostgreSQLAndNotify(ctx, notification, targets)
}

func (wh *WorkflowHandler) dumpPostgreSQLAndNotify(ctx context.Context, notification string, targets []postgresDumpTarget) error {
	const (
		format         = "binary"
		embedType      = "postgres"
		embedText      = "最新バックアップ"
		workerParallel = 3
	)

	creator := wh.GetCreator()
	webhookURL, err := getEnvVars(creator.EnvRepo, EnvKeyDiscordWebhookURLForDailyTemplate)
	if err != nil {
		return fmt.Errorf("resolve Discord webhook URL for PostgreSQL dump: %w", err)
	}
	service := discordWebhook.NewDefaultDiscordWebhookService()
	concurrency := workerParallel

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

		_, minResult, err := postgres.HandleToDumpAllTables(ctx, target.dbURL, creator.Timezone, target.outputDir, format, nil, &concurrency, "markdown", target.name, target.excludeTableData)
		if err != nil {
			return fmt.Errorf("dump %s database: %w", target.name, err)
		}

		wh.logger.WithTags("postgres-dump").Infof("completed %s dump into %s", target.name, target.outputDir)
		dumpSummaries = append(dumpSummaries, minResult)
		extraSummaries, err := wh.dumpExtraSQLTables(ctx, target)
		if err != nil {
			return err
		}
		dumpSummaries = append(dumpSummaries, extraSummaries...)
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

	wh.logger.WithTags("postgres-dump").Infof("dispatched Discord notification for PostgreSQL backups")
	return nil
}

func (wh *WorkflowHandler) dumpExtraSQLTables(ctx context.Context, target postgresDumpTarget) ([]string, error) {
	var summaries []string
	for _, tableName := range target.extraSQLDumpTables {
		result, err := postgres.HandleToDumpTableDataAsSQL(ctx, target.dbURL, wh.GetCreator().Timezone, tableName, target.outputDir)
		if err != nil {
			return nil, fmt.Errorf("dump %s table as SQL for %s database: %w", tableName, target.name, err)
		}

		wh.logger.WithTags("postgres-dump").Infof("completed %s SQL dump for %s into %s: %s", tableName, target.name, target.outputDir, result)
		summaries = append(summaries, formatExtraSQLDumpSummary(target.name, tableName))
	}

	return summaries, nil
}

func formatExtraSQLDumpSummary(targetName, tableName string) string {
	return fmt.Sprintf(
		"## %s extra SQL dump\n\n```markdown\n| Table | Format | Status |\n| --- | --- | --- |\n| `%s` | sql | completed |\n```",
		targetName,
		tableName,
	)
}
