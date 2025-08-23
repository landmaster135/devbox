package domain

import (
	"encoding/json"
	"testing"
	"time"
)

// TestAniListResponse_JSONMarshal はAniListResponseのJSONマーシャルテスト
func TestAniListResponse_JSONMarshal(t *testing.T) {
	// Arrange
	response := AniListResponse{
		Data: &MediaListCollectionData{
			MediaListCollection: &MediaListCollection{
				Lists: []MediaList{
					{
						Entries: []MediaListEntry{
							{
								Media: &Media{
									ID: 1,
									Title: &Title{
										Native: "テストアニメ",
									},
									CoverImage: &Image{
										ExtraLarge: "https://example.com/image.jpg",
									},
									SiteURL: "https://anilist.co/anime/1",
									Studios: &Studios{
										Nodes: []Studio{
											{Name: "テストスタジオ"},
										},
									},
								},
								Score:    85,
								Status:   "COMPLETED",
								Progress: 12,
								CompletedAt: &FuzzyDate{
									Year:  func() *int { y := 2023; return &y }(),
									Month: func() *int { m := 6; return &m }(),
									Day:   func() *int { d := 15; return &d }(),
								},
								Notes:     "面白かった",
								UpdatedAt: 1687123200,
							},
						},
					},
				},
			},
		},
		Errors: []GraphQLError{},
	}

	// Act
	jsonData, err := json.Marshal(response)

	// Assert
	if err != nil {
		t.Errorf("JSONマーシャルでエラーが発生しました: %v", err)
		return
	}

	// JSONが有効であることを確認
	var unmarshaled AniListResponse
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Errorf("JSONアンマーシャルでエラーが発生しました: %v", err)
		return
	}

	// データの整合性を確認
	if unmarshaled.Data == nil {
		t.Error("データがnilです")
		return
	}
	if unmarshaled.Data.MediaListCollection == nil {
		t.Error("MediaListCollectionがnilです")
		return
	}
	if len(unmarshaled.Data.MediaListCollection.Lists) != 1 {
		t.Errorf("期待されるリスト数: 1, 実際: %d", len(unmarshaled.Data.MediaListCollection.Lists))
	}
}

// TestGraphQLError_JSONMarshal はGraphQLErrorのJSONマーシャルテスト
func TestGraphQLError_JSONMarshal(t *testing.T) {
	// Arrange
	error := GraphQLError{
		Message: "User not found",
		Locations: []GraphQLErrorLocation{
			{Line: 2, Column: 3},
		},
		Path: []any{"MediaListCollection"},
	}

	// Act
	jsonData, err := json.Marshal(error)

	// Assert
	if err != nil {
		t.Errorf("JSONマーシャルでエラーが発生しました: %v", err)
		return
	}

	// JSONアンマーシャルで確認
	var unmarshaled GraphQLError
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Errorf("JSONアンマーシャルでエラーが発生しました: %v", err)
		return
	}

	if unmarshaled.Message != error.Message {
		t.Errorf("期待されるメッセージ: %s, 実際: %s", error.Message, unmarshaled.Message)
	}
}

// TestAnimeInfo_JSONMarshal はAnimeInfoのJSONマーシャルテスト
func TestAnimeInfo_JSONMarshal(t *testing.T) {
	// Arrange
	completedAt := time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 6, 20, 12, 0, 0, 0, time.UTC)

	animeInfo := AnimeInfo{
		ID:            1,
		Title:         "テストアニメ",
		Score:         85,
		Status:        "COMPLETED",
		Progress:      12,
		CompletedAt:   &completedAt,
		Notes:         "面白かった",
		CoverImageURL: "https://example.com/image.jpg",
		SiteURL:       "https://anilist.co/anime/1",
		Studio:        "テストスタジオ",
		UpdatedAt:     updatedAt,
	}

	// Act
	jsonData, err := json.Marshal(animeInfo)

	// Assert
	if err != nil {
		t.Errorf("JSONマーシャルでエラーが発生しました: %v", err)
		return
	}

	// JSONアンマーシャルで確認
	var unmarshaled AnimeInfo
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Errorf("JSONアンマーシャルでエラーが発生しました: %v", err)
		return
	}

	// データの整合性を確認
	if unmarshaled.ID != animeInfo.ID {
		t.Errorf("期待されるID: %d, 実際: %d", animeInfo.ID, unmarshaled.ID)
	}
	if unmarshaled.Title != animeInfo.Title {
		t.Errorf("期待されるタイトル: %s, 実際: %s", animeInfo.Title, unmarshaled.Title)
	}
	if unmarshaled.Score != animeInfo.Score {
		t.Errorf("期待されるスコア: %d, 実際: %d", animeInfo.Score, unmarshaled.Score)
	}
}

