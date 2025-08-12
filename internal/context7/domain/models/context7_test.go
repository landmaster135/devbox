package models

import (
	"testing"
	"time"
)

// TestSearchResult_Normal は SearchResult 構造体の正常な初期化をテストします
func TestSearchResult_Normal(t *testing.T) {
	trustScore := 9.5
	stars := 1000
	versions := []string{"1.0.0", "2.0.0"}

	result := SearchResult{
		ID:              "/facebook/react",
		Title:           "React",
		Description:     "A JavaScript library for building user interfaces",
		Branch:          "main",
		LastUpdateDate:  "2024-01-01",
		State:           "finalized",
		TotalTokens:     50000,
		TotalSnippets:   100,
		TotalPages:      50,
		Stars:           &stars,
		TrustScore:      &trustScore,
		Versions:        versions,
	}

	if result.ID != "/facebook/react" {
		t.Errorf("Expected ID '/facebook/react', got '%s'", result.ID)
	}

	if result.Title != "React" {
		t.Errorf("Expected Title 'React', got '%s'", result.Title)
	}

	if result.TotalTokens != 50000 {
		t.Errorf("Expected TotalTokens 50000, got %d", result.TotalTokens)
	}

	if result.Stars == nil || *result.Stars != 1000 {
		t.Errorf("Expected Stars 1000, got %v", result.Stars)
	}

	if result.TrustScore == nil || *result.TrustScore != 9.5 {
		t.Errorf("Expected TrustScore 9.5, got %v", result.TrustScore)
	}

	if len(result.Versions) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(result.Versions))
	}
}

// TestSearchResponse_Normal は SearchResponse 構造体の正常な初期化をテストします
func TestSearchResponse_Normal(t *testing.T) {
	results := []SearchResult{
		{
			ID:    "/facebook/react",
			Title: "React",
		},
		{
			ID:    "/vuejs/vue",
			Title: "Vue.js",
		},
	}

	response := SearchResponse{
		Results: results,
		Error:   nil,
	}

	if len(response.Results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(response.Results))
	}

	if response.Error != nil {
		t.Errorf("Expected no error, got %v", response.Error)
	}
}

// TestSearchResponse_WithError はエラーを含む SearchResponse をテストします
func TestSearchResponse_WithError(t *testing.T) {
	errorMsg := "API呼び出しが失敗しました"
	response := SearchResponse{
		Results: []SearchResult{},
		Error:   &errorMsg,
	}

	if len(response.Results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(response.Results))
	}

	if response.Error == nil {
		t.Fatal("Expected error, got nil")
	}

	if *response.Error != errorMsg {
		t.Errorf("Expected error '%s', got '%s'", errorMsg, *response.Error)
	}
}

// TestDocOptions_Normal は DocOptions 構造体の正常な初期化をテストします
func TestDocOptions_Normal(t *testing.T) {
	options := DocOptions{
		Topic:  "hooks",
		Tokens: 5000,
	}

	if options.Topic != "hooks" {
		t.Errorf("Expected Topic 'hooks', got '%s'", options.Topic)
	}

	if options.Tokens != 5000 {
		t.Errorf("Expected Tokens 5000, got %d", options.Tokens)
	}
}

// TestDocOptions_Empty は空の DocOptions をテストします
func TestDocOptions_Empty(t *testing.T) {
	options := DocOptions{}

	if options.Topic != "" {
		t.Errorf("Expected empty Topic, got '%s'", options.Topic)
	}

	if options.Tokens != 0 {
		t.Errorf("Expected Tokens 0, got %d", options.Tokens)
	}
}

