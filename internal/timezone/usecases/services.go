package usecases

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// エラーメッセージの定数
const (
	ErrEmptyTimezone     = "タイムゾーンが空です"
	ErrInvalidTimezone   = "無効なタイムゾーンです"
	ErrTimezoneNotFound  = "指定されたタイムゾーンが見つかりません"
	ErrInvalidDateFormat = "無効な日付形式です"
)

// 利用可能なタイムゾーンのリスト（重複を削除）
// 実際の実装では、より完全なリストを使用することをお勧めします
var availableTimezones = []string{
	"UTC",
	"Asia/Tokyo",
	"America/New_York",
	"Europe/London",
	"Australia/Sydney",
	"Pacific/Auckland", // 重複を削除
	"Europe/Paris",
	"Asia/Shanghai",
	"America/Los_Angeles",
	"America/Chicago",
	"Africa/Abidjan",
	"Asia/Bangkok",
	"Africa/Bissau",
	"Africa/Cairo",
	"Africa/Casablanca",
	"Africa/Ceuta",
	"Africa/Johannesburg",
	"Africa/Juba",
	"Africa/Khartoum",
	"Africa/Lagos",
	"Africa/Maputo",
	"Africa/Monrovia",
	"Africa/Nairobi",
	"Africa/Ndjamena",
	"Africa/Sao_Tome",
	"Africa/Tripoli",
	"Africa/Tunis",
	"Africa/Windhoek",
	"America/Adak",
	"America/Anchorage",
	"America/Araguaina",
	"America/Argentina/Cordoba",
	"America/Argentina/San_Luis",
	"America/Toronto",
	"Antarctica/Troll",
	"Antarctica/Vostok",
	"Asia/Dubai",
	"Asia/Hong_Kong",
	"Asia/Macau",
	"Asia/Pyongyang",
	"Asia/Qatar",
	"Asia/Seoul",
	"Asia/Singapore",
	"Asia/Taipei",
	"Europe/Berlin",
	"Europe/Helsinki",
	"Europe/Istanbul",
	"Europe/Moscow",
	"Europe/Rome",
	"Europe/Zurich",
	"Indian/Maldives",
	"Pacific/Easter",
	"Pacific/Guam",
	"Pacific/Honolulu",
	"Pacific/Norfolk",
	"Pacific/Port_Moresby",
	"Pacific/Tahiti",
	"Pacific/Tarawa",
}

// 一般的なタイムゾーンのマッピング（ユーザーフレンドリーな名前）
var commonTimezones = map[string]string{
	"jst":       "Asia/Tokyo",
	"est":       "America/New_York",
	"gmt":       "Europe/London",
	"utc":       "UTC",
	"pst":       "America/Los_Angeles",
	"cst":       "America/Chicago",
	"japan":     "Asia/Tokyo",
	"us east":   "America/New_York",
	"us west":   "America/Los_Angeles",
	"uk":        "Europe/London",
	"europe":    "Europe/Paris",
	"china":     "Asia/Shanghai",
	"korea":     "Asia/Seoul",
	"australia": "Australia/Sydney",
}

// #==============================================================#
// ##          TimezoneService                                   ##
// #==============================================================#
// TimezoneService はタイムゾーン関連の機能を提供する構造体です
type TimezoneService struct {
	// 必要に応じてフィールドを追加できます
}

// NewTimezoneService は新しいTimezoneServiceを作成します
func NewTimezoneService() *TimezoneService {
	return &TimezoneService{}
}

// GetCurrentTime は指定されたタイムゾーンの現在時刻を取得するメソッドです
func (ts *TimezoneService) GetCurrentTime(timezone string) (string, error) {
	// 空のタイムゾーンをチェック
	if timezone == "" {
		return "", errors.New(ErrEmptyTimezone)
	}

	// 一般的なタイムゾーン名の変換を試みる
	timezone = ts.NormalizeTimezone(timezone)

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		// より詳細なエラーメッセージを提供
		return "", fmt.Errorf("%s: %s（近いタイムゾーン: %s）", ErrInvalidTimezone, timezone, ts.FindSimilarTimezones(timezone))
	}

	now := time.Now().In(loc)
	return now.Format("2006-01-02 15:04:05 MST (Z07:00)"), nil
}

// ConvertTime は指定された日時を指定されたタイムゾーンに変換するメソッドです
func (ts *TimezoneService) ConvertTime(dateTime, fromTimezone, toTimezone string) (string, error) {
	// 空のタイムゾーンをチェック
	if fromTimezone == "" || toTimezone == "" {
		return "", errors.New(ErrEmptyTimezone)
	}

	// タイムゾーン名の正規化
	fromTimezone = ts.NormalizeTimezone(fromTimezone)
	toTimezone = ts.NormalizeTimezone(toTimezone)

	// タイムゾーンの検証
	fromLoc, err := time.LoadLocation(fromTimezone)
	if err != nil {
		return "", fmt.Errorf("%s: %s（近いタイムゾーン: %s）", ErrInvalidTimezone, fromTimezone, ts.FindSimilarTimezones(fromTimezone))
	}

	toLoc, err := time.LoadLocation(toTimezone)
	if err != nil {
		return "", fmt.Errorf("%s: %s（近いタイムゾーン: %s）", ErrInvalidTimezone, toTimezone, ts.FindSimilarTimezones(toTimezone))
	}

	// 日時のパース
	// 複数の日時フォーマットをサポート
	formats := []string{
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
		"15:04:05",
	}

	var t time.Time
	var parseErr error

	for _, format := range formats {
		t, parseErr = time.ParseInLocation(format, dateTime, fromLoc)
		if parseErr == nil {
			break
		}
	}

	if parseErr != nil {
		return "", fmt.Errorf("%s: %s（サポートされる形式: YYYY-MM-DD HH:MM:SS, YYYY/MM/DD HH:MM:SS, YYYY-MM-DD, HH:MM:SS）", ErrInvalidDateFormat, dateTime)
	}

	// 変換先のタイムゾーンに変換
	convertedTime := t.In(toLoc)
	return convertedTime.Format("2006-01-02 15:04:05 MST (Z07:00)"), nil
}

