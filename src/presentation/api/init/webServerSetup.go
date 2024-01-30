package apiInit

import wsInfra "github.com/speedianet/os/src/infra/webServer"

func WebServerSetup() {
	ws := wsInfra.WebServerSetup{}

	ws.FirstSetup()
	ws.OnStartSetup()
}
