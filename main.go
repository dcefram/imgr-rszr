package main

import (
	"fmt"
	"flag"
	"os"
	"io/ioutil"
	"log"
	"path/filepath"
	"path"
	"image"
	"image/png"
	"image/jpeg"
	"github.com/nfnt/resize"
	"sync"
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
		if !fileInfo.IsDir() && ext == validExt {
			return true
		}
	}

	return false
}

func processImage(imagePath string, outputPath string, imageHeight int, imageWidth int) error {
	if imageHeight == 0 && imageWidth == 0 {
		imageHeight = 720
	}

	defer wg.Done()

	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	// decode the image, first check ext
	ext := filepath.Ext(imagePath)

	var img image.Image

	switch ext {
	case ".png":
		img, err = png.Decode(file)
	case ".jpg":
		img, err = jpeg.Decode(file)
	}

	if err != nil {
		log.Fatal(err)
	}

	newImage := resize.Resize(uint(imageWidth), uint(imageHeight), img, resize.Lanczos3)

	fmt.Printf("Processing %s\n", filepath.Base(file.Name()))
	output, err := os.Create(path.Join(outputPath, filepath.Base(file.Name())))

	file.Close()

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

func main() {
	fileHeight := flag.Int("height", 0, "resize the image to the specified height while retaining aspect ratio")
	fileWidth := flag.Int("width", 0, "resize the image to the specified width while retaining aspect ratio")
	inputPath := flag.String("i", "./", "set input path, could be a folder or a image file")
	outputPath := flag.String("o", "./output/", "set output path")
	flag.Parse()

	inputInfo, err := os.Stat(*inputPath)

	if err != nil {
		log.Fatal(err)
	}

	var files []string

	switch inputMode := inputInfo.Mode(); {
	case inputMode.IsDir():
		dirFiles, err := ioutil.ReadDir(*inputPath)

		if err != nil {
			log.Fatal(err)
		}

		for _, file := range dirFiles {
			filePath := path.Join(*inputPath, file.Name())

			if isValid(filePath) {
				files = append(files, filePath)
			}
		}
	case inputMode.IsRegular():
		if isValid(*inputPath) {
			files = append(files, *inputPath)
		}
	}

	// Create output path
	os.Mkdir(*outputPath, 0777)

	for _, filePath := range files {
		wg.Add(1)
		go processImage(filePath, *outputPath, *fileHeight, *fileWidth)
	}

	wg.Wait()

	fmt.Printf("Files done: %v", files)
}
