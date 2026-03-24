package cliMiddleware

import (
	"log"
	"os"

	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
)

func CheckEnvs() {
	primaryHostname, err := infraHelper.ReadPrimaryVirtualHostHostname()
	if err != nil {
		log.Fatalf("PrimaryHostnameUnidentifiable")
	}
	os.Setenv("PRIMARY_VHOST", primaryHostname.String())

	envFilePath, err := tkValueObject.NewUnixAbsoluteFilePath(
		infraEnvs.InfiniteOsEnvFilePath, false,
	)
	if err != nil {
		log.Fatalf("InvalidEnvFilePath: %v", err)
	}

	requiredEnvVars := []string{
		"ACCOUNT_API_KEY_SECRET",
		"JWT_SECRET",
		"PRIMARY_VHOST",
	}
	autoFillableEnvVars := []string{
		"ACCOUNT_API_KEY_SECRET",
		"JWT_SECRET",
	}

	envsInspector := tkPresentation.NewEnvsInspector(
		&envFilePath, requiredEnvVars, autoFillableEnvVars,
	)
	err = envsInspector.Inspect()
	if err != nil {
		log.Fatalf("EnvsInspectorError: %v", err)
	}
}
