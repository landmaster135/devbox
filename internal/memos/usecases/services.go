package usecases

import (
	"context"
	"net/http"
	"time"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
	"github.com/landmaster135/devbox/internal/memos/usecases/common"
	attachments "github.com/landmaster135/devbox/internal/memos/usecases/operations/attachments"
	creatememo "github.com/landmaster135/devbox/internal/memos/usecases/operations/create_memo"
	deletememo "github.com/landmaster135/devbox/internal/memos/usecases/operations/delete_memo"
	getmemo "github.com/landmaster135/devbox/internal/memos/usecases/operations/get_memo"
	listattachments "github.com/landmaster135/devbox/internal/memos/usecases/operations/list_attachments"
	listmemos "github.com/landmaster135/devbox/internal/memos/usecases/operations/list_memos"
	patchfiles "github.com/landmaster135/devbox/internal/memos/usecases/operations/patch_files"
	updatememo "github.com/landmaster135/devbox/internal/memos/usecases/operations/update_memo"
)

const defaultTimeout = 30 * time.Second

// HTTPClient は http.Client と互換なインターフェース。
type HTTPClient = common.HTTPClient

// ServiceOptions は Service 生成時の入力。
type ServiceOptions struct {
	BaseURL    string
	APIToken   string
	Timeout    time.Duration
	HTTPClient HTTPClient
	FileSystem infrastructures.FileSystem
}

type createMemoOperation interface {
	Execute(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*common.Memo, error)
}

type getMemoOperation interface {
	Execute(ctx context.Context, memo string) (*common.Memo, error)
}

type deleteMemoOperation interface {
	Execute(ctx context.Context, memo string, force bool) (*common.DeleteMemoOutput, error)
}

type listMemosOperation interface {
	Execute(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string) (*common.ListMemosOutput, error)
}

type listAttachmentsOperation interface {
	Execute(ctx context.Context, pageSize int, pageToken string, orderBy string, filter string) (*common.ListAttachmentsOutput, error)
}

type updateMemoOperation interface {
	Execute(ctx context.Context, memo string, content string, contentFile string, visibility string, state string, pinned *bool, updateMask []string, displayTime string) (*common.Memo, error)
}

type patchFilesOperation interface {
	Execute(ctx context.Context, memo string, filePaths []string, replaces bool) (*common.SetMemoAttachmentsOutput, error)
}

type createAttachmentOperation interface {
	Create(ctx context.Context, filename string, content []byte, attachmentType string, memo string) (*common.Attachment, error)
}

type listMemoAttachmentsOperation interface {
	List(ctx context.Context, memo string, pageSize int, pageToken string) (*common.ListMemoAttachmentsOutput, error)
}

type setMemoAttachmentsOperation interface {
	Set(ctx context.Context, memo string, attachments []common.Attachment) (*common.SetMemoAttachmentsOutput, error)
}

// Service は Memos API 呼び出しのユースケースを提供する。
type Service struct {
	createMemoOp          createMemoOperation
	getMemoOp             getMemoOperation
	deleteMemoOp          deleteMemoOperation
	listMemosOp           listMemosOperation
	listAttachmentsOp     listAttachmentsOperation
	updateMemoOp          updateMemoOperation
	patchFilesOp          patchFilesOperation
	createAttachmentOp    createAttachmentOperation
	listMemoAttachmentsOp listMemoAttachmentsOperation
	setMemoAttachmentsOp  setMemoAttachmentsOperation
}

// Memo は CLI/上位層に返すメモ情報。
type Memo = common.Memo

// DeleteMemoOutput は DeleteMemo のレスポンス。
type DeleteMemoOutput = common.DeleteMemoOutput

// ListMemosOutput は ListMemos のレスポンス。
type ListMemosOutput = common.ListMemosOutput

// ListAttachmentsOutput は ListAttachments のレスポンス。
type ListAttachmentsOutput = common.ListAttachmentsOutput

// Attachment は Memos API の添付情報。
type Attachment = common.Attachment

// ListMemoAttachmentsOutput は ListMemoAttachments のレスポンス。
type ListMemoAttachmentsOutput = common.ListMemoAttachmentsOutput

// SetMemoAttachmentsOutput は SetMemoAttachments のレスポンス。
type SetMemoAttachmentsOutput = common.SetMemoAttachmentsOutput

// NewService は Service を生成する。
func NewService(opts ServiceOptions) *Service {
	client := opts.HTTPClient
	if client == nil {
		timeout := opts.Timeout
		if timeout <= 0 {
			timeout = defaultTimeout
		}
		client = &http.Client{Timeout: timeout}
	}

	fileSystem := opts.FileSystem
	if fileSystem == nil {
		fileSystem = infrastructures.NewOSFileSystem()
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    opts.BaseURL,
		APIToken:   opts.APIToken,
		HTTPClient: client,
	})

	attachmentsOp := attachments.New(jsonClient)

	return &Service{
		createMemoOp:          creatememo.New(jsonClient, fileSystem),
		getMemoOp:             getmemo.New(jsonClient),
		deleteMemoOp:          deletememo.New(jsonClient),
		listMemosOp:           listmemos.New(jsonClient),
		listAttachmentsOp:     listattachments.New(jsonClient),
		updateMemoOp:          updatememo.New(jsonClient, fileSystem),
		patchFilesOp:          patchfiles.New(fileSystem, attachmentsOp, attachmentsOp, attachmentsOp),
		createAttachmentOp:    attachmentsOp,
		listMemoAttachmentsOp: attachmentsOp,
		setMemoAttachmentsOp:  attachmentsOp,
	}
}
