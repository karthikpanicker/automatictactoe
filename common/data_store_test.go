package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsFieldSelectedInTrelloDetails(t *testing.T) {
	td := new(TrelloDetails)
	td.FieldsToUse = []string{"hello", "world"}
	value := td.IsFieldSelected("hello")
	assert.Equal(t, value, true)
	value = td.IsFieldSelected("universe")
	assert.Equal(t, value, false)
}
