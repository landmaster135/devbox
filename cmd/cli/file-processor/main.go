package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/landmaster135/devbox/internal/independencies/file_processor/interfaces/repositories"
	"github.com/landmaster135/devbox/internal/independencies/file_processor/usecases/services"
)

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

// run はファイル処理の主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグセットを作成
	fs := flag.NewFlagSet("file-processor", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// コマンドライン引数の定義
	filePath := fs.String("file", "", "処理するファイルのパス")
	startPos := fs.Int("start", 0, "各行の文字列を取得する開始位置（0ベース）")
	endPos := fs.Int("end", 0, "各行の文字列を取得する終了位置")

	// 引数の解析
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// ファイルパスが指定されていない場合はエラー
	if *filePath == "" {
		fmt.Fprintln(stderr, "エラー: ファイルパスを指定してください（-file オプション）")
		fs.Usage()
		return exitCodeError
	}

	// 開始位置と終了位置が指定されていない場合はエラー
	if *startPos == 0 && *endPos == 0 {
		fmt.Fprintln(stderr, "エラー: 開始位置と終了位置を指定してください（-start と -end オプション）")
		fs.Usage()
		return exitCodeError
	}

	// 依存関係の注入
	fileRepo := repositories.NewFileRepository()
	fileService := services.NewFileService(fileRepo)

	// RemoveMatchingLines関数を呼び出す
	count, err := fileService.RemoveMatchingLines(*filePath, *startPos, *endPos)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		return exitCodeError
	}

	fmt.Fprintf(stdout, "処理完了: %d行の重複を削除しました\n", count)
	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
