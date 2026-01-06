package usecases

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// FileInfo はファイル情報を保持する構造体です
type FileInfo struct {
	Path string
	Name string
}

func isImageExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".avif", ".mp4":
		return true
	default:
		return false
	}
}

// findScreenshotFilesForDateTime は--to-datetimeフラグ用にスクリーンショットファイルを検索します
func findScreenshotFilesForDateTime(srcDir string, recursive bool, stdout, stderr io.Writer) ([]string, error) {
	// 既存の関数を使って基本パターンを検索
	files, err := findScreenshotFiles(srcDir, recursive, operationAll, stdout, stderr)
	if err != nil {
		return nil, err
	}

	// Screenshot_パターンを追加で検索
	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			ext := strings.ToLower(filepath.Ext(d.Name()))
			if isImageExt(ext) {
				name := d.Name()
				if strings.HasPrefix(name, "Screenshot_") {
					files = append(files, path)
				}
			}
		}
		return nil
	}

	if recursive {
		if err := filepath.WalkDir(srcDir, walkFunc); err != nil {
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
					if strings.HasPrefix(name, "Screenshot_") {
						files = append(files, filepath.Join(srcDir, name))
					}
				}
			}
		}
	}

	return files, nil
}

// findScreenshotFiles は指定されたディレクトリからスクリーンショットファイルを検索します
func findScreenshotFiles(srcDir string, recursive bool, operation Operation, stdout, stderr io.Writer) ([]string, error) {
	var files []string

	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			name := d.Name()
			origExt := filepath.Ext(name)
			ext := strings.ToLower(origExt)
			if isImageExt(ext) {
				baseName := strings.TrimSuffix(name, origExt)
				matchesXiaomi := false
				if operation == OperationXiaomi || operation == operationAll {
					matchesXiaomi = xiaomiScreenshotRegexp.MatchString(baseName)
				}
				if matchesOperation(name, matchesXiaomi, operation) {
					files = append(files, path)
				}
			}
		}
		return nil
	}

	if recursive {
		if err := filepath.WalkDir(srcDir, walkFunc); err != nil {
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
				name := entry.Name()
				origExt := filepath.Ext(name)
				ext := strings.ToLower(origExt)
				if isImageExt(ext) {
					baseName := strings.TrimSuffix(name, origExt)
					matchesXiaomi := false
					if operation == OperationXiaomi || operation == operationAll {
						matchesXiaomi = xiaomiScreenshotRegexp.MatchString(baseName)
					}
					if matchesOperation(name, matchesXiaomi, operation) {
						files = append(files, filepath.Join(srcDir, name))
					}
				}
			}
		}
	}

	return files, nil
}

func matchesOperation(name string, matchesXiaomi bool, operation Operation) bool {
	switch operation {
	case OperationVLC:
		return strings.HasPrefix(name, "vlcsnap-")
	case OperationWin:
		return strings.HasPrefix(name, "スクリーンショット ")
	case OperationPixel:
		return strings.HasPrefix(name, "screen-")
	case OperationXiaomi:
		return matchesXiaomi
	case operationAll:
		return strings.HasPrefix(name, "vlcsnap-") ||
			strings.HasPrefix(name, "スクリーンショット ") ||
			strings.HasPrefix(name, "screen-") ||
			matchesXiaomi
	default:
		return false
	}
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
			Path: file,
			Name: info.Name(),
		}
	}

	if hasError {
		return fileInfos, fmt.Errorf("一部のファイル情報の取得に失敗しました")
	}
	return fileInfos, nil
}

// ProcessScreenshotRename はスクリーンショットファイルのリネーム処理を統合して実行します
func ProcessScreenshotRename(config Config, stdout, stderr io.Writer) (int, int, error) {
	// 設定の妥当性を検証
	if err := validateConfig(config, stderr); err != nil {
		return 0, 0, err
	}

	// スクリーンショットファイルの検索
	var files []string
	var err error
	if config.ToDateTime {
		files, err = findScreenshotFilesForDateTime(config.SrcDir, config.Recursive, stdout, stderr)
	} else {
		files, err = findScreenshotFiles(config.SrcDir, config.Recursive, config.Operation, stdout, stderr)
	}
	if err != nil {
		return 0, 0, err
	}

	if len(files) == 0 {
		fmt.Fprintln(stdout, "スクリーンショットファイルが見つかりませんでした。")
		return 0, 0, nil
	}

	fmt.Fprintf(stdout, "スクリーンショットファイルが %d 件見つかりました。\n", len(files))

	if config.ToDateTime {
		fmt.Fprintln(stdout, "YYYYMMDDHHMMSS形式でのリネームを実行します。")
	} else {
		switch config.Operation {
		case OperationVLC:
			fmt.Fprintln(stdout, "VLCスナップショットパターンを使用します。")
		case OperationWin:
			fmt.Fprintln(stdout, "Windowsスクリーンショットパターンを使用します。")
		case OperationPixel:
			fmt.Fprintln(stdout, "Pixelスクリーンレコードパターンを使用します。")
		case OperationXiaomi:
			fmt.Fprintln(stdout, "Xiaomiスクリーンショットパターンを使用します。")
		}
	}

	// ファイル情報の取得
	fileInfos, err := getFileInfos(files, stderr)
	if err != nil {
		// エラーがあっても続行するため、ここではエラーコードを返さない
	}

	// リネーム処理の実行
	successCount, errorCount := renameScreenshotFiles(fileInfos, config, stdout, stderr)

	// 処理結果の出力
	fmt.Fprintf(stdout, "✔ ファイルリネームが完了しました\n")
	fmt.Fprintf(stdout, "  成功: %d ファイル\n", successCount)

	if errorCount > 0 {
		fmt.Fprintf(stdout, "  失敗: %d ファイル\n", errorCount)
		return successCount, errorCount, fmt.Errorf("一部のファイルのリネームに失敗しました")
	}

	return successCount, errorCount, nil
}
