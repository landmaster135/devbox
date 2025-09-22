package usecases

import (
	"context"
	"errors"
	"testing"

	steamAPI "github.com/landmaster135/devbox/internal/steam/infrastructure/steam_api"
)

// TestSteamService_GetGamesInfo_Normal はGetGamesInfoの正常系テスト
func TestSteamService_GetGamesInfo_Normal(t *testing.T) {
	// Arrange
	mockClient := &steamAPI.MockSteamClient{
		UsersService: &steamAPI.MockUsersService{
			GetOwnedGamesFunc: func(ctx context.Context, steamID string, includeAppInfo, includePlayedFreeGames bool) ([]steamAPI.OwnedGame, error) {
				return []steamAPI.OwnedGame{
					{
						AppID:                     123,
						Name:                      "Test Game",
						PlaytimeForever:           100,
						PlaytimeDisconnected:      10,
						RtimeLastPlayed:           1234567890,
						ImgIconURL:                "test_icon",
						HasCommunityVisibleStats:  true,
					},
				}, nil
			},
			GetUserRecentlyPlayedGamesFunc: func(ctx context.Context, steamID string) ([]steamAPI.RecentlyPlayedGame, error) {
				return []steamAPI.RecentlyPlayedGame{
					{
						AppID:           123,
						Playtime2Weeks:  50,
					},
				}, nil
			},
		},
		AppsService: &steamAPI.MockAppsService{
			GetUserStatsFunc: func(ctx context.Context, steamID string, appID int) (*steamAPI.UserStats, error) {
				return &steamAPI.UserStats{
					SteamID:  steamID,
					GameName: "Test Game",
					Stats: []steamAPI.Stat{
						{Name: "test_stat", Value: 42},
					},
					Achievements: []steamAPI.Achievement{
						{Name: "test_achievement", Achieved: 1},
					},
				}, nil
			},
			GetUserAchievementsFunc: func(ctx context.Context, steamID string, appID int, language string) ([]steamAPI.Achievement, error) {
				return []steamAPI.Achievement{
					{Name: "test_achievement", Achieved: 1},
				}, nil
			},
		},
	}

	service := NewSteamService(mockClient, nil, nil, nil)
	steamID := "test_steam_id"

	// Act
	games, err := service.GetGamesInfo(context.Background(), steamID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(games) != 1 {
		t.Fatalf("Expected 1 game, got %d", len(games))
	}

	game := games[0]
	if game.Name != "Test Game" {
		t.Errorf("Expected game name 'Test Game', got '%s'", game.Name)
	}
	if game.ID != 123 {
		t.Errorf("Expected game ID 123, got %d", game.ID)
	}
	if game.PlaytimeRecent2Weeks != 50 {
		t.Errorf("Expected recent playtime 50, got %d", game.PlaytimeRecent2Weeks)
	}
	if !game.AchievementsCanRetrieve {
		t.Error("Expected achievements to be retrievable")
	}
	if !game.Stats {
		t.Error("Expected stats to be available")
	}
}

// TestSteamService_GetGamesInfo_OwnedGamesError はGetGamesInfoの所有ゲーム取得エラーテスト
func TestSteamService_GetGamesInfo_OwnedGamesError(t *testing.T) {
	// Arrange
	expectedError := errors.New("owned games error")
	mockClient := &steamAPI.MockSteamClient{
		UsersService: &steamAPI.MockUsersService{
			GetOwnedGamesFunc: func(ctx context.Context, steamID string, includeAppInfo, includePlayedFreeGames bool) ([]steamAPI.OwnedGame, error) {
				return nil, expectedError
			},
		},
	}

	service := NewSteamService(mockClient, nil, nil, nil)
	steamID := "test_steam_id"

	// Act
	games, err := service.GetGamesInfo(context.Background(), steamID)

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if games != nil {
		t.Errorf("Expected nil games, got %v", games)
	}
}

// TestSteamService_GetGamesInfo_RecentGamesError はGetGamesInfoの最近プレイしたゲーム取得エラーテスト
func TestSteamService_GetGamesInfo_RecentGamesError(t *testing.T) {
	// Arrange
	mockClient := &steamAPI.MockSteamClient{
		UsersService: &steamAPI.MockUsersService{
			GetOwnedGamesFunc: func(ctx context.Context, steamID string, includeAppInfo, includePlayedFreeGames bool) ([]steamAPI.OwnedGame, error) {
				return []steamAPI.OwnedGame{
					{
						AppID:                     123,
						Name:                      "Test Game",
						HasCommunityVisibleStats:  false,
					},
				}, nil
			},
			GetUserRecentlyPlayedGamesFunc: func(ctx context.Context, steamID string) ([]steamAPI.RecentlyPlayedGame, error) {
				return nil, errors.New("recent games error")
			},
		},
	}

	service := NewSteamService(mockClient, nil, nil, nil)
	steamID := "test_steam_id"

	// Act
	games, err := service.GetGamesInfo(context.Background(), steamID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(games) != 1 {
		t.Fatalf("Expected 1 game, got %d", len(games))
	}

	game := games[0]
	if game.PlaytimeRecent2Weeks != 0 {
		t.Errorf("Expected recent playtime 0, got %d", game.PlaytimeRecent2Weeks)
	}
}

// TestSteamService_SaveGamesToJSON_Normal はSaveGamesToJSONの正常系テスト
func TestSteamService_SaveGamesToJSON_Normal(t *testing.T) {
	// Arrange
	var capturedData any
	var capturedFilename string
	mockFileWriter := &MockFileWriter{
		WriteToFileFunc: func(data any, filename string) error {
			capturedData = data
			capturedFilename = filename
			return nil
		},
	}

	service := NewSteamService(nil, mockFileWriter, nil, nil)
	games := []SteamGameInfo{
		{Name: "Test Game", ID: 123},
	}
	steamID := "test_steam_id"
	filename := "test.json"

	// Act
	err := service.SaveGamesToJSON(games, steamID, filename)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if capturedFilename != filename {
		t.Errorf("Expected filename '%s', got '%s'", filename, capturedFilename)
	}

	output, ok := capturedData.(GameListOutput)
	if !ok {
		t.Fatal("Expected GameListOutput type")
	}
	if output.SteamID != steamID {
		t.Errorf("Expected SteamID '%s', got '%s'", steamID, output.SteamID)
	}
	if output.TotalGames != 1 {
		t.Errorf("Expected TotalGames 1, got %d", output.TotalGames)
	}
	if len(output.Games) != 1 {
		t.Errorf("Expected 1 game, got %d", len(output.Games))
	}
}

// TestSteamService_SaveGamesToJSON_Error はSaveGamesToJSONのエラーテスト
func TestSteamService_SaveGamesToJSON_Error(t *testing.T) {
	// Arrange
	expectedError := errors.New("write error")
	mockFileWriter := &MockFileWriter{
		WriteToFileFunc: func(data any, filename string) error {
			return expectedError
		},
	}

	service := NewSteamService(nil, mockFileWriter, nil, nil)
	games := []SteamGameInfo{}
	steamID := "test_steam_id"
	filename := "test.json"

	// Act
	err := service.SaveGamesToJSON(games, steamID, filename)

	// Assert
	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}
}

// TestSteamService_GetGamesStats_Normal はGetGamesStatsの正常系テスト
func TestSteamService_GetGamesStats_Normal(t *testing.T) {
	// Arrange
	mockClient := &steamAPI.MockSteamClient{
		UsersService: &steamAPI.MockUsersService{
			GetOwnedGamesFunc: func(ctx context.Context, steamID string, includeAppInfo, includePlayedFreeGames bool) ([]steamAPI.OwnedGame, error) {
				return []steamAPI.OwnedGame{
					{
						AppID: 123,
						Name:  "Test Game",
					},
				}, nil
			},
		},
		AppsService: &steamAPI.MockAppsService{
			GetUserStatsFunc: func(ctx context.Context, steamID string, appID int) (*steamAPI.UserStats, error) {
				return &steamAPI.UserStats{
					SteamID:  steamID,
					GameName: "Test Game",
					Stats: []steamAPI.Stat{
						{Name: "test_stat", Value: 42},
					},
					Achievements: []steamAPI.Achievement{
						{Name: "test_achievement", Achieved: 1},
					},
				}, nil
			},
			GetUserAchievementsFunc: func(ctx context.Context, steamID string, appID int, language string) ([]steamAPI.Achievement, error) {
				return []steamAPI.Achievement{
					{Name: "test_achievement", Achieved: 1},
				}, nil
			},
		},
	}

	service := NewSteamService(mockClient, nil, nil, nil)
	steamID := "test_steam_id"

	// Act
	stats, err := service.GetGamesStats(context.Background(), steamID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(stats) != 1 {
		t.Fatalf("Expected 1 game stats, got %d", len(stats))
	}

	gameStat := stats[0]
	if gameStat.SteamID != steamID {
		t.Errorf("Expected SteamID '%s', got '%s'", steamID, gameStat.SteamID)
	}
	if gameStat.GameName != "Test Game" {
		t.Errorf("Expected game name 'Test Game', got '%s'", gameStat.GameName)
	}
	if gameStat.GameID != 123 {
		t.Errorf("Expected game ID 123, got %d", gameStat.GameID)
	}
	if len(gameStat.Stats) != 1 {
		t.Errorf("Expected 1 stat, got %d", len(gameStat.Stats))
	}
	if len(gameStat.Achievements) != 1 {
		t.Errorf("Expected 1 achievement, got %d", len(gameStat.Achievements))
	}
}

// TestSteamService_GetGamesStats_OwnedGamesError はGetGamesStatsの所有ゲーム取得エラーテスト
func TestSteamService_GetGamesStats_OwnedGamesError(t *testing.T) {
	// Arrange
	expectedError := errors.New("owned games error")
	mockClient := &steamAPI.MockSteamClient{
		UsersService: &steamAPI.MockUsersService{
			GetOwnedGamesFunc: func(ctx context.Context, steamID string, includeAppInfo, includePlayedFreeGames bool) ([]steamAPI.OwnedGame, error) {
				return nil, expectedError
			},
		},
	}

	service := NewSteamService(mockClient, nil, nil, nil)
	steamID := "test_steam_id"

	// Act
	stats, err := service.GetGamesStats(context.Background(), steamID)

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if stats != nil {
		t.Errorf("Expected nil stats, got %v", stats)
	}
}

// TestSteamService_GetGamesStats_StatsAndAchievementsError はGetGamesStatsの統計・実績取得エラーテスト
func TestSteamService_GetGamesStats_StatsAndAchievementsError(t *testing.T) {
	// Arrange
	mockClient := &steamAPI.MockSteamClient{
		UsersService: &steamAPI.MockUsersService{
			GetOwnedGamesFunc: func(ctx context.Context, steamID string, includeAppInfo, includePlayedFreeGames bool) ([]steamAPI.OwnedGame, error) {
				return []steamAPI.OwnedGame{
					{
						AppID: 123,
						Name:  "Test Game",
					},
				}, nil
			},
		},
		AppsService: &steamAPI.MockAppsService{
			GetUserStatsFunc: func(ctx context.Context, steamID string, appID int) (*steamAPI.UserStats, error) {
				return nil, errors.New("stats error")
			},
			GetUserAchievementsFunc: func(ctx context.Context, steamID string, appID int, language string) ([]steamAPI.Achievement, error) {
				return nil, errors.New("achievements error")
			},
		},
	}

	service := NewSteamService(mockClient, nil, nil, nil)
	steamID := "test_steam_id"

	// Act
	stats, err := service.GetGamesStats(context.Background(), steamID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(stats) != 1 {
		t.Fatalf("Expected 1 game stats, got %d", len(stats))
	}

	gameStat := stats[0]
	if len(gameStat.Stats) != 0 {
		t.Errorf("Expected 0 stats, got %d", len(gameStat.Stats))
	}
	if len(gameStat.Achievements) != 0 {
		t.Errorf("Expected 0 achievements, got %d", len(gameStat.Achievements))
	}
}

// TestSteamService_SaveGamesStatsToJSON_Normal はSaveGamesStatsToJSONの正常系テスト
func TestSteamService_SaveGamesStatsToJSON_Normal(t *testing.T) {
	// Arrange
	var capturedData any
	var capturedFilename string
	mockFileWriter := &MockFileWriter{
		WriteToFileFunc: func(data any, filename string) error {
			capturedData = data
			capturedFilename = filename
			return nil
		},
	}

	service := NewSteamService(nil, mockFileWriter, nil, nil)
	gameStats := []*GameStatsInfo{
		{
			SteamID:  "test_steam_id",
			GameName: "Test Game",
			GameID:   123,
		},
	}
	filename := "test_stats.json"

	// Act
	err := service.SaveGamesStatsToJSON(gameStats, filename)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if capturedFilename != filename {
		t.Errorf("Expected filename '%s', got '%s'", filename, capturedFilename)
	}
	if capturedData == nil {
		t.Error("Expected data to be captured")
	}
}

// TestSteamService_SaveGamesStatsToJSON_Error はSaveGamesStatsToJSONのエラーテスト
func TestSteamService_SaveGamesStatsToJSON_Error(t *testing.T) {
	// Arrange
	expectedError := errors.New("write error")
	mockFileWriter := &MockFileWriter{
		WriteToFileFunc: func(data any, filename string) error {
			return expectedError
		},
	}

	service := NewSteamService(nil, mockFileWriter, nil, nil)
	gameStats := []*GameStatsInfo{}
	filename := "test_stats.json"

	// Act
	err := service.SaveGamesStatsToJSON(gameStats, filename)

	// Assert
	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}
}

// TestSteamService_GetterMethods はSteamServiceのゲッターメソッドのテスト
func TestSteamService_GetterMethods(t *testing.T) {
	// Arrange
	mockClient := &steamAPI.MockSteamClient{}
	mockFileWriter := &MockFileWriter{}
	mockLogger := &MockLogger{}
	mockConcurrencyController := &MockConcurrencyManager{}

	service := NewSteamService(mockClient, mockFileWriter, mockLogger, mockConcurrencyController)

	// Act & Assert
	if service.GetClient() != mockClient {
		t.Error("Expected GetClient to return the injected client")
	}
	if service.GetFileWriter() != mockFileWriter {
		t.Error("Expected GetFileWriter to return the injected file writer")
	}
	if service.GetLogger() != mockLogger {
		t.Error("Expected GetLogger to return the injected logger")
	}
	if service.GetConcurrencyController() != mockConcurrencyController {
		t.Error("Expected GetConcurrencyController to return the injected concurrency controller")
	}
}

// TestFileWriter_WriteToFile_CreateFileError はFileWriterのファイル作成エラーテスト
func TestFileWriter_WriteToFile_CreateFileError(t *testing.T) {
	// Arrange
	writer := &FileWriter{}
	testData := map[string]string{"test": "data"}
	// 無効なパスを使用してファイル作成エラーを発生させる
	filename := "/invalid/path/test.json"

	// Act
	err := writer.WriteToFile(testData, filename)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid file path, got nil")
	}
}

// TestFileWriter_WriteToFile_EncodeError はFileWriterのエンコードエラーテスト
func TestFileWriter_WriteToFile_EncodeError(t *testing.T) {
	// Arrange
	writer := &FileWriter{}
	// JSONエンコードできないデータ（循環参照）を作成
	type cyclicStruct struct {
		Self *cyclicStruct
	}
	testData := &cyclicStruct{}
	testData.Self = testData
	filename := "/tmp/test_encode_error.json"

	// Act
	err := writer.WriteToFile(testData, filename)

	// Assert
	if err == nil {
		t.Error("Expected encoding error, got nil")
	}
}

// TestBuildGameInfo_WithoutIcon はbuildGameInfoのアイコンなしテスト
func TestBuildGameInfo_WithoutIcon(t *testing.T) {
	// Arrange
	mockClient := &steamAPI.MockSteamClient{
		AppsService: &steamAPI.MockAppsService{
			GetUserStatsFunc: func(ctx context.Context, steamID string, appID int) (*steamAPI.UserStats, error) {
				return nil, errors.New("stats not available")
			},
			GetUserAchievementsFunc: func(ctx context.Context, steamID string, appID int, language string) ([]steamAPI.Achievement, error) {
				return nil, errors.New("achievements not available")
			},
		},
	}

	service := NewSteamService(mockClient, nil, nil, nil)
	ownedGame := steamAPI.OwnedGame{
		AppID:                     123,
		Name:                      "Test Game",
		ImgIconURL:                "", // アイコンなし
		HasCommunityVisibleStats:  true,
	}
	recentGamesMap := make(map[int]steamAPI.RecentlyPlayedGame)
	steamID := "test_steam_id"

	// Act
	gameInfo := service.buildGameInfo(context.Background(), ownedGame, recentGamesMap, steamID)

	// Assert
	if gameInfo.Icon != "" {
		t.Errorf("Expected empty icon, got '%s'", gameInfo.Icon)
	}
	if gameInfo.Stats {
		t.Error("Expected stats to be false due to error")
	}
	if gameInfo.AchievementsCanRetrieve {
		t.Error("Expected achievements to be false due to error")
	}
}

// TestMockConcurrencyManager_DefaultBehavior はMockConcurrencyManagerのデフォルト動作テスト
func TestMockConcurrencyManager_DefaultBehavior(t *testing.T) {
	// Arrange
	controller := &MockConcurrencyManager{}

	// Act & Assert - デフォルト動作をテスト
	semaphore := controller.GetSemaphore()
	if semaphore == nil {
		t.Error("Expected default semaphore to be created")
	}

	maxConcurrency := controller.GetMaxConcurrency()
	if maxConcurrency != 1 {
		t.Errorf("Expected default max concurrency 1, got %d", maxConcurrency)
	}

	// セマフォの容量をテスト
	select {
	case semaphore <- struct{}{}:
		// 成功
	default:
		t.Error("Expected to be able to send to default semaphore")
	}

	// 容量を超えた場合のテスト
	select {
	case semaphore <- struct{}{}:
		t.Error("Expected default semaphore to be full")
	default:
		// 期待される動作
	}
}

// TestMockSteamService_DefaultBehavior はMockSteamServiceのデフォルト動作テスト
func TestMockSteamService_DefaultBehavior(t *testing.T) {
	// Arrange
	mockService := &MockSteamService{}

	// Act & Assert - デフォルト動作をテスト（すべてnilを返す）
	if mockService.GetClient() != nil {
		t.Error("Expected default GetClient to return nil")
	}
	if mockService.GetFileWriter() != nil {
		t.Error("Expected default GetFileWriter to return nil")
	}
	if mockService.GetLogger() != nil {
		t.Error("Expected default GetLogger to return nil")
	}
	if mockService.GetConcurrencyController() != nil {
		t.Error("Expected default GetConcurrencyController to return nil")
	}

	games, err := mockService.GetGamesInfo(context.Background(), "test")
	if games != nil || err != nil {
		t.Errorf("Expected default GetGamesInfo to return nil, nil, got %v, %v", games, err)
	}

	err = mockService.SaveGamesToJSON(nil, "", "")
	if err != nil {
		t.Errorf("Expected default SaveGamesToJSON to return nil, got %v", err)
	}

	stats, err := mockService.GetGamesStats(context.Background(), "test")
	if stats != nil || err != nil {
		t.Errorf("Expected default GetGamesStats to return nil, nil, got %v, %v", stats, err)
	}

	err = mockService.SaveGamesStatsToJSON(nil, "")
	if err != nil {
		t.Errorf("Expected default SaveGamesStatsToJSON to return nil, got %v", err)
	}
}
