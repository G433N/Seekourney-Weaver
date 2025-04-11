package utils

type Word string
type Path string
type Frequency int
type Score int

type FrequencyMap map[Word]Frequency
type ScoreMap map[Path]Score
type ReverseMap map[Word][]Path

// 0 means no indexing type set
type TypeOfIndexing int
