package config

import (
	"testing"
)

// MockFlagParser はテスト用のFlagParserモック
type MockFlagParser struct {
	stringVars   map[string]*string
	intVars      map[string]*int
	boolVars     map[string]*bool
	stringValues map[string]string
	intValues    map[string]int
	boolValues   map[string]bool
	args         []string
	parseError   error
}

// NewMockFlagParser は新しいMockFlagParserを作成する
func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringVars:   make(map[string]*string),
		intVars:      make(map[string]*int),
		boolVars:     make(map[string]*bool),
		stringValues: make(map[string]string),
		intValues:    make(map[string]int),
		boolValues:   make(map[string]bool),
		args:         []string{},
	}
}

// SetStringValue は文字列フラグの値を事前設定する
func (m *MockFlagParser) SetStringValue(name, value string) {
	m.stringValues[name] = value
}

// SetIntValue は整数フラグの値を事前設定する
func (m *MockFlagParser) SetIntValue(name string, value int) {
	m.intValues[name] = value
}

// SetBoolValue はブールフラグの値を事前設定する
func (m *MockFlagParser) SetBoolValue(name string, value bool) {
	m.boolValues[name] = value
}

// SetArgs は引数を設定する
func (m *MockFlagParser) SetArgs(args []string) {
	m.args = args
}

// SetParseError はParseメソッドで返すエラーを設定する
func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}

// StringVar は文字列フラグを定義する
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	if presetValue, exists := m.stringValues[name]; exists {
		*p = presetValue
	} else {
		*p = value
	}
	m.stringVars[name] = p
}

// IntVar は整数フラグを定義する
func (m *MockFlagParser) IntVar(p *int, name string, value int, usage string) {
	if presetValue, exists := m.intValues[name]; exists {
		*p = presetValue
	} else {
		*p = value
	}
	m.intVars[name] = p
}

// BoolVar はブールフラグを定義する
func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if presetValue, exists := m.boolValues[name]; exists {
		*p = presetValue
	} else {
		*p = value
	}
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

// TestNewConfig_Normal はNewConfigメソッドの正常系テスト
func TestNewConfig_Normal(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		username  string
		userID    *int
		format    string
		limit     int
		status    string
		outputDir string
	}{
		{
			name:      "UsernameOnly_Normal",
			operation: "query-anime",
			username:  "testuser",
			userID:    nil,
			format:    "json",
			limit:     0,
			status:    "",
			outputDir: "",
		},
		{
			name:      "UserIDOnly_Normal",
			operation: "query-anime",
			username:  "",
			userID:    func() *int { id := 12345; return &id }(),
			format:    "table",
			limit:     10,
			status:    "COMPLETED",
			outputDir: "/output",
		},
		{
			name:      "BothUserInfo_Normal",
			operation: "query-anime",
			username:  "testuser",
			userID:    func() *int { id := 12345; return &id }(),
			format:    "",
			limit:     0,
			status:    "CURRENT",
			outputDir: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := NewConfig(tt.operation, tt.username, tt.userID, tt.format, tt.limit, tt.status, tt.outputDir)

			// Assert
			if err != nil {
				t.Errorf("エラーが発生しました: %v", err)
				return
			}
			if result == nil {
				t.Error("結果がnilです")
				return
			}
			if result.Operation != tt.operation {
				t.Errorf("期待される操作: %s, 実際: %s", tt.operation, result.Operation)
			}
			if result.Username != tt.username {
				t.Errorf("期待されるユーザー名: %s, 実際: %s", tt.username, result.Username)
			}
			if (result.UserID == nil) != (tt.userID == nil) {
				t.Errorf("ユーザーIDのnil状態が異なります")
			} else if result.UserID != nil && tt.userID != nil && *result.UserID != *tt.userID {
				t.Errorf("期待されるユーザーID: %d, 実際: %d", *tt.userID, *result.UserID)
			}
			if result.Format != tt.format {
				t.Errorf("期待される形式: %s, 実際: %s", tt.format, result.Format)
			}
			if result.Limit != tt.limit {
				t.Errorf("期待される制限: %d, 実際: %d", tt.limit, result.Limit)
			}
			if result.Status != tt.status {
				t.Errorf("期待されるステータス: %s, 実際: %s", tt.status, result.Status)
			}
			if result.OutputDir != tt.outputDir {
				t.Errorf("期待される出力ディレクトリ: %s, 実際: %s", tt.outputDir, result.OutputDir)
			}
		})
	}
}

