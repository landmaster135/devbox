package usecases

import (
	"errors"
	"os"
	"strings"
	"testing"

	domain "github.com/landmaster135/devbox/internal/anilist/domain"
	infrastructure "github.com/landmaster135/devbox/internal/anilist/infrastructure"
)

// TestAniListService_QueryAnime_Normal はQueryAnimeメソッドの正常系テスト
func TestAniListService_QueryAnime_Normal(t *testing.T) {
	// Arrange
	var queryAnimeListCalled bool
	var queryAnimeListRequest domain.QueryAnimeRequest

	mockRepository := &domain.MockAniListRepository{
		QueryAnimeListFunc: func(req domain.QueryAnimeRequest) (*domain.AniListResponse, error) {
			queryAnimeListCalled = true
			queryAnimeListRequest = req

			year := 2023
			month := 12
			day := 25
			return &domain.AniListResponse{
				Data: &domain.MediaListCollectionData{
					MediaListCollection: &domain.MediaListCollection{
						Lists: []domain.MediaList{
							{
								Entries: []domain.MediaListEntry{
									{
										Media: &domain.Media{
											ID: 1,
											Title: &domain.Title{
												Native: "テストアニメ",
											},
											CoverImage: &domain.Image{
												ExtraLarge: "https://example.com/image.jpg",
											},
											SiteURL: "https://anilist.co/anime/1",
											Studios: &domain.Studios{
												Nodes: []domain.Studio{
													{Name: "テストスタジオ"},
												},
											},
										},
										Score:    9,
										Status:   "COMPLETED",
										Progress: 12,
										Repeat:   2,
										CompletedAt: &domain.FuzzyDate{
											Year:  &year,
											Month: &month,
											Day:   &day,
										},
										Notes:     "テストノート",
										UpdatedAt: 1703505600,
									},
								},
							},
						},
					},
				},
			}, nil
		},
	}

	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{
		MarshalIndentFunc: func(v any, prefix, indent string) ([]byte, error) {
			return []byte(`[{"id":1,"title":"テストアニメ"}]`), nil
		},
	}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryAnime("testuser", nil, "json", 0, "", "")

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if !queryAnimeListCalled {
		t.Error("QueryAnimeListが呼び出されませんでした")
	}
	if queryAnimeListRequest.Username != "testuser" {
		t.Errorf("期待されるユーザー名: testuser, 実際: %s", queryAnimeListRequest.Username)
	}
	if queryAnimeListRequest.UserID != nil {
		t.Errorf("期待されるユーザーID: nil, 実際: %v", queryAnimeListRequest.UserID)
	}
	if result == "" {
		t.Error("結果が空文字列です")
	}
	if !strings.Contains(result, "テストアニメ") {
		t.Error("結果にアニメタイトルが含まれていません")
	}
}

// TestAniListService_QueryAnime_WithUserID はQueryAnimeメソッドのユーザーID指定テスト
func TestAniListService_QueryAnime_WithUserID(t *testing.T) {
	// Arrange
	var queryAnimeListRequest domain.QueryAnimeRequest
	userID := 12345

	mockRepository := &domain.MockAniListRepository{
		QueryAnimeListFunc: func(req domain.QueryAnimeRequest) (*domain.AniListResponse, error) {
			queryAnimeListRequest = req
			return &domain.AniListResponse{
				Data: &domain.MediaListCollectionData{
					MediaListCollection: &domain.MediaListCollection{
						Lists: []domain.MediaList{
							{
								Entries: []domain.MediaListEntry{
									{
										Media: &domain.Media{
											ID: 1,
											Title: &domain.Title{
												Native: "テストアニメ",
											},
										},
										Score:    8,
										Status:   "CURRENT",
										Progress: 5,
										Repeat:   1,
									},
								},
							},
						},
					},
				},
			}, nil
		},
	}

	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{
		MarshalIndentFunc: func(v any, prefix, indent string) ([]byte, error) {
			return []byte(`[{"id":1,"title":"テストアニメ"}]`), nil
		},
	}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryAnime("", &userID, "json", 0, "", "")

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if queryAnimeListRequest.Username != "" {
		t.Errorf("期待されるユーザー名: 空文字列, 実際: %s", queryAnimeListRequest.Username)
	}
	if queryAnimeListRequest.UserID == nil {
		t.Error("ユーザーIDがnilです")
	} else if *queryAnimeListRequest.UserID != 12345 {
		t.Errorf("期待されるユーザーID: 12345, 実際: %d", *queryAnimeListRequest.UserID)
	}
	if result == "" {
		t.Error("結果が空文字列です")
	}
}

