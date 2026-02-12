package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/landmaster135/devbox/internal/qdrant/config"
	"github.com/landmaster135/devbox/internal/qdrant/domain"
	"github.com/landmaster135/devbox/internal/qdrant/usecases"
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	if cfg.Help {
		config.PrintUsage()
		return
	}

	service, err := usecases.NewService(usecases.ServiceOptions{})
	if err != nil {
		exitWithError(fmt.Errorf("サービス初期化に失敗しました: %w", err))
	}

	ctx := context.Background()

	switch cfg.Operation {
	case config.OperationCreateCollection:
		handleCreateCollection(ctx, service, cfg)
	case config.OperationListCollections:
		handleListCollections(ctx, service, cfg)
	case config.OperationDescribeCollection:
		handleDescribeCollection(ctx, service, cfg)
	case config.OperationDeleteCollection:
		handleDeleteCollection(ctx, service, cfg)
	case config.OperationUpsertTexts:
		handleUpsertTexts(ctx, service, cfg)
	case config.OperationQueryPoints:
		handleQueryPoints(ctx, service, cfg)
	case config.OperationOverwritePayload:
		handleOverwritePayload(ctx, service, cfg)
	default:
		exitWithError(fmt.Errorf("未対応の operation です: %s", cfg.Operation))
	}
}

func handleCreateCollection(ctx context.Context, service *usecases.Service, cfg *config.Config) {
	params := domain.CreateCollectionParams{
		DBHost:         cfg.DBHost,
		DBPort:         cfg.DBPort,
		CollectionName: cfg.CollectionName,
		Size:           cfg.Size,
	}
	if err := service.CreateCollection(ctx, params); err != nil {
		exitWithError(err)
	}

	resp := map[string]string{
		"collection_name": cfg.CollectionName,
		"status":          "created",
	}
	if err := printJSON(resp); err != nil {
		exitWithError(err)
	}
}

func handleListCollections(ctx context.Context, service *usecases.Service, cfg *config.Config) {
	params := domain.ListCollectionsParams{DBHost: cfg.DBHost, DBPort: cfg.DBPort}
	collections, err := service.ListCollections(ctx, params)
	if err != nil {
		exitWithError(err)
	}
	if err := printJSON(map[string][]string{"collections": collections}); err != nil {
		exitWithError(err)
	}
}

func handleDescribeCollection(ctx context.Context, service *usecases.Service, cfg *config.Config) {
	params := domain.DescribeCollectionParams{
		DBHost:         cfg.DBHost,
		DBPort:         cfg.DBPort,
		CollectionName: cfg.CollectionName,
	}
	info, err := service.DescribeCollection(ctx, params)
	if err != nil {
		exitWithError(err)
	}
	if err := printJSON(info); err != nil {
		exitWithError(err)
	}
}

func handleDeleteCollection(ctx context.Context, service *usecases.Service, cfg *config.Config) {
	params := domain.DeleteCollectionParams{
		DBHost:         cfg.DBHost,
		DBPort:         cfg.DBPort,
		CollectionName: cfg.CollectionName,
	}
	if err := service.DeleteCollection(ctx, params); err != nil {
		exitWithError(err)
	}
	if err := printJSON(map[string]string{
		"collection_name": cfg.CollectionName,
		"status":          "deleted",
	}); err != nil {
		exitWithError(err)
	}
}

func handleUpsertTexts(ctx context.Context, service *usecases.Service, cfg *config.Config) {
	payloadMap, err := parsePayloadMap(cfg.Payload)
	if err != nil {
		exitWithError(err)
	}

	params := domain.UpsertTextsParams{
		DBHost:         cfg.DBHost,
		DBPort:         cfg.DBPort,
		CollectionName: cfg.CollectionName,
		Embedding: domain.EmbeddingConfig{
			Host:  cfg.EmbeddingHost,
			Port:  cfg.EmbeddingPort,
			Model: cfg.EmbeddingModel,
		},
		Pairs: []domain.InputPayloadPair{
			{
				Input:   cfg.Input,
				Payload: payloadCopy(payloadMap),
			},
		},
	}

	result, err := service.UpsertTexts(ctx, params)
	if err != nil {
		exitWithError(err)
	}
	if err := printJSON(result); err != nil {
		exitWithError(err)
	}
}

func handleQueryPoints(ctx context.Context, service *usecases.Service, cfg *config.Config) {
	filters, err := parsePayloadFilters(cfg.Payload)
	if err != nil {
		exitWithError(err)
	}

	params := domain.QueryPointsParams{
		DBHost:         cfg.DBHost,
		DBPort:         cfg.DBPort,
		CollectionName: cfg.CollectionName,
		Embedding: domain.EmbeddingConfig{
			Host:  cfg.EmbeddingHost,
			Port:  cfg.EmbeddingPort,
			Model: cfg.EmbeddingModel,
		},
		Input:          cfg.Input,
		PayloadFilters: filters,
		Limit:          cfg.QueryLimit,
	}

	points, err := service.QueryPoints(ctx, params)
	if err != nil {
		exitWithError(err)
	}
	if err := printJSON(points); err != nil {
		exitWithError(err)
	}
}

func handleOverwritePayload(ctx context.Context, service *usecases.Service, cfg *config.Config) {
	payloadMap, err := parsePayloadMap(cfg.Payload)
	if err != nil {
		exitWithError(err)
	}

	params := domain.OverwritePayloadParams{
		DBHost:         cfg.DBHost,
		DBPort:         cfg.DBPort,
		CollectionName: cfg.CollectionName,
		Payload:        payloadCopy(payloadMap),
		Filters:        convertOverwriteFilters(cfg),
	}

	result, err := service.OverwritePayload(ctx, params)
	if err != nil {
		exitWithError(err)
	}
	if err := printJSON(result); err != nil {
		exitWithError(err)
	}
}

func parsePayloadMap(raw string) (map[string]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	key, value, err := splitKeyValue(raw)
	if err != nil {
		return nil, err
	}
	return map[string]string{key: value}, nil
}

func payloadCopy(src map[string]string) map[string]string {
	if len(src) == 0 {
		return nil
	}
	dup := make(map[string]string, len(src))
	for k, v := range src {
		dup[k] = v
	}
	return dup
}

func parsePayloadFilters(raw string) ([]domain.PayloadFilter, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	key, value, err := splitKeyValue(raw)
	if err != nil {
		return nil, err
	}
	return []domain.PayloadFilter{{Key: key, Value: value}}, nil
}

func splitKeyValue(raw string) (string, string, error) {
	parts := strings.SplitN(raw, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("payload は key=value 形式で指定してください: %s", raw)
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	if key == "" || value == "" {
		return "", "", fmt.Errorf("payload の key/value には空文字を指定できません")
	}
	return key, value, nil
}

func convertOverwriteFilters(cfg *config.Config) domain.OverwriteFilters {
	return domain.OverwriteFilters{
		Must:           convertKeyValues(cfg.FilterMust),
		MustNot:        convertKeyValues(cfg.FilterMustNot),
		Should:         convertKeyValues(cfg.FilterShould),
		MinShouldCount: cfg.FilterMinShouldCount,
	}
}

func convertKeyValues(list []config.KeyValue) []domain.PayloadFilter {
	if len(list) == 0 {
		return nil
	}
	filters := make([]domain.PayloadFilter, 0, len(list))
	for _, kv := range list {
		filters = append(filters, domain.PayloadFilter{Key: kv.Key, Value: kv.Value})
	}
	return filters
}

func printJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("結果の整形に失敗しました: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
	os.Exit(1)
}
