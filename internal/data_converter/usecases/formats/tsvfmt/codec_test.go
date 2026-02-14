package tsvfmt

import (
	"reflect"
	"testing"
)

func TestParse_Normal(t *testing.T) {
	t.Parallel()

	got, err := Parse([]byte("name\tage\nAlice\t30\nBob\t28\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantKeys := []string{"name", "age"}
	wantRecords := []map[string]string{
		{"name": "Alice", "age": "30"},
		{"name": "Bob", "age": "28"},
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
			name:  "empty header",
			input: "name\t\nAlice\t30\n",
		},
		{
			name:  "duplicate header",
			input: "name\tname\nAlice\t30\n",
		},
		{
			name:  "too many columns",
			input: "name\tage\nAlice\t30\textra\n",
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
			{"name": "Alice", "age": "30"},
			{"name": "Bob", "age": "28"},
		},
		[]string{"name", "age"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "name\tage\nAlice\t30\nBob\t28\n"
	if string(out) != want {
		t.Fatalf("serialized output mismatch: got=%q want=%q", string(out), want)
	}
}

func TestSerialize_EmptyKeys(t *testing.T) {
	t.Parallel()

	if _, err := Serialize([]map[string]string{{"name": "Alice"}}, nil); err == nil {
		t.Fatal("expected error but got nil")
	}
}
