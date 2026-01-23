package dump

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
)

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
