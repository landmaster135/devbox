package models_test

import (
	"testing"

	"github.com/landmaster135/devbox/internal/domain/models"
)

func TestEnvConfig_AddVariable(t *testing.T) {
	// テスト用のEnvConfigを作成
	config := models.NewEnvConfig()

	// 環境変数を追加
	config.AddVariable("TEST_KEY", "test_value")

	// 追加した環境変数が正しく設定されているか確認
	value, exists := config.GetVariable("TEST_KEY")
	if !exists {
		t.Errorf("環境変数が追加されていません: %s", "TEST_KEY")
	}
	if value != "test_value" {
		t.Errorf("環境変数の値が正しくありません: expected=%s, got=%s", "test_value", value)
	}
}

func TestEnvConfig_GetVariable(t *testing.T) {
	// テスト用のEnvConfigを作成
	config := models.NewEnvConfig()
	config.AddVariable("TEST_KEY", "test_value")

	// 存在する環境変数を取得
	value, exists := config.GetVariable("TEST_KEY")
	if !exists {
		t.Errorf("環境変数が存在しません: %s", "TEST_KEY")
	}
	if value != "test_value" {
		t.Errorf("環境変数の値が正しくありません: expected=%s, got=%s", "test_value", value)
	}

	// 存在しない環境変数を取得
	_, exists = config.GetVariable("NON_EXISTENT_KEY")
	if exists {
		t.Errorf("存在しない環境変数が存在すると判定されました: %s", "NON_EXISTENT_KEY")
	}
}

func TestEnvConfig_GetAllVariables(t *testing.T) {
	// テスト用のEnvConfigを作成
	config := models.NewEnvConfig()
	config.AddVariable("KEY1", "value1")
	config.AddVariable("KEY2", "value2")

	// 全ての環境変数を取得
	variables := config.GetAllVariables()

	// 環境変数の数が正しいか確認
	if len(variables) != 2 {
		t.Errorf("環境変数の数が正しくありません: expected=%d, got=%d", 2, len(variables))
	}

	// 各環境変数の値が正しいか確認
	if variables["KEY1"] != "value1" {
		t.Errorf("環境変数の値が正しくありません: expected=%s, got=%s", "value1", variables["KEY1"])
	}
	if variables["KEY2"] != "value2" {
		t.Errorf("環境変数の値が正しくありません: expected=%s, got=%s", "value2", variables["KEY2"])
	}
}
