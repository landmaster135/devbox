package evaluatelinecount

import (
	"errors"
	"os"
	"testing"
)

type mockBufioScanner struct {
	scanResult bool
	err        error
}

func (m *mockBufioScanner) Scan() bool { return m.scanResult }
func (m *mockBufioScanner) Err() error { return m.err }

type mockJSONMarshaler struct {
	err error
}

func (m *mockJSONMarshaler) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []byte("{}"), nil
}

type mockFileOpener struct {
	err error
}

func (m *mockFileOpener) Open(name string) (*os.File, error) {
	if m.err != nil {
		return nil, m.err
	}
	return os.CreateTemp("", "line-count-*.txt")
}

func TestService_Execute_ErrorByMarshaler(t *testing.T) {
	service := NewServiceWithDependencies(&mockFileOpener{}, &mockBufioScanner{scanResult: false}, &mockJSONMarshaler{err: errors.New("json error")})

	_, err := service.Execute("/tmp/test.txt", 1)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestIsGreaterDescription_Normal(t *testing.T) {
	if IsGreaterDescription(true) != "より大きいです。" {
		t.Fatal("unexpected true description")
	}
	if IsGreaterDescription(false) != "以下です。" {
		t.Fatal("unexpected false description")
	}
}
