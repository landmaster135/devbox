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
	Extensions []string
}

// validateConfig は設定の妥当性を検証します
func validateConfig(config Config, stderr io.Writer) error {
	// 並び替え方法のチェック：両方ともfalseならエラー
	if !config.SortByTime && !config.SortByName {
		fmt.Fprintln(stderr, "エラー: -time または -name のいずれかの並べ替え方法を指定する必要があります。")
		fmt.Fprintln(stderr, "例: ./image-renamer -name")
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

var defaultImageExtensions = []string{".jpg", ".jpeg", ".png", ".webp", ".avif"}
var defaultImageExtensionSet = extensionSet(defaultImageExtensions)

func normalizeExtension(ext string) string {
	normalized := strings.TrimSpace(strings.ToLower(ext))
	if normalized == "" {
		return ""
	}
	if !strings.HasPrefix(normalized, ".") {
		normalized = "." + normalized
	}
	return normalized
}

func normalizeExtensions(extensions []string) []string {
	if len(extensions) == 0 {
		defaults := make([]string, len(defaultImageExtensions))
		copy(defaults, defaultImageExtensions)
		return defaults
	}

	normalized := make([]string, 0, len(extensions))
	seen := make(map[string]struct{}, len(extensions))

	for _, extension := range extensions {
		for _, token := range strings.Split(extension, ",") {
			normalizedExt := normalizeExtension(token)
			if normalizedExt == "" {
				continue
			}
			if _, ok := seen[normalizedExt]; ok {
				continue
			}
			seen[normalizedExt] = struct{}{}
			normalized = append(normalized, normalizedExt)
		}
	}

	if len(normalized) == 0 {
		defaults := make([]string, len(defaultImageExtensions))
		copy(defaults, defaultImageExtensions)
		return defaults
	}

	return normalized
}

func extensionSet(extensions []string) map[string]struct{} {
	normalized := normalizeExtensions(extensions)
	set := make(map[string]struct{}, len(normalized))
	for _, ext := range normalized {
		set[ext] = struct{}{}
	}
	return set
}

func isImageExt(ext string) bool {
	_, ok := defaultImageExtensionSet[normalizeExtension(ext)]
	return ok
}

func isTargetExt(ext string, targetExts map[string]struct{}) bool {
	_, ok := targetExts[normalizeExtension(ext)]
	return ok
}

// findImageFiles は指定されたディレクトリから画像ファイルを検索します
func findImageFiles(srcDir string, recursive bool, extensions []string, stdout, stderr io.Writer) ([]string, error) {
	var files []string
	targetExts := extensionSet(extensions)

	if recursive {
		err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				ext := filepath.Ext(d.Name())
				if isTargetExt(ext, targetExts) {
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
				ext := filepath.Ext(entry.Name())
				if isTargetExt(ext, targetExts) {
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
	workerCount = min(len(fileInfos), max(1, workerCount))

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
		newName := buildRenamedFileName(config.Prefix, config.Delimiter, serialStr, ext)
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

func buildRenamedFileName(prefix, delimiter, serial, ext string) string {
	if prefix == "" {
		return serial + ext
	}

	return prefix + delimiter + serial + ext
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
		wg.Go(func() {
			for job := range jobChan {
				processRenameJob(job, stats, stdout, stderr)
			}
		})
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
	files, err := findImageFiles(config.SrcDir, config.Recursive, config.Extensions, stdout, stderr)
	if err != nil {
		return 0, 0, err
	}

	if len(files) == 0 {
		fmt.Fprintln(stdout, "画像ファイルが見つかりませんでした。")
		return 0, 0, nil
	}

	fmt.Fprintf(stdout, "画像ファイルが %d 件見つかりました。\n", len(files))
	if config.Prefix == "" {
		fmt.Fprintln(stdout, "プレフィックス: (なし)")
	} else {
		fmt.Fprintf(stdout, "プレフィックス: %s\n", config.Prefix)
	}
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
