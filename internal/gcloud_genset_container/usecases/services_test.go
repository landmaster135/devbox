package usecases

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_container/config"
)

func TestBuildDeployCloudRunContainerCommand_WithServiceAccountAndNoAuth(t *testing.T) {
	service := NewService()

	command, err := service.BuildDeployCloudRunContainerCommand(DeployCloudRunContainerParams{
		ServiceName:          "app",
		ProjectID:            "my-project",
		Region:               "asia-northeast1",
		Timeout:              "30m",
		RunServiceAccount:    "runner@my-project.iam.gserviceaccount.com",
		AllowUnauthenticated: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud run deploy 'app' --source . --project='my-project' --region='asia-northeast1' --timeout='30m' --service-account='runner@my-project.iam.gserviceaccount.com' --no-allow-unauthenticated"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildDeployCloudRunContainerCommand_DefaultAllowUnauthenticated(t *testing.T) {
	service := NewService()

	command, err := service.BuildDeployCloudRunContainerCommand(DeployCloudRunContainerParams{
		ServiceName:          "svc",
		ProjectID:            "proj",
		Region:               "us-central1",
		Timeout:              "40m",
		AllowUnauthenticated: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(command, "--allow-unauthenticated") {
		t.Fatalf("expected allow-unauthenticated flag: %s", command)
	}
}

func TestBuildUpdateCloudRunContainerEnvCommand(t *testing.T) {
	service := NewService()

	command, err := service.BuildUpdateCloudRunContainerEnvCommand(UpdateCloudRunContainerEnvParams{
		ServiceName: "svc",
		ProjectID:   "proj",
		Region:      "asia-northeast1",
		EnvFile:     "env.yml",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud run deploy 'svc' --image='gcr.io/proj/svc' --region='asia-northeast1' --env-vars-file='env.yml'"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildDeployCloudRunFunctionCommand(t *testing.T) {
	service := NewService()

	command, err := service.BuildDeployCloudRunFunctionCommand(DeployCloudRunFunctionParams{
		FunctionName: "fn",
		Region:       "asia-northeast1",
		EntryPoint:   "Handle",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud functions deploy 'fn' --gen2 --runtime=go122 --region='asia-northeast1' --source=. --entry-point='Handle' --trigger-http --allow-unauthenticated --timeout=180s"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildDeployCloudRunFunctionTriggeredByPubSubCommand_WithEnvVars(t *testing.T) {
	service := NewService()

	command, err := service.BuildDeployCloudRunFunctionTriggeredByPubSubCommand(DeployCloudRunFunctionTriggeredByPubSubParams{
		FunctionName:          "fn",
		ProjectID:             "proj",
		Region:                "us-central1",
		EntryPoint:            "Process",
		TriggerServiceAccount: "svc@proj.iam.gserviceaccount.com",
		TriggerTopic:          "topic-id",
		APIClientID:           "client",
		APIClientSecret:       "secret",
		APIEndpoint:           "https://example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud functions deploy 'fn' --gen2 --runtime=go123 --project='proj' --region='us-central1' --source=. --entry-point='Process' --trigger-service-account='svc@proj.iam.gserviceaccount.com' --trigger-topic='topic-id' --allow-unauthenticated --timeout=180s --set-env-vars='SCRIPT_MANAGER_API_CLIENT_ID=client,SCRIPT_MANAGER_API_CLIENT_SECRET=secret,SCRIPT_MANAGER_API_ENDPOINT=https://example.com'"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildUpdateCloudRunFunctionEnvCommand(t *testing.T) {
	service := NewService()

	command, err := service.BuildUpdateCloudRunFunctionEnvCommand(UpdateCloudRunFunctionEnvParams{
		ServiceName: "fn",
		Region:      "asia-northeast1",
		EnvVars:     "KEY=value,OTHER=value2",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud run services update 'fn' --region='asia-northeast1' --update-env-vars='KEY=value,OTHER=value2'"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildUpdateCloudRunServiceEnvCommand(t *testing.T) {
	service := NewService()

	command, err := service.BuildUpdateCloudRunServiceEnvCommand(UpdateCloudRunServiceEnvParams{
		ServiceName: "svc",
		ProjectID:   "proj",
		Region:      "asia-northeast1",
		EnvFile:     "env.yml",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud run services update 'svc' --project='proj' --region='asia-northeast1' --env-vars-file='env.yml'"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildCreateCloudPubSubTopicCommand(t *testing.T) {
	service := NewService()

	command, err := service.BuildCreateCloudPubSubTopicCommand(CreateCloudPubSubTopicParams{
		TopicName:                "topic",
		MessageRetentionDuration: "2d",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud pubsub topics create 'topic' --message-retention-duration='2d'"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildListCloudPubSubTopicsCommand(t *testing.T) {
	service := NewService()

	command, err := service.BuildListCloudPubSubTopicsCommand(ListCloudPubSubTopicsParams{TopicName: "sample"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud pubsub topics list --filter=\"name.scope(topic):'sample'\""
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildListCloudPubSubSubscriptionsCommand_WithFilterAndURI(t *testing.T) {
	service := NewService()

	command, err := service.BuildListCloudPubSubSubscriptionsCommand(ListCloudPubSubSubscriptionsParams{
		SubscriptionName: "sub",
		ShowURI:          true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud pubsub subscriptions list --filter=\"name.scope(subscription):'sub'\" --uri"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildCreateCloudPubSubSubscriptionCommand_DefaultEndpoint(t *testing.T) {
	service := NewService()

	command, err := service.BuildCreateCloudPubSubSubscriptionCommand(CreateCloudPubSubSubscriptionParams{
		SubscriptionName:         "sub",
		TopicName:                "topic",
		TopicProject:             "proj",
		PushServiceAccount:       "svc@proj.iam.gserviceaccount.com",
		MessageRetentionDuration: "1d",
		ExpirationPeriod:         "never",
		MaxRetryDelay:            "600s",
		MinRetryDelay:            "10s",
		AckDeadline:              "600",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud pubsub subscriptions create 'sub' --topic='topic' --topic-project='proj' --message-retention-duration='1d' --push-auth-service-account='svc@proj.iam.gserviceaccount.com' --push-endpoint='https://proj.appspot.com/sub' --expiration-period='never' --max-retry-delay='600s' --min-retry-delay='10s' --ack-deadline='600'"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildDeleteCloudPubSubSubscriptionsCommand(t *testing.T) {
	service := NewService()

	command, err := service.BuildDeleteCloudPubSubSubscriptionsCommand(DeleteCloudPubSubSubscriptionsParams{
		SubscriptionNames: []string{"sub1", "sub2"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud pubsub subscriptions delete 'sub1' && gcloud pubsub subscriptions delete 'sub2'"
	if command != expected {
		t.Fatalf("unexpected command: %q", command)
	}
}

func TestBuildDeleteCloudPubSubTopicsCommand(t *testing.T) {
	service := NewService()

	command, err := service.BuildDeleteCloudPubSubTopicsCommand(DeleteCloudPubSubTopicsParams{
		TopicNames: []string{"topic1"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud pubsub topics delete 'topic1'"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildDeleteCloudPubSubSubscriptionsAndTopicsCommand(t *testing.T) {
	service := NewService()

	command, err := service.BuildDeleteCloudPubSubSubscriptionsAndTopicsCommand(DeleteCloudPubSubSubscriptionsAndTopicsParams{
		SubscriptionNames: []string{"sub1"},
		TopicNames:        []string{"topic1"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud pubsub subscriptions delete 'sub1' && gcloud pubsub topics delete 'topic1'"
	if command != expected {
		t.Fatalf("unexpected command: %q", command)
	}
}

func TestBuildDeleteCloudRunFunctionCommand(t *testing.T) {
	service := NewService()

	command, err := service.BuildDeleteCloudRunFunctionCommand(DeleteCloudRunFunctionParams{
		ServiceName: "fn",
		Region:      "us-central1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud run services delete 'fn' --region='us-central1'"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildDeleteCloudRunFunctionCommand_MissingParams(t *testing.T) {
	service := NewService()

	if _, err := service.BuildDeleteCloudRunFunctionCommand(DeleteCloudRunFunctionParams{}); err == nil {
		t.Fatal("expected error when parameters are missing")
	}
}

func TestBuildCommandVariants(t *testing.T) {
	service := NewService()

	t.Run("deploy cloud run container", func(t *testing.T) {
		expected, err := service.BuildDeployCloudRunContainerCommand(DeployCloudRunContainerParams{
			ServiceName:          "svc",
			ProjectID:            "proj",
			Region:               "asia-northeast1",
			Timeout:              "45m",
			RunServiceAccount:    "svc@proj.iam.gserviceaccount.com",
			AllowUnauthenticated: false,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		command, err := service.BuildCommand(&cfg.Config{
			Operation:            cfg.OperationDeployCloudRunContainer,
			ServiceName:          "svc",
			ProjectID:            "proj",
			Region:               "asia-northeast1",
			Timeout:              "45m",
			RunServiceAccount:    "svc@proj.iam.gserviceaccount.com",
			AllowUnauthenticated: false,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})

	t.Run("deploy cloud run container missing params", func(t *testing.T) {
		if _, err := service.BuildDeployCloudRunContainerCommand(DeployCloudRunContainerParams{ProjectID: "proj"}); err == nil {
			t.Fatal("expected error for missing parameters")
		}
	})

	t.Run("update cloud run container env", func(t *testing.T) {
		expected, err := service.BuildUpdateCloudRunContainerEnvCommand(UpdateCloudRunContainerEnvParams{
			ServiceName: "svc",
			ProjectID:   "proj",
			Region:      "asia-northeast1",
			EnvFile:     "env.yaml",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		command, err := service.BuildCommand(&cfg.Config{
			Operation:   cfg.OperationUpdateCloudRunContainerEnv,
			ServiceName: "svc",
			ProjectID:   "proj",
			Region:      "asia-northeast1",
			EnvFile:     "env.yaml",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})

	t.Run("update cloud run container env missing params", func(t *testing.T) {
		if _, err := service.BuildUpdateCloudRunContainerEnvCommand(UpdateCloudRunContainerEnvParams{ServiceName: "svc"}); err == nil {
			t.Fatal("expected error for missing parameters")
		}
	})

	t.Run("deploy cloud run function", func(t *testing.T) {
		expected, err := service.BuildDeployCloudRunFunctionCommand(DeployCloudRunFunctionParams{
			FunctionName: "fn",
			Region:       "asia-northeast1",
			EntryPoint:   "Handle",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		command, err := service.BuildCommand(&cfg.Config{
			Operation:    cfg.OperationDeployCloudRunFunction,
			FunctionName: "fn",
			Region:       "asia-northeast1",
			EntryPoint:   "Handle",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})

	t.Run("deploy cloud run function missing params", func(t *testing.T) {
		if _, err := service.BuildDeployCloudRunFunctionCommand(DeployCloudRunFunctionParams{FunctionName: "fn"}); err == nil {
			t.Fatal("expected error when region or entry point missing")
		}
	})

	t.Run("deploy cloud run function triggered by pubsub", func(t *testing.T) {
		expected, err := service.BuildDeployCloudRunFunctionTriggeredByPubSubCommand(DeployCloudRunFunctionTriggeredByPubSubParams{
			FunctionName:          "fn",
			ProjectID:             "proj",
			Region:                "asia-northeast1",
			EntryPoint:            "Process",
			TriggerServiceAccount: "svc@proj.iam.gserviceaccount.com",
			TriggerTopic:          "topic",
			APIClientID:           "client",
			APIClientSecret:       "secret",
			APIEndpoint:           "https://example.com",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		command, err := service.BuildCommand(&cfg.Config{
			Operation:             cfg.OperationDeployCloudRunFunctionTriggeredByPubSub,
			FunctionName:          "fn",
			ProjectID:             "proj",
			Region:                "asia-northeast1",
			EntryPoint:            "Process",
			TriggerServiceAccount: "svc@proj.iam.gserviceaccount.com",
			TriggerTopic:          "topic",
			APIClientID:           "client",
			APIClientSecret:       "secret",
			APIEndpoint:           "https://example.com",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})

	t.Run("deploy cloud run function triggered by pubsub missing params", func(t *testing.T) {
		if _, err := service.BuildDeployCloudRunFunctionTriggeredByPubSubCommand(DeployCloudRunFunctionTriggeredByPubSubParams{FunctionName: "fn"}); err == nil {
			t.Fatal("expected error when required parameters are missing")
		}
	})

	t.Run("update cloud run function env", func(t *testing.T) {
		expected, err := service.BuildUpdateCloudRunFunctionEnvCommand(UpdateCloudRunFunctionEnvParams{
			ServiceName: "fn",
			Region:      "asia-northeast1",
			EnvVars:     "KEY=VALUE",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		command, err := service.BuildCommand(&cfg.Config{
			Operation:   cfg.OperationUpdateCloudRunFunctionEnv,
			ServiceName: "fn",
			Region:      "asia-northeast1",
			EnvVars:     "KEY=VALUE",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})

	t.Run("update cloud run function env missing params", func(t *testing.T) {
		if _, err := service.BuildUpdateCloudRunFunctionEnvCommand(UpdateCloudRunFunctionEnvParams{Region: "asia"}); err == nil {
			t.Fatal("expected error when service name or env vars missing")
		}
	})

	t.Run("update cloud run service env", func(t *testing.T) {
		expected, err := service.BuildUpdateCloudRunServiceEnvCommand(UpdateCloudRunServiceEnvParams{
			ServiceName: "svc",
			ProjectID:   "proj",
			Region:      "asia-northeast1",
			EnvFile:     "env.yaml",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		command, err := service.BuildCommand(&cfg.Config{
			Operation:   cfg.OperationUpdateCloudRunServiceEnv,
			ServiceName: "svc",
			ProjectID:   "proj",
			Region:      "asia-northeast1",
			EnvFile:     "env.yaml",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})

	t.Run("update cloud run service env missing params", func(t *testing.T) {
		if _, err := service.BuildUpdateCloudRunServiceEnvCommand(UpdateCloudRunServiceEnvParams{ServiceName: "svc"}); err == nil {
			t.Fatal("expected error when project or region missing")
		}
	})

	t.Run("create pubsub topic", func(t *testing.T) {
		expected, err := service.BuildCreateCloudPubSubTopicCommand(CreateCloudPubSubTopicParams{
			TopicName:                "topic",
			MessageRetentionDuration: "2d",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		command, err := service.BuildCommand(&cfg.Config{
			Operation:                cfg.OperationCreateCloudPubSubTopic,
			TopicName:                "topic",
			MessageRetentionDuration: "2d",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})

	t.Run("create pubsub topic missing params", func(t *testing.T) {
		if _, err := service.BuildCreateCloudPubSubTopicCommand(CreateCloudPubSubTopicParams{}); err == nil {
			t.Fatal("expected error when parameters missing")
		}
	})

	t.Run("list pubsub topics", func(t *testing.T) {
		expected, err := service.BuildListCloudPubSubTopicsCommand(ListCloudPubSubTopicsParams{TopicName: "topic"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		command, err := service.BuildCommand(&cfg.Config{
			Operation: cfg.OperationListCloudPubSubTopics,
			TopicName: "topic",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})

	t.Run("create pubsub subscription missing params", func(t *testing.T) {
		if _, err := service.BuildCreateCloudPubSubSubscriptionCommand(CreateCloudPubSubSubscriptionParams{SubscriptionName: "sub"}); err == nil {
			t.Fatal("expected error when required parameters missing")
		}
	})

	t.Run("list pubsub subscriptions", func(t *testing.T) {
		expected, err := service.BuildListCloudPubSubSubscriptionsCommand(ListCloudPubSubSubscriptionsParams{
			SubscriptionName: "sub",
			ShowURI:          true,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		command, err := service.BuildCommand(&cfg.Config{
			Operation:        cfg.OperationListCloudPubSubSubscriptions,
			SubscriptionName: "sub",
			ShowURI:          true,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})

	t.Run("delete pubsub subscriptions and topics missing params", func(t *testing.T) {
		if _, err := service.BuildDeleteCloudPubSubSubscriptionsAndTopicsCommand(DeleteCloudPubSubSubscriptionsAndTopicsParams{}); err == nil {
			t.Fatal("expected error when nothing specified")
		}
	})

	t.Run("create pubsub subscription", func(t *testing.T) {
		expected, err := service.BuildCreateCloudPubSubSubscriptionCommand(CreateCloudPubSubSubscriptionParams{
			SubscriptionName:         "sub",
			TopicName:                "topic",
			TopicProject:             "proj",
			PushEndpoint:             "https://custom",
			PushServiceAccount:       "svc@proj",
			MessageRetentionDuration: "2d",
			ExpirationPeriod:         "never",
			MaxRetryDelay:            "600s",
			MinRetryDelay:            "10s",
			AckDeadline:              "120",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		command, err := service.BuildCommand(&cfg.Config{
			Operation:                cfg.OperationCreateCloudPubSubSubscription,
			SubscriptionName:         "sub",
			TopicName:                "topic",
			TopicProject:             "proj",
			PushEndpoint:             "https://custom",
			PushServiceAccount:       "svc@proj",
			MessageRetentionDuration: "2d",
			ExpirationPeriod:         "never",
			MaxRetryDelay:            "600s",
			MinRetryDelay:            "10s",
			AckDeadline:              "120",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})

	t.Run("delete pubsub subscriptions missing params", func(t *testing.T) {
		if _, err := service.BuildDeleteCloudPubSubSubscriptionsCommand(DeleteCloudPubSubSubscriptionsParams{}); err == nil {
			t.Fatal("expected error when subscription names missing")
		}
	})

	t.Run("delete pubsub subscriptions and topics", func(t *testing.T) {
		expected, err := service.BuildDeleteCloudPubSubSubscriptionsAndTopicsCommand(DeleteCloudPubSubSubscriptionsAndTopicsParams{
			SubscriptionNames: []string{"sub1"},
			TopicNames:        []string{"topic1"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		command, err := service.BuildCommand(&cfg.Config{
			Operation:         cfg.OperationDeleteCloudPubSubSubscriptionsAndTopics,
			SubscriptionNames: []string{"sub1"},
			TopicNames:        []string{"topic1"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})

	t.Run("delete pubsub topics missing params", func(t *testing.T) {
		if _, err := service.BuildDeleteCloudPubSubTopicsCommand(DeleteCloudPubSubTopicsParams{}); err == nil {
			t.Fatal("expected error when topic names missing")
		}
	})

	t.Run("delete pubsub subscriptions", func(t *testing.T) {
		expected, err := service.BuildDeleteCloudPubSubSubscriptionsCommand(DeleteCloudPubSubSubscriptionsParams{SubscriptionNames: []string{"sub1", "sub2"}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		command, err := service.BuildCommand(&cfg.Config{
			Operation:         cfg.OperationDeleteCloudPubSubSubscriptions,
			SubscriptionNames: []string{"sub1", "sub2"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})

	t.Run("delete cloud run function missing params", func(t *testing.T) {
		if _, err := service.BuildDeleteCloudRunFunctionCommand(DeleteCloudRunFunctionParams{ServiceName: "svc"}); err == nil {
			t.Fatal("expected error when region missing")
		}
	})

	t.Run("delete pubsub topics", func(t *testing.T) {
		expected, err := service.BuildDeleteCloudPubSubTopicsCommand(DeleteCloudPubSubTopicsParams{TopicNames: []string{"topic1"}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		command, err := service.BuildCommand(&cfg.Config{
			Operation:  cfg.OperationDeleteCloudPubSubTopics,
			TopicNames: []string{"topic1"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})

	t.Run("delete cloud run function", func(t *testing.T) {
		expected, err := service.BuildDeleteCloudRunFunctionCommand(DeleteCloudRunFunctionParams{
			ServiceName: "svc",
			Region:      "asia",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		command, err := service.BuildCommand(&cfg.Config{
			Operation:   cfg.OperationDeleteCloudRunFunction,
			ServiceName: "svc",
			Region:      "asia",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if command != expected {
			t.Fatalf("unexpected command: %s", command)
		}
	})
}

func TestBuildNotificationWrappedCommand_DeployCloudRunContainer(t *testing.T) {
	service := NewService()

	script, ok := service.BuildNotificationWrappedCommand(cfg.OperationDeployCloudRunContainer, "gcloud run deploy 'svc'")
	if !ok {
		t.Fatal("expected notification script to be generated")
	}

	if !strings.Contains(script, "コンテナをデプロイするよ！") {
		t.Fatalf("start notification missing: %s", script)
	}
	if !strings.Contains(script, "if gcloud run deploy 'svc'; then") {
		t.Fatalf("gcloud execution block missing: %s", script)
	}
	if !strings.Contains(script, "google-cloud-run-success") {
		t.Fatalf("success embed type missing: %s", script)
	}
	if !strings.Contains(script, "google-cloud-run-failed") {
		t.Fatalf("failure embed type missing: %s", script)
	}
}

func TestBuildNotificationWrappedCommand_NoTemplate(t *testing.T) {
	service := NewService()

	if _, ok := service.BuildNotificationWrappedCommand(cfg.OperationCreateCloudPubSubTopic, "gcloud pubsub topics create 'topic'"); ok {
		t.Fatal("expected notification template to be unavailable")
	}
}

func TestPrintHighlightedCommand(t *testing.T) {
	service := NewService()

	output := captureStdout(func() {
		service.PrintHighlightedCommand("gcloud run deploy 'svc'")
	})

	if !strings.Contains(output, "生成された gcloud コマンド") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "gcloud run deploy 'svc'") {
		t.Fatalf("command not printed: %s", output)
	}
}

func TestPrintNotificationScript(t *testing.T) {
	service := NewService()

	emptyOutput := captureStdout(func() {
		service.PrintNotificationScript("  ")
	})
	if strings.TrimSpace(emptyOutput) != "" {
		t.Fatalf("expected no output, got %q", emptyOutput)
	}

	scriptOutput := captureStdout(func() {
		service.PrintNotificationScript("echo hello")
	})
	if !strings.Contains(scriptOutput, "通知付きシェルコマンド") || !strings.Contains(scriptOutput, "echo hello") {
		t.Fatalf("unexpected script output: %s", scriptOutput)
	}
}

func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestIndentCommand(t *testing.T) {
	if indentCommand("", "  ") != "" {
		t.Fatal("expected empty command to remain empty")
	}

	result := indentCommand("echo hi\nls", "  ")
	expected := "  echo hi\n  ls"
	if result != expected {
		t.Fatalf("unexpected indentation: %s", result)
	}
}

func TestBuildCommand_Dispatch(t *testing.T) {
	service := NewService()

	conf := &cfg.Config{
		Operation:                cfg.OperationCreateCloudPubSubTopic,
		TopicName:                "topic",
		MessageRetentionDuration: "1d",
	}

	command, err := service.BuildCommand(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud pubsub topics create 'topic' --message-retention-duration='1d'"
	if command != expected {
		t.Fatalf("unexpected command: %s", command)
	}
}

func TestBuildCommand_UnsupportedOperation(t *testing.T) {
	service := NewService()

	_, err := service.BuildCommand(&cfg.Config{Operation: "unknown"})
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}
