package dump_binary

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	infrastructures "github.com/landmaster135/devbox/internal/postgresql/infrastructures"
	writer "github.com/landmaster135/devbox/internal/postgresql/usecases/writer"
)

func noWaitRetryConfig(maxAttempts int) RetryConfig {
	return RetryConfig{
		MaxAttempts: maxAttempts,
		BaseDelay:   time.Millisecond,
		MaxDelay:    time.Millisecond,
		SleepWithContext: func(ctx context.Context, d time.Duration) error {
			return nil
		},
	}
}

func TestDumper_DumpDatabase_Normal(t *testing.T) {
	mockExecutor := &infrastructures.MockCommandExecutor{}
	mockWriter := &writer.MockFileWriter{}
	dumper := NewDumperWithDependencies(mockExecutor, mockWriter, "", RetryConfig{})

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

	result, err := dumper.DumpDatabase(context.Background(), "postgres://user:pass@localhost:5432/testdb", outputDir, nil)
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

func TestDumper_DumpDatabase_ExcludeTableData(t *testing.T) {
	mockExecutor := &infrastructures.MockCommandExecutor{}
	mockWriter := &writer.MockFileWriter{}
	dumper := NewDumperWithDependencies(mockExecutor, mockWriter, "", RetryConfig{})

	outputDir := t.TempDir()
	var capturedArgs []string
	mockExecutor.ExecuteFunc = func(name string, args ...string) ([]byte, error) {
		require.Equal(t, "pg_dump", name)
		capturedArgs = append([]string(nil), args...)
		return []byte("ok"), nil
	}

	result, err := dumper.DumpDatabase(
		context.Background(),
		"postgres://user:pass@localhost:5432/testdb",
		outputDir,
		[]string{"public.attachments"},
	)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, capturedArgs, 6)
	assert.Equal(t, "-Fc", capturedArgs[0])
	assert.Equal(t, "--dbname", capturedArgs[1])
	assert.Equal(t, "postgres://user:pass@localhost:5432/testdb", capturedArgs[2])
	assert.Equal(t, "--exclude-table-data=public.attachments", capturedArgs[3])
	assert.Equal(t, "-f", capturedArgs[4])
	assert.Equal(t, filepath.Join(outputDir, result.FileName), capturedArgs[5])
}

func TestDumper_DumpTableDataAsSQL_Normal(t *testing.T) {
	mockExecutor := &infrastructures.MockCommandExecutor{}
	mockWriter := &writer.MockFileWriter{}
	dumper := NewDumperWithDependencies(mockExecutor, mockWriter, "", RetryConfig{})

	outputDir := t.TempDir()
	var capturedArgs []string
	mockExecutor.ExecuteFunc = func(name string, args ...string) ([]byte, error) {
		require.Equal(t, "pg_dump", name)
		capturedArgs = append([]string(nil), args...)
		return []byte("ok"), nil
	}

	result, err := dumper.DumpTableDataAsSQL(
		context.Background(),
		"postgres://user:pass@localhost:5432/testdb",
		outputDir,
		"public.attachments",
	)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, capturedArgs, 8)
	assert.Equal(t, "--data-only", capturedArgs[0])
	assert.Equal(t, "--table", capturedArgs[1])
	assert.Equal(t, "public.attachments", capturedArgs[2])
	assert.Equal(t, "--column-inserts", capturedArgs[3])
	assert.Equal(t, "--dbname", capturedArgs[4])
	assert.Equal(t, "postgres://user:pass@localhost:5432/testdb", capturedArgs[5])
	assert.Equal(t, "-f", capturedArgs[6])
	assert.Equal(t, filepath.Join(outputDir, result.FileName), capturedArgs[7])
	assert.Equal(t, "public.attachments", result.TableName)
	assert.Equal(t, "sql", result.Format)
	assert.True(t, filepath.Ext(result.FileName) == ".sql")
	assert.True(t, filepath.Base(result.FileName)[:12] == "attachments_")
}

