package repositories_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/landmaster135/devbox/internal/independencies/env_loader/interfaces/repositories"
)
func TestEnvRepositoryImpl_LoadEnvFromYaml(t *testing.T) {
	// テスト用のYAMLファイルを作成
	tempDir := t.TempDir()
	yamlPath := filepath.Join(tempDir, "test_env.yml")
	yamlContent := `
TEST_KEY1: test_value1
TEST_KEY2: test_value2
`
	err := os.WriteFile(yamlPath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("テスト用YAMLファイルの作成に失敗しました: %v", err)
	}

	// テスト対象のリポジトリを作成
	repo := repositories.NewEnvRepository()

	// YAMLファイルから環境変数を読み込む
	config, err := repo.LoadEnvFromYaml(yamlPath)
	if err != nil {
		t.Fatalf("YAMLファイルの読み込みに失敗しました: %v", err)
	}

	// 読み込んだ環境変数が正しいか確認
	variables := config.GetAllVariables()
	if len(variables) != 2 {
		t.Errorf("環境変数の数が正しくありません: expected=%d, got=%d", 2, len(variables))
	}
	if variables["TEST_KEY1"] != "test_value1" {
		t.Errorf("環境変数の値が正しくありません: expected=%s, got=%s", "test_value1", variables["TEST_KEY1"])
	}
	if variables["TEST_KEY2"] != "test_value2" {
		t.Errorf("環境変数の値が正しくありません: expected=%s, got=%s", "test_value2", variables["TEST_KEY2"])
	}
}

func TestEnvRepositoryImpl_SetEnv_GetEnv(t *testing.T) {
	// テスト対象のリポジトリを作成
	repo := repositories.NewEnvRepository()

	// テスト用の環境変数名
	testKey := "TEST_ENV_KEY"
	testValue := "test_env_value"

	// 環境変数を設定
	err := repo.SetEnv(testKey, testValue)
	if err != nil {
		t.Fatalf("環境変数の設定に失敗しました: %v", err)
	}

	// 環境変数を取得
	value := repo.GetEnv(testKey)
	if value != testValue {
		t.Errorf("環境変数の値が正しくありません: expected=%s, got=%s", testValue, value)
	}

	// テスト後のクリーンアップ
	os.Unsetenv(testKey)
}

func TestEnvRepositoryImpl_LoadEnvFromYaml_FileNotFound(t *testing.T) {
	// テスト対象のリポジトリを作成
	repo := repositories.NewEnvRepository()

	// 存在しないファイルを指定
	_, err := repo.LoadEnvFromYaml("non_existent_file.yml")
	if err == nil {
		t.Errorf("存在しないファイルを指定したにもかかわらず、エラーが発生しませんでした")
	}
}
