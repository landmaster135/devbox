package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	usecases "github.com/landmaster135/devbox/internal/image_filterer/usecases"
)

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

// run は画像フィルタリングツールの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグ解析
	procConfig, filterConfig, err := parseFlags(args, stderr)
	if err != nil {
		return exitCodeError
	}

	// 処理情報の表示
	printProcessingInfo(procConfig, filterConfig, stdout)

	// サービス作成（バリデーション含む）
	service, err := usecases.NewImageFilterService(procConfig, filterConfig)
	if err != nil {
		fmt.Fprintf(stderr, "設定エラー: %v\n", err)
		return exitCodeError
	}

	// 処理実行
	result, err := service.ProcessImages()
	if err != nil {
		fmt.Fprintf(stderr, "処理エラー: %v\n", err)
		return exitCodeError
	}

	// 結果出力
	return printResults(result, stdout)
}

// parseFlags はコマンドライン引数を解析してコンフィグを作成します
func parseFlags(args []string, stderr io.Writer) (*usecases.ProcessingConfig, *usecases.FilterConfig, error) {
	// フラグセットを作成
	flagSet := flag.NewFlagSet("image-filterer", flag.ContinueOnError)
	flagSet.SetOutput(stderr)

	// コマンドライン引数の定義
	srcDir := flagSet.String("src", ".", "source directory to scan")
	outDir := flagSet.String("out", "", "output directory (output into src if empty)")
	arcDir := flagSet.String("arc", "./5_original_files", "move processed originals to this directory")
	x1 := flagSet.Int("x1", 0, "X-coordinate on upper left")
	y1 := flagSet.Int("y1", 0, "Y-coordinate on upper left")
	x2 := flagSet.Int("x2", 0, "X-coordinate on lower right")
	y2 := flagSet.Int("y2", 0, "Y-coordinate on lower right")
	suffix := flagSet.String("suffix", "filtered", "suffix to attach to file name to save")
	move := flagSet.Bool("move", false, "move originals instead of copying (effective only with -archive)")
	recursive := flagSet.Bool("r", false, "recursively scan sub-directories")
	workers := flagSet.Int("workers", runtime.NumCPU(), "number of concurrent workers")
	mode := flagSet.String("mode", "blur", "filter mode ('blur' or 'grayscale')")
	radius := flagSet.Float64("radius", 10.0, "blur radius (effective only with -mode=blur)")
	rWeight := flagSet.Float64("r-weight", 0.3, "red weight for grayscale conversion (0.0-1.0)")
	gWeight := flagSet.Float64("g-weight", 0.6, "green weight for grayscale conversion (0.0-1.0)")
	bWeight := flagSet.Float64("b-weight", 0.1, "blue weight for grayscale conversion (0.0-1.0)")

	// 引数の解析
	if err := flagSet.Parse(args); err != nil {
		return nil, nil, err
	}

	// 出力ディレクトリが指定されていない場合は入力ディレクトリと同じに
	if *outDir == "" {
		*outDir = *srcDir
	}

	// フィルターモードの変換
	var filterMode usecases.FilterMode
	switch strings.ToLower(*mode) {
	case "blur":
		filterMode = usecases.BlurMode
	case "grayscale":
		filterMode = usecases.GrayscaleMode
	default:
		return nil, nil, fmt.Errorf("サポートされていないフィルターモードです: %s", *mode)
	}

	// コンフィグ作成
	procConfig := &usecases.ProcessingConfig{
		SrcDir:    *srcDir,
		OutDir:    *outDir,
		ArcDir:    *arcDir,
		X1:        *x1,
		Y1:        *y1,
		X2:        *x2,
		Y2:        *y2,
		Suffix:    *suffix,
		Move:      *move,
		Recursive: *recursive,
		Workers:   *workers,
	}

	filterConfig := &usecases.FilterConfig{
		Mode:    filterMode,
		Radius:  *radius,
		RWeight: *rWeight,
		GWeight: *gWeight,
		BWeight: *bWeight,
	}

	return procConfig, filterConfig, nil
}

// printProcessingInfo は処理情報を表示します
func printProcessingInfo(procConfig *usecases.ProcessingConfig, filterConfig *usecases.FilterConfig, stdout io.Writer) {
	fmt.Fprintf(stdout, "画像フィルタリング処理を開始します\n")

	if procConfig.X1 == 0 && procConfig.Y1 == 0 && procConfig.X2 == 0 && procConfig.Y2 == 0 {
		fmt.Fprintf(stdout, "画像全体にフィルターを適用します\n")
	} else {
		fmt.Fprintf(stdout, "指定された範囲: (%d,%d)-(%d,%d)\n", procConfig.X1, procConfig.Y1, procConfig.X2, procConfig.Y2)
	}

	if filterConfig.Mode == usecases.GrayscaleMode {
		fmt.Fprintf(stdout, "フィルターモード: %s, RGB重み: (%.2f, %.2f, %.2f)\n",
			filterConfig.Mode, filterConfig.RWeight, filterConfig.GWeight, filterConfig.BWeight)
	} else {
		fmt.Fprintf(stdout, "フィルターモード: %s, 半径: %.1f\n", filterConfig.Mode, filterConfig.Radius)
	}
}

// printResults は処理結果を表示します
func printResults(result *usecases.ProcessingResult, stdout io.Writer) exitCode {
	fmt.Fprintf(stdout, "\n✔ 画像フィルタリングが完了しました\n")
	fmt.Fprintf(stdout, "  成功: %d ファイル\n", result.SuccessCount)

	if result.ErrorCount > 0 {
		fmt.Fprintf(stdout, "  失敗: %d ファイル\n", result.ErrorCount)
		return exitCodeError
	}

	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
