package indexAPI

import (
	"seekourney/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidPort(t *testing.T) {
	assert.False(t, isValidPort(42))
	assert.False(t, isValidPort(MININDEXERPORT-1))
	assert.True(t, isValidPort(MININDEXERPORT))
	const middleOfRange utils.Port = MININDEXERPORT +
		((MAXINDEXERPORT - MININDEXERPORT) / 2)
	assert.True(t, isValidPort(middleOfRange))
	assert.True(t, isValidPort(MAXINDEXERPORT))
	assert.False(t, isValidPort(MAXINDEXERPORT+1))
}

// TODO isUnoccupied and with register

func TestRegisterIndexer(t *testing.T) {
	firstIndexerID, err := RegisterIndexer(
		"The text indexer",
		"/home/theorganisation/indexers/text",
		"textindexingprogram.exe",
		[]utils.FileType{"txt"},
		39499,
	)
	assert.Nil(t, err)

	secondIndexerID, err := RegisterIndexer(
		"TheSuperOPIndexer",
		"/home/TSOPIDev/TheSuperOPIndexer",
		"go run main.go",
		[]utils.FileType{"md", "txt", "csv"},
		39498)
	assert.Nil(t, err)

	assert.NotEqual(t, firstIndexerID, secondIndexerID)
}

func TestRegisterIndexerInvalidPort(t *testing.T) {
	const invalidPort = 5000
	_, err := RegisterIndexer(
		"The Best indexer",
		"/home/theorganisation/indexers/the_best",
		"go run main.go",
		[]utils.FileType{"txt"},
		invalidPort,
	)
	assert.NotNil(t, err)
}

func TestUnregisterIndexer(t *testing.T) {
	firstIndexerID, _ := RegisterIndexer(
		"The text indexer",
		"/home/textindex/textindexer/program.exe",
		"program.exe",
		[]utils.FileType{"txt"},
		39499)

	secondIndexerID, _ := RegisterIndexer(
		"TheSuperOPIndexer",
		"/home/TSOPIDev/TheSuperOPIndexer/startup_script.sh",
		"go run main.go",
		[]utils.FileType{"md", "txt", "csv"},
		39498)

	assert.Nil(t, UnregisterIndexer(firstIndexerID))
	assert.NotNil(t, UnregisterIndexer(firstIndexerID))

	assert.Nil(t, UnregisterIndexer(secondIndexerID))
}
