package usecases

import (
	"context"
	"testing"
	"time"

	qdrant "github.com/qdrant/go-client/qdrant"

	"github.com/landmaster135/devbox/internal/qdrant/domain"
	veconfig "github.com/landmaster135/devbox/internal/vector_embedding/config"
	vedomain "github.com/landmaster135/devbox/internal/vector_embedding/domain"
)

func TestService_UpsertTexts(t *testing.T) {
	fakeClient := &fakeQdrantClient{}
	fakeEmbeddings := [][]float64{{0.1, 0.2, 0.3}}

	svc, err := NewService(ServiceOptions{
		ClientFactory:    &fakeClientFactory{client: fakeClient},
		EmbeddingFactory: &fakeEmbeddingFactory{embeddings: fakeEmbeddings},
		RequestTimeout:   time.Second,
	})
	if err != nil {
		t.Fatalf("サービス初期化に失敗: %v", err)
	}

	topic := "greeting"
	params := domain.UpsertTextsParams{
		DBHost:         "localhost",
		DBPort:         6334,
		CollectionName: "demo",
		Embedding: domain.EmbeddingConfig{
			Host: "embed", Port: 11434, Model: "test-model",
		},
		Pairs: []domain.InputPayloadPair{
			{
				Input:   "hello world",
				Payload: map[string]string{"topic": topic},
			},
		},
	}

	opID := uint64(42)
	fakeClient.upsertResult = &qdrant.UpdateResult{OperationId: &opID}

	result, err := svc.UpsertTexts(context.Background(), params)
	if err != nil {
		t.Fatalf("UpsertTexts が失敗: %v", err)
	}
	if result.GetOperationId() != opID {
		t.Fatalf("operation id が一致しません: got %d", result.GetOperationId())
	}

	if fakeClient.upsertReq == nil {
		t.Fatalf("Upsert リクエストが送信されていません")
	}
	if got := fakeClient.upsertReq.CollectionName; got != "demo" {
		t.Fatalf("collection name mismatch: %s", got)
	}
	if len(fakeClient.upsertReq.Points) != 1 {
		t.Fatalf("point 数が不正: %d", len(fakeClient.upsertReq.Points))
	}

	payload := fakeClient.upsertReq.Points[0].GetPayload()
	if payload["topic"].GetStringValue() != topic {
		t.Fatalf("payload topic が不一致: %s", payload["topic"].GetStringValue())
	}
	if payload["text"].GetStringValue() != "hello world" {
		t.Fatalf("payload text が不一致: %s", payload["text"].GetStringValue())
	}

	vector := fakeClient.upsertReq.Points[0].GetVectors().GetVector().GetDense().GetData()
	if len(vector) != len(fakeEmbeddings[0]) {
		t.Fatalf("ベクトル長が一致しない: got %d", len(vector))
	}
	if !fakeClient.closed {
		t.Fatalf("クライアントが Close されていません")
	}
}

