package usecases

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
)

// #==============================================================#
// ##          Mock SQL Types                                    ##
// #==============================================================#

// MockRows は sql.Rows のモック
type MockRows struct {
	columns []string
	data    [][]any
	index   int
	closed  bool
}

func NewMockRows(columns []string, data [][]any) *MockRows {
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

func (m *MockRows) Scan(dest ...any) error {
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
			case *any:
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
	data []any
	err  error
}

func NewMockRow(data []any, err error) *MockRow {
	return &MockRow{data: data, err: err}
}

func (m *MockRow) Scan(dest ...any) error {
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
			case *any:
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
	QueryContextRowsFunc   func(ctx context.Context, query string, args ...any) (model.RowsInterface, error)
	QueryRowContextRowFunc func(ctx context.Context, query string, args ...any) model.RowInterface
}

func (m *MockDatabaseQueryExecutor) QueryContextRows(ctx context.Context, query string, args ...any) (model.RowsInterface, error) {
	if m.QueryContextRowsFunc != nil {
		return m.QueryContextRowsFunc(ctx, query, args)
	}
	return nil, nil
}

func (m *MockDatabaseQueryExecutor) QueryRowContextRow(ctx context.Context, query string, args ...any) model.RowInterface {
	if m.QueryRowContextRowFunc != nil {
		return m.QueryRowContextRowFunc(ctx, query, args)
	}
	return nil
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
	suite.mockWriter.MkdirAllFunc = func(path string, perm os.FileMode) error {
		assert.Equal(t, outputPath, path)
		return nil
	}
	suite.mockWriter.CreateFunc = func(name string) (*os.File, error) {
		return os.CreateTemp("", "dump-all-tables-empty-*.tmp")
	}
	suite.mockWriter.CreateFunc = func(name string) (*os.File, error) {
		return os.CreateTemp("", "dump-all-tables-*.tmp")
	}

	// テーブル一覧取得の成功レスポンス
	// データベース名取得の成功レスポンス
	dbRow := NewMockRow([]any{"testdb"}, nil)

	// 各テーブルのダンプクエリの成功レスポンス
	usersRows := func() model.RowsInterface {
		return NewMockRows([]string{"id", "name"}, [][]any{
			{1, "John"},
			{2, "Jane"},
		})
	}

	productsRows := func() model.RowsInterface {
		return NewMockRows([]string{"id", "title"}, [][]any{
			{1, "Product A"},
		})
	}

	suite.mockExecutor.QueryContextRowsFunc = func(ctx context.Context, query string, args ...any) (model.RowsInterface, error) {
		if query == `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	` {
			return NewMockRows([]string{"table_name"}, [][]any{
				{"users"},
				{"products"},
			}), nil
		} else if query == "SELECT * FROM \"users\"" {
			return usersRows(), nil
		} else if query == "SELECT * FROM \"products\"" {
			return productsRows(), nil
		}
		return nil, errors.New("unexpected query")
	}

	suite.mockExecutor.QueryRowContextRowFunc = func(ctx context.Context, query string, args ...any) model.RowInterface {
		if query == "SELECT current_database()" {
			return dbRow
		}
		return nil
	}

	// ファイル書き込みの成功レスポンス
	suite.mockWriter.WriteFileFunc = func(filename string, data []byte, perm os.FileMode) error {
		return nil
	}

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

func TestTableDumper_DumpAllTables_EmptyDatabase(t *testing.T) {
	// Arrange
	suite := &TestTableDumper_DumpAllTables{}
	suite.setup()

	ctx := context.Background()
	outputPath := "/test/output"
	format := "json"
	concurrency := 2

	// モックの設定
	suite.mockWriter.MkdirAllFunc = func(path string, perm os.FileMode) error {
		assert.Equal(t, outputPath, path)
		return nil
	}

	// データベース名取得
	dbRow := NewMockRow([]any{"testdb"}, nil)

	suite.mockExecutor.QueryContextRowsFunc = func(ctx context.Context, query string, args ...any) (model.RowsInterface, error) {
		if query == `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	` {
			return NewMockRows([]string{"table_name"}, [][]any{}), nil
		}
		return nil, errors.New("unexpected query")
	}

	suite.mockExecutor.QueryRowContextRowFunc = func(ctx context.Context, query string, args ...any) model.RowInterface {
		if query == "SELECT current_database()" {
			return dbRow
		}
		return nil
	}

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
