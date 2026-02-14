package usecases

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestNormalizeToKeyValueList_SupportedFormats(t *testing.T) {
	t.Parallel()

	svc := NewService()

	tests := []struct {
		name        string
		format      string
		input       string
		wantKeys    []string
		wantRecords []map[string]string
	}{
		{
			name:     "json",
			format:   "json",
			input:    `[{"name":"Alice","age":30}]`,
			wantKeys: []string{"age", "name"},
			wantRecords: []map[string]string{
				{"name": "Alice", "age": "30"},
			},
		},
		{
			name:     "yaml",
			format:   "yaml",
			input:    "- name: Alice\n  age: \"30\"\n",
			wantKeys: []string{"age", "name"},
			wantRecords: []map[string]string{
				{"name": "Alice", "age": "30"},
			},
		},
		{
			name:     "csv",
			format:   "csv",
			input:    "name,age\nAlice,30\n",
			wantKeys: []string{"name", "age"},
			wantRecords: []map[string]string{
				{"name": "Alice", "age": "30"},
			},
		},
		{
			name:     "tsv",
			format:   "tsv",
			input:    "name\tage\nAlice\t30\n",
			wantKeys: []string{"name", "age"},
			wantRecords: []map[string]string{
				{"name": "Alice", "age": "30"},
			},
		},
		{
			name:   "html",
			format: "html",
			input: `<table>
  <thead><tr><th>name</th><th>age</th></tr></thead>
  <tbody><tr><td>Alice</td><td>30</td></tr></tbody>
</table>`,
			wantKeys: []string{"name", "age"},
			wantRecords: []map[string]string{
				{"name": "Alice", "age": "30"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := svc.NormalizeToKeyValueList([]byte(tt.input), tt.format)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got.Keys, tt.wantKeys) {
				t.Fatalf("keys mismatch: got=%v want=%v", got.Keys, tt.wantKeys)
			}
			if !reflect.DeepEqual(got.KeyValueList, tt.wantRecords) {
				t.Fatalf("records mismatch: got=%v want=%v", got.KeyValueList, tt.wantRecords)
			}
		})
	}
}

func TestSerializeFromKeyValueList_AllFormats(t *testing.T) {
	t.Parallel()

	svc := NewService()
	source := &NormalizedData{
		Keys: []string{"name", "age"},
		KeyValueList: []map[string]string{
			{"name": "Alice", "age": "30"},
			{"name": "Bob", "age": "28"},
		},
	}

	formats := []string{"json", "yaml", "csv", "tsv", "html"}
	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			t.Parallel()

			out, err := svc.SerializeFromKeyValueList(source, format)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			reparsed, err := svc.NormalizeToKeyValueList(out, format)
			if err != nil {
				t.Fatalf("reparse failed: %v", err)
			}

			if !reflect.DeepEqual(reparsed.KeyValueList, source.KeyValueList) {
				t.Fatalf("reparsed records mismatch: got=%v want=%v", reparsed.KeyValueList, source.KeyValueList)
			}
		})
	}
}

func TestConvertFile_WritesOutputFile(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "input.json")
	outputPath := filepath.Join(tmpDir, "output.csv")
	input := `[{"name":"Alice","age":30},{"name":"Bob","age":28}]`
	if err := os.WriteFile(inputPath, []byte(input), 0o644); err != nil {
		t.Fatalf("failed to write input file: %v", err)
	}

	svc := NewService()
	message, err := svc.ConvertFile(inputPath, outputPath, "json", "csv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(message, "変換完了") {
		t.Fatalf("unexpected message: %s", message)
	}

	out, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	normalized, err := svc.NormalizeToKeyValueList(out, "csv")
	if err != nil {
		t.Fatalf("failed to parse output csv: %v", err)
	}

	want := []map[string]string{
		{"name": "Alice", "age": "30"},
		{"name": "Bob", "age": "28"},
	}
	if !reflect.DeepEqual(normalized.KeyValueList, want) {
		t.Fatalf("output mismatch: got=%v want=%v", normalized.KeyValueList, want)
	}
}

func TestNormalizeToKeyValueList_ErrorCases(t *testing.T) {
	t.Parallel()

	svc := NewService()
	if _, err := svc.NormalizeToKeyValueList([]byte("x"), "txt"); err == nil {
		t.Fatal("expected unsupported format error")
	}

	if _, err := svc.NormalizeToKeyValueList([]byte("<div>no table</div>"), "html"); err == nil {
		t.Fatal("expected html table error")
	}
}
