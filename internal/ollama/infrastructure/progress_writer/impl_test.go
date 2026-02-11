package infrastructure

import (
	"bytes"
	"testing"
)

func TestPullProgressWriter_RewritesSameDigest(t *testing.T) {
	buf := &bytes.Buffer{}
	writer := NewPullProgressWriter(buf)
	writer.isTTY = true

	first := "pulling abc123 10.0% (10/100)"
	second := "pulling abc123 20.0% (20/100)"

	if _, err := writer.Write([]byte(first + "\n")); err != nil {
		t.Fatalf("write first: %v", err)
	}
	if _, err := writer.Write([]byte(second + "\n")); err != nil {
		t.Fatalf("write second: %v", err)
	}

	expected := first + "\n" + "\033[1A\033[2K" + second + "\n"
	if buf.String() != expected {
		t.Fatalf("unexpected output. got %q", buf.String())
	}
}

func TestPullProgressWriter_NonTTY(t *testing.T) {
	buf := &bytes.Buffer{}
	writer := NewPullProgressWriter(buf)
	writer.isTTY = false

	first := "pulling abc123 10.0% (10/100)"
	second := "pulling abc123 20.0% (20/100)"

	if _, err := writer.Write([]byte(first + "\n")); err != nil {
		t.Fatalf("write first: %v", err)
	}
	if _, err := writer.Write([]byte(second + "\n")); err != nil {
		t.Fatalf("write second: %v", err)
	}

	expected := first + "\n" + second + "\n"
	if buf.String() != expected {
		t.Fatalf("unexpected output. got %q", buf.String())
	}
}
