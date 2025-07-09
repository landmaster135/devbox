// internal/history/trend.go (続き)
package history

import (
	"github.com/landmaster135/devbox/internal/code_analyzer/models"
)

// TrendAnalyzer はトレンド分析機能を提供します
type TrendAnalyzer struct{}

// NewTrendAnalyzer は新しいTrendAnalyzerインスタンスを作成します
func NewTrendAnalyzer() *TrendAnalyzer {
	return &TrendAnalyzer{}
}

// AnalyzeTrends は時系列トレンドを分析します
func (a *TrendAnalyzer) AnalyzeTrends(currentMetrics models.ProjectMetrics, history []models.HistoricalData) map[string]models.TrendMetric {
	if len(history) < 2 {
		return nil // 履歴が不足している場合
	}

	previousData := history[len(history)-2] // 最新の前のデータ
	trends := make(map[string]models.TrendMetric)

	// 総行数のトレンド
	trends["total_lines"] = models.TrendMetric{
		Current:    float64(currentMetrics.TotalLines),
		Previous:   float64(previousData.TotalLines),
		Change:     float64(currentMetrics.TotalLines - previousData.TotalLines),
		ChangeRate: a.calculateChangeRate(float64(currentMetrics.TotalLines), float64(previousData.TotalLines)),
	}

	// コード行数のトレンド
	trends["code_lines"] = models.TrendMetric{
		Current:    float64(currentMetrics.TotalCodeLines),
		Previous:   float64(previousData.CodeLines),
		Change:     float64(currentMetrics.TotalCodeLines - previousData.CodeLines),
		ChangeRate: a.calculateChangeRate(float64(currentMetrics.TotalCodeLines), float64(previousData.CodeLines)),
	}

	// コメント率のトレンド
	trends["comment_ratio"] = models.TrendMetric{
		Current:    currentMetrics.CommentRatio,
		Previous:   previousData.CommentRatio,
		Change:     currentMetrics.CommentRatio - previousData.CommentRatio,
		ChangeRate: a.calculateChangeRate(currentMetrics.CommentRatio, previousData.CommentRatio),
	}

	// 平均複雑度のトレンド
	trends["avg_complexity"] = models.TrendMetric{
		Current:    currentMetrics.AvgComplexity,
		Previous:   previousData.AvgComplexity,
		Change:     currentMetrics.AvgComplexity - previousData.AvgComplexity,
		ChangeRate: a.calculateChangeRate(currentMetrics.AvgComplexity, previousData.AvgComplexity),
	}

	// 最大複雑度のトレンド
	trends["max_complexity"] = models.TrendMetric{
		Current:    float64(currentMetrics.MaxComplexity),
		Previous:   float64(previousData.MaxComplexity),
		Change:     float64(currentMetrics.MaxComplexity - previousData.MaxComplexity),
		ChangeRate: a.calculateChangeRate(float64(currentMetrics.MaxComplexity), float64(previousData.MaxComplexity)),
	}

	// クローン率のトレンド
	cloneLineCount := 0
	for _, clone := range currentMetrics.Clones {
		cloneLineCount += clone.LineCount
	}

	currentCloneRatio := 0.0
	if currentMetrics.TotalCodeLines > 0 {
		currentCloneRatio = float64(cloneLineCount) / float64(currentMetrics.TotalCodeLines) * 100.0
	}

	trends["clone_ratio"] = models.TrendMetric{
		Current:    currentCloneRatio,
		Previous:   previousData.CloneLineRatio,
		Change:     currentCloneRatio - previousData.CloneLineRatio,
		ChangeRate: a.calculateChangeRate(currentCloneRatio, previousData.CloneLineRatio),
	}

	return trends
}

// 変化率を計算する関数
func (a *TrendAnalyzer) calculateChangeRate(current, previous float64) float64 {
	if previous == 0 {
		return 0 // ゼロ除算を避ける
	}
	return (current - previous) / previous * 100.0
}
