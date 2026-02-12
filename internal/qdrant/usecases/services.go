package usecases

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	qdrant "github.com/qdrant/go-client/qdrant"

	"github.com/landmaster135/devbox/internal/qdrant/domain"
	veconfig "github.com/landmaster135/devbox/internal/vector_embedding/config"
	vedomain "github.com/landmaster135/devbox/internal/vector_embedding/domain"
	veusecases "github.com/landmaster135/devbox/internal/vector_embedding/usecases"
)

const (
	defaultRequestTimeout   = 60 * time.Second
	defaultEmbeddingTimeout = 60 * time.Second
)

// Service は Qdrant 関連ユースケースを提供する。
type Service struct {
	clientFactory    ClientFactory
	embeddingFactory EmbeddingServiceFactory
	requestTimeout   time.Duration
}

// ServiceOptions は Service の初期化オプション。
type ServiceOptions struct {
	ClientFactory    ClientFactory
	EmbeddingFactory EmbeddingServiceFactory
	RequestTimeout   time.Duration
}

// NewService はユースケースサービスを構築する。
func NewService(opts ServiceOptions) (*Service, error) {
	clientFactory := opts.ClientFactory
	if clientFactory == nil {
		clientFactory = defaultClientFactory{}
	}

	embeddingFactory := opts.EmbeddingFactory
	if embeddingFactory == nil {
		embeddingFactory = defaultEmbeddingFactory{}
	}

	timeout := opts.RequestTimeout
	if timeout <= 0 {
		timeout = defaultRequestTimeout
	}

	return &Service{
		clientFactory:    clientFactory,
		embeddingFactory: embeddingFactory,
		requestTimeout:   timeout,
	}, nil
}

// CreateCollection は Qdrant にコレクションを作成する。
func (s *Service) CreateCollection(ctx context.Context, params domain.CreateCollectionParams) error {
	client, err := s.clientFactory.NewClient(ClientOptions{Host: params.DBHost, Port: params.DBPort})
	if err != nil {
		return fmt.Errorf("Qdrant クライアントの初期化に失敗しました: %w", err)
	}
	defer client.Close()

	reqCtx, cancel := s.withTimeout(ctx)
	defer cancel()

	req := &qdrant.CreateCollection{
		CollectionName: params.CollectionName,
		VectorsConfig: &qdrant.VectorsConfig{
			Config: &qdrant.VectorsConfig_Params{
				Params: &qdrant.VectorParams{
					Size:     uint64(params.Size),
					Distance: qdrant.Distance_Cosine,
				},
			},
		},
	}

	if err := client.CreateCollection(reqCtx, req); err != nil {
		return fmt.Errorf("コレクションの作成に失敗しました: %w", err)
	}
	return nil
}

// ListCollections はコレクション一覧を取得する。
func (s *Service) ListCollections(ctx context.Context, params domain.ListCollectionsParams) ([]string, error) {
	client, err := s.clientFactory.NewClient(ClientOptions{Host: params.DBHost, Port: params.DBPort})
	if err != nil {
		return nil, fmt.Errorf("Qdrant クライアントの初期化に失敗しました: %w", err)
	}
	defer client.Close()

	reqCtx, cancel := s.withTimeout(ctx)
	defer cancel()

	collections, err := client.ListCollections(reqCtx)
	if err != nil {
		return nil, fmt.Errorf("コレクション一覧の取得に失敗しました: %w", err)
	}
	return collections, nil
}

