package usecases

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	domain "github.com/landmaster135/devbox/internal/anilist/domain"
	infrastructure "github.com/landmaster135/devbox/internal/anilist/infrastructure"
)

// TestAniListService_generateFileName_Normal はgenerateFileNameメソッドの正常系テスト
func TestAniListService_generateFileName_Normal(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{
		JoinFunc: func(elem ...string) string {
			return elem[0] + "/" + elem[1]
		},
	}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	// Act
	result := service.generateFileName("/output", "testuser", nil, "json")

	// Assert
	if result == "" {
		t.Error("ファイル名が生成されませんでした")
	}
	if result[:8] != "/output/" {
		t.Errorf("期待されるパス: /output/..., 実際: %s", result)
	}
}

// TestAniListService_saveToFile_Normal はsaveToFileメソッドの正常系テスト
func TestAniListService_saveToFile_Normal(t *testing.T) {
	// Arrange
	var mkdirAllCalled bool
	var writeFileCalled bool
	var createdPath string
	var writtenContent []byte

	mockFS := &infrastructure.MockFileSystem{
		MkdirAllFunc: func(path string, perm os.FileMode) error {
			mkdirAllCalled = true
			createdPath = path
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
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	// Act
	err := service.saveToFile("test content", "/output", "testuser", nil, "json")

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
	if createdPath != "/output" {
		t.Errorf("期待されるパス: /output, 実際: %s", createdPath)
	}
	if string(writtenContent) != "test content" {
		t.Errorf("期待される内容: test content, 実際: %s", string(writtenContent))
	}
}

// TestAniListService_saveToFile_MkdirAllError はsaveToFileメソッドのMkdirAllエラーテスト
func TestAniListService_saveToFile_MkdirAllError(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{
		MkdirAllFunc: func(path string, perm os.FileMode) error {
			return errors.New("ディレクトリ作成エラー")
		},
	}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	// Act
	err := service.saveToFile("test content", "/output", "testuser", nil, "json")

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if err.Error() != "出力ディレクトリの作成に失敗しました: ディレクトリ作成エラー" {
		t.Errorf("期待されるエラーメッセージと異なります: %v", err)
	}
}

// TestAniListService_saveToFile_WriteFileError はsaveToFileメソッドのWriteFileエラーテスト
func TestAniListService_saveToFile_WriteFileError(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{
		MkdirAllFunc: func(path string, perm os.FileMode) error {
			return nil
		},
		WriteFileFunc: func(filename string, data []byte, perm os.FileMode) error {
			return errors.New("ファイル書き込みエラー")
		},
		JoinFunc: func(elem ...string) string {
			return elem[0] + "/" + elem[1]
		},
	}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	// Act
	err := service.saveToFile("test content", "/output", "testuser", nil, "json")

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if err.Error() != "ファイルの書き込みに失敗しました: ファイル書き込みエラー" {
		t.Errorf("期待されるエラーメッセージと異なります: %v", err)
	}
}

// TestAniListService_formatAsJSON_Normal はformatAsJSONメソッドの正常系テスト
func TestAniListService_formatAsJSON_Normal(t *testing.T) {
	// Arrange
	var marshalIndentCalled bool
	var marshalIndentInput any

	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{
		MarshalIndentFunc: func(v any, prefix, indent string) ([]byte, error) {
			marshalIndentCalled = true
			marshalIndentInput = v
			return []byte(`[{"id":1,"title":"test"}]`), nil
		},
	}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	animeList := []domain.AnimeInfo{
		{ID: 1, Title: "test"},
	}

	// Act
	result, err := service.formatAsJSON(animeList)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if !marshalIndentCalled {
		t.Error("MarshalIndentが呼び出されませんでした")
	}
	if result != `[{"id":1,"title":"test"}]` {
		t.Errorf("期待される結果: [{'id':1,'title':'test'}], 実際: %s", result)
	}
	if marshalIndentInput == nil {
		t.Error("MarshalIndentに入力が渡されませんでした")
	}
}

// TestAniListService_formatAsJSON_Error はformatAsJSONメソッドのエラーテスト
func TestAniListService_formatAsJSON_Error(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{
		MarshalIndentFunc: func(v any, prefix, indent string) ([]byte, error) {
			return nil, errors.New("JSONエンコードエラー")
		},
	}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	animeList := []domain.AnimeInfo{
		{ID: 1, Title: "test"},
	}

	// Act
	result, err := service.formatAsJSON(animeList)

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if err.Error() != "JSONエンコードに失敗しました: JSONエンコードエラー" {
		t.Errorf("期待されるエラーメッセージと異なります: %v", err)
	}
	if result != "" {
		t.Errorf("エラー時は空文字列が期待されます, 実際: %s", result)
	}
}

// TestAniListService_formatCompletedAt_Normal はformatCompletedAtメソッドの正常系テスト
func TestAniListService_formatCompletedAt_Normal(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	year := 2023
	month := 12
	day := 25
	fuzzyDate := &domain.FuzzyDate{
		Year:  &year,
		Month: &month,
		Day:   &day,
	}

	// Act
	result := service.formatCompletedAt(fuzzyDate)

	// Assert
	if result == nil {
		t.Error("結果がnilです")
		return
	}
	if result.Year() != 2023 {
		t.Errorf("期待される年: 2023, 実際: %d", result.Year())
	}
	if result.Month() != 12 {
		t.Errorf("期待される月: 12, 実際: %d", result.Month())
	}
	if result.Day() != 25 {
		t.Errorf("期待される日: 25, 実際: %d", result.Day())
	}
}

// TestAniListService_formatCompletedAt_NilDate はformatCompletedAtメソッドのnilテスト
func TestAniListService_formatCompletedAt_NilDate(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	// Act
	result := service.formatCompletedAt(nil)

	// Assert
	if result != nil {
		t.Errorf("nilが期待されましたが、%v が返されました", result)
	}
}

// TestAniListService_formatCompletedAt_NilYear はformatCompletedAtメソッドの年がnilのテスト
func TestAniListService_formatCompletedAt_NilYear(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	month := 12
	day := 25
	fuzzyDate := &domain.FuzzyDate{
		Year:  nil,
		Month: &month,
		Day:   &day,
	}

	// Act
	result := service.formatCompletedAt(fuzzyDate)

	// Assert
	if result != nil {
		t.Errorf("nilが期待されましたが、%v が返されました", result)
	}
}

// TestAniListService_formatCompletedAt_NilMonth はformatCompletedAtメソッドの月がnilのテスト
func TestAniListService_formatCompletedAt_NilMonth(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	year := 2023
	day := 25
	fuzzyDate := &domain.FuzzyDate{
		Year:  &year,
		Month: nil,
		Day:   &day,
	}

	// Act
	result := service.formatCompletedAt(fuzzyDate)

	// Assert
	if result != nil {
		t.Errorf("nilが期待されましたが、%v が返されました", result)
	}
}

// TestAniListService_formatCompletedAt_NilDay はformatCompletedAtメソッドの日がnilのテスト
func TestAniListService_formatCompletedAt_NilDay(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	year := 2023
	month := 12
	fuzzyDate := &domain.FuzzyDate{
		Year:  &year,
		Month: &month,
		Day:   nil,
	}

	// Act
	result := service.formatCompletedAt(fuzzyDate)

	// Assert
	if result != nil {
		t.Errorf("nilが期待されましたが、%v が返されました", result)
	}
}

// TestAniListService_generateFileName_WithUserID はgenerateFileNameメソッドのユーザーID指定テスト
func TestAniListService_generateFileName_WithUserID(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{
		JoinFunc: func(elem ...string) string {
			return elem[0] + "/" + elem[1]
		},
	}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	userID := 12345

	// Act
	result := service.generateFileName("/output", "", &userID, "json")

	// Assert
	if result == "" {
		t.Error("ファイル名が生成されませんでした")
	}
	if result[:8] != "/output/" {
		t.Errorf("期待されるパス: /output/..., 実際: %s", result)
	}
	// ユーザーIDが含まれていることを確認
	if !strings.Contains(result, "user_12345") {
		t.Errorf("ユーザーIDが含まれていません: %s", result)
	}
}

// TestAniListService_generateFileName_UnknownUser はgenerateFileNameメソッドの不明ユーザーテスト
func TestAniListService_generateFileName_UnknownUser(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{
		JoinFunc: func(elem ...string) string {
			return elem[0] + "/" + elem[1]
		},
	}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	// Act
	result := service.generateFileName("/output", "", nil, "json")

	// Assert
	if result == "" {
		t.Error("ファイル名が生成されませんでした")
	}
	if !strings.Contains(result, "unknown") {
		t.Errorf("unknownが含まれていません: %s", result)
	}
}

// TestAniListService_generateFileName_TableFormat はgenerateFileNameメソッドのテーブル形式テスト
func TestAniListService_generateFileName_TableFormat(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{
		JoinFunc: func(elem ...string) string {
			return elem[0] + "/" + elem[1]
		},
	}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	// Act
	result := service.generateFileName("/output", "testuser", nil, "table")

	// Assert
	if result == "" {
		t.Error("ファイル名が生成されませんでした")
	}
	if !strings.HasSuffix(result, ".txt") {
		t.Errorf("期待される拡張子: .txt, 実際: %s", result)
	}
}

// TestAniListService_getJST_Normal はgetJSTメソッドの正常系テスト
func TestAniListService_getJST_Normal(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	// Act
	result := service.getJST()

	// Assert
	if result == nil {
		t.Error("JSTタイムゾーンがnilです")
		return
	}

	// JSTは+9時間のオフセット
	_, offset := time.Now().In(result).Zone()
	expectedOffset := 9 * 60 * 60 // 9時間を秒で表現
	if offset != expectedOffset {
		t.Errorf("期待されるオフセット: %d秒, 実際: %d秒", expectedOffset, offset)
	}
}

// TestAniListService_formatUpdatedAt_Normal はformatUpdatedAtメソッドの正常系テスト
func TestAniListService_formatUpdatedAt_Normal(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	// 2023年12月25日 12:00:00 UTCのUnixタイムスタンプ
	timestamp := int64(1703505600)

	// Act
	result := service.formatUpdatedAt(timestamp)

	// Assert
	// JSTで変換されていることを確認
	expectedTime := time.Unix(timestamp, 0).In(service.getJST())
	if !result.Equal(expectedTime) {
		t.Errorf("期待される時刻: %v, 実際: %v", expectedTime, result)
	}

	// タイムゾーンがJSTであることを確認
	zone, _ := result.Zone()
	if zone != "JST" {
		t.Errorf("期待されるタイムゾーン: JST, 実際: %s", zone)
	}
}

// TestAniListService_filterByStatus_Normal はfilterByStatusメソッドの正常系テスト
func TestAniListService_filterByStatus_Normal(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	animeList := []domain.AnimeInfo{
		{ID: 1, Title: "Anime1", Status: "COMPLETED"},
		{ID: 2, Title: "Anime2", Status: "CURRENT"},
		{ID: 3, Title: "Anime3", Status: "COMPLETED"},
		{ID: 4, Title: "Anime4", Status: "DROPPED"},
	}

	// Act
	result := service.filterByStatus(animeList, "COMPLETED")

	// Assert
	if len(result) != 2 {
		t.Errorf("期待される件数: 2, 実際: %d", len(result))
	}

	for _, anime := range result {
		if anime.Status != "COMPLETED" {
			t.Errorf("フィルタされていないアニメが含まれています: %s", anime.Status)
		}
	}

	// IDの確認
	if result[0].ID != 1 || result[1].ID != 3 {
		t.Errorf("期待されるID: [1, 3], 実際: [%d, %d]", result[0].ID, result[1].ID)
	}
}

// TestAniListService_filterByStatus_NoMatch はfilterByStatusメソッドのマッチなしテスト
func TestAniListService_filterByStatus_NoMatch(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	animeList := []domain.AnimeInfo{
		{ID: 1, Title: "Anime1", Status: "COMPLETED"},
		{ID: 2, Title: "Anime2", Status: "CURRENT"},
	}

	// Act
	result := service.filterByStatus(animeList, "PLANNING")

	// Assert
	if len(result) != 0 {
		t.Errorf("期待される件数: 0, 実際: %d", len(result))
	}
}

// TestAniListService_formatAsTable_Normal はformatAsTableメソッドの正常系テスト
func TestAniListService_formatAsTable_Normal(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	completedAt := time.Date(2023, 12, 25, 0, 0, 0, 0, service.getJST())
	animeList := []domain.AnimeInfo{
		{
			ID:          1,
			Title:       "Test Anime",
			Status:      "COMPLETED",
			Score:       9,
			Progress:    12,
			CompletedAt: &completedAt,
			Studio:      "Test Studio",
		},
	}

	// Act
	result, err := service.formatAsTable(animeList)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}

	if result == "" {
		t.Error("結果が空文字列です")
	}

	// ヘッダーが含まれていることを確認
	if !strings.Contains(result, "ID\tタイトル\tステータス") {
		t.Error("ヘッダーが含まれていません")
	}

	// データが含まれていることを確認
	if !strings.Contains(result, "Test Anime") {
		t.Error("アニメタイトルが含まれていません")
	}

	if !strings.Contains(result, "2023-12-25") {
		t.Error("完了日が含まれていません")
	}
}

// TestAniListService_formatAsTable_Empty はformatAsTableメソッドの空リストテスト
func TestAniListService_formatAsTable_Empty(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	var animeList []domain.AnimeInfo

	// Act
	result, err := service.formatAsTable(animeList)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}

	expected := "アニメが見つかりませんでした。\n"
	if result != expected {
		t.Errorf("期待される結果: %s, 実際: %s", expected, result)
	}
}

