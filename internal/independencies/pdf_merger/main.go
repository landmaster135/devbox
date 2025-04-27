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
	"os"
	"path/filepath"
	"sort"
	"strings"

	api "github.com/pdfcpu/pdfcpu/pkg/api"
	types "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func GetSourceImages(dir string, out string) ([]string, string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, "", err
	}

	// 出力名のデフォルト: <フォルダー名>.pdf
	if out == "" {
		base := filepath.Base(absDir)
		out = filepath.Join(absDir, base+".pdf")
	}

	// ---- 画像ファイルを収集 ----
	var images []string
	err = filepath.WalkDir(absDir, func(p string, d os.DirEntry, _ error) error {
		if d == nil || d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(p))
		if ext == ".jpg" {
			images = append(images, p)
		}
		return nil
	})
	if err != nil {
		return nil, "", err
	}

	if len(images) == 0 {
		fmt.Println("画像が見つかりませんでした。終了します。")
		return nil, "", nil
	}
	sort.Strings(images) // PowerShell の Sort-Object 相当

	return images, out, nil
}

func MergeImagesIntoPDF(images []string, output string) error {
	cfg := api.LoadConfiguration()
	// Unit is used in commands for layout
	cfg.Unit = types.POINTS
	// Compress non-stream object to stream object
	cfg.WriteObjectStream = true
	// Remove unused fonts and images from resource dictionary
	cfg.OptimizeResourceDicts = true
	// Share duplicated streams in all pages
	cfg.OptimizeDuplicateContentStreams = true
	// Set user password and owner password
	// cfg.UserPWNew = "aaaaa"
	// cfg.OwnerPWNew = "aaaaa"
	if err := api.ImportImagesFile(images, output, nil, cfg); err != nil {
		return err
	}
	return nil
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}
}

func main() {
	// ---- コマンドライン引数 ----
	dir := flag.String("dir", ".", "画像を検索するフォルダー (再帰探索)")
	out := flag.String("out", "", "出力 PDF ファイル名 (未指定なら <dir 名>.pdf)")
	flag.Parse()

	images, output, err := GetSourceImages(*dir, *out)
	if err != nil {
		check(err)
	}
	if images != nil {
		check(err)
	}

	fmt.Printf("検出した画像: %d 枚\n", len(images))
	fmt.Printf("出力 PDF   : %s\n", *out)

	// ---- PDF generating ----
	err = MergeImagesIntoPDF(images, output)
	if err != nil {
		check(err)
	}

	fmt.Println("PDF を生成しました。完了です。")
}
