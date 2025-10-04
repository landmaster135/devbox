package config

import (
	"bytes"
	"os"
	"reflect"
	"testing"
	"time"
)

type stringTarget struct {
	name string
	dest *string
	def  string
}

type boolTarget struct {
	name string
	dest *bool
	def  bool
}

type int64Target struct {
	name string
	dest *int64
	def  int64
}

type durationTarget struct {
	name string
	dest *time.Duration
	def  time.Duration
}

type mockFlagParser struct {
	stringTargets   []stringTarget
	boolTargets     []boolTarget
	int64Targets    []int64Target
	durationTargets []durationTarget
	args            []string
	stringValues    map[string]string
	boolValues      map[string]bool
	int64Values     map[string]int64
	durationValues  map[string]time.Duration
	parseErr        error
}

func newMockFlagParser() *mockFlagParser {
	return &mockFlagParser{
		stringValues:   make(map[string]string),
		boolValues:     make(map[string]bool),
		int64Values:    make(map[string]int64),
		durationValues: make(map[string]time.Duration),
	}
}

func (m *mockFlagParser) StringVar(dest *string, name string, value string, usage string) {
	*dest = value
	m.stringTargets = append(m.stringTargets, stringTarget{name: name, dest: dest, def: value})
}

func (m *mockFlagParser) BoolVar(dest *bool, name string, value bool, usage string) {
	*dest = value
	m.boolTargets = append(m.boolTargets, boolTarget{name: name, dest: dest, def: value})
}

func (m *mockFlagParser) IntVar(dest *int, name string, value int, usage string) {
	// not used
}

func (m *mockFlagParser) Int64Var(dest *int64, name string, value int64, usage string) {
	*dest = value
	m.int64Targets = append(m.int64Targets, int64Target{name: name, dest: dest, def: value})
}

func (m *mockFlagParser) DurationVar(dest *time.Duration, name string, value time.Duration, usage string) {
	*dest = value
	m.durationTargets = append(m.durationTargets, durationTarget{name: name, dest: dest, def: value})
}

func (m *mockFlagParser) Parse() error {
	if m.parseErr != nil {
		return m.parseErr
	}
	for _, target := range m.stringTargets {
		if val, ok := m.stringValues[target.name]; ok {
			*target.dest = val
		}
	}
	for _, target := range m.boolTargets {
		if val, ok := m.boolValues[target.name]; ok {
			*target.dest = val
		}
	}
	for _, target := range m.int64Targets {
		if val, ok := m.int64Values[target.name]; ok {
			*target.dest = val
		}
	}
	for _, target := range m.durationTargets {
		if val, ok := m.durationValues[target.name]; ok {
			*target.dest = val
		}
	}
	return nil
}

func (m *mockFlagParser) Args() []string {
	return m.args
}

func TestParseFlagsWithParserDailySummary(t *testing.T) {
	parser := newMockFlagParser()
	parser.stringValues["operation"] = OperationDailySummary
	parser.stringValues["access-token"] = "token"
	parser.stringValues["start-date"] = "2025-10-01"
	parser.stringValues["end-date"] = "2025-10-02"
	parser.stringValues["measure-types"] = "weight,diastolic"
	parser.stringValues["output-file-path"] = "./out.json"
	parser.boolValues["include-activity"] = false
	parser.int64Values["user-id"] = 12345

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Operation != OperationDailySummary {
		t.Fatalf("unexpected operation: %s", cfg.Operation)
	}
	if cfg.AccessToken != "token" || cfg.UserID != 12345 {
		t.Fatalf("unexpected token/user: %s %d", cfg.AccessToken, cfg.UserID)
	}
	if cfg.StartDate.IsZero() || cfg.EndDate.IsZero() {
		t.Fatalf("dates not parsed")
	}
	expected := []int{1, 9}
	if !reflect.DeepEqual(cfg.MeasureTypes, expected) {
		t.Fatalf("unexpected measure types: %v", cfg.MeasureTypes)
	}
	if cfg.IncludeActivity {
		t.Fatalf("includeActivity should be false")
	}
	if cfg.OutputFilePath != "./out.json" {
		t.Fatalf("unexpected output file path: %s", cfg.OutputFilePath)
	}
}

func TestParseFlagsWithParserAllMeasureTypes(t *testing.T) {
	parser := newMockFlagParser()
	parser.stringValues["operation"] = OperationDailySummary
	parser.stringValues["access-token"] = "token"
	parser.stringValues["start-date"] = "2025-10-01"
	parser.stringValues["measure-types"] = "all"
	parser.int64Values["user-id"] = 1

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MeasureTypes != nil {
		t.Fatalf("expected nil measure types for all")
	}
}