// UpsertTexts は入力テキストのリストをベクトル化し、Qdrant に upsert する。
func (s *Service) UpsertTexts(ctx context.Context, params domain.UpsertTextsParams) (*qdrant.UpdateResult, error) {
	if len(params.Pairs) == 0 {
		return nil, fmt.Errorf("upsert する input/payload ペアがありません")
	}

	embeddings, err := s.embedAll(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(embeddings) != len(params.Pairs) {
		return nil, fmt.Errorf("埋め込み結果の件数が一致しません (%d vs %d)", len(embeddings), len(params.Pairs))
	}

	points := make([]*qdrant.PointStruct, 0, len(params.Pairs))
	for i, pair := range params.Pairs {
		vector := float64To32(embeddings[i])
		payload := mergePayload(pair.Input, pair.Payload)

		point := &qdrant.PointStruct{
			Id:      &qdrant.PointId{PointIdOptions: &qdrant.PointId_Uuid{Uuid: uuid.NewString()}},
			Payload: payload,
			Vectors: &qdrant.Vectors{
				VectorsOptions: &qdrant.Vectors_Vector{
					Vector: &qdrant.Vector{
						Vector: &qdrant.Vector_Dense{Dense: &qdrant.DenseVector{Data: vector}},
					},
				},
			},
		}
		points = append(points, point)
	}

	client, err := s.clientFactory.NewClient(ClientOptions{Host: params.DBHost, Port: params.DBPort})
	if err != nil {
		return nil, fmt.Errorf("Qdrant クライアントの初期化に失敗しました: %w", err)
	}
	defer client.Close()

	reqCtx, cancel := s.withTimeout(ctx)
	defer cancel()

	result, err := client.Upsert(reqCtx, &qdrant.UpsertPoints{
		CollectionName: params.CollectionName,
		Wait:           boolPtr(true),
		Points:         points,
	})
	if err != nil {
		return nil, fmt.Errorf("ポイントの upsert に失敗しました: %w", err)
	}
	return result, nil
}

// QueryPoints は1件の入力ベクトルで類似ポイントを検索する。
func (s *Service) QueryPoints(ctx context.Context, params domain.QueryPointsParams) ([]*qdrant.ScoredPoint, error) {
	embeddings, err := s.embedSingle(ctx, params)
	if err != nil {
		return nil, err
	}
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("埋め込み結果が空です")
	}

	vector := float64To32(embeddings[0])
	nearest := &qdrant.Query{
		Variant: &qdrant.Query_Nearest{
			Nearest: &qdrant.VectorInput{
				Variant: &qdrant.VectorInput_Dense{Dense: &qdrant.DenseVector{Data: vector}},
			},
		},
	}

	req := &qdrant.QueryPoints{
		CollectionName: params.CollectionName,
		Query:          nearest,
		WithPayload: &qdrant.WithPayloadSelector{
			SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true},
		},
	}

	if params.Limit > 0 {
		limit := uint64(params.Limit)
		req.Limit = &limit
	}

	if s.requestTimeout > 0 {
		timeout := uint64(math.Ceil(s.requestTimeout.Seconds()))
		if timeout > 0 {
			req.Timeout = &timeout
		}
	}

	if filters := buildFilters(params.PayloadFilters); filters != nil {
		req.Filter = &qdrant.Filter{Must: filters}
	}

	client, err := s.clientFactory.NewClient(ClientOptions{Host: params.DBHost, Port: params.DBPort})
	if err != nil {
		return nil, fmt.Errorf("Qdrant クライアントの初期化に失敗しました: %w", err)
	}
	defer client.Close()

	reqCtx, cancel := s.withTimeout(ctx)
	defer cancel()

	points, err := client.Query(reqCtx, req)
	if err != nil {
		return nil, fmt.Errorf("ポイントのクエリに失敗しました: %w", err)
	}
	return points, nil
}

func (s *Service) embedAll(ctx context.Context, params domain.UpsertTextsParams) ([][]float64, error) {
	svc, err := s.embeddingFactory.New(EmbeddingServiceOptions{
		Host:    params.Embedding.Host,
		Port:    params.Embedding.Port,
		Timeout: timeoutOrDefault(params.Embedding.TimeoutSeconds),
	})
	if err != nil {
		return nil, fmt.Errorf("埋め込みサービスの初期化に失敗しました: %w", err)
	}

	inputs := make([]string, 0, len(params.Pairs))
	for _, pair := range params.Pairs {
		inputs = append(inputs, pair.Input)
	}

	cfg := &veconfig.Config{
		Operation:      veconfig.OperationOllama,
		Host:           params.Embedding.Host,
		Port:           params.Embedding.Port,
		Model:          params.Embedding.Model,
		Inputs:         inputs,
		TimeoutSeconds: clampTimeoutSeconds(params.Embedding.TimeoutSeconds),
	}

	result, err := svc.Embed(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("埋め込みの実行に失敗しました: %w", err)
	}
	return result.Embeddings, nil
}

func (s *Service) embedSingle(ctx context.Context, params domain.QueryPointsParams) ([][]float64, error) {
	svc, err := s.embeddingFactory.New(EmbeddingServiceOptions{
		Host:    params.Embedding.Host,
		Port:    params.Embedding.Port,
		Timeout: timeoutOrDefault(params.Embedding.TimeoutSeconds),
	})
	if err != nil {
		return nil, fmt.Errorf("埋め込みサービスの初期化に失敗しました: %w", err)
	}

	cfg := &veconfig.Config{
		Operation:      veconfig.OperationOllama,
		Host:           params.Embedding.Host,
		Port:           params.Embedding.Port,
		Model:          params.Embedding.Model,
		Inputs:         []string{params.Input},
		TimeoutSeconds: clampTimeoutSeconds(params.Embedding.TimeoutSeconds),
	}

	result, err := svc.Embed(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("埋め込みの実行に失敗しました: %w", err)
	}
	return result.Embeddings, nil
}

func (s *Service) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if s.requestTimeout <= 0 {
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(ctx, s.requestTimeout)
}

func float64To32(in []float64) []float32 {
	out := make([]float32, len(in))
	for i, v := range in {
		out[i] = float32(v)
	}
	return out
}

