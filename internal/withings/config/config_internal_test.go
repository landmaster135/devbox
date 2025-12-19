package config

import (
	"reflect"
	"testing"
	"time"
)

func TestParseMeasureTypesAll(t *testing.T) {
	types, err := parseMeasureTypes("all")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if types != nil {
		t.Fatalf("expected nil slice for all, got %#v", types)
	}
}

func TestParseMeasureTypesAliases(t *testing.T) {
	types, err := parseMeasureTypes(" weight ,diastolic,10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []int{1, 9, 10}
	if !reflect.DeepEqual(types, expected) {
		t.Fatalf("unexpected measure types: want %v, got %v", expected, types)
	}
}

func TestParseMeasureTypesInvalid(t *testing.T) {
	if _, err := parseMeasureTypes("unknown"); err == nil {
		t.Fatal("expected error for unknown measure type")
	}
}

func TestParseDate(t *testing.T) {
	date, err := parseDate("2025-10-02")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if date.Year() != 2025 || date.Month() != time.October || date.Day() != 2 {
		t.Fatalf("unexpected date: %v", date)
	}
	if _, err := parseDate("2025/10/02"); err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestNormalizeCommaSeparated(t *testing.T) {
	got := normalizeCommaSeparated("  a , b ,,, c  ")
	if got != "a,b,c" {
		t.Fatalf("unexpected normalized value: %s", got)
	}
}