// TestAniListService_QueryAnime_APIError はQueryAnimeメソッドのAPIエラーテスト
func TestAniListService_QueryAnime_APIError(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryAnimeListFunc: func(req domain.QueryAnimeRequest) (*domain.AniListResponse, error) {
			return nil, errors.New("API呼び出しエラー")
		},
	}

	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryAnime("testuser", nil, "json", 0, "", "")

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	expectedError := "AniList APIの呼び出しに失敗しました: API呼び出しエラー"
	if err.Error() != expectedError {
		t.Errorf("期待されるエラーメッセージ: %s, 実際: %s", expectedError, err.Error())
	}
	if result != "" {
		t.Errorf("エラー時は空文字列が期待されます, 実際: %s", result)
	}
}

// TestAniListService_QueryAnime_NoData はQueryAnimeメソッドのデータなしエラーテスト
func TestAniListService_QueryAnime_NoData(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryAnimeListFunc: func(req domain.QueryAnimeRequest) (*domain.AniListResponse, error) {
			return &domain.AniListResponse{
				Data: nil, // データがnil
			}, nil
		},
	}

	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryAnime("testuser", nil, "json", 0, "", "")

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	expectedError := "ユーザーのアニメリストが見つかりません"
	if err.Error() != expectedError {
		t.Errorf("期待されるエラーメッセージ: %s, 実際: %s", expectedError, err.Error())
	}
	if result != "" {
		t.Errorf("エラー時は空文字列が期待されます, 実際: %s", result)
	}
}

// TestAniListService_QueryAnime_NoMediaListCollection はQueryAnimeメソッドのMediaListCollectionなしエラーテスト
func TestAniListService_QueryAnime_NoMediaListCollection(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryAnimeListFunc: func(req domain.QueryAnimeRequest) (*domain.AniListResponse, error) {
			return &domain.AniListResponse{
				Data: &domain.MediaListCollectionData{
					MediaListCollection: nil, // MediaListCollectionがnil
				},
			}, nil
		},
	}

	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryAnime("testuser", nil, "json", 0, "", "")

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	expectedError := "ユーザーのアニメリストが見つかりません"
	if err.Error() != expectedError {
		t.Errorf("期待されるエラーメッセージ: %s, 実際: %s", expectedError, err.Error())
	}
	if result != "" {
		t.Errorf("エラー時は空文字列が期待されます, 実際: %s", result)
	}
}

