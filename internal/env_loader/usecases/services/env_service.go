package services

import (
	"fmt"
	"path/filepath"

	"github.com/landmaster135/devbox/internal/env_loader/domain/models"
	"github.com/landmaster135/devbox/internal/env_loader/domain/repositories"
)

// EnvService は環境変数操作のユースケースを実装するサービスです
type EnvService struct {
	envRepo repositories.EnvRepository
}

// NewEnvService は新しいEnvServiceインスタンスを作成します
func NewEnvService(envRepo repositories.EnvRepository) *EnvService {
	return &EnvService{
		envRepo: envRepo,
	}
}

// LoadAndSetEnvFromYaml はYAMLファイルから環境変数を読み込み、設定します
func (s *EnvService) LoadAndSetEnvFromYaml(path string) error {
	// YAMLファイルから環境変数を読み込む
	config, err := s.envRepo.LoadEnvFromYaml(path)
	if err != nil {
		return fmt.Errorf("環境変数の読み込みに失敗しました: %w", err)
	}

	// 環境変数を設定
	return s.SetEnvFromConfig(config)
}

// SetEnvFromConfig はEnvConfigから環境変数を設定します
func (s *EnvService) SetEnvFromConfig(config *models.EnvConfig) error {
	for key, value := range config.GetAllVariables() {
		if err := s.envRepo.SetEnv(key, value); err != nil {
			return fmt.Errorf("環境変数の設定に失敗しました (%s=%s): %w", key, value, err)
		}
	}
	return nil
}

// GetEnv は環境変数の値を取得します
func (s *EnvService) GetEnv(key string) string {
	return s.envRepo.GetEnv(key)
}

// ResolveEnvFilePath は環境変数ファイルのパスを解決します
// pathが空の場合はデフォルトのenv.ymlを使用します
func (s *EnvService) ResolveEnvFilePath(path string) string {
	if path == "" {
		return "env.yml"
	}
	return filepath.Clean(path)
}
