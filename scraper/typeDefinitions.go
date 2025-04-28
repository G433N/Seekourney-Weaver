package scraper

import (
	"regexp"
	"seekourney/utils/Sync"
	"sync"

	"github.com/gocolly/colly/v2"
)

const (

	// Used for enabling debug prints
	_DEBUG_ = true
	// Total amount of worspaces in the collector
	_WORKSPACES_ = 5
	// Buffer size of the link and priority link channels
	_LINKQUEUELEN_ = 20
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

		// buffered queue channel of links scraped from previously visited sites
		linkQueue chan URLString

		// buffered queue channel of links inputed to the scraper
		priorityLinkQueue chan URLString

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
		storageStack storageStack

		priorityQueue PriorityQueue

		inputChan chan linkInputWrap

		outputChan chan URLString

		outputSem Sync.Semaphore

		storedSem       Sync.Semaphore
		quit            bool
		handlersWorking sync.WaitGroup
	}

	URLCompact struct {
		webFileBool bool
		host        int
		inner       innerPath
	}

	filter struct {
		filehosts []hostPath
		webhosts  []hostPath
		filterMap map[hostPath]filterMapInner
	}
	filterMapInner struct {
		webFileBool bool
		index       int
		filterMap   map[innerPath]bool
	}

	linkInputWrap struct {
		prio bool
		URL  *URLCompact
	}

	PriorityQueue struct {
		lock sync.Mutex

		Queue [_PRIOQUEUEMAXLEN_]*URLCompact

		read int

		write int

		len int
	}
	storageStack struct {
		lock sync.Mutex

		Stack []*URLCompact
	}
)

//------------------------------------------------------------//
// _ _ _ _ _ _ _ _ _ _ _Regular Expresions_ _ _ _ _ _ _ _ _ _ //
//------------------------------------------------------------//

var (
	//Matches when the link is in the same main domain as the current webbsite
	shortendLinkRegex = regexp.MustCompile(`^/[a-zA-Z]`)
	//Parts of wikipedia not worth indexing
	WikipediaBadRegex = regexp.MustCompile(
		`/wiki/(File|Wikipedia|Special|User)|/static/|/w/`)
)
