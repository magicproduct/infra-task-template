package processing

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type FileContent struct {
	FileName    string
	Content     string
	OriginalErr error
}

func FilterPII(directory string, cleanedFile chan<- FileContent) {
	var wg sync.WaitGroup

	usernamePattern := regexp.MustCompile(`@user_\d+`)

	dataSetFiles, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println("Could not read dataset:", err)
		close(cleanedFile)
		return
	}

	for _, file := range dataSetFiles {
		if !file.IsDir() {
			wg.Add(1)

			go func(entry os.DirEntry) {
				defer wg.Done()

				filePath := filepath.Join(directory, entry.Name())

				file, err := os.Open(filePath)
				if err != nil {
					cleanedFile <- FileContent{FileName: filePath, OriginalErr: err}
					return
				}
				defer func(file *os.File) {
					err := file.Close()
					if err != nil {
						log.Printf("Failed to close file: %v", err)
					}
				}(file)

				data, err := io.ReadAll(file)
				if err != nil {
					cleanedFile <- FileContent{FileName: filePath, OriginalErr: err}
					return
				}

				updatedData := usernamePattern.ReplaceAllString(string(data), "<FILTERED>")
				cleanedFile <- FileContent{FileName: filePath, Content: updatedData}
			}(file)
		}
	}

	go func() {
		wg.Wait()
		close(cleanedFile)
	}()
}
