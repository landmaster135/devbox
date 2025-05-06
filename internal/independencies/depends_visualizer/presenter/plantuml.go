package presenter

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/landmaster135/devbox/internal/independencies/depends_visualizer/analyzer"
)

// RenderPlantUML は依存関係を PlantUML 形式で出力します
func RenderPlantUML(results []analyzer.AnalysisResult) (string, error) {
	var sb strings.Builder

	sb.WriteString("@startuml\n")
	sb.WriteString("skinparam defaultTextAlignment center\n")

	for _, result := range results {
		shortName := filepath.Base(result.FilePath)
		sb.WriteString(fmt.Sprintf("' File: %s\n", shortName))

		// 依存関係の記述
		for funcName, deps := range result.Dependencies {
			for _, dep := range deps {
				sb.WriteString(fmt.Sprintf("%s --> %s\n", funcName, dep))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString("@enduml\n")
	return sb.String(), nil
}
