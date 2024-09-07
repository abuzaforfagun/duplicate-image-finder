package main

import (
	"bytes"
	"crypto/sha256"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

func main() {
	dirName := "./images"
	jpegFiles, pngFiles, err := getImageFiles(dirName)
	if err != nil {
		log.Panicf("%v", err)
		return
	}

	filesMap := map[[32]byte]string{}
	duplicatedFiles := []string{}

	for _, filename := range jpegFiles {
		file, err := os.Open(dirName + "/" + filename)
		if err != nil {
			log.Printf("Unable to open [filename=%s][Err=%v]\n", filename, err)
			return
		}

		image, _, err := image.Decode(file)
		if err != nil {
			log.Printf("Unable to decode [filename=%s][Err=%v]\n", filename, err)
			return
		}

		var buffer bytes.Buffer
		jpeg.Encode(&buffer, image, &jpeg.Options{
			Quality: 5,
		})
		hash := sha256.Sum256(buffer.Bytes())

		existingFile := filesMap[hash]
		if existingFile != "" {
			duplicatedFiles = append(duplicatedFiles, filename)
			continue
		}
		filesMap[hash] = filename
	}

	for _, filename := range pngFiles {
		file, err := os.Open(dirName + "/" + filename)
		if err != nil {
			log.Printf("Unable to open [filename=%s][Err=%v]\n", filename, err)
			return
		}

		image, _, err := image.Decode(file)
		if err != nil {
			log.Printf("Unable to decode [filename=%s][Err=%v]\n", filename, err)
			return
		}

		var buffer bytes.Buffer
		png.Encode(&buffer, image)
		hash := sha256.Sum256(buffer.Bytes())

		existingFile := filesMap[hash]
		if existingFile != "" {
			duplicatedFiles = append(duplicatedFiles, filename)
			continue
		}
		filesMap[hash] = filename
	}

	log.Println(duplicatedFiles)

}

func getImageFiles(dirPath string) ([]string, []string, error) {
	files, err := os.ReadDir(dirPath)

	if err != nil {
		log.Println("unable to retrieve files", err)
		return nil, nil, err
	}

	jpegImages := []string{}
	pngImages := []string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		isJpgFile := strings.HasSuffix(file.Name(), "jpg") || strings.HasSuffix(file.Name(), "jpeg")
		isPngFile := strings.HasSuffix(file.Name(), "png")
		if isJpgFile {
			jpegImages = append(jpegImages, file.Name())
		}

		if isPngFile {
			pngImages = append(pngImages, file.Name())
		}
	}

	return jpegImages, pngImages, nil
}
