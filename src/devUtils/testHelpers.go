package testHelpers

import (
	"encoding/base64"

	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	cliMiddleware "github.com/goinfinite/os/src/presentation/cli/middleware"
	"github.com/joho/godotenv"
)

func GenerateString(desiredSize int) string {
	desiredSizeBytesLength := float64(desiredSize) * 3
	desiredSizeStringLength := desiredSizeBytesLength / 4
	randomBytes := make([]byte, uint(desiredSizeStringLength))
	return base64.StdEncoding.EncodeToString(randomBytes)
}

func LoadEnvVars() {
	cliMiddleware.CheckEnvs()

	loadEnvErr := godotenv.Load("/infinite/.env")
	if loadEnvErr != nil {
		panic("LoadingEnvFileError: " + loadEnvErr.Error())
	}
}

func GetPersistentDbSvc() *internalDbInfra.PersistentDatabaseService {
	persistentDbSvc, err := internalDbInfra.NewPersistentDatabaseService()
	if err != nil {
		panic("GetPersistentDbSvcError: " + err.Error())
	}
	return persistentDbSvc
}

func GetTrailDbSvc() *internalDbInfra.TrailDatabaseService {
	trailDbSvc, err := internalDbInfra.NewTrailDatabaseService()
	if err != nil {
		panic("GetTrailDbSvcError: " + err.Error())
	}
	return trailDbSvc
}