// TestAniListService_QueryAnime_WithStatusFilter はQueryAnimeメソッドのステータスフィルタテスト
func TestAniListService_QueryAnime_WithStatusFilter(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryAnimeListFunc: func(req domain.QueryAnimeRequest) (*domain.AniListResponse, error) {
			return &domain.AniListResponse{
				Data: &domain.MediaListCollectionData{
					MediaListCollection: &domain.MediaListCollection{
						Lists: []domain.MediaList{
							{
								Entries: []domain.MediaListEntry{
									{
										Media: &domain.Media{
											ID: 1,
											Title: &domain.Title{
												Native: "完了アニメ",
											},
										},
										Score:    9,
										Status:   "COMPLETED",
										Progress: 12,
										Repeat:   3,
									},
									{
										Media: &domain.Media{
											ID: 2,
											Title: &domain.Title{
												Native: "視聴中アニメ",
											},
										},
										Score:    8,
										Status:   "CURRENT",
										Progress: 5,
										Repeat:   0,
									},
								},
							},
						},
					},
				},
			}, nil
		},
	}

	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{
		MarshalIndentFunc: func(v any, prefix, indent string) ([]byte, error) {
			animeList := v.([]domain.AnimeInfo)
			if len(animeList) != 1 {
				t.Errorf("フィルタ後の期待される件数: 1, 実際: %d", len(animeList))
			}
			if animeList[0].Status != "COMPLETED" {
				t.Errorf("フィルタされたアニメのステータス: COMPLETED, 実際: %s", animeList[0].Status)
			}
			return []byte(`[{"id":1,"title":"完了アニメ","status":"COMPLETED"}]`), nil
		},
	}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryAnime("testuser", nil, "json", 0, "COMPLETED", "")

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == "" {
		t.Error("結果が空文字列です")
	}
	if !strings.Contains(result, "完了アニメ") {
		t.Error("フィルタされたアニメが含まれていません")
	}
}

// TestAniListService_QueryAnime_WithLimit はQueryAnimeメソッドの制限値テスト
func TestAniListService_QueryAnime_WithLimit(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryAnimeListFunc: func(req domain.QueryAnimeRequest) (*domain.AniListResponse, error) {
			return &domain.AniListResponse{
				Data: &domain.MediaListCollectionData{
					MediaListCollection: &domain.MediaListCollection{
						Lists: []domain.MediaList{
							{
								Entries: []domain.MediaListEntry{
									{
										Media: &domain.Media{
											ID: 1,
											Title: &domain.Title{
												Native: "アニメ1",
											},
										},
										Score:    9,
										Status:   "COMPLETED",
										Progress: 12,
										Repeat:   1,
									},
									{
										Media: &domain.Media{
											ID: 2,
											Title: &domain.Title{
												Native: "アニメ2",
											},
										},
										Score:    8,
										Status:   "CURRENT",
										Progress: 5,
										Repeat:   0,
									},
									{
										Media: &domain.Media{
											ID: 3,
											Title: &domain.Title{
												Native: "アニメ3",
											},
										},
										Score:    7,
										Status:   "PLANNING",
										Progress: 0,
										Repeat:   0,
									},
								},
							},
						},
					},
				},
			}, nil
		},
	}

	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{
		MarshalIndentFunc: func(v any, prefix, indent string) ([]byte, error) {
			animeList := v.([]domain.AnimeInfo)
			if len(animeList) != 2 {
				t.Errorf("制限後の期待される件数: 2, 実際: %d", len(animeList))
			}
			return []byte(`[{"id":1,"title":"アニメ1"},{"id":2,"title":"アニメ2"}]`), nil
		},
	}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryAnime("testuser", nil, "json", 2, "", "")

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == "" {
		t.Error("結果が空文字列です")
	}
}

// TestAniListService_QueryAnime_TableFormat はQueryAnimeメソッドのテーブル形式テスト
func TestAniListService_QueryAnime_TableFormat(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryAnimeListFunc: func(req domain.QueryAnimeRequest) (*domain.AniListResponse, error) {
			return &domain.AniListResponse{
				Data: &domain.MediaListCollectionData{
					MediaListCollection: &domain.MediaListCollection{
						Lists: []domain.MediaList{
							{
								Entries: []domain.MediaListEntry{
									{
										Media: &domain.Media{
											ID: 1,
											Title: &domain.Title{
												Native: "テストアニメ",
											},
											Studios: &domain.Studios{
												Nodes: []domain.Studio{
													{Name: "テストスタジオ"},
												},
											},
										},
										Score:    9,
										Status:   "COMPLETED",
										Progress: 12,
										Repeat:   2,
									},
								},
							},
						},
					},
				},
			}, nil
		},
	}

	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryAnime("testuser", nil, "table", 0, "", "")

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == "" {
		t.Error("結果が空文字列です")
	}
	// テーブル形式の確認
	if !strings.Contains(result, "ID\tタイトル\tステータス\tスコア\t進行状況\t再視聴回数") {
		t.Error("テーブルヘッダーが含まれていません")
	}
	if !strings.Contains(result, "テストアニメ") {
		t.Error("アニメタイトルが含まれていません")
	}
	if !strings.Contains(result, "\t2\t") {
		t.Error("再視聴回数が含まれていません")
	}
}

