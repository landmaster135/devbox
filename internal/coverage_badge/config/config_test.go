package config

import (
	"bytes"
	"flag"
	"io"
	"os"
	"strings"
	"testing"
)

type TestConfig struct{}

func TestConfigNewConfigCreateBadge_Normal(t *testing.T) {
	cfg, err := NewConfig(
		OperationCreateBadge,
		"Coverage",
		"coverage.out",
		70,
		30,
		"",
		"",
		"",
		"",
		false,
		false,
	)
	if err != nil {
		t.Fatalf("NewConfig() error = %v", err)
	}

	if cfg.TargetFile != "README.md" {
		t.Fatalf("TargetFile = %q, want README.md", cfg.TargetFile)
	}
	if cfg.Operation != OperationCreateBadge {
		t.Fatalf("Operation = %q, want %q", cfg.Operation, OperationCreateBadge)
	}
}

func TestConfigNewConfigPatchBadge_DefaultTarget_Normal(t *testing.T) {
	cfg, err := NewConfig(
		OperationPatchBadge,
		"",
		"coverage.out",
		70,
		30,
		"",
		"",
		"",
		"",
		false,
		false,
	)
	if err != nil {
		t.Fatalf("NewConfig() error = %v", err)
	}

	if cfg.BadgeTitle != "Coverage" {
		t.Fatalf("BadgeTitle = %q, want Coverage", cfg.BadgeTitle)
	}
	if cfg.TargetFile != "README.md" {
		t.Fatalf("TargetFile = %q, want README.md", cfg.TargetFile)
	}
}

func TestConfigNewConfig_InvalidOperation(t *testing.T) {
	_, err := NewConfig(
		"invalid",
		"Coverage",
		"coverage.out",
		70,
		30,
		"",
		"",
		"",
		"",
		false,
		false,
	)
	if err == nil {
		t.Fatal("NewConfig() error = nil, want error")
	}
}

func TestConfigNewConfig_InvalidThresholdRange(t *testing.T) {
	_, err := NewConfig(
		OperationCreateBadge,
		"Coverage",
		"coverage.out",
		101,
		30,
		"",
		"",
		"",
		"",
		false,
		false,
	)
	if err == nil {
		t.Fatal("NewConfig() error = nil, want error")
	}
}

func TestConfigNewConfig_InvalidThresholdRelation(t *testing.T) {
	_, err := NewConfig(
		OperationCreateBadge,
		"Coverage",
		"coverage.out",
		30,
		30,
		"",
		"",
		"",
		"",
		false,
		false,
	)
	if err == nil {
		t.Fatal("NewConfig() error = nil, want error")
	}
}

func TestConfigNewConfig_InvalidForceColor(t *testing.T) {
	_, err := NewConfig(
		OperationCreateBadge,
		"Coverage",
		"coverage.out",
		70,
		30,
		"blue",
		"",
		"",
		"",
		false,
		false,
	)
	if err == nil {
		t.Fatal("NewConfig() error = nil, want error")
	}
}

func TestConfigNewConfig_HelpBypassesValidation(t *testing.T) {
	cfg, err := NewConfig(
		"",
		"",
		"",
		0,
		0,
		"",
		"",
		"",
		"",
		false,
		true,
	)
	if err != nil {
		t.Fatalf("NewConfig() error = %v", err)
	}
	if !cfg.Help {
		t.Fatal("Help = false, want true")
	}
}

type mockFlagParser struct {
	strValues  map[string]string
	intValues  map[string]int
	boolValues map[string]bool
	parseErr   error
}

func (m *mockFlagParser) StringVar(p *string, name string, value string, usage string) {
	if v, ok := m.strValues[name]; ok {
		*p = v
		return
	}
	*p = value
}

func (m *mockFlagParser) IntVar(p *int, name string, value int, usage string) {
	if v, ok := m.intValues[name]; ok {
		*p = v
		return
	}
	*p = value
}

func (m *mockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if v, ok := m.boolValues[name]; ok {
		*p = v
		return
	}
	*p = value
}

func (m *mockFlagParser) Parse() error {
	return m.parseErr
}

func (m *mockFlagParser) Args() []string {
	return []string{}
}

type mockOSArgs struct {
	args []string
}

func (m *mockOSArgs) Args() []string {
	return m.args
}

