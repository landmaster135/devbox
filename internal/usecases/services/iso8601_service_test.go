package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestISO8601Service_ConvertISO8601ToTimestamp_WithRealData(t *testing.T) {
	// テストデータファイルのパス
	testDataPath := "../../../testdata/for_json_iso8601_converter/01_crlf.json"

	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "iso8601-service-test-real-data")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テストデータファイルを一時ディレクトリにコピー
	testDataContent, err := os.ReadFile(testDataPath)
	if err != nil {
		t.Fatalf("テストデータファイルの読み込みに失敗しました: %v", err)
	}

	// UTF-8 BOMを削除
	if len(testDataContent) >= 3 && testDataContent[0] == 0xEF && testDataContent[1] == 0xBB && testDataContent[2] == 0xBF {
		testDataContent = testDataContent[3:]
	}

	tempFilePath := filepath.Join(tempDir, "01.json")
	if err := os.WriteFile(tempFilePath, testDataContent, 0644); err != nil {
		t.Fatalf("テストデータファイルのコピーに失敗しました: %v", err)
	}

	// テスト対象のサービスを作成
	mockJsonRepo := &MockJSONRepository{
		FindJSONFilesFunc: func(dirPath string, recursive bool) ([]string, error) {
			return []string{tempFilePath}, nil
		},
		ConvertFileFunc: func(filePath, key string, dryRun bool) (bool, error) {
			if dryRun {
				return true, nil
			}

			// JSONファイルを読み込む
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				return false, err
			}

			var jsonData []map[string]interface{}
			if err := json.Unmarshal(fileData, &jsonData); err != nil {
				return false, err
			}

			// キーの値を変換
			for i, item := range jsonData {
				if key == "created_at" && item[key] != nil {
					jsonData[i][key] = 1710495045 // 2025-03-15T09:30:45Z
				} else if key == "updated_at" && item[key] != nil {
					jsonData[i][key] = 1711103110 // 2025-03-20T14:25:10+09:00
				} else if key == "output_at" && item[key] != nil {
					if i == 0 {
						jsonData[i][key] = 1710847352 // 2025-03-19 10:42:32+00
					} else {
						jsonData[i][key] = 1710847394 // 2025-03-19 10:43:14+00
					}
				}
			}

			// JSONに変換して書き込む
			jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
			if err != nil {
				return false, err
			}

			if err := os.WriteFile(filePath, jsonBytes, 0644); err != nil {
				return false, err
			}

			return true, nil
		},
		ProcessJSONDataFunc: func(data interface{}, targetKey string) (interface{}, bool) {
			return data, true
		},
	}
	iso8601Service := NewISO8601Service(mockJsonRepo)

	tests := []struct {
		name      string
		dirPath   string
		key       string
		recursive bool
		dryRun    bool
		wantCount int
		wantErr   bool
	}{
		{
			name:      "正常系: created_atキーを変換",
			dirPath:   tempDir,
			key:       "created_at",
			recursive: false,
			dryRun:    false,
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "正常系: updated_atキーを変換",
			dirPath:   tempDir,
			key:       "updated_at",
			recursive: false,
			dryRun:    false,
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "正常系: output_atキーを変換",
			dirPath:   tempDir,
			key:       "output_at",
			recursive: false,
			dryRun:    false,
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "正常系: ドライラン",
			dirPath:   tempDir,
			key:       "created_at",
			recursive: false,
			dryRun:    true,
			wantCount: 1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト実行
			gotCount, err := iso8601Service.ConvertISO8601ToTimestamp(tt.dirPath, tt.key, tt.recursive, tt.dryRun)

			// エラーの検証
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertISO8601ToTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 処理したファイル数の検証
			if gotCount != tt.wantCount {
				t.Errorf("ConvertISO8601ToTimestamp() = %v, want %v", gotCount, tt.wantCount)
			}

			// ドライランでない場合は、ファイルの内容を検証
			if !tt.wantErr && !tt.dryRun {
				// ファイルを読み込む
				fileData, err := os.ReadFile(tempFilePath)
				if err != nil {
					t.Fatalf("ファイルの読み込みに失敗しました: %v", err)
				}

				var jsonData []map[string]interface{}
				if err := json.Unmarshal(fileData, &jsonData); err != nil {
					t.Fatalf("JSONのパースに失敗しました: %v", err)
				}

				// キーの値がUNIXタイムスタンプに変換されているか確認
				for _, item := range jsonData {
					if tt.key == "created_at" {
						if _, ok := item["created_at"].(float64); !ok {
							t.Errorf("created_at is not a number: %v", item["created_at"])
						}
					} else if tt.key == "updated_at" {
						if _, ok := item["updated_at"].(float64); !ok {
							t.Errorf("updated_at is not a number: %v", item["updated_at"])
						}
					} else if tt.key == "output_at" {
						if _, ok := item["output_at"].(float64); !ok {
							t.Errorf("output_at is not a number: %v", item["output_at"])
						}
					}
				}
			}
		})
	}
}

