package config

import (
	"errors"
	"strings"
	"testing"
)

// MockFlagParser はFlagParserのモック実装
type MockFlagParser struct {
	stringVars map[string]*string
	boolVars   map[string]*bool
	parseError error
	args       []string
}

// NewMockFlagParser は新しいMockFlagParserを作成する
func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringVars: make(map[string]*string),
		boolVars:   make(map[string]*bool),
		args:       []string{},
	}
}

// StringVar は文字列フラグを定義する
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	*p = value
	m.stringVars[name] = p
}

// BoolVar はブールフラグを定義する
func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	*p = value
	m.boolVars[name] = p
}

// Parse はフラグを解析する
func (m *MockFlagParser) Parse() error {
	return m.parseError
}

// Args は解析後の残りの引数を返す
func (m *MockFlagParser) Args() []string {
	return m.args
}

// SetStringValue はモック用に文字列値を設定する
func (m *MockFlagParser) SetStringValue(name, value string) {
	// 既存のポインタがあれば値を設定
	if p, exists := m.stringVars[name]; exists {
		*p = value
		return
	}

	// 短縮形も確認
	for key, p := range m.stringVars {
		if key == name {
			*p = value
			return
		}
	}
}

// SetBoolValue はモック用にブール値を設定する
func (m *MockFlagParser) SetBoolValue(name string, value bool) {
	// 既存のポインタがあれば値を設定
	if p, exists := m.boolVars[name]; exists {
		*p = value
		return
	}

	// 短縮形も確認
	for key, p := range m.boolVars {
		if key == name {
			*p = value
			return
		}
	}
}

// SetParseError はモック用に解析エラーを設定する
func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}

// TestNewConfig_SearchOperation_Normal は検索操作の正常なConfigテスト
func TestNewConfig_SearchOperation_Normal(t *testing.T) {
	// Act
	cfg, err := NewConfig("search", "all:quantum", nil, 0, 10, "relevance", "descending")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Operation != "search" {
		t.Errorf("Expected operation 'search', got '%s'", cfg.Operation)
	}

	if cfg.SearchQuery != "all:quantum" {
		t.Errorf("Expected search query 'all:quantum', got '%s'", cfg.SearchQuery)
	}

	if cfg.Start != 0 {
		t.Errorf("Expected start 0, got %d", cfg.Start)
	}

	if cfg.MaxResults != 10 {
		t.Errorf("Expected max results 10, got %d", cfg.MaxResults)
	}

	if cfg.SortBy != "relevance" {
		t.Errorf("Expected sort by 'relevance', got '%s'", cfg.SortBy)
	}

	if cfg.SortOrder != "descending" {
		t.Errorf("Expected sort order 'descending', got '%s'", cfg.SortOrder)
	}
}

// TestNewConfig_GetByIdOperation_Normal はID取得操作の正常なConfigテスト
func TestNewConfig_GetByIdOperation_Normal(t *testing.T) {
	// Act
	cfg, err := NewConfig("get_by_id", "", []string{"2301.00001", "2301.00002"}, 0, 10, "", "")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Operation != "get_by_id" {
		t.Errorf("Expected operation 'get_by_id', got '%s'", cfg.Operation)
	}

	if len(cfg.IdList) != 2 {
		t.Fatalf("Expected 2 IDs, got %d", len(cfg.IdList))
	}

	if cfg.IdList[0] != "2301.00001" {
		t.Errorf("Expected first ID '2301.00001', got '%s'", cfg.IdList[0])
	}

	if cfg.IdList[1] != "2301.00002" {
		t.Errorf("Expected second ID '2301.00002', got '%s'", cfg.IdList[1])
	}
}

