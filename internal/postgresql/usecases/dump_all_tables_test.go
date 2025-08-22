package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// #==============================================================#
// ##          Mock SQL Types                                    ##
// #==============================================================#

// MockRows は sql.Rows のモック
type MockRows struct {
	columns []string
	data    [][]interface{}
	index   int
	closed  bool
}

func NewMockRows(columns []string, data [][]interface{}) *MockRows {
	return &MockRows{
		columns: columns,
		data:    data,
		index:   -1,
		closed:  false,
	}
}

func (m *MockRows) Columns() ([]string, error) {
	return m.columns, nil
}

func (m *MockRows) Next() bool {
	if m.closed || m.index >= len(m.data)-1 {
		return false
	}
	m.index++
	return true
}

func (m *MockRows) Scan(dest ...interface{}) error {
	if m.closed || m.index < 0 || m.index >= len(m.data) {
		return errors.New("no rows")
	}

	row := m.data[m.index]
	for i, val := range row {
		if i < len(dest) {
			switch d := dest[i].(type) {
			case *string:
				if val != nil {
					*d = val.(string)
				}
			case *interface{}:
				*d = val
			}
		}
	}
	return nil
}

func (m *MockRows) Close() error {
	m.closed = true
	return nil
}

func (m *MockRows) Err() error {
	return nil
}

// MockRow は sql.Row のモック
type MockRow struct {
	data []interface{}
	err  error
}

func NewMockRow(data []interface{}, err error) *MockRow {
	return &MockRow{data: data, err: err}
}

func (m *MockRow) Scan(dest ...interface{}) error {
	if m.err != nil {
		return m.err
	}

	for i, val := range m.data {
		if i < len(dest) {
			switch d := dest[i].(type) {
			case *string:
				if val != nil {
					*d = val.(string)
				}
			case *interface{}:
				*d = val
			}
		}
	}
	return nil
}

// #==============================================================#
// ##          Mock Database Query Executor                      ##
// #==============================================================#

// MockDatabaseQueryExecutor は新しいインターフェース用のモック
type MockDatabaseQueryExecutor struct {
	mock.Mock
}

func (m *MockDatabaseQueryExecutor) QueryContextRows(ctx context.Context, query string, args ...interface{}) (RowsInterface, error) {
	arguments := m.Called(ctx, query, args)
	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}
	return arguments.Get(0).(RowsInterface), arguments.Error(1)
}

func (m *MockDatabaseQueryExecutor) QueryRowContextRow(ctx context.Context, query string, args ...interface{}) RowInterface {
	arguments := m.Called(ctx, query, args)
	if arguments.Get(0) == nil {
		return nil
	}
	return arguments.Get(0).(RowInterface)
}

// #==============================================================#
// ##          Test Helper Functions                             ##
// #==============================================================#

func createTestTableDumperForDumpAllTables() (*TableDumper, *MockDatabaseQueryExecutor, *MockFileWriter) {
	mockExecutor := &MockDatabaseQueryExecutor{}
	mockFileWriter := &MockFileWriter{}

	dumper := NewTableDumperWithDependencies(mockExecutor, mockFileWriter)

	return dumper, mockExecutor, mockFileWriter
}

// #==============================================================#
// ##          Test Class                                        ##
// #==============================================================#

type TestTableDumper_DumpAllTables struct {
	dumper       *TableDumper
	mockExecutor *MockDatabaseQueryExecutor
	mockWriter   *MockFileWriter
}

func (suite *TestTableDumper_DumpAllTables) setup() {
	suite.dumper, suite.mockExecutor, suite.mockWriter = createTestTableDumperForDumpAllTables()
}

// #==============================================================#
// ##          DumpAllTables Function Tests                      ##
// #==============================================================#

