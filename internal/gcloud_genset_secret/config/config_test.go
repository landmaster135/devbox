package config

import (
	"bytes"
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
	parseError   error
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
	if preset, ok := m.stringValues[name]; ok {
		*p = preset
	} else {
		*p = value
	}
	m.stringVars[name] = p
}

func (m *mockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if preset, ok := m.boolValues[name]; ok {
		*p = preset
	} else {
		*p = value
	}
	m.boolVars[name] = p
}

func (m *mockFlagParser) Parse() error {
	return m.parseError
}

func (m *mockFlagParser) Args() []string {
	return []string{}
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
	t.Run("create secret automatic", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationCreateSecret)
		parser.setString("secret-name", "test-secret")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.ReplicationPolicy != replicationPolicyAutomatic {
			t.Fatalf("unexpected replication policy: %s", cfg.ReplicationPolicy)
		}
	})

	t.Run("create secret user-managed", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationCreateSecret)
		parser.setString("secret-name", "test-secret")
		parser.setString("replication-policy", replicationPolicyUserManaged)
		parser.setString("locations", "us-central1,us-east1")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Locations != "us-central1,us-east1" {
			t.Fatalf("unexpected locations: %s", cfg.Locations)
		}
	})

	t.Run("add secret version", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationAddSecretVersion)
		parser.setString("secret-name", "app-secret")
		parser.setString("secret-value", "value")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.SecretValue != "value" {
			t.Fatalf("unexpected secret value: %s", cfg.SecretValue)
		}
	})

	t.Run("create and add secret version", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationCreateAndAddSecretVersion)
		parser.setString("secret-name", "combo-secret")
		parser.setString("secret-value", "combo")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.ReplicationPolicy != replicationPolicyAutomatic {
			t.Fatalf("unexpected replication policy: %s", cfg.ReplicationPolicy)
		}
	})

	t.Run("access secret version default", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationAccessSecretVersion)
		parser.setString("secret-name", "read-secret")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Version != defaultVersion {
			t.Fatalf("unexpected version: %s", cfg.Version)
		}
	})

	t.Run("update secret labels", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUpdateSecretLabels)
		parser.setString("secret-name", "labeled-secret")
		parser.setString("labels", "env=prod,team=devops")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Labels != "env=prod,team=devops" {
			t.Fatalf("unexpected labels: %s", cfg.Labels)
		}
	})

	t.Run("update secret version aliases", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUpdateSecretVersionAliases)
		parser.setString("secret-name", "alias-secret")
		parser.setString("alias-option", "--update-version-aliases=prod=5")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.AliasOption != "--update-version-aliases=prod=5" {
			t.Fatalf("unexpected alias option: %s", cfg.AliasOption)
		}
	})
}

func TestParseFlagsWithParser_Errors(t *testing.T) {
	t.Run("missing operation", func(t *testing.T) {
		parser := newMockFlagParser()
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when operation is missing")
		}
	})

	t.Run("missing secret name", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationAddSecretVersion)
		parser.setString("secret-value", "value")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when secret name is missing")
		}
	})

	t.Run("missing secret value", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationAddSecretVersion)
		parser.setString("secret-name", "app-secret")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when secret value is missing")
		}
	})

	t.Run("user-managed without locations", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationCreateSecret)
		parser.setString("secret-name", "missing-locations")
		parser.setString("replication-policy", replicationPolicyUserManaged)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when locations are missing")
		}
	})

	t.Run("invalid replication policy", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationCreateSecret)
		parser.setString("secret-name", "test")
		parser.setString("replication-policy", "invalid")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for invalid replication policy")
		}
	})

	t.Run("invalid alias option", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUpdateSecretVersionAliases)
		parser.setString("secret-name", "alias-secret")
		parser.setString("alias-option", "--unknown")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for invalid alias option")
		}
	})

	t.Run("unsupported operation", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", "nope")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for unsupported operation")
		}
	})
}

func TestParseFlagsWithParser_Help(t *testing.T) {
	parser := newMockFlagParser()
	parser.setBool("help", true)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Help {
		t.Fatalf("expected help flag to remain true")
	}
}

func TestValidateAliasOption(t *testing.T) {
	cases := []struct {
		name    string
		option  string
		wantErr bool
	}{
		{"clear", "--clear-version-aliases", false},
		{"remove", "--remove-version-aliases=prod", false},
		{"update", "--update-version-aliases=prod=5", false},
		{"empty", "", true},
		{"unknown", "--nope", true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := validateAliasOption(tc.option)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error but got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestParseFlags_StandardParser(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{
		"cmd",
		"-operation=add-secret-version",
		"-secret-name=my-secret",
		"-secret-value=VALUE",
	}
	defer func() { os.Args = originalArgs }()

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Operation != OperationAddSecretVersion {
		t.Fatalf("unexpected operation: %s", cfg.Operation)
	}
	if cfg.SecretValue != "VALUE" {
		t.Fatalf("unexpected secret value: %s", cfg.SecretValue)
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

	if str != "value" {
		t.Fatalf("unexpected string value: %s", str)
	}
	if !flag {
		t.Fatalf("expected bool to be true")
	}

	args := parser.Args()
	if len(args) != 1 || args[0] != "extra" {
		t.Fatalf("unexpected remaining args: %v", args)
	}
}

func TestPrintUsage(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"gcloud-genset-secret"}
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
	if !strings.Contains(output, "Google Cloud Secret Manager") {
		t.Fatalf("usage output missing expected description: %s", output)
	}
	if !strings.Contains(output, "update-secret-version-aliases") {
		t.Fatalf("usage output missing alias section: %s", output)
	}
}
