package usecases

import (
	"reflect"
	"strings"
	"testing"
)

func TestServiceKaTeXTable_Normal(t *testing.T) {
	t.Parallel()

	svc := NewService()
	source := &NormalizedData{
		Keys: []string{"service", "issue"},
		KeyValueList: []map[string]string{
			{"service": "Windows", "issue": "Grouping"},
			{"service": "VSCode", "issue": "Popover"},
		},
	}

	out, err := svc.SerializeFromKeyValueList(source, "katex-table")
	if err != nil {
		t.Fatalf("unexpected serialize error: %v", err)
	}
	if !strings.Contains(string(out), `\begin{array}{|l|l|}`) {
		t.Fatalf("unexpected output: %s", string(out))
	}

	reparsed, err := svc.NormalizeToKeyValueList(out, "katex-table")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if !reflect.DeepEqual(reparsed.Keys, source.Keys) {
		t.Fatalf("keys mismatch: got=%v want=%v", reparsed.Keys, source.Keys)
	}
	if !reflect.DeepEqual(reparsed.KeyValueList, source.KeyValueList) {
		t.Fatalf("records mismatch: got=%v want=%v", reparsed.KeyValueList, source.KeyValueList)
	}
}

func TestServiceKaTeXTable_InvalidInput(t *testing.T) {
	t.Parallel()

	svc := NewService()
	if _, err := svc.NormalizeToKeyValueList([]byte(`\text{a}`), "katex-table"); err == nil {
		t.Fatal("expected error but got nil")
	}
}
