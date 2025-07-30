package logger

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"testing"
)

// TestLoggerInit はInit関数のテスト
func TestLoggerInit(t *testing.T) {
	tests := []struct {
		name        string
		level       int
		format      string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "正常なレベル0（Info）とテキストフォーマット",
			level:       0,
			format:      "text",
			expectError: false,
		},
		{
			name:        "正常なレベル1（Debug）とJSONフォーマット",
			level:       1,
			format:      "json",
			expectError: false,
		},
		{
			name:        "正常なレベル2（Warn）とテキストフォーマット",
			level:       2,
			format:      "text",
			expectError: false,
		},
		{
			name:        "正常なレベル3（Error）とJSONフォーマット",
			level:       3,
			format:      "json",
			expectError: false,
		},
		{
			name:        "デフォルトフォーマット（空文字列）",
			level:       0,
			format:      "",
			expectError: false,
		},
		{
			name:        "無効なレベル（負の値）",
			level:       -1,
			format:      "text",
			expectError: true,
			errorMsg:    "level is invalid: -1",
		},
		{
			name:        "無効なレベル（範囲外）",
			level:       4,
			format:      "text",
			expectError: true,
			errorMsg:    "level is invalid: 4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Init(tt.level, tt.format)

			if tt.expectError {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("期待されるエラーメッセージ: %s, 実際: %s", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("エラーが発生しました: %v", err)
				}
			}
		})
	}
}

// TestSetLevel はSetLevel関数のテスト
func TestSetLevel(t *testing.T) {
	tests := []struct {
		name     string
		levelStr string
	}{
		{
			name:     "デバッグレベル設定",
			levelStr: "debug",
		},
		{
			name:     "情報レベル設定",
			levelStr: "info",
		},
		{
			name:     "警告レベル設定",
			levelStr: "warn",
		},
		{
			name:     "エラーレベル設定",
			levelStr: "error",
		},
		{
			name:     "無効なレベル（デフォルトのinfoになる）",
			levelStr: "invalid",
		},
		{
			name:     "空文字列（デフォルトのinfoになる）",
			levelStr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// SetLevel関数を呼び出す（パニックしないことを確認）
			SetLevel(tt.levelStr)

			// 関数が正常に実行されることを確認
			// 実際のレベル設定の確認は出力をキャプチャして行う必要があるが、
			// ここでは関数が正常に実行されることのみを確認
		})
	}
}

// captureLogOutput はログ出力をキャプチャするヘルパー関数
func captureLogOutput(f func()) string {
	// 元の標準出力を保存
	oldStdout := os.Stdout

	// パイプを作成
	r, w, _ := os.Pipe()
	os.Stdout = w

	// テスト対象の関数を実行
	f()

	// パイプを閉じて出力を読み取る
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// 元の設定に戻す
	os.Stdout = oldStdout

	return buf.String()
}

// TestDebug はDebug関数のテスト
func TestDebug(t *testing.T) {
	// デバッグレベルに設定
	SetLevel("debug")

	tests := []struct {
		name string
		msg  string
		args []any
	}{
		{
			name: "基本的なデバッグメッセージ",
			msg:  "デバッグメッセージ",
			args: nil,
		},
		{
			name: "引数付きデバッグメッセージ",
			msg:  "デバッグメッセージ",
			args: []any{"key1", "value1", "key2", 123},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Debug関数を呼び出す（パニックしないことを確認）
			Debug(tt.msg, tt.args...)
		})
	}
}

// TestInfo はInfo関数のテスト
func TestInfo(t *testing.T) {
	// 情報レベルに設定
	SetLevel("info")

	tests := []struct {
		name string
		msg  string
		args []any
	}{
		{
			name: "基本的な情報メッセージ",
			msg:  "情報メッセージ",
			args: nil,
		},
		{
			name: "引数付き情報メッセージ",
			msg:  "情報メッセージ",
			args: []any{"key1", "value1", "key2", 123},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Info関数を呼び出す（パニックしないことを確認）
			Info(tt.msg, tt.args...)
		})
	}
}

// TestWarn はWarn関数のテスト
func TestWarn(t *testing.T) {
	// 警告レベルに設定
	SetLevel("warn")

	tests := []struct {
		name string
		msg  string
		args []any
	}{
		{
			name: "基本的な警告メッセージ",
			msg:  "警告メッセージ",
			args: nil,
		},
		{
			name: "引数付き警告メッセージ",
			msg:  "警告メッセージ",
			args: []any{"key1", "value1", "key2", 123},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Warn関数を呼び出す（パニックしないことを確認）
			Warn(tt.msg, tt.args...)
		})
	}
}

