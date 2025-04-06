package timing

import (
	"log"
	"time"
)

type StopwatchInfo struct {
	active bool
	name   string
}

type Config map[int]StopwatchInfo

var config Config

// TODO: Use global data to have scopes

// Stopwatch is a struct that represents a stopwatch
type Stopwatch struct {
	start time.Time
	id    int
	text  string
}

func Mesure(id int, cxt ...string) *Stopwatch {

	if len(cxt) > 1 {
		log.Fatalf("Context have to many arguments: Max 1, got %d", len(cxt))
	}

	if v, ok := config[id]; !ok || !v.active {
		return &Stopwatch{id: id}
	}

	s := &Stopwatch{
		start: time.Now(),
		id:    id,
	}

	if len(cxt) > 0 {
		s.text = "(" + cxt[0] + ")"
	}

	log.Printf("Started mesuring %s %s\n", s.getName(), s.text)
	return s
}

func (s *Stopwatch) Stop() {

	if v, ok := config[s.id]; !ok || !v.active {
		return
	}

	elapsed := time.Since(s.start)
	log.Printf("%s took %s %s\n", s.getName(), elapsed, s.text)
}

func (s *Stopwatch) getName() string {
	if v, ok := config[s.id]; ok {
		return v.name
	}
	return "Unknown"
}

func Init(c Config) {
	log.Printf("Timing package initialized")
	config = c
}
