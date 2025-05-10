// internal/report/csv.go
package report

import (
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/independencies/code_analyzer/models"
)

// CSVReporter はCSV形式のレポート生成機能を提供します
type CSVReporter struct{}

// NewCSVReporter は新しいCSVReporterインスタンスを作成します
func NewCSVReporter() *CSVReporter {
	return &CSVReporter{}
}

// Generate はCSV形式のレポートを生成します
func (r *CSVReporter) Generate(metrics models.ProjectMetrics) string {
	var output strings.Builder

	// CSVヘッダー
	output.WriteString("Path,TotalLines,CodeLines,CommentLines,BlankLines,FunctionCount,AvgFunctionSize,MaxFunctionSize,Complexity,CommentRatio\n")

	// ファイルごとのデータ
	for _, fm := range metrics.Files {
		output.WriteString(fmt.Sprintf("%s,%d,%d,%d,%d,%d,%.2f,%d,%d,%.2f\n",
			fm.Path, fm.TotalLines, fm.CodeLines, fm.CommentLines, fm.BlankLines,
			fm.FunctionCount, fm.AvgFunctionSize, fm.MaxFunctionSize, fm.Complexity, fm.CommentRatio))
	}

	// プロジェクト全体のサマリー
	output.WriteString(fmt.Sprintf("\nSUMMARY,%d files,%d lines,%d code,%d comments,%d blank,%.2f avg complexity,%d max complexity,%.2f%% comment ratio\n",
		metrics.FileCount, metrics.TotalLines, metrics.TotalCodeLines, metrics.TotalComments,
		metrics.TotalBlankLines, metrics.AvgComplexity, metrics.MaxComplexity, metrics.CommentRatio))

	// クローン情報（検出している場合）
	if len(metrics.Clones) > 0 {
		output.WriteString("\nCode Clones\n")
		output.WriteString("SourceFile,TargetFile,SourceLine,TargetLine,LineCount,Similarity\n")
		for _, clone := range metrics.Clones {
			output.WriteString(fmt.Sprintf("%s,%s,%d,%d,%d,%.2f\n",
				clone.SourceFile, clone.TargetFile, clone.SourceLine, clone.TargetLine,
				clone.LineCount, clone.Similarity))
		}
	}

	return output.String()
}
