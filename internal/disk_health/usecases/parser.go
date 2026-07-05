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
	logicalSectorPattern = regexp.MustCompile(`(?i)(\d+)\s+bytes\s+logical`)
)

const defaultLogicalSectorSizeBytes int64 = 512

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
	s.buildDiskInfo(report)

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
	case "Rotation Rate":
		if rate, ok := s.parseLeadingInt64(value); ok {
			s.ensureDiskInfo(report).RotationRateRPM = rate
		}
	case "Sector Size", "Sector Sizes":
		if size, ok := s.parseLogicalSectorSize(value); ok {
			report.LogicalSectorSizeBytes = size
		}
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

func (s *Service) parseLeadingInt64(text string) (*int64, bool) {
	matches := rawValuePattern.FindStringSubmatch(text)
	if matches == nil {
		return nil, false
	}
	value, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return nil, false
	}
	return &value, true
}

func (s *Service) parseLogicalSectorSize(text string) (*int64, bool) {
	matches := logicalSectorPattern.FindStringSubmatch(text)
	if matches != nil {
		value, err := strconv.ParseInt(matches[1], 10, 64)
		if err == nil {
			return &value, true
		}
	}
	return s.parseLeadingInt64(text)
}

func (s *Service) buildDiskInfo(report *domain.SmartReport) {
	diskInfo := report.DiskInfo
	if diskInfo == nil {
		diskInfo = &domain.DiskInfo{}
	}

	s.setDiskInfoAttributeValues(diskInfo, report.Attributes)
	s.setDiskInfoByteValues(diskInfo, report.LogicalSectorSizeBytes)

	if s.hasDiskInfoValue(diskInfo) {
		report.DiskInfo = diskInfo
		return
	}
	report.DiskInfo = nil
}

func (s *Service) setDiskInfoAttributeValues(diskInfo *domain.DiskInfo, attributes []domain.SmartAttribute) {
	if attribute, ok := s.findAttribute(attributes, 9, "Power_On_Hours"); ok {
		diskInfo.PowerOnHours = s.int64Pointer(attribute.RawValue)
	}
	if attribute, ok := s.findAttribute(attributes, 12, "Power_Cycle_Count"); ok {
		diskInfo.PowerCycleCount = s.int64Pointer(attribute.RawValue)
	}
	if attribute, ok := s.findAttribute(attributes, 194, "Temperature_Celsius"); ok {
		diskInfo.TemperatureCelsius = s.int64Pointer(attribute.RawValue)
	} else if attribute, ok := s.findAttribute(attributes, 190, "Airflow_Temperature_Cel"); ok {
		diskInfo.TemperatureCelsius = s.int64Pointer(attribute.RawValue)
	}
	if attribute, ok := s.findAttribute(attributes, 241, "Total_LBAs_Written"); ok {
		diskInfo.TotalLBAsWritten = s.int64Pointer(attribute.RawValue)
	}
	if attribute, ok := s.findAttribute(attributes, 242, "Total_LBAs_Read"); ok {
		diskInfo.TotalLBAsRead = s.int64Pointer(attribute.RawValue)
	}
}

func (s *Service) setDiskInfoByteValues(diskInfo *domain.DiskInfo, logicalSectorSize *int64) {
	sectorSize := defaultLogicalSectorSizeBytes
	if logicalSectorSize != nil && *logicalSectorSize > 0 {
		sectorSize = *logicalSectorSize
	}
	if diskInfo.TotalLBAsWritten != nil {
		diskInfo.TotalBytesWritten = s.int64Pointer(*diskInfo.TotalLBAsWritten * sectorSize)
	}
	if diskInfo.TotalLBAsRead != nil {
		diskInfo.TotalBytesRead = s.int64Pointer(*diskInfo.TotalLBAsRead * sectorSize)
	}
}

func (s *Service) hasDiskInfoValue(diskInfo *domain.DiskInfo) bool {
	return diskInfo.RotationRateRPM != nil ||
		diskInfo.PowerOnHours != nil ||
		diskInfo.PowerCycleCount != nil ||
		diskInfo.TemperatureCelsius != nil ||
		diskInfo.TotalLBAsWritten != nil ||
		diskInfo.TotalBytesWritten != nil ||
		diskInfo.TotalLBAsRead != nil ||
		diskInfo.TotalBytesRead != nil
}

func (s *Service) ensureDiskInfo(report *domain.SmartReport) *domain.DiskInfo {
	if report.DiskInfo == nil {
		report.DiskInfo = &domain.DiskInfo{}
	}
	return report.DiskInfo
}

func (s *Service) int64Pointer(value int64) *int64 {
	return &value
}
