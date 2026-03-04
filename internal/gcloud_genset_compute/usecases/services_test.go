package usecases

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_compute/config"
	creategceiapsshfirewallrule "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_iap_ssh_firewall_rule"
	creategceingresssshfirewallrule "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_ingress_ssh_firewall_rule"
	creategceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_instance"
	creategcerouterandnat "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_router_and_nat"
	listgcloudinstances "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/list_gcloud_instances"
)

type createGCEInstanceOperationStub struct {
	called bool
	got    creategceinstance.Params
	result string
	err    error
}

func (s *createGCEInstanceOperationStub) Build(params creategceinstance.Params) (string, error) {
	s.called = true
	s.got = params
	return s.result, s.err
}

type createGCERouterAndNATOperationStub struct {
	called bool
	got    creategcerouterandnat.Params
	result string
	err    error
}

func (s *createGCERouterAndNATOperationStub) Build(params creategcerouterandnat.Params) (string, error) {
	s.called = true
	s.got = params
	return s.result, s.err
}

type createGCEIAPSSHFirewallRuleOperationStub struct {
	called bool
	got    creategceiapsshfirewallrule.Params
	result string
	err    error
}

func (s *createGCEIAPSSHFirewallRuleOperationStub) Build(params creategceiapsshfirewallrule.Params) (string, error) {
	s.called = true
	s.got = params
	return s.result, s.err
}

type createGCEIngressSSHFirewallRuleOperationStub struct {
	called bool
	got    creategceingresssshfirewallrule.Params
	result string
	err    error
}

func (s *createGCEIngressSSHFirewallRuleOperationStub) Build(params creategceingresssshfirewallrule.Params) (string, error) {
	s.called = true
	s.got = params
	return s.result, s.err
}

type listGCloudInstancesOperationStub struct {
	called bool
	got    listgcloudinstances.Params
	result string
	err    error
}

func (s *listGCloudInstancesOperationStub) Build(params listgcloudinstances.Params) (string, error) {
	s.called = true
	s.got = params
	return s.result, s.err
}

func TestServiceBuildCommand_DelegatesByOperation(t *testing.T) {
	instanceOp := &createGCEInstanceOperationStub{result: "instance-command"}
	routerNATOp := &createGCERouterAndNATOperationStub{result: "router-nat-command"}
	iapSSHOp := &createGCEIAPSSHFirewallRuleOperationStub{result: "iap-ssh-command"}
	ingressSSHOp := &createGCEIngressSSHFirewallRuleOperationStub{result: "ingress-ssh-command"}
	listOp := &listGCloudInstancesOperationStub{result: "list-command"}

	service := newServiceWithOperations(instanceOp, routerNATOp, iapSSHOp, ingressSSHOp, listOp)

	tests := []struct {
		name     string
		config   *cfg.Config
		expected string
	}{
		{
			name: "create-gce-instance",
			config: &cfg.Config{
				Operation:    cfg.OperationCreateGCEInstance,
				InstanceName: "vm-1",
				Zone:         "us-central1-a",
				MachineType:  "e2-medium",
				BootDiskSize: "100GB",
				BootDiskType: "pd-balanced",
			},
			expected: "instance-command",
		},
		{
			name: "create-gce-router-and-nat",
			config: &cfg.Config{
				Operation:  cfg.OperationCreateGCERouterAndNAT,
				RouterName: "router-1",
				Region:     "us-central1",
				Network:    "default",
				NATName:    "nat1",
			},
			expected: "router-nat-command",
		},
		{
			name: "create-gce-iap-ssh-firewall-rule",
			config: &cfg.Config{
				Operation:    cfg.OperationCreateGCEIAPSSHFirewallRule,
				RuleName:     "allow-ssh-ingress-from-iap",
				Direction:    "INGRESS",
				Action:       "allow",
				Rules:        "tcp:22",
				SourceRanges: "35.235.240.0/20",
			},
			expected: "iap-ssh-command",
		},
		{
			name: "create-gce-ingress-ssh-firewall-rule",
			config: &cfg.Config{
				Operation:    cfg.OperationCreateGCEIngressSSHFirewallRule,
				RuleName:     "allow-ingress-ssh",
				AllowRule:    "tcp:22",
				SourceRanges: "10.0.0.0/8",
			},
			expected: "ingress-ssh-command",
		},
		{
			name: "list-gcloud-instances",
			config: &cfg.Config{
				Operation: cfg.OperationListGCloudInstances,
				Filter:    "zone:us-central1-a",
				Format:    "table(name,status)",
			},
			expected: "list-command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, err := service.BuildCommand(tt.config)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if command != tt.expected {
				t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", tt.expected, command)
			}
		})
	}

	if !instanceOp.called {
		t.Fatal("instance operation was not called")
	}
	if !routerNATOp.called {
		t.Fatal("router and nat operation was not called")
	}
	if !iapSSHOp.called {
		t.Fatal("iap ssh operation was not called")
	}
	if !ingressSSHOp.called {
		t.Fatal("ingress ssh operation was not called")
	}
	if !listOp.called {
		t.Fatal("list operation was not called")
	}
}

func TestServiceBuildCommand_UnknownOperation(t *testing.T) {
	service := NewService()
	if _, err := service.BuildCommand(&cfg.Config{Operation: "unknown"}); err == nil {
		t.Fatal("expected error for unknown operation")
	}
}

func TestServiceBuildCommand_OperationError(t *testing.T) {
	expectedErr := fmt.Errorf("operation error")
	service := newServiceWithOperations(
		&createGCEInstanceOperationStub{err: expectedErr},
		&createGCERouterAndNATOperationStub{},
		&createGCEIAPSSHFirewallRuleOperationStub{},
		&createGCEIngressSSHFirewallRuleOperationStub{},
		&listGCloudInstancesOperationStub{},
	)

	_, err := service.BuildCommand(&cfg.Config{
		Operation:    cfg.OperationCreateGCEInstance,
		InstanceName: "vm-1",
		Zone:         "us-central1-a",
		MachineType:  "e2-medium",
		BootDiskSize: "100GB",
		BootDiskType: "pd-balanced",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != expectedErr.Error() {
		t.Fatalf("error mismatch:\nexpected: %v\nactual:   %v", expectedErr, err)
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
