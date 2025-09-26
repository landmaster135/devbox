package usecases

import (
	"strings"
	"testing"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_billing/config"
)

func TestBuildListBudgetsCommand_WithBillingAccount(t *testing.T) {
	service := NewService()

	command, err := service.BuildListBudgetsCommand(ListBudgetsParams{
		BillingAccount: "0000-1111-2222",
		Limit:          20,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud billing budgets list --billing-account='0000-1111-2222' --limit=20"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildListBudgetsCommand_AutoBillingAccount(t *testing.T) {
	service := NewService()

	command, err := service.BuildListBudgetsCommand(ListBudgetsParams{Limit: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(command, "gcloud billing accounts list") {
		t.Fatalf("command should fetch billing account automatically: %s", command)
	}
	if !strings.Contains(command, "--limit=5") {
		t.Fatalf("command should include limit: %s", command)
	}
}

func TestBuildListBudgetsCommand_InvalidLimit(t *testing.T) {
	service := NewService()

	_, err := service.BuildListBudgetsCommand(ListBudgetsParams{Limit: 0})
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestBuildListProjectsCommand_WithBillingAccountAndFilter(t *testing.T) {
	service := NewService()

	command, err := service.BuildListProjectsCommand(ListProjectsParams{
		BillingAccount: "0000-9999-8888",
		Limit:          15,
		Filter:         "project_id:sample",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud billing projects list --billing-account='0000-9999-8888' --limit=15 --filter='project_id:sample'"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildListProjectsCommand_AutoBillingAccount(t *testing.T) {
	service := NewService()

	command, err := service.BuildListProjectsCommand(ListProjectsParams{Limit: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(command, "gcloud billing accounts list") {
		t.Fatalf("command should fetch billing account automatically: %s", command)
	}
	if !strings.Contains(command, "--limit=3") {
		t.Fatalf("command should include limit: %s", command)
	}
}

func TestBuildListProjectsCommand_InvalidLimit(t *testing.T) {
	service := NewService()

	_, err := service.BuildListProjectsCommand(ListProjectsParams{Limit: -1})
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestBuildDescribeProjectCommand_WithProjectID(t *testing.T) {
	service := NewService()

	command, err := service.BuildDescribeProjectCommand(DescribeProjectParams{ProjectID: "sample-project"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud billing projects describe 'sample-project'"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildDescribeProjectCommand_AutoProject(t *testing.T) {
	service := NewService()

	command, err := service.BuildDescribeProjectCommand(DescribeProjectParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(command, "gcloud config get-value project") {
		t.Fatalf("command should fetch project from gcloud config: %s", command)
	}
	if !strings.Contains(command, "gcloud billing projects describe \"$project_id\"") {
		t.Fatalf("command should describe project using fetched id: %s", command)
	}
}

func TestBuildDescribeBudgetCommand_MissingBudgetID(t *testing.T) {
	service := NewService()

	_, err := service.BuildDescribeBudgetCommand(DescribeBudgetParams{})
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestBuildDescribeBudgetCommand_WithBillingAccount(t *testing.T) {
	service := NewService()

	command, err := service.BuildDescribeBudgetCommand(DescribeBudgetParams{
		BudgetID:       "00AA00-123456-ABCDEF",
		BillingAccount: "0000-9999-8888",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud billing budgets describe '00AA00-123456-ABCDEF' --billing-account='0000-9999-8888'"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildDescribeBudgetCommand_AutoBillingAccount(t *testing.T) {
	service := NewService()

	command, err := service.BuildDescribeBudgetCommand(DescribeBudgetParams{BudgetID: "00BB00-654321-XYZ"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(command, "gcloud billing accounts list") {
		t.Fatalf("command should fetch billing account automatically: %s", command)
	}
	if !strings.Contains(command, "gcloud billing budgets describe '00BB00-654321-XYZ'") {
		t.Fatalf("command should include budget id: %s", command)
	}
}

func TestBuildCommand_ListBudgets(t *testing.T) {
	service := NewService()

	command, err := service.BuildCommand(&cfg.Config{
		Operation:      cfg.OperationListBudgets,
		Limit:          5,
		BillingAccount: "0000-9999-8888",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud billing budgets list --billing-account='0000-9999-8888' --limit=5"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildCommand_ListProjects(t *testing.T) {
	service := NewService()

	command, err := service.BuildCommand(&cfg.Config{
		Operation:      cfg.OperationListProjects,
		Limit:          3,
		Filter:         "project_id:test",
		BillingAccount: "",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(command, "gcloud billing accounts list") {
		t.Fatalf("command should fetch billing account automatically: %s", command)
	}
	if !strings.Contains(command, "--filter='project_id:test'") {
		t.Fatalf("command should include filter: %s", command)
	}
}

func TestBuildCommand_DescribeProject(t *testing.T) {
	service := NewService()

	command, err := service.BuildCommand(&cfg.Config{
		Operation: cfg.OperationDescribeProject,
		ProjectID: "demo-project",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud billing projects describe 'demo-project'"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildCommand_ErrorPropagation(t *testing.T) {
	service := NewService()

	_, err := service.BuildCommand(&cfg.Config{
		Operation: cfg.OperationListBudgets,
		Limit:     0,
	})
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestBuildCommand_InvalidOperation(t *testing.T) {
	service := NewService()

	_, err := service.BuildCommand(&cfg.Config{Operation: "unknown"})
	if err == nil {
		t.Fatal("expected error for unsupported operation")
	}
}
