package indexAPI

import (
	"database/sql"
	"log"
	"seekourney/core/database"
	"seekourney/indexing"
	"seekourney/utils"
)

type UnregisteredCollection = utils.UnregisteredCollection

// Collection is a struct that represents a collection of documents.
// Stored in the database.
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

// CollectionFromDB returns a Collection from the database, with the given ID.
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
