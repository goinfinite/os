package infraHelper

import (
	"log/slog"
	"os"

	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

const maxTerminalSessionsPerAccountDefault uint16 = 30

func ReadMaxTerminalSessionsPerAccount() uint16 {
	rawMaxSessions := os.Getenv(infraEnvs.MaxTerminalSessionsPerAccountEnvKey)
	if rawMaxSessions == "" {
		return maxTerminalSessionsPerAccountDefault
	}

	maxSessions, err := tkVoUtil.InterfaceToUint16(rawMaxSessions)
	if err != nil {
		slog.Debug(
			"ParseMaxTerminalSessionsPerAccountEnvError",
			slog.String("err", err.Error()),
		)
		return maxTerminalSessionsPerAccountDefault
	}

	return maxSessions
}
