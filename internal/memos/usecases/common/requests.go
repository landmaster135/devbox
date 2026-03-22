package common

// MemoMutationRequest は create/update の共通 payload。
type MemoMutationRequest struct {
	Content     string `json:"content,omitempty"`
	Visibility  string `json:"visibility,omitempty"`
	State       string `json:"state,omitempty"`
	Pinned      *bool  `json:"pinned,omitempty"`
	DisplayTime string `json:"displayTime,omitempty"`
}

// CreateAttachmentRequest は添付作成 payload。
type CreateAttachmentRequest struct {
	Filename     string `json:"filename,omitempty"`
	Content      []byte `json:"content,omitempty"`
	ExternalLink string `json:"externalLink,omitempty"`
	Type         string `json:"type,omitempty"`
	Memo         string `json:"memo,omitempty"`
}

// SetMemoAttachmentsRequest は添付更新 payload。
type SetMemoAttachmentsRequest struct {
	Name        string       `json:"name,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// SetMemoRelationsRequest はリレーション更新 payload。
type SetMemoRelationsRequest struct {
	Name      string         `json:"name,omitempty"`
	Relations []MemoRelation `json:"relations,omitempty"`
}
