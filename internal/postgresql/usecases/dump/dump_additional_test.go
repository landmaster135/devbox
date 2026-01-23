package usecases

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
)

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

func TestTableDumper_BuildQuery(t *testing.T) {
	dumper := NewTableDumper(&MockDatabaseQueryExecutor{})

	query, err := dumper.buildQuery(&DumpOptions{TableName: "users", Format: "json"})
	require.NoError(t, err)
	assert.Equal(t, "SELECT * FROM \"users\"", query)

	limit := 5
	query, err = dumper.buildQuery(&DumpOptions{TableName: "users", Format: "json", Limit: &limit})
	require.NoError(t, err)
	assert.Equal(t, "SELECT * FROM \"users\" LIMIT 5", query)
}

func TestStreamRows_WriterNil(t *testing.T) {
	rows := NewMockRows([]string{"id"}, [][]any{})

	dumper := NewTableDumper(&MockDatabaseQueryExecutor{})
	err := dumper.streamRows(rows, []string{"id"}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "writerが初期化されていません")
}
