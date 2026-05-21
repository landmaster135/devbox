package flag_parser

import (
	"fmt"
	"os"
)

// PrintUsage は usage テンプレートを標準エラーに出力します。
func PrintUsage(format string) {
	fmt.Fprintf(os.Stderr, format, os.Args[0])
}
