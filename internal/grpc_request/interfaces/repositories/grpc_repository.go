package repositories

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	"github.com/landmaster135/devbox/internal/grpc_request/config"
	"github.com/landmaster135/devbox/internal/grpc_request/domain/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// GRPCRepository はgRPCリクエストを実行するリポジトリのインターフェースです
type GRPCRepository interface {
	SendRequest(ctx context.Context, request *models.GRPCRequest) (*models.GRPCResponse, error)
	TestConnection(ctx context.Context, serverAddress string, useTLS bool) error
	ListServices(ctx context.Context, serverAddress string, useTLS bool) ([]string, error)
}

// grpcRepository はGRPCRepositoryの実装です
type grpcRepository struct {
	config *config.Config
}

// NewGRPCRepository は新しいgRPCリポジトリインスタンスを作成します
func NewGRPCRepository(cfg *config.Config) GRPCRepository {
	return &grpcRepository{
		config: cfg,
	}
}

// SendRequest はgRPCリクエストを送信します
func (r *grpcRepository) SendRequest(ctx context.Context, request *models.GRPCRequest) (*models.GRPCResponse, error) {
	start := time.Now()

	// gRPC接続を確立
	conn, err := r.createConnection(request.ServerAddress, request.UseTLS)
	if err != nil {
		return nil, fmt.Errorf("接続の確立に失敗しました: %w", err)
	}
	defer conn.Close()

	// タイムアウトの設定
	if request.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, request.Timeout)
		defer cancel()
	}

	// メタデータの設定
	if len(request.Metadata) > 0 {
		md := metadata.New(request.Metadata)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	// リフレクションクライアントを作成
	reflectionClient := grpc_reflection_v1alpha.NewServerReflectionClient(conn)

	// サービスとメソッドの情報を取得
	serviceDesc, methodDesc, err := r.getMethodDescriptor(ctx, reflectionClient, request.Method)
	if err != nil {
		return nil, fmt.Errorf("メソッド情報の取得に失敗しました: %w", err)
	}

	// リクエストメッセージを作成
	reqMsg, err := r.createRequestMessage(methodDesc.Input(), request.Data)
	if err != nil {
		return nil, fmt.Errorf("リクエストメッセージの作成に失敗しました: %w", err)
	}

	// レスポンスメッセージを作成
	respMsg := dynamicpb.NewMessage(methodDesc.Output())

	// gRPCメソッドを呼び出し
	err = conn.Invoke(ctx, fmt.Sprintf("/%s/%s", serviceDesc.FullName(), methodDesc.Name()), reqMsg, respMsg)

	duration := time.Since(start)

	// レスポンスヘッダーとトレーラーを取得
	var responseMetadata map[string][]string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		responseMetadata = map[string][]string(md)
	}

	if err != nil {
		// gRPCエラーの処理
		if st, ok := status.FromError(err); ok {
			return &models.GRPCResponse{
				StatusCode: int(st.Code()),
				StatusMsg:  st.Message(),
				Metadata:   responseMetadata,
				Duration:   duration,
			}, nil
		}
		return nil, fmt.Errorf("gRPCリクエストの実行に失敗しました: %w", err)
	}

	// レスポンスデータをJSONに変換
	respData, err := r.messageToMap(respMsg)
	if err != nil {
		return nil, fmt.Errorf("レスポンスデータの変換に失敗しました: %w", err)
	}

	return &models.GRPCResponse{
		Data:       respData,
		Metadata:   responseMetadata,
		StatusCode: 0, // OK
		StatusMsg:  "OK",
		Duration:   duration,
	}, nil
}

// TestConnection はgRPCサーバーへの接続をテストします
func (r *grpcRepository) TestConnection(ctx context.Context, serverAddress string, useTLS bool) error {
	conn, err := r.createConnection(serverAddress, useTLS)
	if err != nil {
		return fmt.Errorf("接続テストに失敗しました: %w", err)
	}
	defer conn.Close()

	// 簡単な接続テスト（リフレクションサービスの呼び出し）
	reflectionClient := grpc_reflection_v1alpha.NewServerReflectionClient(conn)
	stream, err := reflectionClient.ServerReflectionInfo(ctx)
	if err != nil {
		return fmt.Errorf("リフレクションサービスへの接続に失敗しました: %w", err)
	}
	defer stream.CloseSend()

	return nil
}

// ListServices は利用可能なサービス一覧を取得します
func (r *grpcRepository) ListServices(ctx context.Context, serverAddress string, useTLS bool) ([]string, error) {
	conn, err := r.createConnection(serverAddress, useTLS)
	if err != nil {
		return nil, fmt.Errorf("接続の確立に失敗しました: %w", err)
	}
	defer conn.Close()

	reflectionClient := grpc_reflection_v1alpha.NewServerReflectionClient(conn)
	stream, err := reflectionClient.ServerReflectionInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("リフレクションサービスへの接続に失敗しました: %w", err)
	}
	defer stream.CloseSend()

	// サービス一覧を要求
	req := &grpc_reflection_v1alpha.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_ListServices{
			ListServices: "",
		},
	}

	if err := stream.Send(req); err != nil {
		return nil, fmt.Errorf("サービス一覧の要求に失敗しました: %w", err)
	}

	resp, err := stream.Recv()
	if err != nil {
		return nil, fmt.Errorf("サービス一覧の受信に失敗しました: %w", err)
	}

	listServicesResp := resp.GetListServicesResponse()
	if listServicesResp == nil {
		return nil, fmt.Errorf("サービス一覧の取得に失敗しました")
	}

	var services []string
	for _, service := range listServicesResp.Service {
		services = append(services, service.Name)
	}

	return services, nil
}

// createConnection はgRPC接続を作成します
func (r *grpcRepository) createConnection(serverAddress string, useTLS bool) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	if useTLS {
		// TLS接続の設定
		creds := credentials.NewTLS(&tls.Config{
			ServerName: serverAddress,
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		// 非セキュア接続
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// 接続を確立
	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("gRPC接続の確立に失敗しました: %w", err)
	}

	return conn, nil
}

// getMethodDescriptor はメソッドの記述子を取得します
func (r *grpcRepository) getMethodDescriptor(ctx context.Context, client grpc_reflection_v1alpha.ServerReflectionClient, fullMethodName string) (protoreflect.ServiceDescriptor, protoreflect.MethodDescriptor, error) {
	// この実装は簡略化されています
	// 実際の実装では、リフレクションAPIを使用してサービスとメソッドの情報を動的に取得する必要があります
	return nil, nil, fmt.Errorf("メソッド記述子の取得は未実装です")
}

// createRequestMessage はリクエストメッセージを作成します
func (r *grpcRepository) createRequestMessage(msgDesc protoreflect.MessageDescriptor, data map[string]interface{}) (proto.Message, error) {
	msg := dynamicpb.NewMessage(msgDesc)

	// JSONデータをprotobufメッセージに変換
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("JSONデータのマーシャルに失敗しました: %w", err)
	}

	err = protojson.Unmarshal(jsonData, msg)
	if err != nil {
		return nil, fmt.Errorf("protobufメッセージへの変換に失敗しました: %w", err)
	}

	return msg, nil
}

// messageToMap はprotobufメッセージをmapに変換します
func (r *grpcRepository) messageToMap(msg proto.Message) (map[string]interface{}, error) {
	jsonData, err := protojson.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("protobufメッセージのマーシャルに失敗しました: %w", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, fmt.Errorf("JSONデータのアンマーシャルに失敗しました: %w", err)
	}

	return result, nil
}
