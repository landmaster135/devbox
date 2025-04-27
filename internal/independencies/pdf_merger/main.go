// img2pdf.go
//
// 使い方:
//
//	go run img2pdf.go -dir="C:\\images" -out="merged.pdf"
//
// 依存パッケージ:
//
//	go get github.com/jung-kurt/gofpdf
package main

import (
	"flag"
	"fmt"

	usecases "github.com/landmaster135/devbox/internal/independencies/pdf_merger/usecases"
)

func main() {
	// ---- コマンドライン引数 ----
	dir := flag.String("dir", ".", "画像を検索するフォルダー (再帰探索)")
	out := flag.String("out", "", "出力 PDF ファイル名 (未指定なら <dir 名>.pdf)")
	flag.Parse()

	images, output, err := usecases.GetSourceImages(*dir, *out)
	if err != nil {
		usecases.Check(err)
	}
	if images != nil {
		usecases.Check(err)
	}

	fmt.Printf("検出した画像: %d 枚\n", len(images))
	fmt.Printf("出力 PDF   : %s\n", *out)

	// ---- PDF generating ----
	err = usecases.MergeImagesIntoPDF(images, output)
	if err != nil {
		usecases.Check(err)
	}

	fmt.Println("PDF を生成しました。完了です。")
}
