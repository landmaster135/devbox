package config

import (
	"errors"
	"os"
	"strings"
	"testing"
)

// MockFlagParser はテスト用のモックFlagParser
type MockFlagParser struct {
	stringVars   map[string]*string
	intVars      map[string]*int
	boolVars     map[string]*bool
	stringValues map[string]string // 事前設定された文字列値
	intValues    map[string]int    // 事前設定された整数値
	boolValues   map[string]bool   // 事前設定されたブール値
	args         []string
	parseError   error
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

// IntVar は整数フラグを定義する
func (m *MockFlagParser) IntVar(p *int, name string, value int, usage string) {
	// 事前設定された値があるかチェック
	if presetValue, exists := m.intValues[name]; exists {
		*p = presetValue
	} else {
		*p = value // デフォルト値を設定
	}
	m.intVars[name] = p
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

// SetIntFlag はテスト用に整数フラグの値を設定する
func (m *MockFlagParser) SetIntFlag(name string, value int) {
	m.intValues[name] = value
	if p, exists := m.intVars[name]; exists {
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
		name       string
		apiKey     string
		city       string
		maxDays    int
		webhookURL string
	}{
		{
			name:       "基本的な設定",
			apiKey:     "test-api-key",
			city:       "Tokyo",
			maxDays:    3,
			webhookURL: "https://discord.com/api/webhooks/123456789/abcdefg",
		},
		{
			name:       "最小日数（1日）",
			apiKey:     "test-api-key",
			city:       "Osaka",
			maxDays:    1,
			webhookURL: "https://discord.com/api/webhooks/123456789/abcdefg",
		},
		{
			name:       "最大日数（5日）",
			apiKey:     "test-api-key",
			city:       "Kyoto",
			maxDays:    5,
			webhookURL: "https://discord.com/api/webhooks/123456789/abcdefg",
		},
		{
			name:       "スペースを含む都市名",
			apiKey:     "test-api-key",
			city:       "New York",
			maxDays:    3,
			webhookURL: "https://discord.com/api/webhooks/123456789/abcdefg",
		},
		{
			name:       "国コード付き都市名",
			apiKey:     "test-api-key",
			city:       "Tokyo,JP",
			maxDays:    2,
			webhookURL: "https://discord.com/api/webhooks/123456789/abcdefg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := NewConfig(tt.apiKey, tt.city, tt.maxDays, tt.webhookURL)

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
				return
			}

			if config.APIKey != tt.apiKey {
				t.Errorf("APIKey = %v, want %v", config.APIKey, tt.apiKey)
			}
			if config.City != tt.city {
				t.Errorf("City = %v, want %v", config.City, tt.city)
			}
			if config.MaxDays != tt.maxDays {
				t.Errorf("MaxDays = %v, want %v", config.MaxDays, tt.maxDays)
			}
			if config.WebhookURL != tt.webhookURL {
				t.Errorf("WebhookURL = %v, want %v", config.WebhookURL, tt.webhookURL)
			}
			if config.Help != false {
				t.Errorf("Help = %v, want %v", config.Help, false)
			}
		})
	}
}

