package useCase

const PhpWebServerHtaccessValidationsPerHour int = 30

type PhpWebServerHtaccessWatchdog struct {
}

func NewPhpWebServerHtaccessWatchdog() PhpWebServerHtaccessWatchdog {
	return PhpWebServerHtaccessWatchdog{}
}

func (uc PhpWebServerHtaccessWatchdog) Execute() {}
