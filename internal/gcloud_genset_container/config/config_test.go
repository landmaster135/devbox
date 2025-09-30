package config

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

type mockFlagParser struct {
	stringVars   map[string]*string
	intVars      map[string]*int
	boolVars     map[string]*bool
	stringValues map[string]string
	intValues    map[string]int
	boolValues   map[string]bool
	parseError   error
	args         []string
}

func newMockFlagParser() *mockFlagParser {
	return &mockFlagParser{
		stringVars:   make(map[string]*string),
		intVars:      make(map[string]*int),
		boolVars:     make(map[string]*bool),
		stringValues: make(map[string]string),
		intValues:    make(map[string]int),
		boolValues:   make(map[string]bool),
	}
}

func (m *mockFlagParser) StringVar(p *string, name string, value string, usage string) {
	if v, ok := m.stringValues[name]; ok {
		*p = v
	} else {
		*p = value
	}
	m.stringVars[name] = p
}

func (m *mockFlagParser) IntVar(p *int, name string, value int, usage string) {
	if v, ok := m.intValues[name]; ok {
		*p = v
	} else {
		*p = value
	}
	m.intVars[name] = p
}

func (m *mockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if v, ok := m.boolValues[name]; ok {
		*p = v
	} else {
		*p = value
	}
	m.boolVars[name] = p
}

func (m *mockFlagParser) Parse() error {
	return m.parseError
}

func (m *mockFlagParser) Args() []string {
	if len(m.args) == 0 {
		return []string{}
	}
	return append([]string(nil), m.args...)
}

func (m *mockFlagParser) setString(name, value string) {
	m.stringValues[name] = value
	if ptr, ok := m.stringVars[name]; ok {
		*ptr = value
	}
}

func (m *mockFlagParser) setInt(name string, value int) {
	m.intValues[name] = value
	if ptr, ok := m.intVars[name]; ok {
		*ptr = value
	}
}

func (m *mockFlagParser) setBool(name string, value bool) {
	m.boolValues[name] = value
	if ptr, ok := m.boolVars[name]; ok {
		*ptr = value
	}
}

func (m *mockFlagParser) setArgs(args []string) {
	m.args = append([]string(nil), args...)
}

