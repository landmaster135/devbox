package security

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPathValidator_normalizeAndValidateEncoding_Normal(t *testing.T) {
	validator := NewDefaultPathValidator()

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "通常のパス",
			path:     "/home/user/project",
			expected: "/home/user/project",
		},
		{
			name:     "日本語を含むパス",
			path:     "/home" + "/ユーザー/プロジェクト",
			expected: "/home" + "/ユーザー/プロジェクト",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.normalizeAndValidateEncoding(tt.path)
			if err != nil {
				t.Errorf("normalizeAndValidateEncoding() error = %v, want nil", err)
			}
			if result != tt.expected {
				t.Errorf("normalizeAndValidateEncoding() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPathValidator_normalizeAndValidateEncoding_InvalidUTF8(t *testing.T) {
	validator := NewDefaultPathValidator()

	// 無効なUTF-8バイト列
	invalidUTF8 := string([]byte{0xff, 0xfe, 0xfd})
	_, err := validator.normalizeAndValidateEncoding(invalidUTF8)
	if err == nil {
		t.Error("normalizeAndValidateEncoding() error = nil, want error for invalid UTF-8")
	}

	if !IsSecurityError(err) {
		t.Error("Expected SecurityError for invalid UTF-8")
	}

	if GetSecurityErrorType(err) != ErrorTypeEncodingAttack {
		t.Errorf("Expected ErrorTypeEncodingAttack, got %v", GetSecurityErrorType(err))
	}
}

func TestPathValidator_normalizeAndValidateEncoding_ControlCharacters(t *testing.T) {
	validator := NewDefaultPathValidator()

	controlChars := []rune{
		0x00, // NULL
		0x01, // SOH
		0x02, // STX
		0x1F, // US
		0x7F, // DEL
	}

	for _, char := range controlChars {
		t.Run("control_char_"+string(char), func(t *testing.T) {
			path := "/path/with" + string(char) + "control"
			_, err := validator.normalizeAndValidateEncoding(path)
			if err == nil {
				t.Errorf("normalizeAndValidateEncoding() error = nil, want error for control character U+%04X", char)
			}
			if GetSecurityErrorType(err) != ErrorTypeEncodingAttack {
				t.Errorf("Expected ErrorTypeEncodingAttack, got %v", GetSecurityErrorType(err))
			}
		})
	}
}

func TestPathValidator_normalizeAndValidateEncoding_MultiStageDecoding(t *testing.T) {
	validator := NewDefaultPathValidator()

	path := "/tmp/%252e%252e%252fetc/passwd"
	_, err := validator.normalizeAndValidateEncoding(path)
	if err == nil {
		t.Fatal("expected error for multi-stage encoded traversal, got nil")
	}
	if GetSecurityErrorType(err) != ErrorTypeEncodingAttack {
		t.Fatalf("expected ErrorTypeEncodingAttack, got %v", GetSecurityErrorType(err))
	}
}

func TestPathValidator_checkEncodedDangerousPatterns(t *testing.T) {
	validator := NewDefaultPathValidator()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "URLエンコードされたパストラバーサル",
			path:    "/path/%2e%2e/file",
			wantErr: true,
		},
		{
			name:    "大文字のURLエンコード",
			path:    "/path/%2E%2E/file",
			wantErr: true,
		},
		{
			name:    "URLエンコードされたパストラバーサル（スラッシュ付き）",
			path:    "/path/%2e%2e%2f",
			wantErr: true,
		},
		{
			name:    "URLエンコードされたチルダ",
			path:    "/path/%7e/file",
			wantErr: true,
		},
		{
			name:    "URLエンコードされたドル記号",
			path:    "/path/%24HOME/file",
			wantErr: true,
		},
		{
			name:    "URLエンコードされたバッククォート",
			path:    "/path/%60command%60/file",
			wantErr: true,
		},
		{
			name:    "URLエンコードされたセミコロン",
			path:    "/path/%3bcommand/file",
			wantErr: true,
		},
		{
			name:    "URLエンコードされたパイプ",
			path:    "/path/%7ccommand/file",
			wantErr: true,
		},
		{
			name:    "URLエンコードされたアンパサンド",
			path:    "/path/%26command/file",
			wantErr: true,
		},
		{
			name:    "URLエンコードされたリダイレクト",
			path:    "/path/%3efile",
			wantErr: true,
		},
		{
			name:    "URLエンコードされたワイルドカード",
			path:    "/path/%2a/file",
			wantErr: true,
		},
		{
			name:    "URLエンコードされたNULL文字",
			path:    "/path/%00/file",
			wantErr: true,
		},
		{
			name:    "安全なパス",
			path:    "/home/user/project",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.checkEncodedDangerousPatterns(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkEncodedDangerousPatterns() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && GetSecurityErrorType(err) != ErrorTypeEncodingAttack {
				t.Errorf("Expected ErrorTypeEncodingAttack, got %v", GetSecurityErrorType(err))
			}
		})
	}
}

