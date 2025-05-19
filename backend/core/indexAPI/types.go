package indexAPI

import (
	"database/sql"
	"log"
	"os/exec"
	"seekourney/core/database"
	"seekourney/indexing"
	"seekourney/utils"
	"testing"
	"time"
)

// TODO: Structured log messages and struct
// TODO: Log should be in response body not query

func (indexer *IndexerData) start() (*RunningIndexer, error) {
	// TODO: Ping to make sure it actually started
	// TODO: Ping to make sure it is not already running
	// TODO: Use timeout instead of sleep
	// Basically what startIndexerAlready did
	args := indexer.Args
	// Hack to let us run ls command when testing to mock starting up indexer.
	if !testing.Testing() {
		args = append(indexer.Args, "--port="+indexer.Port.String())
	}

	execCmd := exec.Command(indexer.ExecPath, args...)

	// TODO: Handle output
	execCmd.Stdout = nil
	execCmd.Stderr = nil

	err := execCmd.Start()
	if err != nil {
		log.Fatalf("Failed to start command: %v", err)
	}

	log.Printf("Starting indexer with command: %s %s\n", indexer.ExecPath, args)

	time.Sleep(1 * time.Second)

	return &RunningIndexer{
		ID:   indexer.ID,
		Exec: execCmd,
		Port: indexer.Port,
	}, nil
}

type UnregisteredCollection = utils.UnregisteredCollection

type Collection struct {
	UnregisteredCollection
	ID indexing.CollectionID
}

// SQL

// SQLGetName returns the name of the table in the database
func (col Collection) SQLGetName() string {
	return "collection"
}

// SQLGetFields returns the fields to be inserted into the database
func (col Collection) SQLGetFields() []string {
	return []string{
		"id",
		"path",
		"indexer_id",
		"recursive",
		"source_type",
		"respect_last_modified",
		"normalizer",
	}
}

// SQLGetValues returns the values to be inserted into the database
func (col Collection) SQLGetValues() []any {

	return []database.SQLValue{
		col.ID,
		col.Path,
		col.IndexerID,
		col.Recursive,
		indexing.SourceTypeToStr(col.SourceType),
		col.RespectLastModified,
		col.Normalfunc,
	}
}

// SQLScan scans a row from the database into a Document
func (col Collection) SQLScan(rows *sql.Rows) (Collection, error) {
	var id indexing.CollectionID
	var path utils.Path
	var indexerID IndexerID
	var recursive bool
	var sourceType string
	var respectLastModified bool
	var normalizer utils.Normalizer

	err := rows.Scan(
		&id,
		&path,
		&indexerID,
		&recursive,
		&sourceType,
		&respectLastModified,
		&normalizer,
	)
	if err != nil {
		return Collection{}, err
	}

	sourceTypeEnum, err := indexing.StrToSourceType(sourceType)
	if err != nil {
		log.Printf("Error parsing source type: %s", err)
		return Collection{}, err
	}

	return Collection{
		UnregisteredCollection{
			Path:                path,
			IndexerID:           indexerID,
			SourceType:          sourceTypeEnum,
			Recursive:           recursive,
			RespectLastModified: respectLastModified,
			Normalfunc:          normalizer,
		},
		id,
	}, nil

}

func CollectionFromDB(
	db *sql.DB,
	id indexing.CollectionID,
) (Collection, error) {
	var colloction Collection

	q1 := database.Select().QueryAll()
	query := q1.From(colloction.SQLGetName()).Where("id = $1")

	insert := func(res *Collection, col Collection) {
		*res = col
	}

	err := database.ExecScan(
		db,
		string(query),
		&colloction,
		insert,
		id,
	)

	return colloction, err

}

// TODO move to better place
// RegisterCollection creates a Collection from an UnregisteredCollection
// and adds it to the database, making it available to index.
func RegisterCollection(
	db *sql.DB,
	ureqCol UnregisteredCollection,
) (Collection, error) {
	id := database.GenerateId()
	collection := Collection{
		UnregisteredCollection: ureqCol,
		ID:                     indexing.CollectionID(id),
	}

	_, err := database.InsertInto(db, collection)
	if err != nil {
		log.Printf("Error inserting collection: %s", err)
		return Collection{}, err
	}

	return collection, nil
}
