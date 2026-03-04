package usecases

import (
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
	listgcloudinstances "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/list_gcloud_instances"
	rebootgceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/reboot_gce_instance"
	setgceinstancemetadatafromyaml "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/set_gce_instance_metadata_from_yaml"
	setupgcefirewallandssh "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/setup_gce_firewall_and_ssh"
	startgceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/start_gce_instance"
	stopgceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/stop_gce_instance"
)

type createGCEInstanceOperation interface {
	Build(params creategceinstance.Params) (string, error)
}

type createGCEInstanceWithStartupScriptOperation interface {
	Build(params creategceinstancewithstartupscript.Params) (string, error)
}

type createGCEInstanceAndConfigureOperation interface {
	Build(params creategceinstanceandconfigure.Params) (string, error)
}

type createGCERouterAndNATOperation interface {
	Build(params creategcerouterandnat.Params) (string, error)
}

type createGCEIAPSSHFirewallRuleOperation interface {
	Build(params creategceiapsshfirewallrule.Params) (string, error)
}

type createGCEIngressSSHFirewallRuleOperation interface {
	Build(params creategceingresssshfirewallrule.Params) (string, error)
}

type listGCloudInstancesOperation interface {
	Build(params listgcloudinstances.Params) (string, error)
}

type startGCEInstanceOperation interface {
	Build(params startgceinstance.Params) (string, error)
}

type stopGCEInstanceOperation interface {
	Build(params stopgceinstance.Params) (string, error)
}

type rebootGCEInstanceOperation interface {
	Build(params rebootgceinstance.Params) (string, error)
}

type deleteGCEInstanceOperation interface {
	Build(params deletegceinstance.Params) (string, error)
}

type copyGCESSHKeyOperation interface {
	Build(params copygcesshkey.Params) (string, error)
}

type connectGCEInstanceOperation interface {
	Build(params connectgceinstance.Params) (string, error)
}

type setGCEInstanceMetadataFromYAMLOperation interface {
	Build(params setgceinstancemetadatafromyaml.Params) (string, error)
}

type addStartupScriptToGCEInstanceOperation interface {
	Build(params addstartupscripttogceinstance.Params) (string, error)
}

type setupGCEFirewallAndSSHOperation interface {
	Build(params setupgcefirewallandssh.Params) (string, error)
}

func newServiceWithOperations(
	createGCEInstanceOp createGCEInstanceOperation,
	createGCEInstanceWithStartupScriptOp createGCEInstanceWithStartupScriptOperation,
	createGCEInstanceAndConfigureOp createGCEInstanceAndConfigureOperation,
	createGCERouterAndNATOp createGCERouterAndNATOperation,
	createGCEIAPSSHFirewallRuleOp createGCEIAPSSHFirewallRuleOperation,
	createGCEIngressSSHFirewallRuleOp createGCEIngressSSHFirewallRuleOperation,
	listGCloudInstancesOp listGCloudInstancesOperation,
	startGCEInstanceOp startGCEInstanceOperation,
	stopGCEInstanceOp stopGCEInstanceOperation,
	rebootGCEInstanceOp rebootGCEInstanceOperation,
	deleteGCEInstanceOp deleteGCEInstanceOperation,
	copyGCESSHKeyOp copyGCESSHKeyOperation,
	connectGCEInstanceOp connectGCEInstanceOperation,
	setGCEInstanceMetadataFromYAMLOp setGCEInstanceMetadataFromYAMLOperation,
	addStartupScriptToGCEInstanceOp addStartupScriptToGCEInstanceOperation,
	setupGCEFirewallAndSSHOp setupGCEFirewallAndSSHOperation,
) *Service {
	return &Service{
		createGCEInstanceOperation:                  createGCEInstanceOp,
		createGCEInstanceWithStartupScriptOperation: createGCEInstanceWithStartupScriptOp,
		createGCEInstanceAndConfigureOperation:      createGCEInstanceAndConfigureOp,
		createGCERouterAndNATOperation:              createGCERouterAndNATOp,
		createGCEIAPSSHFirewallRuleOperation:        createGCEIAPSSHFirewallRuleOp,
		createGCEIngressSSHFirewallRuleOperation:    createGCEIngressSSHFirewallRuleOp,
		listGCloudInstancesOperation:                listGCloudInstancesOp,
		startGCEInstanceOperation:                   startGCEInstanceOp,
		stopGCEInstanceOperation:                    stopGCEInstanceOp,
		rebootGCEInstanceOperation:                  rebootGCEInstanceOp,
		deleteGCEInstanceOperation:                  deleteGCEInstanceOp,
		copyGCESSHKeyOperation:                      copyGCESSHKeyOp,
		connectGCEInstanceOperation:                 connectGCEInstanceOp,
		setGCEInstanceMetadataFromYAMLOperation:     setGCEInstanceMetadataFromYAMLOp,
		addStartupScriptToGCEInstanceOperation:      addStartupScriptToGCEInstanceOp,
		setupGCEFirewallAndSSHOperation:             setupGCEFirewallAndSSHOp,
	}
}
