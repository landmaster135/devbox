package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Config はYouTube動画ダウンローダーの設定を保持します
type Config struct {
	URL         string
	OutputDir   string
	Quality     string
	Format      string
	AudioOnly   bool
	Playlist    bool
	MaxRoutines int
	ChunkSize   int64
	Help        bool
}

// FlagParser はフラグ解析のインターフェースです
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	IntVar(p *int, name string, value int, usage string)
	Int64Var(p *int64, name string, value int64, usage string)
	Parse()
	NArg() int
}

// DefaultFlagParser は標準のフラグパーサーです
type DefaultFlagParser struct{}

func (d *DefaultFlagParser) StringVar(p *string, name string, value string, usage string) {
	flag.StringVar(p, name, value, usage)
}

func (d *DefaultFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	flag.BoolVar(p, name, value, usage)
}

func (d *DefaultFlagParser) IntVar(p *int, name string, value int, usage string) {
	flag.IntVar(p, name, value, usage)
}

func (d *DefaultFlagParser) Int64Var(p *int64, name string, value int64, usage string) {
	flag.Int64Var(p, name, value, usage)
}

func (d *DefaultFlagParser) Parse() {
	flag.Parse()
}

func (d *DefaultFlagParser) NArg() int {
	return flag.NArg()
}

// ParseFlags はコマンドライン引数を解析して設定を返します
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(&DefaultFlagParser{})
}

// ParseFlagsWithParser は指定されたパーサーでフラグを解析します
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	cfg := &Config{}

	parser.StringVar(&cfg.URL, "url", "", "YouTube動画またはプレイリストのURL（必須）")
	parser.StringVar(&cfg.OutputDir, "output", "./downloads", "ダウンロード先ディレクトリ")
	parser.StringVar(&cfg.Quality, "quality", "best", "動画品質（best, worst, 720p, 1080p等）")
	parser.StringVar(&cfg.Format, "format", "mp4", "動画形式（mp4, webm等）")
	parser.BoolVar(&cfg.AudioOnly, "audio-only", false, "音声のみダウンロード")
	parser.BoolVar(&cfg.Playlist, "playlist", false, "プレイリスト全体をダウンロード")
	parser.IntVar(&cfg.MaxRoutines, "max-routines", 10, "並列ダウンロード数")
	parser.Int64Var(&cfg.ChunkSize, "chunk-size", 10*1024*1024, "チャンクサイズ（バイト）")
	parser.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")

	parser.Parse()

	if cfg.Help {
		return cfg, nil
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate は設定の妥当性を検証します
func (c *Config) validate() error {
	if c.URL == "" {
		return fmt.Errorf("URLが指定されていません。-url フラグでYouTube URLを指定してください")
	}

	if c.MaxRoutines <= 0 {
		return fmt.Errorf("並列ダウンロード数は1以上である必要があります: %d", c.MaxRoutines)
	}

	if c.ChunkSize <= 0 {
		return fmt.Errorf("チャンクサイズは1以上である必要があります: %d", c.ChunkSize)
	}

	// 出力ディレクトリの絶対パス化
	absPath, err := filepath.Abs(c.OutputDir)
	if err != nil {
		return fmt.Errorf("出力ディレクトリパスの変換に失敗しました: %v", err)
	}
	c.OutputDir = absPath

	return nil
}

// PrintUsage は使用方法を表示します
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "YouTube動画ダウンローダー\n\n")
	fmt.Fprintf(os.Stderr, "使用方法:\n")
	fmt.Fprintf(os.Stderr, "  youtube-downloader -url <YouTube URL> [オプション]\n\n")
	fmt.Fprintf(os.Stderr, "例:\n")
	fmt.Fprintf(os.Stderr, "  # 基本的な使用方法\n")
	fmt.Fprintf(os.Stderr, "  youtube-downloader -url \"https://www.youtube.com/watch?v=VIDEO_ID\"\n\n")
	fmt.Fprintf(os.Stderr, "  # 品質と出力先を指定\n")
	fmt.Fprintf(os.Stderr, "  youtube-downloader -url \"URL\" -quality \"720p\" -output \"./videos\"\n\n")
	fmt.Fprintf(os.Stderr, "  # 音声のみダウンロード\n")
	fmt.Fprintf(os.Stderr, "  youtube-downloader -url \"URL\" -audio-only\n\n")
	fmt.Fprintf(os.Stderr, "  # プレイリスト全体をダウンロード\n")
	fmt.Fprintf(os.Stderr, "  youtube-downloader -url \"PLAYLIST_URL\" -playlist\n\n")
	fmt.Fprintf(os.Stderr, "オプション:\n")
	flag.PrintDefaults()
}
