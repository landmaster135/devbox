package services

import (
	"errors"
	"testing"
)

type mockLineDeduper struct {
	removeMatchingLinesFunc func(filePath string, startPos, endPos int) (int, error)
}

func (m *mockLineDeduper) RemoveMatchingLines(filePath string, startPos, endPos int) (int, error) {
	return m.removeMatchingLinesFunc(filePath, startPos, endPos)
}

func TestCLIService_HandleRemoveMatchingLines_Normal(t *testing.T) {
	service := NewCLIServiceWithFileService(&mockLineDeduper{
		removeMatchingLinesFunc: func(filePath string, startPos, endPos int) (int, error) {
			return 3, nil
		},
	})

	result, err := service.HandleRemoveMatchingLines("test.txt", 1, 5)
	if err != nil {
		t.Fatalf("HandleRemoveMatchingLines() error = %v, want nil", err)
	}

	want := "処理完了: 3行の重複を削除しました\n"
	if result != want {
		t.Errorf("HandleRemoveMatchingLines() = %q, want %q", result, want)
	}
}

func TestCLIService_NewCLIService_Normal(t *testing.T) {
	service := NewCLIService()
	if service == nil {
		t.Fatal("NewCLIService() = nil, want not nil")
	}

	if service.fileService == nil {
		t.Fatal("NewCLIService() fileService = nil, want not nil")
	}
}

func TestCLIService_HandleRemoveMatchingLines_Error(t *testing.T) {
	wantErr := errors.New("remove error")
	service := NewCLIServiceWithFileService(&mockLineDeduper{
		removeMatchingLinesFunc: func(filePath string, startPos, endPos int) (int, error) {
			return 0, wantErr
		},
	})

	_, err := service.HandleRemoveMatchingLines("test.txt", 1, 5)
	if err == nil {
		t.Fatal("HandleRemoveMatchingLines() error = nil, want error")
	}

	if !errors.Is(err, wantErr) {
		t.Errorf("HandleRemoveMatchingLines() error = %v, want %v", err, wantErr)
	}
}
