package runner

import (
	"bytes"
	"context"
	"strings"
	"testing"

	cfg "github.com/landmaster135/devbox/internal/memos/config"
	usecases "github.com/landmaster135/devbox/internal/memos/usecases"
)

func TestRun_ListMemosWithFilter_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	called := false

	exitCode := Run([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=10",
		"-page-token=next-token",
		"-state=NORMAL",
		"-order-by=update_time desc",
		`-filter=created_ts > "2023-01-01T13:00:00Z" && visibility == "PUBLIC"`,
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			ListMemosFunc: func(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string, anyContents []string, allContents []string, allTags []string) (*usecases.ListMemosOutput, error) {
				called = true
				if pageSize != 10 {
					t.Fatalf("pageSize = %d, want 10", pageSize)
				}
				if pageToken != "next-token" {
					t.Fatalf("pageToken = %s, want next-token", pageToken)
				}
				if state != "NORMAL" {
					t.Fatalf("state = %s, want NORMAL", state)
				}
				if orderBy != "update_time desc" {
					t.Fatalf("orderBy = %s, want update_time desc", orderBy)
				}
				wantFilter := `created_ts > "2023-01-01T13:00:00Z" && visibility == "PUBLIC"`
				if filter != wantFilter {
					t.Fatalf("filter = %q, want %q", filter, wantFilter)
				}
				if len(anyContents) != 0 {
					t.Fatalf("anyContents = %v, want empty", anyContents)
				}
				if len(allContents) != 0 || len(allTags) != 0 {
					t.Fatalf("allContents/allTags = %v/%v, want empty", allContents, allTags)
				}
				return &usecases.ListMemosOutput{
					Memos: []usecases.Memo{
						{
							Name: "memos/1",
							Attachments: []usecases.Attachment{
								{Name: "attachments/1"},
							},
						},
					},
				}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("ListMemosFunc was not called")
	}
	if !strings.Contains(stdout.String(), "\"name\": \"memos/1\"") {
		t.Fatalf("stdout = %s, want memo json", stdout.String())
	}
	if !strings.Contains(stdout.String(), "\"attachments\": [") {
		t.Fatalf("stdout = %s, want attachments json", stdout.String())
	}
	if !strings.Contains(stdout.String(), "\"name\": \"attachments/1\"") {
		t.Fatalf("stdout = %s, want attachment name", stdout.String())
	}
}

func TestRun_ListMemosWithAnyContents_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	called := false

	exitCode := Run([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-any-contents= meeting, study ,,",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			ListMemosFunc: func(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string, anyContents []string, allContents []string, allTags []string) (*usecases.ListMemosOutput, error) {
				called = true
				if pageSize != 20 {
					t.Fatalf("pageSize = %d, want 20", pageSize)
				}
				if pageToken != "" {
					t.Fatalf("pageToken = %q, want empty", pageToken)
				}
				if state != "" || orderBy != "" || filter != "" {
					t.Fatalf("state/orderBy/filter = %q/%q/%q, want empty", state, orderBy, filter)
				}
				if len(anyContents) != 2 || anyContents[0] != "meeting" || anyContents[1] != "study" {
					t.Fatalf("anyContents = %v, want [meeting study]", anyContents)
				}
				if len(allContents) != 0 || len(allTags) != 0 {
					t.Fatalf("allContents/allTags = %v/%v, want empty", allContents, allTags)
				}
				return &usecases.ListMemosOutput{
					Memos: []usecases.Memo{
						{Name: "memos/1"},
					},
				}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("ListMemosFunc was not called")
	}
	if !strings.Contains(stdout.String(), "\"name\": \"memos/1\"") {
		t.Fatalf("stdout = %s, want memo json", stdout.String())
	}
}

func TestRun_ListMemosWithAnyTags_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	called := false

	exitCode := Run([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-any-tags= health, book ,,",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			ListMemosFunc: func(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string, anyContents []string, allContents []string, allTags []string) (*usecases.ListMemosOutput, error) {
				called = true
				if filter != "tag in ['health','book']" {
					t.Fatalf("filter = %q, want %q", filter, "tag in ['health','book']")
				}
				return &usecases.ListMemosOutput{
					Memos: []usecases.Memo{
						{Name: "memos/1"},
					},
				}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("ListMemosFunc was not called")
	}
}

func TestRun_ListMemosWithAnyTagsAndFilter_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	called := false

	exitCode := Run([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		`-filter=visibility == "PUBLIC"`,
		"-any-tags=health,book",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			ListMemosFunc: func(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string, anyContents []string, allContents []string, allTags []string) (*usecases.ListMemosOutput, error) {
				called = true
				want := `(visibility == "PUBLIC") && (tag in ['health','book'])`
				if filter != want {
					t.Fatalf("filter = %q, want %q", filter, want)
				}
				return &usecases.ListMemosOutput{
					Memos: []usecases.Memo{
						{Name: "memos/2"},
					},
				}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("ListMemosFunc was not called")
	}
}

func TestRun_ListMemosWithAllContents_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	called := false

	exitCode := Run([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-all-contents= meeting, study ,,",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			ListMemosFunc: func(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string, anyContents []string, allContents []string, allTags []string) (*usecases.ListMemosOutput, error) {
				called = true
				if pageSize != 20 {
					t.Fatalf("pageSize = %d, want 20", pageSize)
				}
				if len(anyContents) != 0 {
					t.Fatalf("anyContents = %v, want empty", anyContents)
				}
				if len(allContents) != 2 || allContents[0] != "meeting" || allContents[1] != "study" {
					t.Fatalf("allContents = %v, want [meeting study]", allContents)
				}
				if len(allTags) != 0 {
					t.Fatalf("allTags = %v, want empty", allTags)
				}
				return &usecases.ListMemosOutput{
					Memos: []usecases.Memo{
						{Name: "memos/2"},
					},
				}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("ListMemosFunc was not called")
	}
	if !strings.Contains(stdout.String(), "\"name\": \"memos/2\"") {
		t.Fatalf("stdout = %s, want memo json", stdout.String())
	}
}

func TestRun_ListMemosWithAllTags_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	called := false

	exitCode := Run([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-all-tags= health, book ,,",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			ListMemosFunc: func(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string, anyContents []string, allContents []string, allTags []string) (*usecases.ListMemosOutput, error) {
				called = true
				if filter != "" {
					t.Fatalf("filter = %q, want empty", filter)
				}
				if len(anyContents) != 0 || len(allContents) != 0 {
					t.Fatalf("anyContents/allContents = %v/%v, want empty", anyContents, allContents)
				}
				if len(allTags) != 2 || allTags[0] != "health" || allTags[1] != "book" {
					t.Fatalf("allTags = %v, want [health book]", allTags)
				}
				return &usecases.ListMemosOutput{
					Memos: []usecases.Memo{
						{Name: "memos/3"},
					},
				}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("ListMemosFunc was not called")
	}
	if !strings.Contains(stdout.String(), "\"name\": \"memos/3\"") {
		t.Fatalf("stdout = %s, want memo json", stdout.String())
	}
}

func TestRun_ListMemosWithAllTagsAndFilter_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	called := false

	exitCode := Run([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		`-filter=visibility == "PUBLIC"`,
		"-all-tags=health,book",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			ListMemosFunc: func(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string, anyContents []string, allContents []string, allTags []string) (*usecases.ListMemosOutput, error) {
				called = true
				if filter != `visibility == "PUBLIC"` {
					t.Fatalf("filter = %q, want base filter", filter)
				}
				if len(anyContents) != 0 || len(allContents) != 0 {
					t.Fatalf("anyContents/allContents = %v/%v, want empty", anyContents, allContents)
				}
				if len(allTags) != 2 || allTags[0] != "health" || allTags[1] != "book" {
					t.Fatalf("allTags = %v, want [health book]", allTags)
				}
				return &usecases.ListMemosOutput{
					Memos: []usecases.Memo{
						{Name: "memos/4"},
					},
				}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("ListMemosFunc was not called")
	}
}

func TestRun_ListMemosWithExcludedTags_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	callCount := 0

	exitCode := Run([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-excluded-tags=health,book",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			ListMemosFunc: func(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string, anyContents []string, allContents []string, allTags []string) (*usecases.ListMemosOutput, error) {
				callCount++

				if pageSize != 20 {
					t.Fatalf("pageSize = %d, want 20", pageSize)
				}
				if len(anyContents) != 0 || len(allContents) != 0 || len(allTags) != 0 {
					t.Fatalf("unexpected list options: anyContents=%v allContents=%v allTags=%v", anyContents, allContents, allTags)
				}

				switch callCount {
				case 1:
					if filter != "" {
						t.Fatalf("filter(1st) = %q, want empty", filter)
					}
					return &usecases.ListMemosOutput{
						Memos: []usecases.Memo{
							{Name: "memos/1"},
							{Name: "memos/2"},
							{Name: "memos/3"},
						},
						TotalSize: 3,
					}, nil
				case 2:
					if filter != "tag in ['health']" {
						t.Fatalf("filter(2nd) = %q, want %q", filter, "tag in ['health']")
					}
					return &usecases.ListMemosOutput{
						Memos: []usecases.Memo{
							{Name: "memos/2"},
						},
						TotalSize: 1,
					}, nil
				case 3:
					if filter != "tag in ['book']" {
						t.Fatalf("filter(3rd) = %q, want %q", filter, "tag in ['book']")
					}
					return &usecases.ListMemosOutput{
						Memos: []usecases.Memo{
							{Name: "memos/3"},
						},
						TotalSize: 1,
					}, nil
				default:
					t.Fatalf("unexpected ListMemos call count: %d", callCount)
					return nil, nil
				}
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if callCount != 3 {
		t.Fatalf("callCount = %d, want 3", callCount)
	}
	if !strings.Contains(stdout.String(), "\"name\": \"memos/1\"") {
		t.Fatalf("stdout = %s, want remaining memo", stdout.String())
	}
	if strings.Contains(stdout.String(), "\"name\": \"memos/2\"") || strings.Contains(stdout.String(), "\"name\": \"memos/3\"") {
		t.Fatalf("stdout = %s, want excluded memos removed", stdout.String())
	}
	if !strings.Contains(stdout.String(), "\"totalSize\": 1") {
		t.Fatalf("stdout = %s, want totalSize=1", stdout.String())
	}
}

func TestRun_ListMemosWithAnyTagsAndExcludedTags_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	callCount := 0

	exitCode := Run([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		`-filter=visibility == "PUBLIC"`,
		"-any-tags=health,book",
		"-excluded-tags=archive",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			ListMemosFunc: func(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string, anyContents []string, allContents []string, allTags []string) (*usecases.ListMemosOutput, error) {
				callCount++
				switch callCount {
				case 1:
					want := `(visibility == "PUBLIC") && (tag in ['health','book'])`
					if filter != want {
						t.Fatalf("filter(1st) = %q, want %q", filter, want)
					}
					return &usecases.ListMemosOutput{
						Memos: []usecases.Memo{
							{Name: "memos/1"},
							{Name: "memos/2"},
						},
						TotalSize: 2,
					}, nil
				case 2:
					if filter != "tag in ['archive']" {
						t.Fatalf("filter(2nd) = %q, want %q", filter, "tag in ['archive']")
					}
					return &usecases.ListMemosOutput{
						Memos: []usecases.Memo{
							{Name: "memos/2"},
						},
						TotalSize: 1,
					}, nil
				default:
					t.Fatalf("unexpected ListMemos call count: %d", callCount)
					return nil, nil
				}
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if callCount != 2 {
		t.Fatalf("callCount = %d, want 2", callCount)
	}
	if !strings.Contains(stdout.String(), "\"name\": \"memos/1\"") {
		t.Fatalf("stdout = %s, want memos/1", stdout.String())
	}
	if strings.Contains(stdout.String(), "\"name\": \"memos/2\"") {
		t.Fatalf("stdout = %s, want memos/2 excluded", stdout.String())
	}
}
