package tsvfmt

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/data_converter/usecases/formats/common"
)

// Parse は TSV を key-value リストへ正規化します。
func Parse(content []byte) (*common.NormalizedData, error) {
	reader := csv.NewReader(bytes.NewReader(content))
	reader.Comma = '\t'

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("TSVの解析に失敗しました: %w", err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("TSVが空です")
	}

	headers := make([]string, 0, len(rows[0]))
	headerExists := map[string]struct{}{}
	for i, h := range rows[0] {
		header := strings.TrimSpace(h)
		if header == "" {
			return nil, fmt.Errorf("TSVヘッダーの%d列目が空です", i+1)
		}
		if _, exists := headerExists[header]; exists {
			return nil, fmt.Errorf("TSVヘッダーに重複があります: %s", header)
		}
		headerExists[header] = struct{}{}
		headers = append(headers, header)
	}

	records := make([]map[string]string, 0, len(rows)-1)
	for rowIndex := 1; rowIndex < len(rows); rowIndex++ {
		row := rows[rowIndex]
		if len(row) > len(headers) {
			return nil, fmt.Errorf("TSVの%d行目の列数がヘッダーより多いです", rowIndex+1)
		}

		record := common.BuildRecordFromValues(headers, row)
		records = append(records, record)
	}

	return common.NewNormalizedData(records, headers), nil
}

// Serialize は key-value リストを TSV へ変換します。
func Serialize(records []map[string]string, keys []string) ([]byte, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("TSV出力に必要なキー情報がありません")
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = '\t'

	if err := writer.Write(keys); err != nil {
		return nil, fmt.Errorf("TSVヘッダーの書き込みに失敗しました: %w", err)
	}

	for _, record := range records {
		if err := writer.Write(common.PickValuesByKeys(record, keys)); err != nil {
			return nil, fmt.Errorf("TSV行の書き込みに失敗しました: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("TSV生成に失敗しました: %w", err)
	}

	return buf.Bytes(), nil
}
