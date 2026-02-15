package usecases

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	summarySectionCLI  = "cli"
	summarySectionMCP  = "mcp"
	summarySectionGRPC = "grpc_handlers"
	summarySectionHTTP = "http_handlers"
)

var summarySectionOrder = []string{
	summarySectionCLI,
	summarySectionMCP,
	summarySectionGRPC,
	summarySectionHTTP,
}

var summarySectionTitleByKey = map[string]string{
	summarySectionCLI:  "## CLI tools",
	summarySectionMCP:  "## MCP tools",
	summarySectionGRPC: "## GRPC/HANDLERS tools",
	summarySectionHTTP: "## HTTP/HANDLERS tools",
}

// AggregateSummaryToFile は target-dirs 配下の README.md から概要を集約し、指定ファイルへ書き込む。
func (s *ServiceImplementingViewerService) AggregateSummaryToFile(writeFilePath string) error {
	summary, err := s.BuildAggregatedSummary()
	if err != nil {
		return err
	}

	if err := s.fileSystem.WriteFile(writeFilePath, []byte(summary), 0o644); err != nil {
		return err
	}

	return nil
}

// BuildAggregatedSummary は target-dirs 配下の README.md から概要を集約し、Markdown文字列を返す。
func (s *ServiceImplementingViewerService) BuildAggregatedSummary() (string, error) {
	summaries := map[string]map[string]string{
		summarySectionCLI:  {},
		summarySectionMCP:  {},
		summarySectionGRPC: {},
		summarySectionHTTP: {},
	}

	for _, targetDir := range s.targetDirs {
		sectionKey, ok := classifySummarySection(targetDir)
		if !ok {
			continue
		}

		dirPath := s.fileSystem.Join(s.rootDir, targetDir)
		tools, err := s.getServicesInDirectory(dirPath)
		if err != nil {
			return "", fmt.Errorf("ディレクトリ %s の読み取りに失敗しました: %w", dirPath, err)
		}

		for _, toolName := range tools {
			summary, err := s.readSummaryLine(dirPath, toolName)
			if err != nil {
				return "", err
			}
			summaries[sectionKey][toolName] = summary
		}
	}

	return formatAggregatedSummary(summaries), nil
}

func (s *ServiceImplementingViewerService) readSummaryLine(targetDirPath, toolName string) (string, error) {
	readmePath := s.fileSystem.Join(targetDirPath, toolName, "README.md")
	content, err := s.fileSystem.ReadFile(readmePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", fmt.Errorf("README.md の読み取りに失敗しました (%s): %w", readmePath, err)
	}

	return extractSummaryLine(string(content)), nil
}

func extractSummaryLine(readmeContent string) string {
	lines := strings.Split(readmeContent, "\n")
	inCodeBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			continue
		}
		if strings.HasPrefix(trimmed, "##") {
			break
		}
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			continue
		}

		return trimmed
	}

	return ""
}

func classifySummarySection(targetDir string) (string, bool) {
	normalized := strings.TrimSpace(strings.ReplaceAll(targetDir, "\\", "/"))

	switch {
	case strings.Contains(normalized, "grpc/handlers"):
		return summarySectionGRPC, true
	case strings.Contains(normalized, "http/handlers"):
		return summarySectionHTTP, true
	case strings.Contains(normalized, "mcp"):
		return summarySectionMCP, true
	case strings.Contains(normalized, "cli"):
		return summarySectionCLI, true
	default:
		return "", false
	}
}

func formatAggregatedSummary(summaries map[string]map[string]string) string {
	var builder strings.Builder

	for idx, sectionKey := range summarySectionOrder {
		builder.WriteString(summarySectionTitleByKey[sectionKey])
		builder.WriteString("\n")

		toolNames := make([]string, 0, len(summaries[sectionKey]))
		for toolName := range summaries[sectionKey] {
			toolNames = append(toolNames, toolName)
		}
		sort.Strings(toolNames)

		for _, toolName := range toolNames {
			builder.WriteString("*")
			builder.WriteString(toolName)
			builder.WriteString("*: ")
			builder.WriteString(summaries[sectionKey][toolName])
			builder.WriteString("\n")
		}

		if idx < len(summarySectionOrder)-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}
