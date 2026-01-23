package config

import (
	"errors"
	"strings"
	"testing"
)

// MockFlagParser はテスト用のモックFlagParser
type MockFlagParser struct {
	stringVars   map[string]*string
	boolVars     map[string]*bool
	stringValues map[string]string // 事前設定された文字列値
	boolValues   map[string]bool   // 事前設定されたブール値
	args         []string
	parseError   error
}

// NewMockFlagParser は新しいMockFlagParserを作成する
func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringVars:   make(map[string]*string),
		boolVars:     make(map[string]*bool),
		stringValues: make(map[string]string),
		boolValues:   make(map[string]bool),
		args:         []string{},
	}
}

// StringVar は文字列フラグを定義する
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	// 事前設定された値があるかチェック
	if presetValue, exists := m.stringValues[name]; exists {
		*p = presetValue
	} else {
		*p = value // デフォルト値を設定
	}
	m.stringVars[name] = p
}

// BoolVar はブールフラグを定義する
func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	// 事前設定された値があるかチェック
	if presetValue, exists := m.boolValues[name]; exists {
		*p = presetValue
	} else {
		*p = value // デフォルト値を設定
	}
	m.boolVars[name] = p
}

// Parse はフラグを解析する
func (m *MockFlagParser) Parse() error {
	return m.parseError
}

// Args は解析後の残りの引数を返す
func (m *MockFlagParser) Args() []string {
	return m.args
}

// SetStringFlag はテスト用に文字列フラグの値を設定する
func (m *MockFlagParser) SetStringFlag(name, value string) {
	m.stringValues[name] = value
	if p, exists := m.stringVars[name]; exists {
		*p = value
	}
}

// SetBoolFlag はテスト用にブールフラグの値を設定する
func (m *MockFlagParser) SetBoolFlag(name string, value bool) {
	m.boolValues[name] = value
	if p, exists := m.boolVars[name]; exists {
		*p = value
	}
}

// SetParseError はモック用にパースエラーを設定する
func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}

// SetArgs はモック用に引数を設定する
func (m *MockFlagParser) SetArgs(args []string) {
	m.args = args
}

// TestNewConfig_Normal は正常系のテストケース
func TestNewConfig_Normal(t *testing.T) {
	tests := []struct {
		name               string
		embedType          string
		webhookURL         string
		botName            string
		contentText        string
		embedText          string
		embedColor         string
		embedURLLinkedText string
		expectError        bool
	}{
		{
			name:        "基本的な設定（embed-type: none）",
			embedType:   "none",
			botName:     "testボット",
			webhookURL:  "https://discord.com/api/webhooks/123456789/abcdefg",
			contentText: "テストメッセージ",
			expectError: false,
		},
		{
			name:        "基本的な設定（embed-type: vscode）",
			embedType:   "vscode",
			botName:     "",
			webhookURL:  "https://discord.com/api/webhooks/123456789/abcdefg",
			contentText: "テストメッセージ",
			expectError: false,
		},
		{
			name:               "全オプション指定",
			embedType:          "vscode",
			webhookURL:         "https://discord.com/api/webhooks/123456789/abcdefg",
			botName:            "",
			contentText:        "テストメッセージ",
			embedText:          "埋め込みタイトル",
			embedColor:         "green",
			embedURLLinkedText: "https://example.com",
			expectError:        false,
		},
		{
			name:        "discordapp.comドメインのwebhook URL",
			embedType:   "none",
			webhookURL:  "https://discordapp.com/api/webhooks/123456789/abcdefg",
			botName:     "",
			contentText: "テストメッセージ",
			expectError: false,
		},
		{
			name:        "OpenWeatherMap embed",
			embedType:   "open-weather-map",
			webhookURL:  "https://discord.com/api/webhooks/123456789/abcdefg",
			botName:     "",
			contentText: "テストメッセージ",
			embedText:   "本日の天気予報",
			embedColor:  "orange",
			expectError: false,
		},
		{
			name:        "Google Compute Engine success embed",
			embedType:   "google-compute-engine-success",
			webhookURL:  "https://discord.com/api/webhooks/123456789/abcdefg",
			botName:     "",
			contentText: "テストメッセージ",
			expectError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			config, err := NewConfig(tt.embedType, tt.webhookURL, tt.botName, tt.contentText, tt.embedText, tt.embedColor, tt.embedURLLinkedText)

			if tt.expectError {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
				return
			}

			if config.EmbedType != tt.embedType {
				t.Errorf("EmbedType = %v, want %v", config.EmbedType, tt.embedType)
			}
			if config.WebhookURL != tt.webhookURL {
				t.Errorf("WebhookURL = %v, want %v", config.WebhookURL, tt.webhookURL)
			}
			if config.ContentText != tt.contentText {
				t.Errorf("ContentText = %v, want %v", config.ContentText, tt.contentText)
			}
			if config.EmbedText != tt.embedText {
				t.Errorf("EmbedText = %v, want %v", config.EmbedText, tt.embedText)
			}
			if config.EmbedColor != tt.embedColor {
				t.Errorf("EmbedColor = %v, want %v", config.EmbedColor, tt.embedColor)
			}
			if config.EmbedURLLinkedText != tt.embedURLLinkedText {
				t.Errorf("EmbedURLLinkedText = %v, want %v", config.EmbedURLLinkedText, tt.embedURLLinkedText)
			}
		})
	}
}

