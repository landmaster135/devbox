package filesystem

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestMockRepository_DefaultBehavior(t *testing.T) {
	mock := &MockRepository{}

	data, err := mock.ReadFile("any")
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if data != nil {
		t.Fatalf("ReadFile() = %v, want nil", data)
	}

	dirs, err := mock.ListDirectories("any")
	if err != nil {
		t.Fatalf("ListDirectories() error = %v", err)
	}
	if !reflect.DeepEqual(dirs, []string{}) {
		t.Fatalf("ListDirectories() = %v, want empty slice", dirs)
	}

	joined := mock.Join("a", "b")
	if joined != filepath.Join("a", "b") {
		t.Fatalf("Join() = %s, want %s", joined, filepath.Join("a", "b"))
	}
}

func TestMockRepository_FunctionOverride(t *testing.T) {
	readErr := errors.New("read err")
	writeErr := errors.New("write err")
	listErr := errors.New("list err")
	mock := &MockRepository{
		ReadFileFunc: func(path string) ([]byte, error) {
			return nil, readErr
		},
		WriteFileFunc: func(path string, data []byte, perm os.FileMode) error {
			return writeErr
		},
		ListDirectoriesFunc: func(path string) ([]string, error) {
			return nil, listErr
		},
		JoinFunc: func(elem ...string) string {
			return "/custom/path"
		},
	}

	if _, err := mock.ReadFile("x"); !errors.Is(err, readErr) {
		t.Fatalf("ReadFile() error = %v, want %v", err, readErr)
	}
	if err := mock.WriteFile("x", []byte("abc"), 0o644); !errors.Is(err, writeErr) {
		t.Fatalf("WriteFile() error = %v, want %v", err, writeErr)
	}
	if _, err := mock.ListDirectories("x"); !errors.Is(err, listErr) {
		t.Fatalf("ListDirectories() error = %v, want %v", err, listErr)
	}
	if got := mock.Join("a"); got != "/custom/path" {
		t.Fatalf("Join() = %s, want /custom/path", got)
	}

	if mock.LastWritePath != "x" {
		t.Fatalf("LastWritePath = %s, want x", mock.LastWritePath)
	}
	if string(mock.LastWriteContent) != "abc" {
		t.Fatalf("LastWriteContent = %q, want %q", string(mock.LastWriteContent), "abc")
	}
	if mock.LastWritePermission != 0o644 {
		t.Fatalf("LastWritePermission = %o, want %o", mock.LastWritePermission, 0o644)
	}
}
