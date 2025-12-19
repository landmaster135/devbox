package config

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

type mockFlagParser struct {
	stringValues map[string]string
	boolValues   map[string]bool
	args         []string
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

func (m *mockFlagParser) Parse() error   { return m.parseErr }
func (m *mockFlagParser) Args() []string { return m.args }

func (m *mockFlagParser) setString(name, value string)    { m.stringValues[name] = value }
func (m *mockFlagParser) setBool(name string, value bool) { m.boolValues[name] = value }
func (m *mockFlagParser) setArgs(args []string)           { m.args = args }

func TestParseFlagsWithParser_UploadFiles(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationUploadFiles)
	parser.setString("local-path", "/tmp/data")
	parser.setString("bucket-url", "gs://bucket")

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Operation != OperationUploadFiles || cfg.LocalPath != "/tmp/data" || cfg.BucketURL != "gs://bucket" {
		t.Fatalf("unexpected upload config: %+v", cfg)
	}
}

func TestParseFlagsWithParser_DownloadFiles(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationDownloadFiles)
	parser.setString("sources", "gs://a, gs://b")
	parser.setString("destination", "./data")

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Sources) != 2 || cfg.Sources[0] != "gs://a" || cfg.Sources[1] != "gs://b" {
		t.Fatalf("unexpected sources: %v", cfg.Sources)
	}
	if cfg.Destination != "./data" {
		t.Fatalf("unexpected destination: %s", cfg.Destination)
	}
}

func TestParseFlagsWithParser_InvalidSources(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationDownloadFiles)
	parser.setString("sources", " , ")
	parser.setString("destination", "./data")
	if _, err := ParseFlagsWithParser(parser); err == nil {
		t.Fatal("expected error for invalid sources")
	}
}

func TestParseFlagsWithParser_OtherOperations(t *testing.T) {
	createBucket := newMockFlagParser()
	createBucket.setString("operation", OperationCreateBucket)
	createBucket.setString("bucket-url", "gs://bucket")
	createBucket.setString("storage-class", "STANDARD")
	createBucket.setString("location", "US")
	if _, err := ParseFlagsWithParser(createBucket); err != nil {
		t.Fatalf("unexpected error create bucket: %v", err)
	}

	listContents := newMockFlagParser()
	listContents.setString("operation", OperationListContents)
	listContents.setString("target", "gs://bucket/path")
	if _, err := ParseFlagsWithParser(listContents); err != nil {
		t.Fatalf("unexpected error list contents: %v", err)
	}

	setACL := newMockFlagParser()
	setACL.setString("operation", OperationSetACL)
	setACL.setString("acl-file", "acl.json")
	setACL.setString("target", "gs://bucket/file")
	if _, err := ParseFlagsWithParser(setACL); err != nil {
		t.Fatalf("unexpected error set acl: %v", err)
	}
}

func TestParseFlagsWithParser_GCSTargetOperations(t *testing.T) {
	ops := []string{OperationShowDetails, OperationListNames, OperationDeleteObject, OperationGetACL, OperationGrantReadAll, OperationRemoveReadAll}
	for _, op := range ops {
		parser := newMockFlagParser()
		parser.setString("operation", op)
		parser.setString("target", "gs://bucket/object")
		if _, err := ParseFlagsWithParser(parser); err != nil {
			t.Fatalf("unexpected error for %s: %v", op, err)
		}
	}
}

func TestParseFlagsWithParser_Errors(t *testing.T) {
	t.Run("missing operation", func(t *testing.T) {
		parser := newMockFlagParser()
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when operation is missing")
		}
	})

	t.Run("download missing destination", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDownloadFiles)
		parser.setString("sources", "gs://a")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected destination error")
		}
	})

	t.Run("download missing sources", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDownloadFiles)
		parser.setString("destination", "./out")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected sources error")
		}
	})

	t.Run("create bucket missing params", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationCreateBucket)
		parser.setString("bucket-url", "gs://bucket")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected storage class error")
		}
	})

	t.Run("target requires gs", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeleteObject)
		parser.setString("target", "./local")
		if _, err := ParseFlagsWithParser(parser); err == nil || !strings.Contains(err.Error(), "gs://") {
			t.Fatalf("expected gs:// error, got: %v", err)
		}
	})

	t.Run("set acl requires acl-file", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationSetACL)
		parser.setString("target", "gs://bucket/file")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected acl-file error")
		}
	})
}

func TestParseFlags_StandardParser(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{
		"cmd",
		"-operation=download-files",
		"-sources=gs://a,gs://b",
		"-destination=./out",
	}
	defer func() { os.Args = originalArgs }()

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Operation != OperationDownloadFiles || len(cfg.Sources) != 2 {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestStandardFlagParser(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"cmd", "-string=value", "-bool", "extra"}
	defer func() { os.Args = originalArgs }()

	parser := NewStandardFlagParser()
	var str string
	var flag bool

	parser.StringVar(&str, "string", "default", "")
	parser.BoolVar(&flag, "bool", false, "")

	if err := parser.Parse(); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if str != "value" || !flag {
		t.Fatalf("unexpected parsed values: str=%s flag=%v", str, flag)
	}
	args := parser.Args()
	if len(args) != 1 || args[0] != "extra" {
		t.Fatalf("unexpected remaining args: %v", args)
	}
}

func TestPrintUsage(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"gcloud-genset-storage"}
	defer func() { os.Args = originalArgs }()

	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stderr = w

	PrintUsage()

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to read usage output: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "Google Cloud Storage") {
		t.Fatalf("usage output missing description: %s", output)
	}
	if !strings.Contains(output, "upload-files") || !strings.Contains(output, "set-acl") {
		t.Fatalf("usage output missing operations: %s", output)
	}
}
