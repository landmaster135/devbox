package config

import (
	"os"
	"testing"
)

func TestStandardFileReader_ReadFile(t *testing.T) {
	tmpFile, err := os.CreateTemp(t.TempDir(), "config-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	content := []byte("sample contents")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	reader := &StandardFileReader{}
	data, err := reader.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadFile() error = %v, want nil", err)
	}

	if string(data) != string(content) {
		t.Errorf("ReadFile() returned %q, want %q", string(data), string(content))
	}
}

func TestStandardFileReader_ReadFile_Error(t *testing.T) {
	reader := &StandardFileReader{}
	if _, err := reader.ReadFile("does-not-exist.txt"); err == nil {
		t.Fatal("ReadFile() error = nil, want error")
	}
}
