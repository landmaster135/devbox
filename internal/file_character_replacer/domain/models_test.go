package domain

import (
	"fmt"
	"testing"
)

// TestEncodingType_IsValid_Normal はEncodingType.IsValid()の正常系をテストします
func TestEncodingType_IsValid_Normal(t *testing.T) {
	tests := []struct {
		name     string
		encoding EncodingType
		expected bool
	}{
		{
			name:     "UTF-8は有効",
			encoding: EncodingUTF8,
			expected: true,
		},
		{
			name:     "Shift_JISは有効",
			encoding: EncodingShiftJIS,
			expected: true,
		},
		{
			name:     "EUC-JPは有効",
			encoding: EncodingEUCJP,
			expected: true,
		},
		{
			name:     "ISO-2022-JPは有効",
			encoding: EncodingISO2022JP,
			expected: true,
		},
		{
			name:     "無効なエンコーディング",
			encoding: EncodingType("invalid"),
			expected: false,
		},
		{
			name:     "空文字列は無効",
			encoding: EncodingType(""),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.encoding.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestReplacementConfig_Validate_Normal はReplacementConfig.Validate()の正常系をテストします
func TestReplacementConfig_Validate_Normal(t *testing.T) {
	config := &ReplacementConfig{
		Target:   "/test/path",
		From:     "old",
		To:       "new",
		Encoding: EncodingUTF8,
	}

	err := config.Validate()
	if err != nil {
		t.Errorf("Validate() returned unexpected error: %v", err)
	}
}

// TestReplacementConfig_Validate_Error はReplacementConfig.Validate()のエラーケースをテストします
func TestReplacementConfig_Validate_Error(t *testing.T) {
	tests := []struct {
		name   string
		config *ReplacementConfig
	}{
		{
			name: "Targetが空",
			config: &ReplacementConfig{
				Target:   "",
				From:     "old",
				To:       "new",
				Encoding: EncodingUTF8,
			},
		},
		{
			name: "Fromが空",
			config: &ReplacementConfig{
				Target:   "/test/path",
				From:     "",
				To:       "new",
				Encoding: EncodingUTF8,
			},
		},
		{
			name: "Toが空",
			config: &ReplacementConfig{
				Target:   "/test/path",
				From:     "old",
				To:       "",
				Encoding: EncodingUTF8,
			},
		},
		{
			name: "無効なエンコーディング",
			config: &ReplacementConfig{
				Target:   "/test/path",
				From:     "old",
				To:       "new",
				Encoding: EncodingType("invalid"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if err == nil {
				t.Error("Validate() should return error but got nil")
			}
		})
	}
}

// TestReplacementConfig_ValidateWithFileSystem_Normal はValidateWithFileSystem()の正常系をテストします
func TestReplacementConfig_ValidateWithFileSystem_Normal(t *testing.T) {
	tests := []struct {
		name        string
		config      *ReplacementConfig
		isDirectory bool
	}{
		{
			name: "ファイル処理（バックアップなし）",
			config: &ReplacementConfig{
				Target:   "/test/file.txt",
				From:     "old",
				To:       "new",
				Encoding: EncodingUTF8,
				Backup:   false,
			},
			isDirectory: false,
		},
		{
			name: "ファイル処理（バックアップあり、BackupDir指定なし）",
			config: &ReplacementConfig{
				Target:   "/test/file.txt",
				From:     "old",
				To:       "new",
				Encoding: EncodingUTF8,
				Backup:   true,
			},
			isDirectory: false,
		},
		{
			name: "ディレクトリ処理（バックアップなし）",
			config: &ReplacementConfig{
				Target:   "/test/dir",
				From:     "old",
				To:       "new",
				Encoding: EncodingUTF8,
				Backup:   false,
			},
			isDirectory: true,
		},
		{
			name: "ディレクトリ処理（バックアップあり、BackupDir指定あり）",
			config: &ReplacementConfig{
				Target:    "/test/dir",
				From:      "old",
				To:        "new",
				Encoding:  EncodingUTF8,
				Backup:    true,
				BackupDir: "/backup",
			},
			isDirectory: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.ValidateWithFileSystem(tt.isDirectory)
			if err != nil {
				t.Errorf("ValidateWithFileSystem() returned unexpected error: %v", err)
			}
		})
	}
}

// TestReplacementConfig_ValidateWithFileSystem_Error はValidateWithFileSystem()のエラーケースをテストします
func TestReplacementConfig_ValidateWithFileSystem_Error(t *testing.T) {
	tests := []struct {
		name        string
		config      *ReplacementConfig
		isDirectory bool
	}{
		{
			name: "基本バリデーションエラー（Targetが空）",
			config: &ReplacementConfig{
				Target:   "",
				From:     "old",
				To:       "new",
				Encoding: EncodingUTF8,
			},
			isDirectory: false,
		},
		{
			name: "ディレクトリ処理でバックアップディレクトリ未指定",
			config: &ReplacementConfig{
				Target:    "/test/dir",
				From:      "old",
				To:        "new",
				Encoding:  EncodingUTF8,
				Backup:    true,
				BackupDir: "",
			},
			isDirectory: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.ValidateWithFileSystem(tt.isDirectory)
			if err == nil {
				t.Error("ValidateWithFileSystem() should return error but got nil")
			}
		})
	}
}

// TestFileProcessResult_AddError_Normal はFileProcessResult.AddError()をテストします
func TestFileProcessResult_AddError_Normal(t *testing.T) {
	result := &FileProcessResult{}

	if len(result.Errors) != 0 {
		t.Errorf("初期状態でエラー数が0でない: %d", len(result.Errors))
	}

	testError := fmt.Errorf("テストエラー")
	result.AddError(testError)

	if len(result.Errors) != 1 {
		t.Errorf("AddError後のエラー数が1でない: %d", len(result.Errors))
	}

	if result.Errors[0] != testError {
		t.Errorf("追加されたエラーが期待値と異なる: %v", result.Errors[0])
	}
}

// TestFileProcessResult_AddMessage_Normal はFileProcessResult.AddMessage()をテストします
func TestFileProcessResult_AddMessage_Normal(t *testing.T) {
	result := &FileProcessResult{}

	if len(result.Messages) != 0 {
		t.Errorf("初期状態でメッセージ数が0でない: %d", len(result.Messages))
	}

	testMessage := "テストメッセージ"
	result.AddMessage(testMessage)

	if len(result.Messages) != 1 {
		t.Errorf("AddMessage後のメッセージ数が1でない: %d", len(result.Messages))
	}

	if result.Messages[0] != testMessage {
		t.Errorf("追加されたメッセージが期待値と異なる: %v", result.Messages[0])
	}
}

// TestFileProcessResult_HasErrors_Normal はFileProcessResult.HasErrors()をテストします
func TestFileProcessResult_HasErrors_Normal(t *testing.T) {
	result := &FileProcessResult{}

	// 初期状態ではエラーなし
	if result.HasErrors() {
		t.Error("初期状態でHasErrors()がtrueを返した")
	}

	// エラーを追加
	result.AddError(fmt.Errorf("テストエラー"))

	// エラーありの状態
	if !result.HasErrors() {
		t.Error("エラー追加後にHasErrors()がfalseを返した")
	}
}
