package config

import (
	"io"
	"os"
	"strings"
	"testing"
)

func withArgs(args []string, fn func()) {
	original := os.Args
	os.Args = args
	defer func() {
		os.Args = original
	}()
	fn()
}

func TestParseFlags_SplitHeadings_Normal(t *testing.T) {
	withArgs([]string{
		"markdown-crafter",
		"--operation", OperationSplitHeadings,
		"--file-path", "./sample.md",
		"--heading-level", "2",
		"--output-dir", "./out",
	}, func() {
		cfg, err := ParseFlags()
		if err != nil {
			t.Fatalf("ParseFlags returned error: %v", err)
		}
		if cfg.Operation != OperationSplitHeadings {
			t.Fatalf("unexpected operation: %s", cfg.Operation)
		}
	})
}

func TestParseFlags_AddFrontMatter_Normal(t *testing.T) {
	withArgs([]string{
		"markdown-crafter",
		"--operation", OperationAddFrontMatter,
		"--file-path", "./sample.md",
		"--kv", "title=doc",
		"--kv", "author=nov",
	}, func() {
		cfg, err := ParseFlags()
		if err != nil {
			t.Fatalf("ParseFlags returned error: %v", err)
		}
		if len(cfg.KVPairs) != 2 {
			t.Fatalf("unexpected kv pairs count: %d", len(cfg.KVPairs))
		}
	})
}

func TestParseFlags_AddTags_Normal(t *testing.T) {
	withArgs([]string{
		"markdown-crafter",
		"--operation", OperationAddTags,
		"--file-path", "./sample.md",
		"--tags", "go,markdown",
	}, func() {
		cfg, err := ParseFlags()
		if err != nil {
			t.Fatalf("ParseFlags returned error: %v", err)
		}
		if cfg.Tags != "go,markdown" {
			t.Fatalf("unexpected tags: %s", cfg.Tags)
		}
	})
}

func TestParseFlags_AddTagsByDir_Normal(t *testing.T) {
	withArgs([]string{
		"markdown-crafter",
		"--operation", OperationAddTags,
		"--dir-path", "./notes",
		"--tags", "go,markdown",
	}, func() {
		cfg, err := ParseFlags()
		if err != nil {
			t.Fatalf("ParseFlags returned error: %v", err)
		}
		if cfg.DirPath != "./notes" {
			t.Fatalf("unexpected dir path: %s", cfg.DirPath)
		}
	})
}

func TestParseFlags_DeleteEmptyFiles_Normal(t *testing.T) {
	withArgs([]string{
		"markdown-crafter",
		"--operation", OperationDeleteEmptyFiles,
		"--dir-path", "./notes",
	}, func() {
		cfg, err := ParseFlags()
		if err != nil {
			t.Fatalf("ParseFlags returned error: %v", err)
		}
		if cfg.DirPath != "./notes" {
			t.Fatalf("unexpected dir path: %s", cfg.DirPath)
		}
	})
}

func TestParseFlags_AddHeading1_Normal(t *testing.T) {
	withArgs([]string{
		"markdown-crafter",
		"--operation", OperationAddHeading1,
		"--file-path", "./sample.md",
		"--heading-text", "概要",
		"--heading-position", HeadingPositionTail,
	}, func() {
		cfg, err := ParseFlags()
		if err != nil {
			t.Fatalf("ParseFlags returned error: %v", err)
		}
		if cfg.HeadingText != "概要" {
			t.Fatalf("unexpected heading text: %s", cfg.HeadingText)
		}
		if cfg.HeadingPosition != HeadingPositionTail {
			t.Fatalf("unexpected heading position: %s", cfg.HeadingPosition)
		}
	})
}

func TestParseFlags_ReplaceImages_Normal(t *testing.T) {
	withArgs([]string{
		"markdown-crafter",
		"--operation", OperationReplaceImages,
		"--file-path", "./sample.md",
		"--replacement-text", "(添付画像)",
	}, func() {
		cfg, err := ParseFlags()
		if err != nil {
			t.Fatalf("ParseFlags returned error: %v", err)
		}
		if cfg.ReplacementText != "(添付画像)" {
			t.Fatalf("unexpected replacement text: %s", cfg.ReplacementText)
		}
	})
}

func TestParseFlags_RemoveHeadingAnnotations_Normal(t *testing.T) {
	withArgs([]string{
		"markdown-crafter",
		"--operation", OperationRemoveHeadingAnnotations,
		"--file-path", "./sample.md",
		"--heading-level", "3",
	}, func() {
		cfg, err := ParseFlags()
		if err != nil {
			t.Fatalf("ParseFlags returned error: %v", err)
		}
		if cfg.HeadingLevel != 3 {
			t.Fatalf("unexpected heading level: %d", cfg.HeadingLevel)
		}
	})
}