func TestTableDumper_DumpAllTables_Normal(t *testing.T) {
	// Arrange
	suite := &TestTableDumper_DumpAllTables{}
	suite.setup()

	ctx := context.Background()
	outputPath := "/test/output"
	format := "json"
	concurrency := 2

	// モックの設定
	suite.mockWriter.On("MkdirAll", outputPath, mock.AnythingOfType("fs.FileMode")).Return(nil)

	// テーブル一覧取得の成功レスポンス
	tableRows := NewMockRows([]string{"table_name"}, [][]interface{}{
		{"users"},
		{"products"},
	})
	suite.mockExecutor.On("QueryContextRows", ctx, mock.MatchedBy(func(query string) bool {
		return query == `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	`
	}), mock.Anything).Return(tableRows, nil)

	// データベース名取得の成功レスポンス
	dbRow := NewMockRow([]interface{}{"testdb"}, nil)
	suite.mockExecutor.On("QueryRowContextRow", ctx, "SELECT current_database()", mock.Anything).Return(dbRow)

	// 各テーブルのダンプクエリの成功レスポンス
	usersRows := NewMockRows([]string{"id", "name"}, [][]interface{}{
		{1, "John"},
		{2, "Jane"},
	})
	suite.mockExecutor.On("QueryContextRows", ctx, "SELECT * FROM users", mock.Anything).Return(usersRows, nil)

	productsRows := NewMockRows([]string{"id", "title"}, [][]interface{}{
		{1, "Product A"},
	})
	suite.mockExecutor.On("QueryContextRows", ctx, "SELECT * FROM products", mock.Anything).Return(productsRows, nil)

	// ファイル書き込みの成功レスポンス
	suite.mockWriter.On("WriteFile", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("fs.FileMode")).Return(nil)

	// Act
	result, err := suite.dumper.DumpAllTables(ctx, outputPath, format, nil, &concurrency)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testdb", result.DatabaseName)
	assert.Equal(t, 2, result.TotalTables)
	assert.Len(t, result.Results, 2)
	assert.Len(t, result.FailedTables, 0)

	// 各テーブルの結果を確認
	tableNames := make(map[string]bool)
	for _, res := range result.Results {
		tableNames[res.TableName] = true
		assert.Equal(t, format, res.Format)
		assert.Equal(t, outputPath, res.OutputPath)
	}
	assert.True(t, tableNames["users"])
	assert.True(t, tableNames["products"])

	// モックの期待値が満たされたことを確認
	suite.mockExecutor.AssertExpectations(t)
	suite.mockWriter.AssertExpectations(t)
}

func TestTableDumper_DumpAllTables_EmptyDatabase(t *testing.T) {
	// Arrange
	suite := &TestTableDumper_DumpAllTables{}
	suite.setup()

	ctx := context.Background()
	outputPath := "/test/output"
	format := "json"
	concurrency := 2

	// モックの設定
	suite.mockWriter.On("MkdirAll", outputPath, mock.AnythingOfType("fs.FileMode")).Return(nil)

	// 空のテーブル一覧
	emptyRows := NewMockRows([]string{"table_name"}, [][]interface{}{})
	suite.mockExecutor.On("QueryContextRows", ctx, mock.MatchedBy(func(query string) bool {
		return query == `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	`
	}), mock.Anything).Return(emptyRows, nil)

	// データベース名取得
	dbRow := NewMockRow([]interface{}{"testdb"}, nil)
	suite.mockExecutor.On("QueryRowContextRow", ctx, "SELECT current_database()", mock.Anything).Return(dbRow)

	// Act
	result, err := suite.dumper.DumpAllTables(ctx, outputPath, format, nil, &concurrency)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testdb", result.DatabaseName)
	assert.Equal(t, 0, result.TotalTables)
	assert.Len(t, result.Results, 0)
	assert.Len(t, result.FailedTables, 0)

	// モックの期待値が満たされたことを確認
	suite.mockExecutor.AssertExpectations(t)
	suite.mockWriter.AssertExpectations(t)
}

func TestTableDumper_DumpAllTables_DatabaseError(t *testing.T) {
	// Arrange
	suite := &TestTableDumper_DumpAllTables{}
	suite.setup()

	ctx := context.Background()
	outputPath := "/test/output"
	format := "json"
	concurrency := 2

	// モックの設定
	suite.mockWriter.On("MkdirAll", outputPath, mock.AnythingOfType("fs.FileMode")).Return(nil)

	// データベースエラーを発生させる
	suite.mockExecutor.On("QueryContextRows", ctx, mock.MatchedBy(func(query string) bool {
		return query == `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	`
	}), mock.Anything).Return(nil, errors.New("データベース接続エラー"))

	// Act
	result, err := suite.dumper.DumpAllTables(ctx, outputPath, format, nil, &concurrency)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "テーブル一覧取得エラー")

	// モックの期待値が満たされたことを確認
	suite.mockExecutor.AssertExpectations(t)
	suite.mockWriter.AssertExpectations(t)
}

