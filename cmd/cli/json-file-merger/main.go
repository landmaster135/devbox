package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/landmaster135/devbox/internal/interfaces/repositories"
	"github.com/landmaster135/devbox/internal/usecases/services"
)

func main() {
	// コマンドライン引数の定義
	dirPath := flag.String("dir", "", "JSONファイルが格納されているディレクトリのパス")
	keyName := flag.String("key", "pc_stats", "JSONデータの配列が入るキーの名前")
	outputPath := flag.String("output", "", "作成したリクエストボディを保存するファイルのパス（省略可）")
	flag.Parse()

	// 必須引数のチェック
	if *dirPath == "" {
		fmt.Println("エラー: ディレクトリパスは必須です")
		fmt.Println("使用方法: json-file-merger -dir <JSONディレクトリパス> [-key <キー名>] [-output <出力ファイルパス>]")
		fmt.Println("例: json-file-merger -dir ../../sample_data -key pc_stats -output output.json")
		os.Exit(1)
	}

	// リポジトリとサービスの作成
	fileRepo := repositories.NewFileRepository()
	fileService := services.NewFileService(fileRepo)

	// リクエストボディの作成
	requestBodyJSON, err := fileService.CreateRequestBodyFromDir(*dirPath, *keyName, *outputPath)
	if err != nil {
		log.Fatalf("リクエストボディの作成に失敗しました: %v", err)
	}

	// 結果の表示
	fmt.Println("リクエストボディが作成されました:")
	fmt.Println(string(requestBodyJSON))

	// 出力ファイルが指定されている場合
	if *outputPath != "" {
		fmt.Printf("リクエストボディが %s に保存されました\n", *outputPath)
	}
}
