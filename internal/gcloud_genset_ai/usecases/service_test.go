package usecases

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestBuildUndeployProcessorVersionCommand(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := NewService()
		command, err := service.BuildUndeployProcessorVersionCommand(UndeployProcessorVersionParams{
			Region:        " us-central1 ",
			ProjectNumber: " 1234567890 ",
			ProcessorID:   " processor-abc ",
			VersionID:     " version-001 ",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := "curl -s -X POST -H \"Authorization: Bearer $(gcloud auth print-access-token)\" -H \"Content-Type: application/json\" \"https://us-central1-documentai.googleapis.com/v1beta3/1234567890/locations/us-central1/processors/processor-abc/processorVersions/version-001:undeploy\""
		if command != expected {
			t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", expected, command)
		}
	})

	for name, params := range map[string]UndeployProcessorVersionParams{
		"missing region":    {ProjectNumber: "123", ProcessorID: "proc", VersionID: "ver"},
		"missing project":   {Region: "us", ProcessorID: "proc", VersionID: "ver"},
		"missing processor": {Region: "us", ProjectNumber: "123", VersionID: "ver"},
		"missing version":   {Region: "us", ProjectNumber: "123", ProcessorID: "proc"},
	} {
		params := params
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			service := NewService()
			if _, err := service.BuildUndeployProcessorVersionCommand(params); err == nil {
				t.Fatalf("expected error for %s", name)
			}
		})
	}
}

func TestPrintHighlightedCommand(t *testing.T) {
	service := NewService()

	output := captureStdout(func() {
		service.PrintHighlightedCommand("gcloud ai example")
	})

	if !strings.Contains(output, "生成された gcloud コマンド") {
		t.Fatalf("expected header in output: %s", output)
	}
	if !strings.Contains(output, "gcloud ai example") {
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
