package config

import (
	"fmt"
	"os"
)

var (
	GcsBucketName   = getEnv("GCS_BUCKET_NAME")
	GcsBucketFolder = getEnv("GCS_BUCKET_FOLDER")
)

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	panic(fmt.Sprintf("Environment variable %s not set", key))
}