// IsValidTimezone はタイムゾーンが有効かどうかを確認するメソッドです
func (ts *TimezoneService) IsValidTimezone(timezone string) bool {
	// 空のタイムゾーンは無効とする
	if timezone == "" {
		return false
	}

	// 一般的なタイムゾーン名の変換を試みる
	timezone = ts.NormalizeTimezone(timezone)

	_, err := time.LoadLocation(timezone)
	return err == nil
}

// NormalizeTimezone は一般的なタイムゾーン名を正規のタイムゾーン名に変換するメソッドです
func (ts *TimezoneService) NormalizeTimezone(timezone string) string {
	// 小文字に変換して検索
	lowercaseTimezone := strings.ToLower(timezone)
	if normalized, ok := commonTimezones[lowercaseTimezone]; ok {
		return normalized
	}
	return timezone
}

// FindSimilarTimezones は指定されたタイムゾーンに似たタイムゾーンを見つけるメソッドです
func (ts *TimezoneService) FindSimilarTimezones(timezone string) string {
	lowercaseTimezone := strings.ToLower(timezone)
	var similarTimezones []string

	// 部分一致するタイムゾーンを探す
	for _, tz := range availableTimezones {
		if strings.Contains(strings.ToLower(tz), lowercaseTimezone) {
			similarTimezones = append(similarTimezones, tz)
		}
	}

	// 最大5つまで表示
	if len(similarTimezones) > 5 {
		similarTimezones = similarTimezones[:5]
	}

	if len(similarTimezones) > 0 {
		return strings.Join(similarTimezones, ", ")
	}

	// 似たタイムゾーンが見つからない場合は、一般的なタイムゾーンを提案
	return "UTC, Asia/Tokyo, America/New_York, Europe/London, Australia/Sydney"
}

// GetAvailableTimezones は利用可能なタイムゾーンのリストを取得するメソッドです
func (ts *TimezoneService) GetAvailableTimezones() []string {
	var validTimezones []string

	// 利用可能なタイムゾーンのリストを作成（重複を削除）
	seen := make(map[string]bool)
	for _, tz := range availableTimezones {
		if !seen[tz] && ts.IsValidTimezone(tz) {
			validTimezones = append(validTimezones, tz)
			seen[tz] = true
		}
	}

	return validTimezones
}

// GetCommonTimezoneAliases は一般的なタイムゾーンの別名を取得するメソッドです
func (ts *TimezoneService) GetCommonTimezoneAliases() map[string]string {
	return commonTimezones
}

// #==============================================================#
// ##          MCPハンドラーメソッド                              ##
// #==============================================================#

// HandleGetCurrentTime は現在時刻取得のMCPリクエストを処理するハンドラーです
func (ts *TimezoneService) HandleGetCurrentTime(timezone string) (string, error) {
	if !ts.IsValidTimezone(timezone) {
		// より詳細なエラーメッセージを提供
		return "", fmt.Errorf("%s: %s（近いタイムゾーン: %s）", ErrInvalidTimezone, timezone, ts.FindSimilarTimezones(timezone))
	}

	// タイムゾーン名の正規化
	normalizedTimezone := ts.NormalizeTimezone(timezone)

	timeStr, err := ts.GetCurrentTime(normalizedTimezone)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s の現在時刻: %s", normalizedTimezone, timeStr), nil
}

// HandleConvertTime は時刻変換のMCPリクエストを処理するハンドラーです
func (ts *TimezoneService) HandleConvertTime(datetime, fromTimezone, toTimezone string) (string, error) {
	// タイムゾーンの検証
	if !ts.IsValidTimezone(fromTimezone) {
		return "", fmt.Errorf("変換元の%s: %s（近いタイムゾーン: %s）", ErrInvalidTimezone, fromTimezone, ts.FindSimilarTimezones(fromTimezone))
	}

	if !ts.IsValidTimezone(toTimezone) {
		return "", fmt.Errorf("変換先の%s: %s（近いタイムゾーン: %s）", ErrInvalidTimezone, toTimezone, ts.FindSimilarTimezones(toTimezone))
	}

	// タイムゾーン名の正規化
	normalizedFromTz := ts.NormalizeTimezone(fromTimezone)
	normalizedToTz := ts.NormalizeTimezone(toTimezone)

	// 時刻変換
	convertedTime, err := ts.ConvertTime(datetime, normalizedFromTz, normalizedToTz)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s の %s から %s への変換結果: %s",
		datetime, normalizedFromTz, normalizedToTz, convertedTime), nil
}

// HandleListAvailableTimezones は利用可能なタイムゾーンリストのMCPリクエストを処理するハンドラーです
func (ts *TimezoneService) HandleListAvailableTimezones() (string, error) {
	validTimezones := ts.GetAvailableTimezones()

	// 一般的なタイムゾーンの別名も表示
	aliases := ts.GetCommonTimezoneAliases()
	var aliasStrings []string
	for alias, tz := range aliases {
		aliasStrings = append(aliasStrings, fmt.Sprintf("%s (%s)", alias, tz))
	}

	result := fmt.Sprintf("利用可能なタイムゾーン: %s\n\n", strings.Join(validTimezones, ", "))
	result += fmt.Sprintf("一般的なタイムゾーンの別名: %s", strings.Join(aliasStrings, ", "))

	return result, nil
}
