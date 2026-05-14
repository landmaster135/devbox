package common

// RepoRecord は repo list の出力レコードです。
type RepoRecord struct {
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	IsPrivate         bool               `json:"is_private"`
	HTTPURL           string             `json:"http_url"`
	OpenIssuesCount   int                `json:"open_issues_count"`
	ClosedIssuesCount int                `json:"closed_issues_count"`
	OpenPullsCount    int                `json:"open_pulls_count"`
	ClosedPullsCount  int                `json:"closed_pulls_count"`
	ForksCount        int                `json:"forks_count"`
	StargazersCount   int                `json:"stargazers_count"`
	SubscribersCount  int                `json:"subscribers_count"`
	Language          string             `json:"language"`
	Languages         map[string]float64 `json:"languages"`
	Size              int                `json:"size"`
	RepoCreatedAt     string             `json:"repo_created_at"`
	RepoUpdatedAt     string             `json:"repo_updated_at"`
	IsArchived        bool               `json:"is_archived"`
	Tags              string             `json:"tags"`
}

// IssueRecord は issue list の出力レコードです。
type IssueRecord struct {
	RepoFullName string   `json:"repo_full_name"`
	Number       int64    `json:"number"`
	Title        string   `json:"title"`
	State        string   `json:"state"`
	HTMLURL      string   `json:"html_url"`
	Author       string   `json:"author"`
	Assignees    []string `json:"assignees"`
	Labels       []string `json:"labels"`
	Comments     int      `json:"comments"`
	IsLocked     bool     `json:"is_locked"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
	ClosedAt     string   `json:"closed_at"`
}
