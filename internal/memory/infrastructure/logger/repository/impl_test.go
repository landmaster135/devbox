package logger

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

// captureOutput は標準出力をキャプチャするヘルパー関数
func captureOutput(f func()) string {
	// 元の標準出力を保存
	oldOutput := os.Stdout

	// パイプを作成
	r, w, _ := os.Pipe()
	os.Stdout = w

	// logパッケージの出力先を変更
	oldLogger := log.Writer()
	log.SetOutput(w)

	// テスト対象の関数を実行
	f()

	// パイプを閉じて出力を読み取る
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// 元の設定に戻す
	os.Stdout = oldOutput
	log.SetOutput(oldLogger)

	return buf.String()
}

// TestDebug はDebugメソッドのテスト
func TestDebug(t *testing.T) {
	logger := NewDefaultLogger()

	// 基本的なメッセージのみのケース
	output := captureOutput(func() {
		logger.Debug("テストメッセージ")
	})

	if !strings.Contains(output, "[DEBUG] テストメッセージ") {
		t.Errorf("期待されるデバッグメッセージが出力されていません。出力: %s", output)
	}

	// キーと値のペアを含むケース
	output = captureOutput(func() {
		logger.Debug("テストメッセージ", "key1", "value1", "key2", 123)
	})

	if !strings.Contains(output, "[DEBUG] テストメッセージ") ||
		!strings.Contains(output, "key1=value1") ||
		!strings.Contains(output, "key2=123") {
		t.Errorf("キーと値のペアが正しく出力されていません。出力: %s", output)
	}

	// 奇数個のキーと値のペア（値がないキーがある）のケース
	output = captureOutput(func() {
		logger.Debug("テストメッセージ", "key1", "value1", "key2")
	})

	if !strings.Contains(output, "[DEBUG] テストメッセージ") ||
		!strings.Contains(output, "key1=value1") ||
		!strings.Contains(output, "key2=<no value>") {
		t.Errorf("値がないキーが正しく処理されていません。出力: %s", output)
	}
}

// TestInfo はInfoメソッドのテスト
func TestInfo(t *testing.T) {
	logger := NewDefaultLogger()

	// 基本的なメッセージのみのケース
	output := captureOutput(func() {
		logger.Info("テストメッセージ")
	})

	if !strings.Contains(output, "[INFO] テストメッセージ") {
		t.Errorf("期待される情報メッセージが出力されていません。出力: %s", output)
	}

	// キーと値のペアを含むケース
	output = captureOutput(func() {
		logger.Info("テストメッセージ", "key1", "value1", "key2", 123)
	})

	if !strings.Contains(output, "[INFO] テストメッセージ") ||
		!strings.Contains(output, "key1=value1") ||
		!strings.Contains(output, "key2=123") {
		t.Errorf("キーと値のペアが正しく出力されていません。出力: %s", output)
	}
}

// TestWarn はWarnメソッドのテスト
func TestWarn(t *testing.T) {
	logger := NewDefaultLogger()

	// 基本的なメッセージのみのケース
	output := captureOutput(func() {
		logger.Warn("テストメッセージ")
	})

	if !strings.Contains(output, "[WARN] テストメッセージ") {
		t.Errorf("期待される警告メッセージが出力されていません。出力: %s", output)
	}

	// キーと値のペアを含むケース
	output = captureOutput(func() {
		logger.Warn("テストメッセージ", "key1", "value1", "key2", 123)
	})

	if !strings.Contains(output, "[WARN] テストメッセージ") ||
		!strings.Contains(output, "key1=value1") ||
		!strings.Contains(output, "key2=123") {
		t.Errorf("キーと値のペアが正しく出力されていません。出力: %s", output)
	}
}

// TestError はErrorメソッドのテスト
func TestError(t *testing.T) {
	logger := NewDefaultLogger()
	testErr := errors.New("テストエラー")

	// 基本的なメッセージとエラーのケース
	output := captureOutput(func() {
		logger.Error("エラーが発生しました", testErr)
	})

	if !strings.Contains(output, "[ERROR] エラーが発生しました") ||
		!strings.Contains(output, "error=テストエラー") {
		t.Errorf("期待されるエラーメッセージが出力されていません。出力: %s", output)
	}

	// キーと値のペアを含むケース
	output = captureOutput(func() {
		logger.Error("エラーが発生しました", testErr, "key1", "value1", "key2", 123)
	})

	if !strings.Contains(output, "[ERROR] エラーが発生しました") ||
		!strings.Contains(output, "error=テストエラー") ||
		!strings.Contains(output, "key1=value1") ||
		!strings.Contains(output, "key2=123") {
		t.Errorf("キーと値のペアが正しく出力されていません。出力: %s", output)
	}

	// nilエラーのケース
	output = captureOutput(func() {
		logger.Error("エラーが発生しました", nil, "key1", "value1")
	})

	if !strings.Contains(output, "[ERROR] エラーが発生しました") ||
		!strings.Contains(output, "key1=value1") {
		t.Errorf("nilエラーが正しく処理されていません。出力: %s", output)
	}
}

