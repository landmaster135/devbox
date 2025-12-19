package domain

import (
	"fmt"
	"time"
)

// DownloadRequest はダウンロード要求を表します
type DownloadRequest struct {
	URL         string
	OutputDir   string
	Quality     string
	Format      string
	AudioOnly   bool
	Playlist    bool
	MaxRoutines int
	ChunkSize   int64
}

// VideoInfo は動画情報を表します
type VideoInfo struct {
	ID          string
	Title       string
	Description string
	Author      string
	ChannelID   string
	Duration    time.Duration
	Views       int
	PublishDate time.Time
	Thumbnails  []Thumbnail
}

// Thumbnail はサムネイル情報を表します
type Thumbnail struct {
	URL    string
	Width  int
	Height int
}

// FormatInfo はフォーマット情報を表します
type FormatInfo struct {
	ItagNo        int
	URL           string
	MimeType      string
	Quality       string
	QualityLabel  string
	Bitrate       int
	ContentLength int64
	AudioChannels int
	Width         int
	Height        int
	FPS           int
}

// DownloadResult はダウンロード結果を表します
type DownloadResult struct {
	VideoInfo    VideoInfo
	FilePath     string
	FileSize     int64
	Duration     time.Duration
	Success      bool
	ErrorMessage string
}

// PlaylistInfo はプレイリスト情報を表します
type PlaylistInfo struct {
	ID          string
	Title       string
	Description string
	Author      string
	VideoCount  int
	Videos      []PlaylistEntry
}

// PlaylistEntry はプレイリストエントリを表します
type PlaylistEntry struct {
	ID       string
	Title    string
	Author   string
	Duration time.Duration
	Index    int
}

// ProgressInfo はダウンロード進捗情報を表します
type ProgressInfo struct {
	VideoID        string
	VideoTitle     string
	Downloaded     int64
	Total          int64
	Percentage     float64
	Speed          int64 // bytes per second
	ETA            time.Duration
	CurrentFile    string
	TotalFiles     int
	CompletedFiles int
}

// DownloadError はダウンロードエラーを表します
type DownloadError struct {
	Type    ErrorType
	Message string
	VideoID string
	Cause   error
}

// ErrorType はエラーの種類を表します
type ErrorType int

const (
	ErrorTypeUnknown ErrorType = iota
	ErrorTypeInvalidURL
	ErrorTypeVideoNotFound
	ErrorTypeVideoPrivate
	ErrorTypeVideoUnavailable
	ErrorTypeAgeRestricted
	ErrorTypeNetworkError
	ErrorTypeFileSystemError
	ErrorTypeFormatNotFound
	ErrorTypeDownloadFailed
)

// Error はDownloadErrorのエラーメッセージを返します
func (e *DownloadError) Error() string {
	if e.VideoID != "" {
		return fmt.Sprintf("[%s] %s: %s", e.VideoID, e.getTypeString(), e.Message)
	}
	return fmt.Sprintf("%s: %s", e.getTypeString(), e.Message)
}

// getTypeString はエラータイプの文字列表現を返します
func (e *DownloadError) getTypeString() string {
	switch e.Type {
	case ErrorTypeInvalidURL:
		return "無効なURL"
	case ErrorTypeVideoNotFound:
		return "動画が見つかりません"
	case ErrorTypeVideoPrivate:
		return "プライベート動画"
	case ErrorTypeVideoUnavailable:
		return "動画が利用できません"
	case ErrorTypeAgeRestricted:
		return "年齢制限"
	case ErrorTypeNetworkError:
		return "ネットワークエラー"
	case ErrorTypeFileSystemError:
		return "ファイルシステムエラー"
	case ErrorTypeFormatNotFound:
		return "フォーマットが見つかりません"
	case ErrorTypeDownloadFailed:
		return "ダウンロード失敗"
	default:
		return "不明なエラー"
	}
}

// Unwrap は元のエラーを返します
func (e *DownloadError) Unwrap() error {
	return e.Cause
}

// NewDownloadError は新しいDownloadErrorを作成します
func NewDownloadError(errorType ErrorType, message string, videoID string, cause error) *DownloadError {
	return &DownloadError{
		Type:    errorType,
		Message: message,
		VideoID: videoID,
		Cause:   cause,
	}
}

// QualitySelector は品質選択の戦略を表します
type QualitySelector interface {
	SelectFormat(formats []FormatInfo, audioOnly bool) (*FormatInfo, error)
}

// BestQualitySelector は最高品質を選択します
type BestQualitySelector struct{}

