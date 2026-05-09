package issue_list

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
				{"id":10,"owner":{"login":"octocat"},"name":"repo-a","full_name":"octocat/repo-a"},
				{"id":11,"owner":{"login":"octocat"},"name":"repo-b","full_name":"octocat/repo-b"}
			]`,
		},
		"GET /api/v1/repos/octocat/repo-a/issues?limit=100&page=1&state=all&type=issues": {
			Status: 200,
			Body: `[
				{
					"id":100,
					"number":5,
					"title":"Bug fix",
					"state":"open",
					"html_url":"https://example.com/octocat/repo-a/issues/5",
					"user":{"login":"alice"},
					"labels":[{"name":"bug"}],
					"assignees":[{"login":"bob"}],
					"comments":2,
					"is_locked":false,
					"created_at":"2024-01-01T00:00:00Z",
					"updated_at":"2024-01-02T00:00:00Z",
					"closed_at":null
				}
			]`,
		},
		"GET /api/v1/repos/octocat/repo-b/issues?limit=100&page=1&state=all&type=issues": {
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
	if record.RepoFullName != "octocat/repo-a" {
		t.Fatalf("RepoFullName = %q, want %q", record.RepoFullName, "octocat/repo-a")
	}
	if record.Number != 5 {
		t.Fatalf("Number = %d, want %d", record.Number, 5)
	}
	if record.Title != "Bug fix" {
		t.Fatalf("Title = %q, want %q", record.Title, "Bug fix")
	}
	if record.Author != "alice" {
		t.Fatalf("Author = %q, want %q", record.Author, "alice")
	}
	if len(record.Assignees) != 1 || record.Assignees[0] != "bob" {
		t.Fatalf("Assignees = %#v, want %#v", record.Assignees, []string{"bob"})
	}
	if len(record.Labels) != 1 || record.Labels[0] != "bug" {
		t.Fatalf("Labels = %#v, want %#v", record.Labels, []string{"bug"})
	}
	if !called("GET /api/v1/repos/octocat/repo-a/issues?limit=100&page=1&state=all&type=issues") &&
		!called("GET /api/v1/repos/octocat/repo-a/issues?state=all&type=issues&limit=100&page=1") &&
		!called("GET /api/v1/repos/octocat/repo-a/issues?type=issues&state=all&limit=100&page=1") {
		t.Fatalf("issue endpoint for repo-a was not requested")
	}
}

func TestExecute_WithPagination(t *testing.T) {
	server, called := testsupport.NewForgejoTestServer(map[string]testsupport.HandlerResponse{
		"GET /api/v1/users/octocat/repos": {
			Status: 200,
			Body:   `[{"id":10,"owner":{"login":"octocat"},"name":"repo-a","full_name":"octocat/repo-a"}]`,
		},
		"GET /api/v1/repos/octocat/repo-a/issues?limit=100&page=1&state=all&type=issues": {
			Status: 200,
			Headers: map[string]string{
				"Link": "</api/v1/repos/octocat/repo-a/issues?page=2&limit=100&state=all&type=issues>; rel=\"next\", </api/v1/repos/octocat/repo-a/issues?page=2&limit=100&state=all&type=issues>; rel=\"last\"",
			},
			Body: `[{"id":100,"number":1,"title":"one","state":"open","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z"}]`,
		},
		"GET /api/v1/repos/octocat/repo-a/issues?limit=100&page=2&state=all&type=issues": {
			Status: 200,
			Body:   `[{"id":101,"number":2,"title":"two","state":"closed","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-02T00:00:00Z","closed_at":"2024-01-03T00:00:00Z"}]`,
		},
	})
	defer server.Close()

	service := newServiceForTest(server, t, "octocat")
	records, err := service.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("len(records) = %d, want 2", len(records))
	}
	if records[1].ClosedAt != "2024-01-03T00:00:00Z" {
		t.Fatalf("ClosedAt = %q, want %q", records[1].ClosedAt, "2024-01-03T00:00:00Z")
	}
	if !called("GET /api/v1/repos/octocat/repo-a/issues?limit=100&page=2&state=all&type=issues") &&
		!called("GET /api/v1/repos/octocat/repo-a/issues?state=all&type=issues&limit=100&page=2") &&
		!called("GET /api/v1/repos/octocat/repo-a/issues?type=issues&state=all&limit=100&page=2") {
		t.Fatalf("second issue page was not requested")
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
		Client:   client,
		Username: username,
	})
}
