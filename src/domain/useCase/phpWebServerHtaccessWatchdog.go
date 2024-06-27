package useCase

import "github.com/speedianet/os/src/domain/repository"

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

func (uc PhpWebServerHtaccessWatchdog) Execute() {}
