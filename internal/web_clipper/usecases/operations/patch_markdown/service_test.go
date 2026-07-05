package patchmarkdown

import (
	"testing"

	"github.com/landmaster135/devbox/internal/web_clipper/infrastructures/filesystem"
)

func TestServiceExecute_Normal(t *testing.T) {
	t.Parallel()

	mockRepo := &filesystem.MockRepository{}
	service := NewService(mockRepo)

	result, err := service.Execute(
		"OpenAI Blog",
		"https://openai.com/blog",
		"## 記事タイトル 要約\r\n\r\n### 見出し1\r\n本文\r\n",
		"",
		"/tmp/out.md",
		2,
	)
	if err != nil {
		t.Fatalf("Execute() returned error: %v", err)
	}
	if result != "出力しました: /tmp/out.md" {
		t.Fatalf("result = %q, want %q", result, "出力しました: /tmp/out.md")
	}
	if len(mockRepo.WriteFileCalls) != 1 {
		t.Fatalf("WriteFile calls = %d, want 1", len(mockRepo.WriteFileCalls))
	}

	want := "## 記事タイトル 要約\n- [OpenAI Blog](https://openai.com/blog)\n\n### 見出し1\n本文\n"
	if string(mockRepo.WriteFileCalls[0].Data) != want {
		t.Fatalf("written data = %q, want %q", string(mockRepo.WriteFileCalls[0].Data), want)
	}
}
