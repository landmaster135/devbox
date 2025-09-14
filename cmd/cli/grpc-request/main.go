package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	config "github.com/landmaster135/devbox/internal/grpc_request/config"
	grpcInfra "github.com/landmaster135/devbox/internal/grpc_request/infrastructure"
	usecases "github.com/landmaster135/devbox/internal/grpc_request/usecases"
)

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

// run はgRPCクライアントの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグセットを作成
	fs := flag.NewFlagSet("grpc-request", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// コマンドライン引数の定義
	server := fs.String("server", "", "gRPCサーバーのアドレス（例: localhost:50051）")
	method := fs.String("method", "", "呼び出すメソッド（例: package.Service/Method）")
	jsonFile := fs.String("data", "", "リクエストデータのJSONファイルパス")
	useTLS := fs.Bool("tls", false, "TLS接続を使用する")
	token := fs.String("token", "", "認証トークン（メタデータとして送信）")
	timeout := fs.Duration("timeout", 30*time.Second, "リクエストタイムアウト")
	testConn := fs.Bool("test", false, "接続テストのみを実行する")
	listServices := fs.Bool("list", false, "利用可能なサービス一覧を表示する")

	// 引数の解析
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// サーバーアドレスが指定されていない場合はエラー
	if *server == "" {
		fmt.Fprintln(stderr, "エラー: サーバーアドレスを指定してください（-server オプション）")
		fs.Usage()
		return exitCodeError
	}

	// 依存関係の注入
	cfg := config.NewConfig()
	grpcRepo := grpcInfra.NewGRPCClient(cfg)
	grpcService := usecases.NewGRPCService(grpcRepo, cfg)

	ctx := context.Background()

	// 接続テストモード
	if *testConn {
		if err := grpcService.GetRepository().TestConnection(ctx, *server, *useTLS); err != nil {
			fmt.Fprintf(stderr, "エラー: 接続テストに失敗しました: %v\n", err)
			return exitCodeError
		}
		fmt.Fprintln(stdout, "接続テスト成功")
		return exitCodeOK
	}

	// サービス一覧表示モード
	if *listServices {
		services, err := grpcService.GetRepository().ListServices(ctx, *server, *useTLS)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: サービス一覧の取得に失敗しました: %v\n", err)
			return exitCodeError
		}

		fmt.Fprintln(stdout, "利用可能なサービス:")
		for _, service := range services {
			fmt.Fprintf(stdout, "  %s\n", service)
		}
		return exitCodeOK
	}

	// 通常のリクエストモード
	if *method == "" {
		fmt.Fprintln(stderr, "エラー: メソッドを指定してください（-method オプション）")
		fs.Usage()
		return exitCodeError
	}

	if *jsonFile == "" {
		fmt.Fprintln(stderr, "エラー: リクエストデータのJSONファイルを指定してください（-data オプション）")
		fs.Usage()
		return exitCodeError
	}

	// メタデータの準備
	metadata := make(map[string]string)
	if *token != "" {
		metadata["authorization"] = "Bearer " + *token
	}

	// gRPCリクエストを送信
	response, err := grpcService.SendRequest(ctx, *server, *method, *jsonFile, metadata, *useTLS, *timeout)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		return exitCodeError
	}

	// レスポンスを整形して表示
	formattedResponse, err := grpcService.FormatResponse(response)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: レスポンスの整形に失敗しました: %v\n", err)
		return exitCodeError
	}

	fmt.Fprintln(stdout, formattedResponse)
	return exitCodeOK
}

func main() {
	// ヘルプメッセージをカスタマイズ
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "gRPCリクエストを送信するCLIツール")
		fmt.Fprintln(os.Stderr, "\nオプション:")
		fmt.Fprintln(os.Stderr, "  -server string")
		fmt.Fprintln(os.Stderr, "        gRPCサーバーのアドレス（例: localhost:50051）")
		fmt.Fprintln(os.Stderr, "  -method string")
		fmt.Fprintln(os.Stderr, "        呼び出すメソッド（例: package.Service/Method）")
		fmt.Fprintln(os.Stderr, "  -data string")
		fmt.Fprintln(os.Stderr, "        リクエストデータのJSONファイルパス")
		fmt.Fprintln(os.Stderr, "  -tls")
		fmt.Fprintln(os.Stderr, "        TLS接続を使用する")
		fmt.Fprintln(os.Stderr, "  -token string")
		fmt.Fprintln(os.Stderr, "        認証トークン（メタデータとして送信）")
		fmt.Fprintln(os.Stderr, "  -timeout duration")
		fmt.Fprintln(os.Stderr, "        リクエストタイムアウト（デフォルト: 30s）")
		fmt.Fprintln(os.Stderr, "  -test")
		fmt.Fprintln(os.Stderr, "        接続テストのみを実行する")
		fmt.Fprintln(os.Stderr, "  -list")
		fmt.Fprintln(os.Stderr, "        利用可能なサービス一覧を表示する")
		fmt.Fprintln(os.Stderr, "\n使用例:")
		fmt.Fprintln(os.Stderr, "  # 接続テスト")
		fmt.Fprintln(os.Stderr, "  grpc-request -server localhost:50051 -test")
		fmt.Fprintln(os.Stderr, "  # サービス一覧表示")
		fmt.Fprintln(os.Stderr, "  grpc-request -server localhost:50051 -list")
		fmt.Fprintln(os.Stderr, "  # gRPCリクエスト送信")
		fmt.Fprintln(os.Stderr, "  grpc-request -server localhost:50051 -method package.Service/Method -data request.json")
		fmt.Fprintln(os.Stderr, "  # TLS接続でリクエスト送信")
		fmt.Fprintln(os.Stderr, "  grpc-request -server example.com:443 -method package.Service/Method -data request.json -tls")
		fmt.Fprintln(os.Stderr, "  # 認証トークン付きでリクエスト送信")
		fmt.Fprintln(os.Stderr, "  grpc-request -server localhost:50051 -method package.Service/Method -data request.json -token your_token")
	}

	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
