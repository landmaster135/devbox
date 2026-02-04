package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// Clock abstracts time retrieval for testing.
type Clock func() time.Time

// StructuredLogger emits timestamps plus tag blocks before the message body.
type StructuredLogger struct {
	base  *log.Logger
	tags  []string
	clock Clock
}

// Option customizes StructuredLogger creation time behaviour.
type Option func(*StructuredLogger)

// New creates a StructuredLogger writing to stdout with no tags by default.
func New(opts ...Option) *StructuredLogger {
	l := &StructuredLogger{
		base:  log.New(os.Stdout, "", 0),
		clock: time.Now,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(l)
		}
	}
	return l
}

// Ensure returns the provided logger or instantiates a new one when nil.
func Ensure(logger *StructuredLogger) *StructuredLogger {
	if logger != nil {
		return logger
	}
	return New()
}

// WithWriter overrides the output writer.
func WithWriter(w io.Writer) Option {
	return func(l *StructuredLogger) {
		if w == nil {
			return
		}
		l.base = log.New(w, "", 0)
	}
}

// WithClock overrides the timestamp source.
func WithClock(clock Clock) Option {
	return func(l *StructuredLogger) {
		if clock == nil {
			return
		}
		l.clock = clock
	}
}

// WithLogger allows reusing an existing log.Logger.
func WithLogger(logger *log.Logger) Option {
	return func(l *StructuredLogger) {
		if logger == nil {
			return
		}
		l.base = logger
	}
}

// WithTags returns a cloned logger with additional tags appended.
func (l *StructuredLogger) WithTags(tags ...string) *StructuredLogger {
	base := Ensure(l)
	cloned := &StructuredLogger{
		base:  base.base,
		clock: base.clock,
	}
	if len(base.tags) > 0 {
		cloned.tags = append(cloned.tags, base.tags...)
	}
	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed == "" {
			continue
		}
		cloned.tags = append(cloned.tags, trimmed)
	}
	return cloned
}

// Infof prints the formatted message with the logger tags.
func (l *StructuredLogger) Infof(format string, args ...any) {
	l.output(format, args...)
}

// Warnf prints the formatted warning message.
func (l *StructuredLogger) Warnf(format string, args ...any) {
	l.output(format, args...)
}

// Errorf prints the formatted error message.
func (l *StructuredLogger) Errorf(format string, args ...any) {
	l.output(format, args...)
}

// output renders the final log line.
func (l *StructuredLogger) output(format string, args ...any) {
	base := Ensure(l)
	ts := base.clock().Format("2006-01-02 15:04:05")
	var builder strings.Builder
	builder.WriteString(ts)
	if len(base.tags) > 0 {
		builder.WriteByte(' ')
		builder.WriteString(formatTags(base.tags))
	}
	message := strings.TrimSpace(fmt.Sprintf(format, args...))
	if message != "" {
		builder.WriteByte(' ')
		builder.WriteString(message)
	}
	base.base.Println(builder.String())
}

func formatTags(tags []string) string {
	var builder strings.Builder
	for i, tag := range tags {
		if i > 0 {
			builder.WriteByte(' ')
		}
		builder.WriteByte('[')
		builder.WriteString(tag)
		builder.WriteByte(']')
	}
	return builder.String()
}
