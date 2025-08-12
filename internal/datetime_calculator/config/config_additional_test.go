package config

import (
	"testing"
)

// TestNewConfigForParseTime_Normal はNewConfigForParseTime関数の正常系テスト
func TestNewConfigForParseTime_Normal(t *testing.T) {
	testCases := []struct {
		name       string
		operation  string
		filePath   string
		textInput  string
		outputUnit string
		wantErr    bool
	}{
		{
			name:       "ファイルパス指定_txtファイル",
			operation:  "parse-time",
			filePath:   "/path/to/test.txt",
			textInput:  "",
			outputUnit: "minute",
			wantErr:    false,
		},
		{
			name:       "ファイルパス指定_mdファイル",
			operation:  "parse-time",
			filePath:   "/path/to/test.md",
			textInput:  "",
			outputUnit: "hour",
			wantErr:    false,
		},
		{
			name:       "テキスト入力指定",
			operation:  "parse-time",
			filePath:   "",
			textInput:  "作業は合計30分掛かった。",
			outputUnit: "minute",
			wantErr:    false,
		},
		{
			name:       "出力単位デフォルト",
			operation:  "parse-time",
			filePath:   "",
			textInput:  "作業は合計30分掛かった。",
			outputUnit: "",
			wantErr:    false,
		},
		{
			name:       "全ての有効な出力単位_year",
			operation:  "parse-time",
			filePath:   "",
			textInput:  "作業は合計30分掛かった。",
			outputUnit: "year",
			wantErr:    false,
		},
		{
			name:       "全ての有効な出力単位_month",
			operation:  "parse-time",
			filePath:   "",
			textInput:  "作業は合計30分掛かった。",
			outputUnit: "month",
			wantErr:    false,
		},
		{
			name:       "全ての有効な出力単位_day",
			operation:  "parse-time",
			filePath:   "",
			textInput:  "作業は合計30分掛かった。",
			outputUnit: "day",
			wantErr:    false,
		},
		{
			name:       "全ての有効な出力単位_hour",
			operation:  "parse-time",
			filePath:   "",
			textInput:  "作業は合計30分掛かった。",
			outputUnit: "hour",
			wantErr:    false,
		},
		{
			name:       "全ての有効な出力単位_minute",
			operation:  "parse-time",
			filePath:   "",
			textInput:  "作業は合計30分掛かった。",
			outputUnit: "minute",
			wantErr:    false,
		},
		{
			name:       "全ての有効な出力単位_second",
			operation:  "parse-time",
			filePath:   "",
			textInput:  "作業は合計30分掛かった。",
			outputUnit: "second",
			wantErr:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			config, err := NewConfigForParseTime(tc.operation, tc.filePath, tc.textInput, tc.outputUnit)

			// Assert
			if (err != nil) != tc.wantErr {
				t.Errorf("NewConfigForParseTime() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				if config == nil {
					t.Fatal("Expected non-nil config")
				}
				if config.Operation != tc.operation {
					t.Errorf("Expected operation %s, got %s", tc.operation, config.Operation)
				}
				if config.FilePath != tc.filePath {
					t.Errorf("Expected FilePath %s, got %s", tc.filePath, config.FilePath)
				}
				if config.TextInput != tc.textInput {
					t.Errorf("Expected TextInput %s, got %s", tc.textInput, config.TextInput)
				}

				expectedOutputUnit := tc.outputUnit
				if expectedOutputUnit == "" {
					expectedOutputUnit = "minute" // デフォルト値
				}
				if config.OutputUnit != expectedOutputUnit {
					t.Errorf("Expected OutputUnit %s, got %s", expectedOutputUnit, config.OutputUnit)
				}
			}
		})
	}
}

