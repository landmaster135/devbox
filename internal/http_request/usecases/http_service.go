package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	models "github.com/landmaster135/devbox/internal/http_request/domain/models"
	interfaces "github.com/landmaster135/devbox/internal/http_request/interfaces"
	sanitizer "github.com/landmaster135/devbox/internal/http_request/usecases/sanitizer"
)

// HTTPService はHTTPリクエストを処理するサービスです
type HTTPService struct {
	httpRepo interfaces.HTTPRepository
}

// NewHTTPService は新しいHTTPServiceインスタンスを作成します
func NewHTTPService(httpRepo interfaces.HTTPRepository) *HTTPService {
	return &HTTPService{
		httpRepo: httpRepo,
	}
}

// SendRequestWithJSONFile はJSONファイルの内容をボディとしてHTTPリクエストを送信します
func (s *HTTPService) SendRequestWithJSONFile(url, method, jsonFilePath string) (*models.HTTPResponse, error) {
	// デフォルトのヘッダーを設定
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	// カスタムヘッダーを指定せずに送信（デフォルトエンコーディング: auto）
	return s.SendRequestWithJSONFileAndHeaders(url, method, jsonFilePath, headers, "auto")
}

// SendRequestWithJSONBody はメモリ上のJSONバイト配列をボディとしてHTTPリクエストを送信します
func (s *HTTPService) SendRequestWithJSONBody(url, method string, jsonBody []byte, headers map[string]string, encoding string) (*models.HTTPResponse, error) {
	// Content-Typeが設定されていない場合は追加
	if _, exists := headers["Content-Type"]; !exists {
		headers["Content-Type"] = "application/json"
	}

	// リクエストを作成
	request := &models.HTTPRequest{
		URL:      url,
		Method:   method,
		Headers:  headers,
		Body:     jsonBody,
		Encoding: encoding,
	}

	// リクエストを送信
	response, err := s.httpRepo.SendRequest(request)
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストの送信に失敗しました: %w", err)
	}

	return response, nil
}

// SendRequestWithJSONFileAndHeaders はJSONファイルの内容をボディとして、指定されたヘッダーを含むHTTPリクエストを送信します
func (s *HTTPService) SendRequestWithJSONFileAndHeaders(url, method, jsonFilePath string, headers map[string]string, encoding string) (*models.HTTPResponse, error) {
	// JSONファイルを読み込む
	jsonBody, err := s.httpRepo.LoadJSONFile(jsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("JSONファイルの読み込みに失敗しました: %w", err)
	}

	return s.SendRequestWithJSONBody(url, method, jsonBody, headers, encoding)
}

func (s *HTTPService) SendRequestWithoutJSONFile(url, method string, headers map[string]string, encoding string) (*models.HTTPResponse, error) {
	// リクエストを作成
	request := &models.HTTPRequest{
		URL:      url,
		Method:   method,
		Headers:  headers,
		Body:     nil,
		Encoding: encoding,
	}

	// リクエストを送信
	response, err := s.httpRepo.SendRequest(request)
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストの送信に失敗しました: %w", err)
	}

	return response, nil
}

// FormatResponse はHTTPレスポンスを整形して文字列として返します
func (s *HTTPService) FormatResponse(response *models.HTTPResponse) (string, error) {
	var prettyJSON bytes.Buffer
	isJSONBody := false
	if len(response.Body) > 0 {
		var jsonObj any
		if err := json.Unmarshal(response.Body, &jsonObj); err == nil {
			encoder := json.NewEncoder(&prettyJSON)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(jsonObj); err != nil {
				return "", fmt.Errorf("JSONの整形に失敗しました: %w", err)
			}
			isJSONBody = true
		} else {
			prettyJSON.Write(response.Body)
		}
	}

	warnings := append([]string{}, response.Warnings...)

	body := prettyJSON.String()
	if !isJSONBody {
		sanitizedBody, mainFound := sanitizer.SanitizeHTMLBody(body, false)
		body = sanitizedBody
		containsTags := strings.Contains(body, "<") && strings.Contains(body, ">")
		if containsTags && !mainFound {
			warnings = append(warnings, "main要素が見つからないため、HTMLボディをそのまま表示しました")
		}
	}

	result := fmt.Sprintf("Status: %d\n", response.StatusCode)
	if len(warnings) > 0 {
		result += "\nWarnings:\n"
		for _, warning := range warnings {
			result += fmt.Sprintf("- %s\n", warning)
		}
	}
	result += "\nHeaders:\n"
	for key, value := range response.Headers {
		result += fmt.Sprintf("%s: %s\n", key, value)
	}

	result += fmt.Sprintf("\nBody:\n%s", body)

	return result, nil
}