func TestParseFlagsWithParserAuthURL(t *testing.T) {
	parser := newMockFlagParser()
	parser.stringValues["operation"] = OperationAuthURL
	parser.stringValues["client-id"] = "cid"
	parser.stringValues["redirect-uri"] = "https://example.com/callback"
	parser.stringValues["scope"] = ""

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Scope != "user.metrics,user.activity" {
		t.Fatalf("scope normalization failed: %s", cfg.Scope)
	}
}

func TestParseFlagsWithParserDailySummaryDefaults(t *testing.T) {
	parser := newMockFlagParser()
	parser.stringValues["operation"] = OperationDailySummary
	parser.stringValues["access-token"] = "token"
	parser.stringValues["start-date"] = "2025-10-05"
	parser.int64Values["user-id"] = 99

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.EndDate.IsZero() || !cfg.StartDate.Equal(cfg.EndDate) {
		t.Fatalf("end date should default to start date: %v %v", cfg.StartDate, cfg.EndDate)
	}
	if !cfg.IncludeActivity {
		t.Fatalf("IncludeActivity should default to true")
	}
	if cfg.MeasureTypes != nil {
		t.Fatalf("MeasureTypes should be nil when not specified")
	}
}

func TestParseFlagsWithParserRequestTokenValidation(t *testing.T) {
	parser := newMockFlagParser()
	parser.stringValues["operation"] = OperationRequestToken
	parser.stringValues["client-id"] = "cid"
	parser.stringValues["client-secret"] = "secret"
	parser.stringValues["redirect-uri"] = "https://example.com"

	if _, err := ParseFlagsWithParser(parser); err == nil {
		t.Fatal("expected error when authorization code missing")
	}

	parser.stringValues["authorization-code"] = "code"
	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AuthorizationCode != "code" {
		t.Fatalf("unexpected authorization code: %s", cfg.AuthorizationCode)
	}
}

func TestParseFlagsWithParserRefreshTokenValidation(t *testing.T) {
	parser := newMockFlagParser()
	parser.stringValues["operation"] = OperationRefreshToken
	parser.stringValues["client-id"] = "cid"
	parser.stringValues["client-secret"] = "secret"

	if _, err := ParseFlagsWithParser(parser); err == nil {
		t.Fatal("expected error when refresh token missing")
	}

	parser.stringValues["refresh-token"] = "refresh"
	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.RefreshToken != "refresh" {
		t.Fatalf("unexpected refresh token: %s", cfg.RefreshToken)
	}
}

func TestParseFlagsWithParserUnknownOperation(t *testing.T) {
	parser := newMockFlagParser()
	parser.stringValues["operation"] = "unknown"
	if _, err := ParseFlagsWithParser(parser); err == nil {
		t.Fatal("expected error for unknown operation")
	}
}

func TestStandardFlagParser(t *testing.T) {
	parser := NewStandardFlagParser()
	sp, ok := parser.(*standardFlagParser)
	if !ok {
		t.Fatal("expected *standardFlagParser type")
	}
	var str string
	var b bool
	var i int
	var i64 int64
	var d time.Duration
	parser.StringVar(&str, "foo", "default", "")
	parser.BoolVar(&b, "bar", false, "")
	parser.IntVar(&i, "baz", 42, "")
	parser.Int64Var(&i64, "qux", 64, "")
	parser.DurationVar(&d, "timeout", time.Second, "")
	origArgs := os.Args
	os.Args = []string{"withings-cli-test", "-baz=100", "-qux=128", "-timeout=2s"}
	defer func() { os.Args = origArgs }()
	if err := parser.Parse(); err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if str != "default" {
		t.Fatalf("unexpected string default: %s", str)
	}
	if b {
		t.Fatalf("unexpected bool default true")
	}
	if i != 100 || i64 != 128 || d != 2*time.Second {
		t.Fatalf("unexpected numeric values: %d %d %v", i, i64, d)
	}
	if sp.Args() == nil {
		t.Fatalf("Args should not be nil")
	}
}

func TestPrintUsage(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe error: %v", err)
	}
	orig := os.Stderr
	os.Stderr = w
	PrintUsage()
	w.Close()
	os.Stderr = orig

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("read error: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("Withings")) {
		t.Fatalf("unexpected usage output: %s", buf.String())
	}
}

func TestParseFlagsHelp(t *testing.T) {
	origArgs := os.Args
	os.Args = []string{"withings", "-help"}
	defer func() { os.Args = origArgs }()

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Help {
		t.Fatalf("expected Help=true when -help specified")
	}
}
