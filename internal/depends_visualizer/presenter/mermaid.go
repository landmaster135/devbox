package presenter

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/landmaster135/devbox/internal/depends_visualizer/analyzer"
)

// RenderMermaid は依存関係を Mermaid クラス図形式で出力します
func RenderMermaid(results []analyzer.AnalysisResult) (string, error) {
	var sb strings.Builder

	sb.WriteString("```mermaid\n")
	sb.WriteString("classDiagram\n")

	for _, result := range results {
		// ファイル名をコメントとして表示
		shortName := filepath.Base(result.FilePath)
		sb.WriteString("  %%% File: " + shortName + "\n")

		// 依存関係の記述
		for funcName, deps := range result.Dependencies {
			for _, dep := range deps {
				// 特殊文字をエスケープ
				escapedFuncName := escapeMermaidIdentifier(funcName)
				escapedDep := escapeMermaidIdentifier(dep)
				sb.WriteString(fmt.Sprintf("  %s <|-- %s\n", escapedDep, escapedFuncName))
			}
		}

		// クラス定義
		for funcName := range result.Dependencies {
			escapedFuncName := escapeMermaidIdentifier(funcName)
			// クラス定義の中に元の関数名を表示
			sb.WriteString(fmt.Sprintf("  class %s {\n", escapedFuncName))
			sb.WriteString(fmt.Sprintf("    %s\n", funcName))
			sb.WriteString("  }\n")
		}

		sb.WriteString("\n")
	}

	sb.WriteString("```\n")
	return sb.String(), nil
}

// RenderMermaidFlowchart は依存関係を Mermaid フローチャート形式で出力します
func RenderMermaidFlowchart(results []analyzer.AnalysisResult) (string, error) {
	var sb strings.Builder

	sb.WriteString("```mermaid\n")
	sb.WriteString("flowchart TD\n")

	// 各ファイルをサブグラフとして表現
	fileIndex := 0
	for _, result := range results {
		shortName := filepath.Base(result.FilePath)

		// ファイルをサブグラフとして定義
		sb.WriteString(fmt.Sprintf("  subgraph file_%d[\"%s\"]\n", fileIndex, shortName))

		// ノードの定義
		for funcName := range result.Dependencies {
			// 関数名をIDとして使用（スペースを含まないID）
			nodeID := sanitizeNodeID(funcName, fileIndex)
			sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", nodeID, funcName))
		}
		sb.WriteString("  end\n")
		fileIndex++
	}

	sb.WriteString("\n") // 空行で分離

	// 依存関係の定義
	fileIndex = 0
	for _, result := range results {
		for funcName, deps := range result.Dependencies {
			funcNodeID := sanitizeNodeID(funcName, fileIndex)
			for _, dep := range deps {
				// 依存関係のノードIDを検索
				depNodeID := findNodeID(dep, results, fileIndex)
				if depNodeID != "" {
					sb.WriteString(fmt.Sprintf("  %s --> %s\n", funcNodeID, depNodeID))
				}
			}
		}
		fileIndex++
	}

	sb.WriteString("```\n")
	return sb.String(), nil
}

// Mermaid識別子をエスケープ
func escapeMermaidIdentifier(id string) string {
	// 特殊文字を含む場合は識別子をエスケープ
	if strings.Contains(id, "__") || strings.Contains(id, "-") ||
	   strings.ContainsAny(id, "!@#$%^&*()+={}[]|\\:;\"'<>,.?/") {
		// アンダースコアを他の文字に置き換え
		safeID := strings.ReplaceAll(id, "__", "DUNDER")
		safeID = strings.ReplaceAll(safeID, "-", "DASH")

		// その他の特殊文字を置換
		for _, char := range "!@#$%^&*()+={}[]|\\:;\"'<>,.?/" {
			safeID = strings.ReplaceAll(safeID, string(char), "X")
		}

		// 先頭が数字の場合、先頭にプレフィックスを追加
		if len(safeID) > 0 && safeID[0] >= '0' && safeID[0] <= '9' {
			safeID = "func_" + safeID
		}

		return safeID
	}
	return id
}

// ノードIDを安全に作成
func sanitizeNodeID(id string, fileIndex int) string {
	// 置換マップ
	replacements := strings.NewReplacer(
		" ", "_",
		"-", "_",
		".", "_",
		"/", "_",
		"\\", "_",
		"(", "_",
		")", "_",
		"[", "_",
		"]", "_",
		"{", "_",
		"}", "_",
		"\"", "_",
		"'", "_",
		"`", "_",
		":", "_",
		";", "_",
		",", "_",
		"<", "_",
		">", "_",
		"!", "_",
		"@", "_",
		"#", "_",
		"$", "_",
		"%", "_",
		"^", "_",
		"&", "_",
		"*", "_",
		"+", "_",
		"=", "_",
		"|", "_",
		"?", "_",
	)

	result := replacements.Replace(id)

	// ファイルインデックスを追加して一意性を確保
	return fmt.Sprintf("func_%d_%s", fileIndex, result)
}

// 依存関係のノードIDを検索
func findNodeID(funcName string, results []analyzer.AnalysisResult, currentFileIndex int) string {
	// まず同じファイル内で検索
	if _, exists := results[currentFileIndex].Dependencies[funcName]; exists {
		return sanitizeNodeID(funcName, currentFileIndex)
	}

	// 他のファイルで検索
	for i, result := range results {
		if i != currentFileIndex {
			if _, exists := result.Dependencies[funcName]; exists {
				return sanitizeNodeID(funcName, i)
			}
		}
	}

	// 見つからない場合は空文字列を返す
	return ""
}
