package createcommonmemos

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
	"github.com/landmaster135/devbox/internal/memos_utility/usecases/common"
)

type mockMemosService struct {
	createMemoFunc       func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error)
	patchFilesFunc       func(ctx context.Context, memo string, filePaths []string, replaces bool) (*memos.SetMemoAttachmentsOutput, error)
	addMemoRelationsFunc func(ctx context.Context, memo string, relatedMemos []string, replaces bool) (*memos.AddMemoRelationsOutput, error)
}

func (m *mockMemosService) CreateMemo(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
	if m.createMemoFunc != nil {
		return m.createMemoFunc(ctx, memoID, content, contentFile, visibility, state, pinned, displayTime)
	}
	return nil, nil
}

func (m *mockMemosService) PatchFiles(ctx context.Context, memo string, filePaths []string, replaces bool) (*memos.SetMemoAttachmentsOutput, error) {
	if m.patchFilesFunc != nil {
		return m.patchFilesFunc(ctx, memo, filePaths, replaces)
	}
	return nil, nil
}

func (m *mockMemosService) AddMemoRelations(ctx context.Context, memo string, relatedMemos []string, replaces bool) (*memos.AddMemoRelationsOutput, error) {
	if m.addMemoRelationsFunc != nil {
		return m.addMemoRelationsFunc(ctx, memo, relatedMemos, replaces)
	}
	return nil, nil
}