func TestParseFlagsWithParser_SuccessCases(t *testing.T) {
	t.Run("deploy cloud run container with defaults", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeployCloudRunContainer)
		parser.setString("service-name", " my-service ")
		parser.setString("project-id", "my-project ")
		parser.setString("timeout", "")
		parser.setBool("allow-unauthenticated", false)

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.ServiceName != "my-service" {
			t.Fatalf("expected trimmed service name, got %q", cfg.ServiceName)
		}
		if cfg.ProjectID != "my-project" {
			t.Fatalf("expected trimmed project id, got %q", cfg.ProjectID)
		}
		if cfg.Region != defaultRegion {
			t.Fatalf("expected default region %q, got %q", defaultRegion, cfg.Region)
		}
		if cfg.AllowUnauthenticated {
			t.Fatalf("expected allow-unauthenticated=false")
		}
		if cfg.Timeout != defaultContainerTimeout {
			t.Fatalf("expected timeout %s, got %s", defaultContainerTimeout, cfg.Timeout)
		}
	})

	t.Run("deploy function triggered by pubsub defaults", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeployCloudRunFunctionTriggeredByPubSub)
		parser.setString("function-name", "fn")
		parser.setString("project-id", "proj")
		parser.setString("trigger-service-account", "svc@proj.iam.gserviceaccount.com")
		parser.setString("trigger-topic", "topic")
		parser.setString("api-client-id", "client")
		parser.setString("api-client-secret", "secret")
		parser.setString("api-endpoint", "https://example.com")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Region != defaultRegion {
			t.Fatalf("expected region %q, got %q", defaultRegion, cfg.Region)
		}
		if cfg.EntryPoint != defaultPubSubEntryPoint {
			t.Fatalf("expected entry point %q, got %q", defaultPubSubEntryPoint, cfg.EntryPoint)
		}
		if cfg.APIClientID != "client" || cfg.APIClientSecret != "secret" || cfg.APIEndpoint != "https://example.com" {
			t.Fatalf("unexpected api config: %+v", cfg)
		}
	})

	t.Run("delete subscriptions and topics parses lists", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeleteCloudPubSubSubscriptionsAndTopics)
		parser.setString("subscription-names", " sub1 , sub2 , ")
		parser.setString("topic-names", " topic1 ")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(cfg.SubscriptionNames) != 2 || cfg.SubscriptionNames[0] != "sub1" || cfg.SubscriptionNames[1] != "sub2" {
			t.Fatalf("unexpected subscription names: %#v", cfg.SubscriptionNames)
		}
		if len(cfg.TopicNames) != 1 || cfg.TopicNames[0] != "topic1" {
			t.Fatalf("unexpected topic names: %#v", cfg.TopicNames)
		}
	})

	t.Run("create pubsub subscription", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationCreateCloudPubSubSubscription)
		parser.setString("subscription-name", "sub")
		parser.setString("topic-name", "topic")
		parser.setString("topic-project", "proj")
		parser.setString("push-service-account", "svc@proj")
		parser.setString("push-endpoint", "https://custom")
		parser.setString("message-retention-duration", "2d")
		parser.setString("expiration-period", "never")
		parser.setString("max-retry-delay", "600s")
		parser.setString("min-retry-delay", "10s")
		parser.setString("ack-deadline", "120")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.PushEndpoint != "https://custom" {
			t.Fatalf("unexpected push endpoint: %s", cfg.PushEndpoint)
		}
	})

	t.Run("create pubsub subscription defaults", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationCreateCloudPubSubSubscription)
		parser.setString("subscription-name", "sub")
		parser.setString("topic-name", "topic")
		parser.setString("topic-project", "proj")
		parser.setString("push-service-account", "svc@proj")
		parser.setString("message-retention-duration", "")
		parser.setString("expiration-period", "")
		parser.setString("max-retry-delay", "")
		parser.setString("min-retry-delay", "")
		parser.setString("ack-deadline", "")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.MessageRetentionDuration != defaultMessageRetentionDuration {
			t.Fatalf("expected default retention %s, got %s", defaultMessageRetentionDuration, cfg.MessageRetentionDuration)
		}
		if cfg.ExpirationPeriod != defaultExpirationPeriod {
			t.Fatalf("expected default expiration %s, got %s", defaultExpirationPeriod, cfg.ExpirationPeriod)
		}
		if cfg.MaxRetryDelay != defaultMaxRetryDelay || cfg.MinRetryDelay != defaultMinRetryDelay || cfg.AckDeadline != defaultAckDeadline {
			t.Fatalf("unexpected defaults: max=%s min=%s ack=%s", cfg.MaxRetryDelay, cfg.MinRetryDelay, cfg.AckDeadline)
		}
	})

	t.Run("update cloud run container env", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUpdateCloudRunContainerEnv)
		parser.setString("service-name", "svc")
		parser.setString("project-id", "proj")
		parser.setString("region", "asia-northeast1")
		parser.setString("env-file", "env.yaml")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.EnvFile != "env.yaml" {
			t.Fatalf("unexpected env file: %s", cfg.EnvFile)
		}
	})

	t.Run("deploy cloud run function", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeployCloudRunFunction)
		parser.setString("function-name", "fn")
		parser.setString("region", "asia-northeast1")
		parser.setString("entry-point", "Handle")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.EntryPoint != "Handle" {
			t.Fatalf("unexpected entry point: %s", cfg.EntryPoint)
		}
	})

	t.Run("update cloud run function env", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUpdateCloudRunFunctionEnv)
		parser.setString("service-name", "fn")
		parser.setString("region", "asia-northeast1")
		parser.setString("env-vars", "KEY=VALUE")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.EnvVars != "KEY=VALUE" {
			t.Fatalf("unexpected env vars: %s", cfg.EnvVars)
		}
	})

	t.Run("update cloud run service env uses default env file", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUpdateCloudRunServiceEnv)
		parser.setString("service-name", "svc")
		parser.setString("project-id", "proj")
		parser.setString("region", "asia")
		parser.setString("env-file", "")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.EnvFile != defaultEnvFile {
			t.Fatalf("expected default env file %s, got %s", defaultEnvFile, cfg.EnvFile)
		}
	})

	t.Run("create pubsub topic defaults", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationCreateCloudPubSubTopic)
		parser.setString("topic-name", "topic")
		parser.setString("message-retention-duration", "")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.MessageRetentionDuration != defaultMessageRetentionDuration {
			t.Fatalf("expected default retention %s, got %s", defaultMessageRetentionDuration, cfg.MessageRetentionDuration)
		}
	})

	t.Run("list pubsub topics", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationListCloudPubSubTopics)
		parser.setString("topic-name", "topic")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.TopicName != "topic" {
			t.Fatalf("unexpected topic name: %s", cfg.TopicName)
		}
	})

	t.Run("list pubsub subscriptions with uri", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationListCloudPubSubSubscriptions)
		parser.setString("subscription-name", "sub")
		parser.setBool("show-uri", true)

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !cfg.ShowURI {
			t.Fatalf("expected show uri to be true")
		}
	})

	t.Run("delete pubsub subscriptions", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeleteCloudPubSubSubscriptions)
		parser.setString("subscription-names", "sub1,sub2")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(cfg.SubscriptionNames) != 2 {
			t.Fatalf("unexpected subscriptions: %#v", cfg.SubscriptionNames)
		}
	})

	t.Run("delete pubsub topics", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeleteCloudPubSubTopics)
		parser.setString("topic-names", "topic1")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(cfg.TopicNames) != 1 || cfg.TopicNames[0] != "topic1" {
			t.Fatalf("unexpected topics: %#v", cfg.TopicNames)
		}
	})

	t.Run("delete cloud run function", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeleteCloudRunFunction)
		parser.setString("service-name", "svc")
		parser.setString("region", "asia")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Region != "asia" {
			t.Fatalf("unexpected region: %s", cfg.Region)
		}
	})

	t.Run("help flag bypasses validation", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setBool("help", true)

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !cfg.Help {
			t.Fatalf("expected help flag to be true")
		}
	})
}

