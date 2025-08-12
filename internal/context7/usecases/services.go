package usecases

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/landmaster135/devbox/internal/context7/domain/models"
	"github.com/landmaster135/devbox/internal/context7/interfaces"
)

// Context7Service はContext7 APIとの通信を担当するサービスです
type Context7Service struct {
	httpClient interfaces.HTTPClient
	baseURL    string
}

// NewContext7Service は新しいContext7Serviceインスタンスを作成します
func NewContext7Service(httpClient interfaces.HTTPClient) *Context7Service {
	return &Context7Service{
		httpClient: httpClient,
		baseURL:    models.Context7APIBaseURL,
	}
}

// NewContext7Service は新しいContext7Serviceインスタンスを作成します
func NewContext7ServiceWithHTTPClient() *Context7Service {
	c := interfaces.NewDefaultHTTPClient()
	return NewContext7Service(c)
}

// ResolveLibraryID はライブラリ名からContext7互換のライブラリIDを解決します
func (s *Context7Service) ResolveLibraryID(libraryName string) (*models.SearchResponse, error) {

	// 検索APIのURLを構築
	searchURL, err := url.Parse(fmt.Sprintf("%s/v1/search", s.baseURL))
	if err != nil {
		return nil, fmt.Errorf("検索URLの構築に失敗しました: %w", err)
	}

	// クエリパラメータを設定
	params := searchURL.Query()
	params.Set("query", libraryName)
	searchURL.RawQuery = params.Encode()

	// HTTPリクエストを作成
	req, err := http.NewRequest("GET", searchURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストの作成に失敗しました: %w", err)
	}

	// リクエストを実行
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストの実行に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスボディを読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンスボディの読み取りに失敗しました: %w", err)
	}

	// エラーレスポンスの処理
	if resp.StatusCode != http.StatusOK {
		errorMsg := fmt.Sprintf("API呼び出しが失敗しました (ステータス: %d)", resp.StatusCode)
		if resp.StatusCode == http.StatusTooManyRequests {
			errorMsg = "レート制限に達しました。しばらく待ってから再試行してください。"
		}
		return &models.SearchResponse{
			Results: []models.SearchResult{},
			Error:   &errorMsg,
		}, nil
	}

	// JSONレスポンスをパース
	var searchResponse models.SearchResponse
	if err := json.Unmarshal(body, &searchResponse); err != nil {
		return nil, fmt.Errorf("JSONレスポンスのパースに失敗しました: %w", err)
	}
	return &searchResponse, nil
}

// GetLibraryDocs はライブラリIDを使用してドキュメントを取得します
func (s *Context7Service) GetLibraryDocs(libraryID string, options models.DocOptions) (string, error) {
	// ライブラリIDの前のスラッシュを除去
	libraryID = strings.TrimPrefix(libraryID, "/")

	// ドキュメントAPIのURLを構築
	docURL, err := url.Parse(fmt.Sprintf("%s/v1/%s", s.baseURL, libraryID))
	if err != nil {
		return "", fmt.Errorf("ドキュメントURLの構築に失敗しました: %w", err)
	}

	// クエリパラメータを設定
	params := docURL.Query()
	if options.Tokens > 0 {
		params.Set("tokens", fmt.Sprintf("%d", options.Tokens))
	} else {
		params.Set("tokens", fmt.Sprintf("%d", models.DefaultTokens))
	}
	if options.Topic != "" {
		params.Set("topic", options.Topic)
	}
	params.Set("type", "txt")
	docURL.RawQuery = params.Encode()

	// HTTPリクエストを作成
	req, err := http.NewRequest("GET", docURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("HTTPリクエストの作成に失敗しました: %w", err)
	}

	// ヘッダーを設定
	req.Header.Set("X-Context7-Source", "go-cli-client")

	// リクエストを実行
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTPリクエストの実行に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスボディを読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("レスポンスボディの読み取りに失敗しました: %w", err)
	}

	// エラーレスポンスの処理
	if resp.StatusCode != http.StatusOK {
		errorMsg := fmt.Sprintf("ドキュメント取得が失敗しました (ステータス: %d)", resp.StatusCode)
		if resp.StatusCode == http.StatusTooManyRequests {
			errorMsg = "レート制限に達しました。しばらく待ってから再試行してください。"
		} else if resp.StatusCode == http.StatusNotFound {
			errorMsg = "指定されたライブラリIDのドキュメントが見つかりません。"
		}
		return errorMsg, nil
	}

	// レスポンステキストを確認
	responseText := string(body)
	if responseText == "" || responseText == "No content available" || responseText == "No context data available" {
		return "このライブラリのドキュメントは利用できません。", nil
	}
	return responseText, nil
}

// FormatSearchResults は検索結果を見やすい形式で整形します
func (s *Context7Service) FormatSearchResults(searchResponse *models.SearchResponse) string {
	if searchResponse.Error != nil {
		return fmt.Sprintf("エラー: %s", *searchResponse.Error)
	}

	if len(searchResponse.Results) == 0 {
		return "検索結果が見つかりませんでした。"
	}

	var buffer bytes.Buffer
	buffer.WriteString("検索結果:\n\n")

	for i, result := range searchResponse.Results {
		buffer.WriteString(fmt.Sprintf("%d. %s\n", i+1, result.Title))
		buffer.WriteString(fmt.Sprintf("   ID: %s\n", result.ID))
		buffer.WriteString(fmt.Sprintf("   説明: %s\n", result.Description))
		buffer.WriteString(fmt.Sprintf("   コードスニペット数: %d\n", result.TotalSnippets))

		if result.TrustScore != nil {
			buffer.WriteString(fmt.Sprintf("   信頼スコア: %.1f\n", *result.TrustScore))
		}

		if result.Stars != nil {
			buffer.WriteString(fmt.Sprintf("   スター数: %d\n", *result.Stars))
		}

		if len(result.Versions) > 0 {
			buffer.WriteString(fmt.Sprintf("   利用可能バージョン: %s\n", strings.Join(result.Versions, ", ")))
		}

		buffer.WriteString(fmt.Sprintf("   最終更新: %s\n", result.LastUpdateDate))
		buffer.WriteString(fmt.Sprintf("   状態: %s\n", result.State))
		buffer.WriteString("\n")
	}

	return buffer.String()
}

// ValidateLibraryID はライブラリIDの形式を検証します
func (s *Context7Service) ValidateLibraryID(libraryID string) error {
	if libraryID == "" {
		return fmt.Errorf("ライブラリIDが空です")
	}

	// スラッシュで始まる場合は除去
	libraryID = strings.TrimPrefix(libraryID, "/")

	// 基本的な形式チェック（org/project または org/project/version）
	parts := strings.Split(libraryID, "/")
	if len(parts) < 2 {
		return fmt.Errorf("ライブラリIDの形式が正しくありません。期待される形式: org/project または org/project/version")
	}

	return nil
}