// TestAniListService_formatAsTable_NilCompletedAt はformatAsTableメソッドの完了日nilテスト
func TestAniListService_formatAsTable_NilCompletedAt(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	animeList := []domain.AnimeInfo{
		{
			ID:          1,
			Title:       "Test Anime",
			Status:      "CURRENT",
			Score:       8,
			Progress:    5,
			CompletedAt: nil,
			Studio:      "Test Studio",
		},
	}

	// Act
	result, err := service.formatAsTable(animeList)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}

	if result == "" {
		t.Error("結果が空文字列です")
	}

	// 完了日が空文字列として表示されることを確認
	lines := strings.Split(result, "\n")
	if len(lines) < 3 {
		t.Error("期待される行数が不足しています")
		return
	}

	dataLine := lines[2] // ヘッダー、区切り線の次がデータ行
	fields := strings.Split(dataLine, "\t")
	if len(fields) < 6 {
		t.Error("期待されるフィールド数が不足しています")
		return
	}

	// 完了日フィールド（インデックス5）が空であることを確認
	if fields[5] != "" {
		t.Errorf("完了日は空文字列が期待されます, 実際: %s", fields[5])
	}
}

// TestAniListService_transformToAnimeInfoList_Normal はtransformToAnimeInfoListメソッドの正常系テスト
func TestAniListService_transformToAnimeInfoList_Normal(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	year := 2023
	month := 12
	day := 25
	collection := &domain.MediaListCollection{
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
						CompletedAt: &domain.FuzzyDate{
							Year:  &year,
							Month: &month,
							Day:   &day,
						},
						Notes:     "テストノート",
						UpdatedAt: 1703505600, // 2023-12-25 12:00:00 UTC
					},
				},
			},
		},
	}

	// Act
	result := service.transformToAnimeInfoList(collection)

	// Assert
	if len(result) != 1 {
		t.Errorf("期待される件数: 1, 実際: %d", len(result))
		return
	}

	anime := result[0]
	if anime.ID != 1 {
		t.Errorf("期待されるID: 1, 実際: %d", anime.ID)
	}
	if anime.Title != "テストアニメ" {
		t.Errorf("期待されるタイトル: テストアニメ, 実際: %s", anime.Title)
	}
	if anime.Score != 9 {
		t.Errorf("期待されるスコア: 9, 実際: %d", anime.Score)
	}
	if anime.Status != "COMPLETED" {
		t.Errorf("期待されるステータス: COMPLETED, 実際: %s", anime.Status)
	}
	if anime.Progress != 12 {
		t.Errorf("期待される進行状況: 12, 実際: %d", anime.Progress)
	}
	if anime.Notes != "テストノート" {
		t.Errorf("期待されるノート: テストノート, 実際: %s", anime.Notes)
	}
	if anime.CoverImageURL != "https://example.com/image.jpg" {
		t.Errorf("期待されるカバー画像URL: https://example.com/image.jpg, 実際: %s", anime.CoverImageURL)
	}
	if anime.SiteURL != "https://anilist.co/anime/1" {
		t.Errorf("期待されるサイトURL: https://anilist.co/anime/1, 実際: %s", anime.SiteURL)
	}
	if anime.Studio != "テストスタジオ" {
		t.Errorf("期待されるスタジオ: テストスタジオ, 実際: %s", anime.Studio)
	}
	if anime.CompletedAt == nil {
		t.Error("完了日がnilです")
	} else {
		if anime.CompletedAt.Year() != 2023 {
			t.Errorf("期待される完了年: 2023, 実際: %d", anime.CompletedAt.Year())
		}
	}
}

