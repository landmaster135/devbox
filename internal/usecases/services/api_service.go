package services

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/landmaster135/devbox/internal/domain/models"
	"github.com/landmaster135/devbox/internal/domain/repositories"
)

// APIService はAPIリクエストを処理するサービスです
type APIService struct {
	apiRepo repositories.APIRepository
}

// NewAPIService は新しいAPIServiceインスタンスを作成します
func NewAPIService(apiRepo repositories.APIRepository) *APIService {
	return &APIService{
		apiRepo: apiRepo,
	}
}

// SendRequestWithJSONFile はJSONファイルの内容をボディとしてAPIリクエストを送信します
func (s *APIService) SendRequestWithJSONFile(url, method, jsonFilePath string) (*models.APIResponse, error) {
	// デフォルトのヘッダーを設定
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	// カスタムヘッダーを指定せずに送信
	return s.SendRequestWithJSONFileAndHeaders(url, method, jsonFilePath, headers)
}

// SendRequestWithJSONFileAndHeaders はJSONファイルの内容をボディとして、指定されたヘッダーを含むAPIリクエストを送信します
func (s *APIService) SendRequestWithJSONFileAndHeaders(url, method, jsonFilePath string, headers map[string]string) (*models.APIResponse, error) {
	// JSONファイルを読み込む
	jsonBody, err := s.apiRepo.LoadJSONFile(jsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("JSONファイルの読み込みに失敗しました: %w", err)
	}

	// Content-Typeが設定されていない場合は追加
	if _, exists := headers["Content-Type"]; !exists {
		headers["Content-Type"] = "application/json"
	}

	// リクエストを作成
	request := &models.APIRequest{
		URL:     url,
		Method:  method,
		Headers: headers,
		Body:    jsonBody,
	}

	// リクエストを送信
	response, err := s.apiRepo.SendRequest(request)
	if err != nil {
		return nil, fmt.Errorf("APIリクエストの送信に失敗しました: %w", err)
	}

	return response, nil
}

// FormatResponse はAPIレスポンスを整形して文字列として返します
func (s *APIService) FormatResponse(response *models.APIResponse) (string, error) {
	// レスポンスボディがJSONの場合は整形する
	var prettyJSON bytes.Buffer
	if len(response.Body) > 0 {
		// JSONとして解析を試みる
		var jsonObj interface{}
		if err := json.Unmarshal(response.Body, &jsonObj); err == nil {
			// 整形したJSONを作成
			encoder := json.NewEncoder(&prettyJSON)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(jsonObj); err != nil {
				return "", fmt.Errorf("JSONの整形に失敗しました: %w", err)
			}
		} else {
			// JSONでない場合はそのまま返す
			prettyJSON.Write(response.Body)
		}
	}

	// レスポンス情報を整形
	result := fmt.Sprintf("Status: %d\n\nHeaders:\n", response.StatusCode)
	for key, value := range response.Headers {
		result += fmt.Sprintf("%s: %s\n", key, value)
	}
	result += fmt.Sprintf("\nBody:\n%s", prettyJSON.String())

	return result, nil
}
