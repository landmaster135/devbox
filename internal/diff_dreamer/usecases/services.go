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
type DiffConfig struct {
	LeftText  string
	RightText string
}

// HTMLテンプレートのデータ
type TemplateData struct {
	Style    string
	Script   string
	ConfigJS string
}

// ReadTextFile は指定したファイルパスからテキストを読み取ります
func ReadTextFile(filePath string) (string, error) {
	if filePath == "" {
		return "", nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("ファイルの読み取りに失敗しました: %v", err)
	}

	return string(data), nil
}

// GenerateConfigJS は設定情報をJavaScript形式に変換します
func GenerateConfigJS(config DiffConfig) string {
	var buf bytes.Buffer
	buf.WriteString("// 自動生成された設定\nconst CONFIG = {\n")
	buf.WriteString(fmt.Sprintf("  leftText: %s,\n", escapeJSString(config.LeftText)))
	buf.WriteString(fmt.Sprintf("  rightText: %s\n", escapeJSString(config.RightText)))
	buf.WriteString("};")
	return buf.String()
}

// escapeJSString はJavaScript文字列リテラル用にエスケープします
func escapeJSString(s string) string {
	if s == "" {
		return `""`
	}

	// 改行やタブ、引用符をエスケープ
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")

	return fmt.Sprintf(`"%s"`, s)
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

// OpenHTMLPage はHTMLファイルを作成してブラウザで開きます
func OpenHTMLPage(html, filePath string) error {
	htmlFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("HTMLファイルの作成に失敗しました: %v", err)
	}

	// HTMLをファイルに書き込む
	if _, err := htmlFile.Write([]byte(html)); err != nil {
		htmlFile.Close()
		return fmt.Errorf("HTMLファイルへの書き込みに失敗しました: %v", err)
	}
	htmlFile.Close()

	// ファイルのパスを表示
	absPath, _ := filepath.Abs(filePath)
	fmt.Fprintf(os.Stdout, "HTMLファイルを作成しました: %s\n", absPath)

	// ブラウザでHTMLファイルを開く
	fmt.Fprintf(os.Stdout, "Diff Dreamerを起動しています...\n")

	// URLを構築
	fileURL := "file://" + absPath
	fmt.Fprintf(os.Stdout, "開こうとしているURL: %s\n", fileURL)

	// ブラウザを開く
	if err := openBrowser(fileURL); err != nil {
		return fmt.Errorf("ブラウザの起動に失敗しました: %v", err)
	}

	return nil
}

// CopyFile はファイルをコピーする関数
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
