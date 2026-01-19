package interfaces

import (
	models "github.com/landmaster135/devbox/internal/http_request/domain/models"
)

// HTTPRepository はHTTPリクエストを実行するためのインターフェースです
type HTTPRepository interface {
	// SendRequest はHTTPリクエストを送信し、レスポンスを返します
	SendRequest(request *models.HTTPRequest) (*models.HTTPResponse, error)

	// LoadJSONFile はJSONファイルを読み込み、バイト配列として返します
	LoadJSONFile(filePath string) ([]byte, error)
}
