package search

import (
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestgetDocByID(t *testing.T) {
	qry := getDocByID()
	log.Println("getDocByID")
	if hasDbName == false {
		require.NoError(t, errors.New("Database name missing"))
	}

	assert.Equal(t, 201, statusRes)
	if assert.NotNil(t, bodyRes) {
		require.NoError(t, errors.New("Query "))
	}
}
