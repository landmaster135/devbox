package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	steamAPI "github.com/landmaster135/devbox/internal/steam/infrastructure/steam_api"
)

// #==============================================================#
// ##       Mocks for FileWriter                                 ##
// #==============================================================#

// #==============================================================#
// ##       Interfaces for FileWriter                            ##
// #==============================================================#
// FileWriterInterface はファイル書き込み機能を抽象化
type FileWriterInterface interface {
	WriteToFile(data any, filename string) error
}

// #==============================================================#
// ##       Implementations for FileWriter                       ##
// #==============================================================#
// FileWriter はデフォルトのファイル書き込み実装
type FileWriter struct{}

// WriteToFile はデータをJSONファイルに保存します
func (w *FileWriter) WriteToFile(data any, filename string) error {
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

// #==============================================================#
// ##       Mocks for Logger                                     ##
// #==============================================================#

// #==============================================================#
// ##       Interfaces for Logger                                ##
// #==============================================================#
// LoggerInterface はログ出力機能を抽象化
type LoggerInterface interface {
	Printf(format string, v ...any)
	Println(v ...any)
}

// #==============================================================#
// ##       Implementations for Logger                           ##
// #==============================================================#
// Logger はデフォルトのログ出力実装
type Logger struct{}

// Printf はフォーマット付きでログを出力します
func (l *Logger) Printf(format string, v ...any) {
	fmt.Printf(format, v...)
}

// Println はログを出力します
func (l *Logger) Println(v ...any) {
	fmt.Println(v...)
}

// #==============================================================#
// ##       Mocks for ConcurrencyManager                         ##
// #==============================================================#

// #==============================================================#
// ##       Interfaces for ConcurrencyManager                    ##
// #==============================================================#
// ConcurrencyManagerInterface は並行処理制御を抽象化
type ConcurrencyManagerInterface interface {
	GetSemaphore() chan struct{}
	GetMaxConcurrency() int
}

// #==============================================================#
// ##       Implementations for ConcurrencyManager               ##
// #==============================================================#
// ConcurrencyManager はデフォルトの並行処理制御実装
type ConcurrencyManager struct {
	maxConcurrency int
	semaphore      chan struct{}
}

// GetSemaphore はセマフォチャネルを返します
func (c *ConcurrencyManager) GetSemaphore() chan struct{} {
	if c.semaphore == nil {
		c.semaphore = make(chan struct{}, c.maxConcurrency)
	}
	return c.semaphore
}

// GetMaxConcurrency は最大並行数を返します
func (c *ConcurrencyManager) GetMaxConcurrency() int {
	return c.maxConcurrency
}

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

// GameStat はゲーム統計情報を格納する構造体
type GameStat struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

// GameAchievement はゲーム実績情報を格納する構造体
type GameAchievement struct {
	Name       string `json:"name"`
	Achieved   int    `json:"achieved"`
	UnlockTime int64  `json:"unlock_time,omitempty"`
}

// GameStatsInfo はゲーム統計と実績の情報を格納する構造体
type GameStatsInfo struct {
	SteamID      string            `json:"steam_id"`
	GameName     string            `json:"game_name"`
	GameID       int               `json:"game_id"`
	Stats        []GameStat        `json:"stats"`
	Achievements []GameAchievement `json:"achievements"`
}

// GameListOutput はJSONファイル出力用の構造体
type GameListOutput struct {
	SteamID     string          `json:"steam_id"`
	GeneratedAt string          `json:"generated_at"`
	TotalGames  int             `json:"total_games"`
	Games       []SteamGameInfo `json:"games"`
}

// #==============================================================#
// ##       Mocks for SteamService                               ##
// #==============================================================#
// MockSteamService はSteamServiceのモック実装
type MockSteamService struct {
	GetGamesInfoFunc             func(ctx context.Context, steamID string) ([]SteamGameInfo, error)
	SaveGamesToJSONFunc          func(games []SteamGameInfo, steamID string, filename string) error
	GetGamesStatsFunc            func(ctx context.Context, steamID string) ([]*GameStatsInfo, error)
	SaveGamesStatsToJSONFunc     func(allGameStats []*GameStatsInfo, filename string) error
	GetClientFunc                func() steamAPI.SteamClientInterface
	GetFileWriterFunc            func() FileWriterInterface
	GetLoggerFunc                func() LoggerInterface
	GetConcurrencyControllerFunc func() ConcurrencyManagerInterface
}

// GetClient はモックのGetClientメソッド
func (m *MockSteamService) GetClient() steamAPI.SteamClientInterface {
	if m.GetClientFunc != nil {
		return m.GetClientFunc()
	}
	return nil
}

// GetFileWriter はモックのGetFileWriterメソッド
func (m *MockSteamService) GetFileWriter() FileWriterInterface {
	if m.GetFileWriterFunc != nil {
		return m.GetFileWriterFunc()
	}
	return nil
}

// GetLogger はモックのGetLoggerメソッド
func (m *MockSteamService) GetLogger() LoggerInterface {
	if m.GetLoggerFunc != nil {
		return m.GetLoggerFunc()
	}
	return nil
}

// GetConcurrencyController はモックのGetConcurrencyControllerメソッド
func (m *MockSteamService) GetConcurrencyController() ConcurrencyManagerInterface {
	if m.GetConcurrencyControllerFunc != nil {
		return m.GetConcurrencyControllerFunc()
	}
	return nil
}

// GetGamesInfo はモックのGetGamesInfoメソッド
func (m *MockSteamService) GetGamesInfo(ctx context.Context, steamID string) ([]SteamGameInfo, error) {
	if m.GetGamesInfoFunc != nil {
		return m.GetGamesInfoFunc(ctx, steamID)
	}
	return nil, nil
}

// SaveGamesToJSON はモックのSaveGamesToJSONメソッド
func (m *MockSteamService) SaveGamesToJSON(games []SteamGameInfo, steamID string, filename string) error {
	if m.SaveGamesToJSONFunc != nil {
		return m.SaveGamesToJSONFunc(games, steamID, filename)
	}
	return nil
}

// GetGamesStats はモックのGetGamesStatsメソッド
func (m *MockSteamService) GetGamesStats(ctx context.Context, steamID string) ([]*GameStatsInfo, error) {
	if m.GetGamesStatsFunc != nil {
		return m.GetGamesStatsFunc(ctx, steamID)
	}
	return nil, nil
}

// SaveGamesStatsToJSON はモックのSaveGamesStatsToJSONメソッド
func (m *MockSteamService) SaveGamesStatsToJSON(allGameStats []*GameStatsInfo, filename string) error {
	if m.SaveGamesStatsToJSONFunc != nil {
		return m.SaveGamesStatsToJSONFunc(allGameStats, filename)
	}
	return nil
}

// #==============================================================#
// ##       Interfaces for SteamService                          ##
// #==============================================================#
// SteamServiceInterface はSteamService全体を抽象化
type SteamServiceInterface interface {
	GetGamesInfo(ctx context.Context, steamID string) ([]SteamGameInfo, error)
	SaveGamesToJSON(games []SteamGameInfo, steamID string, filename string) error
	GetGamesStats(ctx context.Context, steamID string) ([]*GameStatsInfo, error)
	SaveGamesStatsToJSON(allGameStats []*GameStatsInfo, filename string) error
	GetClient() steamAPI.SteamClientInterface
	GetFileWriter() FileWriterInterface
	GetLogger() LoggerInterface
	GetConcurrencyController() ConcurrencyManagerInterface
}

// #==============================================================#
// ##       Implementations for SteamService                     ##
// #==============================================================#
// SteamService はSteam APIを使用するサービス
type SteamService struct {
	client                steamAPI.SteamClientInterface
	fileWriter            FileWriterInterface
	logger                LoggerInterface
	concurrencyController ConcurrencyManagerInterface
}

// NewSteamService は依存性注入対応のコンストラクタ
func NewSteamService(client steamAPI.SteamClientInterface, fileWriter FileWriterInterface, logger LoggerInterface, concurrencyManager ConcurrencyManagerInterface) *SteamService {
	service := &SteamService{
		client: client,
	}

	// nilの場合はデフォルト実装を使用
	service.fileWriter = &FileWriter{}
	if fileWriter != nil {
		service.fileWriter = fileWriter
	}

	service.logger = &Logger{}
	if logger != nil {
		service.logger = logger
	}

	service.concurrencyController = &ConcurrencyManager{maxConcurrency: 10}
	if concurrencyManager != nil {
		service.concurrencyController = concurrencyManager
	}

	return service
}

// NewSteamServiceWithAPIKey はAPIキーから直接SteamServiceを作成するヘルパー関数
func NewSteamServiceWithAPIKey(apiKey string) *SteamService {
	client := steamAPI.NewSteamClient(apiKey, nil)
	return NewSteamService(client, nil, nil, nil)
}

// GetClient はSteamクライアントを返します
func (s *SteamService) GetClient() steamAPI.SteamClientInterface {
	return s.client
}

// GetFileWriter はファイルライターを返します
func (s *SteamService) GetFileWriter() FileWriterInterface {
	return s.fileWriter
}

// GetLogger はロガーを返します
func (s *SteamService) GetLogger() LoggerInterface {
	return s.logger
}

// GetConcurrencyController は並行処理制御を返します
func (s *SteamService) GetConcurrencyController() ConcurrencyManagerInterface {
	return s.concurrencyController
}

// buildGameInfo は個別のゲーム情報を構築します（Steam IDを明示的に渡すバージョン）
func (s *SteamService) buildGameInfo(ctx context.Context, ownedGame steamAPI.OwnedGame, recentGamesMap map[int]steamAPI.RecentlyPlayedGame, steamID string) SteamGameInfo {
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
		_, err := s.GetClient().GetApps().GetUserStats(ctx, steamID, ownedGame.AppID)
		if err != nil {
			gameInfo.Stats = false
		}

		// 実績情報を取得してみる（エラーが発生しても続行）
		_, err = s.GetClient().GetApps().GetUserAchievements(ctx, steamID, ownedGame.AppID, "en")
		if err != nil {
			gameInfo.AchievementsCanRetrieve = false
		}
	}

	return gameInfo
}

