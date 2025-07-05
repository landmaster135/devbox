// Unit tests for the unit_converter/usecases package.
// Run: go test ./internal/unit_converter/...
package usecases

import (
	"math"
	"strings"
	"testing"
)

// tolerance for floating‑point comparisons
const eps = 1e-9

func almostEqual(a, b float64) bool {
	if a == b {
		return true
	}
	diff := math.Abs(a - b)
	// relative tolerance
	return diff <= eps*math.Max(math.Abs(a), math.Abs(b))
}

func TestConvertSuccess(t *testing.T) {
	cases := []struct {
		cat      string
		val      float64
		from, to string
		want     float64
	}{
		{"length", 1, "km", "m", 1000},           // SI prefix
		{"weight", 2.20462262185, "lb", "kg", 1}, // pounds → kilograms
		{"area", 1, "ha", "m2", 10000},           // hectare → m²
		{"volume", 1, "l", "m3", 0.001},          // litre → m³
		{"temp", 0, "c", "k", 273.15},            // Celsius → Kelvin
		{"temperature", 32, "f", "c", 0},         // Fahrenheit → Celsius
	}

	for _, c := range cases {
		got, err := Convert(c.cat, c.val, c.from, c.to)
		if err != nil {
			t.Errorf("Convert(%s, %g, %s, %s) unexpected error: %v", c.cat, c.val, c.from, c.to, err)
			continue
		}
		if !almostEqual(got, c.want) {
			t.Errorf("Convert(%s, %g, %s, %s) = %g, want %g", c.cat, c.val, c.from, c.to, got, c.want)
		}
	}
}

func TestConvertErrors(t *testing.T) {
	_, err := Convert("unknown", 1, "m", "km")
	if err == nil {
		t.Error("expected error for unsupported category")
	}

	_, err = Convert("length", 1, "foo", "m")
	if err == nil {
		t.Error("expected error for unknown from‑unit")
	}

	_, err = Convert("length", 1, "m", "bar")
	if err == nil {
		t.Error("expected error for unknown to‑unit")
	}
}

func TestCategoriesAndPrefixTable(t *testing.T) {
	cats := Categories()
	if len(cats) == 0 {
		t.Fatal("Categories() returned empty map")
	}
	if _, ok := cats["length"]; !ok {
		t.Error("Categories() missing 'length'")
	}

	pt := PrefixTable()
	prefixes := []string{"Y", "Z", "E", "k", "m", "μ", "n", "y"}
	for _, p := range prefixes {
		if !strings.Contains(pt, p) {
			t.Errorf("PrefixTable() missing prefix %s", p)
		}
	}
}
