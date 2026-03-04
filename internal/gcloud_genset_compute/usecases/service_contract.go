package usecases

import (
	creategceiapsshfirewallrule "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_iap_ssh_firewall_rule"
	creategceingresssshfirewallrule "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_ingress_ssh_firewall_rule"
	creategceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_instance"
	creategcerouterandnat "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_router_and_nat"
	deletegceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/delete_gce_instance"
	listgcloudinstances "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/list_gcloud_instances"
	rebootgceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/reboot_gce_instance"
	startgceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/start_gce_instance"
	stopgceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/stop_gce_instance"
)

type createGCEInstanceOperation interface {
	Build(params creategceinstance.Params) (string, error)
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

func newServiceWithOperations(
	createGCEInstanceOp createGCEInstanceOperation,
	createGCERouterAndNATOp createGCERouterAndNATOperation,
	createGCEIAPSSHFirewallRuleOp createGCEIAPSSHFirewallRuleOperation,
	createGCEIngressSSHFirewallRuleOp createGCEIngressSSHFirewallRuleOperation,
	listGCloudInstancesOp listGCloudInstancesOperation,
	startGCEInstanceOp startGCEInstanceOperation,
	stopGCEInstanceOp stopGCEInstanceOperation,
	rebootGCEInstanceOp rebootGCEInstanceOperation,
	deleteGCEInstanceOp deleteGCEInstanceOperation,
) *Service {
	return &Service{
		createGCEInstanceOperation:               createGCEInstanceOp,
		createGCERouterAndNATOperation:           createGCERouterAndNATOp,
		createGCEIAPSSHFirewallRuleOperation:     createGCEIAPSSHFirewallRuleOp,
		createGCEIngressSSHFirewallRuleOperation: createGCEIngressSSHFirewallRuleOp,
		listGCloudInstancesOperation:             listGCloudInstancesOp,
		startGCEInstanceOperation:                startGCEInstanceOp,
		stopGCEInstanceOperation:                 stopGCEInstanceOp,
		rebootGCEInstanceOperation:               rebootGCEInstanceOp,
		deleteGCEInstanceOperation:               deleteGCEInstanceOp,
	}
}
