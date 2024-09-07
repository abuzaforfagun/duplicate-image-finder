package main

import (
	"bytes"
	"crypto/sha256"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dirName := "./images"
	images, err := getImageFiles(dirName)
	if err != nil {
		log.Panicf("%v", err)
		return
	}

	filesMap := map[[32]byte]string{}
	duplicatedFiles := []string{}

	for _, filename := range images {
		file, err := os.Open(filepath.Join(dirName, filename.Name))
		if err != nil {
			log.Printf("Unable to open [filename=%s][Err=%v]\n", filename.Name, err)
			return
		}

		image, _, err := image.Decode(file)
		if err != nil {
			log.Printf("Unable to decode [filename=%s][Err=%v]\n", filename.Name, err)
			return
		}

		var buffer bytes.Buffer
		if filename.Type == JPG {
			jpeg.Encode(&buffer, image, &jpeg.Options{
				Quality: 5,
			})
		}

		if filename.Type == PNG {
			png.Encode(&buffer, image)
		}
		hash := sha256.Sum256(buffer.Bytes())

		existingFile := filesMap[hash]
		if existingFile != "" {
			duplicatedFiles = append(duplicatedFiles, filename.Name)
			continue
		}
		filesMap[hash] = filename.Name
	}

	log.Println(duplicatedFiles)

}

func getImageFiles(dirPath string) ([]*ImageFile, error) {
	files, err := os.ReadDir(dirPath)

	if err != nil {
		log.Println("unable to retrieve files", err)
		return nil, err
	}

	images := []*ImageFile{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		isJpgFile := strings.HasSuffix(file.Name(), "jpg") || strings.HasSuffix(file.Name(), "jpeg")
		isPngFile := strings.HasSuffix(file.Name(), "png")
		if isJpgFile {
			images = append(images, NewJpegImage(file.Name()))
		}

		if isPngFile {
			images = append(images, NewPngImage(file.Name()))
		}
	}

	return images, nil
}
