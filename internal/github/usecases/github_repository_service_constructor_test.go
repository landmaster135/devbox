package usecases

import (
	"testing"
)

// TestNewGitHubRepositoryService はNewGitHubRepositoryService関数をテストする
func TestNewGitHubRepositoryService(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "正常系 - トークンあり",
			token: "test_token",
		},
		{
			name:  "正常系 - トークンなし",
			token: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := NewGitHubRepositoryService(tc.token)

			if service == nil {
				t.Fatal("NewGitHubRepositoryServiceがnilを返しました")
			}

			if service.clientService == nil {
				t.Fatal("clientServiceがnilです")
			}
		})
	}
}
