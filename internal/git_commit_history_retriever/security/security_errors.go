package security

import (
	"fmt"
)

// SecurityErrorType はセキュリティエラーの種類を表す
type SecurityErrorType int

const (
	// ErrorTypeUnknown は不明なエラー
	ErrorTypeUnknown SecurityErrorType = iota
	// ErrorTypeEmptyPath は空のパス
	ErrorTypeEmptyPath
	// ErrorTypePathTooLong はパスが長すぎる
	ErrorTypePathTooLong
	// ErrorTypeDangerousPattern は危険なパターン
	ErrorTypeDangerousPattern
	// ErrorTypePathTraversal はパストラバーサル攻撃
	ErrorTypePathTraversal
	// ErrorTypeUnauthorizedPath は許可されていないパス
	ErrorTypeUnauthorizedPath
	// ErrorTypeDirectoryNotFound はディレクトリが存在しない
	ErrorTypeDirectoryNotFound
	// ErrorTypeNotDirectory はディレクトリではない
	ErrorTypeNotDirectory
	// ErrorTypeNotGitRepository はGitリポジトリではない
	ErrorTypeNotGitRepository
	// ErrorTypeSymlinkAttack はシンボリックリンク攻撃
	ErrorTypeSymlinkAttack
	// ErrorTypeEncodingAttack はエンコーディング攻撃
	ErrorTypeEncodingAttack
	// ErrorTypeSystemPath はシステムパス
	ErrorTypeSystemPath
)

// SecurityError はセキュリティ関連のエラーを表す構造体
type SecurityError struct {
	Type    SecurityErrorType
	Message string
	Path    string
	Details string
}

// Error はerrorインターフェースの実装
func (e *SecurityError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (詳細: %s)", e.Path, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Path, e.Message)
}

// NewSecurityError は新しいSecurityErrorを作成する
func NewSecurityError(errorType SecurityErrorType, path, message, details string) *SecurityError {
	return &SecurityError{
		Type:    errorType,
		Message: message,
		Path:    path,
		Details: details,
	}
}

// IsSecurityError は指定されたエラーがSecurityErrorかどうかを判定する
func IsSecurityError(err error) bool {
	_, ok := err.(*SecurityError)
	return ok
}

// GetSecurityErrorType はSecurityErrorの種類を取得する
func GetSecurityErrorType(err error) SecurityErrorType {
	if secErr, ok := err.(*SecurityError); ok {
		return secErr.Type
	}
	return ErrorTypeUnknown
}
