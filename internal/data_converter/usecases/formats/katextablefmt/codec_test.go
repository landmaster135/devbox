package katextablefmt

import (
	"reflect"
	"strings"
	"testing"
)

func TestParse_Normal(t *testing.T) {
	t.Parallel()

	input := `$$ \def\arraystretch{1.4}
\small
\begin{array}{|l|l|}
\hline
\text{name} & \text{age} \\
\hline
\text{Alice} & \text{30} \\
\hline
\text{Bob} & \text{28} \\
\hline
\end{array}
$$
`

	got, err := Parse([]byte(input))
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

	input := `\begin{array}{|l|l|}
\hline
\text{name} & \text{age} \\
\hline
\end{array}`

	got, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.KeyValueList) != 0 {
		t.Fatalf("records length mismatch: got=%d want=%d", len(got.KeyValueList), 0)
	}
}

func TestParse_ErrorCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "no array block",
			input: `\text{name}`,
		},
		{
			name: "invalid text macro",
			input: `\begin{array}{|l|}
\hline
\text{name \\
\hline
\end{array}`,
		},
		{
			name: "duplicate header",
			input: `\begin{array}{|l|l|}
\hline
\text{name} & \text{name} \\
\hline
\end{array}`,
		},
		{
			name: "too many columns",
			input: `\begin{array}{|l|}
\hline
\text{name} \\
\hline
\text{Alice} & \text{30} \\
\hline
\end{array}`,
		},
		{
			name: "missing row terminator",
			input: `\begin{array}{|l|}
\hline
\text{name}
\hline
\end{array}`,
		},
		{
			name: "column definition not closed",
			input: `\begin{array}{|l|
\hline
\text{name} \\
\hline
\end{array}`,
		},
		{
			name: "invalid text prefix",
			input: `\begin{array}{|l|}
\hline
\text name \\
\hline
\end{array}`,
		},
		{
			name: "text trailing token",
			input: `\begin{array}{|l|}
\hline
\text{name}x \\
\hline
\end{array}`,
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

func TestParse_PlainCellAndUnescape(t *testing.T) {
	t.Parallel()

	input := `\begin{array}{|l|l|}
\hline
name & memo \\
\hline
Alice & \text{A\&B \{x\} \textbackslash{}} \\
\hline
\end{array}`

	got, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []map[string]string{
		{"name": "Alice", "memo": `A&B {x} \`},
	}
	if !reflect.DeepEqual(got.KeyValueList, want) {
		t.Fatalf("records mismatch: got=%v want=%v", got.KeyValueList, want)
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

	want := `$$ \def\arraystretch{1.4}
\small
\begin{array}{|l|l|}
\hline
\text{name} & \text{age} \\
\hline
\text{Alice} & \text{30} \\
\hline
\text{Bob} & \text{28} \\
\hline
\end{array}
$$
`
	if string(out) != want {
		t.Fatalf("serialized output mismatch: got=%q want=%q", string(out), want)
	}
}

func TestSerialize_EscapesSpecialCharacters(t *testing.T) {
	t.Parallel()

	out, err := Serialize(
		[]map[string]string{
			{"memo": `A&B {x} \`},
		},
		[]string{"memo"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := string(out)
	if !strings.Contains(got, `\text{A\&B \{x\} \textbackslash{}}`) {
		t.Fatalf("escaped output not found: %s", got)
	}
}

func TestSerialize_EmptyKeys(t *testing.T) {
	t.Parallel()

	if _, err := Serialize([]map[string]string{{"name": "Alice"}}, nil); err == nil {
		t.Fatal("expected error but got nil")
	}
}
