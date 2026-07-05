package repo_list

import (
	"fmt"
	"strings"
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
			Body: repoListBody(repoBody(
				1,
				"repo1",
				"Repo one",
				false,
				1,
				2,
				3,
				4,
				5,
				123,
				"2022-10-18T00:00:00Z",
				"2022-10-19T00:00:00Z",
			)),
		},
		repoAPIPath("repo1", "topics"): {
			Status: httpStatusOK,
			Body:   `{"topics":["game","demo"]}`,
		},
		repoAPIPath("repo1", "languages"): {
			Status: httpStatusOK,
			Body:   `{"Go":120.0,"C++":40.0}`,
		},
		repoAPIPath("repo1", "pulls"): {
			Status: httpStatusOK,
			Body: `[{
				"id":3,
				"state":"closed",
				"title":"Old PR"
			}]`,
		},
		repoAPIPath("repo1", "issues"): {
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
	if !called(repoAPIPath("repo1", "topics")) {
		t.Fatalf("topics endpoint was not requested")
	}
	if !called(repoAPIPath("repo1", "languages")) {
		t.Fatalf("languages endpoint was not requested")
	}
	if !called(repoAPIPath("repo1", "pulls")) &&
		!called(repoAPIPathWithQuery("repo1", "pulls", "state=closed")) &&
		!called(repoAPIPathWithQuery("repo1", "pulls", "limit=100&state=closed")) &&
		!called(repoAPIPathWithQuery("repo1", "pulls", "state=closed&limit=100")) {
		t.Fatalf("pulls endpoint was not requested")
	}
	if !called(repoAPIPath("repo1", "issues")) &&
		!called(repoAPIPathWithQuery("repo1", "issues", "limit=1&page=1&state=closed&type=issues")) &&
		!called(repoAPIPathWithQuery("repo1", "issues", "state=closed&limit=1&type=issues&page=1")) &&
		!called(repoAPIPathWithQuery("repo1", "issues", "type=issues&state=closed&limit=1&page=1")) &&
		!called(repoAPIPathWithQuery("repo1", "issues", "state=closed&type=issues")) {
		t.Fatalf("issues endpoint was not requested")
	}
}

