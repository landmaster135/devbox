// pdfcrypt.go
//
// ビルド&実行例:
//   go mod init example.com/pdfcrypt
//   go get github.com/pdfcpu/pdfcpu/pkg/api
//   go build -o pdfcrypt.exe pdfcrypt.go
//
//   # 暗号化 (AES-256):
//   pdfcrypt -mode encrypt -in plain.pdf -out locked.pdf -upw reader -opw owner
//
//   # 復号化 (パスワード解除):
//   pdfcrypt -mode decrypt -in locked.pdf -out plain.pdf -pw reader   # ←ユーザPW
//   pdfcrypt -mode decrypt -in locked.pdf -out plain.pdf -pw ownerPW  # ←オーナPW
//
package main

import (
	"errors"
	"flag"
	"log"

	usecases "github.com/landmaster135/devbox/internal/independencies/pdf_encrypter/usecases"
)

func main() {
	mode := flag.String("mode", "encrypt", "encrypt / decrypt")
	in := flag.String("in", "", "入力 PDF ※必須")
	out := flag.String("out", "", "出力 PDF (未指定なら上書き)")
	upw := flag.String("upw", "", "ユーザーパスワード (閲覧用)")
	opw := flag.String("opw", "", "オーナーパスワード (管理用) ※必須")

	flag.Parse()
	if *in == "" {
		flag.Usage()
		usecases.Check(errors.New("入力ファイル (-in) とオーナーパスワード (-opw) は必須です"))
	}
	switch *mode {
	case "encrypt":
		if *opw == "" {
			log.Fatal("-opw (オーナーパスワード) が必須です")
		}
		if err := usecases.EncryptPDF(*in, *out, upw, opw); err != nil {
			usecases.Check(err)
		}
		log.Println("✔ 暗号化が完了しました")
	case "decrypt":
		if err := usecases.DecryptPDF(*in, *out, opw); err != nil {
			usecases.Check(err)
		}
		log.Println("✔ 復号化が完了しました")
	default:
		usecases.Check(errors.New("-mode は encrypt または decrypt を指定してください"))
	}


}
