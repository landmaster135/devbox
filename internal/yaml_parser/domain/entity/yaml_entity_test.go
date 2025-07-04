package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/landmaster135/devbox/internal/yaml_parser/domain/entity"
)

func TestNewYAMLData(t *testing.T) {
	// テストケース
	testCases := []struct {
		name     string
		data     interface{}
		expected interface{}
	}{
		{
			name:     "マップデータの場合",
			data:     map[string]interface{}{"key": "value"},
			expected: map[string]interface{}{"key": "value"},
		},
		{
			name:     "スライスデータの場合",
			data:     []interface{}{"item1", "item2"},
			expected: []interface{}{"item1", "item2"},
		},
		{
			name:     "nilデータの場合",
			data:     nil,
			expected: nil,
		},
	}

	// テスト実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// YAMLDataの作成
			yamlData := entity.NewYAMLData(tc.data)

			// アサーション
			assert.NotNil(t, yamlData)
			assert.Equal(t, tc.expected, yamlData.GetData())
		})
	}
}

func TestGetData(t *testing.T) {
	// テストデータ
	testData := map[string]interface{}{
		"string": "value",
		"number": 123,
		"bool":   true,
		"nested": map[string]interface{}{
			"key": "nested value",
		},
	}

	// YAMLDataの作成
	yamlData := entity.NewYAMLData(testData)

	// GetDataの実行
	result := yamlData.GetData()

	// アサーション
	assert.Equal(t, testData, result)
	assert.Equal(t, "value", result.(map[string]interface{})["string"])
	assert.Equal(t, 123, result.(map[string]interface{})["number"])
	assert.Equal(t, true, result.(map[string]interface{})["bool"])
	assert.Equal(t, "nested value", result.(map[string]interface{})["nested"].(map[string]interface{})["key"])
}
