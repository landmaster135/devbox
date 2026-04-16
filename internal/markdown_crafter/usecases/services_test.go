package usecases

import (
	"strings"
	"testing"
)

type stubRepository struct{}

func (r *stubRepository) ReadFile(filePath string) (string, error) {
	return "", nil
}

func (r *stubRepository) WriteFile(filePath string, content string) error {
	return nil
}

func (r *stubRepository) CreateDir(dirPath string) error {
	return nil
}

func (r *stubRepository) ListMarkdownFiles(dirPath string) ([]string, error) {
	return nil, nil
}

func (r *stubRepository) RemoveFile(filePath string) error {
	return nil
}

func TestNewService_Normal(t *testing.T) {
	customRepo := &stubRepository{}
	service := NewService(customRepo)

	if service == nil {
		t.Fatal("service should not be nil")
	}
	if service.repository == nil {
		t.Fatal("repository should not be nil")
	}
	if service.repository != customRepo {
		t.Fatal("repository should use injected repository")
	}
	if service.splitHeadingsOperation == nil || service.addFrontMatterOperation == nil || service.addTagsOperation == nil {
		t.Fatal("operations should be initialized")
	}
	if service.deleteEmptyFilesOperation == nil || service.addHeading1Operation == nil || service.replaceImagesOperation == nil || service.removeHeadingAnnotationsOperation == nil {
		t.Fatal("operations should be initialized")
	}
	if service.removeTitleHashTagsOperation == nil {
		t.Fatal("operations should be initialized")
	}

	defaultService := NewService(nil)
	if defaultService == nil {
		t.Fatal("default service should not be nil")
	}
	if defaultService.repository == nil {
		t.Fatal("default repository should not be nil")
	}
	if defaultService.splitHeadingsOperation == nil || defaultService.addFrontMatterOperation == nil || defaultService.addTagsOperation == nil {
		t.Fatal("default operations should be initialized")
	}
	if defaultService.deleteEmptyFilesOperation == nil || defaultService.addHeading1Operation == nil || defaultService.replaceImagesOperation == nil || defaultService.removeHeadingAnnotationsOperation == nil {
		t.Fatal("default operations should be initialized")
	}
	if defaultService.removeTitleHashTagsOperation == nil {
		t.Fatal("default operations should be initialized")
	}
}

func TestSplitFrontMatterBlock_NoFrontMatter_Normal(t *testing.T) {
	t.Parallel()

	hasFrontMatter, block, body, err := splitFrontMatterBlock("hello")
	if err != nil {
		t.Fatalf("splitFrontMatterBlock returned error: %v", err)
	}
	if hasFrontMatter {
		t.Fatal("expected hasFrontMatter=false")
	}
	if block != "" {
		t.Fatalf("expected empty block, got %q", block)
	}
	if body != "hello" {
		t.Fatalf("unexpected body: %q", body)
	}
}

func TestSplitFrontMatterBlock_InvalidFormat(t *testing.T) {
	t.Parallel()

	_, _, _, err := splitFrontMatterBlock("---\ntitle: x\n")
	if err == nil {
		t.Fatal("expected error for missing front matter end delimiter")
	}
}

func TestParseFrontMatterMap_InvalidFormat(t *testing.T) {
	t.Parallel()

	_, _, err := parseFrontMatterMap("---\ninvalid-line\n---\n")
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestParseKVPairs_EmptyKey(t *testing.T) {
	t.Parallel()

	_, _, err := parseKVPairs([]string{"=value"})
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestBuildTagLine_Invalid(t *testing.T) {
	t.Parallel()

	_, err := buildTagLine(nil)
	if err == nil {
		t.Fatal("expected error for empty tags")
	}
}

func TestUniqueTrimmedTags_Normal(t *testing.T) {
	t.Parallel()

	tags := uniqueTrimmedTags(" go, #go, markdown ,")
	if len(tags) != 2 {
		t.Fatalf("unexpected tags length: %d", len(tags))
	}
	if strings.Join(tags, ",") != "go,markdown" {
		t.Fatalf("unexpected tags: %v", tags)
	}
}
