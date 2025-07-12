package usecases

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

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
