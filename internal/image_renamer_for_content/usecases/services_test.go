package usecases

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

func TestProcessContentImageRename_DetectsConflicts(t *testing.T) {
	dir := t.TempDir()

	files := map[string]string{
		"MA0001_01.png": "already renamed",
		"000_raw.png":   "needs rename",
	}

	for name, body := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatalf("failed to create file %s: %v", name, err)
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
	if err == nil {
		t.Fatalf("expected conflict error, got nil")
	}
	if success != 0 {
		t.Fatalf("expected 0 successes due to early abort, got %d", success)
	}
	if failed != 1 {
		t.Fatalf("expected error count to be 1, got %d", failed)
	}
	if !strings.Contains(stderr.String(), "衝突") {
		t.Fatalf("expected stderr to mention conflict, got: %s", stderr.String())
	}
	if !strings.Contains(stdout.String(), "スキップ: 1") {
		t.Fatalf("expected stdout to include skip count, got: %s", stdout.String())
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read directory: %v", err)
	}

	existing := map[string]struct{}{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		existing[entry.Name()] = struct{}{}
	}

	if len(existing) != len(files) {
		t.Fatalf("expected files to remain untouched, got %+v", existing)
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

func TestProcessContentImageRename_NoImages(t *testing.T) {
	dir := t.TempDir()

	config := cfg.Config{
		SrcDir:     dir,
		SortByName: true,
		Operation:  "wine",
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	success, failed, err := ProcessContentImageRename(config, stdout, stderr)
	if err != nil {
		t.Fatalf("expected no error when directory is empty, got %v", err)
	}
	if success != 0 || failed != 0 {
		t.Fatalf("expected zero success and failure, got success=%d, failed=%d", success, failed)
	}

	if !strings.Contains(stdout.String(), "画像ファイルが見つかりませんでした") {
		t.Fatalf("expected message about missing images, got %q", stdout.String())
	}
}

func TestProcessContentImageRename_InvalidDirectory(t *testing.T) {
	missingDir := filepath.Join(t.TempDir(), "missing")

	config := cfg.Config{
		SrcDir:     missingDir,
		SortByTime: true,
		Operation:  "date",
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	success, failed, err := ProcessContentImageRename(config, stdout, stderr)
	if err == nil {
		t.Fatalf("expected error for missing directory")
	}
	if success != 0 || failed != 0 {
		t.Fatalf("expected zero counts when initialization fails, got success=%d failed=%d", success, failed)
	}
	if stderr.Len() == 0 {
		t.Fatalf("expected error output for missing directory")
	}
}

func TestFindImageFilesRecursive(t *testing.T) {
	root := t.TempDir()
	topFile := filepath.Join(root, "top.jpg")
	if err := os.WriteFile(topFile, []byte("foo"), 0o644); err != nil {
		t.Fatalf("failed to create top file: %v", err)
	}

	subDir := filepath.Join(root, "nested")
	if err := os.Mkdir(subDir, 0o755); err != nil {
		t.Fatalf("failed to create nested dir: %v", err)
	}
	subFile := filepath.Join(subDir, "inner.webp")
	if err := os.WriteFile(subFile, []byte("bar"), 0o644); err != nil {
		t.Fatalf("failed to create nested file: %v", err)
	}

	filesNonRecursive, err := findImageFiles(root, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filesNonRecursive) != 1 || filesNonRecursive[0] != topFile {
		t.Fatalf("expected only top-level file, got %v", filesNonRecursive)
	}

	filesRecursive, err := findImageFiles(root, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filesRecursive) != 2 {
		t.Fatalf("expected two files with recursion, got %v", filesRecursive)
	}

	paths := map[string]struct{}{}
	for _, path := range filesRecursive {
		paths[path] = struct{}{}
	}
	if _, ok := paths[topFile]; !ok {
		t.Fatalf("expected top file in recursive result")
	}
	if _, ok := paths[subFile]; !ok {
		t.Fatalf("expected nested file in recursive result")
	}
}

func TestBuildFileInfosPartialFailure(t *testing.T) {
	dir := t.TempDir()
	good := filepath.Join(dir, "good.png")
	if err := os.WriteFile(good, []byte("hello"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	missing := filepath.Join(dir, "missing.png")

	stderr := &bytes.Buffer{}
	infos, err := buildFileInfos([]string{good, missing}, stderr)
	if err == nil {
		t.Fatalf("expected error due to missing file")
	}
	if len(infos) != 1 || infos[0].path != good {
		t.Fatalf("expected only existing file info, got %+v", infos)
	}
	if !strings.Contains(stderr.String(), "情報取得に失敗") {
		t.Fatalf("expected error log for missing file, got %q", stderr.String())
	}
}

func TestRenameFilesSkipAndWorkers(t *testing.T) {
	dir := t.TempDir()
	initialPath := filepath.Join(dir, "MA0001_01.jpg")
	if err := os.WriteFile(initialPath, []byte("data"), 0o644); err != nil {
		t.Fatalf("failed to create initial file: %v", err)
	}

	info := fileInfo{path: initialPath, name: filepath.Base(initialPath)}
	config := cfg.Config{
		ContentID: "MA",
		Digits:    4,
		Suffix:    "01",
		Start:     1,
		Workers:   5,
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	success, errors, skipped := renameFiles([]fileInfo{info}, config, stdout, stderr)
	if success != 1 || errors != 0 || skipped != 0 {
		t.Fatalf("expected skip counted as success, got success=%d errors=%d skipped=%d", success, errors, skipped)
	}
	if _, err := os.Stat(initialPath); err != nil {
		t.Fatalf("expected original file to remain, err=%v", err)
	}
}

func TestRenameFilesRenameError(t *testing.T) {
	dir := t.TempDir()
	original := filepath.Join(dir, "orig.jpg")
	if err := os.WriteFile(original, []byte("content"), 0o644); err != nil {
		t.Fatalf("failed to create original file: %v", err)
	}

	targetDir := filepath.Join(dir, "MA0001_01.jpg")
	if err := os.Mkdir(targetDir, 0o755); err != nil {
		t.Fatalf("failed to create blocking directory: %v", err)
	}

	info := fileInfo{path: original, name: filepath.Base(original)}
	config := cfg.Config{
		ContentID: "MA",
		Digits:    4,
		Suffix:    "01",
		Start:     1,
		Workers:   2,
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	success, errors, skipped := renameFiles([]fileInfo{info}, config, stdout, stderr)
	if success != 0 || errors != 1 || skipped != 0 {
		t.Fatalf("expected rename failure to count as error, got success=%d errors=%d skipped=%d", success, errors, skipped)
	}
	if !strings.Contains(stderr.String(), "リネームに失敗") {
		t.Fatalf("expected rename failure logged, got %q", stderr.String())
	}
	if _, err := os.Stat(original); err != nil {
		t.Fatalf("expected original file to remain after failure: %v", err)
	}
}

func TestValidateConfigDigitsAndStartErrors(t *testing.T) {
	dir := t.TempDir()

	buffer := &bytes.Buffer{}
	configDigits := &cfg.Config{
		ContentID:  "MA",
		SortByName: true,
		Digits:     0,
		Start:      1,
		SrcDir:     dir,
	}
	if err := validateConfig(configDigits, buffer); err == nil {
		t.Fatalf("expected error when digits <= 0")
	}

	configStart := &cfg.Config{
		ContentID:  "MA",
		SortByName: true,
		Digits:     4,
		Start:      0,
		SrcDir:     dir,
	}
	buffer.Reset()
	if err := validateConfig(configStart, buffer); err == nil {
		t.Fatalf("expected error when start <= 0")
	}
}

func TestValidateConfigBothSortFlags(t *testing.T) {
	dir := t.TempDir()

	config := &cfg.Config{
		ContentID:  "MA",
		SortByName: true,
		SortByTime: true,
		Digits:     4,
		Start:      1,
		SrcDir:     dir,
	}
	if err := validateConfig(config, &bytes.Buffer{}); err != nil {
		t.Fatalf("did not expect error when both sort flags set, got %v", err)
	}
	if config.SortByTime {
		t.Fatalf("expected SortByTime to be reset to false when both flags provided")
	}
}

func TestProcessContentImageRename_RenameFailure(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.jpg")
	if err := os.WriteFile(input, []byte("sample"), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	blockingDir := filepath.Join(dir, "MA0001_01.jpg")
	if err := os.Mkdir(blockingDir, 0o755); err != nil {
		t.Fatalf("failed to create blocking directory: %v", err)
	}

	config := cfg.Config{
		SrcDir:     dir,
		SortByName: true,
		Operation:  "mackerel",
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	success, failed, err := ProcessContentImageRename(config, stdout, stderr)
	if err == nil {
		t.Fatalf("expected error when rename fails")
	}
	if success != 0 || failed != 1 {
		t.Fatalf("expected zero success and one failure, got success=%d failed=%d", success, failed)
	}
	if !strings.Contains(stderr.String(), "リネームに失敗") {
		t.Fatalf("expected stderr to mention rename failure, got %q", stderr.String())
	}
}

func TestProcessContentImageRename_WarnsMissingFileInfo(t *testing.T) {
	dir := t.TempDir()
	file1 := filepath.Join(dir, "a.jpg")
	if err := os.WriteFile(file1, []byte("a"), 0o644); err != nil {
		t.Fatalf("failed to create file1: %v", err)
	}

	missingTarget := filepath.Join(dir, "does-not-exist.jpg")
	brokenLink := filepath.Join(dir, "broken.jpg")
	if err := os.Symlink(missingTarget, brokenLink); err != nil {
		t.Skipf("symlink not supported on this platform: %v", err)
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
		t.Fatalf("did not expect fatal error, got %v", err)
	}
	if success != 1 || failed != 0 {
		t.Fatalf("expected one success and zero failures, got success=%d failed=%d", success, failed)
	}
	if !strings.Contains(stderr.String(), "警告: 一部のファイル情報の取得に失敗しました") {
		t.Fatalf("expected warning about missing file info, got %q", stderr.String())
	}
}

func TestFindImageFilesError(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing")
	if _, err := findImageFiles(missing, true); err == nil {
		t.Fatalf("expected error when directory missing")
	}
}

func TestBuildNewFileNameVariants(t *testing.T) {
	configDelimiter := cfg.Config{ContentID: "AB", Digits: 3, Delimiter: "-"}
	if got := buildNewFileName(configDelimiter, 7, ".png"); got != "AB-007.png" {
		t.Fatalf("unexpected result with delimiter, got %s", got)
	}

	configSuffix := cfg.Config{ContentID: "CD", Digits: 2, Suffix: "final"}
	if got := buildNewFileName(configSuffix, 3, ".jpg"); got != "CD03_final.jpg" {
		t.Fatalf("unexpected result with suffix, got %s", got)
	}

	configPlain := cfg.Config{ContentID: "EF", Digits: 2}
	if got := buildNewFileName(configPlain, 12, ".webp"); got != "EF12.webp" {
		t.Fatalf("unexpected result without extras, got %s", got)
	}
}

func TestProcessContentImageRename_SourceNotDirectory(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "single.jpg")
	if err := os.WriteFile(filePath, []byte("content"), 0o644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	config := cfg.Config{
		SrcDir:     filePath,
		SortByName: true,
		Operation:  "wine",
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	success, failed, err := ProcessContentImageRename(config, stdout, stderr)
	if err == nil {
		t.Fatalf("expected error when src is not directory")
	}
	if success != 0 || failed != 0 {
		t.Fatalf("expected zero counts on directory error, got success=%d failed=%d", success, failed)
	}
	if stderr.Len() == 0 {
		t.Fatalf("expected stderr to contain error message")
	}
}
