package main

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockSQLiteService struct {
	handleListTablesFunc func(ctx context.Context, dbPath, format string) (string, error)
}

func (m *mockSQLiteService) HandleListTables(ctx context.Context, dbPath, format string) (string, error) {
	return m.handleListTablesFunc(ctx, dbPath, format)
}

func TestRun_NormalListTables(t *testing.T) {
	t.Parallel()

	stdout := &testWriter{}
	stderr := &testWriter{}
	code := run(
		[]string{"--operation=list-tables", "--db-path=/tmp/test.db"},
		stdout,
		stderr,
		func() sqliteService {
			return &mockSQLiteService{
				handleListTablesFunc: func(ctx context.Context, dbPath, format string) (string, error) {
					require.Equal(t, "/tmp/test.db", dbPath)
					require.Equal(t, "text", format)
					return "posts\nusers", nil
				},
			}
		},
	)

	require.Equal(t, 0, code)
	require.Equal(t, "posts\nusers\n", stdout.String())
	require.Equal(t, "", stderr.String())
}

func TestRun_NormalCompatibilityFlag(t *testing.T) {
	t.Parallel()

	stdout := &testWriter{}
	stderr := &testWriter{}
	code := run(
		[]string{"--opearation=list-tables", "--db-path=/tmp/test.db", "--format=json"},
		stdout,
		stderr,
		func() sqliteService {
			return &mockSQLiteService{
				handleListTablesFunc: func(ctx context.Context, dbPath, format string) (string, error) {
					require.Equal(t, "json", format)
					return "[]", nil
				},
			}
		},
	)

	require.Equal(t, 0, code)
	require.Equal(t, "[]\n", stdout.String())
	require.Equal(t, "", stderr.String())
}

func TestRun_NormalHelp(t *testing.T) {
	t.Parallel()

	stdout := &testWriter{}
	stderr := &testWriter{}
	code := run([]string{"--help"}, stdout, stderr, func() sqliteService {
		t.Fatal("service factory should not be called when help requested")
		return nil
	})

	require.Equal(t, 0, code)
	require.Contains(t, stderr.String(), "SQLite CLIツール")
	require.Equal(t, "", stdout.String())
}

func TestRun_ErrorValidation(t *testing.T) {
	t.Parallel()

	stdout := &testWriter{}
	stderr := &testWriter{}
	code := run([]string{"--operation=list-tables"}, stdout, stderr, func() sqliteService {
		t.Fatal("service factory should not be called on validation error")
		return nil
	})

	require.Equal(t, 1, code)
	require.Contains(t, stderr.String(), "--db-path は必須です")
	require.Contains(t, stderr.String(), "SQLite CLIツール")
	require.Equal(t, "", stdout.String())
}

func TestRun_ErrorHandleListTables(t *testing.T) {
	t.Parallel()

	stdout := &testWriter{}
	stderr := &testWriter{}
	code := run(
		[]string{"--operation=list-tables", "--db-path=/tmp/test.db"},
		stdout,
		stderr,
		func() sqliteService {
			return &mockSQLiteService{
				handleListTablesFunc: func(ctx context.Context, dbPath, format string) (string, error) {
					return "", errors.New("test error")
				},
			}
		},
	)

	require.Equal(t, 1, code)
	require.Contains(t, stderr.String(), "エラー: test error")
	require.Equal(t, "", stdout.String())
}

func TestRun_ErrorWriteOutput(t *testing.T) {
	t.Parallel()

	stderr := &testWriter{}
	code := run(
		[]string{"--operation=list-tables", "--db-path=/tmp/test.db"},
		&errorWriter{},
		stderr,
		func() sqliteService {
			return &mockSQLiteService{
				handleListTablesFunc: func(ctx context.Context, dbPath, format string) (string, error) {
					return "users", nil
				},
			}
		},
	)

	require.Equal(t, 1, code)
	require.Contains(t, stderr.String(), "エラー: 出力に失敗しました")
}

func TestDefaultServiceFactory_Normal(t *testing.T) {
	t.Parallel()

	service := defaultServiceFactory()
	require.NotNil(t, service)
}

type testWriter struct {
	buf []byte
}

func (w *testWriter) Write(p []byte) (int, error) {
	w.buf = append(w.buf, p...)
	return len(p), nil
}

func (w *testWriter) String() string {
	return string(w.buf)
}

type errorWriter struct{}

func (w *errorWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write failed")
}
