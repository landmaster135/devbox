package config

import (
	"fmt"
	"testing"
)

type MockFlagParser struct {
	stringValues map[string]string
	intValues    map[string]int
	boolValues   map[string]bool
	parseError   error
}

func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringValues: map[string]string{},
		intValues:    map[string]int{},
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

func (m *MockFlagParser) IntVar(p *int, name string, value int, usage string) {
	if v, ok := m.intValues[name]; ok {
		*p = v
		return
	}
	if *p != 0 {
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

func (m *MockFlagParser) SetInt(name string, value int) {
	m.intValues[name] = value
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
		name           string
		operation      string
		pageType       string
		category       string
		skipsNoSrcBody bool
		srcJSONFile    string
		srcBodyDir     string
		outDir         string
		conNumberStart int
		conNumberEnd   int
		threshold      int
		wantPageType   string
		wantErr        string
	}{
		{
			name:         "distribute normal",
			operation:    OperationDistributeFiles,
			pageType:     PageTypeContent,
			srcJSONFile:  "/tmp/contents.json",
			srcBodyDir:   "/tmp/body",
			outDir:       "/tmp/out",
			wantPageType: PageTypeContent,
		},
		{
			name:           "craft normal",
			operation:      OperationCraftMarkdown,
			pageType:       PageTypeContent,
			srcJSONFile:    "/tmp/contents.json",
			srcBodyDir:     "/tmp/body",
			outDir:         "/tmp/out",
			conNumberStart: 1,
			conNumberEnd:   10,
			wantPageType:   PageTypeContent,
		},
		{
			name:         "check body length normal",
			operation:    OperationCheckBodyLength,
			srcBodyDir:   "/tmp/body",
			threshold:    1000,
			wantPageType: "",
		},
		{
			name:        "missing operation",
			operation:   "",
			pageType:    PageTypeContent,
			srcJSONFile: "/tmp/contents.json",
			srcBodyDir:  "/tmp/body",
			outDir:      "/tmp/out",
			wantErr:     "operation パラメータは必須です",
		},
		{
			name:        "invalid operation",
			operation:   "unknown",
			pageType:    PageTypeContent,
			srcJSONFile: "/tmp/contents.json",
			srcBodyDir:  "/tmp/body",
			outDir:      "/tmp/out",
			wantErr:     "未対応のoperationです: unknown",
		},
		{
			name:        "missing page type",
			operation:   OperationDistributeFiles,
			pageType:    "",
			srcJSONFile: "/tmp/contents.json",
			srcBodyDir:  "/tmp/body",
			outDir:      "/tmp/out",
			wantErr:     "page-type パラメータは必須です",
		},
		{
			name:        "invalid page type",
			operation:   OperationDistributeFiles,
			pageType:    "artifact",
			srcJSONFile: "/tmp/contents.json",
			srcBodyDir:  "/tmp/body",
			outDir:      "/tmp/out",
			wantErr:     "未対応のpage-typeです: artifact",
		},
		{
			name:        "missing src json file",
			operation:   OperationDistributeFiles,
			pageType:    PageTypeContent,
			srcJSONFile: "",
			srcBodyDir:  "/tmp/body",
			outDir:      "/tmp/out",
			wantErr:     "src-json-file パラメータは必須です",
		},
		{
			name:        "missing src body dir",
			operation:   OperationDistributeFiles,
			pageType:    PageTypeContent,
			srcJSONFile: "/tmp/contents.json",
			srcBodyDir:  "",
			outDir:      "/tmp/out",
			wantErr:     "src-body-dir パラメータは必須です",
		},
		{
			name:        "missing out dir",
			operation:   OperationDistributeFiles,
			pageType:    PageTypeContent,
			srcJSONFile: "/tmp/contents.json",
			srcBodyDir:  "/tmp/body",
			outDir:      "",
			wantErr:     "out-dir パラメータは必須です",
		},
		{
			name:      "check body length missing src body dir",
			operation: OperationCheckBodyLength,
			threshold: 1000,
			wantErr:   "src-body-dir パラメータは必須です",
		},
		{
			name:       "check body length invalid threshold",
			operation:  OperationCheckBodyLength,
			srcBodyDir: "/tmp/body",
			threshold:  -1,
			wantErr:    "threshold パラメータは0以上で必須です",
		},
		{
			name:           "craft missing con number start",
			operation:      OperationCraftMarkdown,
			pageType:       PageTypeContent,
			srcJSONFile:    "/tmp/contents.json",
			srcBodyDir:     "/tmp/body",
			outDir:         "/tmp/out",
			conNumberStart: 0,
			conNumberEnd:   10,
			wantPageType:   PageTypeContent,
			wantErr:        "con_number_start パラメータは1以上で必須です",
		},
		{
			name:           "craft missing con number end",
			operation:      OperationCraftMarkdown,
			pageType:       PageTypeContent,
			srcJSONFile:    "/tmp/contents.json",
			srcBodyDir:     "/tmp/body",
			outDir:         "/tmp/out",
			conNumberStart: 1,
			conNumberEnd:   0,
			wantPageType:   PageTypeContent,
			wantErr:        "con_number_end パラメータは1以上で必須です",
		},
		{
			name:           "craft invalid range",
			operation:      OperationCraftMarkdown,
			pageType:       PageTypeContent,
			srcJSONFile:    "/tmp/contents.json",
			srcBodyDir:     "/tmp/body",
			outDir:         "/tmp/out",
			conNumberStart: 10,
			conNumberEnd:   1,
			wantPageType:   PageTypeContent,
			wantErr:        "con_number_start は con_number_end 以下である必要があります",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewConfig(
				tt.operation,
				tt.pageType,
				tt.category,
				tt.skipsNoSrcBody,
				tt.srcJSONFile,
				tt.srcBodyDir,
				tt.outDir,
				tt.conNumberStart,
				tt.conNumberEnd,
				tt.threshold,
			)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("error = %v, want %q", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Operation != tt.operation {
				t.Fatalf("Operation = %q, want %q", got.Operation, tt.operation)
			}
			if got.PageType != tt.wantPageType {
				t.Fatalf("PageType = %q, want %q", got.PageType, tt.wantPageType)
			}
			if got.SkipsNoSrcBody != tt.skipsNoSrcBody {
				t.Fatalf("SkipsNoSrcBody = %v, want %v", got.SkipsNoSrcBody, tt.skipsNoSrcBody)
			}
		})
	}
}

func TestParseFlagsWithParser(t *testing.T) {
	t.Parallel()

	t.Run("distribute normal", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetString("operation", OperationDistributeFiles)
		parser.SetString("page-type", PageTypeContent)
		parser.SetString("src-json-file", "/tmp/contents.json")
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

	t.Run("craft normal", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetString("operation", OperationCraftMarkdown)
		parser.SetString("page-type", PageTypeContent)
		parser.SetString("src-json-file", "/tmp/contents.json")
		parser.SetString("src-body-dir", "/tmp/body")
		parser.SetString("out-dir", "/tmp/out")
		parser.SetInt("con_number_start", 100)
		parser.SetInt("con_number_end", 200)

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.ConNumberStart != 100 || cfg.ConNumberEnd != 200 {
			t.Fatalf("con range = (%d,%d), want (100,200)", cfg.ConNumberStart, cfg.ConNumberEnd)
		}
		if cfg.SkipsNoSrcBody {
			t.Fatalf("SkipsNoSrcBody = true, want false")
		}
	})

	t.Run("craft with category", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetString("operation", OperationCraftMarkdown)
		parser.SetString("page-type", PageTypeContent)
		parser.SetString("category", "software")
		parser.SetString("src-json-file", "/tmp/contents.json")
		parser.SetString("src-body-dir", "/tmp/body")
		parser.SetString("out-dir", "/tmp/out")
		parser.SetInt("con_number_start", 1)
		parser.SetInt("con_number_end", 100)
		parser.SetBool("skips-no-src-body", true)

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Category != "software" {
			t.Fatalf("Category = %q, want software", cfg.Category)
		}
		if !cfg.SkipsNoSrcBody {
			t.Fatalf("SkipsNoSrcBody = false, want true")
		}
	})

	t.Run("src json path alias", func(t *testing.T) {
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
		if cfg.SrcJSONFile != "/tmp/contents.json" {
			t.Fatalf("SrcJSONFile = %q", cfg.SrcJSONFile)
		}
	})

	t.Run("check body length normal", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetString("operation", OperationCheckBodyLength)
		parser.SetString("src-body-dir", "/tmp/body")
		parser.SetInt("threshold", 120)

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Operation != OperationCheckBodyLength {
			t.Fatalf("Operation = %q", cfg.Operation)
		}
		if cfg.Threshold != 120 {
			t.Fatalf("Threshold = %d, want 120", cfg.Threshold)
		}
	})

	t.Run("check body length missing threshold", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetString("operation", OperationCheckBodyLength)
		parser.SetString("src-body-dir", "/tmp/body")

		_, err := ParseFlagsWithParser(parser)
		if err == nil || err.Error() != "threshold パラメータは0以上で必須です" {
			t.Fatalf("error = %v", err)
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
