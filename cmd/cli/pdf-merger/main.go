package main

// unitconv is a command‑line unit‑converter supporting length, weight (mass), temperature, area
// and volume.  It recognises all standard SI prefixes (yotta … yocto) so you may write e.g.
// "µm", "nm2", "pL", "kg"など。"--list" で基底単位一覧と接頭語表を表示します。
//
// Build:
//   go build -o unitconv
//
// Example:
//   unitconv length 10 km mi
//   unitconv -p 3 volume 250 mL cup
//
// Coding style: The main work is done by run(), which takes explicit io.Writer parameters and
// returns an exit code – mirroring the sample PDF‑merger style the user provided.

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

// --------------------- application framing (PDF‑merger style) --------------- //

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

func main() {
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}

// run executes the converter logic and returns an exit code.
func run(args []string, stdout, stderr io.Writer) exitCode {
	// FlagSet isolated from global default set, writes usage/errors to stderr.
	fs := flag.NewFlagSet("unitconv", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var (
		precision int
		listFlag  bool
	)

	fs.IntVar(&precision, "precision", 6, "significant digits to display")
	fs.IntVar(&precision, "p", 6, "significant digits to display (shorthand)")
	fs.BoolVar(&listFlag, "list", false, "list supported categories, base units and prefix table")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "Unit Converter CLI with SI‑prefix support\n\nUsage:\n    unitconv [flags] <category> <value> <from-unit> <to-unit>\n\n")
		fs.PrintDefaults()
		fmt.Fprintln(stderr, "\nRun with --list to see supported base units and prefix table.")
	}

	if err := fs.Parse(args); err != nil {
		// Parse error already written to stderr by FlagSet.
		return exitCodeError
	}

	if listFlag {
		printList(stdout)
		return exitCodeOK
	}

	positional := fs.Args()
	if len(positional) != 4 {
		fs.Usage()
		return exitCodeError
	}

	category := strings.ToLower(positional[0])
	valueStr := positional[1]
	fromUnit := strings.ToLower(positional[2])
	toUnit := strings.ToLower(positional[3])

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		fmt.Fprintf(stderr, "invalid value: %v\n", err)
		return exitCodeError
	}

	var result float64

	switch category {
	case "length":
		result, err = convertLength(value, fromUnit, toUnit)
	case "weight", "mass":
		result, err = convertWeight(value, fromUnit, toUnit)
	case "temp", "temperature":
		result, err = convertTemperature(value, fromUnit, toUnit)
	case "area":
		result, err = convertArea(value, fromUnit, toUnit)
	case "volume":
		result, err = convertVolume(value, fromUnit, toUnit)
	default:
		err = fmt.Errorf("unsupported category: %s", category)
	}

	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return exitCodeError
	}

	fmt.Fprintf(stdout, "%g %s = %.*g %s\n", value, fromUnit, precision, result, toUnit)
	return exitCodeOK
}

// ------------------------------ helpers ------------------------------------ //

func printList(w io.Writer) {
	fmt.Fprintln(w, "Supported categories and *base* units (SI prefixes can be prepended):")
	for cat, units := range categoriesUnits {
		fmt.Fprintf(w, "  %-10s %s\n", cat, strings.Join(units, ", "))
	}
	fmt.Fprintln(w)
	fmt.Fprint(w, prefixInfo)
}

// ------------------------ SI prefix definitions ---------------------------- //

var siPrefixes = map[string]float64{
	"Y":  1e24,  // yotta
	"Z":  1e21,  // zetta
	"E":  1e18,  // exa
	"P":  1e15,  // peta
	"T":  1e12,  // tera
	"G":  1e9,   // giga
	"M":  1e6,   // mega
	"k":  1e3,   // kilo
	"h":  1e2,   // hecto
	"da": 1e1,   // deca
	"d":  1e-1,  // deci
	"c":  1e-2,  // centi
	"m":  1e-3,  // milli
	"u":  1e-6,  // (ASCII) micro
	"μ":  1e-6,  // (Greek mu) micro
	"n":  1e-9,  // nano
	"p":  1e-12, // pico
	"f":  1e-15, // femto
	"a":  1e-18, // atto
	"z":  1e-21, // zepto
	"y":  1e-24, // yocto
}

// Longest prefixes first so "da" beats "d".
var siPrefixKeys = []string{
	"da", "Y", "Z", "E", "P", "T", "G", "M", "k", "h", "d", "c", "m", "u", "μ", "n", "p", "f", "a", "z", "y",
}

const prefixInfo = "" +
	"Metric (SI) prefixes available:\n" +
	"  Y  1e24   yotta\n  Z  1e21   zetta\n  E  1e18   exa\n  P  1e15   peta\n  T  1e12   tera\n  G  1e9    giga\n  M  1e6    mega\n  k  1e3    kilo\n  h  1e2    hecto\n  da 1e1    deca\n  d  1e-1   deci\n  c  1e-2   centi\n  m  1e-3   milli\n  μ/u 1e-6  micro\n  n  1e-9   nano\n  p  1e-12  pico\n  f  1e-15  femto\n  a  1e-18  atto\n  z  1e-21  zepto\n  y  1e-24  yocto\n"

// ------------------------ conversion functions ----------------------------- //

func convertLength(v float64, from, to string) (float64, error) {
	return genericConvert(v, from, to, lengthFactors)
}
func convertWeight(v float64, from, to string) (float64, error) {
	return genericConvert(v, from, to, weightFactors)
}
func convertArea(v float64, from, to string) (float64, error) {
	return genericConvert(v, from, to, areaFactors)
}
func convertVolume(v float64, from, to string) (float64, error) {
	return genericConvert(v, from, to, volumeFactors)
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

// ---------------------------- core engine ---------------------------------- //

func genericConvert(v float64, from, to string, factors map[string]float64) (float64, error) {
	fv, ok1 := resolveUnitFactor(from, factors)
	tv, ok2 := resolveUnitFactor(to, factors)
	if !ok1 || !ok2 {
		return 0, fmt.Errorf("unknown unit: %s or %s", from, to)
	}
	base := v * fv
	return base / tv, nil
}

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

// ----------------------------- factor tables ------------------------------- //

var lengthFactors = map[string]float64{
	"m":  1,
	"km": 1000, // kept for back‑compat
	"cm": 0.01,
	"mm": 0.001,
	"mi": 1609.344,
	"yd": 0.9144,
	"ft": 0.3048,
	"in": 0.0254,
}

var weightFactors = map[string]float64{
	"g":  0.001,
	"kg": 1, // convenience
	"t":  1000,
	"lb": 0.45359237,
	"oz": 0.028349523125,
}

var areaFactors = map[string]float64{
	"m2":  1,
	"km2": 1e6,
	"cm2": 0.0001,
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
	"gal":  0.003785411784,  // US gallon
	"qt":   0.000946352946,  // US quart
	"pt":   0.000473176473,  // US pint
	"cup":  0.0002365882365, // US cup
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
	for i := 1; i < len(keys); i++ {
		j := i
		for j > 0 && keys[j-1] > keys[j] {
			keys[j-1], keys[j] = keys[j], keys[j-1]
			j--
		}
	}
	return keys
}