// TestNewConfig_EmptyOperation はNewConfigメソッドの操作タイプ空エラーテスト
func TestNewConfig_EmptyOperation(t *testing.T) {
	// Act
	result, err := NewConfig("", "testuser", nil, "json", 0, "", "")

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if result != nil {
		t.Error("エラー時は結果がnilである必要があります")
	}
	expectedError := "操作タイプが指定されていません"
	if err.Error() != expectedError {
		t.Errorf("期待されるエラーメッセージ: %s, 実際: %s", expectedError, err.Error())
	}
}

// TestNewConfig_InvalidOperation はNewConfigメソッドの無効操作タイプエラーテスト
func TestNewConfig_InvalidOperation(t *testing.T) {
	// Act
	result, err := NewConfig("invalid-operation", "testuser", nil, "json", 0, "", "")

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if result != nil {
		t.Error("エラー時は結果がnilである必要があります")
	}
	expectedError := "無効な操作タイプです: invalid-operation"
	if err.Error() != expectedError {
		t.Errorf("期待されるエラーメッセージ: %s, 実際: %s", expectedError, err.Error())
	}
}

// TestNewConfig_NoUserInfo はNewConfigメソッドのユーザー情報なしエラーテスト
func TestNewConfig_NoUserInfo(t *testing.T) {
	// Act
	result, err := NewConfig("query-anime", "", nil, "json", 0, "", "")

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if result != nil {
		t.Error("エラー時は結果がnilである必要があります")
	}
	expectedError := "ユーザー名またはユーザーIDのいずれかを指定してください"
	if err.Error() != expectedError {
		t.Errorf("期待されるエラーメッセージ: %s, 実際: %s", expectedError, err.Error())
	}
}

// TestNewConfig_InvalidFormat はNewConfigメソッドの無効形式エラーテスト
func TestNewConfig_InvalidFormat(t *testing.T) {
	// Act
	result, err := NewConfig("query-anime", "testuser", nil, "invalid-format", 0, "", "")

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if result != nil {
		t.Error("エラー時は結果がnilである必要があります")
	}
	expectedError := "無効な出力形式です: invalid-format"
	if err.Error() != expectedError {
		t.Errorf("期待されるエラーメッセージ: %s, 実際: %s", expectedError, err.Error())
	}
}

// TestNewConfig_NegativeLimit はNewConfigメソッドの負の制限値エラーテスト
func TestNewConfig_NegativeLimit(t *testing.T) {
	// Act
	result, err := NewConfig("query-anime", "testuser", nil, "json", -1, "", "")

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if result != nil {
		t.Error("エラー時は結果がnilである必要があります")
	}
	expectedError := "制限値は0以上である必要があります"
	if err.Error() != expectedError {
		t.Errorf("期待されるエラーメッセージ: %s, 実際: %s", expectedError, err.Error())
	}
}

// TestNewConfig_InvalidStatus はNewConfigメソッドの無効ステータスエラーテスト
func TestNewConfig_InvalidStatus(t *testing.T) {
	// Act
	result, err := NewConfig("query-anime", "testuser", nil, "json", 0, "INVALID_STATUS", "")

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if result != nil {
		t.Error("エラー時は結果がnilである必要があります")
	}
	expectedError := "無効なステータスです: INVALID_STATUS"
	if err.Error() != expectedError {
		t.Errorf("期待されるエラーメッセージ: %s, 実際: %s", expectedError, err.Error())
	}
}

