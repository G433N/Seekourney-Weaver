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
	ReturnDirs    bool
	ForbiddenDirs []string
	AllowedExts   []string
}

// NewWalkDirConfig creates a new WalkDirConfig with default values.
func NewWalkDirConfig() *WalkDirConfig {
	return &WalkDirConfig{
		Verbose:       false,
		ReturnDirs:    false,
		ForbiddenDirs: append(alwaysForbiddenDirs(), defaultForbiddenDirs()...),
		AllowedExts:   []string{},
	}
}

// SetVerbose sets the verbose flag for the WalkDirConfig.
func (c *WalkDirConfig) SetVerbose(verbose bool) *WalkDirConfig {
	c.Verbose = verbose
	return c
}

// SetReturnDirs sets the returnDirs flag for the WalkDirConfig.
func (c *WalkDirConfig) SetReturnDirs(returnDirs bool) *WalkDirConfig {
	c.ReturnDirs = returnDirs
	return c
}

// SetForbiddenDirs sets the forbidden directories for the WalkDirConfig.
func (c *WalkDirConfig) SetForbiddenDirs(dirs []string) *WalkDirConfig {
	c.ForbiddenDirs = append(alwaysForbiddenDirs(), dirs...)
	return c
}

// SetAllowedExts sets the forbidden file extensions for the WalkDirConfig.
// If there are no allowed file extensions, all files are returned.
func (c *WalkDirConfig) SetAllowedExts(exts []string) *WalkDirConfig {
	c.AllowedExts = exts
	return c
}

// WalkDir returns a sequence of file paths in the given directory and its subdirectories.
func (c *WalkDirConfig) WalkDir(path string) iter.Seq[string] {

	return func(yield func(string) bool) {

		entries, err := os.ReadDir(path)

		if err != nil {
			log.Println("Error reading directory:", err)
			return
		}

		for _, entry := range entries {

			if entry.IsDir() {
				// Recursively walk the directory
				subDir := path + string(os.PathSeparator) + entry.Name()

				if c.Verbose {
					log.Printf("Walking directory: %s\n", subDir)
				}

				if !c.isValidDir(subDir) {
					if c.Verbose {
						log.Printf("Skipping forbidden directory: %s\n", subDir)
					}
					continue
				}

				for p := range c.WalkDir(subDir) {
					if !yield(p) {
						return
					}
				}
			} else {
				// Yield the file path
				filePath := path + string(os.PathSeparator) + entry.Name()

				if c.Verbose {
					log.Printf("Found file: %s\n", filePath)
				}

				if !c.isValidFile(filePath) {
					if c.Verbose {
						log.Printf("Skipping forbidden file: %s\n", filePath)
					}
					continue
				}

				if !yield(filePath) {
					return
				}
			}
		}
	}
}

// isValidDir checks if the directory is valid based on the forbidden directories.
func (c *WalkDirConfig) isValidDir(path string) bool {

	name := filepath.Base(path)

	for _, dir := range c.ForbiddenDirs {
		if dir == name {
			return false
		}
	}
	return true
}

// isValidFile checks if the file is valid based on the allowed file extensions.
func (c *WalkDirConfig) isValidFile(path string) bool {

	if len(c.AllowedExts) == 0 {
		return true
	}

	fileExt := filepath.Ext(path)

	for _, ext := range c.AllowedExts {
		if ext == fileExt {
			return true
		}
	}
	return false
}
