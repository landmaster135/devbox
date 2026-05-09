package usecases

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestListRepos(t *testing.T) {
	server, called := newForgejoTestServer(map[string]handlerResponse{
		"GET /api/v1/user/repos": {
			status: http.StatusOK,
			body:   `[{"id":1,"owner":{"login":"landmaster135"},"name":"repo1","full_name":"landmaster135/repo1","description":"Repo one","private":false,"html_url":"https://example.com/landmaster135/repo1","open_issues_count":1,"open_pr_counter":2,"forks_count":3,"stars_count":4,"watchers_count":5,"size":123,"archived":false,"created_at":"2022-10-18T00:00:00Z","updated_at":"2022-10-19T00:00:00Z"}]`,
		},
		"GET /api/v1/repos/landmaster135/repo1/topics": {
			status: http.StatusOK,
			body:   `{"topics":["game","demo"]}`,
		},
		"GET /api/v1/repos/landmaster135/repo1/languages": {
			status: http.StatusOK,
			body:   `{"Go":120.0,"C++":40.0}`,
		},
		"GET /api/v1/repos/landmaster135/repo1/pulls": {
			status: http.StatusOK,
			body: `[
				{"id":1,"state":"open","title":"Add README"},
				{"id":2,"state":"open","title":"Fix bug"},
				{"id":3,"state":"closed","title":"Old PR"}
			]`,
		},
	})
	defer server.Close()

	service := newServiceForTest(server, t)
	records, err := service.ListRepos()
	if err != nil {
		t.Fatalf("ListRepos() error = %v", err)
	}

	if len(records) != 1 {
		t.Fatalf("len(records) = %d, want 1", len(records))
	}

	record := records[0]
	if record.Name != "repo1" {
		t.Fatalf("Name = %q, want %q", record.Name, "repo1")
	}
	if record.Language != "Go" {
		t.Fatalf("Language = %q, want %q", record.Language, "Go")
	}
	if record.Tags != "game,demo" {
		t.Fatalf("Tags = %q, want %q", record.Tags, "game,demo")
	}
	if record.OpenPullsCount != 2 {
		t.Fatalf("OpenPullsCount = %d, want %d", record.OpenPullsCount, 2)
	}
	if record.ClosedPullsCount != 1 {
		t.Fatalf("ClosedPullsCount = %d, want %d", record.ClosedPullsCount, 1)
	}
	if record.RepoCreatedAt != "2022-10-18T00:00:00Z" {
		t.Fatalf("RepoCreatedAt = %q, want %q", record.RepoCreatedAt, "2022-10-18T00:00:00Z")
	}
	if !called("GET /api/v1/repos/landmaster135/repo1/topics") {
		t.Fatalf("topics endpoint was not requested")
	}
	if !called("GET /api/v1/repos/landmaster135/repo1/languages") {
		t.Fatalf("languages endpoint was not requested")
	}
	if !called("GET /api/v1/repos/landmaster135/repo1/pulls") {
		t.Fatalf("pulls endpoint was not requested")
	}
}

