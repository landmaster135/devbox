package config

import (
	"testing"
)

func TestNewConfig_SplitHeadings_Normal(t *testing.T) {
	t.Parallel()

	cfg, err := NewConfig(OperationSplitHeadings, "./doc.md", "", "", 2, "./out", nil, "", "", "", "", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
	if cfg.Operation != OperationSplitHeadings {
		t.Fatalf("unexpected operation: %s", cfg.Operation)
	}
}

func TestNewConfig_AddFrontMatter_Normal(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationAddFrontMatter, "./doc.md", "", "", 0, "", []string{"title=test"}, "", "", "", "", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
}

func TestNewConfig_AddTags_Normal(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationAddTags, "./doc.md", "", "", 0, "", nil, "go,markdown", "", "", "", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
}

func TestNewConfig_AddTagsByDir_Normal(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationAddTags, "", "./docs", "", 0, "", nil, "go,markdown", "", "", "", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
}

func TestNewConfig_DeleteEmptyFiles_Normal(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationDeleteEmptyFiles, "", "", "./docs", 0, "", nil, "", "", "", "", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
}

func TestNewConfig_AddHeading1_Normal(t *testing.T) {
	t.Parallel()

	cfg, err := NewConfig(OperationAddHeading1, "./doc.md", "", "", 0, "", nil, "", "概要", HeadingPositionHead, "", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
	if cfg.HeadingText != "概要" {
		t.Fatalf("unexpected heading text: %s", cfg.HeadingText)
	}
}

func TestNewConfig_InvalidSplitHeadings(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationSplitHeadings, "./doc.md", "", "", 7, "./out", nil, "", "", "", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewConfig_InvalidDeleteEmptyFiles(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationDeleteEmptyFiles, "", "", "", 0, "", nil, "", "", "", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewConfig_InvalidAddHeading1(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationAddHeading1, "./doc.md", "", "", 0, "", nil, "", "", HeadingPositionHead, "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}

	_, err = NewConfig(OperationAddHeading1, "./doc.md", "", "", 0, "", nil, "", "概要", "middle", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewConfig_InvalidAddTags(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationAddTags, "", "", "", 0, "", nil, "go,markdown", "", "", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}

	_, err = NewConfig(OperationAddTags, "./doc.md", "./docs", "", 0, "", nil, "go,markdown", "", "", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewConfig_ReplaceImages_Normal(t *testing.T) {
	t.Parallel()

	cfg, err := NewConfig(OperationReplaceImages, "./doc.md", "", "", 0, "", nil, "", "", "", "(添付画像)", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
	if cfg.ReplacementText != "(添付画像)" {
		t.Fatalf("unexpected replacement text: %s", cfg.ReplacementText)
	}
}

func TestNewConfig_RemoveHeadingAnnotations_Normal(t *testing.T) {
	t.Parallel()

	cfg, err := NewConfig(OperationRemoveHeadingAnnotations, "./doc.md", "", "", 3, "", nil, "", "", "", "", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
	if cfg.HeadingLevel != 3 {
		t.Fatalf("unexpected heading level: %d", cfg.HeadingLevel)
	}
}

func TestNewConfig_InvalidReplaceImages(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationReplaceImages, "./doc.md", "", "", 0, "", nil, "", "", "", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewConfig_InvalidRemoveHeadingAnnotations(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationRemoveHeadingAnnotations, "", "", "", 3, "", nil, "", "", "", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}

	_, err = NewConfig(OperationRemoveHeadingAnnotations, "./doc.md", "", "", 0, "", nil, "", "", "", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}
