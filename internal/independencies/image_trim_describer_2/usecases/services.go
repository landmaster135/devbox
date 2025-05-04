package usecases

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// 設定情報を保持する構造体
type Config struct {
	ImageFiles []string
	SrcDir     string
	OutDir     string
	ArcDir     string
	Suffix     string
	Move       bool
}

// HTMLテンプレートのデータ
type TemplateData struct {
	Style      string
	Script     string
	ConfigJS   string
	ImageFiles []string
}

// GetImageFiles は指定したディレクトリ内の画像ファイル一覧を取得します
func GetImageFiles(dirPath string) ([]string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var imageFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
			imageFiles = append(imageFiles, file.Name())
		}
	}

	return imageFiles, nil
}

// formatJSArray は文字列スライスをJavaScript配列形式にフォーマットします
func formatJSArray(items []string) string {
	if len(items) == 0 {
		return "[]"
	}

	quotedItems := make([]string, len(items))
	for i, item := range items {
		quotedItems[i] = fmt.Sprintf("%q", item)
	}

	return fmt.Sprintf("[\n    %s\n  ]", strings.Join(quotedItems, ",\n    "))
}

// GenerateConfigJS は設定情報をJavaScript形式に変換します
func GenerateConfigJS(config Config) string {
	var buf bytes.Buffer
	buf.WriteString("// 自動生成された設定\nconst CONFIG = {\n")
	buf.WriteString(fmt.Sprintf("  imageFiles: %s,\n", formatJSArray(config.ImageFiles)))
	buf.WriteString(fmt.Sprintf("  srcDir: %q,\n", config.SrcDir))
	buf.WriteString(fmt.Sprintf("  outDir: %q,\n", config.OutDir))
	buf.WriteString(fmt.Sprintf("  arcDir: %q,\n", config.ArcDir))
	buf.WriteString(fmt.Sprintf("  suffix: %q,\n", config.Suffix))
	buf.WriteString(fmt.Sprintf("  move: %t\n", config.Move))
	buf.WriteString("};")
	return buf.String()
}

// GenerateHTML はHTMLテンプレートを生成します
func GenerateHTML(tmplStr string, data TemplateData) (string, error) {
	// 直接文字列置換を使用してテンプレートを処理
	html := tmplStr

	// スタイルシートの置き換え
	html = strings.Replace(html,
		`<link rel="stylesheet" href="style.css">`,
		fmt.Sprintf("<style>\n%s\n</style>", data.Style),
		1)

	// スクリプトの置き換え
	html = strings.Replace(html,
		`<script src="script.js"></script>`,
		fmt.Sprintf("<script>\n%s\n</script>", data.Script),
		1)

	// 設定スクリプトの追加
	html = strings.Replace(html,
		`<script src="config.js"></script>`,
		fmt.Sprintf("<script>\n%s\n</script>", data.ConfigJS),
		1)

	return html, nil
}

// openBrowser はURLまたはファイルパスをデフォルトブラウザで開きます
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("このプラットフォームではブラウザを自動的に開くことができません: %s", runtime.GOOS)
	}

	return cmd.Start()
}

func OpenHTMLPage(html, filePath string) error {
	htmlFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("HTMLファイルの作成に失敗しました: %v\n", err)
	}

	// HTMLをファイルに書き込む
	if _, err := htmlFile.Write([]byte(html)); err != nil {
		htmlFile.Close()
		return fmt.Errorf("HTMLファイルへの書き込みに失敗しました: %v\n", err)
	}
	htmlFile.Close()

	// ファイルのパスを表示
	absPath, _ := filepath.Abs(filePath)
	fmt.Fprintf(os.Stdout, "HTMLファイルを作成しました: %s\n", absPath)

	// ブラウザでHTMLファイルを開く
	fmt.Fprintf(os.Stdout, "画像トリミングツールを起動しています...\n")

	// URLを構築
	fileURL := "file://" + absPath
	fmt.Fprintf(os.Stdout, "開こうとしているURL: %s\n", fileURL)

	// ブラウザを開く
	if err := openBrowser(fileURL); err != nil {
		return fmt.Errorf("ブラウザの起動に失敗しました: %v\n", err)
	}

	return nil
}

// ファイルをコピーする関数
func CopyFile(src, dst string) error {
	// 元ファイルを開く
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 宛先ファイルを作成
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// コピー
	_, err = io.Copy(dstFile, srcFile)
	return err
}