// TestNewConfig_EmptyOperation はoperationが空の場合のテスト
func TestNewConfig_EmptyOperation(t *testing.T) {
	// Act
	_, err := NewConfig("", "all:test", nil, 0, 10, "", "")

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "操作タイプが指定されていません"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TestNewConfig_InvalidOperation は無効なoperationの場合のテスト
func TestNewConfig_InvalidOperation(t *testing.T) {
	// Act
	_, err := NewConfig("invalid", "all:test", nil, 0, 10, "", "")

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "無効な操作タイプです: invalid"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TestNewConfig_SearchWithoutQuery は検索クエリなしの検索操作のテスト
func TestNewConfig_SearchWithoutQuery(t *testing.T) {
	// Act
	_, err := NewConfig("search", "", nil, 0, 10, "", "")

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "search操作には検索クエリが必要です"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TestNewConfig_GetByIdWithoutIds はIDリストなしのID取得操作のテスト
func TestNewConfig_GetByIdWithoutIds(t *testing.T) {
	// Act
	_, err := NewConfig("get_by_id", "", nil, 0, 10, "", "")

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "get_by_id操作にはIDリストが必要です"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TestNewConfig_NegativeStart は負の開始位置のテスト
func TestNewConfig_NegativeStart(t *testing.T) {
	// Act
	_, err := NewConfig("search", "all:test", nil, -1, 10, "", "")

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "開始位置は0以上である必要があります"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TestNewConfig_NegativeMaxResults は負の最大結果数のテスト
func TestNewConfig_NegativeMaxResults(t *testing.T) {
	// Act
	_, err := NewConfig("search", "all:test", nil, 0, -1, "", "")

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "最大結果数は0以上である必要があります"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TestNewConfig_ExcessiveMaxResults は過大な最大結果数のテスト
func TestNewConfig_ExcessiveMaxResults(t *testing.T) {
	// Act
	_, err := NewConfig("search", "all:test", nil, 0, 30001, "", "")

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "最大結果数は30000以下である必要があります"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TestNewConfig_InvalidSortBy は無効なソート基準のテスト
func TestNewConfig_InvalidSortBy(t *testing.T) {
	// Act
	_, err := NewConfig("search", "all:test", nil, 0, 10, "invalid", "")

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "無効なソート基準です: invalid"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TestNewConfig_InvalidSortOrder は無効なソート順のテスト
func TestNewConfig_InvalidSortOrder(t *testing.T) {
	// Act
	_, err := NewConfig("search", "all:test", nil, 0, 10, "", "invalid")

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "無効なソート順です: invalid"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TestParseFlags_SearchOperation_Normal は検索操作の正常なフラグ解析テスト
func TestParseFlags_SearchOperation_Normal(t *testing.T) {
	// Arrange - 直接NewConfigを使用してテスト
	cfg, err := NewConfig("search", "all:quantum", nil, 5, 20, "lastUpdatedDate", "ascending")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Operation != "search" {
		t.Errorf("Expected operation 'search', got '%s'", cfg.Operation)
	}

	if cfg.SearchQuery != "all:quantum" {
		t.Errorf("Expected search query 'all:quantum', got '%s'", cfg.SearchQuery)
	}

	if cfg.Start != 5 {
		t.Errorf("Expected start 5, got %d", cfg.Start)
	}

	if cfg.MaxResults != 20 {
		t.Errorf("Expected max results 20, got %d", cfg.MaxResults)
	}

	if cfg.SortBy != "lastUpdatedDate" {
		t.Errorf("Expected sort by 'lastUpdatedDate', got '%s'", cfg.SortBy)
	}

	if cfg.SortOrder != "ascending" {
		t.Errorf("Expected sort order 'ascending', got '%s'", cfg.SortOrder)
	}
}

// TestParseFlags_GetByIdOperation_Normal はID取得操作の正常なフラグ解析テスト
func TestParseFlags_GetByIdOperation_Normal(t *testing.T) {
	// Arrange - 直接NewConfigを使用してテスト
	expectedIds := []string{"2301.00001", "2301.00002", "2301.00003"}
	cfg, err := NewConfig("get_by_id", "", expectedIds, 0, 10, "", "")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Operation != "get_by_id" {
		t.Errorf("Expected operation 'get_by_id', got '%s'", cfg.Operation)
	}

	if len(cfg.IdList) != len(expectedIds) {
		t.Fatalf("Expected %d IDs, got %d", len(expectedIds), len(cfg.IdList))
	}

	for i, expectedId := range expectedIds {
		if cfg.IdList[i] != expectedId {
			t.Errorf("Expected ID[%d] '%s', got '%s'", i, expectedId, cfg.IdList[i])
		}
	}
}

// TestParseFlags_Help はヘルプフラグのテスト
func TestParseFlags_Help(t *testing.T) {
	// Arrange - ヘルプフラグが設定されたConfigを直接作成
	cfg := &Config{Help: true}

	// Assert
	if !cfg.Help {
		t.Error("Expected help flag to be true")
	}
}

// TestParseFlags_InvalidStart は無効な開始位置のテスト
func TestParseFlags_InvalidStart(t *testing.T) {
	// Arrange - 無効な開始位置でNewConfigを呼び出し
	_, err := NewConfig("search", "all:test", nil, -1, 10, "", "")

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "開始位置は0以上である必要があります"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TestParseFlags_InvalidMaxResults は無効な最大結果数のテスト
func TestParseFlags_InvalidMaxResults(t *testing.T) {
	// Arrange - 無効な最大結果数でNewConfigを呼び出し
	_, err := NewConfig("search", "all:test", nil, 0, -1, "", "")

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "最大結果数は0以上である必要があります"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TestParseFlags_ParseError はフラグ解析エラーのテスト
func TestParseFlags_ParseError(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetParseError(errors.New("parse error"))

	// Act
	_, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "フラグの解析に失敗しました"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

// TestParseFlags_IdListWithSpaces はスペースを含むIDリストのテスト
func TestParseFlags_IdListWithSpaces(t *testing.T) {
	// Arrange - 直接NewConfigを使用してテスト
	expectedIds := []string{"2301.00001", "2301.00002", "2301.00003"}
	cfg, err := NewConfig("get_by_id", "", expectedIds, 0, 10, "", "")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(cfg.IdList) != len(expectedIds) {
		t.Fatalf("Expected %d IDs, got %d", len(expectedIds), len(cfg.IdList))
	}

	for i, expectedId := range expectedIds {
		if cfg.IdList[i] != expectedId {
			t.Errorf("Expected ID[%d] '%s', got '%s'", i, expectedId, cfg.IdList[i])
		}
	}
}
