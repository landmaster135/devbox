package config

import (
	"fmt"
	"testing"
)

// MockFlagParser はテスト用のFlagParserモック
type MockFlagParser struct {
	stringValues map[string]string
	boolValues   map[string]bool
	intValues    map[string]int
	parseError   error
	parseCalled  bool
}

// NewMockFlagParser は新しいMockFlagParserを作成する
func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringValues: make(map[string]string),
		boolValues:   make(map[string]bool),
		intValues:    make(map[string]int),
	}
}

// StringVar はstring型の変数を登録する
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	// モックに設定された値があればそれを使用、なければデフォルト値
	if val, exists := m.stringValues[name]; exists {
		*p = val
	} else {
		*p = value
	}
}

// BoolVar はbool型の変数を登録する
func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	// モックに設定された値があればそれを使用、なければデフォルト値
	if val, exists := m.boolValues[name]; exists {
		*p = val
	} else {
		*p = value
	}
}

// IntVar はint型の変数を登録する
func (m *MockFlagParser) IntVar(p *int, name string, value int, usage string) {
	// モックに設定された値があればそれを使用、なければデフォルト値
	if val, exists := m.intValues[name]; exists {
		*p = val
	} else {
		*p = value
	}
}

// Parse はフラグの解析を実行する
func (m *MockFlagParser) Parse() error {
	m.parseCalled = true
	return m.parseError
}

// SetStringValue はstring型の値を事前に設定する
func (m *MockFlagParser) SetStringValue(name, value string) {
	m.stringValues[name] = value
}

// SetBoolValue はbool型の値を事前に設定する
func (m *MockFlagParser) SetBoolValue(name string, value bool) {
	m.boolValues[name] = value
}

// SetIntValue はint型の値を事前に設定する
func (m *MockFlagParser) SetIntValue(name string, value int) {
	m.intValues[name] = value
}

// SetParseError はParse()メソッドで返すエラーを設定する
func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}

func TestParseFlagsWithParser_Normal(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockFlagParser)
		expectedConfig *Config
		wantErr        bool
		errContains    string
	}{
		{
			name: "valid games operation",
			setupMock: func(mock *MockFlagParser) {
				mock.SetStringValue("operation", "games")
				mock.SetStringValue("steam-api-key", "test-api-key")
				mock.SetStringValue("steam-id", "76561198000000000")
			},
			expectedConfig: &Config{
				Operation:   "games",
				SteamAPIKey: "test-api-key",
				SteamID:     "76561198000000000",
				GameID:      0,
				Help:        false,
			},
			wantErr: false,
		},
		{
			name: "valid game-stats operation",
			setupMock: func(mock *MockFlagParser) {
				mock.SetStringValue("operation", "game-stats")
				mock.SetStringValue("steam-api-key", "test-api-key")
				mock.SetStringValue("steam-id", "76561198000000000")
				mock.SetIntValue("game-id", 123)
			},
			expectedConfig: &Config{
				Operation:   "game-stats",
				SteamAPIKey: "test-api-key",
				SteamID:     "76561198000000000",
				GameID:      123,
				Help:        false,
			},
			wantErr: false,
		},
		{
			name: "help flag set",
			setupMock: func(mock *MockFlagParser) {
				mock.SetBoolValue("help", true)
			},
			expectedConfig: &Config{
				Help: true,
			},
			wantErr: false,
		},
		{
			name: "parse error",
			setupMock: func(mock *MockFlagParser) {
				mock.SetParseError(fmt.Errorf("parse failed"))
			},
			wantErr:     true,
			errContains: "フラグの解析に失敗しました",
		},
		{
			name: "validation error - empty operation",
			setupMock: func(mock *MockFlagParser) {
				mock.SetStringValue("steam-api-key", "test-api-key")
				mock.SetStringValue("steam-id", "76561198000000000")
			},
			wantErr:     true,
			errContains: "設定の初期化に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockFlagParser()
			tt.setupMock(mock)

			config, err := ParseFlagsWithParser(mock)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFlagsWithParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errContains != "" && err != nil {
					// エラーメッセージに期待する文字列が含まれているかチェック
					found := false
					if len(err.Error()) >= len(tt.errContains) {
						for i := 0; i <= len(err.Error())-len(tt.errContains); i++ {
							if err.Error()[i:i+len(tt.errContains)] == tt.errContains {
								found = true
								break
							}
						}
					}
					if !found {
						t.Errorf("ParseFlagsWithParser() error = %v, want error containing %v", err, tt.errContains)
					}
				}
				return
			}

			if tt.expectedConfig != nil && tt.expectedConfig.Help {
				// ヘルプフラグの場合は、Helpフラグのみをチェック
				if config.Help != tt.expectedConfig.Help {
					t.Errorf("ParseFlagsWithParser() Help = %v, want %v", config.Help, tt.expectedConfig.Help)
				}
			} else if tt.expectedConfig != nil {
				if config.Operation != tt.expectedConfig.Operation {
					t.Errorf("ParseFlagsWithParser() Operation = %v, want %v", config.Operation, tt.expectedConfig.Operation)
				}
				if config.SteamAPIKey != tt.expectedConfig.SteamAPIKey {
					t.Errorf("ParseFlagsWithParser() SteamAPIKey = %v, want %v", config.SteamAPIKey, tt.expectedConfig.SteamAPIKey)
				}
				if config.SteamID != tt.expectedConfig.SteamID {
					t.Errorf("ParseFlagsWithParser() SteamID = %v, want %v", config.SteamID, tt.expectedConfig.SteamID)
				}
				if config.GameID != tt.expectedConfig.GameID {
					t.Errorf("ParseFlagsWithParser() GameID = %v, want %v", config.GameID, tt.expectedConfig.GameID)
				}
				if config.Help != tt.expectedConfig.Help {
					t.Errorf("ParseFlagsWithParser() Help = %v, want %v", config.Help, tt.expectedConfig.Help)
				}
			}

			if !mock.parseCalled {
				t.Error("ParseFlagsWithParser() did not call Parse() on the mock")
			}
		})
	}
}

func TestStandardFlagParser_Normal(t *testing.T) {
	parser := NewStandardFlagParser()
	if parser == nil {
		t.Error("NewStandardFlagParser() returned nil")
	}

	// 基本的な動作確認のため、変数を設定してみる
	var testString string
	var testBool bool
	var testInt int

	parser.StringVar(&testString, "test-string", "default", "test string")
	parser.BoolVar(&testBool, "test-bool", false, "test bool")
	parser.IntVar(&testInt, "test-int", 42, "test int")

	// Parse()を呼び出してもエラーが発生しないことを確認
	err := parser.Parse()
	if err != nil {
		t.Errorf("StandardFlagParser.Parse() error = %v, want nil", err)
	}
}
