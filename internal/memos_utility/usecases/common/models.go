package common

import memos "github.com/landmaster135/devbox/internal/memos/usecases"

// CreateClipInput は create-web-clip / create-movie-clip の入力。
type CreateClipInput struct {
	Operation   string
	ContentFile string
	Attachments []string
}

// CreateClipOutput は create-web-clip / create-movie-clip の出力。
type CreateClipOutput struct {
	Operation          string                          `json:"operation"`
	DisplayTime        string                          `json:"displayTime"`
	Memo               *memos.Memo                     `json:"memo,omitempty"`
	Attachments        []string                        `json:"attachments,omitempty"`
	SetMemoAttachments *memos.SetMemoAttachmentsOutput `json:"setMemoAttachments,omitempty"`
}

// CreateClipsInput は create-clips の入力。
type CreateClipsInput struct {
	Operation     string
	ContentDir    string
	AttachmentDir string
}

// CreateClipsOutput は create-clips の出力。
type CreateClipsOutput struct {
	Operation     string              `json:"operation"`
	ContentDir    string              `json:"contentDir"`
	AttachmentDir string              `json:"attachmentDir,omitempty"`
	Total         int                 `json:"total"`
	Clips         []*CreateClipOutput `json:"clips"`
}

// CreateClipsProgress は create-clips の進捗通知情報。
type CreateClipsProgress struct {
	Current         int
	Total           int
	Operation       string
	ContentFile     string
	AttachmentCount int
}

const (
	CreateCommonMemosRelationPhaseStart = "start"
	CreateCommonMemosRelationPhaseOK    = "ok"
	CreateCommonMemosRelationPhaseError = "error"
)

// CreateCommonMemosRelationProgress は create-common-memos の relation 処理ログ通知情報。
type CreateCommonMemosRelationProgress struct {
	Phase                        string
	ContentFile                  string
	CurrentMemoIdentifier        string
	CurrentMemoIdentifierSource  string
	PreviousMemoIdentifier       string
	PreviousMemoIdentifierSource string
	ErrorMessage                 string
}

// CreateCommonMemosInput は create-common-memos の入力。
type CreateCommonMemosInput struct {
	Operation     string
	ContentDir    string
	AttachmentDir string
}

// CreateCommonMemoOutput は create-common-memos の各メモ作成結果。
type CreateCommonMemoOutput struct {
	ContentFile         string                          `json:"contentFile"`
	DisplayTime         string                          `json:"displayTime"`
	Memo                *memos.Memo                     `json:"memo,omitempty"`
	Attachments         []string                        `json:"attachments,omitempty"`
	SetMemoAttachments  *memos.SetMemoAttachmentsOutput `json:"setMemoAttachments,omitempty"`
	RelatedToPreviousBy string                          `json:"relatedToPreviousBy,omitempty"`
}

// CreateCommonMemosOutput は create-common-memos の出力。
type CreateCommonMemosOutput struct {
	Operation     string                    `json:"operation"`
	ContentDir    string                    `json:"contentDir"`
	AttachmentDir string                    `json:"attachmentDir,omitempty"`
	Total         int                       `json:"total"`
	Memos         []*CreateCommonMemoOutput `json:"memos"`
}
