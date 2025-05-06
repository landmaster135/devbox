package presenter

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/landmaster135/devbox/internal/independencies/depends_visualizer/analyzer"
)

// RenderDOT は依存関係を DOT 形式で出力します
func RenderDOT(results []analyzer.AnalysisResult) (string, error) {
	var sb strings.Builder

	sb.WriteString("digraph G {\n")
	sb.WriteString("  rankdir=BT;\n")  // 下から上への方向
	sb.WriteString("  node [shape=box, style=filled, fillcolor=lightblue];\n")

	fileCount := 0
	for _, result := range results {
		shortName := filepath.Base(result.FilePath)
		sb.WriteString(fmt.Sprintf("  // File: %s\n", shortName))
		sb.WriteString(fmt.Sprintf("  subgraph cluster_%d {\n", fileCount))
		sb.WriteString(fmt.Sprintf("    label=\"%s\";\n", shortName))

		// ノードの定義
		for funcName := range result.Dependencies {
			sb.WriteString(fmt.Sprintf("    \"%s\" [label=\"%s\"];\n", funcName, funcName))
		}

		// 依存関係の記述
		for funcName, deps := range result.Dependencies {
			for _, dep := range deps {
				sb.WriteString(fmt.Sprintf("    \"%s\" -> \"%s\";\n", funcName, dep))
			}
		}

		sb.WriteString("  }\n")
		fileCount++
	}

	sb.WriteString("}\n")
	return sb.String(), nil
}
