package usecases

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestBuildCreatePubSubJobCommand(t *testing.T) {
	t.Parallel()
	service := NewService()
	cmd, err := service.BuildCreatePubSubJobCommand(CreatePubSubJobParams{
		JobName:     "daily-job",
		ProjectID:   "sample-project",
		PubsubTopic: "projects/sample/topics/task",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantFragments := []string{
		"gcloud scheduler jobs create pubsub 'daily-job'",
		"--project='sample-project'",
		"--topic='projects/sample/topics/task'",
		"--schedule='0 4 * * 0-6'",
		"--time-zone='Asia/Tokyo'",
	}
	for _, frag := range wantFragments {
		if !strings.Contains(cmd, frag) {
			t.Fatalf("expected command to contain %q, got:\n%s", frag, cmd)
		}
	}
}

func TestBuildCreatePubSubJobCommand_WithMessageBody(t *testing.T) {
	t.Parallel()
	service := NewService()
	cmd, err := service.BuildCreatePubSubJobCommand(CreatePubSubJobParams{
		JobName:     "custom",
		ProjectID:   "proj",
		PubsubTopic: "topic",
		MessageBody: "{\n  \"key\": \"value\"\n}",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(cmd, "\\n") {
		t.Fatalf("expected newlines to be removed from message body: %s", cmd)
	}
}

func TestBuildCreatePubSubJobCommand_ErrorWhenMissing(t *testing.T) {
	t.Parallel()
	service := NewService()
	if _, err := service.BuildCreatePubSubJobCommand(CreatePubSubJobParams{}); err == nil {
		t.Fatal("expected error when required parameters missing")
	}
}

func TestBuildCreateHTTPJobCommand(t *testing.T) {
	t.Parallel()
	service := NewService()
	cmd, err := service.BuildCreateHTTPJobCommand(CreateHTTPJobParams{
		JobName:                 "http-job",
		ProjectID:               "proj",
		HTTPMethod:              "post",
		ServiceURL:              "https://example.com",
		Headers:                 "Content-Type=application/json",
		MessageBody:             "{\"ok\":true}",
		OIDCServiceAccountEmail: "svc@example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantFragments := []string{
		"gcloud scheduler jobs create http 'http-job'",
		"--http-method='POST'",
		"--uri='https://example.com'",
		"--headers='Content-Type=application/json'",
		"--message-body='{\"ok\":true}'",
		"--oidc-service-account-email='svc@example.com'",
	}
	for _, frag := range wantFragments {
		if !strings.Contains(cmd, frag) {
			t.Fatalf("expected command to contain %q, got:\n%s", frag, cmd)
		}
	}
}

func TestBuildCreateHTTPJobCommand_ErrorWhenMissing(t *testing.T) {
	t.Parallel()
	service := NewService()
	if _, err := service.BuildCreateHTTPJobCommand(CreateHTTPJobParams{}); err == nil {
		t.Fatal("expected error when required parameters missing")
	}
}

func TestBuildCreateCloudSQLJobCommand(t *testing.T) {
	t.Parallel()
	service := NewService()
	cmd, err := service.BuildCreateCloudSQLJobCommand(CreateCloudSQLJobParams{
		JobName:           "sql-job",
		ProjectID:         "proj",
		PubsubTopic:       "topic",
		DBInstanceID:      "instance",
		DiscordWebhookURL: "https://discord",
		CloudSQLIconURL:   "https://icon",
		Action:            "restart",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(cmd, "\"Action\":\"restart\"") {
		t.Fatalf("expected action to be included in message body: %s", cmd)
	}
	if !strings.Contains(cmd, "\"DiscordWebhookUrl\":\"https://discord\"") {
		t.Fatalf("expected webhook url to be included: %s", cmd)
	}
}

func TestBuildCreateCloudSQLJobCommand_ErrorWhenMissing(t *testing.T) {
	t.Parallel()
	service := NewService()
	if _, err := service.BuildCreateCloudSQLJobCommand(CreateCloudSQLJobParams{}); err == nil {
		t.Fatal("expected error when parameters missing")
	}
}

func TestBuildCreateStartCloudSQLJobCommand_Defaults(t *testing.T) {
	t.Parallel()
	service := NewService()
	cmd, err := service.BuildCreateStartCloudSQLJobCommand(CreateStartStopCloudSQLJobParams{
		ProjectID:    "proj",
		PubsubTopic:  "topic",
		DBInstanceID: "db01",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(cmd, "gcloud scheduler jobs create pubsub 'start-db01-instance'") {
		t.Fatalf("expected auto generated job name, got: %s", cmd)
	}
	if !strings.Contains(cmd, "--schedule='0 4 * * 0-6'") {
		t.Fatalf("expected default start schedule, got: %s", cmd)
	}
	if !strings.Contains(cmd, "\"Action\":\"start\"") {
		t.Fatalf("expected action start in payload, got: %s", cmd)
	}
}

func TestBuildCreateStopCloudSQLJobCommand_Defaults(t *testing.T) {
	t.Parallel()
	service := NewService()
	cmd, err := service.BuildCreateStopCloudSQLJobCommand(CreateStartStopCloudSQLJobParams{
		ProjectID:    "proj",
		PubsubTopic:  "topic",
		DBInstanceID: "db01",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(cmd, "gcloud scheduler jobs create pubsub 'stop-db01-instance'") {
		t.Fatalf("expected auto generated job name, got: %s", cmd)
	}
	if !strings.Contains(cmd, "--schedule='0 7 * * 0-6'") {
		t.Fatalf("expected default stop schedule, got: %s", cmd)
	}
	if !strings.Contains(cmd, "\"Action\":\"stop\"") {
		t.Fatalf("expected action stop in payload, got: %s", cmd)
	}
}

func TestBuildListJobsCommand(t *testing.T) {
	t.Parallel()
	service := NewService()
	cmd, err := service.BuildListJobsCommand(ListJobsParams{Location: "asia-northeast1", Limit: "5"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(cmd, "--location='asia-northeast1'") {
		t.Fatalf("expected location flag, got: %s", cmd)
	}
	if !strings.Contains(cmd, "--limit='5'") {
		t.Fatalf("expected limit flag, got: %s", cmd)
	}
}

func TestBuildUpdateHTTPJobCommand(t *testing.T) {
	t.Parallel()
	service := NewService()
	cmd, err := service.BuildUpdateHTTPJobCommand(UpdateHTTPJobParams{
		JobName:     "http-job",
		Schedule:    "*/5 * * * *",
		MessageBody: "{\n  \"a\":1\n}",
		Headers:     "Content-Type=application/json",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(cmd, "\\n") {
		t.Fatalf("expected newlines removed in message body: %s", cmd)
	}
	expected := []string{
		"gcloud scheduler jobs update http 'http-job'",
		"--schedule='*/5 * * * *'",
		"--message-body='{  \"a\":1}'",
		"--headers='Content-Type=application/json'",
	}
	for _, frag := range expected {
		if !strings.Contains(cmd, frag) {
			t.Fatalf("expected fragment %q in command: %s", frag, cmd)
		}
	}
}

func TestBuildUpdateHTTPJobCommand_Error(t *testing.T) {
	t.Parallel()
	service := NewService()
	if _, err := service.BuildUpdateHTTPJobCommand(UpdateHTTPJobParams{}); err == nil {
		t.Fatal("expected error when job name missing")
	}
}

func TestBuildUpdatePubSubJobCommand(t *testing.T) {
	t.Parallel()
	service := NewService()
	cmd, err := service.BuildUpdatePubSubJobCommand(UpdatePubSubJobParams{
		JobName:  "pubsub-job",
		Schedule: "*/10 * * * *",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(cmd, "gcloud scheduler jobs update pubsub 'pubsub-job'") {
		t.Fatalf("unexpected command: %s", cmd)
	}
}

func TestBuildUpdatePubSubJobCommand_Error(t *testing.T) {
	t.Parallel()
	service := NewService()
	if _, err := service.BuildUpdatePubSubJobCommand(UpdatePubSubJobParams{}); err == nil {
		t.Fatal("expected error when job name missing")
	}
}

func TestJobControlCommands(t *testing.T) {
	t.Parallel()
	service := NewService()

	tests := []struct {
		name    string
		builder func(JobControlParams) (string, error)
		action  string
		extra   string
	}{
		{"pause", service.BuildPauseJobCommand, "pause", ""},
		{"resume", service.BuildResumeJobCommand, "resume", ""},
		{"delete", service.BuildDeleteJobCommand, "delete", "--quiet"},
		{"run", service.BuildRunJobCommand, "run", ""},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cmd, err := tt.builder(JobControlParams{JobName: "job", Location: "asia"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(cmd, "gcloud scheduler jobs "+tt.action+" 'job'") {
				t.Fatalf("unexpected command: %s", cmd)
			}
			if !strings.Contains(cmd, "--location='asia'") {
				t.Fatalf("expected location flag, got: %s", cmd)
			}
			if tt.extra != "" && !strings.Contains(cmd, tt.extra) {
				t.Fatalf("expected extra fragment %q, got: %s", tt.extra, cmd)
			}
		})
	}

	if _, err := service.BuildPauseJobCommand(JobControlParams{}); err == nil {
		t.Fatal("expected error when job name missing")
	}
}

func TestPrintHighlightedCommand(t *testing.T) {
	t.Parallel()
	service := NewService()
	output := captureStdout(func() {
		service.PrintHighlightedCommand("gcloud scheduler jobs list")
	})

	if !strings.Contains(output, "生成された gcloud コマンド") {
		t.Fatalf("expected header in output: %s", output)
	}
	if !strings.Contains(output, "gcloud scheduler jobs list") {
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