// TestNewConfig_Error は異常系のテストケース
func TestNewConfig_Error(t *testing.T) {
	tests := []struct {
		name               string
		embedType          string
		webhookURL         string
		botName            string
		contentText        string
		embedText          string
		embedColor         string
		embedURLLinkedText string
		expectedError      string
	}{
		{
			name:          "webhook-url未指定",
			embedType:     "none",
			webhookURL:    "",
			botName:       "",
			contentText:   "テストメッセージ",
			expectedError: "webhook-urlが指定されていません",
		},
		{
			name:          "content-text未指定",
			embedType:     "none",
			webhookURL:    "https://discord.com/api/webhooks/123456789/abcdefg",
			botName:       "",
			contentText:   "",
			expectedError: "content-textが指定されていません",
		},
		{
			name:          "embed-type未指定",
			embedType:     "",
			webhookURL:    "https://discord.com/api/webhooks/123456789/abcdefg",
			botName:       "",
			contentText:   "テストメッセージ",
			expectedError: "embed-typeが指定されていません",
		},
		{
			name:          "無効なembed-type",
			embedType:     "invalid",
			webhookURL:    "https://discord.com/api/webhooks/123456789/abcdefg",
			botName:       "",
			contentText:   "テストメッセージ",
			expectedError: "無効なembed-typeです: invalid (有効な値: google-cloud-iam-failed, google-cloud-iam-success, google-cloud-run-failed, google-cloud-run-function-failed, google-cloud-run-function-success, google-cloud-run-success, google-cloud-scheduler-failed, google-cloud-scheduler-success, google-cloud-storage-failed, google-cloud-storage-success, google-compute-engine-failed, google-compute-engine-success, google-secret-manager-failed, google-secret-manager-success, none, open-weather-map, postgres, vscode)",
		},
		{
			name:          "無効なwebhook-URL（httpスキーム）",
			embedType:     "none",
			webhookURL:    "http://discord.com/api/webhooks/123456789/abcdefg",
			botName:       "",
			contentText:   "テストメッセージ",
			expectedError: "無効なwebhook-urlです。Discord WebhookのURLを指定してください",
		},
		{
			name:          "無効なwebhook-URL（異なるドメイン）",
			embedType:     "none",
			webhookURL:    "https://example.com/api/webhooks/123456789/abcdefg",
			botName:       "",
			contentText:   "テストメッセージ",
			expectedError: "無効なwebhook-urlです。Discord WebhookのURLを指定してください",
		},
		{
			name:          "content-textの文字数制限超過",
			embedType:     "none",
			webhookURL:    "https://discord.com/api/webhooks/123456789/abcdefg",
			botName:       "testボット",
			contentText:   strings.Repeat("a", 2001),
			expectedError: "content-textは2000文字以下である必要があります",
		},
		{
			name:          "embed-textの文字数制限超過",
			embedType:     "vscode",
			webhookURL:    "https://discord.com/api/webhooks/123456789/abcdefg",
			botName:       "testボット",
			contentText:   "テストメッセージ",
			embedText:     strings.Repeat("a", 257),
			expectedError: "embed-textは256文字以下である必要があります",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConfig(tt.embedType, tt.webhookURL, tt.botName, tt.contentText, tt.embedText, tt.embedColor, tt.embedURLLinkedText)

			if err == nil {
				t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
				return
			}

			if err.Error() != tt.expectedError {
				t.Errorf("エラーメッセージ = %v, want %v", err.Error(), tt.expectedError)
			}
		})
	}
}

