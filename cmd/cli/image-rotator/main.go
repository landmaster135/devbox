package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	usecases "github.com/landmaster135/devbox/internal/image_rotator/usecases"
)

// exitCode は OS に返す終了ステータス
type exitCode int

const (
	exitOK exitCode = iota
	exitErr
)

func main() {
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}

func run(args []string, stdout, stderr io.Writer) exitCode {
	// ---------------------- CLI フラグ定義 ----------------------
	fs := flag.NewFlagSet("image-rotator", flag.ContinueOnError)
	fs.SetOutput(stderr)

	srcDir := fs.String("src", ".", "source directory to scan")
	outDir := fs.String("out", "", "output directory (default: same as src)")
	arcDir := fs.String("arc", "./5_original_files", "archive directory for originals")
	angleDeg := fs.Float64("angle", 0, "rotation angle in degrees clockwise (required)")
	suffix := fs.String("suffix", "rotated", "suffix appended to output file name")
	moveOrig := fs.Bool("move", false, "move originals to archive (instead of copying)")
	recursive := fs.Bool("r", false, "scan directories recursively")
	workers := fs.Int("workers", runtime.NumCPU(), "number of concurrent workers")

	if err := fs.Parse(args); err != nil {
		return exitErr
	}
	if *angleDeg == 0 {
		fmt.Fprintln(stderr, "エラー: -angle は 0 以外を指定してください")
		return exitErr
	}
	if *outDir == "" {
		*outDir = *srcDir
	}

	// ---------------------- 対象ファイル列挙 ----------------------
	paths := make(chan string, 256)
	go func() {
		defer close(paths)
		walk := func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			if usecases.IsImage(path) {
				paths <- path
			}
			return nil
		}
		if *recursive {
			_ = filepath.WalkDir(*srcDir, walk)
		} else {
			entries, err := os.ReadDir(*srcDir)
			if err != nil {
				fmt.Fprintln(stderr, "ディレクトリ読み込み失敗:", err)
				return
			}
			for _, e := range entries {
				if !e.IsDir() && usecases.IsImage(e.Name()) {
					paths <- filepath.Join(*srcDir, e.Name())
				}
			}
		}
	}()

	// ---------------------- ワーカープール ----------------------
	var wg sync.WaitGroup
	var mu sync.Mutex
	success, failed := 0, 0

	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range paths {
				if err := usecases.RotateAndSave(p, *outDir, *suffix, *angleDeg); err != nil {
					mu.Lock()
					failed++
					mu.Unlock()
					fmt.Fprintf(stderr, "失敗 (%s): %v\n", p, err)
					continue
				}
				if *moveOrig {
					if err := usecases.MoveFile(p, *arcDir); err != nil {
						fmt.Fprintf(stderr, "アーカイブ失敗 (%s): %v\n", p, err)
					}
				}
				mu.Lock()
				success++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// ---------------------- 結果表示 ----------------------
	fmt.Fprintf(stdout, "✔ 回転処理完了  成功: %d  失敗: %d\n", success, failed)
	if failed > 0 {
		return exitErr
	}
	return exitOK
}
