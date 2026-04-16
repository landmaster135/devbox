package usecases

import (
	"errors"
	"os"
	"strings"
	"testing"

	filesystem "github.com/landmaster135/devbox/internal/coverage_badge/infrastructures/filesystem"
)

type TestCoverageBadgeService struct{}

func TestCoverageBadgeServiceCreateBadge_FromBadgeValue_Normal(t *testing.T) {
	service := NewCoverageBadgeServiceWithRepository(&filesystem.MockRepository{})

	badge, err := service.CreateBadge(CreateBadgeInput{
		BadgeTitle:      "Coverage",
		BadgeValue:      "58.6",
		GreenThreshold:  70,
		YellowThreshold: 30,
	})
	if err != nil {
		t.Fatalf("CreateBadge() error = %v", err)
	}

	expected := "![Coverage](https://img.shields.io/badge/Coverage-58.6%25-yellow)"
	if badge != expected {
		t.Fatalf("CreateBadge() = %q, want %q", badge, expected)
	}
}

func TestCoverageBadgeServiceCreateBadge_FromCoverageFile_Normal(t *testing.T) {
	repo := &filesystem.MockRepository{
		ReadFileFunc: func(path string) ([]byte, error) {
			if path != "coverage.out" {
				t.Fatalf("ReadFile path = %q, want coverage.out", path)
			}
			return []byte("total: (statements) 72.1%"), nil
		},
	}
	service := NewCoverageBadgeServiceWithRepository(repo)

	badge, err := service.CreateBadge(CreateBadgeInput{
		BadgeTitle:      "Coverage",
		CoverageFile:    "coverage.out",
		GreenThreshold:  70,
		YellowThreshold: 30,
	})
	if err != nil {
		t.Fatalf("CreateBadge() error = %v", err)
	}

	expected := "![Coverage](https://img.shields.io/badge/Coverage-72.1%25-green)"
	if badge != expected {
		t.Fatalf("CreateBadge() = %q, want %q", badge, expected)
	}
}

func TestCoverageBadgeServiceCreateBadge_ForceColorAndLink_Normal(t *testing.T) {
	service := NewCoverageBadgeServiceWithRepository(&filesystem.MockRepository{})

	badge, err := service.CreateBadge(CreateBadgeInput{
		BadgeTitle:      "Coverage",
		BadgeValue:      "95.0",
		GreenThreshold:  70,
		YellowThreshold: 30,
		ForceColor:      "red",
		BadgeLink:       "https://example.com/report",
	})
	if err != nil {
		t.Fatalf("CreateBadge() error = %v", err)
	}

	expected := "[![Coverage](https://img.shields.io/badge/Coverage-95.0%25-red)](https://example.com/report)"
	if badge != expected {
		t.Fatalf("CreateBadge() = %q, want %q", badge, expected)
	}
}

func TestCoverageBadgeServiceCreateBadge_InvalidCoverageReport(t *testing.T) {
	repo := &filesystem.MockRepository{
		ReadFileFunc: func(path string) ([]byte, error) {
			return []byte("mode: count\ninternal/file.go:1.1,2.1 1 1"), nil
		},
	}
	service := NewCoverageBadgeServiceWithRepository(repo)

	_, err := service.CreateBadge(CreateBadgeInput{
		BadgeTitle:      "Coverage",
		CoverageFile:    "coverage.out",
		GreenThreshold:  70,
		YellowThreshold: 30,
	})
	if err == nil {
		t.Fatal("CreateBadge() error = nil, want error")
	}
}

func TestCoverageBadgeServiceNewCoverageBadgeService_Normal(t *testing.T) {
	service := NewCoverageBadgeService()
	if service == nil {
		t.Fatal("NewCoverageBadgeService() = nil, want non-nil")
	}
}

func TestCoverageBadgeServiceNewCoverageBadgeServiceWithRepository_NilRepository(t *testing.T) {
	service := NewCoverageBadgeServiceWithRepository(nil)
	if service == nil {
		t.Fatal("NewCoverageBadgeServiceWithRepository(nil) = nil, want non-nil")
	}
}

