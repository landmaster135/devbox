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

	usecases "github.com/landmaster135/devbox/internal/image_converter/usecases"
)

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

// run は画像変換ツールの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグセットを作成
	flagSet := flag.NewFlagSet("image-converter", flag.ContinueOnError)
	flagSet.SetOutput(stderr)

	// コマンドライン引数の定義
	sourceDir := flagSet.String("src-dir", ".", "source directory to scan")
	outputDir := flagSet.String("output-dir", "", "output directory (required)")
	archiveDir := flagSet.String("archive-dir", "", "move processed originals to this directory (disabled if empty)")
	moveOrig := flagSet.Bool("move", false, "move originals instead of copying (effective only with --archive-dir)")
	outExt := flagSet.String("ext", "png", "target extension (png|jpg|webp|avif)")
	quality := flagSet.Int("q", 80, "quality for lossy formats (1-100)")
	workers := flagSet.Int("workers", runtime.NumCPU(), "number of concurrent workers")
	recursive := flagSet.Bool("R", false, "recursively scan sub-directories")
	lossless := flagSet.Bool("lossless", false, "enable lossless compression (for Webp)")

	// 引数の解析
	if err := flagSet.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	if strings.TrimSpace(*outputDir) == "" {
		fmt.Fprintln(stderr, "エラー: -output-dir は必須です")
		flagSet.Usage()
		return exitCodeError
	}

	// 出力フォーマットの正規化
	*outExt = strings.ToLower(strings.TrimPrefix(*outExt, "."))

	// コーデックテーブルの作成
	codecs := usecases.MakeCodecTable(*quality, *lossless)

	// サポートされていないフォーマットのチェック
	if _, ok := codecs[*outExt]; !ok {
		fmt.Fprintf(stderr, "エラー: サポートされていない出力フォーマット: %s\n", *outExt)
		flagSet.Usage()
		return exitCodeError
	}

	// 変換対象ファイルのパスを収集
	paths := make(chan string, 512)
	go func() {
		defer close(paths)

		if *recursive {
			// 再帰的にディレクトリを走査
			err := filepath.WalkDir(*sourceDir, func(path string, d os.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return nil
				}
				ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
				if _, ok := codecs[ext]; ok {
					paths <- path
				}
				return nil
			})

			if err != nil {
				fmt.Fprintf(stderr, "エラー: ディレクトリの走査中にエラーが発生しました: %v\n", err)
			}
		} else {
			// 単一ディレクトリのみ処理
			entries, err := os.ReadDir(*sourceDir)
			if err != nil {
				fmt.Fprintf(stderr, "エラー: ディレクトリの読み込みに失敗しました: %v\n", err)
				return
			}

			for _, e := range entries {
				if e.IsDir() {
					continue
				}
				ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(e.Name())), ".")
				if _, ok := codecs[ext]; ok {
					paths <- filepath.Join(*sourceDir, e.Name())
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
				if err := usecases.ConvertFile(p, *sourceDir, *outputDir, *outExt, codecs); err != nil {
					fmt.Fprintf(stderr, "警告: %v\n", err)
					countMutex.Lock()
					errorCount++
					countMutex.Unlock()
					continue
				}
				if err := usecases.ArchiveOriginal(p, *sourceDir, *archiveDir, *moveOrig); err != nil {
					fmt.Fprintf(stderr, "警告: %v\n", err)
					countMutex.Lock()
					errorCount++
					countMutex.Unlock()
				} else {
					countMutex.Lock()
					successCount++
					countMutex.Unlock()
				}
			}
		}()
	}

	// すべてのワーカーの完了を待機
	wg.Wait()

	// 処理結果の出力
	fmt.Fprintf(stdout, "✔ 画像変換が完了しました\n")
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
