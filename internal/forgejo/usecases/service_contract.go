package usecases

import (
	forgejo "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"

	commonUsecases "github.com/landmaster135/devbox/internal/forgejo/usecases/common"
)

// RepoRecord は repo list の出力レコードです。
type RepoRecord = commonUsecases.RepoRecord

// IssueRecord は issue list の出力レコードです。
type IssueRecord = commonUsecases.IssueRecord

type repoListOperation interface {
	Execute() ([]RepoRecord, error)
}

type issueListOperation interface {
	Execute() ([]IssueRecord, error)
}

type repoListDependencies struct {
	Client       *forgejo.Client
	Host         string
	Username     string
	Token        string
	HTTPClient   commonUsecases.HTTPClient
	ReposWorkers int
}

type issueListDependencies struct {
	Client   *forgejo.Client
	Username string
}
