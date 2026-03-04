package common

import (
	"strings"
	"testing"
)

func TestFormatTable_Normal(t *testing.T) {
	t.Parallel()

	result := FormatTable(
		[]string{"NAME", "ZONE"},
		[][]string{
			{"pd-balanced", "asia-southeast3-a"},
			{"pd-ssd", "asia-southeast3-a"},
		},
	)

	if !strings.Contains(result, "NAME") {
		t.Fatalf("header missing: %s", result)
	}
	if !strings.Contains(result, "pd-balanced") {
		t.Fatalf("row missing: %s", result)
	}
}

func TestZoneBasename_Normal(t *testing.T) {
	t.Parallel()

	got := ZoneBasename("https://www.googleapis.com/compute/v1/projects/test/zones/asia-southeast3-a")
	if got != "asia-southeast3-a" {
		t.Fatalf("zone basename mismatch: %s", got)
	}
}
