package config

import (
	"fmt"
	"testing"
)

type MockFlagParser struct {
	stringValues map[string]string
	boolValues   map[string]bool
	parseError   error
}

func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringValues: map[string]string{},
		boolValues:   map[string]bool{},
	}
}

func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	if v, ok := m.stringValues[name]; ok {
		*p = v
		return
	}
	if *p != "" {
		return
	}
	*p = value
}

func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if v, ok := m.boolValues[name]; ok {
		*p = v
		return
	}
	if *p {
		return
	}
	*p = value
}

func (m *MockFlagParser) Parse() error {
	return m.parseError
}

func (m *MockFlagParser) SetString(name, value string) {
	m.stringValues[name] = value
}

func (m *MockFlagParser) SetBool(name string, value bool) {
	m.boolValues[name] = value
}

func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}

func TestNewConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		operation string
		pageType  string
		srcJSON   string
		srcBody   string
		outDir    string
		wantErr   string
	}{
		{
			name:      "normal",
			operation: OperationDistributeFiles,
			pageType:  PageTypeContent,
			srcJSON:   "/tmp/contents.json",
			srcBody:   "/tmp/body",
			outDir:    "/tmp/out",
		},
		{
			name:      "missing operation",
			operation: "",
			pageType:  PageTypeContent,
			srcJSON:   "/tmp/contents.json",
			srcBody:   "/tmp/body",
			outDir:    "/tmp/out",
			wantErr:   "operation パラメータは必須です",
		},
		{
			name:      "invalid operation",
			operation: "unknown",
			pageType:  PageTypeContent,
			srcJSON:   "/tmp/contents.json",
			srcBody:   "/tmp/body",
			outDir:    "/tmp/out",
			wantErr:   "未対応のoperationです: unknown",
		},
		{
			name:      "missing page type",
			operation: OperationDistributeFiles,
			pageType:  "",
			srcJSON:   "/tmp/contents.json",
			srcBody:   "/tmp/body",
			outDir:    "/tmp/out",
			wantErr:   "page-type パラメータは必須です",
		},
		{
			name:      "invalid page type",
			operation: OperationDistributeFiles,
			pageType:  "artifact",
			srcJSON:   "/tmp/contents.json",
			srcBody:   "/tmp/body",
			outDir:    "/tmp/out",
			wantErr:   "未対応のpage-typeです: artifact",
		},
		{
			name:      "missing src json",
			operation: OperationDistributeFiles,
			pageType:  PageTypeContent,
			srcJSON:   "",
			srcBody:   "/tmp/body",
			outDir:    "/tmp/out",
			wantErr:   "src-json-path パラメータは必須です",
		},
		{
			name:      "missing src body dir",
			operation: OperationDistributeFiles,
			pageType:  PageTypeContent,
			srcJSON:   "/tmp/contents.json",
			srcBody:   "",
			outDir:    "/tmp/out",
			wantErr:   "src-body-dir パラメータは必須です",
		},
		{
			name:      "missing out dir",
			operation: OperationDistributeFiles,
			pageType:  PageTypeContent,
			srcJSON:   "/tmp/contents.json",
			srcBody:   "/tmp/body",
			outDir:    "",
			wantErr:   "out-dir パラメータは必須です",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewConfig(tt.operation, tt.pageType, tt.srcJSON, tt.srcBody, tt.outDir)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("error = %v, want %q", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Operation != OperationDistributeFiles {
				t.Fatalf("Operation = %q", got.Operation)
			}
			if got.PageType != PageTypeContent {
				t.Fatalf("PageType = %q", got.PageType)
			}
		})
	}
}

func TestParseFlagsWithParser(t *testing.T) {
	t.Parallel()

	t.Run("normal", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetString("operation", OperationDistributeFiles)
		parser.SetString("page-type", PageTypeContent)
		parser.SetString("src-json-path", "/tmp/contents.json")
		parser.SetString("src-body-dir", "/tmp/body")
		parser.SetString("out-dir", "/tmp/out")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Operation != OperationDistributeFiles {
			t.Fatalf("Operation = %q", cfg.Operation)
		}
	})

	t.Run("help", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetBool("help", true)

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !cfg.Help {
			t.Fatalf("Help = false, want true")
		}
	})

	t.Run("parse error", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetParseError(fmt.Errorf("bad flag"))

		_, err := ParseFlagsWithParser(parser)
		if err == nil || err.Error() != "フラグの解析に失敗しました: bad flag" {
			t.Fatalf("error = %v", err)
		}
	})
}
