package main

import (
	"fmt"
	"os"

	"github.com/landmaster135/devbox/internal/file_character_replacer/config"
	"github.com/landmaster135/devbox/internal/file_character_replacer/usecases"
)

func main() {
	// サービスを作成（依存関係注入とコマンドライン引数解析を含む）
	service, err := usecases.NewFileReplacerServiceWithConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	// 文字列置換を実行
	result, err := service.ReplaceStrings()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	summary := service.GetSummary(result)
	fmt.Print(summary)

	// エラーがあった場合は終了コード1で終了
	if result.HasErrors() {
		os.Exit(1)
	}
}