func TestService_Execute_Normal(t *testing.T) {
	contentDir := t.TempDir()
	attachmentDir := t.TempDir()

	content1 := filepath.Join(contentDir, "20260316080301_01.md")
	content2 := filepath.Join(contentDir, "20260316080301_02.md")
	content3 := filepath.Join(contentDir, "20260317010101_01.md")
	for _, contentFile := range []string{content1, content2, content3} {
		if err := os.WriteFile(contentFile, []byte("# content"), 0o644); err != nil {
			t.Fatalf("WriteFile(%s) error = %v", contentFile, err)
		}
	}

	attachment1 := filepath.Join(attachmentDir, "20260316080301_02_11.webp")
	attachment2 := filepath.Join(attachmentDir, "20260316080301_02_21.webp")
	for _, attachmentFile := range []string{attachment1, attachment2} {
		if err := os.WriteFile(attachmentFile, []byte("attachment"), 0o644); err != nil {
			t.Fatalf("WriteFile(%s) error = %v", attachmentFile, err)
		}
	}

	createMemoOrder := make([]string, 0, 3)
	patchMemoID := ""
	patchFiles := []string(nil)
	addRelationMemoID := ""
	addRelationRelatedMemoID := ""

	service := NewService(ServiceOptions{
		FileSystem: infrastructures.NewOSFileSystem(),
		MemosService: &mockMemosService{
			createMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
				baseName := filepath.Base(contentFile)
				createMemoOrder = append(createMemoOrder, baseName)
				if visibility != "PRIVATE" {
					t.Fatalf("visibility = %s, want PRIVATE", visibility)
				}
				if state != "NORMAL" {
					t.Fatalf("state = %s, want NORMAL", state)
				}
				if pinned == nil || *pinned {
					t.Fatalf("pinned = %#v, want pointer(false)", pinned)
				}
				switch baseName {
				case "20260316080301_01.md", "20260316080301_02.md":
					if displayTime != "2026-03-16T08:03:01+09:00" {
						t.Fatalf("displayTime = %s, want 2026-03-16T08:03:01+09:00", displayTime)
					}
				case "20260317010101_01.md":
					if displayTime != "2026-03-17T01:01:01+09:00" {
						t.Fatalf("displayTime = %s, want 2026-03-17T01:01:01+09:00", displayTime)
					}
				default:
					t.Fatalf("unexpected contentFile: %s", contentFile)
				}
				return &memos.Memo{Name: "memos/" + strings.TrimSuffix(baseName, filepath.Ext(baseName))}, nil
			},
			patchFilesFunc: func(ctx context.Context, memo string, filePaths []string, replaces bool) (*memos.SetMemoAttachmentsOutput, error) {
				patchMemoID = memo
				patchFiles = append([]string(nil), filePaths...)
				if replaces {
					t.Fatal("replaces = true, want false")
				}
				return &memos.SetMemoAttachmentsOutput{Name: memo}, nil
			},
			addMemoRelationsFunc: func(ctx context.Context, memo string, relatedMemos []string, replaces bool) (*memos.AddMemoRelationsOutput, error) {
				addRelationMemoID = memo
				if len(relatedMemos) != 1 {
					t.Fatalf("relatedMemos = %#v, want len=1", relatedMemos)
				}
				addRelationRelatedMemoID = relatedMemos[0]
				if replaces {
					t.Fatal("replaces = true, want false")
				}
				return &memos.AddMemoRelationsOutput{Memo: memo}, nil
			},
		},
	})

	result, err := service.Execute(context.Background(), common.CreateCommonMemosInput{
		Operation:     common.OperationCreateCommonMemos,
		ContentDir:    contentDir,
		AttachmentDir: attachmentDir,
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result.Total != 3 {
		t.Fatalf("result.Total = %d, want 3", result.Total)
	}
	if len(result.Memos) != 3 {
		t.Fatalf("len(result.Memos) = %d, want 3", len(result.Memos))
	}

	if len(createMemoOrder) != 3 {
		t.Fatalf("len(createMemoOrder) = %d, want 3", len(createMemoOrder))
	}
	if createMemoOrder[0] != "20260316080301_01.md" || createMemoOrder[1] != "20260316080301_02.md" || createMemoOrder[2] != "20260317010101_01.md" {
		t.Fatalf("createMemoOrder = %#v, want sorted order", createMemoOrder)
	}

	if patchMemoID != "memos/20260316080301_02" {
		t.Fatalf("patchMemoID = %s, want memos/20260316080301_02", patchMemoID)
	}
	if len(patchFiles) != 2 || patchFiles[0] != attachment1 || patchFiles[1] != attachment2 {
		t.Fatalf("patchFiles = %#v, want [%s %s]", patchFiles, attachment1, attachment2)
	}

	if addRelationMemoID != "memos/20260316080301_02" {
		t.Fatalf("addRelationMemoID = %s, want memos/20260316080301_02", addRelationMemoID)
	}
	if addRelationRelatedMemoID != "memos/20260316080301_01" {
		t.Fatalf("addRelationRelatedMemoID = %s, want memos/20260316080301_01", addRelationRelatedMemoID)
	}

	if result.Memos[1].RelatedToPreviousBy != "memos/20260316080301_01" {
		t.Fatalf("relatedToPreviousBy = %s, want memos/20260316080301_01", result.Memos[1].RelatedToPreviousBy)
	}
	if result.Memos[2].RelatedToPreviousBy != "" {
		t.Fatalf("relatedToPreviousBy = %s, want empty", result.Memos[2].RelatedToPreviousBy)
	}
}

func TestService_ExecuteProgressReporter_Normal(t *testing.T) {
	contentDir := t.TempDir()
	attachmentDir := t.TempDir()

	content1 := filepath.Join(contentDir, "20260316080301_01.md")
	content2 := filepath.Join(contentDir, "20260316080301_02.md")
	for _, contentFile := range []string{content1, content2} {
		if err := os.WriteFile(contentFile, []byte("# content"), 0o644); err != nil {
			t.Fatalf("WriteFile(%s) error = %v", contentFile, err)
		}
	}
	attachment := filepath.Join(attachmentDir, "20260316080301_02_11.webp")
	if err := os.WriteFile(attachment, []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile(attachment) error = %v", err)
	}

	progresses := make([]common.CreateClipsProgress, 0, 2)
	service := NewService(ServiceOptions{
		FileSystem: infrastructures.NewOSFileSystem(),
		MemosService: &mockMemosService{
			createMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
				return &memos.Memo{Name: "memos/" + strings.TrimSuffix(filepath.Base(contentFile), ".md")}, nil
			},
			addMemoRelationsFunc: func(ctx context.Context, memo string, relatedMemos []string, replaces bool) (*memos.AddMemoRelationsOutput, error) {
				return &memos.AddMemoRelationsOutput{Memo: memo}, nil
			},
		},
		ProgressReporter: func(progress common.CreateClipsProgress) {
			progresses = append(progresses, progress)
		},
	})

	_, err := service.Execute(context.Background(), common.CreateCommonMemosInput{
		Operation:     common.OperationCreateCommonMemos,
		ContentDir:    contentDir,
		AttachmentDir: attachmentDir,
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
	if progresses[0].Operation != common.OperationCreateCommonMemos || progresses[1].Operation != common.OperationCreateCommonMemos {
		t.Fatalf("progress operations = %#v, want %s", progresses, common.OperationCreateCommonMemos)
	}
	if filepath.Base(progresses[0].ContentFile) != "20260316080301_01.md" {
		t.Fatalf("progresses[0].ContentFile = %s, want 20260316080301_01.md", progresses[0].ContentFile)
	}
	if filepath.Base(progresses[1].ContentFile) != "20260316080301_02.md" {
		t.Fatalf("progresses[1].ContentFile = %s, want 20260316080301_02.md", progresses[1].ContentFile)
	}
	if progresses[0].AttachmentCount != 0 {
		t.Fatalf("progresses[0].AttachmentCount = %d, want 0", progresses[0].AttachmentCount)
	}
	if progresses[1].AttachmentCount != 1 {
		t.Fatalf("progresses[1].AttachmentCount = %d, want 1", progresses[1].AttachmentCount)
	}
}

func TestService_ExecuteContentFilenameInvalid_NoMemoCreated_Error(t *testing.T) {
	contentDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(contentDir, "invalid.md"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	service := NewService(ServiceOptions{
		FileSystem: infrastructures.NewOSFileSystem(),
		MemosService: &mockMemosService{
			createMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
				t.Fatal("CreateMemo should not be called")
				return nil, nil
			},
		},
	})

	_, err := service.Execute(context.Background(), common.CreateCommonMemosInput{
		Operation:  common.OperationCreateCommonMemos,
		ContentDir: contentDir,
	})
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "content-dir 内のファイル名が不正です") {
		t.Fatalf("error = %v, want content-dir validation error", err)
	}
}

func TestService_ExecuteAttachmentFilenameInvalid_NoMemoCreated_Error(t *testing.T) {
	contentDir := t.TempDir()
	attachmentDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(contentDir, "20260316080301_01.md"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile(content) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(attachmentDir, "invalid.webp"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile(attachment) error = %v", err)
	}

	service := NewService(ServiceOptions{
		FileSystem: infrastructures.NewOSFileSystem(),
		MemosService: &mockMemosService{
			createMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
				t.Fatal("CreateMemo should not be called")
				return nil, nil
			},
		},
	})

	_, err := service.Execute(context.Background(), common.CreateCommonMemosInput{
		Operation:     common.OperationCreateCommonMemos,
		ContentDir:    contentDir,
		AttachmentDir: attachmentDir,
	})
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "attachment-dir 内のファイル名が不正です") {
		t.Fatalf("error = %v, want attachment-dir validation error", err)
	}
}

