package valkey

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockLogger はLoggerのモック実装
type MockLogger struct {
	DebugFunc func(msg string, keysAndValues ...any)
	InfoFunc  func(msg string, keysAndValues ...any)
	WarnFunc  func(msg string, keysAndValues ...any)
	ErrorFunc func(msg string, err error, keysAndValues ...any)
	InitFunc  func(level int, format string) error
}

// Debug はデバッグメッセージを出力するモックメソッド
func (m *MockLogger) Debug(msg string, keysAndValues ...any) {
	if m.DebugFunc != nil {
		m.DebugFunc(msg, keysAndValues...)
	}
}

// Info は情報メッセージを出力するモックメソッド
func (m *MockLogger) Info(msg string, keysAndValues ...any) {
	if m.InfoFunc != nil {
		m.InfoFunc(msg, keysAndValues...)
	}
}

// Warn は警告メッセージを出力するモックメソッド
func (m *MockLogger) Warn(msg string, keysAndValues ...any) {
	if m.WarnFunc != nil {
		m.WarnFunc(msg, keysAndValues...)
	}
}

// Error はエラーメッセージを出力するモックメソッド
func (m *MockLogger) Error(msg string, err error, keysAndValues ...any) {
	if m.ErrorFunc != nil {
		m.ErrorFunc(msg, err, keysAndValues...)
	}
}

// Init はロガーを初期化するモックメソッド
func (m *MockLogger) Init(level int, format string) error {
	if m.InitFunc != nil {
		return m.InitFunc(level, format)
	}
	return nil
}

// MockDataRepository はDataRepositoryのモック実装
type MockDataRepository struct {
	GetKeysFunc          func(ctx context.Context, pattern string) ([]string, error)
	GetValueFunc         func(ctx context.Context, key string) (string, error)
	GetValueAsByteFunc   func(ctx context.Context, key string) ([]byte, error)
	GetTypeFunc          func(ctx context.Context, key string) (string, error)
	SetValueFunc         func(ctx context.Context, key string, valueJSON []byte) error
	StartServerFunc      func() error
	StopServerFunc       func() error
	IsServerRunningFunc  func() bool
	GetServerAddressFunc func() string
	DeleteKeyFunc        func(ctx context.Context, key string) (bool, error)
}

// GetKeys はパターンに一致するすべてのキーを取得するモックメソッド
func (m *MockDataRepository) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	return m.GetKeysFunc(ctx, pattern)
}

// GetValue はキーに対応する値を取得するモックメソッド
func (m *MockDataRepository) GetValue(ctx context.Context, key string) (string, error) {
	return m.GetValueFunc(ctx, key)
}

// GetValueAsByte はJSON形式のトークン情報を取得するモックメソッド
func (m *MockDataRepository) GetValueAsByte(ctx context.Context, key string) ([]byte, error) {
	return m.GetValueAsByteFunc(ctx, key)
}

// GetType はキーの型を取得するモックメソッド
func (m *MockDataRepository) GetType(ctx context.Context, key string) (string, error) {
	return m.GetTypeFunc(ctx, key)
}

// SetValue はJSON形式のトークン情報をValkeyに保存するモックメソッド
func (m *MockDataRepository) SetValue(ctx context.Context, key string, valueJSON []byte) error {
	return m.SetValueFunc(ctx, key, valueJSON)
}

// StartServer はDBサーバーを起動するモックメソッド
func (m *MockDataRepository) StartServer() error {
	return m.StartServerFunc()
}

// StopServer はDBサーバーを停止するモックメソッド
func (m *MockDataRepository) StopServer() error {
	return m.StopServerFunc()
}

// IsServerRunning はDBサーバーが起動しているかどうかを返すモックメソッド
func (m *MockDataRepository) IsServerRunning() bool {
	return m.IsServerRunningFunc()
}

// GetServerAddress はDBサーバーのアドレスを返すモックメソッド
func (m *MockDataRepository) GetServerAddress() string {
	return m.GetServerAddressFunc()
}

// DeleteKey はキーを削除するモックメソッド
func (m *MockDataRepository) DeleteKey(ctx context.Context, key string) (bool, error) {
	return m.DeleteKeyFunc(ctx, key)
}

