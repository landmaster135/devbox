package usecases

import (
	"fmt"
	"strings"
	"testing"
)

func TestBuildCreateSecretCommand(t *testing.T) {
	service := NewService()

	t.Run("automatic replication", func(t *testing.T) {
		command, err := service.BuildCreateSecretCommand(CreateSecretParams{
			SecretName: "my-secret",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "gcloud secrets create 'my-secret' --replication-policy='automatic'"
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})

	t.Run("user managed with locations", func(t *testing.T) {
		command, err := service.BuildCreateSecretCommand(CreateSecretParams{
			SecretName:        "my-secret",
			ReplicationPolicy: replicationPolicyUserManaged,
			Locations:         "asia-northeast1,us-east1",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "gcloud secrets create 'my-secret' --replication-policy='user-managed' --locations='asia-northeast1,us-east1'"
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})

	t.Run("user managed without locations", func(t *testing.T) {
		_, err := service.BuildCreateSecretCommand(CreateSecretParams{
			SecretName:        "my-secret",
			ReplicationPolicy: replicationPolicyUserManaged,
		})
		if err == nil {
			t.Fatal("expected error when locations are missing")
		}
	})

	t.Run("invalid replication policy", func(t *testing.T) {
		_, err := service.BuildCreateSecretCommand(CreateSecretParams{
			SecretName:        "my-secret",
			ReplicationPolicy: "invalid",
		})
		if err == nil {
			t.Fatal("expected error for invalid replication policy")
		}
	})
}

func TestBuildAddSecretVersionCommand(t *testing.T) {
	service := NewService()

	t.Run("basic", func(t *testing.T) {
		command, err := service.BuildAddSecretVersionCommand(AddSecretVersionParams{
			SecretName:  "app-secret",
			SecretValue: "super-secret",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "echo -n 'super-secret' | gcloud secrets versions add 'app-secret' --data-file=-"
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})

	t.Run("quote escaping", func(t *testing.T) {
		command, err := service.BuildAddSecretVersionCommand(AddSecretVersionParams{
			SecretName:  "app-secret",
			SecretValue: "pa'ss",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "echo -n 'pa'\"'\"'ss' | gcloud secrets versions add 'app-secret' --data-file=-"
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})

	t.Run("missing value", func(t *testing.T) {
		_, err := service.BuildAddSecretVersionCommand(AddSecretVersionParams{
			SecretName: "app-secret",
		})
		if err == nil {
			t.Fatal("expected error when secret value is missing")
		}
	})
}

func TestBuildCreateAndAddSecretVersionCommand(t *testing.T) {
	service := NewService()

	command, err := service.BuildCreateAndAddSecretVersionCommand(CreateAndAddSecretVersionParams{
		SecretName:  "combo-secret",
		SecretValue: "combo",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud secrets create 'combo-secret' --replication-policy='automatic' && echo -n 'combo' | gcloud secrets versions add 'combo-secret' --data-file=-"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildAccessSecretVersionCommand(t *testing.T) {
	service := NewService()

	t.Run("default version", func(t *testing.T) {
		command, err := service.BuildAccessSecretVersionCommand(AccessSecretVersionParams{
			SecretName: "config-secret",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "gcloud secrets versions access 'latest' --secret='config-secret'"
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})

	t.Run("specific version", func(t *testing.T) {
		command, err := service.BuildAccessSecretVersionCommand(AccessSecretVersionParams{
			SecretName: "config-secret",
			Version:    "5",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "gcloud secrets versions access '5' --secret='config-secret'"
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})
}

func TestBuildUpdateSecretLabelsCommand(t *testing.T) {
	service := NewService()

	command, err := service.BuildUpdateSecretLabelsCommand(UpdateSecretLabelsParams{
		SecretName: "labeled",
		Labels:     "env=prod,team=devops",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud secrets update 'labeled' --update-labels='env=prod,team=devops'"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildUpdateSecretVersionAliasesCommand(t *testing.T) {
	service := NewService()

	t.Run("update aliases", func(t *testing.T) {
		command, err := service.BuildUpdateSecretVersionAliasesCommand(UpdateSecretVersionAliasesParams{
			SecretName:  "alias",
			AliasOption: "--update-version-aliases=prod=5,staging=3",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "gcloud secrets update 'alias' --update-version-aliases=prod=5,staging=3"
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})

	t.Run("invalid option", func(t *testing.T) {
		_, err := service.BuildUpdateSecretVersionAliasesCommand(UpdateSecretVersionAliasesParams{
			SecretName:  "alias",
			AliasOption: "--invalid",
		})
		if err == nil {
			t.Fatal("expected error for invalid alias option")
		}
	})
}

func TestBuildNotificationWrappedCommand(t *testing.T) {
	service := NewService()

	script, ok := service.BuildNotificationWrappedCommand(DiscordNotificationParams{
		Operation:  "create-secret",
		SecretName: "test-secret",
	}, "gcloud secrets create 'test-secret' --replication-policy='automatic'")

	if !ok {
		t.Fatalf("expected notification script")
	}

	expected := strings.Join([]string{
		fmt.Sprintf("%s \\", discordCLIPath),
		fmt.Sprintf("  -webhook-url \"$%s\" \\", discordWebhookEnvVarName),
		"  -content-text 'シークレットを作るよ！' \\",
		"  -embed-type 'none'",
		"if gcloud secrets create 'test-secret' --replication-policy='automatic'; then",
		fmt.Sprintf("  %s \\", discordCLIPath),
		fmt.Sprintf("    -webhook-url \"$%s\" \\", discordWebhookEnvVarName),
		"    -content-text '作ったよ！' \\",
		"    -embed-type 'google-secret-manager-success' \\",
		"    -embed-text 'シークレットを作ったよ！'",
		"else",
		fmt.Sprintf("  %s \\", discordCLIPath),
		fmt.Sprintf("    -webhook-url \"$%s\" \\", discordWebhookEnvVarName),
		"    -content-text '失敗…' \\",
		"    -embed-type 'google-secret-manager-failed' \\",
		"    -embed-text 'シークレットが作れなかったよ…'",
		"fi",
	}, "\n")
	if script != expected {
		t.Fatalf("unexpected notification script:\n%s", script)
	}
}
