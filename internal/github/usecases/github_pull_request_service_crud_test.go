package usecases

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

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
