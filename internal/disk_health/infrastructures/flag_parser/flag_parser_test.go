package flag_parser

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

func TestMockFlagParser_StringVar_Normal(t *testing.T) {
	parser := NewMockFlagParser()
	parser.SetString("operation", "assess-smart")

	var actual string
	parser.StringVar(&actual, "operation", "", "operation")

	if actual != "assess-smart" {
		t.Fatalf("expected assess-smart, got %s", actual)
	}
}

func TestMockFlagParser_BoolVar_Normal(t *testing.T) {
	parser := NewMockFlagParser()
	parser.SetBool("json", true)

	var actual bool
	parser.BoolVar(&actual, "json", false, "json")

	if !actual {
		t.Fatal("expected true, got false")
	}
}

func TestMockFlagParser_Parse_Error(t *testing.T) {
	expectedErr := errors.New("parse failed")
	parser := NewMockFlagParser()
	parser.ParseFunc = func() error {
		return expectedErr
	}

	err := parser.Parse()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestStandardFlagParser_Parse_Normal(t *testing.T) {
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()
	os.Args = []string{"disk-health", "-operation=assess-smart", "-json"}

	parser := NewStandardFlagParser()
	var operation string
	var outputJSON bool
	parser.StringVar(&operation, "operation", "", "operation")
	parser.BoolVar(&outputJSON, "json", false, "json")

	if err := parser.Parse(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if operation != "assess-smart" {
		t.Fatalf("expected assess-smart, got %s", operation)
	}
	if !outputJSON {
		t.Fatal("expected json true, got false")
	}
}

func TestPrintUsage_Normal(t *testing.T) {
	originalArgs := os.Args
	originalStderr := os.Stderr
	defer func() {
		os.Args = originalArgs
		os.Stderr = originalStderr
	}()

	readFile, writeFile, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Args = []string{"disk-health"}
	os.Stderr = writeFile

	PrintUsage("usage: %s\n")
	if err := writeFile.Close(); err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}

	output, readErr := io.ReadAll(readFile)
	if readErr != nil {
		t.Fatalf("failed to read usage output: %v", readErr)
	}
	if !strings.Contains(string(output), "disk-health") {
		t.Fatalf("expected usage output, got %s", string(output))
	}
}
