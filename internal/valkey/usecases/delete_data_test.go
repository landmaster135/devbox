package valkey

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDeleteData_DryRun_SingleKey はドライランモードでの単一キー削除テスト
func TestDeleteData_DryRun_SingleKey(t *testing.T) {
	// Arrange
	ctx := context.Background()
	key := "testKey"
	mockRepo := &MockDataRepository{}
	mockLogger := &MockLogger{}
	service := NewDataService(mockRepo)

	// Act
	result, err := service.DeleteData(ctx, key, []string{}, "", true, mockLogger)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.NotNil(t, result, "結果がnilであってはならない")

	resultMap, ok := result.(map[string]any)
	assert.True(t, ok, "結果がmap[string]any型でない")
	assert.Equal(t, key, resultMap["key"], "キーが正しくない")
	assert.Equal(t, false, resultMap["deleted"], "削除フラグが正しくない")
	assert.Contains(t, resultMap["message"].(string), "ドライランモード", "メッセージが正しくない")
}

// TestDeleteData_DryRun_MultipleKeys はドライランモードでの複数キー削除テスト
func TestDeleteData_DryRun_MultipleKeys(t *testing.T) {
	// Arrange
	ctx := context.Background()
	keys := []string{"key1", "key2", "key3"}
	mockRepo := &MockDataRepository{}
	mockLogger := &MockLogger{}
	service := NewDataService(mockRepo)

	// Act
	result, err := service.DeleteData(ctx, "", keys, "", true, mockLogger)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.NotNil(t, result, "結果がnilであってはならない")

	resultMap, ok := result.(map[string]any)
	assert.True(t, ok, "結果がmap[string]any型でない")
	assert.Equal(t, keys, resultMap["keys"], "キーが正しくない")
	assert.Equal(t, len(keys), resultMap["count"], "カウントが正しくない")
	assert.Equal(t, 0, resultMap["deleted"], "削除カウントが正しくない")
	assert.Contains(t, resultMap["message"].(string), "ドライランモード", "メッセージが正しくない")
}

// TestDeleteData_DryRun_Pattern はドライランモードでのパターン削除テスト
func TestDeleteData_DryRun_Pattern(t *testing.T) {
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
	service := NewDataService(mockRepo)

	// Act
	result, err := service.DeleteData(ctx, "", []string{}, pattern, true, mockLogger)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.NotNil(t, result, "結果がnilであってはならない")

	resultMap, ok := result.(map[string]any)
	assert.True(t, ok, "結果がmap[string]any型でない")
	assert.Equal(t, expectedKeys, resultMap["keys"], "キーが正しくない")
	assert.Equal(t, len(expectedKeys), resultMap["count"], "カウントが正しくない")
	assert.Equal(t, 0, resultMap["deleted"], "削除カウントが正しくない")
	assert.Contains(t, resultMap["message"].(string), "ドライランモード", "メッセージが正しくない")
}

// TestDeleteData_DryRun_PatternNoMatch はドライランモードでパターンに一致するキーがない場合のテスト
func TestDeleteData_DryRun_PatternNoMatch(t *testing.T) {
	// Arrange
	ctx := context.Background()
	pattern := "nomatch*"

	mockRepo := &MockDataRepository{
		GetKeysFunc: func(ctx context.Context, p string) ([]string, error) {
			return []string{}, nil
		},
	}

	mockLogger := &MockLogger{}
	service := NewDataService(mockRepo)

	// Act
	result, err := service.DeleteData(ctx, "", []string{}, pattern, true, mockLogger)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.NotNil(t, result, "結果がnilであってはならない")

	resultMap, ok := result.(map[string]any)
	assert.True(t, ok, "結果がmap[string]any型でない")
	assert.Contains(t, resultMap["message"].(string), "一致するキーが見つかりませんでした", "メッセージが正しくない")
	assert.Equal(t, 0, resultMap["count"], "カウントが正しくない")
}

// TestDeleteData_DryRun_NoArguments はドライランモードで引数が指定されていない場合のエラーテスト
func TestDeleteData_DryRun_NoArguments(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &MockDataRepository{}
	mockLogger := &MockLogger{}
	service := NewDataService(mockRepo)

	// Act
	result, err := service.DeleteData(ctx, "", []string{}, "", true, mockLogger)

	// Assert
	assert.Error(t, err, "エラーが発生するべき")
	assert.Contains(t, err.Error(), "キーまたはパターンを指定してください", "エラーメッセージが正しくない")
	assert.Nil(t, result, "結果はnilであるべき")
}

