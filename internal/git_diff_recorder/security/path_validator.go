package security

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
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
		return "", NewSecurityError(ErrorTypeEmptyPath, path, "パスが空です", "")
	}

	// 長さ制限チェック
	if len(path) > pv.maxPathLength {
		return "", NewSecurityError(ErrorTypePathTooLong, path, "パスが長すぎます", fmt.Sprintf("最大: %d文字", pv.maxPathLength))
	}

	// Unicode正規化とエンコーディング攻撃のチェック
	normalizedPath, err := pv.normalizeAndValidateEncoding(path)
	if err != nil {
		return "", err
	}

	// 危険な文字列のチェック（強化版）
	if err := pv.checkEnhancedDangerousPatterns(normalizedPath); err != nil {
		return "", err
	}

	// 絶対パスに変換
	absPath, err := filepath.Abs(normalizedPath)
	if err != nil {
		return "", fmt.Errorf("絶対パスへの変換に失敗しました: %w", err)
	}

	// パスのクリーニング（../ などを解決）
	cleanPath := filepath.Clean(absPath)

	// シンボリックリンクの検証
	if err := pv.checkSymlinkSafety(cleanPath); err != nil {
		return "", err
	}

	// システムパスの検証
	if err := pv.checkSystemPaths(cleanPath); err != nil {
		return "", err
	}

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

// normalizeAndValidateEncoding はUnicode正規化とエンコーディング攻撃の検証を行う
func (pv *PathValidator) normalizeAndValidateEncoding(path string) (string, error) {
	// UTF-8の妥当性チェック
	if !utf8.ValidString(path) {
		return "", NewSecurityError(ErrorTypeEncodingAttack, path, "無効なUTF-8エンコーディング", "")
	}

	// URLエンコードされた危険文字の検出
	if decoded, err := url.QueryUnescape(path); err == nil && decoded != path {
		// URLデコードできる場合は、デコード後も検証
		if err := pv.checkEncodedDangerousPatterns(decoded); err != nil {
			return "", err
		}
	}

	// 制御文字の検出
	for _, r := range path {
		if unicode.IsControl(r) && r != '\t' && r != '\n' && r != '\r' {
			return "", NewSecurityError(ErrorTypeEncodingAttack, path, "制御文字が含まれています", fmt.Sprintf("文字コード: U+%04X", r))
		}
	}

	return path, nil
}

// checkEncodedDangerousPatterns はエンコードされた危険パターンをチェックする
func (pv *PathValidator) checkEncodedDangerousPatterns(path string) error {
	encodedPatterns := []string{
		"%2e%2e",    // ..
		"%2E%2E",    // ..
		"%2e%2e%2f", // ../
		"%2E%2E%2F", // ../
		"%2e%2e%5c", // ..\
		"%2E%2E%5C", // ..\
		"%7e",       // ~
		"%7E",       // ~
		"%24",       // $
		"%60",       // `
		"%3b",       // ;
		"%3B",       // ;
		"%7c",       // |
		"%7C",       // |
		"%26",       // &
		"%3e",       // >
		"%3E",       // >
		"%3c",       // <
		"%3C",       // <
		"%2a",       // *
		"%2A",       // *
		"%3f",       // ?
		"%3F",       // ?
		"%00",       // NULL
	}

	lowerPath := strings.ToLower(path)
	for _, pattern := range encodedPatterns {
		if strings.Contains(lowerPath, pattern) {
			return NewSecurityError(ErrorTypeEncodingAttack, path, "エンコードされた危険文字が検出されました", pattern)
		}
	}

	return nil
}

