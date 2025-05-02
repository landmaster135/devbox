package services

import (
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/domain/models"
	"github.com/landmaster135/devbox/internal/domain/repositories"
)

// FileService はファイル操作のユースケースを実装するサービスです
type FileService struct {
	fileRepo repositories.FileRepository
}

// NewFileService は新しいFileServiceインスタンスを作成します
func NewFileService(fileRepo repositories.FileRepository) *FileService {
	return &FileService{
		fileRepo: fileRepo,
	}
}

// RemoveMatchingLines は、ファイル内の各行の指定された範囲の文字列が一致する行を削除し、削除した行数を返します
func (s *FileService) RemoveMatchingLines(filePath string, startPos, endPos int) (int, error) {
	// ファイルの存在確認
	if !s.fileRepo.FileExists(filePath) {
		return 0, fmt.Errorf("ファイルが存在しません: %s", filePath)
	}

	// ファイルを読み込む
	content, err := s.fileRepo.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	// 重複行を削除
	removedCount, err := content.RemoveDuplicateLines(startPos, endPos)
	if err != nil {
		return 0, err
	}

	// 結果をファイルに書き込む
	if err := s.fileRepo.WriteFile(filePath, content); err != nil {
		return 0, err
	}

	return removedCount, nil
}

// CreateRequestBodyFromDir は指定されたディレクトリからJSONファイルを読み込み、APIリクエスト用のリクエストボディを作成します
// dirPath: JSONファイルが格納されているディレクトリのパス
// keyName: JSONデータの配列が入るキーの名前
// outputPath: 作成したリクエストボディを保存するファイルのパス（空文字列の場合は保存しない）
func (s *FileService) CreateRequestBodyFromDir(dirPath, keyName, outputPath string) ([]byte, error) {
	// ディレクトリ内のJSONファイルを検索
	jsonFiles, err := s.fileRepo.FindFilesByExt(dirPath, ".json")
	if err != nil {
		return nil, err
	}

	// JSONファイルが存在するか確認
	if len(jsonFiles) == 0 {
		return nil, fmt.Errorf("指定されたディレクトリにJSONファイルが存在しません: %s", dirPath)
	}

	// 各JSONファイルからデータを読み込む
	var jsonDataArray []interface{}
	for _, jsonFile := range jsonFiles {
		// JSONファイルを読み込む
		jsonData, err := s.fileRepo.ReadJSONFile(jsonFile)
		if err != nil {
			return nil, fmt.Errorf("JSONファイルの読み込みに失敗しました (%s): %w", jsonFile, err)
		}

		// 配列に追加
		jsonDataArray = append(jsonDataArray, jsonData)
	}

	// リクエストボディを作成
	requestBody := models.NewRequestBody(keyName, jsonDataArray)

	// JSONに変換
	requestBodyJSON, err := requestBody.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("リクエストボディのJSON変換に失敗しました: %w", err)
	}

	// 出力ファイルパスが指定されている場合は保存
	if outputPath != "" {
		// 出力ディレクトリが存在しない場合は作成
		outputDir := s.fileRepo.GetDirectoryPath(outputPath)
		if err := s.fileRepo.CreateDirectory(outputDir); err != nil {
			return nil, fmt.Errorf("出力ディレクトリの作成に失敗しました: %w", err)
		}

		// JSONをFileContentに変換
		lines := strings.Split(string(requestBodyJSON), "\n")
		content := models.NewFileContent(lines)

		// ファイルに保存
		if err := s.fileRepo.WriteFile(outputPath, content); err != nil {
			return nil, fmt.Errorf("リクエストボディの保存に失敗しました: %w", err)
		}
	}

	return requestBodyJSON, nil
}
