// Package logger はlog/slogを使用した構造化ロギングを提供します
package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

var (
	// デフォルトのロガーを設定
	logger = slog.Default()
)

// init はパッケージの初期化時にデフォルトのロガーを設定します
func init() {
	// デフォルトのハンドラーを設定
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	// デフォルトのロガーを作成
	logger = slog.New(handler)
	slog.SetDefault(logger)
}

// ログレベルの定数
const (
	LevelInfo  = slog.LevelInfo
	LevelDebug = slog.LevelDebug
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

// Init はロガーを初期化します
func Init(level int, format string) error {
	if level < 0 || level > 3 {
		return fmt.Errorf("level is invalid: %d", level)
	}
	var slogLevel slog.Level

	// デフォルトはInfoレベル
	slogLevel = LevelInfo
	if level == 1 {
		slogLevel = LevelDebug
	} else if level == 2 {
		slogLevel = LevelWarn
	} else if level == 3 {
		slogLevel = LevelError
	}

	// フォーマットに基づいてハンドラーを設定
	var handler slog.Handler
	switch format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slogLevel,
		})
	case "text":
		fallthrough
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slogLevel,
		})
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)

	return nil
}

// SetLevel はログレベルを文字列から設定します
func SetLevel(levelStr string) {
	var level slog.Level
	switch levelStr {
	case "debug":
		level = LevelDebug
	case "info":
		level = LevelInfo
	case "warn":
		level = LevelWarn
	case "error":
		level = LevelError
	default:
		level = LevelInfo
	}

	// 新しいレベルでハンドラーを作成
	newHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	// 新しいロガーを作成
	logger = slog.New(newHandler)
	slog.SetDefault(logger)
}

// Debug はデバッグレベルのログを出力します
func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

// Info は情報レベルのログを出力します
func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

// Warn は警告レベルのログを出力します
func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

// Error はエラーレベルのログを出力します
func Error(msg string, err error, args ...any) {
	if err != nil {
		newArgs := append([]any{"error", err}, args...)
		logger.Error(msg, newArgs...)
	} else {
		logger.Error(msg, args...)
	}
}

// Fatal は致命的なエラーをログに記録し、プログラムを終了します
func Fatal(msg string, err error, args ...any) {
	if err != nil {
		newArgs := append([]any{"error", err}, args...)
		logger.Error(msg, newArgs...)
	} else {
		logger.Error(msg, args...)
	}
	os.Exit(1)
}

// WithContext はコンテキスト情報を含むロガーを返します
func WithContext(ctx context.Context) *slog.Logger {
	return logger.With("ctx", ctx)
}

// WithValues は追加の値を含むロガーを返します
func WithValues(args ...any) *slog.Logger {
	return logger.With(args...)
}