// TestAniListService_QueryAnime_WithOutputDir はQueryAnimeメソッドのファイル出力テスト
func TestAniListService_QueryAnime_WithOutputDir(t *testing.T) {
	// Arrange
	var mkdirAllCalled bool
	var writeFileCalled bool
	var writtenContent []byte

	mockRepository := &domain.MockAniListRepository{
		QueryAnimeListFunc: func(req domain.QueryAnimeRequest) (*domain.AniListResponse, error) {
			return &domain.AniListResponse{
				Data: &domain.MediaListCollectionData{
					MediaListCollection: &domain.MediaListCollection{
						Lists: []domain.MediaList{
							{
								Entries: []domain.MediaListEntry{
									{
										Media: &domain.Media{
											ID: 1,
											Title: &domain.Title{
												Native: "テストアニメ",
											},
										},
										Score:    9,
										Status:   "COMPLETED",
										Progress: 12,
										Repeat:   1,
									},
								},
							},
						},
					},
				},
			}, nil
		},
	}

	mockFS := &infrastructure.MockFileSystem{
		MkdirAllFunc: func(path string, perm os.FileMode) error {
			mkdirAllCalled = true
			return nil
		},
		WriteFileFunc: func(filename string, data []byte, perm os.FileMode) error {
			writeFileCalled = true
			writtenContent = data
			return nil
		},
		JoinFunc: func(elem ...string) string {
			return elem[0] + "/" + elem[1]
		},
	}

	mockJSON := &infrastructure.MockJSONProcessor{
		MarshalIndentFunc: func(v any, prefix, indent string) ([]byte, error) {
			return []byte(`[{"id":1,"title":"テストアニメ"}]`), nil
		},
	}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryAnime("testuser", nil, "json", 0, "", "/output")

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if !mkdirAllCalled {
		t.Error("MkdirAllが呼び出されませんでした")
	}
	if !writeFileCalled {
		t.Error("WriteFileが呼び出されませんでした")
	}
	if string(writtenContent) != `[{"id":1,"title":"テストアニメ"}]` {
		t.Errorf("期待される書き込み内容と異なります: %s", string(writtenContent))
	}
	if !strings.Contains(result, "結果をファイルに保存しました") {
		t.Error("ファイル保存メッセージが含まれていません")
	}
}

// TestAniListService_QueryAnime_FileSaveError はQueryAnimeメソッドのファイル保存エラーテスト
func TestAniListService_QueryAnime_FileSaveError(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryAnimeListFunc: func(req domain.QueryAnimeRequest) (*domain.AniListResponse, error) {
			return &domain.AniListResponse{
				Data: &domain.MediaListCollectionData{
					MediaListCollection: &domain.MediaListCollection{
						Lists: []domain.MediaList{
							{
								Entries: []domain.MediaListEntry{
									{
										Media: &domain.Media{
											ID: 1,
											Title: &domain.Title{
												Native: "テストアニメ",
											},
										},
										Score:    9,
										Status:   "COMPLETED",
										Progress: 12,
									},
								},
							},
						},
					},
				},
			}, nil
		},
	}

	mockFS := &infrastructure.MockFileSystem{
		MkdirAllFunc: func(path string, perm os.FileMode) error {
			return errors.New("ディレクトリ作成エラー")
		},
	}

	mockJSON := &infrastructure.MockJSONProcessor{
		MarshalIndentFunc: func(v any, prefix, indent string) ([]byte, error) {
			return []byte(`[{"id":1,"title":"テストアニメ"}]`), nil
		},
	}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryAnime("testuser", nil, "json", 0, "", "/output")

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	expectedError := "ファイル保存に失敗しました: 出力ディレクトリの作成に失敗しました: ディレクトリ作成エラー"
	if err.Error() != expectedError {
		t.Errorf("期待されるエラーメッセージ: %s, 実際: %s", expectedError, err.Error())
	}
	if result != "" {
		t.Errorf("エラー時は空文字列が期待されます, 実際: %s", result)
	}
}

