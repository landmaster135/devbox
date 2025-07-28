package usecases

import (
	"testing"
)

// #==============================================================#
// ##          Test Mock Implementations                        ##
// #==============================================================#

// MockGitBranchProvider はテスト用のGitBranchProviderモックです
type MockGitBranchProvider struct {
	CurrentBranch string
	Error         error
}

// GetCurrentBranchFromPath はモックの実装です
func (m *MockGitBranchProvider) GetCurrentBranchFromPath(absolutePath string) (string, error) {
	if m.Error != nil {
		return "", m.Error
	}
	return m.CurrentBranch, nil
}

// #==============================================================#
// ##          Test Helper Functions                            ##
// #==============================================================#

// createMockPullRequestService はテスト用のGitHubPullRequestServiceを作成します
func createMockPullRequestService(mockClient HTTPClient, jsonMarshaler JSONMarshaler) *GitHubPullRequestService {
	clientService := NewGitHubClientServiceWithDependencies(mockClient, "test-token", jsonMarshaler)
	mockGitBranchProvider := &MockGitBranchProvider{
		CurrentBranch: "test-branch",
		Error:         nil,
	}
	return NewGitHubPullRequestServiceWithDependencies(clientService, mockGitBranchProvider)
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
	mockGitBranchProvider := &MockGitBranchProvider{
		CurrentBranch: "test-branch",
		Error:         nil,
	}

	// Act
	service := NewGitHubPullRequestServiceWithDependencies(clientService, mockGitBranchProvider)

	// Assert
	if service == nil {
		t.Fatal("サービスがnilです")
	}

	if service.clientService != clientService {
		t.Fatal("クライアントサービスが正しく設定されていません")
	}

	if service.gitBranchProvider != mockGitBranchProvider {
		t.Fatal("GitBranchProviderが正しく設定されていません")
	}
}
