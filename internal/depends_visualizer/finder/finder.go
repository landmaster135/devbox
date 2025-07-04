package finder

import (
	"os"
	"path/filepath"
)

// FindFiles は指定されたディレクトリから条件に合うファイルを検索します
func FindFiles(dirPath string, recursive bool, extensions []string) ([]string, error) {
	var files []string

	if recursive {
		// 再帰的にディレクトリを処理
		err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				return nil
			}
			ext := filepath.Ext(path)
			for _, e := range extensions {
				if ext == e {
					files = append(files, path)
					break
				}
			}
			return nil
		})

		if err != nil {
			return nil, err
		}
	} else {
		// ディレクトリ内のファイルのみ処理
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			ext := filepath.Ext(entry.Name())
			for _, e := range extensions {
				if ext == e {
					files = append(files, filepath.Join(dirPath, entry.Name()))
					break
				}
			}
		}
	}

	return files, nil
}
