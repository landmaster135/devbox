package duckduckgo_search

import (
	"testing"

	usecases "github.com/landmaster135/devbox/internal/duckduckgo_search/usecases"
)

// TestServiceIntegration はサービス統合のテストです
func TestServiceIntegration(t *testing.T) {
	service := usecases.NewDuckDuckGoSearchService()
	if service == nil {
		t.Error("NewDuckDuckGoSearchService should not return nil")
	}
}

// TestBuildDuckDuckGoSearchServer はBuildDuckDuckGoSearchServer関数のテストです
func TestBuildDuckDuckGoSearchServer(t *testing.T) {
	// この関数は実際にサーバーを起動するため、テストでは呼び出さない
}
