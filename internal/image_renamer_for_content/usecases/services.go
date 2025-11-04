package usecases

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	cfg "github.com/landmaster135/devbox/internal/image_renamer_for_content/config"
)

type fileInfo struct {
	path    string
	name    string
	modTime int64
}

type renameJob struct {
	info   fileInfo
	serial int
}

var supportedExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".webp": {},
	".avif": {},
	".bmp":  {},
	".gif":  {},
}

// ProcessContentImageRename はコンテンツIDを基に画像ファイルを連番リネームします。
func ProcessContentImageRename(config cfg.Config, stdout, stderr io.Writer) (int, int, error) {
	config.Normalize()

	if err := applyOperationPreset(&config, stderr); err != nil {
		return 0, 0, err
	}

	if err := validateConfig(&config, stderr); err != nil {
		return 0, 0, err
	}

	paths, err := findImageFiles(config.SrcDir, config.Recursive)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: ディレクトリ %s の走査中に失敗しました: %v\n", config.SrcDir, err)
		return 0, 0, err
	}

	if len(paths) == 0 {
		fmt.Fprintln(stdout, "画像ファイルが見つかりませんでした。")
		return 0, 0, nil
	}

	infos, statErr := buildFileInfos(paths, stderr)
	if len(infos) == 0 {
		if statErr != nil {
			return 0, 0, statErr
		}
		fmt.Fprintln(stdout, "画像ファイルが見つかりませんでした。")
		return 0, 0, nil
	}

	if statErr != nil {
		fmt.Fprintf(stderr, "警告: 一部のファイル情報の取得に失敗しました: %v\n", statErr)
	}

	sortFileInfos(infos, config.SortByTime)

	fmt.Fprintf(stdout, "対象ファイル数: %d\n", len(infos))
	fmt.Fprintf(stdout, "コンテンツID: %s\n", config.ContentID)
	fmt.Fprintf(stdout, "開始番号: %d\n", config.Start)
	if config.Delimiter != "" {
		fmt.Fprintf(stdout, "区切り文字: %s\n", config.Delimiter)
	}
	if config.Suffix != "" {
		fmt.Fprintf(stdout, "サフィックス: %s\n", config.Suffix)
	}
	if config.SortByTime {
		fmt.Fprintln(stdout, "更新日時順で処理します。")
	} else {
		fmt.Fprintln(stdout, "ファイル名順で処理します。")
	}

	successCount, errorCount := renameFiles(infos, config, stdout, stderr)

	fmt.Fprintln(stdout, "✔ ファイルリネームが完了しました")
	fmt.Fprintf(stdout, "  成功: %d ファイル\n", successCount)
	if errorCount > 0 {
		fmt.Fprintf(stdout, "  失敗: %d ファイル\n", errorCount)
		return successCount, errorCount, fmt.Errorf("一部のファイルのリネームに失敗しました")
	}

	return successCount, errorCount, nil
}

func applyOperationPreset(config *cfg.Config, stderr io.Writer) error {
	switch config.Operation {
	case "mackerel":
		config.ContentID = "MA"
		config.Digits = 4
	case "web_clip":
		config.ContentID = "WC"
		config.Digits = 9
	case "date":
		config.ContentID = "DA"
		config.Digits = 5
	case "wine":
		config.ContentID = "WI"
		config.Digits = 4
	default:
		if config.Operation == "" {
			fmt.Fprintln(stderr, "エラー: -operation フラグで実行モードを指定してください。例: -operation mackerel / web_clip / date / wine")
		} else {
			fmt.Fprintf(stderr, "エラー: 未対応のoperationが指定されました: %s\n", config.Operation)
		}
		return errors.New("invalid operation")
	}

	return nil
}

func validateConfig(config *cfg.Config, stderr io.Writer) error {
	if config.ContentID == "" {
		fmt.Fprintln(stderr, "エラー: コンテンツIDが設定されていません。operation を確認してください。")
		return errors.New("content id is required")
	}

	if !config.SortByName && !config.SortByTime {
		fmt.Fprintln(stderr, "エラー: -name または -time のいずれかの並べ替え方法を指定する必要があります。")
		return errors.New("sort flag is required")
	}

	if config.SortByName && config.SortByTime {
		fmt.Fprintln(stderr, "警告: -name と -time の両方が指定されています。-name を優先します。")
		config.SortByTime = false
	}

	if config.Digits <= 0 {
		fmt.Fprintln(stderr, "エラー: digits は1以上を指定してください。")
		return errors.New("digits must be positive")
	}

	if config.Start <= 0 {
		fmt.Fprintln(stderr, "エラー: start は1以上を指定してください。")
		return errors.New("start must be positive")
	}

	if _, err := os.Stat(config.SrcDir); err != nil {
		fmt.Fprintf(stderr, "エラー: ディレクトリ %s へのアクセスに失敗しました: %v\n", config.SrcDir, err)
		return err
	}

	return nil
}

