package logger

import (
	"fmt"
	"log"
	"strings"

	logger "github.com/landmaster135/devbox/internal/valkey/infrastructure/logger"
)

// Logger はロギングのためのインターフェース
type Logger interface {
	// Debug はデバッグレベルのログを出力する
	Debug(msg string, keysAndValues ...any)

	// Info は情報レベルのログを出力する
	Info(msg string, keysAndValues ...any)

	// Warn は警告レベルのログを出力する
	Warn(msg string, keysAndValues ...any)

	// Error はエラーレベルのログを出力する
	Error(msg string, err error, keysAndValues ...any)

	// Init はロガーを初期化する
	Init(level int, format string) error
}

// DefaultLogger はLoggerインターフェースの標準実装
type DefaultLogger struct{}

// NewDefaultLogger は新しいDefaultLoggerを作成する
func NewDefaultLogger() Logger {
	return &DefaultLogger{}
}

// Debug はデバッグレベルのログを出力する
func (l *DefaultLogger) Debug(msg string, keysAndValues ...any) {
	l.log("DEBUG", msg, nil, keysAndValues...)
}

// Info は情報レベルのログを出力する
func (l *DefaultLogger) Info(msg string, keysAndValues ...any) {
	l.log("INFO", msg, nil, keysAndValues...)
}

// Warn は警告レベルのログを出力する
func (l *DefaultLogger) Warn(msg string, keysAndValues ...any) {
	l.log("WARN", msg, nil, keysAndValues...)
}

// Error はエラーレベルのログを出力する
func (l *DefaultLogger) Error(msg string, err error, keysAndValues ...any) {
	l.log("ERROR", msg, err, keysAndValues...)
}

// Init はロガーを初期化する
func (l *DefaultLogger) Init(level int, format string) error {
	// infrastructureのloggerパッケージのInit関数を呼び出す
	return logger.Init(level, format)
}

// log は実際のログ出力を行う内部メソッド
func (l *DefaultLogger) log(level, msg string, err error, keysAndValues ...any) {
	// キーと値のペアを文字列に変換
	var kvPairs []string
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			kvPairs = append(kvPairs, fmt.Sprintf("%v=%v", keysAndValues[i], keysAndValues[i+1]))
		} else {
			kvPairs = append(kvPairs, fmt.Sprintf("%v=<no value>", keysAndValues[i]))
		}
	}

	// エラー情報を追加
	if err != nil {
		kvPairs = append(kvPairs, fmt.Sprintf("error=%v", err))
	}

	// ログメッセージを構築
	logMsg := fmt.Sprintf("[%s] %s", level, msg)
	if len(kvPairs) > 0 {
		logMsg = fmt.Sprintf("%s %s", logMsg, strings.Join(kvPairs, " "))
	}

	// ログ出力
	log.Println(logMsg)
}
