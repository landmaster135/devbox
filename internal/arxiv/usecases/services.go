package usecases

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/landmaster135/devbox/internal/arxiv/config"
)

// #==============================================================#
// ##          ArxivService                                      ##
// #==============================================================#

// HTTPClient はHTTPリクエストを行うためのインターフェース
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// DefaultHTTPClient は標準のHTTPクライアント実装
type DefaultHTTPClient struct{}

func (c *DefaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

// ArxivService はArXiv APIを使用して論文を検索・取得するサービス
type ArxivService struct {
	httpClient HTTPClient
	baseURL    string
}

// NewArxivService は新しいArxivServiceを作成する
func NewArxivService() *ArxivService {
	return &ArxivService{
		httpClient: &DefaultHTTPClient{},
		baseURL:    "http://export.arxiv.org/api/query",
	}
}

// NewArxivServiceWithHTTPClient はHTTPクライアントを指定してArxivServiceを作成する（テスト用）
func NewArxivServiceWithHTTPClient(httpClient HTTPClient) *ArxivService {
	return &ArxivService{
		httpClient: httpClient,
		baseURL:    "http://export.arxiv.org/api/query",
	}
}

// Paper はArXiv論文の情報を表す構造体
type Paper struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Authors         []string  `json:"authors"`
	Abstract        string    `json:"abstract"`
	Published       time.Time `json:"published"`
	Updated         time.Time `json:"updated"`
	Categories      []string  `json:"categories"`
	PrimaryCategory string    `json:"primary_category"`
	Comment         string    `json:"comment,omitempty"`
	JournalRef      string    `json:"journal_ref,omitempty"`
	DOI             string    `json:"doi,omitempty"`
	PDFUrl          string    `json:"pdf_url"`
	AbstractUrl     string    `json:"abstract_url"`
}

// ArxivResponse はArXiv APIのレスポンスを表す構造体（XML解析用）
type ArxivResponse struct {
	XMLName      xml.Name     `xml:"feed"`
	Title        string       `xml:"title"`
	ID           string       `xml:"id"`
	Updated      string       `xml:"updated"`
	TotalResults int          `xml:"totalResults"`
	StartIndex   int          `xml:"startIndex"`
	ItemsPerPage int          `xml:"itemsPerPage"`
	Entries      []ArxivEntry `xml:"entry"`
}

// ArxivEntry はArXiv APIの各エントリを表す構造体
type ArxivEntry struct {
	ID              string          `xml:"id"`
	Title           string          `xml:"title"`
	Summary         string          `xml:"summary"`
	Published       string          `xml:"published"`
	Updated         string          `xml:"updated"`
	Authors         []ArxivAuthor   `xml:"author"`
	Categories      []ArxivCategory `xml:"category"`
	PrimaryCategory ArxivCategory   `xml:"primary_category"`
	Comment         string          `xml:"comment"`
	JournalRef      string          `xml:"journal_ref"`
	DOI             string          `xml:"doi"`
	Links           []ArxivLink     `xml:"link"`
}

// ArxivAuthor は著者情報を表す構造体
type ArxivAuthor struct {
	Name string `xml:"name"`
}

// ArxivCategory はカテゴリ情報を表す構造体
type ArxivCategory struct {
	Term   string `xml:"term,attr"`
	Scheme string `xml:"scheme,attr"`
}

// ArxivLink はリンク情報を表す構造体
type ArxivLink struct {
	Href  string `xml:"href,attr"`
	Rel   string `xml:"rel,attr"`
	Type  string `xml:"type,attr"`
	Title string `xml:"title,attr"`
}

// SearchPapers は検索クエリを使用して論文を検索する
func (s *ArxivService) SearchPapers(searchQuery string, start, maxResults int, sortBy, sortOrder string) ([]Paper, error) {
	log.Printf("ArXiv論文検索を開始: query=%s, start=%d, maxResults=%d", searchQuery, start, maxResults)

	params := url.Values{}
	params.Set("search_query", searchQuery)
	params.Set("start", strconv.Itoa(start))
	params.Set("max_results", strconv.Itoa(maxResults))

	if sortBy != "" {
		params.Set("sortBy", sortBy)
	}
	if sortOrder != "" {
		params.Set("sortOrder", sortOrder)
	}

	papers, err := s.makeRequest(params)
	if err != nil {
		log.Printf("ArXiv論文検索に失敗: %v", err)
		return nil, err
	}

	log.Printf("ArXiv論文検索が完了: count=%d", len(papers))
	return papers, nil
}

// GetPapersByIds はIDリストを使用して論文を取得する
func (s *ArxivService) GetPapersByIds(ids []string) ([]Paper, error) {
	log.Printf("ArXiv論文ID取得を開始: ids=%v", ids)

	params := url.Values{}
	params.Set("id_list", strings.Join(ids, ","))

	papers, err := s.makeRequest(params)
	if err != nil {
		log.Printf("ArXiv論文ID取得に失敗: %v", err)
		return nil, err
	}

	log.Printf("ArXiv論文ID取得が完了: count=%d", len(papers))
	return papers, nil
}

