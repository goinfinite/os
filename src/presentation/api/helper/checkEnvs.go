package restApiHelper

import (
	"os"
)

func CheckEnvs() {
	envVars := []string{
		"JWT_SECRET",
		"UAK_SECRET",
	}

	for _, key := range envVars {
		value := os.Getenv(key)
		if value == "" {
			panic("MissingEnvVar: " + key)
		}
	}
}
