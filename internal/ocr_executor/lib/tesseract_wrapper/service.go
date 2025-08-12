package tesseract_wrapper

import (
	"fmt"
	"os/exec"
	"strings"
)

// CommandExecutor はコマンド実行のインターフェース
type CommandExecutor interface {
	LookPath(file string) (string, error)
	Command(name string, arg ...string) *exec.Cmd
}

// RealCommandExecutor は実際のコマンド実行を行う
type RealCommandExecutor struct{}

// LookPath は実際のexec.LookPathを呼び出す
func (r *RealCommandExecutor) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

// Command は実際のexec.Commandを呼び出す
func (r *RealCommandExecutor) Command(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

// TesseractClient はTesseractのCLIラッパー
type TesseractClient struct {
	binPath   string
	languages string
	executor  CommandExecutor
}

// TesseractOption はTesseractClientの設定オプション
type TesseractOption func(*TesseractClient)

// WithBinPath はTesseractバイナリのパスを設定する
func WithBinPath(path string) TesseractOption {
	return func(tc *TesseractClient) {
		tc.binPath = path
	}
}

// WithLanguages は言語設定を行う
func WithLanguages(languages string) TesseractOption {
	return func(tc *TesseractClient) {
		tc.languages = languages
	}
}

// WithExecutor はCommandExecutorを設定する（テスト用）
func WithExecutor(executor CommandExecutor) TesseractOption {
	return func(tc *TesseractClient) {
		tc.executor = executor
	}
}

// NewTesseractClient は新しいTesseractClientを作成する
func NewTesseractClient(options ...TesseractOption) (*TesseractClient, error) {
	tc := &TesseractClient{
		binPath:  "tesseract",            // デフォルトはPATHから検索
		executor: &RealCommandExecutor{}, // デフォルトは実際のコマンド実行
	}

	// オプションを適用
	for _, option := range options {
		option(tc)
	}

	// Tesseractバイナリの存在確認
	if !tc.isAvailable() {
		return nil, fmt.Errorf("tesseractバイナリが見つかりません: %s", tc.binPath)
	}

	return tc, nil
}

// TextFromImageFile は画像ファイルからテキストを抽出する
func (tc *TesseractClient) TextFromImageFile(filePath string) (string, error) {
	// Tesseractコマンドの引数を構築
	args := []string{filePath, "stdout"}

	// 言語設定がある場合は追加
	if tc.languages != "" {
		args = append(args, "-l", tc.languages)
	}

	// Tesseractコマンドを実行
	cmd := tc.executor.Command(tc.binPath, args...)
	output, err := cmd.Output()
	if err != nil {
		// エラーの詳細を取得
		if exitError, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("tesseract実行エラー: %v, stderr: %s", err, string(exitError.Stderr))
		}
		return "", fmt.Errorf("tesseract実行エラー: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// isAvailable はTesseractバイナリが利用可能かを確認する
func (tc *TesseractClient) isAvailable() bool {
	_, err := tc.executor.LookPath(tc.binPath)
	return err == nil
}