// TestParseFlagsWithParser_Normal は正常系のテストケース
func TestParseFlagsWithParser_Normal(t *testing.T) {
	t.Run("ヘルプフラグ", func(t *testing.T) {
		mockParser := NewMockFlagParser()

		// ヘルプフラグを事前設定
		mockParser.SetBoolFlag("help", true)
		config, err := ParseFlagsWithParser(mockParser)

		if err != nil {
			t.Errorf("予期しないエラーが発生しました: %v", err)
			return
		}

		if !config.Help {
			t.Errorf("Help = %v, want %v", config.Help, true)
		}
	})

	t.Run("基本的なフラグ解析", func(t *testing.T) {
		mockParser := NewMockFlagParser()

		// 必要なフラグを事前設定
		mockParser.SetStringFlag("embed-type", "none")
		mockParser.SetStringFlag("webhook-url", "https://discord.com/api/webhooks/123456789/abcdefg")
		mockParser.SetStringFlag("content-text", "テストメッセージ")

		config, err := ParseFlagsWithParser(mockParser)

		if err != nil {
			t.Errorf("予期しないエラーが発生しました: %v", err)
			return
		}

		if config.EmbedType != "none" {
			t.Errorf("EmbedType = %v, want %v", config.EmbedType, "none")
		}
		if config.WebhookURL != "https://discord.com/api/webhooks/123456789/abcdefg" {
			t.Errorf("WebhookURL = %v, want %v", config.WebhookURL, "https://discord.com/api/webhooks/123456789/abcdefg")
		}
		if config.ContentText != "テストメッセージ" {
			t.Errorf("ContentText = %v, want %v", config.ContentText, "テストメッセージ")
		}
	})

	t.Run("全オプション指定", func(t *testing.T) {
		mockParser := NewMockFlagParser()

		// 全てのフラグを事前設定
		mockParser.SetStringFlag("embed-type", "vscode")
		mockParser.SetStringFlag("webhook-url", "https://discord.com/api/webhooks/123456789/abcdefg")
		mockParser.SetStringFlag("content-text", "テストメッセージ")
		mockParser.SetStringFlag("embed-text", "埋め込みタイトル")
		mockParser.SetStringFlag("embed-color", "green")
		mockParser.SetStringFlag("embed-url-linked-text", "https://example.com")

		config, err := ParseFlagsWithParser(mockParser)

		if err != nil {
			t.Errorf("予期しないエラーが発生しました: %v", err)
			return
		}

		if config.EmbedType != "vscode" {
			t.Errorf("EmbedType = %v, want %v", config.EmbedType, "vscode")
		}
		if config.WebhookURL != "https://discord.com/api/webhooks/123456789/abcdefg" {
			t.Errorf("WebhookURL = %v, want %v", config.WebhookURL, "https://discord.com/api/webhooks/123456789/abcdefg")
		}
		if config.ContentText != "テストメッセージ" {
			t.Errorf("ContentText = %v, want %v", config.ContentText, "テストメッセージ")
		}
		if config.EmbedText != "埋め込みタイトル" {
			t.Errorf("EmbedText = %v, want %v", config.EmbedText, "埋め込みタイトル")
		}
		if config.EmbedColor != "green" {
			t.Errorf("EmbedColor = %v, want %v", config.EmbedColor, "green")
		}
		if config.EmbedURLLinkedText != "https://example.com" {
			t.Errorf("EmbedURLLinkedText = %v, want %v", config.EmbedURLLinkedText, "https://example.com")
		}
	})

	t.Run("Google Cloud成功Embed", func(t *testing.T) {
		mockParser := NewMockFlagParser()

		mockParser.SetStringFlag("embed-type", "google-cloud-run-success")
		mockParser.SetStringFlag("webhook-url", "https://discord.com/api/webhooks/123456789/abcdefg")
		mockParser.SetStringFlag("content-text", "テストメッセージ")

		config, err := ParseFlagsWithParser(mockParser)

		if err != nil {
			t.Errorf("予期しないエラーが発生しました: %v", err)
			return
		}

		if config.EmbedType != "google-cloud-run-success" {
			t.Errorf("EmbedType = %v, want %v", config.EmbedType, "google-cloud-run-success")
		}
	})

	t.Run("短縮形フラグ", func(t *testing.T) {
		mockParser := NewMockFlagParser()

		// 短縮形フラグを事前設定
		mockParser.SetStringFlag("et", "vscode")
		mockParser.SetStringFlag("wu", "https://discord.com/api/webhooks/123456789/abcdefg")
		mockParser.SetStringFlag("ct", "テストメッセージ")
		mockParser.SetStringFlag("et-text", "埋め込みタイトル")
		mockParser.SetStringFlag("ec", "green")
		mockParser.SetStringFlag("eult", "https://example.com")

		config, err := ParseFlagsWithParser(mockParser)

		if err != nil {
			t.Errorf("予期しないエラーが発生しました: %v", err)
			return
		}

		if config.EmbedType != "vscode" {
			t.Errorf("EmbedType = %v, want %v", config.EmbedType, "vscode")
		}
		if config.WebhookURL != "https://discord.com/api/webhooks/123456789/abcdefg" {
			t.Errorf("WebhookURL = %v, want %v", config.WebhookURL, "https://discord.com/api/webhooks/123456789/abcdefg")
		}
		if config.ContentText != "テストメッセージ" {
			t.Errorf("ContentText = %v, want %v", config.ContentText, "テストメッセージ")
		}
		if config.EmbedText != "埋め込みタイトル" {
			t.Errorf("EmbedText = %v, want %v", config.EmbedText, "埋め込みタイトル")
		}
		if config.EmbedColor != "green" {
			t.Errorf("EmbedColor = %v, want %v", config.EmbedColor, "green")
		}
		if config.EmbedURLLinkedText != "https://example.com" {
			t.Errorf("EmbedURLLinkedText = %v, want %v", config.EmbedURLLinkedText, "https://example.com")
		}
	})
}

