package usecases

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
	testutil "github.com/landmaster135/devbox/internal/memos/infrastructures/testutil"
	memos "github.com/landmaster135/devbox/internal/memos/usecases"
)

func TestService_CreateClips_Normal(t *testing.T) {
	contentDir := t.TempDir()
	attachmentDir := t.TempDir()

	webContent := filepath.Join(contentDir, "web-summary-20241225-233435-daikokuyu-event-info.md")
	movieContent := filepath.Join(contentDir, "movie-summary-20260319-055716-trump-masako-diplomacy.md")
	if err := os.WriteFile(webContent, []byte("# web"), 0o644); err != nil {
		t.Fatalf("WriteFile(web) error = %v", err)
	}
	if err := os.WriteFile(movieContent, []byte("# movie"), 0o644); err != nil {
		t.Fatalf("WriteFile(movie) error = %v", err)
	}

	webAttachment1 := filepath.Join(attachmentDir, "web-summary-20241225-233435-daikokuyu-event-info_02.webp")
	webAttachment2 := filepath.Join(attachmentDir, "web-summary-20241225-233435-daikokuyu-event-info_11.webp")
	movieAttachment := filepath.Join(attachmentDir, "movie-summary-20260319-055716-trump-masako-diplomacy_01.webp")
	for _, path := range []string{webAttachment1, webAttachment2, movieAttachment} {
		if err := os.WriteFile(path, []byte("attachment"), 0o644); err != nil {
			t.Fatalf("WriteFile(%s) error = %v", path, err)
		}
	}

	patchCalls := make(map[string][]string)
	service := NewService(ServiceOptions{
		FileSystem: infrastructures.NewOSFileSystem(),
		MemosService: &MockMemosService{
			CreateMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
				baseName := filepath.Base(contentFile)
				name := strings.TrimSuffix(baseName, filepath.Ext(baseName))
				return &memos.Memo{Name: "memos/" + name}, nil
			},
			PatchFilesFunc: func(ctx context.Context, memo string, filePaths []string, replaces bool) (*memos.SetMemoAttachmentsOutput, error) {
				patchCalls[memo] = append([]string(nil), filePaths...)
				return &memos.SetMemoAttachmentsOutput{Name: memo}, nil
			},
		},
	})

	result, err := service.CreateClips(context.Background(), CreateClipsInput{
		Operation:     operationCreateClips,
		ContentDir:    contentDir,
		AttachmentDir: attachmentDir,
	})
	if err != nil {
		t.Fatalf("CreateClips() error = %v", err)
	}
	if result.Total != 2 {
		t.Fatalf("result.Total = %d, want 2", result.Total)
	}
	if len(result.Clips) != 2 {
		t.Fatalf("len(result.Clips) = %d, want 2", len(result.Clips))
	}

	webMemo := "memos/web-summary-20241225-233435-daikokuyu-event-info"
	webPaths, ok := patchCalls[webMemo]
	if !ok {
		t.Fatalf("patch call missing for %s", webMemo)
	}
	if len(webPaths) != 2 || webPaths[0] != webAttachment1 || webPaths[1] != webAttachment2 {
		t.Fatalf("web patch files = %#v, want [%s %s]", webPaths, webAttachment1, webAttachment2)
	}

	movieMemo := "memos/movie-summary-20260319-055716-trump-masako-diplomacy"
	moviePaths, ok := patchCalls[movieMemo]
	if !ok {
		t.Fatalf("patch call missing for %s", movieMemo)
	}
	if len(moviePaths) != 1 || moviePaths[0] != movieAttachment {
		t.Fatalf("movie patch files = %#v, want [%s]", moviePaths, movieAttachment)
	}
}

func TestService_CreateClips_ContentFilenameInvalid_NoMemoCreated_Error(t *testing.T) {
	contentDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(contentDir, "invalid.md"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	service := NewService(ServiceOptions{
		FileSystem: infrastructures.NewOSFileSystem(),
		MemosService: &MockMemosService{
			CreateMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
				t.Fatal("CreateMemo should not be called when content filename is invalid")
				return nil, nil
			},
		},
	})

	_, err := service.CreateClips(context.Background(), CreateClipsInput{
		Operation:  operationCreateClips,
		ContentDir: contentDir,
	})
	if err == nil {
		t.Fatal("CreateClips() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "content-dir 内のファイル名が不正です") {
		t.Fatalf("error = %v, want content-dir validation message", err)
	}
}

func TestService_CreateClips_AttachmentFilenameInvalid_NoMemoCreated_Error(t *testing.T) {
	contentDir := t.TempDir()
	attachmentDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(contentDir, "web-summary-20241225-233435-daikokuyu-event-info.md"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile(content) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(attachmentDir, "invalid.webp"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile(attachment) error = %v", err)
	}

	service := NewService(ServiceOptions{
		FileSystem: infrastructures.NewOSFileSystem(),
		MemosService: &MockMemosService{
			CreateMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
				t.Fatal("CreateMemo should not be called when attachment filename is invalid")
				return nil, nil
			},
		},
	})

	_, err := service.CreateClips(context.Background(), CreateClipsInput{
		Operation:     operationCreateClips,
		ContentDir:    contentDir,
		AttachmentDir: attachmentDir,
	})
	if err == nil {
		t.Fatal("CreateClips() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "attachment-dir 内のファイル名が不正です") {
		t.Fatalf("error = %v, want attachment-dir validation message", err)
	}
}

func TestService_CreateClips_AttachmentPrecheckFailed_NoMemoCreated_Error(t *testing.T) {
	contentDir := t.TempDir()
	attachmentDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(contentDir, "web-summary-20241225-233435-daikokuyu-event-info.md"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile(content) error = %v", err)
	}
	attachmentPath := filepath.Join(attachmentDir, "web-summary-20241225-233435-daikokuyu-event-info_01.webp")
	if err := os.WriteFile(attachmentPath, []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile(attachment) error = %v", err)
	}

	service := NewService(ServiceOptions{
		FileSystem: &testutil.MockFileSystem{
			ReadAttachmentFileFunc: func(filePath string) (*infrastructures.AttachmentFile, error) {
				return nil, errors.New("mock read error")
			},
		},
		MemosService: &MockMemosService{
			CreateMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
				t.Fatal("CreateMemo should not be called when attachment precheck fails")
				return nil, nil
			},
		},
	})

	_, err := service.CreateClips(context.Background(), CreateClipsInput{
		Operation:     operationCreateClips,
		ContentDir:    contentDir,
		AttachmentDir: attachmentDir,
	})
	if err == nil {
		t.Fatal("CreateClips() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "メモは作成されませんでした") {
		t.Fatalf("error = %v, want no-memo-created message", err)
	}
}
