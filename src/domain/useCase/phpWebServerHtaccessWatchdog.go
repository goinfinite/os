package useCase

import (
	"log"

	"github.com/speedianet/os/src/domain/repository"
)

const PhpWebServerHtaccessValidationsPerHour int = 30

type PhpWebServerHtaccessWatchdog struct {
	servicesQueryRepo repository.ServicesQueryRepo
	runtimeQueryRepo  repository.RuntimeQueryRepo
	runtimeCmdRepo    repository.RuntimeCmdRepo
}

func NewPhpWebServerHtaccessWatchdog(
	servicesQueyRepo repository.ServicesQueryRepo,
	runtimeQueryRepo repository.RuntimeQueryRepo,
	runtimeCmdRepo repository.RuntimeCmdRepo,
) PhpWebServerHtaccessWatchdog {
	return PhpWebServerHtaccessWatchdog{
		servicesQueryRepo: servicesQueyRepo,
		runtimeQueryRepo:  runtimeQueryRepo,
		runtimeCmdRepo:    runtimeCmdRepo,
	}
}

func (uc PhpWebServerHtaccessWatchdog) Execute() {
	_, err := uc.servicesQueryRepo.GetByName("php-webserver")
	if err != nil {
		return
	}

	if !uc.runtimeQueryRepo.IsHtaccessModifiedRecently() {
		return
	}

	err = uc.runtimeCmdRepo.RestartPhp()
	if err != nil {
		log.Print(err)
	}
}
