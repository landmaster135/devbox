package usecases

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_compute/config"
	addstartupscripttogceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/add_startup_script_to_gce_instance"
	connectgceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/connect_gce_instance"
	copygcesshkey "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/copy_gce_ssh_key"
	creategceiapsshfirewallrule "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_iap_ssh_firewall_rule"
	creategceingresssshfirewallrule "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_ingress_ssh_firewall_rule"
	creategceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_instance"
	creategceinstanceandconfigure "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_instance_and_configure"
	creategceinstancewithstartupscript "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_instance_with_startup_script"
	creategcerouterandnat "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_router_and_nat"
	deletegceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/delete_gce_instance"
	listdisktypes "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/list_disk_types"
	listgcloudinstances "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/list_gcloud_instances"
	listmachinetypes "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/list_machine_types"
	rebootgceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/reboot_gce_instance"
	setgceinstancemetadatafromyaml "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/set_gce_instance_metadata_from_yaml"
	setupgcefirewallandssh "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/setup_gce_firewall_and_ssh"
	startgceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/start_gce_instance"
	stopgceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/stop_gce_instance"
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

type createGCEInstanceWithStartupScriptOperationStub struct {
	called bool
	got    creategceinstancewithstartupscript.Params
	result string
	err    error
}

func (s *createGCEInstanceWithStartupScriptOperationStub) Build(params creategceinstancewithstartupscript.Params) (string, error) {
	s.called = true
	s.got = params
	return s.result, s.err
}

type createGCEInstanceAndConfigureOperationStub struct {
	called bool
	got    creategceinstanceandconfigure.Params
	result string
	err    error
}

