package common

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

// DeleteMemoOutput は DeleteMemo のレスポンス。
type DeleteMemoOutput struct{}

// ListMemosOutput は ListMemos のレスポンス。
type ListMemosOutput struct {
	Memos         []Memo `json:"memos,omitempty"`
	NextPageToken string `json:"nextPageToken,omitempty"`
	TotalSize     int64  `json:"totalSize,omitempty"`
}

// ListAttachmentsOutput は ListAttachments のレスポンス。
type ListAttachmentsOutput struct {
	Attachments   []Attachment `json:"attachments,omitempty"`
	NextPageToken string       `json:"nextPageToken,omitempty"`
	TotalSize     int64        `json:"totalSize,omitempty"`
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