// TestNewConfig_Error は異常系のテストケース
func TestNewConfig_Error(t *testing.T) {
	tests := []struct {
		name          string
		apiKey        string
		city          string
		maxDays       int
		webhookURL    string
		expectedError string
	}{
		{
			name:          "APIキー未指定",
			apiKey:        "",
			city:          "Tokyo",
			maxDays:       3,
			webhookURL:    "https://discord.com/api/webhooks/123456789/abcdefg",
			expectedError: "API キーが指定されていません",
		},
		{
			name:          "都市名未指定",
			apiKey:        "test-api-key",
			city:          "",
			maxDays:       3,
			webhookURL:    "https://discord.com/api/webhooks/123456789/abcdefg",
			expectedError: "都市名が指定されていません",
		},
		{
			name:          "最大日数が0",
			apiKey:        "test-api-key",
			city:          "Tokyo",
			maxDays:       0,
			webhookURL:    "https://discord.com/api/webhooks/123456789/abcdefg",
			expectedError: "最大日数は1以上である必要があります",
		},
		{
			name:          "最大日数が負の値",
			apiKey:        "test-api-key",
			city:          "Tokyo",
			maxDays:       -1,
			webhookURL:    "https://discord.com/api/webhooks/123456789/abcdefg",
			expectedError: "最大日数は1以上である必要があります",
		},
		{
			name:          "最大日数が6日（上限超過）",
			apiKey:        "test-api-key",
			city:          "Tokyo",
			maxDays:       6,
			webhookURL:    "https://discord.com/api/webhooks/123456789/abcdefg",
			expectedError: "最大日数は5日以下である必要があります（OpenWeather API制限）",
		},
		{
			name:          "最大日数が10日（大幅な上限超過）",
			apiKey:        "test-api-key",
			city:          "Tokyo",
			maxDays:       10,
			webhookURL:    "https://discord.com/api/webhooks/123456789/abcdefg",
			expectedError: "最大日数は5日以下である必要があります（OpenWeather API制限）",
		},
		{
			name:          "WebhookURL未指定",
			apiKey:        "test-api-key",
			city:          "Tokyo",
			maxDays:       3,
			webhookURL:    "",
			expectedError: "Discord Webhook URLが指定されていません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConfig(tt.apiKey, tt.city, tt.maxDays, tt.webhookURL)

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

// TestStandardFlagParser_Normal はStandardFlagParserの正常系テスト
func TestStandardFlagParser_Normal(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	t.Run("StringVar解析", func(t *testing.T) {
		os.Args = []string{"program", "-api-key", "test-key", "-city", "Tokyo"}
		parser := NewStandardFlagParser()

		var apiKey, city string
		parser.StringVar(&apiKey, "api-key", "", "API key")
		parser.StringVar(&city, "city", "", "City name")

		if apiKey != "test-key" {
			t.Errorf("apiKey = %v, want %v", apiKey, "test-key")
		}
		if city != "Tokyo" {
			t.Errorf("city = %v, want %v", city, "Tokyo")
		}
	})

	t.Run("IntVar解析", func(t *testing.T) {
		os.Args = []string{"program", "-max-days", "3"}
		parser := NewStandardFlagParser()

		var maxDays int
		parser.IntVar(&maxDays, "max-days", 1, "Max days")

		if maxDays != 3 {
			t.Errorf("maxDays = %v, want %v", maxDays, 3)
		}
	})

	t.Run("BoolVar解析", func(t *testing.T) {
		os.Args = []string{"program", "-help"}
		parser := NewStandardFlagParser()

		var help bool
		parser.BoolVar(&help, "help", false, "Show help")

		if !help {
			t.Errorf("help = %v, want %v", help, true)
		}
	})

	t.Run("Parse関数", func(t *testing.T) {
		os.Args = []string{"program", "arg1", "arg2"}
		parser := NewStandardFlagParser()

		err := parser.Parse()
		if err != nil {
			t.Errorf("Parse() error = %v, want nil", err)
		}

		args := parser.Args()
		expectedArgs := []string{"arg1", "arg2"}
		if len(args) != len(expectedArgs) {
			t.Errorf("Args length = %v, want %v", len(args), len(expectedArgs))
		}
		for i, arg := range args {
			if arg != expectedArgs[i] {
				t.Errorf("Args[%d] = %v, want %v", i, arg, expectedArgs[i])
			}
		}
	})

	t.Run("ロングフラグ形式", func(t *testing.T) {
		os.Args = []string{"program", "--api-key", "long-key", "--city", "Osaka"}
		parser := NewStandardFlagParser()

		var apiKey, city string
		parser.StringVar(&apiKey, "api-key", "", "API key")
		parser.StringVar(&city, "city", "", "City name")

		if apiKey != "long-key" {
			t.Errorf("apiKey = %v, want %v", apiKey, "long-key")
		}
		if city != "Osaka" {
			t.Errorf("city = %v, want %v", city, "Osaka")
		}
	})
}

// TestStandardFlagParser_EdgeCases はStandardFlagParserのエッジケーステスト
func TestStandardFlagParser_EdgeCases(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	t.Run("フラグ値なし", func(t *testing.T) {
		os.Args = []string{"program", "-api-key"}
		parser := NewStandardFlagParser()

		var apiKey string = "default"  // デフォルト値を明示的に設定
		parser.StringVar(&apiKey, "api-key", "default", "API key")

		// StandardFlagParserの実装では、値がない場合は変数を変更しない
		if apiKey != "default" {
			t.Errorf("apiKey = %v, want %v", apiKey, "default")
		}
	})

	t.Run("無効な整数値", func(t *testing.T) {
		os.Args = []string{"program", "-max-days", "invalid"}
		parser := NewStandardFlagParser()

		var maxDays int = 1  // デフォルト値を明示的に設定
		parser.IntVar(&maxDays, "max-days", 1, "Max days")

		// StandardFlagParserの実装では、無効な値の場合は変数を変更しない
		if maxDays != 1 {
			t.Errorf("maxDays = %v, want %v", maxDays, 1)
		}
	})

	t.Run("存在しないフラグ", func(t *testing.T) {
		os.Args = []string{"program", "-nonexistent", "value"}
		parser := NewStandardFlagParser()

		var value string = "default"  // デフォルト値を明示的に設定
		parser.StringVar(&value, "existing", "default", "Existing flag")

		// StandardFlagParserの実装では、存在しないフラグの場合は変数を変更しない
		if value != "default" {
			t.Errorf("value = %v, want %v", value, "default")
		}
	})
}

// TestParseFlagsWithParser_Normal は正常系のテストケース
func TestParseFlagsWithParser_Normal(t *testing.T) {
	t.Run("ヘルプフラグ", func(t *testing.T) {
		mockParser := NewMockFlagParser()
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

	t.Run("ヘルプフラグ短縮形", func(t *testing.T) {
		mockParser := NewMockFlagParser()
		mockParser.SetBoolFlag("h", true)

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
		mockParser.SetStringFlag("api-key", "test-api-key")
		mockParser.SetStringFlag("city", "Tokyo")
		mockParser.SetIntFlag("max-days", 3)
		mockParser.SetStringFlag("webhook-url", "https://discord.com/api/webhooks/123456789/abcdefg")

		config, err := ParseFlagsWithParser(mockParser)

		if err != nil {
			t.Errorf("予期しないエラーが発生しました: %v", err)
			return
		}

		if config.APIKey != "test-api-key" {
			t.Errorf("APIKey = %v, want %v", config.APIKey, "test-api-key")
		}
		if config.City != "Tokyo" {
			t.Errorf("City = %v, want %v", config.City, "Tokyo")
		}
		if config.MaxDays != 3 {
			t.Errorf("MaxDays = %v, want %v", config.MaxDays, 3)
		}
		if config.WebhookURL != "https://discord.com/api/webhooks/123456789/abcdefg" {
			t.Errorf("WebhookURL = %v, want %v", config.WebhookURL, "https://discord.com/api/webhooks/123456789/abcdefg")
		}
		if config.Help {
			t.Errorf("Help = %v, want %v", config.Help, false)
		}
	})

	t.Run("境界値テスト（最小日数）", func(t *testing.T) {
		mockParser := NewMockFlagParser()
		mockParser.SetStringFlag("api-key", "test-api-key")
		mockParser.SetStringFlag("city", "Tokyo")
		mockParser.SetIntFlag("max-days", 1)
		mockParser.SetStringFlag("webhook-url", "https://discord.com/api/webhooks/123456789/abcdefg")

		config, err := ParseFlagsWithParser(mockParser)

		if err != nil {
			t.Errorf("予期しないエラーが発生しました: %v", err)
			return
		}

		if config.MaxDays != 1 {
			t.Errorf("MaxDays = %v, want %v", config.MaxDays, 1)
		}
	})

	t.Run("境界値テスト（最大日数）", func(t *testing.T) {
		mockParser := NewMockFlagParser()
		mockParser.SetStringFlag("api-key", "test-api-key")
		mockParser.SetStringFlag("city", "Tokyo")
		mockParser.SetIntFlag("max-days", 5)
		mockParser.SetStringFlag("webhook-url", "https://discord.com/api/webhooks/123456789/abcdefg")

		config, err := ParseFlagsWithParser(mockParser)

		if err != nil {
			t.Errorf("予期しないエラーが発生しました: %v", err)
			return
		}

		if config.MaxDays != 5 {
			t.Errorf("MaxDays = %v, want %v", config.MaxDays, 5)
		}
	})
}

// TestParseFlagsWithParser_Error は異常系のテストケース
func TestParseFlagsWithParser_Error(t *testing.T) {
	tests := []struct {
		name          string
		parseError    error
		apiKey        string
		city          string
		maxDays       int
		webhookURL    string
		expectedError string
	}{
		{
			name:          "フラグ解析エラー",
			parseError:    errors.New("parse error"),
			expectedError: "フラグの解析に失敗しました:",
		},
		{
			name:          "APIキー不足",
			city:          "Tokyo",
			maxDays:       3,
			webhookURL:    "https://discord.com/api/webhooks/123456789/abcdefg",
			expectedError: "API キーが指定されていません",
		},
		{
			name:          "都市名不足",
			apiKey:        "test-api-key",
			maxDays:       3,
			webhookURL:    "https://discord.com/api/webhooks/123456789/abcdefg",
			expectedError: "都市名が指定されていません",
		},
		{
			name:          "WebhookURL不足",
			apiKey:        "test-api-key",
			city:          "Tokyo",
			maxDays:       3,
			expectedError: "Discord Webhook URLが指定されていません",
		},
		{
			name:          "無効な最大日数（0）",
			apiKey:        "test-api-key",
			city:          "Tokyo",
			maxDays:       0,
			webhookURL:    "https://discord.com/api/webhooks/123456789/abcdefg",
			expectedError: "最大日数は1以上である必要があります",
		},
		{
			name:          "無効な最大日数（6）",
			apiKey:        "test-api-key",
			city:          "Tokyo",
			maxDays:       6,
			webhookURL:    "https://discord.com/api/webhooks/123456789/abcdefg",
			expectedError: "最大日数は5日以下である必要があります（OpenWeather API制限）",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParser := NewMockFlagParser()

			if tt.parseError != nil {
				mockParser.SetParseError(tt.parseError)
			}

			if tt.apiKey != "" {
				mockParser.SetStringFlag("api-key", tt.apiKey)
			}
			if tt.city != "" {
				mockParser.SetStringFlag("city", tt.city)
			}
			// maxDaysは0の場合もテストするため、常に設定する
			mockParser.SetIntFlag("max-days", tt.maxDays)
			if tt.webhookURL != "" {
				mockParser.SetStringFlag("webhook-url", tt.webhookURL)
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

	// IntVarのテスト
	var testInt int
	mockParser.IntVar(&testInt, "int", 42, "usage")
	if testInt != 42 {
		t.Errorf("IntVar failed: got %v, want %v", testInt, 42)
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
	mockParser.SetIntFlag("preset-int", 100)
	mockParser.SetBoolFlag("preset-bool", true)

	// フラグを定義（事前設定値が使用されるはず）
	var stringFlag string
	var intFlag int
	var boolFlag bool
	mockParser.StringVar(&stringFlag, "preset-string", "default", "usage")
	mockParser.IntVar(&intFlag, "preset-int", 0, "usage")
	mockParser.BoolVar(&boolFlag, "preset-bool", false, "usage")

	if stringFlag != "preset-value" {
		t.Errorf("Expected preset string value 'preset-value', got '%s'", stringFlag)
	}
	if intFlag != 100 {
		t.Errorf("Expected preset int value 100, got %d", intFlag)
	}
	if !boolFlag {
		t.Errorf("Expected preset bool value true, got %t", boolFlag)
	}
}

// TestMockFlagParser_DefaultValues はデフォルト値のテスト
func TestMockFlagParser_DefaultValues(t *testing.T) {
	mockParser := NewMockFlagParser()

	var stringFlag string
	var intFlag int
	var boolFlag bool

	// 事前設定なしでフラグを定義（デフォルト値が使用されるはず）
	mockParser.StringVar(&stringFlag, "string-flag", "default-string", "string flag")
	mockParser.IntVar(&intFlag, "int-flag", 42, "int flag")
	mockParser.BoolVar(&boolFlag, "bool-flag", true, "bool flag")

	if stringFlag != "default-string" {
		t.Errorf("Expected default string 'default-string', got '%s'", stringFlag)
	}
	if intFlag != 42 {
		t.Errorf("Expected default int 42, got %d", intFlag)
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
