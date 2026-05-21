package config

import (
	"fmt"
	"regexp"
	"strings"

	flagParser "github.com/landmaster135/devbox/internal/movie_extractor/infrastructures/flag_parser"
)

const (
	defaultOperation = "extract-frames"
	defaultFPS       = 2
	defaultQuality   = 2
)

var supportedOperations = map[string]struct{}{
	"extract-frames": {},
}

var (
	startPositionSecondPattern = regexp.MustCompile(`^\d+(\.\d+)?$`)
	startPositionClockPattern  = regexp.MustCompile(`^\d{1,2}:\d{2}(:\d{2}(\.\d+)?)?$`)
)

// Config は movie-extractor CLI の設定です。
type Config struct {
	Operation     string
	SrcFile       string
	FPS           int
	Quality       int
	StartPosition string
	OutDir        string
	Help          bool
}

// NewConfig は設定を検証して Config を返します。
func NewConfig(
	operation string,
	srcFile string,
	fps int,
	quality int,
	startPosition string,
	outDir string,
) (*Config, error) {
	if _, ok := supportedOperations[operation]; !ok {
		return nil, fmt.Errorf("未対応のoperationです: %s", operation)
	}
	if srcFile == "" {
		return nil, fmt.Errorf("src-file は必須です")
	}
	if fps <= 0 {
		return nil, fmt.Errorf("fps は1以上の整数を指定してください")
	}
	if quality < 1 || quality > 31 {
		return nil, fmt.Errorf("quality は1から31の範囲で指定してください")
	}
	if !isValidStartPosition(startPosition) {
		return nil, fmt.Errorf("start-position の形式が不正です: %s", startPosition)
	}
	if outDir == "" {
		return nil, fmt.Errorf("out-dir は必須です")
	}

	return &Config{
		Operation:     operation,
		SrcFile:       srcFile,
		FPS:           fps,
		Quality:       quality,
		StartPosition: startPosition,
		OutDir:        outDir,
	}, nil
}

// ParseFlags は標準 parser を使って設定を解析します。
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(flagParser.NewStandardFlagParser())
}

// ParseFlagsWithParser は指定 parser で設定を解析します。
func ParseFlagsWithParser(parser flagParser.FlagParser) (*Config, error) {
	var (
		operation     = defaultOperation
		srcFile       = ""
		fps           = defaultFPS
		quality       = defaultQuality
		startPosition = ""
		outDir        = ""
		help          = false
	)

	parser.StringVar(&operation, "operation", operation, "実行する操作。extract-frames のみ対応")
	parser.StringVar(&srcFile, "src-file", srcFile, "入力動画ファイルのパス")
	parser.IntVar(&fps, "fps", fps, "1秒あたりに抽出するフレーム数")
	parser.IntVar(&quality, "quality", quality, "JPEG品質(1-31)。小さいほど高品質")
	parser.StringVar(&startPosition, "start-position", startPosition, "抽出開始位置。秒数または HH:MM:SS[.ms]")
	parser.StringVar(&outDir, "out-dir", outDir, "抽出画像の出力先ディレクトリ")
	parser.BoolVar(&help, "help", help, "ヘルプを表示")
	parser.BoolVar(&help, "h", help, "ヘルプを表示(短縮)")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグ解析に失敗しました: %w", err)
	}
	if help {
		return &Config{Help: true}, nil
	}

	return NewConfig(operation, srcFile, fps, quality, startPosition, outDir)
}

// PrintUsage は CLI 利用方法を表示します。
func PrintUsage() {
	flagParser.PrintUsage(`Movie Extractor CLI

使用方法:
  %[1]s -operation extract-frames -src-file <video file> -out-dir <directory> [options]

必須フラグ:
  -src-file string
      入力動画ファイルのパス
  -out-dir string
      出力ディレクトリ

任意フラグ:
  -operation string
      実行操作（デフォルト: extract-frames）
  -fps int
      抽出フレームレート（デフォルト: 1）
  -quality int
      JPEG品質（1-31, 小さいほど高品質。デフォルト: 2）
  -start-position string
      抽出開始位置（秒数または HH:MM:SS[.ms]）
  -help, -h
      ヘルプを表示

例:
  %[1]s -operation extract-frames -src-file ./movie.mp4 -out-dir ./frames
  %[1]s -operation extract-frames -src-file ./movie.mp4 -fps 5 -quality 3 -out-dir ./frames
  %[1]s -operation extract-frames -src-file ./movie.mp4 -start-position 00:00:10.5 -out-dir ./frames
`)
}

func isValidStartPosition(value string) bool {
	if value == "" {
		return true
	}
	if startPositionSecondPattern.MatchString(value) {
		return true
	}
	if !startPositionClockPattern.MatchString(value) {
		return false
	}

	parts := strings.Split(value, ":")
	if len(parts) < 2 || len(parts) > 3 {
		return false
	}

	minutePart := parts[len(parts)-2]
	secondPart := parts[len(parts)-1]

	minutes := parseTwoDigits(minutePart)
	if minutes < 0 || minutes > 59 {
		return false
	}

	seconds := parseSecondValue(secondPart)
	return seconds >= 0 && seconds < 60
}

func parseTwoDigits(value string) int {
	if len(value) != 2 {
		return -1
	}
	if value[0] < '0' || value[0] > '9' || value[1] < '0' || value[1] > '9' {
		return -1
	}
	return int(value[0]-'0')*10 + int(value[1]-'0')
}

func parseSecondValue(value string) float64 {
	whole := value
	fraction := ""

	if idx := strings.Index(value, "."); idx >= 0 {
		whole = value[:idx]
		fraction = value[idx+1:]
	}

	if len(whole) != 2 {
		return -1
	}
	secInt := parseTwoDigits(whole)
	if secInt < 0 {
		return -1
	}

	result := float64(secInt)
	if fraction == "" {
		return result
	}

	multiplier := 0.1
	for _, ch := range fraction {
		if ch < '0' || ch > '9' {
			return -1
		}
		result += float64(ch-'0') * multiplier
		multiplier *= 0.1
	}
	return result
}
