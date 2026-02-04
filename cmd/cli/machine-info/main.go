package main

import (
	"fmt"
	"os"

	"github.com/landmaster135/devbox/internal/machine_info/config"
	"github.com/landmaster135/devbox/internal/machine_info/usecases"
)

func main() {
	cfg, err := config.ParseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	if cfg.Help {
		config.PrintUsage()
		return
	}

	service := usecases.NewMachineInfoService()

	switch cfg.Operation {
	case config.OperationUbuntu:
		if err := runUbuntuOperation(cfg, service); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応のoperationです: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}

func runUbuntuOperation(cfg *config.Config, service *usecases.MachineInfoService) error {
	fmt.Println("マシン情報を取得中...")

	result, err := service.CollectUbuntuInfo(cfg.NetworkInterface, cfg.MemoryManufacturers, cfg.MemoryNames, "")
	if err != nil {
		return err
	}

	info := result.Info

	for _, warning := range result.Warnings {
		fmt.Printf("警告: %s\n", warning)
	}

	fmt.Printf("CPU名: %s\n", info.CPUName)
	fmt.Printf("CPUコア数: %d\n", info.CPUCores)
	fmt.Printf("論理プロセッサ数: %d\n", info.CPULogicalProcessors)
	fmt.Printf("CPU最大クロック速度: %.2f MHz\n", info.CPUMaxClockSpeedMHz)
	fmt.Printf("CPU現在のクロック速度: %.2f MHz\n", info.CPUCurrentClockSpeedMHz)
	fmt.Printf("CPU温度: %.2f °C\n", info.CPUTemperature)
	fmt.Printf("メモリ総容量: %d MB\n", info.MemoryTotalMB)
	fmt.Printf("メモリ使用量: %d MB\n", info.MemoryUsageMB)
	fmt.Printf("ホスト名: %s\n", info.PCHostname)

	fmt.Println("ネットワーク速度の計測結果:")
	fmt.Printf("平均送信速度: %.2f Kbps\n", info.EthernetAvgSentKbps)
	fmt.Printf("平均受信速度: %.2f Kbps\n", info.EthernetAvgReceivedKbps)

	jsonText, outputPath, err := service.SaveMachineInfoLog(info, cfg.OutputDir, "")
	if err != nil {
		return err
	}

	fmt.Println("\n取得したシステム情報JSON:")
	fmt.Println(jsonText)
	fmt.Printf("\nログファイルに保存しました: %s\n", outputPath)
	return nil
}
