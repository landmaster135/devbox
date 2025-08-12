package usecases

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/arxiv/config"
)

// MockHTTPClient はHTTPクライアントのモック実装
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do はHTTPリクエストを実行する
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// TestArxivService_SearchPapers_Normal は正常な検索のテスト
func TestArxivService_SearchPapers_Normal(t *testing.T) {
	// Arrange
	mockResponse := `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:opensearch="http://a9.com/-/spec/opensearch/1.1/" xmlns:arxiv="http://arxiv.org/schemas/atom">
  <title>ArXiv Query: search_query=all:electron</title>
  <id>http://arxiv.org/api/test</id>
  <updated>2023-01-01T00:00:00-05:00</updated>
  <opensearch:totalResults>1</opensearch:totalResults>
  <opensearch:startIndex>0</opensearch:startIndex>
  <opensearch:itemsPerPage>1</opensearch:itemsPerPage>
  <entry>
    <id>http://arxiv.org/abs/2301.00001</id>
    <title>Test Paper Title</title>
    <summary>This is a test abstract.</summary>
    <published>2023-01-01T12:00:00-05:00</published>
    <updated>2023-01-01T12:00:00-05:00</updated>
    <author>
      <name>Test Author</name>
    </author>
    <arxiv:primary_category term="cs.AI" scheme="http://arxiv.org/schemas/atom"/>
    <category term="cs.AI" scheme="http://arxiv.org/schemas/atom"/>
    <arxiv:comment>Test comment</arxiv:comment>
    <link href="http://arxiv.org/abs/2301.00001v1" rel="alternate" type="text/html"/>
    <link title="pdf" href="http://arxiv.org/pdf/2301.00001v1" rel="related" type="application/pdf"/>
  </entry>
</feed>`

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(mockResponse)),
			}, nil
		},
	}

	service := NewArxivServiceWithHTTPClient(mockClient)

	// Act
	papers, err := service.SearchPapers("all:electron", 0, 1, "", "")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(papers) != 1 {
		t.Fatalf("Expected 1 paper, got %d", len(papers))
	}

	paper := papers[0]
	if paper.ID != "2301.00001" {
		t.Errorf("Expected ID '2301.00001', got '%s'", paper.ID)
	}

	if paper.Title != "Test Paper Title" {
		t.Errorf("Expected title 'Test Paper Title', got '%s'", paper.Title)
	}

	if len(paper.Authors) != 1 || paper.Authors[0] != "Test Author" {
		t.Errorf("Expected authors ['Test Author'], got %v", paper.Authors)
	}

	if paper.Abstract != "This is a test abstract." {
		t.Errorf("Expected abstract 'This is a test abstract.', got '%s'", paper.Abstract)
	}

	if paper.PrimaryCategory != "cs.AI" {
		t.Errorf("Expected primary category 'cs.AI', got '%s'", paper.PrimaryCategory)
	}
}

// TestArxivService_GetPapersByIds_Normal はID指定取得の正常テスト
func TestArxivService_GetPapersByIds_Normal(t *testing.T) {
	// Arrange
	mockResponse := `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:opensearch="http://a9.com/-/spec/opensearch/1.1/" xmlns:arxiv="http://arxiv.org/schemas/atom">
  <title>ArXiv Query: id_list=2301.00001</title>
  <id>http://arxiv.org/api/test</id>
  <updated>2023-01-01T00:00:00-05:00</updated>
  <opensearch:totalResults>1</opensearch:totalResults>
  <opensearch:startIndex>0</opensearch:startIndex>
  <opensearch:itemsPerPage>1</opensearch:itemsPerPage>
  <entry>
    <id>http://arxiv.org/abs/2301.00001</id>
    <title>Test Paper by ID</title>
    <summary>This is a test abstract for ID search.</summary>
    <published>2023-01-01T12:00:00-05:00</published>
    <updated>2023-01-01T12:00:00-05:00</updated>
    <author>
      <name>ID Test Author</name>
    </author>
    <arxiv:primary_category term="cs.LG" scheme="http://arxiv.org/schemas/atom"/>
    <category term="cs.LG" scheme="http://arxiv.org/schemas/atom"/>
    <link href="http://arxiv.org/abs/2301.00001v1" rel="alternate" type="text/html"/>
    <link title="pdf" href="http://arxiv.org/pdf/2301.00001v1" rel="related" type="application/pdf"/>
  </entry>
</feed>`

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(mockResponse)),
			}, nil
		},
	}

	service := NewArxivServiceWithHTTPClient(mockClient)

	// Act
	papers, err := service.GetPapersByIds([]string{"2301.00001"})

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(papers) != 1 {
		t.Fatalf("Expected 1 paper, got %d", len(papers))
	}

	paper := papers[0]
	if paper.ID != "2301.00001" {
		t.Errorf("Expected ID '2301.00001', got '%s'", paper.ID)
	}

	if paper.Title != "Test Paper by ID" {
		t.Errorf("Expected title 'Test Paper by ID', got '%s'", paper.Title)
	}
}