func TestExecute_WithMultiplePullPages(t *testing.T) {
	server, called := testsupport.NewForgejoTestServer(map[string]testsupport.HandlerResponse{
		"GET /api/v1/user/repos": {
			Status: httpStatusOK,
			Body: repoListBody(repoBody(
				1,
				"repo1",
				"Repo one",
				false,
				1,
				2,
				3,
				4,
				5,
				123,
				"2022-10-18T00:00:00Z",
				"2022-10-19T00:00:00Z",
			)),
		},
		repoAPIPath("repo1", "topics"): {
			Status: httpStatusOK,
			Body:   `{"topics":["game","demo"]}`,
		},
		repoAPIPath("repo1", "languages"): {
			Status: httpStatusOK,
			Body:   `{"Go":120.0,"C++":40.0}`,
		},
		repoAPIPath("repo1", "pulls"): {
			Status: httpStatusOK,
			Headers: map[string]string{
				"X-Total-Count": "2",
			},
			Body: `[{"id":3,"state":"closed","title":"Old PR"}]`,
		},
		repoAPIPath("repo1", "issues"): {
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
	if called(repoAPIPathWithQuery("repo1", "pulls", "page=2&limit=1&state=closed")) ||
		called(repoAPIPathWithQuery("repo1", "pulls", "limit=1&page=2&state=closed")) ||
		called(repoAPIPathWithQuery("repo1", "pulls", "state=closed&page=2&limit=1")) ||
		called(repoAPIPathWithQuery("repo1", "pulls", "page=2&state=closed&limit=1")) {
		t.Fatalf("second pulls page was requested")
	}
}

func TestExecute_WithMultipleRepoPages_Normal(t *testing.T) {
	server, called := testsupport.NewForgejoTestServer(map[string]testsupport.HandlerResponse{
		"GET /api/v1/user/repos?limit=100&page=1": {
			Status: httpStatusOK,
			Headers: map[string]string{
				"Link": "</api/v1/user/repos?page=2&limit=100>; rel=\"next\", </api/v1/user/repos?page=2&limit=100>; rel=\"last\"",
			},
			Body: repoListBody(repoBody(
				1,
				"repo1",
				"Repo one",
				false,
				1,
				2,
				3,
				4,
				5,
				123,
				"2022-10-18T00:00:00Z",
				"2022-10-19T00:00:00Z",
			)),
		},
		"GET /api/v1/user/repos?limit=100&page=2": {
			Status: httpStatusOK,
			Body: repoListBody(repoBody(
				2,
				"repo2",
				"Repo two",
				true,
				0,
				1,
				1,
				2,
				3,
				45,
				"2022-10-20T00:00:00Z",
				"2022-10-21T00:00:00Z",
			)),
		},
		repoAPIPath("repo1", "topics"): {
			Status: httpStatusOK,
			Body:   `{"topics":["game"]}`,
		},
		repoAPIPath("repo1", "languages"): {
			Status: httpStatusOK,
			Body:   `{"Go":120.0}`,
		},
		repoAPIPath("repo1", "pulls"): {
			Status:  httpStatusOK,
			Headers: map[string]string{"X-Total-Count": "2"},
			Body:    `[{"id":3,"state":"closed","title":"Old PR"}]`,
		},
		repoAPIPath("repo1", "issues"): {
			Status:  httpStatusOK,
			Headers: map[string]string{"X-Total-Count": "1"},
			Body:    `[{"id":9,"state":"closed","title":"Old issue"}]`,
		},
		repoAPIPath("repo2", "topics"): {
			Status: httpStatusOK,
			Body:   `{"topics":["dev"]}`,
		},
		repoAPIPath("repo2", "languages"): {
			Status: httpStatusOK,
			Body:   `{"Rust":200.0}`,
		},
		repoAPIPath("repo2", "pulls"): {
			Status:  httpStatusOK,
			Headers: map[string]string{"X-Total-Count": "1"},
			Body:    `[{"id":4,"state":"closed","title":"Old PR2"}]`,
		},
		repoAPIPath("repo2", "issues"): {
			Status:  httpStatusOK,
			Headers: map[string]string{"X-Total-Count": "0"},
			Body:    `[]`,
		},
	})
	defer server.Close()

	service := newServiceForTest(server, t, 4)
	records, err := service.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("len(records) = %d, want 2", len(records))
	}
	if records[0].Name != "repo1" {
		t.Fatalf("records[0].Name = %q, want repo1", records[0].Name)
	}
	if records[1].Name != "repo2" {
		t.Fatalf("records[1].Name = %q, want repo2", records[1].Name)
	}
	if !called("GET /api/v1/user/repos?limit=100&page=1") {
		t.Fatalf("first repo page was not requested")
	}
	if !called("GET /api/v1/user/repos?limit=100&page=2") {
		t.Fatalf("second repo page was not requested")
	}
}

func TestExecute_UsesWorkerCount(t *testing.T) {
	paths := map[string]testsupport.HandlerResponse{
		"GET /api/v1/user/repos": {
			Status: httpStatusOK,
			Body: repoListBody(
				repoBody(
					1,
					"repo1",
					"Repo one",
					false,
					1,
					2,
					3,
					4,
					5,
					123,
					"2022-10-18T00:00:00Z",
					"2022-10-19T00:00:00Z",
				),
				repoBody(
					2,
					"repo2",
					"Repo two",
					false,
					0,
					1,
					1,
					2,
					3,
					45,
					"2022-10-20T00:00:00Z",
					"2022-10-21T00:00:00Z",
				),
			),
		},
		repoAPIPath("repo1", "topics"): {
			Status: httpStatusOK,
			Body:   `{"topics":["game"]}`,
		},
		repoAPIPath("repo1", "languages"): {
			Status: httpStatusOK,
			Body:   `{"Go":120.0,"C++":40.0}`,
		},
		repoAPIPath("repo1", "pulls"): {
			Status: httpStatusOK,
			Headers: map[string]string{
				"X-Total-Count": "10",
			},
			Body: `[{"id":3,"state":"closed","title":"Old PR"}]`,
		},
		repoAPIPath("repo1", "issues"): {
			Status: httpStatusOK,
			Headers: map[string]string{
				"X-Total-Count": "21",
			},
			Body: `[{"id":9,"state":"closed","title":"Old issue"}]`,
		},
		repoAPIPath("repo2", "topics"): {
			Status: httpStatusOK,
			Body:   `{"topics":["dev"]}`,
		},
		repoAPIPath("repo2", "languages"): {
			Status: httpStatusOK,
			Body:   `{"Rust":200.0}`,
		},
		repoAPIPath("repo2", "pulls"): {
			Status: httpStatusOK,
			Headers: map[string]string{
				"X-Total-Count": "7",
			},
			Body: `[{"id":4,"state":"closed","title":"Old PR2"}]`,
		},
		repoAPIPath("repo2", "issues"): {
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
		repoAPIPath("repo1", "pulls"): {
			Status: httpStatusOK,
			Headers: map[string]string{
				"Link": fmt.Sprintf("</api/v1/repos/%s/repo1/pulls?page=2&limit=100&state=all>; rel=\"next\", </api/v1/repos/%s/repo1/pulls?page=2&limit=100&state=all>; rel=\"last\"", testUsername, testUsername),
			},
			Body: `[{"id":1,"state":"open","title":"Add README"},{"id":2,"state":"open","title":"Fix bug"}]`,
		},
	})
	defer server.Close()

	service := newServiceForTest(server, t, 4)
	_, response, err := service.client.ListRepoPullRequests(testUsername, "repo1", forgejo.ListPullRequestsOptions{
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

const (
	httpStatusOK    = 200
	testUsername    = "test-user"
	testRepoBaseURL = "https://example.test"
)

func repoAPIPath(repo, resource string) string {
	return fmt.Sprintf("GET /api/v1/repos/%s/%s/%s", testUsername, repo, resource)
}

func repoAPIPathWithQuery(repo, resource, query string) string {
	return fmt.Sprintf("%s?%s", repoAPIPath(repo, resource), query)
}

func repoListBody(repos ...string) string {
	return "[" + strings.Join(repos, ",") + "]"
}

func repoBody(
	id int,
	name string,
	description string,
	private bool,
	openIssuesCount int,
	openPullsCount int,
	forksCount int,
	starsCount int,
	watchersCount int,
	size int,
	createdAt string,
	updatedAt string,
) string {
	fullName := fmt.Sprintf("%s/%s", testUsername, name)
	httpURL := fmt.Sprintf("%s/%s", testRepoBaseURL, fullName)
	return fmt.Sprintf(
		`{"id":%d,"owner":{"login":%q},"name":%q,"full_name":%q,"description":%q,"private":%t,"html_url":%q,"open_issues_count":%d,"open_pr_counter":%d,"forks_count":%d,"stars_count":%d,"watchers_count":%d,"size":%d,"archived":false,"created_at":%q,"updated_at":%q}`,
		id,
		testUsername,
		name,
		fullName,
		description,
		private,
		httpURL,
		openIssuesCount,
		openPullsCount,
		forksCount,
		starsCount,
		watchersCount,
		size,
		createdAt,
		updatedAt,
	)
}

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
		Username:     testUsername,
		Token:        "token",
		HTTPClient:   server.Client(),
		ReposWorkers: reposWorkers,
	})
}
