package useCase

import (
	"log"

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
	sslPairs, err := sslQueryRepo.Read()
	if err != nil {
		log.Printf("ReadSslPairsError: %s", err.Error())
		return
	}

	for _, sslPair := range sslPairs {
		if sslPair.IsPubliclyTrusted() {
			continue
		}

		replaceDto := dto.NewReplaceWithValidSsl(
			sslPair, operatorAccountId, operatorIpAddress,
		)
		err = sslCmdRepo.ReplaceWithValidSsl(replaceDto)
		if err != nil {
			mainSslPairHostname := sslPair.VirtualHostsHostnames[0]
			log.Printf(
				"ReplaceWithValidSslError (%s): %s", mainSslPairHostname.String(), err.Error(),
			)
		}
	}
}