// TestArxivService_HandleSearch_SearchOperation は検索操作のハンドラテスト
func TestArxivService_HandleSearch_SearchOperation(t *testing.T) {
	// Arrange
	mockResponse := `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:opensearch="http://a9.com/-/spec/opensearch/1.1/" xmlns:arxiv="http://arxiv.org/schemas/atom">
  <title>ArXiv Query: search_query=all:test</title>
  <id>http://arxiv.org/api/test</id>
  <updated>2023-01-01T00:00:00-05:00</updated>
  <opensearch:totalResults>1</opensearch:totalResults>
  <opensearch:startIndex>0</opensearch:startIndex>
  <opensearch:itemsPerPage>1</opensearch:itemsPerPage>
  <entry>
    <id>http://arxiv.org/abs/2301.00001</id>
    <title>Handler Test Paper</title>
    <summary>This is a handler test abstract.</summary>
    <published>2023-01-01T12:00:00-05:00</published>
    <updated>2023-01-01T12:00:00-05:00</updated>
    <author>
      <name>Handler Test Author</name>
    </author>
    <arxiv:primary_category term="cs.AI" scheme="http://arxiv.org/schemas/atom"/>
    <category term="cs.AI" scheme="http://arxiv.org/schemas/atom"/>
    <link href="http://arxiv.org/abs/2301.00001v1" rel="alternate" type="text/html"/>
    <link title="pdf" href="http://arxiv.org/pdf/2301.00001v1" rel="related" type="application/pdf"/>
  </entry>
</feed>`

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(mockResponse)),
			}, nil
		},
	}

	service := NewArxivServiceWithHTTPClient(mockClient)

	cfg := &config.Config{
		Operation:   "search",
		SearchQuery: "all:test",
		Start:       0,
		MaxResults:  1,
	}

	// Act
	result, err := service.HandleSearch(cfg)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// JSONの解析
	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	// 結果の検証
	if jsonResult["operation"] != "search" {
		t.Errorf("Expected operation 'search', got '%v'", jsonResult["operation"])
	}

	totalCount, ok := jsonResult["total_count"].(float64)
	if !ok || totalCount != 1 {
		t.Errorf("Expected total_count 1, got %v", jsonResult["total_count"])
	}

	papers, ok := jsonResult["papers"].([]interface{})
	if !ok || len(papers) != 1 {
		t.Errorf("Expected 1 paper in results, got %v", jsonResult["papers"])
	}
}

// TestArxivService_HTTPError はHTTPエラーのテスト
func TestArxivService_HTTPError(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 500,
				Status:     "Internal Server Error",
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	service := NewArxivServiceWithHTTPClient(mockClient)

	// Act
	_, err := service.SearchPapers("all:test", 0, 1, "", "")

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "HTTPエラー: 500 Internal Server Error"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

// TestArxivService_NetworkError はネットワークエラーのテスト
func TestArxivService_NetworkError(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("network error")
		},
	}

	service := NewArxivServiceWithHTTPClient(mockClient)

	// Act
	_, err := service.SearchPapers("all:test", 0, 1, "", "")

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "HTTPリクエストに失敗"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

// TestArxivService_XMLParseError はXML解析エラーのテスト
func TestArxivService_XMLParseError(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("invalid xml")),
			}, nil
		},
	}

	service := NewArxivServiceWithHTTPClient(mockClient)

	// Act
	_, err := service.SearchPapers("all:test", 0, 1, "", "")

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "XML解析に失敗"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

// TestArxivService_APIError はArXiv APIエラーのテスト
func TestArxivService_APIError(t *testing.T) {
	// Arrange
	mockErrorResponse := `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>ArXiv Query Error</title>
  <entry>
    <title>Error</title>
    <summary>incorrect id format for test</summary>
  </entry>
</feed>`

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(mockErrorResponse)),
			}, nil
		},
	}

	service := NewArxivServiceWithHTTPClient(mockClient)

	// Act
	_, err := service.SearchPapers("all:test", 0, 1, "", "")

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "ArXiv APIエラー: incorrect id format for test"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TestArxivService_convertEntryToPaper_InvalidDate は日時解析エラーのテスト
func TestArxivService_convertEntryToPaper_InvalidDate(t *testing.T) {
	// Arrange
	service := NewArxivService()
	entry := ArxivEntry{
		ID:        "http://arxiv.org/abs/2301.00001",
		Title:     "Test Paper",
		Summary:   "Test abstract",
		Published: "invalid-date",
		Updated:   "2023-01-01T12:00:00-05:00",
		Authors:   []ArxivAuthor{{Name: "Test Author"}},
		Categories: []ArxivCategory{{Term: "cs.AI"}},
		PrimaryCategory: ArxivCategory{Term: "cs.AI"},
	}

	// Act
	_, err := service.convertEntryToPaper(entry)

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "公開日時の解析に失敗"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

// TestNewArxivService はコンストラクタのテスト
func TestNewArxivService(t *testing.T) {
	// Act
	service := NewArxivService()

	// Assert
	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}

	if service.baseURL != "http://export.arxiv.org/api/query" {
		t.Errorf("Expected baseURL 'http://export.arxiv.org/api/query', got '%s'", service.baseURL)
	}

	if service.httpClient == nil {
		t.Error("Expected httpClient to be set, got nil")
	}
}