// TestNewConfig_ValidStatuses はNewConfigメソッドの有効ステータステスト
func TestNewConfig_ValidStatuses(t *testing.T) {
	validStatuses := []string{"CURRENT", "PLANNING", "COMPLETED", "DROPPED", "PAUSED", "REPEATING"}

	for _, status := range validStatuses {
		t.Run("Status_"+status, func(t *testing.T) {
			// Act
			result, err := NewConfig("query-anime", "testuser", nil, "json", 0, status, "")

			// Assert
			if err != nil {
				t.Errorf("エラーが発生しました: %v", err)
			}
			if result == nil {
				t.Error("結果がnilです")
				return
			}
			if result.Status != status {
				t.Errorf("期待されるステータス: %s, 実際: %s", status, result.Status)
			}
		})
	}
}

// TestParseFlagsWithParser_Normal はParseFlagsWithParserメソッドの正常系テスト
func TestParseFlagsWithParser_Normal(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringValue("operation", "query-anime")
	mockParser.SetStringValue("username", "testuser")
	mockParser.SetStringValue("format", "json")
	mockParser.SetStringValue("limit", "10")
	mockParser.SetStringValue("status", "COMPLETED")
	mockParser.SetStringValue("output-dir", "/output")

	// Act
	result, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
		return
	}
	if result == nil {
		t.Error("結果がnilです")
		return
	}
	if result.Operation != "query-anime" {
		t.Errorf("期待される操作: query-anime, 実際: %s", result.Operation)
	}
	if result.Username != "testuser" {
		t.Errorf("期待されるユーザー名: testuser, 実際: %s", result.Username)
	}
	if result.Format != "json" {
		t.Errorf("期待される形式: json, 実際: %s", result.Format)
	}
	if result.Limit != 10 {
		t.Errorf("期待される制限: 10, 実際: %d", result.Limit)
	}
	if result.Status != "COMPLETED" {
		t.Errorf("期待されるステータス: COMPLETED, 実際: %s", result.Status)
	}
	if result.OutputDir != "/output" {
		t.Errorf("期待される出力ディレクトリ: /output, 実際: %s", result.OutputDir)
	}
}

// TestParseFlagsWithParser_Help はParseFlagsWithParserメソッドのヘルプテスト
func TestParseFlagsWithParser_Help(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetBoolValue("help", true)

	// Act
	result, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
		return
	}
	if result == nil {
		t.Error("結果がnilです")
		return
	}
	if !result.Help {
		t.Error("ヘルプフラグがtrueである必要があります")
	}
}

// TestParseFlagsWithParser_UserID はParseFlagsWithParserメソッドのユーザーIDテスト
func TestParseFlagsWithParser_UserID(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringValue("operation", "query-anime")
	mockParser.SetStringValue("user-id", "12345")

	// Act
	result, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
		return
	}
	if result == nil {
		t.Error("結果がnilです")
		return
	}
	if result.UserID == nil {
		t.Error("ユーザーIDがnilです")
		return
	}
	if *result.UserID != 12345 {
		t.Errorf("期待されるユーザーID: 12345, 実際: %d", *result.UserID)
	}
}

// TestParseFlagsWithParser_InvalidUserID はParseFlagsWithParserメソッドの無効ユーザーIDエラーテスト
func TestParseFlagsWithParser_InvalidUserID(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringValue("operation", "query-anime")
	mockParser.SetStringValue("user-id", "invalid")

	// Act
	result, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if result != nil {
		t.Error("エラー時は結果がnilである必要があります")
	}
	expectedError := "無効なユーザーIDです: invalid"
	if err.Error() != expectedError {
		t.Errorf("期待されるエラーメッセージ: %s, 実際: %s", expectedError, err.Error())
	}
}

// TestParseFlagsWithParser_InvalidLimit はParseFlagsWithParserメソッドの無効制限値エラーテスト
func TestParseFlagsWithParser_InvalidLimit(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringValue("operation", "query-anime")
	mockParser.SetStringValue("username", "testuser")
	mockParser.SetStringValue("limit", "invalid")

	// Act
	result, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if result != nil {
		t.Error("エラー時は結果がnilである必要があります")
	}
	expectedError := "無効な制限値です: invalid"
	if err.Error() != expectedError {
		t.Errorf("期待されるエラーメッセージ: %s, 実際: %s", expectedError, err.Error())
	}
}

