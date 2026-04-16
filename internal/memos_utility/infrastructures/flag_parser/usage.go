package flag_parser

import (
	"fmt"
	"os"
)

// PrintUsage は標準エラー出力へ利用方法を出力する。
func PrintUsage(format string) {
	fmt.Fprintf(os.Stderr, format, os.Args[0])
}
