package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/landmaster135/devbox/internal/iso8601-converter/usecases"
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
	operation := fs.String("operation", "", "実行する操作を指定 (to-iso|to-unix|now)")
	format := fs.String("format", "all", "--operation now 用の出力形式 (all|unix|iso)")
	isJST := fs.Bool("is-jst", false, "JSTタイムゾーンを使用（日付変換のみ）")
	help := fs.Bool("help", false, "ヘルプメッセージを表示")
	input := fs.String("input", "", "変換する値")

	// 引数の解析
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	op := strings.ToLower(strings.TrimSpace(*operation))

	// ヘルプメッセージの表示
	if *help || op == "" {
		showHelp(stdout)
		return exitCodeOK
	}

	// 現在日時の表示
	if op == "now" {
		fmt.Fprintf(stdout, "現在日時:\n")
		formatType := strings.ToLower(strings.TrimSpace(*format))
		switch formatType {
		case "", "all":
			fmt.Fprintf(stdout, "  ISO-8601形式 (UTC): %s\n", usecases.NowToISO8601InUTC())
			fmt.Fprintf(stdout, "  ISO-8601形式 (JST): %s\n", usecases.NowToISO8601InJST())
			fmt.Fprintf(stdout, "  UNIXタイムスタンプ: %s\n", usecases.NowToUnix())
		case "iso":
			fmt.Fprintf(stdout, "  ISO-8601形式 (UTC): %s\n", usecases.NowToISO8601InUTC())
			fmt.Fprintf(stdout, "  ISO-8601形式 (JST): %s\n", usecases.NowToISO8601InJST())
		case "unix":
			fmt.Fprintf(stdout, "  UNIXタイムスタンプ: %s\n", usecases.NowToUnix())
		default:
			fmt.Fprintf(stderr, "エラー: --format には all, unix, iso のいずれかを指定してください\n")
			return exitCodeError
		}
		return exitCodeOK
	}

	if *input == "" {
		fmt.Fprintln(stderr, "エラー: --input フラグで変換対象を指定してください")
		return exitCodeError
	}

	// 変換処理の実行
	switch op {
	case "to-iso":
		result, err := usecases.UnixToISO8601(*input)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: %v\n", err)
			return exitCodeError
		}
		fmt.Fprintln(stdout, result)
	case "to-unix":
		// 入力に時刻情報が含まれているかチェック
		if strings.Contains(*input, "T") {
			// 入力はISO-8601形式
			result, err := usecases.ISO8601ToUnix(*input, *isJST)
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
	default:
		fmt.Fprintf(stderr, "エラー: --operation には to-iso, to-unix, now のいずれかを指定してください\n")
		return exitCodeError
	}

	return exitCodeOK
}

// showHelp はヘルプメッセージを表示します
func showHelp(w io.Writer) {
	fmt.Fprintln(w, "ISO-8601 コンバーター")
	fmt.Fprintln(w, "使用方法:")
	fmt.Fprintln(w, "  現在日時を表示:")
	fmt.Fprintln(w, "    iso8601-converter --operation now [--format all|unix|iso]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  UNIXタイムスタンプをISO-8601形式に変換:")
	fmt.Fprintln(w, "    iso8601-converter --operation to-iso --input <unix_timestamp>")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  ISO-8601形式をUNIXタイムスタンプに変換:")
	fmt.Fprintln(w, "    iso8601-converter --operation to-unix --input <iso8601_time>")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  日付をUNIXタイムスタンプに変換 (UTC):")
	fmt.Fprintln(w, "    iso8601-converter --operation to-unix --input <date>")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  日付をUNIXタイムスタンプに変換 (JST):")
	fmt.Fprintln(w, "    iso8601-converter --operation to-unix --is-jst --input <date>")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "例:")
	fmt.Fprintln(w, "  iso8601-converter --operation now")
	fmt.Fprintln(w, "  iso8601-converter --operation to-iso --input 1619712000")
	fmt.Fprintln(w, "  iso8601-converter --operation to-unix --input 2021-04-30T00:00:00Z")
	fmt.Fprintln(w, "  iso8601-converter --operation to-unix --input 2021-04-30")
	fmt.Fprintln(w, "  iso8601-converter --operation to-unix --is-jst --input 2021-04-30")
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
