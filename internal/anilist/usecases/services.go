package usecases

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/landmaster135/devbox/internal/anilist/domain"
	"github.com/landmaster135/devbox/internal/anilist/infrastructure"
)

// AniListService はAniListサービス
type AniListService struct {
	client *infrastructure.AniListClient
}

// NewAniListService は新しいAniListServiceを作成する
func NewAniListService() *AniListService {
	return &AniListService{
		client: infrastructure.NewAniListClient(),
	}
}

// NewAniListServiceWithClient はクライアントを指定してAniListServiceを作成する
func NewAniListServiceWithClient(client *infrastructure.AniListClient) *AniListService {
	return &AniListService{
		client: client,
	}
}

// generateFileName はファイル名を生成する
func (s *AniListService) generateFileName(outputDir, username string, userID *int, format string) string {
	now := time.Now()
	timestamp := now.Format("20060102_150405")

	var userIdentifier string
	if username != "" {
		userIdentifier = username
	} else if userID != nil {
		userIdentifier = fmt.Sprintf("user_%d", *userID)
	} else {
		userIdentifier = "unknown"
	}

	var extension string
	switch format {
	case "json":
		extension = "json"
	case "table":
		extension = "txt"
	default:
		extension = "json"
	}

	fileName := fmt.Sprintf("anilist_%s_%s.%s", userIdentifier, timestamp, extension)
	return filepath.Join(outputDir, fileName)
}

// saveToFile は結果をファイルに保存する
func (s *AniListService) saveToFile(content, outputDir, username string, userID *int, format string) error {
	// 出力ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("出力ディレクトリの作成に失敗しました: %v", err)
	}

	// ファイル名を生成
	fileName := s.generateFileName(outputDir, username, userID, format)

	// ファイルに書き込み
	if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました: %v", err)
	}

	return nil
}

// QueryAnime はアニメ情報を取得する
func (s *AniListService) QueryAnime(username string, userID *int, format string, limit int, status, outputDir string) (string, error) {
	// リクエストを構築
	req := domain.QueryAnimeRequest{
		Username: username,
		UserID:   userID,
	}

	// APIを呼び出し
	resp, err := s.client.QueryAnimeList(req)
	if err != nil {
		return "", fmt.Errorf("AniList APIの呼び出しに失敗しました: %v", err)
	}

	// データが存在しない場合
	if resp.Data == nil || resp.Data.MediaListCollection == nil {
		return "", fmt.Errorf("ユーザーのアニメリストが見つかりません")
	}

	// アニメ情報を整形
	animeList := s.transformToAnimeInfoList(resp.Data.MediaListCollection)

	// ステータスフィルタを適用
	if status != "" {
		animeList = s.filterByStatus(animeList, status)
	}

	// 制限を適用
	if limit > 0 && len(animeList) > limit {
		animeList = animeList[:limit]
	}

	// 出力形式に応じて結果を生成
	var result string
	var formatErr error
	switch format {
	case "json":
		result, formatErr = s.formatAsJSON(animeList)
	case "table":
		result, formatErr = s.formatAsTable(animeList)
	default:
		result, formatErr = s.formatAsJSON(animeList)
	}

	if formatErr != nil {
		return "", formatErr
	}

	// 出力ディレクトリが指定されている場合はファイルに保存
	if outputDir != "" {
		saveErr := s.saveToFile(result, outputDir, username, userID, format)
		if saveErr != nil {
			return "", fmt.Errorf("ファイル保存に失敗しました: %v", saveErr)
		}
		return fmt.Sprintf("結果をファイルに保存しました: %s\n", s.generateFileName(outputDir, username, userID, format)), nil
	}

	return result, nil
}

// transformToAnimeInfoList はAPIレスポンスをAnimeInfoのリストに変換する
func (s *AniListService) transformToAnimeInfoList(collection *domain.MediaListCollection) []domain.AnimeInfo {
	var animeList []domain.AnimeInfo

	for _, list := range collection.Lists {
		for _, entry := range list.Entries {
			if entry.Media == nil {
				continue
			}

			animeInfo := domain.AnimeInfo{
				ID:       entry.Media.ID,
				Score:    entry.Score,
				Status:   entry.Status,
				Progress: entry.Progress,
				Notes:    entry.Notes,
			}

			// タイトルを設定
			if entry.Media.Title != nil {
				animeInfo.Title = entry.Media.Title.Native
			}

			// カバー画像URLを設定
			if entry.Media.CoverImage != nil {
				animeInfo.CoverImageURL = entry.Media.CoverImage.ExtraLarge
			}

			// サイトURLを設定
			animeInfo.SiteURL = entry.Media.SiteURL

			// スタジオを設定
			if entry.Media.Studios != nil && len(entry.Media.Studios.Nodes) > 0 {
				animeInfo.Studio = entry.Media.Studios.Nodes[0].Name
			}

			// 完了日を設定
			animeInfo.CompletedAt = s.formatCompletedAt(entry.CompletedAt)

			// 更新日を設定
			animeInfo.UpdatedAt = time.Unix(entry.UpdatedAt, 0)

			animeList = append(animeList, animeInfo)
		}
	}

	return animeList
}

// formatCompletedAt は完了日をフォーマットする
func (s *AniListService) formatCompletedAt(date *domain.FuzzyDate) string {
	if date == nil || date.Year == nil {
		return ""
	}

	year := strconv.Itoa(*date.Year)
	month := "01"
	day := "01"

	if date.Month != nil {
		month = fmt.Sprintf("%02d", *date.Month)
	}
	if date.Day != nil {
		day = fmt.Sprintf("%02d", *date.Day)
	}

	return fmt.Sprintf("%s-%s-%s", year, month, day)
}

// filterByStatus はステータスでフィルタリングする
func (s *AniListService) filterByStatus(animeList []domain.AnimeInfo, status string) []domain.AnimeInfo {
	var filtered []domain.AnimeInfo
	for _, anime := range animeList {
		if anime.Status == status {
			filtered = append(filtered, anime)
		}
	}
	return filtered
}

// formatAsJSON はJSON形式で出力する
func (s *AniListService) formatAsJSON(animeList []domain.AnimeInfo) (string, error) {
	jsonBytes, err := json.MarshalIndent(animeList, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSONエンコードに失敗しました: %v", err)
	}
	return string(jsonBytes), nil
}

// formatAsTable はテーブル形式で出力する
func (s *AniListService) formatAsTable(animeList []domain.AnimeInfo) (string, error) {
	if len(animeList) == 0 {
		return "アニメが見つかりませんでした。\n", nil
	}

	var result strings.Builder

	// ヘッダー
	result.WriteString("ID\tタイトル\tステータス\tスコア\t進行状況\t完了日\tスタジオ\n")
	result.WriteString("---\t---\t---\t---\t---\t---\t---\n")

	// データ行
	for _, anime := range animeList {
		result.WriteString(fmt.Sprintf("%d\t%s\t%s\t%d\t%d\t%s\t%s\n",
			anime.ID,
			anime.Title,
			anime.Status,
			anime.Score,
			anime.Progress,
			anime.CompletedAt,
			anime.Studio,
		))
	}

	return result.String(), nil
}