// TestFuzzyDate_JSONMarshal はFuzzyDateのJSONマーシャルテスト
func TestFuzzyDate_JSONMarshal(t *testing.T) {
	tests := []struct {
		name string
		date FuzzyDate
	}{
		{
			name: "FullDate_Normal",
			date: FuzzyDate{
				Year:  func() *int { y := 2023; return &y }(),
				Month: func() *int { m := 6; return &m }(),
				Day:   func() *int { d := 15; return &d }(),
			},
		},
		{
			name: "YearMonthOnly_Normal",
			date: FuzzyDate{
				Year:  func() *int { y := 2023; return &y }(),
				Month: func() *int { m := 6; return &m }(),
				Day:   nil,
			},
		},
		{
			name: "YearOnly_Normal",
			date: FuzzyDate{
				Year:  func() *int { y := 2023; return &y }(),
				Month: nil,
				Day:   nil,
			},
		},
		{
			name: "AllNil_Normal",
			date: FuzzyDate{
				Year:  nil,
				Month: nil,
				Day:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			jsonData, err := json.Marshal(tt.date)

			// Assert
			if err != nil {
				t.Errorf("JSONマーシャルでエラーが発生しました: %v", err)
				return
			}

			// JSONアンマーシャルで確認
			var unmarshaled FuzzyDate
			if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
				t.Errorf("JSONアンマーシャルでエラーが発生しました: %v", err)
				return
			}

			// データの整合性を確認
			if (tt.date.Year == nil) != (unmarshaled.Year == nil) {
				t.Error("Yearのnil状態が異なります")
			} else if tt.date.Year != nil && unmarshaled.Year != nil && *tt.date.Year != *unmarshaled.Year {
				t.Errorf("期待されるYear: %d, 実際: %d", *tt.date.Year, *unmarshaled.Year)
			}

			if (tt.date.Month == nil) != (unmarshaled.Month == nil) {
				t.Error("Monthのnil状態が異なります")
			} else if tt.date.Month != nil && unmarshaled.Month != nil && *tt.date.Month != *unmarshaled.Month {
				t.Errorf("期待されるMonth: %d, 実際: %d", *tt.date.Month, *unmarshaled.Month)
			}

			if (tt.date.Day == nil) != (unmarshaled.Day == nil) {
				t.Error("Dayのnil状態が異なります")
			} else if tt.date.Day != nil && unmarshaled.Day != nil && *tt.date.Day != *unmarshaled.Day {
				t.Errorf("期待されるDay: %d, 実際: %d", *tt.date.Day, *unmarshaled.Day)
			}
		})
	}
}

// TestGraphQLRequest_JSONMarshal はGraphQLRequestのJSONマーシャルテスト
func TestGraphQLRequest_JSONMarshal(t *testing.T) {
	// Arrange
	request := GraphQLRequest{
		Query: "query { MediaListCollection { lists { entries { media { id } } } } }",
		Variables: map[string]any{
			"username": "testuser",
			"id":       12345,
		},
	}

	// Act
	jsonData, err := json.Marshal(request)

	// Assert
	if err != nil {
		t.Errorf("JSONマーシャルでエラーが発生しました: %v", err)
		return
	}

	// JSONアンマーシャルで確認
	var unmarshaled GraphQLRequest
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Errorf("JSONアンマーシャルでエラーが発生しました: %v", err)
		return
	}

	if unmarshaled.Query != request.Query {
		t.Errorf("期待されるクエリ: %s, 実際: %s", request.Query, unmarshaled.Query)
	}
	if len(unmarshaled.Variables) != len(request.Variables) {
		t.Errorf("期待される変数数: %d, 実際: %d", len(request.Variables), len(unmarshaled.Variables))
	}
}

// TestQueryAnimeRequest_Initialization はQueryAnimeRequestの初期化テスト
func TestQueryAnimeRequest_Initialization(t *testing.T) {
	tests := []struct {
		name     string
		username string
		userID   *int
	}{
		{
			name:     "WithUsername_Normal",
			username: "testuser",
			userID:   nil,
		},
		{
			name:     "WithUserID_Normal",
			username: "",
			userID:   func() *int { id := 12345; return &id }(),
		},
		{
			name:     "WithBoth_Normal",
			username: "testuser",
			userID:   func() *int { id := 12345; return &id }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			request := QueryAnimeRequest{
				Username: tt.username,
				UserID:   tt.userID,
			}

			// Assert
			if request.Username != tt.username {
				t.Errorf("期待されるユーザー名: %s, 実際: %s", tt.username, request.Username)
			}
			if (request.UserID == nil) != (tt.userID == nil) {
				t.Error("ユーザーIDのnil状態が異なります")
			} else if request.UserID != nil && tt.userID != nil && *request.UserID != *tt.userID {
				t.Errorf("期待されるユーザーID: %d, 実際: %d", *tt.userID, *request.UserID)
			}
		})
	}
}