func TestService_QueryPoints(t *testing.T) {
	fakeClient := &fakeQdrantClient{}
	fakeEmbeddings := [][]float64{{0.5, 0.6}}

	svc, err := NewService(ServiceOptions{
		ClientFactory:    &fakeClientFactory{client: fakeClient},
		EmbeddingFactory: &fakeEmbeddingFactory{embeddings: fakeEmbeddings},
		RequestTimeout:   2 * time.Second,
	})
	if err != nil {
		t.Fatalf("サービス初期化に失敗: %v", err)
	}

	fakeClient.queryResult = []*qdrant.ScoredPoint{{Score: 0.9}}

	params := domain.QueryPointsParams{
		DBHost:         "localhost",
		DBPort:         6334,
		CollectionName: "demo",
		Embedding: domain.EmbeddingConfig{
			Host: "embed", Port: 11434, Model: "llm-ja",
		},
		Input:          "hello",
		PayloadFilters: []domain.PayloadFilter{{Key: "topic", Value: "greeting"}},
		Limit:          3,
	}

	points, err := svc.QueryPoints(context.Background(), params)
	if err != nil {
		t.Fatalf("QueryPoints が失敗: %v", err)
	}
	if len(points) != 1 {
		t.Fatalf("取得件数が不正: %d", len(points))
	}

	req := fakeClient.queryReq
	if req == nil {
		t.Fatalf("Query リクエストが送信されていません")
	}
	if req.GetCollectionName() != "demo" {
		t.Fatalf("collection name mismatch")
	}
	if req.GetFilter() == nil || len(req.GetFilter().GetMust()) != 1 {
		t.Fatalf("payload filter が設定されていません")
	}
	cond := req.GetFilter().GetMust()[0].GetField()
	if cond.GetKey() != "topic" {
		t.Fatalf("filter key mismatch: %s", cond.GetKey())
	}
	if cond.GetMatch().GetKeyword() != "greeting" {
		t.Fatalf("filter value mismatch: %s", cond.GetMatch().GetKeyword())
	}
	if !fakeClient.closed {
		t.Fatalf("クライアントが Close されていません")
	}
}

func TestService_DescribeCollection(t *testing.T) {
	fakeClient := &fakeQdrantClient{collectionInfo: &qdrant.CollectionInfo{}}
	svc, err := NewService(ServiceOptions{ClientFactory: &fakeClientFactory{client: fakeClient}})
	if err != nil {
		t.Fatalf("サービス初期化に失敗: %v", err)
	}

	info, err := svc.DescribeCollection(context.Background(), domain.DescribeCollectionParams{
		DBHost:         "localhost",
		DBPort:         6334,
		CollectionName: "demo",
	})
	if err != nil {
		t.Fatalf("DescribeCollection が失敗: %v", err)
	}
	if info != fakeClient.collectionInfo {
		t.Fatalf("期待した collection info が返っていません")
	}
	if fakeClient.describeCollectionName != "demo" {
		t.Fatalf("describe 用 collection name が不正: %s", fakeClient.describeCollectionName)
	}
	if !fakeClient.closed {
		t.Fatalf("クライアントが Close されていません")
	}
}

func TestService_DeleteCollection(t *testing.T) {
	fakeClient := &fakeQdrantClient{}
	svc, err := NewService(ServiceOptions{ClientFactory: &fakeClientFactory{client: fakeClient}})
	if err != nil {
		t.Fatalf("サービス初期化に失敗: %v", err)
	}

	params := domain.DeleteCollectionParams{
		DBHost:         "localhost",
		DBPort:         6334,
		CollectionName: "demo",
	}

	if err := svc.DeleteCollection(context.Background(), params); err != nil {
		t.Fatalf("DeleteCollection が失敗: %v", err)
	}
	if fakeClient.deleteCollectionName != "demo" {
		t.Fatalf("削除対象の collection name が不正: %s", fakeClient.deleteCollectionName)
	}
	if !fakeClient.closed {
		t.Fatalf("クライアントが Close されていません")
	}
}

