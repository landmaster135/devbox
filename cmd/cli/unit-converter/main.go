package main

// unit‑converter is a *thin* CLI wrapper. All conversion logic, tables, and prefix handling live
// in the `github.com/landmaster135/devbox/internal/unit_converter/usecases`
// package. That keeps the command layer free of domain rules, achieving loose coupling.
//
// Build:
//   go build -o unit-converter
//
// Usage examples:
//   unit-converter length 10 km mi          # 10 km → mile (default 6 sig‑digits)
//   unit-converter -p 3 volume 250 mL cup   # 250 mL → US cup with 3 sig‑digits
//   unit-converter --list                   # show supported base units + prefix table

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	uc "github.com/landmaster135/devbox/internal/unit_converter/usecases"
)

// ------------------- application framing (PDF‑merger style) ----------------- //

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

func main() {
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}

// run parses CLI args, delegates all domain work to the use‑case layer, and returns an exit code.
func run(args []string, stdout, stderr io.Writer) exitCode {
	fs := flag.NewFlagSet("unit-converter", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var (
		precision int
		listFlag  bool
	)

	fs.IntVar(&precision, "precision", 6, "significant digits to display")
	fs.IntVar(&precision, "p", 6, "significant digits to display (shorthand)")
	fs.BoolVar(&listFlag, "list", false, "list supported categories, base units and prefix table")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "Unit Converter CLI (presentation layer)\n\nUsage:\n    unit-converter [flags] <category> <value> <from-unit> <to-unit>\n\n")
		fs.PrintDefaults()
		fmt.Fprintln(stderr, "\nRun with --list to see supported base units and prefix table.")
	}

	if err := fs.Parse(args); err != nil {
		// parsing error already reported by FlagSet
		return exitCodeError
	}

	if listFlag {
		printList(stdout)
		return exitCodeOK
	}

	pos := fs.Args()
	if len(pos) != 4 {
		fs.Usage()
		return exitCodeError
	}

	category := strings.ToLower(pos[0])
	valueStr := pos[1]
	fromUnit := strings.ToLower(pos[2])
	toUnit := strings.ToLower(pos[3])

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		fmt.Fprintf(stderr, "invalid value: %v\n", err)
		return exitCodeError
	}

	// Delegate conversion to use‑case layer
	result, err := uc.Convert(category, value, fromUnit, toUnit)
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return exitCodeError
	}

	fmt.Fprintf(stdout, "%g %s = %.*g %s\n", value, fromUnit, precision, result, toUnit)
	return exitCodeOK
}

// printList renders supported base units & the prefix table coming from the use‑case layer.
func printList(w io.Writer) {
	cats := uc.Categories() // map[string][]string
	fmt.Fprintln(w, "Supported categories and *base* units (SI prefixes can be prepended):")
	for cat, units := range cats {
		fmt.Fprintf(w, "  %-10s %s\n", cat, strings.Join(units, ", "))
	}
	fmt.Fprintln(w)
	fmt.Fprint(w, uc.PrefixTable())
}
