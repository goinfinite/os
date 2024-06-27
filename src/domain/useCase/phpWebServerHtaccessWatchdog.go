package useCase

import (
	"log"

	"github.com/speedianet/os/src/domain/repository"
)

const PhpWebServerHtaccessValidationsPerHour int = 30

type PhpWebServerHtaccessWatchdog struct {
	runtimeQueryRepo repository.RuntimeQueryRepo
	runtimeCmdRepo   repository.RuntimeCmdRepo
}

func NewPhpWebServerHtaccessWatchdog(
	runtimeQueryRepo repository.RuntimeQueryRepo,
	runtimeCmdRepo repository.RuntimeCmdRepo,
) PhpWebServerHtaccessWatchdog {
	return PhpWebServerHtaccessWatchdog{
		runtimeQueryRepo: runtimeQueryRepo,
		runtimeCmdRepo:   runtimeCmdRepo,
	}
}

func (uc PhpWebServerHtaccessWatchdog) Execute() {
	if !uc.runtimeQueryRepo.IsHtaccessModifiedRecently() {
		return
	}

	err := uc.runtimeCmdRepo.RestartPhp()
	if err != nil {
		log.Print(err)
	}
}
