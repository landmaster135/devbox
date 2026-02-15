package jsonfmt

import (
	"encoding/json"
	"fmt"

	"github.com/landmaster135/devbox/internal/data_converter/usecases/formats/common"
)

// Parse は JSON を key-value リストへ正規化します。
func Parse(content []byte) (*common.NormalizedData, error) {
	return common.ParseStructuredToKeyValueList(content, json.Unmarshal, "JSON")
}

// Serialize は key-value リストを JSON へ変換します。
func Serialize(records []map[string]string) ([]byte, error) {
	b, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("JSON生成に失敗しました: %w", err)
	}
	return append(b, '\n'), nil
}
