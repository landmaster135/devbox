package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/landmaster135/devbox/internal/secret_detector/domain"
)

// LoadConfig はJSONファイルから設定を読み込む
func LoadConfig(filename string) (*domain.Config, error) {
	// ファイルが存在するかチェック
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filename)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	var config domain.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("invalid JSON format in %s: %w", filename, err)
	}

	return &config, nil
}
