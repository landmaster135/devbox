package usecases

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	models "github.com/landmaster135/devbox/internal/db_server_sync/models"
)

type DbServerSyncServiceTest struct {
	service *DbServerSyncService
}

func TestDbServerSyncService_ProcessAppendAnime_Normal(t *testing.T) {
	test := &DbServerSyncServiceTest{
		service: NewDbServerSyncService(),
	}

	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// テスト用のAniListデータを作成
	anilistData := []models.AniListAnime{
		{
			ID:            5114,
			Title:         "鋼の錬金術師 FULLMETAL ALCHEMIST",
			Score:         94,
			Status:        "COMPLETED",
			Progress:      64,
			CompletedAt:   stringPtr("2022-01-13T00:00:00Z"),
			Notes:         "こういう終わり方もクールだと思います。",
			CoverImageURL: "https://s4.anilist.co/file/anilistcdn/media/anime/cover/large/bx5114-nSWCgQlmOMtj.jpg",
			SiteURL:       "https://anilist.co/anime/5114",
			Studio:        "bones",
			UpdatedAt:     "2024-01-22T16:31:54+09:00",
		},
	}

	// テスト用の追加データを作成
	additionalData := []models.AdditionalAnimeData{
		{
			AniListID:             5114,
			ConID:                 "AN0001",
			CoverImageURLModified: "https://storage.googleapis.com/webclip/20230216_notion_mycontents/AN0001_01.webp",
			FleetingTier:          intPtr(0),
			FunnyTier:             intPtr(2),
			HeartwarmingTier:      intPtr(2),
			MotivatingTier:        intPtr(1),
			NihilisticTier:        intPtr(1),
			TearjerkingTier:       intPtr(2),
		},
	}

	// テストファイルを作成
	inputFilePath := filepath.Join(tempDir, "anilist.json")
	additionalFilePath := filepath.Join(tempDir, "additional.json")
	outputFilePath := filepath.Join(tempDir, "output.json")

	// AniListデータをファイルに書き込み
	anilistJSON, err := json.Marshal(anilistData)
	if err != nil {
		t.Fatalf("AniListデータのJSONマーシャルに失敗しました: %v", err)
	}
	if err := os.WriteFile(inputFilePath, anilistJSON, 0644); err != nil {
		t.Fatalf("AniListデータファイルの作成に失敗しました: %v", err)
	}

	// 追加データをファイルに書き込み
	additionalJSON, err := json.Marshal(additionalData)
	if err != nil {
		t.Fatalf("追加データのJSONマーシャルに失敗しました: %v", err)
	}
	if err := os.WriteFile(additionalFilePath, additionalJSON, 0644); err != nil {
		t.Fatalf("追加データファイルの作成に失敗しました: %v", err)
	}

	// ProcessAppendAnimeを実行
	err = test.service.ProcessAppendAnime(inputFilePath, additionalFilePath, outputFilePath)
	if err != nil {
		t.Fatalf("ProcessAppendAnimeが失敗しました: %v", err)
	}

	// 出力ファイルが作成されたことを確認
	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		t.Fatalf("出力ファイルが作成されませんでした: %s", outputFilePath)
	}

	// 出力ファイルの内容を確認
	outputData, err := os.ReadFile(outputFilePath)
	if err != nil {
		t.Fatalf("出力ファイルの読み込みに失敗しました: %v", err)
	}

	var result models.OutputData
	if err := json.Unmarshal(outputData, &result); err != nil {
		t.Fatalf("出力データのJSONアンマーシャルに失敗しました: %v", err)
	}

	// 結果を検証
	if result.Name != "My Anime List" {
		t.Errorf("期待されるName: 'My Anime List', 実際: '%s'", result.Name)
	}
	if result.Description != "Anime data from AniList" {
		t.Errorf("期待されるDescription: 'Anime data from AniList', 実際: '%s'", result.Description)
	}
	if len(result.Data.Animes) != 1 {
		t.Errorf("期待されるアニメ数: 1, 実際: %d", len(result.Data.Animes))
	}

	anime := result.Data.Animes[0]
	if anime.AniListID != 5114 {
		t.Errorf("期待されるAniListID: 5114, 実際: %d", anime.AniListID)
	}
	if anime.Title != "鋼の錬金術師 FULLMETAL ALCHEMIST" {
		t.Errorf("期待されるTitle: '鋼の錬金術師 FULLMETAL ALCHEMIST', 実際: '%s'", anime.Title)
	}
	if anime.ConID != "AN0001" {
		t.Errorf("期待されるConID: 'AN0001', 実際: '%s'", anime.ConID)
	}
	assert.Equal(t, intPtr(2), anime.FunnyTier)
}

func TestDbServerSyncService_ProcessAppendAnime_WithoutAdditionalData_Normal(t *testing.T) {
	test := &DbServerSyncServiceTest{
		service: NewDbServerSyncService(),
	}

	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// テスト用のAniListデータを作成
	anilistData := []models.AniListAnime{
		{
			ID:            11111,
			Title:         "アナザー",
			Score:         0,
			Status:        "PLANNING",
			Progress:      0,
			CompletedAt:   nil,
			Notes:         "",
			CoverImageURL: "https://s4.anilist.co/file/anilistcdn/media/anime/cover/large/bx11111-gvvE5bBYsyFo.png",
			SiteURL:       "https://anilist.co/anime/11111",
			Studio:        "Lantis",
			UpdatedAt:     "2022-04-21T22:02:34+09:00",
		},
	}

	// テストファイルを作成
	inputFilePath := filepath.Join(tempDir, "anilist.json")
	outputFilePath := filepath.Join(tempDir, "output.json")

	// AniListデータをファイルに書き込み
	anilistJSON, err := json.Marshal(anilistData)
	if err != nil {
		t.Fatalf("AniListデータのJSONマーシャルに失敗しました: %v", err)
	}
	if err := os.WriteFile(inputFilePath, anilistJSON, 0644); err != nil {
		t.Fatalf("AniListデータファイルの作成に失敗しました: %v", err)
	}

	// ProcessAppendAnimeを実行（追加データなし）
	err = test.service.ProcessAppendAnime(inputFilePath, "", outputFilePath)
	if err != nil {
		t.Fatalf("ProcessAppendAnimeが失敗しました: %v", err)
	}

	// 出力ファイルが作成されたことを確認
	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		t.Fatalf("出力ファイルが作成されませんでした: %s", outputFilePath)
	}

	// 出力ファイルの内容を確認
	outputData, err := os.ReadFile(outputFilePath)
	if err != nil {
		t.Fatalf("出力ファイルの読み込みに失敗しました: %v", err)
	}

	var result models.OutputData
	if err := json.Unmarshal(outputData, &result); err != nil {
		t.Fatalf("出力データのJSONアンマーシャルに失敗しました: %v", err)
	}

	// 結果を検証
	if len(result.Data.Animes) != 1 {
		t.Errorf("期待されるアニメ数: 1, 実際: %d", len(result.Data.Animes))
	}

	anime := result.Data.Animes[0]
	if anime.AniListID != 11111 {
		t.Errorf("期待されるAniListID: 11111, 実際: %d", anime.AniListID)
	}
	if anime.ConID != "" {
		t.Errorf("期待されるConID: '', 実際: '%s'", anime.ConID)
	}
	assert.Nil(t, anime.FunnyTier)
}

// stringPtr は文字列のポインタを返すヘルパー関数
func stringPtr(s string) *string {
	return &s
}

// intPtr は整数のポインタを返すヘルパー関数
func intPtr(i int) *int {
	return &i
}
