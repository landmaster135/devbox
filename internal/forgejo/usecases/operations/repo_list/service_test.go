package repo_list

import (
	"sync/atomic"
	"testing"
	"time"

	forgejo "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"

	testsupport "github.com/landmaster135/devbox/internal/forgejo/usecases/testsupport"
)

func TestExecute(t *testing.T) {
	server, called := testsupport.NewForgejoTestServer(map[string]testsupport.HandlerResponse{
		"GET /api/v1/user/repos": {
			Status: httpStatusOK,
			Body:   `[{"id":1,"owner":{"login":"landmaster135"},"name":"repo1","full_name":"landmaster135/repo1","description":"Repo one","private":false,"html_url":"https://example.com/landmaster135/repo1","open_issues_count":1,"open_pr_counter":2,"forks_count":3,"stars_count":4,"watchers_count":5,"size":123,"archived":false,"created_at":"2022-10-18T00:00:00Z","updated_at":"2022-10-19T00:00:00Z"}]`,
		},
		"GET /api/v1/repos/landmaster135/repo1/topics": {
			Status: httpStatusOK,
			Body:   `{"topics":["game","demo"]}`,
		},
		"GET /api/v1/repos/landmaster135/repo1/languages": {
			Status: httpStatusOK,
			Body:   `{"Go":120.0,"C++":40.0}`,
		},
		"GET /api/v1/repos/landmaster135/repo1/pulls": {
			Status: httpStatusOK,
			Body: `[{
				"id":3,
				"state":"closed",
				"title":"Old PR"
			}]`,
		},
		"GET /api/v1/repos/landmaster135/repo1/issues": {
			Status: httpStatusOK,
			Headers: map[string]string{
				"X-Total-Count": "2",
			},
			Body: `[{
				"id":9,
				"state":"closed",
				"title":"Old issue"
			}]`,
		},
	})
	defer server.Close()

	service := newServiceForTest(server, t, 4)
	records, err := service.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
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
	if record.OpenIssuesCount != 1 {
		t.Fatalf("OpenIssuesCount = %d, want %d", record.OpenIssuesCount, 1)
	}
	if record.ClosedPullsCount != 1 {
		t.Fatalf("ClosedPullsCount = %d, want %d", record.ClosedPullsCount, 1)
	}
	if record.ClosedIssuesCount != 2 {
		t.Fatalf("ClosedIssuesCount = %d, want %d", record.ClosedIssuesCount, 2)
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
	if !called("GET /api/v1/repos/landmaster135/repo1/pulls") &&
		!called("GET /api/v1/repos/landmaster135/repo1/pulls?state=closed") &&
		!called("GET /api/v1/repos/landmaster135/repo1/pulls?limit=100&state=closed") &&
		!called("GET /api/v1/repos/landmaster135/repo1/pulls?state=closed&limit=100") {
		t.Fatalf("pulls endpoint was not requested")
	}
	if !called("GET /api/v1/repos/landmaster135/repo1/issues") &&
		!called("GET /api/v1/repos/landmaster135/repo1/issues?limit=1&page=1&state=closed&type=issues") &&
		!called("GET /api/v1/repos/landmaster135/repo1/issues?state=closed&limit=1&type=issues&page=1") &&
		!called("GET /api/v1/repos/landmaster135/repo1/issues?type=issues&state=closed&limit=1&page=1") &&
		!called("GET /api/v1/repos/landmaster135/repo1/issues?state=closed&type=issues") {
		t.Fatalf("issues endpoint was not requested")
	}
}

func TestExecute_WithMultiplePullPages(t *testing.T) {
	server, called := testsupport.NewForgejoTestServer(map[string]testsupport.HandlerResponse{
		"GET /api/v1/user/repos": {
			Status: httpStatusOK,
			Body:   `[{"id":1,"owner":{"login":"landmaster135"},"name":"repo1","full_name":"landmaster135/repo1","description":"Repo one","private":false,"html_url":"https://example.com/landmaster135/repo1","open_issues_count":1,"open_pr_counter":2,"forks_count":3,"stars_count":4,"watchers_count":5,"size":123,"archived":false,"created_at":"2022-10-18T00:00:00Z","updated_at":"2022-10-19T00:00:00Z"}]`,
		},
		"GET /api/v1/repos/landmaster135/repo1/topics": {
			Status: httpStatusOK,
			Body:   `{"topics":["game","demo"]}`,
		},
		"GET /api/v1/repos/landmaster135/repo1/languages": {
			Status: httpStatusOK,
			Body:   `{"Go":120.0,"C++":40.0}`,
		},
		"GET /api/v1/repos/landmaster135/repo1/pulls": {
			Status: httpStatusOK,
			Headers: map[string]string{
				"X-Total-Count": "2",
			},
			Body: `[{"id":3,"state":"closed","title":"Old PR"}]`,
		},
		"GET /api/v1/repos/landmaster135/repo1/issues": {
			Status: httpStatusOK,
			Headers: map[string]string{
				"X-Total-Count": "7",
			},
			Body: `[{"id":9,"state":"closed","title":"Old issue"}]`,
		},
	})
	defer server.Close()

	service := newServiceForTest(server, t, 4)
	records, err := service.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("len(records) = %d, want 1", len(records))
	}

	record := records[0]
	if record.OpenPullsCount != 2 {
		t.Fatalf("OpenPullsCount = %d, want %d", record.OpenPullsCount, 2)
	}
	if record.OpenIssuesCount != 1 {
		t.Fatalf("OpenIssuesCount = %d, want %d", record.OpenIssuesCount, 1)
	}
	if record.ClosedPullsCount != 2 {
		t.Fatalf("ClosedPullsCount = %d, want %d", record.ClosedPullsCount, 2)
	}
	if record.ClosedIssuesCount != 7 {
		t.Fatalf("ClosedIssuesCount = %d, want %d", record.ClosedIssuesCount, 7)
	}
	if called("GET /api/v1/repos/landmaster135/repo1/pulls?page=2&limit=1&state=closed") ||
		called("GET /api/v1/repos/landmaster135/repo1/pulls?limit=1&page=2&state=closed") ||
		called("GET /api/v1/repos/landmaster135/repo1/pulls?state=closed&page=2&limit=1") ||
		called("GET /api/v1/repos/landmaster135/repo1/pulls?page=2&state=closed&limit=1") {
		t.Fatalf("second pulls page was requested")
	}
}

