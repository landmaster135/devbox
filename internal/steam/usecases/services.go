package usecases

import (
	"context"
	"fmt"
	"sync"

	steamAPI "github.com/landmaster135/devbox/internal/steam/infrastructure/steam_api"
)

// SteamGameInfo はゲーム情報を格納する構造体
type SteamGameInfo struct {
	Name                    string `json:"name"`
	ID                      int    `json:"id"`
	Icon                    string `json:"icon"`
	Thumbnail               string `json:"thumbnail"`
	PlaytimeRecent2Weeks    int    `json:"playtime_recent_2_weeks"`
	PlaytimeDisconnected    int    `json:"playtime_disconnected"`
	PlaytimeForever         int    `json:"playtime_forever"`
	RecentTimeLastPlayed    int64  `json:"recent_time_last_played"`
	AchievementsCanRetrieve bool   `json:"achievements_can_be_retrieved"`
	Stats                   bool   `json:"stats"`
}

// SteamService はSteam APIを使用するサービス
type SteamService struct {
	client *steamAPI.SteamClient
}

// NewSteamService は新しいSteamServiceを作成します
func NewSteamService(apiKey string) *SteamService {
	client := steamAPI.NewSteamClient(apiKey, nil)
	return &SteamService{
		client: client,
	}
}

// GetGamesInfo は指定されたSteam IDのゲーム情報を取得します
func (s *SteamService) GetGamesInfo(ctx context.Context, steamID string) ([]SteamGameInfo, error) {
	// 所有ゲーム一覧を取得
	ownedGames, err := s.client.Users.GetOwnedGames(ctx, steamID, true, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get owned games: %w", err)
	}

	// 最近プレイしたゲーム一覧を取得（2週間のプレイ時間情報のため）
	recentGames, err := s.client.Users.GetUserRecentlyPlayedGames(ctx, steamID)
	if err != nil {
		// 最近プレイしたゲームの取得に失敗してもエラーにしない
		fmt.Printf("Warning: failed to get recently played games: %v\n", err)
		recentGames = []steamAPI.RecentlyPlayedGame{}
	}

	// 最近プレイしたゲームのマップを作成（AppIDをキーとする）
	recentGamesMap := make(map[int]steamAPI.RecentlyPlayedGame)
	for _, game := range recentGames {
		recentGamesMap[game.AppID] = game
	}

	// 並行処理でゲーム情報を取得
	gameInfos := make([]SteamGameInfo, len(ownedGames))
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 並行処理数を制限（Steam APIのレート制限を考慮）
	semaphore := make(chan struct{}, 10)

	for i, game := range ownedGames {
		wg.Add(1)
		go func(index int, ownedGame steamAPI.OwnedGame) {
			defer wg.Done()
			semaphore <- struct{}{} // セマフォを取得
			defer func() { <-semaphore }() // セマフォを解放

			gameInfo := s.buildGameInfo(ctx, ownedGame, recentGamesMap)

			mu.Lock()
			gameInfos[index] = gameInfo
			mu.Unlock()
		}(i, game)
	}

	wg.Wait()

	return gameInfos, nil
}

// buildGameInfo は個別のゲーム情報を構築します
func (s *SteamService) buildGameInfo(ctx context.Context, ownedGame steamAPI.OwnedGame, recentGamesMap map[int]steamAPI.RecentlyPlayedGame) SteamGameInfo {
	gameInfo := SteamGameInfo{
		Name:                    ownedGame.Name,
		ID:                      ownedGame.AppID,
		PlaytimeDisconnected:    ownedGame.PlaytimeDisconnected,
		PlaytimeForever:         ownedGame.PlaytimeForever,
		RecentTimeLastPlayed:    ownedGame.RtimeLastPlayed,
		AchievementsCanRetrieve: ownedGame.HasCommunityVisibleStats,
		Stats:                   ownedGame.HasCommunityVisibleStats,
	}

	// アイコンURLを生成
	if ownedGame.ImgIconURL != "" {
		gameInfo.Icon = fmt.Sprintf("https://media.steampowered.com/steamcommunity/public/images/apps/%d/%s.jpg", ownedGame.AppID, ownedGame.ImgIconURL)
	}

	// サムネイルURLを生成
	gameInfo.Thumbnail = fmt.Sprintf("https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/%d/library_600x900.jpg", ownedGame.AppID)

	// 最近プレイしたゲームから2週間のプレイ時間を取得
	if recentGame, exists := recentGamesMap[ownedGame.AppID]; exists {
		gameInfo.PlaytimeRecent2Weeks = recentGame.Playtime2Weeks
	} else {
		gameInfo.PlaytimeRecent2Weeks = 0
	}

	// 実績・統計情報の取得可能性をより詳細にチェック
	if ownedGame.HasCommunityVisibleStats {
		// 実際に統計情報を取得してみる（エラーが発生しても続行）
		_, err := s.client.Apps.GetUserStats(ctx, s.getSteamIDFromContext(ctx), ownedGame.AppID)
		if err != nil {
			gameInfo.Stats = false
		}

		// 実績情報を取得してみる（エラーが発生しても続行）
		_, err = s.client.Apps.GetUserAchievements(ctx, s.getSteamIDFromContext(ctx), ownedGame.AppID, "en")
		if err != nil {
			gameInfo.AchievementsCanRetrieve = false
		}
	}

	return gameInfo
}

