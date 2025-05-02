package indexAPI

import (
	"seekourney/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Temp until we stop using global vars
func resetRegistration() {
	registeredIndexers = make(map[RegisterID]IndexerInfo)
	indexersForFileType = make(map[utils.FileType][]RegisterID)
}

// TODO isUnoccupied and with register

func TestRegisterIndexer(t *testing.T) {
	resetRegistration()
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
	assert.NoError(t, err)

	assert.NotEqual(t, firstIndexerID, secondIndexerID)
}

func TestRegisterIndexerInvalidPort(t *testing.T) {
	resetRegistration()
	const invalidPort = 5000
	_, err := RegisterIndexer(
		"The Best indexer",
		"/home/theorganisation/indexers/the_best",
		"go run main.go",
		[]utils.FileType{"txt"},
		invalidPort,
	)
	assert.Error(t, err)
}

func TestUnregisterIndexer(t *testing.T) {
	resetRegistration()
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

	assert.NoError(t, UnregisterIndexer(firstIndexerID))
	assert.Error(t, UnregisterIndexer(firstIndexerID))

	assert.NoError(t, UnregisterIndexer(secondIndexerID))
}