// TestContext7Request_Normal は Context7Request 構造体の正常な初期化をテストします
func TestContext7Request_Normal(t *testing.T) {
	headers := map[string]string{
		"Content-Type":     "application/json",
		"X-Context7-Source": "go-cli-client",
	}
	body := []byte(`{"query": "react"}`)

	request := Context7Request{
		URL:     "https://context7.com/api/v1/search",
		Method:  "GET",
		Headers: headers,
		Body:    body,
	}

	if request.URL != "https://context7.com/api/v1/search" {
		t.Errorf("Expected URL 'https://context7.com/api/v1/search', got '%s'", request.URL)
	}

	if request.Method != "GET" {
		t.Errorf("Expected Method 'GET', got '%s'", request.Method)
	}

	if len(request.Headers) != 2 {
		t.Errorf("Expected 2 headers, got %d", len(request.Headers))
	}

	if request.Headers["Content-Type"] != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", request.Headers["Content-Type"])
	}

	if string(request.Body) != `{"query": "react"}` {
		t.Errorf("Expected body '{\"query\": \"react\"}', got '%s'", string(request.Body))
	}
}

// TestContext7Response_Normal は Context7Response 構造体の正常な初期化をテストします
func TestContext7Response_Normal(t *testing.T) {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	body := []byte(`{"results": []}`)
	timestamp := time.Now()

	response := Context7Response{
		StatusCode: 200,
		Headers:    headers,
		Body:       body,
		Timestamp:  timestamp,
	}

	if response.StatusCode != 200 {
		t.Errorf("Expected StatusCode 200, got %d", response.StatusCode)
	}

	if len(response.Headers) != 1 {
		t.Errorf("Expected 1 header, got %d", len(response.Headers))
	}

	if response.Headers["Content-Type"] != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", response.Headers["Content-Type"])
	}

	if string(response.Body) != `{"results": []}` {
		t.Errorf("Expected body '{\"results\": []}', got '%s'", string(response.Body))
	}

	if response.Timestamp != timestamp {
		t.Errorf("Expected timestamp %v, got %v", timestamp, response.Timestamp)
	}
}

// TestDocumentState_Constants はDocumentState定数をテストします
func TestDocumentState_Constants(t *testing.T) {
	testCases := []struct {
		state    DocumentState
		expected string
	}{
		{DocumentStateInitial, "initial"},
		{DocumentStateFinalized, "finalized"},
		{DocumentStateError, "error"},
		{DocumentStateDelete, "delete"},
	}

	for _, testCase := range testCases {
		if string(testCase.state) != testCase.expected {
			t.Errorf("Expected DocumentState '%s', got '%s'", testCase.expected, string(testCase.state))
		}
	}
}

// TestConstants は定数の値をテストします
func TestConstants(t *testing.T) {
	if DefaultTokens != 10000 {
		t.Errorf("Expected DefaultTokens 10000, got %d", DefaultTokens)
	}

	if Context7APIBaseURL != "https://context7.com/api" {
		t.Errorf("Expected Context7APIBaseURL 'https://context7.com/api', got '%s'", Context7APIBaseURL)
	}
}

// TestDocumentState_StringConversion はDocumentStateの文字列変換をテストします
func TestDocumentState_StringConversion(t *testing.T) {
	state := DocumentStateFinalized
	stateStr := string(state)

	if stateStr != "finalized" {
		t.Errorf("Expected string 'finalized', got '%s'", stateStr)
	}
}

// TestSearchResult_OptionalFields はオプションフィールドがnilの場合をテストします
func TestSearchResult_OptionalFields(t *testing.T) {
	result := SearchResult{
		ID:             "/test/library",
		Title:          "Test Library",
		Description:    "A test library",
		Branch:         "main",
		LastUpdateDate: "2024-01-01",
		State:          "finalized",
		TotalTokens:    1000,
		TotalSnippets:  10,
		TotalPages:     5,
		Stars:          nil,
		TrustScore:     nil,
		Versions:       nil,
	}

	if result.Stars != nil {
		t.Errorf("Expected Stars to be nil, got %v", result.Stars)
	}

	if result.TrustScore != nil {
		t.Errorf("Expected TrustScore to be nil, got %v", result.TrustScore)
	}

	if result.Versions != nil {
		t.Errorf("Expected Versions to be nil, got %v", result.Versions)
	}
}