// TestNewConfigForParseTime_ErrorCases はNewConfigForParseTime関数のエラーケーステスト
func TestNewConfigForParseTime_ErrorCases(t *testing.T) {
	testCases := []struct {
		name       string
		operation  string
		filePath   string
		textInput  string
		outputUnit string
		wantErr    bool
	}{
		{
			name:       "空の操作タイプ",
			operation:  "",
			filePath:   "/path/to/test.txt",
			textInput:  "",
			outputUnit: "minute",
			wantErr:    true,
		},
		{
			name:       "無効な操作タイプ",
			operation:  "invalid",
			filePath:   "/path/to/test.txt",
			textInput:  "",
			outputUnit: "minute",
			wantErr:    true,
		},
		{
			name:       "両方指定",
			operation:  "parse-time",
			filePath:   "/path/to/test.txt",
			textInput:  "テキスト",
			outputUnit: "minute",
			wantErr:    true,
		},
		{
			name:       "両方未指定",
			operation:  "parse-time",
			filePath:   "",
			textInput:  "",
			outputUnit: "minute",
			wantErr:    true,
		},
		{
			name:       "無効なファイル拡張子_pdf",
			operation:  "parse-time",
			filePath:   "/path/to/test.pdf",
			textInput:  "",
			outputUnit: "minute",
			wantErr:    true,
		},
		{
			name:       "無効なファイル拡張子_doc",
			operation:  "parse-time",
			filePath:   "/path/to/test.doc",
			textInput:  "",
			outputUnit: "minute",
			wantErr:    true,
		},
		{
			name:       "無効なファイル拡張子_拡張子なし",
			operation:  "parse-time",
			filePath:   "/path/to/test",
			textInput:  "",
			outputUnit: "minute",
			wantErr:    true,
		},
		{
			name:       "無効な出力単位",
			operation:  "parse-time",
			filePath:   "",
			textInput:  "テキスト",
			outputUnit: "invalid",
			wantErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			config, err := NewConfigForParseTime(tc.operation, tc.filePath, tc.textInput, tc.outputUnit)

			// Assert
			if (err != nil) != tc.wantErr {
				t.Errorf("NewConfigForParseTime() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if tc.wantErr {
				if config != nil {
					t.Error("Expected nil config for error case")
				}
			}
		})
	}
}

// TestNewConfigForParseTime_FileExtensions はファイル拡張子の詳細テスト
func TestNewConfigForParseTime_FileExtensions(t *testing.T) {
	testCases := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "txtファイル_正常",
			filePath: "/path/to/test.txt",
			wantErr:  false,
		},
		{
			name:     "mdファイル_正常",
			filePath: "/path/to/test.md",
			wantErr:  false,
		},
		{
			name:     "複数ドット_txt_正常",
			filePath: "/path/to/test.backup.txt",
			wantErr:  false,
		},
		{
			name:     "複数ドット_md_正常",
			filePath: "/path/to/test.backup.md",
			wantErr:  false,
		},
		{
			name:     "大文字拡張子_TXT_正常",
			filePath: "/path/to/test.TXT",
			wantErr:  true, // 大文字は無効
		},
		{
			name:     "大文字拡張子_MD_正常",
			filePath: "/path/to/test.MD",
			wantErr:  true, // 大文字は無効
		},
		{
			name:     "pdfファイル_エラー",
			filePath: "/path/to/test.pdf",
			wantErr:  true,
		},
		{
			name:     "docファイル_エラー",
			filePath: "/path/to/test.doc",
			wantErr:  true,
		},
		{
			name:     "docxファイル_エラー",
			filePath: "/path/to/test.docx",
			wantErr:  true,
		},
		{
			name:     "xlsxファイル_エラー",
			filePath: "/path/to/test.xlsx",
			wantErr:  true,
		},
		{
			name:     "拡張子なし_エラー",
			filePath: "/path/to/test",
			wantErr:  true,
		},
		{
			name:     "ドットのみ_エラー",
			filePath: "/path/to/test.",
			wantErr:  true,
		},
		{
			name:     "隠しファイル_txt_正常",
			filePath: "/path/to/.hidden.txt",
			wantErr:  false,
		},
		{
			name:     "隠しファイル_md_正常",
			filePath: "/path/to/.hidden.md",
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			config, err := NewConfigForParseTime("parse-time", tc.filePath, "", "minute")

			// Assert
			if (err != nil) != tc.wantErr {
				t.Errorf("NewConfigForParseTime() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				if config == nil {
					t.Fatal("Expected non-nil config")
				}
				if config.FilePath != tc.filePath {
					t.Errorf("Expected FilePath %s, got %s", tc.filePath, config.FilePath)
				}
			}
		})
	}
}

