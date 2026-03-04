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
	intVars      map[string]*int
	boolVars     map[string]*bool
	stringValues map[string]string
	intValues    map[string]int
	boolValues   map[string]bool
	args         []string
	parseErr     error
}

func newMockFlagParser() *mockFlagParser {
	return &mockFlagParser{
		stringVars:   make(map[string]*string),
		intVars:      make(map[string]*int),
		boolVars:     make(map[string]*bool),
		stringValues: make(map[string]string),
		intValues:    make(map[string]int),
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

func (m *mockFlagParser) IntVar(p *int, name string, value int, usage string) {
	if preset, ok := m.intValues[name]; ok {
		*p = preset
	} else {
		*p = value
	}
	m.intVars[name] = p
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

func (m *mockFlagParser) setInt(name string, value int) {
	m.intValues[name] = value
	if ptr, ok := m.intVars[name]; ok {
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

func TestParseFlagsWithParser_CreateInstanceWithStartupScriptDefaults_Normal(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationCreateGCEInstanceWithStartupScript)
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
	if cfg.StartupScriptPath != defaultStartupScriptPath {
		t.Fatalf("startup-script-path mismatch: %s", cfg.StartupScriptPath)
	}
}

func TestParseFlagsWithParser_CreateInstanceAndConfigureDefaults_Normal(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationCreateGCEInstanceAndConfigure)
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
	if cfg.MetadataYAMLPath != defaultMetadataYAMLPath {
		t.Fatalf("metadata-yaml-path mismatch: %s", cfg.MetadataYAMLPath)
	}
	if cfg.StartupScriptPath != defaultStartupScriptPath {
		t.Fatalf("startup-script-path mismatch: %s", cfg.StartupScriptPath)
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

func TestParseFlagsWithParser_ListDiskTypes_Normal(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationListDiskTypes)
	parser.setString("zones", " asia-southeast3-a , asia-southeast3-b ")
	parser.setInt("min-disk-size-gib", 4)
	parser.setInt("max-disk-size-gib", 65536)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Zones) != 2 {
		t.Fatalf("zones length mismatch: %v", cfg.Zones)
	}
	if cfg.Zones[0] != "asia-southeast3-a" || cfg.Zones[1] != "asia-southeast3-b" {
		t.Fatalf("zones mismatch: %v", cfg.Zones)
	}
	if cfg.MinDiskSizeGiB != 4 {
		t.Fatalf("min-disk-size-gib mismatch: %d", cfg.MinDiskSizeGiB)
	}
	if cfg.MaxDiskSizeGiB != 65536 {
		t.Fatalf("max-disk-size-gib mismatch: %d", cfg.MaxDiskSizeGiB)
	}
}

func TestParseFlagsWithParser_ListMachineTypes_Normal(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationListMachineTypes)
	parser.setString("zones", "asia-southeast3-a")
	parser.setInt("min-disk-size-gib", 1024)
	parser.setInt("max-memory-size-mib", 65536)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Zones) != 1 || cfg.Zones[0] != "asia-southeast3-a" {
		t.Fatalf("zones mismatch: %v", cfg.Zones)
	}
	if cfg.MinDiskSizeGiB != 1024 {
		t.Fatalf("min-disk-size-gib mismatch: %d", cfg.MinDiskSizeGiB)
	}
	if cfg.MaxDiskSizeGiB != 0 {
		t.Fatalf("max-disk-size-gib mismatch: %d", cfg.MaxDiskSizeGiB)
	}
	if cfg.MinMemorySizeMiB != 0 {
		t.Fatalf("min-memory-size-mib mismatch: %d", cfg.MinMemorySizeMiB)
	}
	if cfg.MaxMemorySizeMiB != 65536 {
		t.Fatalf("max-memory-size-mib mismatch: %d", cfg.MaxMemorySizeMiB)
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

func TestParseFlagsWithParser_InstanceLifecycleOperations_Normal(t *testing.T) {
	tests := []struct {
		name      string
		operation string
	}{
		{name: "start", operation: OperationStartGCEInstance},
		{name: "stop", operation: OperationStopGCEInstance},
		{name: "reboot", operation: OperationRebootGCEInstance},
		{name: "delete", operation: OperationDeleteGCEInstance},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := newMockFlagParser()
			parser.setString("operation", tt.operation)
			parser.setString("instance-name", " vm-1 ")
			parser.setString("zone", " us-central1-a ")

			cfg, err := ParseFlagsWithParser(parser)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cfg.InstanceName != "vm-1" {
				t.Fatalf("instance-name mismatch: %s", cfg.InstanceName)
			}
			if cfg.Zone != "us-central1-a" {
				t.Fatalf("zone mismatch: %s", cfg.Zone)
			}
		})
	}
}

func TestParseFlagsWithParser_CopyGCESSHKeyDefaults_Normal(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationCopyGCESSHKey)
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
	if cfg.SSHKeyPath != defaultSSHKeyPath {
		t.Fatalf("ssh-key-path mismatch: %s", cfg.SSHKeyPath)
	}
}

