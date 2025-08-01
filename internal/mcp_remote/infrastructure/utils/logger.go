package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Logger はmcp-remote用のロガー
type Logger struct {
	*log.Logger
	debug       bool
	debugFile   *os.File
	serverHash  string
}

// NewLogger は新しいLoggerを作成する
func NewLogger(prefix string, debug bool, serverHash string) *Logger {
	logger := &Logger{
		Logger:     log.New(os.Stderr, fmt.Sprintf("[%s] ", prefix), log.LstdFlags),
		debug:      debug,
		serverHash: serverHash,
	}

	if debug && serverHash != "" {
		if err := logger.initDebugFile(); err != nil {
			logger.Printf("デバッグファイルの初期化に失敗しました: %v", err)
		}
	}

	return logger
}

// initDebugFile はデバッグファイルを初期化する
func (l *Logger) initDebugFile() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return fmt.Errorf("設定ディレクトリの取得に失敗しました: %v", err)
	}

	// 設定ディレクトリを作成
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("設定ディレクトリの作成に失敗しました: %v", err)
	}

	// デバッグファイルを開く
	debugFilePath := filepath.Join(configDir, fmt.Sprintf("%s_debug.log", l.serverHash))
	debugFile, err := os.OpenFile(debugFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("デバッグファイルの作成に失敗しました: %v", err)
	}

	l.debugFile = debugFile
	l.Printf("デバッグログファイル: %s", debugFilePath)
	return nil
}

// Debug はデバッグメッセージをログ出力する
func (l *Logger) Debug(format string, args ...interface{}) {
	if !l.debug {
		return
	}

	message := fmt.Sprintf(format, args...)

	// 標準エラー出力にも出力
	l.Printf("[DEBUG] %s", message)

	// デバッグファイルにも出力
	if l.debugFile != nil {
		debugMessage := fmt.Sprintf("[%s][DEBUG] %s\n",
			l.Logger.Prefix(), message)
		l.debugFile.WriteString(debugMessage)
	}
}

// Close はロガーのリソースを解放する
func (l *Logger) Close() error {
	if l.debugFile != nil {
		return l.debugFile.Close()
	}
	return nil
}

// GetConfigDir は設定ディレクトリのパスを取得する
func GetConfigDir() (string, error) {
	// 環境変数MCP_REMOTE_CONFIG_DIRが設定されている場合はそれを使用
	if configDir := os.Getenv("MCP_REMOTE_CONFIG_DIR"); configDir != "" {
		return configDir, nil
	}

	// デフォルトは~/.mcp-auth
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("ホームディレクトリの取得に失敗しました: %v", err)
	}

	return filepath.Join(homeDir, ".mcp-auth"), nil
}
