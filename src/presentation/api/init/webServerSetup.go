package apiInit

import (
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	wsInfra "github.com/speedianet/os/src/infra/webServer"
)

func WebServerSetup(
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) {
	ws := wsInfra.NewWebServerSetup(transientDbSvc)

	ws.FirstSetup()
	ws.OnStartSetup()
}
