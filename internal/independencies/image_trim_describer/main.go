package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

//go:embed web/index.html
var indexHTML string

//go:embed web/style.css
var styleCSS string

//go:embed web/script.js
var scriptJS string

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
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

// run は画像トリミング記述ツールの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグセットを作成
	flagSet := flag.NewFlagSet("image-trim-describer", flag.ContinueOnError)
	flagSet.SetOutput(stderr)

	// コマンドライン引数の定義
	srcDir := flagSet.String("src", ".", "画像ファイルのディレクトリ")
	outDir := flagSet.String("out", "", "出力先ディレクトリ (空の場合はsrcと同じ)")
	arcDir := flagSet.String("arc", "./5_original_files", "アーカイブ先ディレクトリ")
	suffix := flagSet.String("suffix", "trimmed", "出力ファイル名に付加するサフィックス")
	move := flagSet.Bool("move", false, "元ファイルを移動する (コピーではなく)")

	// 引数の解析
	if err := flagSet.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// 出力ディレクトリが指定されていない場合は入力ディレクトリと同じに
	if *outDir == "" {
		*outDir = *srcDir
	}

	// 画像ファイル一覧を取得
	imageFiles, err := getImageFiles(*srcDir)
	if err != nil {
		fmt.Fprintf(stderr, "画像ファイル一覧の取得に失敗しました: %v\n", err)
		return exitCodeError
	}

	// 設定情報を作成
	config := Config{
		ImageFiles: imageFiles,
		SrcDir:     *srcDir,
		OutDir:     *outDir,
		ArcDir:     *arcDir,
		Suffix:     *suffix,
		Move:       *move,
	}

	// 設定情報をJavaScript形式に変換
	configJS := generateConfigJS(config)

	// HTMLテンプレートのデータを作成
	data := TemplateData{
		Style:      styleCSS,
		Script:     scriptJS,
		ConfigJS:   configJS,
		ImageFiles: imageFiles,
	}

	// HTMLを生成
	html, err := generateHTML(indexHTML, data)
	if err != nil {
		fmt.Fprintf(stderr, "HTMLの生成に失敗しました: %v\n", err)
		return exitCodeError
	}

	// 画像ファイル用のディレクトリを作成
	imagesDir := "images"
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		fmt.Fprintf(stderr, "画像ディレクトリの作成に失敗しました: %v\n", err)
		return exitCodeError
	}

	// 画像ファイルをコピー
	for _, imgFile := range imageFiles {
		srcPath := filepath.Join(*srcDir, imgFile)
		dstPath := filepath.Join(imagesDir, imgFile)

		// ファイルをコピー
		if err := copyFile(srcPath, dstPath); err != nil {
			fmt.Fprintf(stderr, "画像ファイルのコピーに失敗しました: %v\n", err)
			// エラーがあっても続行
		}
	}

	// カレントディレクトリにHTMLファイルを作成
	htmlFilePath := "image_trimmer.html"
	htmlFile, err := os.Create(htmlFilePath)
	if err != nil {
		fmt.Fprintf(stderr, "HTMLファイルの作成に失敗しました: %v\n", err)
		return exitCodeError
	}

	// HTMLをファイルに書き込む
	if _, err := htmlFile.Write([]byte(html)); err != nil {
		fmt.Fprintf(stderr, "HTMLファイルへの書き込みに失敗しました: %v\n", err)
		htmlFile.Close()
		return exitCodeError
	}
	htmlFile.Close()

	// ファイルのパスを表示
	absPath, _ := filepath.Abs(htmlFilePath)
	fmt.Fprintf(stdout, "HTMLファイルを作成しました: %s\n", absPath)

	// ブラウザでHTMLファイルを開く
	fmt.Fprintf(stdout, "画像トリミングツールを起動しています...\n")

	// URLを構築
	fileURL := "file://" + absPath
	fmt.Fprintf(stdout, "開こうとしているURL: %s\n", fileURL)

	// ブラウザを開く
	if err := openBrowser(fileURL); err != nil {
		fmt.Fprintf(stderr, "ブラウザの起動に失敗しました: %v\n", err)
		return exitCodeError
	}

	fmt.Fprintf(stdout, "画像トリミングツールが起動しました\n")
	fmt.Fprintf(stdout, "ブラウザでトリミング座標を選択し、トリミングを実行してください\n")

	// ユーザーが確認できるように一時停止
	fmt.Fprintf(stdout, "Enterキーを押すとプログラムを終了します（HTMLファイルは削除されません）...\n")
	fmt.Scanln()

	return exitCodeOK
}

// getImageFiles は指定したディレクトリ内の画像ファイル一覧を取得します
func getImageFiles(dirPath string) ([]string, error) {
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

// generateConfigJS は設定情報をJavaScript形式に変換します
func generateConfigJS(config Config) string {
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

// generateHTML はHTMLテンプレートを生成します
func generateHTML(tmplStr string, data TemplateData) (string, error) {
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

// ファイルをコピーする関数
func copyFile(src, dst string) error {
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

// 画像をトリミングして保存する関数
func cropAndSave(inPath, outDir string, x1, y1, x2, y2 int, suffix string) error {
	// 画像を読み込み
	img, err := imgio.Open(inPath)
	if err != nil {
		return fmt.Errorf("open %s: %w", inPath, err)
	}

	// 矩形チェック
	bounds := img.Bounds()
	if x1 < 0 || y1 < 0 || x2 <= x1 || y2 <= y1 ||
		x2 > bounds.Dx() || y2 > bounds.Dy() {
		return fmt.Errorf("invalid rectangle %v for %s",
			image.Rect(x1, y1, x2, y2), inPath)
	}

	// トリミング
	cropped := transform.Crop(img, image.Rect(x1, y1, x2, y2))

	// 保存パス準備
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	base := strings.TrimSuffix(filepath.Base(inPath), filepath.Ext(inPath))
	outPath := filepath.Join(outDir,
		fmt.Sprintf("%s_%s%s", base, suffix, filepath.Ext(inPath)))

	// 形式に応じて保存
	ext := strings.ToLower(filepath.Ext(outPath))
	switch ext {
	case ".jpg", ".jpeg":
		err = imgio.Save(outPath, cropped, imgio.JPEGEncoder(95))
	case ".png":
		err = imgio.Save(outPath, cropped, imgio.PNGEncoder())
	default:
		err = fmt.Errorf("unsupported extension: %s", ext)
	}
	return err
}

// 元画像を移動する関数
func moveOriginal(src, arcDir string) error {
	if err := os.MkdirAll(arcDir, 0o755); err != nil {
		return err
	}
	dst := filepath.Join(arcDir, filepath.Base(src))
	return os.Rename(src, dst)
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
