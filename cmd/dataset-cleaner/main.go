package main

import (
	"fmt"
	"infra-task-solution/pkg/gcs"
	"infra-task-solution/pkg/processing"
	"infra-task-solution/pkg/verification"
	"sync"
)

func main() {
	datasetFolder := "assets/reddit_dataset"
	processed := make(chan processing.FileContent)

	var waitGroup sync.WaitGroup

	go processing.FilterPII(datasetFolder, processed)

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		waitGroup.Wait()
		uploadProcessedFiles(processed)
	}()

	waitGroup.Wait()

	allVerified, err := verification.VerifyFiles(datasetFolder)
	if err != nil {
		fmt.Printf("Verification failed: %v\n", err)
	} else if allVerified {
		fmt.Println("Verification completed successfully")
	} else {
		fmt.Println("Not all files are verified successfully")
	}
}

func uploadProcessedFiles(processed <-chan processing.FileContent) {
	for file := range processed {
		if file.OriginalErr != nil {
			fmt.Printf("Error processing file %s: %v\n", file.FileName, file.OriginalErr)
			continue
		}

		err := gcs.UploadFile(file.Content, file.FileName)
		if err != nil {
			fmt.Printf("Could not upload file %s: %v\n", file.FileName, err)
			return
		}
	}
}
