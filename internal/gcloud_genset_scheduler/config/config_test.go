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
    boolVars     map[string]*bool
    stringValues map[string]string
    boolValues   map[string]bool
    parseErr     error
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
    if val, ok := m.stringValues[name]; ok {
        *p = val
    } else {
        *p = value
    }
    m.stringVars[name] = p
}

func (m *mockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
    if val, ok := m.boolValues[name]; ok {
        *p = val
    } else {
        *p = value
    }
    m.boolVars[name] = p
}

func (m *mockFlagParser) Parse() error {
    return m.parseErr
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
    t.Parallel()

    t.Run("create-pubsub-job", func(t *testing.T) {
        t.Parallel()
        parser := newMockFlagParser()
        parser.setString("operation", OperationCreatePubSubJob)
        parser.setString("job-name", " job-one ")
        parser.setString("project-id", " sample-project ")
        parser.setString("pubsub-topic", " topic ")
        parser.setString("schedule", " * * * * * ")

        cfg, err := ParseFlagsWithParser(parser)
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if cfg.JobName != "job-one" {
            t.Fatalf("expected trimmed job name, got %q", cfg.JobName)
        }
        if cfg.Schedule != "* * * * *" {
            t.Fatalf("expected trimmed schedule, got %q", cfg.Schedule)
        }
    })

    t.Run("create-http-job", func(t *testing.T) {
        t.Parallel()
        parser := newMockFlagParser()
        parser.setString("operation", OperationCreateHTTPJob)
        parser.setString("job-name", "http-job")
        parser.setString("project-id", "proj")
        parser.setString("http-method", "post")
        parser.setString("service-url", " https://example.com ")

        cfg, err := ParseFlagsWithParser(parser)
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if cfg.HTTPMethod != "POST" {
            t.Fatalf("expected method to be upper case, got %q", cfg.HTTPMethod)
        }
        if cfg.ServiceURL != "https://example.com" {
            t.Fatalf("expected trimmed URL, got %q", cfg.ServiceURL)
        }
    })

    t.Run("create-cloud-sql-job", func(t *testing.T) {
        t.Parallel()
        parser := newMockFlagParser()
        parser.setString("operation", OperationCreateCloudSQLJob)
        parser.setString("job-name", "sql-job")
        parser.setString("project-id", "proj")
        parser.setString("pubsub-topic", "topic")
        parser.setString("db-instance-id", "instance")
        parser.setString("action", " Start ")

        cfg, err := ParseFlagsWithParser(parser)
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if cfg.Action != "Start" {
            t.Fatalf("expected action to retain casing, got %q", cfg.Action)
        }
    })

    t.Run("create-start-cloud-sql-job", func(t *testing.T) {
        t.Parallel()
        parser := newMockFlagParser()
        parser.setString("operation", OperationCreateStartCloudSQLJob)
        parser.setString("project-id", "proj")
        parser.setString("pubsub-topic", "topic")
        parser.setString("db-instance-id", "instance")
        parser.setString("limit", "10")

        cfg, err := ParseFlagsWithParser(parser)
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if cfg.Limit != "10" {
            t.Fatalf("expected limit to remain string, got %q", cfg.Limit)
        }
    })

    t.Run("list-jobs", func(t *testing.T) {
        t.Parallel()
        parser := newMockFlagParser()
        parser.setString("operation", OperationListJobs)
        parser.setString("location", "asia-northeast1")
        parser.setString("limit", "25")

        if _, err := ParseFlagsWithParser(parser); err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
    })
}

