package usecases

import (
	"fmt"
	"os"

	api "github.com/pdfcpu/pdfcpu/pkg/api"
	model "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func EncryptPDF(in string, out string, upw *string, opw *string) error {
	// AES 暗号化の設定を生成
	key := 256
	cfg := model.NewAESConfiguration(*upw, *opw, key)
	// Remove unused fonts and images from resource dictionary
	cfg.OptimizeResourceDicts = true
	// Share duplicated streams in all pages
	cfg.OptimizeDuplicateContentStreams = true
	// Set user password and owner password
	cfg.UserPWNew = upw
	cfg.OwnerPWNew = opw

	// EncryptFile(in, out, conf)
	// out を空文字にすると入力ファイルを上書き
	if err := api.EncryptFile(in, out, cfg); err != nil {
		return err
	}

	return nil
}

func DecryptPDF(in string, out string, opw *string) error {
	// AES 暗号化の設定を生成
	cfg := model.NewDefaultConfiguration()
	// Remove unused fonts and images from resource dictionary
	cfg.OptimizeResourceDicts = true
	// Share duplicated streams in all pages
	cfg.OptimizeDuplicateContentStreams = true
	// Set user password and owner password
	cfg.UserPW = *opw
	cfg.OwnerPW = *opw

	// EncryptFile(in, out, conf)
	// out を空文字にすると入力ファイルを上書き
	if err := api.DecryptFile(in, out, cfg); err != nil {
		return err
	}

	return nil
}

func Check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}
}
