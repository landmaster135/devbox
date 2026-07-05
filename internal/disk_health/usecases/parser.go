package usecases

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/landmaster135/devbox/internal/disk_health/domain"
)

var (
	attributeLinePattern = regexp.MustCompile(`^\s*(\d+)\s+(\S+)\s+\S+\s+(\d+)\s+(\d+)\s+(\d+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(.+?)\s*$`)
	rawValuePattern      = regexp.MustCompile(`^\s*(-?\d+)`)
)

func (s *Service) ParseSmartReport(content string) (*domain.SmartReport, error) {
	if strings.TrimSpace(content) == "" {
		return nil, fmt.Errorf("SMART情報が空です")
	}

	report := &domain.SmartReport{}
	for _, line := range strings.Split(content, "\n") {
		s.parseInfoLine(report, line)
		if attribute, ok := s.parseAttributeLine(line); ok {
			report.Attributes = append(report.Attributes, attribute)
		}
	}

	return report, nil
}

func (s *Service) parseInfoLine(report *domain.SmartReport, line string) {
	key, value, ok := strings.Cut(line, ":")
	if !ok {
		return
	}

	value = strings.TrimSpace(value)
	switch strings.TrimSpace(key) {
	case "Device Model":
		report.Model = value
	case "Serial Number":
		report.SerialNumber = value
	case "SMART overall-health self-assessment test result":
		report.OverallHealth = strings.ToUpper(value)
	}
}

func (s *Service) parseAttributeLine(line string) (domain.SmartAttribute, bool) {
	matches := attributeLinePattern.FindStringSubmatch(line)
	if matches == nil {
		return domain.SmartAttribute{}, false
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return domain.SmartAttribute{}, false
	}
	value, err := strconv.Atoi(matches[3])
	if err != nil {
		return domain.SmartAttribute{}, false
	}
	worst, err := strconv.Atoi(matches[4])
	if err != nil {
		return domain.SmartAttribute{}, false
	}
	threshold, err := strconv.Atoi(matches[5])
	if err != nil {
		return domain.SmartAttribute{}, false
	}
	rawValue, err := s.parseRawValue(matches[9])
	if err != nil {
		return domain.SmartAttribute{}, false
	}

	return domain.SmartAttribute{
		ID:         id,
		Name:       matches[2],
		Value:      value,
		Worst:      worst,
		Threshold:  threshold,
		Type:       matches[6],
		Updated:    matches[7],
		WhenFailed: matches[8],
		RawValue:   rawValue,
		RawText:    strings.TrimSpace(matches[9]),
	}, true
}

func (s *Service) parseRawValue(rawText string) (int64, error) {
	matches := rawValuePattern.FindStringSubmatch(rawText)
	if matches == nil {
		return 0, fmt.Errorf("RAW_VALUEの解析に失敗しました: %s", rawText)
	}
	return strconv.ParseInt(matches[1], 10, 64)
}
