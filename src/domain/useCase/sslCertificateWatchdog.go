package useCase

import (
	"log"

	"github.com/speedianet/os/src/domain/repository"
)

const SslValidationsPerHour int = 4

func SslCertificateWatchdog(
	sslQueryRepo repository.SslQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
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

		err = sslCmdRepo.ReplaceWithValidSsl(sslPair)
		if err != nil {
			mainSslPairHostname := sslPair.VirtualHostsHostnames[0]
			log.Printf(
				"ReplaceWithValidSslError (%s): %s", mainSslPairHostname.String(), err.Error(),
			)
		}
	}
}