func TestConfigParserParseFlags_Normal(t *testing.T) {
	parser := &mockFlagParser{
		strValues: map[string]string{
			"operation":     OperationPatchBadge,
			"target-file":   "README.md",
			"badge-title":   "Coverage",
			"coverage-file": "coverage.out",
		},
		intValues: map[string]int{
			"green-threshold":  80,
			"yellow-threshold": 50,
		},
		boolValues: map[string]bool{
			"dry-run": true,
		},
	}
	cfgParser := NewConfigParser(parser, &mockOSArgs{args: []string{"coverage-badge"}})

	cfg, err := cfgParser.ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() error = %v", err)
	}
	if cfg.Operation != OperationPatchBadge {
		t.Fatalf("Operation = %q, want %q", cfg.Operation, OperationPatchBadge)
	}
	if !cfg.DryRun {
		t.Fatal("DryRun = false, want true")
	}
	if cfg.GreenThreshold != 80 || cfg.YellowThreshold != 50 {
		t.Fatalf("thresholds = (%d,%d), want (80,50)", cfg.GreenThreshold, cfg.YellowThreshold)
	}
}

func TestConfigParserParseFlags_DefaultThresholds_Normal(t *testing.T) {
	parser := &mockFlagParser{
		strValues: map[string]string{
			"operation":     OperationCreateBadge,
			"badge-value":   "58.6",
			"badge-title":   "Coverage",
			"coverage-file": "coverage.out",
		},
	}
	cfgParser := NewConfigParser(parser, &mockOSArgs{args: []string{"coverage-badge"}})

	cfg, err := cfgParser.ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() error = %v", err)
	}
	if cfg.GreenThreshold != 70 || cfg.YellowThreshold != 30 {
		t.Fatalf("thresholds = (%d,%d), want (70,30)", cfg.GreenThreshold, cfg.YellowThreshold)
	}
}

func TestConfigParserParseFlags_ParseError(t *testing.T) {
	parser := &mockFlagParser{
		parseErr: io.EOF,
	}
	cfgParser := NewConfigParser(parser, &mockOSArgs{args: []string{"coverage-badge"}})

	cfg, err := cfgParser.ParseFlags()
	if err == nil {
		t.Fatal("ParseFlags() error = nil, want error")
	}
	if cfg != nil {
		t.Fatal("cfg != nil, want nil")
	}
}

func TestPrintUsage_Normal(t *testing.T) {
	output := captureStdout(t, PrintUsage)
	expectedFragments := []string{
		"使用方法:",
		"-operation string",
		"create-badge",
		"patch-badge",
	}

	for _, fragment := range expectedFragments {
		if !strings.Contains(output, fragment) {
			t.Fatalf("PrintUsage() missing fragment: %q", fragment)
		}
	}
}

func TestStandardFlagParserParse_Normal(t *testing.T) {
	originalArgs := os.Args
	originalFlagSet := flag.CommandLine
	defer func() {
		os.Args = originalArgs
		flag.CommandLine = originalFlagSet
	}()

	os.Args = []string{
		"coverage-badge",
		"-operation=create-badge",
		"-green-threshold=75",
		"-dry-run=true",
	}
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)
	flag.CommandLine = flagSet

	parser := NewStandardFlagParser()
	var operation string
	var green int
	var dryRun bool

	parser.StringVar(&operation, "operation", "", "operation")
	parser.IntVar(&green, "green-threshold", 70, "green threshold")
	parser.BoolVar(&dryRun, "dry-run", false, "dry run")

	if err := parser.Parse(); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if operation != OperationCreateBadge {
		t.Fatalf("operation = %q, want %q", operation, OperationCreateBadge)
	}
	if green != 75 {
		t.Fatalf("green = %d, want 75", green)
	}
	if !dryRun {
		t.Fatal("dryRun = false, want true")
	}

	args := parser.Args()
	if len(args) != 0 {
		t.Fatalf("Args() len = %d, want 0", len(args))
	}
}

func TestParseFlagsGlobal_Normal(t *testing.T) {
	originalArgs := os.Args
	originalFlagSet := flag.CommandLine
	defer func() {
		os.Args = originalArgs
		flag.CommandLine = originalFlagSet
	}()

	os.Args = []string{
		"coverage-badge",
		"-operation=create-badge",
		"-badge-value=58.6",
	}
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)
	flag.CommandLine = flagSet

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() error = %v", err)
	}
	if cfg.Operation != OperationCreateBadge {
		t.Fatalf("Operation = %q, want %q", cfg.Operation, OperationCreateBadge)
	}
	if cfg.BadgeValue != "58.6" {
		t.Fatalf("BadgeValue = %q, want 58.6", cfg.BadgeValue)
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	originalStdout := os.Stdout

	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = writePipe

	fn()

	if err := writePipe.Close(); err != nil {
		t.Fatalf("writePipe.Close() error = %v", err)
	}
	os.Stdout = originalStdout

	var output bytes.Buffer
	if _, err := io.Copy(&output, readPipe); err != nil {
		t.Fatalf("io.Copy() error = %v", err)
	}
	if err := readPipe.Close(); err != nil {
		t.Fatalf("readPipe.Close() error = %v", err)
	}

	return output.String()
}
