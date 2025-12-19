package usecases

import (
	"errors"
	"os"
	"strings"
	"testing"

	domain "github.com/landmaster135/devbox/internal/anilist/domain"
	infrastructure "github.com/landmaster135/devbox/internal/anilist/infrastructure"
)

// TestAniListService_QueryManga_Normal はQueryMangaメソッドの正常系テスト
func TestAniListService_QueryManga_Normal(t *testing.T) {
	// Arrange
	var queryMangaListCalled bool
	var queryMangaListRequest domain.QueryMangaRequest

	mockRepository := &domain.MockAniListRepository{
		QueryMangaListFunc: func(req domain.QueryMangaRequest) (*domain.AniListResponse, error) {
			queryMangaListCalled = true
			queryMangaListRequest = req

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
												Native: "テストマンガ",
											},
											CoverImage: &domain.Image{
												ExtraLarge: "https://example.com/manga.jpg",
											},
											SiteURL: "https://anilist.co/manga/1",
										},
										Score:           9,
										Status:          "COMPLETED",
										Progress:        120,
										ProgressVolumes: 12,
										Repeat:          2,
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
			return []byte(`[{"id":1,"title":"テストマンガ"}]`), nil
		},
	}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryManga("testuser", nil, "json", 0, "", "")

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if !queryMangaListCalled {
		t.Error("QueryMangaListが呼び出されませんでした")
	}
	if queryMangaListRequest.Username != "testuser" {
		t.Errorf("期待されるユーザー名: testuser, 実際: %s", queryMangaListRequest.Username)
	}
	if queryMangaListRequest.UserID != nil {
		t.Errorf("期待されるユーザーID: nil, 実際: %v", queryMangaListRequest.UserID)
	}
	if result == "" {
		t.Error("結果が空文字列です")
	}
	if !strings.Contains(result, "テストマンガ") {
		t.Error("結果にマンガタイトルが含まれていません")
	}
}

// TestAniListService_QueryManga_WithUserID はQueryMangaメソッドのユーザーID指定テスト
func TestAniListService_QueryManga_WithUserID(t *testing.T) {
	// Arrange
	var queryMangaListRequest domain.QueryMangaRequest
	userID := 12345

	mockRepository := &domain.MockAniListRepository{
		QueryMangaListFunc: func(req domain.QueryMangaRequest) (*domain.AniListResponse, error) {
			queryMangaListRequest = req
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
												Native: "テストマンガ",
											},
										},
										Score:           8,
										Status:          "CURRENT",
										Progress:        50,
										ProgressVolumes: 5,
										Repeat:          1,
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
			return []byte(`[{"id":1,"title":"テストマンガ"}]`), nil
		},
	}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryManga("", &userID, "json", 0, "", "")

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if queryMangaListRequest.Username != "" {
		t.Errorf("期待されるユーザー名: 空文字列, 実際: %s", queryMangaListRequest.Username)
	}
	if queryMangaListRequest.UserID == nil {
		t.Error("ユーザーIDがnilです")
	} else if *queryMangaListRequest.UserID != 12345 {
		t.Errorf("期待されるユーザーID: 12345, 実際: %d", *queryMangaListRequest.UserID)
	}
	if result == "" {
		t.Error("結果が空文字列です")
	}
}

// TestAniListService_QueryManga_TableFormat はQueryMangaメソッドのテーブル形式テスト
func TestAniListService_QueryManga_TableFormat(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryMangaListFunc: func(req domain.QueryMangaRequest) (*domain.AniListResponse, error) {
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
												Native: "テストマンガ",
											},
										},
										Score:           9,
										Status:          "COMPLETED",
										Progress:        120,
										ProgressVolumes: 12,
										Repeat:          2,
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
	result, err := service.QueryManga("testuser", nil, "table", 0, "", "")

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == "" {
		t.Error("結果が空文字列です")
	}
	// テーブル形式の確認
	if !strings.Contains(result, "ID\tタイトル\tステータス\tスコア\t進行状況\t巻数進行\t再読回数") {
		t.Error("テーブルヘッダーが含まれていません")
	}
	if !strings.Contains(result, "テストマンガ") {
		t.Error("マンガタイトルが含まれていません")
	}
	if !strings.Contains(result, "\t2\t") {
		t.Error("再読回数が含まれていません")
	}
}

