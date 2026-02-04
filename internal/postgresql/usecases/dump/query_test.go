package dump

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// #==============================================================#
// ##          buildQuery Tests                                  ##
// #==============================================================#

func TestTableDumper_buildQuery_WithoutLimit(t *testing.T) {
	// Arrange
	options := DumpOptions{
		TableName: "users",
		Limit:     nil,
	}

	// Act
	query, err := buildQuery(&options)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "SELECT * FROM \"users\"", query)
}

func TestTableDumper_buildQuery_WithLimit(t *testing.T) {
	// Arrange
	limit := 100
	options := DumpOptions{
		TableName: "users",
		Limit:     &limit,
	}

	// Act
	query, err := buildQuery(&options)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "SELECT * FROM \"users\" LIMIT 100", query)
}
