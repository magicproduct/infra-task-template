package verification

import (
	"encoding/json"
	"infra-task-solution/pkg/gcs"
	"os"
	"path/filepath"
	"regexp"
)

type FileVerificationResult struct {
	Exists   bool `json:"exists"`
	Scrubbed bool `json:"scrubbed"`
}

type Result struct {
	Files map[string]FileVerificationResult `json:"files"`
}

func VerifyFiles(localDirectory string) (bool, error) {
	files, err := os.ReadDir(localDirectory)
	if err != nil {
		return false, err
	}

	result := Result{
		Files: make(map[string]FileVerificationResult),
	}

	allVerified := true
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		localFilePath := filepath.Join(localDirectory, file.Name())
		exists, scrubbed := verifyFile(localFilePath)

		result.Files[file.Name()] = FileVerificationResult{
			Exists:   exists,
			Scrubbed: scrubbed,
		}

		if !exists || !scrubbed {
			allVerified = false
		}
	}

	if err := uploadVerificationResult(result); err != nil {
		return false, err
	}

	return allVerified, nil
}

func verifyFile(filePath string) (bool, bool) {
	gcsData, err := gcs.DownloadFile(filePath)
	if err != nil {
		return false, false
	}

	return true, !containsPII(gcsData)
}

func containsPII(data string) bool {
	piiPattern := regexp.MustCompile(`@user_\d+`)
	return piiPattern.MatchString(data)
}

func uploadVerificationResult(result Result) error {
	jsonData, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return gcs.UploadFile(string(jsonData), "verification.json")
}
