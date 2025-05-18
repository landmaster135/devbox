package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	usecases "github.com/landmaster135/devbox/internal/independencies/kana_converter/usecases"
)

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

// run はカナ変換ツールの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグセットを作成
	flagSet := flag.NewFlagSet("kana-converter", flag.ContinueOnError)
	flagSet.SetOutput(stderr)

	// コマンドライン引数の定義
	inputString := flagSet.String("input", "", "Input characters containing Katakana")
	mode := flagSet.String("mode", string(usecases.FullWidthMode),
		"Conversion mode: 'full' for full-width or 'half' for half-width kana")

	// 引数の解析
	if err := flagSet.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// 非フラグ引数がある場合、それを入力として使用
	if flagSet.NArg() > 0 {
		*inputString = strings.Join(flagSet.Args(), " ")
	}

	// 入力が空の場合、処理をキャンセル
	if *inputString == "" {
		fmt.Fprintln(stderr, "No input. This process has cancelled.")
		return exitCodeError
	}

	// モードの検証
	if !usecases.IsValidMode(strings.ToLower(*mode)) {
		fmt.Fprintf(stderr, "Invalid mode: %s. Use 'full' or 'half'.\n", *mode)
		return exitCodeError
	}

	// モードの変換
	convMode := usecases.GetModeFromString(strings.ToLower(*mode))

	// コンバーターの作成と変換の実行
	converter := usecases.NewKanaConverter(convMode)
	convertedString := converter.Convert(*inputString)

	// 結果を出力
	fmt.Fprintf(stdout, "Input characters: %s\n", *inputString)
	fmt.Fprintf(stdout, "Mode: %s\n", usecases.GetModeDescription(convMode))
	fmt.Fprintf(stdout, "Output characters: %s\n", convertedString)

	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