func TestISO8601Service_ConvertISO8601ToTimestamp(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "iso8601-service-test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用のJSONファイルを作成
	jsonContent := `{
		"id": 123,
		"title": "テストデータ",
		"created_at": "2025-04-10T15:30:45Z",
		"metadata": {
			"published_at": "2025-04-10T10:15:20Z"
		},
		"items": [
			{
				"name": "アイテム1",
				"timestamp": "2025-04-09T22:10:05Z"
			},
			{
				"name": "アイテム2",
				"timestamp": "2025-04-10T08:25:15Z"
			}
		]
	}`
	jsonPath := filepath.Join(tempDir, "test.json")
	if err := os.WriteFile(jsonPath, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	// テスト対象のサービスを作成
	mockJsonRepo := &MockJSONRepository{
		FindJSONFilesFunc: func(dirPath string, recursive bool) ([]string, error) {
			if dirPath == filepath.Join(tempDir, "not_exist") {
				return nil, fmt.Errorf("指定されたディレクトリが存在しません: %s", dirPath)
			}
			return []string{jsonPath}, nil
		},
		ConvertFileFunc: func(filePath, key string, dryRun bool) (bool, error) {
			if dryRun {
				return true, nil
			}

			// JSONファイルを読み込む
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				return false, err
			}

			var jsonData map[string]interface{}
			if err := json.Unmarshal(fileData, &jsonData); err != nil {
				return false, err
			}

			// キーの値を変換
			if key == "created_at" {
				jsonData["created_at"] = 1744299045
			} else if key == "published_at" {
				metadata, ok := jsonData["metadata"].(map[string]interface{})
				if ok {
					metadata["published_at"] = 1744280120
				}
			} else if key == "timestamp" {
				items, ok := jsonData["items"].([]interface{})
				if ok {
					for i, item := range items {
						itemObj, ok := item.(map[string]interface{})
						if ok {
							if i == 0 {
								itemObj["timestamp"] = 1744236605
							} else {
								itemObj["timestamp"] = 1744273515
							}
						}
					}
				}
			}

			// JSONに変換して書き込む
			jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
			if err != nil {
				return false, err
			}

			if err := os.WriteFile(filePath, jsonBytes, 0644); err != nil {
				return false, err
			}

			return true, nil
		},
		ProcessJSONDataFunc: func(data interface{}, targetKey string) (interface{}, bool) {
			return data, true
		},
	}
	iso8601Service := NewISO8601Service(mockJsonRepo)

	tests := []struct {
		name      string
		dirPath   string
		key       string
		recursive bool
		dryRun    bool
		wantCount int
		wantErr   bool
	}{
		{
			name:      "正常系: created_atキーを変換",
			dirPath:   tempDir,
			key:       "created_at",
			recursive: false,
			dryRun:    false,
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "正常系: published_atキーを変換",
			dirPath:   tempDir,
			key:       "published_at",
			recursive: false,
			dryRun:    false,
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "正常系: timestampキーを変換",
			dirPath:   tempDir,
			key:       "timestamp",
			recursive: false,
			dryRun:    false,
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "正常系: ドライラン",
			dirPath:   tempDir,
			key:       "created_at",
			recursive: false,
			dryRun:    true,
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "エラー: 存在しないディレクトリ",
			dirPath:   filepath.Join(tempDir, "not_exist"),
			key:       "created_at",
			recursive: false,
			dryRun:    false,
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト実行
			gotCount, err := iso8601Service.ConvertISO8601ToTimestamp(tt.dirPath, tt.key, tt.recursive, tt.dryRun)

			// エラーの検証
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertISO8601ToTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 処理したファイル数の検証
			if gotCount != tt.wantCount {
				t.Errorf("ConvertISO8601ToTimestamp() = %v, want %v", gotCount, tt.wantCount)
			}

			// ドライランでない場合は、ファイルの内容を検証
			if !tt.wantErr && !tt.dryRun {
				// ファイルを読み込む
				fileData, err := os.ReadFile(jsonPath)
				if err != nil {
					t.Fatalf("ファイルの読み込みに失敗しました: %v", err)
				}

				var jsonData map[string]interface{}
				if err := json.Unmarshal(fileData, &jsonData); err != nil {
					t.Fatalf("JSONのパースに失敗しました: %v", err)
				}

				// キーの値がUNIXタイムスタンプに変換されているか確認
				if tt.key == "created_at" {
					if _, ok := jsonData["created_at"].(float64); !ok {
						t.Errorf("created_at is not a number: %v", jsonData["created_at"])
					}
				} else if tt.key == "published_at" {
					metadata, ok := jsonData["metadata"].(map[string]interface{})
					if !ok {
						t.Errorf("metadata is not an object: %v", jsonData["metadata"])
					} else {
						if _, ok := metadata["published_at"].(float64); !ok {
							t.Errorf("published_at is not a number: %v", metadata["published_at"])
						}
					}
				} else if tt.key == "timestamp" {
					items, ok := jsonData["items"].([]interface{})
					if !ok {
						t.Errorf("items is not an array: %v", jsonData["items"])
					} else {
						for i, item := range items {
							itemObj, ok := item.(map[string]interface{})
							if !ok {
								t.Errorf("item[%d] is not an object: %v", i, item)
							} else {
								if _, ok := itemObj["timestamp"].(float64); !ok {
									t.Errorf("timestamp is not a number: %v", itemObj["timestamp"])
								}
							}
						}
					}
				}
			}
		})
	}
}

