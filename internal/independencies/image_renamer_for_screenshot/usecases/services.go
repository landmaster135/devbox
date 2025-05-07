package usecases

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// FileInfo はファイル情報を保持する構造体です
type FileInfo struct {
	Path string
	Name string
}

// Config はプログラムの設定を保持する構造体です
type Config struct {
	SrcDir     string
	Recursive  bool
	Workers    int
	VlcPattern bool
	WinPattern bool
}

// ValidateConfig は設定の妥当性を検証します
func ValidateConfig(config Config, stderr io.Writer) error {
	// パターンのチェック：両方ともfalseならエラー
	if !config.VlcPattern && !config.WinPattern {
		fmt.Fprintln(stderr, "エラー: -vlc または -win のいずれかのパターンを指定する必要があります。")
		fmt.Fprintln(stderr, "例: ./image-renamer-for-screenshot -vlc")
		return fmt.Errorf("パターンが指定されていません")
	}

	// パターンの排他制御：両方がtrueの場合はエラーを表示
	if config.VlcPattern && config.WinPattern {
		fmt.Fprintln(stderr, "エラー: -vlc と -win の両方のフラグを同時に設定することはできません。")
		fmt.Fprintln(stderr, "例: ./image-renamer-for-screenshot -vlc")
		return fmt.Errorf("複数のパターンが指定されています")
	}

	// ディレクトリの存在確認
	_, err := os.Stat(config.SrcDir)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: ディレクトリ %s へのアクセスエラー: %v\n", config.SrcDir, err)
		return err
	}

	return nil
}

// FindScreenshotFiles は指定されたディレクトリからスクリーンショットファイルを検索します
func FindScreenshotFiles(srcDir string, recursive bool, vlcPattern, winPattern bool, stdout, stderr io.Writer) ([]string, error) {
	var files []string

	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			ext := strings.ToLower(filepath.Ext(d.Name()))
			if isImageExt(ext) {
				name := d.Name()
				if (vlcPattern && strings.HasPrefix(name, "vlcsnap-")) ||
					(winPattern && strings.HasPrefix(name, "スクリーンショット ")) {
					files = append(files, path)
				}
			}
		}
		return nil
	}

	if recursive {
		err := filepath.WalkDir(srcDir, walkFunc)
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
					name := entry.Name()
					if (vlcPattern && strings.HasPrefix(name, "vlcsnap-")) ||
						(winPattern && strings.HasPrefix(name, "スクリーンショット ")) {
						files = append(files, filepath.Join(srcDir, name))
					}
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
			Path: file,
			Name: info.Name(),
		}
	}

	if hasError {
		return fileInfos, fmt.Errorf("一部のファイル情報の取得に失敗しました")
	}
	return fileInfos, nil
}

// RenameScreenshotFiles はスクリーンショットファイルをリネームします
func RenameScreenshotFiles(fileInfos []FileInfo, config Config, stdout, stderr io.Writer) (int, int) {
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

	// ジョブチャネル
	jobChan := make(chan FileInfo, len(fileInfos))

	// ワーカーの起動
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range jobChan {
				processScreenshotRename(file, config.VlcPattern, config.WinPattern, &mu, &successCount, &errorCount, stdout, stderr)
			}
		}()
	}

	// ジョブの送信
	for _, file := range fileInfos {
		jobChan <- file
	}
	close(jobChan)

	// すべてのワーカーが完了するのを待つ
	wg.Wait()

	return successCount, errorCount
}

// processScreenshotRename は1つのスクリーンショットファイルをリネームします
func processScreenshotRename(file FileInfo, vlcPattern, winPattern bool, mu *sync.Mutex, successCount, errorCount *int, stdout, stderr io.Writer) {
	oldPath := file.Path
	dir := filepath.Dir(oldPath)
	oldName := filepath.Base(oldPath)
	ext := filepath.Ext(oldName)
	baseName := strings.TrimSuffix(oldName, ext)

	var newName string
	var err error

	if vlcPattern && strings.HasPrefix(baseName, "vlcsnap-") {
		newName, err = renameVlcScreenshot(baseName, ext)
	} else if winPattern && strings.HasPrefix(baseName, "スクリーンショット ") {
		newName, err = renameWindowsScreenshot(baseName, ext)
	} else {
		// パターンに一致しないファイルはスキップ
		return
	}

	if err != nil {
		fmt.Fprintf(stderr, "エラー: %s の解析に失敗しました: %v\n", oldPath, err)
		mu.Lock()
		*errorCount++
		mu.Unlock()
		return
	}

	newPath := filepath.Join(dir, newName)
	fmt.Fprintf(stdout, "処理中: %s -> %s\n", oldPath, newPath)

	err = os.Rename(oldPath, newPath)
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

// renameVlcScreenshot はVLCスクリーンショットファイルをリネームします
func renameVlcScreenshot(baseName, ext string) (string, error) {
	// vlcsnap-YYYY-MM-DD-HH-MM-SS
	re := regexp.MustCompile(`vlcsnap-(\d{4})-(\d{2})-(\d{2})-(\d{2})-(\d{2})-(\d{2})`)
	matches := re.FindStringSubmatch(baseName)
	if len(matches) != 7 {
		return "", fmt.Errorf("VLCスクリーンショットのパターンに一致しません: %s", baseName)
	}

	year, month, day := matches[1], matches[2], matches[3]
	hour, minute, second := matches[4], matches[5], matches[6]

	return fmt.Sprintf("Screenshot_%s%s%s-%s%s%s%s", year, month, day, hour, minute, second, ext), nil
}

// renameWindowsScreenshot はWindowsスクリーンショットファイルをリネームします
func renameWindowsScreenshot(baseName, ext string) (string, error) {
	// スクリーンショット YYYY-MM-DD HH-MM-SS
	re := regexp.MustCompile(`スクリーンショット (\d{4})-(\d{2})-(\d{2}) (\d{2})(\d{2})(\d{2})`)
	matches := re.FindStringSubmatch(baseName)
	if len(matches) != 7 {
		return "", fmt.Errorf("Windowsスクリーンショットのパターンに一致しません: %s", baseName)
	}

	year, month, day := matches[1], matches[2], matches[3]
	hour, minute, second := matches[4], matches[5], matches[6]

	return fmt.Sprintf("Screenshot_%s%s%s-%s%s%s%s", year, month, day, hour, minute, second, ext), nil
}

func isImageExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".avif":
		return true
	default:
		return false
	}
}
