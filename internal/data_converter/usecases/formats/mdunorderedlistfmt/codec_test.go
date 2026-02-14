package mdunorderedlistfmt

import (
	"reflect"
	"testing"
)

func TestParse_Normal(t *testing.T) {
	t.Parallel()

	got, err := Parse([]byte("- Alice\n* Bob\n+ Carol\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantKeys := []string{"item"}
	wantRecords := []map[string]string{
		{"item": "Alice"},
		{"item": "Bob"},
		{"item": "Carol"},
	}
	if !reflect.DeepEqual(got.Keys, wantKeys) {
		t.Fatalf("keys mismatch: got=%v want=%v", got.Keys, wantKeys)
	}
	if !reflect.DeepEqual(got.KeyValueList, wantRecords) {
		t.Fatalf("records mismatch: got=%v want=%v", got.KeyValueList, wantRecords)
	}
}

func TestParse_ErrorCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty",
			input: "",
		},
		{
			name:  "invalid marker",
			input: "1. Alice\n",
		},
		{
			name:  "empty item",
			input: "- \n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if _, err := Parse([]byte(tt.input)); err == nil {
				t.Fatal("expected error but got nil")
			}
		})
	}
}

func TestSerialize_Normal(t *testing.T) {
	t.Parallel()

	out, err := Serialize(
		[]map[string]string{
			{"item": "Alice"},
			{"item": "Bob"},
		},
		[]string{"item"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "- Alice\n- Bob\n"
	if string(out) != want {
		t.Fatalf("serialized output mismatch: got=%q want=%q", string(out), want)
	}
}

func TestSerialize_UsesItemKeyWhenMultipleKeys(t *testing.T) {
	t.Parallel()

	out, err := Serialize(
		[]map[string]string{
			{"item": "Alice", "name": "ignored"},
		},
		[]string{"name", "item"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "- Alice\n"
	if string(out) != want {
		t.Fatalf("serialized output mismatch: got=%q want=%q", string(out), want)
	}
}

func TestSerialize_ErrorCases(t *testing.T) {
	t.Parallel()

	if _, err := Serialize([]map[string]string{{"item": "Alice"}}, nil); err == nil {
		t.Fatal("expected error but got nil")
	}

	if _, err := Serialize(
		[]map[string]string{{"name": "Alice", "age": "30"}},
		[]string{"name", "age"},
	); err == nil {
		t.Fatal("expected error but got nil")
	}
}