func TestParseFlagsWithParser_ConnectGCEInstanceDefaults_Normal(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationConnectGCEInstance)
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
}

func TestParseFlagsWithParser_SetGCEInstanceMetadataFromYAMLDefaults_Normal(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationSetGCEInstanceMetadataFromYAML)
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
	if cfg.MetadataYAMLPath != defaultMetadataYAMLPath {
		t.Fatalf("metadata-yaml-path mismatch: %s", cfg.MetadataYAMLPath)
	}
}

func TestParseFlagsWithParser_AddStartupScriptToGCEInstanceDefaults_Normal(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationAddStartupScriptToGCEInstance)
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
	if cfg.StartupScriptPath != defaultStartupScriptPath {
		t.Fatalf("startup-script-path mismatch: %s", cfg.StartupScriptPath)
	}
}

func TestParseFlagsWithParser_SetupGCEFirewallAndSSHDefaults_Normal(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationSetupGCEFirewallAndSSH)
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
	if cfg.SSHKeyPath != defaultSSHKeyPath {
		t.Fatalf("ssh-key-path mismatch: %s", cfg.SSHKeyPath)
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

	t.Run("instance lifecycle operations require instance-name and zone", func(t *testing.T) {
		tests := []struct {
			name      string
			operation string
			zone      string
		}{
			{name: "start missing instance-name", operation: OperationStartGCEInstance, zone: "us-central1-a"},
			{name: "stop missing instance-name", operation: OperationStopGCEInstance, zone: "us-central1-a"},
			{name: "reboot missing instance-name", operation: OperationRebootGCEInstance, zone: "us-central1-a"},
			{name: "delete missing instance-name", operation: OperationDeleteGCEInstance, zone: "us-central1-a"},
			{name: "start missing zone", operation: OperationStartGCEInstance},
			{name: "stop missing zone", operation: OperationStopGCEInstance},
			{name: "reboot missing zone", operation: OperationRebootGCEInstance},
			{name: "delete missing zone", operation: OperationDeleteGCEInstance},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				parser := newMockFlagParser()
				parser.setString("operation", tt.operation)
				parser.setString("instance-name", "vm-1")
				if tt.zone != "" {
					parser.setString("zone", tt.zone)
				}

				if strings.Contains(tt.name, "missing instance-name") {
					parser.setString("instance-name", "")
				}

				if _, err := ParseFlagsWithParser(parser); err == nil {
					t.Fatal("expected validation error")
				}
			})
		}
	})

	t.Run("copy/connect/setup operations require instance-name", func(t *testing.T) {
		tests := []struct {
			name      string
			operation string
		}{
			{name: "copy missing instance-name", operation: OperationCopyGCESSHKey},
			{name: "connect missing instance-name", operation: OperationConnectGCEInstance},
			{name: "set metadata missing instance-name", operation: OperationSetGCEInstanceMetadataFromYAML},
			{name: "add startup script missing instance-name", operation: OperationAddStartupScriptToGCEInstance},
			{name: "setup missing instance-name", operation: OperationSetupGCEFirewallAndSSH},
			{name: "create with startup script missing instance-name", operation: OperationCreateGCEInstanceWithStartupScript},
			{name: "create and configure missing instance-name", operation: OperationCreateGCEInstanceAndConfigure},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				parser := newMockFlagParser()
				parser.setString("operation", tt.operation)
				if _, err := ParseFlagsWithParser(parser); err == nil {
					t.Fatal("expected validation error")
				}
			})
		}
	})

	t.Run("list operations validate size range", func(t *testing.T) {
		t.Run("list disk types min greater than max", func(t *testing.T) {
			parser := newMockFlagParser()
			parser.setString("operation", OperationListDiskTypes)
			parser.setInt("min-disk-size-gib", 200)
			parser.setInt("max-disk-size-gib", 100)
			if _, err := ParseFlagsWithParser(parser); err == nil {
				t.Fatal("expected validation error")
			}
		})

		t.Run("list machine types negative min", func(t *testing.T) {
			parser := newMockFlagParser()
			parser.setString("operation", OperationListMachineTypes)
			parser.setInt("min-disk-size-gib", -1)
			if _, err := ParseFlagsWithParser(parser); err == nil {
				t.Fatal("expected validation error")
			}
		})

		t.Run("list machine types negative min memory", func(t *testing.T) {
			parser := newMockFlagParser()
			parser.setString("operation", OperationListMachineTypes)
			parser.setInt("min-memory-size-mib", -1)
			if _, err := ParseFlagsWithParser(parser); err == nil {
				t.Fatal("expected validation error")
			}
		})

		t.Run("list machine types min memory greater than max memory", func(t *testing.T) {
			parser := newMockFlagParser()
			parser.setString("operation", OperationListMachineTypes)
			parser.setInt("min-memory-size-mib", 8192)
			parser.setInt("max-memory-size-mib", 4096)
			if _, err := ParseFlagsWithParser(parser); err == nil {
				t.Fatal("expected validation error")
			}
		})

		t.Run("list disk types invalid zone", func(t *testing.T) {
			parser := newMockFlagParser()
			parser.setString("operation", OperationListDiskTypes)
			parser.setString("zones", "asia-southeast3-a,asia_southeast3_b")
			if _, err := ParseFlagsWithParser(parser); err == nil {
				t.Fatal("expected validation error")
			}
		})
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
	os.Args = []string{"cmd", "-string=value", "-int=7", "-bool", "extra"}
	defer func() { os.Args = originalArgs }()

	parser := NewStandardFlagParser()
	var str string
	var i int
	var b bool

	parser.StringVar(&str, "string", "default", "")
	parser.IntVar(&i, "int", 0, "")
	parser.BoolVar(&b, "bool", false, "")

	if err := parser.Parse(); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if str != "value" || i != 7 || !b {
		t.Fatalf("unexpected parsed values: str=%s int=%d bool=%v", str, i, b)
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
		"create-gce-instance-with-startup-script",
		"create-gce-instance-and-configure",
		"create-gce-router-and-nat",
		"list-gcloud-instances",
		"list-disk-types",
		"list-machine-types",
		"min-memory-size-mib",
		"max-memory-size-mib",
		"start-gce-instance",
		"stop-gce-instance",
		"reboot-gce-instance",
		"delete-gce-instance",
		"copy-gce-ssh-key",
		"connect-gce-instance",
		"set-gce-instance-metadata-from-yaml",
		"add-startup-script-to-gce-instance",
		"setup-gce-firewall-and-ssh",
	}
	for _, keyword := range keywords {
		if !strings.Contains(output, keyword) {
			t.Fatalf("usage output missing %q", keyword)
		}
	}
}