func TestParseFlags_RemoveTitleHashTags_Normal(t *testing.T) {
	withArgs([]string{
		"markdown-crafter",
		"--operation", OperationRemoveTitleHashTags,
		"--dir-path", "./notes",
	}, func() {
		cfg, err := ParseFlags()
		if err != nil {
			t.Fatalf("ParseFlags returned error: %v", err)
		}
		if cfg.DirPath != "./notes" {
			t.Fatalf("unexpected dir path: %s", cfg.DirPath)
		}
	})
}

func TestParseFlags_Help_Normal(t *testing.T) {
	withArgs([]string{"markdown-crafter", "--help"}, func() {
		cfg, err := ParseFlags()
		if err != nil {
			t.Fatalf("ParseFlags returned error: %v", err)
		}
		if !cfg.Help {
			t.Fatal("expected help=true")
		}
	})
}

func TestParseFlags_Invalid(t *testing.T) {
	withArgs([]string{
		"markdown-crafter",
		"--operation", OperationSplitHeadings,
		"--file-path", "./sample.md",
		"--heading-level", "0",
		"--output-dir", "./out",
	}, func() {
		_, err := ParseFlags()
		if err == nil {
			t.Fatal("expected validation error")
		}
	})
}

func TestParseFlags_InvalidAddTagsMissingTarget(t *testing.T) {
	withArgs([]string{
		"markdown-crafter",
		"--operation", OperationAddTags,
		"--tags", "go,markdown",
	}, func() {
		_, err := ParseFlags()
		if err == nil {
			t.Fatal("expected validation error")
		}
	})
}

func TestParseFlags_InvalidAddTagsWithBothTargets(t *testing.T) {
	withArgs([]string{
		"markdown-crafter",
		"--operation", OperationAddTags,
		"--file-path", "./sample.md",
		"--dir-path", "./notes",
		"--tags", "go,markdown",
	}, func() {
		_, err := ParseFlags()
		if err == nil {
			t.Fatal("expected validation error")
		}
	})
}

func TestParseFlags_InvalidDeleteEmptyFilesMissingDirPath(t *testing.T) {
	withArgs([]string{
		"markdown-crafter",
		"--operation", OperationDeleteEmptyFiles,
	}, func() {
		_, err := ParseFlags()
		if err == nil {
			t.Fatal("expected validation error")
		}
	})
}

func TestParseFlags_InvalidDeleteEmptyFilesLegacyDirectoryPath(t *testing.T) {
	withArgs([]string{
		"markdown-crafter",
		"--operation", OperationDeleteEmptyFiles,
		"--directory-path", "./notes",
	}, func() {
		_, err := ParseFlags()
		if err == nil {
			t.Fatal("expected parse error")
		}
	})
}

func TestParseFlags_InvalidReplaceImagesMissingReplacement(t *testing.T) {
	withArgs([]string{
		"markdown-crafter",
		"--operation", OperationReplaceImages,
		"--file-path", "./sample.md",
	}, func() {
		_, err := ParseFlags()
		if err == nil {
			t.Fatal("expected validation error")
		}
	})
}

func TestParseFlags_InvalidRemoveHeadingAnnotationsMissingHeadingLevel(t *testing.T) {
	withArgs([]string{
		"markdown-crafter",
		"--operation", OperationRemoveHeadingAnnotations,
		"--file-path", "./sample.md",
	}, func() {
		_, err := ParseFlags()
		if err == nil {
			t.Fatal("expected validation error")
		}
	})
}

func TestParseFlags_InvalidRemoveTitleHashTagsMissingDirPath(t *testing.T) {
	withArgs([]string{
		"markdown-crafter",
		"--operation", OperationRemoveTitleHashTags,
	}, func() {
		_, err := ParseFlags()
		if err == nil {
			t.Fatal("expected validation error")
		}
	})
}

func TestPrintUsage_Normal(t *testing.T) {
	withArgs([]string{"markdown-crafter"}, func() {
		oldStderr := os.Stderr
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("failed to create pipe: %v", err)
		}
		os.Stderr = w
		defer func() {
			os.Stderr = oldStderr
		}()

		PrintUsage()
		_ = w.Close()

		out, err := io.ReadAll(r)
		if err != nil {
			t.Fatalf("failed to read usage output: %v", err)
		}
		text := string(out)
		if !strings.Contains(text, "Markdown Crafter CLI ツール") {
			t.Fatalf("unexpected usage output: %s", text)
		}
	})
}
