package dump

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
)

// #==============================================================#
// ##          validateOptions Tests                             ##
// #==============================================================#

func TestTableDumper_validateOptions_Normal(t *testing.T) {
	// Arrange
	options := DumpOptions{
		TableName:  "users",
		OutputPath: "/tmp/dump",
		Format:     "json",
		Limit:      nil,
	}

	// Act
	err := validateOptions(&options)

	// Assert
	assert.NoError(t, err)
}

func TestTableDumper_validateOptions_EmptyTableName(t *testing.T) {
	// Arrange
	options := DumpOptions{
		TableName:  "",
		OutputPath: "/tmp/dump",
		Format:     "json",
		Limit:      nil,
	}

	// Act
	err := validateOptions(&options)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "テーブル名が指定されていません")
}

func TestTableDumper_validateOptions_InvalidFormat(t *testing.T) {
	// Arrange
	options := DumpOptions{
		TableName:  "users",
		OutputPath: "/tmp/dump",
		Format:     "invalid",
		Limit:      nil,
	}

	// Act
	err := validateOptions(&options)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "サポートされていないフォーマットです")
}

func TestTableDumper_validateOptions_InvalidTableNameCharacters(t *testing.T) {
	// Arrange
	options := DumpOptions{
		TableName:  "users;DROP",
		OutputPath: "/tmp/dump",
		Format:     "json",
		Limit:      nil,
	}

	// Act
	err := validateOptions(&options)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "テーブル名に使用できない文字が含まれています")
}

func TestTableDumper_validateOptions_InvalidLimit(t *testing.T) {
	// Arrange
	invalidLimit := -1
	options := DumpOptions{
		TableName:  "users",
		OutputPath: "/tmp/dump",
		Format:     "json",
		Limit:      &invalidLimit,
	}

	// Act
	err := validateOptions(&options)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "limitは正の数である必要があります")
}

func TestTableDumper_validateOptions_PathTraversal(t *testing.T) {
	// Arrange
	options := DumpOptions{
		TableName:  "users",
		OutputPath: "../../../etc",
		Format:     "json",
		Limit:      nil,
	}

	// Act
	err := validateOptions(&options)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "無効なパスが指定されました")
}

// #==============================================================#
// ##          ensureAllowedTable Tests                          ##
// #==============================================================#

func TestTableDumper_EnsureAllowedTable(t *testing.T) {
	mockExecutor := &MockDatabaseQueryExecutor{}
	mockExecutor.QueryContextRowsFunc = func(ctx context.Context, query string, args ...any) (model.RowsInterface, error) {
		return NewMockRows([]string{"table_name"}, [][]any{{"users"}}), nil
	}

	dumper := NewTableDumperWithDependencies(mockExecutor, &MockFileWriter{})

	err := dumper.ensureAllowedTable(context.Background(), "users")
	assert.NoError(t, err)
}

func TestTableDumper_EnsureAllowedTable_InvalidSchema(t *testing.T) {
	mockExecutor := &MockDatabaseQueryExecutor{}
	dumper := NewTableDumperWithDependencies(mockExecutor, &MockFileWriter{})

	err := dumper.ensureAllowedTable(context.Background(), "admin.users")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "サポートされていないスキーマ")
}

func TestTableDumper_EnsureAllowedTable_NotFound(t *testing.T) {
	mockExecutor := &MockDatabaseQueryExecutor{}
	mockExecutor.QueryContextRowsFunc = func(ctx context.Context, query string, args ...any) (model.RowsInterface, error) {
		return NewMockRows([]string{"table_name"}, [][]any{{"users"}}), nil
	}

	dumper := NewTableDumperWithDependencies(mockExecutor, &MockFileWriter{})

	err := dumper.ensureAllowedTable(context.Background(), "orders")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "指定されたテーブルが存在しません")
}
