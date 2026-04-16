package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Empty Docs Inspector
//
// Usage:
//   go run empty_docs_inspection.go --docs-dir <target-docs-dir>
//
// Example:
//   go run empty_docs_inspection.go --docs-dir docs
//   go run empty_docs_inspection.go --docs-dir /path/to/repo/docs

func main() {
	var docsDir string

	// Parse CLI arguments.
	flag.StringVar(&docsDir, "docs-dir", "docs", "path to docs directory")
	flag.Parse()

	// Run inspection and return non-zero on failure for CI/scripting usage.
	emptyFiles, err := inspectEmptyDocsFiles(docsDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if len(emptyFiles) == 0 {
		fmt.Printf("OK: no empty docs files found: %s\n", docsDir)
		return
	}

	report := "empty docs files detected\n"
	for _, path := range emptyFiles {
		report += fmt.Sprintf("- %s\n", path)
	}

	fmt.Fprint(os.Stderr, report)
	os.Exit(1)
}

// inspectEmptyDocsFiles walks docsDir recursively and returns files
// that are either zero-byte or whitespace-only.
func inspectEmptyDocsFiles(docsDir string) ([]string, error) {
	// Validate docs-dir itself.
	info, err := os.Stat(docsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to access docs directory %q: %w", docsDir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("docs path is not a directory: %s", docsDir)
	}

	emptyFiles := []string{}
	err = filepath.WalkDir(docsDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		// Only inspect regular files in docs.
		if !d.Type().IsRegular() {
			return nil
		}

		isEmpty, err := isEmptyFile(path)
		if err != nil {
			return err
		}
		if isEmpty {
			emptyFiles = append(emptyFiles, toRelativeSlash(docsDir, path))
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to inspect docs directory %q: %w", docsDir, err)
	}

	sort.Strings(emptyFiles)
	return emptyFiles, nil
}

// isEmptyFile returns true when a file is zero-byte
// or contains only whitespace characters.
func isEmptyFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if info.Size() == 0 {
		return true, nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(content)) == "", nil
}

// toRelativeSlash converts absolute path to docs-dir-relative slash path.
func toRelativeSlash(baseDir string, targetPath string) string {
	rel, err := filepath.Rel(baseDir, targetPath)
	if err != nil {
		return filepath.ToSlash(targetPath)
	}
	return filepath.ToSlash(rel)
}
