// internal/report/json.go
package report

import (
	"encoding/json"

	"github.com/landmaster135/devbox/internal/code_analyzer/models"
)

// JSONReporter はJSON形式のレポート生成機能を提供します
type JSONReporter struct{}

// NewJSONReporter は新しいJSONReporterインスタンスを作成します
func NewJSONReporter() *JSONReporter {
	return &JSONReporter{}
}

// Generate はJSON形式のレポートを生成します
func (r *JSONReporter) Generate(metrics models.ProjectMetrics) (string, error) {
	jsonData, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
