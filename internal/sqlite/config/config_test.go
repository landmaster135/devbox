package config

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig_Normal(t *testing.T) {
	t.Parallel()

	cfg, err := NewConfig(OperationListTables, "/tmp/test.db", FormatText)
	require.NoError(t, err)
	require.Equal(t, OperationListTables, cfg.Operation)
	require.Equal(t, "/tmp/test.db", cfg.DBPath)
	require.Equal(t, FormatText, cfg.Format)
}

func TestNewConfig_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		op      string
		dbPath  string
		format  string
		wantErr string
	}{
		{
			name:    "missing operation",
			op:      "",
			dbPath:  "/tmp/test.db",
			format:  FormatText,
			wantErr: "--operation は必須です",
		},
		{
			name:    "invalid operation",
			op:      "unknown",
			dbPath:  "/tmp/test.db",
			format:  FormatText,
			wantErr: "未対応の operation",
		},
		{
			name:    "missing db-path",
			op:      OperationListTables,
			dbPath:  "",
			format:  FormatText,
			wantErr: "--db-path は必須です",
		},
		{
			name:    "invalid format",
			op:      OperationListTables,
			dbPath:  "/tmp/test.db",
			format:  "xml",
			wantErr: "--format は text または json",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewConfig(tt.op, tt.dbPath, tt.format)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestParseFlagsFromArgs_Normal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		args      []string
		operation string
		format    string
	}{
		{
			name:      "operation flag",
			args:      []string{"--operation=list-tables", "--db-path=/tmp/test.db", "--format=json"},
			operation: OperationListTables,
			format:    FormatJSON,
		},
		{
			name:      "opearation compatibility flag",
			args:      []string{"--opearation=list-tables", "--db-path=/tmp/test.db"},
			operation: OperationListTables,
			format:    FormatText,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg, err := ParseFlagsFromArgs(tt.args)
			require.NoError(t, err)
			require.Equal(t, tt.operation, cfg.Operation)
			require.Equal(t, "/tmp/test.db", cfg.DBPath)
			require.Equal(t, tt.format, cfg.Format)
		})
	}
}

func TestParseFlagsFromArgs_Help(t *testing.T) {
	t.Parallel()

	cfg, err := ParseFlagsFromArgs([]string{"--help"})
	require.NoError(t, err)
	require.True(t, cfg.Help)
}

func TestParseFlags_Normal(t *testing.T) {
	t.Parallel()

	originalArgs := os.Args
	t.Cleanup(func() {
		os.Args = originalArgs
	})
	os.Args = []string{"sqlite", "--operation=list-tables", "--db-path=/tmp/test.db"}

	cfg, err := ParseFlags()
	require.NoError(t, err)
	require.Equal(t, OperationListTables, cfg.Operation)
	require.Equal(t, "/tmp/test.db", cfg.DBPath)
}

func TestPrintUsageTo_Normal(t *testing.T) {
	t.Parallel()

	buf := bytes.NewBuffer(nil)
	PrintUsageTo(buf)

	output := buf.String()
	require.Contains(t, output, "SQLite CLIツール")
	require.Contains(t, output, "--operation")
	require.Contains(t, output, "--db-path")
}

func TestPrintUsage_Normal(t *testing.T) {
	t.Parallel()

	originalStderr := os.Stderr
	readPipe, writePipe, err := os.Pipe()
	require.NoError(t, err)

	os.Stderr = writePipe
	t.Cleanup(func() {
		os.Stderr = originalStderr
	})

	PrintUsage()
	require.NoError(t, writePipe.Close())

	data, readErr := io.ReadAll(readPipe)
	require.NoError(t, readErr)
	require.Contains(t, string(data), "SQLite CLIツール")
}
