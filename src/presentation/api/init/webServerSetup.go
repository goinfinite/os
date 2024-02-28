package apiInit

import (
	databaseInfra "github.com/speedianet/os/src/infra/database"
	wsInfra "github.com/speedianet/os/src/infra/webServer"
)

func WebServerSetup(
	transientDbSvc *databaseInfra.TransientDatabaseService,
) {
	ws := wsInfra.NewWebServerSetup(transientDbSvc)

	ws.FirstSetup()
	ws.OnStartSetup()
}