func TestListProjects(t *testing.T) {
	server, called := newForgejoTestServer(map[string]handlerResponse{
		"GET /api/v1/users/octocat/repos": {
			status: http.StatusOK,
			body: `[
				{"id":10,"owner":{"login":"octocat"},"name":"project-repo","full_name":"octocat/project-repo","description":"project repo"},
				{"id":11,"owner":{"login":"octocat"},"name":"empty-project","full_name":"octocat/empty-project","description":"no projects"}
			]`,
		},
		"GET /api/v1/repos/octocat/project-repo/projects?state=all": {
			status: http.StatusOK,
			body:   `[{"name":"Backend","title":"Backend Project","description":"infra","is_private":false,"is_archived":false,"created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-02T00:00:00Z"}]`,
		},
		"GET /api/v1/repos/octocat/empty-project/projects?state=all": {
			status: http.StatusOK,
			body:   `[]`,
		},
	})
	defer server.Close()

	service := newServiceForTestWithUsername(server, "octocat", t)
	records, err := service.ListProjects()
	if err != nil {
		t.Fatalf("ListProjects() error = %v", err)
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

func TestListProjectsNotSupported(t *testing.T) {
	server, _ := newForgejoTestServer(map[string]handlerResponse{
		"GET /api/v1/users/failure/repos": {
			status: http.StatusOK,
			body:   `[{"id":1,"owner":{"login":"failure"},"name":"no-project","full_name":"failure/no-project"}]`,
		},
		"GET /api/v1/repos/failure/no-project/projects?state=all": {
			status: http.StatusNotFound,
			body:   `not found`,
		},
	})
	defer server.Close()

	service := newServiceForTestWithUsername(server, "failure", t)
	_, err := service.ListProjects()
	if err == nil {
		t.Fatalf("ListProjects() error = nil, want error")
	}
	if err.Error() != "project list API is not supported on this server" {
		t.Fatalf("error = %v, want %q", err, "project list API is not supported on this server")
	}
}

func TestDecodeProjects(t *testing.T) {
	plain := `[{"name":"P1","created_at":"2020-01-01T00:00:00Z","updated_at":"2020-01-02T00:00:00Z"}]`
	decoded, err := decodeProjects([]byte(plain))
	if err != nil {
		t.Fatalf("decodeProjects() error = %v", err)
	}
	if len(decoded) != 1 {
		t.Fatalf("len(decoded) = %d, want 1", len(decoded))
	}

	wrapped := `{"data":[{"name":"P2","title":"Title","created_at":"2020-01-03T00:00:00Z","updated_at":"2020-01-04T00:00:00Z"}], "projects":[]}`
	decoded, err = decodeProjects([]byte(wrapped))
	if err != nil {
		t.Fatalf("decodeProjects() wrapped error = %v", err)
	}
	if len(decoded) != 1 {
		t.Fatalf("len(decoded) = %d, want 1", len(decoded))
	}
	if decoded[0].Name != "P2" {
		t.Fatalf("Name = %q, want %q", decoded[0].Name, "P2")
	}
}

func TestIsNotFoundError(t *testing.T) {
	if ok := isNotFoundError(&requestError{status: http.StatusNotFound, body: "not found"}); !ok {
		t.Fatal("isNotFoundError() should return true")
	}
	if ok := isNotFoundError(nil); ok {
		t.Fatal("isNotFoundError(nil) should return false")
	}
	if ok := isNotFoundError(fmt.Errorf("other")); ok {
		t.Fatal("isNotFoundError(other) should return false")
	}
}

func TestNormalizeHost(t *testing.T) {
	if got := normalizeHost("example.com"); got != "https://example.com" {
		t.Fatalf("normalizeHost() = %q, want %q", got, "https://example.com")
	}
	if got := normalizeHost("https://example.com/"); got != "https://example.com" {
		t.Fatalf("normalizeHost() = %q, want %q", got, "https://example.com")
	}
}

func TestFormatDate(t *testing.T) {
	if got := formatDate(time.Time{}); got != "" {
		t.Fatalf("formatDate(zero) = %q, want %q", got, "")
	}
	if got := formatDate(time.Date(2026, 5, 9, 12, 34, 56, 0, time.UTC)); got != "2026-05-09T12:34:56Z" {
		t.Fatalf("formatDate() = %q, want %q", got, "2026-05-09T12:34:56Z")
	}
}

func TestPrimaryLanguage(t *testing.T) {
	if got := primaryLanguage(map[string]float64{"A": 1, "B": 3, "C": 2}); got != "B" {
		t.Fatalf("primaryLanguage() = %q, want %q", got, "B")
	}
	if got := primaryLanguage(map[string]float64{}); got != "" {
		t.Fatalf("primaryLanguage(empty) = %q, want %q", got, "")
	}
}

func newServiceForTest(server *httptest.Server, t *testing.T) *Service {
	t.Helper()
	service, err := NewService(ServiceOptions{
		Host:       server.URL,
		Username:   "landmaster135",
		Token:      "token",
		HTTPClient: server.Client(),
	})
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}
	return service
}

func newServiceForTestWithUsername(server *httptest.Server, username string, t *testing.T) *Service {
	t.Helper()
	service, err := NewService(ServiceOptions{
		Host:       server.URL,
		Username:   username,
		Token:      "token",
		HTTPClient: server.Client(),
	})
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}
	return service
}

// handlerResponse はテストサーバーのレスポンスを定義する。
type handlerResponse struct {
	status  int
	body    string
	headers map[string]string
}

func newForgejoTestServer(paths map[string]handlerResponse) (*httptest.Server, func(string) bool) {
	pathStates := map[string]int{}
	buildKey := func(r *http.Request) string {
		path := strings.TrimSuffix(r.URL.Path, "/")
		if r.URL.RawQuery == "" {
			return fmt.Sprintf("%s %s", r.Method, path)
		}
		return fmt.Sprintf("%s %s?%s", r.Method, path, r.URL.RawQuery)
	}
	buildPathOnlyKey := func(r *http.Request) string {
		return fmt.Sprintf("%s %s", r.Method, strings.TrimSuffix(r.URL.Path, "/"))
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := buildKey(r)
		pathStates[key]++
		if _, ok := pathStates[buildPathOnlyKey(r)]; !ok && key != buildPathOnlyKey(r) {
			pathStates[buildPathOnlyKey(r)] = 0
		}
		if response, ok := paths[key]; ok {
			w.WriteHeader(response.status)
			for headerKey, headerValue := range response.headers {
				w.Header().Set(headerKey, headerValue)
			}
			_, _ = w.Write([]byte(response.body))
			return
		}
		if response, ok := paths[buildPathOnlyKey(r)]; ok {
			pathStates[buildPathOnlyKey(r)]++
			w.WriteHeader(response.status)
			for headerKey, headerValue := range response.headers {
				w.Header().Set(headerKey, headerValue)
			}
			_, _ = w.Write([]byte(response.body))
			return
		}

		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	})

	server := httptest.NewServer(handler)
	called := func(path string) bool {
		_, ok := pathStates[path]
		return ok
	}
	return server, called
}
