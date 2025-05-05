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

	usecases "github.com/landmaster135/devbox/internal/independencies/image_trimmer/usecases"
)

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

// run は画像トリミングツールの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグセットを作成
	flagSet := flag.NewFlagSet("image-trimmer", flag.ContinueOnError)
	flagSet.SetOutput(stderr)

	// コマンドライン引数の定義
	srcDir := flagSet.String("src", ".", "source directory to scan")
	outDir := flagSet.String("out", "", "output directory (output into src if empty)")
	arcDir := flagSet.String("arc", "./5_original_files", "move processed originals to this directory")
	x1 := flagSet.Int("x1", 0, "X-coordinate on upper left")
	y1 := flagSet.Int("y1", 0, "Y-coordinate on upper left")
	x2 := flagSet.Int("x2", 0, "X-coordinate on lower right")
	y2 := flagSet.Int("y2", 0,  "Y-coordinate on lower right")
	suffix := flagSet.String("suffix", "trimmed", "suffix to attach to file name to save")
	move := flagSet.Bool("move", false, "move originals instead of copying (effective only with -archive)")
	recursive := flagSet.Bool("r", false, "recursively scan sub-directories")
	workers := flagSet.Int("workers", runtime.NumCPU(), "number of concurrent workers")

	// 引数の解析
	if err := flagSet.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// 出力ディレクトリが指定されていない場合は入力ディレクトリと同じに
	if *outDir == "" {
		*outDir = *srcDir
	}

	// トリミング座標のバリデーション
	if *x2 <= *x1 || *y2 <= *y1 {
		fmt.Fprintf(stderr, "エラー: 無効なトリミング座標です。x2 > x1, y2 > y1 である必要があります。\n")
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
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
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
				if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
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
				// 画像をトリミングして保存
				if err := usecases.CropAndSave(p, *outDir, *x1, *y1, *x2, *y2, *suffix); err != nil {
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
	fmt.Fprintf(stdout, "✔ 画像トリミングが完了しました\n")
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
