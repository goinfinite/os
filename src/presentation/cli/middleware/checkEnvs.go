package cliMiddleware

import (
	"log/slog"
	"os"

	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
)

func CheckEnvs() {
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
