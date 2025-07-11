package usecases

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

// #==============================================================#
// ##          Mock Structures                                   ##
// #==============================================================#

// MockHTTPClient struct is a mock for HTTP client.
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do method executes HTTP request.
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// MockJSONMarshaler struct is a mock for JSON marshaler.
type MockJSONMarshaler struct {
	MarshalIndentFunc func(v interface{}, prefix, indent string) ([]byte, error)
}

// MarshalIndent method marshals JSON with indentation.
func (m *MockJSONMarshaler) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return m.MarshalIndentFunc(v, prefix, indent)
}

// #==============================================================#
// ##          Helper Functions                                  ##
// #==============================================================#

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

// #==============================================================#
// ##          Error Readers for Testing                        ##
// #==============================================================#

// エラーを返すリーダー
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("読み取りエラー")
}

// エラーを返すReadCloser
type errorReadCloser struct{}

func (e *errorReadCloser) Read(p []byte) (n int, err error) {
	return 0, errors.New("読み取りエラー")
}

func (e *errorReadCloser) Close() error {
	return nil
}
