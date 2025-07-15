package usecases

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// #==============================================================#
// ##          Rate Limiting                                     ##
// #==============================================================#
// レート制限の定義
var rateLimit = struct {
	perSecond int
	perMonth  int
}{
	perSecond: 1,
	perMonth:  15000,
}

// リクエストカウンターの定義
var requestCount = struct {
	second    int
	month     int
	lastReset time.Time
	mu        sync.Mutex
}{
	second:    0,
	month:     0,
	lastReset: time.Now(),
}

// #==============================================================#
// ##          Data Structures                                   ##
// #==============================================================#
// BraveWebResult はWeb検索結果の構造体
type BraveWebResult struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Language    string `json:"language,omitempty"`
	Published   string `json:"published,omitempty"`
	Rank        int    `json:"rank,omitempty"`
}

// BraveWebResponse はWeb検索APIレスポンスの構造体
type BraveWebResponse struct {
	Web struct {
		Results []BraveWebResult `json:"results,omitempty"`
	} `json:"web,omitempty"`
	Locations struct {
		Results []struct {
			ID    string `json:"id"`
			Title string `json:"title,omitempty"`
		} `json:"results,omitempty"`
	} `json:"locations,omitempty"`
}

// BraveLocation はローカル検索結果の構造体
type BraveLocation struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address struct {
		StreetAddress   string `json:"streetAddress,omitempty"`
		AddressLocality string `json:"addressLocality,omitempty"`
		AddressRegion   string `json:"addressRegion,omitempty"`
		PostalCode      string `json:"postalCode,omitempty"`
	} `json:"address"`
	Coordinates struct {
		Latitude  float64 `json:"latitude,omitempty"`
		Longitude float64 `json:"longitude,omitempty"`
	} `json:"coordinates,omitempty"`
	Phone  string `json:"phone,omitempty"`
	Rating struct {
		RatingValue float64 `json:"ratingValue,omitempty"`
		RatingCount int     `json:"ratingCount,omitempty"`
	} `json:"rating,omitempty"`
	OpeningHours []string `json:"openingHours,omitempty"`
	PriceRange   string   `json:"priceRange,omitempty"`
}

// BravePoiResponse はPOI検索APIレスポンスの構造体
type BravePoiResponse struct {
	Results []BraveLocation `json:"results"`
}

// BraveDescription は説明APIレスポンスの構造体
type BraveDescription struct {
	Descriptions map[string]string `json:"descriptions"`
}

// #==============================================================#
// ##          Interfaces                                        ##
// #==============================================================#
// HTTPClient インターフェースを定義
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// DefaultHTTPClient は標準のhttp.Clientを使用する実装
type DefaultHTTPClient struct{}

func (c *DefaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}

// EnvironmentReader インターフェースを定義
type EnvironmentReader interface {
	Getenv(key string) string
}

// DefaultEnvironmentReader は標準のos.Getenvを使用する実装
type DefaultEnvironmentReader struct{}

func (e *DefaultEnvironmentReader) Getenv(key string) string {
	return os.Getenv(key)
}

// #==============================================================#
// ##          BraveSearchService                                ##
// #==============================================================#
// BraveSearchService はBrave検索を行うサービスです
type BraveSearchService struct {
	httpClient HTTPClient
	envReader  EnvironmentReader
}

// NewBraveSearchService は新しいBraveSearchServiceを作成します
func NewBraveSearchService() *BraveSearchService {
	return &BraveSearchService{
		httpClient: &DefaultHTTPClient{},
		envReader:  &DefaultEnvironmentReader{},
	}
}

// NewBraveSearchServiceWithDependencies はテスト用に依存性を注入できるBraveSearchServiceを作成します
func NewBraveSearchServiceWithDependencies(httpClient HTTPClient, envReader EnvironmentReader) *BraveSearchService {
	return &BraveSearchService{
		httpClient: httpClient,
		envReader:  envReader,
	}
}

// #==============================================================#
// ##          Rate Limiting Methods                             ##
// #==============================================================#
// checkRateLimit はレート制限をチェックする関数
func (s *BraveSearchService) checkRateLimit() error {
	requestCount.mu.Lock()
	defer requestCount.mu.Unlock()

	now := time.Now()
	if now.Sub(requestCount.lastReset) > time.Second {
		requestCount.second = 0
		requestCount.lastReset = now
	}

	if requestCount.second >= rateLimit.perSecond || requestCount.month >= rateLimit.perMonth {
		return errors.New("rate limit exceeded")
	}

	requestCount.second++
	requestCount.month++
	return nil
}