func TestCoverageBadgeServicePatchBadge_ReplaceAndWrite_Normal(t *testing.T) {
	before := strings.Join([]string{
		"# devbox",
		"![Coverage](https://img.shields.io/badge/Coverage-10.0%25-red)",
		"本文",
	}, "\n")
	repo := &filesystem.MockRepository{
		ReadFileFunc: func(path string) ([]byte, error) {
			return []byte(before), nil
		},
	}
	service := NewCoverageBadgeServiceWithRepository(repo)

	result, err := service.PatchBadge(PatchBadgeInput{
		CreateBadgeInput: CreateBadgeInput{
			BadgeTitle:      "Coverage",
			BadgeValue:      "58.6",
			GreenThreshold:  70,
			YellowThreshold: 30,
		},
		TargetFile: "README.md",
		DryRun:     false,
	})
	if err != nil {
		t.Fatalf("PatchBadge() error = %v", err)
	}

	if !result.ContentModified {
		t.Fatal("ContentModified = false, want true")
	}
	if !result.FileWritten {
		t.Fatal("FileWritten = false, want true")
	}
	if repo.WriteCallCount != 1 {
		t.Fatalf("WriteCallCount = %d, want 1", repo.WriteCallCount)
	}
	if !strings.Contains(result.PatchedContent, "58.6%25-yellow") {
		t.Fatalf("patched content does not include updated badge: %q", result.PatchedContent)
	}
}

func TestCoverageBadgeServicePatchBadge_InsertAfterHeading_Normal(t *testing.T) {
	before := strings.Join([]string{
		"# devbox",
		"本文",
	}, "\n")
	repo := &filesystem.MockRepository{
		ReadFileFunc: func(path string) ([]byte, error) {
			return []byte(before), nil
		},
	}
	service := NewCoverageBadgeServiceWithRepository(repo)

	result, err := service.PatchBadge(PatchBadgeInput{
		CreateBadgeInput: CreateBadgeInput{
			BadgeTitle:      "Coverage",
			BadgeValue:      "72.1",
			GreenThreshold:  70,
			YellowThreshold: 30,
		},
		TargetFile: "README.md",
	})
	if err != nil {
		t.Fatalf("PatchBadge() error = %v", err)
	}

	lines := strings.Split(result.PatchedContent, "\n")
	if len(lines) < 2 {
		t.Fatalf("patched line count = %d, want >= 2", len(lines))
	}
	if !strings.Contains(lines[1], "img.shields.io/badge/Coverage-72.1%25-green") {
		t.Fatalf("line[1] = %q, want coverage badge", lines[1])
	}
}

func TestCoverageBadgeServicePatchBadge_InsertAtTopWhenNoHeading_Normal(t *testing.T) {
	before := "本文のみ"
	repo := &filesystem.MockRepository{
		ReadFileFunc: func(path string) ([]byte, error) {
			return []byte(before), nil
		},
	}
	service := NewCoverageBadgeServiceWithRepository(repo)

	result, err := service.PatchBadge(PatchBadgeInput{
		CreateBadgeInput: CreateBadgeInput{
			BadgeTitle:      "Coverage",
			BadgeValue:      "61.0",
			GreenThreshold:  70,
			YellowThreshold: 30,
		},
		TargetFile: "README.md",
	})
	if err != nil {
		t.Fatalf("PatchBadge() error = %v", err)
	}

	if !strings.HasPrefix(result.PatchedContent, "![Coverage](https://img.shields.io/badge/Coverage-61.0%25-yellow)") {
		t.Fatalf("patched content = %q, want badge inserted at top", result.PatchedContent)
	}
}

