package apiInit

import webServerInfra "github.com/speedianet/os/src/infra/webServer"

func WebServerSetup() {
	ws := webServerInfra.WebServerSetup{}

	ws.FirstSetup()
	ws.OnStartSetup()
}