// TestNewDataService_Normal は新しいDataServiceを作成するテスト
func TestNewDataService_Normal(t *testing.T) {
	// Arrange
	mockRepo := &MockDataRepository{}
	mockLogger := &MockLogger{}

	// Act
	service := NewDataService(mockRepo, mockLogger)

	// Assert
	assert.NotNil(t, service, "サービスがnilであってはならない")
	assert.Equal(t, mockRepo, service.repo, "リポジトリが正しく設定されていない")
	assert.Equal(t, mockLogger, service.logger, "ロガーが正しく設定されていない")
}

// TestGetRepository_Normal はリポジトリを取得するテスト
func TestGetRepository_Normal(t *testing.T) {
	// Arrange
	mockRepo := &MockDataRepository{}
	mockLogger := &MockLogger{}
	service := NewDataService(mockRepo, mockLogger)

	// Act
	repo := service.GetRepository()

	// Assert
	assert.Equal(t, mockRepo, repo, "取得したリポジトリが正しくない")
}

// TestGetKeys_Normal はパターンに一致するすべてのキーを取得するテスト
func TestGetKeys_Normal(t *testing.T) {
	// Arrange
	ctx := context.Background()
	pattern := "test*"
	expectedKeys := []string{"test1", "test2", "test3"}

	mockRepo := &MockDataRepository{
		GetKeysFunc: func(ctx context.Context, p string) ([]string, error) {
			assert.Equal(t, pattern, p, "パターンが正しく渡されていない")
			return expectedKeys, nil
		},
	}
	mockLogger := &MockLogger{}

	service := NewDataService(mockRepo, mockLogger)

	// Act
	keys, err := service.GetKeys(ctx, pattern)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.Equal(t, expectedKeys, keys, "取得したキーが正しくない")
}

// TestGetValue_Normal はキーに対応する値を取得するテスト
func TestGetValue_Normal(t *testing.T) {
	// Arrange
	ctx := context.Background()
	key := "testKey"
	expectedValue := "testValue"

	mockRepo := &MockDataRepository{
		GetValueFunc: func(ctx context.Context, k string) (string, error) {
			assert.Equal(t, key, k, "キーが正しく渡されていない")
			return expectedValue, nil
		},
	}
	mockLogger := &MockLogger{}

	service := NewDataService(mockRepo, mockLogger)

	// Act
	value, err := service.GetValue(ctx, key)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.Equal(t, expectedValue, value, "取得した値が正しくない")
}

// TestGetType_Normal はキーの型を取得するテスト
func TestGetType_Normal(t *testing.T) {
	// Arrange
	ctx := context.Background()
	key := "testKey"
	expectedType := "string"

	mockRepo := &MockDataRepository{
		GetTypeFunc: func(ctx context.Context, k string) (string, error) {
			assert.Equal(t, key, k, "キーが正しく渡されていない")
			return expectedType, nil
		},
	}
	mockLogger := &MockLogger{}

	service := NewDataService(mockRepo, mockLogger)

	// Act
	typeStr, err := service.GetType(ctx, key)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.Equal(t, expectedType, typeStr, "取得した型が正しくない")
}

// TestSetValue_Normal は値を設定するテスト
func TestSetValue_Normal(t *testing.T) {
	// Arrange
	ctx := context.Background()
	key := "testKey"
	valueJSON := []byte(`{"test": "value"}`)

	mockRepo := &MockDataRepository{
		SetValueFunc: func(ctx context.Context, k string, v []byte) error {
			assert.Equal(t, key, k, "キーが正しく渡されていない")
			assert.Equal(t, valueJSON, v, "値が正しく渡されていない")
			return nil
		},
	}
	mockLogger := &MockLogger{}

	service := NewDataService(mockRepo, mockLogger)

	// Act
	err := service.SetValue(ctx, key, valueJSON)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
}

