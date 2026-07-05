package usecases

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/disk_health/domain"
)

func (s *Service) FormatText(assessment domain.Assessment, verbose bool) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("status: %s\n", assessment.Status))
	builder.WriteString(fmt.Sprintf("score: %d\n", assessment.Score))
	builder.WriteString(fmt.Sprintf("summary: %s\n", assessment.Summary))

	if assessment.Model != "" {
		builder.WriteString(fmt.Sprintf("model: %s\n", assessment.Model))
	}
	if assessment.SerialNumber != "" {
		builder.WriteString(fmt.Sprintf("serial_number: %s\n", assessment.SerialNumber))
	}
	if assessment.OverallHealth != "" {
		builder.WriteString(fmt.Sprintf("overall_health: %s\n", assessment.OverallHealth))
	}

	if len(assessment.Findings) > 0 {
		builder.WriteString("findings:\n")
		for _, finding := range assessment.Findings {
			if finding.AttributeName != "" {
				builder.WriteString(fmt.Sprintf("- [%s] %s raw=%d: %s\n", finding.Severity, finding.AttributeName, finding.RawValue, finding.Message))
				continue
			}
			builder.WriteString(fmt.Sprintf("- [%s] %s\n", finding.Severity, finding.Message))
		}
	}

	if verbose && len(assessment.Attributes) > 0 {
		builder.WriteString("attributes:\n")
		for _, attribute := range assessment.Attributes {
			builder.WriteString(fmt.Sprintf("- id=%d name=%s value=%d worst=%d thresh=%d when_failed=%s raw=%s\n",
				attribute.ID,
				attribute.Name,
				attribute.Value,
				attribute.Worst,
				attribute.Threshold,
				attribute.WhenFailed,
				attribute.RawText,
			))
		}
	}

	return builder.String()
}

func (s *Service) FormatJSON(assessment domain.Assessment, verbose bool) (string, error) {
	if !verbose {
		assessment.Attributes = nil
	}

	output, err := json.MarshalIndent(assessment, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON出力の生成に失敗しました: %w", err)
	}
	return string(output) + "\n", nil
}