func TestPathValidator_checkEnhancedDangerousPatterns(t *testing.T) {
	validator := NewDefaultPathValidator()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "基本的なパストラバーサル",
			path:    "/path/../file",
			wantErr: true,
		},
		{
			name:    "中間パストラバーサル",
			path:    "/path/../other/file",
			wantErr: true,
		},
		{
			name:    "環境変数展開",
			path:    "/path/${HOME}/file",
			wantErr: true,
		},
		{
			name:    "コマンド置換",
			path:    "/path/$(whoami)/file",
			wantErr: true,
		},
		{
			name:    "コマンドインジェクション（セミコロン）",
			path:    "/path; rm -rf /",
			wantErr: true,
		},
		{
			name:    "コマンドインジェクション（パイプ）",
			path:    "/path | cat /etc/passwd",
			wantErr: true,
		},
		{
			name:    "コマンドインジェクション（アンパサンド）",
			path:    "/path & malicious_command",
			wantErr: true,
		},
		{
			name:    "コマンドインジェクション（ダブルアンパサンド）",
			path:    "/path && malicious_command",
			wantErr: true,
		},
		{
			name:    "コマンドインジェクション（ダブルパイプ）",
			path:    "/path || malicious_command",
			wantErr: true,
		},
		{
			name:    "コマンドインジェクション（パイプアンパサンド）",
			path:    "/path |& logger",
			wantErr: true,
		},
		{
			name:    "リダイレクト攻撃",
			path:    "/path > /etc/passwd",
			wantErr: true,
		},
		{
			name:    "安全なパス",
			path:    "/home/user/project",
			wantErr: false,
		},
		{
			name:    "ドット付きファイル名（安全）",
			path:    "/home/user/.bashrc",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.checkEnhancedDangerousPatterns(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkEnhancedDangerousPatterns() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPathValidator_checkSystemPaths(t *testing.T) {
	validator := NewDefaultPathValidator()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "/procへのアクセス",
			path:    "/proc/self/environ",
			wantErr: true,
		},
		{
			name:    "/sysへのアクセス",
			path:    "/sys/kernel/debug",
			wantErr: true,
		},
		{
			name:    "/devへのアクセス",
			path:    "/dev/null",
			wantErr: true,
		},
		{
			name:    "/etc/passwdへのアクセス",
			path:    "/etc/passwd",
			wantErr: true,
		},
		{
			name:    "/etc/shadowへのアクセス",
			path:    "/etc/shadow",
			wantErr: true,
		},
		{
			name:    "/rootへのアクセス",
			path:    "/root/.ssh/id_rsa",
			wantErr: true,
		},
		{
			name:    "/tmpへのアクセス",
			path:    "/tmp/malicious",
			wantErr: true,
		},
		{
			name:    "安全なパス",
			path:    "/home/user/project",
			wantErr: false,
		},
		{
			name:    "相対パス（安全）",
			path:    "./project",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.checkSystemPaths(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkSystemPaths() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && GetSecurityErrorType(err) != ErrorTypeSystemPath {
				t.Errorf("Expected ErrorTypeSystemPath, got %v", GetSecurityErrorType(err))
			}
		})
	}
}

