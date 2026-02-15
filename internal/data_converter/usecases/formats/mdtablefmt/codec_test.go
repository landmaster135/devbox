package mdtablefmt

import (
	"reflect"
	"testing"
)

func TestParse_Normal(t *testing.T) {
	t.Parallel()

	got, err := Parse([]byte("| name | age |\n| --- | :---: |\n| Alice | 30 |\n| Bob | 28 |\n"))
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

func TestParse_NoDataRows(t *testing.T) {
	t.Parallel()

	got, err := Parse([]byte("| name | age |\n| --- | --- |\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.KeyValueList) != 0 {
		t.Fatalf("records length mismatch: got=%d want=%d", len(got.KeyValueList), 0)
	}
}

func TestParse_EscapedPipe(t *testing.T) {
	t.Parallel()

	got, err := Parse([]byte("| item |\n| --- |\n| A\\|B |\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []map[string]string{{"item": "A|B"}}
	if !reflect.DeepEqual(got.KeyValueList, want) {
		t.Fatalf("records mismatch: got=%v want=%v", got.KeyValueList, want)
	}
}

func TestParse_ErrorCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "insufficient rows",
			input: "| name | age |\n",
		},
		{
			name:  "duplicate header",
			input: "| name | name |\n| --- | --- |\n",
		},
		{
			name:  "empty header",
			input: "| name |   |\n| --- | --- |\n",
		},
		{
			name:  "invalid separator",
			input: "| name | age |\n| --- | --x |\n",
		},
		{
			name:  "separator without pipe",
			input: "| name | age |\n--- ---\n",
		},
		{
			name:  "too many columns",
			input: "| name | age |\n| --- | --- |\n| Alice | 30 | extra |\n",
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

	want := "| name | age |\n| --- | --- |\n| Alice | 30 |\n| Bob | 28 |\n"
	if string(out) != want {
		t.Fatalf("serialized output mismatch: got=%q want=%q", string(out), want)
	}
}

func TestSerialize_EscapesPipe(t *testing.T) {
	t.Parallel()

	out, err := Serialize(
		[]map[string]string{
			{"item": "A|B"},
		},
		[]string{"item"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "| item |\n| --- |\n| A\\|B |\n"
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
