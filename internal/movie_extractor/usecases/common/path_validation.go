package common

import (
	"fmt"
	"os"
	"path/filepath"
)

func PrepareOutputDir(outDir string) (string, error) {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return "", fmt.Errorf("出力ディレクトリの作成に失敗しました: %w", err)
	}

	absOutDir, err := filepath.Abs(outDir)
	if err != nil {
		return "", fmt.Errorf("出力ディレクトリの絶対パス変換に失敗しました: %w", err)
	}

	return absOutDir, nil
}
