package usecases

import (
	"strings"
	"testing"
)

func TestBuildAddIamPolicyBindingToProjectCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildAddIamPolicyBindingToProjectCommand(AddIamPolicyBindingToProjectParams{
		ProjectID:        "sample-project",
		ServiceAccountID: "sa-id",
		Role:             "roles/viewer",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud projects add-iam-policy-binding 'sample-project' --member='serviceAccount:sa-id@sample-project.iam.gserviceaccount.com' --role='roles/viewer'"
	if cmd != expected {
		t.Fatalf("unexpected command. want=%q, got=%q", expected, cmd)
	}
}

func TestBuildAddIamPolicyBindingToServiceAccountCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildAddIamPolicyBindingToServiceAccountCommand(AddIamPolicyBindingToServiceAccountParams{
		ServiceAccountEmail: "sa@example.iam.gserviceaccount.com",
		Member:              "user:dev@example.com",
		Role:                "roles/iam.serviceAccountUser",
		Condition:           "expression=request.time < timestamp('2025-01-01T00:00:00Z')",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(cmd, "--condition='expression=request.time < timestamp('") {
		t.Fatalf("condition flag not found in command: %s", cmd)
	}
}

func TestBuildAddIamPolicyBindingToServiceAccountCommand_Errors(t *testing.T) {
	service := NewService()
	_, err := service.BuildAddIamPolicyBindingToServiceAccountCommand(AddIamPolicyBindingToServiceAccountParams{
		ServiceAccountEmail: "sa@example.iam.gserviceaccount.com",
		Member:              "user:dev@example.com",
		Role:                "roles/iam.serviceAccountUser",
		Condition:           "expr",
		ConditionFromFile:   "path",
	})
	if err == nil {
		t.Fatal("expected error when both condition and condition-from-file are specified")
	}
}

func TestBuildAddWorkloadIdentityBindingToServiceAccountCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildAddWorkloadIdentityBindingToServiceAccountCommand(AddWorkloadIdentityBindingToServiceAccountParams{
		ServiceAccountEmail: "sa@example.iam.gserviceaccount.com",
		ProjectNumber:       "1234567890",
		PoolID:              "github-pool",
		RepositoryOwner:     "landmaster135",
		RepositoryName:      "devbox",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedMember := "--member='principalSet://iam.googleapis.com/projects/1234567890/locations/global/workloadIdentityPools/github-pool/attribute.repository/landmaster135/devbox'"
	if !strings.Contains(cmd, expectedMember) {
		t.Fatalf("principalSet not embedded correctly: %s", cmd)
	}
}

func TestBuildCreateServiceAccountCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildCreateServiceAccountCommand(CreateServiceAccountParams{
		ServiceAccountID: "ci-bot",
		ProjectID:        "sample-project",
		Role:             "roles/run.admin",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(cmd, "gcloud iam service-accounts create 'ci-bot'") {
		t.Fatalf("create command missing: %s", cmd)
	}
	if !strings.Contains(cmd, "&&\n") {
		t.Fatalf("commands are not chained with &&: %s", cmd)
	}
}

func TestBuildListServiceAccountsCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildListServiceAccountsCommand(ListServiceAccountsParams{
		Filter: "email:example.com",
		SortBy: "name",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(cmd, "--filter='email:example.com'") {
		t.Fatalf("filter flag missing: %s", cmd)
	}
	if !strings.Contains(cmd, "--sort-by='name'") {
		t.Fatalf("sort-by flag missing: %s", cmd)
	}
}

func TestBuildUpdateServiceAccountCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildUpdateServiceAccountCommand(UpdateServiceAccountParams{
		ServiceAccountEmail: "sa@example.iam.gserviceaccount.com",
		Description:         "updated desc",
		DisplayName:         "Updated Name",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(cmd, "--description='updated desc'") || !strings.Contains(cmd, "--display-name='Updated Name'") {
		t.Fatalf("update options missing: %s", cmd)
	}
}

func TestBuildCreateWorkloadIdentityPoolCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildCreateWorkloadIdentityPoolCommand(CreateWorkloadIdentityPoolParams{
		WorkloadIdentityPoolBaseParams: WorkloadIdentityPoolBaseParams{
			ProjectID: "sample-project",
			PoolID:    "github-pool",
			Location:  "global",
		},
		Description: "GitHub Actions",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(cmd, "--description='GitHub Actions'") {
		t.Fatalf("description flag missing: %s", cmd)
	}
}

func TestBuildUpdateWorkloadIdentityPoolCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildUpdateWorkloadIdentityPoolCommand(UpdateWorkloadIdentityPoolParams{
		WorkloadIdentityPoolBaseParams: WorkloadIdentityPoolBaseParams{
			ProjectID: "sample-project",
			PoolID:    "github-pool",
			Location:  "global",
		},
		Description: "Updated description",
		Disabled:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(cmd, "--disabled") {
		t.Fatalf("disabled flag missing: %s", cmd)
	}
}

func TestBuildCreateOidcWorkloadIdentityPoolProviderCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildCreateOidcWorkloadIdentityPoolProviderCommand(CreateOidcWorkloadIdentityPoolProviderParams{
		WorkloadIdentityPoolProviderBaseParams: WorkloadIdentityPoolProviderBaseParams{
			ProjectID:  "sample-project",
			PoolID:     "github-pool",
			ProviderID: "github-provider",
			Location:   "global",
		},
		IssuerURI:          "https://issuer.example.com",
		AttributeMapping:   "google.subject=assertion.sub",
		AttributeCondition: "assertion.aud=='devbox'",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(cmd, "--issuer-uri='https://issuer.example.com'") {
		t.Fatalf("issuer flag missing: %s", cmd)
	}
}

func TestBuildCreateOidcWorkloadIdentityPoolProviderForGitHubActionsCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildCreateOidcWorkloadIdentityPoolProviderForGitHubActionsCommand(CreateOidcWorkloadIdentityPoolProviderForGitHubActionsParams{
		WorkloadIdentityPoolProviderBaseParams: WorkloadIdentityPoolProviderBaseParams{
			ProjectID:  "sample-project",
			PoolID:     "github-pool",
			ProviderID: "github-provider",
			Location:   "global",
		},
		RepositoryOwner: "landmaster135",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(cmd, "--issuer-uri='https://token.actions.githubusercontent.com/'") {
		t.Fatalf("issuer fixed value missing: %s", cmd)
	}
}

func TestBuildUpdateOidcWorkloadIdentityPoolProviderCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildUpdateOidcWorkloadIdentityPoolProviderCommand(UpdateOidcWorkloadIdentityPoolProviderParams{
		WorkloadIdentityPoolProviderBaseParams: WorkloadIdentityPoolProviderBaseParams{
			ProjectID:  "sample-project",
			PoolID:     "github-pool",
			ProviderID: "github-provider",
			Location:   "global",
		},
		AllowedAudiences:   "https://github.com",
		AttributeCondition: "assertion.aud=='devbox'",
		AttributeMapping:   "google.subject=assertion.sub",
		Description:        "Updated provider",
		Disabled:           true,
		DisplayName:        "GitHub Provider",
		IssuerURI:          "https://issuer.example.com",
		JWKJSONPath:        "./keys.json",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checks := []string{
		"--allowed-audiences='https://github.com'",
		"--attribute-mapping='google.subject=assertion.sub'",
		"--description='Updated provider'",
		"--disabled",
		"--display-name='GitHub Provider'",
		"--issuer-uri='https://issuer.example.com'",
		"--jwk-json-path='./keys.json'",
	}

	for _, c := range checks {
		if !strings.Contains(cmd, c) {
			t.Fatalf("flag %q missing in command: %s", c, cmd)
		}
	}

	if !strings.Contains(cmd, "--attribute-condition=") {
		t.Fatalf("attribute-condition flag missing: %s", cmd)
	}
}

func TestBuildSetupWorkloadIdentityFederationScript(t *testing.T) {
	service := NewService()
	script, err := service.BuildSetupWorkloadIdentityFederationScript(SetupWorkloadIdentityFederationParams{
		ProjectID:        "sample-project",
		PoolID:           "github-pool",
		ProviderID:       "github-provider",
		ServiceAccountID: "gha",
		RepositoryOwner:  "landmaster135",
		RepositoryName:   "devbox",
		Location:         "global",
		PoolDescription:  "GitHub Actions",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(script, "set -e") {
		t.Fatalf("script must enable set -e: %s", script)
	}
	if !strings.Contains(script, "PROJECT_NUMBER=$(gcloud projects describe 'sample-project' --format=value(projectNumber))") {
		t.Fatalf("project number retrieval missing: %s", script)
	}
	if !strings.Contains(script, "--member=\"principalSet://iam.googleapis.com/projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/github-pool/attribute.repository/landmaster135/devbox\"") {
		t.Fatalf("principal set with variable missing: %s", script)
	}
}

func TestBuildCleanupWorkloadIdentityFederationScript(t *testing.T) {
	service := NewService()
	script, err := service.BuildCleanupWorkloadIdentityFederationScript(CleanupWorkloadIdentityFederationParams{
		ProjectID:        "sample-project",
		PoolID:           "github-pool",
		ProviderID:       "github-provider",
		ServiceAccountID: "gha",
		Location:         "global",
		SkipConfirmation: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(script, "read -p") {
		t.Fatalf("confirmation prompt missing: %s", script)
	}

	scriptSkip, err := service.BuildCleanupWorkloadIdentityFederationScript(CleanupWorkloadIdentityFederationParams{
		ProjectID:        "sample-project",
		PoolID:           "github-pool",
		ProviderID:       "github-provider",
		ServiceAccountID: "gha",
		Location:         "global",
		SkipConfirmation: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(scriptSkip, "read -p") {
		t.Fatalf("confirmation prompt should be omitted when skip-confirmation is true: %s", scriptSkip)
	}
}

func TestBuildNotificationWrappedCommand(t *testing.T) {
	service := NewService()
	script, ok := service.BuildNotificationWrappedCommand(
		DiscordNotificationParams{Operation: "add-iam-policy-binding-to-project"},
		"gcloud test",
	)
	if !ok {
		t.Fatal("expected notification script to be generated")
	}

	if !strings.Contains(script, "サービスアカウントにIAMポリシーをバインドするよ！") {
		t.Fatalf("start notification missing: %s", script)
	}
	if !strings.Contains(script, "if gcloud test; then") {
		t.Fatalf("gcloud command not wrapped: %s", script)
	}

	if _, ok := service.BuildNotificationWrappedCommand(DiscordNotificationParams{Operation: "unknown"}, "cmd"); ok {
		t.Fatal("unexpected notification for unknown operation")
	}
}
