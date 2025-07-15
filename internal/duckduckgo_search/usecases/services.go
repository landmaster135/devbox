package usecases

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

// #==============================================================#
// ##          Data Structures                                   ##
// #==============================================================#

// DuckDuckGoWebResult はWeb検索結果の構造体
type DuckDuckGoWebResult struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Rank        int    `json:"rank"`
}

// DuckDuckGoResponse はWeb検索APIレスポンスの構造体
type DuckDuckGoResponse struct {
	Results []DuckDuckGoWebResult `json:"results"`
}

// #==============================================================#
// ##          Interfaces                                        ##
// #==============================================================#

// HTTPClient インターフェースを定義
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// DefaultHTTPClient は標準のhttp.Clientを使用する実装
type DefaultHTTPClient struct {
	client *http.Client
}

func NewDefaultHTTPClient(timeout time.Duration) *DefaultHTTPClient {
	return &DefaultHTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *DefaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// #==============================================================#
// ##          Rate Limiting                                     ##
// #==============================================================#

// RateLimiter はレート制限を管理する構造体
type RateLimiter struct {
	perSecond int
	perMinute int
	second    int
	minute    int
	lastReset time.Time
	mu        sync.Mutex
}

// NewRateLimiter は新しいRateLimiterを作成します
func NewRateLimiter(perSecond, perMinute int) *RateLimiter {
	return &RateLimiter{
		perSecond: perSecond,
		perMinute: perMinute,
		second:    0,
		minute:    0,
		lastReset: time.Now(),
	}
}

// CheckRateLimit はレート制限をチェックします
func (r *RateLimiter) CheckRateLimit() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	if now.Sub(r.lastReset) > time.Second {
		r.second = 0
		r.lastReset = now
	}

	if now.Sub(r.lastReset) > time.Minute {
		r.minute = 0
	}

	if r.second >= r.perSecond || r.minute >= r.perMinute {
		return errors.New("rate limit exceeded")
	}

	r.second++
	r.minute++
	return nil
}

// #==============================================================#
// ##          DuckDuckGoSearchService                           ##
// #==============================================================#

// DuckDuckGoSearchService はDuckDuckGo検索を行うサービスです
type DuckDuckGoSearchService struct {
	httpClient  HTTPClient
	rateLimiter *RateLimiter
}

// NewDuckDuckGoSearchService は新しいDuckDuckGoSearchServiceを作成します
func NewDuckDuckGoSearchService() *DuckDuckGoSearchService {
	return &DuckDuckGoSearchService{
		httpClient:  NewDefaultHTTPClient(30 * time.Second),
		rateLimiter: NewRateLimiter(1, 30), // 1秒に1回、1分に30回
	}
}

// NewDuckDuckGoSearchServiceWithDependencies はテスト用に依存性を注入できるDuckDuckGoSearchServiceを作成します
func NewDuckDuckGoSearchServiceWithDependencies(httpClient HTTPClient, rateLimiter *RateLimiter) *DuckDuckGoSearchService {
	return &DuckDuckGoSearchService{
		httpClient:  httpClient,
		rateLimiter: rateLimiter,
	}
}

// cleanText はテキストをクリーンアップするメソッドです
func (s *DuckDuckGoSearchService) cleanText(text string) string {
	// HTMLエンティティのデコード
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")

	// HTMLタグの除去
	tagRegex := regexp.MustCompile(`<[^>]*>`)
	text = tagRegex.ReplaceAllString(text, "")

	// 余分な空白の除去
	text = strings.TrimSpace(text)
	spaceRegex := regexp.MustCompile(`\s+`)
	text = spaceRegex.ReplaceAllString(text, " ")

	return text
}

// decodeURL はURLをデコードするメソッドです
func (s *DuckDuckGoSearchService) decodeURL(rawURL string) (string, error) {
	// DuckDuckGoのリダイレクトURLの場合、実際のURLを抽出
	if strings.Contains(rawURL, "uddg=") {
		u, err := url.Parse(rawURL)
		if err != nil {
			return rawURL, err
		}

		uddg := u.Query().Get("uddg")
		if uddg != "" {
			decoded, err := url.QueryUnescape(uddg)
			if err == nil {
				return decoded, nil
			}
		}
	}

	// 通常のURL
	decoded, err := url.QueryUnescape(rawURL)
	if err != nil {
		return rawURL, nil
	}
	return decoded, nil
}

