package logging

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func fixedClock() time.Time {
	return time.Date(2026, time.February, 4, 12, 34, 56, 0, time.UTC)
}

func TestStructuredLoggerFormatting(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(WithWriter(buf), WithClock(fixedClock))
	logger = logger.WithTags("CRON workflow").WithTags("daily-heading")

	logger.Infof("dispatched heading content to Discord")

	got := buf.String()
	want := "2026-02-04 12:34:56 [CRON workflow] [daily-heading] dispatched heading content to Discord\n"
	if got != want {
		t.Fatalf("unexpected output:\nwant: %q\n got: %q", want, got)
	}
}

func TestWithTagsDoesNotMutateParent(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(WithWriter(buf), WithClock(fixedClock))

	parent := logger.WithTags("parent")
	child := parent.WithTags("child")

	parent.Infof("parent only")
	child.Infof("child log")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 log lines, got %d", len(lines))
	}
	if strings.Contains(lines[0], "[child]") {
		t.Fatalf("parent log unexpectedly contains child tag: %s", lines[0])
	}
	if !strings.Contains(lines[1], "[child]") {
		t.Fatalf("child log missing child tag: %s", lines[1])
	}
}

func TestEnsureReturnsLoggerWhenNil(t *testing.T) {
	logger := Ensure(nil)
	if logger == nil {
		t.Fatalf("Ensure(nil) returned nil")
	}
}
