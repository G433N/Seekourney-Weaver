package indexing

import (
	"errors"
	"os"
	"seekourney/utils"
)

type SourceType = utils.SourceType

// StrToSourceType converts a string to a SourceType.
func StrToSourceType(str string) (SourceType, error) {
	switch str {
	case "file":
		return utils.FileSource, nil
	case "dir":
		return utils.DirSource, nil
	case "url":
		return utils.UrlSource, nil
	default:
		return 0, errors.New("invalid source type")
	}
}

// SourceTypeToStr converts a SourceType to a string.
func SourceTypeToStr(t SourceType) string {
	switch t {
	case utils.FileSource:
		return "file"
	case utils.DirSource:
		return "dir"
	case utils.UrlSource:
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
		return utils.DirSource, nil
	}

	// TODO: What to do about symlinks and simlar things?
	if stat.Mode().IsRegular() {
		return utils.FileSource, nil
	}

	return 0, errors.New("unknown source type")
}
