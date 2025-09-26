package config

import (
	"strings"
	"testing"
)

type mockFlagParser struct {
	stringValues map[string]string
	boolValues   map[string]bool
	parseErr     error
}

func newMockFlagParser() *mockFlagParser {
	return &mockFlagParser{
		stringValues: make(map[string]string),
		boolValues:   make(map[string]bool),
	}
}

func (m *mockFlagParser) StringVar(p *string, name string, value string, usage string) {
	if v, ok := m.stringValues[name]; ok {
		*p = v
	} else {
		*p = value
	}
}

func (m *mockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if v, ok := m.boolValues[name]; ok {
		*p = v
	} else {
		*p = value
	}
}

func (m *mockFlagParser) Parse() error {
	return m.parseErr
}

func (m *mockFlagParser) Args() []string { return nil }

func (m *mockFlagParser) setString(name, value string)    { m.stringValues[name] = value }
func (m *mockFlagParser) setBool(name string, value bool) { m.boolValues[name] = value }

func TestParseFlagsWithParser_Success(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationUploadFiles)
	parser.setString("local-path", "/tmp/data")
	parser.setString("bucket-url", "gs://bucket")

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Operation != OperationUploadFiles {
		t.Fatalf("unexpected operation: %s", cfg.Operation)
	}
	if cfg.LocalPath != "/tmp/data" {
		t.Fatalf("unexpected local-path: %s", cfg.LocalPath)
	}
	if cfg.BucketURL != "gs://bucket" {
		t.Fatalf("unexpected bucket-url: %s", cfg.BucketURL)
	}
}

func TestParseFlagsWithParser_DownloadRequiresSources(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationDownloadFiles)
	parser.setString("destination", "./out")

	_, err := ParseFlagsWithParser(parser)
	if err == nil || !strings.Contains(err.Error(), "sources") {
		t.Fatalf("expected sources error, got: %v", err)
	}
}

func TestParseFlagsWithParser_TargetRequiresGCS(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationDeleteObject)
	parser.setString("target", "./local")

	_, err := ParseFlagsWithParser(parser)
	if err == nil || !strings.Contains(err.Error(), "gs://") {
		t.Fatalf("expected gs:// error, got: %v", err)
	}
}

func TestParseFlagsWithParser_SourcesParsing(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationDownloadFiles)
	parser.setString("sources", "gs://a, gs://b")
	parser.setString("destination", "./data")

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Sources) != 2 {
		t.Fatalf("expected 2 sources, got %d", len(cfg.Sources))
	}
	if cfg.Sources[0] != "gs://a" || cfg.Sources[1] != "gs://b" {
		t.Fatalf("unexpected sources: %v", cfg.Sources)
	}
}
