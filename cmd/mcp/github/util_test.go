package github

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

// MockHTTPClient struct is a mock for HTTP client.
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do method executes HTTP request.
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// マップの比較用ヘルパー関数
func compareMaps(t *testing.T, expected, actual map[string]interface{}) bool {
	// キーの数が同じか確認
	if len(expected) != len(actual) {
		t.Errorf("マップのサイズが異なります。期待: %d, 実際: %d", len(expected), len(actual))
		return false
	}

	// 各キーと値を個別に比較
	for k, expectedVal := range expected {
		actualVal, exists := actual[k]
		if !exists {
			t.Errorf("キー %s が実際のマップに存在しません", k)
			return false
		}

		// 値の型を確認
		expectedType := reflect.TypeOf(expectedVal)
		actualType := reflect.TypeOf(actualVal)
		if expectedType != actualType {
			t.Errorf("キー %s の値の型が異なります。期待: %v, 実際: %v", k, expectedType, actualType)
			return false
		}

		// 値を文字列に変換して比較
		expectedStr := fmt.Sprintf("%v", expectedVal)
		actualStr := fmt.Sprintf("%v", actualVal)
		if expectedStr != actualStr {
			t.Errorf("キー %s の値が異なります。期待: %v, 実際: %v", k, expectedStr, actualStr)
			return false
		}
	}

	return true
}

func TestReturnJSONResult(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
	}{
		{
			name:        "正常系 - 有効なJSONデータ",
			input:       map[string]interface{}{"key": "value"},
			expectError: false,
		},
		{
			name: "異常系 - JSONにマーシャルできないデータ",
			input: func() interface{} {
				// JSONにマーシャルできない循環参照を持つデータ構造
				type Circular struct {
					Self *Circular
				}
				c := &Circular{}
				c.Self = c
				return c
			}(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := returnJSONResult(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("returnJSONResult() エラーが期待されていましたが、エラーは発生しませんでした")
				}
				// エラーが発生した場合、resultはnilであるべき
				if result != nil {
					t.Errorf("returnJSONResult() エラー時にnilではない結果が返されました: %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("returnJSONResult() error = %v", err)
					return
				}
				// 正常系の場合、resultはnilではないはず
				if result == nil {
					t.Errorf("returnJSONResult() 結果がnilです")
				}
			}
		})
	}
}