func TestParseFlagsWithParser_Errors(t *testing.T) {
	t.Run("parser failure", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.parseError = errors.New("parse error")

		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected parse error")
		}
	})

	t.Run("missing operation", func(t *testing.T) {
		parser := newMockFlagParser()
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when operation is missing")
		}
	})

	t.Run("invalid operation", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", "unknown")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for invalid operation")
		}
	})

	t.Run("missing required param for deploy cloud run container", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeployCloudRunContainer)
		parser.setString("project-id", "proj")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when service-name is missing")
		}
	})

	t.Run("missing project id for deploy cloud run container", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeployCloudRunContainer)
		parser.setString("service-name", "svc")
		parser.setString("region", "asia")
		parser.setString("timeout", "45m")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when project-id is missing")
		}
	})

	t.Run("missing entry point for deploy cloud run function", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeployCloudRunFunction)
		parser.setString("function-name", "fn")
		parser.setString("region", "asia")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when entry-point is missing")
		}
	})

	t.Run("missing region for deploy cloud run function", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeployCloudRunFunction)
		parser.setString("function-name", "fn")
		parser.setString("entry-point", "Handle")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when region is missing")
		}
	})

	t.Run("missing env vars for update cloud run function env", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUpdateCloudRunFunctionEnv)
		parser.setString("service-name", "fn")
		parser.setString("region", "asia")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when env-vars is missing")
		}
	})

	t.Run("missing project for update cloud run service env", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUpdateCloudRunServiceEnv)
		parser.setString("service-name", "svc")
		parser.setString("region", "asia")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when project-id is missing")
		}
	})

	t.Run("missing service name for update cloud run container env", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUpdateCloudRunContainerEnv)
		parser.setString("project-id", "proj")
		parser.setString("region", "asia")
		parser.setString("env-file", "env.yaml")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when service-name is missing")
		}
	})

	t.Run("missing project for update cloud run container env", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUpdateCloudRunContainerEnv)
		parser.setString("service-name", "svc")
		parser.setString("region", "asia")
		parser.setString("env-file", "env.yaml")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when project-id is missing for update container env")
		}
	})

	t.Run("missing region for update cloud run container env", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUpdateCloudRunContainerEnv)
		parser.setString("service-name", "svc")
		parser.setString("project-id", "proj")
		parser.setString("env-file", "env.yaml")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when region is missing for update container env")
		}
	})

	t.Run("missing region for delete cloud run function", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeleteCloudRunFunction)
		parser.setString("service-name", "svc")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when region is missing")
		}
	})

	t.Run("missing trigger service account", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeployCloudRunFunctionTriggeredByPubSub)
		parser.setString("function-name", "fn")
		parser.setString("project-id", "proj")
		parser.setString("trigger-topic", "topic")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when trigger-service-account is missing")
		}
	})

	t.Run("missing trigger topic", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeployCloudRunFunctionTriggeredByPubSub)
		parser.setString("function-name", "fn")
		parser.setString("project-id", "proj")
		parser.setString("trigger-service-account", "svc@proj")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when trigger-topic is missing")
		}
	})

	t.Run("missing function name for pubsub triggered function", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeployCloudRunFunctionTriggeredByPubSub)
		parser.setString("project-id", "proj")
		parser.setString("trigger-service-account", "svc@proj")
		parser.setString("trigger-topic", "topic")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when function-name is missing")
		}
	})

	t.Run("missing required param for pubsub subscription", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationCreateCloudPubSubSubscription)
		parser.setString("topic-name", "topic")
		parser.setString("topic-project", "proj")
		parser.setString("push-service-account", "svc@proj")
		parser.setString("message-retention-duration", "1d")
		parser.setString("expiration-period", "never")
		parser.setString("max-retry-delay", "600s")
		parser.setString("min-retry-delay", "10s")
		parser.setString("ack-deadline", "120")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when subscription-name is missing")
		}
	})

	t.Run("missing topic for pubsub subscription", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationCreateCloudPubSubSubscription)
		parser.setString("subscription-name", "sub")
		parser.setString("topic-project", "proj")
		parser.setString("push-service-account", "svc@proj")
		parser.setString("message-retention-duration", "1d")
		parser.setString("expiration-period", "never")
		parser.setString("max-retry-delay", "600s")
		parser.setString("min-retry-delay", "10s")
		parser.setString("ack-deadline", "120")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when topic-name is missing")
		}
	})

	t.Run("missing topic project for pubsub subscription", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationCreateCloudPubSubSubscription)
		parser.setString("subscription-name", "sub")
		parser.setString("topic-name", "topic")
		parser.setString("push-service-account", "svc@proj")
		parser.setString("message-retention-duration", "1d")
		parser.setString("expiration-period", "never")
		parser.setString("max-retry-delay", "600s")
		parser.setString("min-retry-delay", "10s")
		parser.setString("ack-deadline", "120")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when topic-project is missing")
		}
	})

	t.Run("missing push service account for pubsub subscription", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationCreateCloudPubSubSubscription)
		parser.setString("subscription-name", "sub")
		parser.setString("topic-name", "topic")
		parser.setString("topic-project", "proj")
		parser.setString("message-retention-duration", "1d")
		parser.setString("expiration-period", "never")
		parser.setString("max-retry-delay", "600s")
		parser.setString("min-retry-delay", "10s")
		parser.setString("ack-deadline", "120")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when push-service-account is missing")
		}
	})

	t.Run("invalid comma separated values", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeleteCloudPubSubSubscriptions)
		parser.setString("subscription-names", " , ")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when subscription-names is invalid")
		}
	})

	t.Run("missing subscription and topic names", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeleteCloudPubSubSubscriptionsAndTopics)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when no targets provided")
		}
	})

	t.Run("missing topic names for delete topics", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeleteCloudPubSubTopics)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when topic-names is missing")
		}
	})

	t.Run("missing topic name for create pubsub topic", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationCreateCloudPubSubTopic)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when topic-name is missing")
		}
	})

	t.Run("missing topic name for list pubsub topics", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationListCloudPubSubTopics)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when topic-name is missing for list")
		}
	})

	t.Run("missing service name for delete cloud run function", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDeleteCloudRunFunction)
		parser.setString("region", "asia")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when service-name is missing")
		}
	})
}

