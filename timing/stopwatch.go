package timing

import (
	"log"
	"time"
)

// TODO: Use global data to have scopes

// Stopwatch is a struct that represents a stopwatch
type Stopwatch struct {
	start time.Time
	name  string
}

func Mesure(name string) Stopwatch {
	log.Printf("Started mesuring %s\n", name)
	return Stopwatch{
		start: time.Now(),
		name:  name,
	}
}

func (s Stopwatch) Stop() {
	elapsed := time.Since(s.start)
	log.Printf("%s took %s\n", s.name, elapsed)
}
