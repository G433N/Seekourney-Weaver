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
	FolderFromIter
	PfdToImage
	ImageToText
	PdfWalkHelper
	OCRNew
	OCRRun
)

// Defualt config for the stopwatches
func Default() Config {
	return Config{
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
		PfdToImage: StopwatchInfo{
			print: true,
			name:  "PDF to Image",
		},
		ImageToText: StopwatchInfo{
			print: true,
			name:  "Image to Text",
		},
		PdfWalkHelper: StopwatchInfo{
			print: true,
			name:  "PDF Walk Helper",
		},
		OCRNew: StopwatchInfo{
			print: true,
			name:  "OCR New",
		},
		OCRRun: StopwatchInfo{
			print: true,
			name:  "OCR Run",
		},
	}
}