func TestDumper_DumpDatabase_StripsStatementCacheModeForPgDump(t *testing.T) {
	mockExecutor := &infrastructures.MockCommandExecutor{}
	dumper := NewDumperWithDependencies(mockExecutor, &writer.MockFileWriter{}, "", RetryConfig{})

	var capturedArgs []string
	mockExecutor.ExecuteFunc = func(name string, args ...string) ([]byte, error) {
		capturedArgs = append([]string(nil), args...)
		return []byte("ok"), nil
	}

	_, err := dumper.DumpDatabase(
		context.Background(),
		"postgresql://user:pass@host/db?sslmode=require&statement_cache_mode=describe",
		t.TempDir(),
		nil,
	)
	require.NoError(t, err)
	require.Len(t, capturedArgs, 5)
	assert.Equal(t, "postgresql://user:pass@host/db?sslmode=require", capturedArgs[2])
}

func TestDumper_DumpDatabase_DatabaseURLEmpty_Error(t *testing.T) {
	dumper := NewDumperWithDependencies(&infrastructures.MockCommandExecutor{}, &writer.MockFileWriter{}, "", RetryConfig{})

	result, err := dumper.DumpDatabase(context.Background(), "", t.TempDir(), nil)
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
	dumper := NewDumperWithDependencies(mockExecutor, &writer.MockFileWriter{}, "", noWaitRetryConfig(3))

	result, err := dumper.DumpDatabase(context.Background(), "postgres://user:pass@localhost:5432/testdb", t.TempDir(), nil)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "pg_dump の実行に失敗しました")
	assert.Contains(t, err.Error(), "(attempts=1)")
	assert.Contains(t, err.Error(), "permission denied")
	assert.Len(t, mockExecutor.Calls, 1)
}

func TestDumper_DumpDatabase_RetryOnRetriableErrorThenSuccess(t *testing.T) {
	mockExecutor := &infrastructures.MockCommandExecutor{}
	mockWriter := &writer.MockFileWriter{}
	dumper := NewDumperWithDependencies(mockExecutor, mockWriter, "", noWaitRetryConfig(3))

	attempt := 0
	mockExecutor.ExecuteFunc = func(name string, args ...string) ([]byte, error) {
		attempt++
		if attempt == 1 {
			return []byte("pg_dump: error: connection to server failed: ERROR:  Control plane request failed"), errors.New("exit status 1")
		}
		return []byte("ok"), nil
	}

	result, err := dumper.DumpDatabase(context.Background(), "postgres://user:pass@localhost:5432/testdb", t.TempDir(), nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, mockExecutor.Calls, 2)
}

func TestDumper_DumpDatabase_RetryExhaustedOnRetriableError(t *testing.T) {
	mockExecutor := &infrastructures.MockCommandExecutor{
		ExecuteFunc: func(name string, args ...string) ([]byte, error) {
			return []byte("pg_dump: error: connection to server failed: ERROR:  Control plane request failed"), errors.New("exit status 1")
		},
	}
	dumper := NewDumperWithDependencies(mockExecutor, &writer.MockFileWriter{}, "", noWaitRetryConfig(3))

	result, err := dumper.DumpDatabase(context.Background(), "postgres://user:pass@localhost:5432/testdb", t.TempDir(), nil)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "(attempts=3)")
	assert.Contains(t, err.Error(), "Control plane request failed")
	assert.Len(t, mockExecutor.Calls, 3)
}

func TestDumper_DumpDatabase_ContextCanceledDuringRetryWait(t *testing.T) {
	mockExecutor := &infrastructures.MockCommandExecutor{
		ExecuteFunc: func(name string, args ...string) ([]byte, error) {
			return []byte("pg_dump: error: connection to server failed: ERROR:  Control plane request failed"), errors.New("exit status 1")
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	retryConfig := RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   time.Second,
		MaxDelay:    time.Second,
		SleepWithContext: func(ctx context.Context, d time.Duration) error {
			cancel()
			<-ctx.Done()
			return ctx.Err()
		},
	}
	dumper := NewDumperWithDependencies(mockExecutor, &writer.MockFileWriter{}, "", retryConfig)

	result, err := dumper.DumpDatabase(ctx, "postgres://user:pass@localhost:5432/testdb", t.TempDir(), nil)
	require.ErrorIs(t, err, context.Canceled)
	assert.Nil(t, result)
	assert.Len(t, mockExecutor.Calls, 1)
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
