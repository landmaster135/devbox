package usecases

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

// #==============================================================#
// ##          Test Helper Functions                            ##
// #==============================================================#

// createMockPullRequestService はテスト用のGitHubPullRequestServiceを作成します
func createMockPullRequestService(mockClient HTTPClient, jsonMarshaler JSONMarshaler) *GitHubPullRequestService {
	clientService := NewGitHubClientServiceWithDependencies(mockClient, "test-token", jsonMarshaler)
	return NewGitHubPullRequestServiceWithDependencies(clientService)
}

// #==============================================================#
// ##          Constructor Tests                                ##
// #==============================================================#

func TestNewGitHubPullRequestService(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "トークンあり",
			token: "test-token",
		},
		{
			name:  "トークンなし",
			token: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := NewGitHubPullRequestService(tc.token)

			if service == nil {
				t.Fatal("サービスがnilです")
			}

			if service.clientService == nil {
				t.Fatal("クライアントサービスがnilです")
			}
		})
	}
}

func TestNewGitHubPullRequestServiceWithDependencies(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{}
	clientService := NewGitHubClientServiceWithDependencies(mockClient, "test-token", &DefaultJSONMarshaler{})

	// Act
	service := NewGitHubPullRequestServiceWithDependencies(clientService)

	// Assert
	if service == nil {
		t.Fatal("サービスがnilです")
	}

	if service.clientService != clientService {
		t.Fatal("クライアントサービスが正しく設定されていません")
	}
}

// #==============================================================#
// ##          CreatePullRequest Tests                          ##
// #==============================================================#

func TestGitHubPullRequestService_CreatePullRequest_Normal(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			expectedURL := "https://api.github.com/repos/owner/repo/pulls"
			if req.URL.String() != expectedURL {
				t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
			}

			if req.Method != "POST" {
				t.Errorf("Expected method POST, got %s", req.Method)
			}

			responseBody := `{"id": 1, "title": "Test PR", "state": "open"}`
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	options := map[string]interface{}{
		"title": "Test PR",
		"head":  "feature-branch",
		"base":  "main",
	}

	// Act
	result, err := service.CreatePullRequest("owner", "repo", options)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["id"].(float64) != 1 {
		t.Errorf("Expected id 1, got %v", result["id"])
	}

	if result["title"].(string) != "Test PR" {
		t.Errorf("Expected title 'Test PR', got %v", result["title"])
	}
}

func TestGitHubPullRequestService_CreatePullRequest_Error(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"message": "Validation Failed", "documentation_url": "https://docs.github.com"}`
			return &http.Response{
				StatusCode: 422,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	options := map[string]interface{}{
		"title": "Test PR",
		"head":  "feature-branch",
		"base":  "main",
	}

	// Act
	_, err := service.CreatePullRequest("owner", "repo", options)

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	githubErr, ok := err.(*GitHubError)
	if !ok {
		t.Fatalf("Expected GitHubError, got %T", err)
	}

	if githubErr.StatusCode != 422 {
		t.Errorf("Expected status code 422, got %d", githubErr.StatusCode)
	}
}

// #==============================================================#
// ##          CreatePullRequestReview Tests                    ##
// #==============================================================#

func TestGitHubPullRequestService_CreatePullRequestReview_Normal(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			expectedURL := "https://api.github.com/repos/owner/repo/pulls/1/reviews"
			if req.URL.String() != expectedURL {
				t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
			}

			responseBody := `{"id": 1, "state": "APPROVED", "body": "LGTM"}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	options := map[string]interface{}{
		"event": "APPROVE",
		"body":  "LGTM",
	}

	// Act
	result, err := service.CreatePullRequestReview("owner", "repo", 1, options)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["id"].(float64) != 1 {
		t.Errorf("Expected id 1, got %v", result["id"])
	}
}

// #==============================================================#
// ##          MergePullRequest Tests                           ##
// #==============================================================#