func TestTableDumper_DumpAllTables_PartialFailure(t *testing.T) {
	// Arrange
	suite := &TestTableDumper_DumpAllTables{}
	suite.setup()

	ctx := context.Background()
	outputPath := "/test/output"
	format := "json"
	concurrency := 2

	// モックの設定
	suite.mockWriter.On("MkdirAll", outputPath, mock.AnythingOfType("fs.FileMode")).Return(nil)

	// テーブル一覧取得
	tableRows := NewMockRows([]string{"table_name"}, [][]interface{}{
		{"users"},
		{"products"},
	})
	suite.mockExecutor.On("QueryContextRows", ctx, mock.MatchedBy(func(query string) bool {
		return query == `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	`
	}), mock.Anything).Return(tableRows, nil)

	// データベース名取得
	dbRow := NewMockRow([]interface{}{"testdb"}, nil)
	suite.mockExecutor.On("QueryRowContextRow", ctx, "SELECT current_database()", mock.Anything).Return(dbRow)

	// usersテーブルは成功
	usersRows := NewMockRows([]string{"id", "name"}, [][]interface{}{
		{1, "John"},
	})
	suite.mockExecutor.On("QueryContextRows", ctx, "SELECT * FROM users", mock.Anything).Return(usersRows, nil)

	// productsテーブルは失敗
	suite.mockExecutor.On("QueryContextRows", ctx, "SELECT * FROM products", mock.Anything).Return(nil, errors.New("テーブルアクセスエラー"))

	// ファイル書き込み（成功したテーブルのみ）
	suite.mockWriter.On("WriteFile", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("fs.FileMode")).Return(nil)

	// Act
	result, err := suite.dumper.DumpAllTables(ctx, outputPath, format, nil, &concurrency)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testdb", result.DatabaseName)
	assert.Equal(t, 2, result.TotalTables)
	assert.Len(t, result.Results, 1)
	assert.Len(t, result.FailedTables, 1)

	// 成功したテーブル
	assert.Equal(t, "users", result.Results[0].TableName)

	// 失敗したテーブル
	assert.Equal(t, "products", result.FailedTables[0].TableName)
	assert.Contains(t, result.FailedTables[0].Error, "テーブルアクセスエラー")

	// モックの期待値が満たされたことを確認
	suite.mockExecutor.AssertExpectations(t)
	suite.mockWriter.AssertExpectations(t)
}

func TestTableDumper_DumpAllTables_WithLimit(t *testing.T) {
	// Arrange
	suite := &TestTableDumper_DumpAllTables{}
	suite.setup()

	ctx := context.Background()
	outputPath := "/test/output"
	format := "json"
	limit := 10
	concurrency := 2

	// モックの設定
	suite.mockWriter.On("MkdirAll", outputPath, mock.AnythingOfType("fs.FileMode")).Return(nil)

	// テーブル一覧取得
	tableRows := NewMockRows([]string{"table_name"}, [][]interface{}{
		{"users"},
	})
	suite.mockExecutor.On("QueryContextRows", ctx, mock.MatchedBy(func(query string) bool {
		return query == `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	`
	}), mock.Anything).Return(tableRows, nil)

	// データベース名取得
	dbRow := NewMockRow([]interface{}{"testdb"}, nil)
	suite.mockExecutor.On("QueryRowContextRow", ctx, "SELECT current_database()", mock.Anything).Return(dbRow)

	// usersテーブルのダンプ（LIMIT付き）
	usersRows := NewMockRows([]string{"id", "name"}, [][]interface{}{
		{1, "John"},
		{2, "Jane"},
	})
	suite.mockExecutor.On("QueryContextRows", ctx, "SELECT * FROM users LIMIT 10", mock.Anything).Return(usersRows, nil)

	// ファイル書き込み
	suite.mockWriter.On("WriteFile", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("fs.FileMode")).Return(nil)

	// Act
	result, err := suite.dumper.DumpAllTables(ctx, outputPath, format, &limit, &concurrency)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testdb", result.DatabaseName)
	assert.Equal(t, 1, result.TotalTables)
	assert.Len(t, result.Results, 1)
	assert.Len(t, result.FailedTables, 0)

	// 結果の確認
	assert.Equal(t, "users", result.Results[0].TableName)
	assert.Equal(t, 2, result.Results[0].RecordCount)

	// モックの期待値が満たされたことを確認
	suite.mockExecutor.AssertExpectations(t)
	suite.mockWriter.AssertExpectations(t)
}