// TestDeleteKey_Normal はキーを削除するテスト
func TestDeleteKey_Normal(t *testing.T) {
	// Arrange
	ctx := context.Background()
	key := "testKey"
	expectedResult := true

	mockRepo := &MockDataRepository{
		DeleteKeyFunc: func(ctx context.Context, k string) (bool, error) {
			assert.Equal(t, key, k, "キーが正しく渡されていない")
			return expectedResult, nil
		},
	}
	mockLogger := &MockLogger{}

	service := NewDataService(mockRepo, mockLogger)

	// Act
	result, err := service.DeleteKey(ctx, key)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.Equal(t, expectedResult, result, "削除結果が正しくない")
}

// TestDeleteKeys_Normal は複数のキーを削除するテスト
func TestDeleteKeys_Normal(t *testing.T) {
	// Arrange
	ctx := context.Background()
	keys := []string{"key1", "key2", "key3"}

	deleteResults := map[string]bool{
		"key1": true,
		"key2": true,
		"key3": false,
	}

	mockRepo := &MockDataRepository{
		DeleteKeyFunc: func(ctx context.Context, k string) (bool, error) {
			result, exists := deleteResults[k]
			assert.True(t, exists, "予期しないキーが渡された: "+k)
			return result, nil
		},
	}
	mockLogger := &MockLogger{}

	service := NewDataService(mockRepo, mockLogger)

	// Act
	results, err := service.DeleteKeys(ctx, keys)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.Equal(t, deleteResults, results, "削除結果が正しくない")
}

// TestDeleteKeys_Error は複数のキー削除中にエラーが発生した場合のテスト
func TestDeleteKeys_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	keys := []string{"key1", "key2", "key3"}
	expectedError := assert.AnError

	callCount := 0
	mockRepo := &MockDataRepository{
		DeleteKeyFunc: func(ctx context.Context, k string) (bool, error) {
			callCount++
			if callCount == 2 {
				return false, expectedError
			}
			return true, nil
		},
	}
	mockLogger := &MockLogger{}

	service := NewDataService(mockRepo, mockLogger)

	// Act
	results, err := service.DeleteKeys(ctx, keys)

	// Assert
	assert.Error(t, err, "エラーが発生するべき")
	assert.Equal(t, expectedError, err, "エラーが正しくない")

	// 最初のキーは削除成功、2番目でエラー発生
	expectedResults := map[string]bool{
		"key1": true,
	}
	assert.Equal(t, expectedResults, results, "エラー発生時の結果が正しくない")
}

// TestSelectValues_AllArgumentsEmpty は全ての引数がゼロ値の場合のエラーテスト
func TestSelectValues_AllArgumentsEmpty(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &MockDataRepository{}
	mockLogger := &MockLogger{}
	service := NewDataService(mockRepo, mockLogger)

	// Act
	result, err := service.SelectKeys(ctx, "", []string{}, "", false)

	// Assert
	assert.Error(t, err, "エラーが発生するべき")
	assert.Contains(t, err.Error(), "key、keys、pattern のいずれか、または all=true を指定してください", "エラーメッセージが正しくない")
	assert.Nil(t, result, "結果はnilであるべき")
}

// TestSelectValues_MultipleArgumentsProvided は複数の引数が同時に指定された場合のエラーテスト
func TestSelectValues_MultipleArgumentsProvided(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &MockDataRepository{}
	mockLogger := &MockLogger{}
	service := NewDataService(mockRepo, mockLogger)

	// Act
	result, err := service.SelectKeys(ctx, "testKey", []string{"key1"}, "", false)

	// Assert
	assert.Error(t, err, "エラーが発生するべき")
	assert.Contains(t, err.Error(), "key、keys、pattern、all のうち、同時に指定できるのは1つだけです", "エラーメッセージが正しくない")
	assert.Nil(t, result, "結果はnilであるべき")
}

// TestSelectValues_AllTrueWithOtherArguments はall=trueで他の引数が指定された場合のエラーテスト
func TestSelectValues_AllTrueWithOtherArguments(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &MockDataRepository{}
	mockLogger := &MockLogger{}
	service := NewDataService(mockRepo, mockLogger)

	// Act
	result, err := service.SelectKeys(ctx, "testKey", []string{}, "", true)

	// Assert
	assert.Error(t, err, "エラーが発生するべき")
	assert.Contains(t, err.Error(), "key、keys、pattern、all のうち、同時に指定できるのは1つだけです", "エラーメッセージが正しくない")
	assert.Nil(t, result, "結果はnilであるべき")
}