func TestGitHubPullRequestService_MergePullRequest_Normal(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			expectedURL := "https://api.github.com/repos/owner/repo/pulls/1/merge"
			if req.URL.String() != expectedURL {
				t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
			}

			if req.Method != "PUT" {
				t.Errorf("Expected method PUT, got %s", req.Method)
			}

			responseBody := `{"sha": "abc123", "merged": true, "message": "Pull Request successfully merged"}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	options := map[string]interface{}{
		"commit_title":   "Merge PR",
		"commit_message": "Additional details",
		"merge_method":   "merge",
	}

	// Act
	result, err := service.MergePullRequest("owner", "repo", 1, options)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["merged"].(bool) != true {
		t.Errorf("Expected merged true, got %v", result["merged"])
	}
}

// #==============================================================#
// ##          GetPullRequestFiles Tests                        ##
// #==============================================================#

func TestGitHubPullRequestService_GetPullRequestFiles_Normal(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			expectedURL := "https://api.github.com/repos/owner/repo/pulls/1/files"
			if req.URL.String() != expectedURL {
				t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
			}

			responseBody := `[{"filename": "test.go", "status": "modified", "additions": 10, "deletions": 5}]`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.GetPullRequestFiles("owner", "repo", 1)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 file, got %d", len(result))
	}

	if result[0]["filename"].(string) != "test.go" {
		t.Errorf("Expected filename 'test.go', got %v", result[0]["filename"])
	}
}

// #==============================================================#
// ##          GetPullRequestStatus Tests                       ##
// #==============================================================#

func TestGitHubPullRequestService_GetPullRequestStatus_Normal(t *testing.T) {
	// Arrange
	callCount := 0
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			callCount++

			if callCount == 1 {
				// First call: get PR details
				expectedURL := "https://api.github.com/repos/owner/repo/pulls/1"
				if req.URL.String() != expectedURL {
					t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
				}

				responseBody := `{"head": {"sha": "abc123"}, "state": "open"}`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(responseBody)),
				}, nil
			} else {
				// Second call: get status
				expectedURL := "https://api.github.com/repos/owner/repo/commits/abc123/status"
				if req.URL.String() != expectedURL {
					t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
				}

				responseBody := `{"state": "success", "total_count": 1}`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(responseBody)),
				}, nil
			}
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.GetPullRequestStatus("owner", "repo", 1)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	pr := result["pull_request"].(map[string]interface{})
	if pr["state"].(string) != "open" {
		t.Errorf("Expected PR state 'open', got %v", pr["state"])
	}

	status := result["status"].(map[string]interface{})
	if status["state"].(string) != "success" {
		t.Errorf("Expected status state 'success', got %v", status["state"])
	}
}

// #==============================================================#
// ##          UpdatePullRequestBranch Tests                    ##
// #==============================================================#

func TestGitHubPullRequestService_UpdatePullRequestBranch_Normal(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			expectedURL := "https://api.github.com/repos/owner/repo/pulls/1/update-branch"
			if req.URL.String() != expectedURL {
				t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
			}

			if req.Method != "PUT" {
				t.Errorf("Expected method PUT, got %s", req.Method)
			}

			responseBody := `{"message": "Updating", "url": "https://api.github.com/repos/owner/repo/pulls/1"}`
			return &http.Response{
				StatusCode: 202,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.UpdatePullRequestBranch("owner", "repo", 1, "abc123")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["message"].(string) != "Updating" {
		t.Errorf("Expected message 'Updating', got %v", result["message"])
	}
}

// #==============================================================#
// ##          GetPullRequestComments Tests                     ##
// #==============================================================#

func TestGitHubPullRequestService_GetPullRequestComments_Normal(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			expectedURL := "https://api.github.com/repos/owner/repo/pulls/1/comments"
			if req.URL.String() != expectedURL {
				t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
			}

			responseBody := `[{"id": 1, "body": "Test comment", "user": {"login": "testuser"}}]`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.GetPullRequestComments("owner", "repo", 1)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 comment, got %d", len(result))
	}

	if result[0]["body"].(string) != "Test comment" {
		t.Errorf("Expected body 'Test comment', got %v", result[0]["body"])
	}
}

// #==============================================================#
// ##          GetPullRequestReviews Tests                      ##
// #==============================================================#

func TestGitHubPullRequestService_GetPullRequestReviews_Normal(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			expectedURL := "https://api.github.com/repos/owner/repo/pulls/1/reviews"
			if req.URL.String() != expectedURL {
				t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
			}

			responseBody := `[{"id": 1, "state": "APPROVED", "body": "LGTM", "user": {"login": "reviewer"}}]`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.GetPullRequestReviews("owner", "repo", 1)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 review, got %d", len(result))
	}

	if result[0]["state"].(string) != "APPROVED" {
		t.Errorf("Expected state 'APPROVED', got %v", result[0]["state"])
	}
}

// #==============================================================#
// ##          ListPullRequests Tests                           ##
// #==============================================================#

func TestGitHubPullRequestService_ListPullRequests_Normal(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// URLのベース部分を確認
			expectedBaseURL := "https://api.github.com/repos/owner/repo/pulls"
			if !strings.HasPrefix(req.URL.String(), expectedBaseURL) {
				t.Errorf("Expected URL to start with %s, got %s", expectedBaseURL, req.URL.String())
			}

			// クエリパラメータを個別に確認
			query := req.URL.Query()
			if query.Get("state") != "open" {
				t.Errorf("Expected state=open, got state=%s", query.Get("state"))
			}
			if query.Get("sort") != "created" {
				t.Errorf("Expected sort=created, got sort=%s", query.Get("sort"))
			}

			responseBody := `[{"id": 1, "title": "Test PR 1", "state": "open"}, {"id": 2, "title": "Test PR 2", "state": "open"}]`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	options := map[string]interface{}{
		"state": "open",
		"sort":  "created",
	}

	// Act
	result, err := service.ListPullRequests("owner", "repo", options)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 pull requests, got %d", len(result))
	}

	if result[0]["title"].(string) != "Test PR 1" {
		t.Errorf("Expected title 'Test PR 1', got %v", result[0]["title"])
	}
}

// #==============================================================#
// ##          Handler Method Tests                             ##
// #==============================================================#

func TestGitHubPullRequestService_HandleToCreatePullRequest_Normal(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"id": 1, "title": "Test PR", "state": "open"}`
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.HandleToCreatePullRequest("owner", "repo", "Test PR", "feature", "main", "Test body", false)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	if jsonResult["id"].(float64) != 1 {
		t.Errorf("Expected id 1, got %v", jsonResult["id"])
	}
}

