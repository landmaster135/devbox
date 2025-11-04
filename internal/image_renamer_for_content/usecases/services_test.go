package usecases

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	cfg "github.com/landmaster135/devbox/internal/image_renamer_for_content/config"
)

func TestProcessContentImageRename_SortByName(t *testing.T) {
	dir := t.TempDir()

	files := []string{"b.webp", "a.jpg", "c.png"}
	for _, name := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte("fake image"), 0o644); err != nil {
			t.Fatalf("failed to create test file %s: %v", name, err)
		}
	}

	config := cfg.Config{
		SrcDir:     dir,
		SortByName: true,
		Operation:  "mackerel",
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	success, failed, err := ProcessContentImageRename(config, stdout, stderr)
	if err != nil {
		t.Fatalf("expected no error, got %v (stderr: %s)", err, stderr.String())
	}
	if failed != 0 {
		t.Fatalf("expected no failures, got %d", failed)
	}
	if success != len(files) {
		t.Fatalf("expected %d successes, got %d", len(files), success)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read renamed directory: %v", err)
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		names = append(names, entry.Name())
	}

	sort.Strings(names)
	expected := []string{"MA0001_01.jpg", "MA0002_01.webp", "MA0003_01.png"}

	if len(names) != len(expected) {
		t.Fatalf("expected %d files, got %d (%v)", len(expected), len(names), names)
	}

	for i, name := range expected {
		if names[i] != name {
			to := stdout.String()
			te := stderr.String()
			t.Fatalf("unexpected file order: index %d expected %s got %s\nstdout: %s\nstderr: %s", i, name, names[i], to, te)
		}
	}
}

func TestProcessContentImageRename_SortByTime(t *testing.T) {
	dir := t.TempDir()

	now := time.Now()
	testFiles := []struct {
		name string
		time time.Time
		ext  string
	}{
		{name: "later.png", time: now.Add(-1 * time.Hour), ext: ".png"},
		{name: "earliest.jpg", time: now.Add(-3 * time.Hour), ext: ".jpg"},
		{name: "middle.webp", time: now.Add(-2 * time.Hour), ext: ".webp"},
	}

	for _, file := range testFiles {
		path := filepath.Join(dir, file.name)
		if err := os.WriteFile(path, []byte("fake"), 0o644); err != nil {
			t.Fatalf("failed to create file %s: %v", file.name, err)
		}
		if err := os.Chtimes(path, file.time, file.time); err != nil {
			t.Fatalf("failed to change mod time for %s: %v", file.name, err)
		}
	}

	config := cfg.Config{
		SrcDir:     dir,
		SortByTime: true,
		Delimiter:  "-",
		Suffix:     "99",
		Start:      5,
		Operation:  "mackerel",
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	success, failed, err := ProcessContentImageRename(config, stdout, stderr)
	if err != nil {
		t.Fatalf("expected no error, got %v (stderr: %s)", err, stderr.String())
	}
	if failed != 0 {
		t.Fatalf("expected no failures, got %d", failed)
	}
	if success != len(testFiles) {
		t.Fatalf("expected %d successes, got %d", len(testFiles), success)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read renamed directory: %v", err)
	}

	expected := map[string]struct{}{
		"MA-0005_99.jpg":  {},
		"MA-0006_99.webp": {},
		"MA-0007_99.png":  {},
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if _, ok := expected[entry.Name()]; !ok {
			t.Fatalf("unexpected file renamed to %s", entry.Name())
		}
		delete(expected, entry.Name())
	}

	if len(expected) != 0 {
		t.Fatalf("missing expected files: %v", expected)
	}
}

func TestProcessContentImageRename_MissingSortFlag(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "sample.jpg"), []byte("fake"), 0o644); err != nil {
		t.Fatalf("failed to create sample file: %v", err)
	}

	config := cfg.Config{
		SrcDir:    dir,
		Operation: "mackerel",
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	success, failed, err := ProcessContentImageRename(config, stdout, stderr)
	if err == nil {
		t.Fatalf("expected error due to missing sort flag, got nil")
	}
	if success != 0 {
		t.Fatalf("expected 0 successes, got %d", success)
	}
	if failed != 0 {
		t.Fatalf("expected 0 failures prior to processing, got %d", failed)
	}
}

func TestProcessContentImageRename_InvalidOperation(t *testing.T) {
	config := cfg.Config{
		SrcDir:    "./not-real",
		Operation: "unknown",
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	if _, _, err := ProcessContentImageRename(config, stdout, stderr); err == nil {
		t.Fatalf("expected error for invalid operation")
	}
	if stderr.Len() == 0 {
		t.Fatalf("expected error message to be written to stderr")
	}
}

func TestApplyOperationPresetVariants(t *testing.T) {
	tests := []struct {
		name           string
		operation      string
		expectedID     string
		expectedDigits int
	}{
		{name: "Mackerel", operation: "mackerel", expectedID: "MA", expectedDigits: 4},
		{name: "WebClip", operation: "web_clip", expectedID: "WC", expectedDigits: 9},
		{name: "Date", operation: "date", expectedID: "DA", expectedDigits: 5},
		{name: "Wine", operation: "wine", expectedID: "WI", expectedDigits: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := cfg.Config{Operation: tt.operation}
			err := applyOperationPreset(&config, &bytes.Buffer{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if config.ContentID != tt.expectedID {
				t.Fatalf("expected content ID %s, got %s", tt.expectedID, config.ContentID)
			}
			if config.Digits != tt.expectedDigits {
				t.Fatalf("expected digits %d, got %d", tt.expectedDigits, config.Digits)
			}
		})
	}
}
