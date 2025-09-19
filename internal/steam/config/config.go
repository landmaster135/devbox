package config

import (
	"fmt"
	"os"
)

// Config はSteam CLIの設定を保持する構造体
type Config struct {
	Operation   string // 実行する操作 (games, game-stats)
	SteamAPIKey string // Steam API キー
	SteamID     string // Steam ID
	GameID      int    // ゲームID (game-stats操作で使用)
	Help        bool   // ヘルプ表示フラグ
}

// デフォルト値の定数
const (
	DefaultOperation = ""
	DefaultGameID    = 0
)

// validateConfig は設定を検証します
func validateConfig(operation, steamAPIKey, steamID string, gameID int) error {
	if operation == "" {
		return fmt.Errorf("operation is required")
	}

	if steamAPIKey == "" {
		return fmt.Errorf("steam-api-key is required")
	}

	if steamID == "" {
		return fmt.Errorf("steam-id is required")
	}

	// 現在サポートされている操作をチェック
	switch operation {
	case "games":
		// OK
	case "game-stats":
		// OK
	default:
		return fmt.Errorf("unsupported operation: %s (supported: games, game-stats)", operation)
	}

	// Steam IDの形式を簡単にチェック
	if len(steamID) != 17 {
		return fmt.Errorf("invalid Steam ID format: %s (should be 17 digits)", steamID)
	}

	return nil
}

// NewConfig は新しいConfigを作成する
func NewConfig(operation, steamAPIKey, steamID string, gameID int) (*Config, error) {
	err := validateConfig(operation, steamAPIKey, steamID, gameID)
	if err != nil {
		return nil, fmt.Errorf("設定の初期化に失敗しました: %w", err)
	}

	return &Config{
		Operation:   operation,
		SteamAPIKey: steamAPIKey,
		SteamID:     steamID,
		GameID:      gameID,
	}, nil
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		operation   = DefaultOperation
		steamAPIKey = ""
		steamID     = ""
		gameID      = DefaultGameID
		help        = false
	)

	// 操作関連のフラグ
	parser.StringVar(&operation, "operation", operation, "Operation to perform (required): games, game-stats")
	parser.StringVar(&operation, "o", operation, "Operation to perform (required): games, game-stats (shorthand)")

	// 認証関連のフラグ
	parser.StringVar(&steamAPIKey, "steam-api-key", steamAPIKey, "Steam API key (required)")
	parser.StringVar(&steamAPIKey, "k", steamAPIKey, "Steam API key (required) (shorthand)")
	parser.StringVar(&steamID, "steam-id", steamID, "Steam ID (required)")
	parser.StringVar(&steamID, "s", steamID, "Steam ID (required) (shorthand)")

	// ゲーム関連のフラグ
	parser.IntVar(&gameID, "game-id", gameID, "Game ID (required for game-stats operation)")
	parser.IntVar(&gameID, "g", gameID, "Game ID (required for game-stats operation) (shorthand)")

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

	return NewConfig(operation, steamAPIKey, steamID, gameID)
}

// ParseFlags はコマンドライン引数を解析してConfigを作成する
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "Steam API CLI Tool\n\n")
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -operation, -o         Operation to perform (required): games, game-stats\n")
	fmt.Fprintf(os.Stderr, "  -steam-api-key, -k     Steam API key (required)\n")
	fmt.Fprintf(os.Stderr, "  -steam-id, -s          Steam ID (required)\n")
	fmt.Fprintf(os.Stderr, "  -game-id, -g           Game ID (required for game-stats operation)\n")
	fmt.Fprintf(os.Stderr, "  -help, -h              Show this help message\n")
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  %s --operation games --steam-api-key YOUR_KEY --steam-id 76561198000000000\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -o games -k YOUR_KEY -s 76561198000000000\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s --operation game-stats --steam-api-key YOUR_KEY --steam-id 76561198000000000\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -o game-stats -k YOUR_KEY -s 76561198000000000\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nSupported operations:\n")
	fmt.Fprintf(os.Stderr, "  games       Get user's owned games information and save to JSON file\n")
	fmt.Fprintf(os.Stderr, "  game-stats  Get all games' statistics and achievements and save to JSON file\n")
}
