package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/landmaster135/devbox/internal/api_client/domain/models"
)

// APIRepositoryImpl はAPIRepositoryインターフェースの実装です
type APIRepositoryImpl struct {
	client *http.Client
}

// NewAPIRepository は新しいAPIRepositoryインスタンスを作成します
func NewAPIRepository() *APIRepositoryImpl {
	return &APIRepositoryImpl{
		client: &http.Client{},
	}
}

// SendRequest はAPIリクエストを送信し、レスポンスを返します
func (r *APIRepositoryImpl) SendRequest(request *models.APIRequest) (*models.APIResponse, error) {
	// HTTPリクエストを作成
	req, err := http.NewRequest(request.Method, request.URL, bytes.NewBuffer(request.Body))
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストの作成に失敗しました: %w", err)
	}

	// ヘッダーを設定
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	// リクエストを送信
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("APIリクエストの送信に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスボディを読み込む
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンスボディの読み込みに失敗しました: %w", err)
	}

	// レスポンスヘッダーをマップに変換
	headers := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	// APIResponseを作成して返す
	return &models.APIResponse{
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       body,
	}, nil
}

// LoadJSONFile はJSONファイルを読み込み、バイト配列として返します
func (r *APIRepositoryImpl) LoadJSONFile(filePath string) ([]byte, error) {
	// ファイルを開く
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("JSONファイルを開けませんでした: %w", err)
	}
	defer file.Close()

	// ファイルの内容を読み込む
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("JSONファイルの読み込みに失敗しました: %w", err)
	}

	// BOM（Byte Order Mark）を削除
	if len(content) >= 3 && content[0] == 0xEF && content[1] == 0xBB && content[2] == 0xBF {
		content = content[3:]
	}

	// JSONとして有効かチェック
	var js interface{}
	if err := json.Unmarshal(content, &js); err != nil {
		return nil, fmt.Errorf("無効なJSON形式です: %w", err)
	}

	return content, nil
}
