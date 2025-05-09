package document

import (
	"database/sql"
	"encoding/json"
	"log"
	"seekourney/core/database"
	"seekourney/core/normalize"
	"seekourney/indexing"
	"seekourney/utils"
	"seekourney/utils/timing"
	"sort"
	"time"
)

type udoc = indexing.UnnormalizedDocument

type Document struct {
	udoc
	LastIndexed time.Time
}

// Normalize normalizes the document using the provided normalizer
func Normalize(
	doc indexing.UnnormalizedDocument,
	normalizer normalize.Normalizer,
) Document {

	freqMap := make(utils.FrequencyMap)

	for k, v := range doc.Words {
		k = normalizer.Word(k)
		freqMap[k] += v
	}

	return Document{
		udoc: udoc{
			Path:   doc.Path,
			Source: doc.Source,
			Words:  freqMap,
		},
		// SourceID:   ??? TODO or in event loop
		LastIndexed: time.Now(),
	}
}

// Misc

// DebugPrint prints information about the document
func (doc *Document) DebugPrint() {
	log.Printf(
		"Document = {Path: %s, Type: %d, Length: %d}",
		doc.Path,
		doc.Source,
		len(doc.Words),
	)
}

// Pair
type Pair struct {
	Word utils.Word
	Freq utils.Frequency
}

// GetWords returns a slice of pairs of words and their frequency
func (doc *Document) GetWords() []Pair {
	pairs := make([]Pair, 0)

	for k, v := range doc.Words {
		pairs = append(pairs, Pair{k, v})
	}

	return pairs
}

// GetWordsSorted returns a slice of pairs of words and their frequency
// sorted by frequency in descending order
func (doc *Document) GetWordsSorted() []Pair {
	pairs := doc.GetWords()

	t := timing.Measure(timing.SortWords)
	defer t.Stop()

	sort.Slice(
		pairs,
		func(i, j int) bool { return pairs[i].Freq > pairs[j].Freq },
	)
	return pairs
}

// SQL

// SQLGetName returns the name of the table in the database
func (doc Document) SQLGetName() string {
	return "document"
}

// SQLGetFields returns the fields to be inserted into the database
func (doc Document) SQLGetFields() []string {
	return []string{"path", "type", "words"}
}

// SQLGetValues returns the values to be inserted into the database
func (doc Document) SQLGetValues() []any {

	bytes, err := json.Marshal(doc.Words)

	if err != nil {
		log.Printf("Error marshalling dict: %s", err)
		return []database.SQLValue{doc.Path, "file", nil}
	}

	return []database.SQLValue{doc.Path, "file", bytes}
}

// SQLScan scans a row from the database into a Document
func (doc Document) SQLScan(rows *sql.Rows) (Document, error) {
	var path utils.Path
	var source string
	var words []byte

	err := rows.Scan(&path, &source, &words)
	if err != nil {
		return Document{}, err
	}

	var freqMap utils.FrequencyMap

	err = json.Unmarshal(words, &freqMap)
	if err != nil {
		log.Printf("Error unmarshalling dict: %s", err)
		return Document{}, err
	}

	return Document{
		udoc: udoc{
			Path:   path,
			Source: utils.SourceLocal,
			Words:  freqMap,
		},
	}, nil
}

// GetWordCount returns the total number of words in the document
func (doc *Document) GetWordCount() int {
	count := 0
	for _, v := range doc.Words {
		count += int(v)
	}
	return count
}

// CalculateTf calculates the term frequency of a word in the document
// See: https://en.wikipedia.org/wiki/Tf%E2%80%93idf#Term_frequency
func (doc *Document) CalculateTf(word utils.Word) float64 {
	// this will return 0 if the word is not in the document
	freq := doc.Words[word]
	return float64(freq) / float64(doc.GetWordCount())

}

func (doc *Document) UpdateDB(db *sql.DB) error {

	pairs := []string{}
	// Skip the first one, it's the path/primary key
	q1 := database.Update("document").Set(pairs[1:]...)
	query := q1.Where("path=$1")

	log.Printf("Update query: %s", string(query))

	_, err := db.Exec(
		string(query),
		doc.Path,
	)

	return err
}

// DocumentFromDB retrieves a document from the database
func DocumentFromDB(db *sql.DB, path utils.Path) (Document, error) {

	var doc Document

	q1 := database.Select().Queries(doc.SQLGetFields()...)
	query := q1.From("document").Where("path = $1")

	insert := func(res *Document, doc Document) {
		*res = doc
	}

	err := database.ExecScan(
		db,
		string(query),
		&doc,
		insert,
		path,
	)

	return doc, err

}

// DocumentExsitsDB checks if a document exists in the database
func DocumentExsitsDB(db *sql.DB, path utils.Path) (bool, error) {

	// TODO: A better way to do this

	doc, err := DocumentFromDB(db, path)
	if err != nil {
		return false, err
	}

	if doc.Path == "" {
		return false, nil
	}

	if doc.Path != path {
		// This should never happen
		log.Panicf("Document path mismatch: %s != %s", doc.Path, path)
	}

	return true, nil

}