func TestGitHubPullRequestService_HandleToListPullRequests_Normal(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `[{"id": 1, "title": "Test PR", "state": "open"}]`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.HandleToListPullRequests("owner", "repo", "open", "created", "desc", "", "", 30, 1)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var jsonResult []map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	if len(jsonResult) != 1 {
		t.Errorf("Expected 1 pull request, got %d", len(jsonResult))
	}
}

func TestGitHubPullRequestService_HandleToCreatePullRequestReview_Normal(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"id": 1, "state": "APPROVED", "body": "LGTM"}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.HandleToCreatePullRequestReview("owner", "repo", 1, "APPROVE", "LGTM")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	if jsonResult["id"].(float64) != 1 {
		t.Errorf("Expected id 1, got %v", jsonResult["id"])
	}
}

func TestGitHubPullRequestService_HandleToCreatePullRequestReview_EmptyParams(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"id": 1, "state": "COMMENTED"}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.HandleToCreatePullRequestReview("owner", "repo", 1, "", "")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	if jsonResult["id"].(float64) != 1 {
		t.Errorf("Expected id 1, got %v", jsonResult["id"])
	}
}

func TestGitHubPullRequestService_HandleToMergePullRequest_Normal(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"sha": "abc123", "merged": true, "message": "Pull Request successfully merged"}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.HandleToMergePullRequest("owner", "repo", 1, "Merge PR", "Additional details", "merge")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	if jsonResult["merged"].(bool) != true {
		t.Errorf("Expected merged true, got %v", jsonResult["merged"])
	}
}

func TestGitHubPullRequestService_HandleToMergePullRequest_EmptyParams(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"sha": "abc123", "merged": true, "message": "Pull Request successfully merged"}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.HandleToMergePullRequest("owner", "repo", 1, "", "", "")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	if jsonResult["merged"].(bool) != true {
		t.Errorf("Expected merged true, got %v", jsonResult["merged"])
	}
}

func TestGitHubPullRequestService_HandleToGetPullRequestFiles_Normal(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `[{"filename": "test.go", "status": "modified", "additions": 10, "deletions": 5}]`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.HandleToGetPullRequestFiles("owner", "repo", 1)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var jsonResult []map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	if len(jsonResult) != 1 {
		t.Errorf("Expected 1 file, got %d", len(jsonResult))
	}

	if jsonResult[0]["filename"].(string) != "test.go" {
		t.Errorf("Expected filename 'test.go', got %v", jsonResult[0]["filename"])
	}
}

