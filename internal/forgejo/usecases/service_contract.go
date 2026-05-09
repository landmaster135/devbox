package usecases

import (
	forgejo "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"

	commonUsecases "github.com/landmaster135/devbox/internal/forgejo/usecases/common"
)

// RepoRecord は repo list の出力レコードです。
type RepoRecord = commonUsecases.RepoRecord

// ProjectRecord は project list の出力レコードです。
type ProjectRecord = commonUsecases.ProjectRecord

type repoListOperation interface {
	Execute() ([]RepoRecord, error)
}

type projectListOperation interface {
	Execute() ([]ProjectRecord, error)
}

type repoListDependencies struct {
	Client       *forgejo.Client
	Host         string
	Username     string
	Token        string
	HTTPClient   commonUsecases.HTTPClient
	ReposWorkers int
}

type projectListDependencies struct {
	Client     *forgejo.Client
	Host       string
	Username   string
	Token      string
	HTTPClient commonUsecases.HTTPClient
}