// TestDeleteData_Actual_SingleKey は実際の削除モードでの単一キー削除テスト
func TestDeleteData_Actual_SingleKey(t *testing.T) {
	// Arrange
	ctx := context.Background()
	key := "testKey"
	expectedDeleted := true

	mockRepo := &MockDataRepository{
		DeleteKeyFunc: func(ctx context.Context, k string) (bool, error) {
			assert.Equal(t, key, k, "キーが正しく渡されていない")
			return expectedDeleted, nil
		},
	}

	mockLogger := &MockLogger{}
	service := NewDataService(mockRepo)

	// Act
	result, err := service.DeleteData(ctx, key, []string{}, "", false, mockLogger)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.NotNil(t, result, "結果がnilであってはならない")

	resultMap, ok := result.(map[string]any)
	assert.True(t, ok, "結果がmap[string]any型でない")
	assert.Equal(t, key, resultMap["key"], "キーが正しくない")
	assert.Equal(t, expectedDeleted, resultMap["deleted"], "削除フラグが正しくない")
}

// TestDeleteData_Actual_SingleKey_Error は実際の削除モードでの単一キー削除エラーテスト
func TestDeleteData_Actual_SingleKey_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	key := "testKey"
	expectedError := assert.AnError

	mockRepo := &MockDataRepository{
		DeleteKeyFunc: func(ctx context.Context, k string) (bool, error) {
			return false, expectedError
		},
	}

	mockLogger := &MockLogger{}
	service := NewDataService(mockRepo)

	// Act
	result, err := service.DeleteData(ctx, key, []string{}, "", false, mockLogger)

	// Assert
	assert.Error(t, err, "エラーが発生するべき")
	assert.Contains(t, err.Error(), "削除に失敗しました", "エラーメッセージが正しくない")
	assert.Nil(t, result, "結果はnilであるべき")
}

// TestDeleteData_Actual_MultipleKeys は実際の削除モードでの複数キー削除テスト
func TestDeleteData_Actual_MultipleKeys(t *testing.T) {
	// Arrange
	ctx := context.Background()
	keys := []string{"key1", "key2", "key3"}

	deleteResults := map[string]bool{
		"key1": true,
		"key2": false,
		"key3": true,
	}

	mockRepo := &MockDataRepository{
		DeleteKeyFunc: func(ctx context.Context, k string) (bool, error) {
			result, exists := deleteResults[k]
			assert.True(t, exists, "予期しないキーが渡された: "+k)
			return result, nil
		},
	}

	mockLogger := &MockLogger{}
	service := NewDataService(mockRepo)

	// Act
	result, err := service.DeleteData(ctx, "", keys, "", false, mockLogger)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.NotNil(t, result, "結果がnilであってはならない")

	resultMap, ok := result.(map[string]any)
	assert.True(t, ok, "結果がmap[string]any型でない")
	assert.Equal(t, keys, resultMap["keys"], "キーが正しくない")
	assert.Equal(t, deleteResults, resultMap["results"], "削除結果が正しくない")
	assert.Equal(t, len(keys), resultMap["count"], "カウントが正しくない")
	assert.Equal(t, 2, resultMap["deleted"], "削除カウントが正しくない") // key1とkey3が削除成功
}

// TestDeleteData_Actual_MultipleKeys_Error は実際の削除モードでの複数キー削除エラーテスト
func TestDeleteData_Actual_MultipleKeys_Error(t *testing.T) {
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
	service := NewDataService(mockRepo)

	// Act
	result, err := service.DeleteData(ctx, "", keys, "", false, mockLogger)

	// Assert
	assert.Error(t, err, "エラーが発生するべき")
	assert.Contains(t, err.Error(), "削除に失敗しました", "エラーメッセージが正しくない")
	assert.Nil(t, result, "結果はnilであるべき")
}

