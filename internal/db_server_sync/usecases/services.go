package usecases

import (
	"encoding/json"
	"fmt"
	"os"

	models "github.com/landmaster135/devbox/internal/db_server_sync/models"
)

// DbServerSyncService はデータベースサーバー同期サービスを表す構造体
type DbServerSyncService struct{}

// NewDbServerSyncService は新しいDbServerSyncServiceを作成する
func NewDbServerSyncService() *DbServerSyncService {
	return &DbServerSyncService{}
}

// ProcessAppendAnime はappend-anime操作を実行する
func (s *DbServerSyncService) ProcessAppendAnime(inputFilePath, additionalInputFilePath, outputFilePath string) error {
	// メインファイルを読み込み
	anilistData, err := s.loadAniListData(inputFilePath)
	if err != nil {
		return fmt.Errorf("メインファイルの読み込みに失敗しました: %w", err)
	}

	// 追加ファイルを読み込み（任意）
	var additionalData []models.AdditionalAnimeData
	if additionalInputFilePath != "" {
		additionalData, err = s.loadAdditionalData(additionalInputFilePath)
		if err != nil {
			return fmt.Errorf("追加ファイルの読み込みに失敗しました: %w", err)
		}
	}

	// データをマージして変換
	outputData, err := s.mergeAndConvertData(anilistData, additionalData)
	if err != nil {
		return fmt.Errorf("データの変換に失敗しました: %w", err)
	}

	// 出力ファイルに書き込み
	if err := s.writeOutputData(outputData, outputFilePath); err != nil {
		return fmt.Errorf("出力ファイルの書き込みに失敗しました: %w", err)
	}

	return nil
}

// loadAniListData はAniListデータファイルを読み込む
func (s *DbServerSyncService) loadAniListData(filePath string) ([]models.AniListAnime, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var anilistData []models.AniListAnime
	if err := json.Unmarshal(data, &anilistData); err != nil {
		return nil, err
	}

	return anilistData, nil
}

// loadAdditionalData は追加データファイルを読み込む
func (s *DbServerSyncService) loadAdditionalData(filePath string) ([]models.AdditionalAnimeData, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var additionalData []models.AdditionalAnimeData
	if err := json.Unmarshal(data, &additionalData); err != nil {
		return nil, err
	}

	return additionalData, nil
}

// mergeAndConvertData はデータをマージして出力形式に変換する
func (s *DbServerSyncService) mergeAndConvertData(anilistData []models.AniListAnime, additionalData []models.AdditionalAnimeData) (*models.OutputData, error) {
	// 追加データをマップに変換（高速検索のため）
	additionalMap := make(map[int]*models.AdditionalAnimeData)
	for i := range additionalData {
		additionalMap[additionalData[i].AniListID] = &additionalData[i]
	}

	// 変換処理
	var outputAnimes []models.OutputAnime
	for _, anime := range anilistData {
		// 追加データを検索
		additional := additionalMap[anime.ID]

		// 変換
		outputAnime, err := anime.ConvertToOutputAnime(additional)
		if err != nil {
			return nil, fmt.Errorf("アニメID %d の変換に失敗しました: %w", anime.ID, err)
		}

		outputAnimes = append(outputAnimes, *outputAnime)
	}

	// 最終的な出力データを作成
	outputData := &models.OutputData{
		Data: models.OutputDataContent{
			Animes: outputAnimes,
		},
		Description: "Anime data from AniList",
		Name:        "My Anime List",
	}

	return outputData, nil
}

// writeOutputData は出力データをファイルに書き込む
func (s *DbServerSyncService) writeOutputData(outputData *models.OutputData, filePath string) error {
	data, err := json.MarshalIndent(outputData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}
