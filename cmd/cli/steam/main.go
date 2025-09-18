package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/landmaster135/devbox/internal/steam/usecases"
)

// Config はCLIの設定を格納する構造体
type Config struct {
	Operation   string
	SteamAPIKey string
	SteamID     string
}

// GameListOutput はJSONファイル出力用の構造体
type GameListOutput struct {
	SteamID     string                      `json:"steam_id"`
	GeneratedAt string                      `json:"generated_at"`
	TotalGames  int                         `json:"total_games"`
	Games       []usecases.SteamGameInfo    `json:"games"`
}

func main() {
	config := parseFlags()

	if err := validateConfig(config); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	ctx := context.Background()

	switch config.Operation {
	case "games":
		if err := handleGamesOperation(ctx, config); err != nil {
			log.Fatalf("Failed to execute games operation: %v", err)
		}
	default:
		log.Fatalf("Unknown operation: %s", config.Operation)
	}
}

// parseFlags はコマンドライン引数を解析します
func parseFlags() Config {
	var config Config

	flag.StringVar(&config.Operation, "operation", "", "Operation to perform (required): games")
	flag.StringVar(&config.Operation, "o", "", "Operation to perform (required): games (shorthand)")
	flag.StringVar(&config.SteamAPIKey, "steam-api-key", "", "Steam API key (required)")
	flag.StringVar(&config.SteamAPIKey, "k", "", "Steam API key (required) (shorthand)")
	flag.StringVar(&config.SteamID, "steam-id", "", "Steam ID (required)")
	flag.StringVar(&config.SteamID, "s", "", "Steam ID (required) (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Steam API CLI Tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s --operation games --steam-api-key YOUR_KEY --steam-id 76561198000000000\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -o games -k YOUR_KEY -s 76561198000000000\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nSupported operations:\n")
		fmt.Fprintf(os.Stderr, "  games    Get user's owned games information and save to JSON file\n")
	}

	flag.Parse()

	return config
}

// validateConfig は設定を検証します
func validateConfig(config Config) error {
	if config.Operation == "" {
		return fmt.Errorf("operation is required")
	}

	if config.SteamAPIKey == "" {
		return fmt.Errorf("steam-api-key is required")
	}

	if config.SteamID == "" {
		return fmt.Errorf("steam-id is required")
	}

	// 現在サポートされている操作をチェック
	switch config.Operation {
	case "games":
		// OK
	default:
		return fmt.Errorf("unsupported operation: %s (supported: games)", config.Operation)
	}

	// Steam IDの形式を簡単にチェック
	if len(config.SteamID) != 17 {
		return fmt.Errorf("invalid Steam ID format: %s (should be 17 digits)", config.SteamID)
	}

	return nil
}

// handleGamesOperation はgames操作を処理します
func handleGamesOperation(ctx context.Context, config Config) error {
	fmt.Printf("Starting games operation for Steam ID: %s\n", config.SteamID)

	// Steam サービスを作成
	steamService := usecases.NewSteamService(config.SteamAPIKey)

	// ゲーム情報を取得
	fmt.Println("Fetching games information...")
	games, err := steamService.GetGamesInfoWithSteamID(ctx, config.SteamID)
	if err != nil {
		return fmt.Errorf("failed to get games info: %w", err)
	}

	fmt.Printf("Successfully retrieved %d games\n", len(games))

	// JSON出力用のデータを準備
	output := GameListOutput{
		SteamID:     config.SteamID,
		GeneratedAt: time.Now().Format(time.RFC3339),
		TotalGames:  len(games),
		Games:       games,
	}

	// JSONファイルに出力
	filename := fmt.Sprintf("steam_games_%s_%s.json", config.SteamID, time.Now().Format("20060102_150405"))
	if err := saveToJSONFile(output, filename); err != nil {
		return fmt.Errorf("failed to save to JSON file: %w", err)
	}

	fmt.Printf("Games information saved to: %s\n", filename)

	// 簡単な統計情報を表示
	displayStatistics(games)

	return nil
}

// saveToJSONFile はデータをJSONファイルに保存します
func saveToJSONFile(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // 読みやすい形式でインデント

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// displayStatistics は簡単な統計情報を表示します
func displayStatistics(games []usecases.SteamGameInfo) {
	if len(games) == 0 {
		fmt.Println("No games found.")
		return
	}

	fmt.Println("\n=== Statistics ===")
	fmt.Printf("Total games: %d\n", len(games))

	// プレイ時間の統計
	var totalPlaytime int
	var gamesWithPlaytime int
	var gamesWithRecentPlaytime int
	var achievementsAvailable int
	var statsAvailable int

	for _, game := range games {
		if game.PlaytimeForever > 0 {
			totalPlaytime += game.PlaytimeForever
			gamesWithPlaytime++
		}
		if game.PlaytimeRecent2Weeks > 0 {
			gamesWithRecentPlaytime++
		}
		if game.AchievementsCanRetrieve {
			achievementsAvailable++
		}
		if game.Stats {
			statsAvailable++
		}
	}

	fmt.Printf("Games with playtime: %d\n", gamesWithPlaytime)
	fmt.Printf("Total playtime: %.1f hours\n", float64(totalPlaytime)/60.0)
	if gamesWithPlaytime > 0 {
		fmt.Printf("Average playtime per played game: %.1f hours\n", float64(totalPlaytime)/float64(gamesWithPlaytime)/60.0)
	}
	fmt.Printf("Games played in recent 2 weeks: %d\n", gamesWithRecentPlaytime)
	fmt.Printf("Games with achievements available: %d\n", achievementsAvailable)
	fmt.Printf("Games with stats available: %d\n", statsAvailable)

	// 最もプレイしたゲームトップ5
	fmt.Println("\n=== Top 5 Most Played Games ===")

	// ゲームをプレイ時間でソート（簡単な実装）
	topGames := make([]usecases.SteamGameInfo, len(games))
	copy(topGames, games)

	// バブルソート（簡単な実装）
	for i := 0; i < len(topGames)-1; i++ {
		for j := 0; j < len(topGames)-i-1; j++ {
			if topGames[j].PlaytimeForever < topGames[j+1].PlaytimeForever {
				topGames[j], topGames[j+1] = topGames[j+1], topGames[j]
			}
		}
	}

	// トップ5を表示
	for i, game := range topGames {
		if i >= 5 || game.PlaytimeForever == 0 {
			break
		}
		fmt.Printf("%d. %s - %.1f hours\n", i+1, game.Name, float64(game.PlaytimeForever)/60.0)
	}
}
