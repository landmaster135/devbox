package evaluatelinecount

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type FileOpener interface {
	Open(name string) (*os.File, error)
}

type DefaultFileOpener struct{}

func (o *DefaultFileOpener) Open(name string) (*os.File, error) {
	return os.Open(name)
}

type BufioScanner interface {
	Scan() bool
	Err() error
}

type JSONMarshaler interface {
	MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
}

type DefaultJSONMarshaler struct{}

func (m *DefaultJSONMarshaler) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

type Service struct {
	fileOpener    FileOpener
	bufioScanner  BufioScanner
	jsonMarshaler JSONMarshaler
}

func NewService() *Service {
	return &Service{
		fileOpener:    &DefaultFileOpener{},
		bufioScanner:  &bufio.Scanner{},
		jsonMarshaler: &DefaultJSONMarshaler{},
	}
}

func NewServiceWithDependencies(fileOpener FileOpener, bufioScanner BufioScanner, jsonMarshaler JSONMarshaler) *Service {
	return &Service{
		fileOpener:    fileOpener,
		bufioScanner:  bufioScanner,
		jsonMarshaler: jsonMarshaler,
	}
}

func (s *Service) CountLines(filePath string) (int, error) {
	file, err := s.fileOpener.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("ファイルを開けませんでした: %w", err)
	}
	defer file.Close()

	var scanner BufioScanner
	if err := s.bufioScanner.Err(); err != nil {
		scanner = s.bufioScanner
	} else {
		scanner = bufio.NewScanner(file)
		s.bufioScanner = scanner
	}

	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}

func (s *Service) IsLineCountGreaterThan(filePath string, threshold int) (bool, int, error) {
	lineCount, err := s.CountLines(filePath)
	if err != nil {
		return false, 0, err
	}
	return lineCount > threshold, lineCount, nil
}

func IsGreaterDescription(isGreater bool) string {
	if isGreater {
		return "より大きいです。"
	}
	return "以下です。"
}

func (s *Service) HandleToEvaluateLineCount(filePath string, threshold int) (string, error) {
	isGreater, lineCount, err := s.IsLineCountGreaterThan(filePath, threshold)
	if err != nil {
		return "", err
	}

	result := map[string]interface{}{
		"is_greater":  isGreater,
		"line_count":  lineCount,
		"threshold":   threshold,
		"file_path":   filePath,
		"description": fmt.Sprintf("ファイル '%s' の行数は %d 行で、閾値 %d %s", filePath, lineCount, threshold, IsGreaterDescription(isGreater)),
	}

	jsonResult, err := s.jsonMarshaler.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonResult), nil
}

func (s *Service) Execute(filePath string, threshold int) (string, error) {
	return s.HandleToEvaluateLineCount(filePath, threshold)
}
