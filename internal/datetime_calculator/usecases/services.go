package usecases

import (
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/datetime_calculator/config"
)

// #==============================================================#
// ##          DatetimeCalculatorService                         ##
// #==============================================================#
// DatetimeCalculatorService は日時計算を行うサービスです
type DatetimeCalculatorService struct {
	calculator *DatetimeCalculator
	fileReader config.FileReader
}

// NewDatetimeCalculatorService は新しいDatetimeCalculatorServiceを作成します
func NewDatetimeCalculatorService() *DatetimeCalculatorService {
	return &DatetimeCalculatorService{
		calculator: &DatetimeCalculator{},
		fileReader: &config.StandardFileReader{},
	}
}

// NewDatetimeCalculatorServiceWithFileReader はFileReaderを注入した新しいDatetimeCalculatorServiceを作成します
func NewDatetimeCalculatorServiceWithFileReader(fileReader config.FileReader) *DatetimeCalculatorService {
	return &DatetimeCalculatorService{
		calculator: &DatetimeCalculator{},
		fileReader: fileReader,
	}
}

// HandleDatetimeCalc はMCPリクエストを処理して日時計算結果を返すハンドラーです
func (s *DatetimeCalculatorService) HandleDatetimeCalc(op string, year1, month1, day1, hour1, minute1, second1 float64, durationYear, durationMonth, durationDay, durationHour, durationMinute, durationSecond float64) (string, error) {
	var result string
	switch op {
	case "add":
		result = s.calculator.addDatetimeFloat(year1, month1, day1, hour1, minute1, second1, durationYear, durationMonth, durationDay, durationHour, durationMinute, durationSecond)
	case "subtract":
		result = s.calculator.subtractDatetimeFloat(year1, month1, day1, hour1, minute1, second1, durationYear, durationMonth, durationDay, durationHour, durationMinute, durationSecond)
	default:
		return "", fmt.Errorf("unsupported operation: %s", op)
	}
	return result, nil
}

// HandleTimeExtraction はファイルまたはテキストから時間を抽出し合計を計算する
func (s *DatetimeCalculatorService) HandleTimeExtraction(filePath, textInput, outputUnit string) (float64, error) {
	// 排他制御
	if filePath != "" && textInput != "" {
		return 0, fmt.Errorf("ファイルパスとテキスト入力は同時に指定できません")
	}
	if filePath == "" && textInput == "" {
		return 0, fmt.Errorf("ファイルパスまたはテキスト入力のいずれかを指定してください")
	}

	var content string
	if filePath != "" {
		// ファイル拡張子の検証
		if !strings.HasSuffix(filePath, ".md") && !strings.HasSuffix(filePath, ".txt") {
			return 0, fmt.Errorf("ファイルは.mdまたは.txt形式である必要があります")
		}

		// ファイル読み込み（依存性注入されたFileReaderを使用）
		data, err := s.fileReader.ReadFile(filePath)
		if err != nil {
			return 0, fmt.Errorf("ファイル読み込みエラー: %v", err)
		}
		content = string(data)
	} else {
		content = textInput
	}

	return s.calculator.extractTimeFromText(content, outputUnit)
}

// HandleTimeUnitSum は時間単位での合計計算を処理するハンドラーです
func (s *DatetimeCalculatorService) HandleTimeUnitSum(figures []float64, inputUnit, outputUnit string) (float64, error) {
	if len(figures) == 0 {
		return 0, fmt.Errorf("数値の配列が空です")
	}

	// 入力単位で合計を計算
	sum := s.calculator.sumTimeFloat(figures)

	// 時間単位変換を実行
	result, err := s.calculator.convertTimeUnit(sum, inputUnit, outputUnit)
	if err != nil {
		return 0, fmt.Errorf("時間単位変換に失敗しました: %v", err)
	}

	return result, nil
}
