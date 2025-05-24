package pdftotext

import (
	"fmt"
	"github.com/tiagomelo/go-ocr/ocr"
	"os"
	"path/filepath"
	"regexp"
	"seekourney/timing"
	"github.com/gen2brain/go-fitz"
	"image/jpeg"
	"seekourney/utils"
)

/*
pdftoimg
Converts a pdf to a series of jpegs where each page is its own image using mupdf

# Parameters:

# Returns:
  - nothing

# Example usage:

pdftoimg("~/pdf/EXAMPLE.pdf", "~/covpdf/example", "-png") 

*/

func pdftoimg(pdfpath utils.Path, outputDir utils.Path) {
	doc, err := fitz.New(string(pdfpath))
	if err != nil {
		panic(err)
	}
	defer doc.Close()

	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.Image(n)
		if err != nil {
			panic(err)
		}

		f, err := os.Create(filepath.Join(string(outputDir), fmt.Sprintf("page-%d.jpeg", n+1)))
		if err != nil {
			panic(err)
		}

		err = jpeg.Encode(f, img, &jpeg.Options{jpeg.DefaultQuality})
		if err != nil {
			panic(err)
		}

		f.Close()
	}

}

/*
clearOutputDir
Clears an output directory of all files with prefix + "page-" in them.

# Parameters:
  - outputDir string
The directory and the prefix to clear as a string

# Returns:
  - nothing
*/
func clearOutputDir(outputDir utils.Path) {
	toRemove,_ := filepath.Glob(string(outputDir) + "*page-*")
	for _, file := range toRemove {
		os.Remove(file)
	}
}

/*
ocrInit
Initialises an ocr for image to text recognition.

# Parameters:
  - none

# Returns:
  - ocr.Ocr
an initialised ocr
*/
func ocrInit() ocr.Ocr {
	ocr, err := ocr.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return ocr
}

var ocrEngine ocr.Ocr = ocrInit()

/*
imgToText
Uses Tesseract to recognise text from an image and returns all text as a string.


# Parameters:
  - image string
The path to the image to analyse as a string.

# Returns:
  - string
All text found in the image.
*/
func imgToText(image utils.Path) string {
	sw := timing.Measure(timing.OCRRun, string(image))
	defer sw.Stop()
	extractedText, err := ocrEngine.TextFromImageFile(string(image))
	if err != nil {
		fmt.Println(err)
	}
	return extractedText
}

/*
imagesToText
Converts multiple images from a directory to text

# Returns:
  - type
desc.
*/
func imagesToText(inputDir utils.Path, outputDir utils.Path) []string {

	regex, err := regexp.Compile(string(inputDir) + "page-.*")
	var txt []string
	if err != nil {
		fmt.Println("regex is kill", err)
		return txt
	}

	walkHelper := func(path string, info os.FileInfo, err error) error {
		sw := timing.Measure(timing.PdfWalkHelper)
		defer sw.Stop()
		if regex.MatchString(info.Name()) {
			txt = append(txt, imgToText(utils.Path(path)))
		}
		if err != nil {
			fmt.Println("something went wrong when accessing file " + info.Name())
			return err
		}
		return nil
	}

	err = filepath.Walk(string(outputDir), walkHelper)
	{
		if err != nil {
			fmt.Println("something went wrong")
		}
	}
	return txt
}

/*
imagesToTextParallel
desc.

newline continued desc.

# Parameters:

# Returns:
  - type
desc.
*/
func imagesToTextParallel(image utils.Path, dir utils.Path) []string {

	regex, err := regexp.Compile(string(image) + "page-.*")
	var txt []string
	if err != nil {
		fmt.Println("regex is kill", err)
		return txt
	}

	channel := make(chan string)
	amount := 0

	walkHelper := func(path string, info os.FileInfo, err error) error {
		if regex.MatchString(info.Name()) {
			go func(path string) {
				channel <- imgToText(utils.Path(path))
			}(path)

			amount++
		}
		return nil
	}

	err = filepath.Walk(string(dir), walkHelper)
	{
		if err != nil {
			fmt.Println("something went wrong with filepath walk")
		}
	}

	for range amount {
		result := <-channel
		txt = append(txt, result)
	}

	return txt
}

func Run() {
	prefix := "pdftotxt/"
	sw := timing.Measure(timing.PfdToImage)
	pdftoimg(utils.Path(prefix+"pdf/EXAMPLE.pdf"), utils.Path(prefix+"covpdf/")) //kör pdftoimg först på din pdf, lägg pdf i pdf folder och byt ut "EXAMPLE" med dess namn
	sw.Stop()
	sw = timing.Measure(timing.ImageToText)
	defer sw.Stop()
	// test := imgToText("covpdf/page-1.png")
	// imagesToText("", prefix+"covpdf/")
	imagesToTextParallel("", utils.Path(prefix+"covpdf/"))
	// fmt.Println(test)
	// clearOutputDir(prefix + "/covpdf/")
}
