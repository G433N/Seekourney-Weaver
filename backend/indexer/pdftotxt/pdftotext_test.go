package main

import (
	"fmt"
	"path/filepath"
	"seekourney/utils"
	"slices"
	"testing"
)

func TestPDFtoimgOnePage(t *testing.T) {
	err := pdftoimg("./pdf/sample.pdf", "./covpdf/sample")
	if err != nil {
		t.Errorf("pdftoimg failed: %v", err)
	}
	file_exists, err := filepath.Glob("./covpdf/sample/page-1.jpeg")
	if !slices.Contains(file_exists, "./covpdf/sample/page-1.jpeg") ||
		err != nil {
		t.Errorf(`pdftoimg = %q, %v, want "./covpdf/page-1.jpeg", error`,
			file_exists, err)
	}
	err = clearOutputDir("./covpdf/sample_multiple_pages/")
	if err != nil {
		t.Errorf("clearOutputDir failed: %v", err)
	}
}

func TestPDFtoimgNoPDF(t *testing.T) {
	err := pdftoimg("./pdf/doesntexist.pdf", "./covpdf/")
	if err == nil {
		t.Errorf("pdftoimg should have returned an error for non-existent PDF")
	}
	file_exists, err := filepath.Glob("./covpdf/page*")
	if slices.Contains(file_exists, "covpdf/page*") || err != nil {
		t.Errorf(`pdftoimg = %q, %v, want "", error`, file_exists, err)
	}
}

func TestPDFtoimgMultiplePages(t *testing.T) {
	err := pdftoimg("./pdf/sample_multiple_pages.pdf",
		"./covpdf/sample_multiple_pages")
	if err != nil {
		t.Errorf("pdftoimg failed: %v", err)
	}
	file_exists, err := filepath.Glob("./covpdf/sample_multiple_pages/page*")
	if !slices.Contains(file_exists,
		"covpdf/sample_multiple_pages/page-1.jpeg") || err != nil {
		t.Errorf(`pdftoimg = %q, %v, want "covpdf/page-1.jpeg", error`,
			file_exists, err)
	}
	if !slices.Contains(file_exists,
		"covpdf/sample_multiple_pages/page-2.jpeg") {
		t.Errorf(`pdftoimg = %q, %v, want "covpdf/page-2.jpeg", error`,
			file_exists, err)
	}
	if !slices.Contains(file_exists,
		"covpdf/sample_multiple_pages/page-3.jpeg") {
		t.Errorf(`pdftoimg = %q, %v, want "covpdf/page-3.jpeg", error`,
			file_exists, err)
	}
	err = clearOutputDir("./covpdf/sample_multiple_pages/")
	if err != nil {
		t.Errorf("clearOutputDir failed: %v", err)
	}
}

func TestImgtotext(t *testing.T) {
	// Test with a single image
	if testing.Short() {
		t.Skip("skipping long test")
	}
	err := pdftoimg("./pdf/sample.pdf", "./covpdf/sample")
	if err != nil {
		t.Errorf("pdftoimg failed: %v", err)
	}
	text, err := imagesToText("", utils.Path("./covpdf/sample/"))
	if err != nil {
		t.Errorf("imgToText failed: %v", err)
	}
	if text == nil {
		t.Error("imgToText returned empty text")
	}
	fmt.Printf("Extracted text from single image: %s\n", text)
	err = clearOutputDir("./covpdf/sample_multiple_pages/")
	if err != nil {
		t.Errorf("clearOutputDir failed: %v", err)
	}
}
