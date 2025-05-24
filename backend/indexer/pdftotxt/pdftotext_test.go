package pdftotext

import "testing"
import "path/filepath"
import "slices"

/*
func TestPDFtoimgOnePage(t *testing.T){
	pdftoimg("./pdf/sample.pdf", "./covpdf/sample", "-png")
	file_exists, err := filepath.Glob("./covpdf/samplepage-1.png");
	if !slices.Contains(file_exists, "./covpdf/samplepage-1.png") || err != nil{
		t.Errorf(`pdftoimg = %q, %v, want "./covpdf/page-1.png", error`, file_exists, err)
	}
	clearOutputDir("./covpdf/")
}

func TestPDFtoimgNoPDF(t *testing.T){
	pdftoimg("./pdf/doesntexist.pdf", "./covpdf/", "-png")
	file_exists, err := filepath.Glob("./covpdf/page*");
	if slices.Contains(file_exists, "covpdf/page*") || err != nil{
		t.Errorf(`pdftoimg = %q, %v, want "", error`, file_exists, err)
	}
	clearOutputDir("./covpdf/")
}
*/

func TestPDFtoimgMultiplePages(t *testing.T){
	pdftoimg("./pdf/sample_multiple_pages.pdf", "./covpdf/")
	file_exists, err := filepath.Glob("./covpdf/page*");
	if !slices.Contains(file_exists, "covpdf/page-1.png") || err != nil{
		t.Errorf(`pdftoimg = %q, %v, want "covpdf/page-1.jpeg", error`, file_exists, err)
	}
	if !slices.Contains(file_exists, "covpdf/page-2.png"){
		t.Errorf(`pdftoimg = %q, %v, want "covpdf/page-2.jpeg", error`, file_exists, err)
	}
	if !slices.Contains(file_exists, "covpdf/page-3.png"){
		t.Errorf(`pdftoimg = %q, %v, want "covpdf/page-3.jpeg", error`, file_exists, err)
	}
	clearOutputDir("./covpdf/")
	
}
/*
func TestTest(t *testing.T){
	pdftoimg("./pdf/sample_book.pdf", "./covpdf/book", "-png")
	//sw := timing.Measure(timing.ImageToText)
	//defer sw.Stop()
	//fmt.Println(imagesToTextAsync("book", "./covpdf/"))
	clearOutputDir("./covpdf/")
}
	*/

func TestImgtotext(t *testing.T){
}

func TestImgtotextNoImg(t *testing.T){
}

func TestImgtotextNoText(t *testing.T){
}

func testImgTotextMultipledifferentimages(t *testing.T){

}

func TestImagestotextOneImage(t *testing.T){
}

func TestImagestotextNoImage(t *testing.T){
}

func TestImagestotextMultipleImages(t *testing.T){

}

