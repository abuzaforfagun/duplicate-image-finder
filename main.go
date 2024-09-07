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
	var folderPath string
	if len(os.Args) != 2 {
		log.Println("Using current directory to find duplicate")
		folderPath = "./"
	} else {
		folderPath = os.Args[1]

		log.Printf("Using '%s' folder to find the duplicate\n", folderPath)
	}

	start := time.Now()
	duplicatedFiles := findDuplicateImageAsync(folderPath)

	for _, file := range duplicatedFiles {
		log.Println(file)
	}
	duration := time.Since(start)
	log.Printf("IT takes %f seconds to find '%d' duplicate images\n", duration.Seconds(), len(duplicatedFiles))
}

func findDuplicateImageAsync(folderPath string) []string {
	files, err := getImageFilesAsync(folderPath)
	if err != nil {
		log.Panicf("%v", err)
		return nil
	}

	filesMap := map[[32]byte]string{}
	duplicatedFiles := []string{}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	for _, filename := range files {
		wg.Add(1)
		go func(fileName *ImageFile) {
			defer wg.Done()
			file, err := os.Open(filename.Name)
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
				duplicatedFiles = append(duplicatedFiles, existingFile)
				mutex.Unlock()
				return
			}
			filesMap[hash] = filename.Name
			mutex.Unlock()
		}(filename)

	}
	wg.Wait()
	return duplicatedFiles
}

func findDuplicateImageSync(folderPath string) []string {
	images, err := getImageFiles(folderPath)
	if err != nil {
		log.Panicf("%v", err)
		return nil
	}

	filesMap := map[[32]byte]string{}
	duplicatedFiles := []string{}

	for _, filename := range images {
		file, err := os.Open(filepath.Join(folderPath, filename.Name))
		if err != nil {
			log.Printf("Unable to open [filename=%s][Err=%v]\n", filename.Name, err)
			return nil
		}

		image, _, err := image.Decode(file)
		if err != nil {
			log.Printf("Unable to decode [filename=%s][Err=%v]\n", filename.Name, err)
			return nil
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

	return duplicatedFiles
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
			imagePath := filepath.Join(dirPath, file.Name())
			if isJpgFile {
				imageChan <- NewJpegImage(imagePath)
			}

			if isPngFile {
				imageChan <- NewPngImage(imagePath)
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