// TestDeleteData_Actual_Pattern は実際の削除モードでのパターン削除テスト
func TestDeleteData_Actual_Pattern(t *testing.T) {
	// Arrange
	ctx := context.Background()
	pattern := "test*"
	expectedKeys := []string{"test1", "test2"}

	deleteResults := map[string]bool{
		"test1": true,
		"test2": false,
	}

	mockRepo := &MockDataRepository{
		GetKeysFunc: func(ctx context.Context, p string) ([]string, error) {
			assert.Equal(t, pattern, p, "パターンが正しく渡されていない")
			return expectedKeys, nil
		},
		DeleteKeyFunc: func(ctx context.Context, k string) (bool, error) {
			result, exists := deleteResults[k]
			assert.True(t, exists, "予期しないキーが渡された: "+k)
			return result, nil
		},
	}

	mockLogger := &MockLogger{}
	service := NewDataService(mockRepo)

	// Act
	result, err := service.DeleteData(ctx, "", []string{}, pattern, false, mockLogger)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.NotNil(t, result, "結果がnilであってはならない")

	resultMap, ok := result.(map[string]any)
	assert.True(t, ok, "結果がmap[string]any型でない")
	assert.Equal(t, expectedKeys, resultMap["keys"], "キーが正しくない")
	assert.Equal(t, deleteResults, resultMap["results"], "削除結果が正しくない")
	assert.Equal(t, len(expectedKeys), resultMap["count"], "カウントが正しくない")
	assert.Equal(t, 1, resultMap["deleted"], "削除カウントが正しくない") // test1のみ削除成功
}

// TestDeleteData_Actual_Pattern_GetKeysError は実際の削除モードでのパターン削除時のGetKeysエラーテスト
func TestDeleteData_Actual_Pattern_GetKeysError(t *testing.T) {
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
	service := NewDataService(mockRepo)

	// Act
	result, err := service.DeleteData(ctx, "", []string{}, pattern, false, mockLogger)

	// Assert
	assert.Error(t, err, "エラーが発生するべき")
	assert.Contains(t, err.Error(), "に一致するキーの取得に失敗しました", "エラーメッセージが正しくない")
	assert.Nil(t, result, "結果はnilであるべき")
}

// TestDeleteData_Actual_Pattern_NoMatch は実際の削除モードでパターンに一致するキーがない場合のテスト
func TestDeleteData_Actual_Pattern_NoMatch(t *testing.T) {
	// Arrange
	ctx := context.Background()
	pattern := "nomatch*"

	mockRepo := &MockDataRepository{
		GetKeysFunc: func(ctx context.Context, p string) ([]string, error) {
			return []string{}, nil
		},
	}

	mockLogger := &MockLogger{}
	service := NewDataService(mockRepo)

	// Act
	result, err := service.DeleteData(ctx, "", []string{}, pattern, false, mockLogger)

	// Assert
	assert.NoError(t, err, "エラーが発生してはならない")
	assert.NotNil(t, result, "結果がnilであってはならない")

	resultMap, ok := result.(map[string]any)
	assert.True(t, ok, "結果がmap[string]any型でない")
	assert.Contains(t, resultMap["message"].(string), "一致するキーが見つかりませんでした", "メッセージが正しくない")
	assert.Equal(t, 0, resultMap["count"], "カウントが正しくない")
}

// TestDeleteData_Actual_Pattern_DeleteKeysError は実際の削除モードでのパターン削除時のDeleteKeysエラーテスト
func TestDeleteData_Actual_Pattern_DeleteKeysError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	pattern := "test*"
	expectedKeys := []string{"test1", "test2"}
	expectedError := assert.AnError

	callCount := 0
	mockRepo := &MockDataRepository{
		GetKeysFunc: func(ctx context.Context, p string) ([]string, error) {
			return expectedKeys, nil
		},
		DeleteKeyFunc: func(ctx context.Context, k string) (bool, error) {
			callCount++
			if callCount == 2 {
				return false, expectedError
			}
			return true, nil
		},
	}

	mockLogger := &MockLogger{}
	service := NewDataService(mockRepo)

	// Act
	result, err := service.DeleteData(ctx, "", []string{}, pattern, false, mockLogger)

	// Assert
	assert.Error(t, err, "エラーが発生するべき")
	assert.Contains(t, err.Error(), "削除に失敗しました", "エラーメッセージが正しくない")
	assert.Nil(t, result, "結果はnilであるべき")
}

// TestDeleteData_Actual_NoArguments は実際の削除モードで引数が指定されていない場合のエラーテスト
func TestDeleteData_Actual_NoArguments(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &MockDataRepository{}
	mockLogger := &MockLogger{}
	service := NewDataService(mockRepo)

	// Act
	result, err := service.DeleteData(ctx, "", []string{}, "", false, mockLogger)

	// Assert
	assert.Error(t, err, "エラーが発生するべき")
	assert.Contains(t, err.Error(), "キーまたはパターンを指定してください", "エラーメッセージが正しくない")
	assert.Nil(t, result, "結果はnilであるべき")
}
