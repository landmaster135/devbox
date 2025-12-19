package config

import (
	"testing"
)

// MockFlagParser はテスト用のフラグパーサーです
type MockFlagParser struct {
	stringValues map[string]string
	boolValues   map[string]bool
	intValues    map[string]int
	int64Values  map[string]int64
	stringVars   map[string]*string
	boolVars     map[string]*bool
	intVars      map[string]*int
	int64Vars    map[string]*int64
	nArgValue    int
}

func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringValues: make(map[string]string),
		boolValues:   make(map[string]bool),
		intValues:    make(map[string]int),
		int64Values:  make(map[string]int64),
		stringVars:   make(map[string]*string),
		boolVars:     make(map[string]*bool),
		intVars:      make(map[string]*int),
		int64Vars:    make(map[string]*int64),
		nArgValue:    0,
	}
}

func (m *MockFlagParser) SetString(name, value string) {
	m.stringValues[name] = value
}

func (m *MockFlagParser) SetBool(name string, value bool) {
	m.boolValues[name] = value
}

func (m *MockFlagParser) SetInt(name string, value int) {
	m.intValues[name] = value
}

func (m *MockFlagParser) SetInt64(name string, value int64) {
	m.int64Values[name] = value
}

func (m *MockFlagParser) SetNArg(value int) {
	m.nArgValue = value
}

func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	if presetValue, exists := m.stringValues[name]; exists {
		*p = presetValue
	} else {
		*p = value
	}
	m.stringVars[name] = p
}

func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if presetValue, exists := m.boolValues[name]; exists {
		*p = presetValue
	} else {
		*p = value
	}
	m.boolVars[name] = p
}

func (m *MockFlagParser) IntVar(p *int, name string, value int, usage string) {
	if presetValue, exists := m.intValues[name]; exists {
		*p = presetValue
	} else {
		*p = value
	}
	m.intVars[name] = p
}

func (m *MockFlagParser) Int64Var(p *int64, name string, value int64, usage string) {
	if presetValue, exists := m.int64Values[name]; exists {
		*p = presetValue
	} else {
		*p = value
	}
	m.int64Vars[name] = p
}

func (m *MockFlagParser) Parse() {
	// モックでは何もしない
}

func (m *MockFlagParser) NArg() int {
	return m.nArgValue
}

func TestParseFlagsWithParser_ValidURL_Normal(t *testing.T) {
	const (
		testURL       = "https://www.youtube.com/watch?v=test123"
		testOutputDir = "./test-downloads"
		testQuality   = "720p"
		testFormat    = "mp4"
	)

	tests := []struct {
		name              string
		setupMock         func(*MockFlagParser)
		expectedURL       string
		expectedOutput    string
		expectedQuality   string
		expectedFormat    string
		expectedAudioOnly bool
		expectedPlaylist  bool
		expectError       bool
	}{
		{
			name: "ValidConfiguration_Normal",
			setupMock: func(mock *MockFlagParser) {
				mock.SetString("url", testURL)
				mock.SetString("output", testOutputDir)
				mock.SetString("quality", testQuality)
				mock.SetString("format", testFormat)
				mock.SetBool("audio-only", false)
				mock.SetBool("playlist", false)
				mock.SetInt("max-routines", 10)
				mock.SetInt64("chunk-size", 10*1024*1024)
				mock.SetBool("help", false)
			},
			expectedURL:       testURL,
			expectedOutput:    testOutputDir,
			expectedQuality:   testQuality,
			expectedFormat:    testFormat,
			expectedAudioOnly: false,
			expectedPlaylist:  false,
			expectError:       false,
		},
		{
			name: "AudioOnlyConfiguration_Normal",
			setupMock: func(mock *MockFlagParser) {
				mock.SetString("url", testURL)
				mock.SetString("output", testOutputDir)
				mock.SetString("quality", "best")
				mock.SetString("format", "m4a")
				mock.SetBool("audio-only", true)
				mock.SetBool("playlist", false)
				mock.SetInt("max-routines", 5)
				mock.SetInt64("chunk-size", 5*1024*1024)
				mock.SetBool("help", false)
			},
			expectedURL:       testURL,
			expectedOutput:    testOutputDir,
			expectedQuality:   "best",
			expectedFormat:    "m4a",
			expectedAudioOnly: true,
			expectedPlaylist:  false,
			expectError:       false,
		},
		{
			name: "PlaylistConfiguration_Normal",
			setupMock: func(mock *MockFlagParser) {
				mock.SetString("url", "https://www.youtube.com/playlist?list=test123")
				mock.SetString("output", testOutputDir)
				mock.SetString("quality", "best")
				mock.SetString("format", "mp4")
				mock.SetBool("audio-only", false)
				mock.SetBool("playlist", true)
				mock.SetInt("max-routines", 10)
				mock.SetInt64("chunk-size", 10*1024*1024)
				mock.SetBool("help", false)
			},
			expectedURL:       "https://www.youtube.com/playlist?list=test123",
			expectedOutput:    testOutputDir,
			expectedQuality:   "best",
			expectedFormat:    "mp4",
			expectedAudioOnly: false,
			expectedPlaylist:  true,
			expectError:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParser := NewMockFlagParser()
			tt.setupMock(mockParser)

			cfg, err := ParseFlagsWithParser(mockParser)

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

			if cfg.URL != tt.expectedURL {
				t.Errorf("URL = %v, 期待値 %v", cfg.URL, tt.expectedURL)
			}

			if cfg.Quality != tt.expectedQuality {
				t.Errorf("Quality = %v, 期待値 %v", cfg.Quality, tt.expectedQuality)
			}

			if cfg.Format != tt.expectedFormat {
				t.Errorf("Format = %v, 期待値 %v", cfg.Format, tt.expectedFormat)
			}

			if cfg.AudioOnly != tt.expectedAudioOnly {
				t.Errorf("AudioOnly = %v, 期待値 %v", cfg.AudioOnly, tt.expectedAudioOnly)
			}

			if cfg.Playlist != tt.expectedPlaylist {
				t.Errorf("Playlist = %v, 期待値 %v", cfg.Playlist, tt.expectedPlaylist)
			}
		})
	}
}

