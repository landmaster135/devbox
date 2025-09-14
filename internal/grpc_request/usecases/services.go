package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	config "github.com/landmaster135/devbox/internal/grpc_request/config"
	grpcDomain "github.com/landmaster135/devbox/internal/grpc_request/domain"
	grpcInfra "github.com/landmaster135/devbox/internal/grpc_request/infrastructure"
)

// #==============================================================#
// ##       Interfaces for DiscordClient                         ##
// #==============================================================#
// GRPCServiceRepository はgRPCリクエストのサービス層インターフェースです
type GRPCServiceRepository interface {
	GetRepository() grpcInfra.GRPCClientRepository
	GetConfig() *config.Config
	GetFileReader() grpcInfra.FileReaderRepository
	SendRequest(ctx context.Context, serverAddress, method, jsonFile string, metadata map[string]string, useTLS bool, timeout time.Duration) (*grpcDomain.GRPCResponse, error)
	SendRequestWithData(ctx context.Context, request *grpcDomain.GRPCRequest) (*grpcDomain.GRPCResponse, error)
	FormatResponse(response *grpcDomain.GRPCResponse) (string, error)
	ExecuteCLICommand(ctx context.Context, options *CLIOptions) (string, error)
}

// #==============================================================#
// ##       Implementations for DiscordClient                    ##
// #==============================================================#
// GRPCService はGRPCServiceRepositoryの実装です
type GRPCService struct {
	client     grpcInfra.GRPCClientRepository
	config     *config.Config
	fileReader grpcInfra.FileReaderRepository
}

// NewGRPCService は新しいgRPCサービスインスタンスを作成します
func NewGRPCService(client grpcInfra.GRPCClientRepository, cfg *config.Config, fileReader grpcInfra.FileReaderRepository) GRPCServiceRepository {
	return &GRPCService{
		client:     client,
		config:     cfg,
		fileReader: fileReader,
	}
}

func (g *GRPCService) GetRepository() grpcInfra.GRPCClientRepository {
	return g.client
}

func (g *GRPCService) GetConfig() *config.Config {
	return g.config
}

func (g *GRPCService) GetFileReader() grpcInfra.FileReaderRepository {
	return g.fileReader
}

// SendRequest はJSONファイルを使用してgRPCリクエストを送信します
func (s *GRPCService) SendRequest(ctx context.Context, serverAddress, method, jsonFile string, metadata map[string]string, useTLS bool, timeout time.Duration) (*grpcDomain.GRPCResponse, error) {
	// JSONファイルからデータを読み込み
	data, err := s.fileReader.LoadJSONFile(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("JSONファイルの読み込みに失敗しました: %w", err)
	}

	// リクエストオブジェクトを作成
	request := grpcDomain.NewGRPCRequest(serverAddress, method, data, metadata, useTLS, timeout)

	return s.SendRequestWithData(ctx, request)
}

// TODO: これって必要なの？
// SendRequestWithData はリクエストオブジェクトを使用してgRPCリクエストを送信します
func (s *GRPCService) SendRequestWithData(ctx context.Context, request *grpcDomain.GRPCRequest) (*grpcDomain.GRPCResponse, error) {
	// タイムアウトのデフォルト値を設定
	if request.Timeout == 0 {
		request.Timeout = s.config.DefaultTimeout
	}

	// リクエストの検証
	if err := s.validateRequest(request); err != nil {
		return nil, fmt.Errorf("リクエストの検証に失敗しました: %w", err)
	}

	// リポジトリを通じてリクエストを送信
	response, err := s.client.SendRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("gRPCリクエストの送信に失敗しました: %w", err)
	}

	return response, nil
}

// FormatResponse はレスポンスを整形して文字列として返します
func (s *GRPCService) FormatResponse(response *grpcDomain.GRPCResponse) (string, error) {
	var result strings.Builder

	// ステータス情報
	result.WriteString(fmt.Sprintf("Status: %d - %s\n", response.StatusCode, response.StatusMsg))
	result.WriteString(fmt.Sprintf("Duration: %v\n", response.Duration))
	result.WriteString("\n")

	// メタデータ
	if len(response.Metadata) > 0 {
		result.WriteString("Metadata:\n")
		for key, values := range response.Metadata {
			for _, value := range values {
				result.WriteString(fmt.Sprintf("  %s: %s\n", key, value))
			}
		}
		result.WriteString("\n")
	}

	// レスポンスデータ
	if response.Data != nil {
		result.WriteString("Response Data:\n")
		jsonData, err := json.MarshalIndent(response.Data, "", "  ")
		if err != nil {
			return "", fmt.Errorf("レスポンスデータのJSON変換に失敗しました: %w", err)
		}
		result.WriteString(string(jsonData))
		result.WriteString("\n")
	}

	return result.String(), nil
}

// validateRequest はリクエストの妥当性を検証します
func (s *GRPCService) validateRequest(request *grpcDomain.GRPCRequest) error {
	if request.ServerAddress == "" {
		return fmt.Errorf("サーバーアドレスが指定されていません")
	}

	if request.Method == "" {
		return fmt.Errorf("メソッドが指定されていません")
	}

	// メソッド名の形式をチェック（package.Service/Method）
	parts := strings.Split(request.Method, "/")
	if len(parts) != 2 {
		return fmt.Errorf("メソッド名の形式が正しくありません。package.Service/Method の形式で指定してください")
	}

	serviceParts := strings.Split(parts[0], ".")
	if len(serviceParts) < 2 {
		return fmt.Errorf("サービス名の形式が正しくありません。package.Service の形式で指定してください")
	}

	if parts[1] == "" {
		return fmt.Errorf("メソッド名が空です")
	}

	return nil
}

// ExecuteCLICommand はCLIオプションに基づいてコマンドを実行します
func (s *GRPCService) ExecuteCLICommand(ctx context.Context, options *CLIOptions) (string, error) {
	// オプションの検証
	if err := options.Validate(); err != nil {
		return "", err
	}

	// 接続テストモード
	if options.TestConn {
		if err := s.client.TestConnection(ctx, options.Server, options.UseTLS); err != nil {
			return "", fmt.Errorf("接続テストに失敗しました: %w", err)
		}
		return "接続テスト成功", nil
	}

	// サービス一覧表示モード
	if options.ListServices {
		services, err := s.client.ListServices(ctx, options.Server, options.UseTLS)
		if err != nil {
			return "", fmt.Errorf("サービス一覧の取得に失敗しました: %w", err)
		}

		var result strings.Builder
		result.WriteString("利用可能なサービス:\n")
		for _, service := range services {
			result.WriteString(fmt.Sprintf("  %s\n", service))
		}
		return result.String(), nil
	}

	// 通常のリクエストモード
	// メタデータの準備
	metadata := make(map[string]string)
	if options.Token != "" {
		metadata["authorization"] = "Bearer " + options.Token
	}

	// gRPCリクエストを送信
	response, err := s.SendRequest(ctx, options.Server, options.Method, options.JSONFile, metadata, options.UseTLS, options.Timeout)
	if err != nil {
		return "", err
	}

	// レスポンスを整形して返す
	return s.FormatResponse(response)
}
