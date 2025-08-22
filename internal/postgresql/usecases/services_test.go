package usecases

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// #==============================================================#
// ##          Mock Implementations                              ##
// #==============================================================#

// MockDatabaseExecutor はテスト用のDatabaseExecutorモック
type MockDatabaseExecutor struct {
	mock.Mock
}

func (m *MockDatabaseExecutor) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	arguments := m.Called(ctx, query, args)
	return arguments.Get(0).(*sql.Rows), arguments.Error(1)
}

func (m *MockDatabaseExecutor) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	arguments := m.Called(ctx, query, args)
	return arguments.Get(0).(*sql.Row)
}

func (m *MockDatabaseExecutor) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	arguments := m.Called(ctx, opts)
	return arguments.Get(0).(*sql.Tx), arguments.Error(1)
}

func (m *MockDatabaseExecutor) Ping() error {
	arguments := m.Called()
	return arguments.Error(0)
}

func (m *MockDatabaseExecutor) Close() error {
	arguments := m.Called()
	return arguments.Error(0)
}

// MockTemplateRenderer はテスト用のTemplateRendererモック
type MockTemplateRenderer struct {
	mock.Mock
}

func (m *MockTemplateRenderer) RenderTableDetail(detail *TableDetail) (string, error) {
	arguments := m.Called(detail)
	return arguments.String(0), arguments.Error(1)
}

func (m *MockTemplateRenderer) RenderTableList(data ListTablesData) (string, error) {
	arguments := m.Called(data)
	return arguments.String(0), arguments.Error(1)
}

// MockJSONMarshaler はテスト用のJSONMarshalerモック
type MockJSONMarshaler struct {
	mock.Mock
}

func (m *MockJSONMarshaler) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	arguments := m.Called(v, prefix, indent)
	return arguments.Get(0).([]byte), arguments.Error(1)
}

// #==============================================================#
// ##          Test Helper Functions                             ##
// #==============================================================#

func createTestPostgreSQLService() *PostgreSQLService {
	mockExecutor := &MockDatabaseExecutor{}
	mockRenderer := &MockTemplateRenderer{}
	mockMarshaler := &MockJSONMarshaler{}
	tableDumper := NewTableDumper(mockExecutor)

	return NewPostgreSQLServiceWithDependencies(
		mockExecutor,
		mockRenderer,
		mockMarshaler,
		tableDumper,
		"postgres://test:test@localhost/testdb",
		"postgres://test@localhost/testdb",
	)
}

// #==============================================================#
// ##          PostgreSQLService Tests                           ##
// #==============================================================#

func TestNewPostgreSQLService_Normal(t *testing.T) {
	// このテストは実際のデータベース接続が必要なため、スキップ
	t.Skip("実際のデータベース接続が必要なため、スキップします")
}

func TestNewPostgreSQLServiceWithDependencies_Normal(t *testing.T) {
	// Arrange
	mockExecutor := &MockDatabaseExecutor{}
	mockRenderer := &MockTemplateRenderer{}
	mockMarshaler := &MockJSONMarshaler{}
	tableDumper := NewTableDumper(mockExecutor)
	databaseURL := "postgres://test:test@localhost/testdb"
	resourceBase := "postgres://test@localhost/testdb"

	// Act
	service := NewPostgreSQLServiceWithDependencies(
		mockExecutor,
		mockRenderer,
		mockMarshaler,
		tableDumper,
		databaseURL,
		resourceBase,
	)

	// Assert
	assert.NotNil(t, service)
	assert.Equal(t, databaseURL, service.databaseURL)
	assert.Equal(t, resourceBase, service.resourceBase)
}

func TestPostgreSQLService_Close_Normal(t *testing.T) {
	// Arrange
	service := createTestPostgreSQLService()
	mockExecutor := service.executor.(*MockDatabaseExecutor)
	mockExecutor.On("Close").Return(nil)

	// Act
	err := service.Close()

	// Assert
	assert.NoError(t, err)
	mockExecutor.AssertExpectations(t)
}

func TestPostgreSQLService_Close_Error(t *testing.T) {
	// Arrange
	service := createTestPostgreSQLService()
	mockExecutor := service.executor.(*MockDatabaseExecutor)
	expectedError := errors.New("close error")
	mockExecutor.On("Close").Return(expectedError)

	// Act
	err := service.Close()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockExecutor.AssertExpectations(t)
}

// #==============================================================#
// ##          Handler Methods Tests                             ##
// #==============================================================#

func TestPostgreSQLService_HandleToQuery_Normal(t *testing.T) {
	// このテストは実際のデータベース操作が複雑なため、スキップ
	t.Skip("sql.Rowsのモックが複雑なため、統合テストで実施")
}

func TestPostgreSQLService_HandleToGetTableSchema_Normal(t *testing.T) {
	// このテストは実際のデータベース操作が必要なため、スキップ
	t.Skip("実際のデータベース操作が必要なため、統合テストで実施")
}

func TestPostgreSQLService_HandleToListTables_Normal(t *testing.T) {
	// このテストは実際のデータベース操作が必要なため、スキップ
	t.Skip("実際のデータベース操作が必要なため、統合テストで実施")
}

func TestPostgreSQLService_HandleToGetTableSchemaMinimum_Normal(t *testing.T) {
	// このテストは実際のデータベース操作が必要なため、スキップ
	t.Skip("実際のデータベース操作が必要なため、統合テストで実施")
}

func TestPostgreSQLService_HandleToListTablesMinimum_Normal(t *testing.T) {
	// このテストは実際のデータベース操作が必要なため、スキップ
	t.Skip("実際のデータベース操作が必要なため、統合テストで実施")
}

// #==============================================================#
// ##          Utility Function Tests                            ##
// #==============================================================#

func TestCreateResourceBaseURL_Normal(t *testing.T) {
	// Arrange
	databaseURL := "postgres://user:password@localhost:5432/dbname"
	expected := "postgres://user@localhost:5432/dbname"

	// Act
	result, err := createResourceBaseURL(databaseURL)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCreateResourceBaseURL_InvalidURL(t *testing.T) {
	// Arrange
	databaseURL := "invalid-url"

	// Act
	result, err := createResourceBaseURL(databaseURL)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
}

// #==============================================================#
// ##          Default Implementation Tests                       ##
// #==============================================================#

func TestDefaultJSONMarshaler_MarshalIndent_Normal(t *testing.T) {
	// Arrange
	marshaler := &DefaultJSONMarshaler{}
	data := map[string]interface{}{
		"name": "test",
		"id":   1,
	}
	prefix := ""
	indent := "  "

	// Act
	result, err := marshaler.MarshalIndent(data, prefix, indent)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, string(result), "test")
	assert.Contains(t, string(result), "1")
}

func TestDefaultJSONMarshaler_MarshalIndent_InvalidData(t *testing.T) {
	// Arrange
	marshaler := &DefaultJSONMarshaler{}
	// 循環参照を作成してJSONエラーを発生させる
	data := make(map[string]interface{})
	data["self"] = data
	prefix := ""
	indent := "  "

	// Act
	result, err := marshaler.MarshalIndent(data, prefix, indent)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}