// TestAniListService_transformToAnimeInfoList_NilMedia はtransformToAnimeInfoListメソッドのMediaがnilのテスト
func TestAniListService_transformToAnimeInfoList_NilMedia(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	collection := &domain.MediaListCollection{
		Lists: []domain.MediaList{
			{
				Entries: []domain.MediaListEntry{
					{
						Media:    nil, // Mediaがnil
						Score:    9,
						Status:   "COMPLETED",
						Progress: 12,
					},
					{
						Media: &domain.Media{
							ID: 2,
							Title: &domain.Title{
								Native: "有効なアニメ",
							},
						},
						Score:    8,
						Status:   "CURRENT",
						Progress: 5,
					},
				},
			},
		},
	}

	// Act
	result := service.transformToAnimeInfoList(collection)

	// Assert
	// Mediaがnilのエントリはスキップされ、有効なエントリのみが返される
	if len(result) != 1 {
		t.Errorf("期待される件数: 1, 実際: %d", len(result))
		return
	}

	anime := result[0]
	if anime.ID != 2 {
		t.Errorf("期待されるID: 2, 実際: %d", anime.ID)
	}
	if anime.Title != "有効なアニメ" {
		t.Errorf("期待されるタイトル: 有効なアニメ, 実際: %s", anime.Title)
	}
}

