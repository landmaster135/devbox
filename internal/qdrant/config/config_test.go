package config

import (
	"os"
	"testing"
)

func TestParseFlags_CreateCollection(t *testing.T) {
	cfg := parseWithArgs(t,
		"-operation", OperationCreateCollection,
		"-collection-name", "docs",
		"-size", "1024",
	)

	if cfg.Operation != OperationCreateCollection {
		t.Fatalf("operation mismatch: %s", cfg.Operation)
	}
	if cfg.CollectionName != "docs" {
		t.Fatalf("collection name mismatch: %s", cfg.CollectionName)
	}
	if cfg.Size != 1024 {
		t.Fatalf("size mismatch: %d", cfg.Size)
	}
}

func TestParseFlags_UpsertTexts(t *testing.T) {
	cfg := parseWithArgs(t,
		"-operation", OperationUpsertTexts,
		"-collection-name", "docs",
		"-embedding-model", "nomic",
		"-input", "hello",
		"-payload", "topic=travel",
	)

	if cfg.Input != "hello" {
		t.Fatalf("input mismatch: %s", cfg.Input)
	}
	if cfg.Payload != "topic=travel" {
		t.Fatalf("payload mismatch: %s", cfg.Payload)
	}
	if cfg.EmbeddingModel != "nomic" {
		t.Fatalf("embedding model mismatch: %s", cfg.EmbeddingModel)
	}
}

func TestParseFlags_QueryPointsPayloadFormatError(t *testing.T) {
	_, err := parseWithArgsExpectError(
		"-operation", OperationQueryPoints,
		"-collection-name", "docs",
		"-embedding-model", "nomic",
		"-input", "q",
		"-payload", "invalid",
	)
	if err == nil {
		t.Fatalf("expected payload format error")
	}
	if err.Error() == "" {
		t.Fatalf("error should have message")
	}
}

func TestValidatePayloadFormat(t *testing.T) {
	if err := validatePayloadFormat("key=value"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := validatePayloadFormat("badvalue"); err == nil {
		t.Fatalf("expected error for invalid payload")
	}
}

func parseWithArgs(t *testing.T, args ...string) *Config {
	t.Helper()
	cfg, err := parseWithArgsExpectError(args...)
	if err != nil {
		t.Fatalf("ParseFlags returned error: %v", err)
	}
	return cfg
}

func parseWithArgsExpectError(args ...string) (*Config, error) {
	orig := os.Args
	defer func() { os.Args = orig }()
	os.Args = append([]string{"qdrant"}, args...)
	return ParseFlags()
}
