package models

import (
	"encoding/json"
	"fmt"
)

// FileContent はファイルの内容を表す構造体です
type FileContent struct {
	Lines []string
}

// NewFileContent は新しいFileContentインスタンスを作成します
func NewFileContent(lines []string) *FileContent {
	return &FileContent{
		Lines: lines,
	}
}

// RemoveDuplicateLines は指定された範囲の文字列が一致する行を削除します
func (fc *FileContent) RemoveDuplicateLines(startPos, endPos int) (int, error) {
	// 特殊なケース：開始位置が負の場合または終了位置が開始位置より小さい場合は何も削除しない
	if startPos < 0 {
		return 0, fmt.Errorf("開始位置が負の値です: %d", startPos)
	}

	if endPos <= startPos {
		return 0, fmt.Errorf("終了位置が開始位置以下です: 開始位置=%d, 終了位置=%d", startPos, endPos)
	}

	// 各行の指定された範囲の文字列を取得し、一致する行を削除する
	substrings := make(map[string]bool)
	var uniqueLines []string
	removedCount := 0

	for _, line := range fc.Lines {
		// 行の長さを確認
		lineLen := len(line)

		// 開始位置が行の長さを超えている場合は、その行をそのまま保持
		if startPos >= lineLen {
			uniqueLines = append(uniqueLines, line)
			continue
		}

		// 終了位置を調整（行の長さを超えないようにする）
		actualEndPos := endPos
		if endPos > lineLen {
			actualEndPos = lineLen
		}

		// 指定された範囲の文字列を取得
		substring := line[startPos:actualEndPos]

		// この部分文字列が既に見つかっていない場合は、行を保持
		if !substrings[substring] {
			substrings[substring] = true
			uniqueLines = append(uniqueLines, line)
		} else {
			// 既に見つかっている場合は、行を削除（カウントを増やす）
			removedCount++
		}
	}

	// 結果を更新
	fc.Lines = uniqueLines

	return removedCount, nil
}

// RequestBody はAPIリクエスト用のリクエストボディを表す構造体です
type RequestBody struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Data        map[string]interface{} `json:"data"`
}

// NewRequestBody は新しいRequestBodyインスタンスを作成します
func NewRequestBody(keyName string, jsonData []interface{}) *RequestBody {
	return &RequestBody{
		Name:        "manual_request",
		Description: "By manually bulk request",
		Data: map[string]interface{}{
			keyName: jsonData,
		},
	}
}

// ToJSON はRequestBodyをJSON形式に変換します
func (rb *RequestBody) ToJSON() ([]byte, error) {
	return json.MarshalIndent(rb, "", "  ")
}
