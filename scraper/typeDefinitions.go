package scraper

import (
	"regexp"
	"seekourney/utils/Sync"
	"sync"

	"github.com/gocolly/colly/v2"
)

const (

	// Used for enabling debug prints
	_DEBUG_ = false
	// Total amount of worspaces in the collector
	_WORKSPACES_ = 5
	// Key used to access workspace ID in a map
	_IDKEY_ = `WorkspaceID`
	_HOST_  = `HostPath`

	_PRIOQUEUEMAXLEN_ = 100
	_STACKMAXLEN_     = 1_000_000
)

type (
	WorkspaceID  int
	URLString    string
	HtmlFileType string
	hostPath     string
	innerPath    string

	CollectorStruct struct {

		// struct holding all context
		// to make the inteface with the collector as simple as possible
		context context

		// used to keep track of number of workspaces in use
		counter counter

		// the colly collector used for webb scraping and formatting
		collectorColly *colly.Collector

		settings settings
	}
	context struct {
		linkHandler *linkHandler

		// worspace for the collectors async requests
		// each request gets their own index
		// and then appends the text they recieve to the slice
		workspaceBuffer [_WORKSPACES_][]string

		// channel of indexes in the `workspace` array
		// ready to be assigned to a request
		// (buffer of size 'QUEUEMAXLEN')
		emptyIndexes chan WorkspaceID

		//  channel of indexes in the `workspace` array ready to be read
		// (buffer of size 'QUEUEMAXLEN')
		finishedIndexes chan WorkspaceID
	}

	counter struct {

		// mutex used to sync changes to the two counters in context
		counterLock sync.Mutex

		// currently working on amount
		workingCounter int

		// currently finished amount should be in sync with len(finishedIndexes)
		finishedCounter int
	}

	settings struct {
		HtmlFileType HtmlFileType

		async bool
	}
	//------------------------------------------------------------//
	// _ _ _ _ _ _ _ _ _ _ _ _Link Handler_ _ _ _ _ _ _ _ _ _ _ _ //
	//------------------------------------------------------------//

	linkHandler struct {
		filter       filter
		storageStack Sync.Stack[*URLCompact]

		priorityQueue Sync.CyclicQueue[*URLCompact]

		inputChan chan linkInputWrap

		outputChan chan URLString

		outputSem Sync.Semaphore

		storedSem       Sync.Semaphore
		quit            bool
		handlersWorking sync.WaitGroup
	}

	URLCompact struct {
		webFileBool bool
		host        Sync.ArrayPlusIndex
		inner       innerPath
	}

	filter struct {
		webhosts  Sync.ArrayPlus[hostPath]
		filterMap map[hostPath]filterMapInner
	}
	filterMapInner struct {
		webFileBool bool
		index       Sync.ArrayPlusIndex
		filterMap   map[innerPath]bool
	}

	linkInputWrap struct {
		prio bool
		URL  *URLCompact
	}
)

//------------------------------------------------------------//
// _ _ _ _ _ _ _ _ _ _ _Regular Expresions_ _ _ _ _ _ _ _ _ _ //
//------------------------------------------------------------//

var (
	//Matches when the link does have the same host as the current webbsite
	shortendLinkRegex = regexp.MustCompile(`^/[a-zA-Z]`)
	//Parts of wikipedia not worth indexing
	WikipediaBadRegex = regexp.MustCompile(
		`/wiki/(File|Wikipedia|Special|User)|/static/|/w/`)
	//Matches when the link doesn't have the same host as the current webbsite
	fullLocalLinkRegex = regexp.MustCompile(`^file://`)
	fullWebLinkRegex   = regexp.MustCompile(`^https://`)
)
