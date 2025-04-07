package pdftotext

import "os/exec"
import "os"
import "fmt"
import "github.com/tiagomelo/go-ocr/ocr"
import "path/filepath"
import "regexp"
import "sync"
import "strconv"

var currentPage = 1
var cond = sync.NewCond(&sync.Mutex{})

//converts one pdf to images, replace exec.command later
func pdftoimg(pdf string, outputDir string, imgFormat string) {
	_, err := exec.Command("pdftoppm", imgFormat, pdf, outputDir+"page").CombinedOutput()
	if err != nil {
		fmt.Println("Error running pdftoppm:", err)
	}
}

func clearOutputDir(outputDir string) {
	exec.Command("rm", "-rf", outputDir+"*page-*")
}

func imgToText(image string) string {
	ocr, err := ocr.New()
	if err != nil {
		fmt.Println(err)
	}
	extractedText, err := ocr.TextFromImageFile(image)
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

func imagesToTextAsync(image string, dir string, format string) []string {

	regex, err := regexp.Compile(image + "page-.*")
	var txt []string
	if err != nil {
		fmt.Println("regex is kill", err)
		return txt
	}

	idRegex, err := regexp.Compile(image + "page-(\\d+)\\" + format)
	if err != nil {
		fmt.Println("regex is kill", err)
		return txt
	}

	var wg sync.WaitGroup

	walkHelperHelper := func(path string, info os.FileInfo, wg *sync.WaitGroup) error {
		if regex.MatchString(info.Name()) {
			id := idRegex.FindStringSubmatch(info.Name())[1]
			idInt, erro := strconv.Atoi(id)
			if erro != nil {
				fmt.Println("strconv fail", erro)
				return erro
			}
			toAppend := imgToText(path)
			cond.L.Lock()
			for currentPage != idInt {
				cond.Wait()
			}
			txt = append(txt, toAppend)
			currentPage++
			cond.Broadcast()
			cond.L.Unlock()
		}
		if err != nil {
			fmt.Println("something went wrong when accessing file " + info.Name())
			return err
		}
		wg.Done()
		return nil
	}

	walkHelper := func(path string, info os.FileInfo, err error) error {
		wg.Add(1)
		go walkHelperHelper(path, info, &wg)
		return nil
	}

	err = filepath.Walk(dir, walkHelper)
	{
		if err != nil {
			fmt.Println("something went wrong with filepath walk")
		}
	}
	wg.Wait()

	return txt
}

func Run() {
	prefix := "pdftotxt/"
	pdftoimg(prefix+"pdf/EXAMPLE.pdf", prefix+"covpdf/", "-png") //kör pdftoimg först på din pdf, lägg pdf i pdf folder och byt ut "EXAMPLE" med dess namn
	// test := imgToText("covpdf/page-1.png")
	test := imagesToText("", prefix+"covpdf/")
	// test := imagesToTextAsync("", prefix+"covpdf/", ".png")
	fmt.Println(test)
	clearOutputDir(prefix + "/covpdf/")
}
