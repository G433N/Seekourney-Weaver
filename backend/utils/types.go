package utils

import "strconv"

// Query is a string containing plain words separated by spaces.
// E.g. "1.24.2 golang documentation".
type Query string

// Word is a plain text word.
type Word string

// Path is a file path or web address.
type Path string

// Frequency is the number of occurrences that a Word appears.
type Frequency int

// Score is a number representing how relevant a given Word is when searching.
type Score float64

// Normalizer is a type that represents a normalizer.
type Normalizer int

const (
	ToLower Normalizer = iota
	Stemming
)

const AMOUNT_NORMALIZERS = 1

// Source denotes the type of source indexed.
// E.g. a local file or a web page.
type Source int

const (
	// Source is the source of the document
	// SourceLocal is a local file
	SourceLocal Source = iota
	// SourceWeb is a web page
	SourceWeb
)

// TODO: Should probably use utils.Source instead of SourceType or rename it

// SourceType is an enumeration of the different source types.
type SourceType int

const (
	FileSource SourceType = iota
	DirSource
	UrlSource
)

// FrequencyMap gives the frequency of a given word.
type FrequencyMap map[Word]Frequency

// ScoreMap gives the relevance score of a given document
// with respect to a current search query.
type ScoreMap map[Path]Score

// ReverseMap gives all paths to files that contain the given word.
type ReverseMap map[Word][]Path

// WordFrequencyMap gives the frequency of a specific word
// for every document that contains it.
// Used when searching.
type WordFrequencyMap map[Path]Frequency

// SearchResult is information about a single document
// with respect to a current search query.
type SearchResult struct {
	Path   Path
	Score  Score
	Source Source
}

// SearchResponse is the format an HTTP search response
// from Core has after unmarshalling JSON.
type SearchResponse struct {
	Query   Query
	Results []SearchResult
}

// Result is a tuple used when handling database data.
type Result[T any] struct {
	Value T
	Err   error
}

// FileType-s are lower-case letters without dot, e.g. "html" or "md".
type FileType string

// Port is a port number used in http request-responses to-from
// indexers, Core, and DB.
// Value for indexing API must be within range
// [MININDEXERPORT, MAXINDEXERPORT].
type Port uint

func (p Port) String() string {
	return strconv.Itoa(int(p))
}

// Address including port acting as endpoint for http request.
// E.g. "http://localhost:39010".
type Endpoint string

type ObjectId string

// IndexerID is a unique identifier for an indexer.
type IndexerID ObjectId

type UnregisteredCollection struct {
	// Root path / start of reqursive indexing
	Path Path

	// Indexer used to index this collection
	IndexerID IndexerID

	// Type of source
	SourceType SourceType

	// If true, the indexer will index Recursivevly
	Recursive bool

	// If false will always index when reindexing is requested
	RespectLastModified bool

	// What function to normalize all documents with
	Normalfunc Normalizer
}
