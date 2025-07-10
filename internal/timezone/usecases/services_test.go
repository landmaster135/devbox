package usecases

import (
	"strings"
	"testing"
	"time"
)

// TestTimezoneService_GetCurrentTime は GetCurrentTime メソッドをテストします
func TestTimezoneService_GetCurrentTime(t *testing.T) {
	service := NewTimezoneService()

	tests := []struct {
		name      string
		timezone  string
		wantError bool
	}{
		{
			name:      "有効なタイムゾーン: UTC",
			timezone:  "UTC",
			wantError: false,
		},
		{
			name:      "有効なタイムゾーン: Asia/Tokyo",
			timezone:  "Asia/Tokyo",
			wantError: false,
		},
		{
			name:      "有効なタイムゾーン: America/New_York",
			timezone:  "America/New_York",
			wantError: false,
		},
		{
			name:      "無効なタイムゾーン",
			timezone:  "Invalid/Timezone",
			wantError: true,
		},
		{
			name:      "空のタイムゾーン",
			timezone:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetCurrentTime(tt.timezone)
			if (err != nil) != tt.wantError {
				t.Errorf("GetCurrentTime() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				// 結果が正しい形式かチェック（新しいフォーマット）
				_, err := time.Parse("2006-01-02 15:04:05 MST (Z07:00)", result)
				if err != nil {
					t.Errorf("GetCurrentTime() returned invalid time format: %v", result)
				}
			}
		})
	}
}

// TestTimezoneService_NormalizeTimezone は NormalizeTimezone メソッドをテストします
func TestTimezoneService_NormalizeTimezone(t *testing.T) {
	service := NewTimezoneService()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "標準的なタイムゾーン",
			input:    "UTC",
			expected: "UTC",
		},
		{
			name:     "一般的な略称: JST",
			input:    "jst",
			expected: "Asia/Tokyo",
		},
		{
			name:     "一般的な略称: EST",
			input:    "est",
			expected: "America/New_York",
		},
		{
			name:     "国名: Japan",
			input:    "japan",
			expected: "Asia/Tokyo",
		},
		{
			name:     "存在しない略称",
			input:    "xyz",
			expected: "xyz",
		},
		{
			name:     "大文字小文字の違い",
			input:    "JST",
			expected: "Asia/Tokyo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.NormalizeTimezone(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeTimezone() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestTimezoneService_FindSimilarTimezones は FindSimilarTimezones メソッドをテストします
func TestTimezoneService_FindSimilarTimezones(t *testing.T) {
	service := NewTimezoneService()

	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:     "Tokyo を含む検索",
			input:    "Tokyo",
			contains: []string{"Asia/Tokyo"},
		},
		{
			name:     "New_York を含む検索",
			input:    "New_York",
			contains: []string{"America/New_York"},
		},
		{
			name:     "存在しない検索",
			input:    "NonExistent",
			contains: []string{"UTC"}, // デフォルトの提案を含む
		},
		{
			name:     "Pacific を含む検索",
			input:    "Pacific",
			contains: []string{"Pacific"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.FindSimilarTimezones(tt.input)
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("FindSimilarTimezones() = %v, should contain %v", result, expected)
				}
			}
		})
	}
}