// TestSelectValues_KeyProvided_Normal はkey指定時の正常系テスト
func TestSelectValues_KeyProvided_Normal(t *testing.T) {
	// Arrange
	ctx := context.Background()
	key := "testKey"
	expectedValue := "testValue"
	expectedType := "string"

	mockRepo := &MockDataRepository{
		GetValueFunc: func(ctx context.Context, k string) (string, error) {
			assert.Equal(t, key, k, "キーが正しく渡されていない")
			return expectedValue, nil
		},
		GetTypeFunc: func(ctx context.Context, k string) (string, error) {
			assert.Equal(t, key, k, "キーが正しく渡されていない")
			return expectedType, nil
		},
	}
	mockLogger := &MockLogger{}

	service := NewDataService(mockRepo, mockLogger)

	// Act
	result, err := service.SelectKeys(ctx, key, []string{}, "", false)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.NotNil(t, result, "結果がnilであってはならない")

	resultMap, ok := result.(map[string]any)
	assert.True(t, ok, "結果がmap[string]any型でない")
	assert.Equal(t, key, resultMap["key"], "キーが正しくない")
	assert.Equal(t, expectedValue, resultMap["value"], "値が正しくない")
	assert.Equal(t, expectedType, resultMap["type"], "型が正しくない")
}

// TestSelectValues_KeyProvided_GetValueError はkey指定時のGetValueエラーテスト
func TestSelectValues_KeyProvided_GetValueError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	key := "testKey"
	expectedError := assert.AnError

	mockRepo := &MockDataRepository{
		GetValueFunc: func(ctx context.Context, k string) (string, error) {
			return "", expectedError
		},
	}
	mockLogger := &MockLogger{}

	service := NewDataService(mockRepo, mockLogger)

	// Act
	result, err := service.SelectKeys(ctx, key, []string{}, "", false)

	// Assert
	assert.Error(t, err, "エラーが発生するべき")
	assert.Contains(t, err.Error(), "値の取得に失敗しました", "エラーメッセージが正しくない")
	assert.Nil(t, result, "結果はnilであるべき")
}

// TestSelectValues_KeyProvided_GetTypeError はkey指定時のGetTypeエラーテスト
func TestSelectValues_KeyProvided_GetTypeError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	key := "testKey"
	expectedValue := "testValue"
	expectedError := assert.AnError

	mockRepo := &MockDataRepository{
		GetValueFunc: func(ctx context.Context, k string) (string, error) {
			return expectedValue, nil
		},
		GetTypeFunc: func(ctx context.Context, k string) (string, error) {
			return "", expectedError
		},
	}
	mockLogger := &MockLogger{}

	service := NewDataService(mockRepo, mockLogger)

	// Act
	result, err := service.SelectKeys(ctx, key, []string{}, "", false)

	// Assert
	assert.Error(t, err, "エラーが発生するべき")
	assert.Contains(t, err.Error(), "型の取得に失敗しました", "エラーメッセージが正しくない")
	assert.Nil(t, result, "結果はnilであるべき")
}

// TestSelectValues_KeysProvided_Normal はkeys指定時の正常系テスト
func TestSelectValues_KeysProvided_Normal(t *testing.T) {
	// Arrange
	ctx := context.Background()
	keys := []string{"pattern1*", "pattern2*"}
	expectedKeys1 := []string{"pattern1_key1", "pattern1_key2"}
	expectedKeys2 := []string{"pattern2_key1"}

	mockRepo := &MockDataRepository{
		GetKeysFunc: func(ctx context.Context, pattern string) ([]string, error) {
			switch pattern {
			case "pattern1*":
				return expectedKeys1, nil
			case "pattern2*":
				return expectedKeys2, nil
			default:
				return []string{}, nil
			}
		},
	}
	mockLogger := &MockLogger{}

	service := NewDataService(mockRepo, mockLogger)

	// Act
	result, err := service.SelectKeys(ctx, "", keys, "", false)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.NotNil(t, result, "結果がnilであってはならない")

	resultMap, ok := result.(map[string]any)
	assert.True(t, ok, "結果がmap[string]any型でない")

	retrievedKeys, ok := resultMap["keys"].([]string)
	assert.True(t, ok, "keysがstring配列でない")
	expectedAllKeys := append(expectedKeys1, expectedKeys2...)
	assert.Equal(t, expectedAllKeys, retrievedKeys, "取得したキーが正しくない")
	assert.Equal(t, len(expectedAllKeys), resultMap["count"], "カウントが正しくない")
}

