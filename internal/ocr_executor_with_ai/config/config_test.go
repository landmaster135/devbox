package config

import (
	"os"
	"testing"
)

// TestConfig はConfig構造体のテストクラス
type TestConfig struct{}

func (t *TestConfig) TestNewConfig_ValidParams_Normal(test *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		test.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	const (
		testModel             = "gemini-2.0-flash"
		testPrompt            = "テスト用プロンプト"
		testSystemInstruction = "テスト用システム指示"
		testTemperature       = 0.8
		testMaxTokens         = 4096
		testAiType            = "gemini"
		testApiKey            = "test-api-key"
		testProject           = "test-project"
		testLocation          = "us-central1"
	)

	config, err := NewConfig(
		tempDir,
		true,
		testModel,
		testPrompt,
		testSystemInstruction,
		testTemperature,
		testMaxTokens,
		testAiType,
		testApiKey,
		testProject,
		testLocation,
		false,
	)

	if err != nil {
		test.Errorf("NewConfig() でエラーが発生しました: %v", err)
	}
	if config == nil {
		test.Fatal("config が nil です")
	}
	if config.Path != tempDir {
		test.Errorf("Path = %v, want %v", config.Path, tempDir)
	}
	if config.Recursive != true {
		test.Errorf("Recursive = %v, want %v", config.Recursive, true)
	}
	if config.Model != testModel {
		test.Errorf("Model = %v, want %v", config.Model, testModel)
	}
	if config.Prompt != testPrompt {
		test.Errorf("Prompt = %v, want %v", config.Prompt, testPrompt)
	}
	if config.SystemInstruction != testSystemInstruction {
		test.Errorf("SystemInstruction = %v, want %v", config.SystemInstruction, testSystemInstruction)
	}
	if config.Temperature != testTemperature {
		test.Errorf("Temperature = %v, want %v", config.Temperature, testTemperature)
	}
	if config.MaxTokens != testMaxTokens {
		test.Errorf("MaxTokens = %v, want %v", config.MaxTokens, testMaxTokens)
	}
	if config.AiType != testAiType {
		test.Errorf("AiType = %v, want %v", config.AiType, testAiType)
	}
	if config.APIKey != testApiKey {
		test.Errorf("APIKey = %v, want %v", config.APIKey, testApiKey)
	}
	if config.Project != testProject {
		test.Errorf("Project = %v, want %v", config.Project, testProject)
	}
	if config.Location != testLocation {
		test.Errorf("Location = %v, want %v", config.Location, testLocation)
	}
	if config.GeneratesMarkdownTable != false {
		test.Errorf("GeneratesMarkdownTable = %v, want %v", config.GeneratesMarkdownTable, false)
	}
}

func (t *TestConfig) TestNewConfig_InvalidPath_Error(test *testing.T) {
	const nonExistentPath = "/path/that/does/not/exist"

	config, err := NewConfig(
		nonExistentPath,
		false,
		DefaultModel,
		DefaultPrompt,
		DefaultSystemInstruction,
		DefaultTemperature,
		DefaultMaxTokens,
		DefaultAiType,
		"test-api-key",
		"",
		DefaultLocation,
		false,
	)

	if err == nil {
		test.Error("存在しないパスでエラーが発生しませんでした")
	}
	if config != nil {
		test.Error("エラー時にconfigがnilではありません")
	}
}

func (t *TestConfig) TestNewConfig_InvalidTemperature_Error(test *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		test.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		temperature float64
	}{
		{"NegativeTemperature", -0.1},
		{"TooHighTemperature", 2.1},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			config, err := NewConfig(
				tempDir,
				false,
				DefaultModel,
				DefaultPrompt,
				DefaultSystemInstruction,
				tt.temperature,
				DefaultMaxTokens,
				DefaultAiType,
				"test-api-key",
				"",
				DefaultLocation,
				false,
			)

			if err == nil {
				t.Error("無効な温度でエラーが発生しませんでした")
			}
			if config != nil {
				t.Error("エラー時にconfigがnilではありません")
			}
		})
	}
}

