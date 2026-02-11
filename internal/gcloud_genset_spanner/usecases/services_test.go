package usecases

import "testing"

func TestBuildInstanceListCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildInstanceListCommand()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "gcloud spanner instances list"
	if cmd != expected {
		t.Fatalf("unexpected command: %s", cmd)
	}
}

func TestBuildInstanceCreateCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildInstanceCreateCommand(InstanceCreateParams{
		InstanceID:     " my-instance ",
		InstanceConfig: "regional-asia-northeast1",
		Description:    "Spanner for Payments",
		Nodes:          3,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "gcloud spanner instances create 'my-instance' \\\n    --config='regional-asia-northeast1' \\\n    --description='Spanner for Payments' \\\n    --nodes=3"
	if cmd != expected {
		t.Fatalf("unexpected command: %s", cmd)
	}
}

func TestBuildInstanceCreateCommand_InvalidNodes(t *testing.T) {
	service := NewService()
	_, err := service.BuildInstanceCreateCommand(InstanceCreateParams{
		InstanceID:     "prod",
		InstanceConfig: "regional-asia-northeast1",
		Description:    "desc",
		Nodes:          0,
	})
	if err == nil {
		t.Fatalf("expected error when nodes is invalid")
	}
}

func TestBuildDatabaseCreateCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildDatabaseCreateCommand(DatabaseCreateParams{
		InstanceID:  "primary",
		DatabaseID:  "orders",
		DDLFilePath: "schema.sql",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "gcloud spanner databases create 'orders' \\\n    --instance='primary' \\\n    --ddl-file='schema.sql'"
	if cmd != expected {
		t.Fatalf("unexpected command: %s", cmd)
	}
}

func TestBuildDatabaseListCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildDatabaseListCommand(DatabaseListParams{InstanceID: " prod "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "gcloud spanner databases list --instance='prod'"
	if cmd != expected {
		t.Fatalf("unexpected command: %s", cmd)
	}
}

func TestBuildDatabaseDescribeCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildDatabaseDescribeCommand(DatabaseDescribeParams{InstanceID: "prod", DatabaseID: "orders"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "gcloud spanner databases describe 'orders' --instance='prod'"
	if cmd != expected {
		t.Fatalf("unexpected command: %s", cmd)
	}
}