func TestExecute_UsesWorkerCount(t *testing.T) {
	paths := map[string]testsupport.HandlerResponse{
		"GET /api/v1/user/repos": {
			Status: httpStatusOK,
			Body: `[{"id":1,"owner":{"login":"landmaster135"},"name":"repo1","full_name":"landmaster135/repo1","description":"Repo one","private":false,"html_url":"https://example.com/landmaster135/repo1","open_issues_count":1,"open_pr_counter":2,"forks_count":3,"stars_count":4,"watchers_count":5,"size":123,"archived":false,"created_at":"2022-10-18T00:00:00Z","updated_at":"2022-10-19T00:00:00Z"},
				{"id":2,"owner":{"login":"landmaster135"},"name":"repo2","full_name":"landmaster135/repo2","description":"Repo two","private":false,"html_url":"https://example.com/landmaster135/repo2","open_issues_count":0,"open_pr_counter":1,"forks_count":1,"stars_count":2,"watchers_count":3,"size":45,"archived":false,"created_at":"2022-10-20T00:00:00Z","updated_at":"2022-10-21T00:00:00Z"}]`,
		},
		"GET /api/v1/repos/landmaster135/repo1/topics": {
			Status: httpStatusOK,
			Body:   `{"topics":["game"]}`,
		},
		"GET /api/v1/repos/landmaster135/repo1/languages": {
			Status: httpStatusOK,
			Body:   `{"Go":120.0,"C++":40.0}`,
		},
		"GET /api/v1/repos/landmaster135/repo1/pulls": {
			Status: httpStatusOK,
			Headers: map[string]string{
				"X-Total-Count": "10",
			},
			Body: `[{"id":3,"state":"closed","title":"Old PR"}]`,
		},
		"GET /api/v1/repos/landmaster135/repo1/issues": {
			Status: httpStatusOK,
			Headers: map[string]string{
				"X-Total-Count": "21",
			},
			Body: `[{"id":9,"state":"closed","title":"Old issue"}]`,
		},
		"GET /api/v1/repos/landmaster135/repo2/topics": {
			Status: httpStatusOK,
			Body:   `{"topics":["dev"]}`,
		},
		"GET /api/v1/repos/landmaster135/repo2/languages": {
			Status: httpStatusOK,
			Body:   `{"Rust":200.0}`,
		},
		"GET /api/v1/repos/landmaster135/repo2/pulls": {
			Status: httpStatusOK,
			Headers: map[string]string{
				"X-Total-Count": "7",
			},
			Body: `[{"id":4,"state":"closed","title":"Old PR2"}]`,
		},
		"GET /api/v1/repos/landmaster135/repo2/issues": {
			Status: httpStatusOK,
			Headers: map[string]string{
				"X-Total-Count": "3",
			},
			Body: `[{"id":10,"state":"closed","title":"Closed issue"}]`,
		},
	}

	var activeCount int64
	var maxActiveCount int64
	server, _ := testsupport.NewForgejoTestServerWithRequestDelay(paths, &activeCount, &maxActiveCount, 40*time.Millisecond)
	defer server.Close()

	{
		service := newServiceForTest(server, t, 1)
		if _, err := service.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if got := atomic.LoadInt64(&maxActiveCount); got != 1 {
			t.Fatalf("worker=1 -> maxActiveCount = %d, want 1", got)
		}
	}

	atomic.StoreInt64(&activeCount, 0)
	atomic.StoreInt64(&maxActiveCount, 0)

	{
		service := newServiceForTest(server, t, 4)
		if _, err := service.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if got := atomic.LoadInt64(&maxActiveCount); got < 2 {
			t.Fatalf("worker=4 -> maxActiveCount = %d, want >= 2", got)
		}
	}
}

