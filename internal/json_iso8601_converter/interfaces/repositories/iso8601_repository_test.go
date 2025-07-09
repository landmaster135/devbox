package repositories

import (
	"testing"
	"time"
)

// TestNewISO8601Repository は NewISO8601Repository 関数をテストします
func TestNewISO8601Repository(t *testing.T) {
	repo := NewISO8601Repository()
	if repo == nil {
		t.Error("NewISO8601Repository() がnilを返しました")
	}
}

// TestISO8601RepositoryImpl_ParseISO8601 は ParseISO8601 メソッドをテストします
func TestISO8601RepositoryImpl_ParseISO8601(t *testing.T) {
	// テスト対象のインスタンスを作成
	repo := NewISO8601Repository()

	// テストケース
	tests := []struct {
		name     string
		dateStr  string
		expected int64
		wantErr  bool
	}{
		{
			name:     "RFC3339形式の日時文字列",
			dateStr:  "2023-04-01T12:34:56Z",
			expected: time.Date(2023, 4, 1, 12, 34, 56, 0, time.UTC).Unix(),
			wantErr:  false,
		},
		{
			name:     "RFC3339形式の日時文字列（タイムゾーン付き）",
			dateStr:  "2023-04-01T12:34:56+09:00",
			expected: time.Date(2023, 4, 1, 12, 34, 56, 0, time.FixedZone("", 9*60*60)).Unix(),
			wantErr:  false,
		},
		{
			name:     "RFC3339Nano形式の日時文字列",
			dateStr:  "2023-04-01T12:34:56.789Z",
			expected: time.Date(2023, 4, 1, 12, 34, 56, 789000000, time.UTC).Unix(),
			wantErr:  false,
		},
		{
			name:     "無効な日時文字列",
			dateStr:  "2023-04-01 12:34:56",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "空の文字列",
			dateStr:  "",
			expected: 0,
			wantErr:  true,
		},
	}

	// 各テストケースを実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト実行
			got, err := repo.ParseISO8601(tt.dateStr)

			// エラーの検証
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseISO8601() エラー = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			// エラーが期待される場合は、ここで終了
			if tt.wantErr {
				return
			}

			// 結果の検証
			if got != tt.expected {
				t.Errorf("ParseISO8601() = %v, want = %v", got, tt.expected)
			}
		})
	}
}