// TestLog は内部メソッドlogのテスト
func TestLog(t *testing.T) {
	logger := &DefaultLogger{}
	testErr := errors.New("テストエラー")

	// 複雑なケース：様々なタイプの値を含む
	output := captureOutput(func() {
		logger.log("TEST", "複雑なテスト", testErr,
			"string", "文字列",
			"int", 42,
			"bool", true,
			"float", 3.14,
			"nil", nil)
	})

	expectedContents := []string{
		"[TEST] 複雑なテスト",
		"string=文字列",
		"int=42",
		"bool=true",
		"float=3.14",
		"nil=<nil>",
		"error=テストエラー",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(output, expected) {
			t.Errorf("期待される内容 '%s' が出力に含まれていません。出力: %s", expected, output)
		}
	}
}

// TestInit はInit関数のテスト
func TestInit(t *testing.T) {
	logger := NewDefaultLogger()

	tests := []struct {
		name        string
		level       int
		format      string
		expectError bool
	}{
		{
			name:        "正常なレベル0とテキストフォーマット",
			level:       0,
			format:      "text",
			expectError: false,
		},
		{
			name:        "正常なレベル1とJSONフォーマット",
			level:       1,
			format:      "json",
			expectError: false,
		},
		{
			name:        "正常なレベル2",
			level:       2,
			format:      "text",
			expectError: false,
		},
		{
			name:        "正常なレベル3",
			level:       3,
			format:      "json",
			expectError: false,
		},
		{
			name:        "無効なレベル（負の値）",
			level:       -1,
			format:      "text",
			expectError: true,
		},
		{
			name:        "無効なレベル（範囲外）",
			level:       4,
			format:      "text",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := logger.Init(tt.level, tt.format)

			if tt.expectError {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("エラーが発生しました: %v", err)
				}
			}
		})
	}
}

// TestNewDefaultLogger はNewDefaultLogger関数のテスト
func TestNewDefaultLogger(t *testing.T) {
	logger := NewDefaultLogger()

	if logger == nil {
		t.Error("NewDefaultLogger関数がnilを返しました")
	}

	// 返されたロガーがLogger interface を実装していることを確認
	var _ Logger = logger
}

// TestEdgeCases はエッジケースのテスト
func TestEdgeCases(t *testing.T) {
	logger := NewDefaultLogger()

	t.Run("空文字列メッセージ", func(t *testing.T) {
		output := captureOutput(func() {
			logger.Info("")
		})

		if !strings.Contains(output, "[INFO]") {
			t.Errorf("空文字列メッセージが正しく処理されていません。出力: %s", output)
		}
	})

	t.Run("特殊文字を含むメッセージ", func(t *testing.T) {
		specialMsg := "特殊文字: !@#$%^&*()_+-=[]{}|;':\",./<>?"
		output := captureOutput(func() {
			logger.Info(specialMsg)
		})

		if !strings.Contains(output, "[INFO]") || !strings.Contains(output, "特殊文字") {
			t.Errorf("特殊文字を含むメッセージが正しく処理されていません。出力: %s", output)
		}
	})

	t.Run("非常に長いメッセージ", func(t *testing.T) {
		longMsg := strings.Repeat("長いメッセージ", 100)
		output := captureOutput(func() {
			logger.Info(longMsg)
		})

		if !strings.Contains(output, "[INFO]") {
			t.Errorf("長いメッセージが正しく処理されていません。出力: %s", output)
		}
	})

	t.Run("大量のキーと値のペア", func(t *testing.T) {
		var args []any
		for i := 0; i < 50; i++ {
			args = append(args, fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
		}

		output := captureOutput(func() {
			logger.Info("大量のキーと値", args...)
		})

		if !strings.Contains(output, "[INFO]") || !strings.Contains(output, "key0=value0") {
			t.Errorf("大量のキーと値のペアが正しく処理されていません。出力: %s", output)
		}
	})
}

// TestConcurrency は並行性のテスト
func TestConcurrency(t *testing.T) {
	logger := NewDefaultLogger()
	const numGoroutines = 10
	const numMessages = 100

	// 並行してログを出力
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < numMessages; j++ {
				logger.Info(fmt.Sprintf("Goroutine %d - Message %d", id, j))
			}
			done <- true
		}(i)
	}

	// すべてのgoroutineの完了を待つ
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// パニックが発生しなかったことを確認（テストが完了すれば成功）
}