func TestPullsListResponsePaginationHeader(t *testing.T) {
	server, _ := testsupport.NewForgejoTestServer(map[string]testsupport.HandlerResponse{
		"GET /api/v1/repos/landmaster135/repo1/pulls": {
			Status: httpStatusOK,
			Headers: map[string]string{
				"Link": "</api/v1/repos/landmaster135/repo1/pulls?page=2&limit=100&state=all>; rel=\"next\", </api/v1/repos/landmaster135/repo1/pulls?page=2&limit=100&state=all>; rel=\"last\"",
			},
			Body: `[{"id":1,"state":"open","title":"Add README"},{"id":2,"state":"open","title":"Fix bug"}]`,
		},
	})
	defer server.Close()

	service := newServiceForTest(server, t, 4)
	_, response, err := service.client.ListRepoPullRequests("landmaster135", "repo1", forgejo.ListPullRequestsOptions{
		ListOptions: forgejo.ListOptions{
			Page:     1,
			PageSize: 100,
		},
		State: forgejo.StateAll,
	})
	if err != nil {
		t.Fatalf("ListRepoPullRequests() error = %v", err)
	}
	if response == nil {
		t.Fatal("response = nil")
	}
	if response.LastPage != 2 {
		t.Fatalf("response.LastPage = %d (want 2)", response.LastPage)
	}
}

const httpStatusOK = 200

func newServiceForTest(server *testsupport.TestServer, t *testing.T, reposWorkers int) *Service {
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
		Client:       client,
		Host:         server.URL,
		Username:     "landmaster135",
		Token:        "token",
		HTTPClient:   server.Client(),
		ReposWorkers: reposWorkers,
	})
}