func mergePayload(input string, base map[string]string) map[string]*qdrant.Value {
	result := make(map[string]*qdrant.Value)
	for k, v := range base {
		if k == "" {
			continue
		}
		result[k] = stringValue(v)
	}
	if input != "" {
		if _, exists := result["text"]; !exists {
			result["text"] = stringValue(input)
		}
	}
	return result
}

func stringValue(v string) *qdrant.Value {
	return &qdrant.Value{Kind: &qdrant.Value_StringValue{StringValue: v}}
}

func buildFilters(filters []domain.PayloadFilter) []*qdrant.Condition {
	if len(filters) == 0 {
		return nil
	}
	conditions := make([]*qdrant.Condition, 0, len(filters))
	for _, f := range filters {
		if f.Key == "" || f.Value == "" {
			continue
		}
		conditions = append(conditions, &qdrant.Condition{
			ConditionOneOf: &qdrant.Condition_Field{
				Field: &qdrant.FieldCondition{
					Key: f.Key,
					Match: &qdrant.Match{
						MatchValue: &qdrant.Match_Text{Text: f.Value},
					},
				},
			},
		})
	}
	if len(conditions) == 0 {
		return nil
	}
	return conditions
}

func timeoutOrDefault(seconds int) time.Duration {
	if seconds <= 0 {
		return defaultEmbeddingTimeout
	}
	return time.Duration(seconds) * time.Second
}

func clampTimeoutSeconds(value int) int {
	if value <= 0 {
		return int(defaultEmbeddingTimeout.Seconds())
	}
	return value
}

func boolPtr(v bool) *bool {
	return &v
}

// ----- Qdrant クライアントファクトリ -----

// ClientOptions は Qdrant クライアントを作る際の情報。
type ClientOptions struct {
	Host   string
	Port   int
	APIKey string
	UseTLS bool
}

// QdrantClient は必要な Qdrant API を表す。
type QdrantClient interface {
	Close() error
	CreateCollection(ctx context.Context, request *qdrant.CreateCollection) error
	ListCollections(ctx context.Context) ([]string, error)
	Upsert(ctx context.Context, request *qdrant.UpsertPoints) (*qdrant.UpdateResult, error)
	Query(ctx context.Context, request *qdrant.QueryPoints) ([]*qdrant.ScoredPoint, error)
}

// ClientFactory は Qdrant クライアントを生成する。
type ClientFactory interface {
	NewClient(opts ClientOptions) (QdrantClient, error)
}

type defaultClientFactory struct{}

func (defaultClientFactory) NewClient(opts ClientOptions) (QdrantClient, error) {
	cfg := &qdrant.Config{
		Host:   opts.Host,
		Port:   opts.Port,
		APIKey: opts.APIKey,
		UseTLS: opts.UseTLS,
	}
	client, err := qdrant.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &wrappedQdrantClient{client: client}, nil
}

type wrappedQdrantClient struct {
	client *qdrant.Client
}

func (w *wrappedQdrantClient) Close() error {
	return w.client.Close()
}

func (w *wrappedQdrantClient) CreateCollection(ctx context.Context, request *qdrant.CreateCollection) error {
	return w.client.CreateCollection(ctx, request)
}

func (w *wrappedQdrantClient) ListCollections(ctx context.Context) ([]string, error) {
	return w.client.ListCollections(ctx)
}

func (w *wrappedQdrantClient) Upsert(ctx context.Context, request *qdrant.UpsertPoints) (*qdrant.UpdateResult, error) {
	return w.client.Upsert(ctx, request)
}

func (w *wrappedQdrantClient) Query(ctx context.Context, request *qdrant.QueryPoints) ([]*qdrant.ScoredPoint, error) {
	return w.client.Query(ctx, request)
}

// ----- 埋め込みサービスファクトリ -----

// EmbeddingService は Embed を提供するインターフェース。
type EmbeddingService interface {
	Embed(ctx context.Context, cfg *veconfig.Config) (*vedomain.EmbedResult, error)
}

// EmbeddingServiceOptions は埋め込みサービス生成用オプション。
type EmbeddingServiceOptions struct {
	Host    string
	Port    int
	Timeout time.Duration
}

// EmbeddingServiceFactory は埋め込みサービスを生成する。
type EmbeddingServiceFactory interface {
	New(opts EmbeddingServiceOptions) (EmbeddingService, error)
}

type defaultEmbeddingFactory struct{}

func (defaultEmbeddingFactory) New(opts EmbeddingServiceOptions) (EmbeddingService, error) {
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = defaultEmbeddingTimeout
	}
	return veusecases.NewService(veusecases.Options{
		Host:    opts.Host,
		Port:    opts.Port,
		Timeout: timeout,
	})
}