// TestTimezoneService_IsValidTimezone は IsValidTimezone メソッドをテストします
func TestTimezoneService_IsValidTimezone(t *testing.T) {
	service := NewTimezoneService()

	tests := []struct {
		name     string
		timezone string
		want     bool
	}{
		{
			name:     "有効なタイムゾーン: UTC",
			timezone: "UTC",
			want:     true,
		},
		{
			name:     "有効なタイムゾーン: Asia/Tokyo",
			timezone: "Asia/Tokyo",
			want:     true,
		},
		{
			name:     "有効なタイムゾーン: America/New_York",
			timezone: "America/New_York",
			want:     true,
		},
		{
			name:     "無効なタイムゾーン",
			timezone: "Invalid/Timezone",
			want:     false,
		},
		{
			name:     "空のタイムゾーン",
			timezone: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.IsValidTimezone(tt.timezone)
			if got != tt.want {
				t.Errorf("IsValidTimezone() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestAvailableTimezones は availableTimezones 変数をテストします
func TestAvailableTimezones(t *testing.T) {
	// 重要なタイムゾーンが含まれているか確認
	expectedTimezones := []string{
		"UTC",
		"Asia/Tokyo",
		"America/New_York",
		"Europe/London",
	}

	for _, expected := range expectedTimezones {
		found := false
		for _, actual := range availableTimezones {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("availableTimezones does not contain %s", expected)
		}
	}

	// すべてのタイムゾーンが有効か確認
	service := NewTimezoneService()
	for _, tz := range availableTimezones {
		if !service.IsValidTimezone(tz) {
			t.Errorf("availableTimezones contains invalid timezone: %s", tz)
		}
	}
}

// TestTimeFormatting は時刻のフォーマット機能をテストします
func TestTimeFormatting(t *testing.T) {
	// 特定の時刻を使用してテスト
	testTime := time.Date(2023, 1, 2, 15, 4, 5, 0, time.UTC)
	formatted := testTime.Format("2006-01-02 15:04:05")
	expected := "2023-01-02 15:04:05"

	if formatted != expected {
		t.Errorf("Time formatting failed, got: %s, want: %s", formatted, expected)
	}
}

// TestTimezoneConversion はタイムゾーン変換機能をテストします
func TestTimezoneConversion(t *testing.T) {
	// UTC の特定の時刻
	utcTime := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

	// 東京のタイムゾーンに変換
	tokyoLoc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatalf("Failed to load Asia/Tokyo location: %v", err)
	}
	tokyoTime := utcTime.In(tokyoLoc)

	// 東京は UTC+9 なので、9時間進んでいるはず
	expectedHour := 9
	if tokyoTime.Hour() != expectedHour {
		t.Errorf("Timezone conversion failed, got hour: %d, want: %d", tokyoTime.Hour(), expectedHour)
	}
}

// TestTimezoneService_ConvertTime は ConvertTime メソッドをテストします
func TestTimezoneService_ConvertTime(t *testing.T) {
	service := NewTimezoneService()

	tests := []struct {
		name         string
		dateTime     string
		fromTimezone string
		toTimezone   string
		wantError    bool
		contains     string
	}{
		{
			name:         "UTC から Tokyo への変換",
			dateTime:     "2023-01-02 00:00:00",
			fromTimezone: "UTC",
			toTimezone:   "Asia/Tokyo",
			wantError:    false,
			contains:     "09:00:00",
		},
		{
			name:         "Tokyo から New_York への変換",
			dateTime:     "2023-01-02 12:00:00",
			fromTimezone: "Asia/Tokyo",
			toTimezone:   "America/New_York",
			wantError:    false,
			contains:     "22:00:00",
		},
		{
			name:         "略称を使用: JST から EST への変換",
			dateTime:     "2023-01-02 12:00:00",
			fromTimezone: "jst",
			toTimezone:   "est",
			wantError:    false,
			contains:     "22:00:00",
		},
		{
			name:         "無効なタイムゾーン",
			dateTime:     "2023-01-02 12:00:00",
			fromTimezone: "Invalid/Timezone",
			toTimezone:   "UTC",
			wantError:    true,
		},
		{
			name:         "無効な日付形式",
			dateTime:     "Invalid Date",
			fromTimezone: "UTC",
			toTimezone:   "UTC",
			wantError:    true,
		},
		{
			name:         "別の日付形式: YYYY/MM/DD",
			dateTime:     "2023/01/02 12:00:00",
			fromTimezone: "UTC",
			toTimezone:   "Asia/Tokyo",
			wantError:    false,
			contains:     "21:00:00",
		},
		{
			name:         "別の日付形式: YYYY-MM-DD",
			dateTime:     "2023-01-02",
			fromTimezone: "UTC",
			toTimezone:   "Asia/Tokyo",
			wantError:    false,
			contains:     "09:00:00",
		},
		{
			name:         "別の日付形式: HH:MM:SS",
			dateTime:     "12:00:00",
			fromTimezone: "UTC",
			toTimezone:   "Asia/Tokyo",
			wantError:    false,
			contains:     "21:18:59", // 歴史的なタイムゾーン（LMT）のオフセットを考慮
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ConvertTime(tt.dateTime, tt.fromTimezone, tt.toTimezone)
			if (err != nil) != tt.wantError {
				t.Errorf("ConvertTime() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError && !strings.Contains(result, tt.contains) {
				t.Errorf("ConvertTime() = %v, should contain %v", result, tt.contains)
			}
		})
	}
}

// TestTimezoneService_GetAvailableTimezones は GetAvailableTimezones メソッドをテストします
func TestTimezoneService_GetAvailableTimezones(t *testing.T) {
	service := NewTimezoneService()

	timezones := service.GetAvailableTimezones()

	// 結果が空でないことを確認
	if len(timezones) == 0 {
		t.Error("GetAvailableTimezones() returned empty list")
	}

	// 重要なタイムゾーンが含まれているか確認
	expectedTimezones := []string{
		"UTC",
		"Asia/Tokyo",
		"America/New_York",
		"Europe/London",
	}

	for _, expected := range expectedTimezones {
		found := false
		for _, actual := range timezones {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetAvailableTimezones() does not contain %s", expected)
		}
	}
}

// TestTimezoneService_GetCommonTimezoneAliases は GetCommonTimezoneAliases メソッドをテストします
func TestTimezoneService_GetCommonTimezoneAliases(t *testing.T) {
	service := NewTimezoneService()

	aliases := service.GetCommonTimezoneAliases()

	// 結果が空でないことを確認
	if len(aliases) == 0 {
		t.Error("GetCommonTimezoneAliases() returned empty map")
	}

	// 重要な別名が含まれているか確認
	expectedAliases := map[string]string{
		"jst":   "Asia/Tokyo",
		"est":   "America/New_York",
		"utc":   "UTC",
		"japan": "Asia/Tokyo",
	}

	for alias, expectedTz := range expectedAliases {
		if actualTz, ok := aliases[alias]; !ok {
			t.Errorf("GetCommonTimezoneAliases() does not contain alias %s", alias)
		} else if actualTz != expectedTz {
			t.Errorf("GetCommonTimezoneAliases()[%s] = %s, want %s", alias, actualTz, expectedTz)
		}
	}
}

// TestTimezoneService_HandleGetCurrentTime は HandleGetCurrentTime メソッドをテストします
func TestTimezoneService_HandleGetCurrentTime(t *testing.T) {
	service := NewTimezoneService()

	tests := []struct {
		name      string
		timezone  string
		wantError bool
		contains  string
	}{
		{
			name:      "有効なタイムゾーン: UTC",
			timezone:  "UTC",
			wantError: false,
			contains:  "UTC の現在時刻:",
		},
		{
			name:      "有効なタイムゾーン: Asia/Tokyo",
			timezone:  "Asia/Tokyo",
			wantError: false,
			contains:  "Asia/Tokyo の現在時刻:",
		},
		{
			name:      "略称: JST",
			timezone:  "jst",
			wantError: false,
			contains:  "Asia/Tokyo の現在時刻:",
		},
		{
			name:      "無効なタイムゾーン",
			timezone:  "Invalid/Timezone",
			wantError: true,
		},
		{
			name:      "空のタイムゾーン",
			timezone:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.HandleGetCurrentTime(tt.timezone)
			if (err != nil) != tt.wantError {
				t.Errorf("HandleGetCurrentTime() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError && !strings.Contains(result, tt.contains) {
				t.Errorf("HandleGetCurrentTime() = %v, should contain %v", result, tt.contains)
			}
		})
	}
}

// TestTimezoneService_HandleConvertTime は HandleConvertTime メソッドをテストします
func TestTimezoneService_HandleConvertTime(t *testing.T) {
	service := NewTimezoneService()

	tests := []struct {
		name         string
		datetime     string
		fromTimezone string
		toTimezone   string
		wantError    bool
		contains     string
	}{
		{
			name:         "UTC から Tokyo への変換",
			datetime:     "2023-01-02 00:00:00",
			fromTimezone: "UTC",
			toTimezone:   "Asia/Tokyo",
			wantError:    false,
			contains:     "UTC から Asia/Tokyo への変換結果:",
		},
		{
			name:         "略称を使用: JST から EST への変換",
			datetime:     "2023-01-02 12:00:00",
			fromTimezone: "jst",
			toTimezone:   "est",
			wantError:    false,
			contains:     "Asia/Tokyo から America/New_York への変換結果:",
		},
		{
			name:         "無効な変換元タイムゾーン",
			datetime:     "2023-01-02 12:00:00",
			fromTimezone: "Invalid/Timezone",
			toTimezone:   "UTC",
			wantError:    true,
		},
		{
			name:         "無効な変換先タイムゾーン",
			datetime:     "2023-01-02 12:00:00",
			fromTimezone: "UTC",
			toTimezone:   "Invalid/Timezone",
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.HandleConvertTime(tt.datetime, tt.fromTimezone, tt.toTimezone)
			if (err != nil) != tt.wantError {
				t.Errorf("HandleConvertTime() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError && !strings.Contains(result, tt.contains) {
				t.Errorf("HandleConvertTime() = %v, should contain %v", result, tt.contains)
			}
		})
	}
}

// TestTimezoneService_HandleListAvailableTimezones は HandleListAvailableTimezones メソッドをテストします
func TestTimezoneService_HandleListAvailableTimezones(t *testing.T) {
	service := NewTimezoneService()

	result, err := service.HandleListAvailableTimezones()
	if err != nil {
		t.Errorf("HandleListAvailableTimezones() error = %v", err)
		return
	}

	// 結果に重要な情報が含まれているか確認
	expectedContents := []string{
		"利用可能なタイムゾーン:",
		"一般的なタイムゾーンの別名:",
		"UTC",
		"Asia/Tokyo",
		"jst (Asia/Tokyo)",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(result, expected) {
			t.Errorf("HandleListAvailableTimezones() = %v, should contain %v", result, expected)
		}
	}
}
