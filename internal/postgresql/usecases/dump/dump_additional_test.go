package dump

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTableDumper_BuildQuery(t *testing.T) {
	query, err := buildQuery(&DumpOptions{TableName: "users", Format: "json"})
	require.NoError(t, err)
	assert.Equal(t, "SELECT * FROM \"users\"", query)

	limit := 5
	query, err = buildQuery(&DumpOptions{TableName: "users", Format: "json", Limit: &limit})
	require.NoError(t, err)
	assert.Equal(t, "SELECT * FROM \"users\" LIMIT 5", query)
}

func TestStreamRows_WriterNil(t *testing.T) {
	rows := NewMockRows([]string{"id"}, [][]any{})

	err := streamRows(rows, []string{"id"}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "writerが初期化されていません")
}
