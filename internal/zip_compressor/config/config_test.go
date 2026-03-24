package config

import (
	"errors"
	"testing"

	flagParser "github.com/landmaster135/devbox/internal/zip_compressor/infrastructures/flag_parser"
)

func TestNewConfig_Normal(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		path      string
		wantErr   bool
	}{
		{
			name:      "正常ケース_compress",
			operation: "compress",
			path:      "/path/to/file",
			wantErr:   false,
		},
		{
			name:      "正常ケース_decompress",
			operation: "decompress",
			path:      "/path/to/archive.zip",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := NewConfig(tc.operation, tc.path)
			if (err != nil) != tc.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !tc.wantErr {
				if cfg.Operation != tc.operation {
					t.Errorf("NewConfig() Operation = %v, want %v", cfg.Operation, tc.operation)
				}
				if cfg.Path != tc.path {
					t.Errorf("NewConfig() Path = %v, want %v", cfg.Path, tc.path)
				}
			}
		})
	}
}

func TestNewConfig_Error(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		path      string
		wantErr   string
	}{
		{
			name:      "操作タイプが空",
			operation: "",
			path:      "/path/to/file",
			wantErr:   "操作タイプが指定されていません",
		},
		{
			name:      "無効な操作タイプ",
			operation: "invalid",
			path:      "/path/to/file",
			wantErr:   "無効な操作タイプです: invalid",
		},
		{
			name:      "パスが空",
			operation: "compress",
			path:      "",
			wantErr:   "パスが指定されていません",
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewConfig(tc.operation, tc.path)
			if err == nil {
				t.Errorf("NewConfig() expected error but got nil")
				return
			}
			if err.Error() != tc.wantErr {
				t.Errorf("NewConfig() error = %v, want %v", err.Error(), tc.wantErr)
			}
		})
	}
}

func TestParseFlagsWithParser_Normal(t *testing.T) {
	tests := []struct {
		name         string
		operation    string
		path         string
		help         bool
		args         []string
		expectedOp   string
		expectedPath string
		expectedHelp bool
	}{
		{
			name:         "フラグ指定_compress",
			operation:    "compress",
			path:         "/path/to/file",
			help:         false,
			args:         []string{},
			expectedOp:   "compress",
			expectedPath: "/path/to/file",
			expectedHelp: false,
		},
		{
			name:         "フラグ指定_decompress",
			operation:    "decompress",
			path:         "/path/to/archive.zip",
			help:         false,
			args:         []string{},
			expectedOp:   "decompress",
			expectedPath: "/path/to/archive.zip",
			expectedHelp: false,
		},
		{
			name:         "位置引数指定",
			operation:    "",
			path:         "",
			help:         false,
			args:         []string{"compress", "/path/to/file"},
			expectedOp:   "compress",
			expectedPath: "/path/to/file",
			expectedHelp: false,
		},
		{
			name:         "ヘルプ指定",
			operation:    "",
			path:         "",
			help:         true,
			args:         []string{},
			expectedOp:   "",
			expectedPath: "",
			expectedHelp: true,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			parser := flagParser.NewMockFlagParser()
			parser.SetArgs(tc.args)

			parser.SetStringFlag("operation", tc.operation)
			parser.SetStringFlag("o", tc.operation)
			parser.SetStringFlag("path", tc.path)
			parser.SetStringFlag("p", tc.path)
			parser.SetBoolFlag("help", tc.help)
			parser.SetBoolFlag("h", tc.help)

			cfg, err := ParseFlagsWithParser(parser)
			if err != nil {
				t.Errorf("ParseFlagsWithParser() error = %v", err)
				return
			}

			if tc.expectedHelp {
				if !cfg.Help {
					t.Errorf("ParseFlagsWithParser() Help = %v, want %v", cfg.Help, tc.expectedHelp)
				}
				return
			}

			if cfg.Operation != tc.expectedOp {
				t.Errorf("ParseFlagsWithParser() Operation = %v, want %v", cfg.Operation, tc.expectedOp)
			}
			if cfg.Path != tc.expectedPath {
				t.Errorf("ParseFlagsWithParser() Path = %v, want %v", cfg.Path, tc.expectedPath)
			}
		})
	}
}

func TestParseFlagsWithParser_ParseError(t *testing.T) {
	t.Parallel()

	parser := flagParser.NewMockFlagParser()
	parser.SetParseError(errors.New("parse failed"))

	_, err := ParseFlagsWithParser(parser)
	if err == nil {
		t.Fatalf("ParseFlagsWithParser() error = nil, want non-nil")
	}
	if err.Error() != "フラグの解析に失敗しました: parse failed" {
		t.Fatalf("ParseFlagsWithParser() error = %q", err.Error())
	}
}
