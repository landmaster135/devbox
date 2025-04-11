package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// テスト用のYAMLファイルを作成する関数
func createTestYamlFile(t *testing.T) string {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	yamlPath := filepath.Join(tempDir, "test_env.yml")

	// テスト用のYAML内容
	yamlContent := `
TEST_KEY1: test_value1
TEST_KEY2: test_value2
`
	// ファイルに書き込む
	err := os.WriteFile(yamlPath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("テスト用YAMLファイルの作成に失敗しました: %v", err)
	}

	return yamlPath
}

func TestRun_Success(t *testing.T) {
	// テスト用のYAMLファイルを作成
	yamlPath := createTestYamlFile(t)

	// 標準出力と標準エラー出力をキャプチャするためのバッファを作成
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// 引数を設定
	args := []string{"-env", yamlPath}

	// run関数を実行
	code := run(args, stdout, stderr)

	// 終了コードが正常であることを確認
	if code != exitCodeOK {
		t.Errorf("終了コードが正しくありません: expected=%d, got=%d", exitCodeOK, code)
	}

	// 標準出力に成功メッセージが含まれていることを確認
	if !strings.Contains(stdout.String(), "環境変数を正常に読み込みました") {
		t.Errorf("標準出力に成功メッセージが含まれていません: %s", stdout.String())
	}

	// 標準エラー出力が空であることを確認
	if stderr.String() != "" {
		t.Errorf("標準エラー出力が空ではありません: %s", stderr.String())
	}
}

func TestRun_InvalidFlag(t *testing.T) {
	// 標準出力と標準エラー出力をキャプチャするためのバッファを作成
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// 無効な引数を設定
	args := []string{"-invalid"}

	// run関数を実行
	code := run(args, stdout, stderr)

	// 終了コードがエラーであることを確認
	if code != exitCodeError {
		t.Errorf("終了コードが正しくありません: expected=%d, got=%d", exitCodeError, code)
	}

	// 標準エラー出力にエラーメッセージが含まれていることを確認
	if stderr.String() == "" {
		t.Errorf("標準エラー出力にエラーメッセージが含まれていません")
	}
}

func TestRun_DefaultEnvFile(t *testing.T) {
	// テスト用のYAMLファイルを作成
	yamlPath := createTestYamlFile(t)

	// 現在のディレクトリを保存
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("現在のディレクトリの取得に失敗しました: %v", err)
	}

	// テスト用ディレクトリに移動
	dir := filepath.Dir(yamlPath)
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("ディレクトリの変更に失敗しました: %v", err)
	}

	// テスト終了後に元のディレクトリに戻る
	defer os.Chdir(currentDir)

	// env.ymlにリネーム
	newPath := filepath.Join(dir, "env.yml")
	if err := os.Rename(yamlPath, newPath); err != nil {
		t.Fatalf("ファイルのリネームに失敗しました: %v", err)
	}

	// 標準出力と標準エラー出力をキャプチャするためのバッファを作成
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// 引数を設定（環境変数ファイルを指定しない）
	args := []string{}

	// run関数を実行
	code := run(args, stdout, stderr)

	// 終了コードが正常であることを確認
	if code != exitCodeOK {
		t.Errorf("終了コードが正しくありません: expected=%d, got=%d", exitCodeOK, code)
	}

	// 標準出力に成功メッセージが含まれていることを確認
	if !strings.Contains(stdout.String(), "環境変数を正常に読み込みました") {
		t.Errorf("標準出力に成功メッセージが含まれていません: %s", stdout.String())
	}

	// 標準エラー出力が空であることを確認
	if stderr.String() != "" {
		t.Errorf("標準エラー出力が空ではありません: %s", stderr.String())
	}
}