// TestComplexStructure_JSONMarshal は複雑な構造体のJSONマーシャルテスト
func TestComplexStructure_JSONMarshal(t *testing.T) {
	// Arrange
	response := AniListResponse{
		Data: &MediaListCollectionData{
			MediaListCollection: &MediaListCollection{
				Lists: []MediaList{
					{
						Entries: []MediaListEntry{
							{
								Media: &Media{
									ID: 1,
									Title: &Title{
										Native: "進撃の巨人",
									},
									CoverImage: &Image{
										ExtraLarge: "https://s4.anilist.co/file/anilistcdn/media/anime/cover/large/bx16498-73IhOXpJZiMF.jpg",
									},
									SiteURL: "https://anilist.co/anime/16498",
									Studios: &Studios{
										Nodes: []Studio{
											{Name: "Mappa"},
											{Name: "WIT Studio"},
										},
									},
								},
								Score:    95,
								Status:   "COMPLETED",
								Progress: 25,
								CompletedAt: &FuzzyDate{
									Year:  func() *int { y := 2023; return &y }(),
									Month: func() *int { m := 4; return &m }(),
									Day:   func() *int { d := 4; return &d }(),
								},
								Notes:     "最高のアニメ",
								UpdatedAt: 1680566400,
							},
							{
								Media: &Media{
									ID: 2,
									Title: &Title{
										Native: "鬼滅の刃",
									},
									CoverImage: &Image{
										ExtraLarge: "https://s4.anilist.co/file/anilistcdn/media/anime/cover/large/bx101922-PEn1CTc93blC.jpg",
									},
									SiteURL: "https://anilist.co/anime/101922",
									Studios: &Studios{
										Nodes: []Studio{
											{Name: "Ufotable"},
										},
									},
								},
								Score:    90,
								Status:   "COMPLETED",
								Progress: 26,
								CompletedAt: &FuzzyDate{
									Year:  func() *int { y := 2023; return &y }(),
									Month: func() *int { m := 5; return &m }(),
									Day:   func() *int { d := 28; return &d }(),
								},
								Notes:     "美しいアニメーション",
								UpdatedAt: 1685232000,
							},
						},
					},
				},
			},
		},
		Errors: []GraphQLError{},
	}

	// Act
	jsonData, err := json.Marshal(response)

	// Assert
	if err != nil {
		t.Errorf("JSONマーシャルでエラーが発生しました: %v", err)
		return
	}

	// JSONアンマーシャルで確認
	var unmarshaled AniListResponse
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Errorf("JSONアンマーシャルでエラーが発生しました: %v", err)
		return
	}

	// データの整合性を確認
	if unmarshaled.Data == nil || unmarshaled.Data.MediaListCollection == nil {
		t.Error("データ構造が正しくありません")
		return
	}

	entries := unmarshaled.Data.MediaListCollection.Lists[0].Entries
	if len(entries) != 2 {
		t.Errorf("期待されるエントリ数: 2, 実際: %d", len(entries))
		return
	}

	// 最初のエントリを確認
	firstEntry := entries[0]
	if firstEntry.Media.Title.Native != "進撃の巨人" {
		t.Errorf("期待されるタイトル: 進撃の巨人, 実際: %s", firstEntry.Media.Title.Native)
	}
	if firstEntry.Score != 95 {
		t.Errorf("期待されるスコア: 95, 実際: %d", firstEntry.Score)
	}
	if len(firstEntry.Media.Studios.Nodes) != 2 {
		t.Errorf("期待されるスタジオ数: 2, 実際: %d", len(firstEntry.Media.Studios.Nodes))
	}
}

// TestEmptyStructures_JSONMarshal は空の構造体のJSONマーシャルテスト
func TestEmptyStructures_JSONMarshal(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "EmptyAniListResponse",
			data: AniListResponse{},
		},
		{
			name: "EmptyMediaListCollection",
			data: MediaListCollection{},
		},
		{
			name: "EmptyAnimeInfo",
			data: AnimeInfo{},
		},
		{
			name: "EmptyGraphQLRequest",
			data: GraphQLRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			jsonData, err := json.Marshal(tt.data)

			// Assert
			if err != nil {
				t.Errorf("JSONマーシャルでエラーが発生しました: %v", err)
				return
			}

			// JSONが有効であることを確認
			var result map[string]interface{}
			if err := json.Unmarshal(jsonData, &result); err != nil {
				t.Errorf("JSONアンマーシャルでエラーが発生しました: %v", err)
			}
		})
	}
}
