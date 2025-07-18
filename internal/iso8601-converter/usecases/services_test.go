package usecases

import (
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
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
		isJST    bool
		expected string
	}{
		{
			name:     "正常なISO8601形式（UTC）",
			input:    "2021-04-30T00:00:00Z",
			isJST:    false,
			expected: "1619740800",
		},
		{
			name:     "タイムゾーン付きISO8601形式",
			input:    "2021-04-30T09:00:00+09:00",
			isJST:    false,
			expected: "1619740800",
		},
		{
			name:     "タイムゾーンなしISO8601形式（UTC扱い）",
			input:    "2021-04-30T00:00:00",
			isJST:    false,
			expected: "1619740800",
		},
		{
			name:     "タイムゾーンなしISO8601形式（JST扱い）",
			input:    "2021-04-30T00:00:00",
			isJST:    true,
			expected: "1619708400",
		},
	}

	// テストの実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ISO8601ToUnix(tc.input, tc.isJST)
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
		isJST bool
	}{
		{
			name:  "空の入力",
			input: "",
			isJST: false,
		},
		{
			name:  "無効な日付",
			input: "2021-13-30T00:00:00Z",
			isJST: false,
		},
	}

	// テストの実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ISO8601ToUnix(tc.input, tc.isJST)
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

func TestISO8601ToUnix_AdditionalFormats(t *testing.T) {
	// 追加のISO8601フォーマットテスト
	testCases := []struct {
		name     string
		input    string
		isJST    bool
		expected string
	}{
		{
			name:     "ミリ秒付きISO8601形式（タイムゾーンなし、UTC扱い）",
			input:    "2021-04-30T00:00:00.000",
			isJST:    false,
			expected: "1619740800",
		},
		{
			name:     "ミリ秒付きISO8601形式（タイムゾーンなし、JST扱い）",
			input:    "2021-04-30T00:00:00.000",
			isJST:    true,
			expected: "1619708400",
		},
		{
			name:     "マイクロ秒付きISO8601形式（タイムゾーンなし、UTC扱い）",
			input:    "2021-04-30T00:00:00.000000",
			isJST:    false,
			expected: "1619740800",
		},
		{
			name:     "マイクロ秒付きISO8601形式（タイムゾーンなし、JST扱い）",
			input:    "2021-04-30T00:00:00.000000",
			isJST:    true,
			expected: "1619708400",
		},
		{
			name:     "ナノ秒付きISO8601形式（タイムゾーンなし、UTC扱い）",
			input:    "2021-04-30T00:00:00.000000000",
			isJST:    false,
			expected: "1619740800",
		},
		{
			name:     "ナノ秒付きISO8601形式（タイムゾーンなし、JST扱い）",
			input:    "2021-04-30T00:00:00.000000000",
			isJST:    true,
			expected: "1619708400",
		},
		{
			name:     "RFC3339Nano形式（タイムゾーン付き）",
			input:    "2021-04-30T09:00:00.123456789+09:00",
			isJST:    false,
			expected: "1619740800",
		},
	}

	// テストの実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ISO8601ToUnix(tc.input, tc.isJST)
			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
			if result != tc.expected {
				t.Errorf("期待値: %s, 実際の値: %s", tc.expected, result)
			}
		})
	}
}

func TestNowToUnix_Normal(t *testing.T) {
	// 現在時刻のUnixタイムスタンプを取得
	result := NowToUnix()

	// 結果が空でないことを確認
	if result == "" {
		t.Error("結果が空です")
	}

	// 結果が数値文字列であることを確認
	_, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		t.Errorf("結果が有効な数値文字列ではありません: %s, エラー: %v", result, err)
	}

	// 現在時刻に近い値であることを確認（±10秒以内）
	currentUnix := time.Now().Unix()
	resultUnix, _ := strconv.ParseInt(result, 10, 64)
	diff := currentUnix - resultUnix
	if diff < -10 || diff > 10 {
		t.Errorf("現在時刻との差が大きすぎます。差: %d秒", diff)
	}
}

func TestNowToISO8601InUTC_Normal(t *testing.T) {
	// 現在時刻のISO8601形式（UTC）を取得
	result := NowToISO8601InUTC()

	// 結果が空でないことを確認
	if result == "" {
		t.Error("結果が空です")
	}

	// 結果がISO8601形式であることを確認
	_, err := time.Parse(time.RFC3339, result)
	if err != nil {
		t.Errorf("結果が有効なISO8601形式ではありません: %s, エラー: %v", result, err)
	}

	// UTCタイムゾーンであることを確認（Zで終わる）
	if !strings.HasSuffix(result, "Z") {
		t.Errorf("UTC形式ではありません（Zで終わっていません）: %s", result)
	}

	// 現在時刻に近い値であることを確認
	parsedTime, _ := time.Parse(time.RFC3339, result)
	currentTime := time.Now().UTC()
	diff := currentTime.Sub(parsedTime).Seconds()
	if diff < -10 || diff > 10 {
		t.Errorf("現在時刻との差が大きすぎます。差: %.2f秒", diff)
	}
}

func TestNowToISO8601InJST_Normal(t *testing.T) {
	// 現在時刻のISO8601形式（JST）を取得
	result := NowToISO8601InJST()

	// 結果が空でないことを確認
	if result == "" {
		t.Error("結果が空です")
	}

	// 結果がISO8601形式であることを確認
	_, err := time.Parse(time.RFC3339, result)
	if err != nil {
		t.Errorf("結果が有効なISO8601形式ではありません: %s, エラー: %v", result, err)
	}

	// JSTタイムゾーンであることを確認（+09:00で終わる）
	if !strings.HasSuffix(result, "+09:00") {
		t.Errorf("JST形式ではありません（+09:00で終わっていません）: %s", result)
	}

	// 現在時刻に近い値であることを確認
	parsedTime, _ := time.Parse(time.RFC3339, result)
	jst := time.FixedZone("JST", 9*60*60)
	currentTime := time.Now().In(jst)
	diff := currentTime.Sub(parsedTime).Seconds()
	if diff < -10 || diff > 10 {
		t.Errorf("現在時刻との差が大きすぎます。差: %.2f秒", diff)
	}
}

func TestUnixToISO8601_EdgeCases(t *testing.T) {
	// 境界値テスト
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Unix epoch（1970-01-01 00:00:00 UTC）",
			input:    "0",
			expected: "1970-01-01T00:00:00Z",
		},
		{
			name:     "負の値（1969年）",
			input:    "-86400", // 1969-12-31 00:00:00 UTC
			expected: "1969-12-31T00:00:00Z",
		},
		{
			name:     "大きな値（2038年問題前）",
			input:    "2147483647", // 2038-01-19T03:14:07Z
			expected: "2038-01-19T03:14:07Z",
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

func TestUnixToISO8601_AdditionalErrors(t *testing.T) {
	// 追加のエラーケース
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "浮動小数点数",
			input: "123.456",
		},
		{
			name:  "16進数",
			input: "0x123",
		},
		{
			name:  "非常に大きな数値（オーバーフロー）",
			input: strings.Repeat("9", 1000), // 1000桁の9
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
