package renameattachments

import (
	"testing"
	"time"

	"github.com/landmaster135/devbox/internal/web_clipper/infrastructures/filesystem"
)

func TestServiceExecute_Normal(t *testing.T) {
	t.Parallel()

	mockRepo := &filesystem.MockRepository{
		ListDirectoryFunc: func(path string) ([]filesystem.FileInfo, error) {
			return []filesystem.FileInfo{
				{Name: "b.jpg", Path: "/tmp/attachments/b.jpg", ModTime: time.Date(2026, 7, 5, 10, 0, 0, 0, time.UTC)},
				{Name: "a.png", Path: "/tmp/attachments/a.png", ModTime: time.Date(2026, 7, 5, 9, 0, 0, 0, time.UTC)},
			}, nil
		},
	}
	service := NewService(mockRepo)

	result, err := service.Execute(Options{
		SrcDir:     "/tmp/attachments",
		Slug:       "openai-blog",
		Start:      1,
		Digits:     2,
		SortByName: true,
	}, time.Date(2026, 7, 5, 17, 12, 34, 0, time.UTC))
	if err != nil {
		t.Fatalf("Execute() returned error: %v", err)
	}
	if result != "リネームしました: 2件" {
		t.Fatalf("result = %q, want %q", result, "リネームしました: 2件")
	}

	wantCalls := []filesystem.RenameCall{
		{OldPath: "/tmp/attachments/a.png", NewPath: "/tmp/attachments/web-summary-20260705-171234-openai-blog_01.png"},
		{OldPath: "/tmp/attachments/b.jpg", NewPath: "/tmp/attachments/web-summary-20260705-171234-openai-blog_02.jpg"},
	}
	if len(mockRepo.RenameCalls) != len(wantCalls) {
		t.Fatalf("Rename calls = %d, want %d", len(mockRepo.RenameCalls), len(wantCalls))
	}
	for i, want := range wantCalls {
		if mockRepo.RenameCalls[i] != want {
			t.Fatalf("RenameCalls[%d] = %+v, want %+v", i, mockRepo.RenameCalls[i], want)
		}
	}
}