func (t *TestConfig) TestNewConfig_InvalidMaxTokens_Error(test *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		test.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name      string
		maxTokens int
	}{
		{"ZeroMaxTokens", 0},
		{"NegativeMaxTokens", -1},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			config, err := NewConfig(
				tempDir,
				false,
				DefaultModel,
				DefaultPrompt,
				DefaultSystemInstruction,
				DefaultTemperature,
				tt.maxTokens,
				DefaultAiType,
				"test-api-key",
				"",
				DefaultLocation,
				false,
			)

			if err == nil {
				t.Error("無効な最大トークン数でエラーが発生しませんでした")
			}
			if config != nil {
				t.Error("エラー時にconfigがnilではありません")
			}
		})
	}
}

func (t *TestConfig) TestNewConfig_InvalidAiType_Error(test *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		test.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	const invalidAiType = "invalid-ai-type"

	config, err := NewConfig(
		tempDir,
		false,
		DefaultModel,
		DefaultPrompt,
		DefaultSystemInstruction,
		DefaultTemperature,
		DefaultMaxTokens,
		invalidAiType,
		"test-api-key",
		"",
		DefaultLocation,
		false,
	)

	if err == nil {
		test.Error("無効なAIタイプでエラーが発生しませんでした")
	}
	if config != nil {
		test.Error("エラー時にconfigがnilではありません")
	}
}

func (t *TestConfig) TestNewConfig_MissingApiKey_Error(test *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		test.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config, err := NewConfig(
		tempDir,
		false,
		DefaultModel,
		DefaultPrompt,
		DefaultSystemInstruction,
		DefaultTemperature,
		DefaultMaxTokens,
		"gemini",
		"", // 空のAPIキー
		"",
		DefaultLocation,
		false,
	)

	if err == nil {
		test.Error("APIキー不足でエラーが発生しませんでした")
	}
	if config != nil {
		test.Error("エラー時にconfigがnilではありません")
	}
}

func (t *TestConfig) TestNewConfig_MissingProject_Error(test *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		test.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		project     string
		location    string
		expectError bool
	}{
		{
			name:        "BothEmpty_Error",
			project:     "",
			location:    "",
			expectError: true,
		},
		{
			name:        "ProjectOnly_Normal",
			project:     "test-project",
			location:    "",
			expectError: false,
		},
		{
			name:        "LocationOnly_Normal",
			project:     "",
			location:    "us-central1",
			expectError: false,
		},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			config, err := NewConfig(
				tempDir,
				false,
				DefaultModel,
				DefaultPrompt,
				DefaultSystemInstruction,
				DefaultTemperature,
				DefaultMaxTokens,
				"vertex",
				"",
				tt.project,
				tt.location,
				false,
			)

			if tt.expectError {
				if err == nil {
					t.Error("プロジェクト・ロケーション不足でエラーが発生しませんでした")
				}
				if config != nil {
					t.Error("エラー時にconfigがnilではありません")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
				if config == nil {
					t.Error("正常時にconfigがnilです")
				}
			}
		})
	}
}