// #==============================================================#
// ##          Web Search Methods                                ##
// #==============================================================#
// performWebSearch はWeb検索を実行する関数
func (s *BraveSearchService) performWebSearch(query string, count int, offset int) (string, error) {
	if err := s.checkRateLimit(); err != nil {
		return "", err
	}

	apiKey := s.envReader.Getenv("BRAVE_API_KEY")
	if apiKey == "" {
		return "", errors.New("BRAVE_API_KEY environment variable is required")
	}

	// URLの構築
	baseURL := "https://api.search.brave.com/res/v1/web/search"
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// クエリパラメータの設定
	q := u.Query()
	q.Set("q", query)
	if count > 20 {
		count = 20 // API制限
	}
	q.Set("count", fmt.Sprintf("%d", count))
	q.Set("offset", fmt.Sprintf("%d", offset))
	u.RawQuery = q.Encode()

	// リクエストの作成
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}

	// ヘッダーの設定
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Subscription-Token", apiKey)

	// リクエストの実行
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error of Brave API: %d %s\n%s", resp.StatusCode, resp.Status, string(body))
	}

	// レスポンスの解析
	var data BraveWebResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	// 結果のフォーマット
	var results []string
	for _, result := range data.Web.Results {
		formattedResult := fmt.Sprintf("Title: %s\nDescription: %s\nURL: %s",
			result.Title, result.Description, result.URL)
		results = append(results, formattedResult)
	}

	return strings.Join(results, "\n\n"), nil
}

// HandleWebSearch はWeb検索のMCPリクエストを処理するハンドラーです
func (s *BraveSearchService) HandleWebSearch(query string, count, offset int) (string, error) {
	return s.performWebSearch(query, count, offset)
}

// #==============================================================#
// ##          Local Search Methods                              ##
// #==============================================================#
// performLocalSearch はローカル検索を実行する関数
func (s *BraveSearchService) performLocalSearch(query string, count int) (string, error) {
	if err := s.checkRateLimit(); err != nil {
		return "", err
	}

	apiKey := s.envReader.Getenv("BRAVE_API_KEY")
	if apiKey == "" {
		return "", errors.New("BRAVE_API_KEY environment variable is required")
	}

	// 最初の検索でロケーションIDを取得
	webURL, err := url.Parse("https://api.search.brave.com/res/v1/web/search")
	if err != nil {
		return "", err
	}

	q := webURL.Query()
	q.Set("q", query)
	q.Set("search_lang", "en")
	q.Set("result_filter", "locations")
	if count > 20 {
		count = 20 // API制限
	}
	q.Set("count", fmt.Sprintf("%d", count))
	webURL.RawQuery = q.Encode()

	// リクエストの作成
	webReq, err := http.NewRequest("GET", webURL.String(), nil)
	if err != nil {
		return "", err
	}

	// ヘッダーの設定
	webReq.Header.Add("Accept", "application/json")
	webReq.Header.Add("X-Subscription-Token", apiKey)

	// リクエストの実行
	webResp, err := s.httpClient.Do(webReq)
	if err != nil {
		return "", err
	}
	defer webResp.Body.Close()

	if webResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(webResp.Body)
		return "", fmt.Errorf("error of Brave API: %d %s\n%s", webResp.StatusCode, webResp.Status, string(body))
	}

	// レスポンスの解析
	var webData BraveWebResponse
	if err := json.NewDecoder(webResp.Body).Decode(&webData); err != nil {
		return "", err
	}

	// ロケーションIDの抽出
	var locationIDs []string
	for _, result := range webData.Locations.Results {
		if result.ID != "" {
			locationIDs = append(locationIDs, result.ID)
		}
	}

	// ロケーションIDがない場合はWeb検索にフォールバック
	if len(locationIDs) == 0 {
		return s.performWebSearch(query, count, 0)
	}

	// POIデータと説明を並行して取得
	poisCh := make(chan BravePoiResponse, 1)
	descCh := make(chan BraveDescription, 1)
	errCh := make(chan error, 2)

	// POIデータの取得
	go func() {
		pois, err := s.getPoisData(locationIDs)
		if err != nil {
			errCh <- err
			return
		}
		poisCh <- pois
	}()

	// 説明データの取得
	go func() {
		desc, err := s.getDescriptionsData(locationIDs)
		if err != nil {
			errCh <- err
			return
		}
		descCh <- desc
	}()

	// 結果の待機
	var poisData BravePoiResponse
	var descData BraveDescription
	for i := 0; i < 2; i++ {
		select {
		case err := <-errCh:
			return "", err
		case poisData = <-poisCh:
		case descData = <-descCh:
		}
	}

	// 結果のフォーマット
	return s.formatLocalResults(poisData, descData), nil
}

// HandleLocalSearch はローカル検索のMCPリクエストを処理するハンドラーです
func (s *BraveSearchService) HandleLocalSearch(query string, count int) (string, error) {
	return s.performLocalSearch(query, count)
}

