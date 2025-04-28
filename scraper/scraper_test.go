package scraper

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

//"file:///home/vilma/Coding/School/OSPP/Seekourney-Weaver/new.html"

func TestLongLocalWikipediaHtml(t *testing.T) {
	dir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	wa := regexp.MustCompile(`.*/`)
	path := "file://" + wa.FindString(dir) + "testingFiles/"
	newScraper := NewCollector(true, true)
	newScraper.RequestVisitToSite(
		path + "htmlTest1.html",
	)

	newScraper.CollectorRepopulateFixedNumber(1)

	readStrings := newScraper.ReadFinished()
	var fullString string

	for _, x := range readStrings {
		fullString += x

	}

	iterator := strings.SplitSeq(fullString, "	")
	fullString = ""
	for x := range iterator {
		if x != "" {
			fullString += x + " "
		}
	}
	iterator = strings.SplitSeq(fullString, "\n")
	fullString = ""
	for x := range iterator {
		if x != "" {
			fullString += x + " "
		}
	}
	iterator = strings.SplitSeq(fullString, " ")
	fullString = ""
	for x := range iterator {
		if x != "" {
			fullString += x + " "
		}
	}
	testStrings := []string{
		"The fruit of typical cultivars of cucumber is roughly cylindrical",
		"preliminary research to identify whether" +
			" cucumbers are able to deter herbivores and",
		"Herbivore defense",
		"Description",
		"References",
		"Seekourney-Weaver/testingFiles/htmlTest1.html",
		"Cucumbers grown to eat fresh are called slicing cucumbers.",
		"Shoots Cucumber shoots are regularly consumed as a vegetable," +
			" especially in rural areas. In Thailand they are often served" +
			" with a crab meat sauce." +
			" They can also be stir fried or used in soups.",
	}
	for _, x := range testStrings {
		if !strings.Contains(fullString, x) {
			t.Error("did not contain: ", x)
		}
	}
}

func TestLocalLinkHopping(t *testing.T) {
	dir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	wa := regexp.MustCompile(`.*/`)

	path := "file://" + wa.FindString(dir) + "testingFiles/"
	newScraper := NewCollector(true, true)
	newScraper.RequestVisitToSite(
		path + "htmlTest2.html",
	)
	newScraper.CollectorRepopulateFixedNumber(4)
	txt := strings.Join(newScraper.ReadFinished(), " ")
	txt += strings.Join(newScraper.ReadFinished(), " ")
	txt += strings.Join(newScraper.ReadFinished(), " ")
	txt += strings.Join(newScraper.ReadFinished(), " ")

	if !strings.Contains(txt, "Index.html") {
		t.Error("did not find or did not parse correctly file:",
			" HtmlTest2",
		)
	}
	if !strings.Contains(txt, "Child Page One") {
		t.Error("did not find or did not parse correctly file:",
			"HtmlTest2Child/one.html",
		)
	}
	if !strings.Contains(txt, "Child Page Two") {
		t.Error("did not find or did not parse correctly file:",
			"HtmlTest2Child/two.html",
		)
	}
	if !strings.Contains(txt, "Child Page Three") {
		t.Error("did not find or did not parse correctly file:",
			"HtmlTest2Child/three.html",
		)
	}

}