// TestParseFlagsWithParser_PositionalArgs はParseFlagsWithParserメソッドの位置引数テスト
func TestParseFlagsWithParser_PositionalArgs(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringValue("username", "testuser")
	mockParser.SetArgs([]string{"query-anime"})

	// Act
	result, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
		return
	}
	if result == nil {
		t.Error("結果がnilです")
		return
	}
	if result.Operation != "query-anime" {
		t.Errorf("期待される操作: query-anime, 実際: %s", result.Operation)
	}
}

// TestParseFlagsWithParser_ShortFlags はParseFlagsWithParserメソッドの短縮フラグテスト
func TestParseFlagsWithParser_ShortFlags(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringValue("o", "query-anime")
	mockParser.SetStringValue("u", "testuser")
	mockParser.SetStringValue("f", "table")
	mockParser.SetStringValue("l", "5")
	mockParser.SetStringValue("s", "CURRENT")
	mockParser.SetStringValue("d", "/tmp")
	mockParser.SetBoolValue("h", false)

	// Act
	result, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
		return
	}
	if result == nil {
		t.Error("結果がnilです")
		return
	}
	if result.Operation != "query-anime" {
		t.Errorf("期待される操作: query-anime, 実際: %s", result.Operation)
	}
	if result.Username != "testuser" {
		t.Errorf("期待されるユーザー名: testuser, 実際: %s", result.Username)
	}
	if result.Format != "table" {
		t.Errorf("期待される形式: table, 実際: %s", result.Format)
	}
	if result.Limit != 5 {
		t.Errorf("期待される制限: 5, 実際: %d", result.Limit)
	}
	if result.Status != "CURRENT" {
		t.Errorf("期待されるステータス: CURRENT, 実際: %s", result.Status)
	}
	if result.OutputDir != "/tmp" {
		t.Errorf("期待される出力ディレクトリ: /tmp, 実際: %s", result.OutputDir)
	}
}

// #==============================================================#
// ##          StandardFlagParser Tests                         ##
// #==============================================================#

// TestNewStandardFlagParser_Normal はNewStandardFlagParserメソッドの正常系テスト
func TestNewStandardFlagParser_Normal(t *testing.T) {
	// Act
	result := NewStandardFlagParser()

	// Assert
	if result == nil {
		t.Error("結果がnilです")
		return
	}
	if result.flagSet == nil {
		t.Error("flagSetがnilです")
	}
}

// TestStandardFlagParser_StringVar はStandardFlagParserのStringVarメソッドテスト
func TestStandardFlagParser_StringVar(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var testValue string

	// Act
	parser.StringVar(&testValue, "test-string-var", "default", "test usage")

	// Assert
	if testValue != "default" {
		t.Errorf("期待されるデフォルト値: default, 実際: %s", testValue)
	}
}

// TestStandardFlagParser_IntVar はStandardFlagParserのIntVarメソッドテスト
func TestStandardFlagParser_IntVar(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var testValue int

	// Act
	parser.IntVar(&testValue, "test-int-var", 42, "test usage")

	// Assert
	if testValue != 42 {
		t.Errorf("期待されるデフォルト値: 42, 実際: %d", testValue)
	}
}

// TestStandardFlagParser_BoolVar はStandardFlagParserのBoolVarメソッドテスト
func TestStandardFlagParser_BoolVar(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var testValue bool

	// Act
	parser.BoolVar(&testValue, "test-bool-var", true, "test usage")

	// Assert
	if testValue != true {
		t.Errorf("期待されるデフォルト値: true, 実際: %t", testValue)
	}
}

// TestStandardFlagParser_Args はStandardFlagParserのArgsメソッドテスト
func TestStandardFlagParser_Args(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()

	// Act
	result := parser.Args()

	// Assert
	if result == nil {
		t.Error("結果がnilです")
	}
}
