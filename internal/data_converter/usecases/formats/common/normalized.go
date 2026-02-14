package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
)

// NormalizedData は全入力形式を統一した key-value リスト表現です。
type NormalizedData struct {
	KeyValueList []map[string]string
	Keys         []string
}

// ParseStructuredToKeyValueList は JSON/YAML などの構造化データを key-value リストへ正規化します。
func ParseStructuredToKeyValueList(
	content []byte,
	unmarshal func([]byte, any) error,
	formatName string,
) (*NormalizedData, error) {
	var raw any
	if err := unmarshal(content, &raw); err != nil {
		return nil, fmt.Errorf("%sの解析に失敗しました: %w", formatName, err)
	}

	recordsAny, err := parseRecordObjects(raw)
	if err != nil {
		return nil, fmt.Errorf("%sの構造が不正です: %w", formatName, err)
	}

	return newNormalizedDataFromAnyRecords(recordsAny), nil
}

// NewNormalizedData は正規化済みレコードとキー一覧から NormalizedData を生成します。
func NewNormalizedData(records []map[string]string, keys []string) *NormalizedData {
	return &NormalizedData{
		KeyValueList: records,
		Keys:         keys,
	}
}

// BuildRecordFromValues はヘッダー順に値を詰めてレコードを作成します。
func BuildRecordFromValues(headers, values []string) map[string]string {
	record := make(map[string]string, len(headers))
	for i, key := range headers {
		if i < len(values) {
			record[key] = values[i]
			continue
		}
		record[key] = ""
	}
	return record
}

// PickValuesByKeys はキー順の値配列を返します。
func PickValuesByKeys(record map[string]string, keys []string) []string {
	row := make([]string, 0, len(keys))
	for _, key := range keys {
		row = append(row, record[key])
	}
	return row
}

// CollectSortedKeys はレコード群からキーを収集し、レコード内は辞書順で安定化します。
func CollectSortedKeys(records []map[string]string) []string {
	seen := map[string]struct{}{}
	keys := make([]string, 0)

	for _, record := range records {
		recordKeys := make([]string, 0, len(record))
		for key := range record {
			recordKeys = append(recordKeys, key)
		}
		sort.Strings(recordKeys)

		for _, key := range recordKeys {
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			keys = append(keys, key)
		}
	}

	return keys
}

func newNormalizedDataFromAnyRecords(recordsAny []map[string]any) *NormalizedData {
	records := normalizeRecords(recordsAny)
	return &NormalizedData{
		KeyValueList: records,
		Keys:         CollectSortedKeys(records),
	}
}

func parseRecordObjects(raw any) ([]map[string]any, error) {
	switch value := raw.(type) {
	case []any:
		records := make([]map[string]any, 0, len(value))
		for i, elem := range value {
			record, ok := toStringKeyMap(elem)
			if !ok {
				return nil, fmt.Errorf("%d件目がオブジェクトではありません", i+1)
			}
			records = append(records, record)
		}
		return records, nil
	case map[string]any, map[any]any:
		record, ok := toStringKeyMap(value)
		if !ok {
			return nil, errors.New("オブジェクト形式として解釈できません")
		}
		return []map[string]any{record}, nil
	default:
		return nil, errors.New("配列またはオブジェクト形式を指定してください")
	}
}

func toStringKeyMap(value any) (map[string]any, bool) {
	switch v := value.(type) {
	case map[string]any:
		converted := make(map[string]any, len(v))
		for key, val := range v {
			converted[key] = normalizeAnyValue(val)
		}
		return converted, true
	case map[any]any:
		converted := make(map[string]any, len(v))
		for key, val := range v {
			converted[fmt.Sprint(key)] = normalizeAnyValue(val)
		}
		return converted, true
	default:
		return nil, false
	}
}

func normalizeAnyValue(value any) any {
	switch v := value.(type) {
	case map[string]any:
		converted := make(map[string]any, len(v))
		for key, child := range v {
			converted[key] = normalizeAnyValue(child)
		}
		return converted
	case map[any]any:
		converted := make(map[string]any, len(v))
		for key, child := range v {
			converted[fmt.Sprint(key)] = normalizeAnyValue(child)
		}
		return converted
	case []any:
		items := make([]any, 0, len(v))
		for _, item := range v {
			items = append(items, normalizeAnyValue(item))
		}
		return items
	default:
		return v
	}
}

func normalizeRecords(recordsAny []map[string]any) []map[string]string {
	records := make([]map[string]string, 0, len(recordsAny))
	for _, recordAny := range recordsAny {
		record := make(map[string]string, len(recordAny))
		for key, val := range recordAny {
			record[key] = stringifyValue(val)
		}
		records = append(records, record)
	}
	return records
}

func stringifyValue(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return fmt.Sprint(v)
	default:
		b, err := json.Marshal(v)
		if err == nil {
			return string(b)
		}
		return fmt.Sprint(v)
	}
}
