package testHelpers

import (
	"encoding/base64"
	"fmt"
	"path"
	"runtime"

	"github.com/joho/godotenv"
)

func GenerateString(desiredSize int) string {
	desiredByteSize := uint((float64(desiredSize) / 4) * 3)
	randomBytes := make([]byte, desiredByteSize)
	return base64.StdEncoding.EncodeToString(randomBytes)
}

func LoadEnvVars() {
	_, fileDirectory, _, _ := runtime.Caller(0)
	rootDir := path.Dir(path.Dir(path.Dir(fileDirectory)))
	testEnvPath := rootDir + "/.env"

	loadEnvErr := godotenv.Load(testEnvPath)
	if loadEnvErr != nil {
		panic(fmt.Errorf("Error loading .env file: %s", loadEnvErr))
	}
}