func TestPathValidator_checkSymlinkSafety(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// 安全なシンボリックリンクを作成
	safeTarget := filepath.Join(tempDir, "safe_target")
	if err := os.Mkdir(safeTarget, 0755); err != nil {
		t.Fatalf("Failed to create safe target directory: %v", err)
	}

	safeLink := filepath.Join(tempDir, "safe_link")
	if err := os.Symlink("safe_target", safeLink); err != nil {
		t.Fatalf("Failed to create safe symlink: %v", err)
	}

	validator := NewPathValidator([]string{tempDir}, 4096)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "安全なシンボリックリンク",
			path:    safeLink,
			wantErr: false,
		},
		{
			name:    "通常のディレクトリ",
			path:    safeTarget,
			wantErr: false,
		},
		{
			name:    "存在しないパス",
			path:    filepath.Join(tempDir, "nonexistent"),
			wantErr: false, // 存在しないパスは後でチェックされる
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.checkSymlinkSafety(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkSymlinkSafety() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// generateUniqueTestDir はハッシュ値を使用してユニークなテストディレクトリ名を生成する
func generateUniqueTestDir(baseDir, prefix string) string {
	// 現在時刻とテスト名からハッシュを生成
	data := fmt.Sprintf("%s_%d_%s", prefix, time.Now().UnixNano(), "test_git_validator")
	hash := sha256.Sum256([]byte(data))
	hashStr := fmt.Sprintf("%x", hash)[:16] // 16文字のハッシュを使用
	return filepath.Join(baseDir, fmt.Sprintf("%s_%s", prefix, hashStr))
}

func TestPathValidator_ValidateWorkingDirectory_Enhanced_Integration(t *testing.T) {
	// ホームディレクトリ配下にテスト用ディレクトリを作成
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	testDir := generateUniqueTestDir(homeDir, "test_git_validator")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// テスト用のGitリポジトリを作成
	gitDir := filepath.Join(testDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	validator := NewPathValidator([]string{homeDir}, 4096)

	tests := []struct {
		name    string
		path    string
		wantErr bool
		errType SecurityErrorType
	}{
		{
			name:    "正常なパス",
			path:    testDir,
			wantErr: false,
		},
		{
			name:    "空のパス",
			path:    "",
			wantErr: true,
			errType: ErrorTypeEmptyPath,
		},
		{
			name:    "長すぎるパス",
			path:    strings.Repeat("a", 5000),
			wantErr: true,
			errType: ErrorTypePathTooLong,
		},
		{
			name:    "パストラバーサル攻撃",
			path:    testDir + "/../../../etc/passwd",
			wantErr: true,
			errType: ErrorTypeDangerousPattern, // ..が先に検出される
		},
		{
			name:    "URLエンコード攻撃",
			path:    testDir + "/%2e%2e/file",
			wantErr: true,
			errType: ErrorTypeEncodingAttack,
		},
		{
			name:    "システムパス攻撃",
			path:    "/proc/self/environ",
			wantErr: true,
			errType: ErrorTypeSymlinkAttack, // 実際はシンボリックリンク攻撃として検出される
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validator.ValidateWorkingDirectory(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateWorkingDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.wantErr {
				if !IsSecurityError(err) {
					t.Error("Expected SecurityError")
				}
				if GetSecurityErrorType(err) != tt.errType {
					t.Errorf("Expected error type %v, got %v", tt.errType, GetSecurityErrorType(err))
				}
			}
		})
	}
}

func TestPathValidator_checkAllowedCharacters(t *testing.T) {
	validator := NewDefaultPathValidator()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "英数字と日本語を含むパス",
			path:    "/home" + "/ユーザー/project",
			wantErr: false,
		},
		{
			name:    "シングルクォートを含むパス",
			path:    "/home/user/'malicious'",
			wantErr: true,
		},
		{
			name:    "全角セミコロンを含むパス",
			path:    "/home/user/；/project",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.checkAllowedCharacters(tt.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("checkAllowedCharacters() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && GetSecurityErrorType(err) != ErrorTypeDangerousPattern {
				t.Fatalf("expected ErrorTypeDangerousPattern, got %v", GetSecurityErrorType(err))
			}
		})
	}
}

func TestPathValidator_checkSymlinkSafety_ParentSymlinkTraversal(t *testing.T) {
	baseDir := t.TempDir()
	outsideDir := t.TempDir()

	symlinkParent := filepath.Join(baseDir, "link")
	if err := os.Symlink(outsideDir, symlinkParent); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	validator := NewPathValidator([]string{baseDir}, 4096)
	attackPath := filepath.Join(symlinkParent, "nested", "dir")

	err := validator.checkSymlinkSafety(attackPath)
	if err == nil {
		t.Fatal("expected symlink safety error, got nil")
	}
	if GetSecurityErrorType(err) != ErrorTypeSymlinkAttack {
		t.Fatalf("expected ErrorTypeSymlinkAttack, got %v", GetSecurityErrorType(err))
	}
}

func TestSecurityError_Methods(t *testing.T) {
	err := NewSecurityError(ErrorTypeDangerousPattern, "/test/path", "テストエラー", "詳細情報")

	// Error()メソッドのテスト
	expectedMsg := "/test/path: テストエラー (詳細: 詳細情報)"
	if err.Error() != expectedMsg {
		t.Errorf("Error() = %v, want %v", err.Error(), expectedMsg)
	}

	// IsSecurityError()のテスト
	if !IsSecurityError(err) {
		t.Error("IsSecurityError() = false, want true")
	}

	// GetSecurityErrorType()のテスト
	if GetSecurityErrorType(err) != ErrorTypeDangerousPattern {
		t.Errorf("GetSecurityErrorType() = %v, want %v", GetSecurityErrorType(err), ErrorTypeDangerousPattern)
	}

	// 詳細なしのエラー
	errNoDetails := NewSecurityError(ErrorTypeEmptyPath, "/test", "エラー", "")
	expectedMsgNoDetails := "/test: エラー"
	if errNoDetails.Error() != expectedMsgNoDetails {
		t.Errorf("Error() = %v, want %v", errNoDetails.Error(), expectedMsgNoDetails)
	}
}

func TestSecurityError_NonSecurityError(t *testing.T) {
	// 通常のエラー
	normalErr := fmt.Errorf("通常のエラー")

	if IsSecurityError(normalErr) {
		t.Error("IsSecurityError() = true, want false for normal error")
	}

	if GetSecurityErrorType(normalErr) != ErrorTypeUnknown {
		t.Errorf("GetSecurityErrorType() = %v, want %v for normal error", GetSecurityErrorType(normalErr), ErrorTypeUnknown)
	}
}
