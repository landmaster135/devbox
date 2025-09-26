package usecases

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestBuildListDashboardsCommand(t *testing.T) {
	service := NewService()

	cmd, err := service.BuildListDashboardsCommand(ListDashboardsParams{
		Project:  "my-project",
		Filter:   "displayName:test",
		Format:   "json",
		PageSize: 50,
		SortBy:   "name",
		Limit:    100,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud monitoring dashboards list --project=my-project --filter=displayName:test --format=\"json\" --page-size=50 --sort-by=name --limit=100"
	if cmd != expected {
		t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", expected, cmd)
	}
}

func TestBuildListDashboardsCommand_Errors(t *testing.T) {
	service := NewService()

	t.Run("zero limit", func(t *testing.T) {
		_, err := service.BuildListDashboardsCommand(ListDashboardsParams{Limit: 0, PageSize: -1})
		if err == nil {
			t.Fatal("expected error for zero limit")
		}
	})

	t.Run("zero page size", func(t *testing.T) {
		_, err := service.BuildListDashboardsCommand(ListDashboardsParams{PageSize: 0, Limit: -1})
		if err == nil {
			t.Fatal("expected error for zero page size")
		}
	})
}

func TestBuildDescribeDashboardCommand(t *testing.T) {
	service := NewService()

	cmd, err := service.BuildDescribeDashboardCommand(DescribeDashboardParams{
		DashboardID: "dashboards/123",
		Project:     "my-project",
		Format:      "yaml",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud monitoring dashboards describe dashboards/123 --project=my-project --format=\"yaml\""
	if cmd != expected {
		t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", expected, cmd)
	}
}

func TestBuildDescribeDashboardCommand_Errors(t *testing.T) {
	service := NewService()

	if _, err := service.BuildDescribeDashboardCommand(DescribeDashboardParams{}); err == nil {
		t.Fatal("expected error when dashboard id is missing")
	}
}

func TestBuildListSnoozesCommand(t *testing.T) {
	service := NewService()

	cmd, err := service.BuildListSnoozesCommand(ListSnoozesParams{
		Project:    "my-project",
		Filter:     "displayName:maintenance",
		Format:     "json",
		PageSize:   20,
		SortBy:     "displayName",
		Limit:      5,
		IncludeURI: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud monitoring snoozes list --project=my-project --filter=displayName:maintenance --format=\"json\" --page-size=20 --sort-by=displayName --limit=5 --uri"
	if cmd != expected {
		t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", expected, cmd)
	}
}

func TestBuildListUptimeConfigsCommand(t *testing.T) {
	service := NewService()

	cmd, err := service.BuildListUptimeConfigsCommand(ListUptimeConfigsParams{
		Project:    "my-project",
		Filter:     "displayName:api",
		Format:     "yaml",
		PageSize:   30,
		Limit:      40,
		SortBy:     "displayName",
		IncludeURI: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud monitoring uptime list-configs --project=my-project --filter=displayName:api --format=\"yaml\" --page-size=30 --sort-by=displayName --limit=40"
	if cmd != expected {
		t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", expected, cmd)
	}
}

func TestBuildListSnoozesCommand_Errors(t *testing.T) {
	service := NewService()

	if _, err := service.BuildListSnoozesCommand(ListSnoozesParams{Limit: 0, PageSize: -1}); err == nil {
		t.Fatal("expected error for zero limit")
	}
}

func TestBuildListUptimeConfigsCommand_Errors(t *testing.T) {
	service := NewService()

	if _, err := service.BuildListUptimeConfigsCommand(ListUptimeConfigsParams{PageSize: 0, Limit: -1}); err == nil {
		t.Fatal("expected error for zero page size")
	}
}

func TestPrintHighlightedCommand(t *testing.T) {
	service := NewService()

	output := captureStdout(func() {
		service.PrintHighlightedCommand("gcloud monitoring dashboards list")
	})

	if !strings.Contains(output, "生成された gcloud コマンド") {
		t.Fatalf("expected header in output: %s", output)
	}
	if !strings.Contains(output, "gcloud monitoring dashboards list") {
		t.Fatalf("expected command in output: %s", output)
	}
}

func captureStdout(fn func()) string {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		panic(err)
	}

	return buf.String()
}
