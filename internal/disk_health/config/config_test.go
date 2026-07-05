package config

import (
	"errors"
	"os"
	"testing"

	flagParser "github.com/landmaster135/devbox/internal/disk_health/infrastructures/flag_parser"
)

func TestConfig_ParseFlagsWithParser_Normal(t *testing.T) {
	parser := flagParser.NewMockFlagParser()
	parser.SetString("operation", OperationAssessSmart)
	parser.SetString("src-file", "smart.log")
	parser.SetBool("json", true)
	parser.SetBool("verbose", true)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Operation != OperationAssessSmart {
		t.Fatalf("expected %s, got %s", OperationAssessSmart, cfg.Operation)
	}
	if cfg.SrcFile != "smart.log" {
		t.Fatalf("expected smart.log, got %s", cfg.SrcFile)
	}
	if !cfg.JSON {
		t.Fatal("expected JSON true, got false")
	}
	if !cfg.Verbose {
		t.Fatal("expected Verbose true, got false")
	}
}

func TestConfig_ParseFlagsWithParser_Help(t *testing.T) {
	parser := flagParser.NewMockFlagParser()
	parser.SetBool("help", true)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !cfg.Help {
		t.Fatal("expected Help true, got false")
	}
}

func TestConfig_ParseFlagsWithParser_ParseError(t *testing.T) {
	expectedErr := errors.New("parse failed")
	parser := flagParser.NewMockFlagParser()
	parser.ParseFunc = func() error {
		return expectedErr
	}

	_, err := ParseFlagsWithParser(parser)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestConfig_Validate_MissingOperation(t *testing.T) {
	cfg := &Config{SrcFile: "smart.log"}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestConfig_Validate_UnsupportedOperation(t *testing.T) {
	cfg := &Config{Operation: "unknown", SrcFile: "smart.log"}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestConfig_Validate_MissingSrcFile(t *testing.T) {
	cfg := &Config{Operation: OperationAssessSmart}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestConfig_ParseFlags_Normal(t *testing.T) {
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()
	os.Args = []string{"disk-health", "-operation=assess-smart", "-src-file=smart.log", "-json"}

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.SrcFile != "smart.log" {
		t.Fatalf("expected smart.log, got %s", cfg.SrcFile)
	}
	if !cfg.JSON {
		t.Fatal("expected JSON true, got false")
	}
}
