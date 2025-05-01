package utils

type Query string

type Word string
type Path string
type Frequency int
type Score float64

type Source int

const (
	// Source is the source of the document
	// SourceLocal is a local file
	SourceLocal Source = iota
	// SourceWeb is a web page
	SourceWeb
)

type FrequencyMap map[Word]Frequency
type ScoreMap map[Path]Score
type ReverseMap map[Word][]Path

// WordFrequencyMap is maps paths to their frequency for every document
type WordFrequencyMap map[Path]Frequency

type SearchResult struct {
	Path   Path
	Score  Score
	Source Source
}

type SearchResponse struct {
	Query   Query
	Results []SearchResult
}

type Result[T any] struct {
	Value T
	Err   error
}

// All FileType-s are lower-case letters without dot, e.g. "html" or "md"
type FileType string

// Port is a port number used in http request-responses to-from
// indexers, Core, and DB.
// Value for indexing API must be within range
// [MININDEXERPORT, MAXINDEXERPORT]
type Port uint
