package main

import (
	"fmt"
	"github.com/gen2brain/go-fitz"
	"github.com/tiagomelo/go-ocr/ocr"
	"image/jpeg"
	"os"
	"path/filepath"
	"regexp"
	"seekourney/utils"
)

// Text is a type alias for a string that represents extracted text
type Text string

/*
pdftoimg
Converts a pdf to a series of jpegs where each page is its own image using mupdf
*/
func pdftoimg(pdfpath utils.Path, outputDir utils.Path) error {
	doc, err := fitz.New(string(pdfpath))
	if err != nil {
		return err
	}
	_, err = os.Stat(string(outputDir))
	if os.IsNotExist(err) {
		err = os.MkdirAll(string(outputDir), 0777)
		//0777 is the file permission for the directory created
		if err != nil {
			err2 := doc.Close()
			if err2 != nil {
				return err2
			}
			return err
		}
	}

	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.Image(n)
		if err != nil {
			err2 := doc.Close()
			if err2 != nil {
				return err2
			}
			return err
		}

		file, err := os.Create(filepath.Join(string(outputDir),
			fmt.Sprintf("page-%d.jpeg", n+1)))

		if err != nil {
			err2 := file.Close()
			if err2 != nil {
				return err2
			}
			err2 = doc.Close()
			if err2 != nil {
				return err2
			}
			return err
		}

		err = jpeg.Encode(file, img,
			&jpeg.Options{Quality: jpeg.DefaultQuality})
		if err != nil {
			err2 := file.Close()
			if err2 != nil {
				return err2
			}
			err2 = doc.Close()
			if err2 != nil {
				return err2
			}
			return err
		}

		err = file.Close()
		if err != nil {
			err2 := doc.Close()
			if err2 != nil {
				return err2
			}
			return err
		}
	}
	err2 := doc.Close()
	if err2 != nil {
		return err2
	}
	return nil
}

/*
clearOutputDir
Clears an output directory of all files with prefix + "page-" in them.
*/
func clearOutputDir(outputDir utils.Path) error {
	toRemove, err := filepath.Glob(string(outputDir) + "*page-*")
	if err != nil {
		return err
	}
	for _, file := range toRemove {
		err = os.Remove(file)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
ocrInit
Initialises an ocr for image to text recognition.
*/
func ocrInit() (ocr.Ocr, error) {
	ocr, err := ocr.New()
	if err != nil {
		return nil, err
	}
	return ocr, nil
}

var ocrEngine ocr.Ocr

/*
imgToText
Uses Tesseract to recognise text from an image and returns all text as a string.
*/
func imgToText(image utils.Path) (Text, error) {
	extractedText, err := ocrEngine.TextFromImageFile(string(image))
	if err != nil {
		return "", err
	}
	return Text(extractedText), nil
}

/*
imagesToText
Converts multiple images from a directory to text
*/
func imagesToText(inputDir utils.Path, outputDir utils.Path) ([]Text, error) {
	regex, err := regexp.Compile(string(inputDir) + "page-.*")
	var txt []Text
	if err != nil {
		return txt, err
	}
	ocrEngine, err = ocrInit()
	if err != nil {
		return txt, err
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(string(outputDir), 0777)
		//0777 is the file permission for the directory created
		if err != nil {
			return txt, err
		}
	}
	walkHelper := func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return fmt.Errorf("file info is nil for path: %s", path)
		}
		if regex.MatchString(info.Name()) {
			var newText Text
			newText, err = imgToText(utils.Path(path))
			if err != nil {
				return err
			}
			txt = append(txt, newText)
		}
		if err != nil {
			return err
		}
		return nil
	}
	err = filepath.Walk(string(outputDir), walkHelper)
	{
		if err != nil {
			return txt, err
		}
	}
	return txt, nil
}

/*
imagesToTextParallel
Converts multiple images from a directory to text in parallel.
*/
func imagesToTextParallel(image utils.Path, outputDir utils.Path) ([]Text,
	error) {

	regex, err := regexp.Compile(string(image) + "page-.*")
	var txt []Text
	if err != nil {
		return txt, err
	}
	ocrEngine, err = ocrInit()
	if err != nil {
		return txt, err
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(string(outputDir), 0777)
		//0777 is the file permission for the directory created
		if err != nil {
			return txt, err
		}
	}

	channel := make(chan utils.Result[Text])
	amount := 0

	walkHelper := func(path string, info os.FileInfo, err error) error {
		if regex.MatchString(info.Name()) {
			go func() {
				var newText Text
				newText, err = imgToText(utils.Path(path))
				channel <- utils.Result[Text]{
					Value: newText,
					Err:   err,
				}
			}()

			amount++
		}
		return nil
	}

	err = filepath.Walk(string(outputDir), walkHelper)
	{
		if err != nil {
			return txt, err
		}
	}

	for range amount {
		result := <-channel
		if result.Err != nil {
			return txt, result.Err
		}
		txt = append(txt, result.Value)
	}

	return txt, nil
}
