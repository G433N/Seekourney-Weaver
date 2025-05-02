package indexing

import (
	"errors"
	"os"
	"seekourney/utils"
)

// TODO: Should probably use utils.Source instead of SourceType or rename it

// SourceType is an enumeration of the different source types.
type SourceType int

const (
	FileSource SourceType = iota
	DirSource
	UrlSource
)

// StrToSourceType converts a string to a SourceType.
func StrToSourceType(str string) (SourceType, error) {
	switch str {
	case "file":
		return FileSource, nil
	case "dir":
		return DirSource, nil
	case "url":
		return UrlSource, nil
	default:
		return 0, errors.New("invalid source type")
	}
}

// SourceTypeToStr converts a SourceType to a string.
func SourceTypeToStr(t SourceType) string {
	switch t {
	case FileSource:
		return "file"
	case DirSource:
		return "dir"
	case UrlSource:
		return "url"
	default:
		return "unknown"
	}
}

// SourceTypeFromPath determines the source type from a path.
// By looking it up in the local filesystem.
// used for simplifying the tui code
func SourceTypeFromPath(path utils.Path) (SourceType, error) {

	stat, err := os.Stat(string(path))
	if err != nil {
		return 0, err
	}

	if stat.IsDir() {
		return DirSource, nil
	}

	// TODO: What to do about symlinks and simlar things?
	if stat.Mode().IsRegular() {
		return FileSource, nil
	}

	return 0, errors.New("unknown source type")
}
