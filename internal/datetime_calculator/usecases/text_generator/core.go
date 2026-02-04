package textGenerator

import (
	"fmt"
	"strings"
	"time"
)

var timeNow = time.Now

// GenerateDailyHeading は日次見出しテキストを生成して返す
func GenerateDailyHeading(dayOffset int, timezone string) string {
	loc := resolveLocation(timezone)
	now := timeNow().In(loc)
	startDate := now.AddDate(0, 0, dayOffset)
	endDate := startDate.AddDate(0, 0, 1)

	// 日付と曜日の文字列を構築
	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")
	startWeekday := getWeekdayJapanese(startDate.Weekday())
	endWeekday := getWeekdayJapanese(endDate.Weekday())

	heading := fmt.Sprintf("## %s(%s)から%s(%s)にかけて（）\n", startDateStr, startWeekday, endDateStr, endWeekday)
	checkpoint := fmt.Sprintf("- [ ]  %s(%s)から%s(%s)にかけて進める。・・・合計0分掛かった。\n", startDateStr, startWeekday, endDateStr, endWeekday)

	return heading + checkpoint
}

func resolveLocation(timezone string) *time.Location {
	tz := strings.TrimSpace(timezone)
	if tz == "" || strings.EqualFold(tz, "local") {
		return time.Local
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Local
	}
	return loc
}

// getWeekdayJapanese は曜日を英語略称で返す
func getWeekdayJapanese(weekday time.Weekday) string {
	switch weekday {
	case time.Sunday:
		return "Sun"
	case time.Monday:
		return "Mon"
	case time.Tuesday:
		return "Tue"
	case time.Wednesday:
		return "Wed"
	case time.Thursday:
		return "Thu"
	case time.Friday:
		return "Fri"
	case time.Saturday:
		return "Sat"
	default:
		return ""
	}
}
