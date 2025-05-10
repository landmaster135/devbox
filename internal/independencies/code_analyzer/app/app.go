// internal/app/app.go
package app

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/landmaster135/devbox/internal/independencies/code_analyzer/analyzer"
	"github.com/landmaster135/devbox/internal/independencies/code_analyzer/config"
	"github.com/landmaster135/devbox/internal/independencies/code_analyzer/history"
	"github.com/landmaster135/devbox/internal/independencies/code_analyzer/models"
	"github.com/landmaster135/devbox/internal/independencies/code_analyzer/report"
)

// App はアプリケーションの中心的な処理を行います
type App struct {
	Config          *config.AppConfig
	ProjectAnalyzer *analyzer.ProjectAnalyzer
	HistoryManager  *history.HistoryManager
	TrendAnalyzer   *history.TrendAnalyzer
}

// NewApp は新しいAppインスタンスを作成します
func NewApp(config *config.AppConfig) *App {
	fileAnalyzer := analyzer.NewFileAnalyzer()
	cloneDetector := analyzer.NewCloneDetector(config.MinBlockSize, config.MinSimilarity)
	projectAnalyzer := analyzer.NewProjectAnalyzer(fileAnalyzer, cloneDetector)
	historyManager := history.NewHistoryManager(config.HistoryPath)
	trendAnalyzer := history.NewTrendAnalyzer()

	return &App{
		Config:          config,
		ProjectAnalyzer: projectAnalyzer,
		HistoryManager:  historyManager,
		TrendAnalyzer:   trendAnalyzer,
	}
}

// Run はアプリケーションのメインロジックを実行します
func (a *App) Run(stdout, stderr io.Writer) int {
	// プロジェクト分析
	fmt.Fprintf(stdout, "Analyzing project at: %s\n", a.Config.ProjectPath)

	metrics, err := a.ProjectAnalyzer.AnalyzeProject(
		a.Config.ProjectPath,
		a.Config.Extensions,
		a.Config.DetectClones,
	)

	if err != nil {
		fmt.Fprintf(stderr, "Error analyzing project: %v\n", err)
		return 1
	}

	// 履歴データの読み込み
	var history []models.HistoricalData
	if a.Config.HistoryPath != "" {
		history, err = a.HistoryManager.LoadHistory()
		if err != nil {
			fmt.Fprintf(stderr, "Error loading history: %v\n", err)
		} else {
			// トレンド分析
			metrics.Trends = a.TrendAnalyzer.AnalyzeTrends(metrics, history)

			// 履歴の更新
			if err := a.HistoryManager.SaveHistory(history, metrics); err != nil {
				fmt.Fprintf(stderr, "Error saving history: %v\n", err)
			}
		}
	}

	// レポート生成
	var output string
	var reportErr error

	switch a.Config.OutputFormat {
	case "json":
		jsonReporter := report.NewJSONReporter()
		output, reportErr = jsonReporter.Generate(metrics)

	case "csv":
		csvReporter := report.NewCSVReporter()
		output = csvReporter.Generate(metrics)

	default: // text
		textReporter := report.NewTextReporter()
		output = textReporter.Generate(metrics)
	}

	if reportErr != nil {
		fmt.Fprintf(stderr, "Error generating report: %v\n", reportErr)
		return 1
	}

	// ビジュアルレポート生成
	if a.Config.VisualReport {
		htmlReporter := report.NewHTMLReporter()
		htmlOutput, err := htmlReporter.Generate(metrics, history)
		if err != nil {
			fmt.Fprintf(stderr, "Error generating HTML report: %v\n", err)
		} else {
			htmlPath := "code_metrics_report.html"
			if a.Config.OutputFile != "" {
				// 出力ファイルの拡張子を .html に変更
				ext := filepath.Ext(a.Config.OutputFile)
				if ext != "" {
					htmlPath = a.Config.OutputFile[:len(a.Config.OutputFile)-len(ext)] + ".html"
				} else {
					htmlPath = a.Config.OutputFile + ".html"
				}
			}

			if err := ioutil.WriteFile(htmlPath, []byte(htmlOutput), 0644); err != nil {
				fmt.Fprintf(stderr, "Error writing HTML report: %v\n", err)
			} else {
				fmt.Fprintf(stdout, "Visual report generated successfully: %s\n", htmlPath)
			}
		}
	}

	// テキスト出力
	if a.Config.OutputFile != "" && !a.Config.VisualReport {
		if err := ioutil.WriteFile(a.Config.OutputFile, []byte(output), 0644); err != nil {
			fmt.Fprintf(stderr, "Error writing to file: %v\n", err)
			return 1
		}
		fmt.Fprintf(stdout, "Results written to: %s\n", a.Config.OutputFile)
	} else if !a.Config.VisualReport {
		fmt.Fprintln(stdout, output)
	}

	return 0
}