// TestParseFlagsWithParser_ParseTimeOperation はparse-time操作のフラグ解析テスト
func TestParseFlagsWithParser_ParseTimeOperation(t *testing.T) {
	testCases := []struct {
		name       string
		flags      map[string]string
		expected   *Config
		wantErr    bool
	}{
		{
			name: "ファイルパス指定_正常",
			flags: map[string]string{
				"operation":  "parse-time",
				"file-path":  "/path/to/test.txt",
				"output-unit": "minute",
			},
			expected: &Config{
				Operation:  "parse-time",
				FilePath:   "/path/to/test.txt",
				TextInput:  "",
				OutputUnit: "minute",
			},
			wantErr: false,
		},
		{
			name: "テキスト入力指定_正常",
			flags: map[string]string{
				"operation":   "parse-time",
				"text-input":  "作業は合計30分掛かった。",
				"output-unit": "hour",
			},
			expected: &Config{
				Operation:  "parse-time",
				FilePath:   "",
				TextInput:  "作業は合計30分掛かった。",
				OutputUnit: "hour",
			},
			wantErr: false,
		},
		{
			name: "短縮形フラグ_正常",
			flags: map[string]string{
				"o":  "parse-time",
				"fp": "/path/to/test.md",
				"ou": "second",
			},
			expected: &Config{
				Operation:  "parse-time",
				FilePath:   "/path/to/test.md",
				TextInput:  "",
				OutputUnit: "second",
			},
			wantErr: false,
		},
		{
			name: "短縮形フラグ_テキスト入力",
			flags: map[string]string{
				"o":  "parse-time",
				"ti": "会議は合計90分掛かった。",
				"ou": "hour",
			},
			expected: &Config{
				Operation:  "parse-time",
				FilePath:   "",
				TextInput:  "会議は合計90分掛かった。",
				OutputUnit: "hour",
			},
			wantErr: false,
		},
		{
			name: "出力単位デフォルト",
			flags: map[string]string{
				"operation":  "parse-time",
				"text-input": "作業は合計30分掛かった。",
			},
			expected: &Config{
				Operation:  "parse-time",
				FilePath:   "",
				TextInput:  "作業は合計30分掛かった。",
				OutputUnit: "minute", // デフォルト値
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockParser := NewMockFlagParser()
			for key, value := range tc.flags {
				mockParser.SetStringFlag(key, value)
			}

			// Act
			config, err := ParseFlagsWithParser(mockParser)

			// Assert
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseFlagsWithParser() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				if config == nil {
					t.Fatal("Expected non-nil config")
				}
				if config.Operation != tc.expected.Operation {
					t.Errorf("Expected Operation %s, got %s", tc.expected.Operation, config.Operation)
				}
				if config.FilePath != tc.expected.FilePath {
					t.Errorf("Expected FilePath %s, got %s", tc.expected.FilePath, config.FilePath)
				}
				if config.TextInput != tc.expected.TextInput {
					t.Errorf("Expected TextInput %s, got %s", tc.expected.TextInput, config.TextInput)
				}
				if config.OutputUnit != tc.expected.OutputUnit {
					t.Errorf("Expected OutputUnit %s, got %s", tc.expected.OutputUnit, config.OutputUnit)
				}
			}
		})
	}
}

