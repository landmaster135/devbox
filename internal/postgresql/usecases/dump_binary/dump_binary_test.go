package dump_binary

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	infrastructures "github.com/landmaster135/devbox/internal/postgresql/infrastructures"
	writer "github.com/landmaster135/devbox/internal/postgresql/usecases/writer"
)

func TestDumper_DumpDatabase_Normal(t *testing.T) {
	mockExecutor := &infrastructures.MockCommandExecutor{}
	mockWriter := &writer.MockFileWriter{}
	dumper := NewDumperWithDependencies(mockExecutor, mockWriter, "")

	outputDir := t.TempDir()
	mockWriter.MkdirAllFunc = func(path string, perm os.FileMode) error {
		return nil
	}

	var capturedArgs []string
	mockExecutor.ExecuteFunc = func(name string, args ...string) ([]byte, error) {
		require.Equal(t, "pg_dump", name)
		capturedArgs = append([]string(nil), args...)
		return []byte("ok"), nil
	}

	result, err := dumper.DumpDatabase(context.Background(), "postgres://user:pass@localhost:5432/testdb", outputDir)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, capturedArgs, 5)
	assert.Equal(t, "-Fc", capturedArgs[0])
	assert.Equal(t, "--dbname", capturedArgs[1])
	assert.Equal(t, "postgres://user:pass@localhost:5432/testdb", capturedArgs[2])
	assert.Equal(t, "-f", capturedArgs[3])
	assert.Equal(t, filepath.Join(outputDir, result.FileName), capturedArgs[4])
	assert.Equal(t, "all_tables", result.TableName)
	assert.Equal(t, "binary", result.Format)
	assert.True(t, filepath.Ext(result.FileName) == ".dump")
}

func TestDumper_DumpDatabase_StripsStatementCacheModeForPgDump(t *testing.T) {
	mockExecutor := &infrastructures.MockCommandExecutor{}
	dumper := NewDumperWithDependencies(mockExecutor, &writer.MockFileWriter{}, "")

	var capturedArgs []string
	mockExecutor.ExecuteFunc = func(name string, args ...string) ([]byte, error) {
		capturedArgs = append([]string(nil), args...)
		return []byte("ok"), nil
	}

	_, err := dumper.DumpDatabase(
		context.Background(),
		"postgresql://user:pass@host/db?sslmode=require&statement_cache_mode=describe",
		t.TempDir(),
	)
	require.NoError(t, err)
	require.Len(t, capturedArgs, 5)
	assert.Equal(t, "postgresql://user:pass@host/db?sslmode=require", capturedArgs[2])
}

func TestDumper_DumpDatabase_DatabaseURLEmpty_Error(t *testing.T) {
	dumper := NewDumperWithDependencies(&infrastructures.MockCommandExecutor{}, &writer.MockFileWriter{}, "")

	result, err := dumper.DumpDatabase(context.Background(), "", t.TempDir())
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database-url が設定されていません")
}

func TestDumper_DumpDatabase_CommandError(t *testing.T) {
	mockExecutor := &infrastructures.MockCommandExecutor{
		ExecuteFunc: func(name string, args ...string) ([]byte, error) {
			return []byte("permission denied"), errors.New("exit status 1")
		},
	}
	dumper := NewDumperWithDependencies(mockExecutor, &writer.MockFileWriter{}, "")

	result, err := dumper.DumpDatabase(context.Background(), "postgres://user:pass@localhost:5432/testdb", t.TempDir())
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "pg_dump の実行に失敗しました")
	assert.Contains(t, err.Error(), "permission denied")
}

func TestResolveDatabaseName_Normal(t *testing.T) {
	assert.Equal(t, "testdb", resolveDatabaseName("postgres://user:pass@localhost:5432/testdb"))
	assert.Equal(t, "database", resolveDatabaseName("postgres://user:pass@localhost:5432/"))
	assert.Equal(t, "database", resolveDatabaseName("::invalid::"))
}

func TestSanitizeDatabaseURLForPgDump_Normal(t *testing.T) {
	assert.Equal(
		t,
		"postgresql://user:pass@host/db?sslmode=require",
		sanitizeDatabaseURLForPgDump("postgresql://user:pass@host/db?sslmode=require&statement_cache_mode=describe"),
	)
	assert.Equal(
		t,
		"postgresql://user:pass@host/db?sslmode=require",
		sanitizeDatabaseURLForPgDump("postgresql://user:pass@host/db?statement_cache_mode=describe&sslmode=require"),
	)
	assert.Equal(t, "not-a-url", sanitizeDatabaseURLForPgDump("not-a-url"))
}
