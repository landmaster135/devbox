package models

import "time"

// AniListAnime はAniListからのアニメデータを表す構造体
type AniListAnime struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	Score         int     `json:"score"`
	Status        string  `json:"status"`
	Progress      int     `json:"progress"`
	CompletedAt   *string `json:"completed_at"`
	Notes         string  `json:"notes"`
	CoverImageURL string  `json:"cover_image_url"`
	SiteURL       string  `json:"site_url"`
	Studio        string  `json:"studio"`
	UpdatedAt     string  `json:"updated_at"`
}

// AdditionalAnimeData は追加のアニメデータを表す構造体
type AdditionalAnimeData struct {
	AniListID           int    `json:"anilist_id"`
	ConID               string `json:"con_id"`
	CoverImageURLEdited string `json:"cover_image_url_edited"`
	FleetingTier        *int   `json:"fleeting_tier"`
	FunnyTier           *int   `json:"funny_tier"`
	HeartwarmingTier    *int   `json:"heartwarming_tier"`
	MotivatingTier      *int   `json:"motivating_tier"`
	NihilisticTier      *int   `json:"nihilistic_tier"`
	TearjerkingTier     *int   `json:"tearjerking_tier"`
}

// OutputAnime は出力用のアニメデータを表す構造体
type OutputAnime struct {
	AniListID           int     `json:"anilist_id"`
	CompletedAt         *string `json:"completed_at"`
	ConID               string  `json:"con_id"`
	CoverImageURL       string  `json:"cover_image_url"`
	CoverImageURLEdited string  `json:"cover_image_url_edited"`
	FleetingTier        *int    `json:"fleeting_tier"`
	FunnyTier           *int    `json:"funny_tier"`
	HeartwarmingTier    *int    `json:"heartwarming_tier"`
	MotivatingTier      *int    `json:"motivating_tier"`
	NihilisticTier      *int    `json:"nihilistic_tier"`
	Notes               string  `json:"notes"`
	Progress            int     `json:"progress"`
	Score               int     `json:"score"`
	SiteURL             string  `json:"site_url"`
	Status              string  `json:"status"`
	Studio              string  `json:"studio"`
	TearjerkingTier     *int    `json:"tearjerking_tier"`
	Title               string  `json:"title"`
	UpdatedAt           int64   `json:"updated_at"`
}

// OutputData は最終的な出力データを表す構造体
type OutputData struct {
	Data        OutputDataContent `json:"data"`
	Description string            `json:"description"`
	Name        string            `json:"name"`
}

// OutputDataContent は出力データの内容を表す構造体
type OutputDataContent struct {
	Animes []OutputAnime `json:"animes"`
}

// ConvertToOutputAnime はAniListAnimeとAdditionalAnimeDataをOutputAnimeに変換する
func (a *AniListAnime) ConvertToOutputAnime(additional *AdditionalAnimeData) (*OutputAnime, error) {
	output := &OutputAnime{
		AniListID:     a.ID,
		Title:         a.Title,
		Score:         a.Score,
		Status:        a.Status,
		Progress:      a.Progress,
		Notes:         a.Notes,
		CoverImageURL: a.CoverImageURL,
		SiteURL:       a.SiteURL,
		Studio:        a.Studio,
	}

	// CompletedAtの変換
	if a.CompletedAt != nil && *a.CompletedAt != "" {
		output.CompletedAt = a.CompletedAt
	}

	// UpdatedAtの変換（ISO8601文字列からUnixタイムスタンプへ）
	if a.UpdatedAt != "" {
		t, err := time.Parse(time.RFC3339, a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		output.UpdatedAt = t.Unix()
	}

	// 追加データがある場合はマージ
	if additional != nil {
		output.ConID = additional.ConID
		output.CoverImageURLEdited = additional.CoverImageURLEdited

		// tierの値を設定（nullの場合はnullのまま）
		output.FleetingTier = additional.FleetingTier
		output.FunnyTier = additional.FunnyTier
		output.HeartwarmingTier = additional.HeartwarmingTier
		output.MotivatingTier = additional.MotivatingTier
		output.NihilisticTier = additional.NihilisticTier
		output.TearjerkingTier = additional.TearjerkingTier
	}

	return output, nil
}
