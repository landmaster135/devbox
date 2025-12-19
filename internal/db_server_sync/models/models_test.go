package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAniListAnime_ConvertToOutputAnime_Normal(t *testing.T) {
	// テスト用のAniListAnimeデータを作成
	completedAt := "2022-01-13T00:00:00Z"
	anilistAnime := &AniListAnime{
		ID:            5114,
		Title:         "鋼の錬金術師 FULLMETAL ALCHEMIST",
		Score:         94,
		Status:        "COMPLETED",
		Progress:      64,
		CompletedAt:   &completedAt,
		Notes:         "こういう終わり方もクールだと思います。",
		CoverImageURL: "https://s4.anilist.co/file/anilistcdn/media/anime/cover/large/bx5114-nSWCgQlmOMtj.jpg",
		SiteURL:       "https://anilist.co/anime/5114",
		Studio:        "bones",
		UpdatedAt:     "2024-01-22T07:31:54Z",
	}

	// テスト用の追加データを作成
	additionalData := &AdditionalAnimeData{
		AniListID:             5114,
		ConID:                 "AN0001",
		CoverImageURLModified: "https://storage.googleapis.com/webclip/20230216_notion_mycontents/AN0001_01.webp",
		FleetingTier:          intPtr(0),
		FunnyTier:             intPtr(2),
		HeartwarmingTier:      intPtr(2),
		MotivatingTier:        intPtr(1),
		NihilisticTier:        intPtr(1),
		TearjerkingTier:       intPtr(2),
	}

	// 変換を実行
	result, err := anilistAnime.ConvertToOutputAnime(additionalData)
	require.NoError(t, err)

	// 基本フィールドの検証
	assert.Equal(t, 5114, result.AniListID)
	assert.Equal(t, "鋼の錬金術師 FULLMETAL ALCHEMIST", result.Title)
	assert.Equal(t, 94, result.Score)
	assert.Equal(t, "COMPLETED", result.Status)
	assert.Equal(t, 64, result.Progress)
	assert.Equal(t, "こういう終わり方もクールだと思います。", result.Notes)
	assert.Equal(t, "https://s4.anilist.co/file/anilistcdn/media/anime/cover/large/bx5114-nSWCgQlmOMtj.jpg", result.CoverImageURL)
	assert.Equal(t, "https://anilist.co/anime/5114", result.SiteURL)
	assert.Equal(t, "bones", result.Studio)

	// CompletedAtの検証
	assert.Equal(t, &completedAt, result.CompletedAt)

	// UpdatedAtのタイムスタンプ変換を検証
	expectedTime, _ := time.Parse(time.RFC3339, "2024-01-22T07:31:54Z")
	assert.Equal(t, expectedTime.Unix(), result.UpdatedAt)

	// 追加データの検証
	assert.Equal(t, "AN0001", result.ConID)
	assert.Equal(t, "https://storage.googleapis.com/webclip/20230216_notion_mycontents/AN0001_01.webp", result.CoverImageURLModified)
	assert.Equal(t, intPtr(0), result.FleetingTier)
	assert.Equal(t, intPtr(2), result.FunnyTier)
	assert.Equal(t, intPtr(2), result.HeartwarmingTier)
	assert.Equal(t, intPtr(1), result.MotivatingTier)
	assert.Equal(t, intPtr(1), result.NihilisticTier)
	assert.Equal(t, intPtr(2), result.TearjerkingTier)
}

func TestAniListAnime_ConvertToOutputAnime_WithoutAdditionalData_Normal(t *testing.T) {
	// テスト用のAniListAnimeデータを作成（追加データなし）
	anilistAnime := &AniListAnime{
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
		UpdatedAt:     "2022-04-21T13:02:34Z",
	}

	// 変換を実行（追加データなし）
	result, err := anilistAnime.ConvertToOutputAnime(nil)
	require.NoError(t, err)

	// 基本フィールドの検証
	assert.Equal(t, 11111, result.AniListID)
	assert.Equal(t, "アナザー", result.Title)

	// CompletedAtがnullの場合の検証
	assert.Nil(t, result.CompletedAt)

	// 追加データがない場合のデフォルト値を検証
	assert.Equal(t, "", result.ConID)
	assert.Equal(t, "", result.CoverImageURLModified)
	assert.Nil(t, result.FleetingTier)
	assert.Nil(t, result.FunnyTier)
	assert.Nil(t, result.HeartwarmingTier)
	assert.Nil(t, result.MotivatingTier)
	assert.Nil(t, result.NihilisticTier)
	assert.Nil(t, result.TearjerkingTier)
}

func TestAniListAnime_ConvertToOutputAnime_WithNullTiers_Normal(t *testing.T) {
	// テスト用のAniListAnimeデータを作成
	anilistAnime := &AniListAnime{
		ID:            138425,
		Title:         "テストアニメ",
		Score:         85,
		Status:        "COMPLETED",
		Progress:      12,
		CompletedAt:   nil,
		Notes:         "",
		CoverImageURL: "https://example.com/cover.jpg",
		SiteURL:       "https://anilist.co/anime/138425",
		Studio:        "TestStudio",
		UpdatedAt:     "2023-01-01T00:00:00Z",
	}

	// テスト用の追加データを作成（一部のtierがnull）
	additionalData := &AdditionalAnimeData{
		AniListID:             138425,
		ConID:                 "AN0610",
		CoverImageURLModified: "https://storage.googleapis.com/webclip/20230216_notion_mycontents/AN0610_01.webp",
		FleetingTier:          nil,
		FunnyTier:             nil,
		HeartwarmingTier:      nil,
		MotivatingTier:        nil,
		NihilisticTier:        nil,
		TearjerkingTier:       nil,
	}

	// 変換を実行
	result, err := anilistAnime.ConvertToOutputAnime(additionalData)
	require.NoError(t, err)

	// nullのtierはnullのまま出力されることを検証
	assert.Nil(t, result.FleetingTier)
	assert.Nil(t, result.FunnyTier)
	assert.Nil(t, result.HeartwarmingTier)
	assert.Nil(t, result.MotivatingTier)
	assert.Nil(t, result.NihilisticTier)
	assert.Nil(t, result.TearjerkingTier)

	// その他のフィールドは正常に設定されることを検証
	assert.Equal(t, "AN0610", result.ConID)
}

// intPtr は整数のポインタを返すヘルパー関数
func intPtr(i int) *int {
	return &i
}