func (s *createGCEInstanceAndConfigureOperationStub) Build(params creategceinstanceandconfigure.Params) (string, error) {
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

type listDiskTypesOperationStub struct {
	executeCalled bool
	executeGot    listdisktypes.Params
	executeResult string
	executeErr    error
}

func (s *listDiskTypesOperationStub) Execute(params listdisktypes.Params) (string, error) {
	s.executeCalled = true
	s.executeGot = params
	return s.executeResult, s.executeErr
}

type listMachineTypesOperationStub struct {
	executeCalled bool
	executeGot    listmachinetypes.Params
	executeResult string
	executeErr    error
}

func (s *listMachineTypesOperationStub) Execute(params listmachinetypes.Params) (string, error) {
	s.executeCalled = true
	s.executeGot = params
	return s.executeResult, s.executeErr
}

type startGCEInstanceOperationStub struct {
	called bool
	got    startgceinstance.Params
	result string
	err    error
}

func (s *startGCEInstanceOperationStub) Build(params startgceinstance.Params) (string, error) {
	s.called = true
	s.got = params
	return s.result, s.err
}

type stopGCEInstanceOperationStub struct {
	called bool
	got    stopgceinstance.Params
	result string
	err    error
}

func (s *stopGCEInstanceOperationStub) Build(params stopgceinstance.Params) (string, error) {
	s.called = true
	s.got = params
	return s.result, s.err
}

type rebootGCEInstanceOperationStub struct {
	called bool
	got    rebootgceinstance.Params
	result string
	err    error
}

func (s *rebootGCEInstanceOperationStub) Build(params rebootgceinstance.Params) (string, error) {
	s.called = true
	s.got = params
	return s.result, s.err
}

type deleteGCEInstanceOperationStub struct {
	called bool
	got    deletegceinstance.Params
	result string
	err    error
}

func (s *deleteGCEInstanceOperationStub) Build(params deletegceinstance.Params) (string, error) {
	s.called = true
	s.got = params
	return s.result, s.err
}

type copyGCESSHKeyOperationStub struct {
	called bool
	got    copygcesshkey.Params
	result string
	err    error
}

func (s *copyGCESSHKeyOperationStub) Build(params copygcesshkey.Params) (string, error) {
	s.called = true
	s.got = params
	return s.result, s.err
}

type connectGCEInstanceOperationStub struct {
	called bool
	got    connectgceinstance.Params
	result string
	err    error
}

func (s *connectGCEInstanceOperationStub) Build(params connectgceinstance.Params) (string, error) {
	s.called = true
	s.got = params
	return s.result, s.err
}

type setGCEInstanceMetadataFromYAMLOperationStub struct {
	called bool
	got    setgceinstancemetadatafromyaml.Params
	result string
	err    error
}

func (s *setGCEInstanceMetadataFromYAMLOperationStub) Build(params setgceinstancemetadatafromyaml.Params) (string, error) {
	s.called = true
	s.got = params
	return s.result, s.err
}

type addStartupScriptToGCEInstanceOperationStub struct {
	called bool
	got    addstartupscripttogceinstance.Params
	result string
	err    error
}

func (s *addStartupScriptToGCEInstanceOperationStub) Build(params addstartupscripttogceinstance.Params) (string, error) {
	s.called = true
	s.got = params
	return s.result, s.err
}

type setupGCEFirewallAndSSHOperationStub struct {
	called bool
	got    setupgcefirewallandssh.Params
	result string
	err    error
}

func (s *setupGCEFirewallAndSSHOperationStub) Build(params setupgcefirewallandssh.Params) (string, error) {
	s.called = true
	s.got = params
	return s.result, s.err
}

func TestServiceBuildCommand_DelegatesByOperation(t *testing.T) {
	instanceOp := &createGCEInstanceOperationStub{result: "instance-command"}
	instanceWithStartupScriptOp := &createGCEInstanceWithStartupScriptOperationStub{result: "instance-with-startup-script-command"}
	instanceAndConfigureOp := &createGCEInstanceAndConfigureOperationStub{result: "instance-and-configure-command"}
	routerNATOp := &createGCERouterAndNATOperationStub{result: "router-nat-command"}
	iapSSHOp := &createGCEIAPSSHFirewallRuleOperationStub{result: "iap-ssh-command"}
	ingressSSHOp := &createGCEIngressSSHFirewallRuleOperationStub{result: "ingress-ssh-command"}
	listOp := &listGCloudInstancesOperationStub{result: "list-command"}
	listDiskTypesOp := &listDiskTypesOperationStub{}
	listMachineTypesOp := &listMachineTypesOperationStub{}
	startOp := &startGCEInstanceOperationStub{result: "start-command"}
	stopOp := &stopGCEInstanceOperationStub{result: "stop-command"}
	rebootOp := &rebootGCEInstanceOperationStub{result: "reboot-command"}
	deleteOp := &deleteGCEInstanceOperationStub{result: "delete-command"}
	copySSHKeyOp := &copyGCESSHKeyOperationStub{result: "copy-ssh-key-command"}
	connectOp := &connectGCEInstanceOperationStub{result: "connect-command"}
	setMetadataOp := &setGCEInstanceMetadataFromYAMLOperationStub{result: "set-metadata-command"}
	addStartupScriptOp := &addStartupScriptToGCEInstanceOperationStub{result: "add-startup-script-command"}
	setupOp := &setupGCEFirewallAndSSHOperationStub{result: "setup-command"}

	service := newServiceWithOperations(
		instanceOp,
		instanceWithStartupScriptOp,
		instanceAndConfigureOp,
		routerNATOp,
		iapSSHOp,
		ingressSSHOp,
		listOp,
		listDiskTypesOp,
		listMachineTypesOp,
		startOp,
		stopOp,
		rebootOp,
		deleteOp,
		copySSHKeyOp,
		connectOp,
		setMetadataOp,
		addStartupScriptOp,
		setupOp,
	)

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
			name: "create-gce-instance-with-startup-script",
			config: &cfg.Config{
				Operation:         cfg.OperationCreateGCEInstanceWithStartupScript,
				InstanceName:      "vm-1",
				Zone:              "us-central1-a",
				MachineType:       "e2-medium",
				BootDiskSize:      "100GB",
				BootDiskType:      "pd-balanced",
				StartupScriptPath: "startup-script.sh",
			},
			expected: "instance-with-startup-script-command",
		},
		{
			name: "create-gce-instance-and-configure",
			config: &cfg.Config{
				Operation:         cfg.OperationCreateGCEInstanceAndConfigure,
				InstanceName:      "vm-1",
				Zone:              "us-central1-a",
				MachineType:       "e2-medium",
				BootDiskSize:      "100GB",
				BootDiskType:      "pd-balanced",
				MetadataYAMLPath:  "env.yml",
				StartupScriptPath: "startup-script.sh",
			},
			expected: "instance-and-configure-command",
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
		{
			name: "start-gce-instance",
			config: &cfg.Config{
				Operation:    cfg.OperationStartGCEInstance,
				InstanceName: "vm-1",
				Zone:         "us-central1-a",
			},
			expected: "start-command",
		},
		{
			name: "stop-gce-instance",
			config: &cfg.Config{
				Operation:    cfg.OperationStopGCEInstance,
				InstanceName: "vm-1",
				Zone:         "us-central1-a",
			},
			expected: "stop-command",
		},
		{
			name: "reboot-gce-instance",
			config: &cfg.Config{
				Operation:    cfg.OperationRebootGCEInstance,
				InstanceName: "vm-1",
				Zone:         "us-central1-a",
			},
			expected: "reboot-command",
		},
		{
			name: "delete-gce-instance",
			config: &cfg.Config{
				Operation:    cfg.OperationDeleteGCEInstance,
				InstanceName: "vm-1",
				Zone:         "us-central1-a",
			},
			expected: "delete-command",
		},
		{
			name: "copy-gce-ssh-key",
			config: &cfg.Config{
				Operation:    cfg.OperationCopyGCESSHKey,
				InstanceName: "vm-1",
				Zone:         "us-central1-a",
				SSHKeyPath:   "$HOME/.ssh/google_compute_engine",
			},
			expected: "copy-ssh-key-command",
		},
		{
			name: "connect-gce-instance",
			config: &cfg.Config{
				Operation:    cfg.OperationConnectGCEInstance,
				InstanceName: "vm-1",
				Zone:         "us-central1-a",
			},
			expected: "connect-command",
		},
		{
			name: "set-gce-instance-metadata-from-yaml",
			config: &cfg.Config{
				Operation:        cfg.OperationSetGCEInstanceMetadataFromYAML,
				InstanceName:     "vm-1",
				Zone:             "us-central1-a",
				MetadataYAMLPath: "env.yml",
			},
			expected: "set-metadata-command",
		},
		{
			name: "add-startup-script-to-gce-instance",
			config: &cfg.Config{
				Operation:         cfg.OperationAddStartupScriptToGCEInstance,
				InstanceName:      "vm-1",
				Zone:              "us-central1-a",
				StartupScriptPath: "startup-script.sh",
			},
			expected: "add-startup-script-command",
		},
		{
			name: "setup-gce-firewall-and-ssh",
			config: &cfg.Config{
				Operation:    cfg.OperationSetupGCEFirewallAndSSH,
				InstanceName: "vm-1",
				Zone:         "us-central1-a",
				SSHKeyPath:   "$HOME/.ssh/google_compute_engine",
			},
			expected: "setup-command",
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
	if !instanceWithStartupScriptOp.called {
		t.Fatal("instance with startup script operation was not called")
	}
	if !instanceAndConfigureOp.called {
		t.Fatal("instance and configure operation was not called")
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
	if !startOp.called {
		t.Fatal("start operation was not called")
	}
	if !stopOp.called {
		t.Fatal("stop operation was not called")
	}
	if !rebootOp.called {
		t.Fatal("reboot operation was not called")
	}
	if !deleteOp.called {
		t.Fatal("delete operation was not called")
	}
	if !copySSHKeyOp.called {
		t.Fatal("copy ssh key operation was not called")
	}
	if !connectOp.called {
		t.Fatal("connect operation was not called")
	}
	if !setMetadataOp.called {
		t.Fatal("set metadata operation was not called")
	}
	if !addStartupScriptOp.called {
		t.Fatal("add startup script operation was not called")
	}
	if !setupOp.called {
		t.Fatal("setup operation was not called")
	}
}

func TestServiceBuildCommand_UnknownOperation(t *testing.T) {
	service := NewService()
	if _, err := service.BuildCommand(&cfg.Config{Operation: "unknown"}); err == nil {
		t.Fatal("expected error for unknown operation")
	}
}

func TestServiceBuildCommand_ListOperationsAreUnsupported(t *testing.T) {
	service := NewService()
	tests := []string{
		cfg.OperationListDiskTypes,
		cfg.OperationListMachineTypes,
	}

	for _, operation := range tests {
		if _, err := service.BuildCommand(&cfg.Config{Operation: operation}); err == nil {
			t.Fatalf("expected error for operation: %s", operation)
		}
	}
}

func TestServiceBuildCommand_OperationError(t *testing.T) {
	expectedErr := fmt.Errorf("operation error")
	service := newServiceWithOperations(
		&createGCEInstanceOperationStub{err: expectedErr},
		&createGCEInstanceWithStartupScriptOperationStub{},
		&createGCEInstanceAndConfigureOperationStub{},
		&createGCERouterAndNATOperationStub{},
		&createGCEIAPSSHFirewallRuleOperationStub{},
		&createGCEIngressSSHFirewallRuleOperationStub{},
		&listGCloudInstancesOperationStub{},
		&listDiskTypesOperationStub{},
		&listMachineTypesOperationStub{},
		&startGCEInstanceOperationStub{},
		&stopGCEInstanceOperationStub{},
		&rebootGCEInstanceOperationStub{},
		&deleteGCEInstanceOperationStub{},
		&copyGCESSHKeyOperationStub{},
		&connectGCEInstanceOperationStub{},
		&setGCEInstanceMetadataFromYAMLOperationStub{},
		&addStartupScriptToGCEInstanceOperationStub{},
		&setupGCEFirewallAndSSHOperationStub{},
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

func TestServiceExecuteListDiskTypes_Normal(t *testing.T) {
	listDiskTypesOp := &listDiskTypesOperationStub{executeResult: "disk-types-result"}
	service := newServiceWithOperations(
		&createGCEInstanceOperationStub{},
		&createGCEInstanceWithStartupScriptOperationStub{},
		&createGCEInstanceAndConfigureOperationStub{},
		&createGCERouterAndNATOperationStub{},
		&createGCEIAPSSHFirewallRuleOperationStub{},
		&createGCEIngressSSHFirewallRuleOperationStub{},
		&listGCloudInstancesOperationStub{},
		listDiskTypesOp,
		&listMachineTypesOperationStub{},
		&startGCEInstanceOperationStub{},
		&stopGCEInstanceOperationStub{},
		&rebootGCEInstanceOperationStub{},
		&deleteGCEInstanceOperationStub{},
		&copyGCESSHKeyOperationStub{},
		&connectGCEInstanceOperationStub{},
		&setGCEInstanceMetadataFromYAMLOperationStub{},
		&addStartupScriptToGCEInstanceOperationStub{},
		&setupGCEFirewallAndSSHOperationStub{},
	)

	result, err := service.ExecuteListDiskTypes(ListDiskTypesParams{
		Zones:          []string{"asia-southeast3-a"},
		MinDiskSizeGiB: 4,
		MaxDiskSizeGiB: 65536,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "disk-types-result" {
		t.Fatalf("result mismatch: %s", result)
	}
	if !listDiskTypesOp.executeCalled {
		t.Fatal("execute was not called")
	}
}

func TestServiceExecuteListMachineTypes_Normal(t *testing.T) {
	listMachineTypesOp := &listMachineTypesOperationStub{executeResult: "machine-types-result"}
	service := newServiceWithOperations(
		&createGCEInstanceOperationStub{},
		&createGCEInstanceWithStartupScriptOperationStub{},
		&createGCEInstanceAndConfigureOperationStub{},
		&createGCERouterAndNATOperationStub{},
		&createGCEIAPSSHFirewallRuleOperationStub{},
		&createGCEIngressSSHFirewallRuleOperationStub{},
		&listGCloudInstancesOperationStub{},
		&listDiskTypesOperationStub{},
		listMachineTypesOp,
		&startGCEInstanceOperationStub{},
		&stopGCEInstanceOperationStub{},
		&rebootGCEInstanceOperationStub{},
		&deleteGCEInstanceOperationStub{},
		&copyGCESSHKeyOperationStub{},
		&connectGCEInstanceOperationStub{},
		&setGCEInstanceMetadataFromYAMLOperationStub{},
		&addStartupScriptToGCEInstanceOperationStub{},
		&setupGCEFirewallAndSSHOperationStub{},
	)

	result, err := service.ExecuteListMachineTypes(ListMachineTypesParams{
		Zones:          []string{"asia-southeast3-a"},
		MinDiskSizeGiB: 1024,
		MaxDiskSizeGiB: 65536,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "machine-types-result" {
		t.Fatalf("result mismatch: %s", result)
	}
	if !listMachineTypesOp.executeCalled {
		t.Fatal("execute was not called")
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
