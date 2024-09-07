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
	"sync"
	"time"
)

func main() {
	start := time.Now()
	findDuplicateImageAsync()
	// findDuplicateImageSync()
	duration := time.Since(start)
	log.Printf("IT takes %f seconds", duration.Seconds())
}

func findDuplicateImageAsync() {

	dirName := "./images"
	images, err := getImageFilesAsync(dirName)
	if err != nil {
		log.Panicf("%v", err)
		return
	}

	filesMap := map[[32]byte]string{}
	duplicatedFiles := []string{}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	for _, filename := range images {
		wg.Add(1)
		go func(fileName *ImageFile) {
			defer wg.Done()
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

			mutex.Lock()
			existingFile := filesMap[hash]
			if existingFile != "" {
				duplicatedFiles = append(duplicatedFiles, filename.Name)
				mutex.Unlock()
				return
			}
			filesMap[hash] = filename.Name
			mutex.Unlock()
		}(filename)

	}
	wg.Wait()

	log.Println(duplicatedFiles)
}

func findDuplicateImageSync() {
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

func getImageFilesAsync(dirPath string) ([]*ImageFile, error) {
	files, err := os.ReadDir(dirPath)

	if err != nil {
		log.Println("unable to retrieve files", err)
		return nil, err
	}

	images := []*ImageFile{}
	imageChan := make(chan *ImageFile)
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(file os.DirEntry) {
			defer wg.Done()
			if file.IsDir() {
				return
			}

			isJpgFile := strings.HasSuffix(file.Name(), "jpg") || strings.HasSuffix(file.Name(), "jpeg")
			isPngFile := strings.HasSuffix(file.Name(), "png")
			if isJpgFile {
				// images = append(images, NewJpegImage(file.Name()))
				imageChan <- NewJpegImage(file.Name())
			}

			if isPngFile {

				// images = append(images, NewPngImage(file.Name()))
				imageChan <- NewPngImage(file.Name())
			}
		}(file)

	}

	go func() {
		wg.Wait()
		close(imageChan)
	}()

	for image := range imageChan {
		images = append(images, image)
	}

	return images, nil
}
