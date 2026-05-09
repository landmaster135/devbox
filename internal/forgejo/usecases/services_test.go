package usecases

import (
	"errors"
	"testing"
)

func TestNormalizeHost(t *testing.T) {
	if got := normalizeHost("example.com"); got != "https://example.com" {
		t.Fatalf("normalizeHost() = %q, want %q", got, "https://example.com")
	}
	if got := normalizeHost("https://example.com/"); got != "https://example.com" {
		t.Fatalf("normalizeHost() = %q, want %q", got, "https://example.com")
	}
}

func TestNewService_HostRequired(t *testing.T) {
	_, err := NewService(ServiceOptions{
		Host: "",
	})
	if err == nil {
		t.Fatal("NewService() error = nil, want error")
	}
}

func TestListRepos_DelegatesToOperation(t *testing.T) {
	oldFactory := newRepoListOperation
	defer func() { newRepoListOperation = oldFactory }()

	expected := []RepoRecord{{Name: "repo1"}}
	newRepoListOperation = func(dependencies repoListDependencies) repoListOperation {
		return repoListOperationMock{
			execute: func() ([]RepoRecord, error) {
				if dependencies.ReposWorkers != 3 {
					t.Fatalf("ReposWorkers = %d, want 3", dependencies.ReposWorkers)
				}
				return expected, nil
			},
		}
	}

	service := &Service{reposWorkers: 3}
	got, err := service.ListRepos()
	if err != nil {
		t.Fatalf("ListRepos() error = %v", err)
	}
	if len(got) != 1 || got[0].Name != "repo1" {
		t.Fatalf("ListRepos() = %#v, want %#v", got, expected)
	}
}

func TestListProjects_DelegatesToOperation(t *testing.T) {
	oldFactory := newProjectListOperation
	defer func() { newProjectListOperation = oldFactory }()

	expectedErr := errors.New("boom")
	newProjectListOperation = func(_ projectListDependencies) projectListOperation {
		return projectListOperationMock{
			execute: func() ([]ProjectRecord, error) {
				return nil, expectedErr
			},
		}
	}

	service := &Service{}
	_, err := service.ListProjects()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("ListProjects() error = %v, want %v", err, expectedErr)
	}
}

type repoListOperationMock struct {
	execute func() ([]RepoRecord, error)
}

func (m repoListOperationMock) Execute() ([]RepoRecord, error) {
	return m.execute()
}

type projectListOperationMock struct {
	execute func() ([]ProjectRecord, error)
}

func (m projectListOperationMock) Execute() ([]ProjectRecord, error) {
	return m.execute()
}
