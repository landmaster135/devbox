package usecases

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	filesystem "github.com/landmaster135/devbox/internal/service_implementing_viewer/infrastructures/filesystem"
)

func TestBuildAggregatedSummary_Normal(t *testing.T) {
	rootDir := filepath.Join("tmp", "project", "cmd")
	cliPath := filepath.Join(rootDir, "cli")
	mcpPath := filepath.Join(rootDir, "mcp")
	grpcPath := filepath.Join(rootDir, "grpc", "handlers")
	httpPath := filepath.Join(rootDir, "http", "handlers")

	readmeByPath := map[string]string{
		filepath.Join(cliPath, "anilist", "README.md"): `# AniList CLI ツール

AniListからアニメ・マンガ情報を取得するためのコマンドラインツールです。

## 概要
この文章は抽出対象外`,
		filepath.Join(cliPath, "no-summary-tool", "README.md"): `# no-summary-tool

## 概要
本文なし扱い`,
		filepath.Join(mcpPath, "arithmetic-calculator", "README.md"): `# Arithmetic Calculator

A simple arithmetic calculator.

## Features
- add`,
		filepath.Join(grpcPath, "cron_workflow", "README.md"): `# Cron Workflow

Serves GUI for CRON workflow.

## Build`,
		filepath.Join(httpPath, "cron_workflow", "README.md"): `# Cron Workflow HTTP Handler

Serves GUI for CRON workflow.

## Prerequisites`,
	}

	mockRepo := &filesystem.MockRepository{
		ListDirectoriesFunc: func(path string) ([]string, error) {
			switch path {
			case cliPath:
				return []string{"anilist", "missing-readme", "no-summary-tool"}, nil
			case mcpPath:
				return []string{"arithmetic-calculator"}, nil
			case grpcPath:
				return []string{"cron_workflow"}, nil
			case httpPath:
				return []string{"cron_workflow"}, nil
			default:
				return []string{}, nil
			}
		},
		ReadFileFunc: func(path string) ([]byte, error) {
			if content, exists := readmeByPath[path]; exists {
				return []byte(content), nil
			}
			return nil, os.ErrNotExist
		},
	}

	service := NewServiceImplementingViewerServiceWithDependencies(rootDir, []string{"cli", "mcp", "grpc/handlers", "http/handlers", "powershell"}, mockRepo)

	got, err := service.BuildAggregatedSummary()
	if err != nil {
		t.Fatalf("BuildAggregatedSummary() error = %v", err)
	}

	want := strings.Join([]string{
		"## CLI tools",
		"*anilist*: AniListからアニメ・マンガ情報を取得するためのコマンドラインツールです。",
		"*missing-readme*: ",
		"*no-summary-tool*: ",
		"",
		"## MCP tools",
		"*arithmetic-calculator*: A simple arithmetic calculator.",
		"",
		"## GRPC/HANDLERS tools",
		"*cron_workflow*: Serves GUI for CRON workflow.",
		"",
		"## HTTP/HANDLERS tools",
		"*cron_workflow*: Serves GUI for CRON workflow.",
		"",
	}, "\n")
	if got != want {
		t.Fatalf("BuildAggregatedSummary() mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestBuildAggregatedSummary_ListDirectoriesError(t *testing.T) {
	rootDir := "/workspace/cmd"
	cliPath := filepath.Join(rootDir, "cli")
	mockRepo := &filesystem.MockRepository{
		ListDirectoriesFunc: func(path string) ([]string, error) {
			if path == cliPath {
				return nil, errors.New("list failed")
			}
			return []string{}, nil
		},
	}

	service := NewServiceImplementingViewerServiceWithDependencies(rootDir, []string{"cli"}, mockRepo)

	_, err := service.BuildAggregatedSummary()
	if err == nil {
		t.Fatal("エラーが発生しませんでした")
	}
	if !strings.Contains(err.Error(), "ディレクトリ") {
		t.Fatalf("エラーメッセージが期待値と異なります: %v", err)
	}
}

func TestBuildAggregatedSummary_ReadFileError(t *testing.T) {
	rootDir := "/workspace/cmd"
	cliPath := filepath.Join(rootDir, "cli")
	readmePath := filepath.Join(cliPath, "anilist", "README.md")
	mockRepo := &filesystem.MockRepository{
		ListDirectoriesFunc: func(path string) ([]string, error) {
			return []string{"anilist"}, nil
		},
		ReadFileFunc: func(path string) ([]byte, error) {
			if path == readmePath {
				return nil, errors.New("permission denied")
			}
			return nil, os.ErrNotExist
		},
	}

	service := NewServiceImplementingViewerServiceWithDependencies(rootDir, []string{"cli"}, mockRepo)

	_, err := service.BuildAggregatedSummary()
	if err == nil {
		t.Fatal("エラーが発生しませんでした")
	}
	if !strings.Contains(err.Error(), "README.md の読み取りに失敗しました") {
		t.Fatalf("エラーメッセージが期待値と異なります: %v", err)
	}
}

func TestAggregateSummaryToFile_Normal(t *testing.T) {
	rootDir := "/workspace/cmd"
	cliPath := filepath.Join(rootDir, "cli")
	mockRepo := &filesystem.MockRepository{
		ListDirectoriesFunc: func(path string) ([]string, error) {
			if path == cliPath {
				return []string{"anilist"}, nil
			}
			return []string{}, nil
		},
		ReadFileFunc: func(path string) ([]byte, error) {
			return []byte("# Title\n\nsummary line\n\n## details"), nil
		},
	}

	service := NewServiceImplementingViewerServiceWithDependencies(rootDir, []string{"cli"}, mockRepo)
	err := service.AggregateSummaryToFile("docs/service_summary.md")
	if err != nil {
		t.Fatalf("AggregateSummaryToFile() error = %v", err)
	}

	if mockRepo.LastWritePath != "docs/service_summary.md" {
		t.Fatalf("LastWritePath = %s, want %s", mockRepo.LastWritePath, "docs/service_summary.md")
	}
	if mockRepo.LastWritePermission != 0o644 {
		t.Fatalf("LastWritePermission = %o, want %o", mockRepo.LastWritePermission, 0o644)
	}
	if !strings.Contains(string(mockRepo.LastWriteContent), "*anilist*: summary line") {
		t.Fatalf("LastWriteContent is unexpected: %s", string(mockRepo.LastWriteContent))
	}
}

func TestAggregateSummaryToFile_WriteError(t *testing.T) {
	rootDir := "/workspace/cmd"
	mockRepo := &filesystem.MockRepository{
		WriteFileFunc: func(path string, data []byte, perm os.FileMode) error {
			return errors.New("write failed")
		},
	}

	service := NewServiceImplementingViewerServiceWithDependencies(rootDir, []string{}, mockRepo)
	err := service.AggregateSummaryToFile("docs/service_summary.md")
	if err == nil {
		t.Fatal("エラーが発生しませんでした")
	}
	if err.Error() != "write failed" {
		t.Fatalf("エラーメッセージが期待値と異なります: %v", err)
	}
}

func TestExtractSummaryLine(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "先頭の本文を抽出できる",
			input: "# title\n\nfirst summary\n\n## details\nbody",
			want:  "first summary",
		},
		{
			name:  "README本文が見つからない場合は空文字",
			input: "# title\n\n## details\nbody",
			want:  "",
		},
		{
			name:  "コードブロックを無視できる",
			input: "# title\n\n```md\nsummary in code\n```\n\nactual summary\n\n## details",
			want:  "actual summary",
		},
		{
			name:  "コードブロック内の見出しは判定対象外",
			input: "# title\n\n```md\n## fake heading\n```\n\nactual summary\n\n## details",
			want:  "actual summary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractSummaryLine(tt.input)
			if got != tt.want {
				t.Fatalf("extractSummaryLine() = %q, want %q", got, tt.want)
			}
		})
	}
}