// TestSelectValues_KeysProvided_Error はkeys指定時のエラーテスト
func TestSelectValues_KeysProvided_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	keys := []string{"pattern1*", "pattern2*"}
	expectedError := assert.AnError

	mockRepo := &MockDataRepository{
		GetKeysFunc: func(ctx context.Context, pattern string) ([]string, error) {
			if pattern == "pattern1*" {
				return []string{"key1"}, nil
			}
			return nil, expectedError
		},
	}
	mockLogger := &MockLogger{}

	service := NewDataService(mockRepo, mockLogger)

	// Act
	result, err := service.SelectKeys(ctx, "", keys, "", false)

	// Assert
	assert.Error(t, err, "エラーが発生するべき")
	assert.Contains(t, err.Error(), "に一致するキーの取得に失敗しました", "エラーメッセージが正しくない")
	assert.Nil(t, result, "結果はnilであるべき")
}

// TestSelectValues_PatternProvided_Normal はpattern指定時の正常系テスト
func TestSelectValues_PatternProvided_Normal(t *testing.T) {
	// Arrange
	ctx := context.Background()
	pattern := "test*"
	expectedKeys := []string{"test1", "test2", "test3"}

	mockRepo := &MockDataRepository{
		GetKeysFunc: func(ctx context.Context, p string) ([]string, error) {
			assert.Equal(t, pattern, p, "パターンが正しく渡されていない")
			return expectedKeys, nil
		},
	}
	mockLogger := &MockLogger{}

	service := NewDataService(mockRepo, mockLogger)

	// Act
	result, err := service.SelectKeys(ctx, "", []string{}, pattern, false)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.NotNil(t, result, "結果がnilであってはならない")

	resultMap, ok := result.(map[string]any)
	assert.True(t, ok, "結果がmap[string]any型でない")

	retrievedKeys, ok := resultMap["keys"].([]string)
	assert.True(t, ok, "keysがstring配列でない")
	assert.Equal(t, expectedKeys, retrievedKeys, "取得したキーが正しくない")
	assert.Equal(t, len(expectedKeys), resultMap["count"], "カウントが正しくない")
}

// TestSelectValues_AllTrue_Normal はall=true指定時の正常系テスト
func TestSelectValues_AllTrue_Normal(t *testing.T) {
	// Arrange
	ctx := context.Background()
	expectedKeys := []string{"key1", "key2", "key3"}

	mockRepo := &MockDataRepository{
		GetKeysFunc: func(ctx context.Context, pattern string) ([]string, error) {
			assert.Equal(t, "*", pattern, "パターンが正しく渡されていない")
			return expectedKeys, nil
		},
	}
	mockLogger := &MockLogger{}

	service := NewDataService(mockRepo, mockLogger)

	// Act
	result, err := service.SelectKeys(ctx, "", []string{}, "", true)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.NotNil(t, result, "結果がnilであってはならない")

	resultMap, ok := result.(map[string]any)
	assert.True(t, ok, "結果がmap[string]any型でない")

	retrievedKeys, ok := resultMap["keys"].([]string)
	assert.True(t, ok, "keysがstring配列でない")
	assert.Equal(t, expectedKeys, retrievedKeys, "取得したキーが正しくない")
	assert.Equal(t, len(expectedKeys), resultMap["count"], "カウントが正しくない")
}