// makeRequest はArXiv APIにリクエストを送信し、結果を解析する
func (s *ArxivService) makeRequest(params url.Values) ([]Paper, error) {
	reqURL := s.baseURL + "?" + params.Encode()
	log.Printf("ArXiv APIリクエスト: url=%s", reqURL)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成に失敗: %w", err)
	}

	req.Header.Set("User-Agent", "ArXiv-CLI-Tool/1.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストに失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTPエラー: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンス読み取りに失敗: %w", err)
	}

	var arxivResp ArxivResponse
	if err := xml.Unmarshal(body, &arxivResp); err != nil {
		return nil, fmt.Errorf("XML解析に失敗: %w", err)
	}

	// エラーエントリのチェック
	if len(arxivResp.Entries) == 1 && strings.Contains(arxivResp.Entries[0].Title, "Error") {
		return nil, fmt.Errorf("ArXiv APIエラー: %s", arxivResp.Entries[0].Summary)
	}

	papers := make([]Paper, 0, len(arxivResp.Entries))
	for _, entry := range arxivResp.Entries {
		paper, err := s.convertEntryToPaper(entry)
		if err != nil {
			log.Printf("エントリ変換に失敗: id=%s, error=%v", entry.ID, err)
			continue
		}
		papers = append(papers, paper)
	}

	return papers, nil
}

// convertEntryToPaper はArxivEntryをPaper構造体に変換する
func (s *ArxivService) convertEntryToPaper(entry ArxivEntry) (Paper, error) {
	// IDからarXiv IDを抽出（http://arxiv.org/abs/を除去）
	id := strings.TrimPrefix(entry.ID, "http://arxiv.org/abs/")

	// 著者リストを作成
	authors := make([]string, len(entry.Authors))
	for i, author := range entry.Authors {
		authors[i] = strings.TrimSpace(author.Name)
	}

	// カテゴリリストを作成
	categories := make([]string, len(entry.Categories))
	for i, category := range entry.Categories {
		categories[i] = category.Term
	}

	// 日時の解析
	published, err := time.Parse(time.RFC3339, entry.Published)
	if err != nil {
		return Paper{}, fmt.Errorf("公開日時の解析に失敗: %w", err)
	}

	updated, err := time.Parse(time.RFC3339, entry.Updated)
	if err != nil {
		return Paper{}, fmt.Errorf("更新日時の解析に失敗: %w", err)
	}

	// PDFとAbstractのURLを取得
	var pdfURL, abstractURL string
	for _, link := range entry.Links {
		switch {
		case link.Title == "pdf":
			pdfURL = link.Href
		case link.Rel == "alternate":
			abstractURL = link.Href
		}
	}

	return Paper{
		ID:              id,
		Title:           strings.TrimSpace(entry.Title),
		Authors:         authors,
		Abstract:        strings.TrimSpace(entry.Summary),
		Published:       published,
		Updated:         updated,
		Categories:      categories,
		PrimaryCategory: entry.PrimaryCategory.Term,
		Comment:         strings.TrimSpace(entry.Comment),
		JournalRef:      strings.TrimSpace(entry.JournalRef),
		DOI:             strings.TrimSpace(entry.DOI),
		PDFUrl:          pdfURL,
		AbstractUrl:     abstractURL,
	}, nil
}

// HandleSearch は設定に基づいて検索を実行し、JSON形式で結果を返す
func (s *ArxivService) HandleSearch(cfg *config.Config) (string, error) {
	log.Printf("ArXiv検索ハンドラを開始: operation=%s", cfg.Operation)

	var papers []Paper
	var err error

	switch cfg.Operation {
	case "search":
		papers, err = s.SearchPapers(cfg.SearchQuery, cfg.Start, cfg.MaxResults, cfg.SortBy, cfg.SortOrder)
	case "get_by_id":
		papers, err = s.GetPapersByIds(cfg.IdList)
	default:
		return "", fmt.Errorf("未対応の操作タイプです: %s", cfg.Operation)
	}

	if err != nil {
		return "", err
	}

	// 結果をJSON形式で返す
	result := map[string]interface{}{
		"operation":   cfg.Operation,
		"total_count": len(papers),
		"papers":      papers,
		"query_info": map[string]interface{}{
			"search_query": cfg.SearchQuery,
			"id_list":      cfg.IdList,
			"start":        cfg.Start,
			"max_results":  cfg.MaxResults,
			"sort_by":      cfg.SortBy,
			"sort_order":   cfg.SortOrder,
		},
	}

	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON変換に失敗: %w", err)
	}

	log.Printf("ArXiv検索ハンドラが完了: total_count=%d", len(papers))
	return string(jsonResult), nil
}
