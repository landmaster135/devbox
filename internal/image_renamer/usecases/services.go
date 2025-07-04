package usecases

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// FileInfo はファイル情報を保持する構造体です
type FileInfo struct {
	Path    string
	ModTime int64
	Name    string
}

// Job はリネームジョブを表す構造体です
type Job struct {
	File      FileInfo
	NewSerial int
}

// Config はプログラムの設定を保持する構造体です
type Config struct {
	SrcDir     string
	SortByName bool
	SortByTime bool
	Prefix     string
	Delimiter  string
	Digits     int
	StartCount int
	Recursive  bool
	Workers    int
}

// ValidateConfig は設定の妥当性を検証します
func ValidateConfig(config Config, stderr io.Writer) error {
	// プレフィックスが指定されていない場合はエラーを表示して終了
	if config.Prefix == "" {
		fmt.Fprintln(stderr, "エラー: プレフィックスは必須です。-prefix フラグを使用して記事番号を指定してください。")
		fmt.Fprintln(stderr, "例: ./image-renamer -prefix \"20250507\" -time")
		return fmt.Errorf("プレフィックスが指定されていません")
	}

	// 並び替え方法のチェック：両方ともfalseならエラー
	if !config.SortByTime && !config.SortByName {
		fmt.Fprintln(stderr, "エラー: -time または -name のいずれかの並べ替え方法を指定する必要があります。")
		fmt.Fprintln(stderr, "例: ./image-renamer -prefix \"20250507\" -time")
		return fmt.Errorf("並べ替え方法が指定されていません")
	}

	// 並び替え方法の排他制御：両方がtrueの場合は警告を表示
	if config.SortByTime && config.SortByName {
		fmt.Fprintln(stderr, "警告: -time と -name の両方のフラグが設定されています。-name の並べ替え方法を使用します。")
		config.SortByTime = false
	}

	// ディレクトリの存在確認
	_, err := os.Stat(config.SrcDir)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: ディレクトリ %s へのアクセスエラー: %v\n", config.SrcDir, err)
		return err
	}

	return nil
}

// FindImageFiles は指定されたディレクトリから画像ファイルを検索します
func FindImageFiles(srcDir string, recursive bool, stdout, stderr io.Writer) ([]string, error) {
	var files []string

	if recursive {
		err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
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
			fmt.Fprintf(stderr, "エラー: ディレクトリ %s の走査中にエラーが発生しました: %v\n", srcDir, err)
			return nil, err
		}
	} else {
		entries, err := os.ReadDir(srcDir)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: ディレクトリ %s の読み込みに失敗しました: %v\n", srcDir, err)
			return nil, err
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				ext := strings.ToLower(filepath.Ext(entry.Name()))
				if isImageExt(ext) {
					files = append(files, filepath.Join(srcDir, entry.Name()))
				}
			}
		}
	}

	return files, nil
}

// GetFileInfos はファイルパスのリストからファイル情報を取得します
func GetFileInfos(files []string, stderr io.Writer) ([]FileInfo, error) {
	fileInfos := make([]FileInfo, len(files))
	var hasError bool

	for i, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: ファイル %s の情報取得に失敗しました: %v\n", file, err)
			hasError = true
			continue
		}
		fileInfos[i] = FileInfo{
			Path:    file,
			ModTime: info.ModTime().Unix(),
			Name:    info.Name(),
		}
	}

	if hasError {
		return fileInfos, fmt.Errorf("一部のファイル情報の取得に失敗しました")
	}
	return fileInfos, nil
}

// SortFiles はファイル情報を指定された方法で並べ替えます
func SortFiles(fileInfos []FileInfo, sortByTime bool, stdout io.Writer) {
	if sortByTime {
		fmt.Fprintln(stdout, "ファイルを更新日時順に並べ替えています（古い順）")
		sort.Slice(fileInfos, func(i, j int) bool {
			return fileInfos[i].ModTime < fileInfos[j].ModTime
		})
	} else {
		fmt.Fprintln(stdout, "ファイルを名前順に並べ替えています")
		sort.Slice(fileInfos, func(i, j int) bool {
			return fileInfos[i].Name < fileInfos[j].Name
		})
	}
}

// RenameFiles はファイルをリネームします
func RenameFiles(fileInfos []FileInfo, config Config, stdout, stderr io.Writer) (int, int) {
	// ワーカープールの設定
	workerCount := config.Workers
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
	successCount := 0
	errorCount := 0

	// ジョブの準備（シリアル番号を事前に割り当て）
	jobs := prepareJobs(fileInfos, config.StartCount)

	// ジョブチャネル
	jobChan := make(chan Job, len(jobs))

	// ワーカーの起動
	startWorkers(workerCount, jobChan, &wg, &mu, &successCount, &errorCount, config.Digits, config.Prefix, config.Delimiter, stdout, stderr)

	// ジョブの送信
	for _, job := range jobs {
		jobChan <- job
	}
	close(jobChan)

	// すべてのワーカーが完了するのを待つ
	wg.Wait()

	return successCount, errorCount
}

// prepareJobs はリネームジョブを準備します
func prepareJobs(fileInfos []FileInfo, startCount int) []Job {
	jobs := make([]Job, len(fileInfos))
	for i, file := range fileInfos {
		jobs[i] = Job{
			File:      file,
			NewSerial: startCount + i,
		}
	}
	return jobs
}

// startWorkers はリネームワーカーを起動します
func startWorkers(workerCount int, jobChan <-chan Job, wg *sync.WaitGroup, mu *sync.Mutex, successCount, errorCount *int, digits int, prefix string, delimiter string, stdout, stderr io.Writer) {
	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobChan {
				processRenameJob(job, digits, prefix, delimiter, mu, successCount, errorCount, stdout, stderr)
			}
		}()
	}
}

// processRenameJob は1つのリネームジョブを処理します
func processRenameJob(job Job, digits int, prefix string, delimiter string, mu *sync.Mutex, successCount, errorCount *int, stdout, stderr io.Writer) {
	// ファイルのリネーム
	oldPath := job.File.Path
	dir := filepath.Dir(oldPath)
	oldName := filepath.Base(oldPath)
	ext := filepath.Ext(oldName)

	// シリアル番号を指定桁数になるようにフォーマット
	formatStr := fmt.Sprintf("%%0%dd", digits)
	serial := fmt.Sprintf(formatStr, job.NewSerial)
	newName := fmt.Sprintf("%s%s%s%s", prefix, delimiter, serial, ext)
	newPath := filepath.Join(dir, newName)

	fmt.Fprintf(stdout, "処理中: %s -> %s\n", oldPath, newPath)

	err := os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %s のリネームに失敗しました: %v\n", oldPath, err)
		mu.Lock()
		*errorCount++
		mu.Unlock()
	} else {
		mu.Lock()
		*successCount++
		mu.Unlock()
	}
}

func isImageExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".avif":
		return true
	default:
		return false
	}
}
