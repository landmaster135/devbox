package usecases

import (
	"net/http"
	"strings"
	"time"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
)

const defaultTimeout = 30 * time.Second

// HTTPClient は http.Client と互換なインターフェース。
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ServiceOptions は Service 生成時の入力。
type ServiceOptions struct {
	BaseURL    string
	APIToken   string
	Timeout    time.Duration
	HTTPClient HTTPClient
	FileSystem infrastructures.FileSystem
}

// Service は Memos API 呼び出しのユースケースを提供する。
type Service struct {
	baseURL    string
	apiToken   string
	client     HTTPClient
	fileSystem infrastructures.FileSystem
}

// Memo は CLI/上位層に返すメモ情報。
type Memo struct {
	Name        string `json:"name,omitempty"`
	UID         string `json:"uid,omitempty"`
	ID          int64  `json:"id,omitempty"`
	CreateTime  string `json:"createTime,omitempty"`
	UpdateTime  string `json:"updateTime,omitempty"`
	DisplayTime string `json:"displayTime,omitempty"`
	Content     string `json:"content,omitempty"`
	Visibility  string `json:"visibility,omitempty"`
	State       string `json:"state,omitempty"`
	Pinned      bool   `json:"pinned,omitempty"`
}

// ListMemosOutput は ListMemos のレスポンス。
type ListMemosOutput struct {
	Memos         []Memo `json:"memos,omitempty"`
	NextPageToken string `json:"nextPageToken,omitempty"`
	TotalSize     int64  `json:"totalSize,omitempty"`
}

// Attachment は Memos API の添付情報。
type Attachment struct {
	Name         string `json:"name,omitempty"`
	Filename     string `json:"filename,omitempty"`
	ExternalLink string `json:"externalLink,omitempty"`
	Type         string `json:"type,omitempty"`
	Memo         string `json:"memo,omitempty"`
}

// ListMemoAttachmentsOutput は ListMemoAttachments のレスポンス。
type ListMemoAttachmentsOutput struct {
	Attachments   []Attachment `json:"attachments,omitempty"`
	NextPageToken string       `json:"nextPageToken,omitempty"`
}

// SetMemoAttachmentsOutput は SetMemoAttachments のレスポンス。
type SetMemoAttachmentsOutput struct {
	Name        string       `json:"name,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type memoMutationRequest struct {
	Content     string `json:"content,omitempty"`
	Visibility  string `json:"visibility,omitempty"`
	State       string `json:"state,omitempty"`
	Pinned      *bool  `json:"pinned,omitempty"`
	DisplayTime string `json:"displayTime,omitempty"`
}

type createAttachmentRequest struct {
	Filename     string `json:"filename,omitempty"`
	Content      []byte `json:"content,omitempty"`
	ExternalLink string `json:"externalLink,omitempty"`
	Type         string `json:"type,omitempty"`
	Memo         string `json:"memo,omitempty"`
}

type setMemoAttachmentsRequest struct {
	Name        string       `json:"name,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

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

	return &Service{
		baseURL:    normalizeBaseURL(opts.BaseURL),
		apiToken:   strings.TrimSpace(opts.APIToken),
		client:     client,
		fileSystem: fileSystem,
	}
}