func TestService_OverwritePayload(t *testing.T) {
	t.Run("PointIDsSelector", func(t *testing.T) {
		fakeClient := &fakeQdrantClient{}
		svc, err := NewService(ServiceOptions{ClientFactory: &fakeClientFactory{client: fakeClient}})
		if err != nil {
			t.Fatalf("サービス初期化に失敗: %v", err)
		}

		payload := map[string]string{"topic": "travel", "lang": "ja"}
		params := domain.OverwritePayloadParams{
			DBHost:         "localhost",
			DBPort:         6334,
			CollectionName: "demo",
			Payload:        payload,
			PointIDs:       []string{"point-1"},
			Filters:        domain.OverwriteFilters{},
		}

		if _, err := svc.OverwritePayload(context.Background(), params); err != nil {
			t.Fatalf("OverwritePayload が失敗: %v", err)
		}

		req := fakeClient.overwriteReq
		if req == nil {
			t.Fatalf("OverwritePayload リクエストが送られていません")
		}
		if got := req.GetCollectionName(); got != "demo" {
			t.Fatalf("collection name mismatch: %s", got)
		}
		points := req.GetPointsSelector().GetPoints()
		if points == nil || len(points.GetIds()) != 1 {
			t.Fatalf("ポイント ID が設定されていません")
		}
		if points.GetIds()[0].GetUuid() != "point-1" {
			t.Fatalf("ポイント ID の変換に失敗しています: %s", points.GetIds()[0].GetUuid())
		}
		if req.GetPayload()["topic"].GetStringValue() != "travel" {
			t.Fatalf("payload topic が一致しません")
		}
		if !fakeClient.closed {
			t.Fatalf("クライアントが Close されていません")
		}
	})

	t.Run("FilterSelectorFallback", func(t *testing.T) {
		fakeClient := &fakeQdrantClient{}
		svc, err := NewService(ServiceOptions{ClientFactory: &fakeClientFactory{client: fakeClient}})
		if err != nil {
			t.Fatalf("サービス初期化に失敗: %v", err)
		}

		params := domain.OverwritePayloadParams{
			DBHost:         "localhost",
			DBPort:         6334,
			CollectionName: "demo",
			Payload:        map[string]string{" status ": "active"},
			Filters: domain.OverwriteFilters{
				Must: []domain.PayloadFilter{{Key: "topic", Value: "travel"}},
			},
		}

		if _, err := svc.OverwritePayload(context.Background(), params); err != nil {
			t.Fatalf("OverwritePayload が失敗: %v", err)
		}

		req := fakeClient.overwriteReq
		if req.GetPointsSelector().GetFilter() == nil {
			t.Fatalf("フィルタセレクタが設定されていません")
		}
		if len(req.GetPointsSelector().GetFilter().GetMust()) != 1 {
			t.Fatalf("フィルタ条件の数が不正です")
		}
		if _, exists := req.GetPayload()["status"]; !exists {
			t.Fatalf("payload のキー整形に失敗しています")
		}
	})

	t.Run("ShouldConditions", func(t *testing.T) {
		fakeClient := &fakeQdrantClient{}
		svc, err := NewService(ServiceOptions{ClientFactory: &fakeClientFactory{client: fakeClient}})
		if err != nil {
			t.Fatalf("サービス初期化に失敗: %v", err)
		}

		params := domain.OverwritePayloadParams{
			DBHost:         "localhost",
			DBPort:         6334,
			CollectionName: "demo",
			Payload:        map[string]string{"status": "active"},
			Filters: domain.OverwriteFilters{
				Should:         []domain.PayloadFilter{{Key: "topic", Value: "travel"}},
				MinShouldCount: 1,
			},
		}

		if _, err := svc.OverwritePayload(context.Background(), params); err != nil {
			t.Fatalf("OverwritePayload が失敗: %v", err)
		}

		selector := fakeClient.overwriteReq.GetPointsSelector()
		filter := selector.GetFilter()
		if filter == nil {
			t.Fatalf("フィルタが設定されていません")
		}
		if len(filter.GetShould()) != 1 {
			t.Fatalf("should 条件が設定されていません")
		}
		if filter.GetMinShould() == nil || filter.GetMinShould().GetMinCount() != 1 {
			t.Fatalf("min_should が正しく設定されていません")
		}
	})

	t.Run("MissingPayload", func(t *testing.T) {
		fakeClient := &fakeQdrantClient{}
		svc, err := NewService(ServiceOptions{ClientFactory: &fakeClientFactory{client: fakeClient}})
		if err != nil {
			t.Fatalf("サービス初期化に失敗: %v", err)
		}
		_, err = svc.OverwritePayload(context.Background(), domain.OverwritePayloadParams{
			DBHost:         "localhost",
			DBPort:         6334,
			CollectionName: "demo",
		})
		if err == nil {
			t.Fatalf("payload 未指定エラーを検知できていません")
		}
	})
}

