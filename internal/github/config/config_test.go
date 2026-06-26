package config

import (
	"strings"
	"testing"
)

// MockFlagParser はテスト用のFlagParserモック
type MockFlagParser struct {
	stringVars map[string]*string
	intVars    map[string]*int
	boolVars   map[string]*bool
	parseError error
	// 事前設定値
	stringValues map[string]string
	intValues    map[string]int
	boolValues   map[string]bool
}

// NewMockFlagParser は新しいMockFlagParserを作成する
func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringVars:   make(map[string]*string),
		intVars:      make(map[string]*int),
		boolVars:     make(map[string]*bool),
		stringValues: make(map[string]string),
		intValues:    make(map[string]int),
		boolValues:   make(map[string]bool),
	}
}

// SetStringValue は事前設定値を設定する
func (m *MockFlagParser) SetStringValue(name, value string) {
	m.stringValues[name] = value
}

// SetIntValue は事前設定値を設定する
func (m *MockFlagParser) SetIntValue(name string, value int) {
	m.intValues[name] = value
}

// SetBoolValue は事前設定値を設定する
func (m *MockFlagParser) SetBoolValue(name string, value bool) {
	m.boolValues[name] = value
}

// SetParseError はParseメソッドで返すエラーを設定する
func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}

// StringVar は文字列フラグを定義する
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	if presetValue, exists := m.stringValues[name]; exists {
		*p = presetValue // 事前設定値を優先
	} else {
		*p = value // デフォルト値
	}
	m.stringVars[name] = p
}

// IntVar は整数フラグを定義する
func (m *MockFlagParser) IntVar(p *int, name string, value int, usage string) {
	if presetValue, exists := m.intValues[name]; exists {
		*p = presetValue // 事前設定値を優先
	} else {
		*p = value // デフォルト値
	}
	m.intVars[name] = p
}

// BoolVar はブールフラグを定義する
func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if presetValue, exists := m.boolValues[name]; exists {
		*p = presetValue // 事前設定値を優先
	} else {
		*p = value // デフォルト値
	}
	m.boolVars[name] = p
}

// Parse はフラグを解析する
func (m *MockFlagParser) Parse() error {
	return m.parseError
}

// TestNewConfig_Normal はNewConfigの正常系テスト
func TestNewConfig_Normal(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		token     string
		owner     string
		repo      string
		state     string
		sort      string
		direction string
		perPage   int
		page      int
		expected  *Config
	}{
		{
			name:      "基本的なパラメータ",
			operation: "list-issues",
			token:     "test-token",
			owner:     "test-owner",
			repo:      "test-repo",
			state:     "",
			sort:      "",
			direction: "",
			perPage:   30,
			page:      1,
			expected: &Config{
				Operation: "list-issues",
				Token:     "test-token",
				Owner:     "test-owner",
				Repo:      "test-repo",
				State:     "",
				Sort:      "",
				Direction: "",
				PerPage:   30,
				Page:      1,
			},
		},
		{
			name:      "全パラメータ指定",
			operation: "list-issues",
			token:     "test-token",
			owner:     "test-owner",
			repo:      "test-repo",
			state:     "open",
			sort:      "created",
			direction: "desc",
			perPage:   50,
			page:      2,
			expected: &Config{
				Operation: "list-issues",
				Token:     "test-token",
				Owner:     "test-owner",
				Repo:      "test-repo",
				State:     "open",
				Sort:      "created",
				Direction: "desc",
				PerPage:   50,
				Page:      2,
			},
		},
		{
			name:      "デフォルト値の適用",
			operation: "list-issues",
			token:     "test-token",
			owner:     "test-owner",
			repo:      "test-repo",
			state:     "",
			sort:      "",
			direction: "",
			perPage:   0, // 0以下の場合は30に設定される
			page:      0, // 0以下の場合は1に設定される
			expected: &Config{
				Operation: "list-issues",
				Token:     "test-token",
				Owner:     "test-owner",
				Repo:      "test-repo",
				State:     "",
				Sort:      "",
				Direction: "",
				PerPage:   30,
				Page:      1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewConfig(tt.operation, tt.token, tt.owner, tt.repo, tt.state, tt.sort, tt.direction, tt.perPage, tt.page, 0)
			if err != nil {
				t.Errorf("NewConfig() error = %v, wantErr false", err)
				return
			}

			if result.Operation != tt.expected.Operation {
				t.Errorf("NewConfig().Operation = %v, want %v", result.Operation, tt.expected.Operation)
			}
			if result.Token != tt.expected.Token {
				t.Errorf("NewConfig().Token = %v, want %v", result.Token, tt.expected.Token)
			}
			if result.Owner != tt.expected.Owner {
				t.Errorf("NewConfig().Owner = %v, want %v", result.Owner, tt.expected.Owner)
			}
			if result.Repo != tt.expected.Repo {
				t.Errorf("NewConfig().Repo = %v, want %v", result.Repo, tt.expected.Repo)
			}
			if result.State != tt.expected.State {
				t.Errorf("NewConfig().State = %v, want %v", result.State, tt.expected.State)
			}
			if result.Sort != tt.expected.Sort {
				t.Errorf("NewConfig().Sort = %v, want %v", result.Sort, tt.expected.Sort)
			}
			if result.Direction != tt.expected.Direction {
				t.Errorf("NewConfig().Direction = %v, want %v", result.Direction, tt.expected.Direction)
			}
			if result.PerPage != tt.expected.PerPage {
				t.Errorf("NewConfig().PerPage = %v, want %v", result.PerPage, tt.expected.PerPage)
			}
			if result.Page != tt.expected.Page {
				t.Errorf("NewConfig().Page = %v, want %v", result.Page, tt.expected.Page)
			}
		})
	}
}