func TestParseFlagsWithParser_Errors(t *testing.T) {
    t.Parallel()

    t.Run("missing operation", func(t *testing.T) {
        parser := newMockFlagParser()
        if _, err := ParseFlagsWithParser(parser); err == nil {
            t.Fatal("expected error when operation missing")
        }
    })

    t.Run("unsupported operation", func(t *testing.T) {
        parser := newMockFlagParser()
        parser.setString("operation", "unknown")
        if _, err := ParseFlagsWithParser(parser); err == nil {
            t.Fatal("expected error for unsupported operation")
        }
    })

    t.Run("missing required parameter", func(t *testing.T) {
        parser := newMockFlagParser()
        parser.setString("operation", OperationCreatePubSubJob)
        parser.setString("job-name", "job")
        if _, err := ParseFlagsWithParser(parser); err == nil {
            t.Fatal("expected error when required parameter missing")
        }
    })

    t.Run("invalid http method", func(t *testing.T) {
        parser := newMockFlagParser()
        parser.setString("operation", OperationCreateHTTPJob)
        parser.setString("job-name", "job")
        parser.setString("project-id", "proj")
        parser.setString("service-url", "https://example.com")
        parser.setString("http-method", "trace")
        if _, err := ParseFlagsWithParser(parser); err == nil {
            t.Fatal("expected error for invalid http method")
        }
    })

    t.Run("invalid limit", func(t *testing.T) {
        parser := newMockFlagParser()
        parser.setString("operation", OperationListJobs)
        parser.setString("limit", "not-number")
        if _, err := ParseFlagsWithParser(parser); err == nil {
            t.Fatal("expected error for invalid limit")
        }
    })
}

func TestParseFlagsWithParser_HelpSkipsValidation(t *testing.T) {
    parser := newMockFlagParser()
    parser.setBool("help", true)

    cfg, err := ParseFlagsWithParser(parser)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !cfg.Help {
        t.Fatal("help flag should be true")
    }
}

func TestParseFlagsWithParser_ParseError(t *testing.T) {
    parser := newMockFlagParser()
    parser.parseErr = errors.New("parse failure")

    if _, err := ParseFlagsWithParser(parser); err == nil {
        t.Fatal("expected parse error")
    }
}

func TestParseFlags_StandardParser(t *testing.T) {
    origArgs := os.Args
    defer func() { os.Args = origArgs }()

    os.Args = []string{
        "gcloud-genset-scheduler",
        "-operation", OperationCreatePubSubJob,
        "-job-name", "job",
        "-project-id", "proj",
        "-pubsub-topic", "topic",
    }

    cfg, err := ParseFlags()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if cfg.Operation != OperationCreatePubSubJob {
        t.Fatalf("operation mismatch: %s", cfg.Operation)
    }
    if cfg.JobName != "job" {
        t.Fatalf("job name mismatch: %s", cfg.JobName)
    }
}

func TestStandardFlagParser(t *testing.T) {
    origArgs := os.Args
    defer func() { os.Args = origArgs }()

    os.Args = []string{"cmd", "-string-flag=val", "-bool-flag", "positional"}

    parser := NewStandardFlagParser()
    var stringFlag string
    var boolFlag bool

    parser.StringVar(&stringFlag, "string-flag", "", "string flag")
    parser.BoolVar(&boolFlag, "bool-flag", false, "bool flag")

    if err := parser.Parse(); err != nil {
        t.Fatalf("unexpected parse error: %v", err)
    }

    if stringFlag != "val" {
        t.Fatalf("expected val, got %s", stringFlag)
    }
    if !boolFlag {
        t.Fatalf("expected boolFlag to be true")
    }

    if args := parser.Args(); len(args) != 1 || args[0] != "positional" {
        t.Fatalf("unexpected args: %v", args)
    }
}

func TestPrintUsage(t *testing.T) {
    origStderr := os.Stderr
    defer func() { os.Stderr = origStderr }()

    r, w, err := os.Pipe()
    if err != nil {
        t.Fatalf("failed to create pipe: %v", err)
    }
    os.Stderr = w

    PrintUsage()

    w.Close()
    var buf bytes.Buffer
    if _, err := io.Copy(&buf, r); err != nil {
        t.Fatalf("failed to read usage output: %v", err)
    }

    output := buf.String()
    if !strings.Contains(output, "Cloud Scheduler") {
        t.Fatalf("unexpected usage output: %s", output)
    }
}
