package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

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
	y2 := flagSet.Int("y2", 0,  "Y-coordinate on lower right")
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
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// 出力ディレクトリが指定されていない場合は入力ディレクトリと同じに
	if *outDir == "" {
		*outDir = *srcDir
	}

	// フィルタリング座標のバリデーション
	if *x2 <= *x1 || *y2 <= *y1 {
		fmt.Fprintf(stderr, "エラー: 無効な座標です。x2 > x1, y2 > y1 である必要があります。\n")
		flagSet.Usage()
		return exitCodeError
	}

	// 処理情報の表示
	fmt.Fprintf(stdout, "画像フィルタリング処理を開始します\n")
	fmt.Fprintf(stdout, "指定された範囲: (%d,%d)-(%d,%d)\n", *x1, *y1, *x2, *y2)
	if strings.ToLower(*mode) == "grayscale" {
		fmt.Fprintf(stdout, "フィルターモード: %s, RGB重み: (%.2f, %.2f, %.2f)\n", *mode, *rWeight, *gWeight, *bWeight)
	} else {
		fmt.Fprintf(stdout, "フィルターモード: %s, 半径: %.1f\n", *mode, *radius)
	}

	// フィルターモードのバリデーション
	var filterMode usecases.FilterMode
	switch strings.ToLower(*mode) {
	case "blur":
		filterMode = usecases.BlurMode
	case "grayscale":
		filterMode = usecases.GrayscaleMode
	default:
		fmt.Fprintf(stderr, "エラー: サポートされていないフィルターモードです: %s\n", *mode)
		flagSet.Usage()
		return exitCodeError
	}

	// RGB重みのバリデーション
	if *rWeight < 0.0 || *rWeight > 1.0 {
		fmt.Fprintf(stderr, "エラー: r-weightは0.0-1.0の範囲で指定してください: %.2f\n", *rWeight)
		flagSet.Usage()
		return exitCodeError
	}
	if *gWeight < 0.0 || *gWeight > 1.0 {
		fmt.Fprintf(stderr, "エラー: g-weightは0.0-1.0の範囲で指定してください: %.2f\n", *gWeight)
		flagSet.Usage()
		return exitCodeError
	}
	if *bWeight < 0.0 || *bWeight > 1.0 {
		fmt.Fprintf(stderr, "エラー: b-weightは0.0-1.0の範囲で指定してください: %.2f\n", *bWeight)
		flagSet.Usage()
		return exitCodeError
	}

	// 処理対象ファイルのパスを収集
	paths := make(chan string, 512)
	go func() {
		defer close(paths)

		walkFunc := func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" {
				paths <- path
			}
			return nil
		}

		if *recursive {
			// 再帰的にディレクトリを走査
			err := filepath.WalkDir(*srcDir, walkFunc)
			if err != nil {
				fmt.Fprintf(stderr, "エラー: ディレクトリの走査中にエラーが発生しました: %v\n", err)
			}
		} else {
			// 単一ディレクトリのみ処理
			entries, err := os.ReadDir(*srcDir)
			if err != nil {
				fmt.Fprintf(stderr, "エラー: ディレクトリの読み込みに失敗しました: %v\n", err)
				return
			}

			for _, e := range entries {
				if e.IsDir() {
					continue
				}
				ext := strings.ToLower(filepath.Ext(e.Name()))
				if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" {
					paths <- filepath.Join(*srcDir, e.Name())
				}
			}
		}
	}()

	// ワーカープールの作成と実行
	var wg sync.WaitGroup
	errorCount := 0
	successCount := 0
	var countMutex sync.Mutex

	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range paths {
				// 画像にフィルターを適用して保存
				if err := usecases.ApplyFilterAndSave(p, *outDir, *x1, *y1, *x2, *y2, *suffix, filterMode, *radius, *rWeight, *gWeight, *bWeight); err != nil {
					fmt.Fprintf(stderr, "警告: %v\n", err)
					countMutex.Lock()
					errorCount++
					countMutex.Unlock()
					continue
				}

				// 元ファイルを移動（オプションが有効な場合）
				if *move {
					if err := usecases.MoveOriginal(p, *arcDir); err != nil {
						fmt.Fprintf(stderr, "警告: %v\n", err)
						countMutex.Lock()
						errorCount++
						countMutex.Unlock()
					}
				}

				countMutex.Lock()
				successCount++
				countMutex.Unlock()
			}
		}()
	}

	// すべてのワーカーの完了を待機
	wg.Wait()

	// 処理結果の出力
	fmt.Fprintf(stdout, "\n✔ 画像フィルタリングが完了しました\n")
	fmt.Fprintf(stdout, "  成功: %d ファイル\n", successCount)

	if errorCount > 0 {
		fmt.Fprintf(stdout, "  失敗: %d ファイル\n", errorCount)
		return exitCodeError
	}

	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
