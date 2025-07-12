package usecases

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

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