// TestAniListService_QueryAnime_JSONFormatError はQueryAnimeメソッドのJSON形式エラーテスト
func TestAniListService_QueryAnime_JSONFormatError(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryAnimeListFunc: func(req domain.QueryAnimeRequest) (*domain.AniListResponse, error) {
			return &domain.AniListResponse{
				Data: &domain.MediaListCollectionData{
					MediaListCollection: &domain.MediaListCollection{
						Lists: []domain.MediaList{
							{
								Entries: []domain.MediaListEntry{
									{
										Media: &domain.Media{
											ID: 1,
											Title: &domain.Title{
												Native: "テストアニメ",
											},
										},
										Score:    9,
										Status:   "COMPLETED",
										Progress: 12,
									},
								},
							},
						},
					},
				},
			}, nil
		},
	}

	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{
		MarshalIndentFunc: func(v any, prefix, indent string) ([]byte, error) {
			return nil, errors.New("JSONエンコードエラー")
		},
	}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryAnime("testuser", nil, "json", 0, "", "")

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	expectedError := "JSONエンコードに失敗しました: JSONエンコードエラー"
	if err.Error() != expectedError {
		t.Errorf("期待されるエラーメッセージ: %s, 実際: %s", expectedError, err.Error())
	}
	if result != "" {
		t.Errorf("エラー時は空文字列が期待されます, 実際: %s", result)
	}
}

// TestAniListService_QueryAnime_DefaultFormat はQueryAnimeメソッドのデフォルト形式テスト
func TestAniListService_QueryAnime_DefaultFormat(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryAnimeListFunc: func(req domain.QueryAnimeRequest) (*domain.AniListResponse, error) {
			return &domain.AniListResponse{
				Data: &domain.MediaListCollectionData{
					MediaListCollection: &domain.MediaListCollection{
						Lists: []domain.MediaList{
							{
								Entries: []domain.MediaListEntry{
									{
										Media: &domain.Media{
											ID: 1,
											Title: &domain.Title{
												Native: "テストアニメ",
											},
										},
										Score:    9,
										Status:   "COMPLETED",
										Progress: 12,
									},
								},
							},
						},
					},
				},
			}, nil
		},
	}

	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{
		MarshalIndentFunc: func(v any, prefix, indent string) ([]byte, error) {
			return []byte(`[{"id":1,"title":"テストアニメ"}]`), nil
		},
	}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act - 無効な形式を指定してデフォルト（JSON）が使用されることを確認
	result, err := service.QueryAnime("testuser", nil, "invalid-format", 0, "", "")

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == "" {
		t.Error("結果が空文字列です")
	}
	// JSON形式で出力されることを確認
	if !strings.Contains(result, `"id":1`) {
		t.Error("JSON形式で出力されていません")
	}
}