// TestParseFlagsWithParser_Error は異常系のテストケース
func TestParseFlagsWithParser_Error(t *testing.T) {
	tests := []struct {
		name          string
		parseError    error
		embedType     string
		webhookURL    string
		contentText   string
		expectedError string
	}{
		{
			name:          "フラグ解析エラー",
			parseError:    errors.New("parse error"),
			expectedError: "フラグの解析に失敗しました:",
		},
		{
			name:          "必須パラメータ不足（webhook-url）",
			embedType:     "none",
			contentText:   "テストメッセージ",
			expectedError: "webhook-urlが指定されていません",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockParser := NewMockFlagParser()

			if tt.parseError != nil {
				mockParser.SetParseError(tt.parseError)
			}

			if tt.embedType != "" {
				mockParser.SetStringFlag("embed-type", tt.embedType)
			}
			if tt.webhookURL != "" {
				mockParser.SetStringFlag("webhook-url", tt.webhookURL)
			}
			if tt.contentText != "" {
				mockParser.SetStringFlag("content-text", tt.contentText)
			}

			_, err := ParseFlagsWithParser(mockParser)

			if err == nil {
				t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
				return
			}

			if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("エラーメッセージに '%v' が含まれていません。実際: %v", tt.expectedError, err.Error())
			}
		})
	}
}

// TestPrintUsage はPrintUsage関数の存在を確認するテスト
func TestPrintUsage(t *testing.T) {
	// PrintUsage関数が存在し、パニックしないことを確認
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintUsage() caused panic: %v", r)
		}
	}()

	PrintUsage()
}

