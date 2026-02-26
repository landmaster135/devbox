package usecases

import (
	"testing"

	"github.com/landmaster135/devbox/internal/mcp_infra_mark3labs/config"
	"github.com/landmaster135/devbox/internal/mcp_infra_mark3labs/infrastructures/mark3labs"
)

func TestService_NewService_DefaultAdapter_Normal(t *testing.T) {
	t.Parallel()

	service := NewService(nil)
	if service == nil {
		t.Fatal("service is nil")
	}
	if service.Mark3labs() == nil {
		t.Fatal("adapter is nil")
	}
}

func TestService_NewService_InjectedAdapter_Normal(t *testing.T) {
	t.Parallel()

	expectedAdapter := mark3labs.NewAdapterWithServeStdioFunc(func(s *config.MCPServer) error {
		return nil
	})

	service := NewService(expectedAdapter)
	if service == nil {
		t.Fatal("service is nil")
	}
	if service.Mark3labs() != expectedAdapter {
		t.Fatal("injected adapter was not used")
	}
}
