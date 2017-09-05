package main

import (
	"fmt"
	"flag"
	"os"
	"log"
	"path"
	"image"
	"sync"
	"strings"

	"io/ioutil"
	"path/filepath"
	"image/png"
	"image/jpeg"

	"github.com/nfnt/resize"
)

var wg sync.WaitGroup

func isValid(path string) bool {
	validExts := []string{".png", ".jpg", ".jpeg"}
	fileInfo, err := os.Stat(path)

	if err != nil {
		log.Fatal(err)
	}

	ext := filepath.Ext(path)

	for _, validExt := range validExts {
		if !fileInfo.IsDir() && strings.ToLower(ext) == validExt {
			return true
		}
	}

	return false
}

func processImage(imagePath string, outputPath string, imageHeight int, imageWidth int, ch chan struct{}) error {
	if imageHeight == 0 && imageWidth == 0 {
		imageHeight = 720
	}

	// Add a new struct to our channel. The thing is that our go routines would be blocked if we already reached
	// the maximum amount of data in our channel, ie. 10
	ch <- struct{}{}

	// but of course, we should make sure that we always free up what we added once we're done with this goroutine
	defer func() { <-ch }()
	defer wg.Done()

	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	// check extension, then decode the file according to the ext.
	ext := strings.ToLower(filepath.Ext(imagePath))

	var img image.Image

	switch ext {
	case ".png":
		img, err = png.Decode(file)
	case ".jpg":
		fallthrough
	case ".jpeg":
		img, err = jpeg.Decode(file)
	}

	if err != nil {
		log.Fatal(err)
	}

	// create the new resized image
	newImage := resize.Resize(uint(imageWidth), uint(imageHeight), img, resize.Lanczos3)

	fmt.Printf("Processing %s\n", filepath.Base(file.Name()))
	output, err := os.Create(path.Join(outputPath, filepath.Base(file.Name())))

	defer file.Close()
	defer output.Close()

	if err != nil {
		log.Fatal(err)
	}

	switch ext {
	case ".png":
		png.Encode(output, newImage)
	case ".jpg":
		jpeg.Encode(output, newImage, nil)
	}

	return nil
}

func getFiles(input string) []string {
	var fileList []string

	info, err := os.Stat(input)

	if err != nil {
		log.Fatal(err)
	}

	switch mode := info.Mode(); {
	case mode.IsDir():
		files, err := ioutil.ReadDir(input)

		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			filePath := path.Join(input, file.Name())

			if isValid(filePath) {
				fileList = append(fileList, filePath)
			}
		}
	case mode.IsRegular():
		if isValid(input) {
			fileList = append(fileList, input)
		}
	}

	return fileList
}

func main() {
	// Parse flags
	fileHeight := flag.Int("height", 0, "resize the image to the specified height while retaining aspect ratio")
	fileWidth := flag.Int("width", 0, "resize the image to the specified width while retaining aspect ratio")
	inputPath := flag.String("i", "./", "set input path, could be a folder or a image file")
	outputPath := flag.String("o", "./output/", "set output path")

	flag.Parse()

	// Get images paths, and put them in a list
	files := getFiles(*inputPath)

	// Create output path
	os.Mkdir(*outputPath, 0777)

	// Use goroutines for processing images. Use WaitGroup to wait for all of them to complete
	wg.Add(len(files))

	// Just in case we are trying to process a directory with loads of images, we should only limit ourselves to 10
	// simultaneous processes.
	ch := make(chan struct{}, 10)

	for _, filePath := range files {
		go processImage(filePath, *outputPath, *fileHeight, *fileWidth, ch)
	}

	wg.Wait()
	close(ch)

	fmt.Printf("Files done: %v", files)
}
