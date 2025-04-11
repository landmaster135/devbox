package repositories

import "github.com/landmaster135/devbox/internal/domain/models"

// EnvRepository は環境変数操作を行うインターフェースです
type EnvRepository interface {
	// LoadEnvFromYaml はYAMLファイルから環境変数を読み込みます
	LoadEnvFromYaml(path string) (*models.EnvConfig, error)

	// SetEnv は環境変数を設定します
	SetEnv(key, value string) error

	// GetEnv は環境変数の値を取得します
	GetEnv(key string) string
}
