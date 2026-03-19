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
