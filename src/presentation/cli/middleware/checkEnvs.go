package cliMiddleware

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
	"slices"

	infraHelper "github.com/goinfinite/os/src/infra/helper"
	"github.com/joho/godotenv"
)

var requiredEnvVars = []string{
	"ACCOUNT_API_KEY_SECRET",
	"ACCOUNT_SECURE_ACCESS_KEY_SECRET",
	"JWT_SECRET",
	"PRIMARY_VHOST",
}

var envVarsToGenerateIfEmpty = []string{
	"ACCOUNT_API_KEY_SECRET",
	"ACCOUNT_SECURE_ACCESS_KEY_SECRET",
	"JWT_SECRET",
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
	primaryHostname, err := infraHelper.GetPrimaryVirtualHost()
	if err != nil {
		log.Fatalf("PrimaryHostnameUnidentifiable")
	}

	envFilePath := "/infinite/.env"

	envFile, err := os.OpenFile(envFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0400)
	if err != nil {
		log.Fatalf("EnvOpenFileError: %v", err)
	}
	defer envFile.Close()

	err = godotenv.Load(envFilePath)
	if err != nil {
		log.Fatalf("EnvLoadError: %v", err)
	}

	for _, key := range requiredEnvVars {
		value := os.Getenv(key)
		if value != "" {
			continue
		}

		if key == "PRIMARY_VHOST" {
			value = primaryHostname.String()
			os.Setenv(key, value)
			continue
		}

		if !slices.Contains(envVarsToGenerateIfEmpty, key) {
			log.Fatalf("MissingEnvVar: %s", key)
		}

		value, err = genSecret()
		if err != nil {
			log.Fatalf("GenSecretError: %v", err)
		}

		_, err = envFile.WriteString(key + "=" + value + "\n")
		if err != nil {
			log.Fatalf("EnvWriteFileError: %v", err)
		}

		os.Setenv(key, value)
	}
}