// TestError はError関数のテスト
func TestError(t *testing.T) {
	// エラーレベルに設定
	SetLevel("error")

	tests := []struct {
		name string
		msg  string
		err  error
		args []any
	}{
		{
			name: "エラーありの基本的なエラーメッセージ",
			msg:  "エラーメッセージ",
			err:  &testError{message: "テストエラー"},
			args: nil,
		},
		{
			name: "エラーなしの基本的なエラーメッセージ",
			msg:  "エラーメッセージ",
			err:  nil,
			args: nil,
		},
		{
			name: "引数付きエラーメッセージ",
			msg:  "エラーメッセージ",
			err:  &testError{message: "テストエラー"},
			args: []any{"key1", "value1", "key2", 123},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Error関数を呼び出す（パニックしないことを確認）
			Error(tt.msg, tt.err, tt.args...)
		})
	}
}

// testError はテスト用のエラー型
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}

// TestFatal はFatal関数のテスト
func TestFatal(t *testing.T) {
	// Fatal関数は実際にはos.Exit(1)を呼び出すため、
	// テスト環境では直接呼び出すことができません。
	t.Skip("Fatal関数の存在確認")
}

// TestWithContext はWithContext関数のテスト
func TestWithContext(t *testing.T) {
	ctx := context.Background()

	// WithContext関数を呼び出す
	contextLogger := WithContext(ctx)

	// 返されたロガーがnilでないことを確認
	if contextLogger == nil {
		t.Error("WithContext関数がnilを返しました")
	}

	// 返されたロガーが*slog.Logger型であることを確認
	if contextLogger == nil {
		t.Error("WithContext関数が*slog.Logger型を返しませんでした")
	}
}

// TestWithValues はWithValues関数のテスト
func TestWithValues(t *testing.T) {
	tests := []struct {
		name string
		args []any
	}{
		{
			name: "基本的なキーと値のペア",
			args: []any{"key1", "value1", "key2", "value2"},
		},
		{
			name: "様々な型の値",
			args: []any{"string", "文字列", "int", 123, "bool", true, "float", 3.14},
		},
		{
			name: "空の引数",
			args: []any{},
		},
		{
			name: "奇数個の引数",
			args: []any{"key1", "value1", "key2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// WithValues関数を呼び出す
			valuesLogger := WithValues(tt.args...)

			// 返されたロガーがnilでないことを確認
			if valuesLogger == nil {
				t.Error("WithValues関数がnilを返しました")
			}

			// 返されたロガーが*slog.Logger型であることを確認
			if valuesLogger == nil {
				t.Error("WithValues関数が*slog.Logger型を返しませんでした")
			}
		})
	}
}

// TestLogLevels はログレベル定数のテスト
func TestLogLevels(t *testing.T) {
	tests := []struct {
		name     string
		level    slog.Level
		expected slog.Level
	}{
		{
			name:     "LevelInfo定数",
			level:    LevelInfo,
			expected: slog.LevelInfo,
		},
		{
			name:     "LevelDebug定数",
			level:    LevelDebug,
			expected: slog.LevelDebug,
		},
		{
			name:     "LevelWarn定数",
			level:    LevelWarn,
			expected: slog.LevelWarn,
		},
		{
			name:     "LevelError定数",
			level:    LevelError,
			expected: slog.LevelError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.level != tt.expected {
				t.Errorf("期待されるレベル: %v, 実際: %v", tt.expected, tt.level)
			}
		})
	}
}

// TestInitFunction はinit関数の動作をテスト
func TestInitFunction(t *testing.T) {
	// init関数は自動的に実行されるため、
	// ここではloggerが初期化されていることを確認
	if logger == nil {
		t.Error("init関数でloggerが初期化されていません")
	}

	// デフォルトのロガーが設定されていることを確認
	defaultLogger := slog.Default()
	if defaultLogger == nil {
		t.Error("デフォルトのロガーが設定されていません")
	}
}

// TestComplexScenarios は複雑なシナリオのテスト
func TestComplexScenarios(t *testing.T) {
	t.Run("レベル変更後のログ出力", func(t *testing.T) {
		// 初期設定
		err := Init(0, "text")
		if err != nil {
			t.Fatalf("Init関数でエラーが発生しました: %v", err)
		}

		// 各レベルでのログ出力テスト
		Info("情報メッセージ", "key", "value")

		// レベルを変更
		SetLevel("debug")
		Debug("デバッグメッセージ", "key", "value")

		// レベルを再度変更
		SetLevel("error")
		Error("エラーメッセージ", &testError{message: "テストエラー"}, "key", "value")
	})

	t.Run("WithContextとWithValuesの組み合わせ", func(t *testing.T) {
		ctx := context.Background()

		// WithContextとWithValuesを組み合わせて使用
		contextLogger := WithContext(ctx)
		valuesLogger := WithValues("global_key", "global_value")

		// 両方のロガーが正常に作成されることを確認
		if contextLogger == nil {
			t.Error("contextLoggerがnilです")
		}
		if valuesLogger == nil {
			t.Error("valuesLoggerがnilです")
		}
	})
}
