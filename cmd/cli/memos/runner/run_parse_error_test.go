package runner

import (
	"bytes"
	"strings"
	"testing"

	cfg "github.com/landmaster135/devbox/internal/memos/config"
	usecases "github.com/landmaster135/devbox/internal/memos/usecases"
)

func TestRun_ParseError_Error(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run([]string{
		"-operation=get-memo",
		"-base-url=https://memos.example.com",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		t.Fatal("factory should not be called")
		return nil
	})

	if exitCode != 1 {
		t.Fatalf("exitCode = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr.String(), "エラー") {
		t.Fatalf("stderr = %s, want エラー", stderr.String())
	}
}
