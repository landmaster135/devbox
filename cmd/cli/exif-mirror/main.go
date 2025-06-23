package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	usecases "github.com/landmaster135/devbox/internal/independencies/exif_mirror/usecases"
)

const version = "1.0.0"

// コマンドライン引数
var (
	sourceFolderPath = flag.String("source-dir", "", "ソース画像ファイルがあるディレクトリのパス")
	targetFolderPath = flag.String("target-dir", "", "ターゲット画像ファイルがあるディレクトリのパス")
	sourceExtension  = flag.String("source-ext", "", "ソース拡張子 (例: jpg, jpeg, png)")
	targetExtension  = flag.String("target-ext", "", "ターゲット拡張子 (例: webp, png, jpg)")
	recursive        = flag.Bool("recursive", false, "サブフォルダも再帰的に処理する")
	dryRun           = flag.Bool("dry-run", false, "実際には変更せず、処理対象ファイルのみ表示")
	verbose          = flag.Bool("verbose", false, "詳細な出力を表示")
	showHelp         = flag.Bool("help", false, "ヘルプメッセージを表示")
	showVersion      = flag.Bool("version", false, "バージョン情報を表示")
	workerCount      = flag.Int("workers", runtime.NumCPU(), "並行処理のワーカー数")
)

func main() {
	// カスタムUsage関数を設定
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "exif-mirror - EXIF Data Mirror Tool\n\n")
		fmt.Fprintf(os.Stderr, "USAGE:\n")
		fmt.Fprintf(os.Stderr, "    exif-mirror [OPTIONS]\n\n")
		fmt.Fprintf(os.Stderr, "DESCRIPTION:\n")
		fmt.Fprintf(os.Stderr, "    指定されたソースフォルダ内の画像ファイルから、ターゲットフォルダ内の同名ファイルにEXIFデータをコピーします。\n")
		fmt.Fprintf(os.Stderr, "    PowerShellのCopy-ExifFromSrcコマンドレットと同じ機能を提供します。\n\n")
		fmt.Fprintf(os.Stderr, "OPTIONS:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "    # JPGファイルからWebPファイルにEXIFをコピー\n")
		fmt.Fprintf(os.Stderr, "    exif-mirror --source-dir ./originals --target-dir ./converted --source-ext jpg --target-ext webp\n\n")
		fmt.Fprintf(os.Stderr, "    # PNGファイルからWebPファイルにEXIFをコピー（再帰処理）\n")
		fmt.Fprintf(os.Stderr, "    exif-mirror --source-dir ./photos --target-dir ./photos --source-ext png --target-ext webp --recursive\n\n")
		fmt.Fprintf(os.Stderr, "    # ドライラン（実際の処理は行わず、対象ファイルのみ表示）\n")
		fmt.Fprintf(os.Stderr, "    exif-mirror --source-dir ./src --target-dir ./dst --source-ext jpg --target-ext png --dry-run\n\n")
		fmt.Fprintf(os.Stderr, "    # 詳細ログ付きで実行\n")
		fmt.Fprintf(os.Stderr, "    exif-mirror --source-dir ./src --target-dir ./dst --source-ext jpg --target-ext webp --verbose\n\n")
		fmt.Fprintf(os.Stderr, "NOTES:\n")
		fmt.Fprintf(os.Stderr, "    - ソースとターゲットのフォルダは異なっていても同じでも構いません\n")
		fmt.Fprintf(os.Stderr, "    - ファイル名（拡張子除く）が同じファイル同士でEXIFデータがコピーされます\n")
		fmt.Fprintf(os.Stderr, "    - 現在サポートしているフォーマット: JPEG, PNG（限定的にTIFF, WebP）\n")
	}

	flag.Parse()

	// ヘルプまたはバージョン表示
	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("exif-mirror version %s\n", version)
		os.Exit(0)
	}

	// 必須引数の検証
	if *sourceFolderPath == "" {
		fmt.Fprintf(os.Stderr, "Error: --source-dir パラメータが必要です\n")
		flag.Usage()
		os.Exit(1)
	}

	if *targetFolderPath == "" {
		fmt.Fprintf(os.Stderr, "Error: --target-dir パラメータが必要です\n")
		flag.Usage()
		os.Exit(1)
	}

	if *sourceExtension == "" {
		fmt.Fprintf(os.Stderr, "Error: --source-ext パラメータが必要です\n")
		flag.Usage()
		os.Exit(1)
	}

	if *targetExtension == "" {
		fmt.Fprintf(os.Stderr, "Error: --target-ext パラメータが必要です\n")
		flag.Usage()
		os.Exit(1)
	}

	// ディレクトリの存在確認とバリデーション
	if err := usecases.ValidateDirectory(*sourceFolderPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error (source-dir): %v\n", err)
		os.Exit(1)
	}

	if err := usecases.ValidateDirectory(*targetFolderPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error (target-dir): %v\n", err)
		os.Exit(1)
	}

	// 拡張子のバリデーション
	if err := usecases.ValidateExtension(*sourceExtension); err != nil {
		fmt.Fprintf(os.Stderr, "Error (source-ext): %v\n", err)
		os.Exit(1)
	}

	if err := usecases.ValidateExtension(*targetExtension); err != nil {
		fmt.Fprintf(os.Stderr, "Error (target-ext): %v\n", err)
		os.Exit(1)
	}

	// ログ設定
	if *verbose {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	} else {
		log.SetFlags(log.Ldate | log.Ltime)
	}

	// 設定を作成
	config := &usecases.Config{
		SourceFolderPath: *sourceFolderPath,
		TargetFolderPath: *targetFolderPath,
		SourceExtension:  *sourceExtension,
		TargetExtension:  *targetExtension,
		Recursive:        *recursive,
		DryRun:           *dryRun,
		Verbose:          *verbose,
		WorkerCount:      *workerCount,
	}

	// 実行情報を表示
	printExecutionInfo(config)

	// ExifMirrorServiceを作成
	service := usecases.NewExifMirrorService()

	// EXIFデータをミラーリング
	processedCount, errorCount, skipCount, err := service.MirrorExifData(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error mirroring EXIF data: %v\n", err)
		os.Exit(1)
	}

	// 結果を表示
	fmt.Printf("\n処理完了: %d個のファイルを処理しました", processedCount)
	fmt.Printf(" (エラー: %d個, スキップ: %d個)", errorCount, skipCount)
	fmt.Println()

	if errorCount > 0 {
		os.Exit(1)
	}
}

// 実行情報を表示
func printExecutionInfo(config *usecases.Config) {
	fmt.Printf("ソースディレクトリ: %s\n", config.SourceFolderPath)
	fmt.Printf("ターゲットディレクトリ: %s\n", config.TargetFolderPath)
	fmt.Printf("ソース拡張子: %s\n", config.SourceExtension)
	fmt.Printf("ターゲット拡張子: %s\n", config.TargetExtension)
	fmt.Printf("再帰処理: %t\n", config.Recursive)
	fmt.Printf("ドライラン: %t\n", config.DryRun)
	fmt.Printf("詳細モード: %t\n", config.Verbose)
	fmt.Printf("ワーカー数: %d\n", config.WorkerCount)
	fmt.Println()
}
