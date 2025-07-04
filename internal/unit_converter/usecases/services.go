// Package usecases contains all domain logic for the unit‑converter tool.
// It exposes a minimal façade (Convert, Categories, PrefixTable) so the CLI layer
// depends only on behaviour, not implementation details.
//
// Conventions
// -----------
//   - Base units are kept in small factor maps per category.
//   - Metric (SI) prefixes are recognised dynamically – the factor maps do **not** need
//     to include "km", "cm2", etc.
//   - Temperatures are handled via dedicated formulas rather than factors.
//   - The public API is deliberately tiny to keep coupling low.
package usecases

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// Convert converts the given value from one unit to another within a category.
// It returns the converted value, or an error if the category or unit is unsupported.
func Convert(category string, value float64, fromUnit, toUnit string) (float64, error) {
	category = strings.ToLower(category)
	fromUnit = strings.ToLower(fromUnit)
	toUnit = strings.ToLower(toUnit)

	switch category {
	case "length":
		return genericConvert(value, fromUnit, toUnit, lengthFactors)
	case "weight", "mass":
		return genericConvert(value, fromUnit, toUnit, weightFactors)
	case "area":
		return genericConvert(value, fromUnit, toUnit, areaFactors)
	case "volume":
		return genericConvert(value, fromUnit, toUnit, volumeFactors)
	case "temp", "temperature":
		return convertTemperature(value, fromUnit, toUnit)
	default:
		return 0, fmt.Errorf("unsupported category: %s", category)
	}
}

// Categories returns a copy of the map of supported categories and their *base* units.
// The slice values are sorted for stable presentation in the CLI.
func Categories() map[string][]string {
	out := make(map[string][]string, len(categoriesUnits))
	for k, v := range categoriesUnits {
		vv := make([]string, len(v))
		copy(vv, v)
		out[k] = vv
	}
	return out
}

// PrefixTable returns a formatted string that lists all recognised SI prefixes.
func PrefixTable() string { return prefixInfo }

// ---------------------------------------------------------------------------
// Internal engine (not exported)
// ---------------------------------------------------------------------------

func genericConvert(v float64, from, to string, factors map[string]float64) (float64, error) {
	fv, ok1 := resolveUnitFactor(from, factors)
	tv, ok2 := resolveUnitFactor(to, factors)
	if !ok1 || !ok2 {
		return 0, fmt.Errorf("unknown unit: %s or %s", from, to)
	}
	base := v * fv
	return base / tv, nil
}

func convertTemperature(v float64, from, to string) (float64, error) {
	normalize := map[string]func(float64) float64{
		"c": func(x float64) float64 { return x },
		"k": func(x float64) float64 { return x - 273.15 },
		"f": func(x float64) float64 { return (x - 32) * 5 / 9 },
	}
	denormalize := map[string]func(float64) float64{
		"c": func(x float64) float64 { return x },
		"k": func(x float64) float64 { return x + 273.15 },
		"f": func(x float64) float64 { return x*9/5 + 32 },
	}

	nf, ok1 := normalize[from]
	df, ok2 := denormalize[to]
	if !ok1 || !ok2 {
		return 0, fmt.Errorf("unknown temperature unit: %s or %s", from, to)
	}
	return df(nf(v)), nil
}

// ----------------------------- SI prefixes ----------------------------------

var siPrefixes = map[string]float64{
	"Y":  1e24, // yotta
	"Z":  1e21,
	"E":  1e18,
	"P":  1e15,
	"T":  1e12,
	"G":  1e9,
	"M":  1e6,
	"k":  1e3,
	"h":  1e2,
	"da": 1e1,
	"d":  1e-1,
	"c":  1e-2,
	"m":  1e-3,
	"u":  1e-6,
	"μ":  1e-6,
	"n":  1e-9,
	"p":  1e-12,
	"f":  1e-15,
	"a":  1e-18,
	"z":  1e-21,
	"y":  1e-24,
}

var siPrefixKeys = func() []string {
	keys := make([]string, 0, len(siPrefixes))
	for k := range siPrefixes {
		keys = append(keys, k)
	}
	// ensure longer prefixes ("da") come before shorter ones ("d")
	sort.Slice(keys, func(i, j int) bool { return len(keys[i]) > len(keys[j]) })
	return keys
}()

const prefixInfo = "" +
	"Metric (SI) prefixes available:\n" +
	"  Y  1e24   yotta\n  Z  1e21   zetta\n  E  1e18   exa\n  P  1e15   peta\n  T  1e12   tera\n  G  1e9    giga\n  M  1e6    mega\n  k  1e3    kilo\n  h  1e2    hecto\n  da 1e1    deca\n  d  1e-1   deci\n  c  1e-2   centi\n  m  1e-3   milli\n  μ/u 1e-6  micro\n  n  1e-9   nano\n  p  1e-12  pico\n  f  1e-15  femto\n  a  1e-18  atto\n  z  1e-21  zepto\n  y  1e-24  yocto\n"

func resolveUnitFactor(u string, base map[string]float64) (float64, bool) {
	if f, ok := base[u]; ok {
		return f, true
	}
	for _, p := range siPrefixKeys {
		if strings.HasPrefix(u, p) {
			suffix := u[len(p):]
			b, ok := base[suffix]
			if !ok {
				continue
			}
			pf := siPrefixes[p]
			exp := 1
			if n := len(suffix); n > 0 {
				if d := suffix[n-1]; d >= '2' && d <= '9' {
					exp = int(d - '0')
				}
			}
			return b * math.Pow(pf, float64(exp)), true
		}
	}
	return 0, false
}

// --------------------------- factor tables ----------------------------------

var lengthFactors = map[string]float64{
	"m":  1,
	"km": 1000, // convenience alias; also handled by prefix
	"cm": 0.01,
	"mm": 0.001,
	"mi": 1609.344,
	"yd": 0.9144,
	"ft": 0.3048,
	"in": 0.0254,
}

var weightFactors = map[string]float64{
	"g":  0.001,
	"kg": 1,
	"t":  1000,
	"lb": 0.45359237,
	"oz": 0.028349523125,
}

var areaFactors = map[string]float64{
	"m2":  1,
	"km2": 1e6,
	"cm2": 1e-4,
	"mm2": 1e-6,
	"ha":  10000,
	"ac":  4046.8564224,
	"ft2": 0.09290304,
	"yd2": 0.83612736,
	"mi2": 2589988.110336,
}

var volumeFactors = map[string]float64{
	"m3":   1,
	"l":    0.001,
	"ml":   1e-6,
	"cm3":  1e-6,
	"mm3":  1e-9,
	"gal":  0.003785411784,
	"qt":   0.000946352946,
	"pt":   0.000473176473,
	"cup":  0.0002365882365,
	"floz": 2.95735295625e-5,
	"in3":  1.6387064e-5,
	"ft3":  0.028316846592,
}

var categoriesUnits = map[string][]string{
	"length": sortedKeys(lengthFactors),
	"weight": sortedKeys(weightFactors),
	"area":   sortedKeys(areaFactors),
	"volume": sortedKeys(volumeFactors),
	"temp":   {"C", "F", "K"},
}

func sortedKeys(m map[string]float64) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
