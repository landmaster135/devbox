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
		baseURL        string
		apiToken       string
		category       string
		skipsNoSrcBody bool
		srcJSONFile    string
		srcBodyDir     string
		srcResourceDir string
		outDir         string
		targetStr      string
		conNumberStart int
		conNumberEnd   int
		threshold      int
		wantPageType   string
		wantTargetStr  string
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
			name:           "rename bodies by category id normal",
			operation:      OperationRenameBodiesByCategoryID,
			pageType:       PageTypeContent,
			srcJSONFile:    "/tmp/contents.json",
			srcResourceDir: "/tmp/resource",
			conNumberStart: 100,
			conNumberEnd:   200,
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
			name:          "grep str normal",
			operation:     OperationGrepStr,
			srcBodyDir:    "/tmp/body",
			targetStr:     "  TODO  ",
			wantPageType:  "",
			wantTargetStr: "TODO",
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
			name:       "grep str missing target str",
			operation:  OperationGrepStr,
			srcBodyDir: "/tmp/body",
			targetStr:  " ",
			wantErr:    "target-str パラメータは必須です",
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
		{
			name:           "rename missing src resource dir",
			operation:      OperationRenameBodiesByCategoryID,
			pageType:       PageTypeContent,
			srcJSONFile:    "/tmp/contents.json",
			srcResourceDir: "",
			conNumberStart: 1,
			conNumberEnd:   10,
			wantErr:        "src-resource-dir パラメータは必須です",
		},
		{
			name:           "rename invalid range",
			operation:      OperationRenameBodiesByCategoryID,
			pageType:       PageTypeContent,
			srcJSONFile:    "/tmp/contents.json",
			srcResourceDir: "/tmp/resource",
			conNumberStart: 10,
			conNumberEnd:   1,
			wantErr:        "con_number_start は con_number_end 以下である必要があります",
		},
		{
			name:           "migrate to memos normal",
			operation:      OperationMigrateToMemos,
			pageType:       PageTypeContent,
			baseURL:        "https://memos.example.com",
			apiToken:       "token",
			srcBodyDir:     "/tmp/body",
			srcResourceDir: "/tmp/resource",
			wantPageType:   PageTypeContent,
		},
		{
			name:           "migrate to memos missing base url",
			operation:      OperationMigrateToMemos,
			pageType:       PageTypeContent,
			baseURL:        " ",
			apiToken:       "token",
			srcBodyDir:     "/tmp/body",
			srcResourceDir: "/tmp/resource",
			wantErr:        "base-url パラメータは必須です",
		},
		{
			name:           "migrate to memos missing api token",
			operation:      OperationMigrateToMemos,
			pageType:       PageTypeContent,
			baseURL:        "https://memos.example.com",
			apiToken:       " ",
			srcBodyDir:     "/tmp/body",
			srcResourceDir: "/tmp/resource",
			wantErr:        "api-token パラメータは必須です",
		},
		{
			name:           "migrate to memos missing src body dir",
			operation:      OperationMigrateToMemos,
			pageType:       PageTypeContent,
			baseURL:        "https://memos.example.com",
			apiToken:       "token",
			srcBodyDir:     " ",
			srcResourceDir: "/tmp/resource",
			wantErr:        "src-body-dir パラメータは必須です",
		},
		{
			name:           "migrate to memos missing src resource dir",
			operation:      OperationMigrateToMemos,
			pageType:       PageTypeContent,
			baseURL:        "https://memos.example.com",
			apiToken:       "token",
			srcBodyDir:     "/tmp/body",
			srcResourceDir: " ",
			wantErr:        "src-resource-dir パラメータは必須です",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewConfig(
				tt.operation,
				tt.pageType,
				tt.baseURL,
				tt.apiToken,
				tt.category,
				tt.skipsNoSrcBody,
				tt.srcJSONFile,
				tt.srcBodyDir,
				tt.srcResourceDir,
				tt.outDir,
				tt.targetStr,
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
			if tt.baseURL != "" && got.BaseURL != tt.baseURL {
				t.Fatalf("BaseURL = %q, want %q", got.BaseURL, tt.baseURL)
			}
			if tt.apiToken != "" && got.APIToken != tt.apiToken {
				t.Fatalf("APIToken = %q, want %q", got.APIToken, tt.apiToken)
			}
			if got.SkipsNoSrcBody != tt.skipsNoSrcBody {
				t.Fatalf("SkipsNoSrcBody = %v, want %v", got.SkipsNoSrcBody, tt.skipsNoSrcBody)
			}
			if tt.wantTargetStr != "" && got.TargetStr != tt.wantTargetStr {
				t.Fatalf("TargetStr = %q, want %q", got.TargetStr, tt.wantTargetStr)
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

	t.Run("grep str normal", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetString("operation", OperationGrepStr)
		parser.SetString("src-body-dir", "/tmp/body")
		parser.SetString("target-str", "TODO")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Operation != OperationGrepStr {
			t.Fatalf("Operation = %q", cfg.Operation)
		}
		if cfg.TargetStr != "TODO" {
			t.Fatalf("TargetStr = %q, want TODO", cfg.TargetStr)
		}
	})

	t.Run("grep str missing target str", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetString("operation", OperationGrepStr)
		parser.SetString("src-body-dir", "/tmp/body")

		_, err := ParseFlagsWithParser(parser)
		if err == nil || err.Error() != "target-str パラメータは必須です" {
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

	t.Run("rename bodies by category id normal", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetString("operation", OperationRenameBodiesByCategoryID)
		parser.SetString("page-type", PageTypeContent)
		parser.SetString("src-json-file", "/tmp/contents.json")
		parser.SetString("src-resource-dir", "/tmp/resource")
		parser.SetInt("con_number_start", 100)
		parser.SetInt("con_number_end", 200)

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Operation != OperationRenameBodiesByCategoryID {
			t.Fatalf("Operation = %q", cfg.Operation)
		}
		if cfg.SrcResourceDir != "/tmp/resource" {
			t.Fatalf("SrcResourceDir = %q, want /tmp/resource", cfg.SrcResourceDir)
		}
		if cfg.ConNumberStart != 100 || cfg.ConNumberEnd != 200 {
			t.Fatalf("con range = (%d,%d), want (100,200)", cfg.ConNumberStart, cfg.ConNumberEnd)
		}
	})

	t.Run("migrate to memos normal", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetString("operation", OperationMigrateToMemos)
		parser.SetString("page-type", PageTypeContent)
		parser.SetString("base-url", "https://memos.example.com")
		parser.SetString("api-token", "token")
		parser.SetString("src-body-dir", "/tmp/body")
		parser.SetString("src-resource-dir", "/tmp/resource")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Operation != OperationMigrateToMemos {
			t.Fatalf("Operation = %q", cfg.Operation)
		}
		if cfg.BaseURL != "https://memos.example.com" {
			t.Fatalf("BaseURL = %q, want https://memos.example.com", cfg.BaseURL)
		}
		if cfg.APIToken != "token" {
			t.Fatalf("APIToken = %q, want token", cfg.APIToken)
		}
		if cfg.SrcBodyDir != "/tmp/body" {
			t.Fatalf("SrcBodyDir = %q, want /tmp/body", cfg.SrcBodyDir)
		}
		if cfg.SrcResourceDir != "/tmp/resource" {
			t.Fatalf("SrcResourceDir = %q, want /tmp/resource", cfg.SrcResourceDir)
		}
	})

	t.Run("migrate to memos missing api token", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetString("operation", OperationMigrateToMemos)
		parser.SetString("page-type", PageTypeContent)
		parser.SetString("base-url", "https://memos.example.com")
		parser.SetString("src-body-dir", "/tmp/body")
		parser.SetString("src-resource-dir", "/tmp/resource")

		_, err := ParseFlagsWithParser(parser)
		if err == nil || err.Error() != "api-token パラメータは必須です" {
			t.Fatalf("error = %v", err)
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
