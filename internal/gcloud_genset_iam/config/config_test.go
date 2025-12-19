package config

import "testing"

type mockFlagParser struct {
	stringVars   map[string]*string
	boolVars     map[string]*bool
	stringValues map[string]string
	boolValues   map[string]bool
	parseError   error
}

func newMockFlagParser() *mockFlagParser {
	return &mockFlagParser{
		stringVars:   make(map[string]*string),
		boolVars:     make(map[string]*bool),
		stringValues: make(map[string]string),
		boolValues:   make(map[string]bool),
	}
}

func (m *mockFlagParser) StringVar(p *string, name string, value string, usage string) {
	if preset, ok := m.stringValues[name]; ok {
		*p = preset
	} else {
		*p = value
	}
	m.stringVars[name] = p
}

func (m *mockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if preset, ok := m.boolValues[name]; ok {
		*p = preset
	} else {
		*p = value
	}
	m.boolVars[name] = p
}

func (m *mockFlagParser) Parse() error {
	return m.parseError
}

func (m *mockFlagParser) Args() []string {
	return nil
}

func (m *mockFlagParser) setString(name, value string) {
	m.stringValues[name] = value
	if ptr, ok := m.stringVars[name]; ok {
		*ptr = value
	}
}

func (m *mockFlagParser) setBool(name string, value bool) {
	m.boolValues[name] = value
	if ptr, ok := m.boolVars[name]; ok {
		*ptr = value
	}
}

func TestParseFlagsWithParser_Success(t *testing.T) {
	t.Run("add iam policy binding to project", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationAddIamPolicyBindingToProject)
		parser.setString("project-id", "sample-project")
		parser.setString("service-account-id", "sa")
		parser.setString("role", "roles/viewer")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.ProjectID != "sample-project" || cfg.ServiceAccountID != "sa" || cfg.Role != "roles/viewer" {
			t.Fatalf("unexpected config: %+v", cfg)
		}
	})

	t.Run("list workload identity pools with flags", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationListWorkloadIdentityPools)
		parser.setString("project-id", "sample-project")
		parser.setString("location", "asia-northeast1")
		parser.setBool("show-deleted", true)
		parser.setString("filter", "displayName:github")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !cfg.ShowDeleted {
			t.Fatalf("show-deleted flag should be true")
		}
		if cfg.Filter != "displayName:github" {
			t.Fatalf("unexpected filter: %s", cfg.Filter)
		}
	})

	t.Run("setup workload identity federation", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationSetupWorkloadIdentityFederation)
		parser.setString("project-id", "sample-project")
		parser.setString("pool-id", "gha-pool")
		parser.setString("provider-id", "gha-provider")
		parser.setString("service-account-id", "gha")
		parser.setString("repository-owner", "landmaster135")
		parser.setString("repository-name", "devbox")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Location != defaultLocation {
			t.Fatalf("location default mismatch: %s", cfg.Location)
		}
	})
}

func TestParseFlagsWithParser_Errors(t *testing.T) {
	t.Run("missing operation", func(t *testing.T) {
		parser := newMockFlagParser()
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when operation is missing")
		}
	})

	t.Run("condition conflict", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationAddIamPolicyBindingToServiceAccount)
		parser.setString("service-account-email", "sa@example")
		parser.setString("member", "user:dev@example.com")
		parser.setString("role", "roles/iam.serviceAccountUser")
		parser.setString("condition", "expr")
		parser.setString("condition-from-file", "file")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for conflicting condition flags")
		}
	})

	t.Run("update service account missing fields", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUpdateServiceAccount)
		parser.setString("service-account-email", "sa@example")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when update fields are missing")
		}
	})

	t.Run("update workload identity pool requires target fields", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUpdateWorkloadIdentityPool)
		parser.setString("project-id", "sample-project")
		parser.setString("pool-id", "gha-pool")
		parser.setString("location", "global")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when update pool has no fields")
		}
	})

	t.Run("update provider requires options", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUpdateOidcWorkloadIdentityPoolProvider)
		parser.setString("project-id", "sample-project")
		parser.setString("pool-id", "gha-pool")
		parser.setString("provider-id", "gha-provider")
		parser.setString("location", "global")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when no provider update options specified")
		}
	})

	t.Run("unsupported operation", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", "unknown-op")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for unsupported operation")
		}
	})
}