func TestCoverageBadgeServicePatchBadge_DryRunDoesNotWrite_Normal(t *testing.T) {
	before := "# devbox\n本文\n"
	repo := &filesystem.MockRepository{
		ReadFileFunc: func(path string) ([]byte, error) {
			return []byte(before), nil
		},
	}
	service := NewCoverageBadgeServiceWithRepository(repo)

	result, err := service.PatchBadge(PatchBadgeInput{
		CreateBadgeInput: CreateBadgeInput{
			BadgeTitle:      "Coverage",
			BadgeValue:      "72.1",
			GreenThreshold:  70,
			YellowThreshold: 30,
		},
		TargetFile: "README.md",
		DryRun:     true,
	})
	if err != nil {
		t.Fatalf("PatchBadge() error = %v", err)
	}
	if result.FileWritten {
		t.Fatal("FileWritten = true, want false")
	}
	if repo.WriteCallCount != 0 {
		t.Fatalf("WriteCallCount = %d, want 0", repo.WriteCallCount)
	}
}

func TestCoverageBadgeServicePatchBadge_RemoveDuplicateCoverageBadges_Normal(t *testing.T) {
	before := strings.Join([]string{
		"# devbox",
		"![Coverage](https://img.shields.io/badge/Coverage-10.0%25-red)",
		"![coverage](https://img.shields.io/badge/coverage-11.0%25-red)",
		"本文",
	}, "\n")
	repo := &filesystem.MockRepository{
		ReadFileFunc: func(path string) ([]byte, error) {
			return []byte(before), nil
		},
	}
	service := NewCoverageBadgeServiceWithRepository(repo)

	result, err := service.PatchBadge(PatchBadgeInput{
		CreateBadgeInput: CreateBadgeInput{
			BadgeTitle:      "Coverage",
			BadgeValue:      "58.6",
			GreenThreshold:  70,
			YellowThreshold: 30,
		},
		TargetFile: "README.md",
	})
	if err != nil {
		t.Fatalf("PatchBadge() error = %v", err)
	}

	if strings.Count(result.PatchedContent, "img.shields.io/badge/") != 1 {
		t.Fatalf("coverage badge count = %d, want 1", strings.Count(result.PatchedContent, "img.shields.io/badge/"))
	}
}

func TestCoverageBadgeServicePatchBadge_NoChangeDoesNotWrite_Normal(t *testing.T) {
	before := strings.Join([]string{
		"# devbox",
		"![Coverage](https://img.shields.io/badge/Coverage-58.6%25-yellow)",
		"本文",
	}, "\n")
	repo := &filesystem.MockRepository{
		ReadFileFunc: func(path string) ([]byte, error) {
			return []byte(before), nil
		},
	}
	service := NewCoverageBadgeServiceWithRepository(repo)

	result, err := service.PatchBadge(PatchBadgeInput{
		CreateBadgeInput: CreateBadgeInput{
			BadgeTitle:      "Coverage",
			BadgeValue:      "58.6",
			GreenThreshold:  70,
			YellowThreshold: 30,
		},
		TargetFile: "README.md",
	})
	if err != nil {
		t.Fatalf("PatchBadge() error = %v", err)
	}
	if result.ContentModified {
		t.Fatal("ContentModified = true, want false")
	}
	if result.FileWritten {
		t.Fatal("FileWritten = true, want false")
	}
	if repo.WriteCallCount != 0 {
		t.Fatalf("WriteCallCount = %d, want 0", repo.WriteCallCount)
	}
}

func TestCoverageBadgeServicePatchBadge_PreserveSingleTrailingNewline_Normal(t *testing.T) {
	before := strings.Join([]string{
		"# devbox",
		"![Coverage](https://img.shields.io/badge/Coverage-10.0%25-red)",
		"本文",
	}, "\n") + "\n"
	repo := &filesystem.MockRepository{
		ReadFileFunc: func(path string) ([]byte, error) {
			return []byte(before), nil
		},
	}
	service := NewCoverageBadgeServiceWithRepository(repo)

	result, err := service.PatchBadge(PatchBadgeInput{
		CreateBadgeInput: CreateBadgeInput{
			BadgeTitle:      "Coverage",
			BadgeValue:      "58.6",
			GreenThreshold:  70,
			YellowThreshold: 30,
		},
		TargetFile: "README.md",
	})
	if err != nil {
		t.Fatalf("PatchBadge() error = %v", err)
	}

	if strings.Count(result.PatchedContent, "\n") != strings.Count(before, "\n") {
		t.Fatalf("newline count = %d, want %d", strings.Count(result.PatchedContent, "\n"), strings.Count(before, "\n"))
	}
	if !strings.HasSuffix(result.PatchedContent, "\n") {
		t.Fatal("patched content must keep trailing newline")
	}
}

