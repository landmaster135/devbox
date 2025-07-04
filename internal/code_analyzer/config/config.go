// internal/config/config.go
package config

import "strings"

// AppConfig はアプリケーションの設定情報を保持します
type AppConfig struct {
	ProjectPath   string   // 分析対象のプロジェクトパス
	OutputFormat  string   // 出力フォーマット (text, json, csv)
	OutputFile    string   // 出力ファイルパス
	Extensions    []string // 分析対象の拡張子リスト
	VisualReport  bool     // ビジュアルレポート生成フラグ
	HistoryPath   string   // 履歴データファイルのパス
	DetectClones  bool     // コードクローン検出フラグ
	MinBlockSize  int      // クローン検出の最小ブロックサイズ
	MinSimilarity float64  // クローン検出の最小類似度
	Verbose       bool     // 詳細ログ出力フラグ
}

// NewAppConfig は新しいAppConfigインスタンスを作成します
func NewAppConfig() *AppConfig {
	return &AppConfig{
		ProjectPath:   ".",
		OutputFormat:  "text",
		Extensions:    []string{".go"},
		MinBlockSize:  30,
		MinSimilarity: 0.8,
	}
}

// SetExtensions は拡張子リストを設定します
func (c *AppConfig) SetExtensions(extensions string) {
	// カンマ区切りの拡張子を処理
	if extensions == "" {
		return
	}

	c.Extensions = []string{}
	for _, ext := range strings.Split(extensions, ",") {
		ext = strings.TrimSpace(ext)
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		c.Extensions = append(c.Extensions, ext)
	}
}