func TestParseFlagsWithParser_InvalidConfiguration_Error(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*MockFlagParser)
		errorMsg  string
	}{
		{
			name: "EmptyURL_Error",
			setupMock: func(mock *MockFlagParser) {
				mock.SetString("url", "")
				mock.SetBool("help", false)
			},
			errorMsg: "URLが指定されていません",
		},
		{
			name: "InvalidMaxRoutines_Error",
			setupMock: func(mock *MockFlagParser) {
				mock.SetString("url", "https://www.youtube.com/watch?v=test123")
				mock.SetInt("max-routines", 0)
				mock.SetBool("help", false)
			},
			errorMsg: "並列ダウンロード数は1以上である必要があります",
		},
		{
			name: "InvalidChunkSize_Error",
			setupMock: func(mock *MockFlagParser) {
				mock.SetString("url", "https://www.youtube.com/watch?v=test123")
				mock.SetInt64("chunk-size", 0)
				mock.SetBool("help", false)
			},
			errorMsg: "チャンクサイズは1以上である必要があります",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParser := NewMockFlagParser()
			tt.setupMock(mockParser)

			_, err := ParseFlagsWithParser(mockParser)

			if err == nil {
				t.Errorf("期待されたエラーが発生しませんでした")
				return
			}

			if !containsString(err.Error(), tt.errorMsg) {
				t.Errorf("エラーメッセージ = %v, 期待される文字列 %v", err.Error(), tt.errorMsg)
			}
		})
	}
}

func TestParseFlagsWithParser_HelpFlag_Normal(t *testing.T) {
	mockParser := NewMockFlagParser()
	mockParser.SetBool("help", true)

	cfg, err := ParseFlagsWithParser(mockParser)

	if err != nil {
		t.Errorf("予期しないエラーが発生しました: %v", err)
		return
	}

	if !cfg.Help {
		t.Errorf("Help = %v, 期待値 true", cfg.Help)
	}
}

func TestConfig_validate_Normal(t *testing.T) {
	cfg := &Config{
		URL:         "https://www.youtube.com/watch?v=test123",
		OutputDir:   "./downloads",
		Quality:     "720p",
		Format:      "mp4",
		AudioOnly:   false,
		Playlist:    false,
		MaxRoutines: 10,
		ChunkSize:   10 * 1024 * 1024,
		Help:        false,
	}

	err := cfg.validate()

	if err != nil {
		t.Errorf("予期しないエラーが発生しました: %v", err)
	}

	// 出力ディレクトリが絶対パスに変換されているかチェック
	if !isAbsolutePath(cfg.OutputDir) {
		t.Errorf("OutputDir が絶対パスに変換されていません: %v", cfg.OutputDir)
	}
}

// ヘルパー関数
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func isAbsolutePath(path string) bool {
	if len(path) == 0 {
		return false
	}
	// Unix系の絶対パス
	if path[0] == '/' {
		return true
	}
	// Windows系の絶対パス
	if len(path) >= 3 && path[1] == ':' && (path[2] == '\\' || path[2] == '/') {
		return true
	}
	return false
}
