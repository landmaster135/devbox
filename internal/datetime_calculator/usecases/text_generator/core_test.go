package textGenerator

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateDailyHeadingWithTimezone(t *testing.T) {
	fixedNow := time.Date(2024, 12, 31, 15, 4, 5, 0, time.UTC)
	originalNow := timeNow
	timeNow = func() time.Time {
		return fixedNow
	}
	defer func() {
		timeNow = originalNow
	}()

	timezone := "Asia/Tokyo"
	got := GenerateDailyHeading(0, timezone)

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		t.Fatalf("failed to load location: %v", err)
	}
	localized := fixedNow.In(loc)
	startDate := localized.AddDate(0, 0, 0)
	endDate := localized.AddDate(0, 0, 1)

	expected := fmt.Sprintf(
		"## %s(%s)から%s(%s)にかけて（）\n- [ ]  %s(%s)から%s(%s)にかけて進める。・・・合計0分掛かった。\n",
		startDate.Format("2006-01-02"),
		getWeekdayJapanese(startDate.Weekday()),
		endDate.Format("2006-01-02"),
		getWeekdayJapanese(endDate.Weekday()),
		startDate.Format("2006-01-02"),
		getWeekdayJapanese(startDate.Weekday()),
		endDate.Format("2006-01-02"),
		getWeekdayJapanese(endDate.Weekday()),
	)

	if got != expected {
		t.Fatalf("unexpected heading.\nwant: %q\n got: %q", expected, got)
	}
}

func TestGenerateDailyHeadingFallbackToLocal(t *testing.T) {
	fixedNow := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	originalNow := timeNow
	timeNow = func() time.Time {
		return fixedNow
	}
	defer func() {
		timeNow = originalNow
	}()

	originalLocal := time.Local
	customLocal := time.FixedZone("Test/Zone", 3*60*60)
	time.Local = customLocal
	defer func() {
		time.Local = originalLocal
	}()

	got := GenerateDailyHeading(1, "Invalid/Timezone")

	localized := fixedNow.In(customLocal)
	startDate := localized.AddDate(0, 0, 1)
	endDate := localized.AddDate(0, 0, 2)

	expected := fmt.Sprintf(
		"## %s(%s)から%s(%s)にかけて（）\n- [ ]  %s(%s)から%s(%s)にかけて進める。・・・合計0分掛かった。\n",
		startDate.Format("2006-01-02"),
		getWeekdayJapanese(startDate.Weekday()),
		endDate.Format("2006-01-02"),
		getWeekdayJapanese(endDate.Weekday()),
		startDate.Format("2006-01-02"),
		getWeekdayJapanese(startDate.Weekday()),
		endDate.Format("2006-01-02"),
		getWeekdayJapanese(endDate.Weekday()),
	)

	if got != expected {
		t.Fatalf("unexpected heading for fallback.\nwant: %q\n got: %q", expected, got)
	}
}
