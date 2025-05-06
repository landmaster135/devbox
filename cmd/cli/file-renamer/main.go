package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
)

func main() {
	// コマンドラインフラグの定義
	srcDir := flag.String("src", ".", "スキャンするソースディレクトリ")
	sortByTime := flag.Bool("time", false, "画像ファイルを更新日時順に並べ替え")
	sortByName := flag.Bool("name", false, "画像ファイルをファイル名順に並べ替え")
	prefix := flag.String("prefix", "", "記事番号のプレフィックス (必須)")
	digits := flag.Int("digits", 4, "シリアル番号の桁数")
	startCount := flag.Int("start", 1, "リネーム操作の開始番号")
	recursive := flag.Bool("r", false, "サブディレクトリを再帰的にスキャン")
	workers := flag.Int("workers", runtime.NumCPU(), "並行ワーカー数")

	flag.Parse()

	// プレフィックスが指定されていない場合はエラーを表示して終了
	if *prefix == "" {
		fmt.Println("Error: prefix is required. Use -prefix flag to specify the article number.")
		fmt.Println("Example: ./rename-images -prefix \"20250507\" -time")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// 並び替え方法のチェック：両方ともfalseならエラー
	if !*sortByTime && !*sortByName {
		fmt.Println("Error: You must specify either -time or -name sorting method.")
		fmt.Println("Example: ./rename-images -prefix \"20250507\" -time")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// 並び替え方法の排他制御：両方がtrueの場合は警告を表示
	if *sortByTime && *sortByName {
		fmt.Println("Warning: Both -time and -name flags are set. Using -name sorting method.")
		*sortByTime = false
	}

	// ディレクトリの存在確認
	_, err := os.Stat(*srcDir)
	if err != nil {
		fmt.Printf("Error accessing directory %s: %v\n", *srcDir, err)
		os.Exit(1)
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
			fmt.Printf("Error walking directory %s: %v\n", *srcDir, err)
			os.Exit(1)
		}
	} else {
		entries, err := os.ReadDir(*srcDir)
		if err != nil {
			fmt.Printf("Error reading directory %s: %v\n", *srcDir, err)
			os.Exit(1)
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
		fmt.Println("No image files found.")
		os.Exit(0)
	}

	fmt.Printf("Found %d image files.\n", len(files))
	fmt.Printf("Using prefix: %s\n", *prefix)
	fmt.Printf("Starting count: %d\n", *startCount)

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
			fmt.Printf("Error getting file info for %s: %v\n", file, err)
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
		fmt.Println("Sorting files by modification time (oldest first)")
		sort.Slice(fileInfos, func(i, j int) bool {
			return fileInfos[i].modTime < fileInfos[j].modTime
		})
	} else {
		fmt.Println("Sorting files by name")
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

	fmt.Printf("Using %d workers for renaming operation.\n", workerCount)

	// カウンターとワーカーの同期用
	var mu sync.Mutex
	var wg sync.WaitGroup
	currentCount := *startCount

	// ワーカーチャネル
	jobs := make(chan FileInfo, len(fileInfos))

	// ワーカーの起動
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range jobs {
				// カウンター値の取得（排他制御）
				mu.Lock()
				thisCount := currentCount
				currentCount++
				mu.Unlock()

				// ファイルのリネーム
				oldPath := file.path
				dir := filepath.Dir(oldPath)
				oldName := filepath.Base(oldPath)
				ext := filepath.Ext(oldName)

				// シリアル番号を指定桁数になるようにフォーマット
				formatStr := fmt.Sprintf("%%0%dd", *digits)
				serial := fmt.Sprintf(formatStr, thisCount)
				newName := fmt.Sprintf("%s_%s%s", *prefix, serial, ext)
				newPath := filepath.Join(dir, newName)

				fmt.Printf("Processing %s -> %s\n", oldPath, newPath)

				err := os.Rename(oldPath, newPath)
				if err != nil {
					fmt.Printf("Error renaming %s: %v\n", oldPath, err)
				}
			}
		}()
	}

	// ジョブの送信
	for _, file := range fileInfos {
		jobs <- file
	}
	close(jobs)

	// すべてのワーカーが完了するのを待つ
	wg.Wait()

	fmt.Println("All done!")
}

func isImageExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".avif":
		return true
	default:
		return false
	}
}