// TestSelectValues_PatternProvided_Error はpattern指定時のエラーテスト
func TestSelectValues_PatternProvided_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	pattern := "test*"
	expectedError := assert.AnError

	mockRepo := &MockDataRepository{
		GetKeysFunc: func(ctx context.Context, p string) ([]string, error) {
			return nil, expectedError
		},
	}
	mockLogger := &MockLogger{}

	service := NewDataService(mockRepo, mockLogger)

	// Act
	result, err := service.SelectKeys(ctx, "", []string{}, pattern, false)

	// Assert
	assert.Error(t, err, "エラーが発生するべき")
	assert.Contains(t, err.Error(), "に一致するキーの取得に失敗しました", "エラーメッセージが正しくない")
	assert.Nil(t, result, "結果はnilであるべき")
}

// TestGetAllValues_EmptyKeys は空のキー配列の場合のテスト
func TestGetAllValues_EmptyKeys(t *testing.T) {
	// Arrange
	ctx := context.Background()
	keys := []string{}
	mockRepo := &MockDataRepository{}
	mockLogger := &MockLogger{}
	service := NewDataService(mockRepo, mockLogger)

	// Act
	result, err := service.GetAllValues(ctx, keys)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.NotNil(t, result, "結果がnilであってはならない")

	resultMap, ok := result.(map[string]any)
	assert.True(t, ok, "結果がmap[string]any型でない")

	values, ok := resultMap["values"].([]any)
	assert.True(t, ok, "valuesが配列でない")
	assert.Empty(t, values, "valuesが空でない")
	assert.Equal(t, 0, resultMap["count"], "カウントが正しくない")
}

// TestGetAllValues_Normal は正常系のテスト
func TestGetAllValues_Normal(t *testing.T) {
	// Arrange
	ctx := context.Background()
	keys := []string{"key1", "key2", "key3"}

	keyValueMap := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	keyTypeMap := map[string]string{
		"key1": "string",
		"key2": "string",
		"key3": "hash",
	}

	mockRepo := &MockDataRepository{
		GetValueFunc: func(ctx context.Context, key string) (string, error) {
			value, exists := keyValueMap[key]
			assert.True(t, exists, "予期しないキーが渡された: "+key)
			return value, nil
		},
		GetTypeFunc: func(ctx context.Context, key string) (string, error) {
			keyType, exists := keyTypeMap[key]
			assert.True(t, exists, "予期しないキーが渡された: "+key)
			return keyType, nil
		},
	}

	mockLogger := &MockLogger{}
	service := NewDataService(mockRepo, mockLogger)

	// Act
	result, err := service.GetAllValues(ctx, keys)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.NotNil(t, result, "結果がnilであってはならない")

	resultMap, ok := result.(map[string]any)
	assert.True(t, ok, "結果がmap[string]any型でない")

	values, ok := resultMap["values"].([]any)
	assert.True(t, ok, "valuesが配列でない")
	assert.Len(t, values, 3, "valuesの長さが正しくない")
	assert.Equal(t, 3, resultMap["count"], "カウントが正しくない")

	// 各値の検証
	for i, value := range values {
		keyValue, ok := value.(map[string]any)
		assert.True(t, ok, "値がmap[string]any型でない")

		key := keys[i]
		assert.Equal(t, key, keyValue["key"], "キーが正しくない")
		assert.Equal(t, keyValueMap[key], keyValue["value"], "値が正しくない")
		assert.Equal(t, keyTypeMap[key], keyValue["type"], "型が正しくない")
	}
}

