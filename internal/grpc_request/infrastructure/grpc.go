package repositories

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"

	config "github.com/landmaster135/devbox/internal/grpc_request/config"
	grpcDomain "github.com/landmaster135/devbox/internal/grpc_request/domain"
)

// #==============================================================#
// ##       Mocks for HTTPClient                                 ##
// #==============================================================#
// MockGRPCRepository はGRPCRepositoryのモック実装
type MockGRPCRepository struct {
	SendRequestFunc    func(ctx context.Context, request *grpcDomain.GRPCRequest) (*grpcDomain.GRPCResponse, error)
	TestConnectionFunc func(ctx context.Context, serverAddress string, useTLS bool) error
	ListServicesFunc   func(ctx context.Context, serverAddress string, useTLS bool) ([]string, error)
}

func (m *MockGRPCRepository) SendRequest(ctx context.Context, request *grpcDomain.GRPCRequest) (*grpcDomain.GRPCResponse, error) {
	return m.SendRequestFunc(ctx, request)
}

func (m *MockGRPCRepository) TestConnection(ctx context.Context, serverAddress string, useTLS bool) error {
	return m.TestConnectionFunc(ctx, serverAddress, useTLS)
}

func (m *MockGRPCRepository) ListServices(ctx context.Context, serverAddress string, useTLS bool) ([]string, error) {
	return m.ListServicesFunc(ctx, serverAddress, useTLS)
}

// #==============================================================#
// ##       Interfaces for DiscordClient                         ##
// #==============================================================#
// GRPCRepository はgRPCリクエストを実行するリポジトリのインターフェースです
type GRPCRepository interface {
	SendRequest(ctx context.Context, request *grpcDomain.GRPCRequest) (*grpcDomain.GRPCResponse, error)
	TestConnection(ctx context.Context, serverAddress string, useTLS bool) error
	ListServices(ctx context.Context, serverAddress string, useTLS bool) ([]string, error)
}

// #==============================================================#
// ##       Implementations for DiscordClient                    ##
// #==============================================================#
// GRPCClient はGRPCClientの実装です
type GRPCClient struct {
	config *config.Config
}

// NewGRPCClient は新しいgRPCクライアントインスタンスを作成します
func NewGRPCClient(cfg *config.Config) GRPCRepository {
	return &GRPCClient{
		config: cfg,
	}
}

// SendRequest はgRPCリクエストを送信します
func (r *GRPCClient) SendRequest(ctx context.Context, request *grpcDomain.GRPCRequest) (*grpcDomain.GRPCResponse, error) {
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
			return &grpcDomain.GRPCResponse{
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

	return &grpcDomain.GRPCResponse{
		Data:       respData,
		Metadata:   responseMetadata,
		StatusCode: 0, // OK
		StatusMsg:  "OK",
		Duration:   duration,
	}, nil
}

// TestConnection はgRPCサーバーへの接続をテストします
func (r *GRPCClient) TestConnection(ctx context.Context, serverAddress string, useTLS bool) error {
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
func (r *GRPCClient) ListServices(ctx context.Context, serverAddress string, useTLS bool) ([]string, error) {
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
func (r *GRPCClient) createConnection(serverAddress string, useTLS bool) (*grpc.ClientConn, error) {
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
func (r *GRPCClient) getMethodDescriptor(ctx context.Context, client grpc_reflection_v1alpha.ServerReflectionClient, fullMethodName string) (protoreflect.ServiceDescriptor, protoreflect.MethodDescriptor, error) {
	// メソッド名を解析 (例: "weather_notificator.WeatherNotificatorService/SendWeatherNotification")
	parts := strings.Split(fullMethodName, "/")
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("無効なメソッド名の形式です: %s", fullMethodName)
	}

	serviceName := parts[0]
	methodName := parts[1]

	// リフレクションストリームを作成
	stream, err := client.ServerReflectionInfo(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("リフレクションストリームの作成に失敗しました: %w", err)
	}
	defer stream.CloseSend()

	// サービス記述子を要求
	req := &grpc_reflection_v1alpha.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_FileContainingSymbol{
			FileContainingSymbol: serviceName,
		},
	}

	if err := stream.Send(req); err != nil {
		return nil, nil, fmt.Errorf("サービス記述子の要求に失敗しました: %w", err)
	}

	resp, err := stream.Recv()
	if err != nil {
		return nil, nil, fmt.Errorf("サービス記述子の受信に失敗しました: %w", err)
	}

	// エラーレスポンスをチェック
	if errorResp := resp.GetErrorResponse(); errorResp != nil {
		return nil, nil, fmt.Errorf("サービス記述子の取得でエラーが発生しました: %s", errorResp.ErrorMessage)
	}

	// ファイル記述子レスポンスを取得
	fileDescResp := resp.GetFileDescriptorResponse()
	if fileDescResp == nil {
		return nil, nil, fmt.Errorf("ファイル記述子レスポンスが取得できませんでした")
	}

	// ファイル記述子を解析
	registry := &protoregistry.Files{}
	var fileDescriptors []*descriptorpb.FileDescriptorProto

	// まず全てのファイル記述子をアンマーシャル
	for _, fdBytes := range fileDescResp.FileDescriptorProto {
		fd := &descriptorpb.FileDescriptorProto{}
		if err := proto.Unmarshal(fdBytes, fd); err != nil {
			return nil, nil, fmt.Errorf("ファイル記述子のアンマーシャルに失敗しました: %w", err)
		}
		fileDescriptors = append(fileDescriptors, fd)
	}

	// 依存関係を考慮してファイル記述子を作成・登録
	opts := protodesc.FileOptions{AllowUnresolvable: true}
	for _, fd := range fileDescriptors {
		// ファイル記述子を作成
		fileDesc, err := opts.New(fd, registry)
		if err != nil {
			return nil, nil, fmt.Errorf("ファイル記述子の作成に失敗しました: %w", err)
		}

		// レジストリに登録
		if err := registry.RegisterFile(fileDesc); err != nil {
			return nil, nil, fmt.Errorf("ファイル記述子の登録に失敗しました: %w", err)
		}
	}

	// サービス記述子を検索
	serviceDesc, err := registry.FindDescriptorByName(protoreflect.FullName(serviceName))
	if err != nil {
		return nil, nil, fmt.Errorf("サービス記述子が見つかりません: %s", serviceName)
	}

	service, ok := serviceDesc.(protoreflect.ServiceDescriptor)
	if !ok {
		return nil, nil, fmt.Errorf("記述子がサービス記述子ではありません: %s", serviceName)
	}

	// メソッド記述子を検索
	methods := service.Methods()
	for i := 0; i < methods.Len(); i++ {
		method := methods.Get(i)
		if string(method.Name()) == methodName {
			return service, method, nil
		}
	}

	return nil, nil, fmt.Errorf("メソッドが見つかりません: %s", methodName)
}

// createRequestMessage はリクエストメッセージを作成します
func (r *GRPCClient) createRequestMessage(msgDesc protoreflect.MessageDescriptor, data map[string]interface{}) (proto.Message, error) {
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
func (r *GRPCClient) messageToMap(msg proto.Message) (map[string]interface{}, error) {
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
