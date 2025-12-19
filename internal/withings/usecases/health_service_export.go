package usecases

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExportDailySummary は日次サマリの結果を Withings Health Mate 形式の JSON として保存します。
func (s *HealthService) ExportDailySummary(resp *DailySummaryResponse, path string) error {
	if resp == nil {
		return fmt.Errorf("daily summary レスポンスが空です")
	}
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("出力ファイルパスが指定されていません")
	}

	flattened := FlattenDailySummaryResponse(resp)
	export := healthMateExport{
		Data:        healthMateExportData{HealthMates: flattened.Summaries},
		Description: "Health Mate data from Withings",
		Name:        "My Health Mate Data",
	}

	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("ディレクトリの作成に失敗しました: %w", err)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("ファイルの作成に失敗しました: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(export); err != nil {
		return fmt.Errorf("JSON のエンコードに失敗しました: %w", err)
	}

	return nil
}

type healthMateExport struct {
	Data        healthMateExportData `json:"data"`
	Description string               `json:"description"`
	Name        string               `json:"name"`
}

type healthMateExportData struct {
	HealthMates []FlattenedDailySummary `json:"health_mates"`
}
