// internal/analyzer/project.go
package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/landmaster135/devbox/internal/code_analyzer/models"
)

// ProjectAnalyzer はプロジェクト分析機能を提供します
type ProjectAnalyzer struct {
	FileAnalyzer  *FileAnalyzer
	CloneDetector *CloneDetector
}

// NewProjectAnalyzer は新しいProjectAnalyzerインスタンスを作成します
func NewProjectAnalyzer(fileAnalyzer *FileAnalyzer, cloneDetector *CloneDetector) *ProjectAnalyzer {
	return &ProjectAnalyzer{
		FileAnalyzer:  fileAnalyzer,
		CloneDetector: cloneDetector,
	}
}

// AnalyzeProject はプロジェクト全体を分析します
func (a *ProjectAnalyzer) AnalyzeProject(path string, extensions []string, detectCloneFlag bool) (models.ProjectMetrics, error) {
	projMetrics := models.ProjectMetrics{
		ProjectPath: path,
		AnalyzedAt:  time.Now(),
	}

	var fileList []string
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(filePath))
			for _, validExt := range extensions {
				if ext == validExt {
					fileList = append(fileList, filePath)
					break
				}
			}
		}
		return nil
	})

	if err != nil {
		return projMetrics, err
	}

	// 並列処理でファイル分析
	var wg sync.WaitGroup
	var mutex sync.Mutex

	projMetrics.Files = make([]models.FileMetrics, len(fileList))
	fileContents := make(map[string]string)

	maxWorkers := 4
	if len(fileList) < maxWorkers {
		maxWorkers = len(fileList)
	}

	jobs := make(chan int, len(fileList))

	// ワーカーの起動
	for w := 0; w < maxWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				fileMetrics, content, err := a.FileAnalyzer.AnalyzeFile(fileList[idx])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error analyzing %s: %v\n", fileList[idx], err)
					continue
				}

				mutex.Lock()
				projMetrics.Files[idx] = fileMetrics
				fileContents[fileList[idx]] = content
				mutex.Unlock()
			}
		}()
	}

	// ジョブの送信
	for i := range fileList {
		jobs <- i
	}
	close(jobs)

	// すべてのワーカーの完了を待つ
	wg.Wait()

	// プロジェクト全体の集計
	projMetrics.FileCount = len(projMetrics.Files)
	var totalComplexity int

	for _, fm := range projMetrics.Files {
		projMetrics.TotalLines += fm.TotalLines
		projMetrics.TotalCodeLines += fm.CodeLines
		projMetrics.TotalComments += fm.CommentLines
		projMetrics.TotalBlankLines += fm.BlankLines
		totalComplexity += fm.Complexity

		if fm.Complexity > projMetrics.MaxComplexity {
			projMetrics.MaxComplexity = fm.Complexity
		}
	}

	if projMetrics.FileCount > 0 {
		projMetrics.AvgComplexity = float64(totalComplexity) / float64(projMetrics.FileCount)
	}

	if projMetrics.TotalCodeLines > 0 {
		projMetrics.CommentRatio = float64(projMetrics.TotalComments) / float64(projMetrics.TotalCodeLines) * 100.0
	}

	// コードクローン検出
	if detectCloneFlag && a.CloneDetector != nil {
		projMetrics.Clones = a.CloneDetector.DetectClones(projMetrics.Files, fileContents)
	}

	return projMetrics, nil
}
