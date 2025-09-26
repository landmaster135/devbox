package config

import (
	"strings"
	"testing"
)

func TestParseArgs_ListDeploymentsDefaults(t *testing.T) {
	cfg, err := ParseArgs([]string{"--operation=list-deployments"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Operation != OperationListDeployments {
		t.Fatalf("unexpected operation: %s", cfg.Operation)
	}

	expectedFormat := "table(name,insertTime,operation.operationType,operation.status,description)"
	if cfg.ListDeployments.Format != expectedFormat {
		t.Fatalf("unexpected format: %s", cfg.ListDeployments.Format)
	}
}

func TestParseArgs_RequiresOperation(t *testing.T) {
	if _, err := ParseArgs([]string{}); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestParseArgs_SimpleOverridesFormat(t *testing.T) {
	cfg, err := ParseArgs([]string{
		"--operation=list-deployments",
		"--format=custom",
		"--simple",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "table(name,insertTime)"
	if cfg.ListDeployments.Format != expected {
		t.Fatalf("expected format %s but got %s", expected, cfg.ListDeployments.Format)
	}

	if !cfg.ListDeployments.Simple {
		t.Fatal("expected simple flag to be true")
	}
}

func TestParseArgs_UnsupportedOperation(t *testing.T) {
	if _, err := ParseArgs([]string{"--operation=unknown"}); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestUsageContainsKeySections(t *testing.T) {
	usage := Usage()

	keywords := []string{
		"使用方法",
		"--operation=list-deployments",
		"--project=<PROJECT_ID>",
		"--simple",
	}

	for _, keyword := range keywords {
		if !strings.Contains(usage, keyword) {
			t.Fatalf("usage should contain %q", keyword)
		}
	}
}
