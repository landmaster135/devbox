package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	config "github.com/landmaster135/devbox/internal/image_filterer_v2/config"
	usecases "github.com/landmaster135/devbox/internal/image_filterer_v2/usecases"
)

func main() {
	supportedModes := []string{
		string(config.FilterModeGrayscale),
		string(config.FilterModeColorize),
		string(config.FilterModeVignette),
	}

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	inputPath := flags.String("input", "", "入力画像ファイルのパス")
	outputPath := flags.String("output", "", "出力画像ファイルのパス (省略時は自動生成)")
	outputFormat := flags.String("format", "", "出力フォーマット (png|jpg)。outputの拡張子が無い場合に使用")
	filterMode := flags.String("mode", string(config.FilterModeGrayscale), fmt.Sprintf("適用するフィルター (%s)", strings.Join(supportedModes, ", ")))
	strength := flags.Float64("strength", 1.0, "フィルター強度 (0.0-1.0)")
	tintHex := flags.String("tint", "#ffffff", "colorizeモード用のティントカラー (例: #ff8800)")

	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "画像フィルターをEbitengineシェーダーで適用するCLI\n\n")
		fmt.Fprintf(flags.Output(), "使用方法:\n  %s [オプション]\n\n", os.Args[0])
		fmt.Fprintf(flags.Output(), "オプション:\n")
		flags.PrintDefaults()
		fmt.Fprintf(flags.Output(), "\nmode=colorize の場合は tint を、mode=vignette の場合は strength を調整してください。\n")
	}

	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalf("引数の解析に失敗しました: %v", err)
	}

	cfg := config.Config{
		InputPath:    *inputPath,
		OutputPath:   *outputPath,
		OutputFormat: *outputFormat,
		Mode:         config.FilterMode(*filterMode),
		Strength:     *strength,
		TintHex:      *tintHex,
	}

	svc, err := usecases.NewService(cfg)
	if err != nil {
		log.Fatalf("設定エラー: %v", err)
	}

	output, err := svc.Process()
	if err != nil {
		log.Fatalf("フィルター適用中にエラーが発生しました: %v", err)
	}

	fmt.Println(output)
}
