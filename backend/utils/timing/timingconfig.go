package timing

// WatchID is the ID for each used stopwatch.
type WatchID = int

// TODO: Make this a config file instead

// All the different stopwatches
const (
	Search WatchID = iota
	Main
	IndexBytes
	SortWords
	DocFromFile
	ReverseMapLocal
	FolderFromIter
)

// Default creates the default config for the stopwatches.
func Default() WatchesConfig {
	return WatchesConfig{
		Search: StopwatchInfo{
			print: true,
			name:  "Search",
		},
		Main: StopwatchInfo{
			print: true,
			name:  "Main",
		},
		IndexBytes: StopwatchInfo{
			print: false,
			name:  "Indexing Bytes",
		},
		SortWords: StopwatchInfo{
			print: true,
			name:  "Sorting Words",
		},
		DocFromFile: StopwatchInfo{
			print: false,
			name:  "Document From File",
		},
		ReverseMapLocal: StopwatchInfo{
			print: true,
			name:  "Reverse Map Local",
		},
		FolderFromIter: StopwatchInfo{
			print: true,
			name:  "Folder From Dir",
		},
	}
}
