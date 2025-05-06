package indexAPI

import (
	"os/exec"
	"seekourney/core/normalize"
	"seekourney/indexing"
	"seekourney/utils"
)

// TODO: Structured log messages and struct
// TODO: Log should be in response body not query

type StartUpCMD struct {
	path utils.Path
	args []string
}

func CMDFromString(str string) StartUpCMD {
	return StartUpCMD{
		path: utils.Path(str),
		args: []string{},
	}
}

type IndexerID uint

type IndexerData struct {
	ID       IndexerID
	Name     string
	ExecPath string
}

func (indexer *IndexerData) GetPort() utils.Port {
	port := utils.MININDEXERPORT + utils.Port(indexer.ID)
	if !indexing.IsValidPort(port) {
		panic("Invalid port number: to many indexers registered")
	}
	return port
}

type ActiveIndexer struct {
	ID   IndexerID
	Exec *exec.Cmd
}

type SourceCollectionID uint

// TODO: Rename
type SourceCollecton struct {
	ID                  SourceCollectionID
	Path                utils.Path
	IndexerID           IndexerID
	Recursive           bool
	RespectLastModified bool
	Normalfunc          normalize.Normalizer
}
