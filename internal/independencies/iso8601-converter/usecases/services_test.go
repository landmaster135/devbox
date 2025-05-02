package usecases

import (
	"os"
	"testing"
)

var originalTZ string

func TestMain(m *testing.M) {
	// 元のTZ環境変数を保存
	originalTZ = os.Getenv("TZ")

	// テスト実行時にタイムゾーンをUTCに設定
	os.Setenv("TZ", "UTC")

	// テストを実行
	code := m.Run()

	// 元のTZ環境変数を復元
	if originalTZ != "" {
		os.Setenv("TZ", originalTZ)
	} else {
		os.Unsetenv("TZ")
	}

	os.Exit(code)
}

func TestUnixToISO8601_Normal(t *testing.T) {
	// テストケース
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "正常なUNIXタイムスタンプ",
			input:    "1619712000",
			expected: "2021-04-29T16:00:00Z",
		},
		{
			name:     "現在時刻に近いUNIXタイムスタンプ",
			input:    "1714154400", // 2024-04-27T00:00:00Z
			expected: "2024-04-26T18:00:00Z",
		},
	}

	// テストの実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := UnixToISO8601(tc.input)
			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
			if result != tc.expected {
				t.Errorf("期待値: %s, 実際の値: %s", tc.expected, result)
			}
		})
	}
}

func TestUnixToISO8601_Error(t *testing.T) {
	// エラーケース
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "無効な入力（文字列）",
			input: "invalid",
		},
		{
			name:  "空の入力",
			input: "",
		},
	}

	// テストの実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := UnixToISO8601(tc.input)
			if err == nil {
				t.Error("エラーが期待されましたが、エラーは発生しませんでした")
			}
		})
	}
}

func TestISO8601ToUnix_Normal(t *testing.T) {
	// テストケース
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "正常なISO8601形式",
			input:    "2021-04-30T00:00:00Z",
			expected: "1619740800",
		},
		{
			name:     "タイムゾーン付きISO8601形式",
			input:    "2021-04-30T09:00:00+09:00",
			expected: "1619740800",
		},
	}

	// テストの実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ISO8601ToUnix(tc.input)
			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
			if result != tc.expected {
				t.Errorf("期待値: %s, 実際の値: %s", tc.expected, result)
			}
		})
	}
}

func TestISO8601ToUnix_Error(t *testing.T) {
	// エラーケース
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "無効な形式",
			input: "2021-04-30",
		},
		{
			name:  "空の入力",
			input: "",
		},
		{
			name:  "無効な日付",
			input: "2021-13-30T00:00:00Z",
		},
	}

	// テストの実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ISO8601ToUnix(tc.input)
			if err == nil {
				t.Error("エラーが期待されましたが、エラーは発生しませんでした")
			}
		})
	}
}

func TestDateToUnix_Normal(t *testing.T) {
	// テストケース
	testCases := []struct {
		name     string
		input    string
		isJST    bool
		expected string
	}{
		{
			name:     "UTC日付（ハイフン区切り）",
			input:    "2021-04-30",
			isJST:    false,
			expected: "1619740800", // 2021-04-30 00:00:00 UTC
		},
		{
			name:     "JST日付（ハイフン区切り）",
			input:    "2021-04-30",
			isJST:    true,
			expected: "1619708400", // 2021-04-30 00:00:00 JST
		},
		{
			name:     "UTC日付（スラッシュ区切り）",
			input:    "2021/04/30",
			isJST:    false,
			expected: "1619740800", // 2021-04-30 00:00:00 UTC
		},
		{
			name:     "JST日付（スラッシュ区切り）",
			input:    "2021/04/30",
			isJST:    true,
			expected: "1619708400", // 2021-04-30 00:00:00 JST
		},
		{
			name:     "UTC日付（区切りなし）",
			input:    "20210430",
			isJST:    false,
			expected: "1619740800", // 2021-04-30 00:00:00 UTC
		},
		{
			name:     "JST日付（区切りなし）",
			input:    "20210430",
			isJST:    true,
			expected: "1619708400", // 2021-04-30 00:00:00 JST
		},
	}

	// テストの実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := DateToUnix(tc.input, tc.isJST)
			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
			if result != tc.expected {
				t.Errorf("期待値: %s, 実際の値: %s", tc.expected, result)
			}
		})
	}
}

func TestDateToUnix_Error(t *testing.T) {
	// エラーケース
	testCases := []struct {
		name  string
		input string
		isJST bool
	}{
		{
			name:  "無効な日付形式",
			input: "invalid",
			isJST: false,
		},
		{
			name:  "空の入力",
			input: "",
			isJST: false,
		},
		{
			name:  "時刻情報を含む入力",
			input: "2021-04-30T00:00:00Z",
			isJST: false,
		},
		{
			name:  "無効な日付（存在しない日付）",
			input: "2021-02-30",
			isJST: false,
		},
	}

	// テストの実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DateToUnix(tc.input, tc.isJST)
			if err == nil {
				t.Error("エラーが期待されましたが、エラーは発生しませんでした")
			}
		})
	}
}
