package config

import (
	"fmt"
	"runtime"
	"strings"

	flagParser "github.com/landmaster135/devbox/internal/image_converter_by_libwebp/infrastructures/flag_parser"
)

const (
	defaultSrcDir  = "."
	defaultOutDir  = "./999_converted_images"
	defaultOutExt  = "webp"
	defaultQuality = 99
	minQuality     = 1
	maxQuality     = 100
)

// Config は image-converter-by-libwebp CLI の設定を保持します。
type Config struct {
	SrcDir     string
	OutDir     string
	ArchiveDir string
	Move       bool
	OutExt     string
	Quality    int
	Workers    int
	Recursive  bool
	Lossless   bool
	Help       bool
}

// NewConfig は値検証済みの Config を作成します。
func NewConfig(srcDir, outDir, archiveDir, outExt string, move bool, quality, workers int, recursive, lossless bool) (*Config, error) {
	normalizedOutExt := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(outExt), "."))
	if normalizedOutExt == "" {
		return nil, fmt.Errorf("出力形式が指定されていません")
	}
	if normalizedOutExt != defaultOutExt {
		return nil, fmt.Errorf("サポートされていない出力フォーマットです: %s", normalizedOutExt)
	}
	if strings.TrimSpace(srcDir) == "" {
		return nil, fmt.Errorf("入力ディレクトリが指定されていません")
	}
	if strings.TrimSpace(outDir) == "" {
		return nil, fmt.Errorf("出力ディレクトリが指定されていません")
	}
	if quality < minQuality || quality > maxQuality {
		return nil, fmt.Errorf("品質は%dから%dの範囲で指定してください: %d", minQuality, maxQuality, quality)
	}
	if workers < 1 {
		return nil, fmt.Errorf("workers は1以上で指定してください: %d", workers)
	}

	return &Config{
		SrcDir:     srcDir,
		OutDir:     outDir,
		ArchiveDir: archiveDir,
		Move:       move,
		OutExt:     normalizedOutExt,
		Quality:    quality,
		Workers:    workers,
		Recursive:  recursive,
		Lossless:   lossless,
	}, nil
}

// ParseFlags はコマンドライン引数から Config を作成します。
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(flagParser.NewStandardFlagParser())
}

// ParseFlagsWithArgs は指定された引数から Config を作成します。
func ParseFlagsWithArgs(args []string) (*Config, error) {
	return ParseFlagsWithParser(flagParser.NewStandardFlagParserWithArgs(args))
}

// ParseFlagsWithParser は指定された parser から Config を作成します。
func ParseFlagsWithParser(parser flagParser.FlagParser) (*Config, error) {
	srcDir := defaultSrcDir
	outDir := defaultOutDir
	archiveDir := ""
	move := false
	outExt := defaultOutExt
	quality := defaultQuality
	workers := runtime.NumCPU()
	recursive := false
	lossless := false
	help := false

	parser.StringVar(&srcDir, "src-dir", srcDir, "入力画像を探索するディレクトリ")
	parser.StringVar(&outDir, "out-dir", outDir, "変換後画像の出力ディレクトリ")
	parser.StringVar(&archiveDir, "archive-dir", archiveDir, "処理済み元ファイルの退避先ディレクトリ")
	parser.BoolVar(&move, "move", move, "退避時にコピーではなく移動する")
	parser.StringVar(&outExt, "ext", outExt, "出力形式 (webp)")
	parser.IntVar(&quality, "q", quality, "WebP 品質 (1-100)")
	parser.IntVar(&workers, "workers", workers, "並列ワーカー数")
	parser.BoolVar(&recursive, "recursive", recursive, "サブディレクトリを再帰的に走査する")
	parser.BoolVar(&recursive, "R", recursive, "サブディレクトリを再帰的に走査する")
	parser.BoolVar(&lossless, "lossless", lossless, "cwebp の lossless 圧縮を有効にする")
	parser.BoolVar(&help, "help", help, "ヘルプを表示する")
	parser.BoolVar(&help, "h", help, "ヘルプを表示する")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}
	if help {
		return &Config{Help: true}, nil
	}

	return NewConfig(srcDir, outDir, archiveDir, outExt, move, quality, workers, recursive, lossless)
}

// PrintUsage は usage を標準エラーへ出力します。
func PrintUsage() {
	flagParser.PrintUsage(usageTemplate)
}