// #==============================================================#
// ##          POI and Description Methods                       ##
// #==============================================================#
// getPoisData はPOIデータを取得する関数
func (s *BraveSearchService) getPoisData(ids []string) (BravePoiResponse, error) {
	if err := s.checkRateLimit(); err != nil {
		return BravePoiResponse{}, err
	}

	apiKey := s.envReader.Getenv("BRAVE_API_KEY")
	if apiKey == "" {
		return BravePoiResponse{}, errors.New("BRAVE_API_KEY environment variable is required")
	}

	// URLの構築
	u, err := url.Parse("https://api.search.brave.com/res/v1/local/pois")
	if err != nil {
		return BravePoiResponse{}, err
	}

	// クエリパラメータの設定
	q := u.Query()
	for _, id := range ids {
		if id != "" {
			q.Add("ids", id)
		}
	}
	u.RawQuery = q.Encode()

	// リクエストの作成
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return BravePoiResponse{}, err
	}

	// ヘッダーの設定
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Subscription-Token", apiKey)

	// リクエストの実行
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return BravePoiResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return BravePoiResponse{}, fmt.Errorf("error of Brave API: %d %s\n%s", resp.StatusCode, resp.Status, string(body))
	}

	// レスポンスの解析
	var poisResponse BravePoiResponse
	if err := json.NewDecoder(resp.Body).Decode(&poisResponse); err != nil {
		return BravePoiResponse{}, err
	}

	return poisResponse, nil
}

// getDescriptionsData は説明データを取得する関数
func (s *BraveSearchService) getDescriptionsData(ids []string) (BraveDescription, error) {
	if err := s.checkRateLimit(); err != nil {
		return BraveDescription{}, err
	}

	apiKey := s.envReader.Getenv("BRAVE_API_KEY")
	if apiKey == "" {
		return BraveDescription{}, errors.New("BRAVE_API_KEY environment variable is required")
	}

	// URLの構築
	u, err := url.Parse("https://api.search.brave.com/res/v1/local/descriptions")
	if err != nil {
		return BraveDescription{}, err
	}

	// クエリパラメータの設定
	q := u.Query()
	for _, id := range ids {
		if id != "" {
			q.Add("ids", id)
		}
	}
	u.RawQuery = q.Encode()

	// リクエストの作成
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return BraveDescription{}, err
	}

	// ヘッダーの設定
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Subscription-Token", apiKey)

	// リクエストの実行
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return BraveDescription{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return BraveDescription{}, fmt.Errorf("error of Brave API: %d %s\n%s", resp.StatusCode, resp.Status, string(body))
	}

	// レスポンスの解析
	var descData BraveDescription
	if err := json.NewDecoder(resp.Body).Decode(&descData); err != nil {
		return BraveDescription{}, err
	}

	return descData, nil
}

// #==============================================================#
// ##          Utility Methods                                   ##
// #==============================================================#
// formatLocalResults はローカル検索結果をフォーマットする関数
func (s *BraveSearchService) formatLocalResults(poisData BravePoiResponse, descData BraveDescription) string {
	if len(poisData.Results) == 0 {
		return "No local results found"
	}

	var formattedResults []string
	for _, poi := range poisData.Results {
		// 住所のフォーマット
		addressParts := []string{
			poi.Address.StreetAddress,
			poi.Address.AddressLocality,
			poi.Address.AddressRegion,
			poi.Address.PostalCode,
		}
		var addressParts2 []string
		for _, part := range addressParts {
			if part != "" {
				addressParts2 = append(addressParts2, part)
			}
		}
		address := "N/A"
		if len(addressParts2) > 0 {
			address = strings.Join(addressParts2, ", ")
		}

		// 営業時間のフォーマット
		hours := "N/A"
		if len(poi.OpeningHours) > 0 {
			hours = strings.Join(poi.OpeningHours, ", ")
		}

		// 評価のフォーマット
		rating := "N/A"
		if poi.Rating.RatingValue > 0 {
			rating = fmt.Sprintf("%.1f (%d reviews)", poi.Rating.RatingValue, poi.Rating.RatingCount)
		}

		// 説明の取得
		description := "No description available"
		if desc, ok := descData.Descriptions[poi.ID]; ok && desc != "" {
			description = desc
		}

		// 結果のフォーマット
		formattedResult := fmt.Sprintf(`Name: %s
Address: %s
Phone: %s
Rating: %s
Price Range: %s
Hours: %s
Description: %s`,
			poi.Name, address, s.getOrDefault(poi.Phone, "N/A"), rating,
			s.getOrDefault(poi.PriceRange, "N/A"), hours, description)

		formattedResults = append(formattedResults, formattedResult)
	}

	return strings.Join(formattedResults, "\n---\n")
}

// getOrDefault は空文字列の場合にデフォルト値を返すヘルパー関数
func (s *BraveSearchService) getOrDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
