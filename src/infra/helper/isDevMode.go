package infraHelper

import (
	"log/slog"
	"os"

	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

func IsDevMode() bool {
	if os.Getenv(infraEnvs.DevModeEnvKey) == "" {
		return false
	}

	isDevModeEnabled, err := tkVoUtil.InterfaceToBool(os.Getenv(infraEnvs.DevModeEnvKey))
	if err != nil {
		slog.Debug("ParseDevModeEnvError", slog.String("err", err.Error()))
		return false
	}

	return isDevModeEnabled
}

func IsSilentExitMode() bool {
	if os.Getenv(infraEnvs.SilentExitModeEnvKey) == "" {
		return false
	}

	isSilentExitMode, err := tkVoUtil.InterfaceToBool(os.Getenv(infraEnvs.SilentExitModeEnvKey))
	if err != nil {
		slog.Debug("ParseSilentExitModeEnvError", slog.String("err", err.Error()))
		return false
	}

	return isSilentExitMode
}
