package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/service_implementing_viewer/config"
	usecases "github.com/landmaster135/devbox/internal/service_implementing_viewer/usecases"
)

func main() {
	// コマンドライン引数を解析
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	// サービスを作成
	service := usecases.NewServiceImplementingViewerService(cfg.RootDir, cfg.TargetDirs)

	// サービス実装状況を取得
	result, statistics, err := service.GetServiceImplementingStatus()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Print(result)

	// 統計情報を出力
	fmt.Printf("\n## 統計情報\n\n")
	fmt.Printf("- **総サービス数**: %d\n", statistics.TotalServices)
	fmt.Printf("- **CLIツール実装数**: %d\n", statistics.CLICount)
	fmt.Printf("- **MCPツール実装数**: %d\n", statistics.MCPCount)
	fmt.Printf("- **gRPCハンドラ実装数**: %d\n", statistics.GRPCCount)
	fmt.Printf("- **HTTPハンドラ実装数**: %d\n", statistics.HTTPCount)
	fmt.Printf("- **CLIのみ実装**: %d\n", statistics.CLIOnlyCount)
	fmt.Printf("- **MCPのみ実装**: %d\n", statistics.MCPOnlyCount)
	fmt.Printf("- **gRPCハンドラのみ実装**: %d\n", statistics.GRPCOnlyCount)
	fmt.Printf("- **HTTPハンドラのみ実装**: %d\n", statistics.HTTPOnlyCount)
	fmt.Printf("- **CLI+MCP両方実装**: %d\n", statistics.BothCLIMCPCount)
	fmt.Printf("- **全て実装済み**: %d\n", statistics.AllImplementedCount)
}
