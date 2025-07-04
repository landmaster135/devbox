package config

import (
	"encoding/json"
	"os"
)

// AppConfig はアプリケーションの設定を表します
type AppConfig struct {
	ConfigPath  string
	SourceFile  string
	Extension   string
	OutputPath  string
	Format      string
	Recursive   bool
	Verbose     bool
	Directory   string
}

// 言語設定を表す構造体
type Language struct {
	FunctionHeader   string `json:"function_header"`
	FunctionTail     string `json:"function_tail"`
	MainMarker       string `json:"main_marker"`
	CommentPrefix    string `json:"comment_prefix"`
	MultilineComment bool   `json:"multiline_comment"`
}

// 全体設定を表す構造体
type Config struct {
	Spaces       []string            `json:"spaces"`
	Languages    map[string]Language `json:"languages"`
	OutputFormat string              `json:"output_format"`
}

// デフォルト設定
var DefaultConfig = Config{
	Spaces: []string{" ", "\t"},
	Languages: map[string]Language{
		".go": {
			FunctionHeader:   "func ",
			FunctionTail:     "(",
			MainMarker:       "func main() {",
			CommentPrefix:    "//",
			MultilineComment: true,
		},
		".py": {
			FunctionHeader:   "def ",
			FunctionTail:     "(",
			MainMarker:       "if __name__ == \"__main__\":",
			CommentPrefix:    "#",
			MultilineComment: true,
		},
		".js": {
			FunctionHeader:   "function ",
			FunctionTail:     "(",
			MainMarker:       "}",
			CommentPrefix:    "//",
			MultilineComment: true,
		},
	},
	OutputFormat: "mermaid",
}

// 設定を外部ファイルから読み込む
func LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	// グローバル設定を更新
	DefaultConfig = config
	return nil
}

// 言語の設定を取得する
func GetLanguageConfig(extension string) (Language, bool) {
	lang, ok := DefaultConfig.Languages[extension]
	return lang, ok
}

// スペースの設定を取得する
func GetSpaces() []string {
	return DefaultConfig.Spaces
}