// TestMockFlagParser_Methods はモックFlagParserのメソッドをテストする
func TestMockFlagParser_Methods(t *testing.T) {
	mockParser := NewMockFlagParser()

	// StringVarのテスト
	var testString string
	mockParser.StringVar(&testString, "test", "default", "usage")
	if testString != "default" {
		t.Errorf("StringVar failed: got %v, want %v", testString, "default")
	}

	// BoolVarのテスト
	var testBool bool
	mockParser.BoolVar(&testBool, "bool", true, "usage")
	if !testBool {
		t.Errorf("BoolVar failed: got %v, want %v", testBool, true)
	}

	// Parseのテスト
	err := mockParser.Parse()
	if err != nil {
		t.Errorf("Parse failed: %v", err)
	}

	// Argsのテスト
	testArgs := []string{"arg1", "arg2"}
	mockParser.SetArgs(testArgs)
	args := mockParser.Args()
	if len(args) != len(testArgs) {
		t.Errorf("Args length mismatch: got %d, want %d", len(args), len(testArgs))
	}
	for i, arg := range args {
		if arg != testArgs[i] {
			t.Errorf("Args[%d] = %v, want %v", i, arg, testArgs[i])
		}
	}
}

// TestMockFlagParser_PresetValues は事前設定値のテスト
func TestMockFlagParser_PresetValues(t *testing.T) {
	mockParser := NewMockFlagParser()

	// 事前に値を設定
	mockParser.SetStringFlag("preset-string", "preset-value")
	mockParser.SetBoolFlag("preset-bool", true)

	// フラグを定義（事前設定値が使用されるはず）
	var stringFlag string
	var boolFlag bool
	mockParser.StringVar(&stringFlag, "preset-string", "default", "usage")
	mockParser.BoolVar(&boolFlag, "preset-bool", false, "usage")

	if stringFlag != "preset-value" {
		t.Errorf("Expected preset string value 'preset-value', got '%s'", stringFlag)
	}
	if !boolFlag {
		t.Errorf("Expected preset bool value true, got %t", boolFlag)
	}
}

// TestMockFlagParser_DefaultValues はデフォルト値のテスト
func TestMockFlagParser_DefaultValues(t *testing.T) {
	mockParser := NewMockFlagParser()

	var stringFlag string
	var boolFlag bool

	// 事前設定なしでフラグを定義（デフォルト値が使用されるはず）
	mockParser.StringVar(&stringFlag, "string-flag", "default-string", "string flag")
	mockParser.BoolVar(&boolFlag, "bool-flag", true, "bool flag")

	if stringFlag != "default-string" {
		t.Errorf("Expected default string 'default-string', got '%s'", stringFlag)
	}
	if !boolFlag {
		t.Error("Expected default bool true, got false")
	}
}

// TestMockFlagParser_EmptyArgs は空の引数リストのテスト
func TestMockFlagParser_EmptyArgs(t *testing.T) {
	mockParser := NewMockFlagParser()

	args := mockParser.Args()
	if args == nil {
		t.Error("Expected non-nil args, got nil")
	}
	if len(args) != 0 {
		t.Errorf("Expected empty args, got %d items", len(args))
	}
}

// TestMockFlagParser_ParseError はパースエラーのテスト
func TestMockFlagParser_ParseError(t *testing.T) {
	mockParser := NewMockFlagParser()
	testError := errors.New("test parse error")
	mockParser.SetParseError(testError)

	err := mockParser.Parse()
	if err != testError {
		t.Errorf("Expected parse error %v, got %v", testError, err)
	}
}

// TestParseFlags は実際のParseFlags関数のテスト
func TestParseFlags(t *testing.T) {
	// この関数は実際のコマンドライン引数を使用するため、
	// 単体テストでは直接テストしにくいが、関数が存在することを確認
	_, err := ParseFlags()
	// エラーが発生することは期待される（必須パラメータが不足しているため）
	if err == nil {
		t.Log("ParseFlags()が正常に実行されました（テスト環境では通常エラーになります）")
	}
}
