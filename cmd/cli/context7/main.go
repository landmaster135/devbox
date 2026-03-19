package main

import (
	"flag"
	"fmt"
	"os"

	models "github.com/landmaster135/devbox/internal/context7/domain/models"
	usecases "github.com/landmaster135/devbox/internal/context7/usecases"
)

func handleSearchCommand(service *usecases.Context7Service, args []string) {
	if len(args) == 0 {
		fmt.Println("エラー: 検索するライブラリ名を指定してください")
		fmt.Println("使用例: context7 search react")
		os.Exit(1)
	}

	libraryName := args[0]

	// ライブラリを検索
	searchResponse, err := service.ResolveLibraryID(libraryName)
	if err != nil {
		fmt.Printf("検索エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を表示
	result := service.FormatSearchResults(searchResponse)
	fmt.Print(result)
}

func handleDocsCommand(service *usecases.Context7Service, args []string) {
	if len(args) == 0 {
		fmt.Println("エラー: ライブラリIDを指定してください")
		fmt.Println("使用例: context7 docs /facebook/react")
		os.Exit(1)
	}

	// フラグを定義
	var topic string
	var tokens int

	flagSet := flag.NewFlagSet("docs", flag.ExitOnError)
	flagSet.StringVar(&topic, "topic", "", "特定のトピックに焦点を当てる（例: hooks, routing）")
	flagSet.IntVar(&tokens, "tokens", models.DefaultTokens, "取得する最大トークン数")

	// ライブラリIDを取得（最初の引数）
	libraryID := args[0]

	// 残りの引数をフラグとして解析
	if len(args) > 1 {
		err := flagSet.Parse(args[1:])
		if err != nil {
			fmt.Printf("フラグ解析エラー: %v\n", err)
			os.Exit(1)
		}
	}

	// ライブラリIDの形式を検証
	if err := service.ValidateLibraryID(libraryID); err != nil {
		fmt.Printf("ライブラリID検証エラー: %v\n", err)
		os.Exit(1)
	}

	// ドキュメントオプションを設定
	options := models.DocOptions{
		Topic:  topic,
		Tokens: tokens,
	}

	// ドキュメントを取得
	docs, err := service.GetLibraryDocs(libraryID, options)
	if err != nil {
		fmt.Printf("ドキュメント取得エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を表示
	fmt.Printf("=== %s のドキュメント ===\n\n", libraryID)
	if topic != "" {
		fmt.Printf("トピック: %s\n", topic)
	}
	if tokens != models.DefaultTokens {
		fmt.Printf("トークン数: %d\n", tokens)
	}
	fmt.Println()
	fmt.Print(docs)
}

func printUsage() {
	fmt.Println("Context7 CLI - 最新のライブラリドキュメントを取得")
	fmt.Println()
	fmt.Println("使用方法:")
	fmt.Println("  context7 <command> [arguments]")
	fmt.Println()
	fmt.Println("利用可能なコマンド:")
	fmt.Println("  search <library_name>     ライブラリを検索してContext7互換IDを取得")
	fmt.Println("  docs <library_id> [flags] ライブラリIDからドキュメントを取得")
	fmt.Println("  help                      このヘルプメッセージを表示")
	fmt.Println()
	fmt.Println("searchコマンドの例:")
	fmt.Println("  context7 search react")
	fmt.Println("  context7 search \"next.js\"")
	fmt.Println()
	fmt.Println("docsコマンドの例:")
	fmt.Println("  context7 docs /facebook/react")
	fmt.Println("  context7 docs /vercel/next.js -topic=routing")
	fmt.Println("  context7 docs /mongodb/docs -tokens=5000")
	fmt.Println("  context7 docs /supabase/supabase -topic=auth -tokens=8000")
	fmt.Println()
	fmt.Println("docsコマンドのフラグ:")
	fmt.Println("  -topic string    特定のトピックに焦点を当てる（例: hooks, routing）")
	fmt.Printf("  -tokens int      取得する最大トークン数（デフォルト: %d）\n", models.DefaultTokens)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// HTTPクライアントとサービスを初期化
	service := usecases.NewContext7ServiceWithHTTPClient()

	switch command {
	case "search":
		handleSearchCommand(service, os.Args[2:])
	case "docs":
		handleDocsCommand(service, os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("不明なコマンド: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}
