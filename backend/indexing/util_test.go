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

func TestSettingsIntoURL(t *testing.T) {

	settings := &Settings{
		Path:      utils.Path("test"),
		Type:      FileSource,
		Recursive: true,
		Parrallel: false,
	}

	expected := "http://localhost:1234/index" +
		"?path=test&type=file&recursive=true&parallel=false"

	url, err := settings.IntoURL(1234)
	assert.NoError(t, err)
	assert.Equal(t, expected, url)

}
