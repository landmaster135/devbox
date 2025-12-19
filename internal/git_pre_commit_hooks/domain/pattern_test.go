package domain

import (
	"testing"
)

// TestGetSuspiciousPatterns_Normal は疑わしいパターンの取得テスト
func TestGetSuspiciousPatterns_Normal(t *testing.T) {
	patterns := GetSuspiciousPatterns()

	if len(patterns) == 0 {
		t.Error("GetSuspiciousPatterns() returned empty slice, expected non-empty")
	}
}

// TestGetAllowedPlaceholders_Normal は許可されたプレースホルダーの取得テスト
func TestGetAllowedPlaceholders_Normal(t *testing.T) {
	placeholders := GetAllowedPlaceholders()

	if len(placeholders) == 0 {
		t.Error("GetAllowedPlaceholders() returned empty slice, expected non-empty")
	}
}

// TestGetRealSecretPatterns_Normal は実際のシークレットパターンの取得テスト
func TestGetRealSecretPatterns_Normal(t *testing.T) {
	patterns := GetRealSecretPatterns()

	if len(patterns) == 0 {
		t.Error("GetRealSecretPatterns() returned empty slice, expected non-empty")
	}
}

// TestGetConfigFilePatterns_Normal は設定ファイルパターンの取得テスト
func TestGetConfigFilePatterns_Normal(t *testing.T) {
	patterns := GetConfigFilePatterns()

	if len(patterns) == 0 {
		t.Error("GetConfigFilePatterns() returned empty slice, expected non-empty")
	}
}

// TestGetTestPatterns_Normal はテストパターンの取得テスト
func TestGetTestPatterns_Normal(t *testing.T) {
	patterns := GetTestPatterns()

	if len(patterns) == 0 {
		t.Error("GetTestPatterns() returned empty slice, expected non-empty")
	}
}

// TestGetProtocolPrefixes_Normal はプロトコルプレフィックスの取得テスト
func TestGetProtocolPrefixes_Normal(t *testing.T) {
	prefixes := GetProtocolPrefixes()

	if len(prefixes) == 0 {
		t.Error("GetProtocolPrefixes() returned empty slice, expected non-empty")
	}
}

// TestGetHomePathPattern_Normal はホームパスパターンの取得テスト
func TestGetHomePathPattern_Normal(t *testing.T) {
	pattern := GetHomePathPattern()

	if len(pattern) == 0 {
		t.Error("GetHomePathPattern() returned empty string, expected non-empty")
	}
}

// TestGetAllowedHomePathPatterns_Normal は許可されたホームパスパターンの取得テスト
func TestGetAllowedHomePathPatterns_Normal(t *testing.T) {
	patterns := GetAllowedHomePathPatterns()

	if len(patterns) == 0 {
		t.Error("GetAllowedHomePathPatterns() returned empty slice, expected non-empty")
	}
}

// TestGetBinaryFileExtensions_Normal はバイナリファイル拡張子の取得テスト
func TestGetBinaryFileExtensions_Normal(t *testing.T) {
	extensions := GetBinaryFileExtensions()

	if len(extensions) == 0 {
		t.Error("GetBinaryFileExtensions() returned empty slice, expected non-empty")
	}
}

// TestColorConstants_Normal は色定数のテスト
func TestColorConstants_Normal(t *testing.T) {
	// 色定数が空でないことを確認
	if Red == "" {
		t.Error("Red constant is empty")
	}
	if Green == "" {
		t.Error("Green constant is empty")
	}
	if Yellow == "" {
		t.Error("Yellow constant is empty")
	}
	if Blue == "" {
		t.Error("Blue constant is empty")
	}
	if Reset == "" {
		t.Error("Reset constant is empty")
	}
}
