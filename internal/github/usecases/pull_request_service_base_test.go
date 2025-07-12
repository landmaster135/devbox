package usecases

import (
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
	clientService := NewGitHubClientService("test-token")

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
