package indexing

import (
	"seekourney/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	_MINPORT_ utils.Port = utils.MININDEXERPORT
	_MAXPORT_ utils.Port = utils.MAXINDEXERPORT
)

func TestIsValidPort(t *testing.T) {
	assert.False(t, IsValidPort(42))
	assert.False(t, IsValidPort(_MINPORT_-1))
	assert.True(t, IsValidPort(_MINPORT_))
	const middleOfRange utils.Port = _MINPORT_ +
		((_MAXPORT_ - _MINPORT_) / 2)
	assert.True(t, IsValidPort(middleOfRange))
	assert.True(t, IsValidPort(_MAXPORT_))
	assert.False(t, IsValidPort(_MAXPORT_+1))
}
