package timing

import (
	"log"
	"time"
)

// StopwatchInfo is a struct that represents the information of a stopwatch.
// It contains the name of the stopwatch and whether it is active or not.
type StopwatchInfo struct {
	print bool
	name  string
}

// WatchesConfig is a map of stopwatch ids to their information
// It is used to configure the stopwatches.
type WatchesConfig map[int]StopwatchInfo

// A global read-only variable to store the configuration.
var config WatchesConfig

// TODO: Use global data to have scopes

// Stopwatch is a struct that represents a stopwatch.
type Stopwatch struct {
	start time.Time
	id    int
	text  string
}

// Measure starts a stopwatch.
// id is the ID of the stopwatch type, specified in the Init(c Config) function.
// cxt is an optional context string(s) that will be printed with the stopwatch.
func Measure(id int, cxt ...string) *Stopwatch {
	sw := &Stopwatch{
		start: time.Now(),
		id:    id,
	}

	if info, ok := config[id]; ok && info.print {
		if len(cxt) > 0 {
			sw.text = "("
			for _, s := range cxt {
				sw.text += s
			}
			sw.text += ")"
		}
		log.Printf("Started mesuring %s %s\n", sw.getName(), sw.text)
	}

	return sw
}

// Stop stops the stopwatch.
// It will print the elapsed time and the context string.
// If the stopwatch is not active, it will do nothing.
// It is recommended to use a defer statement to stop the stopwatch.
func (s *Stopwatch) Stop() time.Duration {
	elapsed := time.Since(s.start)
	if v, ok := config[s.id]; ok && v.print {
		log.Printf("%s took %s %s\n", s.getName(), elapsed, s.text)
	}

	return elapsed
}

// getName gets the name of the stopwatch.
func (s *Stopwatch) getName() string {
	if v, ok := config[s.id]; ok {
		return v.name
	}
	return "Unknown"
}

// Init initializes the timing package.
func Init(conf WatchesConfig) {
	log.Printf("Timing package initialized")
	config = conf
}
