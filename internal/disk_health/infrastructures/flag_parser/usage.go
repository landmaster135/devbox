package flag_parser

import (
	"fmt"
	"os"
)

func PrintUsage(format string) {
	fmt.Fprintf(os.Stderr, format, os.Args[0])
}
