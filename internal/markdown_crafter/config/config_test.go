package config

import (
	"testing"

	"github.com/landmaster135/devbox/internal/markdown_crafter/domain"
)

func TestNewConfig_SplitHeadings_Normal(t *testing.T) {
	t.Parallel()

	cfg, err := NewConfig(domain.OperationSplitHeadings, "./doc.md", 2, "./out", nil, "", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
	if cfg.Operation != domain.OperationSplitHeadings {
		t.Fatalf("unexpected operation: %s", cfg.Operation)
	}
}

func TestNewConfig_AddFrontMatter_Normal(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(domain.OperationAddFrontMatter, "./doc.md", 0, "", []string{"title=test"}, "", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
}

func TestNewConfig_AddTags_Normal(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(domain.OperationAddTags, "./doc.md", 0, "", nil, "go,markdown", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
}

func TestNewConfig_InvalidSplitHeadings(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(domain.OperationSplitHeadings, "./doc.md", 7, "./out", nil, "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}
