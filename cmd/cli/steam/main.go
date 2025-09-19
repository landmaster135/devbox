package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	config "github.com/landmaster135/devbox/internal/steam/config"
	usecases "github.com/landmaster135/devbox/internal/steam/usecases"
)

// handleGamesOperation はgames操作を処理します
func handleGamesOperation(ctx context.Context, cfg *config.Config) error {
	fmt.Printf("Starting games operation for Steam ID: %s\n", cfg.SteamID)

	// Steam サービスを作成
	steamService := usecases.NewSteamServiceWithAPIKey(cfg.SteamAPIKey)

	// ゲーム情報を取得
	fmt.Println("Fetching games information...")
	games, err := steamService.GetGamesInfo(ctx, cfg.SteamID)
	if err != nil {
		return fmt.Errorf("failed to get games info: %w", err)
	}

	fmt.Printf("Successfully retrieved %d games\n", len(games))

	// JSONファイルに出力（サービスのメソッドを使用）
	filename := fmt.Sprintf("steam_games_%s_%s.json", cfg.SteamID, time.Now().Format("20060102_150405"))
	if err := steamService.SaveGamesToJSON(games, cfg.SteamID, filename); err != nil {
		return fmt.Errorf("failed to save to JSON file: %w", err)
	}

	fmt.Printf("Games information saved to: %s\n", filename)

	// 簡単な統計情報を表示
	displayStatistics(games)

	return nil
}

// handleGameStatsOperation はgame-stats操作を処理します
func handleGameStatsOperation(ctx context.Context, cfg *config.Config) error {
	fmt.Printf("Starting game-stats operation for Steam ID: %s\n", cfg.SteamID)

	// Steam サービスを作成
	steamService := usecases.NewSteamServiceWithAPIKey(cfg.SteamAPIKey)

	// 全ゲームの統計情報を取得
	fmt.Println("Fetching all games statistics and achievements...")
	allGameStats, err := steamService.GetGamesStats(ctx, cfg.SteamID)
	if err != nil {
		return fmt.Errorf("failed to get all games stats: %w", err)
	}

	fmt.Printf("Successfully retrieved stats for %d games\n", len(allGameStats))

	// JSONファイルに出力（サービスのメソッドを使用）
	filename := fmt.Sprintf("steam_games_stats_%s_%s.json", cfg.SteamID, time.Now().Format("20060102_150405"))
	if err := steamService.SaveGamesStatsToJSON(allGameStats, filename); err != nil {
		return fmt.Errorf("failed to save to JSON file: %w", err)
	}

	fmt.Printf("All games statistics saved to: %s\n", filename)

	// 簡単な統計情報を表示
	displayGamesStatsInfo(allGameStats)

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

// displayGamesStatsInfo は全ゲームの統計情報を表示します
func displayGamesStatsInfo(allGameStats []*usecases.GameStatsInfo) {
	if len(allGameStats) == 0 {
		fmt.Println("No games found.")
		return
	}

	fmt.Printf("\n=== All Games Statistics ===\n")
	fmt.Printf("Total games: %d\n", len(allGameStats))

	var totalStats int
	var totalAchievements int
	var gamesWithStats int
	var gamesWithAchievements int

	for _, gameStats := range allGameStats {
		if len(gameStats.Stats) > 0 {
			totalStats += len(gameStats.Stats)
			gamesWithStats++
		}
		if len(gameStats.Achievements) > 0 {
			totalAchievements += len(gameStats.Achievements)
			gamesWithAchievements++
		}
	}

	fmt.Printf("Games with stats: %d\n", gamesWithStats)
	fmt.Printf("Total stats: %d\n", totalStats)
	fmt.Printf("Games with achievements: %d\n", gamesWithAchievements)
	fmt.Printf("Total achievements: %d\n", totalAchievements)

	// 統計・実績が多いゲームトップ5
	fmt.Println("\n=== Top 5 Games with Most Stats ===")

	// ゲームを統計数でソート
	topStatGames := make([]*usecases.GameStatsInfo, len(allGameStats))
	copy(topStatGames, allGameStats)

	// バブルソート（統計数で）
	for i := 0; i < len(topStatGames)-1; i++ {
		for j := 0; j < len(topStatGames)-i-1; j++ {
			if len(topStatGames[j].Stats) < len(topStatGames[j+1].Stats) {
				topStatGames[j], topStatGames[j+1] = topStatGames[j+1], topStatGames[j]
			}
		}
	}

	// トップ5を表示
	for i, gameStats := range topStatGames {
		if i >= 5 || len(gameStats.Stats) == 0 {
			break
		}
		fmt.Printf("%d. %s - %d stats, %d achievements\n", i+1, gameStats.GameName, len(gameStats.Stats), len(gameStats.Achievements))
	}
}

func main() {
	// コマンドライン引数を解析
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	// ヘルプが要求された場合
	if cfg.Help {
		config.PrintUsage()
		return
	}

	ctx := context.Background()

	switch cfg.Operation {
	case "games":
		if err := handleGamesOperation(ctx, cfg); err != nil {
			log.Fatalf("Failed to execute games operation: %v", err)
		}
	case "game-stats":
		if err := handleGameStatsOperation(ctx, cfg); err != nil {
			log.Fatalf("Failed to execute game-stats operation: %v", err)
		}
	default:
		log.Fatalf("Unknown operation: %s", cfg.Operation)
	}
}
