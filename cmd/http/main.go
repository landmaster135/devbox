package main

import (
	"errors"
	"log"

	httpServer "github.com/landmaster135/devbox/cmd/http/server"

	workflow "github.com/landmaster135/devbox/cmd/cli/cron-workflow/workflow"
	schedulerService "github.com/landmaster135/devbox/internal/cron_workflow/usecases/scheduler_service"
)

func runCronWorkflow() error {
	workflows, err := workflow.List()
	if err != nil {
		return err
	}
	if len(workflows) == 0 {
		return errors.New("no workflows configured")
	}

	if err := schedulerService.Schedule(workflows); err != nil {
		return err
	}

	log.Printf("scheduler stopped cleanly")
	return nil
}

func main() {
	log.Println("HTTP REST API サーバーを初期化しています...")

	// サーバーを作成
	server := httpServer.NewHTTPServer()

	// グレースフルシャットダウンのゴルーチンを開始
	go server.GracefulShutdown()
	log.Printf("registered graceful shutdown")

	// cron workflowのゴルーチンを開始
	go runCronWorkflow()
	log.Printf("registered cron workflow")

	// サーバーを開始
	server.Start()
}
