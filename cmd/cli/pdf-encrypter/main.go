package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	usecases "github.com/landmaster135/devbox/internal/pdf_encrypter/usecases"
)

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

// run はPDF Encrypterツールの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグセットを作成
	fs := flag.NewFlagSet("pdf-encrypter", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// コマンドライン引数の定義
	mode := fs.String("mode", "encrypt", "encrypt / decrypt")
	in := fs.String("in", "", "入力 PDF ※必須")
	out := fs.String("out", "", "出力 PDF (未指定なら上書き)")
	upw := fs.String("upw", "", "ユーザーパスワード (閲覧用)")
	opw := fs.String("opw", "", "オーナーパスワード (管理用) ※必須")

	// 引数の解析
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// 入力ファイルが指定されていない場合はエラー
	if *in == "" {
		fmt.Fprintln(stderr, "エラー: 入力ファイル (-in) は必須です")
		fs.Usage()
		return exitCodeError
	}

	switch *mode {
	case "encrypt":
		// 暗号化モードの場合、オーナーパスワードは必須
		if *opw == "" {
			fmt.Fprintln(stderr, "エラー: 暗号化モードでは -opw (オーナーパスワード) が必須です")
			fs.Usage()
			return exitCodeError
		}

		// PDFファイルを暗号化
		if err := usecases.EncryptPDF(*in, *out, upw, opw); err != nil {
			fmt.Fprintf(stderr, "エラー: %v\n", err)
			return exitCodeError
		}

		fmt.Fprintf(stdout, "✔ PDFファイル '%s' の暗号化が完了しました\n", *in)

	case "decrypt":
		// PDFファイルを復号化
		if err := usecases.DecryptPDF(*in, *out, opw); err != nil {
			fmt.Fprintf(stderr, "エラー: %v\n", err)
			return exitCodeError
		}

		fmt.Fprintf(stdout, "✔ PDFファイル '%s' の復号化が完了しました\n", *in)

	default:
		fmt.Fprintf(stderr, "エラー: -mode は encrypt または decrypt を指定してください\n")
		fs.Usage()
		return exitCodeError
	}

	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