// TestAniListService_QueryAnime_WithRepeatField はQueryAnimeメソッドのRepeatフィールド特化テスト
func TestAniListService_QueryAnime_WithRepeatField(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryAnimeListFunc: func(req domain.QueryAnimeRequest) (*domain.AniListResponse, error) {
			return &domain.AniListResponse{
				Data: &domain.MediaListCollectionData{
					MediaListCollection: &domain.MediaListCollection{
						Lists: []domain.MediaList{
							{
								Entries: []domain.MediaListEntry{
									{
										Media: &domain.Media{
											ID: 1,
											Title: &domain.Title{
												Native: "再視聴なしアニメ",
											},
										},
										Score:    8,
										Status:   "COMPLETED",
										Progress: 12,
										Repeat:   0,
									},
									{
										Media: &domain.Media{
											ID: 2,
											Title: &domain.Title{
												Native: "1回再視聴アニメ",
											},
										},
										Score:    9,
										Status:   "COMPLETED",
										Progress: 24,
										Repeat:   1,
									},
									{
										Media: &domain.Media{
											ID: 3,
											Title: &domain.Title{
												Native: "複数回再視聴アニメ",
											},
										},
										Score:    10,
										Status:   "COMPLETED",
										Progress: 12,
										Repeat:   5,
									},
								},
							},
						},
					},
				},
			}, nil
		},
	}

	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{
		MarshalIndentFunc: func(v any, prefix, indent string) ([]byte, error) {
			animeList := v.([]domain.AnimeInfo)
			// Repeatフィールドが正しく設定されていることを確認
			if len(animeList) != 3 {
				t.Errorf("期待されるアニメ数: 3, 実際: %d", len(animeList))
			}
			if animeList[0].Repeat != 0 {
				t.Errorf("1番目のアニメのRepeat値: 0, 実際: %d", animeList[0].Repeat)
			}
			if animeList[1].Repeat != 1 {
				t.Errorf("2番目のアニメのRepeat値: 1, 実際: %d", animeList[1].Repeat)
			}
			if animeList[2].Repeat != 5 {
				t.Errorf("3番目のアニメのRepeat値: 5, 実際: %d", animeList[2].Repeat)
			}
			return []byte(`[{"id":1,"repeat":0},{"id":2,"repeat":1},{"id":3,"repeat":5}]`), nil
		},
	}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryAnime("testuser", nil, "json", 0, "", "")

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == "" {
		t.Error("結果が空文字列です")
	}
	if !strings.Contains(result, `"repeat":0`) {
		t.Error("Repeat値0が含まれていません")
	}
	if !strings.Contains(result, `"repeat":1`) {
		t.Error("Repeat値1が含まれていません")
	}
	if !strings.Contains(result, `"repeat":5`) {
		t.Error("Repeat値5が含まれていません")
	}
}

// TestAniListService_QueryAnime_RepeatFieldInTable はQueryAnimeメソッドのテーブル形式でのRepeatフィールドテスト
func TestAniListService_QueryAnime_RepeatFieldInTable(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryAnimeListFunc: func(req domain.QueryAnimeRequest) (*domain.AniListResponse, error) {
			return &domain.AniListResponse{
				Data: &domain.MediaListCollectionData{
					MediaListCollection: &domain.MediaListCollection{
						Lists: []domain.MediaList{
							{
								Entries: []domain.MediaListEntry{
									{
										Media: &domain.Media{
											ID: 1,
											Title: &domain.Title{
												Native: "高再視聴アニメ",
											},
											Studios: &domain.Studios{
												Nodes: []domain.Studio{
													{Name: "人気スタジオ"},
												},
											},
										},
										Score:    10,
										Status:   "COMPLETED",
										Progress: 12,
										Repeat:   3,
									},
								},
							},
						},
					},
				},
			}, nil
		},
	}

	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryAnime("testuser", nil, "table", 0, "", "")

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == "" {
		t.Error("結果が空文字列です")
	}
	// テーブル形式でRepeatフィールドが正しく表示されることを確認
	if !strings.Contains(result, "再視聴回数") {
		t.Error("テーブルヘッダーに「再視聴回数」が含まれていません")
	}
	if !strings.Contains(result, "高再視聴アニメ") {
		t.Error("アニメタイトルが含まれていません")
	}
	// Repeat値3がテーブルに含まれることを確認
	lines := strings.Split(result, "\n")
	found := false
	for _, line := range lines {
		if strings.Contains(line, "高再視聴アニメ") && strings.Contains(line, "\t3\t") {
			found = true
			break
		}
	}
	if !found {
		t.Error("テーブルにRepeat値3が正しく表示されていません")
	}
}
