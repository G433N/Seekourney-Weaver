package timing

type WatchID = int

// All the different stopwatches
const (
	Search WatchID = iota
	Main
	IndexBytes
	SortWords
	DocFromFile
	ReverseMapLocal
	FolderFromDir
)

// Defualt config for the stopwatches
func Default() Config {
	return Config{
		Search: StopwatchInfo{
			active: true,
			name:   "Search",
		},
		Main: StopwatchInfo{
			active: true,
			name:   "Main",
		},
		IndexBytes: StopwatchInfo{
			active: false,
			name:   "Indexing Bytes",
		},
		SortWords: StopwatchInfo{
			active: true,
			name:   "Sorting Words",
		},
		DocFromFile: StopwatchInfo{
			active: false,
			name:   "Document From File",
		},
		ReverseMapLocal: StopwatchInfo{
			active: true,
			name:   "Reverse Map Local",
		},
		FolderFromDir: StopwatchInfo{
			active: true,
			name:   "Folder From Dir",
		},
	}
}