// getSteamIDFromContext はコンテキストからSteam IDを取得します（簡易実装）
func (s *SteamService) getSteamIDFromContext(ctx context.Context) string {
	// 実際の実装では、コンテキストからSteam IDを取得するか、
	// SteamServiceにSteam IDを保持するなどの方法を使用します
	// ここでは簡易的に空文字を返します（実際の使用時は適切に実装する必要があります）
	return ""
}

// GetGamesInfoWithSteamID は指定されたSteam IDのゲーム情報を取得します（Steam IDを明示的に渡すバージョン）
func (s *SteamService) GetGamesInfoWithSteamID(ctx context.Context, steamID string) ([]SteamGameInfo, error) {
	// 所有ゲーム一覧を取得
	ownedGames, err := s.client.Users.GetOwnedGames(ctx, steamID, true, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get owned games: %w", err)
	}

	// 最近プレイしたゲーム一覧を取得（2週間のプレイ時間情報のため）
	recentGames, err := s.client.Users.GetUserRecentlyPlayedGames(ctx, steamID)
	if err != nil {
		// 最近プレイしたゲームの取得に失敗してもエラーにしない
		fmt.Printf("Warning: failed to get recently played games: %v\n", err)
		recentGames = []steamAPI.RecentlyPlayedGame{}
	}

	// 最近プレイしたゲームのマップを作成（AppIDをキーとする）
	recentGamesMap := make(map[int]steamAPI.RecentlyPlayedGame)
	for _, game := range recentGames {
		recentGamesMap[game.AppID] = game
	}

	// 並行処理でゲーム情報を取得
	gameInfos := make([]SteamGameInfo, len(ownedGames))
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 並行処理数を制限（Steam APIのレート制限を考慮）
	semaphore := make(chan struct{}, 10)

	for i, game := range ownedGames {
		wg.Add(1)
		go func(index int, ownedGame steamAPI.OwnedGame) {
			defer wg.Done()
			semaphore <- struct{}{} // セマフォを取得
			defer func() { <-semaphore }() // セマフォを解放

			gameInfo := s.buildGameInfoWithSteamID(ctx, ownedGame, recentGamesMap, steamID)

			mu.Lock()
			gameInfos[index] = gameInfo
			mu.Unlock()
		}(i, game)
	}

	wg.Wait()

	return gameInfos, nil
}

// buildGameInfoWithSteamID は個別のゲーム情報を構築します（Steam IDを明示的に渡すバージョン）
func (s *SteamService) buildGameInfoWithSteamID(ctx context.Context, ownedGame steamAPI.OwnedGame, recentGamesMap map[int]steamAPI.RecentlyPlayedGame, steamID string) SteamGameInfo {
	gameInfo := SteamGameInfo{
		Name:                    ownedGame.Name,
		ID:                      ownedGame.AppID,
		PlaytimeDisconnected:    ownedGame.PlaytimeDisconnected,
		PlaytimeForever:         ownedGame.PlaytimeForever,
		RecentTimeLastPlayed:    ownedGame.RtimeLastPlayed,
		AchievementsCanRetrieve: ownedGame.HasCommunityVisibleStats,
		Stats:                   ownedGame.HasCommunityVisibleStats,
	}

	// アイコンURLを生成
	if ownedGame.ImgIconURL != "" {
		gameInfo.Icon = fmt.Sprintf("https://media.steampowered.com/steamcommunity/public/images/apps/%d/%s.jpg", ownedGame.AppID, ownedGame.ImgIconURL)
	}

	// サムネイルURLを生成
	gameInfo.Thumbnail = fmt.Sprintf("https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/%d/library_600x900.jpg", ownedGame.AppID)

	// 最近プレイしたゲームから2週間のプレイ時間を取得
	if recentGame, exists := recentGamesMap[ownedGame.AppID]; exists {
		gameInfo.PlaytimeRecent2Weeks = recentGame.Playtime2Weeks
	} else {
		gameInfo.PlaytimeRecent2Weeks = 0
	}

	// 実績・統計情報の取得可能性をより詳細にチェック
	if ownedGame.HasCommunityVisibleStats {
		// 実際に統計情報を取得してみる（エラーが発生しても続行）
		_, err := s.client.Apps.GetUserStats(ctx, steamID, ownedGame.AppID)
		if err != nil {
			gameInfo.Stats = false
		}

		// 実績情報を取得してみる（エラーが発生しても続行）
		_, err = s.client.Apps.GetUserAchievements(ctx, steamID, ownedGame.AppID, "en")
		if err != nil {
			gameInfo.AchievementsCanRetrieve = false
		}
	}

	return gameInfo
}
