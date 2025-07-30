package config

import (
	"flag"
	"fmt"
	"os"
)

// Config はOCR実行CLIの設定を保持する構造体
type Config struct {
	Path              string  // 画像ファイルまたはディレクトリのパス
	Recursive         bool    // ディレクトリを再帰的に検索するか
	Model             string  // 使用するGeminiモデル
	Prompt            string  // OCR用プロンプト
	SystemInstruction string  // システム指示
	Temperature       float64 // 生成パラメータ
	MaxTokens         int     // 最大トークン数
	AiType            string  // AIタイプ ("gemini" または "vertex")
	APIKey            string  // Gemini API キー
	Project           string  // Google Cloud プロジェクトID
	Location          string  // Google Cloud ロケーション
	Help              bool    // ヘルプ表示フラグ
}

// デフォルト値の定数
const (
	DefaultModel = "gemini-2.5-flash-lite"
	DefaultPrompt = "OCRして。補足や説明は不要です。"
	DefaultSystemInstruction = "OCRして。"
	DefaultTemperature = 1.0
	DefaultMaxTokens   = 8192
	DefaultAiType      = "gemini"
	DefaultLocation    = "us-central1"
)

// NewConfig は新しいConfigを作成する
func NewConfig(path string, recursive bool, model, prompt, systemInstruction string, temperature float64, maxTokens int, aiType, apiKey, project, location string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf("パスが指定されていません")
	}

	// パスの存在確認
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("指定されたパスが存在しません: %s", path)
	}

	// 温度の範囲チェック
	if temperature < 0.0 || temperature > 2.0 {
		return nil, fmt.Errorf("温度は0.0から2.0の範囲で指定してください: %f", temperature)
	}

	// 最大トークン数のチェック
	if maxTokens <= 0 {
		return nil, fmt.Errorf("最大トークン数は正の値で指定してください: %d", maxTokens)
	}

	// AIタイプの検証
	if aiType != "gemini" && aiType != "vertex" {
		return nil, fmt.Errorf("無効なAIタイプです: %s (gemini または vertex を指定してください)", aiType)
	}

	// 認証情報の検証
	if aiType == "gemini" && apiKey == "" {
		return nil, fmt.Errorf("Gemini API使用時はapi-keyが必要です")
	}
	if aiType == "vertex" && project == "" {
		return nil, fmt.Errorf("Vertex AI使用時はprojectが必要です")
	}

	return &Config{
		Path:              path,
		Recursive:         recursive,
		Model:             model,
		Prompt:            prompt,
		SystemInstruction: systemInstruction,
		Temperature:       temperature,
		MaxTokens:         maxTokens,
		AiType:            aiType,
		APIKey:            apiKey,
		Project:           project,
		Location:          location,
	}, nil
}

// StandardFlagParser の実装を更新
func (p *StandardFlagParser) StringVar(ptr *string, name string, value string, usage string) {
	flag.StringVar(ptr, name, value, usage)
}

func (p *StandardFlagParser) BoolVar(ptr *bool, name string, value bool, usage string) {
	flag.BoolVar(ptr, name, value, usage)
}

func (p *StandardFlagParser) Float64Var(ptr *float64, name string, value float64, usage string) {
	flag.Float64Var(ptr, name, value, usage)
}

func (p *StandardFlagParser) IntVar(ptr *int, name string, value int, usage string) {
	flag.IntVar(ptr, name, value, usage)
}

