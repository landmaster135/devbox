package filter

import (
	"strings"
	"testing"
)

func TestNormalizeFilter_RFC3339UTC_Normal(t *testing.T) {
	got, err := NormalizeFilter(`created_ts > '2023-01-01T13:00:00Z' && visibility == "PUBLIC"`)
	if err != nil {
		t.Fatalf("NormalizeFilter() error = %v", err)
	}

	want := `created_ts > 1672578000 && visibility == "PUBLIC"`
	if got != want {
		t.Fatalf("NormalizeFilter() = %q, want %q", got, want)
	}
}

func TestNormalizeFilter_RFC3339Offset_Normal(t *testing.T) {
	got, err := NormalizeFilter(`updated_ts >= "2023-01-01T13:00:00+09:00"`)
	if err != nil {
		t.Fatalf("NormalizeFilter() error = %v", err)
	}

	want := `updated_ts >= 1672545600`
	if got != want {
		t.Fatalf("NormalizeFilter() = %q, want %q", got, want)
	}
}

func TestNormalizeFilter_RFC3339Nano_Normal(t *testing.T) {
	got, err := NormalizeFilter(`updated_ts == "2023-01-01T13:00:00.123456789Z"`)
	if err != nil {
		t.Fatalf("NormalizeFilter() error = %v", err)
	}

	want := `updated_ts == 1672578000`
	if got != want {
		t.Fatalf("NormalizeFilter() = %q, want %q", got, want)
	}
}

func TestNormalizeFilter_NonTimestampFieldUnchanged_Normal(t *testing.T) {
	input := `visibility == "PUBLIC" && content.contains("memos")`
	got, err := NormalizeFilter(input)
	if err != nil {
		t.Fatalf("NormalizeFilter() error = %v", err)
	}
	if got != input {
		t.Fatalf("NormalizeFilter() = %q, want %q", got, input)
	}
}

func TestNormalizeFilter_UnixTimestampUnchanged_Normal(t *testing.T) {
	input := `created_ts > 1672578000 && updated_ts < now()`
	got, err := NormalizeFilter(input)
	if err != nil {
		t.Fatalf("NormalizeFilter() error = %v", err)
	}
	if got != input {
		t.Fatalf("NormalizeFilter() = %q, want %q", got, input)
	}
}

func TestNormalizeFilter_InvalidDateTime_Error(t *testing.T) {
	_, err := NormalizeFilter(`created_ts > '2023-01-01T13:00:00'`)
	if err == nil {
		t.Fatal("NormalizeFilter() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "タイムゾーン必須") {
		t.Fatalf("error = %v, want タイムゾーン必須", err)
	}
}