// TestAniListService_QueryManga_DefaultFormat はQueryMangaメソッドのデフォルト形式テスト
func TestAniListService_QueryManga_DefaultFormat(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryMangaListFunc: func(req domain.QueryMangaRequest) (*domain.AniListResponse, error) {
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
												Native: "テストマンガ",
											},
										},
										Score:           9,
										Status:          "COMPLETED",
										Progress:        120,
										ProgressVolumes: 12,
										Repeat:          1,
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
			return []byte(`[{"id":1,"title":"テストマンガ"}]`), nil
		},
	}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act - 無効な形式を指定してデフォルト（JSON）が使用されることを確認
	result, err := service.QueryManga("testuser", nil, "invalid-format", 0, "", "")

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

// TestAniListService_QueryManga_WithOutputDir はQueryMangaメソッドのファイル出力テスト
func TestAniListService_QueryManga_WithOutputDir(t *testing.T) {
	// Arrange
	var mkdirAllCalled bool
	var writeFileCalled bool
	var writtenContent []byte

	mockRepository := &domain.MockAniListRepository{
		QueryMangaListFunc: func(req domain.QueryMangaRequest) (*domain.AniListResponse, error) {
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
												Native: "テストマンガ",
											},
										},
										Score:           9,
										Status:          "COMPLETED",
										Progress:        120,
										ProgressVolumes: 12,
										Repeat:          1,
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
			return []byte(`[{"id":1,"title":"テストマンガ"}]`), nil
		},
	}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryManga("testuser", nil, "json", 0, "", "/output")

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
	if string(writtenContent) != `[{"id":1,"title":"テストマンガ"}]` {
		t.Errorf("期待される書き込み内容と異なります: %s", string(writtenContent))
	}
	if !strings.Contains(result, "結果をファイルに保存しました") {
		t.Error("ファイル保存メッセージが含まれていません")
	}
}

// TestAniListService_QueryManga_FileSaveError はQueryMangaメソッドのファイル保存エラーテスト
func TestAniListService_QueryManga_FileSaveError(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryMangaListFunc: func(req domain.QueryMangaRequest) (*domain.AniListResponse, error) {
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
												Native: "テストマンガ",
											},
										},
										Score:           9,
										Status:          "COMPLETED",
										Progress:        120,
										ProgressVolumes: 12,
										Repeat:          0,
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
			return []byte(`[{"id":1,"title":"テストマンガ"}]`), nil
		},
	}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryManga("testuser", nil, "json", 0, "", "/output")

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

// TestAniListService_QueryManga_JSONFormatError はQueryMangaメソッドのJSON形式エラーテスト
func TestAniListService_QueryManga_JSONFormatError(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryMangaListFunc: func(req domain.QueryMangaRequest) (*domain.AniListResponse, error) {
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
												Native: "テストマンガ",
											},
										},
										Score:           9,
										Status:          "COMPLETED",
										Progress:        120,
										ProgressVolumes: 12,
										Repeat:          0,
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
	result, err := service.QueryManga("testuser", nil, "json", 0, "", "")

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

// TestAniListService_QueryManga_APIError はQueryMangaメソッドのAPIエラーテスト
func TestAniListService_QueryManga_APIError(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryMangaListFunc: func(req domain.QueryMangaRequest) (*domain.AniListResponse, error) {
			return nil, errors.New("API呼び出しエラー")
		},
	}

	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryManga("testuser", nil, "json", 0, "", "")

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