func (p *StandardFlagParser) Parse() error {
	flag.Parse()
	return nil
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		path              = ""
		recursive         = false
		model             = DefaultModel
		prompt            = DefaultPrompt
		systemInstruction = DefaultSystemInstruction
		temperature       = DefaultTemperature
		maxTokens         = DefaultMaxTokens
		aiType            = DefaultAiType
		apiKey            = ""
		project           = ""
		location          = DefaultLocation
		help              = false
	)

	// パス関連のフラグ
	parser.StringVar(&path, "path", path, "画像ファイルまたはディレクトリのパス")
	parser.StringVar(&path, "p", path, "パスの短縮形")
	parser.BoolVar(&recursive, "recursive", recursive, "ディレクトリを再帰的に検索")
	parser.BoolVar(&recursive, "r", recursive, "再帰の短縮形")

	// モデル設定
	parser.StringVar(&model, "model", model, "使用するGeminiモデル")
	parser.StringVar(&model, "m", model, "モデルの短縮形")

	// プロンプト設定
	parser.StringVar(&prompt, "prompt", prompt, "OCR用プロンプト")
	parser.StringVar(&prompt, "pr", prompt, "プロンプトの短縮形")
	parser.StringVar(&systemInstruction, "system-instruction", systemInstruction, "システム指示")
	parser.StringVar(&systemInstruction, "si", systemInstruction, "システム指示の短縮形")

	// 生成パラメータ
	parser.Float64Var(&temperature, "temperature", temperature, "生成パラメータ (0.0-2.0)")
	parser.Float64Var(&temperature, "t", temperature, "温度の短縮形")
	parser.IntVar(&maxTokens, "max-tokens", maxTokens, "最大トークン数")
	parser.IntVar(&maxTokens, "mt", maxTokens, "最大トークンの短縮形")

	// 認証関連のフラグ
	parser.StringVar(&aiType, "ai-type", aiType, "AIタイプ (gemini, vertex)")
	parser.StringVar(&aiType, "at", aiType, "AIタイプの短縮形")
	parser.StringVar(&apiKey, "api-key", apiKey, "Gemini API キー")
	parser.StringVar(&apiKey, "ak", apiKey, "APIキーの短縮形")
	parser.StringVar(&project, "project", project, "Google Cloud プロジェクトID")
	parser.StringVar(&project, "pj", project, "プロジェクトの短縮形")
	parser.StringVar(&location, "location", location, "Google Cloud ロケーション")
	parser.StringVar(&location, "loc", location, "ロケーションの短縮形")

	// ヘルプ
	parser.BoolVar(&help, "help", help, "ヘルプを表示")
	parser.BoolVar(&help, "h", help, "ヘルプの短縮形")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	// ヘルプが要求された場合
	if help {
		return &Config{Help: true}, nil
	}

	return NewConfig(path, recursive, model, prompt, systemInstruction, temperature, maxTokens, aiType, apiKey, project, location)
}

// ParseFlags はコマンドライン引数を解析してConfigを作成する
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, `AI OCR実行CLIツール

使用方法:
  Gemini API使用:
    %s -path /path/to/image.webp -ai-type gemini -api-key "your-api-key"

  Vertex AI使用:
    %s -path /path/to/image.webp -ai-type vertex -project "your-project-id"

  ディレクトリ内の画像（再帰）:
    %s -path /path/to/directory -recursive -ai-type gemini -api-key "your-api-key"

  詳細設定:
    %s -path /path/to/screenshots -recursive -ai-type gemini -api-key "your-api-key" -model gemini-2.0-flash -temperature 0.8 -max-tokens 4096

  短縮形:
    %s -p /path/to/image.webp -at gemini -ak "your-api-key" -m gemini-1.5-pro-002 -pr "テキストを抽出" -t 0.5 -mt 2048

オプション:
  -path, -p              画像ファイルまたはディレクトリのパス (必須)
  -recursive, -r         ディレクトリを再帰的に検索 (デフォルト: false)
  -ai-type, -at          AIタイプ (gemini, vertex) (デフォルト: %s)
  -api-key, -ak          Gemini API キー (Gemini使用時必須)
  -project, -pj          Google Cloud プロジェクトID (Vertex AI使用時必須)
  -location, -loc        Google Cloud ロケーション (デフォルト: %s)
  -model, -m             使用するGeminiモデル (デフォルト: %s)
  -prompt, -pr           OCR用プロンプト
  -system-instruction, -si システム指示
  -temperature, -t       生成パラメータ 0.0-2.0 (デフォルト: %.1f)
  -max-tokens, -mt       最大トークン数 (デフォルト: %d)
  -help, -h              このヘルプを表示

サポートされる画像形式:
  .jpg, .jpeg, .png, .gif, .bmp, .webp

認証について:
  - Gemini API使用時: -api-key が必要
  - Vertex AI使用時: -project が必要、-location はオプション

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], DefaultAiType, DefaultLocation, DefaultModel, DefaultTemperature, DefaultMaxTokens)
}