// checkEnhancedDangerousPatterns は強化された危険パターンをチェックする
func (pv *PathValidator) checkEnhancedDangerousPatterns(path string) error {
	// 基本的な危険パターン
	basicPatterns := []string{
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

	for _, pattern := range basicPatterns {
		if strings.Contains(path, pattern) {
			return NewSecurityError(ErrorTypeDangerousPattern, path, "危険な文字が含まれています", pattern)
		}
	}

	// 正規表現による高度なパターンマッチング
	dangerousRegexes := []*regexp.Regexp{
		regexp.MustCompile(`\.\.[\\/]`),                    // パストラバーサル
		regexp.MustCompile(`[\\/]\.\.[\\/]`),               // 中間パストラバーサル
		regexp.MustCompile(`[\\/]\.\.[^[\\/\w]`),           // 偽装パストラバーサル
		regexp.MustCompile(`\$\{[^}]*\}`),                  // 環境変数展開
		regexp.MustCompile(`\$\([^)]*\)`),                  // コマンド置換
		regexp.MustCompile(`[;&|><]\s*[a-zA-Z0-9_\\/.-]+`), // コマンドインジェクション
	}

	for _, regex := range dangerousRegexes {
		if regex.MatchString(path) {
			return NewSecurityError(ErrorTypePathTraversal, path, "高度な攻撃パターンが検出されました", regex.String())
		}
	}

	return nil
}

// checkSymlinkSafety はシンボリックリンクの安全性をチェックする
func (pv *PathValidator) checkSymlinkSafety(path string) error {
	// パスの各コンポーネントをチェック
	currentPath := path
	for {
		info, err := os.Lstat(currentPath)
		if err != nil {
			if os.IsNotExist(err) {
				// 存在しないパスは後でチェックされるのでここではスキップ
				break
			}
			return fmt.Errorf("パス情報の取得に失敗しました: %w", err)
		}

		// シンボリックリンクの場合
		if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(currentPath)
			if err != nil {
				return NewSecurityError(ErrorTypeSymlinkAttack, path, "シンボリックリンクの読み取りに失敗", err.Error())
			}

			// 絶対パスの場合
			if filepath.IsAbs(target) {
				// 許可されたベースパス外を指している場合は危険
				if err := pv.checkAllowedBasePaths(target); err != nil {
					return NewSecurityError(ErrorTypeSymlinkAttack, path, "危険なシンボリックリンク", fmt.Sprintf("リンク先: %s", target))
				}
			} else {
				// 相対パスの場合、解決後のパスをチェック
				resolvedTarget := filepath.Join(filepath.Dir(currentPath), target)
				cleanTarget := filepath.Clean(resolvedTarget)
				if err := pv.checkAllowedBasePaths(cleanTarget); err != nil {
					return NewSecurityError(ErrorTypeSymlinkAttack, path, "危険なシンボリックリンク", fmt.Sprintf("解決後のリンク先: %s", cleanTarget))
				}
			}
		}

		// 親ディレクトリに移動
		parent := filepath.Dir(currentPath)
		if parent == currentPath {
			break
		}
		currentPath = parent
	}

	return nil
}

// checkSystemPaths はシステムパスへのアクセスをチェックする
func (pv *PathValidator) checkSystemPaths(path string) error {
	dangerousSystemPaths := []string{
		"/proc",
		"/sys",
		"/dev",
		"/etc/passwd",
		"/etc/shadow",
		"/etc/hosts",
		"/boot",
		"/root",
		"/var/log",
		"/tmp",
		"/var/tmp",
	}

	cleanPath := filepath.Clean(path)
	for _, dangerousPath := range dangerousSystemPaths {
		if strings.HasPrefix(cleanPath, dangerousPath) {
			return NewSecurityError(ErrorTypeSystemPath, path, "システムパスへのアクセスは禁止されています", dangerousPath)
		}
	}

	return nil
}

// checkDangerousPatterns は危険なパターンをチェックする（後方互換性のため残す）
func (pv *PathValidator) checkDangerousPatterns(path string) error {
	return pv.checkEnhancedDangerousPatterns(path)
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

	return NewSecurityError(ErrorTypeUnauthorizedPath, path, "許可されていないディレクトリです", "")
}

// checkDirectoryExists はディレクトリの存在を確認する
func (pv *PathValidator) checkDirectoryExists(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewSecurityError(ErrorTypeDirectoryNotFound, path, "ディレクトリが存在しません", "")
		}
		return fmt.Errorf("ディレクトリの確認に失敗しました: %w", err)
	}

	if !info.IsDir() {
		return NewSecurityError(ErrorTypeNotDirectory, path, "指定されたパスはディレクトリではありません", "")
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
			return NewSecurityError(ErrorTypeNotGitRepository, path, "gitリポジトリではありません", "")
		}
		return fmt.Errorf("gitリポジトリの確認に失敗しました: %w", err)
	}

	return nil
}
