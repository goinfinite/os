package useCase

import (
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

const SslValidationsPerHour int = 4

func SslCertificateWatchdog(
	sslQueryRepo repository.SslQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) {
	sslPairEntities, err := sslQueryRepo.Read()
	if err != nil {
		slog.Error("ReadSslPairsInfraError", slog.String("error", err.Error()))
		return
	}

	for _, sslPairEntity := range sslPairEntities {
		if sslPairEntity.IsPubliclyTrusted() {
			continue
		}

		err = sslCmdRepo.ReplaceWithValidSsl(dto.NewReplaceWithValidSsl(
			sslPairEntity, operatorAccountId, operatorIpAddress,
		))
		if err != nil {
			mainSslPairHostname := sslPairEntity.VirtualHostsHostnames[0]
			slog.Error(
				"ReplaceWithValidSslInfraError",
				slog.String("error", err.Error()),
				slog.String("mainSslPairHostname", mainSslPairHostname.String()),
			)
		}
	}
}