// GetGamesInfo は指定されたSteam IDのゲーム情報を取得します
func (s *SteamService) GetGamesInfo(ctx context.Context, steamID string) ([]SteamGameInfo, error) {
	// 所有ゲーム一覧を取得
	ownedGames, err := s.GetClient().GetUsers().GetOwnedGames(ctx, steamID, true, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get owned games: %w", err)
	}

	// 最近プレイしたゲーム一覧を取得（2週間のプレイ時間情報のため）
	recentGames, err := s.GetClient().GetUsers().GetUserRecentlyPlayedGames(ctx, steamID)
	if err != nil {
		// 最近プレイしたゲームの取得に失敗してもエラーにしない
		s.GetLogger().Printf("Warning: failed to get recently played games: %v\n", err)
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
	semaphore := s.GetConcurrencyController().GetSemaphore()

	for i, game := range ownedGames {
		wg.Add(1)
		go func(index int, ownedGame steamAPI.OwnedGame) {
			defer wg.Done()
			semaphore <- struct{}{}        // セマフォを取得
			defer func() { <-semaphore }() // セマフォを解放

			gameInfo := s.buildGameInfo(ctx, ownedGame, recentGamesMap, steamID)

			mu.Lock()
			gameInfos[index] = gameInfo
			mu.Unlock()
		}(i, game)
	}

	wg.Wait()

	return gameInfos, nil
}

// saveToJSONFile はデータをJSONファイルに保存します
func (s *SteamService) saveToJSONFile(data any, filename string) error {
	return s.GetFileWriter().WriteToFile(data, filename)
}

// SaveGamesToJSON はゲーム一覧をJSONファイルに保存します
func (s *SteamService) SaveGamesToJSON(games []SteamGameInfo, steamID string, filename string) error {
	output := GameListOutput{
		SteamID:     steamID,
		GeneratedAt: time.Now().Format(time.RFC3339),
		TotalGames:  len(games),
		Games:       games,
	}
	return s.saveToJSONFile(output, filename)
}

// buildGameStatsInfo は個別のゲーム統計情報を構築します
func (s *SteamService) buildGameStatsInfo(ctx context.Context, game steamAPI.OwnedGame, steamID string) *GameStatsInfo {
	gameStats := &GameStatsInfo{
		SteamID:      steamID,
		GameName:     game.Name,
		GameID:       game.AppID,
		Stats:        []GameStat{},
		Achievements: []GameAchievement{},
	}

	// 統計情報を取得
	userStats, err := s.GetClient().GetApps().GetUserStats(ctx, steamID, game.AppID)
	if err == nil {
		for _, stat := range userStats.Stats {
			gameStats.Stats = append(gameStats.Stats, GameStat{
				Name:  stat.Name,
				Value: stat.Value,
			})
		}
	} else {
		// エラーログを出力（デバッグ用）
		s.GetLogger().Printf("Warning: failed to get user stats for game %s (ID: %d): %v\n", game.Name, game.AppID, err)
	}

	// 実績情報を取得
	achievements, err := s.GetClient().GetApps().GetUserAchievements(ctx, steamID, game.AppID, "en")
	if err == nil {
		for _, achievement := range achievements {
			gameStats.Achievements = append(gameStats.Achievements, GameAchievement{
				Name:       achievement.Name,
				Achieved:   achievement.Achieved,
				UnlockTime: 0,
			})
		}
	} else {
		// エラーログを出力（デバッグ用）
		s.GetLogger().Printf("Warning: failed to get user achievements for game %s (ID: %d): %v\n", game.Name, game.AppID, err)
	}

	return gameStats
}

// GetGamesStats はゲームの統計情報を取得します
func (s *SteamService) GetGamesStats(ctx context.Context, steamID string) ([]*GameStatsInfo, error) {
	// 所有ゲーム一覧を取得
	ownedGames, err := s.GetClient().GetUsers().GetOwnedGames(ctx, steamID, true, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get owned games: %w", err)
	}

	// 並行処理でゲーム統計情報を取得
	allGameStats := make([]*GameStatsInfo, len(ownedGames))
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 並行処理数を制限（Steam APIのレート制限を考慮）
	semaphore := s.GetConcurrencyController().GetSemaphore()

	for i, game := range ownedGames {
		wg.Add(1)
		go func(index int, ownedGame steamAPI.OwnedGame) {
			defer wg.Done()
			semaphore <- struct{}{}        // セマフォを取得
			defer func() { <-semaphore }() // セマフォを解放

			gameStats := s.buildGameStatsInfo(ctx, ownedGame, steamID)

			mu.Lock()
			allGameStats[index] = gameStats
			mu.Unlock()
		}(i, game)
	}

	wg.Wait()

	return allGameStats, nil
}

// SaveGamesStatsToJSON はゲームの統計情報をJSONファイルに保存します
func (s *SteamService) SaveGamesStatsToJSON(allGameStats []*GameStatsInfo, filename string) error {
	return s.saveToJSONFile(allGameStats, filename)
}
