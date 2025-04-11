package pdftotext

import (
	"fmt"
	"github.com/tiagomelo/go-ocr/ocr"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"seekourney/timing"
)

// converts one pdf to images, replace exec.command later
func pdftoimg(pdf string, outputDir string, imgFormat string) {
	_, err := exec.Command("pdftoppm", imgFormat, pdf, outputDir+"page").CombinedOutput()
	if err != nil {
		fmt.Println("Error running pdftoppm:", err)
	}
}

func clearOutputDir(outputDir string) {
	exec.Command("rm", "-rf", outputDir+"*page-*")
}

func ocrInit() ocr.Ocr {
	ocr, err := ocr.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return ocr

}

var ocrEngine ocr.Ocr = ocrInit()

func imgToText(image string) string {
	sw := timing.Measure(timing.OCRRun, image)
	defer sw.Stop()
	extractedText, err := ocrEngine.TextFromImageFile(image)
	if err != nil {
		fmt.Println(err)
	}
	return extractedText
}

func imagesToText(image string, dir string) []string {

	regex, err := regexp.Compile(image + "page-.*")
	var txt []string
	if err != nil {
		fmt.Println("regex is kill", err)
		return txt
	}

	walkHelper := func(path string, info os.FileInfo, err error) error {
		sw := timing.Measure(timing.PdfWalkHelper)
		defer sw.Stop()
		if regex.MatchString(info.Name()) {
			txt = append(txt, imgToText(path))
		}
		if err != nil {
			fmt.Println("something went wrong when accessing file " + info.Name())
			return err
		}
		return nil
	}

	err = filepath.Walk(dir, walkHelper)
	{
		if err != nil {
			fmt.Println("something went wrong")
		}
	}
	return txt
}

func imagesToTextAsync(image string, dir string) []string {

	regex, err := regexp.Compile(image + "page-.*")
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
				channel <- imgToText(path)
			}(path)

			amount++
		}
		return nil
	}

	err = filepath.Walk(dir, walkHelper)
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
	pdftoimg(prefix+"pdf/EXAMPLE.pdf", prefix+"covpdf/", "-png") //kör pdftoimg först på din pdf, lägg pdf i pdf folder och byt ut "EXAMPLE" med dess namn
	sw.Stop()
	sw = timing.Measure(timing.ImageToText)
	defer sw.Stop()
	// test := imgToText("covpdf/page-1.png")
	// imagesToText("", prefix+"covpdf/")
	imagesToTextAsync("", prefix+"covpdf/")
	// fmt.Println(test)
	// clearOutputDir(prefix + "/covpdf/")
}