// TestAniListService_QueryManga_WithStatusFilter はQueryMangaメソッドのステータスフィルタテスト
func TestAniListService_QueryManga_WithStatusFilter(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryMangaListFunc: func(req domain.QueryMangaRequest) (*domain.AniListResponse, error) {
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
												Native: "完了マンガ",
											},
										},
										Score:           9,
										Status:          "COMPLETED",
										Progress:        100,
										ProgressVolumes: 10,
										Repeat:          3,
									},
									{
										Media: &domain.Media{
											ID: 2,
											Title: &domain.Title{
												Native: "読書中マンガ",
											},
										},
										Score:           8,
										Status:          "CURRENT",
										Progress:        50,
										ProgressVolumes: 5,
										Repeat:          0,
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
			mangaList := v.([]domain.MangaInfo)
			if len(mangaList) != 1 {
				t.Errorf("フィルタ後の期待される件数: 1, 実際: %d", len(mangaList))
			}
			if mangaList[0].Status != "COMPLETED" {
				t.Errorf("フィルタされたマンガのステータス: COMPLETED, 実際: %s", mangaList[0].Status)
			}
			return []byte(`[{"id":1,"title":"完了マンガ","status":"COMPLETED"}]`), nil
		},
	}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryManga("testuser", nil, "json", 0, "COMPLETED", "")

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == "" {
		t.Error("結果が空文字列です")
	}
	if !strings.Contains(result, "完了マンガ") {
		t.Error("フィルタされたマンガが含まれていません")
	}
}

// TestAniListService_QueryManga_WithLimit はQueryMangaメソッドの制限値テスト
func TestAniListService_QueryManga_WithLimit(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryMangaListFunc: func(req domain.QueryMangaRequest) (*domain.AniListResponse, error) {
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
												Native: "マンガ1",
											},
										},
										Score:           9,
										Status:          "COMPLETED",
										Progress:        120,
										ProgressVolumes: 12,
										Repeat:          1,
									},
									{
										Media: &domain.Media{
											ID: 2,
											Title: &domain.Title{
												Native: "マンガ2",
											},
										},
										Score:           8,
										Status:          "CURRENT",
										Progress:        50,
										ProgressVolumes: 5,
										Repeat:          0,
									},
									{
										Media: &domain.Media{
											ID: 3,
											Title: &domain.Title{
												Native: "マンガ3",
											},
										},
										Score:           7,
										Status:          "PLANNING",
										Progress:        0,
										ProgressVolumes: 0,
										Repeat:          0,
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
			mangaList := v.([]domain.MangaInfo)
			if len(mangaList) != 2 {
				t.Errorf("制限後の期待される件数: 2, 実際: %d", len(mangaList))
			}
			return []byte(`[{"id":1,"title":"マンガ1"},{"id":2,"title":"マンガ2"}]`), nil
		},
	}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryManga("testuser", nil, "json", 2, "", "")

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == "" {
		t.Error("結果が空文字列です")
	}
}

// TestAniListService_QueryManga_NoData はQueryMangaメソッドのデータなしエラーテスト
func TestAniListService_QueryManga_NoData(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryMangaListFunc: func(req domain.QueryMangaRequest) (*domain.AniListResponse, error) {
			return &domain.AniListResponse{
				Data: nil, // データがnil
			}, nil
		},
	}

	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}

	service := NewAniListServiceWithDependencies(mockRepository, mockFS, mockJSON)

	// Act
	result, err := service.QueryManga("testuser", nil, "json", 0, "", "")

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	expectedError := "ユーザーのマンガリストが見つかりません"
	if err.Error() != expectedError {
		t.Errorf("期待されるエラーメッセージ: %s, 実際: %s", expectedError, err.Error())
	}
	if result != "" {
		t.Errorf("エラー時は空文字列が期待されます, 実際: %s", result)
	}
}

// TestAniListService_QueryManga_NoMediaListCollection はQueryMangaメソッドのMediaListCollectionなしエラーテスト
func TestAniListService_QueryManga_NoMediaListCollection(t *testing.T) {
	// Arrange
	mockRepository := &domain.MockAniListRepository{
		QueryMangaListFunc: func(req domain.QueryMangaRequest) (*domain.AniListResponse, error) {
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
	result, err := service.QueryManga("testuser", nil, "json", 0, "", "")

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	expectedError := "ユーザーのマンガリストが見つかりません"
	if err.Error() != expectedError {
		t.Errorf("期待されるエラーメッセージ: %s, 実際: %s", expectedError, err.Error())
	}
	if result != "" {
		t.Errorf("エラー時は空文字列が期待されます, 実際: %s", result)
	}
}
