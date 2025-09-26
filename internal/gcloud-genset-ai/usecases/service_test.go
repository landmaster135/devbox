package usecases

import "testing"

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