// calculateQualityScore はフォーマットの品質スコアを計算します
func (s *BestQualitySelector) calculateQualityScore(format *FormatInfo, audioOnly bool) int {
	if audioOnly {
		// 音声のみの場合はビットレートを重視
		return format.Bitrate
	}

	// 動画の場合は解像度とビットレートを考慮
	score := format.Width * format.Height
	if score == 0 {
		// 解像度情報がない場合はビットレートのみ
		score = format.Bitrate
	} else {
		// 解像度とビットレートの組み合わせ
		score += format.Bitrate / 1000 // ビットレートの重みを調整
	}

	return score
}

// selectAudioOnlyFormat は音声のみフォーマットを選択します
func (s *BestQualitySelector) selectAudioOnlyFormat(formats []FormatInfo) (*FormatInfo, error) {
	var bestFormat *FormatInfo
	var bestScore int

	for i := range formats {
		format := &formats[i]

		if format.AudioChannels == 0 {
			continue
		}

		score := s.calculateQualityScore(format, true)
		if bestFormat == nil || score > bestScore {
			bestFormat = format
			bestScore = score
		}
	}

	if bestFormat == nil {
		return nil, NewDownloadError(ErrorTypeFormatNotFound, "音声フォーマットが見つかりません", "", nil)
	}

	return bestFormat, nil
}

// selectAudioVideoFormat は音声付き動画フォーマットを選択します
func (s *BestQualitySelector) selectAudioVideoFormat(formats []FormatInfo) *FormatInfo {
	var bestFormat *FormatInfo
	var bestScore int

	for i := range formats {
		format := &formats[i]

		// 音声付き動画フォーマットのみを対象
		if format.AudioChannels == 0 || format.Width == 0 {
			continue
		}

		score := s.calculateQualityScore(format, false)
		if bestFormat == nil || score > bestScore {
			bestFormat = format
			bestScore = score
		}
	}

	return bestFormat
}

// selectVideoOnlyFormat は映像のみフォーマットを選択します
func (s *BestQualitySelector) selectVideoOnlyFormat(formats []FormatInfo) *FormatInfo {
	var bestFormat *FormatInfo
	var bestScore int

	for i := range formats {
		format := &formats[i]

		// 映像のみフォーマットを対象
		if format.AudioChannels > 0 || format.Width == 0 {
			continue
		}

		score := s.calculateQualityScore(format, false)
		if bestFormat == nil || score > bestScore {
			bestFormat = format
			bestScore = score
		}
	}

	return bestFormat
}

func (s *BestQualitySelector) SelectFormat(formats []FormatInfo, audioOnly bool) (*FormatInfo, error) {
	if len(formats) == 0 {
		return nil, NewDownloadError(ErrorTypeFormatNotFound, "利用可能なフォーマットがありません", "", nil)
	}

	if audioOnly {
		return s.selectAudioOnlyFormat(formats)
	}

	// 動画フォーマットの場合、音声付き一体型を優先
	audioVideoFormat := s.selectAudioVideoFormat(formats)
	if audioVideoFormat != nil {
		return audioVideoFormat, nil
	}

	// 音声付き一体型が見つからない場合は、映像のみフォーマットを選択
	// （後でFFmpegで音声と結合する）
	videoOnlyFormat := s.selectVideoOnlyFormat(formats)
	if videoOnlyFormat != nil {
		return videoOnlyFormat, nil
	}

	return nil, NewDownloadError(ErrorTypeFormatNotFound, "適切な動画フォーマットが見つかりません", "", nil)
}

// SpecificQualitySelector は特定の品質を選択します
type SpecificQualitySelector struct {
	Quality string
}

// matchesQuality は品質が一致するかチェックします
func (s *SpecificQualitySelector) matchesQuality(format *FormatInfo, quality string) bool {
	return format.Quality == quality ||
		format.QualityLabel == quality ||
		(quality == "720p" && format.Height == 720) ||
		(quality == "1080p" && format.Height == 1080) ||
		(quality == "480p" && format.Height == 480) ||
		(quality == "360p" && format.Height == 360)
}

func (s *SpecificQualitySelector) SelectFormat(formats []FormatInfo, audioOnly bool) (*FormatInfo, error) {
	if len(formats) == 0 {
		return nil, NewDownloadError(ErrorTypeFormatNotFound, "利用可能なフォーマットがありません", "", nil)
	}

	// 特定の品質を検索
	for i := range formats {
		format := &formats[i]

		if audioOnly && format.AudioChannels == 0 {
			continue
		}

		if !audioOnly && format.AudioChannels == 0 && format.Width == 0 {
			continue
		}

		if s.matchesQuality(format, s.Quality) {
			return format, nil
		}
	}

	// 見つからない場合は最高品質を選択
	bestSelector := &BestQualitySelector{}
	return bestSelector.SelectFormat(formats, audioOnly)
}
