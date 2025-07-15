package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/landmaster135/devbox/internal/http_request/domain/models"
	"github.com/landmaster135/devbox/internal/http_request/interfaces/repositories"
	"github.com/landmaster135/devbox/internal/http_request/usecases/services"
)

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

// run はAPIクライアントの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグセットを作成
	fs := flag.NewFlagSet("http-request", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// コマンドライン引数の定義
	url := fs.String("url", "", "リクエスト先のURL")
	method := fs.String("method", "GET", "HTTPメソッド（GET, POST, PUT, DELETE, etc.）")
	jsonFile := fs.String("json", "", "リクエストボディとして送信するJSONファイルのパス")
	token := fs.String("token", "", "認証トークン（Bearer トークンとして送信されます）")
	encoding := fs.String("encoding", "auto", "文字エンコーディング（shift_jis, utf-8, euc-jp, auto）")

	// 引数の解析
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// URLが指定されていない場合はエラー
	if *url == "" {
		fmt.Fprintln(stderr, "エラー: URLを指定してください（-url オプション）")
		fs.Usage()
		return exitCodeError
	}

	// POSTやPUTの場合はJSONファイルが必要
	if (*method == "POST" || *method == "PUT" || *method == "PATCH") && *jsonFile == "" {
		fmt.Fprintf(stderr, "エラー: %sメソッドにはJSONファイルが必要です（-json オプション）\n", *method)
		fs.Usage()
		return exitCodeError
	}

	// 依存関係の注入
	apiRepo := repositories.NewHTTPRepository()
	apiService := services.NewHTTPService(apiRepo)

	// APIリクエストを送信
	var response *models.HTTPResponse
	var err error

	// ヘッダーの準備
	headers := map[string]string{"Accept": "application/json"}

	// トークンが指定されている場合は、Authorizationヘッダーを追加
	if *token != "" {
		headers["Authorization"] = "Bearer " + *token
	}

	if *jsonFile != "" {
		// JSONファイルを使用してリクエストを送信
		response, err = apiService.SendRequestWithJSONFileAndHeaders(*url, *method, *jsonFile, headers, *encoding)
	} else {
		// JSONファイルなしでリクエストを送信（GETなど）
		response, err = apiService.SendRequestWithoutJSONFile(*url, *method, headers, *encoding)
	}

	if err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		return exitCodeError
	}

	// レスポンスを整形して表示
	formattedResponse, err := apiService.FormatResponse(response)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: レスポンスの整形に失敗しました: %v\n", err)
		return exitCodeError
	}

	fmt.Fprintln(stdout, formattedResponse)
	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
