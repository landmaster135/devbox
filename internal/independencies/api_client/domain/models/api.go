package models

// APIRequest はAPIリクエストの情報を表します
type APIRequest struct {
	URL     string            // リクエスト先のURL
	Method  string            // HTTPメソッド（GET, POST, PUT, DELETE, etc.）
	Headers map[string]string // HTTPヘッダー
	Body    []byte            // リクエストボディ
}

// APIResponse はAPIレスポンスの情報を表します
type APIResponse struct {
	StatusCode int               // HTTPステータスコード
	Headers    map[string]string // HTTPヘッダー
	Body       []byte            // レスポンスボディ
}