func TestGitHubPullRequestService_HandleToGetPullRequestStatus_Normal(t *testing.T) {
	// Arrange
	callCount := 0
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			callCount++

			if callCount == 1 {
				responseBody := `{"head": {"sha": "abc123"}, "state": "open"}`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(responseBody)),
				}, nil
			} else {
				responseBody := `{"state": "success", "total_count": 1}`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(responseBody)),
				}, nil
			}
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.HandleToGetPullRequestStatus("owner", "repo", 1)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	pr := jsonResult["pull_request"].(map[string]interface{})
	if pr["state"].(string) != "open" {
		t.Errorf("Expected PR state 'open', got %v", pr["state"])
	}
}

func TestGitHubPullRequestService_HandleToUpdatePullRequestBranch_Normal(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"message": "Updating", "url": "https://api.github.com/repos/owner/repo/pulls/1"}`
			return &http.Response{
				StatusCode: 202,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.HandleToUpdatePullRequestBranch("owner", "repo", 1, "abc123")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	if jsonResult["message"].(string) != "Updating" {
		t.Errorf("Expected message 'Updating', got %v", jsonResult["message"])
	}
}

func TestGitHubPullRequestService_HandleToUpdatePullRequestBranch_EmptyExpectedSHA(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"message": "Updating", "url": "https://api.github.com/repos/owner/repo/pulls/1"}`
			return &http.Response{
				StatusCode: 202,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.HandleToUpdatePullRequestBranch("owner", "repo", 1, "")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	if jsonResult["message"].(string) != "Updating" {
		t.Errorf("Expected message 'Updating', got %v", jsonResult["message"])
	}
}

func TestGitHubPullRequestService_HandleToGetPullRequestComments_Normal(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `[{"id": 1, "body": "Test comment", "user": {"login": "testuser"}}]`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.HandleToGetPullRequestComments("owner", "repo", 1)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var jsonResult []map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	if len(jsonResult) != 1 {
		t.Errorf("Expected 1 comment, got %d", len(jsonResult))
	}

	if jsonResult[0]["body"].(string) != "Test comment" {
		t.Errorf("Expected body 'Test comment', got %v", jsonResult[0]["body"])
	}
}

func TestGitHubPullRequestService_HandleToGetPullRequestReviews_Normal(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `[{"id": 1, "state": "APPROVED", "body": "LGTM", "user": {"login": "reviewer"}}]`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.HandleToGetPullRequestReviews("owner", "repo", 1)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var jsonResult []map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	if len(jsonResult) != 1 {
		t.Errorf("Expected 1 review, got %d", len(jsonResult))
	}

	if jsonResult[0]["state"].(string) != "APPROVED" {
		t.Errorf("Expected state 'APPROVED', got %v", jsonResult[0]["state"])
	}
}

// #==============================================================#
// ##          Error Cases and Edge Cases Tests                ##
// #==============================================================#

func TestGitHubPullRequestService_HandleToCreatePullRequest_JSONMarshalError(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"id": 1, "title": "Test PR", "state": "open"}`
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}
	mockMarshaler := &MockJSONMarshaler{
		MarshalIndentFunc: func(v interface{}, prefix, indent string) ([]byte, error) {
			return nil, fmt.Errorf("marshal error")
		},
	}

	service := createMockPullRequestService(mockClient, mockMarshaler)

	// Act
	_, err := service.HandleToCreatePullRequest("owner", "repo", "Test PR", "feature", "main", "Test body", false)

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "marshal error" {
		t.Errorf("Expected 'marshal error', got %v", err.Error())
	}
}

func TestGitHubPullRequestService_CreatePullRequest_JSONUnmarshalError(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `invalid json`
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	options := map[string]interface{}{
		"title": "Test PR",
	}

	// Act
	_, err := service.CreatePullRequest("owner", "repo", options)

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGitHubPullRequestService_GetPullRequestStatus_InvalidHeadSHA(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// PRレスポンスでhead.shaが不正な形式
			responseBody := `{"head": {"invalid": "data"}, "state": "open"}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	_, err := service.GetPullRequestStatus("owner", "repo", 1)

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "could not get head SHA") {
		t.Errorf("Expected error about head SHA, got %v", err.Error())
	}
}

