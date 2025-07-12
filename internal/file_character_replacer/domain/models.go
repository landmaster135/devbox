package domain

import "fmt"

// EncodingType は文字エンコーディングの種別を表します
type EncodingType string

const (
	EncodingUTF8     EncodingType = "utf-8"
	EncodingShiftJIS EncodingType = "shift_jis"
	EncodingEUCJP    EncodingType = "euc-jp"
	EncodingISO2022JP EncodingType = "iso-2022-jp"
)

// IsValid はエンコーディングタイプが有効かどうかを確認します
func (e EncodingType) IsValid() bool {
	switch e {
	case EncodingUTF8, EncodingShiftJIS, EncodingEUCJP, EncodingISO2022JP:
		return true
	default:
		return false
	}
}

// ReplacementConfig は文字列置換の設定を保持するドメインモデルです
type ReplacementConfig struct {
	Target    string       // 対象パス（ファイルまたはディレクトリ）
	From      string       // 置換元文字列
	To        string       // 置換先文字列
	Encoding  EncodingType // 文字エンコーディング
	Recursive bool         // 再帰的処理
	Backup    bool         // バックアップ作成
	DryRun    bool         // ドライラン
}

// Validate は設定の妥当性を検証します
func (c *ReplacementConfig) Validate() error {
	if c.Target == "" {
		return fmt.Errorf("対象パスが指定されていません")
	}
	if c.From == "" {
		return fmt.Errorf("置換元文字列が指定されていません")
	}
	if c.To == "" {
		return fmt.Errorf("置換先文字列が指定されていません")
	}
	if !c.Encoding.IsValid() {
		return fmt.Errorf("無効な文字エンコーディングです: %s", c.Encoding)
	}
	return nil
}

// FileProcessResult は処理結果を表すドメインモデルです
type FileProcessResult struct {
	ProcessedFiles int      // 処理されたファイル数
	ReplacedCount  int      // 置換された箇所数
	Errors         []error  // 発生したエラー
	Messages       []string // 処理メッセージ
}

// AddError はエラーを追加します
func (r *FileProcessResult) AddError(err error) {
	r.Errors = append(r.Errors, err)
}

// AddMessage はメッセージを追加します
func (r *FileProcessResult) AddMessage(msg string) {
	r.Messages = append(r.Messages, msg)
}

// HasErrors はエラーが発生したかどうかを確認します
func (r *FileProcessResult) HasErrors() bool {
	return len(r.Errors) > 0
}

// FileInfo はファイル情報を表すドメインモデルです
type FileInfo struct {
	Path     string // ファイルパス
	IsDir    bool   // ディレクトリかどうか
	Size     int64  // ファイルサイズ
	Encoding EncodingType // 文字エンコーディング
}
