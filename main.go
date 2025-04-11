package main

import (
	"fmt"
	"io"
	"os"
)

// run はアプリケーションのメイン処理を実行します
func run(w io.Writer) int {
	// メインメッセージを出力
	fmt.Fprint(w, "Devbox: ファイル処理ユーティリティ\n")
	fmt.Fprint(w, "詳細なコマンドは cmd/file-processor を使用してください。\n")
	return 0
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Stdout)
	os.Exit(code)
}
