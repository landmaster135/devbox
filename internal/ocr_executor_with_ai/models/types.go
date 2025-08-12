package models

import "fmt"

// OcrResult は単一画像のOCR結果を保持する構造体
type OcrResult struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
	Error    string `json:"error,omitempty"`
}

// OcrExecutionResult はOCR実行結果全体を保持する構造体
type OcrExecutionResult struct {
	Results []OcrResult `json:"results"`
	Total   int         `json:"total"`
	Success int         `json:"success"`
	Failed  int         `json:"failed"`
}

// AddResult は結果を追加する
func (r *OcrExecutionResult) AddResult(result OcrResult) {
	r.Results = append(r.Results, result)
	r.Total++
	if result.Error == "" {
		r.Success++
	} else {
		r.Failed++
	}
}

// FormatAsText は結果をテキスト形式で出力する
func (r *OcrExecutionResult) FormatAsText() string {
	if r.Total == 0 {
		return "処理対象の画像ファイルが見つかりませんでした。"
	}

	var output string
	output += "=== AI OCR実行結果 ===\n"
	output += fmt.Sprintf("処理総数: %d件 (成功: %d件, 失敗: %d件)\n\n", r.Total, r.Success, r.Failed)

	for i, result := range r.Results {
		output += fmt.Sprintf("[%d] %s\n", i+1, result.FilePath)
		if result.Error != "" {
			output += fmt.Sprintf("エラー: %s\n", result.Error)
		} else {
			output += fmt.Sprintf("OCR結果:\n%s\n", result.Content)
		}
		output += "\n"
	}

	return output
}
