package usecases

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	youtube "github.com/kkdai/youtube/v2"
	"github.com/landmaster135/devbox/internal/youtube_downloader/domain"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// Service はYouTube動画ダウンロードサービスです
type Service struct {
	client *youtube.Client
}

// NewService は新しいサービスインスタンスを作成します
func NewService() *Service {
	client := &youtube.Client{
		HTTPClient:  http.DefaultClient,
		MaxRoutines: 10,
		ChunkSize:   10 * 1024 * 1024, // 10MB
	}
	return &Service{
		client: client,
	}
}

// NewServiceWithClient は指定されたクライアントでサービスを作成します
func NewServiceWithClient(client *youtube.Client) *Service {
	return &Service{
		client: client,
	}
}

// DownloadVideo は単一の動画をダウンロードします
func (s *Service) DownloadVideo(ctx context.Context, request domain.DownloadRequest) (string, error) {
	// クライアント設定を更新
	s.client.MaxRoutines = request.MaxRoutines
	s.client.ChunkSize = request.ChunkSize

	if request.Playlist {
		return s.downloadPlaylist(ctx, request)
	}

	return s.downloadSingleVideo(ctx, request)
}

// downloadSingleVideo は単一動画をダウンロードします
func (s *Service) downloadSingleVideo(ctx context.Context, request domain.DownloadRequest) (string, error) {
	// 動画情報を取得
	video, err := s.client.GetVideoContext(ctx, request.URL)
	if err != nil {
		return "", s.convertError(err, "")
	}

	// 動画情報を変換
	videoInfo := s.convertVideoInfo(video)

	// フォーマットを選択
	formats := s.convertFormats(video.Formats)
	selectedFormat, err := s.selectFormat(formats, request.Quality, request.AudioOnly)
	if err != nil {
		return "", err
	}

	// 出力ディレクトリを作成
	if err := os.MkdirAll(request.OutputDir, 0755); err != nil {
		return "", domain.NewDownloadError(
			domain.ErrorTypeFileSystemError,
			fmt.Sprintf("出力ディレクトリの作成に失敗しました: %v", err),
			video.ID,
			err,
		)
	}

	// ファイル名を生成
	fileName := s.generateSingleFileName(videoInfo, selectedFormat, request.AudioOnly)
	filePath := filepath.Join(request.OutputDir, fileName)

	// 映像のみフォーマットの場合は音声と結合
	if !request.AudioOnly && selectedFormat.AudioChannels == 0 && selectedFormat.Width > 0 {
		return s.downloadAndMergeVideoAudio(ctx, video, videoInfo, formats, selectedFormat, filePath)
	}

	// 通常のダウンロード実行
	err = s.downloadFile(ctx, video, selectedFormat, filePath)
	if err != nil {
		return "", err
	}

	// 結果を生成
	result := "ダウンロード完了:\n"
	result += fmt.Sprintf("  タイトル: %s\n", videoInfo.Title)
	result += fmt.Sprintf("  作者: %s\n", videoInfo.Author)
	result += fmt.Sprintf("  再生時間: %s\n", videoInfo.Duration)
	result += fmt.Sprintf("  品質: %s\n", selectedFormat.QualityLabel)
	result += fmt.Sprintf("  ファイル: %s\n", filePath)

	// ファイルサイズを取得
	if fileInfo, err := os.Stat(filePath); err == nil {
		result += fmt.Sprintf("  サイズ: %.2f MB\n", float64(fileInfo.Size())/(1024*1024))
	}

	return result, nil
}

