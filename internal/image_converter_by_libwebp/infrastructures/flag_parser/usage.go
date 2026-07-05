package flag_parser

import (
	"fmt"
	"os"
)

// PrintUsage は usage テンプレートへ実行ファイル名を注入して出力します。
func PrintUsage(format string) {
	fmt.Fprintf(os.Stderr, format, os.Args[0])
}
