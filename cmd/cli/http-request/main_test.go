package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun_MissingURL(t *testing.T) {
	// 標準出力と標準エラー出力をキャプチャするためのバッファを作成
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	// URLを指定せずにrun関数を呼び出す
	code := run([]string{}, stdout, stderr)

	// 検証
	assert.Equal(t, exitCodeError, code)
	assert.Contains(t, stderr.String(), "エラー: URLを指定してください")
}

func TestRun_MissingJSONFileForPOST(t *testing.T) {
	// 標準出力と標準エラー出力をキャプチャするためのバッファを作成
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	// POSTメソッドを指定するが、JSONファイルを指定しない
	code := run([]string{"-url", "https://example.com/api", "-method", "POST"}, stdout, stderr)

	// 検証
	assert.Equal(t, exitCodeError, code)
	assert.Contains(t, stderr.String(), "エラー: POSTメソッドにはJSONファイルが必要です")
}

func TestRun_InvalidFlag(t *testing.T) {
	// 標準出力と標準エラー出力をキャプチャするためのバッファを作成
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	// 無効なフラグを指定
	code := run([]string{"-invalid"}, stdout, stderr)

	// 検証
	assert.Equal(t, exitCodeError, code)
	assert.Contains(t, stderr.String(), "flag provided but not defined")
}