func findImageFiles(root string, recursive bool) ([]string, error) {
	var files []string

	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if !recursive && path != root {
				return filepath.SkipDir
			}
			return nil
		}

		if isSupportedImage(d.Name()) {
			files = append(files, path)
		}
		return nil
	}

	if recursive {
		if err := filepath.WalkDir(root, walkFunc); err != nil {
			return nil, err
		}
		return files, nil
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if isSupportedImage(entry.Name()) {
			files = append(files, filepath.Join(root, entry.Name()))
		}
	}

	return files, nil
}

func isSupportedImage(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	_, ok := supportedExtensions[ext]
	return ok
}

func buildFileInfos(paths []string, stderr io.Writer) ([]fileInfo, error) {
	infos := make([]fileInfo, 0, len(paths))
	var failed int

	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: ファイル %s の情報取得に失敗しました: %v\n", path, err)
			failed++
			continue
		}

		infos = append(infos, fileInfo{
			path:    path,
			name:    stat.Name(),
			modTime: stat.ModTime().UnixNano(),
		})
	}

	if failed > 0 {
		return infos, fmt.Errorf("%d 件のファイル情報の取得に失敗しました", failed)
	}

	return infos, nil
}

func sortFileInfos(infos []fileInfo, sortByTime bool) {
	if sortByTime {
		sort.SliceStable(infos, func(i, j int) bool {
			return infos[i].modTime < infos[j].modTime
		})
		return
	}

	sort.SliceStable(infos, func(i, j int) bool {
		return infos[i].name < infos[j].name
	})
}

func renameFiles(infos []fileInfo, config cfg.Config, stdout, stderr io.Writer) (int, int) {
	if len(infos) == 0 {
		return 0, 0
	}

	workerCount := config.Workers
	if workerCount > len(infos) {
		workerCount = len(infos)
	}
	if workerCount < 1 {
		workerCount = 1
	}

	fmt.Fprintf(stdout, "リネーム操作に %d ワーカーを使用します。\n", workerCount)

	jobs := make(chan renameJob, len(infos))
	var wg sync.WaitGroup
	var mu sync.Mutex
	successCount := 0
	errorCount := 0

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				if err := executeRename(job, config, stdout, stderr); err != nil {
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

	for idx, info := range infos {
		jobs <- renameJob{info: info, serial: config.Start + idx}
	}
	close(jobs)

	wg.Wait()

	return successCount, errorCount
}

func executeRename(job renameJob, config cfg.Config, stdout, stderr io.Writer) error {
	oldPath := job.info.path
	dir := filepath.Dir(oldPath)
	ext := strings.ToLower(filepath.Ext(job.info.name))

	if ext == "" {
		fmt.Fprintf(stderr, "エラー: %s は拡張子がないためスキップします。\n", oldPath)
		return errors.New("missing file extension")
	}

	newName := buildNewFileName(config, job.serial, ext)
	newPath := filepath.Join(dir, newName)

	if oldPath == newPath {
		fmt.Fprintf(stdout, "スキップ: %s は既に目的のファイル名です。\n", oldPath)
		return nil
	}

	fmt.Fprintf(stdout, "処理中: %s -> %s\n", oldPath, newPath)

	if err := os.Rename(oldPath, newPath); err != nil {
		fmt.Fprintf(stderr, "エラー: %s のリネームに失敗しました: %v\n", oldPath, err)
		return err
	}

	return nil
}

func buildNewFileName(config cfg.Config, serial int, ext string) string {
	serialStr := fmt.Sprintf("%0*d", config.Digits, serial)
	base := config.ContentID + serialStr
	if config.Delimiter != "" {
		base = config.ContentID + config.Delimiter + serialStr
	}

	if config.Suffix != "" {
		return fmt.Sprintf("%s_%s%s", base, config.Suffix, ext)
	}

	return base + ext
}
