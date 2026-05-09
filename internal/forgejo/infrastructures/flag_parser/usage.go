package flag_parser

import (
	"fmt"
	"os"
)

// PrintUsage はヘルプ本文を標準エラー出力に表示します。
func PrintUsage(format string) {
	fmt.Fprint(os.Stderr, format)
}
