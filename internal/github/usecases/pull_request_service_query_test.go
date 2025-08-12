package usecases

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

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
	result, err := service.getPullRequestFiles("owner", "repo", 1)

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
	result, err := service.getPullRequestStatus("owner", "repo", 1)

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
	result, err := service.getPullRequestComments("owner", "repo", 1)

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
	result, err := service.getPullRequestReviews("owner", "repo", 1)

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
	result, err := service.listPullRequests("owner", "repo", options)

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
	result, err := service.listPullRequests("owner", "repo", map[string]interface{}{})

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 pull requests, got %d", len(result))
	}
}