// TestAniListService_transformToAnimeInfoList_NilFields はtransformToAnimeInfoListメソッドの各種フィールドがnilのテスト
func TestAniListService_transformToAnimeInfoList_NilFields(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	collection := &domain.MediaListCollection{
		Lists: []domain.MediaList{
			{
				Entries: []domain.MediaListEntry{
					{
						Media: &domain.Media{
							ID:         1,
							Title:      nil, // Titleがnil
							CoverImage: nil, // CoverImageがnil
							SiteURL:    "",
							Studios:    nil, // Studiosがnil
						},
						Score:       7,
						Status:      "CURRENT",
						Progress:    3,
						CompletedAt: nil, // CompletedAtがnil
						Notes:       "",
						UpdatedAt:   1703505600,
					},
				},
			},
		},
	}

	// Act
	result := service.transformToAnimeInfoList(collection)

	// Assert
	if len(result) != 1 {
		t.Errorf("期待される件数: 1, 実際: %d", len(result))
		return
	}

	anime := result[0]
	if anime.ID != 1 {
		t.Errorf("期待されるID: 1, 実際: %d", anime.ID)
	}
	if anime.Title != "" {
		t.Errorf("期待されるタイトル: 空文字列, 実際: %s", anime.Title)
	}
	if anime.CoverImageURL != "" {
		t.Errorf("期待されるカバー画像URL: 空文字列, 実際: %s", anime.CoverImageURL)
	}
	if anime.Studio != "" {
		t.Errorf("期待されるスタジオ: 空文字列, 実際: %s", anime.Studio)
	}
	if anime.CompletedAt != nil {
		t.Errorf("期待される完了日: nil, 実際: %v", anime.CompletedAt)
	}
}

