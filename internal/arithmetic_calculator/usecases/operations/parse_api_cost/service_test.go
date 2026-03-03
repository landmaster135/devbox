package parseapicost

import (
	"errors"
	"testing"
)

type mockFileReader struct {
	data []byte
	err  error
}

func (m *mockFileReader) ReadFile(filename string) ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.data, nil
}

func TestService_Execute_Normal(t *testing.T) {
	service := NewServiceWithFileReader(&mockFileReader{data: []byte("API料金が100円掛かった。API料金が200円掛かった。")})

	result, err := service.Execute("file.md", "")
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if result != "抽出されたAPI料金の合計: 300円\n" {
		t.Fatalf("unexpected result: %q", result)
	}
}

func TestService_Execute_Error(t *testing.T) {
	service := NewServiceWithFileReader(&mockFileReader{err: errors.New("read error")})

	_, err := service.Execute("file.md", "")
	if err == nil {
		t.Fatal("expected error")
	}
}
