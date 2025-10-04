package usecases

import (
	"testing"
)

func assertFloat(t *testing.T, value any, expected float64) {
	t.Helper()
	v, ok := value.(float64)
	if !ok {
		t.Fatalf("value is not float64: %#v", value)
	}
	if v != expected {
		t.Fatalf("unexpected float value: got %v want %v", v, expected)
	}
}
