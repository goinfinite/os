package testHelpers

import (
	"fmt"
	"path"
	"runtime"

	"github.com/joho/godotenv"
)

func LoadEnvVars() {
	_, fileDirectory, _, _ := runtime.Caller(0)
	rootDir := path.Dir(path.Dir(path.Dir(fileDirectory)))
	testEnvPath := rootDir + "/.env"

	loadEnvErr := godotenv.Load(testEnvPath)
	if loadEnvErr != nil {
		panic(fmt.Errorf("Error loading .env file: %s", loadEnvErr))
	}
}
