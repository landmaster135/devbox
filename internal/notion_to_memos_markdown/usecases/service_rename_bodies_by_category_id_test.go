package usecases

import (
	"errors"
	"testing"
)

type mockRenameBodiesByCategoryIDOperation struct {
	ExecuteFunc func(pageType string, conNumberStart, conNumberEnd int, srcJSONFile, srcResourceDir string) (string, error)
}

func (m *mockRenameBodiesByCategoryIDOperation) Execute(pageType string, conNumberStart, conNumberEnd int, srcJSONFile, srcResourceDir string) (string, error) {
	return m.ExecuteFunc(pageType, conNumberStart, conNumberEnd, srcJSONFile, srcResourceDir)
}

func TestService_RenameBodiesByCategoryID_Normal(t *testing.T) {
	t.Parallel()

	op := &mockRenameBodiesByCategoryIDOperation{
		ExecuteFunc: func(pageType string, conNumberStart, conNumberEnd int, srcJSONFile, srcResourceDir string) (string, error) {
			if pageType != "content" ||
				conNumberStart != 100 ||
				conNumberEnd != 200 ||
				srcJSONFile != "/tmp/contents.json" ||
				srcResourceDir != "/tmp/resource" {
				t.Fatalf("unexpected args: pageType=%s start=%d end=%d json=%s dir=%s", pageType, conNumberStart, conNumberEnd, srcJSONFile, srcResourceDir)
			}
			return "ok", nil
		},
	}

	service := newServiceWithOperations(nil, nil, nil, nil, op, nil)
	got, err := service.RenameBodiesByCategoryID("content", 100, 200, "/tmp/contents.json", "/tmp/resource")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ok" {
		t.Fatalf("got = %q, want ok", got)
	}
}

func TestService_RenameBodiesByCategoryID_Error(t *testing.T) {
	t.Parallel()

	op := &mockRenameBodiesByCategoryIDOperation{
		ExecuteFunc: func(pageType string, conNumberStart, conNumberEnd int, srcJSONFile, srcResourceDir string) (string, error) {
			return "", errors.New("failed")
		},
	}

	service := newServiceWithOperations(nil, nil, nil, nil, op, nil)
	_, err := service.RenameBodiesByCategoryID("content", 100, 200, "/tmp/contents.json", "/tmp/resource")
	if err == nil || err.Error() != "failed" {
		t.Fatalf("error = %v", err)
	}
}
