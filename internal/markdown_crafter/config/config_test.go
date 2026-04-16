package config

import (
	"testing"
)

func TestNewConfig_SplitHeadings_Normal(t *testing.T) {
	t.Parallel()

	cfg, err := NewConfig(OperationSplitHeadings, "./doc.md", "", 2, 0, 0, "./out", nil, "", "", "", "", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
	if cfg.Operation != OperationSplitHeadings {
		t.Fatalf("unexpected operation: %s", cfg.Operation)
	}
}

func TestNewConfig_AddFrontMatter_Normal(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationAddFrontMatter, "./doc.md", "", 0, 0, 0, "", []string{"title=test"}, "", "", "", "", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
}

func TestNewConfig_AddTags_Normal(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationAddTags, "./doc.md", "", 0, 0, 0, "", nil, "go,markdown", "", "", "", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
}

func TestNewConfig_AddTagsByDir_Normal(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationAddTags, "", "./docs", 0, 0, 0, "", nil, "go,markdown", "", "", "", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
}

func TestNewConfig_DeleteEmptyFiles_Normal(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationDeleteEmptyFiles, "", "./docs", 0, 0, 0, "", nil, "", "", "", "", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
}

func TestNewConfig_AddHeading1_Normal(t *testing.T) {
	t.Parallel()

	cfg, err := NewConfig(OperationAddHeading1, "./doc.md", "", 0, 0, 0, "", nil, "", "概要", HeadingPositionHead, "", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
	if cfg.HeadingText != "概要" {
		t.Fatalf("unexpected heading text: %s", cfg.HeadingText)
	}
}

func TestNewConfig_InvalidSplitHeadings(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationSplitHeadings, "./doc.md", "", 7, 0, 0, "./out", nil, "", "", "", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewConfig_InvalidDeleteEmptyFiles(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationDeleteEmptyFiles, "", "", 0, 0, 0, "", nil, "", "", "", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewConfig_InvalidAddHeading1(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationAddHeading1, "./doc.md", "", 0, 0, 0, "", nil, "", "", HeadingPositionHead, "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}

	_, err = NewConfig(OperationAddHeading1, "./doc.md", "", 0, 0, 0, "", nil, "", "概要", "middle", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewConfig_InvalidAddTags(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationAddTags, "", "", 0, 0, 0, "", nil, "go,markdown", "", "", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}

	_, err = NewConfig(OperationAddTags, "./doc.md", "./docs", 0, 0, 0, "", nil, "go,markdown", "", "", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewConfig_ReplaceImages_Normal(t *testing.T) {
	t.Parallel()

	cfg, err := NewConfig(OperationReplaceImages, "./doc.md", "", 0, 0, 0, "", nil, "", "", "", "(添付画像)", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
	if cfg.ReplacementText != "(添付画像)" {
		t.Fatalf("unexpected replacement text: %s", cfg.ReplacementText)
	}
}

func TestNewConfig_RemoveHeadingAnnotations_Normal(t *testing.T) {
	t.Parallel()

	cfg, err := NewConfig(OperationRemoveHeadingAnnotations, "./doc.md", "", 3, 0, 0, "", nil, "", "", "", "", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
	if cfg.HeadingLevel != 3 {
		t.Fatalf("unexpected heading level: %d", cfg.HeadingLevel)
	}
}

func TestNewConfig_RemoveTitleHashTags_Normal(t *testing.T) {
	t.Parallel()

	cfg, err := NewConfig(OperationRemoveTitleHashTags, "", "./docs", 0, 1, 2, "", nil, "", "", "", "", false)
	if err != nil {
		t.Fatalf("NewConfig returned error: %v", err)
	}
	if cfg.DirPath != "./docs" {
		t.Fatalf("unexpected dir path: %s", cfg.DirPath)
	}
	if cfg.StartLine != 1 || cfg.EndLine != 2 {
		t.Fatalf("unexpected line range: start=%d end=%d", cfg.StartLine, cfg.EndLine)
	}
}

func TestNewConfig_InvalidReplaceImages(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationReplaceImages, "./doc.md", "", 0, 0, 0, "", nil, "", "", "", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewConfig_InvalidRemoveHeadingAnnotations(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationRemoveHeadingAnnotations, "", "", 3, 0, 0, "", nil, "", "", "", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}

	_, err = NewConfig(OperationRemoveHeadingAnnotations, "./doc.md", "", 0, 0, 0, "", nil, "", "", "", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewConfig_InvalidRemoveTitleHashTags(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationRemoveTitleHashTags, "", "", 0, 0, 0, "", nil, "", "", "", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewConfig_InvalidRemoveTitleHashTags_MissingStartLine(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationRemoveTitleHashTags, "", "./docs", 0, 0, 2, "", nil, "", "", "", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewConfig_InvalidRemoveTitleHashTags_MissingEndLine(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationRemoveTitleHashTags, "", "./docs", 0, 1, 0, "", nil, "", "", "", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewConfig_InvalidRemoveTitleHashTags_StartLineGreaterThanEndLine(t *testing.T) {
	t.Parallel()

	_, err := NewConfig(OperationRemoveTitleHashTags, "", "./docs", 0, 3, 2, "", nil, "", "", "", "", false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}
