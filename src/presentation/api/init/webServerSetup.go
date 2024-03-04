package apiInit

import (
	internalDatabaseInfra "github.com/speedianet/os/src/infra/internalDatabase"
	wsInfra "github.com/speedianet/os/src/infra/webServer"
)

func WebServerSetup(
	transientDbSvc *internalDatabaseInfra.TransientDatabaseService,
) {
	ws := wsInfra.NewWebServerSetup(transientDbSvc)

	ws.FirstSetup()
	ws.OnStartSetup()
}
