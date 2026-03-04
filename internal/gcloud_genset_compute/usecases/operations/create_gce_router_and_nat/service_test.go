package creategcerouterandnat

import (
	"strings"
	"testing"
)

func TestServiceBuild_Normal(t *testing.T) {
	t.Parallel()

	service := NewService()
	command, err := service.Build(Params{
		RouterName: "router-1",
		Region:     "us-central1",
		Network:    "default",
		NATName:    "nat1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedParts := []string{
		"gcloud compute routers create 'router-1' --region='us-central1' --network='default'",
		"gcloud compute routers nats create 'nat1' --router='router-1' --region='us-central1' --auto-allocate-nat-external-ips --nat-all-subnet-ip-ranges",
	}
	for _, expected := range expectedParts {
		if !strings.Contains(command, expected) {
			t.Fatalf("command should contain %q: %s", expected, command)
		}
	}
	if !strings.Contains(command, "&&") {
		t.Fatalf("command should include conditional chain: %s", command)
	}
}

func TestServiceBuild_ValidationError(t *testing.T) {
	t.Parallel()

	service := NewService()

	tests := []Params{
		{Region: "us-central1", Network: "default", NATName: "nat1"},
		{RouterName: "router-1", Network: "default", NATName: "nat1"},
		{RouterName: "router-1", Region: "us-central1", NATName: "nat1"},
		{RouterName: "router-1", Region: "us-central1", Network: "default"},
	}

	for i, params := range tests {
		if _, err := service.Build(params); err == nil {
			t.Fatalf("expected validation error at case %d", i)
		}
	}
}
