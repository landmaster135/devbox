package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/landmaster135/devbox/internal/independencies/iso8601-converter/usecases"
)

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

// run はISO-8601コンバーターツールの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグセットを作成
	fs := flag.NewFlagSet("iso8601-converter", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// コマンドライン引数の定義
	toISO := fs.Bool("to-iso", false, "UNIXタイムスタンプをISO-8601形式に変換")
	toUnix := fs.Bool("to-unix", false, "ISO-8601形式をUNIXタイムスタンプに変換")
	isJST := fs.Bool("is-jst", false, "JSTタイムゾーンを使用（日付変換のみ）")
	help := fs.Bool("help", false, "ヘルプメッセージを表示")
	input := fs.String("input", "", "変換する値")

	// 引数の解析
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// ヘルプメッセージの表示
	if *help || *input == "" {
		showHelp(stdout)
		return exitCodeOK
	}

	// 変換処理の実行
	if *toISO {
		result, err := usecases.UnixToISO8601(*input)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: %v\n", err)
			return exitCodeError
		}
		fmt.Fprintln(stdout, result)
	} else if *toUnix {
		// 入力に時刻情報が含まれているかチェック
		if strings.Contains(*input, "T") {
			// 入力はISO-8601形式
			result, err := usecases.ISO8601ToUnix(*input)
			if err != nil {
				fmt.Fprintf(stderr, "エラー: %v\n", err)
				return exitCodeError
			}
			fmt.Fprintln(stdout, result)
		} else {
			// 入力は日付形式
			result, err := usecases.DateToUnix(*input, *isJST)
			if err != nil {
				fmt.Fprintf(stderr, "エラー: %v\n", err)
				return exitCodeError
			}
			fmt.Fprintln(stdout, result)
		}
	} else {
		fmt.Fprintf(stderr, "エラー: --to-iso または --to-unix フラグを指定する必要があります\n")
		return exitCodeError
	}

	return exitCodeOK
}

// showHelp はヘルプメッセージを表示します
func showHelp(w io.Writer) {
	fmt.Fprintln(w, "ISO-8601 コンバーター")
	fmt.Fprintln(w, "使用方法:")
	fmt.Fprintln(w, "  UNIXタイムスタンプをISO-8601形式に変換:")
	fmt.Fprintln(w, "    iso8601-converter --to-iso --input <unix_timestamp>")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  ISO-8601形式をUNIXタイムスタンプに変換:")
	fmt.Fprintln(w, "    iso8601-converter --to-unix --input <iso8601_time>")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  日付をUNIXタイムスタンプに変換 (UTC):")
	fmt.Fprintln(w, "    iso8601-converter --to-unix --input <date>")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  日付をUNIXタイムスタンプに変換 (JST):")
	fmt.Fprintln(w, "    iso8601-converter --to-unix --is-jst --input <date>")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "例:")
	fmt.Fprintln(w, "  iso8601-converter --to-iso --input 1619712000")
	fmt.Fprintln(w, "  iso8601-converter --to-unix --input 2021-04-30T00:00:00Z")
	fmt.Fprintln(w, "  iso8601-converter --to-unix --input 2021-04-30")
	fmt.Fprintln(w, "  iso8601-converter --to-unix --is-jst --input 2021-04-30")
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
