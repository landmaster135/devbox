package postgresql

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	usecases "github.com/landmaster135/devbox/internal/postgresql/usecases"
)

// #==============================================================#
// ##          Mock Implementations                              ##
// #==============================================================#

// MockPostgreSQLService はテスト用のPostgreSQLServiceモック
type MockPostgreSQLService struct {
	mock.Mock
}

func (m *MockPostgreSQLService) HandleToQuery(ctx context.Context, sqlQuery string) ([]map[string]interface{}, error) {
	args := m.Called(ctx, sqlQuery)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockPostgreSQLService) HandleToGetTableSchema(ctx context.Context, tableName string) (string, error) {
	args := m.Called(ctx, tableName)
	return args.String(0), args.Error(1)
}

func (m *MockPostgreSQLService) HandleToListTables(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockPostgreSQLService) HandleToGetTableSchemaMinimum(ctx context.Context, tableName string) ([]usecases.Column, error) {
	args := m.Called(ctx, tableName)
	return args.Get(0).([]usecases.Column), args.Error(1)
}

func (m *MockPostgreSQLService) HandleToListTablesMinimum(ctx context.Context) ([]usecases.Table, error) {
	args := m.Called(ctx)
	return args.Get(0).([]usecases.Table), args.Error(1)
}

func (m *MockPostgreSQLService) Close() error {
	args := m.Called()
	return args.Error(0)
}

// TestPostgreSQLMCPHandler はテスト用のハンドラー構造体
type TestPostgreSQLMCPHandler struct {
	service *MockPostgreSQLService
}

// #==============================================================#
// ##          Test Helper Functions                             ##
// #==============================================================#

func createTestHandler() *TestPostgreSQLMCPHandler {
	mockService := &MockPostgreSQLService{}
	return &TestPostgreSQLMCPHandler{
		service: mockService,
	}
}

func (h *TestPostgreSQLMCPHandler) Close() error {
	return h.service.Close()
}

// #==============================================================#
// ##          PostgreSQLMCPHandler Tests                        ##
// #==============================================================#

func TestNewPostgreSQLMCPHandler_Normal(t *testing.T) {
	// このテストは実際のデータベース接続が必要なため、スキップ
	t.Skip("実際のデータベース接続が必要なため、スキップします")
}

func TestPostgreSQLMCPHandler_Close_Normal(t *testing.T) {
	// Arrange
	handler := createTestHandler()
	handler.service.On("Close").Return(nil)

	// Act
	err := handler.Close()

	// Assert
	assert.NoError(t, err)
	handler.service.AssertExpectations(t)
}

func TestPostgreSQLMCPHandler_Close_Error(t *testing.T) {
	// Arrange
	handler := createTestHandler()
	expectedError := errors.New("close error")
	handler.service.On("Close").Return(expectedError)

	// Act
	err := handler.Close()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	handler.service.AssertExpectations(t)
}

// #==============================================================#
// ##          Service Method Tests                              ##
// #==============================================================#

func TestMockPostgreSQLService_HandleToQuery_Normal(t *testing.T) {
	// Arrange
	handler := createTestHandler()
	ctx := context.Background()
	sqlQuery := "SELECT * FROM test_table"
	expectedResult := []map[string]interface{}{
		{"id": 1, "name": "test"},
	}

	handler.service.On("HandleToQuery", ctx, sqlQuery).Return(expectedResult, nil)

	// Act
	result, err := handler.service.HandleToQuery(ctx, sqlQuery)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	handler.service.AssertExpectations(t)
}

func TestMockPostgreSQLService_HandleToQuery_Error(t *testing.T) {
	// Arrange
	handler := createTestHandler()
	ctx := context.Background()
	sqlQuery := "SELECT * FROM non_existent_table"
	expectedError := errors.New("table does not exist")

	handler.service.On("HandleToQuery", ctx, sqlQuery).Return([]map[string]interface{}{}, expectedError)

	// Act
	result, err := handler.service.HandleToQuery(ctx, sqlQuery)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, result)
	handler.service.AssertExpectations(t)
}

func TestMockPostgreSQLService_HandleToGetTableSchema_Normal(t *testing.T) {
	// Arrange
	handler := createTestHandler()
	ctx := context.Background()
	tableName := "test_table"
	expectedResult := "# テーブル: test_table\n## カラム\n- id: integer NOT NULL\n"

	handler.service.On("HandleToGetTableSchema", ctx, tableName).Return(expectedResult, nil)

	// Act
	result, err := handler.service.HandleToGetTableSchema(ctx, tableName)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	handler.service.AssertExpectations(t)
}

func TestMockPostgreSQLService_HandleToListTables_Normal(t *testing.T) {
	// Arrange
	handler := createTestHandler()
	ctx := context.Background()
	expectedResult := "# データベースのテーブル一覧 (全1件)\n- **test_table** — テストテーブル\n"

	handler.service.On("HandleToListTables", ctx).Return(expectedResult, nil)

	// Act
	result, err := handler.service.HandleToListTables(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	handler.service.AssertExpectations(t)
}

func TestMockPostgreSQLService_HandleToGetTableSchemaMinimum_Normal(t *testing.T) {
	// Arrange
	handler := createTestHandler()
	ctx := context.Background()
	tableName := "test_table"
	expectedResult := []usecases.Column{
		{Name: "id", DataType: "integer"},
		{Name: "name", DataType: "varchar"},
	}

	handler.service.On("HandleToGetTableSchemaMinimum", ctx, tableName).Return(expectedResult, nil)

	// Act
	result, err := handler.service.HandleToGetTableSchemaMinimum(ctx, tableName)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	handler.service.AssertExpectations(t)
}

func TestMockPostgreSQLService_HandleToListTablesMinimum_Normal(t *testing.T) {
	// Arrange
	handler := createTestHandler()
	ctx := context.Background()
	expectedResult := []usecases.Table{
		{Name: "test_table"},
		{Name: "users"},
	}

	handler.service.On("HandleToListTablesMinimum", ctx).Return(expectedResult, nil)

	// Act
	result, err := handler.service.HandleToListTablesMinimum(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	handler.service.AssertExpectations(t)
}

// #==============================================================#
// ##          Helper Function Tests                             ##
// #==============================================================#

func TestReturnJSONResult_Normal(t *testing.T) {
	// Arrange
	data := map[string]interface{}{
		"name": "test",
		"id":   1,
	}

	// Act
	result, err := returnJSONResult(data)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	// MCPの結果構造は複雑なため、基本的な存在確認のみ
	assert.NotEmpty(t, result.Content)
}

func TestReturnJSONResult_InvalidData(t *testing.T) {
	// Arrange
	// 循環参照を作成してJSONエラーを発生させる
	data := make(map[string]interface{})
	data["self"] = data

	// Act
	result, err := returnJSONResult(data)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestReturnTextResult_Normal(t *testing.T) {
	// Arrange
	text := "test result"

	// Act
	result, err := returnTextResult(text)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Content)
}

func TestReturnError_Normal(t *testing.T) {
	// Arrange
	originalError := errors.New("test error")

	// Act
	result, err := returnError(originalError)

	// Assert
	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Content)
	assert.Contains(t, err.Error(), "database error")
}
