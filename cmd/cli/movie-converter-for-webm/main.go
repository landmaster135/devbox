package main

import (
	"flag"
	"fmt"
	"log"

	usecases "github.com/landmaster135/devbox/internal/movie_converter_for_webm/usecases"
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
	flag.StringVar(&batchConfig.OutputExt, "output-ext", "", "出力ファイルの拡張子（バッチ処理時、例: webm）")
	flag.BoolVar(&batchConfig.Recursive, "recursive", false, "サブディレクトリも再帰的に処理するか（バッチ処理時）")

	// WEBM変換固有のオプション
	flag.StringVar(&config.VideoBitrate, "video-bitrate", "", "ビデオビットレート（例: 1M, 1500k）")
	flag.StringVar(&config.AudioBitrate, "audio-bitrate", "128k", "オーディオビットレート（デフォルト: 128k）")
	flag.IntVar(&config.CRF, "crf", 0, "Constant Rate Factor（0-63、低いほど高品質、デフォルト: 30）")
	flag.StringVar(&config.AudioCodec, "audio-codec", "opus", "オーディオコーデック（opus/vorbis、デフォルト: opus）")
	flag.StringVar(&config.ConversionMode, "conversion-mode", "crf", "変換モード（crf/cbr、デフォルト: crf）")
	flag.IntVar(&config.VideoQuality, "video-quality", 75, "ビデオ品質（CBRモード時、デフォルト: 75）")

	help := flag.Bool("help", false, "ヘルプを表示")
	flag.Parse()

	// ヘルプの表示
	if *help {
		fmt.Println("Movie Converter for WEBM - WEBMとMP4の相互変換ツール")
		fmt.Println()
		fmt.Println("使用方法:")
		fmt.Println("  # 単一ファイル処理")
		fmt.Println("  movie-converter-for-webm -input <入力ファイル> [オプション]")
		fmt.Println()
		fmt.Println("  # バッチ処理")
		fmt.Println("  movie-converter-for-webm -input-dir <入力ディレクトリ> -input-ext <拡張子> -output-dir <出力ディレクトリ> -output-ext <拡張子> [オプション]")
		fmt.Println()
		fmt.Println("単一ファイル処理の例:")
		fmt.Println("  # MP4からWEBM（デフォルト設定、CRFモード）")
		fmt.Println("  movie-converter-for-webm -input video.mp4")
		fmt.Println()
		fmt.Println("  # WEBMからMP4")
		fmt.Println("  movie-converter-for-webm -input video.webm")
		fmt.Println()
		fmt.Println("  # MP4からWEBM（カスタムCRF設定）")
		fmt.Println("  movie-converter-for-webm -input video.mp4 -output custom.webm -crf 25 -audio-codec vorbis")
		fmt.Println()
		fmt.Println("  # MP4からWEBM（CBRモード）")
		fmt.Println("  movie-converter-for-webm -input video.mp4 -conversion-mode cbr -video-bitrate 2M -video-quality 80")
		fmt.Println()
		fmt.Println("バッチ処理の例:")
		fmt.Println("  # MP4ファイルを一括でWEBMに変換")
		fmt.Println("  movie-converter-for-webm -input-dir ./videos -input-ext mp4 -output-dir ./webm -output-ext webm")
		fmt.Println()
		fmt.Println("  # WEBMファイルを一括でMP4に変換")
		fmt.Println("  movie-converter-for-webm -input-dir ./webm -input-ext webm -output-dir ./videos -output-ext mp4")
		fmt.Println()
		fmt.Println("  # 再帰的にサブディレクトリも処理（CRFモード）")
		fmt.Println("  movie-converter-for-webm -input-dir ./media -input-ext mp4 -output-dir ./converted -output-ext webm -recursive -crf 28")
		fmt.Println()
		fmt.Println("サポートされている拡張子:")
		extensions := usecases.GetSupportedExtensions()
		fmt.Printf("  入力: %v\n", extensions["input"])
		fmt.Printf("  出力: %v\n", extensions["output"])
		fmt.Println()
		fmt.Println("WEBM変換オプション:")
		fmt.Println("  -crf <値>              CRF値（0-63、低いほど高品質、デフォルト: 30）")
		fmt.Println("  -video-bitrate <値>    ビデオビットレート（例: 1M, 1500k）")
		fmt.Println("  -audio-bitrate <値>    オーディオビットレート（デフォルト: 128k）")
		fmt.Println("  -audio-codec <codec>   オーディオコーデック（opus/vorbis、デフォルト: opus）")
		fmt.Println("  -conversion-mode <mode> 変換モード（crf/cbr、デフォルト: crf）")
		fmt.Println("  -video-quality <値>    ビデオ品質（CBRモード時、デフォルト: 75）")
		fmt.Println()
		fmt.Println("変換モードについて:")
		fmt.Println("  CRF (Constant Rate Factor): 品質重視、ファイルサイズは可変")
		fmt.Println("  CBR (Constant Bit Rate): ビットレート固定、ファイルサイズ予測可能")
		fmt.Println()
		fmt.Println("その他のオプション:")
		flag.PrintDefaults()
		return
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

	// 統合されたサービスを使用して処理を実行
	executeUnifiedConversion(config, batchConfig, isBatchMode)
}

// displayConversionResult displays the conversion result
func displayConversionResult(result *usecases.UnifiedConversionResult) {
	if result.Mode == usecases.SingleFileMode {
		// 単一ファイル処理の結果表示
		if result.Success && result.SingleResult != nil {
			fmt.Printf("変換が正常に完了しました: %s -> %s\n",
				result.SingleResult.InputFile, result.SingleResult.OutputFile)
		}
	} else {
		// バッチ処理の結果表示
		if result.BatchResult != nil {
			fmt.Printf("\n=== バッチ変換結果 ===\n")
			fmt.Printf("総ファイル数: %d\n", result.BatchResult.TotalFiles)
			fmt.Printf("成功: %d\n", result.BatchResult.SuccessCount)
			fmt.Printf("失敗: %d\n", result.BatchResult.FailureCount)

			if result.BatchResult.FailureCount > 0 {
				fmt.Printf("\n失敗したファイル:\n")
				for _, failedFile := range result.BatchResult.FailedFiles {
					fmt.Printf("  - %s\n", failedFile)
				}

				fmt.Printf("\n詳細なエラー情報:\n")
				for _, convResult := range result.BatchResult.Results {
					if !convResult.Success {
						fmt.Printf("  %s: %v\n", convResult.InputFile, convResult.Error)
					}
				}
			}

			if result.BatchResult.SuccessCount > 0 {
				fmt.Printf("\nバッチ変換が完了しました\n")
			}
		}
	}
}

// executeUnifiedConversion executes conversion using the unified service
func executeUnifiedConversion(config usecases.ConversionConfig, batchConfig usecases.BatchConversionConfig, isBatchMode bool) {
	var singleConfig *usecases.ConversionConfig
	var batchConfigPtr *usecases.BatchConversionConfig

	if isBatchMode {
		// 共通オプションをバッチ設定にコピー
		batchConfig.VideoBitrate = config.VideoBitrate
		batchConfig.AudioBitrate = config.AudioBitrate
		batchConfig.CRF = config.CRF
		batchConfig.AudioCodec = config.AudioCodec
		batchConfig.ConversionMode = config.ConversionMode
		batchConfig.VideoQuality = config.VideoQuality
		batchConfigPtr = &batchConfig
	} else {
		singleConfig = &config
	}

	// 統合サービスの作成と実行
	unifiedService := usecases.NewUnifiedMovieConverterService(singleConfig, batchConfigPtr)
	result, err := unifiedService.ProcessConversion()
	if err != nil {
		log.Fatalf("変換エラー: %v", err)
	}

	// 結果の表示
	displayConversionResult(result)
}
