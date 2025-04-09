package utils

import (
	"iter"
	"log"
	"os"
	"path/filepath"
)

// TODO: WalkDir is not tested with symbolic/hard links.

// NOTE: WalkDir is also a STDLib function in Go, but this is a custom iterator implementation.

// Constants for the always forbidden directories
func alwaysForbiddenDirs() []string {
	return []string{".", ".."}
}

// Constants for the default forbidden directories
func defaultForbiddenDirs() []string {
	return []string{".git", ".svn", ".hg", ".idea", ".vscode"}
}

// WalkDirConfig is a configuration struct for the WalkDir function.
// It allows setting options for verbosity, returning directories, and filtering by forbidden directories and allowed file extensions.
// If there is no allowed file extension, all files are returned.
type WalkDirConfig struct {
	Verbose       bool
	ForbiddenDirs []string
	AllowedExts   []string
}

// NewWalkDirConfig creates a new WalkDirConfig with default values.
func NewWalkDirConfig() *WalkDirConfig {
	return &WalkDirConfig{
		Verbose:       false,
		ForbiddenDirs: append(alwaysForbiddenDirs(), defaultForbiddenDirs()...),
		AllowedExts:   []string{},
	}
}

// SetVerbose sets the verbose flag for the WalkDirConfig.
func (config *WalkDirConfig) SetVerbose(verbose bool) *WalkDirConfig {
	config.Verbose = verbose
	return config
}

// SetForbiddenDirs sets the forbidden directories for the WalkDirConfig.
func (config *WalkDirConfig) SetForbiddenDirs(dirs []string) *WalkDirConfig {
	config.ForbiddenDirs = append(alwaysForbiddenDirs(), dirs...)
	return config
}

// SetAllowedExts sets the forbidden file extensions for the WalkDirConfig.
// If there are no allowed file extensions, all files are returned.
func (config *WalkDirConfig) SetAllowedExts(exts []string) *WalkDirConfig {
	config.AllowedExts = exts
	return config
}

// WalkDir returns a sequence of file paths in the given directory and its subdirectories.
func (config *WalkDirConfig) WalkDir(path string) iter.Seq[string] {

	return func(yield func(string) bool) {
		config.walkDirIter(yield, path)
	}
}

// walkDirIter is a helper function that recursively walks the directory and yields file paths.
func (config *WalkDirConfig) walkDirIter(yield func(string) bool, path string) {
	entries, err := os.ReadDir(path)

	if err != nil {
		log.Println("Error reading directory:", err)
		return
	}

	for _, entry := range entries {

		if entry.IsDir() {
			// Recursively walk the directory
			subDir := path + string(os.PathSeparator) + entry.Name()

			if !config.forEachDir(yield, subDir) {
				return // Stop iteration if yield returns false
			}

		} else {
			filePath := path + string(os.PathSeparator) + entry.Name()
			if !config.forEachFile(yield, filePath) {
				return // Stop iteration if yield returns false
			}
		}
	}
}

// forEachDir yields each file in the directory and its subdirectories.
func (config *WalkDirConfig) forEachDir(yield func(string) bool, dirPath string) bool {

	if config.Verbose {
		log.Printf("Found directory: %s\n", dirPath)
	}

	if !config.isValidDir(dirPath) {
		if config.Verbose {
			log.Printf("Skipping forbidden directory: %s\n", dirPath)
		}
		return true
	}

	for p := range config.WalkDir(dirPath) {
		if !yield(p) {
			return false
		}
	}

	return true
}

// forEachFile yields a file if it is valid based on the allowed file extensions.
func (config *WalkDirConfig) forEachFile(yield func(string) bool, filePath string) bool {

	if config.Verbose {
		log.Printf("Found file: %s\n", filePath)
	}

	if !config.isValidFile(filePath) {
		if config.Verbose {
			log.Printf("Skipping forbidden file: %s\n", filePath)
		}
		return true
	}

	return yield(filePath)
}

// isValidDir checks if the directory is valid based on the forbidden directories.
func (config *WalkDirConfig) isValidDir(path string) bool {

	name := filepath.Base(path)

	for _, dir := range config.ForbiddenDirs {
		if dir == name {
			return false
		}
	}
	return true
}

// isValidFile checks if the file is valid based on the allowed file extensions.
func (config *WalkDirConfig) isValidFile(path string) bool {

	if len(config.AllowedExts) == 0 {
		return true
	}

	fileExt := filepath.Ext(path)

	for _, ext := range config.AllowedExts {
		if ext == fileExt {
			return true
		}
	}
	return false
}
