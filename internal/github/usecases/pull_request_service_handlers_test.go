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

func TestGitHubPullRequestService_HandleToUpdatePullRequest_Normal(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// リクエストの検証
			expectedURL := "https://api.github.com/repos/owner/repo/pulls/1"
			if req.URL.String() != expectedURL {
				t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
			}

			if req.Method != "PATCH" {
				t.Errorf("Expected method PATCH, got %s", req.Method)
			}

			if req.Header.Get("Accept") != "application/vnd.github.v3+json" {
				t.Errorf("Expected Accept header application/vnd.github.v3+json, got %s", req.Header.Get("Accept"))
			}

			if req.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type header application/json, got %s", req.Header.Get("Content-Type"))
			}

			// リクエストボディの検証
			body, _ := io.ReadAll(req.Body)
			var requestData map[string]interface{}
			json.Unmarshal(body, &requestData)

			if requestData["title"].(string) != "Updated Title" {
				t.Errorf("Expected title 'Updated Title', got %v", requestData["title"])
			}

			if requestData["body"].(string) != "Updated Body" {
				t.Errorf("Expected body 'Updated Body', got %v", requestData["body"])
			}

			responseBody := `{"id": 1, "title": "Updated Title", "body": "Updated Body", "state": "open"}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.HandleToUpdatePullRequest("owner", "repo", 1, "Updated Title", "Updated Body")

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

	if jsonResult["title"].(string) != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got %v", jsonResult["title"])
	}

	if jsonResult["body"].(string) != "Updated Body" {
		t.Errorf("Expected body 'Updated Body', got %v", jsonResult["body"])
	}
}

func TestGitHubPullRequestService_HandleToUpdatePullRequest_TitleOnly(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// リクエストボディの検証
			body, _ := io.ReadAll(req.Body)
			var requestData map[string]interface{}
			json.Unmarshal(body, &requestData)

			if requestData["title"].(string) != "New Title Only" {
				t.Errorf("Expected title 'New Title Only', got %v", requestData["title"])
			}

			// bodyフィールドが含まれていないことを確認
			if _, exists := requestData["body"]; exists {
				t.Errorf("Expected body field to be absent, but it was present")
			}

			responseBody := `{"id": 1, "title": "New Title Only", "state": "open"}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.HandleToUpdatePullRequest("owner", "repo", 1, "New Title Only", "")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	if jsonResult["title"].(string) != "New Title Only" {
		t.Errorf("Expected title 'New Title Only', got %v", jsonResult["title"])
	}
}

func TestGitHubPullRequestService_HandleToUpdatePullRequest_BodyOnly(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// リクエストボディの検証
			body, _ := io.ReadAll(req.Body)
			var requestData map[string]interface{}
			json.Unmarshal(body, &requestData)

			if requestData["body"].(string) != "New Body Only" {
				t.Errorf("Expected body 'New Body Only', got %v", requestData["body"])
			}

			// titleフィールドが含まれていないことを確認
			if _, exists := requestData["title"]; exists {
				t.Errorf("Expected title field to be absent, but it was present")
			}

			responseBody := `{"id": 1, "body": "New Body Only", "state": "open"}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.HandleToUpdatePullRequest("owner", "repo", 1, "", "New Body Only")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	if jsonResult["body"].(string) != "New Body Only" {
		t.Errorf("Expected body 'New Body Only', got %v", jsonResult["body"])
	}
}

func TestGitHubPullRequestService_HandleToUpdatePullRequest_BothEmpty(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// このテストではHTTPリクエストは送信されないはず
			t.Error("HTTP request should not be made when both title and body are empty")
			return nil, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.HandleToUpdatePullRequest("owner", "repo", 1, "", "")

	// Assert
	if err == nil {
		t.Fatal("Expected error for empty title and body, got nil")
	}

	expectedError := "at least one of title or body must be provided"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if result != "" {
		t.Errorf("Expected empty result on error, got '%s'", result)
	}
}

func TestGitHubPullRequestService_HandleToUpdatePullRequest_HTTPError(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"message": "Not Found", "documentation_url": "https://docs.github.com"}`
			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}

	service := createMockPullRequestService(mockClient, &DefaultJSONMarshaler{})

	// Act
	result, err := service.HandleToUpdatePullRequest("owner", "repo", 999, "New Title", "New Body")

	// Assert
	if err == nil {
		t.Fatal("Expected error for HTTP 404, got nil")
	}

	// GitHubErrorの検証
	if githubErr, ok := err.(*GitHubError); ok {
		if githubErr.StatusCode != 404 {
			t.Errorf("Expected status code 404, got %d", githubErr.StatusCode)
		}
		if githubErr.Message != "Not Found" {
			t.Errorf("Expected message 'Not Found', got '%s'", githubErr.Message)
		}
	} else {
		t.Errorf("Expected GitHubError, got %T", err)
	}

	if result != "" {
		t.Errorf("Expected empty result on error, got '%s'", result)
	}
}