// downloadPlaylist はプレイリストをダウンロードします
func (s *Service) downloadPlaylist(ctx context.Context, request domain.DownloadRequest) (string, error) {
	// プレイリスト情報を取得
	playlist, err := s.client.GetPlaylistContext(ctx, request.URL)
	if err != nil {
		return "", s.convertError(err, "")
	}

	result := "プレイリストダウンロード開始:\n"
	result += fmt.Sprintf("  タイトル: %s\n", playlist.Title)
	result += fmt.Sprintf("  作者: %s\n", playlist.Author)
	result += fmt.Sprintf("  動画数: %d\n\n", len(playlist.Videos))

	successCount := 0
	errorCount := 0

	// 各動画をダウンロード
	for i, entry := range playlist.Videos {
		result += fmt.Sprintf("(%d/%d) %s をダウンロード中...\n", i+1, len(playlist.Videos), entry.Title)

		// 動画を取得
		video, err := s.client.VideoFromPlaylistEntryContext(ctx, entry)
		if err != nil {
			result += fmt.Sprintf("  エラー: %v\n", err)
			errorCount++
			continue
		}

		// 動画情報を変換
		videoInfo := s.convertVideoInfo(video)

		// フォーマットを選択
		formats := s.convertFormats(video.Formats)
		selectedFormat, err := s.selectFormat(formats, request.Quality, request.AudioOnly)
		if err != nil {
			result += fmt.Sprintf("  エラー: %v\n", err)
			errorCount++
			continue
		}

		// ファイル名を生成（プレイリスト用）
		fileName := s.generatePlaylistFileName(i+1, videoInfo, selectedFormat, request.AudioOnly)
		filePath := filepath.Join(request.OutputDir, fileName)

		// ダウンロード実行
		err = s.downloadFile(ctx, video, selectedFormat, filePath)
		if err != nil {
			result += fmt.Sprintf("  エラー: %v\n", err)
			errorCount++
			continue
		}

		result += fmt.Sprintf("  完了: %s\n", fileName)
		successCount++
	}

	result += "\nプレイリストダウンロード完了:\n"
	result += fmt.Sprintf("  成功: %d件\n", successCount)
	result += fmt.Sprintf("  失敗: %d件\n", errorCount)
	result += fmt.Sprintf("  出力先: %s\n", request.OutputDir)

	return result, nil
}

// downloadFile はファイルをダウンロードします
func (s *Service) downloadFile(ctx context.Context, video *youtube.Video, format *domain.FormatInfo, filePath string) error {
	// YouTubeフォーマットに変換
	youtubeFormat := s.convertToYouTubeFormat(format)

	// ストリームを取得
	stream, size, err := s.client.GetStreamContext(ctx, video, youtubeFormat)
	if err != nil {
		return domain.NewDownloadError(
			domain.ErrorTypeDownloadFailed,
			fmt.Sprintf("ストリームの取得に失敗しました: %v", err),
			video.ID,
			err,
		)
	}
	defer stream.Close()

	// ファイルを作成
	file, err := os.Create(filePath)
	if err != nil {
		return domain.NewDownloadError(
			domain.ErrorTypeFileSystemError,
			fmt.Sprintf("ファイルの作成に失敗しました: %v", err),
			video.ID,
			err,
		)
	}
	defer file.Close()

	// ダウンロード実行
	_, err = io.Copy(file, stream)
	if err != nil {
		// 失敗した場合はファイルを削除
		os.Remove(filePath)
		return domain.NewDownloadError(
			domain.ErrorTypeDownloadFailed,
			fmt.Sprintf("ダウンロードに失敗しました: %v", err),
			video.ID,
			err,
		)
	}

	// サイズ検証
	if size > 0 {
		fileInfo, err := file.Stat()
		if err == nil && fileInfo.Size() != size {
			os.Remove(filePath)
			return domain.NewDownloadError(
				domain.ErrorTypeDownloadFailed,
				fmt.Sprintf("ファイルサイズが一致しません: 期待値=%d, 実際=%d", size, fileInfo.Size()),
				video.ID,
				nil,
			)
		}
	}

	return nil
}

// selectFormat はフォーマットを選択します
func (s *Service) selectFormat(formats []domain.FormatInfo, quality string, audioOnly bool) (*domain.FormatInfo, error) {
	var selector domain.QualitySelector

	if quality == "best" || quality == "" {
		selector = &domain.BestQualitySelector{}
	} else {
		selector = &domain.SpecificQualitySelector{Quality: quality}
	}

	return selector.SelectFormat(formats, audioOnly)
}

// convertVideoInfo はYouTube動画情報をドメインモデルに変換します
func (s *Service) convertVideoInfo(video *youtube.Video) domain.VideoInfo {
	thumbnails := make([]domain.Thumbnail, len(video.Thumbnails))
	for i, thumb := range video.Thumbnails {
		thumbnails[i] = domain.Thumbnail{
			URL:    thumb.URL,
			Width:  int(thumb.Width),
			Height: int(thumb.Height),
		}
	}

	return domain.VideoInfo{
		ID:          video.ID,
		Title:       video.Title,
		Description: video.Description,
		Author:      video.Author,
		ChannelID:   video.ChannelID,
		Duration:    video.Duration,
		Views:       video.Views,
		PublishDate: video.PublishDate,
		Thumbnails:  thumbnails,
	}
}

// convertFormats はYouTubeフォーマットをドメインモデルに変換します
func (s *Service) convertFormats(formats youtube.FormatList) []domain.FormatInfo {
	result := make([]domain.FormatInfo, len(formats))
	for i, format := range formats {
		result[i] = domain.FormatInfo{
			ItagNo:        format.ItagNo,
			URL:           format.URL,
			MimeType:      format.MimeType,
			Quality:       format.Quality,
			QualityLabel:  format.QualityLabel,
			Bitrate:       format.Bitrate,
			ContentLength: format.ContentLength,
			AudioChannels: format.AudioChannels,
			Width:         format.Width,
			Height:        format.Height,
			FPS:           format.FPS,
		}
	}
	return result
}

