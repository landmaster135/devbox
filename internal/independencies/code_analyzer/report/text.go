// internal/report/text.go
package report

import (
	"fmt"
	"sort"
	"strings"

	"github.com/landmaster135/devbox/internal/independencies/code_analyzer/models"
)

// TextReporter はテキスト形式のレポート生成機能を提供します
type TextReporter struct{}

// NewTextReporter は新しいTextReporterインスタンスを作成します
func NewTextReporter() *TextReporter {
	return &TextReporter{}
}

// Generate はテキスト形式のレポートを生成します
func (r *TextReporter) Generate(metrics models.ProjectMetrics) string {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("Project Analysis for: %s\n", metrics.ProjectPath))
	output.WriteString(fmt.Sprintf("Analyzed at: %s\n\n", metrics.AnalyzedAt.Format("2006-01-02 15:04:05")))
	output.WriteString(fmt.Sprintf("Files analyzed: %d\n", metrics.FileCount))
	output.WriteString(fmt.Sprintf("Total lines: %d\n", metrics.TotalLines))

	if metrics.TotalLines > 0 {
		output.WriteString(fmt.Sprintf("  - Code lines: %d (%.2f%%)\n",
			metrics.TotalCodeLines, float64(metrics.TotalCodeLines)/float64(metrics.TotalLines)*100))
		output.WriteString(fmt.Sprintf("  - Comment lines: %d (%.2f%%)\n",
			metrics.TotalComments, float64(metrics.TotalComments)/float64(metrics.TotalLines)*100))
		output.WriteString(fmt.Sprintf("  - Blank lines: %d (%.2f%%)\n",
			metrics.TotalBlankLines, float64(metrics.TotalBlankLines)/float64(metrics.TotalLines)*100))
	}

	output.WriteString(fmt.Sprintf("Average complexity: %.2f\n", metrics.AvgComplexity))
	output.WriteString(fmt.Sprintf("Maximum complexity: %d\n", metrics.MaxComplexity))
	output.WriteString(fmt.Sprintf("Comment-to-code ratio: %.2f%%\n\n", metrics.CommentRatio))

	// トレンド情報
	if metrics.Trends != nil {
		output.WriteString("Trends (compared to previous analysis):\n")
		for name, trend := range metrics.Trends {
			output.WriteString(fmt.Sprintf("  - %s: %.2f -> %.2f (%.2f%% change)\n",
				name, trend.Previous, trend.Current, trend.ChangeRate))
		}
		output.WriteString("\n")
	}

	// 複雑度でソート
	type complexityEntry struct {
		path       string
		complexity int
	}

	entries := make([]complexityEntry, 0, len(metrics.Files))
	for _, fm := range metrics.Files {
		entries = append(entries, complexityEntry{fm.Path, fm.Complexity})
	}

	// 複雑度順にソート
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].complexity > entries[j].complexity
	})

	// 複雑度の高いファイル
	output.WriteString("Files with highest complexity:\n")
	limit := 10
	if len(entries) < limit {
		limit = len(entries)
	}

	for i := 0; i < limit; i++ {
		output.WriteString(fmt.Sprintf("  %s: %d\n", entries[i].path, entries[i].complexity))
	}
	output.WriteString("\n")

	// クローン情報
	if len(metrics.Clones) > 0 {
		output.WriteString(fmt.Sprintf("Code Clones Detected: %d\n", len(metrics.Clones)))
		output.WriteString("Top code clones:\n")

		// 類似度でソート
		sort.Slice(metrics.Clones, func(i, j int) bool {
			return metrics.Clones[i].Similarity > metrics.Clones[j].Similarity
		})

		showLimit := 5
		if len(metrics.Clones) < showLimit {
			showLimit = len(metrics.Clones)
		}

		for i := 0; i < showLimit; i++ {
			c := metrics.Clones[i]
			output.WriteString(fmt.Sprintf("  - Source: %s (line %d)\n", c.SourceFile, c.SourceLine))
			output.WriteString(fmt.Sprintf("    Target: %s (line %d)\n", c.TargetFile, c.TargetLine))
			output.WriteString(fmt.Sprintf("    Size: %d lines, Similarity: %.2f%%\n", c.LineCount, c.Similarity*100))
			output.WriteString(fmt.Sprintf("    Preview: %s\n\n", c.Code))
		}
	}

	return output.String()
}
