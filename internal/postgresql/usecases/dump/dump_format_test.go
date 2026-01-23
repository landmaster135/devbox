package dump

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testStringer struct {
	value string
}

func (s testStringer) String() string {
	return "stringer:" + s.value
}

func TestFormatCSVValue_VariousTypes(t *testing.T) {
	sampleTime := time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name  string
		input any
		want  string
	}{
		{name: "nil", input: nil, want: ""},
		{name: "string", input: "hello", want: "hello"},
		{name: "bytesEmpty", input: []byte{}, want: ""},
		{name: "bytes", input: []byte("hello"), want: base64.StdEncoding.EncodeToString([]byte("hello"))},
		{name: "time", input: sampleTime, want: sampleTime.Format(time.RFC3339Nano)},
		{name: "bool", input: true, want: "true"},
		{name: "int", input: int(-42), want: "-42"},
		{name: "int8", input: int8(-8), want: "-8"},
		{name: "int16", input: int16(-16), want: "-16"},
		{name: "int32", input: int32(-32), want: "-32"},
		{name: "int64", input: int64(-64), want: "-64"},
		{name: "uint", input: uint(42), want: "42"},
		{name: "uint8", input: uint8(8), want: "8"},
		{name: "uint16", input: uint16(16), want: "16"},
		{name: "uint32", input: uint32(32), want: "32"},
		{name: "uint64", input: uint64(64), want: "64"},
		{name: "float32", input: float32(1.5), want: "1.5"},
		{name: "float64", input: float64(2.25), want: "2.25"},
		{name: "stringer", input: testStringer{value: "csv"}, want: "stringer:csv"},
		{name: "default", input: map[string]int{"a": 1}, want: "map[a:1]"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := formatCSVValue(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFormatSQLValue_VariousTypes(t *testing.T) {
	sampleTime := time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name  string
		input any
		want  string
	}{
		{name: "nil", input: nil, want: "NULL"},
		{name: "string", input: "hello", want: "'hello'"},
		{name: "stringEscaped", input: "O'Reilly", want: "'O''Reilly'"},
		{name: "bytesEmpty", input: []byte{}, want: "decode('', 'hex')"},
		{name: "bytes", input: []byte{0xBE, 0xEF}, want: "decode('beef','hex')"},
		{name: "time", input: sampleTime, want: "'2024-01-02T03:04:05Z'"},
		{name: "boolTrue", input: true, want: "TRUE"},
		{name: "boolFalse", input: false, want: "FALSE"},
		{name: "int", input: int(-42), want: "-42"},
		{name: "int8", input: int8(-8), want: "-8"},
		{name: "int16", input: int16(-16), want: "-16"},
		{name: "int32", input: int32(-32), want: "-32"},
		{name: "int64", input: int64(-64), want: "-64"},
		{name: "uint", input: uint(42), want: "42"},
		{name: "uint8", input: uint8(8), want: "8"},
		{name: "uint16", input: uint16(16), want: "16"},
		{name: "uint32", input: uint32(32), want: "32"},
		{name: "uint64", input: uint64(64), want: "64"},
		{name: "float32", input: float32(1.5), want: "1.5"},
		{name: "float64", input: float64(2.25), want: "2.25"},
		{name: "stringer", input: testStringer{value: "sql"}, want: "'stringer:sql'"},
		{name: "default", input: map[string]int{"a": 1}, want: "'map[a:1]'"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := formatSQLValue(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}
