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
		return utils.FILE_SOURCE, nil
	case "dir":
		return utils.DIR_SOURCE, nil
	case "url":
		return utils.URL_SOURCE, nil
	default:
		return 0, errors.New("invalid source type")
	}
}

// SourceTypeToStr converts a SourceType to a string.
func SourceTypeToStr(t SourceType) string {
	switch t {
	case utils.FILE_SOURCE:
		return "file"
	case utils.DIR_SOURCE:
		return "dir"
	case utils.URL_SOURCE:
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
		return utils.DIR_SOURCE, nil
	}

	// TODO: What to do about symlinks and simlar things?
	if stat.Mode().IsRegular() {
		return utils.FILE_SOURCE, nil
	}

	return 0, errors.New("unknown source type")
}
