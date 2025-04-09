package timing

import (
	"log"
	"time"
)

// StopwatchInfo is a struct that represents the information of a stopwatch
// It contains the name of the stopwatch and whether it is active or not
type StopwatchInfo struct {
	active bool
	name   string
}

// Config is a map of stopwatch ids to their information
// It is used to configure the stopwatches
type Config map[int]StopwatchInfo

// Global read-only variable to store the configuration
var config Config

// TODO: Use global data to have scopes

// Stopwatch is a struct that represents a stopwatch
type Stopwatch struct {
	start time.Time
	id    int
	text  string
}

// Start a stopwatch
// Id is the id of the stopwatch type, specified in the Init(c Config) function
// Cxt is an optional context string that will be printed with the stopwatch
func Mesure(id int, cxt ...string) *Stopwatch {

	if len(cxt) > 1 {
		log.Fatalf("Context have to many arguments: Max 1, got %d", len(cxt))
	}

	if info, ok := config[id]; !ok || !info.active {
		return &Stopwatch{id: id}
	}

	sw := &Stopwatch{
		start: time.Now(),
		id:    id,
	}

	if len(cxt) > 0 {
		sw.text = "(" + cxt[0] + ")"
	}

	log.Printf("Started mesuring %s %s\n", sw.getName(), sw.text)
	return sw
}

// Stop the stopwatch
// It will print the elapsed time and the context string
// If the stopwatch is not active, it will do nothing
// It is recommended to use a defer statement to stop the stopwatch
func (s *Stopwatch) Stop() {

	if v, ok := config[s.id]; !ok || !v.active {
		return
	}

	elapsed := time.Since(s.start)
	log.Printf("%s took %s %s\n", s.getName(), elapsed, s.text)
}

// Gets the name of the stopwatch
func (s *Stopwatch) getName() string {
	if v, ok := config[s.id]; ok {
		return v.name
	}
	return "Unknown"
}

// Init initializes the timing package
func Init(conf Config) {
	log.Printf("Timing package initialized")
	config = conf
}