func TestTableDumper_DumpAllTables_DirectoryCreationError(t *testing.T) {
	// Arrange
	suite := &TestTableDumper_DumpAllTables{}
	suite.setup()

	ctx := context.Background()
	outputPath := "/invalid/path"
	format := "json"
	concurrency := 2

	// ディレクトリ作成エラーを発生させる
	suite.mockWriter.On("MkdirAll", outputPath, mock.AnythingOfType("fs.FileMode")).Return(errors.New("permission denied"))

	// Act
	result, err := suite.dumper.DumpAllTables(ctx, outputPath, format, nil, &concurrency)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "出力ディレクトリ作成エラー")

	// モックの期待値が満たされたことを確認
	suite.mockWriter.AssertExpectations(t)
}

// #==============================================================#
// ##          Error Path Tests                                  ##
// #==============================================================#

func TestTableDumper_DumpAllTables_Success_Normal(t *testing.T) {
	suite := &TestTableDumper_DumpAllTables{}
	suite.setup()

	// Arrange
	ctx := context.Background()
	outputPath := "/test/output"
	format := "json"
	concurrency := 2

	// モックの設定
	suite.mockWriter.On("MkdirAll", outputPath, mock.AnythingOfType("fs.FileMode")).Return(nil)

	// テーブル一覧取得の成功レスポンス
	tableRows := NewMockRows([]string{"table_name"}, [][]interface{}{
		{"users"},
		{"products"},
	})
	suite.mockExecutor.On("QueryContextRows", ctx, mock.MatchedBy(func(query string) bool {
		return query == `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	`
	}), mock.Anything).Return(tableRows, nil)

	// データベース名取得の成功レスポンス
	dbRow := NewMockRow([]interface{}{"testdb"}, nil)
	suite.mockExecutor.On("QueryRowContextRow", ctx, "SELECT current_database()", mock.Anything).Return(dbRow)

	// 各テーブルのダンプクエリの成功レスポンス
	usersRows := NewMockRows([]string{"id", "name"}, [][]interface{}{
		{1, "John"},
		{2, "Jane"},
	})
	suite.mockExecutor.On("QueryContextRows", ctx, "SELECT * FROM users", mock.Anything).Return(usersRows, nil)

	productsRows := NewMockRows([]string{"id", "title"}, [][]interface{}{
		{1, "Product A"},
	})
	suite.mockExecutor.On("QueryContextRows", ctx, "SELECT * FROM products", mock.Anything).Return(productsRows, nil)

	// ファイル書き込みの成功レスポンス
	suite.mockWriter.On("WriteFile", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("fs.FileMode")).Return(nil)

	// Act
	result, err := suite.dumper.DumpAllTables(ctx, outputPath, format, nil, &concurrency)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testdb", result.DatabaseName)
	assert.Equal(t, 2, result.TotalTables)
	assert.Len(t, result.Results, 2)
	assert.Len(t, result.FailedTables, 0)

	// 各テーブルの結果を確認
	tableNames := make(map[string]bool)
	for _, res := range result.Results {
		tableNames[res.TableName] = true
		assert.Equal(t, format, res.Format)
		assert.Equal(t, outputPath, res.OutputPath)
	}
	assert.True(t, tableNames["users"])
	assert.True(t, tableNames["products"])
}

func TestTableDumper_DumpAllTables_Success_EmptyTables(t *testing.T) {
	suite := &TestTableDumper_DumpAllTables{}
	suite.setup()

	// Arrange
	ctx := context.Background()
	outputPath := "/test/output"
	format := "json"
	concurrency := 2

	// モックの設定
	suite.mockWriter.On("MkdirAll", outputPath, mock.AnythingOfType("fs.FileMode")).Return(nil)

	// 空のテーブル一覧
	emptyRows := NewMockRows([]string{"table_name"}, [][]interface{}{})
	suite.mockExecutor.On("QueryContextRows", ctx, mock.MatchedBy(func(query string) bool {
		return query == `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	`
	}), mock.Anything).Return(emptyRows, nil)

	// データベース名取得
	dbRow := NewMockRow([]interface{}{"testdb"}, nil)
	suite.mockExecutor.On("QueryRowContextRow", ctx, "SELECT current_database()", mock.Anything).Return(dbRow)

	// Act
	result, err := suite.dumper.DumpAllTables(ctx, outputPath, format, nil, &concurrency)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testdb", result.DatabaseName)
	assert.Equal(t, 0, result.TotalTables)
	assert.Len(t, result.Results, 0)
	assert.Len(t, result.FailedTables, 0)
}
