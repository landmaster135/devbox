package config

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

type mockFlagParser struct {
	stringVars   map[string]*string
	boolVars     map[string]*bool
	stringValues map[string]string
	boolValues   map[string]bool
	parseErr     error
}

func newMockFlagParser() *mockFlagParser {
	return &mockFlagParser{
		stringVars:   make(map[string]*string),
		boolVars:     make(map[string]*bool),
		stringValues: make(map[string]string),
		boolValues:   make(map[string]bool),
	}
}

func (m *mockFlagParser) StringVar(p *string, name string, value string, usage string) {
	if val, ok := m.stringValues[name]; ok {
		*p = val
	} else {
		*p = value
	}
	m.stringVars[name] = p
}

func (m *mockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if val, ok := m.boolValues[name]; ok {
		*p = val
	} else {
		*p = value
	}
	m.boolVars[name] = p
}

func (m *mockFlagParser) Parse() error {
	return m.parseErr
}

func (m *mockFlagParser) Args() []string {
	return nil
}

func (m *mockFlagParser) setString(name, value string) {
	m.stringValues[name] = value
	if ptr, ok := m.stringVars[name]; ok {
		*ptr = value
	}
}

func (m *mockFlagParser) setBool(name string, value bool) {
	m.boolValues[name] = value
	if ptr, ok := m.boolVars[name]; ok {
		*ptr = value
	}
}

func TestParseFlagsWithParser_Success(t *testing.T) {
	t.Parallel()

	t.Run("delete-instance", func(t *testing.T) {
		t.Parallel()
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeleteInstance)
		parser.setString("instance-name", " my-instance ")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Operation != OperationDeleteInstance {
			t.Fatalf("operation mismatch: %s", cfg.Operation)
		}
		if cfg.InstanceName != "my-instance" {
			t.Fatalf("instance name mismatch: %s", cfg.InstanceName)
		}
	})

	t.Run("patch-deletion-protection", func(t *testing.T) {
		t.Parallel()
		parser := newMockFlagParser()
		parser.setString("operation", OperationPatchDeletionProtection)
		parser.setString("instance-name", "db-instance")
		parser.setString("deletion-protection-mode", " ENABLE ")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.DeletionProtectionMode != "enable" {
			t.Fatalf("expected normalized deletion protection mode: %s", cfg.DeletionProtectionMode)
		}
	})

	t.Run("patch-activation-policy", func(t *testing.T) {
		t.Parallel()
		parser := newMockFlagParser()
		parser.setString("operation", OperationPatchActivationPolicy)
		parser.setString("instance-name", "db-instance")
		parser.setString("activation-policy", " ALWAYS ")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.ActivationPolicy != "always" {
			t.Fatalf("expected normalized activation policy: %s", cfg.ActivationPolicy)
		}
	})

	t.Run("start-instance", func(t *testing.T) {
		t.Parallel()
		parser := newMockFlagParser()
		parser.setString("operation", OperationStartInstance)
		parser.setString("instance-name", "start-me")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Operation != OperationStartInstance {
			t.Fatalf("operation mismatch: %s", cfg.Operation)
		}
	})

	t.Run("stop-instance", func(t *testing.T) {
		t.Parallel()
		parser := newMockFlagParser()
		parser.setString("operation", OperationStopInstance)
		parser.setString("instance-name", "stop-me")

		if _, err := ParseFlagsWithParser(parser); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestParseFlagsWithParser_HelpSkipsValidation(t *testing.T) {
	parser := newMockFlagParser()
	parser.setBool("help", true)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Help {
		t.Fatalf("help flag should be true")
	}
}

func TestParseFlagsWithParser_ParseError(t *testing.T) {
	parser := newMockFlagParser()
	parser.parseErr = errors.New("parse failure")

	if _, err := ParseFlagsWithParser(parser); err == nil {
		t.Fatal("expected parse error")
	}
}

func TestParseFlags_StandardParser(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{
		"gcloud-genset-cloudsql",
		"-operation", OperationStartInstance,
		"-instance-name", "std-parser-instance",
	}

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Operation != OperationStartInstance {
		t.Fatalf("operation mismatch: %s", cfg.Operation)
	}
	if cfg.InstanceName != "std-parser-instance" {
		t.Fatalf("instance-name mismatch: %s", cfg.InstanceName)
	}
}

func TestStandardFlagParser(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"cmd", "-string-flag=val", "-bool-flag", "positional"}

	parser := NewStandardFlagParser()
	var stringFlag string
	var boolFlag bool
	parser.StringVar(&stringFlag, "string-flag", "", "usage")
	parser.BoolVar(&boolFlag, "bool-flag", false, "usage")

	if err := parser.Parse(); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if stringFlag != "val" {
		t.Fatalf("string flag mismatch: %s", stringFlag)
	}
	if !boolFlag {
		t.Fatal("expected bool flag to be true")
	}

	args := parser.Args()
	if len(args) != 1 || args[0] != "positional" {
		t.Fatalf("unexpected args: %v", args)
	}
}

func TestParseFlagsWithParser_ValidationErrors(t *testing.T) {
	t.Run("missing operation", func(t *testing.T) {
		parser := newMockFlagParser()
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when operation missing")
		}
	})

	t.Run("unsupported operation", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", "unknown")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for unsupported operation")
		}
	})

	t.Run("missing instance", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationStartInstance)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when instance-name missing")
		}
	})

	t.Run("missing deletion protection mode", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationPatchDeletionProtection)
		parser.setString("instance-name", "db")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when deletion mode missing")
		}
	})

	t.Run("invalid deletion protection mode", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationPatchDeletionProtection)
		parser.setString("instance-name", "db")
		parser.setString("deletion-protection-mode", "invalid")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for invalid deletion mode")
		}
	})

	t.Run("missing activation policy", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationPatchActivationPolicy)
		parser.setString("instance-name", "db")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when activation policy missing")
		}
	})

	t.Run("invalid activation policy", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationPatchActivationPolicy)
		parser.setString("instance-name", "db")
		parser.setString("activation-policy", "sometimes")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for invalid activation policy")
		}
	})
}

func TestRequiresInstance(t *testing.T) {
	if !requiresInstance(OperationDeleteInstance) {
		t.Fatal("expected delete-instance to require instance")
	}
	if requiresInstance("no-instance-op") {
		t.Fatal("expected custom operation to not require instance")
	}
}

func TestPrintUsage(t *testing.T) {
	output := captureStderr(func() {
		PrintUsage()
	})

	for _, expected := range []string{
		"Cloud SQL 関連の gcloud コマンド生成ツール",
		OperationDeleteInstance,
		OperationPatchDeletionProtection,
		"-instance-name",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected usage output to contain %q, but got: %s", expected, output)
		}
	}
}

func captureStderr(fn func()) string {
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	os.Stderr = w
	fn()
	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		panic(err)
	}
	return buf.String()
}
