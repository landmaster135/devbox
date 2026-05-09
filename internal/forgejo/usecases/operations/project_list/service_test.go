package project_list

import (
	"testing"

	forgejo "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"

	testsupport "github.com/landmaster135/devbox/internal/forgejo/usecases/testsupport"
)

func TestExecute(t *testing.T) {
	server, called := testsupport.NewForgejoTestServer(map[string]testsupport.HandlerResponse{
		"GET /api/v1/users/octocat/repos": {
			Status: 200,
			Body: `[
				{"id":10,"owner":{"login":"octocat"},"name":"project-repo","full_name":"octocat/project-repo","description":"project repo"},
				{"id":11,"owner":{"login":"octocat"},"name":"empty-project","full_name":"octocat/empty-project","description":"no projects"}
			]`,
		},
		"GET /api/v1/repos/octocat/project-repo/projects?state=all": {
			Status: 200,
			Body:   `[{"name":"Backend","title":"Backend Project","description":"infra","is_private":false,"is_archived":false,"created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-02T00:00:00Z"}]`,
		},
		"GET /api/v1/repos/octocat/empty-project/projects?state=all": {
			Status: 200,
			Body:   `[]`,
		},
	})
	defer server.Close()

	service := newServiceForTest(server, t, "octocat")
	records, err := service.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("len(records) = %d, want 1", len(records))
	}

	record := records[0]
	if record.Name != "Backend" {
		t.Fatalf("Name = %q, want %q", record.Name, "Backend")
	}
	if record.RepoFullName != "octocat/project-repo" {
		t.Fatalf("RepoFullName = %q, want %q", record.RepoFullName, "octocat/project-repo")
	}
	if record.CreatedAt != "2024-01-01T00:00:00Z" {
		t.Fatalf("CreatedAt = %q, want %q", record.CreatedAt, "2024-01-01T00:00:00Z")
	}
	if !called("GET /api/v1/repos/octocat/project-repo/projects?state=all") {
		t.Fatalf("project endpoint for project-repo was not requested")
	}
	if !called("GET /api/v1/repos/octocat/empty-project/projects?state=all") {
		t.Fatalf("project endpoint for empty-project was not requested")
	}
}

func TestExecute_NotSupported(t *testing.T) {
	server, _ := testsupport.NewForgejoTestServer(map[string]testsupport.HandlerResponse{
		"GET /api/v1/users/failure/repos": {
			Status: 200,
			Body:   `[{"id":1,"owner":{"login":"failure"},"name":"no-project","full_name":"failure/no-project"}]`,
		},
		"GET /api/v1/repos/failure/no-project/projects?state=all": {
			Status: 404,
			Body:   `not found`,
		},
	})
	defer server.Close()

	service := newServiceForTest(server, t, "failure")
	_, err := service.Execute()
	if err == nil {
		t.Fatalf("Execute() error = nil, want error")
	}
	if err.Error() != "project list API is not supported on this server" {
		t.Fatalf("error = %v, want %q", err, "project list API is not supported on this server")
	}
}

func newServiceForTest(server *testsupport.TestServer, t *testing.T, username string) *Service {
	t.Helper()
	client, err := forgejo.NewClient(
		server.URL,
		forgejo.SetHTTPClient(server.Client()),
		forgejo.SetToken("token"),
		forgejo.SetForgejoVersion(""),
	)
	if err != nil {
		t.Fatalf("forgejo.NewClient() error = %v", err)
	}
	return NewService(Options{
		Client:     client,
		Host:       server.URL,
		Username:   username,
		Token:      "token",
		HTTPClient: server.Client(),
	})
}
