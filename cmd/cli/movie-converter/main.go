package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	usecases "github.com/landmaster135/devbox/internal/movie_converter/usecases"
)

func main() {
	// 単一ファイル処理用の設定
	var config usecases.ConversionConfig

	// バッチ処理用の設定
	var batchConfig usecases.BatchConversionConfig

	// 単一ファイル処理用のコマンドライン引数
	flag.StringVar(&config.InputFile, "input", "", "入力ファイルのパス（単一ファイル処理時）")
	flag.StringVar(&config.OutputFile, "output", "", "出力ファイルのパス（単一ファイル処理時、省略時は自動生成）")

	// バッチ処理用のコマンドライン引数
	flag.StringVar(&batchConfig.InputDir, "input-dir", "", "入力ディレクトリのパス（バッチ処理時）")
	flag.StringVar(&batchConfig.InputExt, "input-ext", "", "入力ファイルの拡張子（バッチ処理時、例: mp4）")
	flag.StringVar(&batchConfig.OutputDir, "output-dir", "", "出力ディレクトリのパス（バッチ処理時）")
	flag.StringVar(&batchConfig.OutputExt, "output-ext", "", "出力ファイルの拡張子（バッチ処理時、例: gif）")
	flag.BoolVar(&batchConfig.Recursive, "recursive", false, "サブディレクトリも再帰的に処理するか（バッチ処理時）")

	// 共通のオプション
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
		fmt.Println("  # 単一ファイル処理")
		fmt.Println("  movie-converter -input <入力ファイル> [オプション]")
		fmt.Println()
		fmt.Println("  # バッチ処理")
		fmt.Println("  movie-converter -input-dir <入力ディレクトリ> -input-ext <拡張子> -output-dir <出力ディレクトリ> -output-ext <拡張子> [オプション]")
		fmt.Println()
		fmt.Println("単一ファイル処理の例:")
		fmt.Println("  # MP4からGIF（デフォルト設定）")
		fmt.Println("  movie-converter -input video.mp4")
		fmt.Println()
		fmt.Println("  # GIFからMP4（カスタムFPS）")
		fmt.Println("  movie-converter -input video.gif -fps 24")
		fmt.Println()
		fmt.Println("  # MP4からGIF（カスタム設定）")
		fmt.Println("  movie-converter -input video.mp4 -output custom.gif -fps 30 -width 320 -speed 1.5")
		fmt.Println()
		fmt.Println("バッチ処理の例:")
		fmt.Println("  # MP4ファイルを一括でGIFに変換")
		fmt.Println("  movie-converter -input-dir ./videos -input-ext mp4 -output-dir ./gifs -output-ext gif")
		fmt.Println()
		fmt.Println("  # GIFファイルを一括でMP4に変換")
		fmt.Println("  movie-converter -input-dir ./animations -input-ext gif -output-dir ./videos -output-ext mp4")
		fmt.Println()
		fmt.Println("  # 再帰的にサブディレクトリも処理")
		fmt.Println("  movie-converter -input-dir ./media -input-ext mp4 -output-dir ./converted -output-ext gif -recursive")
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
	var speed float64
	if *speedStr != "0" {
		var err error
		speed, err = strconv.ParseFloat(*speedStr, 64)
		if err != nil {
			log.Fatalf("エラー: 速度の値が不正です: %s", *speedStr)
		}
		config.Speed = speed
	}

	// 処理モードの判定
	isBatchMode := batchConfig.InputDir != "" || batchConfig.InputExt != "" || batchConfig.OutputDir != "" || batchConfig.OutputExt != ""
	isSingleMode := config.InputFile != ""

	if isBatchMode && isSingleMode {
		log.Fatal("エラー: 単一ファイル処理とバッチ処理のオプションを同時に指定することはできません")
	}

	if !isBatchMode && !isSingleMode {
		log.Fatal("エラー: 入力ファイル（-input）または入力ディレクトリ（-input-dir）を指定してください。-help でヘルプを表示します。")
	}

	if isBatchMode {
		// バッチ処理モード
		executeBatchConversion(batchConfig, config.FPS, config.Width, speed, config.Loop, config.UseItsScale)
	} else {
		// 単一ファイル処理モード
		executeSingleConversion(config)
	}
}

// executeSingleConversion executes single file conversion
func executeSingleConversion(config usecases.ConversionConfig) {
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

// executeBatchConversion executes batch conversion
func executeBatchConversion(batchConfig usecases.BatchConversionConfig, fps, width int, speed float64, loop int, useItsScale bool) {
	// 共通オプションをバッチ設定にコピー
	batchConfig.FPS = fps
	batchConfig.Width = width
	batchConfig.Speed = speed
	batchConfig.Loop = loop
	batchConfig.UseItsScale = useItsScale

	// バッチ設定の検証
	if err := usecases.ValidateBatchConfig(&batchConfig); err != nil {
		log.Fatalf("バッチ設定エラー: %v", err)
	}

	// バッチ変換サービスの作成と実行
	batchService := usecases.NewBatchMovieConverterService(batchConfig)
	result, err := batchService.BatchConvert()
	if err != nil {
		log.Fatalf("バッチ変換エラー: %v", err)
	}

	// 結果の表示
	fmt.Printf("\n=== バッチ変換結果 ===\n")
	fmt.Printf("総ファイル数: %d\n", result.TotalFiles)
	fmt.Printf("成功: %d\n", result.SuccessCount)
	fmt.Printf("失敗: %d\n", result.FailureCount)

	if result.FailureCount > 0 {
		fmt.Printf("\n失敗したファイル:\n")
		for _, failedFile := range result.FailedFiles {
			fmt.Printf("  - %s\n", failedFile)
		}

		fmt.Printf("\n詳細なエラー情報:\n")
		for _, convResult := range result.Results {
			if !convResult.Success {
				fmt.Printf("  %s: %v\n", convResult.InputFile, convResult.Error)
			}
		}
	}

	if result.SuccessCount > 0 {
		fmt.Printf("\nバッチ変換が完了しました: %s (%s) -> %s (%s)\n",
			batchConfig.InputDir, batchConfig.InputExt, batchConfig.OutputDir, batchConfig.OutputExt)
	}
}
