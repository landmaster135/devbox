package creategceiapsshfirewallrule

import "testing"

func TestServiceBuild_Normal(t *testing.T) {
	t.Parallel()

	service := NewService()
	command, err := service.Build(Params{
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

func TestServiceBuild_ValidationError(t *testing.T) {
	t.Parallel()

	service := NewService()

	tests := []Params{
		{Direction: "INGRESS", Action: "allow", Rules: "tcp:22", SourceRanges: "35.235.240.0/20"},
		{RuleName: "r", Action: "allow", Rules: "tcp:22", SourceRanges: "35.235.240.0/20"},
		{RuleName: "r", Direction: "INGRESS", Rules: "tcp:22", SourceRanges: "35.235.240.0/20"},
		{RuleName: "r", Direction: "INGRESS", Action: "allow", SourceRanges: "35.235.240.0/20"},
		{RuleName: "r", Direction: "INGRESS", Action: "allow", Rules: "tcp:22"},
	}

	for i, params := range tests {
		if _, err := service.Build(params); err == nil {
			t.Fatalf("expected validation error at case %d", i)
		}
	}
}
