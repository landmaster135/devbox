package usecases

import (
	"encoding/json"
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
