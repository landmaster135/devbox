package usecases

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_compute/config"
)

func TestServiceBuildCommand_CreateGCEInstance_Normal(t *testing.T) {
	service := NewService()

	command, err := service.BuildCommand(&cfg.Config{
		Operation:    cfg.OperationCreateGCEInstance,
		InstanceName: "vm-1",
		Zone:         "us-central1-a",
		MachineType:  "e2-medium",
		BootDiskSize: "100GB",
		BootDiskType: "pd-balanced",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud compute instances create 'vm-1' --zone='us-central1-a' --machine-type='e2-medium' --no-address --boot-disk-size='100GB' --boot-disk-type='pd-balanced'"
	if command != expected {
		t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", expected, command)
	}
}

func TestServiceBuildCommand_UnknownOperation(t *testing.T) {
	service := NewService()
	if _, err := service.BuildCommand(&cfg.Config{Operation: "unknown"}); err == nil {
		t.Fatal("expected error for unknown operation")
	}
}

func TestServiceBuildCommand_OtherOperations_Normal(t *testing.T) {
	service := NewService()

	tests := []struct {
		name   string
		config *cfg.Config
	}{
		{
			name: "create router and nat",
			config: &cfg.Config{
				Operation:  cfg.OperationCreateGCERouterAndNAT,
				RouterName: "router-1",
				Region:     "us-central1",
				Network:    "default",
				NATName:    "nat1",
			},
		},
		{
			name: "create iap firewall",
			config: &cfg.Config{
				Operation:    cfg.OperationCreateGCEIAPSSHFirewallRule,
				RuleName:     "allow-ssh-ingress-from-iap",
				Direction:    "INGRESS",
				Action:       "allow",
				Rules:        "tcp:22",
				SourceRanges: "35.235.240.0/20",
			},
		},
		{
			name: "create ingress firewall",
			config: &cfg.Config{
				Operation:    cfg.OperationCreateGCEIngressSSHFirewallRule,
				RuleName:     "allow-ingress-ssh",
				AllowRule:    "tcp:22",
				SourceRanges: "10.0.0.0/8",
			},
		},
		{
			name: "list instances",
			config: &cfg.Config{
				Operation: cfg.OperationListGCloudInstances,
				Format:    "table(name,status)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, err := service.BuildCommand(tt.config)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if strings.TrimSpace(command) == "" {
				t.Fatal("command should not be empty")
			}
		})
	}
}

func TestServiceBuildCreateGCEInstanceCommand_QuoteEscape_Normal(t *testing.T) {
	service := NewService()

	command, err := service.BuildCreateGCEInstanceCommand(CreateGCEInstanceParams{
		InstanceName: "vm'o",
		Zone:         "us-central1-a",
		MachineType:  "e2-medium",
		BootDiskSize: "100GB",
		BootDiskType: "pd-balanced",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(command, "'vm'\"'\"'o'") {
		t.Fatalf("quote escaping is not applied: %s", command)
	}
}

func TestServiceBuildCreateGCEInstanceCommand_ValidationError(t *testing.T) {
	service := NewService()

	if _, err := service.BuildCreateGCEInstanceCommand(CreateGCEInstanceParams{}); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestServiceBuildCreateGCEInstanceCommand_ValidationDetails(t *testing.T) {
	service := NewService()

	tests := []CreateGCEInstanceParams{
		{InstanceName: "vm-1", MachineType: "e2-medium", BootDiskSize: "100GB", BootDiskType: "pd-balanced"},
		{InstanceName: "vm-1", Zone: "us-central1-a", BootDiskSize: "100GB", BootDiskType: "pd-balanced"},
		{InstanceName: "vm-1", Zone: "us-central1-a", MachineType: "e2-medium", BootDiskType: "pd-balanced"},
		{InstanceName: "vm-1", Zone: "us-central1-a", MachineType: "e2-medium", BootDiskSize: "100GB"},
	}

	for i, params := range tests {
		if _, err := service.BuildCreateGCEInstanceCommand(params); err == nil {
			t.Fatalf("expected validation error at case %d", i)
		}
	}
}

func TestServiceBuildCreateGCERouterAndNATCommand_Normal(t *testing.T) {
	service := NewService()

	command, err := service.BuildCreateGCERouterAndNATCommand(CreateGCERouterAndNATParams{
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

func TestServiceBuildCreateGCERouterAndNATCommand_ValidationError(t *testing.T) {
	service := NewService()

	tests := []CreateGCERouterAndNATParams{
		{Region: "us-central1", Network: "default", NATName: "nat1"},
		{RouterName: "router-1", Network: "default", NATName: "nat1"},
		{RouterName: "router-1", Region: "us-central1", NATName: "nat1"},
		{RouterName: "router-1", Region: "us-central1", Network: "default"},
	}

	for i, params := range tests {
		if _, err := service.BuildCreateGCERouterAndNATCommand(params); err == nil {
			t.Fatalf("expected validation error at case %d", i)
		}
	}
}

func TestServiceBuildCreateGCEIAPSSHFirewallRuleCommand_Normal(t *testing.T) {
	service := NewService()

	command, err := service.BuildCreateGCEIAPSSHFirewallRuleCommand(CreateGCEIAPSSHFirewallRuleParams{
		RuleName:     "allow-ssh-ingress-from-iap",
		Direction:    "INGRESS",
		Action:       "allow",
		Rules:        "tcp:22",
		SourceRanges: "35.235.240.0/20",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud compute firewall-rules create 'allow-ssh-ingress-from-iap' --direction='INGRESS' --action='allow' --rules='tcp:22' --source-ranges='35.235.240.0/20'"
	if command != expected {
		t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", expected, command)
	}
}

func TestServiceBuildCreateGCEIAPSSHFirewallRuleCommand_ValidationError(t *testing.T) {
	service := NewService()

	tests := []CreateGCEIAPSSHFirewallRuleParams{
		{Direction: "INGRESS", Action: "allow", Rules: "tcp:22", SourceRanges: "35.235.240.0/20"},
		{RuleName: "r", Action: "allow", Rules: "tcp:22", SourceRanges: "35.235.240.0/20"},
		{RuleName: "r", Direction: "INGRESS", Rules: "tcp:22", SourceRanges: "35.235.240.0/20"},
		{RuleName: "r", Direction: "INGRESS", Action: "allow", SourceRanges: "35.235.240.0/20"},
		{RuleName: "r", Direction: "INGRESS", Action: "allow", Rules: "tcp:22"},
	}

	for i, params := range tests {
		if _, err := service.BuildCreateGCEIAPSSHFirewallRuleCommand(params); err == nil {
			t.Fatalf("expected validation error at case %d", i)
		}
	}
}

func TestServiceBuildCreateGCEIngressSSHFirewallRuleCommand_Normal(t *testing.T) {
	service := NewService()

	command, err := service.BuildCreateGCEIngressSSHFirewallRuleCommand(CreateGCEIngressSSHFirewallRuleParams{
		RuleName:     "allow-ingress-ssh",
		AllowRule:    "tcp:22",
		SourceRanges: "10.0.0.0/8",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud compute firewall-rules create 'allow-ingress-ssh' --allow='tcp:22' --source-ranges='10.0.0.0/8'"
	if command != expected {
		t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", expected, command)
	}
}

func TestServiceBuildCreateGCEIngressSSHFirewallRuleCommand_ValidationError(t *testing.T) {
	service := NewService()

	tests := []CreateGCEIngressSSHFirewallRuleParams{
		{AllowRule: "tcp:22", SourceRanges: "10.0.0.0/8"},
		{RuleName: "allow-ingress-ssh", SourceRanges: "10.0.0.0/8"},
		{RuleName: "allow-ingress-ssh", AllowRule: "tcp:22"},
	}

	for i, params := range tests {
		if _, err := service.BuildCreateGCEIngressSSHFirewallRuleCommand(params); err == nil {
			t.Fatalf("expected validation error at case %d", i)
		}
	}
}

func TestServiceBuildListGCloudInstancesCommand_Normal(t *testing.T) {
	service := NewService()

	command, err := service.BuildListGCloudInstancesCommand(ListGCloudInstancesParams{
		Filter: "zone:us-central1-a",
		Format: "table(name,status)",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud compute instances list --filter='zone:us-central1-a' --format='table(name,status)'"
	if command != expected {
		t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", expected, command)
	}
}

func TestServiceBuildListGCloudInstancesCommand_WithoutFilter_Normal(t *testing.T) {
	service := NewService()

	command, err := service.BuildListGCloudInstancesCommand(ListGCloudInstancesParams{Format: "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud compute instances list --format='json'"
	if command != expected {
		t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", expected, command)
	}
}

func TestServiceBuildListGCloudInstancesCommand_ValidationError(t *testing.T) {
	service := NewService()
	if _, err := service.BuildListGCloudInstancesCommand(ListGCloudInstancesParams{}); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestServicePrintHighlightedCommand(t *testing.T) {
	service := NewService()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	service.PrintHighlightedCommand("gcloud compute instances list")

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "生成された gcloud コマンド") {
		t.Fatalf("header not found: %s", output)
	}
	if !strings.Contains(output, "gcloud compute instances list") {
		t.Fatalf("command not found: %s", output)
	}
}
