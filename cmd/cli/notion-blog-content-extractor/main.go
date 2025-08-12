package main

import (
	"fmt"
	"os"

	"github.com/landmaster135/devbox/internal/notion_blog_content_extractor/config"
	"github.com/landmaster135/devbox/internal/notion_blog_content_extractor/usecases"
)

func main() {
	// コマンドライン引数の解析
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	// ヘルプが要求された場合
	if cfg.Help {
		config.PrintUsage()
		return
	}

	// サービスの作成
	service := usecases.NewService()

	// ブログコンテンツの抽出
	result, err := service.ExtractBlogContent(cfg.SrcDir, cfg.DestDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を標準出力に表示
	fmt.Print(result)
}
