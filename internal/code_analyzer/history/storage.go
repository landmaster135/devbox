// internal/history/storage.go
package history

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/landmaster135/devbox/internal/code_analyzer/models"
)

// HistoryManager は履歴データの管理機能を提供します
type HistoryManager struct {
	HistoryPath string
}

// NewHistoryManager は新しいHistoryManagerインスタンスを作成します
func NewHistoryManager(historyPath string) *HistoryManager {
	return &HistoryManager{
		HistoryPath: historyPath,
	}
}

// LoadHistory は過去のメトリクスデータを読み込みます
func (m *HistoryManager) LoadHistory() ([]models.HistoricalData, error) {
	if m.HistoryPath == "" {
		return []models.HistoricalData{}, nil
	}

	if _, err := os.Stat(m.HistoryPath); os.IsNotExist(err) {
		return []models.HistoricalData{}, nil
	}

	data, err := ioutil.ReadFile(m.HistoryPath)
	if err != nil {
		return nil, err
	}

	var history []models.HistoricalData
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, err
	}

	return history, nil
}

// SaveHistory は履歴データを保存します
func (m *HistoryManager) SaveHistory(history []models.HistoricalData, metrics models.ProjectMetrics) error {
	if m.HistoryPath == "" {
		return nil
	}

	// 現在のメトリクスを履歴に追加
	cloneLineCount := 0
	for _, clone := range metrics.Clones {
		cloneLineCount += clone.LineCount
	}

	cloneLineRatio := 0.0
	if metrics.TotalCodeLines > 0 {
		cloneLineRatio = float64(cloneLineCount) / float64(metrics.TotalCodeLines) * 100.0
	}

	currentData := models.HistoricalData{
		Date:           metrics.AnalyzedAt,
		TotalLines:     metrics.TotalLines,
		CodeLines:      metrics.TotalCodeLines,
		CommentLines:   metrics.TotalComments,
		BlankLines:     metrics.TotalBlankLines,
		FileCount:      metrics.FileCount,
		AvgComplexity:  metrics.AvgComplexity,
		MaxComplexity:  metrics.MaxComplexity,
		CommentRatio:   metrics.CommentRatio,
		CloneCount:     len(metrics.Clones),
		CloneLineRatio: cloneLineRatio,
	}

	history = append(history, currentData)

	// 履歴を最大100エントリーに制限
	if len(history) > 100 {
		history = history[len(history)-100:]
	}

	// 履歴を保存
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(m.HistoryPath, data, 0644)
}
