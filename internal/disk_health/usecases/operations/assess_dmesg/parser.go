package assessdmesg

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/landmaster135/devbox/internal/disk_health/domain"
)

var (
	dmesgDevPattern     = regexp.MustCompile(`(?i)\bdev\s+((?:sd|hd|vd|xvd)[a-z][0-9]*|nvme[0-9]+n[0-9]+|mmcblk[0-9]+)\b`)
	dmesgBracketPattern = regexp.MustCompile(`(?i)\[((?:sd|hd|vd|xvd)[a-z][0-9]*|nvme[0-9]+n[0-9]+|mmcblk[0-9]+)\]`)
)

func (s *Service) ParseDmesgLog(content string) ([]domain.DmesgEvent, error) {
	if strings.TrimSpace(content) == "" {
		return nil, fmt.Errorf("dmesgログが空です")
	}

	events := make([]domain.DmesgEvent, 0)
	seen := map[string]struct{}{}
	for _, line := range strings.Split(content, "\n") {
		event, ok := s.parseDmesgLine(line)
		if !ok {
			continue
		}
		key := string(event.Severity) + "\x00" + event.Device + "\x00" + event.Message + "\x00" + event.Line
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		events = append(events, event)
	}
	return events, nil
}

func (s *Service) parseDmesgLine(line string) (domain.DmesgEvent, bool) {
	trimmedLine := strings.TrimSpace(line)
	if trimmedLine == "" {
		return domain.DmesgEvent{}, false
	}

	normalizedLine := strings.ToLower(trimmedLine)
	if s.isDmesgNoiseLine(normalizedLine) {
		return domain.DmesgEvent{}, false
	}

	device := s.extractDmesgDevice(trimmedLine)
	if severity, message, ok := s.classifyDmesgLine(normalizedLine); ok {
		if device == "" && !s.isDeviceOptionalDmesgEvent(normalizedLine) {
			return domain.DmesgEvent{}, false
		}
		return domain.DmesgEvent{
			Severity: severity,
			Device:   device,
			Message:  message,
			Line:     trimmedLine,
		}, true
	}
	return domain.DmesgEvent{}, false
}

func (s *Service) isDmesgNoiseLine(normalizedLine string) bool {
	noisePatterns := []string{
		"acpi bios error",
		"acpi error",
		"directed i/o",
		"symmetric i/o mode",
		"i/o scheduler",
		"correctable errors collector",
		"spinning up disk",
		"attached scsi disk",
		"write protect is off",
		"mode sense:",
		"write cache:",
		"preferred minimum i/o size",
		"optimal transfer size",
	}
	for _, pattern := range noisePatterns {
		if strings.Contains(normalizedLine, pattern) {
			return true
		}
	}
	return false
}

func (s *Service) extractDmesgDevice(line string) string {
	if matches := dmesgDevPattern.FindStringSubmatch(line); matches != nil {
		return matches[1]
	}
	if matches := dmesgBracketPattern.FindStringSubmatch(line); matches != nil {
		return matches[1]
	}
	return ""
}

func (s *Service) classifyDmesgLine(normalizedLine string) (domain.Severity, string, bool) {
	switch {
	case strings.Contains(normalizedLine, "critical medium error"):
		return domain.SeverityCritical, "critical medium error を検出しました", true
	case strings.Contains(normalizedLine, "unrecovered read error"):
		return domain.SeverityCritical, "Unrecovered read error を検出しました", true
	case strings.Contains(normalizedLine, "sense key : medium error") || strings.Contains(normalizedLine, "sense key: medium error"):
		return domain.SeverityCritical, "Medium Error を検出しました", true
	case strings.Contains(normalizedLine, "medium error"):
		return domain.SeverityCritical, "Medium Error を検出しました", true
	case strings.Contains(normalizedLine, "buffer i/o error"):
		return domain.SeverityWarning, "Buffer I/O error を検出しました", true
	case strings.Contains(normalizedLine, "i/o error"):
		return domain.SeverityWarning, "I/O error を検出しました", true
	case strings.Contains(normalizedLine, "failed result"):
		return domain.SeverityWarning, "FAILED Result を検出しました", true
	case strings.Contains(normalizedLine, "failed command"):
		return domain.SeverityWarning, "failed command を検出しました", true
	case strings.Contains(normalizedLine, "timeout"):
		return domain.SeverityWarning, "timeout を検出しました", true
	case strings.Contains(normalizedLine, "link is slow to respond"):
		return domain.SeverityWarning, "link is slow to respond を検出しました", true
	case strings.Contains(normalizedLine, "reset") && strings.Contains(normalizedLine, "link"):
		return domain.SeverityWarning, "link reset を検出しました", true
	default:
		return "", "", false
	}
}

func (s *Service) isDeviceOptionalDmesgEvent(normalizedLine string) bool {
	return strings.Contains(normalizedLine, "critical medium error")
}
