package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
)

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

// run は画像ファイルリネームツールの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグセットを作成
	flagSet := flag.NewFlagSet("image-renamer", flag.ContinueOnError)
	flagSet.SetOutput(stderr)

	// コマンドラインフラグの定義
	srcDir := flagSet.String("src", ".", "スキャンするソースディレクトリ")
	sortByName := flagSet.Bool("name", false, "画像ファイルをファイル名順に並べ替え")
	sortByTime := flagSet.Bool("time", false, "画像ファイルを更新日時順に並べ替え")
	prefix := flagSet.String("prefix", "", "記事番号のプレフィックス (必須)")
	digits := flagSet.Int("digits", 4, "シリアル番号の桁数")
	startCount := flagSet.Int("start", 1, "リネーム操作の開始番号")
	recursive := flagSet.Bool("r", false, "サブディレクトリを再帰的にスキャン")
	workers := flagSet.Int("workers", runtime.NumCPU(), "並行ワーカー数")

	// 引数の解析
	if err := flagSet.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// プレフィックスが指定されていない場合はエラーを表示して終了
	if *prefix == "" {
		fmt.Fprintln(stderr, "エラー: プレフィックスは必須です。-prefix フラグを使用して記事番号を指定してください。")
		fmt.Fprintln(stderr, "例: ./image-renamer -prefix \"20250507\" -time")
		flagSet.Usage()
		return exitCodeError
	}

	// 並び替え方法のチェック：両方ともfalseならエラー
	if !*sortByTime && !*sortByName {
		fmt.Fprintln(stderr, "エラー: -time または -name のいずれかの並べ替え方法を指定する必要があります。")
		fmt.Fprintln(stderr, "例: ./image-renamer -prefix \"20250507\" -time")
		flagSet.Usage()
		return exitCodeError
	}

	// 並び替え方法の排他制御：両方がtrueの場合は警告を表示
	if *sortByTime && *sortByName {
		fmt.Fprintln(stderr, "警告: -time と -name の両方のフラグが設定されています。-name の並べ替え方法を使用します。")
		*sortByTime = false
	}

	// ディレクトリの存在確認
	_, err := os.Stat(*srcDir)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: ディレクトリ %s へのアクセスエラー: %v\n", *srcDir, err)
		return exitCodeError
	}

	// 画像ファイルの検索
	var files []string
	if *recursive {
		err = filepath.WalkDir(*srcDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				ext := strings.ToLower(filepath.Ext(d.Name()))
				if isImageExt(ext) {
					files = append(files, path)
				}
			}
			return nil
		})
		if err != nil {
			fmt.Fprintf(stderr, "エラー: ディレクトリ %s の走査中にエラーが発生しました: %v\n", *srcDir, err)
			return exitCodeError
		}
	} else {
		entries, err := os.ReadDir(*srcDir)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: ディレクトリ %s の読み込みに失敗しました: %v\n", *srcDir, err)
			return exitCodeError
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				ext := strings.ToLower(filepath.Ext(entry.Name()))
				if isImageExt(ext) {
					files = append(files, filepath.Join(*srcDir, entry.Name()))
				}
			}
		}
	}

	if len(files) == 0 {
		fmt.Fprintln(stdout, "画像ファイルが見つかりませんでした。")
		return exitCodeOK
	}

	fmt.Fprintf(stdout, "画像ファイルが %d 件見つかりました。\n", len(files))
	fmt.Fprintf(stdout, "プレフィックス: %s\n", *prefix)
	fmt.Fprintf(stdout, "開始番号: %d\n", *startCount)

	// ファイル情報の取得
	type FileInfo struct {
		path    string
		modTime int64
		name    string
	}

	fileInfos := make([]FileInfo, len(files))
	for i, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: ファイル %s の情報取得に失敗しました: %v\n", file, err)
			continue
		}
		fileInfos[i] = FileInfo{
			path:    file,
			modTime: info.ModTime().Unix(),
			name:    info.Name(),
		}
	}

	// 選択された方法でファイルをソート
	if *sortByTime {
		fmt.Fprintln(stdout, "ファイルを更新日時順に並べ替えています（古い順）")
		sort.Slice(fileInfos, func(i, j int) bool {
			return fileInfos[i].modTime < fileInfos[j].modTime
		})
	} else {
		fmt.Fprintln(stdout, "ファイルを名前順に並べ替えています")
		sort.Slice(fileInfos, func(i, j int) bool {
			return fileInfos[i].name < fileInfos[j].name
		})
	}

	// ワーカープールの設定
	workerCount := *workers
	if workerCount < 1 {
		workerCount = 1
	}
	if workerCount > len(fileInfos) {
		workerCount = len(fileInfos)
	}

	fmt.Fprintf(stdout, "リネーム操作に %d ワーカーを使用します。\n", workerCount)

	// カウンターとワーカーの同期用
	var mu sync.Mutex
	var wg sync.WaitGroup
	currentCount := *startCount
	successCount := 0
	errorCount := 0

	// ジョブ構造体の定義
	type Job struct {
		file      FileInfo
		newSerial int
	}

	// ジョブの準備（シリアル番号を事前に割り当て）
	jobs := make([]Job, len(fileInfos))
	for i, file := range fileInfos {
		jobs[i] = Job{
			file:      file,
			newSerial: currentCount + i,
		}
	}

	// ジョブチャネル
	jobChan := make(chan Job, len(jobs))

	// ワーカーの起動
	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobChan {
				// ファイルのリネーム
				oldPath := job.file.path
				dir := filepath.Dir(oldPath)
				oldName := filepath.Base(oldPath)
				ext := filepath.Ext(oldName)

				// シリアル番号を指定桁数になるようにフォーマット
				formatStr := fmt.Sprintf("%%0%dd", *digits)
				serial := fmt.Sprintf(formatStr, job.newSerial)
				newName := fmt.Sprintf("%s_%s%s", *prefix, serial, ext)
				newPath := filepath.Join(dir, newName)

				fmt.Fprintf(stdout, "処理中: %s -> %s\n", oldPath, newPath)

				err := os.Rename(oldPath, newPath)
				if err != nil {
					fmt.Fprintf(stderr, "エラー: %s のリネームに失敗しました: %v\n", oldPath, err)
					mu.Lock()
					errorCount++
					mu.Unlock()
				} else {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
			}
		}()
	}

	// ジョブの送信
	for _, job := range jobs {
		jobChan <- job
	}
	close(jobChan)

	// すべてのワーカーが完了するのを待つ
	wg.Wait()

	// 処理結果の出力
	fmt.Fprintf(stdout, "✔ ファイルリネームが完了しました\n")
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

func isImageExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".avif":
		return true
	default:
		return false
	}
}
