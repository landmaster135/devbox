package creategceingresssshfirewallrule

import "testing"

func TestServiceBuild_Normal(t *testing.T) {
	t.Parallel()

	service := NewService()
	command, err := service.Build(Params{
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

func TestServiceBuild_ValidationError(t *testing.T) {
	t.Parallel()

	service := NewService()

	tests := []Params{
		{AllowRule: "tcp:22", SourceRanges: "10.0.0.0/8"},
		{RuleName: "allow-ingress-ssh", SourceRanges: "10.0.0.0/8"},
		{RuleName: "allow-ingress-ssh", AllowRule: "tcp:22"},
	}

	for i, params := range tests {
		if _, err := service.Build(params); err == nil {
			t.Fatalf("expected validation error at case %d", i)
		}
	}
}
