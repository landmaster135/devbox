package main

import (
	"errors"

	httpServer "github.com/landmaster135/devbox/cmd/http/server"

	workflow "github.com/landmaster135/devbox/cmd/cli/cron-workflow/workflow"
	schedulerService "github.com/landmaster135/devbox/internal/cron_workflow/usecases/scheduler_service"
	logging "github.com/landmaster135/devbox/internal/logging"
)

func runCronWorkflow(logger *logging.StructuredLogger) error {
	cronLogger := logging.Ensure(logger)
	workflows, err := workflow.List(cronLogger)
	if err != nil {
		return err
	}
	if len(workflows) == 0 {
		return errors.New("no workflows configured")
	}

	if err := schedulerService.Schedule(cronLogger, workflows); err != nil {
		return err
	}

	cronLogger.WithTags("scheduler").Infof("scheduler stopped cleanly")
	return nil
}

func main() {
	baseLogger := logging.New()
	httpLogger := baseLogger.WithTags("HTTP server")
	cronLogger := baseLogger.WithTags("CRON workflow")

	tag := "lifecycle"
	httpLogger.WithTags(tag).Infof("HTTP REST API サーバーを初期化しています...")

	// サーバーを作成
	server := httpServer.NewHTTPServer(httpLogger)

	// グレースフルシャットダウンのゴルーチンを開始
	go server.GracefulShutdown()
	httpLogger.WithTags(tag).Infof("registered graceful shutdown")

	// cron workflowのゴルーチンを開始
	go func() {
		if err := runCronWorkflow(cronLogger); err != nil {
			cronLogger.WithTags("scheduler").Errorf("cron workflow exited: %v", err)
		}
	}()
	httpLogger.WithTags(tag).Infof("registered cron workflow")

	// サーバーを開始
	server.Start()
}