func TestGitHubPullRequestService_GetPullRequestStatus_FirstRequestError(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"message": "Not Found"}`
			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	_, err := service.GetPullRequestStatus("owner", "repo", 1)

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	githubErr, ok := err.(*GitHubError)
	if !ok {
		t.Fatalf("Expected GitHubError, got %T", err)
	}

	if githubErr.StatusCode != 404 {
		t.Errorf("Expected status code 404, got %d", githubErr.StatusCode)
	}
}

func TestGitHubPullRequestService_ListPullRequests_EmptyOptions(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			expectedURL := "https://api.github.com/repos/owner/repo/pulls"
			if req.URL.String() != expectedURL {
				t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
			}

			responseBody := `[]`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.ListPullRequests("owner", "repo", map[string]interface{}{})

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 pull requests, got %d", len(result))
	}
}

func TestGitHubPullRequestService_HandleToCreatePullRequest_EmptyBody(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"id": 1, "title": "Test PR", "state": "open"}`
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.HandleToCreatePullRequest("owner", "repo", "Test PR", "feature", "main", "", true)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	if jsonResult["id"].(float64) != 1 {
		t.Errorf("Expected id 1, got %v", jsonResult["id"])
	}
}

func TestGitHubPullRequestService_HandleToListPullRequests_EmptyParams(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			expectedURL := "https://api.github.com/repos/owner/repo/pulls"
			if req.URL.String() != expectedURL {
				t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
			}

			responseBody := `[{"id": 1, "title": "Test PR", "state": "open"}]`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.HandleToListPullRequests("owner", "repo", "", "", "", "", "", 0, 0)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var jsonResult []map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	if len(jsonResult) != 1 {
		t.Errorf("Expected 1 pull request, got %d", len(jsonResult))
	}
}

// #==============================================================#
// ##          Additional Error Cases                          ##
// #==============================================================#

func TestGitHubPullRequestService_HandleToCreatePullRequestReview_JSONMarshalError(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"id": 1, "state": "APPROVED", "body": "LGTM"}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}
	mockMarshaler := &MockJSONMarshaler{
		MarshalIndentFunc: func(v interface{}, prefix, indent string) ([]byte, error) {
			return nil, fmt.Errorf("marshal error")
		},
	}

	service := createMockPullRequestService(mockClient, mockMarshaler)

	// Act
	_, err := service.HandleToCreatePullRequestReview("owner", "repo", 1, "APPROVE", "LGTM")

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "marshal error" {
		t.Errorf("Expected 'marshal error', got %v", err.Error())
	}
}

func TestGitHubPullRequestService_HandleToMergePullRequest_JSONMarshalError(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"sha": "abc123", "merged": true, "message": "Pull Request successfully merged"}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}
	mockMarshaler := &MockJSONMarshaler{
		MarshalIndentFunc: func(v interface{}, prefix, indent string) ([]byte, error) {
			return nil, fmt.Errorf("marshal error")
		},
	}

	service := createMockPullRequestService(mockClient, mockMarshaler)

	// Act
	_, err := service.HandleToMergePullRequest("owner", "repo", 1, "Merge PR", "Additional details", "merge")

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "marshal error" {
		t.Errorf("Expected 'marshal error', got %v", err.Error())
	}
}

func TestGitHubPullRequestService_HandleToListPullRequests_JSONMarshalError(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `[{"id": 1, "title": "Test PR", "state": "open"}]`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}
	mockMarshaler := &MockJSONMarshaler{
		MarshalIndentFunc: func(v interface{}, prefix, indent string) ([]byte, error) {
			return nil, fmt.Errorf("marshal error")
		},
	}

	service := createMockPullRequestService(mockClient, mockMarshaler)

	// Act
	_, err := service.HandleToListPullRequests("owner", "repo", "open", "created", "desc", "", "", 30, 1)

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "marshal error" {
		t.Errorf("Expected 'marshal error', got %v", err.Error())
	}
}

func TestGitHubPullRequestService_GetPullRequestFiles_JSONUnmarshalError(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `invalid json`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	_, err := service.GetPullRequestFiles("owner", "repo", 1)

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGitHubPullRequestService_GetPullRequestComments_JSONUnmarshalError(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `invalid json`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	_, err := service.GetPullRequestComments("owner", "repo", 1)

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGitHubPullRequestService_GetPullRequestReviews_JSONUnmarshalError(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `invalid json`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	_, err := service.GetPullRequestReviews("owner", "repo", 1)

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGitHubPullRequestService_ListPullRequests_JSONUnmarshalError(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `invalid json`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	_, err := service.ListPullRequests("owner", "repo", map[string]interface{}{})

	// Assert
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