// --- テスト用フェイク ---

type fakeClientFactory struct {
	client *fakeQdrantClient
}

func (f *fakeClientFactory) NewClient(opts ClientOptions) (QdrantClient, error) {
	f.client.opts = opts
	f.client.closed = false
	return f.client, nil
}

type fakeQdrantClient struct {
	opts                   ClientOptions
	closed                 bool
	upsertReq              *qdrant.UpsertPoints
	queryReq               *qdrant.QueryPoints
	createReq              *qdrant.CreateCollection
	upsertResult           *qdrant.UpdateResult
	queryResult            []*qdrant.ScoredPoint
	collectionInfo         *qdrant.CollectionInfo
	describeCollectionName string
	deleteCollectionName   string
	overwriteReq           *qdrant.SetPayloadPoints
	overwriteResult        *qdrant.UpdateResult
}

func (f *fakeQdrantClient) Close() error {
	f.closed = true
	return nil
}

func (f *fakeQdrantClient) CreateCollection(ctx context.Context, request *qdrant.CreateCollection) error {
	f.createReq = request
	return nil
}

func (f *fakeQdrantClient) ListCollections(ctx context.Context) ([]string, error) {
	return []string{"demo"}, nil
}

func (f *fakeQdrantClient) GetCollectionInfo(ctx context.Context, collectionName string) (*qdrant.CollectionInfo, error) {
	f.describeCollectionName = collectionName
	if f.collectionInfo != nil {
		return f.collectionInfo, nil
	}
	return &qdrant.CollectionInfo{}, nil
}

func (f *fakeQdrantClient) DeleteCollection(ctx context.Context, collectionName string) error {
	f.deleteCollectionName = collectionName
	return nil
}

func (f *fakeQdrantClient) Upsert(ctx context.Context, request *qdrant.UpsertPoints) (*qdrant.UpdateResult, error) {
	f.upsertReq = request
	if f.upsertResult != nil {
		return f.upsertResult, nil
	}
	return &qdrant.UpdateResult{}, nil
}

func (f *fakeQdrantClient) Query(ctx context.Context, request *qdrant.QueryPoints) ([]*qdrant.ScoredPoint, error) {
	f.queryReq = request
	if f.queryResult != nil {
		return f.queryResult, nil
	}
	return nil, nil
}

func (f *fakeQdrantClient) OverwritePayload(ctx context.Context, request *qdrant.SetPayloadPoints) (*qdrant.UpdateResult, error) {
	f.overwriteReq = request
	if f.overwriteResult != nil {
		return f.overwriteResult, nil
	}
	return &qdrant.UpdateResult{}, nil
}

type fakeEmbeddingFactory struct {
	embeddings [][]float64
	lastConfig *veconfig.Config
}

func (f *fakeEmbeddingFactory) New(opts EmbeddingServiceOptions) (EmbeddingService, error) {
	return &fakeEmbeddingService{factory: f}, nil
}

type fakeEmbeddingService struct {
	factory *fakeEmbeddingFactory
}

func (f *fakeEmbeddingService) Embed(ctx context.Context, cfg *veconfig.Config) (*vedomain.EmbedResult, error) {
	// copy config for inspection
	copied := *cfg
	copied.Inputs = append([]string{}, cfg.Inputs...)
	f.factory.lastConfig = &copied
	return &vedomain.EmbedResult{Embeddings: cloneEmbeddings(f.factory.embeddings)}, nil
}

func cloneEmbeddings(in [][]float64) [][]float64 {
	out := make([][]float64, len(in))
	for i, vec := range in {
		copied := make([]float64, len(vec))
		copy(copied, vec)
		out[i] = copied
	}
	return out
}
