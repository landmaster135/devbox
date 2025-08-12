package models

import "time"

// SearchResult はContext7のライブラリ検索結果を表します
type SearchResult struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Branch          string    `json:"branch"`
	LastUpdateDate  string    `json:"lastUpdateDate"`
	State           string    `json:"state"`
	TotalTokens     int       `json:"totalTokens"`
	TotalSnippets   int       `json:"totalSnippets"`
	TotalPages      int       `json:"totalPages"`
	Stars           *int      `json:"stars,omitempty"`
	TrustScore      *float64  `json:"trustScore,omitempty"`
	Versions        []string  `json:"versions,omitempty"`
}

// SearchResponse はContext7のライブラリ検索APIレスポンスを表します
type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Error   *string        `json:"error,omitempty"`
}

// DocOptions はドキュメント取得時のオプションを表します
type DocOptions struct {
	Topic  string // 特定のトピックに焦点を当てる（例: "hooks", "routing"）
	Tokens int    // 取得する最大トークン数（デフォルト: 10000）
}

// Context7Request はContext7 APIへのリクエストを表します
type Context7Request struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    []byte
}

// Context7Response はContext7 APIからのレスポンスを表します
type Context7Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
	Timestamp  time.Time
}

// DocumentState はドキュメントの状態を表します
type DocumentState string

const (
	DocumentStateInitial   DocumentState = "initial"
	DocumentStateFinalized DocumentState = "finalized"
	DocumentStateError     DocumentState = "error"
	DocumentStateDelete    DocumentState = "delete"
)

// DefaultTokens はデフォルトのトークン数です
const DefaultTokens = 10000

// Context7APIBaseURL はContext7 APIのベースURLです
const Context7APIBaseURL = "https://context7.com/api"