func TestService_ExecuteAttachmentPrecheckFailed_NoMemoCreated_Error(t *testing.T) {
	contentDir := t.TempDir()
	attachmentDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(contentDir, "20260316080301_01.md"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile(content) error = %v", err)
	}
	attachmentPath := filepath.Join(attachmentDir, "20260316080301_01_01.webp")
	if err := os.WriteFile(attachmentPath, []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile(attachment) error = %v", err)
	}

	service := NewService(ServiceOptions{
		FileSystem: &testutil.MockFileSystem{
			ReadAttachmentFileFunc: func(filePath string) (*infrastructures.AttachmentFile, error) {
				return nil, errors.New("mock read error")
			},
		},
		MemosService: &mockMemosService{
			createMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
				t.Fatal("CreateMemo should not be called")
				return nil, nil
			},
		},
	})

	_, err := service.Execute(context.Background(), common.CreateCommonMemosInput{
		Operation:     common.OperationCreateCommonMemos,
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

func TestService_ExecutePatchFilesFailed_Error(t *testing.T) {
	contentDir := t.TempDir()
	attachmentDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(contentDir, "20260316080301_01.md"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile(content) error = %v", err)
	}
	attachmentPath := filepath.Join(attachmentDir, "20260316080301_01_01.webp")
	if err := os.WriteFile(attachmentPath, []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile(attachment) error = %v", err)
	}

	service := NewService(ServiceOptions{
		FileSystem: infrastructures.NewOSFileSystem(),
		MemosService: &mockMemosService{
			createMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
				return &memos.Memo{Name: "memos/1"}, nil
			},
			patchFilesFunc: func(ctx context.Context, memo string, filePaths []string, replaces bool) (*memos.SetMemoAttachmentsOutput, error) {
				return nil, errors.New("patch failed")
			},
		},
	})

	_, err := service.Execute(context.Background(), common.CreateCommonMemosInput{
		Operation:     common.OperationCreateCommonMemos,
		ContentDir:    contentDir,
		AttachmentDir: attachmentDir,
	})
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "添付の追加に失敗しました") {
		t.Fatalf("error = %v, want patch failure message", err)
	}
}

func TestService_ExecuteCreateMemoFailed_Error(t *testing.T) {
	contentDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(contentDir, "20260316080301_01.md"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile(content) error = %v", err)
	}

	service := NewService(ServiceOptions{
		FileSystem: infrastructures.NewOSFileSystem(),
		MemosService: &mockMemosService{
			createMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
				return nil, errors.New("create failed")
			},
		},
	})

	_, err := service.Execute(context.Background(), common.CreateCommonMemosInput{
		Operation:  common.OperationCreateCommonMemos,
		ContentDir: contentDir,
	})
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "メモ作成に失敗しました") {
		t.Fatalf("error = %v, want create failure message", err)
	}
}

func TestService_ExecuteAddMemoRelationsFailed_Error(t *testing.T) {
	contentDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(contentDir, "20260316080301_01.md"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile(content1) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(contentDir, "20260316080301_02.md"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile(content2) error = %v", err)
	}

	service := NewService(ServiceOptions{
		FileSystem: infrastructures.NewOSFileSystem(),
		MemosService: &mockMemosService{
			createMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
				return &memos.Memo{Name: "memos/" + strings.TrimSuffix(filepath.Base(contentFile), ".md")}, nil
			},
			addMemoRelationsFunc: func(ctx context.Context, memo string, relatedMemos []string, replaces bool) (*memos.AddMemoRelationsOutput, error) {
				return nil, errors.New("relation failed")
			},
		},
	})

	_, err := service.Execute(context.Background(), common.CreateCommonMemosInput{
		Operation:  common.OperationCreateCommonMemos,
		ContentDir: contentDir,
	})
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "relation 追加に失敗しました") {
		t.Fatalf("error = %v, want relation failure message", err)
	}
}
