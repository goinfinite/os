package restApiHelper

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"

	"golang.org/x/exp/slices"
)

var requiredEnvVars = []string{
	"VIRTUAL_HOST",
	"JWT_SECRET",
	"UAK_SECRET",
}

var envVarsToGenerateIfEmpty = []string{
	"JWT_SECRET",
	"UAK_SECRET",
}

func genSecret() (string, error) {
	secretLength := 32
	secretBytes := make([]byte, secretLength)

	_, err := rand.Read(secretBytes)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(secretBytes), nil
}

func CheckEnvs() {
	file, err := os.OpenFile(".env", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("EnvOpenFileError: %v", err)
	}
	defer file.Close()

	for _, key := range requiredEnvVars {
		value := os.Getenv(key)
		if value != "" {
			continue
		}

		if !slices.Contains(envVarsToGenerateIfEmpty, key) {
			log.Fatalf("MissingEnvVar: %s", key)
		}

		value, err = genSecret()
		if err != nil {
			log.Fatalf("GenSecretError: %v", err)
		}

		os.Setenv(key, value)

		_, err = file.WriteString(key + "=" + value + "\n")
		if err != nil {
			log.Fatalf("EnvWriteFileError: %v", err)
		}
	}
}
