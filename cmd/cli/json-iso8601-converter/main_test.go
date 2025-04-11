package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "iso8601-test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用のJSONファイルを作成
	jsonContent := `{
		"id": 123,
		"title": "テストデータ",
		"created_at": "2025-04-10T15:30:45Z",
		"metadata": {
			"published_at": "2025-04-10T10:15:20Z"
		}
	}`
	jsonPath := filepath.Join(tempDir, "test.json")
	if err := os.WriteFile(jsonPath, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	tests := []struct {
		name     string
		args     []string
		wantCode exitCode
		wantOut  string
		wantErr  string
	}{
		{
			name:     "正常系: created_atキーを変換",
			args:     []string{"-dir", tempDir, "-key", "created_at"},
			wantCode: exitCodeOK,
			wantOut:  "1 個のJSONファイルのキー 'created_at' の値をISO8601形式からUNIXタイムスタンプに変換しました",
			wantErr:  "",
		},
		{
			name:     "正常系: ドライラン",
			args:     []string{"-dir", tempDir, "-key", "published_at", "-dry-run"},
			wantCode: exitCodeOK,
			wantOut:  "ドライラン: 1 個のJSONファイルで変換対象のキー 'published_at' が見つかりました",
			wantErr:  "",
		},
		{
			name:     "エラー: キーが指定されていない",
			args:     []string{"-dir", tempDir},
			wantCode: exitCodeError,
			wantOut:  "",
			wantErr:  "エラー: 変換対象のキー名を指定してください（-key オプション）",
		},
		{
			name:     "エラー: 存在しないディレクトリ",
			args:     []string{"-dir", filepath.Join(tempDir, "not_exist"), "-key", "created_at"},
			wantCode: exitCodeError,
			wantOut:  "",
			wantErr:  "エラー:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			code := run(tt.args, stdout, stderr)

			if code != tt.wantCode {
				t.Errorf("run() = %v, want %v", code, tt.wantCode)
			}

			if tt.wantOut != "" && !bytes.Contains(stdout.Bytes(), []byte(tt.wantOut)) {
				t.Errorf("stdout = %v, want to contain %v", stdout.String(), tt.wantOut)
			}

			if tt.wantErr != "" && !bytes.Contains(stderr.Bytes(), []byte(tt.wantErr)) {
				t.Errorf("stderr = %v, want to contain %v", stderr.String(), tt.wantErr)
			}
		})
	}
}
