package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	config "github.com/landmaster135/devbox/internal/forgejo/config"
	usecases "github.com/landmaster135/devbox/internal/forgejo/usecases"
)

func main() {
	os.Exit(run())
}

func run() int {
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		return 1
	}

	if cfg.Help {
		config.PrintUsage()
		return 0
	}

	service, err := usecases.NewService(usecases.ServiceOptions{
		Host:     cfg.Host,
		Username: cfg.Username,
		Token:    cfg.Token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		ReposWorkers: cfg.ReposWorkers,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		return 1
	}

	switch cfg.Operation {
	case "repo list":
		records, err := service.ListRepos()
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			return 1
		}
		return outputRecords(records, cfg.JSON)
	case "issue list":
		records, err := service.ListIssues()
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			return 1
		}
		return outputRecords(records, cfg.JSON)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作タイプです: %s\n", cfg.Operation)
		config.PrintUsage()
		return 1
	}
}

func outputRecords(records any, asJSON bool) int {
	if asJSON {
		data, err := json.MarshalIndent(records, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: JSON の整形に失敗しました: %v\n", err)
			return 1
		}
		fmt.Println(string(data))
		return 0
	}

	switch typed := records.(type) {
	case []usecases.RepoRecord:
		return printRepos(typed)
	case []usecases.IssueRecord:
		return printIssues(typed)
	default:
		data, err := json.Marshal(records)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: 出力に失敗しました: %v\n", err)
			return 1
		}
		fmt.Println(string(data))
		return 0
	}
}

func printRepos(records []usecases.RepoRecord) int {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "name\tdescription\tis_private\thttp_url\topen_issues_count\tclosed_issues_count\topen_pulls_count\tclosed_pulls_count\tforks_count\tstargazers_count\tsubscribers_count\tlanguage\tlanguages\tsize\trepo_created_at\trepo_updated_at\tis_archived\ttags")
	for _, record := range records {
		languageJSON, err := json.Marshal(record.Languages)
		if err != nil {
			_ = w.Flush()
			fmt.Fprintf(os.Stderr, "エラー: 言語情報の整形に失敗しました: %v\n", err)
			return 1
		}
		_, _ = fmt.Fprintf(w, "%s\t%s\t%t\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%s\t%s\t%d\t%s\t%s\t%t\t%s\n",
			record.Name,
			record.Description,
			record.IsPrivate,
			record.HTTPURL,
			record.OpenIssuesCount,
			record.ClosedIssuesCount,
			record.OpenPullsCount,
			record.ClosedPullsCount,
			record.ForksCount,
			record.StargazersCount,
			record.SubscribersCount,
			record.Language,
			string(languageJSON),
			record.Size,
			record.RepoCreatedAt,
			record.RepoUpdatedAt,
			record.IsArchived,
			record.Tags,
		)
	}

	_ = w.Flush()
	return 0
}

func printIssues(records []usecases.IssueRecord) int {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "repo_full_name\tnumber\ttitle\tstate\thtml_url\tauthor\tassignees\tlabels\tcomments\tis_locked\tcreated_at\tupdated_at\tclosed_at")
	for _, record := range records {
		_, _ = fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%d\t%t\t%s\t%s\t%s\n",
			record.RepoFullName,
			record.Number,
			record.Title,
			record.State,
			record.HTMLURL,
			record.Author,
			strings.Join(record.Assignees, ","),
			strings.Join(record.Labels, ","),
			record.Comments,
			record.IsLocked,
			record.CreatedAt,
			record.UpdatedAt,
			record.ClosedAt,
		)
	}

	_ = w.Flush()
	return 0
}