func TestCoverageBadgeServicePatchBadge_PreserveNoTrailingNewline_Normal(t *testing.T) {
	before := strings.Join([]string{
		"# devbox",
		"![Coverage](https://img.shields.io/badge/Coverage-10.0%25-red)",
		"本文",
	}, "\n")
	repo := &filesystem.MockRepository{
		ReadFileFunc: func(path string) ([]byte, error) {
			return []byte(before), nil
		},
	}
	service := NewCoverageBadgeServiceWithRepository(repo)

	result, err := service.PatchBadge(PatchBadgeInput{
		CreateBadgeInput: CreateBadgeInput{
			BadgeTitle:      "Coverage",
			BadgeValue:      "58.6",
			GreenThreshold:  70,
			YellowThreshold: 30,
		},
		TargetFile: "README.md",
	})
	if err != nil {
		t.Fatalf("PatchBadge() error = %v", err)
	}

	if strings.HasSuffix(result.PatchedContent, "\n") {
		t.Fatal("patched content must not add trailing newline")
	}
}

func TestCoverageBadgeServicePatchBadge_WriteError(t *testing.T) {
	repo := &filesystem.MockRepository{
		ReadFileFunc: func(path string) ([]byte, error) {
			return []byte("# devbox\n"), nil
		},
		WriteFileFunc: func(path string, data []byte, perm os.FileMode) error {
			return errors.New("write failed")
		},
	}
	service := NewCoverageBadgeServiceWithRepository(repo)

	_, err := service.PatchBadge(PatchBadgeInput{
		CreateBadgeInput: CreateBadgeInput{
			BadgeTitle:      "Coverage",
			BadgeValue:      "58.6",
			GreenThreshold:  70,
			YellowThreshold: 30,
		},
		TargetFile: "README.md",
	})
	if err == nil {
		t.Fatal("PatchBadge() error = nil, want error")
	}
}

func TestParseCoverageValue_Normal(t *testing.T) {
	got, err := parseCoverageValue("58.6%")
	if err != nil {
		t.Fatalf("parseCoverageValue() error = %v", err)
	}
	if got != 58.6 {
		t.Fatalf("parseCoverageValue() = %v, want 58.6", got)
	}
}

func TestParseCoverageValue_InvalidRange(t *testing.T) {
	if _, err := parseCoverageValue("120"); err == nil {
		t.Fatal("parseCoverageValue() error = nil, want error")
	}
}

func TestParseCoverageValue_InvalidFormat(t *testing.T) {
	if _, err := parseCoverageValue("abc"); err == nil {
		t.Fatal("parseCoverageValue() error = nil, want error")
	}
}

func TestParseCoverageFromReport_Normal(t *testing.T) {
	content := "pkg/file.go:10:\tFunc\t100.0%\ntotal: (statements) 61.9%"
	got, err := parseCoverageFromReport(content)
	if err != nil {
		t.Fatalf("parseCoverageFromReport() error = %v", err)
	}
	if got != 61.9 {
		t.Fatalf("parseCoverageFromReport() = %v, want 61.9", got)
	}
}

func TestParseCoverageFromReport_InvalidTotalValue(t *testing.T) {
	content := "total: (statements) abc%"
	if _, err := parseCoverageFromReport(content); err == nil {
		t.Fatal("parseCoverageFromReport() error = nil, want error")
	}
}
