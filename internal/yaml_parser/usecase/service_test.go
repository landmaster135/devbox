package usecase

import (
	"errors"
	"testing"
)

type stubFileAccessor struct {
	data          []byte
	readErr       error
	writeErr      error
	lastReadPath  string
	lastWritePath string
	written       []byte
}

func (s *stubFileAccessor) ReadFile(path string) ([]byte, error) {
	s.lastReadPath = path
	if s.readErr != nil {
		return nil, s.readErr
	}
	return s.data, nil
}

func (s *stubFileAccessor) WriteFile(path string, data []byte) error {
	s.lastWritePath = path
	s.written = append([]byte(nil), data...)
	return s.writeErr
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

func TestService_ReadFromFile_UsesInjectedAccessor(t *testing.T) {
	accessor := &stubFileAccessor{data: []byte("key: value\n")}
	svc := NewServiceWithFileAccessor(accessor)

	result, err := svc.ReadFromFile("./testdata/config.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if accessor.lastReadPath != "./testdata/config.yaml" {
		t.Fatalf("expected reader to be called with './testdata/config.yaml', got %s", accessor.lastReadPath)
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
	accessor := &stubFileAccessor{readErr: expectedErr}
	svc := NewServiceWithFileAccessor(accessor)

	if _, err := svc.ReadFromFile("config.yaml"); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestService_EditFile_UpdatesNestedValues(t *testing.T) {
	initial := "server:\n  port: 8080\ninfo:\n  env: dev\n"
	accessor := &stubFileAccessor{data: []byte(initial)}
	svc := NewServiceWithFileAccessor(accessor)

	updates := "server.port=9090\ninfo.region=\"asia\"\ndebug=true"

	result, err := svc.EditFile("app.yaml", updates)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	root, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("result is not a map: %#v", result)
	}

	server := root["server"].(map[string]any)
	port, ok := server["port"].(int)
	if !ok || port != 9090 {
		t.Fatalf("expected updated port, got %v", server["port"])
	}

	info := root["info"].(map[string]any)
	region, ok := info["region"].(string)
	if !ok || region != "asia" {
		t.Fatalf("expected region to be asia, got %v", info["region"])
	}

	debug, ok := root["debug"].(bool)
	if !ok || !debug {
		t.Fatalf("expected debug flag to be true")
	}

	if accessor.lastWritePath != "app.yaml" {
		t.Fatalf("expected write call with app.yaml, got %s", accessor.lastWritePath)
	}

	if len(accessor.written) == 0 {
		t.Fatalf("expected file to be written")
	}
}

func TestService_EditFile_ReturnsErrorWhenWriteFails(t *testing.T) {
	accessor := &stubFileAccessor{
		data:     []byte("key: value\n"),
		writeErr: errors.New("write error"),
	}
	svc := NewServiceWithFileAccessor(accessor)

	if _, err := svc.EditFile("config.yaml", "key=updated"); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestService_EditFile_ReturnsErrorWhenNestedConflict(t *testing.T) {
	accessor := &stubFileAccessor{data: []byte("server: 1\n")}
	svc := NewServiceWithFileAccessor(accessor)

	if _, err := svc.EditFile("config.yaml", "server.port=9000"); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestService_EditFile_ReturnsErrorWhenKeyValueListInvalid(t *testing.T) {
	accessor := &stubFileAccessor{data: []byte("key: value\n")}
	svc := NewServiceWithFileAccessor(accessor)

	if _, err := svc.EditFile("config.yaml", "invalid-entry"); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestService_EditFile_SupportsArrayIndexUpdates(t *testing.T) {
	initial := "servers:\n  - port: 8000\n"
	accessor := &stubFileAccessor{data: []byte(initial)}
	svc := NewServiceWithFileAccessor(accessor)

	updates := "servers.0.port=9000\nservers.1.port=9100"

	result, err := svc.EditFile("servers.yaml", updates)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	root, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("result is not a map: %#v", result)
	}

	servers, ok := root["servers"].([]any)
	if !ok {
		t.Fatalf("servers is not an array: %#v", root["servers"])
	}
	if len(servers) != 2 {
		t.Fatalf("expected 2 servers, got %d", len(servers))
	}

	first := servers[0].(map[string]any)
	if first["port"].(int) != 9000 {
		t.Fatalf("expected first port 9000, got %v", first["port"])
	}

	second := servers[1].(map[string]any)
	if second["port"].(int) != 9100 {
		t.Fatalf("expected second port 9100, got %v", second["port"])
	}
}
