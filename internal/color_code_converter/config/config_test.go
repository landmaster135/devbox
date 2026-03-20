package config

import (
	"testing"
)

// MockFlagParser はテスト用のFlagParserモック
type MockFlagParser struct {
	stringValues map[string]string
	boolValues   map[string]bool
	stringVars   map[string]*string
	boolVars     map[string]*bool
	args         []string
}

// NewMockFlagParser は新しいMockFlagParserを作成する
func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringValues: make(map[string]string),
		boolValues:   make(map[string]bool),
		stringVars:   make(map[string]*string),
		boolVars:     make(map[string]*bool),
		args:         []string{},
	}
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

func (m *MockFlagParser) Parse() error {
	return nil
}

func (m *MockFlagParser) Args() []string {
	return m.args
}

// SetStringValue はモックに文字列値を設定する
func (m *MockFlagParser) SetStringValue(name, value string) {
	m.stringValues[name] = value
}

// SetBoolValue はモックにブール値を設定する
func (m *MockFlagParser) SetBoolValue(name string, value bool) {
	m.boolValues[name] = value
}

// SetArgs はモックに引数を設定する
func (m *MockFlagParser) SetArgs(args []string) {
	m.args = args
}

func TestNewConfig_Normal(t *testing.T) {
	const (
		testSrcFormat  = "hex"
		testDestFormat = "rgb"
		testValue      = "#FF0000"
		testDecValue   = "16711680"
	)

	tests := []struct {
		name           string
		srcFormat      string
		destFormat     string
		value          string
		expectError    bool
		expectedConfig *Config
	}{
		{
			name:        "ValidHexToRgb_Normal",
			srcFormat:   testSrcFormat,
			destFormat:  testDestFormat,
			value:       testValue,
			expectError: false,
			expectedConfig: &Config{
				SrcFormat:  testSrcFormat,
				DestFormat: testDestFormat,
				Value:      testValue,
			},
		},
		{
			name:        "ValidDecToHex_Normal",
			srcFormat:   "dec",
			destFormat:  "hex",
			value:       testDecValue,
			expectError: false,
			expectedConfig: &Config{
				SrcFormat:  "dec",
				DestFormat: "hex",
				Value:      testDecValue,
			},
		},
		{
			name:        "EmptySrcFormat_Error",
			srcFormat:   "",
			destFormat:  testDestFormat,
			value:       testValue,
			expectError: true,
		},
		{
			name:        "EmptyDestFormat_Error",
			srcFormat:   testSrcFormat,
			destFormat:  "",
			value:       testValue,
			expectError: true,
		},
		{
			name:        "EmptyValue_Error",
			srcFormat:   testSrcFormat,
			destFormat:  testDestFormat,
			value:       "",
			expectError: true,
		},
		{
			name:        "InvalidSrcFormat_Error",
			srcFormat:   "invalid",
			destFormat:  testDestFormat,
			value:       testValue,
			expectError: true,
		},
		{
			name:        "InvalidDestFormat_Error",
			srcFormat:   testSrcFormat,
			destFormat:  "invalid",
			value:       testValue,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := NewConfig(tt.srcFormat, tt.destFormat, tt.value)

			if tt.expectError {
				if err == nil {
					t.Errorf("期待されたエラーが発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
				return
			}

			if config.SrcFormat != tt.expectedConfig.SrcFormat {
				t.Errorf("SrcFormat = %v, want %v", config.SrcFormat, tt.expectedConfig.SrcFormat)
			}
			if config.DestFormat != tt.expectedConfig.DestFormat {
				t.Errorf("DestFormat = %v, want %v", config.DestFormat, tt.expectedConfig.DestFormat)
			}
			if config.Value != tt.expectedConfig.Value {
				t.Errorf("Value = %v, want %v", config.Value, tt.expectedConfig.Value)
			}
		})
	}
}

func TestParseFlags_Normal(t *testing.T) {
	const (
		testSrcFormat  = "hex"
		testDestFormat = "rgb"
		testValue      = "#FF0000"
		testDecValue   = "16711680"
	)

	tests := []struct {
		name           string
		setupMock      func(*MockFlagParser)
		expectError    bool
		expectedConfig *Config
	}{
		{
			name: "ValidFlags_Normal",
			setupMock: func(mock *MockFlagParser) {
				mock.SetStringValue("src-format", testSrcFormat)
				mock.SetStringValue("dest-format", testDestFormat)
				mock.SetStringValue("value", testValue)
			},
			expectError: false,
			expectedConfig: &Config{
				SrcFormat:  testSrcFormat,
				DestFormat: testDestFormat,
				Value:      testValue,
			},
		},
		{
			name: "HelpFlag_Normal",
			setupMock: func(mock *MockFlagParser) {
				mock.SetBoolValue("help", true)
			},
			expectError: false,
			expectedConfig: &Config{
				Help: true,
			},
		},
		{
			name: "PositionalArgs_Normal",
			setupMock: func(mock *MockFlagParser) {
				mock.SetArgs([]string{testSrcFormat, testDestFormat, testValue})
			},
			expectError: false,
			expectedConfig: &Config{
				SrcFormat:  testSrcFormat,
				DestFormat: testDestFormat,
				Value:      testValue,
			},
		},
		{
			name: "ShortFlags_Normal",
			setupMock: func(mock *MockFlagParser) {
				mock.SetStringValue("s", testSrcFormat)
				mock.SetStringValue("d", testDestFormat)
				mock.SetStringValue("v", testValue)
			},
			expectError: false,
			expectedConfig: &Config{
				SrcFormat:  testSrcFormat,
				DestFormat: testDestFormat,
				Value:      testValue,
			},
		},
		{
			name: "DecFormatFlags_Normal",
			setupMock: func(mock *MockFlagParser) {
				mock.SetStringValue("src-format", "dec")
				mock.SetStringValue("dest-format", "hex")
				mock.SetStringValue("value", testDecValue)
			},
			expectError: false,
			expectedConfig: &Config{
				SrcFormat:  "dec",
				DestFormat: "hex",
				Value:      testDecValue,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParser := NewMockFlagParser()
			tt.setupMock(mockParser)

			config, err := ParseFlagsWithParser(mockParser)

			if tt.expectError {
				if err == nil {
					t.Errorf("期待されたエラーが発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
				return
			}

			if tt.expectedConfig.Help {
				if !config.Help {
					t.Errorf("Help = %v, want %v", config.Help, tt.expectedConfig.Help)
				}
				return
			}

			if config.SrcFormat != tt.expectedConfig.SrcFormat {
				t.Errorf("SrcFormat = %v, want %v", config.SrcFormat, tt.expectedConfig.SrcFormat)
			}
			if config.DestFormat != tt.expectedConfig.DestFormat {
				t.Errorf("DestFormat = %v, want %v", config.DestFormat, tt.expectedConfig.DestFormat)
			}
			if config.Value != tt.expectedConfig.Value {
				t.Errorf("Value = %v, want %v", config.Value, tt.expectedConfig.Value)
			}
		})
	}
}

func TestIsValidFormat_Normal(t *testing.T) {
	validFormats := []string{"hex", "rgb", "hsl", "hsv", "dec"}

	tests := []struct {
		name     string
		format   string
		expected bool
	}{
		{
			name:     "ValidHex_Normal",
			format:   "hex",
			expected: true,
		},
		{
			name:     "ValidRgb_Normal",
			format:   "rgb",
			expected: true,
		},
		{
			name:     "ValidHsl_Normal",
			format:   "hsl",
			expected: true,
		},
		{
			name:     "ValidHsv_Normal",
			format:   "hsv",
			expected: true,
		},
		{
			name:     "ValidUpperCase_Normal",
			format:   "HEX",
			expected: true,
		},
		{
			name:     "ValidDec_Normal",
			format:   "dec",
			expected: true,
		},
		{
			name:     "InvalidFormat_Normal",
			format:   "invalid",
			expected: false,
		},
		{
			name:     "EmptyFormat_Normal",
			format:   "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidFormat(tt.format, validFormats)
			if result != tt.expected {
				t.Errorf("isValidFormat(%v) = %v, want %v", tt.format, result, tt.expected)
			}
		})
	}
}