// TestAniListService_transformToAnimeInfoList_EmptyStudios はtransformToAnimeInfoListメソッドのスタジオが空のテスト
func TestAniListService_transformToAnimeInfoList_EmptyStudios(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	collection := &domain.MediaListCollection{
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
								Nodes: []domain.Studio{}, // 空のスタジオリスト
							},
						},
						Score:    8,
						Status:   "CURRENT",
						Progress: 5,
					},
				},
			},
		},
	}

	// Act
	result := service.transformToAnimeInfoList(collection)

	// Assert
	if len(result) != 1 {
		t.Errorf("期待される件数: 1, 実際: %d", len(result))
		return
	}

	anime := result[0]
	if anime.Studio != "" {
		t.Errorf("期待されるスタジオ: 空文字列, 実際: %s", anime.Studio)
	}
}

// TestAniListService_transformToAnimeInfoList_MultipleEntries はtransformToAnimeInfoListメソッドの複数エントリテスト
func TestAniListService_transformToAnimeInfoList_MultipleEntries(t *testing.T) {
	// Arrange
	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{}
	service := NewAniListServiceWithDependencies(nil, mockFS, mockJSON)

	collection := &domain.MediaListCollection{
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
					},
				},
			},
			{
				Entries: []domain.MediaListEntry{
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
					},
				},
			},
		},
	}

	// Act
	result := service.transformToAnimeInfoList(collection)

	// Assert
	if len(result) != 3 {
		t.Errorf("期待される件数: 3, 実際: %d", len(result))
		return
	}

	// 各エントリの確認
	expectedTitles := []string{"アニメ1", "アニメ2", "アニメ3"}
	expectedIDs := []int{1, 2, 3}

	for i, anime := range result {
		if anime.ID != expectedIDs[i] {
			t.Errorf("エントリ%d - 期待されるID: %d, 実際: %d", i, expectedIDs[i], anime.ID)
		}
		if anime.Title != expectedTitles[i] {
			t.Errorf("エントリ%d - 期待されるタイトル: %s, 実際: %s", i, expectedTitles[i], anime.Title)
		}
	}
}
