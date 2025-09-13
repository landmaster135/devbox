package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/landmaster135/devbox/internal/grpc_request/config"
	"github.com/landmaster135/devbox/internal/grpc_request/domain/models"
	"github.com/landmaster135/devbox/internal/grpc_request/interfaces/repositories"
)

// GRPCService はgRPCリクエストのサービス層インターフェースです
type GRPCService interface {
	SendRequest(ctx context.Context, serverAddress, method, jsonFile string, metadata map[string]string, useTLS bool, timeout time.Duration) (*models.GRPCResponse, error)
	SendRequestWithData(ctx context.Context, request *models.GRPCRequest) (*models.GRPCResponse, error)
	TestConnection(ctx context.Context, serverAddress string, useTLS bool) error
	ListServices(ctx context.Context, serverAddress string, useTLS bool) ([]string, error)
	FormatResponse(response *models.GRPCResponse) (string, error)
	LoadJSONFile(filePath string) (map[string]interface{}, error)
}

// grpcService はGRPCServiceの実装です
type grpcService struct {
	repository repositories.GRPCRepository
	config     *config.Config
}

// NewGRPCService は新しいgRPCサービスインスタンスを作成します
func NewGRPCService(repo repositories.GRPCRepository, cfg *config.Config) GRPCService {
	return &grpcService{
		repository: repo,
		config:     cfg,
	}
}

// SendRequest はJSONファイルを使用してgRPCリクエストを送信します
func (s *grpcService) SendRequest(ctx context.Context, serverAddress, method, jsonFile string, metadata map[string]string, useTLS bool, timeout time.Duration) (*models.GRPCResponse, error) {
	// JSONファイルからデータを読み込み
	data, err := s.LoadJSONFile(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("JSONファイルの読み込みに失敗しました: %w", err)
	}

	// リクエストオブジェクトを作成
	request := &models.GRPCRequest{
		ServerAddress: serverAddress,
		Method:        method,
		Data:          data,
		Metadata:      metadata,
		UseTLS:        useTLS,
		Timeout:       timeout,
	}

	return s.SendRequestWithData(ctx, request)
}

// SendRequestWithData はリクエストオブジェクトを使用してgRPCリクエストを送信します
func (s *grpcService) SendRequestWithData(ctx context.Context, request *models.GRPCRequest) (*models.GRPCResponse, error) {
	// タイムアウトのデフォルト値を設定
	if request.Timeout == 0 {
		request.Timeout = s.config.DefaultTimeout
	}

	// リクエストの検証
	if err := s.validateRequest(request); err != nil {
		return nil, fmt.Errorf("リクエストの検証に失敗しました: %w", err)
	}

	// リポジトリを通じてリクエストを送信
	response, err := s.repository.SendRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("gRPCリクエストの送信に失敗しました: %w", err)
	}

	return response, nil
}

// TestConnection はgRPCサーバーへの接続をテストします
func (s *grpcService) TestConnection(ctx context.Context, serverAddress string, useTLS bool) error {
	return s.repository.TestConnection(ctx, serverAddress, useTLS)
}

// ListServices は利用可能なサービス一覧を取得します
func (s *grpcService) ListServices(ctx context.Context, serverAddress string, useTLS bool) ([]string, error) {
	return s.repository.ListServices(ctx, serverAddress, useTLS)
}

// FormatResponse はレスポンスを整形して文字列として返します
func (s *grpcService) FormatResponse(response *models.GRPCResponse) (string, error) {
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

// LoadJSONFile はJSONファイルを読み込んでmapとして返します
func (s *grpcService) LoadJSONFile(filePath string) (map[string]interface{}, error) {
	// ファイルを開く
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("ファイルのオープンに失敗しました: %w", err)
	}
	defer file.Close()

	// ファイル内容を読み込み
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("ファイルの読み込みに失敗しました: %w", err)
	}

	// JSONをパース
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("JSONのパースに失敗しました: %w", err)
	}

	return result, nil
}

// validateRequest はリクエストの妥当性を検証します
func (s *grpcService) validateRequest(request *models.GRPCRequest) error {
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
