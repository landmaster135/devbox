package main

import (
	"bytes"
	"strings"
	"testing"
)

// TestRun はrun関数のテストを行います
func TestRun(t *testing.T) {
	// 出力先のバッファを用意
	buf := &bytes.Buffer{}

	// run関数を実行
	code := run(buf)

	// 終了コードの確認
	if code != 0 {
		t.Errorf("期待する終了コード: 0, 取得した終了コード: %d", code)
	}

	// 出力内容の確認
	output := buf.String()
	if !strings.Contains(output, "Devbox") {
		t.Errorf("出力にDevboxという文字列が含まれていません。出力: %s", output)
	}

	if !strings.Contains(output, "file-processor") {
		t.Errorf("出力にfile-processorという文字列が含まれていません。出力: %s", output)
	}
}
