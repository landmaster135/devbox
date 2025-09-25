package usecases

import "testing"

func TestBuildLoggingReadCommand_WithFilterParts(t *testing.T) {
	service := NewService()

	cmd, err := service.BuildLoggingReadCommand(LoggingReadParams{
		Severity:       "ERROR",
		ResourceType:   "gce_instance",
		Query:          `textPayload:"Database"`,
		Limit:          15,
		AdditionalArgs: " --format=json ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `gcloud logging read "severity>=ERROR AND resource.type=gce_instance AND textPayload:"Database"" --limit=15 --format=json`
	if cmd != expected {
		t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", expected, cmd)
	}
}

func TestBuildLoggingReadCommand_WithDirectFilter(t *testing.T) {
	service := NewService()

	cmd, err := service.BuildLoggingReadCommand(LoggingReadParams{
		Filter: "resource.type=gce_instance AND severity>=WARNING",
		Limit:  5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `gcloud logging read "resource.type=gce_instance AND severity>=WARNING" --limit=5`
	if cmd != expected {
		t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", expected, cmd)
	}
}

func TestBuildLoggingReadCommand_Errors(t *testing.T) {
	service := NewService()

	t.Run("missing filter", func(t *testing.T) {
		_, err := service.BuildLoggingReadCommand(LoggingReadParams{Limit: 10})
		if err == nil {
			t.Fatal("expected error when filter parts are missing")
		}
	})

	t.Run("invalid limit", func(t *testing.T) {
		_, err := service.BuildLoggingReadCommand(LoggingReadParams{Filter: "severity>=ERROR", Limit: 0})
		if err == nil {
			t.Fatal("expected error for invalid limit")
		}
	})
}

func TestBuildCreateSinkCommand(t *testing.T) {
	service := NewService()

	cmd, err := service.BuildCreateSinkCommand(CreateSinkParams{
		SinkName:       "my-sink",
		Destination:    "storage.googleapis.com/my-bucket",
		LogFilter:      "resource.type=gce_instance",
		AdditionalArgs: " --format=json ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `gcloud logging sinks create my-sink storage.googleapis.com/my-bucket --log-filter="resource.type=gce_instance" --format=json`
	if cmd != expected {
		t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", expected, cmd)
	}
}

func TestBuildCreateSinkCommand_Errors(t *testing.T) {
	service := NewService()

	t.Run("missing sink name", func(t *testing.T) {
		_, err := service.BuildCreateSinkCommand(CreateSinkParams{Destination: "storage.googleapis.com/my-bucket"})
		if err == nil {
			t.Fatal("expected error when sink name is missing")
		}
	})

	t.Run("missing destination", func(t *testing.T) {
		_, err := service.BuildCreateSinkCommand(CreateSinkParams{SinkName: "my-sink"})
		if err == nil {
			t.Fatal("expected error when destination is missing")
		}
	})
}
