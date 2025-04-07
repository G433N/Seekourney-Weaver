package utils

type Word string
type Path string
type Frequency int
type Score float64

type FrequencyMap map[Word]Frequency
type ScoreMap map[Path]Score
type ReverseMap map[Word][]Path
