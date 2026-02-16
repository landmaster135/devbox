package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRun_NoArgs は引数なしで実行した場合のテスト
func TestRun_NoArgs(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := run([]string{}, stdout, stderr)

	if code != exitCodeError {
		t.Errorf("期待する終了コード %d, 取得した終了コード %d", exitCodeError, code)
	}

	if !strings.Contains(stderr.String(), "ファイルパスを指定してください") {
		t.Errorf("エラーメッセージが期待と異なります。取得したメッセージ: %s", stderr.String())
	}
}

// TestRun_FilePathOnly はファイルパスのみ指定した場合のテスト
func TestRun_FilePathOnly(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := run([]string{"-file", "test.txt"}, stdout, stderr)

	if code != exitCodeError {
		t.Errorf("期待する終了コード %d, 取得した終了コード %d", exitCodeError, code)
	}

	if !strings.Contains(stderr.String(), "開始位置と終了位置を指定してください") {
		t.Errorf("エラーメッセージが期待と異なります。取得したメッセージ: %s", stderr.String())
	}
}

// TestRun_InvalidFile は存在しないファイルを指定した場合のテスト
func TestRun_InvalidFile(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := run([]string{"-file", "存在しないファイル.txt", "-start", "1", "-end", "5"}, stdout, stderr)

	if code != exitCodeError {
		t.Errorf("期待する終了コード %d, 取得した終了コード %d", exitCodeError, code)
	}

	if !strings.Contains(stderr.String(), "エラー:") {
		t.Errorf("エラーメッセージが期待と異なります。取得したメッセージ: %s", stderr.String())
	}
}

// TestRun_ValidFile は有効なファイルと引数を指定した場合のテスト
func TestRun_ValidFile(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "file-line-deduper-test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用のファイルを作成
	testFilePath := filepath.Join(tempDir, "test.txt")
	testContent := []string{
		"abc123def",
		"abc456def",
		"abc123def",
		"xyz789uvw",
	}

	file, err := os.Create(testFilePath)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	for _, line := range testContent {
		file.WriteString(line + "\n")
	}
	file.Close()

	// テスト実行
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// 文字位置3-6を指定（"123", "456", "123", "789"の部分）
	code := run([]string{"-file", testFilePath, "-start", "3", "-end", "6"}, stdout, stderr)

	if code != exitCodeOK {
		t.Errorf("期待する終了コード %d, 取得した終了コード %d", exitCodeOK, code)
	}

	if !strings.Contains(stdout.String(), "処理完了:") {
		t.Errorf("成功メッセージが期待と異なります。取得したメッセージ: %s", stdout.String())
	}

	// 処理後のファイル内容を確認
	content, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("処理後のファイル読み込みに失敗しました: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) != 3 {
		t.Errorf("処理後のファイル行数が期待と異なります。期待: 3, 取得: %d", len(lines))
	}

	// "abc123def"は重複しているため、1つだけ残るはず
	count123 := 0
	for _, line := range lines {
		if strings.Contains(line, "123") {
			count123++
		}
	}

	if count123 != 1 {
		t.Errorf("重複行の削除が正しく行われていません。'123'を含む行数: %d, 期待: 1", count123)
	}
}

// TestRun_InvalidArgs は無効な引数を指定した場合のテスト
func TestRun_InvalidArgs(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// 無効なフラグを指定
	code := run([]string{"-invalid"}, stdout, stderr)

	if code != exitCodeError {
		t.Errorf("期待する終了コード %d, 取得した終了コード %d", exitCodeError, code)
	}

	if stderr.String() == "" {
		t.Error("エラーメッセージが出力されていません")
	}
}
