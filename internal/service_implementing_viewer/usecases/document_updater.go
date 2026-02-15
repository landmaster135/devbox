package usecases

import (
	"fmt"
	"strings"

	filesystem "github.com/landmaster135/devbox/internal/service_implementing_viewer/infrastructures/filesystem"
)

const (
	implementationHeading = "### 実装状況一覧"
	statisticsHeading     = "### 統計情報"
)

// UpdateDocumentationFile はMarkdownドキュメント内の対象セクションを書き換える
func UpdateDocumentationFile(filePath, table string, stats *ServiceStatistics) error {
	return updateDocumentationFileWithRepository(filesystem.NewRepository(), filePath, table, stats)
}

func updateDocumentationFileWithRepository(fileSystem filesystem.Repository, filePath, table string, stats *ServiceStatistics) error {
	if fileSystem == nil {
		fileSystem = filesystem.NewRepository()
	}

	content, err := fileSystem.ReadFile(filePath)
	if err != nil {
		return err
	}

	updated, err := buildUpdatedDocument(string(content), table, stats)
	if err != nil {
		return err
	}

	if err := fileSystem.WriteFile(filePath, []byte(updated), 0o644); err != nil {
		return err
	}
	return nil
}

func buildUpdatedDocument(original, table string, stats *ServiceStatistics) (string, error) {
	lines := strings.Split(original, "\n")

	var err error
	lines, err = replaceSectionWithContent(lines, implementationHeading, buildTableSection(table))
	if err != nil {
		return "", err
	}

	lines, err = replaceSectionWithContent(lines, statisticsHeading, buildStatisticsSection(stats))
	if err != nil {
		return "", err
	}

	return strings.Join(lines, "\n"), nil
}

func replaceSectionWithContent(lines []string, heading string, section []string) ([]string, error) {
	idx := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == heading {
			idx = i
			break
		}
	}
	if idx == -1 {
		return nil, fmt.Errorf("見出しが見つかりません: %s", heading)
	}

	end := len(lines)
	for i := idx + 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "### ") || strings.HasPrefix(trimmed, "## ") {
			end = i
			break
		}
	}

	newLines := append([]string{}, lines[:idx+1]...)
	newLines = append(newLines, section...)
	newLines = append(newLines, lines[end:]...)
	return newLines, nil
}

func buildTableSection(table string) []string {
	trimmed := strings.TrimSpace(table)
	if trimmed == "" {
		return []string{""}
	}

	rows := strings.Split(trimmed, "\n")
	section := make([]string, 0, len(rows)+2)
	section = append(section, "")
	section = append(section, rows...)
	section = append(section, "")
	return section
}

func buildStatisticsSection(stats *ServiceStatistics) []string {
	if stats == nil {
		stats = &ServiceStatistics{}
	}

	return []string{
		"",
		fmt.Sprintf("- **総サービス数**: %d", stats.TotalServices),
		fmt.Sprintf("- **CLIツール実装数**: %d", stats.CLICount),
		fmt.Sprintf("- **MCPツール実装数**: %d", stats.MCPCount),
		fmt.Sprintf("- **gRPCハンドラ実装数**: %d", stats.GRPCCount),
		fmt.Sprintf("- **HTTPハンドラ実装数**: %d", stats.HTTPCount),
		fmt.Sprintf("- **CLIのみ実装**: %d", stats.CLIOnlyCount),
		fmt.Sprintf("- **MCPのみ実装**: %d", stats.MCPOnlyCount),
		fmt.Sprintf("- **gRPCハンドラのみ実装**: %d", stats.GRPCOnlyCount),
		fmt.Sprintf("- **HTTPハンドラのみ実装**: %d", stats.HTTPOnlyCount),
		fmt.Sprintf("- **CLI+MCP両方実装**: %d", stats.BothCLIMCPCount),
		fmt.Sprintf("- **全て実装済み**: %d", stats.AllImplementedCount),
		"",
	}
}