func TestParseFlags_RealParser(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{
		"cmd",
		"-operation", OperationDeployCloudRunContainer,
		"-service-name", "svc",
		"-project-id", "proj",
	}

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Operation != OperationDeployCloudRunContainer {
		t.Fatalf("unexpected operation: %s", cfg.Operation)
	}
	if cfg.Region != defaultRegion {
		t.Fatalf("expected default region %q, got %q", defaultRegion, cfg.Region)
	}
}

func TestPrintUsage(t *testing.T) {
	output := captureStderr(func() {
		PrintUsage()
	})

	if !strings.Contains(output, "Cloud Run / Cloud Functions") {
		t.Fatalf("usage output missing expected text: %s", output)
	}
}

func captureStderr(fn func()) string {
	old := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stderr = w
	fn()
	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestStandardFlagParser(t *testing.T) {
	parser := NewStandardFlagParser()
	if parser == nil {
		t.Fatal("expected parser to be created")
	}

	var name string
	var num int
	var flag bool

	parser.StringVar(&name, "name", "default", "")
	parser.IntVar(&num, "num", 1, "")
	parser.BoolVar(&flag, "flag", false, "")

	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"cmd", "-name", "value", "-num", "3", "-flag", "extra"}

	if err := parser.Parse(); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if name != "value" {
		t.Fatalf("expected name=\"value\", got %q", name)
	}
	if num != 3 {
		t.Fatalf("expected num=3, got %d", num)
	}
	if !flag {
		t.Fatal("expected flag to be true")
	}

	args := parser.Args()
	if len(args) != 1 || args[0] != "extra" {
		t.Fatalf("unexpected remaining args: %#v", args)
	}
}
