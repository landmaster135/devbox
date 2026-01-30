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
	NewName   string
	NewPath   string
}

// renameStats はリネーム結果のカウンタを保持します
type renameStats struct {
	mu           sync.Mutex
	successCount int
	errorCount   int
}

func (s *renameStats) recordSuccess() {
	s.mu.Lock()
	s.successCount++
	s.mu.Unlock()
}

func (s *renameStats) recordError() {
	s.mu.Lock()
	s.errorCount++
	s.mu.Unlock()
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

// validateConfig は設定の妥当性を検証します
func validateConfig(config Config, stderr io.Writer) error {
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

func isImageExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".avif":
		return true
	default:
		return false
	}
}

// findImageFiles は指定されたディレクトリから画像ファイルを検索します
func findImageFiles(srcDir string, recursive bool, stdout, stderr io.Writer) ([]string, error) {
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

// getFileInfos はファイルパスのリストからファイル情報を取得します
func getFileInfos(files []string, stderr io.Writer) ([]FileInfo, error) {
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

// sortFiles はファイル情報を指定された方法で並べ替えます
func sortFiles(fileInfos []FileInfo, sortByTime bool, stdout io.Writer) {
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

// renameFiles はファイルをリネームします
func renameFiles(fileInfos []FileInfo, config Config, stdout, stderr io.Writer) (int, int, error) {
	if len(fileInfos) == 0 {
		return 0, 0, nil
	}

	// ワーカープールの設定
	workerCount := config.Workers
	if workerCount < 1 {
		workerCount = 1
	}
	if workerCount > len(fileInfos) {
		workerCount = len(fileInfos)
	}

	// ジョブの準備（シリアル番号と新パスを事前に割り当て）
	jobs := prepareJobs(fileInfos, config)

	// リネーム前後のファイル名衝突チェック
	if err := detectRenameConflicts(fileInfos, jobs, stderr); err != nil {
		return 0, 0, err
	}

	fmt.Fprintf(stdout, "リネーム操作に %d ワーカーを使用します。\n", workerCount)

	stats := &renameStats{}
	var wg sync.WaitGroup
	jobChan := make(chan Job, len(jobs))

	// ワーカーの起動
	startWorkers(workerCount, jobChan, &wg, stats, stdout, stderr)

	// ジョブの送信
	for _, job := range jobs {
		jobChan <- job
	}
	close(jobChan)

	// すべてのワーカーが完了するのを待つ
	wg.Wait()

	return stats.successCount, stats.errorCount, nil
}

// prepareJobs はリネームジョブを準備します
func prepareJobs(fileInfos []FileInfo, config Config) []Job {
	jobs := make([]Job, len(fileInfos))
	formatStr := fmt.Sprintf("%%0%dd", config.Digits)
	for i, file := range fileInfos {
		serial := config.StartCount + i
		serialStr := fmt.Sprintf(formatStr, serial)
		ext := filepath.Ext(file.Path)
		newName := fmt.Sprintf("%s%s%s%s", config.Prefix, config.Delimiter, serialStr, ext)
		newPath := filepath.Join(filepath.Dir(file.Path), newName)
		jobs[i] = Job{
			File:      file,
			NewSerial: serial,
			NewName:   newName,
			NewPath:   newPath,
		}
	}
	return jobs
}

// detectRenameConflicts はリネーム前後のパス衝突を検出します
func detectRenameConflicts(fileInfos []FileInfo, jobs []Job, stderr io.Writer) error {
	if len(jobs) == 0 {
		return nil
	}

	existingPaths := make(map[string]struct{}, len(fileInfos))
	for _, info := range fileInfos {
		existingPaths[info.Path] = struct{}{}
	}

	plannedPaths := make(map[string]string, len(jobs))
	var conflicts []string

	for _, job := range jobs {
		// 既存ファイルとの衝突チェック
		if _, ok := existingPaths[job.NewPath]; ok {
			conflicts = append(conflicts, fmt.Sprintf("%s -> %s (既存ファイルと衝突)", job.File.Path, job.NewPath))
			continue
		}

		// リネーム先同士の重複チェック
		if prev, ok := plannedPaths[job.NewPath]; ok {
			conflicts = append(conflicts, fmt.Sprintf("%s と %s が同じリネーム先 %s を要求", job.File.Path, prev, job.NewPath))
			continue
		}

		plannedPaths[job.NewPath] = job.File.Path
	}

	if len(conflicts) > 0 {
		fmt.Fprintln(stderr, "エラー: リネーム予定のファイル名が既存ファイルと衝突しています。以下を確認してください:")
		for _, msg := range conflicts {
			fmt.Fprintf(stderr, "  - %s\n", msg)
		}
		return fmt.Errorf("リネーム先のファイルパスが衝突しています")
	}

	return nil
}

// startWorkers はリネームワーカーを起動します
func startWorkers(workerCount int, jobChan <-chan Job, wg *sync.WaitGroup, stats *renameStats, stdout, stderr io.Writer) {
	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobChan {
				processRenameJob(job, stats, stdout, stderr)
			}
		}()
	}
}

// processRenameJob は1つのリネームジョブを処理します
func processRenameJob(job Job, stats *renameStats, stdout, stderr io.Writer) {
	oldPath := job.File.Path
	newPath := job.NewPath

	fmt.Fprintf(stdout, "処理中: %s -> %s\n", oldPath, newPath)

	if err := os.Rename(oldPath, newPath); err != nil {
		fmt.Fprintf(stderr, "エラー: %s のリネームに失敗しました: %v\n", oldPath, err)
		stats.recordError()
		return
	}

	stats.recordSuccess()
}

// ProcessImageRename は画像ファイルのリネーム処理全体を実行します
func ProcessImageRename(config Config, stdout, stderr io.Writer) (int, int, error) {
	// 設定の検証
	if err := validateConfig(config, stderr); err != nil {
		return 0, 0, err
	}

	// 画像ファイルの検索
	files, err := findImageFiles(config.SrcDir, config.Recursive, stdout, stderr)
	if err != nil {
		return 0, 0, err
	}

	if len(files) == 0 {
		fmt.Fprintln(stdout, "画像ファイルが見つかりませんでした。")
		return 0, 0, nil
	}

	fmt.Fprintf(stdout, "画像ファイルが %d 件見つかりました。\n", len(files))
	fmt.Fprintf(stdout, "プレフィックス: %s\n", config.Prefix)
	fmt.Fprintf(stdout, "区切り文字: %s\n", config.Delimiter)
	fmt.Fprintf(stdout, "開始番号: %d\n", config.StartCount)

	// ファイル情報の取得と並べ替え
	fileInfos, err := getFileInfos(files, stderr)
	if err != nil {
		// エラーがあっても続行するため、ここではエラーコードを返さない
	}

	// ファイルの並べ替え
	sortFiles(fileInfos, config.SortByTime, stdout)

	// リネーム処理の実行
	successCount, errorCount, err := renameFiles(fileInfos, config, stdout, stderr)
	if err != nil {
		return successCount, errorCount, err
	}

	return successCount, errorCount, nil
}