// parseSearchResults はHTMLから検索結果を抽出するメソッドです
func (s *DuckDuckGoSearchService) parseSearchResults(html string, maxResults int) ([]DuckDuckGoWebResult, error) {
	var results []DuckDuckGoWebResult

	// DuckDuckGoのHTML構造に基づく正規表現パターン
	// 結果のリンクを抽出
	linkPattern := `<a class="result__a" href="([^"]+)"[^>]*>([^<]+)</a>`
	linkRegex := regexp.MustCompile(linkPattern)
	linkMatches := linkRegex.FindAllStringSubmatch(html, -1)

	// 結果の説明を抽出
	descPattern := `<a class="result__snippet"[^>]*>([^<]+)</a>`
	descRegex := regexp.MustCompile(descPattern)
	descMatches := descRegex.FindAllStringSubmatch(html, -1)

	// より柔軟なパターンも試す
	if len(linkMatches) == 0 {
		// 代替パターン1
		linkPattern = `<h2[^>]*><a[^>]+href="([^"]+)"[^>]*>([^<]+)</a></h2>`
		linkRegex = regexp.MustCompile(linkPattern)
		linkMatches = linkRegex.FindAllStringSubmatch(html, -1)
	}

	if len(linkMatches) == 0 {
		// 代替パターン2 - より汎用的なリンクパターン
		linkPattern = `<a[^>]+href="([^"]+)"[^>]*>([^<]+)</a>`
		linkRegex = regexp.MustCompile(linkPattern)
		linkMatches = linkRegex.FindAllStringSubmatch(html, -1)

		// フィルタリング - 明らかに検索結果ではないリンクを除外
		var filteredMatches [][]string
		for _, match := range linkMatches {
			url := match[1]
			title := match[2]
			// 内部リンクや画像リンクを除外
			if !strings.Contains(url, "duckduckgo.com") &&
				!strings.HasPrefix(url, "#") &&
				!strings.HasPrefix(url, "javascript:") &&
				len(title) > 5 { // タイトルが短すぎるものを除外
				filteredMatches = append(filteredMatches, match)
			}
		}
		linkMatches = filteredMatches
	}

	// 結果を構築
	for i, linkMatch := range linkMatches {
		if i >= maxResults {
			break
		}

		url := linkMatch[1]
		title := linkMatch[2]

		// URLのデコード
		decodedURL, err := s.decodeURL(url)
		if err == nil {
			url = decodedURL
		}

		// タイトルのクリーンアップ
		title = s.cleanText(title)

		// 説明の取得（可能な場合）
		description := ""
		if i < len(descMatches) {
			description = s.cleanText(descMatches[i][1])
		}

		result := DuckDuckGoWebResult{
			Title:       title,
			Description: description,
			URL:         url,
			Rank:        i + 1,
		}

		results = append(results, result)
	}

	return results, nil
}

// HandleWebSearch はWeb検索を実行するハンドラーです
func (s *DuckDuckGoSearchService) HandleWebSearch(query string, count int) (string, error) {
	if err := s.rateLimiter.CheckRateLimit(); err != nil {
		return "", err
	}

	// URLの構築 - DuckDuckGoのHTML版を使用
	baseURL := "https://html.duckduckgo.com/html/"
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// クエリパラメータの設定
	q := u.Query()
	q.Set("q", query)
	if count > 50 {
		count = 50 // 制限
	}

	u.RawQuery = q.Encode()

	// リクエストの作成
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}

	// ヘッダーの設定 - ブラウザのように見せかける
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")

	// リクエストの実行
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error from DuckDuckGo: %d %s\n%s", resp.StatusCode, resp.Status, string(body))
	}

	// レスポンスの解析
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// HTMLから検索結果を抽出
	results, err := s.parseSearchResults(string(body), count)
	if err != nil {
		return "", err
	}

	// 結果のフォーマット
	var formattedResults []string
	for i, result := range results {
		if i >= count {
			break
		}
		formattedResult := fmt.Sprintf("Title: %s\nDescription: %s\nURL: %s",
			result.Title, result.Description, result.URL)
		formattedResults = append(formattedResults, formattedResult)
	}

	if len(formattedResults) == 0 {
		return "No search results found", nil
	}

	return strings.Join(formattedResults, "\n\n"), nil
}

// HandleInstantSearch は簡易検索を実行するハンドラーです（Instant Answer API使用）
func (s *DuckDuckGoSearchService) HandleInstantSearch(query string) (string, error) {
	if err := s.rateLimiter.CheckRateLimit(); err != nil {
		return "", err
	}

	// DuckDuckGo Instant Answer API
	baseURL := "https://api.duckduckgo.com/"
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// クエリパラメータの設定
	q := u.Query()
	q.Set("q", query)
	q.Set("format", "json")
	q.Set("no_html", "1")
	q.Set("skip_disambig", "1")
	u.RawQuery = q.Encode()

	// リクエストの作成
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}

	// ヘッダーの設定
	req.Header.Set("User-Agent", "DuckDuckGo-MCP-Server/1.0")
	req.Header.Set("Accept", "application/json")

	// リクエストの実行
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error from DuckDuckGo API: %d %s\n%s", resp.StatusCode, resp.Status, string(body))
	}

	// レスポンスの解析
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	// 結果のフォーマット
	var results []string

	// Abstract（要約）
	if abstract, ok := data["Abstract"].(string); ok && abstract != "" {
		results = append(results, fmt.Sprintf("Abstract: %s", abstract))

		if source, ok := data["AbstractSource"].(string); ok && source != "" {
			results = append(results, fmt.Sprintf("Source: %s", source))
		}

		if url, ok := data["AbstractURL"].(string); ok && url != "" {
			results = append(results, fmt.Sprintf("URL: %s", url))
		}
	}

	// Definition（定義）
	if definition, ok := data["Definition"].(string); ok && definition != "" {
		results = append(results, fmt.Sprintf("Definition: %s", definition))

		if source, ok := data["DefinitionSource"].(string); ok && source != "" {
			results = append(results, fmt.Sprintf("Source: %s", source))
		}
	}

	// Answer（回答）
	if answer, ok := data["Answer"].(string); ok && answer != "" {
		results = append(results, fmt.Sprintf("Answer: %s", answer))

		if answerType, ok := data["AnswerType"].(string); ok && answerType != "" {
			results = append(results, fmt.Sprintf("Type: %s", answerType))
		}
	}

	// Related Topics（関連トピック）
	if relatedTopics, ok := data["RelatedTopics"].([]interface{}); ok && len(relatedTopics) > 0 {
		results = append(results, "\nRelated Topics:")
		for i, topic := range relatedTopics {
			if i >= 3 { // 最初の3つのみ
				break
			}
			if topicMap, ok := topic.(map[string]interface{}); ok {
				if text, ok := topicMap["Text"].(string); ok && text != "" {
					results = append(results, fmt.Sprintf("- %s", text))
				}
			}
		}
	}

	if len(results) == 0 {
		return "No instant answer found", nil
	}

	return strings.Join(results, "\n"), nil
}
