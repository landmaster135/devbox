package createclips

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
	testUtil "github.com/landmaster135/devbox/internal/memos/infrastructures/testutil"
	memos "github.com/landmaster135/devbox/internal/memos/usecases"
	common "github.com/landmaster135/devbox/internal/memos_utility/usecases/common"
)

type mockCreateClipService struct {
	executeFunc func(ctx context.Context, input common.CreateClipInput) (*common.CreateClipOutput, error)
}

func (m *mockCreateClipService) Execute(ctx context.Context, input common.CreateClipInput) (*common.CreateClipOutput, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, input)
	}
	return nil, nil
}

func TestService_Execute_Normal(t *testing.T) {
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

	createClipCalls := make(map[string]common.CreateClipInput)
	service := NewService(ServiceOptions{
		FileSystem: infrastructures.NewOSFileSystem(),
		CreateClipService: &mockCreateClipService{
			executeFunc: func(ctx context.Context, input common.CreateClipInput) (*common.CreateClipOutput, error) {
				createClipCalls[filepath.Base(input.ContentFile)] = input
				return &common.CreateClipOutput{
					Operation: input.Operation,
					Memo:      &memos.Memo{Name: "memos/" + strings.TrimSuffix(filepath.Base(input.ContentFile), filepath.Ext(input.ContentFile))},
				}, nil
			},
		},
	})

	result, err := service.Execute(context.Background(), common.CreateClipsInput{
		Operation:     common.OperationCreateClips,
		ContentDir:    contentDir,
		AttachmentDir: attachmentDir,
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result.Total != 2 {
		t.Fatalf("result.Total = %d, want 2", result.Total)
	}
	if len(result.Clips) != 2 {
		t.Fatalf("len(result.Clips) = %d, want 2", len(result.Clips))
	}

	webInput, ok := createClipCalls[filepath.Base(webContent)]
	if !ok {
		t.Fatalf("create clip call missing for %s", filepath.Base(webContent))
	}
	if webInput.Operation != common.OperationCreateWebClip {
		t.Fatalf("web operation = %s, want %s", webInput.Operation, common.OperationCreateWebClip)
	}
	if len(webInput.Attachments) != 2 || webInput.Attachments[0] != webAttachment1 || webInput.Attachments[1] != webAttachment2 {
		t.Fatalf("web attachments = %#v, want [%s %s]", webInput.Attachments, webAttachment1, webAttachment2)
	}

	movieInput, ok := createClipCalls[filepath.Base(movieContent)]
	if !ok {
		t.Fatalf("create clip call missing for %s", filepath.Base(movieContent))
	}
	if movieInput.Operation != common.OperationCreateMovieClip {
		t.Fatalf("movie operation = %s, want %s", movieInput.Operation, common.OperationCreateMovieClip)
	}
	if len(movieInput.Attachments) != 1 || movieInput.Attachments[0] != movieAttachment {
		t.Fatalf("movie attachments = %#v, want [%s]", movieInput.Attachments, movieAttachment)
	}
}

func TestService_ExecuteProgressReporter_Normal(t *testing.T) {
	contentDir := t.TempDir()
	webContent := filepath.Join(contentDir, "web-summary-20241225-233435-daikokuyu-event-info.md")
	movieContent := filepath.Join(contentDir, "movie-summary-20260319-055716-trump-masako-diplomacy.md")
	for _, path := range []string{webContent, movieContent} {
		if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
			t.Fatalf("WriteFile(%s) error = %v", path, err)
		}
	}

	progresses := make([]common.CreateClipsProgress, 0, 2)
	service := NewService(ServiceOptions{
		FileSystem: infrastructures.NewOSFileSystem(),
		CreateClipService: &mockCreateClipService{
			executeFunc: func(ctx context.Context, input common.CreateClipInput) (*common.CreateClipOutput, error) {
				return &common.CreateClipOutput{Operation: input.Operation}, nil
			},
		},
		ProgressReporter: func(progress common.CreateClipsProgress) {
			progresses = append(progresses, progress)
		},
	})

	_, err := service.Execute(context.Background(), common.CreateClipsInput{
		Operation:  common.OperationCreateClips,
		ContentDir: contentDir,
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if len(progresses) != 2 {
		t.Fatalf("len(progresses) = %d, want 2", len(progresses))
	}
	if progresses[0].Current != 1 || progresses[0].Total != 2 {
		t.Fatalf("progresses[0] = %#v, want current=1 total=2", progresses[0])
	}
	if progresses[1].Current != 2 || progresses[1].Total != 2 {
		t.Fatalf("progresses[1] = %#v, want current=2 total=2", progresses[1])
	}

	if filepath.Base(progresses[0].ContentFile) != filepath.Base(movieContent) {
		t.Fatalf("progresses[0].ContentFile = %s, want %s", progresses[0].ContentFile, movieContent)
	}
	if progresses[0].Operation != common.OperationCreateMovieClip {
		t.Fatalf("progresses[0].Operation = %s, want %s", progresses[0].Operation, common.OperationCreateMovieClip)
	}
	if progresses[0].AttachmentCount != 0 {
		t.Fatalf("progresses[0].AttachmentCount = %d, want 0", progresses[0].AttachmentCount)
	}
	if filepath.Base(progresses[1].ContentFile) != filepath.Base(webContent) {
		t.Fatalf("progresses[1].ContentFile = %s, want %s", progresses[1].ContentFile, webContent)
	}
	if progresses[1].Operation != common.OperationCreateWebClip {
		t.Fatalf("progresses[1].Operation = %s, want %s", progresses[1].Operation, common.OperationCreateWebClip)
	}
	if progresses[1].AttachmentCount != 0 {
		t.Fatalf("progresses[1].AttachmentCount = %d, want 0", progresses[1].AttachmentCount)
	}
}

func TestService_ExecuteContentFilenameInvalid_NoCreateClipCall_Error(t *testing.T) {
	contentDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(contentDir, "invalid.md"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	service := NewService(ServiceOptions{
		FileSystem: infrastructures.NewOSFileSystem(),
		CreateClipService: &mockCreateClipService{
			executeFunc: func(ctx context.Context, input common.CreateClipInput) (*common.CreateClipOutput, error) {
				t.Fatal("Execute should not be called when content filename is invalid")
				return nil, nil
			},
		},
	})

	_, err := service.Execute(context.Background(), common.CreateClipsInput{
		Operation:  common.OperationCreateClips,
		ContentDir: contentDir,
	})
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "content-dir 内のファイル名が不正です") {
		t.Fatalf("error = %v, want content-dir validation message", err)
	}
}

func TestService_ExecuteAttachmentFilenameInvalid_NoCreateClipCall_Error(t *testing.T) {
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
		CreateClipService: &mockCreateClipService{
			executeFunc: func(ctx context.Context, input common.CreateClipInput) (*common.CreateClipOutput, error) {
				t.Fatal("Execute should not be called when attachment filename is invalid")
				return nil, nil
			},
		},
	})

	_, err := service.Execute(context.Background(), common.CreateClipsInput{
		Operation:     common.OperationCreateClips,
		ContentDir:    contentDir,
		AttachmentDir: attachmentDir,
	})
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "attachment-dir 内のファイル名が不正です") {
		t.Fatalf("error = %v, want attachment-dir validation message", err)
	}
}

func TestService_ExecuteAttachmentPrecheckFailed_NoCreateClipCall_Error(t *testing.T) {
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
		FileSystem: &testUtil.MockFileSystem{
			ReadAttachmentFileFunc: func(filePath string) (*infrastructures.AttachmentFile, error) {
				return nil, errors.New("mock read error")
			},
		},
		CreateClipService: &mockCreateClipService{
			executeFunc: func(ctx context.Context, input common.CreateClipInput) (*common.CreateClipOutput, error) {
				t.Fatal("Execute should not be called when attachment precheck fails")
				return nil, nil
			},
		},
	})

	_, err := service.Execute(context.Background(), common.CreateClipsInput{
		Operation:     common.OperationCreateClips,
		ContentDir:    contentDir,
		AttachmentDir: attachmentDir,
	})
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "メモは作成されませんでした") {
		t.Fatalf("error = %v, want no-memo-created message", err)
	}
}
