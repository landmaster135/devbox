package services

import (
	"errors"
	"testing"

	"github.com/landmaster135/devbox/internal/domain/models"
)

func TestEnvService_LoadAndSetEnvFromYaml(t *testing.T) {
	// テスト用のEnvConfigを作成
	config := models.NewEnvConfig()
	config.AddVariable("TEST_KEY1", "test_value1")
	config.AddVariable("TEST_KEY2", "test_value2")

	// モックリポジトリを作成
	mockRepo := &MockEnvRepository{
		LoadEnvFromYamlFunc: func(path string) (*models.EnvConfig, error) {
			if path != "test_env.yml" {
				return nil, errors.New("unexpected path")
			}
			return config, nil
		},
		SetEnvFunc: func(key, value string) error {
			return nil
		},
		GetEnvFunc: func(key string) string {
			return ""
		},
	}

	// テスト対象のサービスを作成
	service := NewEnvService(mockRepo)

	// YAMLファイルから環境変数を読み込み、設定
	err := service.LoadAndSetEnvFromYaml("test_env.yml")
	if err != nil {
		t.Fatalf("環境変数の読み込みと設定に失敗しました: %v", err)
	}
}

func TestEnvService_LoadAndSetEnvFromYaml_LoadError(t *testing.T) {
	// モックリポジトリを作成（読み込みエラーを返す）
	mockRepo := &MockEnvRepository{
		LoadEnvFromYamlFunc: func(path string) (*models.EnvConfig, error) {
			return nil, errors.New("load error")
		},
		SetEnvFunc: func(key, value string) error {
			return nil
		},
		GetEnvFunc: func(key string) string {
			return ""
		},
	}

	// テスト対象のサービスを作成
	service := NewEnvService(mockRepo)

	// YAMLファイルから環境変数を読み込み、設定
	err := service.LoadAndSetEnvFromYaml("test_env.yml")
	if err == nil {
		t.Errorf("エラーが発生するはずなのに、エラーが発生しませんでした")
	}
}

func TestEnvService_SetEnvFromConfig(t *testing.T) {
	// テスト用のEnvConfigを作成
	config := models.NewEnvConfig()
	config.AddVariable("TEST_KEY1", "test_value1")
	config.AddVariable("TEST_KEY2", "test_value2")

	// 設定された環境変数を追跡するためのマップ
	setEnvCalls := make(map[string]string)

	// モックリポジトリを作成
	mockRepo := &MockEnvRepository{
		LoadEnvFromYamlFunc: func(path string) (*models.EnvConfig, error) {
			return nil, nil
		},
		SetEnvFunc: func(key, value string) error {
			setEnvCalls[key] = value
			return nil
		},
		GetEnvFunc: func(key string) string {
			return ""
		},
	}

	// テスト対象のサービスを作成
	service := NewEnvService(mockRepo)

	// EnvConfigから環境変数を設定
	err := service.SetEnvFromConfig(config)
	if err != nil {
		t.Fatalf("環境変数の設定に失敗しました: %v", err)
	}

	// 環境変数が正しく設定されたか確認
	if len(setEnvCalls) != 2 {
		t.Errorf("設定された環境変数の数が正しくありません: expected=%d, got=%d", 2, len(setEnvCalls))
	}
	if setEnvCalls["TEST_KEY1"] != "test_value1" {
		t.Errorf("環境変数の値が正しくありません: expected=%s, got=%s", "test_value1", setEnvCalls["TEST_KEY1"])
	}
	if setEnvCalls["TEST_KEY2"] != "test_value2" {
		t.Errorf("環境変数の値が正しくありません: expected=%s, got=%s", "test_value2", setEnvCalls["TEST_KEY2"])
	}
}

func TestEnvService_SetEnvFromConfig_Error(t *testing.T) {
	// テスト用のEnvConfigを作成
	config := models.NewEnvConfig()
	config.AddVariable("TEST_KEY", "test_value")

	// モックリポジトリを作成（設定エラーを返す）
	mockRepo := &MockEnvRepository{
		LoadEnvFromYamlFunc: func(path string) (*models.EnvConfig, error) {
			return nil, nil
		},
		SetEnvFunc: func(key, value string) error {
			return errors.New("set error")
		},
		GetEnvFunc: func(key string) string {
			return ""
		},
	}

	// テスト対象のサービスを作成
	service := NewEnvService(mockRepo)

	// EnvConfigから環境変数を設定
	err := service.SetEnvFromConfig(config)
	if err == nil {
		t.Errorf("エラーが発生するはずなのに、エラーが発生しませんでした")
	}
}

func TestEnvService_GetEnv(t *testing.T) {
	// モックリポジトリを作成
	mockRepo := &MockEnvRepository{
		LoadEnvFromYamlFunc: func(path string) (*models.EnvConfig, error) {
			return nil, nil
		},
		SetEnvFunc: func(key, value string) error {
			return nil
		},
		GetEnvFunc: func(key string) string {
			if key == "TEST_KEY" {
				return "test_value"
			}
			return ""
		},
	}

	// テスト対象のサービスを作成
	service := NewEnvService(mockRepo)

	// 環境変数を取得
	value := service.GetEnv("TEST_KEY")
	if value != "test_value" {
		t.Errorf("環境変数の値が正しくありません: expected=%s, got=%s", "test_value", value)
	}
}

func TestEnvService_ResolveEnvFilePath(t *testing.T) {
	// モックリポジトリを作成
	mockRepo := &MockEnvRepository{
		LoadEnvFromYamlFunc: func(path string) (*models.EnvConfig, error) {
			return nil, nil
		},
		SetEnvFunc: func(key, value string) error {
			return nil
		},
		GetEnvFunc: func(key string) string {
			return ""
		},
	}

	// テスト対象のサービスを作成
	service := NewEnvService(mockRepo)

	// 空のパスを指定した場合
	path := service.ResolveEnvFilePath("")
	if path != "env.yml" {
		t.Errorf("デフォルトのパスが正しくありません: expected=%s, got=%s", "env.yml", path)
	}

	// パスを指定した場合
	path = service.ResolveEnvFilePath("custom_env.yml")
	if path != "custom_env.yml" {
		t.Errorf("指定したパスが正しくありません: expected=%s, got=%s", "custom_env.yml", path)
	}
}
