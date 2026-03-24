package flag_parser

import (
	"fmt"
	"os"
)

// PrintUsage はusageテンプレートを標準エラーへ出力する
func PrintUsage(format string) {
	fmt.Fprintf(os.Stderr, format, os.Args[0])
}
