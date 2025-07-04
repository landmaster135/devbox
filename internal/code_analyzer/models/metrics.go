// internal/models/metrics.go
package models

import (
	"time"
)

// FileMetrics は1つのファイルのメトリクス情報を保持します
type FileMetrics struct {
	Path            string            `json:"path"`
	TotalLines      int               `json:"total_lines"`
	CodeLines       int               `json:"code_lines"`
	CommentLines    int               `json:"comment_lines"`
	BlankLines      int               `json:"blank_lines"`
	FunctionCount   int               `json:"function_count"`
	AvgFunctionSize float64           `json:"avg_function_size"`
	MaxFunctionSize int               `json:"max_function_size"`
	Complexity      int               `json:"complexity"`
	CommentRatio    float64           `json:"comment_ratio"`
	TokenHashes     map[string]string `json:"token_hashes,omitempty"` // コードクローン検出用
}

// ProjectMetrics はプロジェクト全体のメトリクス情報を保持します
type ProjectMetrics struct {
	ProjectPath     string                 `json:"project_path"`
	AnalyzedAt      time.Time              `json:"analyzed_at"`
	FileCount       int                    `json:"file_count"`
	TotalLines      int                    `json:"total_lines"`
	TotalCodeLines  int                    `json:"total_code_lines"`
	TotalComments   int                    `json:"total_comments"`
	TotalBlankLines int                    `json:"total_blank_lines"`
	AvgComplexity   float64                `json:"avg_complexity"`
	MaxComplexity   int                    `json:"max_complexity"`
	CommentRatio    float64                `json:"comment_ratio"`
	Files           []FileMetrics          `json:"files"`
	Clones          []CodeClone            `json:"clones,omitempty"`
	Trends          map[string]TrendMetric `json:"trends,omitempty"`
}

// CodeClone はコードクローンの情報を保持します
type CodeClone struct {
	SourceFile string  `json:"source_file"`
	TargetFile string  `json:"target_file"`
	SourceLine int     `json:"source_line"`
	TargetLine int     `json:"target_line"`
	LineCount  int     `json:"line_count"`
	Similarity float64 `json:"similarity"`
	Code       string  `json:"code"`
}

// TrendMetric は時系列分析のための指標を保持します
type TrendMetric struct {
	Current    float64 `json:"current"`
	Previous   float64 `json:"previous"`
	Change     float64 `json:"change"`
	ChangeRate float64 `json:"change_rate"`
}

// HistoricalData は過去の分析結果を保持します
type HistoricalData struct {
	Date           time.Time `json:"date"`
	TotalLines     int       `json:"total_lines"`
	CodeLines      int       `json:"code_lines"`
	CommentLines   int       `json:"comment_lines"`
	BlankLines     int       `json:"blank_lines"`
	FileCount      int       `json:"file_count"`
	AvgComplexity  float64   `json:"avg_complexity"`
	MaxComplexity  int       `json:"max_complexity"`
	CommentRatio   float64   `json:"comment_ratio"`
	CloneCount     int       `json:"clone_count"`
	CloneLineRatio float64   `json:"clone_line_ratio"`
}
