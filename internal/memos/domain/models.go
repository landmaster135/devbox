package domain

const (
	VisibilityUnspecified = "VISIBILITY_UNSPECIFIED"
	VisibilityPrivate     = "PRIVATE"
	VisibilityProtected   = "PROTECTED"
	VisibilityPublic      = "PUBLIC"

	StateUnspecified = "STATE_UNSPECIFIED"
	StateNormal      = "NORMAL"
	StateArchived    = "ARCHIVED"
)

// Memo は Memos API のメモ情報を表す。
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

// CreateMemoParams は CreateMemo の入力。
type CreateMemoParams struct {
	MemoID      string
	Content     string
	Visibility  string
	State       string
	Pinned      *bool
	DisplayTime string
}

// ListMemosParams は ListMemos の入力。
type ListMemosParams struct {
	PageSize  int
	PageToken string
	State     string
	OrderBy   string
	Filter    string
}

// UpdateMemoParams は UpdateMemo の入力。
type UpdateMemoParams struct {
	Memo       string
	Content    string
	Visibility string
	State      string
	Pinned     *bool
	UpdateMask []string
}

// ListMemosResponse は ListMemos のレスポンス。
type ListMemosResponse struct {
	Memos         []Memo `json:"memos,omitempty"`
	NextPageToken string `json:"nextPageToken,omitempty"`
	TotalSize     int64  `json:"totalSize,omitempty"`
}