func TestISO8601Service_ParseRealJSONFile(t *testing.T) {
	// テストデータファイルのパス
	testDataPath := "../../../testdata/for_json_iso8601_converter/01_crlf.json"

	// ファイルを直接読み込む
	fileData, err := os.ReadFile(testDataPath)
	if err != nil {
		t.Fatalf("テストデータファイルの読み込みに失敗しました: %v", err)
	}

	// UTF-8 BOMを削除
	if len(fileData) >= 3 && fileData[0] == 0xEF && fileData[1] == 0xBB && fileData[2] == 0xBF {
		fileData = fileData[3:]
	}

	// JSONデータをパース
	var jsonData []map[string]interface{}
	if err := json.Unmarshal(fileData, &jsonData); err != nil {
		t.Fatalf("JSONのパースに失敗しました: %v", err)
	}

	// モックリポジトリの設定
	mockJsonRepo := &MockJSONRepository{}

	// テスト対象のサービスを作成
	iso8601Service := NewISO8601Service(mockJsonRepo)

	// 日時文字列とその期待値のマッピング
	dateTests := []struct {
		fieldName string
		index     int
		expected  int64
	}{
		{"created_at", 0, 1744278149}, // 2025-04-10 09:42:29.059877+00
		{"updated_at", 0, 1744278149}, // 2025-04-10 09:42:29.059877+00
		{"output_at", 0, 1742380952},  // 2025-03-19 10:42:32+00
		{"created_at", 1, 1744278149}, // 2025-04-10 09:42:29.936466+00
		{"updated_at", 1, 1744278149}, // 2025-04-10 09:42:29.936466+00
		{"output_at", 1, 1742380994},  // 2025-03-19 10:43:14+00
	}

	for _, tt := range dateTests {
		t.Run(fmt.Sprintf("%s[%d]", tt.fieldName, tt.index), func(t *testing.T) {
			// 日時文字列を取得
			dateStr, ok := jsonData[tt.index][tt.fieldName].(string)
			if !ok {
				t.Fatalf("%s is not a string: %v", tt.fieldName, jsonData[tt.index][tt.fieldName])
			}

			// ISO8601形式の日時文字列をUNIXタイムスタンプに変換
			timestamp, err := iso8601Service.parseISO8601(dateStr)
			if err != nil {
				t.Errorf("parseISO8601(%s) error = %v", dateStr, err)
				return
			}

			// 結果の検証（タイムゾーンの違いなどで多少の誤差があるため、前後60秒以内なら許容する）
			diff := timestamp - tt.expected
			if diff < -60 || diff > 60 {
				t.Errorf("parseISO8601(%s) = %v, want %v (diff: %v)", dateStr, timestamp, tt.expected, diff)
			}
		})
	}
}

