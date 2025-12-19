package usecases

import (
	"bytes"
	"fmt"
	"io"
	"os"
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

	t.Run("missing name", func(t *testing.T) {
		_, err := service.BuildAddSecretVersionCommand(AddSecretVersionParams{SecretValue: "value"})
		if err == nil {
			t.Fatal("expected error when secret name is missing")
		}
	})
}

func TestBuildCreateAndAddSecretVersionCommand(t *testing.T) {
	service := NewService()

	t.Run("success", func(t *testing.T) {
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
	})

	t.Run("create secret error", func(t *testing.T) {
		_, err := service.BuildCreateAndAddSecretVersionCommand(CreateAndAddSecretVersionParams{
			SecretValue: "value",
		})
		if err == nil {
			t.Fatal("expected error when secret name is missing")
		}
	})

	t.Run("add version error", func(t *testing.T) {
		_, err := service.BuildCreateAndAddSecretVersionCommand(CreateAndAddSecretVersionParams{
			SecretName: "combo-secret",
		})
		if err == nil {
			t.Fatal("expected error when secret value is missing")
		}
	})
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

	t.Run("missing secret", func(t *testing.T) {
		_, err := service.BuildAccessSecretVersionCommand(AccessSecretVersionParams{})
		if err == nil {
			t.Fatal("expected error when secret name is missing")
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

	if _, err := service.BuildUpdateSecretLabelsCommand(UpdateSecretLabelsParams{Labels: "env=prod"}); err == nil {
		t.Fatal("expected error when secret name is missing")
	}
	if _, err := service.BuildUpdateSecretLabelsCommand(UpdateSecretLabelsParams{SecretName: "labeled"}); err == nil {
		t.Fatal("expected error when labels are missing")
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

	t.Run("clear aliases", func(t *testing.T) {
		command, err := service.BuildUpdateSecretVersionAliasesCommand(UpdateSecretVersionAliasesParams{
			SecretName:  "alias",
			AliasOption: "--clear-version-aliases",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "gcloud secrets update 'alias' --clear-version-aliases"
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

	t.Run("missing secret", func(t *testing.T) {
		_, err := service.BuildUpdateSecretVersionAliasesCommand(UpdateSecretVersionAliasesParams{AliasOption: "--clear-version-aliases"})
		if err == nil {
			t.Fatal("expected error when secret name is missing")
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

	t.Run("unknown operation", func(t *testing.T) {
		if _, ok := service.BuildNotificationWrappedCommand(DiscordNotificationParams{Operation: "unknown"}, "cmd"); ok {
			t.Fatalf("expected false when template is missing")
		}
	})
}

func TestPrintHighlightedCommand(t *testing.T) {
	service := NewService()

	output := captureStdout(func() {
		service.PrintHighlightedCommand("gcloud secrets list")
	})

	if !strings.Contains(output, "生成された gcloud コマンド") {
		t.Fatalf("expected header in output: %s", output)
	}
	if !strings.Contains(output, "gcloud secrets list") {
		t.Fatalf("expected command in output: %s", output)
	}
}

func TestPrintNotificationScript(t *testing.T) {
	service := NewService()

	output := captureStdout(func() {
		service.PrintNotificationScript("echo hi")
	})
	if !strings.Contains(output, "通知付きシェルコマンド") {
		t.Fatalf("expected header in output: %s", output)
	}
	if !strings.Contains(output, "echo hi") {
		t.Fatalf("expected script content in output: %s", output)
	}

	// 空文字の場合は出力しない
	output = captureStdout(func() {
		service.PrintNotificationScript("   ")
	})
	if output != "" {
		t.Fatalf("expected no output for blank script, got: %q", output)
	}
}

func TestIndentCommand(t *testing.T) {
	indented := indentCommand("line1\nline2", "  ")
	expected := "  line1\n  line2"
	if indented != expected {
		t.Fatalf("unexpected indented result: %s", indented)
	}

	if indentCommand("", "  ") != "" {
		t.Fatalf("expected empty string when command is empty")
	}
}

func TestShelled(t *testing.T) {
	if got := shelled("text"); got != "'text'" {
		t.Fatalf("unexpected quoted text: %s", got)
	}
	if got := shelled(""); got != "''" {
		t.Fatalf("expected empty shell quote, got: %s", got)
	}
}

func TestValidateAliasOptionFunc(t *testing.T) {
	cases := []struct {
		name    string
		option  string
		wantErr bool
	}{
		{"clear", "--clear-version-aliases", false},
		{"remove", "--remove-version-aliases=prod", false},
		{"update", "--update-version-aliases=prod=5", false},
		{"empty", "", true},
		{"invalid", "--nope", true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := validateAliasOption(tc.option)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error but got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func captureStdout(fn func()) string {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		panic(err)
	}

	return buf.String()
}
