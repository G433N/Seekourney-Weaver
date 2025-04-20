package utils

type Word string
type Path string
type Frequency int
type Score int

type FrequencyMap map[Word]Frequency
type ScoreMap map[Path]Score
type ReverseMap map[Word][]Path

// All FileType-s are lower-case letters without dot, e.g. "html" or "md"
type FileType string

// Port number for http request-responses.
// Value for indexing API must be within range
// [indexAPI.MININDEXERPORT, indexAPI.MAXINDEXERPORT]
type Port uint