// convertToYouTubeFormat はドメインフォーマットをYouTubeフォーマットに変換します
func (s *Service) convertToYouTubeFormat(format *domain.FormatInfo) *youtube.Format {
	return &youtube.Format{
		ItagNo:        format.ItagNo,
		URL:           format.URL,
		MimeType:      format.MimeType,
		Quality:       format.Quality,
		QualityLabel:  format.QualityLabel,
		Bitrate:       format.Bitrate,
		ContentLength: format.ContentLength,
		AudioChannels: format.AudioChannels,
		Width:         format.Width,
		Height:        format.Height,
		FPS:           format.FPS,
	}
}

func (s *Service) generatePrefixOfFileName(index int) string {
	// 現在日付を取得（YYYYMMDD形式）
	jst, _ := time.LoadLocation("Asia/Tokyo")
	datePrefix := time.Now().In(jst).Format("20060102")

	if index >= 0 {
		return fmt.Sprintf("%s_%03d", datePrefix, index)
	}
	return datePrefix
}

func (s *Service) generateFullFileName(prefixOfFileName, title, suffixOfFileName string) string {
	return fmt.Sprintf("%s_%s_%s", prefixOfFileName, title, suffixOfFileName)
}

func (s *Service) generateFileName(index int, video domain.VideoInfo, format *domain.FormatInfo, audioOnly bool) string {
	// タイトルをファイル名に適した形に変換
	title := s.sanitizeFileName(video.Title)

	// 拡張子を決定
	ext := s.getFileExtension(format.MimeType, audioOnly)

	// ファイル名を生成（日付_連番_タイトル_動画ID.拡張子）
	prefixOfFileName := s.generatePrefixOfFileName(index)
	suffixOfFileName := fmt.Sprintf("%s.%s", video.ID, ext)
	fileName := s.generateFullFileName(prefixOfFileName, title, suffixOfFileName)

	// ファイル名の長さ制限
	if len(fileName) > 200 {
		maxTitleLen := 200 - len(fmt.Sprintf("%s_%s", prefixOfFileName, suffixOfFileName))
		if maxTitleLen > 0 {
			title = title[:maxTitleLen]
			fileName = s.generateFullFileName(prefixOfFileName, title, suffixOfFileName)
		}
	}

	return fileName
}

// generateSingleFileName は単一のファイル名を生成します
func (s *Service) generateSingleFileName(video domain.VideoInfo, format *domain.FormatInfo, audioOnly bool) string {
	return s.generateFileName(-1, video, format, audioOnly)
}

// generatePlaylistFileName はプレイリスト用のファイル名を生成します
func (s *Service) generatePlaylistFileName(index int, video domain.VideoInfo, format *domain.FormatInfo, audioOnly bool) string {
	return s.generateFileName(index, video, format, audioOnly)
}

// sanitizeFileName はファイル名を安全な形に変換します
func (s *Service) sanitizeFileName(name string) string {
	// 無効な文字を置換
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		"\n", "_",
		"\r", "_",
		"\t", "_",
		" ", "_", // 半角スペース
		"　", "_", // 全角スペース
	)

	result := replacer.Replace(name)

	// 連続するアンダースコアを単一のアンダースコアに変換
	for strings.Contains(result, "__") {
		result = strings.ReplaceAll(result, "__", "_")
	}

	// 前後のアンダースコアを削除
	result = strings.Trim(result, "_")

	// 空の場合はデフォルト名を使用
	if result == "" {
		result = "video"
	}

	return result
}

// getFileExtension はMIMEタイプから拡張子を取得します
func (s *Service) getFileExtension(mimeType string, audioOnly bool) string {
	if audioOnly {
		if strings.Contains(mimeType, "mp4") {
			return "m4a"
		}
		if strings.Contains(mimeType, "webm") {
			return "webm"
		}
		return "m4a"
	}

	if strings.Contains(mimeType, "mp4") {
		return "mp4"
	}
	if strings.Contains(mimeType, "webm") {
		return "webm"
	}
	return "mp4"
}

