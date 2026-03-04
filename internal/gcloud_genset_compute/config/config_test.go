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
	args         []string
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

func (m *mockFlagParser) Parse() error   { return m.parseErr }
func (m *mockFlagParser) Args() []string { return m.args }

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

func (m *mockFlagParser) setArgs(args []string) {
	m.args = args
}

func TestParseFlagsWithParser_CreateInstanceDefaults_Normal(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationCreateGCEInstance)
	parser.setString("instance-name", " vm-1 ")

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.InstanceName != "vm-1" {
		t.Fatalf("instance-name mismatch: %s", cfg.InstanceName)
	}
	if cfg.Zone != defaultZone {
		t.Fatalf("zone mismatch: %s", cfg.Zone)
	}
	if cfg.MachineType != defaultMachineType {
		t.Fatalf("machine-type mismatch: %s", cfg.MachineType)
	}
	if cfg.BootDiskSize != defaultBootDiskSize {
		t.Fatalf("boot-disk-size mismatch: %s", cfg.BootDiskSize)
	}
	if cfg.BootDiskType != defaultBootDiskType {
		t.Fatalf("boot-disk-type mismatch: %s", cfg.BootDiskType)
	}
}

func TestParseFlagsWithParser_CreateRouterAndNATDefaults_Normal(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationCreateGCERouterAndNAT)
	parser.setString("router-name", " router-1 ")

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.RouterName != "router-1" {
		t.Fatalf("router-name mismatch: %s", cfg.RouterName)
	}
	if cfg.Region != defaultRegion {
		t.Fatalf("region mismatch: %s", cfg.Region)
	}
	if cfg.Network != defaultNetwork {
		t.Fatalf("network mismatch: %s", cfg.Network)
	}
	if cfg.NATName != defaultNATName {
		t.Fatalf("nat-name mismatch: %s", cfg.NATName)
	}
}

func TestParseFlagsWithParser_ListInstancesDefaultFormat_Normal(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationListGCloudInstances)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Format != defaultInstanceListFormat {
		t.Fatalf("format mismatch: %s", cfg.Format)
	}
}

func TestParseFlagsWithParser_IAPFirewallDefaults_Normal(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationCreateGCEIAPSSHFirewallRule)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.RuleName != defaultIAPRuleName {
		t.Fatalf("rule-name mismatch: %s", cfg.RuleName)
	}
	if cfg.Direction != defaultDirection {
		t.Fatalf("direction mismatch: %s", cfg.Direction)
	}
	if cfg.Action != defaultAction {
		t.Fatalf("action mismatch: %s", cfg.Action)
	}
	if cfg.Rules != defaultRules {
		t.Fatalf("rules mismatch: %s", cfg.Rules)
	}
	if cfg.SourceRanges != defaultIAPSourceRanges {
		t.Fatalf("source-ranges mismatch: %s", cfg.SourceRanges)
	}
}

func TestParseFlagsWithParser_IngressFirewallDefaults_Normal(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationCreateGCEIngressSSHFirewallRule)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.RuleName != defaultIngressRuleName {
		t.Fatalf("rule-name mismatch: %s", cfg.RuleName)
	}
	if cfg.AllowRule != defaultAllowRule {
		t.Fatalf("allow-rule mismatch: %s", cfg.AllowRule)
	}
	if cfg.SourceRanges != defaultIngressSourceRanges {
		t.Fatalf("source-ranges mismatch: %s", cfg.SourceRanges)
	}
}

func TestParseFlagsWithParser_Errors(t *testing.T) {
	t.Run("missing operation", func(t *testing.T) {
		parser := newMockFlagParser()
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when operation is missing")
		}
	})

	t.Run("unsupported operation", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", "unknown")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when operation is unknown")
		}
	})

	t.Run("create instance requires instance-name", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationCreateGCEInstance)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when instance-name is missing")
		}
	})

	t.Run("create router and nat requires router-name", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationCreateGCERouterAndNAT)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when router-name is missing")
		}
	})

	t.Run("reject positional args", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationListGCloudInstances)
		parser.setArgs([]string{"positional"})
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when positional arg exists")
		}
	})

	t.Run("help bypasses validation", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setBool("help", true)
		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !cfg.Help {
			t.Fatal("help should be true")
		}
	})

	t.Run("parse error", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.parseErr = errors.New("parse failed")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected parse error")
		}
	})
}

func TestParseFlags_StandardParser(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{
		"cmd",
		"-operation=create-gce-instance",
		"-instance-name=my-vm",
	}
	defer func() { os.Args = originalArgs }()

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Operation != OperationCreateGCEInstance {
		t.Fatalf("operation mismatch: %s", cfg.Operation)
	}
	if cfg.InstanceName != "my-vm" {
		t.Fatalf("instance-name mismatch: %s", cfg.InstanceName)
	}
}

func TestStandardFlagParser(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"cmd", "-string=value", "-bool", "extra"}
	defer func() { os.Args = originalArgs }()

	parser := NewStandardFlagParser()
	var str string
	var b bool

	parser.StringVar(&str, "string", "default", "")
	parser.BoolVar(&b, "bool", false, "")

	if err := parser.Parse(); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if str != "value" || !b {
		t.Fatalf("unexpected parsed values: str=%s bool=%v", str, b)
	}
	args := parser.Args()
	if len(args) != 1 || args[0] != "extra" {
		t.Fatalf("unexpected remaining args: %v", args)
	}
}

func TestPrintUsage(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"gcloud-genset-compute"}
	defer func() { os.Args = originalArgs }()

	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stderr = w

	PrintUsage()

	_ = w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to read usage output: %v", err)
	}

	output := buf.String()
	keywords := []string{
		"Google Compute Engine",
		"create-gce-instance",
		"create-gce-router-and-nat",
		"list-gcloud-instances",
	}
	for _, keyword := range keywords {
		if !strings.Contains(output, keyword) {
			t.Fatalf("usage output missing %q", keyword)
		}
	}
}
