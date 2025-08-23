package usecases

import (
	"errors"
	"os"
	"testing"

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
	var marshalIndentInput interface{}

	mockFS := &infrastructure.MockFileSystem{}
	mockJSON := &infrastructure.MockJSONProcessor{
		MarshalIndentFunc: func(v interface{}, prefix, indent string) ([]byte, error) {
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
		MarshalIndentFunc: func(v interface{}, prefix, indent string) ([]byte, error) {
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
