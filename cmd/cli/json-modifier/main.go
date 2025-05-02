package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/landmaster135/devbox/internal/interfaces/repositories"
	"github.com/landmaster135/devbox/internal/usecases/services"
)

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

// run はJSON修正ツールの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグセットを作成
	fs := flag.NewFlagSet("json-modifier", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// コマンドライン引数の定義
	filePath := fs.String("file", "", "操作するJSONファイルのパス")
	dirPath := fs.String("dir", "", "操作するJSONファイルが含まれるディレクトリのパス")
	recursive := fs.Bool("recursive", false, "ディレクトリを再帰的に処理する（-dirオプションと共に使用）")
	key := fs.String("key", "", "追加または取得するキー")
	value := fs.String("set", "", "追加する値")
	getFlag := fs.Bool("get", false, "指定されたキーの値を取得する")
	getAllFlag := fs.Bool("get-all", false, "すべてのキーと値を取得する")

	// 引数の解析
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// ファイルパスとディレクトリパスの両方が指定されていない場合はエラー
	if *filePath == "" && *dirPath == "" {
		fmt.Fprintln(stderr, "エラー: JSONファイルのパス（-file）またはディレクトリのパス（-dir）を指定してください")
		fs.Usage()
		return exitCodeError
	}

	// ファイルパスとディレクトリパスの両方が指定されている場合はエラー
	if *filePath != "" && *dirPath != "" {
		fmt.Fprintln(stderr, "エラー: -fileと-dirオプションは同時に指定できません")
		fs.Usage()
		return exitCodeError
	}

	// 依存関係の注入
	fileRepo := repositories.NewFileRepository()
	jsonService := services.NewJSONService(fileRepo)

	// 単一ファイルモードとディレクトリモードで処理を分岐
	if *filePath != "" {
		// 単一ファイルモード
		// コマンドに応じた処理を実行
		if *getFlag {
			// キーが指定されていない場合はエラー
			if *key == "" {
				fmt.Fprintln(stderr, "エラー: 取得するキーを指定してください（-key オプション）")
				fs.Usage()
				return exitCodeError
			}

			// 指定されたキーの値を取得
			value, err := jsonService.GetValue(*filePath, *key)
			if err != nil {
				fmt.Fprintf(stderr, "エラー: %v\n", err)
				return exitCodeError
			}

			// 値を表示
			fmt.Fprintf(stdout, "%v\n", value)
		} else if *getAllFlag {
			// すべてのキーと値を取得
			data, err := jsonService.GetAllData(*filePath)
			if err != nil {
				fmt.Fprintf(stderr, "エラー: %v\n", err)
				return exitCodeError
			}

			// JSONとして整形して表示
			jsonData, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				fmt.Fprintf(stderr, "エラー: JSONの整形に失敗しました: %v\n", err)
				return exitCodeError
			}

			fmt.Fprintln(stdout, string(jsonData))
		} else {
			// キーが指定されていない場合はエラー
			if *key == "" {
				fmt.Fprintln(stderr, "エラー: 追加するキーを指定してください（-key オプション）")
				fs.Usage()
				return exitCodeError
			}

			// 値が指定されていない場合はエラー
			if *value == "" {
				fmt.Fprintln(stderr, "エラー: 追加する値を指定してください（-set オプション）")
				fs.Usage()
				return exitCodeError
			}

			// 値が整数として解釈できるか確認
			var valueInterface interface{} = *value
			if intValue, err := strconv.ParseInt(*value, 10, 64); err == nil {
				valueInterface = intValue
			} else if floatValue, err := strconv.ParseFloat(*value, 64); err == nil {
				valueInterface = floatValue
			}

			// キーと値を追加
			err := jsonService.AddKeyValue(*filePath, *key, valueInterface)
			if err != nil {
				fmt.Fprintf(stderr, "エラー: %v\n", err)
				return exitCodeError
			}

			fmt.Fprintf(stdout, "JSONファイル '%s' にキー '%s' と値を追加しました\n", *filePath, *key)
		}
	} else {
		// ディレクトリモード
		// ディレクトリが存在するか確認
		fileInfo, err := os.Stat(*dirPath)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: ディレクトリ '%s' にアクセスできません: %v\n", *dirPath, err)
			return exitCodeError
		}
		if !fileInfo.IsDir() {
			fmt.Fprintf(stderr, "エラー: '%s' はディレクトリではありません\n", *dirPath)
			return exitCodeError
		}

		// getとget-allはディレクトリモードでは使用できない
		if *getFlag || *getAllFlag {
			fmt.Fprintln(stderr, "エラー: -getと-get-allオプションはディレクトリモードでは使用できません")
			fs.Usage()
			return exitCodeError
		}

		// キーが指定されていない場合はエラー
		if *key == "" {
			fmt.Fprintln(stderr, "エラー: 追加するキーを指定してください（-key オプション）")
			fs.Usage()
			return exitCodeError
		}

		// 値が指定されていない場合はエラー
		if *value == "" {
			fmt.Fprintln(stderr, "エラー: 追加する値を指定してください（-set オプション）")
			fs.Usage()
			return exitCodeError
		}

		// 値が整数として解釈できるか確認
		var valueInterface interface{} = *value
		if intValue, err := strconv.ParseInt(*value, 10, 64); err == nil {
			valueInterface = intValue
		} else if floatValue, err := strconv.ParseFloat(*value, 64); err == nil {
			valueInterface = floatValue
		}

		// ディレクトリ内の全てのJSONファイルにキーと値を追加
		count, err := jsonService.AddKeyValueToAllFiles(*dirPath, *key, valueInterface, *recursive)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: %v\n", err)
			return exitCodeError
		}

		fmt.Fprintf(stdout, "ディレクトリ '%s' 内の %d 個のJSONファイルにキー '%s' と値を追加しました\n", *dirPath, count, *key)
	}

	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