// TestGetAllValues_GetValueError はGetValueエラー時の警告ログ出力テスト
func TestGetAllValues_GetValueError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	keys := []string{"key1", "key2", "key3"}
	expectedError := assert.AnError

	warnMessages := []string{}
	mockRepo := &MockDataRepository{
		GetValueFunc: func(ctx context.Context, key string) (string, error) {
			if key == "key2" {
				return "", expectedError
			}
			return "value", nil
		},
		GetTypeFunc: func(ctx context.Context, key string) (string, error) {
			return "string", nil
		},
	}

	mockLogger := &MockLogger{
		WarnFunc: func(msg string, keysAndValues ...any) {
			warnMessages = append(warnMessages, msg)
		},
	}

	service := NewDataService(mockRepo, mockLogger)

	// Act
	result, err := service.GetAllValues(ctx, keys)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.NotNil(t, result, "結果がnilであってはならない")

	resultMap, ok := result.(map[string]any)
	assert.True(t, ok, "結果がmap[string]any型でない")

	values, ok := resultMap["values"].([]any)
	assert.True(t, ok, "valuesが配列でない")
	assert.Len(t, values, 2, "valuesの長さが正しくない（エラーのキーは除外される）")
	assert.Equal(t, 2, resultMap["count"], "カウントが正しくない")

	// 警告ログが出力されたことを確認
	assert.Len(t, warnMessages, 1, "警告メッセージが出力されていない")
	assert.Contains(t, warnMessages[0], "key2", "警告メッセージにキー名が含まれていない")
	assert.Contains(t, warnMessages[0], "値の取得に失敗しました", "警告メッセージが正しくない")
}

// TestGetAllValues_GetTypeError はGetTypeエラー時の警告ログ出力テスト
func TestGetAllValues_GetTypeError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	keys := []string{"key1", "key2", "key3"}
	expectedError := assert.AnError

	warnMessages := []string{}
	mockRepo := &MockDataRepository{
		GetValueFunc: func(ctx context.Context, key string) (string, error) {
			return "value", nil
		},
		GetTypeFunc: func(ctx context.Context, key string) (string, error) {
			if key == "key3" {
				return "", expectedError
			}
			return "string", nil
		},
	}

	mockLogger := &MockLogger{
		WarnFunc: func(msg string, keysAndValues ...any) {
			warnMessages = append(warnMessages, msg)
		},
	}

	service := NewDataService(mockRepo, mockLogger)

	// Act
	result, err := service.GetAllValues(ctx, keys)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.NotNil(t, result, "結果がnilであってはならない")

	resultMap, ok := result.(map[string]any)
	assert.True(t, ok, "結果がmap[string]any型でない")

	values, ok := resultMap["values"].([]any)
	assert.True(t, ok, "valuesが配列でない")
	assert.Len(t, values, 2, "valuesの長さが正しくない（エラーのキーは除外される）")
	assert.Equal(t, 2, resultMap["count"], "カウントが正しくない")

	// 警告ログが出力されたことを確認
	assert.Len(t, warnMessages, 1, "警告メッセージが出力されていない")
	assert.Contains(t, warnMessages[0], "key3", "警告メッセージにキー名が含まれていない")
	assert.Contains(t, warnMessages[0], "型の取得に失敗しました", "警告メッセージが正しくない")
}

// TestGetAllValues_BothErrors は両方のエラーが発生した場合のテスト
func TestGetAllValues_BothErrors(t *testing.T) {
	// Arrange
	ctx := context.Background()
	keys := []string{"key1", "key2", "key3", "key4"}
	expectedError := assert.AnError

	warnMessages := []string{}
	mockRepo := &MockDataRepository{
		GetValueFunc: func(ctx context.Context, key string) (string, error) {
			if key == "key2" {
				return "", expectedError
			}
			return "value", nil
		},
		GetTypeFunc: func(ctx context.Context, key string) (string, error) {
			if key == "key3" {
				return "", expectedError
			}
			return "string", nil
		},
	}

	mockLogger := &MockLogger{
		WarnFunc: func(msg string, keysAndValues ...any) {
			warnMessages = append(warnMessages, msg)
		},
	}

	service := NewDataService(mockRepo, mockLogger)

	// Act
	result, err := service.GetAllValues(ctx, keys)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.NotNil(t, result, "結果がnilであってはならない")

	resultMap, ok := result.(map[string]any)
	assert.True(t, ok, "結果がmap[string]any型でない")

	values, ok := resultMap["values"].([]any)
	assert.True(t, ok, "valuesが配列でない")
	assert.Len(t, values, 2, "valuesの長さが正しくない（エラーのキーは除外される）")
	assert.Equal(t, 2, resultMap["count"], "カウントが正しくない")

	// 両方の警告ログが出力されたことを確認
	assert.Len(t, warnMessages, 2, "警告メッセージが2つ出力されていない")
}