// TestParseFlagsWithParser_ParseTimeOperationErrorCases はparse-time操作のエラーケーステスト
func TestParseFlagsWithParser_ParseTimeOperationErrorCases(t *testing.T) {
	testCases := []struct {
		name    string
		flags   map[string]string
		wantErr bool
	}{
		{
			name: "両方指定_エラー",
			flags: map[string]string{
				"operation":  "parse-time",
				"file-path":  "/path/to/test.txt",
				"text-input": "テキスト",
			},
			wantErr: true,
		},
		{
			name: "両方未指定_エラー",
			flags: map[string]string{
				"operation": "parse-time",
			},
			wantErr: true,
		},
		{
			name: "無効なファイル拡張子_エラー",
			flags: map[string]string{
				"operation": "parse-time",
				"file-path": "/path/to/test.pdf",
			},
			wantErr: true,
		},
		{
			name: "無効な出力単位_エラー",
			flags: map[string]string{
				"operation":   "parse-time",
				"text-input":  "テキスト",
				"output-unit": "invalid",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockParser := NewMockFlagParser()
			for key, value := range tc.flags {
				mockParser.SetStringFlag(key, value)
			}

			// Act
			config, err := ParseFlagsWithParser(mockParser)

			// Assert
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseFlagsWithParser() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if tc.wantErr {
				if config != nil {
					t.Error("Expected nil config for error case")
				}
			}
		})
	}
}

// TestNewConfigForParseTime_EdgeCases はエッジケースのテスト
func TestNewConfigForParseTime_EdgeCases(t *testing.T) {
	testCases := []struct {
		name       string
		operation  string
		filePath   string
		textInput  string
		outputUnit string
		wantErr    bool
	}{
		{
			name:       "長いファイルパス",
			operation:  "parse-time",
			filePath:   "/very/long/path/to/some/deeply/nested/directory/structure/with/many/levels/test.txt",
			textInput:  "",
			outputUnit: "minute",
			wantErr:    false,
		},
		{
			name:       "長いテキスト入力",
			operation:  "parse-time",
			filePath:   "",
			textInput:  "これは非常に長いテキスト入力です。作業は合計30分掛かった。この文章は複数の文を含んでおり、様々な内容が記述されています。別の作業は合計45分掛かった。さらに追加の作業は合計60分掛かった。",
			outputUnit: "hour",
			wantErr:    false,
		},
		{
			name:       "特殊文字を含むファイルパス",
			operation:  "parse-time",
			filePath:   "/path/to/test-file_with.special@chars.txt",
			textInput:  "",
			outputUnit: "minute",
			wantErr:    false,
		},
		{
			name:       "日本語を含むファイルパス",
			operation:  "parse-time",
			filePath:   "/path/to/テストファイル.txt",
			textInput:  "",
			outputUnit: "minute",
			wantErr:    false,
		},
		{
			name:       "空白を含むテキスト入力",
			operation:  "parse-time",
			filePath:   "",
			textInput:  "   作業は合計30分掛かった。   ",
			outputUnit: "minute",
			wantErr:    false,
		},
		{
			name:       "改行を含むテキスト入力",
			operation:  "parse-time",
			filePath:   "",
			textInput:  "第1段階の作業は合計30分掛かった。\n第2段階の作業は合計45分掛かった。",
			outputUnit: "minute",
			wantErr:    false,
		},
		{
			name:       "タブを含むテキスト入力",
			operation:  "parse-time",
			filePath:   "",
			textInput:  "作業A\t合計30分掛かった。\n作業B\t合計45分掛かった。",
			outputUnit: "minute",
			wantErr:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			config, err := NewConfigForParseTime(tc.operation, tc.filePath, tc.textInput, tc.outputUnit)

			// Assert
			if (err != nil) != tc.wantErr {
				t.Errorf("NewConfigForParseTime() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				if config == nil {
					t.Fatal("Expected non-nil config")
				}
				if config.Operation != tc.operation {
					t.Errorf("Expected operation %s, got %s", tc.operation, config.Operation)
				}
				if config.FilePath != tc.filePath {
					t.Errorf("Expected FilePath %s, got %s", tc.filePath, config.FilePath)
				}
				if config.TextInput != tc.textInput {
					t.Errorf("Expected TextInput %s, got %s", tc.textInput, config.TextInput)
				}

				expectedOutputUnit := tc.outputUnit
				if expectedOutputUnit == "" {
					expectedOutputUnit = "minute"
				}
				if config.OutputUnit != expectedOutputUnit {
					t.Errorf("Expected OutputUnit %s, got %s", expectedOutputUnit, config.OutputUnit)
				}
			}
		})
	}
}
