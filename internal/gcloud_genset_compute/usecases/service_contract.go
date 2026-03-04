package usecases

import (
	creategceiapsshfirewallrule "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_iap_ssh_firewall_rule"
	creategceingresssshfirewallrule "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_ingress_ssh_firewall_rule"
	creategceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_instance"
	creategcerouterandnat "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_router_and_nat"
	listgcloudinstances "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/list_gcloud_instances"
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

func newServiceWithOperations(
	createGCEInstanceOp createGCEInstanceOperation,
	createGCERouterAndNATOp createGCERouterAndNATOperation,
	createGCEIAPSSHFirewallRuleOp createGCEIAPSSHFirewallRuleOperation,
	createGCEIngressSSHFirewallRuleOp createGCEIngressSSHFirewallRuleOperation,
	listGCloudInstancesOp listGCloudInstancesOperation,
) *Service {
	return &Service{
		createGCEInstanceOperation:               createGCEInstanceOp,
		createGCERouterAndNATOperation:           createGCERouterAndNATOp,
		createGCEIAPSSHFirewallRuleOperation:     createGCEIAPSSHFirewallRuleOp,
		createGCEIngressSSHFirewallRuleOperation: createGCEIngressSSHFirewallRuleOp,
		listGCloudInstancesOperation:             listGCloudInstancesOp,
	}
}
