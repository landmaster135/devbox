package usecases

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dbExecutor "github.com/landmaster135/devbox/internal/postgresql/domain/executor"
	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
	dump "github.com/landmaster135/devbox/internal/postgresql/usecases/dump"
)

// #==============================================================#
// ##          Mock Implementations                              ##
// #==============================================================#

// MockDatabaseExecutor はテスト用のDatabaseExecutorモック
type MockDatabaseExecutor struct {
	QueryContextFunc       func(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContextFunc    func(ctx context.Context, query string, args ...any) *sql.Row
	BeginTxFunc            func(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	PingFunc               func() error
	CloseFunc              func() error
	QueryContextRowsFunc   func(ctx context.Context, query string, args ...any) (model.RowsInterface, error)
	QueryRowContextRowFunc func(ctx context.Context, query string, args ...any) model.RowInterface
}

func (m *MockDatabaseExecutor) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if m.QueryContextFunc != nil {
		return m.QueryContextFunc(ctx, query, args)
	}
	return nil, nil
}

func (m *MockDatabaseExecutor) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if m.QueryRowContextFunc != nil {
		return m.QueryRowContextFunc(ctx, query, args)
	}
	return nil
}

func (m *MockDatabaseExecutor) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	if m.BeginTxFunc != nil {
		return m.BeginTxFunc(ctx, opts)
	}
	return nil, nil
}

func (m *MockDatabaseExecutor) Ping() error {
	if m.PingFunc != nil {
		return m.PingFunc()
	}
	return nil
}

func (m *MockDatabaseExecutor) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

// QueryContextRows は新しいインターフェース用のメソッド
func (m *MockDatabaseExecutor) QueryContextRows(ctx context.Context, query string, args ...any) (model.RowsInterface, error) {
	if m.QueryContextRowsFunc != nil {
		return m.QueryContextRowsFunc(ctx, query, args)
	}
	return nil, nil
}

// QueryRowContextRow は新しいインターフェース用のメソッド
func (m *MockDatabaseExecutor) QueryRowContextRow(ctx context.Context, query string, args ...any) model.RowInterface {
	if m.QueryRowContextRowFunc != nil {
		return m.QueryRowContextRowFunc(ctx, query, args)
	}
	return nil
}

// MockTemplateRenderer はテスト用のTemplateRendererモック
type MockTemplateRenderer struct {
	RenderTableDetailFunc func(detail *model.TableDetail) (string, error)
	RenderTableListFunc   func(data model.ListTablesData) (string, error)
}

func (m *MockTemplateRenderer) RenderTableDetail(detail *model.TableDetail) (string, error) {
	if m.RenderTableDetailFunc != nil {
		return m.RenderTableDetailFunc(detail)
	}
	return "", nil
}

func (m *MockTemplateRenderer) RenderTableList(data model.ListTablesData) (string, error) {
	if m.RenderTableListFunc != nil {
		return m.RenderTableListFunc(data)
	}
	return "", nil
}

// MockJSONMarshaler はテスト用のJSONMarshalerモック
type MockJSONMarshaler struct {
	MarshalIndentFunc func(v any, prefix, indent string) ([]byte, error)
}

func (m *MockJSONMarshaler) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	if m.MarshalIndentFunc != nil {
		return m.MarshalIndentFunc(v, prefix, indent)
	}
	return nil, nil
}

// #==============================================================#
// ##          Test Helper Functions                             ##
// #==============================================================#

func createTestPostgreSQLService() *PostgreSQLService {
	mockExecutor := &MockDatabaseExecutor{}
	mockRenderer := &MockTemplateRenderer{}
	mockMarshaler := &MockJSONMarshaler{}
	tableDumper := dump.NewTableDumper(mockExecutor)

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
	tableDumper := dump.NewTableDumper(mockExecutor)
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
	mockExecutor.CloseFunc = func() error {
		return nil
	}

	// Act
	err := service.Close()

	// Assert
	assert.NoError(t, err)
}

func TestPostgreSQLService_Close_Error(t *testing.T) {
	// Arrange
	service := createTestPostgreSQLService()
	mockExecutor := service.executor.(*MockDatabaseExecutor)
	expectedError := errors.New("close error")
	mockExecutor.CloseFunc = func() error {
		return expectedError
	}

	// Act
	err := service.Close()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
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
	data := map[string]any{
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
	data := make(map[string]any)
	data["self"] = data
	prefix := ""
	indent := "  "

	// Act
	result, err := marshaler.MarshalIndent(data, prefix, indent)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDefaultDatabaseExecutor_QueryHelpers(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	executor := &dbExecutor.DefaultDatabaseExecutor{DB: db}
	ctx := context.Background()

	mock.ExpectQuery("SELECT 1").
		WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(1))

	rows, err := executor.QueryContextRows(ctx, "SELECT 1")
	require.NoError(t, err)

	cols, err := rows.Columns()
	require.NoError(t, err)
	assert.Equal(t, []string{"value"}, cols)

	if assert.True(t, rows.Next()) {
		var value int
		require.NoError(t, rows.Scan(&value))
		assert.Equal(t, 1, value)
	}
	assert.False(t, rows.Next())
	assert.NoError(t, rows.Err())
	assert.NoError(t, rows.Close())

	mock.ExpectQuery("SELECT 'hello'").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("hello"))

	row := executor.QueryRowContextRow(ctx, "SELECT 'hello'")
	var name string
	require.NoError(t, row.Scan(&name))
	assert.Equal(t, "hello", name)

	mock.ExpectPing()
	mock.ExpectClose()

	assert.NoError(t, executor.Ping())
	assert.NoError(t, executor.Close())

	require.NoError(t, mock.ExpectationsWereMet())
}
