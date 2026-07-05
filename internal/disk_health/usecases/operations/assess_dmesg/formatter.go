package assessdmesg

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/disk_health/domain"
)

func (s *Service) FormatText(assessment domain.DmesgAssessment, verbose bool) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("status: %s\n", assessment.Status))
	builder.WriteString(fmt.Sprintf("score: %d\n", assessment.Score))
	builder.WriteString(fmt.Sprintf("summary: %s\n", assessment.Summary))

	if len(assessment.Findings) > 0 {
		builder.WriteString("findings:\n")
		for _, finding := range assessment.Findings {
			if finding.Device != "" {
				builder.WriteString(fmt.Sprintf("- [%s] %s: %s\n", finding.Severity, finding.Device, finding.Message))
			} else {
				builder.WriteString(fmt.Sprintf("- [%s] %s\n", finding.Severity, finding.Message))
			}
			if verbose && finding.Line != "" {
				builder.WriteString(fmt.Sprintf("  raw: %s\n", finding.Line))
			}
		}
	}

	return builder.String()
}

func (s *Service) FormatJSON(assessment domain.DmesgAssessment, verbose bool) (string, error) {
	if !verbose {
		for i := range assessment.Findings {
			assessment.Findings[i].Line = ""
		}
	}

	output, err := json.MarshalIndent(assessment, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON出力の生成に失敗しました: %w", err)
	}
	return string(output) + "\n", nil
}
