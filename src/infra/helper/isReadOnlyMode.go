package infraHelper

import (
	"log/slog"
	"os"

	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

func IsReadOnlyMode() bool {
	if os.Getenv(infraEnvs.ReadOnlyModeEnvKey) == "" {
		return false
	}

	isReadOnlyModeEnabled, err := tkVoUtil.InterfaceToBool(
		os.Getenv(infraEnvs.ReadOnlyModeEnvKey),
	)
	if err != nil {
		slog.Debug("ParseReadOnlyModeEnvError", slog.String("err", err.Error()))
		return false
	}

	return isReadOnlyModeEnabled
}
