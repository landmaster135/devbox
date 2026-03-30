package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Docs Directory Validator
//
// Usage:
//   go run directory_validate.go --docs-dir <target-docs-dir>
//
// Example:
//   go run directory_validate.go --docs-dir docs
//   go run directory_validate.go --docs-dir /path/to/repo/docs

type expectedEntry struct {
	// base_dir is a base relative path from docs-dir.
	// Current baseline uses "." to mean direct child directories of docs-dir.
	base_dir string
	// dir_name is the required directory name under path.
	dir_name string
	// file_names are required files under dir_name.
	file_names []string
}

// expectedLayout defines the canonical minimum docs structure.
var expectedLayout = []expectedEntry{
	{base_dir: ".", dir_name: "changelog", file_names: []string{"README.md"}},
	{base_dir: ".", dir_name: "docs_management", file_names: []string{"index.md"}},
	{base_dir: ".", dir_name: "exec_plans", file_names: []string{"index.md"}},
	{base_dir: ".", dir_name: "project_overview", file_names: []string{"index.md"}},
	{base_dir: ".", dir_name: "project_status", file_names: []string{"index.md"}},
	{base_dir: ".", dir_name: "task_implementation", file_names: []string{"index.md"}},
	{base_dir: ".", dir_name: "tool_implementation", file_names: []string{"index.md"}},
	{base_dir: ".", dir_name: "user_prompt", file_names: []string{"index.md"}},
}

func main() {
	var docsDir string

	// Parse CLI arguments.
	flag.StringVar(&docsDir, "docs-dir", "docs", "path to docs directory")
	flag.Parse()

	// Run validation and return non-zero on failure for CI/scripting usage.
	if err := validateLayout(docsDir); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf("OK: docs layout is valid: %s\n", docsDir)
}

// validateLayout verifies that required directories/files exist
// and have expected types.
func validateLayout(docsDir string) error {
	// Validate docs-dir itself.
	info, err := os.Stat(docsDir)
	if err != nil {
		return fmt.Errorf("failed to access docs directory %q: %w", docsDir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("docs path is not a directory: %s", docsDir)
	}

	missing := collectMissing(docsDir)
	typeMismatch := collectTypeMismatch(docsDir)

	if len(missing) == 0 && len(typeMismatch) == 0 {
		return nil
	}

	report := "docs layout validation failed\n"
	if len(missing) > 0 {
		report += formatList("missing entries", missing)
	}
	if len(typeMismatch) > 0 {
		report += formatList("type mismatch entries", typeMismatch)
	}

	return fmt.Errorf("%s", report)
}

// collectMissing returns required directories/files that do not exist.
func collectMissing(docsDir string) []string {
	missing := []string{}
	for _, entry := range expectedLayout {
		dirPath := filepath.Join(docsDir, filepath.FromSlash(entry.base_dir), entry.dir_name)
		dirRel := toRelativeSlash(docsDir, dirPath)
		dirInfo, err := os.Stat(dirPath)
		if os.IsNotExist(err) {
			missing = append(missing, dirRel)
			continue
		}
		if err != nil || !dirInfo.IsDir() {
			continue
		}

		for _, fileName := range entry.file_names {
			filePath := filepath.Join(dirPath, fileName)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				missing = append(missing, toRelativeSlash(docsDir, filePath))
			}
		}
	}
	sort.Strings(missing)
	return missing
}

// collectTypeMismatch returns paths where directory/file type is incorrect.
func collectTypeMismatch(docsDir string) []string {
	mismatch := []string{}
	for _, entry := range expectedLayout {
		dirPath := filepath.Join(docsDir, filepath.FromSlash(entry.base_dir), entry.dir_name)
		dirInfo, err := os.Stat(dirPath)
		if err != nil {
			continue
		}

		dirRel := toRelativeSlash(docsDir, dirPath)
		if !dirInfo.IsDir() {
			mismatch = append(mismatch, fmt.Sprintf("%s expected=dir actual=file", dirRel))
			continue
		}

		for _, fileName := range entry.file_names {
			filePath := filepath.Join(dirPath, fileName)
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				continue
			}
			if fileInfo.IsDir() {
				mismatch = append(mismatch, fmt.Sprintf("%s expected=file actual=dir", toRelativeSlash(docsDir, filePath)))
			}
		}
	}
	sort.Strings(mismatch)
	return mismatch
}

// toRelativeSlash converts absolute path to docs-dir-relative slash path.
func toRelativeSlash(baseDir string, targetPath string) string {
	rel, err := filepath.Rel(baseDir, targetPath)
	if err != nil {
		return filepath.ToSlash(targetPath)
	}
	return filepath.ToSlash(rel)
}

// formatList formats validation results for readable multi-line output.
func formatList(title string, items []string) string {
	out := fmt.Sprintf("- %s:\n", title)
	for _, item := range items {
		out += fmt.Sprintf("  - %s\n", item)
	}
	return out
}
