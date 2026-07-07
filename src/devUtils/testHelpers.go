package testHelpers

import (
	"encoding/base64"
	"log/slog"
	"os"

	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
)

func GenerateString(desiredSize int) string {
	desiredSizeBytesLength := float64(desiredSize) * 3
	desiredSizeStringLength := desiredSizeBytesLength / 4
	randomBytes := make([]byte, uint(desiredSizeStringLength))
	return base64.StdEncoding.EncodeToString(randomBytes)
}

func LoadEnvVars() {
	envFilePath, err := tkValueObject.NewUnixAbsoluteFilePath(
		infraEnvs.InfiniteOsEnvFilePath, false,
	)
	if err != nil {
		slog.Error("InvalidEnvFilePath", slog.String("error", err.Error()))
		os.Exit(1)
	}

	requiredEnvVars := []string{
		"ACCOUNT_API_KEY_SECRET",
		"JWT_SECRET",
	}
	autoFillableEnvVars := []string{
		"ACCOUNT_API_KEY_SECRET",
		"JWT_SECRET",
	}

	envsInspector := tkPresentation.NewEnvsInspector(
		&envFilePath, requiredEnvVars, autoFillableEnvVars,
	)
	inspectErr := envsInspector.Inspect()
	if inspectErr != nil {
		slog.Error("EnvsInspectorError", slog.String("error", inspectErr.Error()))
		os.Exit(1)
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