// TestNewConfig_Error はNewConfigのエラー系テスト
func TestNewConfig_Error(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		token     string
		owner     string
		repo      string
		wantErr   string
	}{
		{
			name:      "操作タイプが空",
			operation: "",
			token:     "test-token",
			owner:     "test-owner",
			repo:      "test-repo",
			wantErr:   "操作タイプが指定されていません",
		},
		{
			name:      "無効な操作タイプ",
			operation: "invalid-operation",
			token:     "test-token",
			owner:     "test-owner",
			repo:      "test-repo",
			wantErr:   "無効な操作タイプです: invalid-operation",
		},
		{
			name:      "トークンが空",
			operation: "list-issues",
			token:     "",
			owner:     "test-owner",
			repo:      "test-repo",
			wantErr:   "GitHubトークンが指定されていません",
		},
		{
			name:      "オーナーが空",
			operation: "list-issues",
			token:     "test-token",
			owner:     "",
			repo:      "test-repo",
			wantErr:   "リポジトリオーナーが指定されていません",
		},
		{
			name:      "リポジトリ名が空",
			operation: "list-issues",
			token:     "test-token",
			owner:     "test-owner",
			repo:      "",
			wantErr:   "リポジトリ名が指定されていません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConfig(tt.operation, tt.token, tt.owner, tt.repo, "", "", "", 30, 1, 0)
			if err == nil {
				t.Errorf("NewConfig() error = nil, wantErr %v", tt.wantErr)
				return
			}
			if err.Error() != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err.Error(), tt.wantErr)
			}
		})
	}
}

// TestParseFlagsWithParser_Normal はParseFlagsWithParserの正常系テスト
func TestParseFlagsWithParser_Normal(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		setup    func(*MockFlagParser)
		expected *Config
	}{
		{
			name:  "基本的なパラメータ",
			token: "test-token",
			setup: func(m *MockFlagParser) {
				m.SetStringValue("operation", "list-issues")
				m.SetStringValue("owner", "test-owner")
				m.SetStringValue("repo", "test-repo")
			},
			expected: &Config{
				Operation: "list-issues",
				Token:     "test-token",
				Owner:     "test-owner",
				Repo:      "test-repo",
				State:     "",
				Sort:      "",
				Direction: "",
				PerPage:   30,
				Page:      1,
			},
		},
		{
			name:  "短縮形パラメータ",
			token: "test-token",
			setup: func(m *MockFlagParser) {
				m.SetStringValue("o", "list-issues")
				m.SetStringValue("ow", "test-owner")
				m.SetStringValue("r", "test-repo")
			},
			expected: &Config{
				Operation: "list-issues",
				Token:     "test-token",
				Owner:     "test-owner",
				Repo:      "test-repo",
				State:     "",
				Sort:      "",
				Direction: "",
				PerPage:   30,
				Page:      1,
			},
		},
		{
			name: "ヘルプフラグ",
			setup: func(m *MockFlagParser) {
				m.SetBoolValue("help", true)
			},
			expected: &Config{
				Help: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(envKeyGitHubToken, tt.token)

			mockParser := NewMockFlagParser()
			tt.setup(mockParser)

			result, err := ParseFlagsWithParser(mockParser)
			if err != nil {
				t.Errorf("ParseFlagsWithParser() error = %v, wantErr false", err)
				return
			}

			if result.Help != tt.expected.Help {
				t.Errorf("ParseFlagsWithParser().Help = %v, want %v", result.Help, tt.expected.Help)
			}

			if !tt.expected.Help {
				if result.Operation != tt.expected.Operation {
					t.Errorf("ParseFlagsWithParser().Operation = %v, want %v", result.Operation, tt.expected.Operation)
				}
				if result.Token != tt.expected.Token {
					t.Errorf("ParseFlagsWithParser().Token = %v, want %v", result.Token, tt.expected.Token)
				}
				if result.Owner != tt.expected.Owner {
					t.Errorf("ParseFlagsWithParser().Owner = %v, want %v", result.Owner, tt.expected.Owner)
				}
				if result.Repo != tt.expected.Repo {
					t.Errorf("ParseFlagsWithParser().Repo = %v, want %v", result.Repo, tt.expected.Repo)
				}
			}
		})
	}
}

// TestParseFlagsWithParser_Error はParseFlagsWithParserのエラー系テスト
func TestParseFlagsWithParser_Error(t *testing.T) {
	t.Run("GitHubトークン環境変数が空", func(t *testing.T) {
		t.Setenv(envKeyGitHubToken, "")

		mockParser := NewMockFlagParser()
		mockParser.SetStringValue("operation", "list-issues")
		mockParser.SetStringValue("owner", "test-owner")
		mockParser.SetStringValue("repo", "test-repo")

		_, err := ParseFlagsWithParser(mockParser)
		if err == nil {
			t.Fatal("ParseFlagsWithParser() error = nil, want error")
		}
		if !strings.Contains(err.Error(), envKeyGitHubToken) {
			t.Fatalf("ParseFlagsWithParser() error = %v, want %s", err, envKeyGitHubToken)
		}
	})

	t.Run("ヘルプフラグはGitHubトークン環境変数が空でも成功", func(t *testing.T) {
		t.Setenv(envKeyGitHubToken, "")

		mockParser := NewMockFlagParser()
		mockParser.SetBoolValue("help", true)

		result, err := ParseFlagsWithParser(mockParser)
		if err != nil {
			t.Fatalf("ParseFlagsWithParser() error = %v, want nil", err)
		}
		if !result.Help {
			t.Fatal("ParseFlagsWithParser().Help = false, want true")
		}
	})
}
