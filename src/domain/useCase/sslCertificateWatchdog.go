package useCase

import (
	"log"

	"github.com/speedianet/os/src/domain/repository"
)

const SslValidationsPerHour int = 4

type SslCertificateWatchdog struct {
	sslQueryRepo   repository.SslQueryRepo
	sslCmdRepo     repository.SslCmdRepo
	vhostQueryRepo repository.VirtualHostQueryRepo
	vhostCmdRepo   repository.VirtualHostCmdRepo
}

func NewSslCertificateWatchdog(
	sslQueryRepo repository.SslQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	vhostCmdRepo repository.VirtualHostCmdRepo,
) SslCertificateWatchdog {
	return SslCertificateWatchdog{
		sslQueryRepo:   sslQueryRepo,
		sslCmdRepo:     sslCmdRepo,
		vhostQueryRepo: vhostQueryRepo,
		vhostCmdRepo:   vhostCmdRepo,
	}
}

func (uc SslCertificateWatchdog) Execute() {
	sslPairs, err := uc.sslQueryRepo.GetSslPairs()
	if err != nil {
		log.Printf("FailedToGetSslPairs: %s", err.Error())
		return
	}

	for _, sslPair := range sslPairs {
		if sslPair.IsPubliclyTrusted() {
			continue
		}

		err = uc.sslCmdRepo.ReplaceWithValidSsl(sslPair)
		if err != nil {
			firstVhost := sslPair.VirtualHosts[0]
			log.Printf("FailedToReplaceWithValidSsl (%s): %s", firstVhost.String(), err.Error())
		}
	}
}
