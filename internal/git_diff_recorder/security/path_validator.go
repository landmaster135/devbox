package security

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// #==============================================================#
// ##       Mocks for NotionConfig                               ##
// #==============================================================#
// MockPathValidator はテスト用のPathValidatorモック
type MockPathValidator struct {
	ValidateFunc func(string) (string, error)
}

func (m *MockPathValidator) ValidateWorkingDirectory(path string) (string, error) {
	if m.ValidateFunc != nil {
		return m.ValidateFunc(path)
	}
	return path, nil
}

// #==============================================================#
// ##       Interfaces for NotionConfig                          ##
// #==============================================================#
// PathValidatorInterface はテスト用のインターフェース
type PathValidatorInterface interface {
	ValidateWorkingDirectory(path string) (string, error)
}

// #==============================================================#
// ##       Implementations for NotionConfig                     ##
// #==============================================================#
// PathValidator はパス検証を行う構造体
type PathValidator struct {
	allowedBasePaths []string
	maxPathLength    int
}

// NewPathValidator は新しいPathValidatorを作成する
func NewPathValidator(allowedBasePaths []string, maxPathLength int) *PathValidator {
	if maxPathLength <= 0 {
		maxPathLength = 4096 // デフォルトの最大パス長
	}

	return &PathValidator{
		allowedBasePaths: allowedBasePaths,
		maxPathLength:    maxPathLength,
	}
}

// NewDefaultPathValidator はデフォルト設定でPathValidatorを作成する
func NewDefaultPathValidator() *PathValidator {
	// デフォルトでは現在のユーザーのホームディレクトリ配下のみ許可
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// ホームディレクトリが取得できない場合は空のリストで初期化
		return NewPathValidator([]string{}, 4096)
	}

	return NewPathValidator([]string{homeDir}, 4096)
}

// ValidateWorkingDirectory はワーキングディレクトリのパスを検証する
func (pv *PathValidator) ValidateWorkingDirectory(path string) (string, error) {
	if path == "" {
		return "", errors.New("パスが空です")
	}

	// 長さ制限チェック
	if len(path) > pv.maxPathLength {
		return "", fmt.Errorf("パスが長すぎます (最大: %d文字)", pv.maxPathLength)
	}

	// 危険な文字列のチェック
	if err := pv.checkDangerousPatterns(path); err != nil {
		return "", err
	}

	// 絶対パスに変換
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("絶対パスへの変換に失敗しました: %w", err)
	}

	// パスのクリーニング（../ などを解決）
	cleanPath := filepath.Clean(absPath)

	// 許可されたベースパス配下かチェック
	if err := pv.checkAllowedBasePaths(cleanPath); err != nil {
		return "", err
	}

	// ディレクトリの存在確認
	if err := pv.checkDirectoryExists(cleanPath); err != nil {
		return "", err
	}

	// Gitリポジトリかどうかの確認
	if err := pv.checkIsGitRepository(cleanPath); err != nil {
		return "", err
	}

	return cleanPath, nil
}

// checkDangerousPatterns は危険なパターンをチェックする
func (pv *PathValidator) checkDangerousPatterns(path string) error {
	dangerousPatterns := []string{
		"..",   // パストラバーサル
		"~",    // ホームディレクトリ展開
		"$",    // 環境変数展開
		"`",    // コマンド実行
		";",    // コマンド区切り
		"|",    // パイプ
		"&",    // バックグラウンド実行
		">",    // リダイレクト
		"<",    // リダイレクト
		"*",    // ワイルドカード
		"?",    // ワイルドカード
		"\n",   // 改行
		"\r",   // キャリッジリターン
		"\t",   // タブ
		"\x00", // NULL文字
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(path, pattern) {
			return fmt.Errorf("危険な文字が含まれています: %s", pattern)
		}
	}

	return nil
}

// checkAllowedBasePaths は許可されたベースパス配下かチェックする
func (pv *PathValidator) checkAllowedBasePaths(path string) error {
	if len(pv.allowedBasePaths) == 0 {
		// 許可されたベースパスが設定されていない場合はスキップ
		return nil
	}

	for _, basePath := range pv.allowedBasePaths {
		absBasePath, err := filepath.Abs(basePath)
		if err != nil {
			continue
		}

		cleanBasePath := filepath.Clean(absBasePath)

		// パスがベースパス配下にあるかチェック
		if strings.HasPrefix(path, cleanBasePath+string(filepath.Separator)) || path == cleanBasePath {
			return nil
		}
	}

	return fmt.Errorf("許可されていないディレクトリです: %s", path)
}

// checkDirectoryExists はディレクトリの存在を確認する
func (pv *PathValidator) checkDirectoryExists(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("ディレクトリが存在しません: %s", path)
		}
		return fmt.Errorf("ディレクトリの確認に失敗しました: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("指定されたパスはディレクトリではありません: %s", path)
	}

	return nil
}

// checkIsGitRepository はGitリポジトリかどうかを確認する
func (pv *PathValidator) checkIsGitRepository(path string) error {
	gitDir := filepath.Join(path, ".git")

	// .gitディレクトリまたはファイルの存在確認
	_, err := os.Stat(gitDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("gitリポジトリではありません: %s", path)
		}
		return fmt.Errorf("gitリポジトリの確認に失敗しました: %w", err)
	}

	return nil
}
