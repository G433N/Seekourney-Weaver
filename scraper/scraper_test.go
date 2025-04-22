package scraper

import (
	"path/filepath"
	"regexp"
	"testing"
)

//"file:///home/vilma/Coding/School/OSPP/Seekourney-Weaver/new.html"

func TestWaw(t *testing.T) {
	dir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	wa := regexp.MustCompile(`.*/Seekourney-Weaver/`)

	path := "file://" + wa.FindString(dir) + "testingFiles/"
	newScraper := NewCollector(true, true)
	newScraper.RequestVisitToSite(
		path + "htmlTest1.html",
	)

	newScraper.CollectorRepopulateFixedNumber(1)

	goola := newScraper.ReadFinished()

	println(goola[0])
}
