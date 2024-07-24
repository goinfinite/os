package cliMiddleware

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
	"slices"

	"github.com/joho/godotenv"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

var requiredEnvVars = []string{
	"PRIMARY_VHOST",
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
	primaryHostname, err := infraHelper.GetPrimaryVirtualHost()
	if err != nil {
		log.Fatalf("PrimaryHostnameUnidentifiable")
	}

	envFilePath := "/speedia/.env"

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
