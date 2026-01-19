package models

// HTTPRequest はHTTPリクエストの情報を表します
type HTTPRequest struct {
	URL      string            // リクエスト先のURL
	Method   string            // HTTPメソッド（GET, POST, PUT, DELETE, etc.）
	Headers  map[string]string // HTTPヘッダー
	Body     []byte            // リクエストボディ
	Encoding string            // 文字エンコーディング（shift_jis, utf-8, euc-jp, auto）
}

// HTTPResponse はHTTPレスポンスの情報を表します
type HTTPResponse struct {
	StatusCode int               // HTTPステータスコード
	Headers    map[string]string // HTTPヘッダー
	Body       []byte            // レスポンスボディ
	Warnings   []string          // レスポンスに関する注意事項
}
