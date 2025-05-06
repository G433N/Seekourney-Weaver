package indexAPI

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"seekourney/core/normalize"
	"seekourney/indexing"
	"seekourney/utils"
	"strings"
	"time"
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

func (id IndexerID) GetPort() utils.Port {
	port := utils.MININDEXERPORT + utils.Port(id)
	if !indexing.IsValidPort(port) {
		panic("Invalid port number: to many indexers registered")
	}
	return port
}

type IndexerData struct {
	ID       IndexerID
	Name     string
	ExecPath string
	Args     []string
}

func (cmd StartUpCMD) appendPort(port utils.Port) StartUpCMD {
	cmd.args = append(cmd.args, "--port="+port.String())
	return cmd
}

func (indexer *IndexerData) start() *ActiveIndexer {
	args := append(indexer.Args, "--port="+indexer.ID.GetPort().String())

	execCmd := exec.Command(indexer.ExecPath, args...)

	out, err := os.Create("a.out")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}

	execCmd.Stdout = out
	execCmd.Stderr = out

	err = execCmd.Start()
	if err != nil {
		log.Fatalf("Failed to start command: %v", err)
	}

	log.Printf("Starting indexer with command: %s %s\n", indexer.ExecPath, args)

	time.Sleep(1 * time.Second)

	return &ActiveIndexer{
		ID:   indexer.ID,
		Exec: execCmd,
	}
}

type ActiveIndexer struct {
	ID   IndexerID
	Exec *exec.Cmd
}

func (indexer *ActiveIndexer) GetRequest(args ...string) (string, error) {

	port := indexer.ID.GetPort()

	url := _ENDPOINTPREFIX_ + port.String() + "/" + strings.Join(args, "/")

	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New("indexer did not respond to request: " + err.Error())
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New("indexer did not respond to request, " +
			"alternatively did not respond with ok statuscode")
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("indexer did not respond to request")
	}

	respStr := string(respByte)

	return string(respStr), nil
}

func (indexer *ActiveIndexer) Wait() error {
	return indexer.Exec.Wait()
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
