package common

import "time"

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

// ProjectRecord は project list の出力レコードです。
type ProjectRecord struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	IsPrivate    bool   `json:"is_private"`
	IsArchived   bool   `json:"is_archived"`
	RepoFullName string `json:"repo_full_name"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// ProjectResponse は Forgejo project API 応答の最小形です。
type ProjectResponse struct {
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsPrivate   bool      `json:"is_private"`
	Archived    bool      `json:"is_archived"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