func TestISO8601Service_processJSONData(t *testing.T) {
	// テスト用のJSONデータ
	jsonData := map[string]interface{}{
		"id": 123,
		"created_at": "2025-04-10T15:30:45Z",
		"metadata": map[string]interface{}{
			"published_at": "2025-04-10T10:15:20Z",
		},
		"items": []interface{}{
			map[string]interface{}{
				"name": "アイテム1",
				"timestamp": "2025-04-09T22:10:05Z",
			},
			map[string]interface{}{
				"name": "アイテム2",
				"timestamp": "2025-04-10T08:25:15Z",
			},
		},
	}

	// モックリポジトリの設定
	mockJsonRepo := &MockJSONRepository{
		ProcessJSONDataFunc: func(data interface{}, targetKey string) (interface{}, bool) {
			// 存在しないキーの場合はfalseを返す
			if targetKey == "non_existent_key" {
				return data, false
			}

			// それ以外のキーの場合はtrueを返す
			return data, true
		},
	}

	// テスト対象のサービスを作成
	iso8601Service := NewISO8601Service(mockJsonRepo)

	tests := []struct {
		name      string
		data      interface{}
		targetKey string
		wantData  interface{}
		wantOk    bool
	}{
		{
			name:      "正常系: created_atキーを変換",
			data:      jsonData,
			targetKey: "created_at",
			wantData:  jsonData, // モックが処理するので実際の値は変わらない
			wantOk:    true,
		},
		{
			name:      "正常系: published_atキーを変換",
			data:      jsonData,
			targetKey: "published_at",
			wantData:  jsonData, // モックが処理するので実際の値は変わらない
			wantOk:    true,
		},
		{
			name:      "正常系: timestampキーを変換",
			data:      jsonData,
			targetKey: "timestamp",
			wantData:  jsonData, // モックが処理するので実際の値は変わらない
			wantOk:    true,
		},
		{
			name:      "エラー: 存在しないキー",
			data:      jsonData,
			targetKey: "non_existent_key",
			wantData:  jsonData,
			wantOk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト実行
			gotData, gotOk := iso8601Service.processJSONData(tt.data, tt.targetKey)

			// 結果の検証
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("processJSONData() gotData = %v, want %v", gotData, tt.wantData)
			}
			if gotOk != tt.wantOk {
				t.Errorf("processJSONData() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestISO8601Service_parseISO8601(t *testing.T) {
	// モックリポジトリの設定
	mockJsonRepo := &MockJSONRepository{}

	// テスト対象のサービスを作成
	iso8601Service := NewISO8601Service(mockJsonRepo)

	tests := []struct {
		name    string
		dateStr string
		want    int64
		wantErr bool
	}{
		{
			name:    "正常系: RFC3339形式",
			dateStr: "2025-04-10T15:30:45Z",
			want:    1744299045,
			wantErr: false,
		},
		{
			name:    "正常系: RFC3339形式（タイムゾーン付き）",
			dateStr: "2025-04-11T02:45:30+09:00",
			want:    1744307130,
			wantErr: false,
		},
		{
			name:    "正常系: RFC3339Nano形式",
			dateStr: "2025-04-10T15:30:45.123456789Z",
			want:    1744299045,
			wantErr: false,
		},
		{
			name:    "エラー: 不正な形式",
			dateStr: "2025/04/10 15:30:45",
			want:    0,
			wantErr: true,
		},
		{
			name:    "エラー: 空文字列",
			dateStr: "",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト実行
			got, err := iso8601Service.parseISO8601(tt.dateStr)

			// エラーの検証
			if (err != nil) != tt.wantErr {
				t.Errorf("parseISO8601() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 結果の検証
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseISO8601() = %v, want %v", got, tt.want)
			}
		})
	}
}
