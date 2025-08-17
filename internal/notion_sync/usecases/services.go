package usecases

import (
	"encoding/json"
	"fmt"

	"github.com/landmaster135/devbox/internal/http_request/interfaces/repositories"
	"github.com/landmaster135/devbox/internal/http_request/usecases/services"

	config "github.com/landmaster135/devbox/internal/notion_sync/config"
)

// PatchRequest はページパッチAPIのリクエスト構造体
type PatchRequest struct {
	Token           string `json:"token"`
	ConID           string `json:"con_id,omitempty"`
	PageID          string `json:"page_id,omitempty"`
	MarkdownContent string `json:"markdown_content"`
	ToggleH1        bool   `json:"toggle_h1"`
	ToggleH2        bool   `json:"toggle_h2"`
	ToggleH3        bool   `json:"toggle_h3"`
}

func NewPatchRequest(token, conID, pageID, markdownContent string, toggleH1, toggleH2, toggleH3 bool) *PatchRequest {
	return &PatchRequest{
		Token:           token,
		ConID:           conID,
		PageID:          pageID,
		MarkdownContent: markdownContent,
		ToggleH1:        toggleH1,
		ToggleH2:        toggleH2,
		ToggleH3:        toggleH3,
	}
}

// WebClipPatchRequest はWebクリップページパッチAPIのリクエスト構造体
type WebClipPatchRequest struct {
	Token           string `json:"token"`
	MarkdownContent string `json:"markdown_content"`
	Date            string `json:"date,omitempty"`
	Title           string `json:"title"`
	URL             string `json:"url"`
	ToggleH1        bool   `json:"toggle_h1"`
	ToggleH2        bool   `json:"toggle_h2"`
	ToggleH3        bool   `json:"toggle_h3"`
}

func NewWebClipPatchRequest(token, markdownContent, date, title, url string, toggleH1, toggleH2, toggleH3 bool) *WebClipPatchRequest {
	return &WebClipPatchRequest{
		Token:           token,
		MarkdownContent: markdownContent,
		Date:            date,
		Title:           title,
		URL:             url,
		ToggleH1:        toggleH1,
		ToggleH2:        toggleH2,
		ToggleH3:        toggleH3,
	}
}

// NotionSyncService はNotion同期を行うサービス
type NotionSyncService struct {
	httpService *services.HTTPService
}

// NewNotionSyncService は新しいNotionSyncServiceを作成する
func NewNotionSyncService() *NotionSyncService {
	httpRepo := repositories.NewHTTPRepository()
	httpService := services.NewHTTPService(httpRepo)

	return &NotionSyncService{
		httpService: httpService,
	}
}

// NewNotionSyncServiceWithDependencies はテスト用に依存性を注入できるNotionSyncServiceを作成する
func NewNotionSyncServiceWithDependencies(httpService *services.HTTPService) *NotionSyncService {
	return &NotionSyncService{
		httpService: httpService,
	}
}

// SendPatchRequest はページパッチリクエストを送信する
func (s *NotionSyncService) SendPatchRequest(cfg *config.Config) (string, error) {
	// PatchRequest構造体を作成
	patchReq := NewPatchRequest(
		cfg.Token,
		cfg.ConID,
		cfg.PageID,
		cfg.MarkdownContent,
		cfg.ToggleH1,
		cfg.ToggleH2,
		cfg.ToggleH3,
	)

	// JSONにマーシャル
	jsonBody, err := json.Marshal(patchReq)
	if err != nil {
		return "", fmt.Errorf("JSONマーシャルエラー: %v", err)
	}

	// ヘッダーを設定
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	// HTTPServiceの新しいメソッドを使用してリクエスト送信
	response, err := s.httpService.SendRequestWithJSONBody(cfg.EndpointURL, "POST", jsonBody, headers, "utf-8")
	if err != nil {
		return "", fmt.Errorf("HTTPリクエスト送信エラー: %v", err)
	}

	// レスポンスを整形
	formattedResponse, err := s.httpService.FormatResponse(response)
	if err != nil {
		return "", fmt.Errorf("レスポンス整形エラー: %v", err)
	}

	return formattedResponse, nil
}

// SendWebClipPatchRequest はWebクリップページパッチリクエストを送信する
func (s *NotionSyncService) SendWebClipPatchRequest(cfg *config.Config) (string, error) {
	// WebClipPatchRequest構造体を作成
	webClipReq := NewWebClipPatchRequest(
		cfg.Token,
		cfg.MarkdownContent,
		cfg.Date,
		cfg.Title,
		cfg.URL,
		cfg.ToggleH1,
		cfg.ToggleH2,
		cfg.ToggleH3,
	)

	// JSONにマーシャル
	jsonBody, err := json.Marshal(webClipReq)
	if err != nil {
		return "", fmt.Errorf("JSONマーシャルエラー: %v", err)
	}

	// ヘッダーを設定
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	// HTTPServiceの新しいメソッドを使用してリクエスト送信
	response, err := s.httpService.SendRequestWithJSONBody(cfg.EndpointURL, "POST", jsonBody, headers, "utf-8")
	if err != nil {
		return "", fmt.Errorf("HTTPリクエスト送信エラー: %v", err)
	}

	// レスポンスを整形
	formattedResponse, err := s.httpService.FormatResponse(response)
	if err != nil {
		return "", fmt.Errorf("レスポンス整形エラー: %v", err)
	}

	return formattedResponse, nil
}

// HandleNotionSync はNotion同期のメイン処理を行う
func (s *NotionSyncService) HandleNotionSync(cfg *config.Config) (string, error) {
	// バリデーション
	if cfg == nil {
		return "", fmt.Errorf("設定が指定されていません")
	}

	// 操作タイプによる分岐処理
	switch cfg.Operation {
	case "patch":
		// 通常のパッチリクエストを送信
		result, err := s.SendPatchRequest(cfg)
		if err != nil {
			return "", fmt.Errorf("パッチリクエスト送信に失敗しました: %v", err)
		}
		return result, nil

	case "patch-web-clip":
		// WebClipパッチリクエストを送信
		result, err := s.SendWebClipPatchRequest(cfg)
		if err != nil {
			return "", fmt.Errorf("WebClipパッチリクエスト送信に失敗しました: %v", err)
		}
		return result, nil

	default:
		return "", fmt.Errorf("サポートされていない操作タイプです: %s", cfg.Operation)
	}
}