// downloadAndMergeVideoAudio は映像のみフォーマットと音声フォーマットを別々にダウンロードして結合します
func (s *Service) downloadAndMergeVideoAudio(ctx context.Context, video *youtube.Video, videoInfo domain.VideoInfo, formats []domain.FormatInfo, videoFormat *domain.FormatInfo, outputPath string) (string, error) {
	// 音声フォーマットを選択
	audioSelector := &domain.BestQualitySelector{}
	audioFormat, err := audioSelector.SelectFormat(formats, true)
	if err != nil {
		return "", domain.NewDownloadError(
			domain.ErrorTypeFormatNotFound,
			fmt.Sprintf("音声フォーマットが見つかりません: %v", err),
			video.ID,
			err,
		)
	}

	// 一時ファイルパスを生成
	tempDir := filepath.Dir(outputPath)
	videoTempPath := filepath.Join(tempDir, fmt.Sprintf("temp_video_%s.tmp", video.ID))
	audioTempPath := filepath.Join(tempDir, fmt.Sprintf("temp_audio_%s.tmp", video.ID))

	// 映像ファイルをダウンロード
	err = s.downloadFile(ctx, video, videoFormat, videoTempPath)
	if err != nil {
		return "", err
	}

	// 音声ファイルをダウンロード
	err = s.downloadFile(ctx, video, audioFormat, audioTempPath)
	if err != nil {
		// 映像ファイルを削除
		os.Remove(videoTempPath)
		return "", err
	}

	// 映像と音声を結合
	err = s.mergeVideoAudio(videoTempPath, audioTempPath, outputPath, video.ID)

	// 一時ファイルを削除
	os.Remove(videoTempPath)
	os.Remove(audioTempPath)

	if err != nil {
		return "", err
	}

	// 結果を生成
	result := "ダウンロード完了（映像+音声結合）:\n"
	result += fmt.Sprintf("  タイトル: %s\n", videoInfo.Title)
	result += fmt.Sprintf("  作者: %s\n", videoInfo.Author)
	result += fmt.Sprintf("  再生時間: %s\n", videoInfo.Duration)
	result += fmt.Sprintf("  映像品質: %s\n", videoFormat.QualityLabel)
	result += fmt.Sprintf("  音声品質: %dkbps\n", audioFormat.Bitrate/1000)
	result += fmt.Sprintf("  ファイル: %s\n", outputPath)

	// ファイルサイズを取得
	if fileInfo, err := os.Stat(outputPath); err == nil {
		result += fmt.Sprintf("  サイズ: %.2f MB\n", float64(fileInfo.Size())/(1024*1024))
	}

	return result, nil
}

// mergeVideoAudio は映像ファイルと音声ファイルを結合します
func (s *Service) mergeVideoAudio(videoPath, audioPath, outputPath, videoID string) error {
	videoInput := ffmpeg.Input(videoPath)
	audioInput := ffmpeg.Input(audioPath)

	err := ffmpeg.Output([]*ffmpeg.Stream{videoInput, audioInput}, outputPath, ffmpeg.KwArgs{
		"c:v": "copy",                     // 映像をそのままコピー（再エンコードなし）
		"c:a": "aac",                      // 音声をAACでエンコード
		"map": []string{"0:v:0", "1:a:0"}, // 映像と音声のマッピング
	}).OverWriteOutput().Run()

	if err != nil {
		return domain.NewDownloadError(
			domain.ErrorTypeDownloadFailed,
			fmt.Sprintf("映像と音声の結合に失敗しました: %v", err),
			videoID,
			err,
		)
	}

	return nil
}

// convertError はYouTubeライブラリのエラーをドメインエラーに変換します
func (s *Service) convertError(err error, videoID string) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	// エラーメッセージに基づいてエラータイプを判定
	switch {
	case strings.Contains(errStr, "video is private"):
		return domain.NewDownloadError(domain.ErrorTypeVideoPrivate, "この動画はプライベートです", videoID, err)
	case strings.Contains(errStr, "video not found"):
		return domain.NewDownloadError(domain.ErrorTypeVideoNotFound, "動画が見つかりません", videoID, err)
	case strings.Contains(errStr, "age restricted"):
		return domain.NewDownloadError(domain.ErrorTypeAgeRestricted, "年齢制限のある動画です", videoID, err)
	case strings.Contains(errStr, "unavailable"):
		return domain.NewDownloadError(domain.ErrorTypeVideoUnavailable, "動画が利用できません", videoID, err)
	case strings.Contains(errStr, "invalid") && strings.Contains(errStr, "url"):
		return domain.NewDownloadError(domain.ErrorTypeInvalidURL, "無効なURLです", videoID, err)
	case strings.Contains(errStr, "network") || strings.Contains(errStr, "connection"):
		return domain.NewDownloadError(domain.ErrorTypeNetworkError, "ネットワークエラーが発生しました", videoID, err)
	default:
		return domain.NewDownloadError(domain.ErrorTypeUnknown, fmt.Sprintf("不明なエラー: %v", err), videoID, err)
	}
}
