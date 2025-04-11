package repositories

import (
	"github.com/landmaster135/devbox/internal/domain/models"
)

// APIRepository はAPIリクエストを実行するためのインターフェースです
type APIRepository interface {
	// SendRequest はAPIリクエストを送信し、レスポンスを返します
	SendRequest(request *models.APIRequest) (*models.APIResponse, error)

	// LoadJSONFile はJSONファイルを読み込み、バイト配列として返します
	LoadJSONFile(filePath string) ([]byte, error)
}
