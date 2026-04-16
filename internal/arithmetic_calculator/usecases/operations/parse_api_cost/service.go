package parseapicost

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/landmaster135/devbox/internal/arithmetic_calculator/config"
)

type Service struct {
	fileReader config.FileReader
}

func NewService() *Service {
	return &Service{fileReader: &config.StandardFileReader{}}
}

func NewServiceWithFileReader(fileReader config.FileReader) *Service {
	return &Service{fileReader: fileReader}
}

func (s *Service) ExtractAPICostFromText(text string) (float64, error) {
	pattern := `API料金が(\d+)円掛かった`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(text, -1)

	total := 0.0
	for _, match := range matches {
		if len(match) <= 1 {
			continue
		}
		cost, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		total += float64(cost)
	}

	return total, nil
}

func (s *Service) HandleApiCostExtraction(filePath, textInput string) (float64, error) {
	if filePath != "" && textInput != "" {
		return 0, fmt.Errorf("ファイルパスとテキスト入力は同時に指定できません")
	}
	if filePath == "" && textInput == "" {
		return 0, fmt.Errorf("ファイルパスまたはテキスト入力のいずれかを指定してください")
	}

	var content string
	if filePath != "" {
		if !strings.HasSuffix(filePath, ".md") && !strings.HasSuffix(filePath, ".txt") {
			return 0, fmt.Errorf("ファイルは.mdまたは.txt形式である必要があります")
		}

		data, err := s.fileReader.ReadFile(filePath)
		if err != nil {
			return 0, fmt.Errorf("ファイル読み込みエラー: %v", err)
		}
		content = string(data)
	} else {
		content = textInput
	}

	return s.ExtractAPICostFromText(content)
}

func (s *Service) Execute(filePath, textInput string) (string, error) {
	result, err := s.HandleApiCostExtraction(filePath, textInput)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("抽出されたAPI料金の合計: %.0f円\n", result), nil
}
