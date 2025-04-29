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
)

type Document indexing.UnnormalizedDocument

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
		Path:   doc.Path,
		Source: doc.Source,
		Words:  freqMap,
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
		Path:   path,
		Source: utils.SourceLocal,
		Words:  freqMap,
	}, nil
}

func (doc *Document) GetWordCount() int {
	count := 0
	for _, v := range doc.Words {
		count += int(v)
	}
	return count
}

func (doc *Document) CalculateTf(word utils.Word) float64 {
	if _, ok := doc.Words[word]; !ok {
		return 0
	}
	return float64(doc.Words[word]) / float64(doc.GetWordCount())

}

func DocumentFromDB(db *sql.DB, path utils.Path) (Document, error) {

	var doc Document

	query := database.Select().Queries(doc.SQLGetFields()...).From("document").Where("path = $1")

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
