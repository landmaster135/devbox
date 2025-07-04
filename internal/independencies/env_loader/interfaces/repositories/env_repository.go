package repositories

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/landmaster135/devbox/internal/independencies/env_loader/domain/models"
	domainRepo "github.com/landmaster135/devbox/internal/independencies/env_loader/domain/repositories"
)

// EnvRepositoryImpl はEnvRepositoryインターフェースの実装です
type EnvRepositoryImpl struct{}

// NewEnvRepository は新しいEnvRepositoryImplインスタンスを作成します
func NewEnvRepository() domainRepo.EnvRepository {
	return &EnvRepositoryImpl{}
}

// LoadEnvFromYaml はYAMLファイルから環境変数を読み込みます
func (r *EnvRepositoryImpl) LoadEnvFromYaml(path string) (*models.EnvConfig, error) {
	// ファイルを開く
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("YAMLファイルを開けませんでした: %w", err)
	}
	defer file.Close()

	// YAMLをデコード
	var yamlData map[string]string
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&yamlData); err != nil {
		return nil, fmt.Errorf("YAMLのデコードに失敗しました: %w", err)
	}

	// EnvConfigを作成
	config := models.NewEnvConfig()
	for key, value := range yamlData {
		config.AddVariable(key, value)
	}

	return config, nil
}

// SetEnv は環境変数を設定します
func (r *EnvRepositoryImpl) SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

// GetEnv は環境変数の値を取得します
func (r *EnvRepositoryImpl) GetEnv(key string) string {
	return os.Getenv(key)
}
