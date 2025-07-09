package usecases

import (
	"os"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func createTestPDF(t *testing.T) string {
	tmpfile, err := os.CreateTemp("", "test*.pdf")
	if err != nil {
		t.Fatalf("一時ファイルの作成に失敗: %v", err)
	}
	tmpfile.Close()

	// デフォルト設定でコンテキストを作成
	cfg := model.NewDefaultConfiguration()

	// ページサイズを指定してコンテキストを作成（A4サイズ）
	ctx, err := pdfcpu.CreateContextWithXRefTable(cfg, types.PaperSize["A4"])
	if err != nil {
		os.Remove(tmpfile.Name())
		t.Fatalf("XRefテーブルコンテキストの作成に失敗: %v", err)
	}

	// PDFをファイルに書き込む
	f, err := os.Create(tmpfile.Name())
	if err != nil {
		os.Remove(tmpfile.Name())
		t.Fatalf("ファイルの作成に失敗: %v", err)
	}
	defer f.Close()

	err = api.WriteContext(ctx, f)
	if err != nil {
		os.Remove(tmpfile.Name())
		t.Fatalf("PDFコンテキストの書き込みに失敗: %v", err)
	}

	return tmpfile.Name()
}

func TestEncryptPDF(t *testing.T) {
	tests := []struct {
		name          string
		setupPDF      func() string
		inputPath     string
		outputPath    string
		userPassword  string
		ownerPassword string
		expectError   bool
		expectedError string
	}{
		{
			name: "正常な暗号化",
			setupPDF: func() string {
				return createTestPDF(t)
			},
			outputPath:    "",  // 空文字は入力ファイルを上書き
			userPassword:  "userPass123",
			ownerPassword: "ownerPass123",
			expectError:   false,
		},
		{
			name: "存在しないファイルの暗号化",
			setupPDF: func() string {
				return "/non/existent/file.pdf"
			},
			outputPath:    "output.pdf",
			userPassword:  "userPass123",
			ownerPassword: "ownerPass123",
			expectError:   true,
		},
		{
			name: "オーナーパスワードのみの暗号化",
			setupPDF: func() string {
				return createTestPDF(t)
			},
			outputPath:    "",
			userPassword:  "",
			ownerPassword: "ownerPass123",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputFile := tt.setupPDF()
			if inputFile != "/non/existent/file.pdf" {
				defer os.Remove(inputFile)
			}

			upw := tt.userPassword
			opw := tt.ownerPassword
			err := EncryptPDF(inputFile, tt.outputPath, &upw, &opw)

			if tt.expectError && err == nil {
				t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
			}
			if !tt.expectError && err != nil {
				t.Errorf("エラーが発生しました: %v", err)
			}

			// 暗号化が成功した場合、PDFが実際に暗号化されているか確認
			if !tt.expectError && err == nil {
				cfg := model.NewDefaultConfiguration()
				cfg.UserPW = tt.userPassword
				cfg.OwnerPW = tt.ownerPassword

				outputFile := tt.outputPath
				if outputFile == "" {
					outputFile = inputFile
				}

				// 正しいパスワードで復号化を試す
				if err := api.DecryptFile(outputFile, "", cfg); err != nil {
					t.Errorf("暗号化されたPDFの復号化に失敗: %v", err)
				}
			}
		})
	}
}

func TestDecryptPDF(t *testing.T) {
	tests := []struct {
		name          string
		setupPDF      func() (string, string)
		inputPath     string
		outputPath    string
		ownerPassword string
		expectError   bool
	}{
		{
			name: "正常な復号化",
			setupPDF: func() (string, string) {
				// テスト用のPDFを作成して暗号化
				originalFile := createTestPDF(t)

				password := "testPassword123"
				encryptCfg := model.NewAESConfiguration(password, password, 256)
				encryptCfg.OptimizeResourceDicts = true
				encryptCfg.OptimizeDuplicateContentStreams = true
				encryptCfg.UserPWNew = &password
				encryptCfg.OwnerPWNew = &password

				encryptedFile := originalFile + ".encrypted"
				if err := api.EncryptFile(originalFile, encryptedFile, encryptCfg); err != nil {
					t.Fatalf("テストPDFの暗号化に失敗: %v", err)
				}

				os.Remove(originalFile) // オリジナルは削除
				return encryptedFile, password
			},
			outputPath:  "",
			expectError: false,
		},
		{
			name: "間違ったパスワードで復号化",
			setupPDF: func() (string, string) {
				originalFile := createTestPDF(t)

				password := "testPassword123"
				encryptCfg := model.NewAESConfiguration(password, password, 256)
				encryptedFile := originalFile + ".encrypted"
				if err := api.EncryptFile(originalFile, encryptedFile, encryptCfg); err != nil {
					t.Fatalf("テストPDFの暗号化に失敗: %v", err)
				}

				os.Remove(originalFile) // オリジナルは削除
				return encryptedFile, "wrongPassword"
			},
			outputPath:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputFile, password := tt.setupPDF()
			defer os.Remove(inputFile)

			opw := password
			err := DecryptPDF(inputFile, tt.outputPath, &opw)

			if tt.expectError && err == nil {
				t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
			}
			if !tt.expectError && err != nil {
				t.Errorf("エラーが発生しました: %v", err)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	// Checkのテストは出力とexit処理を検証する必要があるため、
	// 実際のコマンドラインテストや統合テストで行うことが適切
	t.Run("エラーハンドリング", func(t *testing.T) {
		// これは実際にos.Exit(1)を呼び出すため、テストは困難
		// 必要に応じてモックを使用するか、外部コマンドでテストする
		t.Skip("Check関数は直接のユニットテストが難しいためスキップ")
	})
}

func TestIntegration(t *testing.T) {
	t.Run("暗号化→復号化フロー", func(t *testing.T) {
		// 一時PDFファイルを作成
		tmpfile := createTestPDF(t)
		defer os.Remove(tmpfile)

		userPassword := "user123"
		ownerPassword := "owner123"

		// 暗号化
		encryptedFile := tmpfile + ".encrypted"
		err := EncryptPDF(tmpfile, encryptedFile, &userPassword, &ownerPassword)
		if err != nil {
			t.Fatalf("暗号化に失敗: %v", err)
		}
		defer os.Remove(encryptedFile)

		// 復号化
		decryptedFile := encryptedFile + ".decrypted"
		err = DecryptPDF(encryptedFile, decryptedFile, &ownerPassword)
		if err != nil {
			t.Fatalf("復号化に失敗: %v", err)
		}
		defer os.Remove(decryptedFile)

		// 復号化したファイルが存在することを確認
		if _, err := os.Stat(decryptedFile); os.IsNotExist(err) {
			t.Errorf("復号化したファイルが存在しません")
		}
	})
}