func (t *TestConfig) TestNewConfig_ModelDefaultsByAiType(test *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		test.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testCases := []struct {
		name          string
		aiType        string
		model         string
		expectedModel string
		apiKey        string
		project       string
		location      string
		expectNoError bool
	}{
		{
			name:          "GeminiDefaultModelConverted",
			aiType:        "gemini",
			model:         DefaultModel,
			expectedModel: DefaultGeminiModel,
			apiKey:        "test-api-key",
			expectNoError: true,
		},
		{
			name:          "GeminiCustomModelMaintained",
			aiType:        "gemini",
			model:         "gemini-custom-model",
			expectedModel: "gemini-custom-model",
			apiKey:        "test-api-key",
			expectNoError: true,
		},
		{
			name:          "VertexDefaultModelConverted",
			aiType:        "vertex",
			model:         DefaultModel,
			expectedModel: DefaultVertexModel,
			project:       "test-project",
			location:      DefaultLocation,
			expectNoError: true,
		},
		{
			name:          "VertexCustomModelMaintained",
			aiType:        "vertex",
			model:         "vertex-custom-model",
			expectedModel: "vertex-custom-model",
			project:       "test-project",
			location:      DefaultLocation,
			expectNoError: true,
		},
		{
			name:          "OllamaDefaultModel",
			aiType:        "ollama",
			model:         "",
			expectedModel: DefaultModel,
			expectNoError: true,
		},
		{
			name:          "OllamaCustomModelMaintained",
			aiType:        "ollama",
			model:         "ollama-custom",
			expectedModel: "ollama-custom",
			expectNoError: true,
		},
	}

	for _, tt := range testCases {
		test.Run(tt.name, func(t *testing.T) {
			config, err := NewConfig(
				tempDir,
				false,
				tt.model,
				DefaultPrompt,
				DefaultSystemInstruction,
				DefaultTemperature,
				DefaultMaxTokens,
				tt.aiType,
				tt.apiKey,
				tt.project,
				tt.location,
				false,
			)

			if !tt.expectNoError {
				if err == nil {
					t.Fatal("期待されたエラーが発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Fatalf("NewConfig() がエラーを返しました: %v", err)
			}
			if config.Model != tt.expectedModel {
				t.Errorf("Model = %v, want %v", config.Model, tt.expectedModel)
			}
		})
	}
}

func (t *TestConfig) TestNewConfig_MarkdownTableConflict_Error(test *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		test.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	const customPrompt = "カスタムプロンプト"

	config, err := NewConfig(
		tempDir,
		false,
		DefaultModel,
		customPrompt, // デフォルト以外のプロンプト
		DefaultSystemInstruction,
		DefaultTemperature,
		DefaultMaxTokens,
		DefaultAiType,
		"test-api-key",
		"",
		DefaultLocation,
		true, // Markdownテーブルフラグを有効
	)

	if err == nil {
		test.Error("Markdownテーブルフラグとプロンプトの競合でエラーが発生しませんでした")
	}
	if config != nil {
		test.Error("エラー時にconfigがnilではありません")
	}
}

func (t *TestConfig) TestNewConfig_MarkdownTableFlag_Normal(test *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		test.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config, err := NewConfig(
		tempDir,
		false,
		DefaultModel,
		DefaultPrompt, // デフォルトプロンプト
		DefaultSystemInstruction,
		DefaultTemperature,
		DefaultMaxTokens,
		DefaultAiType,
		"test-api-key",
		"",
		DefaultLocation,
		true, // Markdownテーブルフラグを有効
	)

	if err != nil {
		test.Errorf("NewConfig() でエラーが発生しました: %v", err)
	}
	if config == nil {
		test.Fatal("config が nil です")
	}
	if config.Prompt != DefaultMarkdownTablePrompt {
		test.Errorf("Prompt = %v, want %v", config.Prompt, DefaultMarkdownTablePrompt)
	}
	if config.GeneratesMarkdownTable != true {
		test.Errorf("GeneratesMarkdownTable = %v, want %v", config.GeneratesMarkdownTable, true)
	}
}

// TestParseFlagsWithParser はParseFlagsWithParser関数のテストクラス
type TestParseFlagsWithParser struct{}

func (t *TestParseFlagsWithParser) TestParseFlagsWithParser_ValidFlags_Normal(test *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		test.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// モックフラグパーサーを作成
	mockParser := &MockFlagParser{
		stringValues:  make(map[string]string),
		boolValues:    make(map[string]bool),
		float64Values: make(map[string]float64),
		intValues:     make(map[string]int),
		stringVars:    make(map[string]*string),
		boolVars:      make(map[string]*bool),
		float64Vars:   make(map[string]*float64),
		intVars:       make(map[string]*int),
	}

	// テスト値を設定
	mockParser.stringValues["path"] = tempDir
	mockParser.stringValues["model"] = "gemini-2.0-flash"
	mockParser.stringValues["prompt"] = "テスト用プロンプト"
	mockParser.stringValues["ai-type"] = "gemini"
	mockParser.stringValues["api-key"] = "test-api-key"
	mockParser.boolValues["recursive"] = true
	mockParser.float64Values["temperature"] = 0.8
	mockParser.intValues["max-tokens"] = 4096

	config, err := ParseFlagsWithParser(mockParser)

	if err != nil {
		test.Errorf("ParseFlagsWithParser() でエラーが発生しました: %v", err)
	}
	if config == nil {
		test.Fatal("config が nil です")
	}
	if config.Path != tempDir {
		test.Errorf("Path = %v, want %v", config.Path, tempDir)
	}
	if config.Model != "gemini-2.0-flash" {
		test.Errorf("Model = %v, want %v", config.Model, "gemini-2.0-flash")
	}
	if config.Prompt != "テスト用プロンプト" {
		test.Errorf("Prompt = %v, want %v", config.Prompt, "テスト用プロンプト")
	}
	if config.Recursive != true {
		test.Errorf("Recursive = %v, want %v", config.Recursive, true)
	}
	if config.Temperature != 0.8 {
		test.Errorf("Temperature = %v, want %v", config.Temperature, 0.8)
	}
	if config.MaxTokens != 4096 {
		test.Errorf("MaxTokens = %v, want %v", config.MaxTokens, 4096)
	}
}

func (t *TestParseFlagsWithParser) TestParseFlagsWithParser_HelpFlag_Normal(test *testing.T) {
	// モックフラグパーサーを作成
	mockParser := &MockFlagParser{
		stringValues:  make(map[string]string),
		boolValues:    make(map[string]bool),
		float64Values: make(map[string]float64),
		intValues:     make(map[string]int),
		stringVars:    make(map[string]*string),
		boolVars:      make(map[string]*bool),
		float64Vars:   make(map[string]*float64),
		intVars:       make(map[string]*int),
	}

	// ヘルプフラグを設定
	mockParser.boolValues["help"] = true

	config, err := ParseFlagsWithParser(mockParser)

	if err != nil {
		test.Errorf("ParseFlagsWithParser() でエラーが発生しました: %v", err)
	}
	if config == nil {
		test.Fatal("config が nil です")
	}
	if config.Help != true {
		test.Errorf("Help = %v, want %v", config.Help, true)
	}
}

func (t *TestParseFlagsWithParser) TestParseFlagsWithParser_ShortFlags_Normal(test *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		test.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// モックフラグパーサーを作成
	mockParser := &MockFlagParser{
		stringValues:  make(map[string]string),
		boolValues:    make(map[string]bool),
		float64Values: make(map[string]float64),
		intValues:     make(map[string]int),
		stringVars:    make(map[string]*string),
		boolVars:      make(map[string]*bool),
		float64Vars:   make(map[string]*float64),
		intVars:       make(map[string]*int),
	}

	// 短縮フラグを設定
	mockParser.stringValues["p"] = tempDir
	mockParser.stringValues["m"] = "gemini-1.5-pro"
	mockParser.stringValues["at"] = "gemini"
	mockParser.stringValues["ak"] = "test-api-key"
	mockParser.boolValues["r"] = true
	mockParser.float64Values["t"] = 0.5
	mockParser.intValues["mt"] = 2048

	config, err := ParseFlagsWithParser(mockParser)

	if err != nil {
		test.Errorf("ParseFlagsWithParser() でエラーが発生しました: %v", err)
	}
	if config == nil {
		test.Fatal("config が nil です")
	}
	if config.Path != tempDir {
		test.Errorf("Path = %v, want %v", config.Path, tempDir)
	}
	if config.Model != "gemini-1.5-pro" {
		test.Errorf("Model = %v, want %v", config.Model, "gemini-1.5-pro")
	}
	if config.AiType != "gemini" {
		test.Errorf("AiType = %v, want %v", config.AiType, "gemini")
	}
	if config.APIKey != "test-api-key" {
		test.Errorf("APIKey = %v, want %v", config.APIKey, "test-api-key")
	}
	if config.Recursive != true {
		test.Errorf("Recursive = %v, want %v", config.Recursive, true)
	}
	if config.Temperature != 0.5 {
		test.Errorf("Temperature = %v, want %v", config.Temperature, 0.5)
	}
	if config.MaxTokens != 2048 {
		test.Errorf("MaxTokens = %v, want %v", config.MaxTokens, 2048)
	}
}

func (t *TestParseFlagsWithParser) TestParseFlagsWithParser_MarkdownTableFlag_Normal(test *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		test.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// モックフラグパーサーを作成
	mockParser := &MockFlagParser{
		stringValues:  make(map[string]string),
		boolValues:    make(map[string]bool),
		float64Values: make(map[string]float64),
		intValues:     make(map[string]int),
		stringVars:    make(map[string]*string),
		boolVars:      make(map[string]*bool),
		float64Vars:   make(map[string]*float64),
		intVars:       make(map[string]*int),
	}

	// Markdownテーブルフラグを設定
	mockParser.stringValues["path"] = tempDir
	mockParser.stringValues["ai-type"] = "gemini"
	mockParser.stringValues["api-key"] = "test-api-key"
	mockParser.boolValues["generates-markdown-table"] = true

	config, err := ParseFlagsWithParser(mockParser)

	if err != nil {
		test.Errorf("ParseFlagsWithParser() でエラーが発生しました: %v", err)
	}
	if config == nil {
		test.Fatal("config が nil です")
	}
	if config.GeneratesMarkdownTable != true {
		test.Errorf("GeneratesMarkdownTable = %v, want %v", config.GeneratesMarkdownTable, true)
	}
	if config.Prompt != DefaultMarkdownTablePrompt {
		test.Errorf("Prompt = %v, want %v", config.Prompt, DefaultMarkdownTablePrompt)
	}
}

// テスト実行用の関数
func TestNewConfig_ValidParams_Normal(t *testing.T) {
	testInstance := &TestConfig{}
	testInstance.TestNewConfig_ValidParams_Normal(t)
}

func TestNewConfig_InvalidPath_Error(t *testing.T) {
	testInstance := &TestConfig{}
	testInstance.TestNewConfig_InvalidPath_Error(t)
}

func TestNewConfig_InvalidTemperature_Error(t *testing.T) {
	testInstance := &TestConfig{}
	testInstance.TestNewConfig_InvalidTemperature_Error(t)
}

func TestNewConfig_InvalidMaxTokens_Error(t *testing.T) {
	testInstance := &TestConfig{}
	testInstance.TestNewConfig_InvalidMaxTokens_Error(t)
}

func TestNewConfig_InvalidAiType_Error(t *testing.T) {
	testInstance := &TestConfig{}
	testInstance.TestNewConfig_InvalidAiType_Error(t)
}

func TestNewConfig_MissingApiKey_Error(t *testing.T) {
	testInstance := &TestConfig{}
	testInstance.TestNewConfig_MissingApiKey_Error(t)
}

func TestNewConfig_MissingProject_Error(t *testing.T) {
	testInstance := &TestConfig{}
	testInstance.TestNewConfig_MissingProject_Error(t)
}

func TestNewConfig_MarkdownTableConflict_Error(t *testing.T) {
	testInstance := &TestConfig{}
	testInstance.TestNewConfig_MarkdownTableConflict_Error(t)
}

func TestNewConfig_MarkdownTableFlag_Normal(t *testing.T) {
	testInstance := &TestConfig{}
	testInstance.TestNewConfig_MarkdownTableFlag_Normal(t)
}

func TestParseFlagsWithParser_ValidFlags_Normal(t *testing.T) {
	testInstance := &TestParseFlagsWithParser{}
	testInstance.TestParseFlagsWithParser_ValidFlags_Normal(t)
}

func TestParseFlagsWithParser_HelpFlag_Normal(t *testing.T) {
	testInstance := &TestParseFlagsWithParser{}
	testInstance.TestParseFlagsWithParser_HelpFlag_Normal(t)
}

func TestParseFlagsWithParser_ShortFlags_Normal(t *testing.T) {
	testInstance := &TestParseFlagsWithParser{}
	testInstance.TestParseFlagsWithParser_ShortFlags_Normal(t)
}

func TestParseFlagsWithParser_MarkdownTableFlag_Normal(t *testing.T) {
	testInstance := &TestParseFlagsWithParser{}
	testInstance.TestParseFlagsWithParser_MarkdownTableFlag_Normal(t)
}

// MockFlagParser はテスト用のFlagParserモック実装
type MockFlagParser struct {
	stringValues  map[string]string
	boolValues    map[string]bool
	float64Values map[string]float64
	intValues     map[string]int
	stringVars    map[string]*string
	boolVars      map[string]*bool
	float64Vars   map[string]*float64
	intVars       map[string]*int
	parseError    error
}

func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	if presetValue, exists := m.stringValues[name]; exists {
		*p = presetValue // 事前設定値を優先
	} else {
		*p = value // デフォルト値
	}
	m.stringVars[name] = p
}

func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if presetValue, exists := m.boolValues[name]; exists {
		*p = presetValue // 事前設定値を優先
	} else {
		*p = value // デフォルト値
	}
	m.boolVars[name] = p
}

func (m *MockFlagParser) Float64Var(p *float64, name string, value float64, usage string) {
	if presetValue, exists := m.float64Values[name]; exists {
		*p = presetValue // 事前設定値を優先
	} else {
		*p = value // デフォルト値
	}
	m.float64Vars[name] = p
}

func (m *MockFlagParser) IntVar(p *int, name string, value int, usage string) {
	if presetValue, exists := m.intValues[name]; exists {
		*p = presetValue // 事前設定値を優先
	} else {
		*p = value // デフォルト値
	}
	m.intVars[name] = p
}

func (m *MockFlagParser) Parse() error {
	return m.parseError
}
