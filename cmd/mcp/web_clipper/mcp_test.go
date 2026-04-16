package web_clipper

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	mcpInfraUsecases "github.com/landmaster135/devbox/internal/mcp_infra_mark3labs/usecases"
)

func createMockCallToolRequest(arguments map[string]interface{}) mcpInfraUsecases.CallToolRequest {
	return mcpInfraUsecases.CallToolRequest{
		Request: mcpInfraUsecases.Request{
			Method: "tools/call",
		},
		Params: mcpInfraUsecases.CallToolParams{
			Name:      "patch_markdown",
			Arguments: arguments,
		},
	}
}

func TestCreateConfigFromRequest_Normal(t *testing.T) {
	t.Parallel()

	outFilePath := filepath.Join(t.TempDir(), "out.md")
	request := createMockCallToolRequest(map[string]interface{}{
		"target_title":         "OpenAI Blog",
		"target_url":           "https://openai.com/blog",
		"src_markdown_content": "## 記事タイトル 要約\n\n### 見出し\n本文\n",
		"out_file_path":        outFilePath,
		"top_heading_level":    2,
	})

	cfg, err := createConfigFromRequest(request)
	if err != nil {
		t.Fatalf("createConfigFromRequest returned error: %v", err)
	}

	if cfg.Operation != "patch-markdown" {
		t.Fatalf("unexpected operation: got=%q want=%q", cfg.Operation, "patch-markdown")
	}
	if cfg.TargetTitle != "OpenAI Blog" {
		t.Fatalf("unexpected target title: got=%q", cfg.TargetTitle)
	}
	if cfg.TargetURL != "https://openai.com/blog" {
		t.Fatalf("unexpected target url: got=%q", cfg.TargetURL)
	}
	if cfg.OutFilePath != outFilePath {
		t.Fatalf("unexpected out file path: got=%q want=%q", cfg.OutFilePath, outFilePath)
	}
	if cfg.TopHeadingLevel != 2 {
		t.Fatalf("unexpected top heading level: got=%d want=%d", cfg.TopHeadingLevel, 2)
	}
}

func TestCreateConfigFromRequest_MissingRequiredField(t *testing.T) {
	t.Parallel()

	request := createMockCallToolRequest(map[string]interface{}{
		"target_url":        "https://openai.com/blog",
		"out_file_path":     "/tmp/out.md",
		"top_heading_level": 2,
	})

	_, err := createConfigFromRequest(request)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if !strings.Contains(err.Error(), `required argument "target_title" not found`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateConfigFromRequest_DualInputSources(t *testing.T) {
	t.Parallel()

	request := createMockCallToolRequest(map[string]interface{}{
		"target_title":         "OpenAI Blog",
		"target_url":           "https://openai.com/blog",
		"src_markdown_content": "## 見出し\n本文\n",
		"src_markdown_file":    "/tmp/in.md",
		"out_file_path":        "/tmp/out.md",
		"top_heading_level":    2,
	})

	_, err := createConfigFromRequest(request)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if !strings.Contains(err.Error(), "--src-markdown-content と --src-markdown-file は同時に指定できません") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandlePatchMarkdown_Normal(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	outFilePath := filepath.Join(tempDir, "out.md")
	request := createMockCallToolRequest(map[string]interface{}{
		"target_title":         "OpenAI Blog",
		"target_url":           "https://openai.com/blog",
		"src_markdown_content": "## 記事タイトル 要約\n\n### 見出し\n本文\n",
		"out_file_path":        outFilePath,
		"top_heading_level":    2,
	})

	result, err := handlePatchMarkdown(context.Background(), request)
	if err != nil {
		t.Fatalf("handlePatchMarkdown returned error: %v", err)
	}
	if result == nil {
		t.Fatal("result is nil")
	}
	if len(result.Content) == 0 {
		t.Fatal("result content is empty")
	}

	textContent, ok := result.Content[0].(mcpInfraUsecases.TextContent)
	if !ok {
		t.Fatalf("unexpected content type: %T", result.Content[0])
	}
	if textContent.Text != "出力しました: "+outFilePath {
		t.Fatalf("unexpected result text: got=%q", textContent.Text)
	}

	output, err := os.ReadFile(outFilePath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if !strings.Contains(string(output), "- [OpenAI Blog](https://openai.com/blog)") {
		t.Fatalf("output does not contain inserted link: %s", string(output))
	}
}

func TestHandlePatchMarkdown_Error(t *testing.T) {
	t.Parallel()

	outFilePath := filepath.Join(t.TempDir(), "out.md")
	request := createMockCallToolRequest(map[string]interface{}{
		"target_title":         "OpenAI Blog",
		"target_url":           "https://openai.com/blog",
		"src_markdown_content": "### 見出し\n本文\n",
		"out_file_path":        outFilePath,
		"top_heading_level":    2,
	})

	_, err := handlePatchMarkdown(context.Background(), request)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if !strings.Contains(err.Error(), "MarkdownへのWeb記事リンク挿入に失敗しました") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddPromptIntoServer_Normal(t *testing.T) {
	t.Parallel()

	s := mcpInfra.NewMCPServer(
		"Test Web Clipper",
		"1.0.0",
		mcpInfra.WithResourceCapabilities(true, true),
		mcpInfra.WithPromptCapabilities(true),
	)

	result := addPromptIntoServer(s)
	if result == nil {
		t.Fatal("expected server but got nil")
	}
}

func TestCreateWebClipperServer_Normal(t *testing.T) {
	t.Parallel()

	s := createWebClipperServer()
	if s == nil {
		t.Fatal("expected server but got nil")
	}
}
