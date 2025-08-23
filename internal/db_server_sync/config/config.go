package config

import (
	"flag"
	"fmt"
	"os"
)

const (
	helpMsgOfOperation              = "実行する操作 (append-anime)"
	helpMsgOfInputFilePath          = "入力ファイルのパス（必須）"
	helpMsgOfOutputFilePath         = "出力ファイルのパス（必須）"
	helpMsgOfAdditionalInputFilePath = "追加入力ファイルのパス（任意）"
)

// Config はCLIツールの設定を保持する構造体
type Config struct {
	Operation                string // append-anime
	InputFilePath            string // 入力ファイルパス（必須）
	OutputFilePath           string // 出力ファイルパス（必須）
	AdditionalInputFilePath  string // 追加入力ファイルパス（任意）
	Help                     bool   // ヘルプ表示
}

// ParseFlags はコマンドライン引数を解析してConfigを返す
func ParseFlags() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.Operation, "operation", "", helpMsgOfOperation)
	flag.StringVar(&cfg.InputFilePath, "input-file-path", "", helpMsgOfInputFilePath)
	flag.StringVar(&cfg.OutputFilePath, "output-file-path", "", helpMsgOfOutputFilePath)
	flag.StringVar(&cfg.AdditionalInputFilePath, "additional-input-file-path", "", helpMsgOfAdditionalInputFilePath)
	flag.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")

	flag.Parse()

	// バリデーション
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate は設定の妥当性を検証する
func (c *Config) validate() error {
	// ヘルプが要求された場合はバリデーションをスキップ
	if c.Help {
		return nil
	}

	// 必須パラメータのチェック
	if c.Operation == "" {
		return fmt.Errorf("operationは必須です")
	}
	if c.InputFilePath == "" {
		return fmt.Errorf("input-file-pathは必須です")
	}
	if c.OutputFilePath == "" {
		return fmt.Errorf("output-file-pathは必須です")
	}

	// 操作の妥当性チェック
	if !isValidOperation(c.Operation) {
		return fmt.Errorf("未対応の操作です: %s %s", c.Operation, helpMsgOfOperation)
	}

	return nil
}

// isValidOperation は操作が有効かどうかを判定する
func isValidOperation(operation string) bool {
	validOperations := []string{"append-anime"}
	for _, valid := range validOperations {
		if operation == valid {
			return true
		}
	}
	return false
}

// PrintUsage はヘルプメッセージを表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: db-server-sync [オプション]\n\n")
	fmt.Fprintf(os.Stderr, "オプション:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n")
	fmt.Fprintf(os.Stderr, "        "+helpMsgOfOperation+"\n")
	fmt.Fprintf(os.Stderr, "  -input-file-path string\n")
	fmt.Fprintf(os.Stderr, "        "+helpMsgOfInputFilePath+"\n")
	fmt.Fprintf(os.Stderr, "  -output-file-path string\n")
	fmt.Fprintf(os.Stderr, "        "+helpMsgOfOutputFilePath+"\n")
	fmt.Fprintf(os.Stderr, "  -additional-input-file-path string\n")
	fmt.Fprintf(os.Stderr, "        "+helpMsgOfAdditionalInputFilePath+"\n")
	fmt.Fprintf(os.Stderr, "  -help\n")
	fmt.Fprintf(os.Stderr, "        ヘルプを表示\n\n")
	fmt.Fprintf(os.Stderr, "使用例:\n")
	fmt.Fprintf(os.Stderr, "  # AniListデータをリクエストボディ形式に変換\n")
	fmt.Fprintf(os.Stderr, "  db-server-sync -operation=append-anime -input-file-path=anilist.json -output-file-path=output.json\n\n")
	fmt.Fprintf(os.Stderr, "  # 追加データと結合してリクエストボディ形式に変換\n")
	fmt.Fprintf(os.Stderr, "  db-server-sync -operation=append-anime -input-file-path=anilist.json -additional-input-file-path=additional.json -output-file-path=output.json\n\n")
	fmt.Fprintf(os.Stderr, "対応する操作:\n")
	fmt.Fprintf(os.Stderr, "  - append-anime: AniListデータをリクエストボディ形式に変換\n")
}
