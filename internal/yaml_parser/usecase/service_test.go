package usecase

import (
	"errors"
	"testing"
)

type stubFileReader struct {
	data     []byte
	err      error
	lastPath string
}

func (s *stubFileReader) ReadFile(path string) ([]byte, error) {
	s.lastPath = path
	if s.err != nil {
		return nil, s.err
	}
	return s.data, nil
}

func TestService_ParseFromContent_ReturnsNormalizedMap(t *testing.T) {
	svc := NewService()

	input := "name: Hanako\nage: 29\ndetails:\n  city: Tokyo\n  zip: 100-0001\n"

	result, err := svc.ParseFromContent(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("result is not a map: %#v", result)
	}

	if data["name"] != "Hanako" {
		t.Errorf("expected name to be Hanako, got %v", data["name"])
	}

	details, ok := data["details"].(map[string]any)
	if !ok {
		t.Fatalf("details is not a map: %#v", data["details"])
	}

	if details["city"] != "Tokyo" {
		t.Errorf("expected city to be Tokyo, got %v", details["city"])
	}
}

func TestService_ReadFromFile_UsesInjectedReader(t *testing.T) {
	reader := &stubFileReader{data: []byte("key: value\n")}
	svc := NewServiceWithFileReader(reader)

	result, err := svc.ReadFromFile("./testdata/config.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if reader.lastPath != "./testdata/config.yaml" {
		t.Fatalf("expected reader to be called with './testdata/config.yaml', got %s", reader.lastPath)
	}

	data, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("result is not a map: %#v", result)
	}

	if data["key"] != "value" {
		t.Errorf("expected key to equal 'value', got %v", data["key"])
	}
}

func TestParseYAML_ReturnsSliceForMultipleDocuments(t *testing.T) {
	yamlText := "---\nname: doc1\n---\nname: doc2\n"
	result, err := parseYAML([]byte(yamlText))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	docs, ok := result.([]any)
	if !ok {
		t.Fatalf("result is not a slice: %#v", result)
	}

	if len(docs) != 2 {
		t.Fatalf("expected two documents, got %d", len(docs))
	}
}

func TestParseYAML_InvalidContent(t *testing.T) {
	if _, err := parseYAML([]byte("list: [1, 2")); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestService_ReadFromFile_ReturnsErrorWhenReaderFails(t *testing.T) {
	expectedErr := errors.New("read error")
	reader := &stubFileReader{err: expectedErr}
	svc := NewServiceWithFileReader(reader)

	if _, err := svc.ReadFromFile("config.yaml"); err == nil {
		t.Fatal("expected error but got nil")
	}
}
