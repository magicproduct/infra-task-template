package gcs

import (
	"context"
	"io"
	"log"
	"path/filepath"
	"sync"

	"cloud.google.com/go/storage"
	"infra-task-solution/pkg/config"
)

var (
	client *storage.Client
	once   sync.Once
)

func getClient() (*storage.Client, error) {
	var err error
	once.Do(func() {
		client, err = storage.NewClient(context.Background())
	})
	return client, err
}

func UploadFile(fileContent, remoteFilePath string) error {
	gcsClient, err := getClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	writer := gcsClient.Bucket(config.GcsBucketName).Object(filepath.Join(config.GcsBucketFolder, remoteFilePath)).NewWriter(ctx)
	if _, err = writer.Write([]byte(fileContent)); err != nil {
		return err
	}

	log.Printf("Uploaded file to path %s", remoteFilePath)
	return writer.Close()
}

func DownloadFile(remoteFilePath string) (string, error) {
	gcsClient, err := getClient()
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	reader, err := gcsClient.Bucket(config.GcsBucketName).Object(filepath.Join(config.GcsBucketFolder, remoteFilePath)).NewReader(ctx)
	if err != nil {
		return "", err
	}
	defer func(reader *storage.Reader) {
		err := reader.Close()
		if err != nil {
			log.Printf("Failed to close GCS reader: %v", err)
		}
	}(reader)

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
