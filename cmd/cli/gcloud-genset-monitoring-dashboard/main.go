package main

import (
	"fmt"
	"log"
	"os"

	"github.com/landmaster135/devbox/internal/gcloud_genset_monitoring_dashboard/config"
	"github.com/landmaster135/devbox/internal/gcloud_genset_monitoring_dashboard/usecases"
)

func main() {
	// コマンドライン引数の解析
	cfg, err := config.ParseFlags()
	if err != nil {
		log.Printf("エラー: %v", err)
		config.PrintUsage()
		os.Exit(1)
	}

	// ヘルプが要求された場合
	if cfg.Help {
		config.PrintUsage()
		return
	}

	// サービスの作成
	service := usecases.NewService(cfg.Project, cfg.Location, cfg.Service, cfg.ServiceAccountID)

	// 操作の実行
	switch cfg.Operation {
	case "create-dashboard-for-cloud-run":
		result, err := service.CreateDashboardForCloudRun()
		if err != nil {
			log.Printf("ダッシュボードの作成に失敗しました: %v", err)
			os.Exit(1)
		}
		fmt.Println(result)
	default:
		log.Printf("未対応の操作です: %s", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}
