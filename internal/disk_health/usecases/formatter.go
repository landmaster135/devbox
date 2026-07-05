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
	s.writeDiskInfoText(&builder, assessment.DiskInfo)

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

func (s *Service) writeDiskInfoText(builder *strings.Builder, diskInfo *domain.DiskInfo) {
	if diskInfo == nil {
		return
	}

	lines := s.formatDiskInfoLines(diskInfo)
	if len(lines) == 0 {
		return
	}

	builder.WriteString("disk_info:\n")
	for _, line := range lines {
		builder.WriteString(line)
	}
}

func (s *Service) formatDiskInfoLines(diskInfo *domain.DiskInfo) []string {
	lines := make([]string, 0, 8)
	lines = s.appendInt64Line(lines, "rotation_rate_rpm", diskInfo.RotationRateRPM)
	lines = s.appendInt64Line(lines, "power_on_hours", diskInfo.PowerOnHours)
	lines = s.appendInt64Line(lines, "power_cycle_count", diskInfo.PowerCycleCount)
	lines = s.appendInt64Line(lines, "temperature_celsius", diskInfo.TemperatureCelsius)
	lines = s.appendInt64Line(lines, "total_lbas_written", diskInfo.TotalLBAsWritten)
	lines = s.appendInt64Line(lines, "total_bytes_written", diskInfo.TotalBytesWritten)
	lines = s.appendInt64Line(lines, "total_lbas_read", diskInfo.TotalLBAsRead)
	lines = s.appendInt64Line(lines, "total_bytes_read", diskInfo.TotalBytesRead)
	return lines
}

func (s *Service) appendInt64Line(lines []string, name string, value *int64) []string {
	if value == nil {
		return lines
	}
	return append(lines, fmt.Sprintf("  %s: %d\n", name, *value))
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
