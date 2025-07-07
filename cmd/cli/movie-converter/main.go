package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	usecases "github.com/landmaster135/devbox/internal/movie_converter/usecases"
)

func main() {
	var config usecases.ConversionConfig

	// コマンドライン引数の定義
	flag.StringVar(&config.InputFile, "input", "", "入力ファイルのパス（必須）")
	flag.StringVar(&config.OutputFile, "output", "", "出力ファイルのパス（省略時は自動生成）")
	flag.IntVar(&config.FPS, "fps", 0, "出力のFPS（MP4→GIF: デフォルト60, GIF→MP4: デフォルト15）")
	flag.IntVar(&config.Width, "width", 0, "GIF出力時の幅（0=デフォルト品質）")
	speedStr := flag.String("speed", "0", "GIF出力時の速度倍率（デフォルト2.0）")
	flag.IntVar(&config.Loop, "loop", 0, "GIF出力時のループ設定（0=無限ループ, -1=ループなし）")
	flag.BoolVar(&config.UseItsScale, "use-itsscale", true, "GIF出力時にitsscaleを使用するか")

	help := flag.Bool("help", false, "ヘルプを表示")
	flag.Parse()

	// ヘルプの表示
	if *help {
		fmt.Println("Movie Converter - GIFとMP4の相互変換ツール")
		fmt.Println()
		fmt.Println("使用方法:")
		fmt.Println("  movie-converter -input <入力ファイル> [オプション]")
		fmt.Println()
		fmt.Println("例:")
		fmt.Println("  # MP4からGIF（デフォルト設定）")
		fmt.Println("  movie-converter -input video.mp4")
		fmt.Println()
		fmt.Println("  # GIFからMP4（カスタムFPS）")
		fmt.Println("  movie-converter -input video.gif -fps 24")
		fmt.Println()
		fmt.Println("  # MP4からGIF（カスタム設定）")
		fmt.Println("  movie-converter -input video.mp4 -output custom.gif -fps 30 -width 320 -speed 1.5")
		fmt.Println()
		fmt.Println("サポートされている拡張子:")
		extensions := usecases.GetSupportedExtensions()
		fmt.Printf("  入力: %v\n", extensions["input"])
		fmt.Printf("  出力: %v\n", extensions["output"])
		fmt.Println()
		fmt.Println("オプション:")
		flag.PrintDefaults()
		return
	}

	// 速度の変換
	if *speedStr != "0" {
		speed, err := strconv.ParseFloat(*speedStr, 64)
		if err != nil {
			log.Fatalf("エラー: 速度の値が不正です: %s", *speedStr)
		}
		config.Speed = speed
	}

	// 設定の検証
	if err := usecases.ValidateConfig(config); err != nil {
		log.Fatalf("設定エラー: %v", err)
	}

	// 出力ファイルの自動生成
	if config.OutputFile == "" {
		config.OutputFile = usecases.GenerateOutputFile(config.InputFile)
		log.Printf("出力ファイル名を自動生成しました: %s", config.OutputFile)
	}

	// 変換サービスの作成と実行
	service := usecases.NewMovieConverterService(config)
	if err := service.Convert(); err != nil {
		log.Fatalf("変換エラー: %v", err)
	}

	fmt.Printf("変換が正常に完了しました: %s -> %s\n", config.InputFile, config.OutputFile)
}
